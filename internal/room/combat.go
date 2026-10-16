package room

// Combat mode, attacks and the DM log (docs/features/combat-log.md).

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"projectdnd/internal/rules"
)

const (
	maxLog    = 300   // DM log entries kept in memory
	maxDamage = 999   // attack.resolve, npc.act
	maxTalk   = 200   // npc.act talk text (runes)
	noteQuiet = 60000 // ms: at most one "edited note" log line per note in this window
	moveQuiet = 3000  // ms: moves of one token by one seat in this window share one log line
)

// Attack is a declared attack waiting for the DM's hit/miss. DM/co-DM only; never part of State.
type Attack struct {
	ID            string  `json:"id"`
	ByCharacterID string  `json:"byCharacterId"`
	ByName        rules.L `json:"byName"`
	TargetTokenID string  `json:"targetTokenId"`
	TargetName    string  `json:"targetName"`
	At            int64   `json:"at"` // unix milliseconds
}

// MyAttack is what a player sees of their own pending attack.
type MyAttack struct {
	ID         string `json:"id"`
	TargetName string `json:"targetName"`
}

// LogEntry is one line of the DM log. Actor is the sender's seat name, never a clientId.
type LogEntry struct {
	At    int64   `json:"at"` // unix milliseconds
	Actor string  `json:"actor"`
	Text  rules.L `json:"text"`

	move  string // token.move lines: seat id + "/" + token id, else ""
	from  string // token.move lines: the cell the line's first move started from
	first int64  // token.move lines: unix ms of the line's first move
}

// Banner is the big message every connection gets once after attack.resolve or npc.act.
type Banner struct {
	Kind string  `json:"kind"` // "hit" | "miss" | "talk" | "gift" | "surrender"
	Text rules.L `json:"text"`
}

var (
	errCombatOff = errors.New("combat mode is off")
	errNoTarget  = errors.New("no such target") // also for hidden tokens: players must not learn they exist
	errNoAttack  = errors.New("no such attack")
)

// addLog queues one DM log line for the current event; Apply stamps the actor and keeps it only on success.
func (r *Room) addLog(th, en string) {
	r.fresh = append(r.fresh, LogEntry{At: r.now().UnixMilli(), Text: rules.L{TH: th, EN: en}})
}

// Attacks returns the pending attacks (DM view, never nil). Callers must not mutate it.
func (r *Room) Attacks() []Attack {
	if r.attacks == nil {
		return []Attack{}
	}
	return r.attacks
}

// MyAttack returns the pending attack of clientID's character, or nil.
func (r *Room) MyAttack(clientID string) *MyAttack {
	s := r.seatOf(clientID)
	if s == nil || s.CharacterID == "" {
		return nil
	}
	if a := r.attackBy(s.CharacterID); a != nil {
		return &MyAttack{ID: a.ID, TargetName: a.TargetName}
	}
	return nil
}

// DMLog returns the DM log, oldest first (never nil). Callers must not mutate it.
func (r *Room) DMLog() []LogEntry {
	if r.dmLog == nil {
		return []LogEntry{}
	}
	return r.dmLog
}

func (r *Room) attackBy(characterID string) *Attack {
	for i := range r.attacks {
		if r.attacks[i].ByCharacterID == characterID {
			return &r.attacks[i]
		}
	}
	return nil
}

// downed adds the knocked-out / back-up line when hp crosses 0. A hidden NPC's line goes to the DM log only.
func (r *Room) downed(name rules.L, old, now int, char, hidden bool) {
	var th, en string
	switch {
	case old > 0 && now == 0 && char:
		th, en = "💀 %s หมดสติ (HP 0)", "💀 %s is knocked out (HP 0)"
	case old > 0 && now == 0:
		th, en = "💀 %s ล้มลง (HP 0)", "💀 %s falls (HP 0)"
	case old == 0 && now > 0:
		th, en = "%s ฟื้นขึ้น", "%s is back up"
	default:
		return
	}
	th, en = fmt.Sprintf(th, name.TH), fmt.Sprintf(en, name.EN)
	if hidden {
		r.addLog(th, en)
		return
	}
	r.addFeed(th, en)
}

func (r *Room) combatSet(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		On bool `json:"on"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	if p.On == r.st.Combat {
		return nil // repeat click: nothing to change or log
	}
	r.st.Combat = p.On
	if p.On {
		r.addFeed("เข้าสู่โหมดต่อสู้", "Combat started")
		return nil
	}
	if n := len(r.attacks); n > 0 {
		r.addLog(fmt.Sprintf("ยกเลิกการโจมตีที่ค้าง %d รายการ", n), fmt.Sprintf("cancelled %d pending attack(s)", n))
	}
	r.attacks = nil
	r.addFeed("จบการต่อสู้", "Combat ended")
	return nil
}

func (r *Room) attackDeclare(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		TargetTokenID string `json:"targetTokenId"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(StatusPlaying); err != nil {
		return err
	}
	if !r.st.Combat {
		return errCombatOff
	}
	c, err := r.ownCharacter(s)
	if err != nil {
		return err
	}
	if r.attackBy(c.ID) != nil {
		return errors.New("you already have a pending attack")
	}
	i := slices.IndexFunc(r.st.Tokens, func(t Token) bool { return t.CharacterID == c.ID && t.SceneID == r.st.CurrentSceneID })
	if i < 0 {
		return errors.New("your token is not in this scene")
	}
	me := r.st.Tokens[i]
	t := r.token(p.TargetTokenID)
	if t == nil || t.Hidden {
		return errNoTarget
	}
	if t.ID == me.ID {
		return errors.New("you cannot attack yourself")
	}
	if t.SceneID != me.SceneID {
		return errors.New("target is not in this scene")
	}
	r.attacks = append(r.attacks, Attack{
		ID: r.id("a"), ByCharacterID: c.ID, ByName: c.Name, TargetTokenID: t.ID, TargetName: t.Name, At: r.now().UnixMilli(),
	})
	r.addLog(fmt.Sprintf("%s โจมตี %s", c.Name.TH, t.Name), fmt.Sprintf("%s attacks %s", c.Name.EN, t.Name))
	return nil
}

func (r *Room) attackCancel(s *Seat) error {
	if err := r.needStatus(StatusPlaying); err != nil {
		return err
	}
	c, err := r.ownCharacter(s)
	if err != nil {
		return err
	}
	a := r.attackBy(c.ID)
	if a == nil {
		return errors.New("you have no pending attack")
	}
	r.addLog("ยกเลิกการโจมตี "+a.TargetName, "cancelled the attack on "+a.TargetName)
	r.attacks = slices.DeleteFunc(r.attacks, func(x Attack) bool { return x.ByCharacterID == c.ID })
	return nil
}

func (r *Room) attackResolve(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		AttackID string `json:"attackId"`
		Hit      bool   `json:"hit"`
		Damage   *int   `json:"damage"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	i := slices.IndexFunc(r.attacks, func(a Attack) bool { return a.ID == p.AttackID })
	if i < 0 {
		return errNoAttack
	}
	dmg, err := checkDamage(p.Hit, p.Damage)
	if err != nil {
		return err
	}
	a := r.attacks[i]
	r.attacks = slices.Delete(slices.Clone(r.attacks), i, i+1)

	// token.remove drops attacks on that token, so t is never nil here; stay safe anyway.
	t := r.token(a.TargetTokenID)
	var c *Character
	if t != nil && t.CharacterID != "" {
		c = r.character(t.CharacterID)
	}
	target := rules.L{TH: a.TargetName, EN: a.TargetName}
	if c != nil {
		target = c.Name
	}
	th, en := a.ByName.TH+" → "+target.TH, a.ByName.EN+" → "+target.EN
	// Hidden since the declare: DM log only, no feed and no banner (players must not learn the result).
	hidden := t != nil && t.Hidden
	say, sayPublic := r.addFeed, r.pushFeed
	if hidden {
		say, sayPublic = r.addLog, func(string, string) {}
	}
	kind := "miss"
	if p.Hit {
		kind = "hit"
	}
	if !hidden {
		r.banner = &Banner{Kind: kind, Text: rules.L{TH: th, EN: en}}
	}
	if !p.Hit {
		say(th+": ไม่โดน", en+": miss")
		return nil
	}
	switch {
	case c != nil:
		old := c.HP
		c.HP = max(c.HP-dmg, 0)
		if dmg > 0 {
			d := signed(c.HP - old)
			say(fmt.Sprintf("%s: โดน! %s (เลือด %d/%d)", th, d, c.HP, c.HPMax), fmt.Sprintf("%s: hit! %s (HP %d/%d)", en, d, c.HP, c.HPMax))
		} else {
			say(th+": โดน!", en+": hit!")
		}
		r.downed(c.Name, old, c.HP, true, hidden)
	case t != nil && t.HP != nil && t.HPMax != nil && dmg > 0:
		old := *t.HP
		*t.HP = max(old-dmg, 0)
		sayPublic(th+": โดน!", en+": hit!") // NPC: players never see numbers
		r.addLog(fmt.Sprintf("%s: โดน ดาเมจ %d (เลือด %d/%d)", th, dmg, *t.HP, *t.HPMax),
			fmt.Sprintf("%s: hit, %d damage (HP %d/%d)", en, dmg, *t.HP, *t.HPMax))
		r.downed(target, old, *t.HP, false, hidden)
	case dmg > 0: // no HP to take it from: players see a plain hit, the DM log keeps the number
		sayPublic(th+": โดน!", en+": hit!")
		r.addLog(fmt.Sprintf("%s: โดน ดาเมจ %d (ไม่มีเลือด)", th, dmg), fmt.Sprintf("%s: hit, %d damage (no HP)", en, dmg))
	default:
		say(th+": โดน!", en+": hit!")
	}
	return nil
}

// checkDamage validates the optional damage of a hit/miss (attack.resolve, npc.act attack).
func checkDamage(hit bool, d *int) (int, error) {
	if d == nil {
		return 0, nil
	}
	if *d < 0 || *d > maxDamage {
		return 0, fmt.Errorf("damage must be 0-%d", maxDamage)
	}
	if !hit && *d > 0 {
		return 0, errors.New("a miss deals no damage")
	}
	return *d, nil
}

// npcAct is the DM making a visible NPC token talk to, attack, gift or surrender to a player character in combat
// (docs/features/npc-actions.md). Feed and banner name the NPC without npc.place's "#n"; the DM log keeps it.
func (r *Room) npcAct(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		TokenID  string `json:"tokenId"`
		TargetID string `json:"targetId"`
		Kind     string `json:"kind"`
		Text     string `json:"text"`
		Hit      bool   `json:"hit"`
		Damage   *int   `json:"damage"`
		Item     struct {
			Name string `json:"name"`
			Qty  int    `json:"qty"`
		} `json:"item"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	if !r.st.Combat {
		return errCombatOff
	}
	t := r.token(p.TokenID)
	if t == nil || t.CharacterID != "" {
		return errors.New("no such NPC token")
	}
	if t.Hidden {
		return errors.New("reveal the NPC first")
	}
	if t.SceneID != r.st.CurrentSceneID {
		return errors.New("the NPC is not in this scene")
	}
	tg := r.token(p.TargetID)
	var c *Character
	if tg != nil && tg.CharacterID != "" && !tg.Hidden {
		c = r.character(tg.CharacterID)
	}
	if c == nil {
		return errNoTarget
	}
	if tg.SceneID != t.SceneID {
		return errors.New("target is not in this scene")
	}
	// Every line is "<NPC name>" + tail; the tails never go through a format string (talk text is user input).
	who, kind := c.Name, p.Kind
	var th, en string
	var downed func()
	switch p.Kind {
	case "talk":
		text, ok := rules.CleanText(p.Text, 0, maxTalk, false)
		if !ok {
			return fmt.Errorf("talk text must be one line of at most %d characters", maxTalk)
		}
		th, en = " คุยกับ "+who.TH, " talks with "+who.EN
		if text != "" {
			th, en = th+": “"+text+"”", en+": “"+text+"”"
		}
	case "attack":
		dmg, err := checkDamage(p.Hit, p.Damage)
		if err != nil {
			return err
		}
		th, en = " โจมตี "+who.TH+" — ไม่โดน", " attacks "+who.EN+" — miss"
		kind = "miss"
		if p.Hit {
			kind = "hit"
			th, en = " โจมตี "+who.TH+" — โดน!", " attacks "+who.EN+" — hit!"
			if dmg > 0 {
				th, en = th+" "+signed(-dmg), en+" "+signed(-dmg)
			}
			old := c.HP
			c.HP = max(c.HP-dmg, 0)
			downed = func() { r.downed(c.Name, old, c.HP, true, false) }
		}
	case "gift":
		it, err := giveItem(c, p.Item.Name, p.Item.Qty)
		if err != nil {
			return err
		}
		th = fmt.Sprintf(" มอบ %s ×%d ให้ %s", it.TH, p.Item.Qty, who.TH)
		en = fmt.Sprintf(" gives %s %s ×%d", who.EN, it.EN, p.Item.Qty)
	case "surrender":
		th, en = " ยอมแพ้ต่อ "+who.TH, " surrenders to "+who.EN
	default:
		return errors.New("kind must be talk, attack, gift or surrender")
	}
	name := r.npcPublicName(t)
	r.pushFeed(name+th, name+en)
	r.addLog(t.Name+th, t.Name+en)
	r.banner = &Banner{Kind: kind, Text: rules.L{TH: name + th, EN: name + en}}
	if kind == "hit" || kind == "miss" { // same banner text as attack.resolve: the big word says hit/miss
		r.banner.Text = rules.L{TH: name + " → " + who.TH, EN: name + " → " + who.EN}
	}
	if downed != nil {
		downed()
	}
	return nil
}

// npcPublicName is t's pool NPC name when t.Name is that name plus npc.place's " #n" suffix, else t.Name.
func (r *Room) npcPublicName(t *Token) string {
	if n, ok := r.npcs.Get(t.NpcID); ok && strings.HasPrefix(t.Name, n.Name+" #") {
		return n.Name
	}
	return t.Name
}

// moveLog adds "moved <name> from <from> to <to>". A move of the same token by the same seat within moveQuiet
// of that line's first move rewrites the line in place (first from, latest to, latest time) and adds no new
// line, so it sends no log frame and spamming moves cannot push older lines out of the maxLog window.
func (r *Room) moveLog(key, name, from, to string) {
	now := r.now().UnixMilli()
	for i := len(r.dmLog) - 1; i >= 0; i-- {
		if e := &r.dmLog[i]; e.move == key {
			if now-e.first < moveQuiet {
				e.At, e.Text = now, moveText(name, e.from, to)
				return
			}
			break
		}
	}
	t := moveText(name, from, to)
	r.addLog(t.TH, t.EN)
	e := &r.fresh[len(r.fresh)-1]
	e.move, e.from, e.first = key, from, now
}

func moveText(name, from, to string) rules.L {
	return rules.L{TH: fmt.Sprintf("ย้าย %s จาก %s ไป %s", name, from, to), EN: fmt.Sprintf("moved %s from %s to %s", name, from, to)}
}

// LogPresence adds "joined the room" / "left the room" for a seated client and returns the line, or nil when
// clientID holds no seat. The Hub calls it only on a clientId's first connection and its last disconnection.
func (r *Room) LogPresence(clientID string, joined bool) *LogEntry {
	s := r.seatOf(clientID)
	if s == nil {
		return nil
	}
	e := LogEntry{At: r.now().UnixMilli(), Actor: s.Name, Text: rules.L{TH: "ออกจากห้อง", EN: "left the room"}}
	if joined {
		e.Text = rules.L{TH: "เข้าห้อง", EN: "joined the room"}
	}
	r.keepLog(e)
	return &e
}

// keepLog appends lines to the DM log, keeping the newest maxLog.
func (r *Room) keepLog(es ...LogEntry) {
	r.dmLog = append(r.dmLog, es...)
	if n := len(r.dmLog); n > maxLog {
		r.dmLog = slices.Clone(r.dmLog[n-maxLog:])
	}
}

// noteLog adds "edited note <title>" at most once per note per noteQuiet.
func (r *Room) noteLog(id, title string) {
	now := r.now().UnixMilli()
	if last, ok := r.noteLogged[id]; ok && now-last < noteQuiet {
		return
	}
	if r.noteLogged == nil {
		r.noteLogged = map[string]int64{}
	}
	r.noteLogged[id] = now
	r.addLog("แก้โน้ต "+title, "edited note "+title)
}
