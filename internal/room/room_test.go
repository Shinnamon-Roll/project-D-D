package room

import (
	"encoding/json"
	"strings"
	"testing"
)

func apply(r *Room, client, typ string, payload any) error {
	b, _ := json.Marshal(payload)
	return r.Apply(client, typ, b)
}

func ok(t *testing.T, r *Room, client, typ string, payload any) {
	t.Helper()
	if err := apply(r, client, typ, payload); err != nil {
		t.Fatalf("%s by %s: %v", typ, client, err)
	}
}

func fail(t *testing.T, r *Room, client, typ string, payload any) {
	t.Helper()
	before, _ := json.Marshal(r.State())
	if err := apply(r, client, typ, payload); err == nil {
		t.Fatalf("%s by %s: want error", typ, client)
	}
	after, _ := json.Marshal(r.State())
	if string(before) != string(after) {
		t.Fatalf("%s by %s: state changed on error", typ, client)
	}
}

type m = map[string]any

var goodStats = m{"str": 15, "dex": 15, "con": 13, "int": 12, "wis": 10, "cha": 10, "hp": 10, "ac": 10}

func goodChar(name string) m {
	return m{"name": name, "portraitUrl": "", "stats": goodStats, "skillIds": []string{"stealth", "arcana"}}
}

// lobby: dm in DM seat, p1 seated in the first player seat.
func lobby(t *testing.T) (*Room, string) {
	t.Helper()
	r := New()
	ok(t, r, "dm", "campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	p1Seat := r.State().Seats[1].ID
	ok(t, r, "p1", "seat.sit", m{"seatId": p1Seat, "name": "Pat"})
	return r, p1Seat
}

// playing: lobby + started + p1 character + one scene (10x5) selected.
func playing(t *testing.T) (*Room, string) {
	t.Helper()
	r, _ := lobby(t)
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "game.start", m{})
	ok(t, r, "p1", "character.submit", goodChar("Mira"))
	ok(t, r, "dm", "scene.add", m{"name": "Beach", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 10, "rows": 5}})
	sc := r.State().Scenes[0].ID
	ok(t, r, "dm", "scene.select", m{"sceneId": sc})
	return r, sc
}

func TestStatusFlow(t *testing.T) {
	r := New()
	if r.State().Status != StatusMenu {
		t.Fatal("not menu")
	}
	fail(t, r, "dm", "campaign.select", m{"campaignId": "salem", "name": "Dee"})
	fail(t, r, "dm", "campaign.select", m{"campaignId": "nope", "name": "Dee"})
	fail(t, r, "dm", "campaign.select", m{"campaignId": "starter-set", "name": ""})
	fail(t, r, "dm", "ruleset.loadPreset", m{"presetId": "dnd5e-srd"})
	ok(t, r, "dm", "campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	st := r.State()
	if st.Status != StatusLobby || st.CampaignID != "starter-set" || len(st.Seats) != 2 ||
		st.Seats[0].Kind != KindDM || st.Seats[1].Kind != KindPlayer || st.Seats[1].Occupied || st.Ruleset.StatPoints != 27 {
		t.Fatalf("after select: %+v", st)
	}
	fail(t, r, "x", "campaign.select", m{"campaignId": "starter-set", "name": "X"})

	p1 := st.Seats[1].ID
	ok(t, r, "dm", "seat.add", m{}) // second, empty player seat removed on start
	fail(t, r, "dm", "game.start", m{})
	ok(t, r, "p1", "seat.sit", m{"seatId": p1, "name": "Pat"})
	fail(t, r, "dm", "game.start", m{}) // not ready
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	fail(t, r, "p1", "game.start", m{})
	fail(t, r, "p1", "character.submit", goodChar("Mira")) // still lobby
	ok(t, r, "dm", "game.start", m{})
	st = r.State()
	if st.Status != StatusCreating || len(st.Seats) != 2 {
		t.Fatalf("after start: status %s, %d seats", st.Status, len(st.Seats))
	}
	fail(t, r, "dm", "seat.add", m{})
	fail(t, r, "p1", "seat.leave", m{})
	fail(t, r, "dm", "scene.select", m{"sceneId": "x"})
	ok(t, r, "p1", "character.submit", goodChar("Mira"))
	st = r.State()
	if st.Status != StatusPlaying {
		t.Fatalf("status %s", st.Status)
	}
	c := st.Characters[0]
	if c.SeatID != p1 || st.Seats[1].CharacterID != c.ID || len(c.Inventory) != len(st.Ruleset.StartingItems) {
		t.Fatalf("character: %+v", c)
	}
	fail(t, r, "p1", "character.submit", goodChar("Again"))
}

func TestCreatingWaitsForAllPlayers(t *testing.T) {
	r, _ := lobby(t)
	ok(t, r, "dm", "seat.add", m{})
	ok(t, r, "p2", "seat.sit", m{"seatId": r.State().Seats[2].ID, "name": "Sam"})
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "p2", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "game.start", m{})
	ok(t, r, "p1", "character.submit", goodChar("Mira"))
	if r.State().Status != StatusCreating {
		t.Fatal("moved to playing before p2 submitted")
	}
	fail(t, r, "dm", "character.submit", goodChar("DM char"))
	ok(t, r, "p2", "character.submit", goodChar("Bo"))
	if r.State().Status != StatusPlaying {
		t.Fatal("not playing")
	}
}

func TestSeats(t *testing.T) {
	r, p1 := lobby(t)
	fail(t, r, "p2", "seat.sit", m{"seatId": p1, "name": "Taken"})
	fail(t, r, "p1", "seat.sit", m{"seatId": p1, "name": "Twice"})
	fail(t, r, "p2", "seat.sit", m{"seatId": r.State().Seats[0].ID, "name": "Usurper"}) // dm seat
	fail(t, r, "p1", "seat.add", m{})
	fail(t, r, "dm", "seat.ready", m{"ready": true})

	for i := 0; i < 5; i++ {
		ok(t, r, "dm", "seat.add", m{})
	}
	fail(t, r, "dm", "seat.add", m{}) // 6 player seats max

	ok(t, r, "p2", "seat.sit", m{"seatId": r.State().Seats[2].ID, "name": strings.Repeat("n", 32)})
	fail(t, r, "p3", "seat.sit", m{"seatId": r.State().Seats[3].ID, "name": strings.Repeat("n", 33)})
	fail(t, r, "p3", "seat.sit", m{"seatId": r.State().Seats[3].ID, "name": "   "})

	// co-DM: max 1, can do DM things but cannot remove the dm seat.
	fail(t, r, "p1", "seat.openCoDm", m{})
	ok(t, r, "dm", "seat.openCoDm", m{})
	fail(t, r, "dm", "seat.openCoDm", m{})
	co := r.State().Seats[len(r.State().Seats)-1]
	if co.Kind != KindCoDM {
		t.Fatalf("last seat %+v", co)
	}
	ok(t, r, "co", "seat.sit", m{"seatId": co.ID, "name": "Cora"})
	fail(t, r, "co", "seat.remove", m{"seatId": r.State().Seats[0].ID})
	fail(t, r, "dm", "seat.remove", m{"seatId": r.State().Seats[0].ID})
	fail(t, r, "dm", "seat.leave", m{})
	fail(t, r, "dm", "seat.remove", m{"seatId": "nope"})

	// co-DM kicks p2; p2 is no longer seated.
	p2 := r.State().Seats[2].ID
	ok(t, r, "co", "seat.remove", m{"seatId": p2})
	if r.SeatID("p2") != "" {
		t.Fatal("p2 still seated")
	}
	fail(t, r, "p2", "seat.ready", m{"ready": true})

	// leave frees the seat and clears ready.
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "p1", "seat.leave", m{})
	if s := r.seat(p1); s.Occupied || s.Ready || s.Name != "" || r.SeatID("p1") != "" {
		t.Fatalf("after leave: %+v", s)
	}
	ok(t, r, "p4", "seat.sit", m{"seatId": p1, "name": "New"})
}

func TestNeverLeakClientIDs(t *testing.T) {
	r, _ := lobby(t)
	b, _ := json.Marshal(r.State())
	for _, secret := range []string{"p1", "client"} { // "dm" is also a seat kind, so check p1 only
		if strings.Contains(string(b), secret) {
			t.Fatalf("state JSON contains %q: %s", secret, b)
		}
	}
	if strings.Contains(string(b), "null") {
		t.Fatalf("state JSON contains null: %s", b)
	}
	if r.SeatID("p1") == "" || r.SeatID("stranger") != "" || r.SeatID("") != "" {
		t.Fatal("SeatID")
	}
}

func TestPermissionsAndUnknown(t *testing.T) {
	r, sc := playing(t)
	dmOnly := []struct {
		typ     string
		payload any
	}{
		{"ruleset.loadPreset", m{"presetId": "dnd5e-srd"}},
		{"ruleset.update", m{"ruleset": m{"name": "Custom"}}},
		{"character.grant", m{"characterId": r.State().Characters[0].ID, "setStats": m{"str": 20}}},
		{"scene.add", m{"name": "X", "imageUrl": "/uploads/b.png", "grid": m{"cellSize": 50, "cols": 5, "rows": 5}}},
		{"scene.update", m{"sceneId": sc, "name": "X", "grid": m{"cellSize": 50, "cols": 5, "rows": 5}}},
		{"scene.select", m{"sceneId": sc}},
		{"token.addNpc", m{"sceneId": sc, "name": "Goblin", "imageUrl": ""}},
		{"token.remove", m{"tokenId": r.State().Tokens[0].ID}},
		{"fog.add", m{"sceneId": sc, "col": 0, "row": 0, "w": 1, "h": 1}},
	}
	for _, e := range dmOnly {
		fail(t, r, "p1", e.typ, e.payload)
		fail(t, r, "stranger", e.typ, e.payload)
	}
	fail(t, r, "dm", "nope", m{})
	fail(t, r, "dm", "dice.roll", m{"notation": "1d20"}) // removed 2026-10-16: physical dice
	fail(t, r, "dm", "hello", m{"clientId": "x"})
	if err := r.Apply("dm", "ruleset.loadPreset", json.RawMessage(`{"presetId":`)); err == nil {
		t.Fatal("malformed payload accepted")
	}
	fail(t, r, "dm", "ruleset.loadPreset", m{"presetId": 5})
	fail(t, r, "dm", "ruleset.loadPreset", m{"presetId": "other"})
	fail(t, r, "dm", "ruleset.update", m{"ruleset": m{"name": "Bad", "statPoints": -1}})
	ok(t, r, "dm", "ruleset.update", m{"ruleset": m{"name": "Custom"}})
	if rs := r.State().Ruleset; rs.Name != "Custom" || rs.Stats == nil || rs.StartingItems == nil {
		t.Fatalf("ruleset: %+v", rs)
	}
	ok(t, r, "dm", "ruleset.loadPreset", m{"presetId": "dnd5e-srd"})
}

func TestCharacterValidation(t *testing.T) {
	r, _ := lobby(t)
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "game.start", m{})
	bad := []m{
		{"name": "", "stats": goodStats, "skillIds": []string{"stealth", "arcana"}},
		{"name": "A", "portraitUrl": "http://evil/x.png", "stats": goodStats, "skillIds": []string{"stealth", "arcana"}},
		{"name": "A", "portraitUrl": "/uploads/../x", "stats": goodStats, "skillIds": []string{"stealth", "arcana"}},
		{"name": "A", "stats": m{"str": 15}, "skillIds": []string{"stealth", "arcana"}},
		{"name": "A", "stats": goodStats, "skillIds": []string{"stealth"}},
		{"name": "A", "stats": goodStats, "skillIds": []string{"stealth", "stealth"}},
		{"name": "A", "stats": goodStats, "skillIds": []string{"stealth", "flying"}},
	}
	over := m{}
	for k, v := range goodStats {
		over[k] = v
	}
	over["cha"] = 11
	bad = append(bad, m{"name": "A", "stats": over, "skillIds": []string{"stealth", "arcana"}})
	for _, p := range bad {
		fail(t, r, "p1", "character.submit", p)
	}
	ok(t, r, "p1", "character.submit", m{"name": "A", "portraitUrl": "/uploads/p.png", "stats": goodStats, "skillIds": []string{"stealth", "arcana"}})
}

func TestGrant(t *testing.T) {
	r, _ := playing(t)
	cid := r.State().Characters[0].ID
	fail(t, r, "dm", "character.grant", m{"characterId": "nope"})
	fail(t, r, "dm", "character.grant", m{"characterId": cid, "setStats": m{"luck": 1}})
	fail(t, r, "dm", "character.grant", m{"characterId": cid, "addSkillIds": []string{"flying"}})
	fail(t, r, "dm", "character.grant", m{"characterId": cid, "addItems": []m{{"itemId": "sword", "qty": 1}}})
	fail(t, r, "dm", "character.grant", m{"characterId": cid, "addItems": []m{{"itemId": "rope", "qty": 0}}})
	// One bad part rejects the whole grant.
	fail(t, r, "dm", "character.grant", m{"characterId": cid, "setStats": m{"str": 20}, "removeItems": []m{{"itemId": "rope", "qty": -1}}})

	ok(t, r, "dm", "character.grant", m{
		"characterId": cid, "setStats": m{"str": 20, "hp": 12},
		"addSkillIds": []string{"history", "stealth"}, "removeSkillIds": []string{"arcana"},
		"addItems":    []m{{"itemId": "healing-potion", "qty": 2}, {"itemId": "torch", "qty": 1}},
		"removeItems": []m{{"itemId": "rope", "qty": 5}, {"itemId": "rations", "qty": 2}},
	})
	c := r.State().Characters[0]
	if c.Stats["str"] != 20 || c.Stats["hp"] != 12 {
		t.Fatalf("stats %v", c.Stats)
	}
	if strings.Join(c.SkillIDs, ",") != "stealth,history" {
		t.Fatalf("skills %v", c.SkillIDs)
	}
	qty := map[string]int{}
	for _, it := range c.Inventory {
		qty[it.ItemID] = it.Qty
	}
	if _, has := qty["rope"]; has || qty["rations"] != 3 || qty["torch"] != 6 || qty["healing-potion"] != 2 {
		t.Fatalf("inventory %v", qty)
	}
}

func TestScenesTokensFog(t *testing.T) {
	r, sc := playing(t)
	st := r.State()
	if st.CurrentSceneID != sc || len(st.Tokens) != 1 || st.Tokens[0].Col != 0 || st.Tokens[0].Row != 0 || st.Tokens[0].Name != "Mira" {
		t.Fatalf("auto place: %+v", st.Tokens)
	}
	// Selecting again does not duplicate.
	ok(t, r, "dm", "scene.select", m{"sceneId": sc})
	if len(r.State().Tokens) != 1 {
		t.Fatal("duplicate token")
	}
	tok := r.State().Tokens[0].ID

	// scene.add validation.
	grid := m{"cellSize": 50, "cols": 10, "rows": 5}
	fail(t, r, "dm", "scene.add", m{"name": "X", "imageUrl": "https://x/y.png", "grid": grid})
	fail(t, r, "dm", "scene.add", m{"name": "X", "imageUrl": "/uploads/", "grid": grid})
	for _, g := range []m{
		{"cellSize": 9, "cols": 10, "rows": 5}, {"cellSize": 401, "cols": 10, "rows": 5},
		{"cellSize": 50, "cols": 0, "rows": 5}, {"cellSize": 50, "cols": 10, "rows": 201},
	} {
		fail(t, r, "dm", "scene.add", m{"name": "X", "imageUrl": "/uploads/x.png", "grid": g})
	}

	// Moves.
	ok(t, r, "p1", "token.move", m{"tokenId": tok, "col": 9, "row": 4})
	fail(t, r, "p1", "token.move", m{"tokenId": tok, "col": 10, "row": 0})
	fail(t, r, "p1", "token.move", m{"tokenId": tok, "col": 0, "row": 5})
	fail(t, r, "p1", "token.move", m{"tokenId": tok, "col": -1, "row": 0})
	fail(t, r, "stranger", "token.move", m{"tokenId": tok, "col": 1, "row": 1})
	ok(t, r, "dm", "token.addNpc", m{"sceneId": sc, "name": "Goblin", "imageUrl": "/uploads/g.png"})
	fail(t, r, "dm", "token.addNpc", m{"sceneId": sc, "name": "Goblin", "imageUrl": "/etc/passwd"})
	fail(t, r, "dm", "token.addNpc", m{"sceneId": "nope", "name": "Goblin"})
	npc := r.State().Tokens[1].ID
	fail(t, r, "p1", "token.move", m{"tokenId": npc, "col": 1, "row": 1}) // NPC
	ok(t, r, "dm", "token.move", m{"tokenId": npc, "col": 1, "row": 1})
	ok(t, r, "dm", "token.move", m{"tokenId": tok, "col": 2, "row": 2})

	// Shrinking the grid pulls tokens in.
	ok(t, r, "dm", "scene.update", m{"sceneId": sc, "name": "Beach 2", "grid": m{"cellSize": 40, "cols": 2, "rows": 2}})
	if tk := r.State().Tokens[0]; tk.Col != 1 || tk.Row != 1 {
		t.Fatalf("token after shrink %+v", tk)
	}
	ok(t, r, "dm", "scene.update", m{"sceneId": sc, "name": "Beach", "grid": grid})

	// Fog bounds.
	ok(t, r, "dm", "fog.add", m{"sceneId": sc, "col": 0, "row": 0, "w": 10, "h": 5})
	ok(t, r, "dm", "fog.add", m{"sceneId": sc, "col": 9, "row": 4, "w": 1, "h": 1})
	for _, f := range []m{
		{"col": 0, "row": 0, "w": 0, "h": 1}, {"col": 0, "row": 0, "w": 1, "h": 0},
		{"col": 0, "row": 0, "w": 11, "h": 1}, {"col": 9, "row": 0, "w": 2, "h": 1},
		{"col": 0, "row": 4, "w": 1, "h": 2}, {"col": -1, "row": 0, "w": 1, "h": 1},
	} {
		f["sceneId"] = sc
		fail(t, r, "dm", "fog.add", f)
	}
	fog := r.State().Scenes[0].Fog[0].ID
	fail(t, r, "p1", "fog.remove", m{"sceneId": sc, "fogId": fog})
	fail(t, r, "dm", "fog.remove", m{"sceneId": sc, "fogId": "nope"})
	ok(t, r, "dm", "fog.remove", m{"sceneId": sc, "fogId": fog})
	if len(r.State().Scenes[0].Fog) != 1 {
		t.Fatal("fog not removed")
	}

	ok(t, r, "dm", "token.remove", m{"tokenId": npc})
	fail(t, r, "dm", "token.remove", m{"tokenId": npc})
}

func TestSceneAutoPlaceWraps(t *testing.T) {
	r, _ := lobby(t)
	for i := 0; i < 2; i++ {
		ok(t, r, "dm", "seat.add", m{})
	}
	for i, c := range []string{"p2", "p3"} {
		ok(t, r, c, "seat.sit", m{"seatId": r.State().Seats[2+i].ID, "name": c})
	}
	for _, c := range []string{"p1", "p2", "p3"} {
		ok(t, r, c, "seat.ready", m{"ready": true})
	}
	ok(t, r, "dm", "game.start", m{})
	for _, c := range []string{"p1", "p2", "p3"} {
		ok(t, r, c, "character.submit", goodChar(c))
	}
	ok(t, r, "dm", "scene.add", m{"name": "Hall", "imageUrl": "/uploads/h.png", "grid": m{"cellSize": 50, "cols": 2, "rows": 3}})
	ok(t, r, "dm", "scene.select", m{"sceneId": r.State().Scenes[0].ID})
	var got []string
	for _, tk := range r.State().Tokens {
		got = append(got, string(rune('0'+tk.Col))+string(rune('0'+tk.Row)))
	}
	if strings.Join(got, " ") != "00 10 01" {
		t.Fatalf("placements %v", got)
	}
}

func TestSummary(t *testing.T) {
	if s := New().Summary(); s.Status != StatusMenu || s.Campaign != "" || s.Seats != 0 {
		t.Fatalf("empty room %+v", s)
	}
	r, _ := lobby(t)
	ok(t, r, "dm", "seat.add", m{})
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "seat.openCoDm", m{})
	for _, x := range r.State().Seats {
		if x.Kind == KindCoDM {
			ok(t, r, "co", "seat.sit", m{"seatId": x.ID, "name": "Cody"})
		}
	}
	s := r.Summary()
	if s.Status != StatusLobby || s.Campaign == "" || s.DM != "Dee" || s.CoDM != "Cody" || s.Players != 1 || s.Seats != 2 || s.Ready != 1 {
		t.Fatalf("lobby %+v", s)
	}
	r, _ = playing(t)
	s = r.Summary()
	if s.Status != StatusPlaying || s.Characters != 1 || s.Scenes != 1 || s.Scene != "Beach" {
		t.Fatalf("playing %+v", s)
	}
}

func TestSceneUpdateShrinkClampsAndClips(t *testing.T) {
	type pos struct{ col, row int }
	type rect struct{ col, row, w, h int }
	cases := []struct {
		name       string
		cols, rows int
		token      pos    // start 9,4 on the 10x5 grid
		fog        []rect // after update; start rects below
	}{
		{"no change", 10, 5, pos{9, 4}, []rect{{0, 0, 10, 5}, {8, 3, 2, 2}, {9, 0, 1, 1}}},
		{"grow", 20, 20, pos{9, 4}, []rect{{0, 0, 10, 5}, {8, 3, 2, 2}, {9, 0, 1, 1}}},
		{"shrink cols", 9, 5, pos{8, 4}, []rect{{0, 0, 9, 5}, {8, 3, 1, 2}}},
		{"shrink rows", 10, 4, pos{9, 3}, []rect{{0, 0, 10, 4}, {8, 3, 2, 1}, {9, 0, 1, 1}}},
		{"shrink both", 3, 2, pos{2, 1}, []rect{{0, 0, 3, 2}}},
		{"to 1x1", 1, 1, pos{0, 0}, []rect{{0, 0, 1, 1}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			ok(t, r, "dm", "token.move", m{"tokenId": r.State().Tokens[0].ID, "col": 9, "row": 4})
			ok(t, r, "dm", "token.addNpc", m{"sceneId": sc, "name": "Gob"}) // at 0,0: stays
			for _, f := range []rect{{0, 0, 10, 5}, {8, 3, 2, 2}, {9, 0, 1, 1}} {
				ok(t, r, "dm", "fog.add", m{"sceneId": sc, "col": f.col, "row": f.row, "w": f.w, "h": f.h})
			}
			ok(t, r, "dm", "scene.update", m{"sceneId": sc, "name": "Beach", "grid": m{"cellSize": 50, "cols": tc.cols, "rows": tc.rows}})
			st := r.State()
			if tk := st.Tokens[0]; tk.Col != tc.token.col || tk.Row != tc.token.row {
				t.Errorf("token %d,%d want %+v", tk.Col, tk.Row, tc.token)
			}
			if npc := st.Tokens[1]; npc.Col != 0 || npc.Row != 0 {
				t.Errorf("npc moved to %d,%d", npc.Col, npc.Row)
			}
			var got []rect
			for _, f := range st.Scenes[0].Fog {
				got = append(got, rect{f.Col, f.Row, f.W, f.H})
			}
			if len(got) != len(tc.fog) {
				t.Fatalf("fog %v want %v", got, tc.fog)
			}
			for i := range got {
				if got[i] != tc.fog[i] {
					t.Fatalf("fog %v want %v", got, tc.fog)
				}
			}
		})
	}
}
