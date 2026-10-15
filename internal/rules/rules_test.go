package rules

import (
	"strings"
	"testing"
)

func TestPresetDnD5e(t *testing.T) {
	r := PresetDnD5e()
	if err := ValidateRuleset(r); err != nil {
		t.Fatalf("preset invalid: %v", err)
	}
	if len(r.Skills) != 18 || r.StatPoints != 27 || r.StartingSkillPicks != 2 {
		t.Fatalf("preset shape: %d skills, %d points, %d picks", len(r.Skills), r.StatPoints, r.StartingSkillPicks)
	}
	if !strings.Contains(r.Attribution, `System Reference Document 5.2.1 ("SRD 5.2.1")`) {
		t.Fatalf("attribution missing: %q", r.Attribution)
	}
	for _, s := range r.Stats {
		switch s.Key {
		case "hp", "ac":
			if s.Max != s.Base {
				t.Errorf("%s: max %d != base %d", s.Key, s.Max, s.Base)
			}
		default:
			if s.Base != 8 || s.Min != 8 || s.Max != 15 {
				t.Errorf("%s: %+v", s.Key, s)
			}
		}
	}
	// Copies are independent.
	a, b := PresetDnD5e(), PresetDnD5e()
	a.Stats[0].Max = 99
	if b.Stats[0].Max == 99 {
		t.Fatal("preset shares slices between calls")
	}
}

func TestCampaigns(t *testing.T) {
	c := Campaigns()
	if len(c) != 2 || c[0].ID != "starter-set" || c[0].Status != "available" || c[1].ID != "salem" || c[1].Status != "coming_soon" {
		t.Fatalf("campaigns: %+v", c)
	}
}

func TestValidateRuleset(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(r *Ruleset)
		ok     bool
	}{
		{"preset", func(r *Ruleset) {}, true},
		{"empty custom", func(r *Ruleset) { *r = Ruleset{Name: "Custom"} }, true},
		{"no name", func(r *Ruleset) { r.Name = "" }, false},
		{"dup stat", func(r *Ruleset) { r.Stats[1].Key = "str" }, false},
		{"base above max", func(r *Ruleset) { r.Stats[0].Base = 16 }, false},
		{"min above base", func(r *Ruleset) { r.Stats[0].Min = 9 }, false},
		{"negative points", func(r *Ruleset) { r.StatPoints = -1 }, false},
		{"negative picks", func(r *Ruleset) { r.StartingSkillPicks = -1 }, false},
		{"points unreachable", func(r *Ruleset) { r.StatPoints = 43 }, false},
		{"picks exceed skills", func(r *Ruleset) { r.StartingSkillPicks = 19 }, false},
		{"dup skill", func(r *Ruleset) { r.Skills[1].ID = "athletics" }, false},
		{"skill unknown stat", func(r *Ruleset) { r.Skills[0].StatKey = "luck" }, false},
		{"skill no stat ok", func(r *Ruleset) { r.Skills[0].StatKey = "" }, true},
		{"dup item", func(r *Ruleset) { r.Items[1].ID = "backpack" }, false},
		{"starting unknown item", func(r *Ruleset) { r.StartingItems[0].ItemID = "sword" }, false},
		{"starting zero qty", func(r *Ruleset) { r.StartingItems[0].Qty = 0 }, false},
		{"long stat name", func(r *Ruleset) { r.Stats[0].Name = strings.Repeat("x", 65) }, false},
		{"skill desc 500 runes", func(r *Ruleset) { r.Skills[0].Description = strings.Repeat("ก", 500) }, true},
		{"skill desc 501 runes", func(r *Ruleset) { r.Skills[0].Description = strings.Repeat("ก", 501) }, false},
		{"skill desc newline ok", func(r *Ruleset) { r.Skills[0].Description = "line 1\nline 2" }, true},
		{"skill desc control char", func(r *Ruleset) { r.Skills[0].Description = "bad\x00" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PresetDnD5e()
			tt.mutate(&r)
			if err := ValidateRuleset(r); (err == nil) != tt.ok {
				t.Fatalf("ok=%v, err=%v", tt.ok, err)
			}
		})
	}
}

func TestNormalizeTrimsSkillDescription(t *testing.T) {
	r := PresetDnD5e()
	r.Skills[0].Description = "  " + strings.Repeat("x", 500) + " \n "
	r.Normalize()
	if err := ValidateRuleset(r); err != nil || len(r.Skills[0].Description) != 500 {
		t.Fatalf("trimmed len %d, err %v", len(r.Skills[0].Description), err)
	}
	for _, s := range PresetDnD5e().Skills {
		if s.Description != "" {
			t.Fatalf("preset skill %q has a description", s.ID)
		}
	}
}

func validStats() map[string]int {
	// 7+7+5+4+2+2 = 27 points above base 8.
	return map[string]int{"str": 15, "dex": 15, "con": 13, "int": 12, "wis": 10, "cha": 10, "hp": 10, "ac": 10}
}

func TestValidateCharacter(t *testing.T) {
	r := PresetDnD5e()
	tests := []struct {
		name   string
		stats  func(m map[string]int)
		skills []string
		ok     bool
	}{
		{"valid", func(m map[string]int) {}, []string{"stealth", "arcana"}, true},
		{"too few points", func(m map[string]int) { m["cha"] = 9 }, []string{"stealth", "arcana"}, false},
		{"too many points", func(m map[string]int) { m["cha"] = 11 }, []string{"stealth", "arcana"}, false},
		{"above max", func(m map[string]int) { m["str"] = 16; m["dex"] = 14 }, []string{"stealth", "arcana"}, false},
		{"below min", func(m map[string]int) { m["cha"] = 7; m["wis"] = 13 }, []string{"stealth", "arcana"}, false},
		{"points into hp", func(m map[string]int) { m["hp"] = 11; m["cha"] = 9 }, []string{"stealth", "arcana"}, false},
		{"missing stat", func(m map[string]int) { delete(m, "ac") }, []string{"stealth", "arcana"}, false},
		{"extra stat", func(m map[string]int) { m["luck"] = 0 }, []string{"stealth", "arcana"}, false},
		{"one skill", func(m map[string]int) {}, []string{"stealth"}, false},
		{"three skills", func(m map[string]int) {}, []string{"stealth", "arcana", "history"}, false},
		{"duplicate skill", func(m map[string]int) {}, []string{"stealth", "stealth"}, false},
		{"unknown skill", func(m map[string]int) {}, []string{"stealth", "flying"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validStats()
			tt.stats(m)
			if err := ValidateCharacter(r, m, tt.skills); (err == nil) != tt.ok {
				t.Fatalf("ok=%v, err=%v", tt.ok, err)
			}
		})
	}
}
