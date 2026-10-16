package rules

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

// Class is one of the 12 base classes (docs/design/classes/<id>.md, SRD 5.2.1 based). Read-only game data.
type Class struct {
	ID           string       `json:"id"`
	Name         L            `json:"name"`
	HitDie       int          `json:"hitDie"`  // 6 | 8 | 10 | 12
	Primary      []string     `json:"primary"` // ability ids
	Saves        []string     `json:"saves"`   // 2 ability ids
	Armor        L            `json:"armor"`
	Weapons      L            `json:"weapons"`
	Tools        L            `json:"tools"` // {"",""} = none
	SkillChoices SkillChoices `json:"skillChoices"`
	Caster       string       `json:"caster"`      // "none" | "full" | "half" | "pact"
	CastAbility  string       `json:"castAbility"` // ability id; "" when caster is "none"
	Tables       []LevelTable `json:"tables"`
	Resources    []Resource   `json:"resources"`
	Features     []Feature    `json:"features"`
}

type SkillChoices struct {
	Count int      `json:"count"`
	From  []string `json:"from"` // skill ids
}

// LevelTable is a display-only value per class level; Values[i] is level i+1 and "—" means none yet.
type LevelTable struct {
	ID     string   `json:"id"`
	Name   L        `json:"name"`
	Values []string `json:"values"` // 20 entries
}

// Resource is a use counter. Its max comes from Max (20 entries, level i+1) or, when Max is nil, MaxFormula.
type Resource struct {
	ID             string `json:"id"`
	Name           L      `json:"name"`
	Max            []int  `json:"max"`            // nil = use MaxFormula
	MaxFormula     string `json:"maxFormula"`     // "abilityMod:<ability>" (min 1) | "level*<n>" | ""
	Refill         string `json:"refill"`         // "short" | "long": rest that refills it fully
	ShortRegain    int    `json:"shortRegain"`    // uses a short rest gives back when Refill is "long"
	ShortFromLevel int    `json:"shortFromLevel"` // from this level a short rest refills it fully; 0 = never
	FromLevel      int    `json:"fromLevel"`      // a MaxFormula resource exists from this level; 0 = level 1
}

type Feature struct {
	Level int `json:"level"`
	Name  L   `json:"name"`
	Desc  L   `json:"desc"`
}

// Abilities are the six ability scores (1-30).
type Abilities struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
}

// Score returns the score for ability id ("str".."cha") and false for an unknown id.
func (a Abilities) Score(id string) (int, bool) {
	switch id {
	case "str":
		return a.Str, true
	case "dex":
		return a.Dex, true
	case "con":
		return a.Con, true
	case "int":
		return a.Int, true
	case "wis":
		return a.Wis, true
	case "cha":
		return a.Cha, true
	}
	return 0, false
}

func (a Abilities) valid() bool {
	for _, v := range []int{a.Str, a.Dex, a.Con, a.Int, a.Wis, a.Cha} {
		if v < MinAbility || v > MaxAbility {
			return false
		}
	}
	return true
}

// Skill maps a skill id to the ability it uses.
type Skill struct {
	ID      string `json:"id"`
	Ability string `json:"ability"`
}

// PactSlots is the Warlock's Pact Magic: Count slots, all of spell level Level.
type PactSlots struct {
	Count int `json:"count"`
	Level int `json:"level"`
}

const (
	MinLevel, MaxLevel     = 1, 20
	MinAbility, MaxAbility = 1, 30
	MinAC, MaxAC           = 1, 30
)

// AbilityIDs lists the six abilities in sheet order.
var AbilityIDs = []string{"str", "dex", "con", "int", "wis", "cha"}

var skills = []Skill{
	{"athletics", "str"},
	{"acrobatics", "dex"}, {"sleight-of-hand", "dex"}, {"stealth", "dex"},
	{"arcana", "int"}, {"history", "int"}, {"investigation", "int"}, {"nature", "int"}, {"religion", "int"},
	{"animal-handling", "wis"}, {"insight", "wis"}, {"medicine", "wis"}, {"perception", "wis"}, {"survival", "wis"},
	{"deception", "cha"}, {"intimidation", "cha"}, {"performance", "cha"}, {"persuasion", "cha"},
}

// Skills returns the 18 skills in sheet order (grouped by ability).
func Skills() []Skill { return slices.Clone(skills) }

// SkillAbility returns the ability id skill id uses.
func SkillAbility(id string) (string, bool) {
	if i := slices.IndexFunc(skills, func(s Skill) bool { return s.ID == id }); i >= 0 {
		return skills[i].Ability, true
	}
	return "", false
}

// Shared tables from docs/design/classes/_shared.md (TestSharedTablesMatchDoc). Index = class level - 1.
var (
	proficiency = [MaxLevel]int{2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6}
	fullSlots   = [MaxLevel][9]int{
		{2}, {3}, {4, 2}, {4, 3}, {4, 3, 2}, {4, 3, 3}, {4, 3, 3, 1}, {4, 3, 3, 2}, {4, 3, 3, 3, 1}, {4, 3, 3, 3, 2},
		{4, 3, 3, 3, 2, 1}, {4, 3, 3, 3, 2, 1}, {4, 3, 3, 3, 2, 1, 1}, {4, 3, 3, 3, 2, 1, 1},
		{4, 3, 3, 3, 2, 1, 1, 1}, {4, 3, 3, 3, 2, 1, 1, 1}, {4, 3, 3, 3, 2, 1, 1, 1, 1},
		{4, 3, 3, 3, 3, 1, 1, 1, 1}, {4, 3, 3, 3, 3, 2, 1, 1, 1}, {4, 3, 3, 3, 3, 2, 2, 1, 1},
	}
	halfSlots = [MaxLevel][9]int{
		{2}, {2}, {3}, {3}, {4, 2}, {4, 2}, {4, 3}, {4, 3}, {4, 3, 2}, {4, 3, 2},
		{4, 3, 3}, {4, 3, 3}, {4, 3, 3, 1}, {4, 3, 3, 1}, {4, 3, 3, 2}, {4, 3, 3, 2},
		{4, 3, 3, 3, 1}, {4, 3, 3, 3, 1}, {4, 3, 3, 3, 2}, {4, 3, 3, 3, 2},
	}
	pactSlots = [MaxLevel]PactSlots{
		{1, 1}, {2, 1}, {2, 2}, {2, 2}, {2, 3}, {2, 3}, {2, 4}, {2, 4}, {2, 5}, {2, 5},
		{3, 5}, {3, 5}, {3, 5}, {3, 5}, {3, 5}, {3, 5}, {4, 5}, {4, 5}, {4, 5}, {4, 5},
	}
	classes = classList() // never mutated
)

// Classes returns the 12 classes. Shared data: callers must not mutate it.
func Classes() []Class { return classes }

// ClassByID returns class id. Its slices are shared: callers must not mutate them.
func ClassByID(id string) (Class, bool) {
	if i := slices.IndexFunc(classes, func(c Class) bool { return c.ID == id }); i >= 0 {
		return classes[i], true
	}
	return Class{}, false
}

// idx turns a class level into a table index, clamping to 1..20 so a bad level never panics.
func idx(level int) int { return min(max(level, MinLevel), MaxLevel) - 1 }

// AbilityMod is floor((score-10)/2).
func AbilityMod(score int) int {
	d := score - 10
	if d < 0 {
		return (d - 1) / 2
	}
	return d / 2
}

// ProfBonus is the proficiency bonus at class level (clamped to 1..20).
func ProfBonus(level int) int { return proficiency[idx(level)] }

// HPMax is the full hit die + CON mod at level 1, then the rounded-up average + CON mod per level; each level
// gives at least 1.
func HPMax(c Class, level, conScore int) int {
	mod, n := AbilityMod(conScore), idx(level)
	return max(1, c.HitDie+mod) + n*max(1, c.HitDie/2+1+mod)
}

// SlotsMax returns the spell slots per spell level (index 0 = 1st level) for full and half casters.
// Pact casters get zeros here; see PactMax.
func SlotsMax(c Class, level int) [9]int {
	switch c.Caster {
	case "full":
		return fullSlots[idx(level)]
	case "half":
		return halfSlots[idx(level)]
	}
	return [9]int{}
}

// PactMax returns a pact caster's Pact Magic slots; zero for every other class.
func PactMax(c Class, level int) PactSlots {
	if c.Caster != "pact" {
		return PactSlots{}
	}
	return pactSlots[idx(level)]
}

// ResourceMax is the counter's max at class level: 0 before FromLevel, else Max[level-1] or MaxFormula.
// An unknown formula gives 0 (TestClassesValid keeps the data free of them).
func ResourceMax(r Resource, level int, a Abilities) int {
	if level < r.FromLevel {
		return 0
	}
	if r.Max != nil {
		return r.Max[idx(level)]
	}
	if ab, ok := strings.CutPrefix(r.MaxFormula, "abilityMod:"); ok {
		if s, ok := a.Score(ab); ok {
			return max(1, AbilityMod(s))
		}
	} else if n, ok := strings.CutPrefix(r.MaxFormula, "level*"); ok {
		if k, err := strconv.Atoi(n); err == nil && k > 0 {
			return (idx(level) + 1) * k
		}
	}
	return 0
}

// SaveBonus is the ability mod + PB when ability is in saves (the proficient saves, see ProficientSaves).
// Unknown ability = 0.
func SaveBonus(a Abilities, ability string, level int, saves []string) int {
	s, ok := a.Score(ability)
	if !ok {
		return 0
	}
	if slices.Contains(saves, ability) {
		return AbilityMod(s) + ProfBonus(level)
	}
	return AbilityMod(s)
}

// ProficientSaves is the class saves plus the extra ones (class features such as Slippery Mind), in sheet order.
func ProficientSaves(c Class, extra []string) []string {
	out := []string{}
	for _, ab := range AbilityIDs {
		if slices.Contains(c.Saves, ab) || slices.Contains(extra, ab) {
			out = append(out, ab)
		}
	}
	return out
}

// CleanSaves checks extra save proficiencies: known ability ids. Duplicates are dropped; the result is fresh,
// non-nil and in sheet order.
func CleanSaves(saves []string) ([]string, error) {
	for _, ab := range saves {
		if !slices.Contains(AbilityIDs, ab) {
			return nil, errors.New("unknown ability id")
		}
	}
	return ProficientSaves(Class{}, saves), nil
}

// SkillBonus is the skill's ability mod + PB if proficient, + 2×PB with expertise. Unknown skill = 0.
func SkillBonus(a Abilities, skill string, level int, proficient, expertise []string) int {
	ab, ok := SkillAbility(skill)
	if !ok {
		return 0
	}
	s, _ := a.Score(ab)
	switch {
	case slices.Contains(expertise, skill):
		return AbilityMod(s) + 2*ProfBonus(level)
	case slices.Contains(proficient, skill):
		return AbilityMod(s) + ProfBonus(level)
	}
	return AbilityMod(s)
}

// Features returns the class features gained at or below level, in data order. Shared: do not mutate.
func Features(c Class, level int) []Feature {
	i := slices.IndexFunc(c.Features, func(f Feature) bool { return f.Level > level })
	if i < 0 {
		return c.Features
	}
	return c.Features[:i:i]
}
