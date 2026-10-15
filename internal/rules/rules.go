// Package rules holds the white-label ruleset, character validation and the built-in presets.
// JSON names follow docs/design/protocol.md.
package rules

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type StatDef struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Base int    `json:"base"`
}

type SkillDef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	StatKey string `json:"statKey"` // may be ""
	// Description is optional DM text; preset skills leave it "" (the frontend shows SRD text).
	Description string `json:"description"`
}

type ItemDef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ItemStack struct {
	ItemID string `json:"itemId"`
	Qty    int    `json:"qty"`
}

type Ruleset struct {
	Name               string      `json:"name"`
	Attribution        string      `json:"attribution"`
	Stats              []StatDef   `json:"stats"`
	StatPoints         int         `json:"statPoints"`
	Skills             []SkillDef  `json:"skills"`
	StartingSkillPicks int         `json:"startingSkillPicks"`
	Items              []ItemDef   `json:"items"`
	StartingItems      []ItemStack `json:"startingItems"`
}

type Campaign struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "available" | "coming_soon"
}

const (
	PresetDnD5eID = "dnd5e-srd"
	maxText       = 64
	maxAttrib     = 1000
	maxSkillDesc  = 500
)

// SRDAttribution is the verbatim notice required by SRD 5.2.1 (docs/research/dnd-basics.md section 6).
const SRDAttribution = `This work includes material from the System Reference Document 5.2.1 ("SRD 5.2.1") by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd. The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International License, available at https://creativecommons.org/licenses/by/4.0/legalcode.`

// Campaigns is the hard-coded campaign menu. No adventure content is bundled.
func Campaigns() []Campaign {
	return []Campaign{
		{ID: "starter-set", Name: "Starter Set: Dragons of Stormwreck Isle", Status: "available"},
		{ID: "salem", Name: "Salem", Status: "coming_soon"},
	}
}

// PresetDnD5e returns a fresh copy of the 5E-compatible preset built from SRD 5.2.1 data only.
func PresetDnD5e() Ruleset {
	r := Ruleset{
		Name:               "D&D 5e (SRD 5.2.1)",
		Attribution:        SRDAttribution,
		StatPoints:         27,
		StartingSkillPicks: 2,
	}
	for _, k := range []string{"str", "dex", "con", "int", "wis", "cha"} {
		r.Stats = append(r.Stats, StatDef{Key: k, Name: upper(k), Min: 8, Max: 15, Base: 8})
	}
	// HP and AC: max == base so a player cannot spend points there; the DM sets them with character.grant.
	r.Stats = append(r.Stats,
		StatDef{Key: "hp", Name: "HP", Min: 10, Max: 10, Base: 10},
		StatDef{Key: "ac", Name: "AC", Min: 10, Max: 10, Base: 10},
	)
	skills := []struct{ id, name, stat string }{
		{"athletics", "Athletics", "str"},
		{"acrobatics", "Acrobatics", "dex"},
		{"sleight-of-hand", "Sleight of Hand", "dex"},
		{"stealth", "Stealth", "dex"},
		{"arcana", "Arcana", "int"},
		{"history", "History", "int"},
		{"investigation", "Investigation", "int"},
		{"nature", "Nature", "int"},
		{"religion", "Religion", "int"},
		{"animal-handling", "Animal Handling", "wis"},
		{"insight", "Insight", "wis"},
		{"medicine", "Medicine", "wis"},
		{"perception", "Perception", "wis"},
		{"survival", "Survival", "wis"},
		{"deception", "Deception", "cha"},
		{"intimidation", "Intimidation", "cha"},
		{"performance", "Performance", "cha"},
		{"persuasion", "Persuasion", "cha"},
	}
	for _, s := range skills {
		r.Skills = append(r.Skills, SkillDef{ID: s.id, Name: s.name, StatKey: s.stat})
	}
	r.Items = []ItemDef{
		{ID: "backpack", Name: "Backpack"},
		{ID: "bedroll", Name: "Bedroll"},
		{ID: "rations", Name: "Rations (1 day)"},
		{ID: "rope", Name: "Rope (50 feet)"},
		{ID: "torch", Name: "Torch"},
		{ID: "waterskin", Name: "Waterskin"},
		{ID: "healing-potion", Name: "Potion of Healing"},
	}
	r.StartingItems = []ItemStack{
		{ItemID: "backpack", Qty: 1},
		{ItemID: "bedroll", Qty: 1},
		{ItemID: "rations", Qty: 5},
		{ItemID: "rope", Qty: 1},
		{ItemID: "torch", Qty: 5},
		{ItemID: "waterskin", Qty: 1},
	}
	return r
}

func upper(s string) string {
	b := []byte(s)
	for i := range b {
		b[i] -= 'a' - 'A'
	}
	return string(b)
}

func textOK(s string) bool {
	n := utf8.RuneCountInString(s)
	return n >= 1 && n <= maxText
}

// descOK allows 0-maxSkillDesc runes, newlines included, but no other control characters.
func descOK(s string) bool {
	if utf8.RuneCountInString(s) > maxSkillDesc {
		return false
	}
	for _, c := range s {
		if c != '\n' && unicode.IsControl(c) {
			return false
		}
	}
	return true
}

// Normalize replaces nil slices with empty ones so JSON never carries null arrays,
// and trims skill descriptions.
func (r *Ruleset) Normalize() {
	for i := range r.Skills {
		r.Skills[i].Description = strings.TrimSpace(r.Skills[i].Description)
	}
	if r.Stats == nil {
		r.Stats = []StatDef{}
	}
	if r.Skills == nil {
		r.Skills = []SkillDef{}
	}
	if r.Items == nil {
		r.Items = []ItemDef{}
	}
	if r.StartingItems == nil {
		r.StartingItems = []ItemStack{}
	}
}

// ValidateRuleset checks a DM-submitted ruleset: unique keys/ids, min <= base <= max, counts >= 0,
// startingItems reference known items, and that a valid character can exist at all.
func ValidateRuleset(r Ruleset) error {
	if !textOK(r.Name) {
		return fmt.Errorf("ruleset name must be 1-%d characters", maxText)
	}
	if utf8.RuneCountInString(r.Attribution) > maxAttrib {
		return fmt.Errorf("attribution longer than %d characters", maxAttrib)
	}
	if r.StatPoints < 0 || r.StartingSkillPicks < 0 {
		return errors.New("statPoints and startingSkillPicks must be >= 0")
	}
	stats := map[string]bool{}
	room := 0 // total points that can be spent above base
	for _, s := range r.Stats {
		if !textOK(s.Key) || !textOK(s.Name) {
			return fmt.Errorf("stat key and name must be 1-%d characters", maxText)
		}
		if stats[s.Key] {
			return fmt.Errorf("duplicate stat key %q", s.Key)
		}
		stats[s.Key] = true
		if !(s.Min <= s.Base && s.Base <= s.Max) {
			return fmt.Errorf("stat %q needs min <= base <= max", s.Key)
		}
		room += s.Max - s.Base
	}
	if r.StatPoints > room {
		return fmt.Errorf("statPoints %d exceed the %d points the stats can hold", r.StatPoints, room)
	}
	skills := map[string]bool{}
	for _, s := range r.Skills {
		if !textOK(s.ID) || !textOK(s.Name) {
			return fmt.Errorf("skill id and name must be 1-%d characters", maxText)
		}
		if !descOK(s.Description) {
			return fmt.Errorf("skill %q description must be at most %d characters with no control characters", s.ID, maxSkillDesc)
		}
		if skills[s.ID] {
			return fmt.Errorf("duplicate skill id %q", s.ID)
		}
		skills[s.ID] = true
		if s.StatKey != "" && !stats[s.StatKey] {
			return fmt.Errorf("skill %q uses unknown stat %q", s.ID, s.StatKey)
		}
	}
	if r.StartingSkillPicks > len(r.Skills) {
		return fmt.Errorf("startingSkillPicks %d exceeds %d skills", r.StartingSkillPicks, len(r.Skills))
	}
	items := map[string]bool{}
	for _, it := range r.Items {
		if !textOK(it.ID) || !textOK(it.Name) {
			return fmt.Errorf("item id and name must be 1-%d characters", maxText)
		}
		if items[it.ID] {
			return fmt.Errorf("duplicate item id %q", it.ID)
		}
		items[it.ID] = true
	}
	seen := map[string]bool{}
	for _, st := range r.StartingItems {
		if !items[st.ItemID] {
			return fmt.Errorf("starting item %q is not in items", st.ItemID)
		}
		if seen[st.ItemID] {
			return fmt.Errorf("duplicate starting item %q", st.ItemID)
		}
		seen[st.ItemID] = true
		if st.Qty < 1 {
			return fmt.Errorf("starting item %q needs qty >= 1", st.ItemID)
		}
	}
	return nil
}

// ValidateCharacter checks a player's creation choices: every stat present and within [min, max],
// sum(v - base) == statPoints, exactly startingSkillPicks distinct known skills.
func ValidateCharacter(r Ruleset, stats map[string]int, skillIDs []string) error {
	if len(stats) != len(r.Stats) {
		return fmt.Errorf("need exactly %d stats", len(r.Stats))
	}
	spent := 0
	for _, s := range r.Stats {
		v, ok := stats[s.Key]
		if !ok {
			return fmt.Errorf("missing stat %q", s.Key)
		}
		if v < s.Min || v > s.Max {
			return fmt.Errorf("stat %q must be %d-%d", s.Key, s.Min, s.Max)
		}
		spent += v - s.Base
	}
	if spent != r.StatPoints {
		return fmt.Errorf("spend exactly %d stat points (spent %d)", r.StatPoints, spent)
	}
	if len(skillIDs) != r.StartingSkillPicks {
		return fmt.Errorf("pick exactly %d skills", r.StartingSkillPicks)
	}
	picked := map[string]bool{}
	for _, id := range skillIDs {
		if picked[id] {
			return fmt.Errorf("skill %q picked twice", id)
		}
		picked[id] = true
		if !r.HasSkill(id) {
			return fmt.Errorf("unknown skill %q", id)
		}
	}
	return nil
}

func (r Ruleset) HasStat(key string) bool {
	for _, s := range r.Stats {
		if s.Key == key {
			return true
		}
	}
	return false
}

func (r Ruleset) HasSkill(id string) bool {
	for _, s := range r.Skills {
		if s.ID == id {
			return true
		}
	}
	return false
}

func (r Ruleset) HasItem(id string) bool {
	for _, it := range r.Items {
		if it.ID == id {
			return true
		}
	}
	return false
}
