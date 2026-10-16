package room

import (
	"encoding/json"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"projectdnd/internal/rules"
)

// failWith is fail plus an exact error message.
func failWith(t *testing.T, r *Room, client, typ string, payload any, want string) {
	t.Helper()
	before, _ := json.Marshal(r.State())
	err := apply(r, client, typ, payload)
	if err == nil || err.Error() != want {
		t.Fatalf("%s by %s: err %v, want %q", typ, client, err, want)
	}
	if after, _ := json.Marshal(r.State()); string(before) != string(after) {
		t.Fatalf("%s by %s: state changed on error", typ, client)
	}
}

// creating: dm, co-DM "co", players p1 (Pat) and p2 (Sam); status creating, nobody has picked.
func creating(t *testing.T) *Room {
	t.Helper()
	r, _ := lobby(t)
	ok(t, r, "dm", "seat.add", m{})
	ok(t, r, "p2", "seat.sit", m{"seatId": r.State().Seats[2].ID, "name": "Sam"})
	ok(t, r, "dm", "seat.openCoDm", m{})
	ok(t, r, "co", "seat.sit", m{"seatId": r.State().Seats[3].ID, "name": "Cora"})
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "p2", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "game.start", m{})
	return r
}

// party: creating + p1 knight (HP 36, MP 4) + p2 wizard (HP 16, MP 20); status playing.
func party(t *testing.T) (r *Room, knight, wizard string) {
	t.Helper()
	r = creating(t)
	ok(t, r, "p1", "character.pick", pick("knight"))
	ok(t, r, "p2", "character.pick", pick("wizard"))
	cs := r.State().Characters
	return r, cs[0].ID, cs[1].ID
}

func lastFeed(r *Room) rules.L {
	f := r.State().Feed
	return f[len(f)-1].Text
}

func TestCharacterPick(t *testing.T) {
	r, _ := lobby(t)
	failWith(t, r, "p1", "character.pick", pick("knight"), "not allowed while the room is lobby")

	r = creating(t)
	failWith(t, r, "dm", "character.pick", pick("knight"), "only a seated player can pick a character")
	failWith(t, r, "co", "character.pick", pick("knight"), "only a seated player can pick a character")
	failWith(t, r, "stranger", "character.pick", pick("knight"), "only a seated player can pick a character")
	failWith(t, r, "p1", "character.pick", pick("pirate"), "no such class")
	failWith(t, r, "p1", "character.pick", pick(""), "no such class")
	failWith(t, r, "p1", "character.pick", m{"classId": 3}, "malformed payload")

	// A roster edit before the pick is what the player gets.
	e := rules.DefaultRoster()[0]
	e.HeroName.EN = "Sir Edit"
	ok(t, r, "dm", "roster.update", m{"classId": "knight", "entry": e})
	ok(t, r, "p1", "character.pick", pick("knight"))
	failWith(t, r, "p1", "character.pick", pick("wizard"), "you already have a character")
	failWith(t, r, "p2", "character.pick", pick("knight"), "that class is already taken")
	if r.State().Status != StatusCreating {
		t.Fatal("playing before p2 picked")
	}
	c := r.State().Characters[0]
	if c.ClassID != "knight" || c.Name.EN != "Sir Edit" || c.Name.TH != "อัลดริก เวน" || c.Art != "" ||
		c.HP != 36 || c.HPMax != 36 || c.MP != 4 || c.MPMax != 4 || c.Str != 5 || c.End != 8 || c.Agi != 2 ||
		len(c.Skills) != 3 || c.Skills[0].Kind != "innate" || c.Skills[0].ID != "stand-guard" || len(c.Items) != 4 {
		t.Fatalf("character %+v", c)
	}
	if f := lastFeed(r); f.EN != "Pat picked Sir Edit (Knight)" || f.TH != "Pat เลือก อัลดริก เวน (อัศวิน)" {
		t.Fatalf("feed %+v", f)
	}

	// A roster edit after the pick leaves the picked character alone.
	e.Stats.HPMax, e.Items[0].Qty = 99, 9
	ok(t, r, "dm", "roster.update", m{"classId": "knight", "entry": e})
	ok(t, r, "dm", "item.remove", m{"characterId": c.ID, "index": 3, "qty": 1})
	if c := r.State().Characters[0]; c.HPMax != 36 || c.Items[0].Qty != 1 || c.Items[3].Qty != 2 {
		t.Fatalf("picked character changed by roster edit: %+v", c)
	}
	if r.State().Roster[0].Items[3].Qty != 3 {
		t.Fatal("item.remove changed the roster")
	}

	ok(t, r, "p2", "character.pick", pick("wizard"))
	if r.State().Status != StatusPlaying {
		t.Fatal("not playing after every player picked")
	}
	failWith(t, r, "p2", "character.pick", pick("rogue"), "not allowed while the room is playing")
}

func TestCharacterTokenUsesArt(t *testing.T) {
	r := creating(t)
	img := withUpload(t, r)
	e := rules.DefaultRoster()[0]
	e.Art = img
	ok(t, r, "dm", "roster.update", m{"classId": "knight", "entry": e})
	ok(t, r, "p1", "character.pick", pick("knight"))
	ok(t, r, "p2", "character.pick", pick("bard"))
	ok(t, r, "dm", "scene.add", m{"name": "Beach", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 10, "rows": 5}})
	ok(t, r, "dm", "scene.select", m{"sceneId": r.State().Scenes[0].ID})
	st := r.State()
	if st.Characters[0].Art != img || st.Tokens[0].ImageURL != img || st.Tokens[0].Name != "Aldric Vane" ||
		st.Tokens[1].ImageURL != "" || st.Tokens[1].Name != "Cosimo Lark" {
		t.Fatalf("tokens %+v", st.Tokens)
	}
}

func TestCustomSkills(t *testing.T) {
	r := creating(t)
	ok(t, r, "p1", "character.pick", pick("knight"))
	kn := r.State().Characters[0].ID
	add := func(name, desc string, cost any) m { return m{"name": name, "desc": desc, "cost": cost} }

	for _, tc := range []struct {
		name, client string
		payload      m
		want         string
	}{
		{"no character", "p2", add("x", "", 1), "you have no character"},
		{"dm has none", "dm", add("x", "", 1), "you have no character"},
		{"stranger", "stranger", add("x", "", 1), "you have no character"},
		{"empty name", "p1", add("   ", "", 1), "skill name must be 1-40 characters"},
		{"long name", "p1", add(strings.Repeat("ก", 41), "", 1), "skill name must be 1-40 characters"},
		{"name newline", "p1", add("a\nb", "", 1), "skill name must be 1-40 characters"},
		{"long desc", "p1", add("x", strings.Repeat("d", 301), 1), "skill desc must be at most 300 characters"},
		{"desc control", "p1", add("x", "a\x07", 1), "skill desc must be at most 300 characters"},
		{"cost 11", "p1", add("x", "", 11), "skill cost must be 0-10"},
		{"cost -1", "p1", add("x", "", -1), "skill cost must be 0-10"},
		{"cost float", "p1", add("x", "", 1.5), "malformed payload"},
	} {
		t.Run(tc.name, func(t *testing.T) { failWith(t, r, tc.client, "skill.addCustom", tc.payload, tc.want) })
	}

	ok(t, r, "p1", "skill.addCustom", add(" Kick ", " a\nb ", 0))
	ok(t, r, "p1", "skill.addCustom", add(strings.Repeat("ก", 40), strings.Repeat("d", 300), 10))
	ok(t, r, "p1", "skill.addCustom", add("Third", "", 3))
	failWith(t, r, "p1", "skill.addCustom", add("Fourth", "", 1), "at most 3 custom skills")
	sk := slices.Clone(r.State().Characters[0].Skills) // State shares the slice that removeCustom shifts
	k := sk[3]
	if len(sk) != 6 || k.Kind != "custom" || k.Name != (rules.L{TH: "Kick", EN: "Kick"}) || k.Desc.EN != "a\nb" || k.Cost != 0 ||
		!strings.HasPrefix(k.ID, "cs_") {
		t.Fatalf("skills %+v", sk)
	}

	edit := func(char, skill, name string, cost int) m {
		return m{"characterId": char, "skillId": skill, "name": name, "desc": "", "cost": cost}
	}
	const notOwn = "you can only change your own character's skills"
	failWith(t, r, "p2", "skill.editCustom", edit(kn, k.ID, "Hack", 1), notOwn)
	failWith(t, r, "stranger", "skill.editCustom", edit(kn, k.ID, "Hack", 1), notOwn)
	failWith(t, r, "p2", "skill.removeCustom", m{"characterId": kn, "skillId": k.ID}, notOwn)
	failWith(t, r, "p1", "skill.editCustom", edit(kn, "stand-guard", "Hack", 1), "only custom skills can be changed")
	failWith(t, r, "p1", "skill.removeCustom", m{"characterId": kn, "skillId": "stand-guard"}, "only custom skills can be changed")
	failWith(t, r, "p1", "skill.editCustom", edit(kn, "nope", "Hack", 1), "no such skill")
	failWith(t, r, "dm", "skill.editCustom", edit("nope", k.ID, "Hack", 1), "no such character")
	failWith(t, r, "p1", "skill.editCustom", edit(kn, k.ID, "", 1), "skill name must be 1-40 characters")
	failWith(t, r, "co", "skill.editCustom", edit(kn, k.ID, "Hack", 11), "skill cost must be 0-10")

	ok(t, r, "p1", "skill.editCustom", edit(kn, k.ID, "Kick 2", 2))
	ok(t, r, "co", "skill.editCustom", edit(kn, k.ID, "Kick 3", 5)) // DM caps an overpowered skill
	if got := r.State().Characters[0].Skills[3]; got.ID != k.ID || got.Name.EN != "Kick 3" || got.Cost != 5 {
		t.Fatalf("edited %+v", got)
	}
	ok(t, r, "dm", "skill.removeCustom", m{"characterId": kn, "skillId": k.ID})
	ok(t, r, "p1", "skill.removeCustom", m{"characterId": kn, "skillId": sk[4].ID})
	failWith(t, r, "p1", "skill.removeCustom", m{"characterId": kn, "skillId": k.ID}, "no such skill")
	ok(t, r, "p2", "character.pick", pick("wizard")) // playing: custom skills still allowed
	ok(t, r, "p1", "skill.addCustom", add("Again", "", 1))
	ok(t, r, "p1", "skill.addCustom", add("Again 2", "", 1))
	if sk := r.State().Characters[0].Skills; len(sk) != 6 || sk[3].Name.EN != "Third" {
		t.Fatalf("after remove/add %+v", sk)
	}
}

func TestSkillUse(t *testing.T) {
	r := creating(t)
	ok(t, r, "p1", "character.pick", pick("knight"))
	failWith(t, r, "p1", "skill.use", m{"skillId": "stand-guard"}, "not allowed while the room is creating")
	ok(t, r, "p1", "skill.addCustom", m{"name": "Free", "desc": "", "cost": 0})
	free := r.State().Characters[0].Skills[3].ID
	ok(t, r, "p2", "character.pick", pick("wizard"))

	failWith(t, r, "dm", "skill.use", m{"skillId": "stand-guard"}, "you have no character")
	failWith(t, r, "stranger", "skill.use", m{"skillId": "stand-guard"}, "you have no character")
	failWith(t, r, "p1", "skill.use", m{"skillId": "nope"}, "no such skill")
	failWith(t, r, "p2", "skill.use", m{"skillId": "stand-guard"}, "no such skill") // p1's skill, not p2's
	ok(t, r, "p1", "skill.use", m{"skillId": "last-stand"})                         // cost 4 of 4
	if c := r.State().Characters[0]; c.MP != 0 {
		t.Fatalf("mp %d", c.MP)
	}
	if f := lastFeed(r); f.EN != "Aldric Vane used Last Stand (−4 MP)" || f.TH != "อัลดริก เวน ใช้ ยืนหยัดครั้งสุดท้าย (−4 มานา)" {
		t.Fatalf("feed %+v", f)
	}
	failWith(t, r, "p1", "skill.use", m{"skillId": "stand-guard"}, "not enough mana")
	ok(t, r, "p1", "skill.use", m{"skillId": free}) // custom, cost 0
	if f := lastFeed(r); f.EN != "Aldric Vane used Free (0 MP)" || !strings.HasSuffix(f.TH, "(0 มานา)") {
		t.Fatalf("free feed %+v", f)
	}
	ok(t, r, "p2", "skill.use", m{"skillId": "starfall"})
	if c := r.State().Characters[1]; c.MP != 15 {
		t.Fatalf("wizard mp %d", c.MP)
	}
}

func TestCharacterAdjust(t *testing.T) {
	ip := func(n int) *int { return &n }
	for _, tc := range []struct {
		name   string
		client string
		hp, mp *int
		wantHP int // knight starts HP 36/36, MP 4/4
		wantMP int
		err    string
		feedEN string
	}{
		{"hp down", "dm", ip(-5), nil, 31, 4, "", "DM: Aldric Vane HP −5 (31/36)"},
		{"hp clamps at max: no change, no feed", "dm", ip(100), nil, 36, 4, "", ""},
		{"unchanged field omitted", "dm", ip(5), ip(-1), 36, 3, "", "DM: Aldric Vane MP −1 (3/4)"},
		{"zero deltas: no feed", "dm", ip(0), ip(0), 36, 4, "", ""},
		{"hp clamps at 0", "co", ip(-999), nil, 0, 4, "", "DM: Aldric Vane HP −36 (0/36)"},
		{"mp clamps at 0", "dm", nil, ip(-10), 36, 0, "", "DM: Aldric Vane MP −4 (0/4)"},
		{"both", "dm", ip(-1), ip(-2), 35, 2, "", "DM: Aldric Vane HP −1 (35/36), MP −2 (2/4)"},
		{"player", "p1", ip(-1), nil, 0, 0, "only the DM can do that", ""},
		{"nothing", "dm", nil, nil, 0, 0, "hp or mp delta is required", ""},
		{"hp too big", "dm", ip(1000), nil, 0, 0, "delta must be -999 to 999", ""},
		{"mp too small", "dm", nil, ip(-1000), 0, 0, "delta must be -999 to 999", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, kn, _ := party(t)
			p := m{"characterId": kn}
			if tc.hp != nil {
				p["hp"] = *tc.hp
			}
			if tc.mp != nil {
				p["mp"] = *tc.mp
			}
			if tc.err != "" {
				failWith(t, r, tc.client, "character.adjust", p, tc.err)
				return
			}
			before := len(r.State().Feed)
			ok(t, r, tc.client, "character.adjust", p)
			if c := r.State().Characters[0]; c.HP != tc.wantHP || c.MP != tc.wantMP {
				t.Fatalf("hp %d mp %d", c.HP, c.MP)
			}
			if tc.feedEN == "" {
				if n := len(r.State().Feed); n != before {
					t.Fatalf("feed grew %d -> %d: %+v", before, n, lastFeed(r))
				}
				return
			}
			f := lastFeed(r)
			if tc.wantHP == 0 { // knocked out: the skull line follows the DM line
				if f.EN != "💀 Aldric Vane is knocked out (HP 0)" {
					t.Fatalf("death line %+v", f)
				}
				f = r.State().Feed[len(r.State().Feed)-2].Text
			}
			if f.EN != tc.feedEN || !strings.HasPrefix(f.TH, "DM: อัลดริก เวน ") {
				t.Fatalf("feed %+v", f)
			}
		})
	}
	r, _ := lobby(t)
	failWith(t, r, "dm", "character.adjust", m{"characterId": "c1", "hp": 1}, "not allowed while the room is lobby")
	r = creating(t)
	failWith(t, r, "dm", "character.adjust", m{"characterId": "nope", "hp": 1}, "no such character")
}

func TestCharacterSetStats(t *testing.T) {
	stats := func(id string, hpMax, mpMax, str, end, agi int) m {
		return m{"characterId": id, "hpMax": hpMax, "mpMax": mpMax, "str": str, "end": end, "agi": agi}
	}
	r, kn, _ := party(t)
	failWith(t, r, "p1", "character.setStats", stats(kn, 40, 4, 5, 8, 2), "only the DM can do that")
	failWith(t, r, "dm", "character.setStats", stats("nope", 40, 4, 5, 8, 2), "no such character")
	failWith(t, r, "dm", "character.setStats", stats(kn, 0, 4, 5, 8, 2), "hpMax must be 1-999")
	failWith(t, r, "dm", "character.setStats", stats(kn, 1000, 4, 5, 8, 2), "hpMax must be 1-999")
	failWith(t, r, "dm", "character.setStats", stats(kn, 40, 100, 5, 8, 2), "mpMax must be 0-99")
	failWith(t, r, "dm", "character.setStats", stats(kn, 40, -1, 5, 8, 2), "mpMax must be 0-99")
	failWith(t, r, "dm", "character.setStats", stats(kn, 40, 4, 21, 8, 2), "str, end and agi must be 0-20")
	failWith(t, r, "dm", "character.setStats", stats(kn, 40, 4, 5, 8, -1), "str, end and agi must be 0-20")

	ok(t, r, "co", "character.setStats", stats(kn, 10, 2, 20, 0, 1)) // lowering max clamps hp/mp
	c := r.State().Characters[0]
	if c.HP != 10 || c.HPMax != 10 || c.MP != 2 || c.MPMax != 2 || c.Str != 20 || c.End != 0 || c.Agi != 1 {
		t.Fatalf("stats %+v", c)
	}
	if f := lastFeed(r); f.EN != "DM: set Aldric Vane: HP max 10, MP max 2, STR 20, END 0, AGI 1" {
		t.Fatalf("feed %+v", f)
	}
	ok(t, r, "dm", "character.setStats", stats(kn, 50, 9, 5, 8, 2)) // raising max keeps current hp/mp
	if c := r.State().Characters[0]; c.HP != 10 || c.MP != 2 || c.HPMax != 50 {
		t.Fatalf("raise %+v", c)
	}
}

func TestPartyRest(t *testing.T) {
	r := creating(t)
	ok(t, r, "p1", "character.pick", pick("knight"))
	failWith(t, r, "dm", "party.rest", m{}, "not allowed while the room is creating")
	ok(t, r, "p2", "character.pick", pick("wizard"))
	ok(t, r, "p1", "skill.use", m{"skillId": "last-stand"})
	ok(t, r, "p2", "skill.use", m{"skillId": "starfall"})
	ok(t, r, "dm", "character.adjust", m{"characterId": r.State().Characters[0].ID, "hp": -10})
	failWith(t, r, "p1", "party.rest", m{}, "only the DM can do that")
	ok(t, r, "co", "party.rest", m{})
	cs := r.State().Characters
	if cs[0].MP != 4 || cs[1].MP != 20 || cs[0].HP != 26 {
		t.Fatalf("after rest %+v %+v", cs[0], cs[1])
	}
	if f := lastFeed(r); f != (rules.L{TH: "DM: พักผ่อน — มานาเต็ม", EN: "DM: rest — mana refilled"}) {
		t.Fatalf("feed %+v", f)
	}
}

func TestItemGive(t *testing.T) {
	give := func(id, name string, qty any) m { return m{"characterId": id, "name": name, "qty": qty} }
	r, kn, wz := party(t)
	failWith(t, r, "p1", "item.give", give(kn, "Rope", 1), "only the DM can do that")
	failWith(t, r, "dm", "item.give", give("nope", "Rope", 1), "no such character")
	failWith(t, r, "dm", "item.give", give(kn, " ", 1), "item name must be 1-40 characters")
	failWith(t, r, "dm", "item.give", give(kn, strings.Repeat("x", 41), 1), "item name must be 1-40 characters")
	failWith(t, r, "dm", "item.give", give(kn, "Rope", 0), "item qty must be 1-99")
	failWith(t, r, "dm", "item.give", give(kn, "Rope", 100), "item qty must be 1-99")

	ok(t, r, "dm", "item.give", give(kn, " Rope ", 2))
	ok(t, r, "co", "item.give", give(kn, "Rope", 3))
	ok(t, r, "dm", "item.give", give(kn, "Field Bandage Roll", 1)) // roster item, matched by en name
	ok(t, r, "dm", "item.give", give(kn, "ม้วนผ้าพันแผลสนาม", 1))  // and by th name
	items := r.State().Characters[0].Items
	if len(items) != 5 || items[4] != (rules.RosterItem{Name: rules.L{TH: "Rope", EN: "Rope"}, Qty: 5}) || items[3].Qty != 5 {
		t.Fatalf("items %+v", items)
	}
	if f := lastFeed(r); f.EN != "DM: gave Aldric Vane Field Bandage Roll ×1" || f.TH != "DM: ให้ ม้วนผ้าพันแผลสนาม ×1 แก่ อัลดริก เวน" {
		t.Fatalf("feed %+v", f)
	}
	for range 10 {
		ok(t, r, "dm", "item.give", give(kn, "Rope", 99)) // 995
	}
	failWith(t, r, "dm", "item.give", give(kn, "Rope", 5), "an item stack holds at most 999")
	for i := len(r.State().Characters[1].Items); i < maxItemStacks; i++ {
		ok(t, r, "dm", "item.give", give(wz, "thing "+string(rune('A'+i)), 1))
	}
	failWith(t, r, "dm", "item.give", give(wz, "one too many", 1), "at most 50 different items")
	ok(t, r, "dm", "item.give", give(wz, "thing E", 1)) // merging is still fine when full
}

func TestItemRemove(t *testing.T) {
	rm := func(id string, index, qty int) m { return m{"characterId": id, "index": index, "qty": qty} }
	r, kn, _ := party(t) // knight: shield 1, sword 1, whetstone 1, bandage 3
	failWith(t, r, "p1", "item.remove", rm(kn, 3, 1), "only the DM can do that")
	failWith(t, r, "dm", "item.remove", rm("nope", 3, 1), "no such character")
	failWith(t, r, "dm", "item.remove", rm(kn, 4, 1), "no such item")
	failWith(t, r, "dm", "item.remove", rm(kn, -1, 1), "no such item")
	failWith(t, r, "dm", "item.remove", rm(kn, 3, 0), "qty must be at least 1")

	ok(t, r, "dm", "item.remove", rm(kn, 3, 1))
	if it := r.State().Characters[0].Items; len(it) != 4 || it[3].Qty != 2 {
		t.Fatalf("−1: %+v", it)
	}
	if f := lastFeed(r); f.EN != "DM: removed Field Bandage Roll ×1 from Aldric Vane" || f.TH != "DM: เอา ม้วนผ้าพันแผลสนาม ×1 ออกจาก อัลดริก เวน" {
		t.Fatalf("feed %+v", f)
	}
	ok(t, r, "co", "item.remove", rm(kn, 3, 50)) // more than the stack: whole stack, feed says 2
	if f := lastFeed(r); f.EN != "DM: removed Field Bandage Roll ×2 from Aldric Vane" {
		t.Fatalf("feed %+v", f)
	}
	ok(t, r, "dm", "item.remove", rm(kn, 0, 1)) // exactly the stack
	if it := r.State().Characters[0].Items; len(it) != 2 || it[0].Name.EN != "Knight's Longsword" {
		t.Fatalf("items %+v", it)
	}
}

func TestFeedIsCappedAndHasNoClientIDs(t *testing.T) {
	r, kn, _ := party(t)
	for i := range 60 {
		ok(t, r, "dm", "character.adjust", m{"characterId": kn, "hp": 2*(i%2) - 1}) // −1, +1, ...: every call changes HP
	}
	f := r.State().Feed
	if len(f) != maxFeed || f[0].At == 0 {
		t.Fatalf("feed len %d", len(f))
	}
	ok(t, r, "dm", "party.rest", m{})
	if lastFeed(r).EN != "DM: rest — mana refilled" || len(r.State().Feed) != maxFeed {
		t.Fatal("newest entry not last")
	}
	b, _ := json.Marshal(r.State().Feed)
	for _, secret := range []string{"p1", "p2", "stranger"} {
		if strings.Contains(string(b), secret) {
			t.Fatalf("feed contains client id %q", secret)
		}
	}
	if b, _ := json.Marshal(New().State()); !strings.Contains(string(b), `"feed":[]`) {
		t.Fatalf("empty feed not []: %s", b)
	}
}

func TestRosterEdit(t *testing.T) {
	r, _ := lobby(t)
	e := rules.DefaultRoster()[0]
	e.ID = "wizard" // ignored: classId is the key
	e.HeroName.EN = " Sir Test "
	upd := func(class string, e rules.RosterEntry) m { return m{"classId": class, "entry": e} }

	failWith(t, r, "p1", "roster.update", upd("knight", e), "only the DM can do that")
	failWith(t, r, "stranger", "roster.reset", m{}, "only the DM can do that")
	failWith(t, r, "dm", "roster.update", upd("pirate", e), "no such class")
	bad := e
	bad.Stats.HPMax = 0
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "hpMax must be 1-999")
	bad = e
	bad.Skills = nil
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "a class needs 1-5 skills")
	bad = e
	bad.Art = "/uploads/" + strings.Repeat("a", 32) + ".png"
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "art must be empty or an uploaded image") // no file
	img := withUpload(t, r)
	bad.Art = "/uploads/" + strings.Repeat("d", 32) + ".png" // a directory
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "art must be empty or an uploaded image")

	ok(t, r, "dm", "roster.update", upd("knight", e)) // lobby
	ro := r.State().Roster
	if len(ro) != 10 || ro[0].ID != "knight" || ro[0].HeroName.EN != "Sir Test" || ro[1].ID != "wizard" || ro[1].HeroName.EN != "Soraya Quill" {
		t.Fatalf("roster %+v %+v", ro[0].HeroName, ro[1].HeroName)
	}
	e.Art = img
	ok(t, r, "dm", "roster.update", upd("rogue", e))
	if ro := r.State().Roster; ro[2].ID != "rogue" || ro[2].Art != img {
		t.Fatalf("rogue %+v", ro[2])
	}
	if ps := r.PlayerState().Roster; len(ps) != 10 || ps[2].Art != img {
		t.Fatal("players do not get the roster in use")
	}
	failWith(t, r, "dm", "roster.reset", m{"classId": "pirate"}, "no such class")
	ok(t, r, "dm", "roster.reset", m{"classId": "knight"})
	if ro := r.State().Roster; ro[0].HeroName.EN != "Aldric Vane" || ro[2].Art != img {
		t.Fatal("reset one")
	}
	ok(t, r, "dm", "roster.reset", m{})
	if ro := r.State().Roster; ro[2].Art != "" {
		t.Fatal("reset all")
	}

	// co-DM, while playing, persisted through the store the Hub was given.
	p, _, _ := party(t)
	dir := t.TempDir()
	store, err := rules.LoadRoster(dir, "", slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	p.roster = store
	e = rules.DefaultRoster()[9]
	e.Blurb.EN = "Edited"
	ok(t, p, "co", "roster.update", upd("necromancer", e))
	again, _ := rules.LoadRoster(dir, "", slog.New(slog.DiscardHandler))
	if got, _ := again.Get("necromancer"); got.Blurb.EN != "Edited" {
		t.Fatal("roster edit not persisted")
	}
}
