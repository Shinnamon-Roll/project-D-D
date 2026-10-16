package room

// Character, custom skill, item, feed and roster events (docs/features/character-roster.md, "Protocol contract").

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"projectdnd/internal/rules"
)

const (
	maxCustomSkills = 3
	maxFeed         = 50
	maxDelta        = 999 // character.adjust
	maxStack        = 999 // qty of one item stack after item.give merges
	maxItemStacks   = 50
	kindInnate      = "innate"
	kindCustom      = "custom"
)

var (
	errNoCharacter = errors.New("you have no character")
	errNoSuchChar  = errors.New("no such character")
	errNoSkill     = errors.New("no such skill")
	errNoItem      = errors.New("no such item")
)

// addFeed appends one bilingual log line, keeping the newest maxFeed.
func (r *Room) addFeed(th, en string) {
	r.st.Feed = append(r.st.Feed, FeedEntry{At: time.Now().UnixMilli(), Text: rules.L{TH: th, EN: en}})
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

func (r *Room) characterPick(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		ClassID string `json:"classId"`
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
	e, ok := r.roster.Get(p.ClassID)
	if !ok {
		return rules.ErrNoClass
	}
	if slices.ContainsFunc(r.st.Characters, func(c Character) bool { return c.ClassID == e.ID }) {
		return errors.New("that class is already taken")
	}
	skills := make([]CharSkill, 0, len(e.Skills)+maxCustomSkills)
	for _, k := range e.Skills {
		skills = append(skills, CharSkill{ID: k.ID, Kind: kindInnate, Name: k.Name, Desc: k.Desc, Cost: k.Cost})
	}
	st := e.Stats
	c := Character{
		ID: r.id("c"), SeatID: s.ID, ClassID: e.ID, Name: e.HeroName, Art: e.Art,
		HP: st.HPMax, HPMax: st.HPMax, MP: st.MPMax, MPMax: st.MPMax, Str: st.Str, End: st.End, Agi: st.Agi,
		Skills: skills, Items: append([]rules.RosterItem{}, e.Items...), // own copy: item events change it
	}
	r.st.Characters = append(r.st.Characters, c)
	s.CharacterID = c.ID
	r.addFeed(fmt.Sprintf("%s เลือก %s (%s)", s.Name, e.HeroName.TH, e.ClassName.TH),
		fmt.Sprintf("%s picked %s (%s)", s.Name, e.HeroName.EN, e.ClassName.EN))
	for _, x := range r.st.Seats {
		if x.Kind == KindPlayer && x.Occupied && x.CharacterID == "" {
			return nil
		}
	}
	r.st.Status = StatusPlaying
	return nil
}

// customSkill handles skill.addCustom (owner), skill.editCustom and skill.removeCustom (owner or DM).
func (r *Room) customSkill(s *Seat, typ string, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		SkillID     string `json:"skillId"`
		Name        string `json:"name"`
		Desc        string `json:"desc"`
		Cost        int    `json:"cost"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(inGame...); err != nil {
		return err
	}
	var c *Character
	if typ == "skill.addCustom" {
		if c, err = r.ownCharacter(s); err != nil {
			return err
		}
	} else {
		if !isDM(s) && (s == nil || s.Kind != KindPlayer || p.CharacterID == "" || s.CharacterID != p.CharacterID) {
			return errors.New("you can only change your own character's skills")
		}
		if c = r.character(p.CharacterID); c == nil {
			return errNoSuchChar
		}
	}
	i := slices.IndexFunc(c.Skills, func(k CharSkill) bool { return k.ID == p.SkillID })
	if typ != "skill.addCustom" {
		if i < 0 {
			return errNoSkill
		}
		if c.Skills[i].Kind != kindCustom {
			return errors.New("only custom skills can be changed")
		}
		if typ == "skill.removeCustom" {
			c.Skills = slices.Delete(c.Skills, i, i+1)
			return nil
		}
	}
	name, ok := rules.CleanText(p.Name, 1, rules.MaxName, false)
	if !ok {
		return fmt.Errorf("skill name must be 1-%d characters", rules.MaxName)
	}
	desc, ok := rules.CleanText(p.Desc, 0, rules.MaxDesc, true)
	if !ok {
		return fmt.Errorf("skill desc must be at most %d characters", rules.MaxDesc)
	}
	if p.Cost < 0 || p.Cost > rules.MaxCost {
		return fmt.Errorf("skill cost must be 0-%d", rules.MaxCost)
	}
	k := CharSkill{Kind: kindCustom, Name: rules.L{TH: name, EN: name}, Desc: rules.L{TH: desc, EN: desc}, Cost: p.Cost}
	if typ == "skill.editCustom" {
		k.ID = c.Skills[i].ID
		c.Skills[i] = k
		return nil
	}
	custom := 0
	for _, x := range c.Skills {
		if x.Kind == kindCustom {
			custom++
		}
	}
	if custom >= maxCustomSkills {
		return fmt.Errorf("at most %d custom skills", maxCustomSkills)
	}
	k.ID = r.id("cs_") // "_" never appears in roster skill ids, so custom ids cannot collide with innate ones
	c.Skills = append(c.Skills, k)
	return nil
}

func (r *Room) skillUse(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		SkillID string `json:"skillId"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(StatusPlaying); err != nil {
		return err
	}
	c, err := r.ownCharacter(s)
	if err != nil {
		return err
	}
	i := slices.IndexFunc(c.Skills, func(k CharSkill) bool { return k.ID == p.SkillID })
	if i < 0 {
		return errNoSkill
	}
	k := c.Skills[i]
	if c.MP < k.Cost {
		return errors.New("not enough mana")
	}
	c.MP -= k.Cost
	r.addFeed(fmt.Sprintf("%s ใช้ %s (%s มานา)", c.Name.TH, k.Name.TH, signed(-k.Cost)),
		fmt.Sprintf("%s used %s (%s MP)", c.Name.EN, k.Name.EN, signed(-k.Cost)))
	return nil
}

func (r *Room) characterAdjust(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		HP          *int   `json:"hp"`
		MP          *int   `json:"mp"`
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	if p.HP == nil && p.MP == nil {
		return errors.New("hp or mp delta is required")
	}
	for _, d := range []*int{p.HP, p.MP} {
		if d != nil && (*d < -maxDelta || *d > maxDelta) {
			return fmt.Errorf("delta must be -%d to %d", maxDelta, maxDelta)
		}
	}
	var th, en []string
	if p.HP != nil {
		old := c.HP
		c.HP = min(max(c.HP+*p.HP, 0), c.HPMax)
		if d := c.HP - old; d != 0 {
			th = append(th, fmt.Sprintf("เลือด %s (%d/%d)", signed(d), c.HP, c.HPMax))
			en = append(en, fmt.Sprintf("HP %s (%d/%d)", signed(d), c.HP, c.HPMax))
		}
	}
	if p.MP != nil {
		old := c.MP
		c.MP = min(max(c.MP+*p.MP, 0), c.MPMax)
		if d := c.MP - old; d != 0 {
			th = append(th, fmt.Sprintf("มานา %s (%d/%d)", signed(d), c.MP, c.MPMax))
			en = append(en, fmt.Sprintf("MP %s (%d/%d)", signed(d), c.MP, c.MPMax))
		}
	}
	if len(en) == 0 {
		return nil // clamped to no change: nothing to log
	}
	r.addFeed("DM: "+c.Name.TH+" "+strings.Join(th, ", "), "DM: "+c.Name.EN+" "+strings.Join(en, ", "))
	return nil
}

func (r *Room) characterSetStats(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CharacterID string `json:"characterId"`
		rules.Stats
	}](payload)
	if err != nil {
		return err
	}
	c, err := r.dmCharacter(s, p.CharacterID, inGame...)
	if err != nil {
		return err
	}
	if err := p.Stats.Check(); err != nil {
		return err
	}
	c.HPMax, c.MPMax, c.Str, c.End, c.Agi = p.HPMax, p.MPMax, p.Str, p.End, p.Agi
	c.HP, c.MP = min(c.HP, c.HPMax), min(c.MP, c.MPMax)
	stats := fmt.Sprintf("STR %d, END %d, AGI %d", c.Str, c.End, c.Agi)
	r.addFeed(fmt.Sprintf("DM: ตั้งค่า %s: เลือดสูงสุด %d, มานาสูงสุด %d, %s", c.Name.TH, c.HPMax, c.MPMax, stats),
		fmt.Sprintf("DM: set %s: HP max %d, MP max %d, %s", c.Name.EN, c.HPMax, c.MPMax, stats))
	return nil
}

func (r *Room) partyRest(s *Seat) error {
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	for i := range r.st.Characters {
		r.st.Characters[i].MP = r.st.Characters[i].MPMax
	}
	r.addFeed("DM: พักผ่อน — มานาเต็ม", "DM: rest — mana refilled")
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
	name, ok := rules.CleanText(p.Name, 1, rules.MaxName, false)
	if !ok {
		return fmt.Errorf("item name must be 1-%d characters", rules.MaxName)
	}
	if p.Qty < 1 || p.Qty > rules.MaxQty {
		return fmt.Errorf("item qty must be 1-%d", rules.MaxQty)
	}
	// Same name in either language = same stack.
	i := slices.IndexFunc(c.Items, func(it rules.RosterItem) bool { return it.Name.TH == name || it.Name.EN == name })
	switch {
	case i >= 0 && c.Items[i].Qty+p.Qty > maxStack:
		return fmt.Errorf("an item stack holds at most %d", maxStack)
	case i >= 0:
		c.Items[i].Qty += p.Qty
	case len(c.Items) >= maxItemStacks:
		return fmt.Errorf("at most %d different items", maxItemStacks)
	default:
		i = len(c.Items)
		c.Items = append(c.Items, rules.RosterItem{Name: rules.L{TH: name, EN: name}, Qty: p.Qty})
	}
	it := c.Items[i].Name
	r.addFeed(fmt.Sprintf("DM: ให้ %s ×%d แก่ %s", it.TH, p.Qty, c.Name.TH),
		fmt.Sprintf("DM: gave %s %s ×%d", c.Name.EN, it.EN, p.Qty))
	return nil
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
		Entry   rules.RosterEntry `json:"entry"`
	}](payload)
	if err != nil {
		return err
	}
	if typ == "roster.reset" {
		return r.roster.Reset(p.ClassID)
	}
	if !rules.IsClass(p.ClassID) {
		return rules.ErrNoClass
	}
	p.Entry.ID = p.ClassID // the key is classId; entry.id is not trusted
	if p.Entry.Art != "" && !r.uploadedFile(p.Entry.Art) {
		return errors.New("art must be empty or an uploaded image")
	}
	return r.roster.Update(p.Entry)
}
