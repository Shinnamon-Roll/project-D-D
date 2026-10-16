package room

import (
	"slices"
	"strings"
	"testing"

	"projectdnd/internal/rules"
)

// Rules v2 events: character.level, character.setAbilities, slot.set, resource.set, hitDice.set, party.rest,
// ability.add/update/remove. Hand-computed numbers: knight = fighter d10, CON 15 (+2): HP 12, then +8 a level.

func TestCharacterLevel(t *testing.T) {
	r, kn, wz := party(t)
	lv := func(id string, d any) m { return m{"characterId": id, "delta": d} }
	failWith(t, r, "p1", "character.level", lv(kn, 1), "only the DM can do that")
	failWith(t, r, "stranger", "character.level", lv(kn, 1), "only the DM can do that")
	failWith(t, r, "dm", "character.level", lv("nope", 1), "no such character")
	for _, d := range []any{0, 2, -2} {
		failWith(t, r, "dm", "character.level", lv(kn, d), "delta must be 1 or -1")
	}
	failWith(t, r, "dm", "character.level", lv(kn, "1"), "malformed payload")
	failWith(t, r, "dm", "character.level", lv(kn, -1), "level must be 1-20")

	ok(t, r, "dm", "character.adjust", m{"characterId": kn, "hp": -5}) // 7/12
	ok(t, r, "co", "character.level", lv(kn, 1))
	c := r.State().Characters[0]
	if c.Level != 2 || c.HPMax != 20 || c.HP != 15 || c.HitDiceMax != 2 || c.ResourcesMax["action-surge"] != 1 || c.PB != 2 {
		t.Fatalf("level 2: %+v", c)
	}
	if f := lastFeed(r); f.EN != "Aldric Vane reaches level 2" || f.TH != "อัลดริก เวน ขึ้นเลเวล 2" {
		t.Fatalf("feed %+v", f)
	}
	ok(t, r, "p1", "resource.set", m{"characterId": kn, "resourceId": "action-surge", "used": 1})
	ok(t, r, "p1", "hitDice.set", m{"characterId": kn, "used": 2})
	ok(t, r, "dm", "character.level", lv(kn, -1)) // max 0 again: used clamps
	c = r.State().Characters[0]
	if c.Level != 1 || c.HPMax != 12 || c.HP != 7 || c.ResourcesUsed["action-surge"] != 0 || c.HitDiceUsed != 1 {
		t.Fatalf("back to 1: %+v", c)
	}
	if f := lastFeed(r); f.EN != "Aldric Vane drops to level 1" || f.TH != "อัลดริก เวน ลดเป็นเลเวล 1" {
		t.Fatalf("feed %+v", f)
	}

	// Level down with less HP than the lost max: HP 0 and the knocked-out line.
	ok(t, r, "dm", "character.level", lv(kn, 1))                        // 15/20
	ok(t, r, "dm", "character.adjust", m{"characterId": kn, "hp": -12}) // 3/20
	ok(t, r, "dm", "character.level", lv(kn, -1))
	if c := r.State().Characters[0]; c.HP != 0 || c.HPMax != 12 || lastFeed(r).EN != "💀 Aldric Vane is knocked out (HP 0)" {
		t.Fatalf("knocked out: %d/%d %+v", c.HP, c.HPMax, lastFeed(r))
	}

	// Slots follow the level; used slots clamp when the max drops.
	ok(t, r, "dm", "character.level", lv(wz, 1))
	ok(t, r, "dm", "character.level", lv(wz, 1))
	w := r.State().Characters[1]
	if w.SlotsMax != [9]int{4, 2} || w.HPMax != 7+5+5 { // d6 average 4 + CON 1 a level
		t.Fatalf("wizard 3: %v %d", w.SlotsMax, w.HPMax)
	}
	ok(t, r, "p2", "slot.set", m{"characterId": wz, "level": 2, "used": 2})
	ok(t, r, "p2", "slot.set", m{"characterId": wz, "level": 1, "used": 4})
	ok(t, r, "dm", "character.level", lv(wz, -1))
	if w := r.State().Characters[1]; w.SlotsMax != [9]int{3} || w.SlotsUsed != [9]int{3} {
		t.Fatalf("wizard 2: max %v used %v", w.SlotsMax, w.SlotsUsed)
	}

	// Up to 20 and no further; PB 6 at 17+.
	for range 18 {
		ok(t, r, "dm", "character.level", lv(wz, 1))
	}
	failWith(t, r, "dm", "character.level", lv(wz, 1), "level must be 1-20")
	if w := r.State().Characters[1]; w.Level != 20 || w.PB != 6 || w.SlotsMax != [9]int{4, 3, 3, 3, 3, 2, 2, 1, 1} {
		t.Fatalf("wizard 20: %+v", w)
	}

	r, _ = lobby(t)
	failWith(t, r, "dm", "character.level", lv("c1", 1), "not allowed while the room is lobby")
}

func TestCharacterSetAbilities(t *testing.T) {
	ab := func(id string, a rules.Abilities, ac int) m { return m{"characterId": id, "abilities": a, "ac": ac} }
	base := rules.Abilities{Str: 17, Dex: 10, Con: 15, Int: 8, Wis: 13, Cha: 12}
	r, ids := partyOf(t, "knight", "bard")
	kn, bd := ids[0], ids[1]
	failWith(t, r, "p1", "character.setAbilities", ab(kn, base, 18), "only the DM can do that")
	failWith(t, r, "dm", "character.setAbilities", ab("nope", base, 18), "no such character")
	for _, bad := range []rules.Abilities{{}, {Str: 31, Dex: 10, Con: 15, Int: 8, Wis: 13, Cha: 12}, {Str: 17, Dex: 10, Con: 15, Int: 8, Wis: 13, Cha: 0}} {
		failWith(t, r, "dm", "character.setAbilities", ab(kn, bad, 18), "abilities must be 1-30")
	}
	failWith(t, r, "dm", "character.setAbilities", m{"characterId": kn, "ac": 18}, "abilities must be 1-30")
	failWith(t, r, "dm", "character.setAbilities", ab(kn, base, 0), "ac must be 1-30")
	failWith(t, r, "dm", "character.setAbilities", ab(kn, base, 31), "ac must be 1-30")

	ok(t, r, "dm", "character.adjust", m{"characterId": kn, "hp": -2}) // 10/12
	feed := len(r.State().Feed)
	low := base
	low.Con = 8 // −1: hpMax 10 − 1 = 9, keeps the 2 damage
	ok(t, r, "dm", "character.setAbilities", ab(kn, low, 20))
	c := r.State().Characters[0]
	if c.Abilities != low || c.AC != 20 || c.HPMax != 9 || c.HP != 7 {
		t.Fatalf("set: %+v", c)
	}
	if len(r.State().Feed) != feed {
		t.Fatal("setAbilities wrote to the feed")
	}
	if l := lastLog(r); l.Text.EN != "set Aldric Vane: STR 17, DEX 10, CON 8, INT 8, WIS 13, CHA 12, AC 20 (HP 7/9)" || l.Actor != "Dee" {
		t.Fatalf("log %+v", l)
	}

	// abilityMod resources follow the score (min 1) and used clamps.
	ok(t, r, "p2", "resource.set", m{"characterId": bd, "resourceId": "bardic-inspiration", "used": 3})
	b := r.State().Characters[1].Abilities
	b.Cha = 10
	ok(t, r, "dm", "character.setAbilities", ab(bd, b, 13))
	if c := r.State().Characters[1]; c.ResourcesMax["bardic-inspiration"] != 1 || c.ResourcesUsed["bardic-inspiration"] != 1 {
		t.Fatalf("bard: max %v used %v", c.ResourcesMax, c.ResourcesUsed)
	}
}

func TestCharacterSetSkills(t *testing.T) {
	sk := func(id string, skills, exp, saves any) m {
		return m{"characterId": id, "skills": skills, "expertise": exp, "saves": saves}
	}
	r, kn, _ := party(t) // knight = fighter: skills athletics, perception; saves str, con; WIS 13
	good := sk(kn, []string{"athletics"}, []string{}, []string{})
	for _, who := range []string{"p1", "p2", "stranger"} {
		failWith(t, r, who, "character.setSkills", good, "only the DM can do that")
	}
	for _, tt := range []struct {
		name    string
		payload any
		want    string
	}{
		{"missing character", sk("nope", []string{}, []string{}, []string{}), "no such character"},
		{"unknown skill", sk(kn, []string{"Stealth"}, nil, nil), "unknown skill id"},
		{"duplicate skill", sk(kn, []string{"stealth", "stealth"}, nil, nil), "duplicate skill"},
		{"expertise not proficient", sk(kn, []string{"athletics"}, []string{"stealth"}, nil), "expertise must be a proficient skill"},
		{"duplicate expertise", sk(kn, []string{"athletics"}, []string{"athletics", "athletics"}, nil), "duplicate expertise"},
		{"unknown save", sk(kn, []string{}, nil, []string{"luck"}), "unknown ability id"},
		{"save as skill id", sk(kn, []string{}, nil, []string{"perception"}), "unknown ability id"},
		{"saves not a list", sk(kn, []string{}, nil, "wis"), "malformed payload"},
		{"skills not strings", sk(kn, []int{1}, nil, nil), "malformed payload"},
	} {
		t.Run(tt.name, func(t *testing.T) { failWith(t, r, "dm", "character.setSkills", tt.payload, tt.want) })
	}

	feed := len(r.State().Feed)
	ok(t, r, "co", "character.setSkills",
		sk(kn, []string{"athletics", "perception", "stealth"}, []string{"perception"}, []string{"cha", "wis", "cha"}))
	c := r.State().Characters[0]
	if !slices.Equal(c.Skills, []string{"athletics", "perception", "stealth"}) || !slices.Equal(c.Expertise, []string{"perception"}) ||
		!slices.Equal(c.ExtraSaves, []string{"wis", "cha"}) || !slices.Equal(c.Saves, []string{"str", "con", "wis", "cha"}) {
		t.Fatalf("set: skills %v exp %v extra %v saves %v", c.Skills, c.Expertise, c.ExtraSaves, c.Saves)
	}
	if got := rules.SaveBonus(c.Abilities, "wis", c.Level, c.Saves); got != 1+2 {
		t.Fatalf("wis save %d, want 3", got)
	}
	if got := rules.SkillBonus(c.Abilities, "perception", c.Level, c.Skills, c.Expertise); got != 1+4 {
		t.Fatalf("perception %d, want 5", got)
	}
	if len(r.State().Feed) != feed {
		t.Fatal("setSkills wrote to the feed")
	}
	if l := lastLog(r); l.Text.EN != "set skills of Aldric Vane: athletics, perception, stealth; expertise: perception; extra saves: wis, cha" ||
		l.Text.TH != "ตั้ง skill ของ อัลดริก เวน: athletics, perception, stealth; expertise: perception; save เพิ่ม: wis, cha" || l.Actor != "Cora" {
		t.Fatalf("log %+v", l)
	}

	// A level change keeps the extra saves; PB follows the level.
	ok(t, r, "dm", "character.level", m{"characterId": kn, "delta": 1})
	if c := r.State().Characters[0]; !slices.Equal(c.Saves, []string{"str", "con", "wis", "cha"}) {
		t.Fatalf("after level: %v", c.Saves)
	}

	// Empty lists clear everything back to the class saves; missing keys count as empty.
	ok(t, r, "dm", "character.setSkills", m{"characterId": kn})
	c = r.State().Characters[0]
	if c.Skills == nil || len(c.Skills) != 0 || len(c.Expertise) != 0 || c.ExtraSaves == nil || len(c.ExtraSaves) != 0 ||
		!slices.Equal(c.Saves, []string{"str", "con"}) {
		t.Fatalf("cleared: %+v", c)
	}
	if lastLog(r).Text.EN != "set skills of Aldric Vane: —; expertise: —; extra saves: —" {
		t.Fatalf("log %+v", lastLog(r))
	}

	r, _ = lobby(t)
	failWith(t, r, "dm", "character.setSkills", sk("c1", []string{}, []string{}, []string{}), "not allowed while the room is lobby")
	r = creating(t) // allowed while others still pick
	ok(t, r, "p1", "character.pick", pick("knight"))
	ok(t, r, "dm", "character.setSkills", sk(r.State().Characters[0].ID, []string{"stealth"}, []string{"stealth"}, rules.AbilityIDs))
	if c := r.State().Characters[0]; !slices.Equal(c.Saves, rules.AbilityIDs) {
		t.Fatalf("all saves: %v", c.Saves)
	}
}

func TestSlotSet(t *testing.T) {
	r, ids := partyOf(t, "knight", "wizard", "legend-warlock")
	kn, wz, wl := ids[0], ids[1], ids[2]
	sl := func(id string, lv any, used any) m { return m{"characterId": id, "level": lv, "used": used} }
	for _, tc := range []struct {
		name, client string
		p            m
		want         string
	}{
		{"other player", "p1", sl(wz, 1, 1), "you can only change your own character"},
		{"stranger", "stranger", sl(wz, 1, 1), "you can only change your own character"},
		{"player, missing id", "p2", sl("nope", 1, 1), "you can only change your own character"},
		{"dm, missing id", "dm", sl("nope", 1, 1), "no such character"},
		{"level 0", "p2", sl(wz, 0, 1), `level must be 1-9 or "pact"`},
		{"level 10", "p2", sl(wz, 10, 1), `level must be 1-9 or "pact"`},
		{"level text", "p2", sl(wz, "one", 1), `level must be 1-9 or "pact"`},
		{"level missing", "p2", m{"characterId": wz, "used": 1}, `level must be 1-9 or "pact"`},
		{"level 2 not yet", "p2", sl(wz, 2, 0), "no such spell slot"},
		{"wizard has no pact", "p2", sl(wz, "pact", 0), "no such spell slot"},
		{"fighter has none", "p1", sl(kn, 1, 0), "no such spell slot"},
		{"warlock has no level-1 slots", "p3", sl(wl, 1, 0), "no such spell slot"},
		{"over max", "p2", sl(wz, 1, 3), "used must be 0-2"},
		{"negative", "p2", sl(wz, 1, -1), "used must be 0-2"},
		{"pact over max", "p3", sl(wl, "pact", 2), "used must be 0-1"},
		{"used float", "p2", sl(wz, 1, 1.5), "malformed payload"},
	} {
		t.Run(tc.name, func(t *testing.T) { failWith(t, r, tc.client, "slot.set", tc.p, tc.want) })
	}
	ok(t, r, "p2", "slot.set", sl(wz, 1, 2))
	if l := lastLog(r); l.Text.EN != "Soraya Quill: level 1 spell slots used 2/2" || l.Actor != "P2" {
		t.Fatalf("log %+v", l)
	}
	ok(t, r, "dm", "slot.set", sl(wz, 1, 1)) // the DM may fix it
	ok(t, r, "p3", "slot.set", sl(wl, "pact", 1))
	if l := lastLog(r); l.Text.EN != "Vesna Morrow: pact slots used 1/1" {
		t.Fatalf("log %+v", l)
	}
	cs := r.State().Characters
	if cs[1].SlotsUsed != [9]int{1} || cs[2].PactUsed != 1 || cs[2].SlotsUsed != [9]int{} {
		t.Fatalf("used %v %d", cs[1].SlotsUsed, cs[2].PactUsed)
	}
	r, _ = lobby(t)
	failWith(t, r, "dm", "slot.set", sl("c1", 1, 1), "not allowed while the room is lobby")
}

func TestResourceSet(t *testing.T) {
	r, kn, wz := party(t)
	rs := func(id, res string, used any) m { return m{"characterId": id, "resourceId": res, "used": used} }
	failWith(t, r, "p2", "resource.set", rs(kn, "second-wind", 1), "you can only change your own character")
	failWith(t, r, "p1", "resource.set", rs(kn, "rage", 1), "no such resource")         // not a fighter resource
	failWith(t, r, "p1", "resource.set", rs(kn, "action-surge", 0), "no such resource") // from level 2
	failWith(t, r, "p1", "resource.set", rs(kn, "", 0), "no such resource")
	failWith(t, r, "p1", "resource.set", rs(kn, "second-wind", 3), "used must be 0-2")
	failWith(t, r, "p1", "resource.set", rs(kn, "second-wind", -1), "used must be 0-2")
	failWith(t, r, "p2", "resource.set", rs(wz, "second-wind", 1), "no such resource")

	ok(t, r, "p1", "resource.set", rs(kn, "second-wind", 2))
	fighter, _ := rules.ClassByID("fighter")
	name := fighter.Resources[0].Name
	if l := lastLog(r); l.Text.EN != "Aldric Vane: "+name.EN+" used 2/2" || l.Text.TH != "อัลดริก เวน: ใช้ "+name.TH+" 2/2" {
		t.Fatalf("log %+v", l)
	}
	ok(t, r, "dm", "resource.set", rs(kn, "second-wind", 0))
	if u := r.State().Characters[0].ResourcesUsed["second-wind"]; u != 0 {
		t.Fatalf("used %d", u)
	}
}

func TestHitDiceSet(t *testing.T) {
	r, kn, _ := party(t)
	hd := func(id string, used any) m { return m{"characterId": id, "used": used} }
	failWith(t, r, "p2", "hitDice.set", hd(kn, 1), "you can only change your own character")
	failWith(t, r, "p1", "hitDice.set", hd(kn, 2), "used must be 0-1")
	failWith(t, r, "p1", "hitDice.set", hd(kn, -1), "used must be 0-1")
	ok(t, r, "p1", "hitDice.set", hd(kn, 1))
	if l := lastLog(r); l.Text.EN != "Aldric Vane: hit dice used 1/1" {
		t.Fatalf("log %+v", l)
	}
	ok(t, r, "dm", "character.level", m{"characterId": kn, "delta": 1})
	ok(t, r, "dm", "hitDice.set", hd(kn, 2))
	if c := r.State().Characters[0]; c.HitDiceUsed != 2 || c.HitDiceMax != 2 {
		t.Fatalf("hit dice %d/%d", c.HitDiceUsed, c.HitDiceMax)
	}
}

func TestPartyRest(t *testing.T) {
	r := creating(t)
	ok(t, r, "p1", "character.pick", pick("knight"))
	failWith(t, r, "dm", "party.rest", m{"kind": "long"}, "not allowed while the room is creating")

	r, ids := partyOf(t, "knight", "bard", "legend-warlock", "wizard")
	kn, bd, wl, wz := ids[0], ids[1], ids[2], ids[3]
	failWith(t, r, "p1", "party.rest", m{"kind": "short"}, "only the DM can do that")
	for _, k := range []any{"", "nap", nil} {
		failWith(t, r, "dm", "party.rest", m{"kind": k}, "kind must be short or long")
	}
	failWith(t, r, "dm", "party.rest", m{}, "kind must be short or long")

	ok(t, r, "dm", "character.level", m{"characterId": kn, "delta": 1}) // action surge 1 (short)
	use := func(client, id, res string, n int) {
		t.Helper()
		ok(t, r, client, "resource.set", m{"characterId": id, "resourceId": res, "used": n})
	}
	use("p1", kn, "second-wind", 2) // long, shortRegain 1
	use("p1", kn, "action-surge", 1)
	use("p2", bd, "bardic-inspiration", 2) // long, short from level 5
	ok(t, r, "p3", "slot.set", m{"characterId": wl, "level": "pact", "used": 1})
	ok(t, r, "p4", "slot.set", m{"characterId": wz, "level": 1, "used": 2})
	use("p4", wz, "arcane-recovery", 1) // long, no short regain
	ok(t, r, "p1", "hitDice.set", m{"characterId": kn, "used": 1})
	ok(t, r, "dm", "character.adjust", m{"characterId": kn, "hp": -5})

	ok(t, r, "dm", "party.rest", m{"kind": "short"})
	if f := lastFeed(r); f.EN != "DM: short rest — short-rest uses and pact slots back" || f.TH != "DM: พักสั้น — คืนความสามารถแบบพักสั้นและ pact slot" {
		t.Fatalf("feed %+v", f)
	}
	cs := r.State().Characters
	if u := cs[0].ResourcesUsed; u["second-wind"] != 1 || u["action-surge"] != 0 || cs[0].HitDiceUsed != 1 || cs[0].HP != 15 {
		t.Fatalf("fighter after short: %v hd %d hp %d", u, cs[0].HitDiceUsed, cs[0].HP)
	}
	if cs[1].ResourcesUsed["bardic-inspiration"] != 2 || cs[2].PactUsed != 0 || cs[3].SlotsUsed != [9]int{2} ||
		cs[3].ResourcesUsed["arcane-recovery"] != 1 {
		t.Fatalf("after short: bard %v pact %d wizard %v %v", cs[1].ResourcesUsed, cs[2].PactUsed, cs[3].SlotsUsed, cs[3].ResourcesUsed)
	}
	ok(t, r, "dm", "party.rest", m{"kind": "short"}) // second short: shortRegain again, never below 0
	ok(t, r, "dm", "party.rest", m{"kind": "short"})
	if u := r.State().Characters[0].ResourcesUsed["second-wind"]; u != 0 {
		t.Fatalf("second wind %d", u)
	}

	// Bard from level 5 (Font of Inspiration): a short rest refills it.
	for range 4 {
		ok(t, r, "dm", "character.level", m{"characterId": bd, "delta": 1})
	}
	use("p2", bd, "bardic-inspiration", 3)
	ok(t, r, "dm", "party.rest", m{"kind": "short"})
	if u := r.State().Characters[1].ResourcesUsed["bardic-inspiration"]; u != 0 {
		t.Fatalf("bard at 5 after short: %d", u)
	}

	// Long: everything full; a knocked-out character is back up.
	use("p2", bd, "bardic-inspiration", 1)
	ok(t, r, "dm", "character.adjust", m{"characterId": kn, "hp": -99})
	ok(t, r, "dm", "party.rest", m{"kind": "long"})
	cs = r.State().Characters
	for _, c := range cs {
		for id, u := range c.ResourcesUsed {
			if u != 0 {
				t.Fatalf("%s %s used %d after long rest", c.Name.EN, id, u)
			}
		}
		if c.HP != c.HPMax || c.SlotsUsed != [9]int{} || c.PactUsed != 0 || c.HitDiceUsed != 0 {
			t.Fatalf("%s after long: %+v", c.Name.EN, c)
		}
	}
	f := r.State().Feed
	if f[len(f)-2].Text.EN != "DM: long rest — HP, spell slots, uses and hit dice refilled" || f[len(f)-1].Text.EN != "Aldric Vane is back up" {
		t.Fatalf("feed %+v %+v", f[len(f)-2], f[len(f)-1])
	}
}

func TestCustomAbilities(t *testing.T) {
	r, kn, wz := party(t)
	add := func(name, kind string, level any, desc string) m {
		return m{"name": name, "kind": kind, "level": level, "desc": desc}
	}
	for _, tc := range []struct {
		name, client string
		p            m
		want         string
	}{
		{"dm has none", "dm", add("x", "spell", 1, ""), "you have no character"},
		{"stranger", "stranger", add("x", "spell", 1, ""), "you have no character"},
		{"empty name", "p1", add("  ", "spell", 1, ""), "ability name must be 1-40 characters"},
		{"long name", "p1", add(strings.Repeat("ก", 41), "spell", 1, ""), "ability name must be 1-40 characters"},
		{"name newline", "p1", add("a\nb", "spell", 1, ""), "ability name must be 1-40 characters"},
		{"long desc", "p1", add("x", "spell", 1, strings.Repeat("d", 301)), "ability desc must be at most 300 characters"},
		{"desc control", "p1", add("x", "spell", 1, "a\x07"), "ability desc must be at most 300 characters"},
		{"bad kind", "p1", add("x", "item", 1, ""), "kind must be spell or feature"},
		{"no kind", "p1", m{"name": "x"}, "kind must be spell or feature"},
		{"spell level 10", "p1", add("x", "spell", 10, ""), "spell level must be 0-9"},
		{"spell level -1", "p1", add("x", "spell", -1, ""), "spell level must be 0-9"},
		{"level float", "p1", add("x", "spell", 1.5, ""), "malformed payload"},
	} {
		t.Run(tc.name, func(t *testing.T) { failWith(t, r, tc.client, "ability.add", tc.p, tc.want) })
	}

	ok(t, r, "p1", "ability.add", add(" Fire Bolt ", "spell", 0, " a\nb "))
	ok(t, r, "p1", "ability.add", add("Tough", "feature", 7, "")) // a feature has no spell level
	got := r.State().Characters[0].AbilitiesCustom
	if len(got) != 2 || got[0] != (CustomAbility{ID: got[0].ID, Name: "Fire Bolt", Kind: "spell", Level: 0, Desc: "a\nb"}) ||
		got[1].Level != 0 || got[1].Kind != "feature" || !strings.HasPrefix(got[0].ID, "ab") || got[0].ID == got[1].ID {
		t.Fatalf("abilities %+v", got)
	}
	if l := lastLog(r); l.Text.EN != "added feature Tough to Aldric Vane" {
		t.Fatalf("log %+v", l)
	}
	id := got[0].ID

	upd := func(cid, aid, name string) m {
		return m{"characterId": cid, "abilityId": aid, "name": name, "kind": "spell", "level": 3, "desc": "d"}
	}
	failWith(t, r, "p2", "ability.update", upd(kn, id, "Hack"), "no such ability") // only the sender's own list
	failWith(t, r, "dm", "ability.update", upd(kn, id, "Hack"), "you have no character")
	failWith(t, r, "p1", "ability.update", upd(kn, "nope", "Hack"), "no such ability")
	failWith(t, r, "p1", "ability.update", upd(kn, id, ""), "ability name must be 1-40 characters")
	ok(t, r, "p1", "ability.update", upd(kn, id, "Fireball"))
	if a := r.State().Characters[0].AbilitiesCustom[0]; a != (CustomAbility{ID: id, Name: "Fireball", Kind: "spell", Level: 3, Desc: "d"}) {
		t.Fatalf("updated %+v", a)
	}
	if l := lastLog(r); l.Text.EN != "edited Fireball of Aldric Vane" {
		t.Fatalf("log %+v", l)
	}

	rm := func(cid, aid string) m { return m{"characterId": cid, "abilityId": aid} }
	failWith(t, r, "p2", "ability.remove", rm(kn, id), "you can only change your own character")
	failWith(t, r, "p2", "ability.remove", rm(wz, id), "no such ability")
	failWith(t, r, "dm", "ability.remove", rm("nope", id), "no such character")
	ok(t, r, "co", "ability.remove", rm(kn, id)) // the DM may remove
	if l := lastLog(r); l.Text.EN != "removed Fireball from Aldric Vane" {
		t.Fatalf("log %+v", l)
	}
	ok(t, r, "p1", "ability.remove", rm(kn, r.State().Characters[0].AbilitiesCustom[0].ID))
	if n := len(r.State().Characters[0].AbilitiesCustom); n != 0 {
		t.Fatalf("left %d", n)
	}

	for i := range maxCustomAbilities {
		ok(t, r, "p2", "ability.add", add("Spell "+string(rune('A'+i)), "spell", i%10, ""))
	}
	failWith(t, r, "p2", "ability.add", add("One more", "spell", 1, ""), "at most 40 custom abilities")
	ids := []string{}
	for _, a := range r.State().Characters[1].AbilitiesCustom {
		ids = append(ids, a.ID)
	}
	slices.Sort(ids)
	if len(slices.Compact(ids)) != maxCustomAbilities {
		t.Fatal("duplicate ability ids")
	}

	r, _ = lobby(t)
	failWith(t, r, "p1", "ability.add", add("x", "spell", 1, ""), "not allowed while the room is lobby")
}
