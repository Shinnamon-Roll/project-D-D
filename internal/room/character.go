package room

// Character, rules v2 counters, custom ability, item, feed and roster events (docs/features/rules-v2.md,
// docs/design/protocol.md "Character roster").

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"projectdnd/internal/rules"
)

const (
	maxCustomAbilities = 40
	maxFeed            = 50
	maxDelta           = 999 // character.adjust
	maxStack           = 999 // qty of one item stack after item.give merges
	maxItemStacks      = 50
	kindSpell          = "spell"
	kindFeature        = "feature"
)

var (
	errNoCharacter = errors.New("you have no character")
	errNoSuchChar  = errors.New("no such character")
	errNotOwner    = errors.New("you can only change your own character")
	errNoAbility   = errors.New("no such ability")
	errNoItem      = errors.New("no such item")
)

// addFeed appends one bilingual line to the public feed and the same line to the DM log.
func (r *Room) addFeed(th, en string) {
	r.pushFeed(th, en)
	r.addLog(th, en)
}

// pushFeed appends one line to the public feed only, keeping the newest maxFeed.
func (r *Room) pushFeed(th, en string) {
	r.st.Feed = append(r.st.Feed, FeedEntry{At: r.now().UnixMilli(), Text: rules.L{TH: th, EN: en}})
	if n := len(r.st.Feed); n > maxFeed {
		r.st.Feed = slices.Clone(r.st.Feed[n-maxFeed:])
	}
}

// signed formats a change as "+3", "−3" (U+2212 minus, as in the feed examples) or "0".
func signed(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "−" + strconv.Itoa(-n)
	}
	return "+" + strconv.Itoa(n)
}

// ownCharacter is the sender's own character; only players have one.
func (r *Room) ownCharacter(s *Seat) (*Character, error) {
	if s == nil || s.Kind != KindPlayer || s.CharacterID == "" {
		return nil, errNoCharacter
	}
	if c := r.character(s.CharacterID); c != nil {
		return c, nil
	}
	return nil, errNoCharacter
}

// dmCharacter checks the DM role and status, then finds the target character.
func (r *Room) dmCharacter(s *Seat, id string, status ...string) (*Character, error) {
	if err := r.needDM(s, status...); err != nil {
		return nil, err
	}
	if c := r.character(id); c != nil {
		return c, nil
	}
	return nil, errNoSuchChar
}

// ownerOrDM finds character id for a DM/co-DM or for the player seated on it (creating/playing).
// A player gets errNotOwner for a missing id too, so the error never says which ids exist.
func (r *Room) ownerOrDM(s *Seat, id string) (*Character, error) {
	if err := r.needStatus(inGame...); err != nil {
		return nil, err
	}
	if !isDM(s) && (s == nil || s.Kind != KindPlayer || id == "" || s.CharacterID != id) {
		return nil, errNotOwner
	}
	if c := r.character(id); c != nil {
		return c, nil
	}
	return nil, errNoSuchChar
}

// derive recomputes hpMax and every max from class, level and abilities. HP keeps the damage taken (it moves by
// the hpMax change, clamped to 0..hpMax) and every used counter clamps to its new max. Maps are replaced, never
// mutated, so a State taken before the event keeps its values.
func derive(c *Character) {
	cl, _ := rules.ClassByID(c.ClassID) // pick only creates characters of a known class
	oldMax := c.HPMax
	c.HPMax = rules.HPMax(cl, c.Level, c.Abilities.Con)
	c.HP = min(max(c.HP+c.HPMax-oldMax, 0), c.HPMax)
	c.PB = rules.ProfBonus(c.Level)
	c.Saves = rules.ProficientSaves(cl, c.ExtraSaves)
	c.SlotsMax = rules.SlotsMax(cl, c.Level)
	for i := range c.SlotsUsed {
		c.SlotsUsed[i] = min(c.SlotsUsed[i], c.SlotsMax[i])
	}
	c.Pact = rules.PactMax(cl, c.Level)
	c.PactUsed = min(c.PactUsed, c.Pact.Count)
	maxes, used := make(map[string]int, len(cl.Resources)), make(map[string]int, len(cl.Resources))
	for _, res := range cl.Resources {
		maxes[res.ID] = rules.ResourceMax(res, c.Level, c.Abilities)
		used[res.ID] = min(c.ResourcesUsed[res.ID], maxes[res.ID])
	}
	c.ResourcesMax, c.ResourcesUsed = maxes, used
	c.HitDiceMax = c.Level
	c.HitDiceUsed = min(c.HitDiceUsed, c.HitDiceMax)
}

func (r *Room) characterPick(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		HeroID  string `json:"heroId"`
		ClassID string `json:"classId"` // pre-v2 name of heroId, still accepted
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(StatusCreating); err != nil {
		return err
	}
	if s == nil || s.Kind != KindPlayer {
		return errors.New("only a seated player can pick a character")
	}
	if s.CharacterID != "" {
		return errors.New("you already have a character")
	}
	if p.HeroID == "" {
		p.HeroID = p.ClassID
	}
	e, ok := r.roster.Get(p.HeroID)
	if !ok {
		return rules.ErrNoClass
	}
	cl, ok := rules.ClassByID(e.ClassID)
	if !ok {
		return rules.ErrNoClass
	}
	if !slices.Contains(r.st.ActiveSets, e.Set) {
		return errors.New("that set is not in play")
	}
	if slices.ContainsFunc(r.st.Characters, func(c Character) bool { return c.HeroID == e.ID }) {
		return errors.New("that class is already taken")
	}
	c := Character{
		ID: r.id("c"), SeatID: s.ID, ClassID: cl.ID, HeroID: e.ID, Name: e.HeroName, Art: e.Art, Level: 1,
		Abilities: e.Abilities, AC: e.AC,
		// own copies, never nil: events change them
		Skills: append([]string{}, e.Skills...), Expertise: append([]string{}, e.Expertise...),
		ExtraSaves:      []string{},
		AbilitiesCustom: []CustomAbility{}, Items: append([]rules.RosterItem{}, e.Items...),
	}
	derive(&c) // hp 0 + hpMax gained = full HP
	r.st.Characters = append(r.st.Characters, c)
	s.CharacterID = c.ID
	r.addFeed(fmt.Sprintf("%s เลือก %s (%s)", s.Name, e.HeroName.TH, cl.Name.TH),
		fmt.Sprintf("%s picked %s (%s)", s.Name, e.HeroName.EN, cl.Name.EN))
	for _, x := range r.st.Seats {
		if x.Kind == KindPlayer && x.Occupied && x.CharacterID == "" {
			return nil
		}
	}
	r.st.Status = StatusPlaying
	return nil
}

func (r *Room) characterAdjust(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		HP          *int   `json:"hp"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	if p.HP == nil {
		return errors.New("hp delta is required")
	}
	if *p.HP < -maxDelta || *p.HP > maxDelta {
		return fmt.Errorf("delta must be -%d to %d", maxDelta, maxDelta)
	}
	old := c.HP
	c.HP = min(max(c.HP+*p.HP, 0), c.HPMax)
	d := c.HP - old
	if d == 0 {
		return nil // clamped to no change: nothing to log
	}
	r.addFeed(fmt.Sprintf("DM: %s เลือด %s (%d/%d)", c.Name.TH, signed(d), c.HP, c.HPMax),
		fmt.Sprintf("DM: %s HP %s (%d/%d)", c.Name.EN, signed(d), c.HP, c.HPMax))
	r.downed(c.Name, old, c.HP, true, false)
	return nil
}

func (r *Room) characterLevel(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		Delta       int    `json:"delta"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	if p.Delta != 1 && p.Delta != -1 {
		return errors.New("delta must be 1 or -1")
	}
	lv := c.Level + p.Delta
	if lv < rules.MinLevel || lv > rules.MaxLevel {
		return fmt.Errorf("level must be %d-%d", rules.MinLevel, rules.MaxLevel)
	}
	old := c.HP
	c.Level = lv
	derive(c)
	if p.Delta > 0 {
		r.addFeed(fmt.Sprintf("%s ขึ้นเลเวล %d", c.Name.TH, lv), fmt.Sprintf("%s reaches level %d", c.Name.EN, lv))
	} else {
		r.addFeed(fmt.Sprintf("%s ลดเป็นเลเวล %d", c.Name.TH, lv), fmt.Sprintf("%s drops to level %d", c.Name.EN, lv))
	}
	r.downed(c.Name, old, c.HP, true, false)
	return nil
}

func (r *Room) characterSetAbilities(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string          `json:"characterId"`
		Abilities   rules.Abilities `json:"abilities"`
		AC          int             `json:"ac"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	if err := rules.CheckAbilities(p.Abilities, p.AC); err != nil {
		return err
	}
	old := c.HP
	c.Abilities, c.AC = p.Abilities, p.AC
	derive(c)
	a := c.Abilities
	scores := fmt.Sprintf("STR %d, DEX %d, CON %d, INT %d, WIS %d, CHA %d, AC %d", a.Str, a.Dex, a.Con, a.Int, a.Wis, a.Cha, c.AC)
	r.addLog(fmt.Sprintf("ตั้งค่า %s: %s (เลือด %d/%d)", c.Name.TH, scores, c.HP, c.HPMax),
		fmt.Sprintf("set %s: %s (HP %d/%d)", c.Name.EN, scores, c.HP, c.HPMax))
	r.downed(c.Name, old, c.HP, true, false)
	return nil
}

// characterSetSkills replaces the skill/save proficiencies class features add after the pick (Expertise,
// Slippery Mind, ...). Same skill rules as roster.update; extra saves are ability ids, duplicates dropped.
func (r *Room) characterSetSkills(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string   `json:"characterId"`
		Skills      []string `json:"skills"`
		Expertise   []string `json:"expertise"`
		Saves       []string `json:"saves"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	sk, ex, err := rules.CleanSkills(p.Skills, p.Expertise)
	if err != nil {
		return err
	}
	sv, err := rules.CleanSaves(p.Saves)
	if err != nil {
		return err
	}
	c.Skills, c.Expertise, c.ExtraSaves = sk, ex, sv
	derive(c)
	list := func(xs []string) string {
		if len(xs) == 0 {
			return "—"
		}
		return strings.Join(xs, ", ")
	}
	r.addLog(fmt.Sprintf("ตั้ง skill ของ %s: %s; expertise: %s; save เพิ่ม: %s", c.Name.TH, list(sk), list(ex), list(sv)),
		fmt.Sprintf("set skills of %s: %s; expertise: %s; extra saves: %s", c.Name.EN, list(sk), list(ex), list(sv)))
	return nil
}

// slotSet handles slot.set: level 1-9 sets slotsUsed[level-1], level "pact" sets pactUsed.
func (r *Room) slotSet(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string          `json:"characterId"`
		Level       json.RawMessage `json:"level"`
		Used        int             `json:"used"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.ownerOrDM(s, p.CharacterID)
	if err != nil {
		return err
	}
	var lv int
	var name string
	pact := false
	if json.Unmarshal(p.Level, &name) == nil && name == "pact" {
		pact = true
	} else if json.Unmarshal(p.Level, &lv) != nil || lv < 1 || lv > 9 {
		return errors.New(`level must be 1-9 or "pact"`)
	}
	used, mx := &c.PactUsed, c.Pact.Count
	if !pact {
		used, mx = &c.SlotsUsed[lv-1], c.SlotsMax[lv-1]
	}
	if mx == 0 {
		return errors.New("no such spell slot")
	}
	if p.Used < 0 || p.Used > mx {
		return fmt.Errorf("used must be 0-%d", mx)
	}
	*used = p.Used
	if pact {
		r.addLog(fmt.Sprintf("%s: ใช้ pact slot %d/%d", c.Name.TH, p.Used, mx), fmt.Sprintf("%s: pact slots used %d/%d", c.Name.EN, p.Used, mx))
	} else {
		r.addLog(fmt.Sprintf("%s: ใช้ช่องคาถาระดับ %d %d/%d", c.Name.TH, lv, p.Used, mx),
			fmt.Sprintf("%s: level %d spell slots used %d/%d", c.Name.EN, lv, p.Used, mx))
	}
	return nil
}

func (r *Room) resourceSet(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		ResourceID  string `json:"resourceId"`
		Used        int    `json:"used"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.ownerOrDM(s, p.CharacterID)
	if err != nil {
		return err
	}
	mx := c.ResourcesMax[p.ResourceID]
	if mx == 0 { // unknown id, or not gained at this level
		return errors.New("no such resource")
	}
	if p.Used < 0 || p.Used > mx {
		return fmt.Errorf("used must be 0-%d", mx)
	}
	c.ResourcesUsed[p.ResourceID] = p.Used
	cl, _ := rules.ClassByID(c.ClassID)
	i := slices.IndexFunc(cl.Resources, func(x rules.Resource) bool { return x.ID == p.ResourceID })
	n := cl.Resources[i].Name // ResourcesMax only holds the class's own resources
	r.addLog(fmt.Sprintf("%s: ใช้ %s %d/%d", c.Name.TH, n.TH, p.Used, mx), fmt.Sprintf("%s: %s used %d/%d", c.Name.EN, n.EN, p.Used, mx))
	return nil
}

func (r *Room) hitDiceSet(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		Used        int    `json:"used"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.ownerOrDM(s, p.CharacterID)
	if err != nil {
		return err
	}
	if p.Used < 0 || p.Used > c.HitDiceMax {
		return fmt.Errorf("used must be 0-%d", c.HitDiceMax)
	}
	c.HitDiceUsed = p.Used
	r.addLog(fmt.Sprintf("%s: ใช้ hit dice %d/%d", c.Name.TH, p.Used, c.HitDiceMax),
		fmt.Sprintf("%s: hit dice used %d/%d", c.Name.EN, p.Used, c.HitDiceMax))
	return nil
}

// partyRest handles party.rest for every character. Short: resources with refill "short" (or at/after their
// shortFromLevel) reset, other resources give back shortRegain uses, pact slots reset. Long: everything full.
func (r *Room) partyRest(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		Kind string `json:"kind"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	if p.Kind != "short" && p.Kind != "long" {
		return errors.New("kind must be short or long")
	}
	if p.Kind == "short" {
		r.addFeed("DM: พักสั้น — คืนความสามารถแบบพักสั้นและ pact slot", "DM: short rest — short-rest uses and pact slots back")
	} else {
		r.addFeed("DM: พักยาว — เลือด ช่องคาถา ความสามารถ และ hit dice เต็ม", "DM: long rest — HP, spell slots, uses and hit dice refilled")
	}
	for i := range r.st.Characters {
		c := &r.st.Characters[i]
		cl, _ := rules.ClassByID(c.ClassID)
		used := make(map[string]int, len(c.ResourcesUsed)) // replaced, never mutated (see derive)
		c.PactUsed = 0
		if p.Kind == "long" {
			for id := range c.ResourcesUsed {
				used[id] = 0
			}
			old := c.HP
			c.HP, c.SlotsUsed, c.HitDiceUsed, c.ResourcesUsed = c.HPMax, [9]int{}, 0, used
			r.downed(c.Name, old, c.HP, true, false)
			continue
		}
		for _, res := range cl.Resources {
			u := c.ResourcesUsed[res.ID]
			if res.Refill == "short" || (res.ShortFromLevel > 0 && c.Level >= res.ShortFromLevel) {
				u = 0
			}
			used[res.ID] = max(0, u-res.ShortRegain)
		}
		c.ResourcesUsed = used
	}
	return nil
}

// customAbility handles ability.add and ability.update (the sender's own character) and ability.remove
// (owner or DM). Name and desc are the player's own text in one language, so plain strings, not rules.L.
func (r *Room) customAbility(s *Seat, typ string, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		AbilityID   string `json:"abilityId"`
		Name        string `json:"name"`
		Kind        string `json:"kind"`
		Level       int    `json:"level"`
		Desc        string `json:"desc"`
	}](payload)
	if err != nil {
		return err
	}
	var c *Character
	if typ == "ability.remove" {
		c, err = r.ownerOrDM(s, p.CharacterID)
	} else if err = r.needStatus(inGame...); err == nil {
		c, err = r.ownCharacter(s)
	}
	if err != nil {
		return err
	}
	i := slices.IndexFunc(c.AbilitiesCustom, func(a CustomAbility) bool { return a.ID == p.AbilityID })
	if typ != "ability.add" && i < 0 {
		return errNoAbility
	}
	if typ == "ability.remove" {
		n := c.AbilitiesCustom[i].Name
		c.AbilitiesCustom = slices.Delete(c.AbilitiesCustom, i, i+1)
		r.addLog(fmt.Sprintf("ลบ %s ของ %s", n, c.Name.TH), fmt.Sprintf("removed %s from %s", n, c.Name.EN))
		return nil
	}
	name, ok := rules.CleanText(p.Name, 1, rules.MaxName, false)
	if !ok {
		return fmt.Errorf("ability name must be 1-%d characters", rules.MaxName)
	}
	desc, ok := rules.CleanText(p.Desc, 0, rules.MaxDesc, true)
	if !ok {
		return fmt.Errorf("ability desc must be at most %d characters", rules.MaxDesc)
	}
	switch p.Kind {
	case kindSpell:
		if p.Level < 0 || p.Level > 9 {
			return errors.New("spell level must be 0-9")
		}
	case kindFeature:
		p.Level = 0 // a feature has no spell level
	default:
		return errors.New("kind must be spell or feature")
	}
	a := CustomAbility{Name: name, Kind: p.Kind, Level: p.Level, Desc: desc}
	if typ == "ability.update" {
		a.ID = c.AbilitiesCustom[i].ID
		c.AbilitiesCustom[i] = a
		r.addLog(fmt.Sprintf("แก้ %s ของ %s", name, c.Name.TH), fmt.Sprintf("edited %s of %s", name, c.Name.EN))
		return nil
	}
	if len(c.AbilitiesCustom) >= maxCustomAbilities {
		return fmt.Errorf("at most %d custom abilities", maxCustomAbilities)
	}
	a.ID = r.id("ab")
	c.AbilitiesCustom = append(c.AbilitiesCustom, a)
	if a.Kind == kindSpell {
		r.addLog(fmt.Sprintf("เพิ่มคาถา %s ให้ %s", name, c.Name.TH), fmt.Sprintf("added spell %s to %s", name, c.Name.EN))
	} else {
		r.addLog(fmt.Sprintf("เพิ่มความสามารถ %s ให้ %s", name, c.Name.TH), fmt.Sprintf("added feature %s to %s", name, c.Name.EN))
	}
	return nil
}

func (r *Room) itemGive(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		Name        string `json:"name"`
		Qty         int    `json:"qty"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	it, err := giveItem(c, p.Name, p.Qty)
	if err != nil {
		return err
	}
	r.addFeed(fmt.Sprintf("DM: ให้ %s ×%d แก่ %s", it.TH, p.Qty, c.Name.TH),
		fmt.Sprintf("DM: gave %s %s ×%d", c.Name.EN, it.EN, p.Qty))
	return nil
}

// giveItem validates name and qty, then adds them to c's bag (item.give, npc.act gift) and returns the stack's name.
// On error c is unchanged.
func giveItem(c *Character, rawName string, qty int) (rules.L, error) {
	name, ok := rules.CleanText(rawName, 1, rules.MaxName, false)
	if !ok {
		return rules.L{}, fmt.Errorf("item name must be 1-%d characters", rules.MaxName)
	}
	if qty < 1 || qty > rules.MaxQty {
		return rules.L{}, fmt.Errorf("item qty must be 1-%d", rules.MaxQty)
	}
	// Same name in either language = same stack.
	i := slices.IndexFunc(c.Items, func(it rules.RosterItem) bool { return it.Name.TH == name || it.Name.EN == name })
	switch {
	case i >= 0 && c.Items[i].Qty+qty > maxStack:
		return rules.L{}, fmt.Errorf("an item stack holds at most %d", maxStack)
	case i >= 0:
		c.Items[i].Qty += qty
	case len(c.Items) >= maxItemStacks:
		return rules.L{}, fmt.Errorf("at most %d different items", maxItemStacks)
	default:
		i = len(c.Items)
		c.Items = append(c.Items, rules.RosterItem{Name: rules.L{TH: name, EN: name}, Qty: qty})
	}
	return c.Items[i].Name, nil
}

// itemRemove takes qty from the stack at index; qty at or above the stack size removes the whole stack.
func (r *Room) itemRemove(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		Index       int    `json:"index"`
		Qty         int    `json:"qty"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	if p.Index < 0 || p.Index >= len(c.Items) {
		return errNoItem
	}
	if p.Qty < 1 {
		return errors.New("qty must be at least 1")
	}
	it := c.Items[p.Index]
	n := min(p.Qty, it.Qty)
	if n == it.Qty {
		c.Items = slices.Delete(c.Items, p.Index, p.Index+1)
	} else {
		c.Items[p.Index].Qty -= n
	}
	r.addFeed(fmt.Sprintf("DM: เอา %s ×%d ออกจาก %s", it.Name.TH, n, c.Name.TH),
		fmt.Sprintf("DM: removed %s ×%d from %s", it.Name.EN, n, c.Name.EN))
	return nil
}

// rosterEdit handles roster.update and roster.reset (DM, any status). Characters already picked keep their copy.
func (r *Room) rosterEdit(typ string, payload json.RawMessage) error {
	p, err := decode[struct {
		ClassID string            `json:"classId"`
		Set     string            `json:"set"`
		Entry   rules.RosterEntry `json:"entry"`
	}](payload)
	if err != nil {
		return err
	}
	if typ == "roster.reset" && p.Set != "" {
		if p.ClassID != "" {
			return errors.New("send classId or set, not both")
		}
		if err := r.roster.ResetSet(p.Set); err != nil {
			return err
		}
		r.addLog("รีเซ็ตเซ็ต "+p.Set+" ใน roster", "reset roster set "+p.Set)
		return nil
	}
	if typ == "roster.reset" {
		if err := r.roster.Reset(p.ClassID); err != nil {
			return err
		}
		if p.ClassID == "" {
			r.addLog("รีเซ็ต roster ทุกคลาส", "reset every roster class")
		} else {
			r.addLog("รีเซ็ตคลาส "+p.ClassID+" ใน roster", "reset roster class "+p.ClassID)
		}
		return nil
	}
	if !rules.IsClass(p.ClassID) {
		return rules.ErrNoClass
	}
	p.Entry.ID = p.ClassID // the key is classId; entry.id is not trusted
	if p.Entry.Art != "" && !r.uploadedFile(p.Entry.Art) {
		return errors.New("art must be empty or an uploaded image")
	}
	if err := r.roster.Update(p.Entry); err != nil {
		return err
	}
	r.addLog("แก้คลาส "+p.ClassID+" ใน roster", "edited roster class "+p.ClassID)
	return nil
}
