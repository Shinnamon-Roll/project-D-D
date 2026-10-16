package rules

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"
)

var classIDs = []string{"barbarian", "bard", "cleric", "druid", "fighter", "monk", "paladin", "ranger", "rogue", "sorcerer", "warlock", "wizard"}

func TestClassesMatchDoc(t *testing.T) {
	syncGo(t, "classes_data.go", "docs/design/classes/*.md", "classList", classList(), docClasses(t))
}

func TestSharedTablesMatchDoc(t *testing.T) {
	var doc struct {
		Proficiency [MaxLevel]int `json:"proficiency"`
		Slots       struct {
			Full [MaxLevel][9]int    `json:"full"`
			Half [MaxLevel][9]int    `json:"half"`
			Pact [MaxLevel]PactSlots `json:"pact"`
		} `json:"slots"`
	}
	docJSON(t, "classes/_shared.md", &doc)
	if doc.Proficiency != proficiency || doc.Slots.Full != fullSlots || doc.Slots.Half != halfSlots || doc.Slots.Pact != pactSlots {
		t.Fatal("shared tables in classes.go differ from docs/design/classes/_shared.md")
	}
}

// TestClassesValid checks the class data invariants the helpers and the frontend rely on.
func TestClassesValid(t *testing.T) {
	cs := Classes()
	if got := slices.Collect(func(yield func(string) bool) {
		for _, c := range cs {
			yield(c.ID)
		}
	}); !slices.Equal(got, classIDs) {
		t.Fatalf("class ids %v", got)
	}
	isAb := func(id string) bool { return slices.Contains(AbilityIDs, id) }
	for _, c := range cs {
		fail := func(f string, a ...any) { t.Helper(); t.Errorf(c.ID+": "+f, a...) }
		if !slices.Contains([]int{6, 8, 10, 12}, c.HitDie) {
			fail("hitDie %d", c.HitDie)
		}
		if c.Name.TH == "" || c.Name.EN == "" || len(c.Primary) == 0 || !slices.ContainsFunc(c.Primary, isAb) {
			fail("name or primary")
		}
		for _, a := range c.Primary {
			if !isAb(a) {
				fail("primary %q", a)
			}
		}
		if len(c.Saves) != 2 || !isAb(c.Saves[0]) || !isAb(c.Saves[1]) || c.Saves[0] == c.Saves[1] {
			fail("saves %v", c.Saves)
		}
		switch {
		case !slices.Contains([]string{"none", "full", "half", "pact"}, c.Caster):
			fail("caster %q", c.Caster)
		case (c.Caster == "none") != (c.CastAbility == ""), c.CastAbility != "" && !isAb(c.CastAbility):
			fail("castAbility %q for caster %q", c.CastAbility, c.Caster)
		}
		sc := c.SkillChoices
		if sc.Count < 1 || sc.Count > len(sc.From) {
			fail("skill count %d of %d", sc.Count, len(sc.From))
		}
		for i, s := range sc.From {
			if _, ok := SkillAbility(s); !ok || slices.Contains(sc.From[:i], s) {
				fail("skill choice %q", s)
			}
		}
		ids := map[string]bool{}
		for _, tb := range c.Tables {
			if ids["t:"+tb.ID] || tb.ID == "" || len(tb.Values) != MaxLevel {
				fail("table %q: dup/empty id or %d values", tb.ID, len(tb.Values))
			}
			ids["t:"+tb.ID] = true
			for _, v := range tb.Values {
				if v == "" || v == "-" || v == "–" {
					fail("table %s: placeholder %q, use —", tb.ID, v)
				}
			}
		}
		for _, r := range c.Resources {
			if ids["r:"+r.ID] || r.ID == "" || r.Name.TH == "" || r.Name.EN == "" {
				fail("resource %q: dup/empty", r.ID)
			}
			ids["r:"+r.ID] = true
			if (r.Max == nil) == (r.MaxFormula == "") || (r.Max != nil && len(r.Max) != MaxLevel) {
				fail("resource %s: need exactly one of 20 max values or a formula", r.ID)
			}
			if r.Max != nil && r.FromLevel != 0 {
				fail("resource %s: fromLevel is only for formulas (use 0 in max)", r.ID)
			}
			if r.MaxFormula != "" && ResourceMax(r, MaxLevel, Abilities{10, 10, 10, 10, 10, 10}) < 1 {
				fail("resource %s: bad formula %q", r.ID, r.MaxFormula)
			}
			if r.Refill != "short" && r.Refill != "long" || r.ShortRegain < 0 || r.FromLevel < 0 || r.FromLevel > MaxLevel ||
				r.ShortFromLevel < 0 || r.ShortFromLevel > MaxLevel {
				fail("resource %s: refill %q shortRegain %d fromLevel %d shortFromLevel %d", r.ID, r.Refill, r.ShortRegain, r.FromLevel, r.ShortFromLevel)
			}
		}
		for i, f := range c.Features {
			if f.Level < MinLevel || f.Level > MaxLevel || (i > 0 && f.Level < c.Features[i-1].Level) {
				fail("feature %q level %d out of order", f.Name.EN, f.Level)
			}
			if f.Name.TH == "" || f.Name.EN == "" || f.Desc.TH == "" || f.Desc.EN == "" {
				fail("feature %q at %d has empty text", f.Name.EN, f.Level)
			}
		}
		if len(c.Features) == 0 || c.Features[0].Level != 1 || c.Features[len(c.Features)-1].Level != MaxLevel {
			fail("features must cover level 1 and 20")
		}
	}
}

// TestClassJSONSize logs the wire size of all class data (decides whether RoomState carries it or only hello).
func TestClassJSONSize(t *testing.T) {
	b, err := json.Marshal(Classes())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("12 classes marshal to %d bytes", len(b))
	if len(b) > 1<<20 {
		t.Fatalf("class JSON %d bytes", len(b))
	}
}

func TestAbilityMod(t *testing.T) {
	for _, tt := range []struct{ score, want int }{
		{1, -5}, {2, -4}, {3, -4}, {7, -2}, {8, -1}, {9, -1}, {10, 0}, {11, 0}, {12, 1}, {15, 2}, {20, 5}, {29, 9}, {30, 10},
	} {
		if got := AbilityMod(tt.score); got != tt.want {
			t.Errorf("AbilityMod(%d) = %d, want %d", tt.score, got, tt.want)
		}
	}
}

func TestProfBonus(t *testing.T) {
	for _, tt := range []struct{ level, want int }{{-3, 2}, {0, 2}, {1, 2}, {4, 2}, {5, 3}, {8, 3}, {9, 4}, {13, 5}, {16, 5}, {17, 6}, {20, 6}, {99, 6}} {
		if got := ProfBonus(tt.level); got != tt.want {
			t.Errorf("ProfBonus(%d) = %d, want %d", tt.level, got, tt.want)
		}
	}
}

func mustClass(t *testing.T, id string) Class {
	t.Helper()
	c, ok := ClassByID(id)
	if !ok {
		t.Fatalf("no class %s", id)
	}
	return c
}

func TestClassByID(t *testing.T) {
	if _, ok := ClassByID("knight"); ok {
		t.Fatal("knight is a hero, not a class")
	}
	if _, ok := ClassByID(""); ok {
		t.Fatal("empty id")
	}
}

func TestHPMax(t *testing.T) {
	for _, tt := range []struct {
		class      string
		level, con int
		want       int
	}{
		{"wizard", 1, 10, 6},     // d6
		{"wizard", 2, 10, 10},    // 6 + 4
		{"fighter", 1, 15, 12},   // 10 + 2
		{"fighter", 5, 15, 44},   // 12 + 4×(6+2)
		{"barbarian", 1, 16, 15}, // 12 + 3
		{"barbarian", 20, 16, 15 + 19*10},
		{"rogue", 3, 8, 7 + 2*4}, // d8: 8-1, then 5-1
		{"wizard", 3, 1, 1 + 2},  // CON 1 (-5): every level gives at least 1
		{"cleric", 0, 10, 8},     // level clamps to 1
		{"cleric", 25, 10, 8 + 19*5},
	} {
		if got := HPMax(mustClass(t, tt.class), tt.level, tt.con); got != tt.want {
			t.Errorf("HPMax(%s, %d, con %d) = %d, want %d", tt.class, tt.level, tt.con, got, tt.want)
		}
	}
}

func TestSlotsMax(t *testing.T) {
	for _, tt := range []struct {
		class string
		level int
		want  [9]int
		pact  PactSlots
	}{
		{"wizard", 1, [9]int{2}, PactSlots{}},
		{"bard", 5, [9]int{4, 3, 2}, PactSlots{}},
		{"cleric", 20, [9]int{4, 3, 3, 3, 3, 2, 2, 1, 1}, PactSlots{}},
		{"paladin", 1, [9]int{2}, PactSlots{}},
		{"ranger", 9, [9]int{4, 3, 2}, PactSlots{}},
		{"ranger", 20, [9]int{4, 3, 3, 3, 2}, PactSlots{}},
		{"fighter", 20, [9]int{}, PactSlots{}},
		{"rogue", 5, [9]int{}, PactSlots{}},
		{"warlock", 1, [9]int{}, PactSlots{1, 1}},
		{"warlock", 5, [9]int{}, PactSlots{2, 3}},
		{"warlock", 11, [9]int{}, PactSlots{3, 5}},
		{"warlock", 20, [9]int{}, PactSlots{4, 5}},
		{"wizard", 0, [9]int{2}, PactSlots{}}, // clamps
	} {
		c := mustClass(t, tt.class)
		if got := SlotsMax(c, tt.level); got != tt.want {
			t.Errorf("SlotsMax(%s, %d) = %v, want %v", tt.class, tt.level, got, tt.want)
		}
		if got := PactMax(c, tt.level); got != tt.pact {
			t.Errorf("PactMax(%s, %d) = %v, want %v", tt.class, tt.level, got, tt.pact)
		}
	}
}

func resource(t *testing.T, class, id string) Resource {
	t.Helper()
	c := mustClass(t, class)
	i := slices.IndexFunc(c.Resources, func(r Resource) bool { return r.ID == id })
	if i < 0 {
		t.Fatalf("%s has no resource %s", class, id)
	}
	return c.Resources[i]
}

func TestResourceMax(t *testing.T) {
	ab := func(wis, cha int) Abilities { return Abilities{10, 10, 10, 10, wis, cha} }
	for _, tt := range []struct {
		class, res string
		level      int
		a          Abilities
		want       int
	}{
		{"barbarian", "rage", 1, ab(10, 10), 2},
		{"barbarian", "rage", 20, ab(10, 10), 6},
		{"barbarian", "persistent-rage", 14, ab(10, 10), 0},
		{"barbarian", "persistent-rage", 15, ab(10, 10), 1},
		{"bard", "bardic-inspiration", 1, ab(10, 16), 3},
		{"bard", "bardic-inspiration", 1, ab(10, 8), 1}, // formula minimum 1
		{"paladin", "lay-on-hands", 1, ab(10, 10), 5},
		{"paladin", "lay-on-hands", 20, ab(10, 10), 100},
		{"ranger", "tireless", 9, ab(18, 10), 0}, // fromLevel 10
		{"ranger", "tireless", 10, ab(18, 10), 4},
		{"ranger", "tireless", 10, ab(6, 10), 1},
		{"ranger", "natures-veil", 13, ab(14, 10), 0},
		{"ranger", "natures-veil", 14, ab(14, 10), 2},
		{"fighter", "action-surge", 1, ab(10, 10), 0},
		{"fighter", "second-wind", 0, ab(10, 10), 2}, // level clamps
	} {
		if got := ResourceMax(resource(t, tt.class, tt.res), tt.level, tt.a); got != tt.want {
			t.Errorf("ResourceMax(%s/%s, %d) = %d, want %d", tt.class, tt.res, tt.level, got, tt.want)
		}
	}
	for _, f := range []string{"", "abilityMod:luck", "level*x", "level*0", "level*-2", "nope"} {
		if got := ResourceMax(Resource{MaxFormula: f}, 5, ab(20, 20)); got != 0 {
			t.Errorf("formula %q = %d, want 0", f, got)
		}
	}
}

func TestSaveBonus(t *testing.T) {
	fighter := mustClass(t, "fighter") // saves str, con
	a := Abilities{Str: 16, Dex: 14, Con: 12, Int: 10, Wis: 8, Cha: 1}
	for _, tt := range []struct {
		ability string
		level   int
		want    int
	}{
		{"str", 1, 3 + 2}, {"con", 5, 1 + 3}, {"dex", 1, 2}, {"wis", 20, -1}, {"cha", 9, -5}, {"luck", 1, 0}, {"", 1, 0},
	} {
		if got := SaveBonus(a, tt.ability, tt.level, fighter.Saves); got != tt.want {
			t.Errorf("SaveBonus(%s, %d) = %d, want %d", tt.ability, tt.level, got, tt.want)
		}
	}
}

func TestProficientSaves(t *testing.T) {
	monk := mustClass(t, "monk") // saves str, dex
	for _, tt := range []struct {
		name  string
		extra []string
		want  []string
	}{
		{"class only", nil, []string{"str", "dex"}},
		{"slippery mind", []string{"wis", "cha"}, []string{"str", "dex", "wis", "cha"}},
		{"overlap and order", []string{"cha", "str", "con"}, []string{"str", "dex", "con", "cha"}},
		{"disciplined survivor", AbilityIDs, AbilityIDs},
	} {
		got := ProficientSaves(monk, tt.extra)
		if !slices.Equal(got, tt.want) {
			t.Errorf("%s: %v, want %v", tt.name, got, tt.want)
		}
	}
	a := Abilities{Str: 10, Dex: 10, Con: 10, Int: 10, Wis: 14, Cha: 10}
	if got := SaveBonus(a, "wis", 15, ProficientSaves(monk, []string{"wis"})); got != 2+5 {
		t.Errorf("extra wis save at 15 = %d, want 7", got)
	}
	if got := SaveBonus(a, "wis", 15, monk.Saves); got != 2 {
		t.Errorf("class-only wis save = %d, want 2", got)
	}
}

func TestCleanSaves(t *testing.T) {
	for _, tt := range []struct {
		in      []string
		want    []string
		wantErr bool
	}{
		{nil, []string{}, false},
		{[]string{"cha", "wis", "cha"}, []string{"wis", "cha"}, false},
		{[]string{"str", "dex", "con", "int", "wis", "cha", "str"}, AbilityIDs, false},
		{[]string{"luck"}, nil, true},
		{[]string{"WIS"}, nil, true},
		{[]string{""}, nil, true},
	} {
		got, err := CleanSaves(tt.in)
		if (err != nil) != tt.wantErr || (err != nil && err.Error() != "unknown ability id") {
			t.Errorf("CleanSaves(%v) err %v", tt.in, err)
			continue
		}
		if !tt.wantErr && (got == nil || !slices.Equal(got, tt.want)) {
			t.Errorf("CleanSaves(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSkills(t *testing.T) {
	sk := Skills()
	if len(sk) != 18 {
		t.Fatalf("%d skills", len(sk))
	}
	count := map[string]int{}
	for _, s := range sk {
		count[s.Ability]++
	}
	if !reflect.DeepEqual(count, map[string]int{"str": 1, "dex": 3, "int": 5, "wis": 5, "cha": 4}) {
		t.Fatalf("skills per ability %v", count)
	}
	sk[0].ID = "x" // a copy
	if Skills()[0].ID != "athletics" {
		t.Fatal("Skills shares its slice")
	}
	a := Abilities{Str: 8, Dex: 16, Con: 10, Int: 13, Wis: 12, Cha: 15}
	prof, exp := []string{"stealth", "deception", "arcana"}, []string{"deception"}
	for _, tt := range []struct {
		skill string
		level int
		want  int
	}{
		{"athletics", 1, -1},      // not proficient
		{"stealth", 1, 3 + 2},     // proficient
		{"stealth", 17, 3 + 6},    // PB 6
		{"deception", 1, 2 + 4},   // expertise 2×PB
		{"deception", 9, 2 + 8},   // 2×4
		{"arcana", 5, 1 + 3},      // INT
		{"perception", 1, 1},      // WIS
		{"Stealth", 1, 0},         // unknown id
		{"sleight-of-hand", 1, 3}, // DEX, not proficient
	} {
		if got := SkillBonus(a, tt.skill, tt.level, prof, exp); got != tt.want {
			t.Errorf("SkillBonus(%s, %d) = %d, want %d", tt.skill, tt.level, got, tt.want)
		}
	}
}

func TestFeatures(t *testing.T) {
	b := mustClass(t, "barbarian")
	names := func(fs []Feature) string {
		var s []string
		for _, f := range fs {
			s = append(s, f.Name.EN)
		}
		return strings.Join(s, ",")
	}
	if got := names(Features(b, 1)); got != "Rage,Unarmored Defense,Weapon Mastery" {
		t.Fatalf("L1: %s", got)
	}
	if got := Features(b, 2); len(got) != 5 || got[4].Name.EN != "Reckless Attack" {
		t.Fatalf("L2: %s", names(got))
	}
	if got := Features(b, 20); len(got) != len(b.Features) {
		t.Fatalf("L20: %d of %d", len(got), len(b.Features))
	}
	if got := Features(b, 0); len(got) != 0 {
		t.Fatalf("L0: %s", names(got))
	}
	if got := Features(b, 1); cap(got) != len(got) {
		t.Fatal("appending to Features would overwrite class data")
	}
}
