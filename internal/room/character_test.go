package room

import (
	"encoding/json"
	"log/slog"
	"maps"
	"slices"
	"strconv"
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

// party: creating + p1 knight (fighter, HP 12) + p2 wizard (HP 7, two 1st-level slots); status playing.
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

func className(id string) rules.L {
	c, _ := rules.ClassByID(id)
	return c.Name
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
	failWith(t, r, "p1", "character.pick", m{"heroId": 3}, "malformed payload")

	// A roster edit before the pick is what the player gets.
	e := rules.DefaultRoster()[0]
	e.HeroName.EN = "Sir Edit"
	ok(t, r, "dm", "roster.update", m{"classId": "knight", "entry": e})
	ok(t, r, "p1", "character.pick", m{"classId": "knight"}) // pre-v2 key still works
	failWith(t, r, "p1", "character.pick", pick("wizard"), "you already have a character")
	failWith(t, r, "p2", "character.pick", pick("knight"), "that class is already taken")
	if r.State().Status != StatusCreating {
		t.Fatal("playing before p2 picked")
	}
	c := r.State().Characters[0]
	want := rules.Abilities{Str: 17, Dex: 10, Con: 15, Int: 8, Wis: 13, Cha: 12}
	if c.ClassID != "fighter" || c.HeroID != "knight" || c.Name.EN != "Sir Edit" || c.Name.TH != "อัลดริก เวน" || c.Art != "" ||
		c.Level != 1 || c.HP != 12 || c.HPMax != 12 || c.Abilities != want || c.AC != 18 || c.PB != 2 || c.HitDiceMax != 1 ||
		!slices.Equal(c.Skills, []string{"athletics", "perception"}) || len(c.Items) != 6 || c.SlotsMax != [9]int{} ||
		c.Pact != (rules.PactSlots{}) || !maps.Equal(c.ResourcesMax, map[string]int{"second-wind": 2, "action-surge": 0, "indomitable": 0}) ||
		!maps.Equal(c.ResourcesUsed, map[string]int{"second-wind": 0, "action-surge": 0, "indomitable": 0}) {
		t.Fatalf("character %+v", c)
	}
	b, _ := json.Marshal(c)
	for _, k := range []string{`"expertise":[]`, `"abilitiesCustom":[]`, `"slotsUsed":[0,0,0,0,0,0,0,0,0]`, `"pactUsed":0`,
		`"pact":{"count":0,"level":0}`, `"hitDiceUsed":0`, `"pb":2`} {
		if !strings.Contains(string(b), k) {
			t.Fatalf("character JSON has no %s: %s", k, b)
		}
	}
	fighter := className("fighter")
	if f := lastFeed(r); f.EN != "Pat picked Sir Edit ("+fighter.EN+")" || f.TH != "Pat เลือก อัลดริก เวน ("+fighter.TH+")" {
		t.Fatalf("feed %+v", f)
	}

	// A roster edit after the pick leaves the picked character alone.
	e.Abilities.Con, e.Items[0].Qty = 30, 9
	ok(t, r, "dm", "roster.update", m{"classId": "knight", "entry": e})
	ok(t, r, "dm", "item.remove", m{"characterId": c.ID, "index": 4, "qty": 1})
	if c := r.State().Characters[0]; c.HPMax != 12 || c.Abilities.Con != 15 || c.Items[0].Qty != 1 || c.Items[4].Qty != 2 {
		t.Fatalf("picked character changed by roster edit: %+v", c)
	}
	if r.State().Roster[0].Items[4].Qty != 3 {
		t.Fatal("item.remove changed the roster")
	}

	ok(t, r, "p2", "character.pick", m{"heroId": "wizard"})
	if r.State().Status != StatusPlaying {
		t.Fatal("not playing after every player picked")
	}
	failWith(t, r, "p2", "character.pick", pick("rogue"), "not allowed while the room is playing")
}

// partyOf: both sets active, players p1..pn (n ≤ 6) pick heroes[i] in order; status playing.
func partyOf(t *testing.T, heroes ...string) (*Room, []string) {
	t.Helper()
	r, _ := lobby(t)
	ok(t, r, "dm", "sets.set", m{"sets": []string{"core", "legend"}})
	for i := 2; i <= len(heroes); i++ {
		ok(t, r, "dm", "seat.add", m{})
		seats := r.State().Seats
		ok(t, r, "p"+strconv.Itoa(i), "seat.sit", m{"seatId": seats[len(seats)-1].ID, "name": "P" + strconv.Itoa(i)})
	}
	for i := range heroes {
		ok(t, r, "p"+strconv.Itoa(i+1), "seat.ready", m{"ready": true})
	}
	ok(t, r, "dm", "game.start", m{})
	ids := make([]string, len(heroes))
	for i, h := range heroes {
		ok(t, r, "p"+strconv.Itoa(i+1), "character.pick", pick(h))
		ids[i] = r.State().Characters[i].ID
	}
	return r, ids
}

// Level 1 values by hand: HP = hit die + CON mod; resources from the class tables.
func TestPickLevelOneValues(t *testing.T) {
	for _, tc := range []struct {
		hero, class string
		hp          int
		slots       [9]int
		pact        rules.PactSlots
		res         map[string]int
	}{
		{"knight", "fighter", 12, [9]int{}, rules.PactSlots{}, map[string]int{"second-wind": 2, "action-surge": 0, "indomitable": 0}},
		{"wizard", "wizard", 7, [9]int{2}, rules.PactSlots{}, map[string]int{"arcane-recovery": 1, "signature-spells": 0}},        // d6 + 1
		{"bard", "bard", 9, [9]int{2}, rules.PactSlots{}, map[string]int{"bardic-inspiration": 3}},                                // CHA 17 → 3
		{"barbarian", "barbarian", 14, [9]int{}, rules.PactSlots{}, map[string]int{"rage": 2, "persistent-rage": 0}},              // d12 + 2
		{"legend-warlock", "warlock", 10, [9]int{}, rules.PactSlots{Count: 1, Level: 1}, nil},                                     // d8 + 2
		{"legend-paladin", "paladin", 13, [9]int{2}, rules.PactSlots{}, map[string]int{"lay-on-hands": 5, "channel-divinity": 0}}, // d10 + 3, level*5
		{"legend-ranger", "ranger", 12, [9]int{2}, rules.PactSlots{}, map[string]int{"favored-enemy": 2, "tireless": 0, "natures-veil": 0}},
		{"legend-monk", "monk", 9, [9]int{}, rules.PactSlots{}, map[string]int{"focus-points": 0, "uncanny-metabolism": 0}},
	} {
		t.Run(tc.hero, func(t *testing.T) {
			r, _ := partyOf(t, tc.hero)
			c := r.State().Characters[0]
			if c.ClassID != tc.class || c.HP != tc.hp || c.HPMax != tc.hp || c.SlotsMax != tc.slots || c.Pact != tc.pact ||
				c.PB != 2 || c.Level != 1 || c.HitDiceMax != 1 {
				t.Fatalf("got class %s hp %d/%d slots %v pact %v pb %d", c.ClassID, c.HP, c.HPMax, c.SlotsMax, c.Pact, c.PB)
			}
			for id, v := range tc.res {
				if got, ok := c.ResourcesMax[id]; !ok || got != v {
					t.Fatalf("resource %s max %d (%v), want %d", id, got, ok, v)
				}
			}
		})
	}
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

func TestRemovedV1Events(t *testing.T) {
	r, kn, _ := party(t)
	for _, typ := range []string{"skill.addCustom", "skill.editCustom", "skill.removeCustom", "skill.use", "character.setStats"} {
		failWith(t, r, "dm", typ, m{"characterId": kn}, `unknown event type "`+typ+`"`)
	}
}

func TestCharacterAdjust(t *testing.T) {
	for _, tc := range []struct {
		name   string
		client string
		hp     any
		wantHP int // knight starts HP 12/12
		err    string
		feedEN string
	}{
		{"hp down", "dm", -5, 7, "", "DM: Aldric Vane HP −5 (7/12)"},
		{"hp clamps at max: no change, no feed", "dm", 100, 12, "", ""},
		{"zero: no feed", "dm", 0, 12, "", ""},
		{"hp clamps at 0", "co", -999, 0, "", "DM: Aldric Vane HP −12 (0/12)"},
		{"player", "p1", -1, 0, "only the DM can do that", ""},
		{"nothing", "dm", nil, 0, "hp delta is required", ""},
		{"too big", "dm", 1000, 0, "delta must be -999 to 999", ""},
		{"too small", "dm", -1000, 0, "delta must be -999 to 999", ""},
		{"float", "dm", 1.5, 0, "malformed payload", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, kn, _ := party(t)
			p := m{"characterId": kn}
			if tc.hp != nil {
				p["hp"] = tc.hp
			}
			if tc.err != "" {
				failWith(t, r, tc.client, "character.adjust", p, tc.err)
				return
			}
			before := len(r.State().Feed)
			ok(t, r, tc.client, "character.adjust", p)
			if c := r.State().Characters[0]; c.HP != tc.wantHP {
				t.Fatalf("hp %d", c.HP)
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
	ok(t, r, "dm", "item.give", give(kn, "Field Bandages", 1)) // roster item, matched by en name
	ok(t, r, "dm", "item.give", give(kn, "ผ้าพันแผลสนาม", 1))  // and by th name
	items := r.State().Characters[0].Items
	if len(items) != 7 || items[6] != (rules.RosterItem{Name: rules.L{TH: "Rope", EN: "Rope"}, Qty: 5}) || items[4].Qty != 5 {
		t.Fatalf("items %+v", items)
	}
	if f := lastFeed(r); f.EN != "DM: gave Aldric Vane Field Bandages ×1" || f.TH != "DM: ให้ ผ้าพันแผลสนาม ×1 แก่ อัลดริก เวน" {
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
	ok(t, r, "dm", "item.give", give(wz, "thing Z", 1)) // merging is still fine when full
}

func TestItemRemove(t *testing.T) {
	rm := func(id string, index, qty int) m { return m{"characterId": id, "index": index, "qty": qty} }
	r, kn, _ := party(t) // knight: shield, chain mail, longsword, javelin 8, bandages 3, pack
	failWith(t, r, "p1", "item.remove", rm(kn, 4, 1), "only the DM can do that")
	failWith(t, r, "dm", "item.remove", rm("nope", 4, 1), "no such character")
	failWith(t, r, "dm", "item.remove", rm(kn, 6, 1), "no such item")
	failWith(t, r, "dm", "item.remove", rm(kn, -1, 1), "no such item")
	failWith(t, r, "dm", "item.remove", rm(kn, 4, 0), "qty must be at least 1")

	ok(t, r, "dm", "item.remove", rm(kn, 4, 1))
	if it := r.State().Characters[0].Items; len(it) != 6 || it[4].Qty != 2 {
		t.Fatalf("−1: %+v", it)
	}
	if f := lastFeed(r); f.EN != "DM: removed Field Bandages ×1 from Aldric Vane" || f.TH != "DM: เอา ผ้าพันแผลสนาม ×1 ออกจาก อัลดริก เวน" {
		t.Fatalf("feed %+v", f)
	}
	ok(t, r, "co", "item.remove", rm(kn, 4, 50)) // more than the stack: whole stack, feed says 2
	if f := lastFeed(r); f.EN != "DM: removed Field Bandages ×2 from Aldric Vane" {
		t.Fatalf("feed %+v", f)
	}
	ok(t, r, "dm", "item.remove", rm(kn, 0, 1)) // exactly the stack
	if it := r.State().Characters[0].Items; len(it) != 4 || it[0].Name.EN != "Chain Mail" {
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
	ok(t, r, "dm", "party.rest", m{"kind": "long"})
	if lastFeed(r).EN != "DM: long rest — HP, spell slots, uses and hit dice refilled" || len(r.State().Feed) != maxFeed {
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
	bad.Abilities.Str = 0
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "abilities must be 1-30")
	bad = e
	bad.Skills = []string{"flying"}
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "unknown skill id")
	bad = e
	bad.Art = "/uploads/" + strings.Repeat("a", 32) + ".png"
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "art must be empty or an uploaded image") // no file
	img := withUpload(t, r)
	bad.Art = "/uploads/" + strings.Repeat("d", 32) + ".png" // a directory
	failWith(t, r, "dm", "roster.update", upd("knight", bad), "art must be empty or an uploaded image")

	ok(t, r, "dm", "roster.update", upd("knight", e)) // lobby
	ro := r.State().Roster
	if len(ro) != len(rules.DefaultRoster()) || ro[0].ID != "knight" || ro[0].HeroName.EN != "Sir Test" || ro[1].ID != "wizard" || ro[1].HeroName.EN != "Soraya Quill" {
		t.Fatalf("roster %+v %+v", ro[0].HeroName, ro[1].HeroName)
	}
	e.Art = img
	ok(t, r, "dm", "roster.update", upd("rogue", e))
	if ro := r.State().Roster; ro[2].ID != "rogue" || ro[2].Art != img {
		t.Fatalf("rogue %+v", ro[2])
	}
	if ps := r.PlayerState().Roster; len(ps) != 10 || ps[2].Art != img { // only the core set is active
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

	// Reset by set drops only that set's edits.
	leg := rules.DefaultRoster()[10]
	leg.HeroName.EN = "Legend Test"
	ok(t, r, "dm", "roster.update", upd(leg.ID, leg))
	ok(t, r, "dm", "roster.update", upd("rogue", e))
	failWith(t, r, "dm", "roster.reset", m{"set": "pirate"}, "no such set")
	failWith(t, r, "dm", "roster.reset", m{"set": "legend", "classId": "rogue"}, "send classId or set, not both")
	failWith(t, r, "p1", "roster.reset", m{"set": "legend"}, "only the DM can do that")
	ok(t, r, "dm", "roster.reset", m{"set": "legend"})
	if ro := r.State().Roster; ro[10].HeroName.EN == "Legend Test" || ro[2].Art != img {
		t.Fatalf("reset legend: %q %q", ro[10].HeroName.EN, ro[2].Art)
	}
	if l := lastLog(r); l.Text.EN != "reset roster set legend" {
		t.Fatalf("log %+v", l)
	}
	ok(t, r, "dm", "roster.reset", m{"set": "core"})
	if ro := r.State().Roster; ro[2].Art != "" {
		t.Fatal("reset core")
	}
	// The client cannot move an entry to another set.
	leg.Set = "core"
	ok(t, r, "dm", "roster.update", upd(leg.ID, leg))
	if got := r.State().Roster[10]; got.Set != "legend" {
		t.Fatalf("set changed to %q", got.Set)
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
