package room

import (
	"slices"
	"strings"
	"testing"
)

func TestSplitName(t *testing.T) {
	for _, tc := range []struct {
		in   string
		base string
		n    int
	}{
		{"Goblin", "Goblin", 1},
		{"Goblin #2", "Goblin", 2},
		{"Goblin #15", "Goblin", 15},
		{"Goblin #2 #3", "Goblin #2", 3},
		{"Goblin #0", "Goblin #0", 1},   // no number 0
		{"Goblin #07", "Goblin #07", 1}, // no leading zero
		{"Goblin#2", "Goblin#2", 1},     // needs the space
		{"#2", "#2", 1},                 // needs a base
		{"Goblin #", "Goblin #", 1},
		{"Goblin #-2", "Goblin #-2", 1},
		{"Goblin #999999999", "Goblin", 999999999},
		{"Goblin #1000000000", "Goblin #1000000000", 1}, // past 9 digits: a bare base
		{"ก๊อบลิน #3", "ก๊อบลิน", 3},
	} {
		if b, n := splitName(tc.in); b != tc.base || n != tc.n {
			t.Errorf("splitName(%q) = %q, %d; want %q, %d", tc.in, b, n, tc.base, tc.n)
		}
	}
}

// sceneNames lists the token names of a scene in token order (the DM's full view).
func sceneNames(r *Room, sc string) []string {
	var out []string
	for _, tk := range r.State().Tokens {
		if tk.SceneID == sc {
			out = append(out, tk.Name)
		}
	}
	return out
}

// lastID is the id of the newest token.
func lastID(r *Room) string { return r.st.Tokens[len(r.st.Tokens)-1].ID }

// idOf is the id of the token named name in scene sc.
func idOf(t *testing.T, r *Room, sc, name string) string {
	t.Helper()
	for _, tk := range r.st.Tokens {
		if tk.SceneID == sc && tk.Name == name {
			return tk.ID
		}
	}
	t.Fatalf("no token %q in %s: %v", name, sc, sceneNames(r, sc))
	return ""
}

// Visible tokens only (spec examples): the next number is the highest of that base + 1, never a count.
func TestNpcNameRules(t *testing.T) {
	const hero = "Aldric Vane" // playing(): p1's knight token
	type step struct {
		place string // pool NPC name to npc.place (added to the pool on first use), or
		dup   string // token name to token.duplicate, or
		del   string // token name to token.remove
	}
	for _, tc := range []struct {
		name  string
		steps []step
		want  []string // scene token names after the steps, in order (the hero first)
	}{
		{"first is bare, then #2 #3 #4", []step{{place: "Goblin"}, {place: "Goblin"}, {place: "Goblin"}, {place: "Goblin"}},
			[]string{hero, "Goblin", "Goblin #2", "Goblin #3", "Goblin #4"}},
		{"duplicate the bare one", []step{{place: "Goblin"}, {dup: "Goblin"}, {dup: "Goblin"}},
			[]string{hero, "Goblin", "Goblin #2", "Goblin #3"}},
		{"#5 is the highest: next #6", []step{{place: "Goblin"}, {place: "Goblin #5"}, {dup: "Goblin"}},
			[]string{hero, "Goblin", "Goblin #5", "Goblin #6"}},
		{"5 goblins, duplicate #2: #6 at the end", []step{
			{place: "Goblin"}, {place: "Goblin"}, {place: "Goblin"}, {place: "Goblin"}, {place: "Goblin"}, {dup: "Goblin #2"}},
			[]string{hero, "Goblin", "Goblin #2", "Goblin #3", "Goblin #4", "Goblin #5", "Goblin #6"}},
		{"pool name with #5 starts at 5", []step{{place: "Goblin #5"}, {place: "Goblin #5"}, {dup: "Goblin #5"}},
			[]string{hero, "Goblin #5", "Goblin #6", "Goblin #7"}},
		{"pool #5 after a bare one is still 5", []step{{place: "Goblin"}, {place: "Goblin #5"}},
			[]string{hero, "Goblin", "Goblin #5"}},
		{"pool #5 below the highest: highest + 1", []step{{place: "Goblin #5"}, {place: "Goblin #5"}, {place: "Goblin #5"}, {place: "Goblin"}},
			[]string{hero, "Goblin #5", "Goblin #6", "Goblin #7", "Goblin #8"}},
		{"case-insensitive base", []step{{place: "Goblin"}, {place: "goblin #2"}, {dup: "goblin #2"}, {place: "GOBLIN #3"}},
			[]string{hero, "Goblin", "goblin #2", "goblin #3", "GOBLIN #4"}},
		{"case-insensitive against a character", []step{{place: "aldric VANE"}},
			[]string{hero, "aldric VANE #2"}},
		{"a freed top number is reused, gaps are not filled", []step{
			{place: "Goblin"}, {place: "Goblin"}, {place: "Goblin"}, {del: "Goblin #3"}, {del: "Goblin"}, {place: "Goblin"}},
			[]string{hero, "Goblin #2", "Goblin #3"}},
		{"a character name counts as taken and is kept", []step{{place: hero}, {dup: hero + " #2"}},
			[]string{hero, hero + " #2", hero + " #3"}},
		{"other bases are separate", []step{{place: "Goblin"}, {place: "Wolf"}, {place: "Goblin Boss"}, {place: "Wolf"}},
			[]string{hero, "Goblin", "Wolf", "Goblin Boss", "Wolf #2"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			pool := map[string]string{}
			for _, s := range tc.steps {
				switch {
				case s.place != "":
					if _, ok2 := pool[s.place]; !ok2 {
						ok(t, r, "dm", "npc.add", m{"name": s.place, "hpMax": 5})
						pool[s.place] = r.npcs.List()[len(r.npcs.List())-1].ID
					}
					ok(t, r, "dm", "npc.place", m{"npcId": pool[s.place], "sceneId": sc})
				case s.dup != "":
					ok(t, r, "dm", "token.duplicate", m{"tokenId": idOf(t, r, sc, s.dup)})
				default:
					ok(t, r, "dm", "token.remove", m{"tokenId": idOf(t, r, sc, s.del)})
				}
			}
			if got := sceneNames(r, sc); !slices.Equal(got, tc.want) {
				t.Fatalf("names %q, want %q", got, tc.want)
			}
			for _, e := range r.DMLog() {
				if strings.HasPrefix(e.Text.EN, "renamed") {
					t.Fatalf("visible-only steps renamed a token: %+v", e)
				}
			}
		})
	}
}

func TestNpcNamesPerScene(t *testing.T) {
	r, beach := playing(t)
	ok(t, r, "dm", "scene.add", m{"name": "Cave", "imageUrl": "/uploads/b.png", "grid": m{"cellSize": 50, "cols": 4, "rows": 4}})
	cave := r.State().Scenes[1].ID
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": beach})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": beach})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": cave})
	ok(t, r, "dm", "token.duplicate", m{"tokenId": idOf(t, r, beach, "Goblin #2")})
	if got := sceneNames(r, cave); !slices.Equal(got, []string{"Goblin"}) {
		t.Fatalf("cave %q", got)
	}
	if got := sceneNames(r, beach); !slices.Equal(got, []string{"Aldric Vane", "Goblin", "Goblin #2", "Goblin #3"}) {
		t.Fatalf("beach %q", got)
	}
}

// renames lists the DM log's automatic-rename lines (EN).
func renames(r *Room) []string {
	var out []string
	for _, e := range r.DMLog() {
		if strings.HasPrefix(e.Text.EN, "renamed ") {
			out = append(out, e.Text.EN)
		}
	}
	return out
}

// playerNames lists the token names players see, and fails if any player-visible text holds a hidden token's name.
func playerNames(t *testing.T, r *Room) []string {
	t.Helper()
	var out []string
	for _, tk := range r.PlayerState().Tokens {
		out = append(out, tk.Name)
	}
	pl := string(mustJSON(r.PlayerState()))
	for _, tk := range r.st.Tokens {
		if tk.Hidden && strings.Contains(pl, `"`+tk.Name+`"`) {
			t.Fatalf("player state names hidden %q: %s", tk.Name, pl)
		}
	}
	return out
}

// A player must never infer a hidden token from the numbers (spec "ความลับ").
func TestNpcNamesHiddenSecrecy(t *testing.T) {
	type step struct {
		place  string // "v" or "h": npc.place Goblin visible or hidden
		dup    string // token name to duplicate
		reveal string // token name to reveal
		hide   string // token name to hide
	}
	for _, tc := range []struct {
		name    string
		steps   []step
		want    []string // DM view, token order (the hero first)
		players []string // player view
		renames []string // DM log rename lines
	}{
		{"hidden gets the scene-wide next", []step{{place: "v"}, {place: "h"}},
			[]string{"Goblin", "Goblin #2"}, []string{"Goblin"}, nil},
		{"visible after hidden: numbered from visible, hidden renamed", []step{{place: "v"}, {place: "h"}, {place: "v"}},
			[]string{"Goblin", "Goblin #3", "Goblin #2"}, []string{"Goblin", "Goblin #2"},
			[]string{"renamed Goblin #2 to Goblin #3"}},
		{"only hidden ones: visible is bare", []step{{place: "h"}, {place: "h"}, {place: "v"}},
			[]string{"Goblin #3", "Goblin #2", "Goblin"}, []string{"Goblin"},
			[]string{"renamed Goblin to Goblin #3"}},
		{"duplicate of hidden is hidden, scene-wide next", []step{{place: "v"}, {place: "h"}, {dup: "Goblin #2"}},
			[]string{"Goblin", "Goblin #2", "Goblin #3"}, []string{"Goblin"}, nil},
		{"duplicate of visible skips hidden numbers and displaces", []step{{place: "v"}, {place: "h"}, {dup: "Goblin"}},
			[]string{"Goblin", "Goblin #3", "Goblin #2"}, []string{"Goblin", "Goblin #2"},
			[]string{"renamed Goblin #2 to Goblin #3"}},
		{"reveal keeps a name that is already next", []step{{place: "v"}, {place: "h"}, {reveal: "Goblin #2"}},
			[]string{"Goblin", "Goblin #2"}, []string{"Goblin", "Goblin #2"}, nil},
		{"reveal renames past a gap", []step{{place: "v"}, {place: "h"}, {place: "h"}, {reveal: "Goblin #3"}},
			[]string{"Goblin", "Goblin #3", "Goblin #2"}, []string{"Goblin", "Goblin #2"},
			[]string{"renamed Goblin #3 to Goblin #2", "renamed Goblin #2 to Goblin #3"}},
		{"reveal with no visible one: bare", []step{{place: "h"}, {place: "h"}, {reveal: "Goblin #2"}},
			[]string{"Goblin #2", "Goblin"}, []string{"Goblin"}, // the two swap names
			[]string{"renamed Goblin #2 to Goblin", "renamed Goblin to Goblin #2"}},
		{"hide keeps the name", []step{{place: "v"}, {place: "v"}, {hide: "Goblin"}},
			[]string{"Goblin", "Goblin #2"}, []string{"Goblin #2"}, nil},
		{"hidden then revealed again: next visible", []step{{place: "v"}, {place: "v"}, {hide: "Goblin #2"}, {place: "v"}, {reveal: "Goblin #3"}},
			[]string{"Goblin", "Goblin #3", "Goblin #2"}, []string{"Goblin", "Goblin #3", "Goblin #2"},
			[]string{"renamed Goblin #2 to Goblin #3"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5})
			for _, s := range tc.steps {
				switch {
				case s.place != "":
					ok(t, r, "dm", "npc.update", m{"id": "npc1", "name": "Goblin", "hpMax": 5, "hidden": s.place == "h"})
					ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
				case s.dup != "":
					ok(t, r, "dm", "token.duplicate", m{"tokenId": idOf(t, r, sc, s.dup)})
				case s.reveal != "":
					ok(t, r, "dm", "token.setHidden", m{"tokenId": idOf(t, r, sc, s.reveal), "hidden": false})
				default:
					ok(t, r, "dm", "token.setHidden", m{"tokenId": idOf(t, r, sc, s.hide), "hidden": true})
				}
			}
			if got := sceneNames(r, sc)[1:]; !slices.Equal(got, tc.want) {
				t.Fatalf("DM names %q, want %q", got, tc.want)
			}
			if got := playerNames(t, r)[1:]; !slices.Equal(got, tc.players) {
				t.Fatalf("player names %q, want %q", got, tc.players)
			}
			if got := renames(r); !slices.Equal(got, tc.renames) {
				t.Fatalf("rename lines %q, want %q", got, tc.renames)
			}
			for _, f := range r.State().Feed { // no rename or hidden number in the public feed
				if strings.Contains(f.Text.EN, "renamed") || strings.Contains(f.Text.EN, "hidden") {
					t.Fatalf("feed leaks: %+v", f)
				}
			}
			seen := map[string]bool{} // names stay unique, case-insensitive
			for _, n := range sceneNames(r, sc) {
				if k := strings.ToLower(n); seen[k] {
					t.Fatalf("duplicate name %q: %q", n, sceneNames(r, sc))
				} else {
					seen[k] = true
				}
			}
		})
	}
}

// A revealed token starts at its pool NPC's own number, like npc.place.
func TestRevealKeepsPoolNumber(t *testing.T) {
	r, sc := playing(t)
	ok(t, r, "dm", "npc.add", m{"name": "Orc #5", "hpMax": 5, "hidden": true})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	ok(t, r, "dm", "token.setHidden", m{"tokenId": lastID(r), "hidden": false})
	if got := sceneNames(r, sc)[1:]; !slices.Equal(got, []string{"Orc #5"}) || renames(r) != nil {
		t.Fatalf("names %q, renames %q", got, renames(r))
	}
}

// Hidden-only changes send players nothing new: their view and feed stay byte-identical.
func TestHiddenNamingSendsPlayersNothing(t *testing.T) {
	r, sc := playing(t)
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5, "hidden": true})
	ok(t, r, "dm", "npc.add", m{"name": "Wolf", "hpMax": 5})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": sc})
	before := string(mustJSON(r.PlayerState()))
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	ok(t, r, "dm", "token.duplicate", m{"tokenId": lastID(r)}) // hidden copy
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	if after := string(mustJSON(r.PlayerState())); after != before {
		t.Fatalf("player view changed:\n%s\n%s", before, after)
	}
	if got := sceneNames(r, sc)[2:]; !slices.Equal(got, []string{"Goblin", "Goblin #2", "Goblin #3"}) {
		t.Fatalf("hidden names %q", got)
	}
}

// A pending attack keeps the name the player saw while its target is hidden, and takes the new name on reveal.
func TestRevealRenamesPendingAttack(t *testing.T) {
	r, sc := playing(t)
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	target := lastID(r) // "Goblin #2"
	ok(t, r, "dm", "combat.set", m{"on": true})
	ok(t, r, "p1", "attack.declare", m{"targetTokenId": target})
	ok(t, r, "dm", "token.setHidden", m{"tokenId": target, "hidden": true})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc}) // visible "Goblin #2": the hidden one becomes #3
	if n := r.token(target).Name; n != "Goblin #3" {
		t.Fatalf("hidden target %q", n)
	}
	if a := r.MyAttack("p1"); a == nil || a.TargetName != "Goblin #2" {
		t.Fatalf("myAttack leaks the hidden rename: %+v", a)
	}
	ok(t, r, "dm", "token.setHidden", m{"tokenId": target, "hidden": false})
	if a := r.MyAttack("p1"); a == nil || a.TargetName != "Goblin #3" || r.Attacks()[0].TargetName != "Goblin #3" {
		t.Fatalf("attack not renamed on reveal: %+v %+v", a, r.Attacks())
	}
}

func TestTokenDuplicateRefused(t *testing.T) {
	r, sc := playing(t)
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	gob, hero := lastID(r), r.st.Tokens[0].ID
	for _, tc := range []struct {
		name, client string
		payload      any
		want         string
	}{
		{"player", "p1", m{"tokenId": gob}, "only the DM can do that"},
		{"unseated", "stranger", m{"tokenId": gob}, "only the DM can do that"},
		{"malformed", "dm", `{"tokenId":7}`, "malformed payload"},
		{"missing", "dm", m{"tokenId": "t999"}, "no such token"},
		{"no id", "dm", m{}, "no such token"},
		{"character token", "dm", m{"tokenId": hero}, "only NPC tokens from the pool"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			failWith(t, r, tc.client, "token.duplicate", tc.payload, tc.want)
			if len(r.Fresh()) != 0 {
				t.Fatal("log left after a failed event")
			}
		})
	}
	failWith(t, creating(t), "dm", "token.duplicate", m{"tokenId": "t1"}, "not allowed while the room is creating")
}

func TestTokenDuplicate(t *testing.T) {
	r, sc := playing(t) // hero at 0,0 of 10x5
	img := withUpload(t, r)
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 9, "imageUrl": img})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc}) // (1,0)
	orig := lastID(r)
	ok(t, r, "dm", "token.adjustHp", m{"tokenId": orig, "delta": -4})
	ok(t, r, "dm", "npc.update", m{"id": "npc1", "name": "Goblin", "hpMax": 50}) // the copy follows the token, not the pool
	feed := len(r.st.Feed)

	ok(t, r, "dm", "token.duplicate", m{"tokenId": orig})
	c := r.st.Tokens[len(r.st.Tokens)-1]
	if c.ID == orig || c.Name != "Goblin #2" || c.ImageURL != img || *c.HPMax != 9 || *c.HP != 9 || c.Hidden ||
		c.NpcID != "npc1" || c.CharacterID != "" || c.SceneID != sc || c.Col != 1 || c.Row != 0 {
		t.Fatalf("copy %+v", c)
	}
	if *r.token(orig).HP != 5 {
		t.Fatal("original changed")
	}
	if got := r.st.Feed[feed:]; len(got) != 1 || got[0].Text.EN != "added Goblin #2" || got[0].Text.TH != "เพิ่ม Goblin #2" {
		t.Fatalf("feed %+v", got)
	}
	if l := lastLog(r); l.Text.EN != "added Goblin #2" || l.Actor != "Dee" {
		t.Fatalf("dm log %+v", l)
	}

	// The original's cell fills to MaxPerCell, then the nearest cell (Chebyshev, then row-major): (0,0) with the hero.
	ok(t, r, "dm", "token.duplicate", m{"tokenId": orig})
	ok(t, r, "dm", "token.duplicate", m{"tokenId": orig})
	ok(t, r, "dm", "token.duplicate", m{"tokenId": orig})
	if c := r.st.Tokens[len(r.st.Tokens)-1]; c.Name != "Goblin #5" || c.Col != 0 || c.Row != 0 {
		t.Fatalf("5th goblin %+v", c)
	}

	// Hidden copies: DM log only, and they count every token for the cell.
	ok(t, r, "dm", "token.setHidden", m{"tokenId": orig, "hidden": true})
	feed = len(r.st.Feed)
	ok(t, r, "dm", "token.duplicate", m{"tokenId": orig})
	h := r.st.Tokens[len(r.st.Tokens)-1]
	if !h.Hidden || h.Name != "Goblin #6" || h.Col != 0 || h.Row != 0 {
		t.Fatalf("hidden copy %+v", h) // (1,0) holds 4 tokens with the hidden original; (0,0) holds 2
	}
	if len(r.st.Feed) != feed || lastLog(r).Text.EN != "added Goblin #6 (hidden)" || lastLog(r).Text.TH != "เพิ่ม Goblin #6 (ซ่อน)" {
		t.Fatalf("hidden copy logged publicly: feed %+v log %+v", r.st.Feed[feed:], lastLog(r))
	}
	if strings.Contains(string(mustJSON(r.PlayerState())), "Goblin #6") {
		t.Fatal("player state names the hidden copy")
	}

	// A visible copy ignores hidden tokens for the cell (3 visible on (1,0)) and the number: visible max #5 gives #6,
	// so the hidden #6 is renamed to the scene-wide next.
	ok(t, r, "dm", "token.duplicate", m{"tokenId": idOf(t, r, sc, "Goblin #2")})
	if v := r.st.Tokens[len(r.st.Tokens)-1]; v.Name != "Goblin #6" || v.Hidden || v.Col != 1 || v.Row != 0 {
		t.Fatalf("visible copy %+v", v)
	}
	if h := r.token(h.ID); h.Name != "Goblin #7" || !h.Hidden {
		t.Fatalf("displaced hidden copy %+v", h)
	}
	if got := renames(r); !slices.Equal(got, []string{"renamed Goblin #6 to Goblin #7"}) {
		t.Fatalf("renames %q", got)
	}
	if f := r.st.Feed[len(r.st.Feed)-1].Text.EN; f != "added Goblin #6" {
		t.Fatalf("feed %q", f)
	}
}

func TestTokenDuplicateFullScene(t *testing.T) {
	r, _ := playing(t)
	ok(t, r, "dm", "scene.add", m{"name": "Tiny", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 1, "rows": 1}})
	tiny := r.State().Scenes[1].ID
	ok(t, r, "dm", "npc.add", m{"name": "Wolf", "hpMax": 3})
	ok(t, r, "dm", "npc.add", m{"name": "Ghost", "hpMax": 3, "hidden": true})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": tiny})
	wolf := lastID(r)
	for range MaxPerCell - 1 {
		ok(t, r, "dm", "token.duplicate", m{"tokenId": wolf})
	}
	failWith(t, r, "dm", "token.duplicate", m{"tokenId": wolf}, "scene has no free cell")
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": tiny})
	ok(t, r, "dm", "token.duplicate", m{"tokenId": lastID(r)}) // hidden always fits
	if c := r.st.Tokens[len(r.st.Tokens)-1]; !c.Hidden || c.Name != "Ghost #2" || c.Col != 0 || c.Row != 0 {
		t.Fatalf("hidden copy %+v", c)
	}
}

func TestTokenDuplicateByCoDM(t *testing.T) {
	r, ids := arena(t) // Beach: Wolf (visible), Ghost (hidden); co-DM "co"
	beach := r.st.CurrentSceneID
	ok(t, r, "co", "token.duplicate", m{"tokenId": ids["Wolf"]})
	ok(t, r, "co", "token.duplicate", m{"tokenId": ids["Ghost"]})
	if got := sceneNames(r, beach)[2:]; !slices.Equal(got, []string{"Wolf", "Ghost", "Wolf #2", "Ghost #2"}) {
		t.Fatalf("names %q", got)
	}
	if l := lastLog(r); l.Actor != "Cora" || l.Text.EN != "added Ghost #2 (hidden)" {
		t.Fatalf("log %+v", l)
	}
}

// A hidden target renamed by displace: DM views (attacks, the resolve/cancel log line) name the token as it is now;
// the attacker's myAttack and the player view keep what players saw.
func TestHiddenRenameDMViews(t *testing.T) {
	for _, tc := range []struct {
		name, client, typ string
		payload           func(r *Room) m
		wantLog           string
	}{
		{"resolve", "dm", "attack.resolve", func(r *Room) m { return m{"attackId": r.Attacks()[0].ID, "hit": true, "damage": 1} },
			"Aldric Vane → Goblin #3: hit, 1 damage (HP 4/5)"},
		{"cancel", "p1", "attack.cancel", func(*Room) m { return m{} }, "cancelled the attack on Goblin #3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			target := lastID(r) // "Goblin #2"
			ok(t, r, "dm", "combat.set", m{"on": true})
			ok(t, r, "p1", "attack.declare", m{"targetTokenId": target})
			ok(t, r, "dm", "token.setHidden", m{"tokenId": target, "hidden": true})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc}) // visible "Goblin #2": the hidden one becomes #3
			if a := r.Attacks(); len(a) != 1 || a[0].TargetName != "Goblin #3" || a[0].TargetTokenID != target {
				t.Fatalf("dm attacks %+v", a)
			}
			if a := r.MyAttack("p1"); a == nil || a.TargetName != "Goblin #2" {
				t.Fatalf("myAttack %+v", a)
			}
			before, feed := string(mustJSON(r.PlayerState())), len(r.st.Feed)
			ok(t, r, tc.client, tc.typ, tc.payload(r))
			if !hasLog(r, tc.wantLog) {
				t.Fatalf("log missing %q: %+v", tc.wantLog, r.dmLog)
			}
			if after := string(mustJSON(r.PlayerState())); after != before || len(r.st.Feed) != feed || r.Banner() != nil {
				t.Fatalf("players learned something: feed %+v banner %+v", r.st.Feed[feed:], r.Banner())
			}
		})
	}
}

// scene.select adds character tokens: an NPC token holding a character's name is renamed (visible: next visible
// name, displacing a hidden holder of it; hidden: scene-wide next), logged for the DM.
func TestSceneSelectRenamesClashingNpcs(t *testing.T) {
	const hero = "Aldric Vane"
	for _, tc := range []struct {
		name    string
		place   []string // pool NPCs to place in Cave, in order: name, "!" prefix = hide the token after placing
		want    []string // Cave token names after scene.select (the hero last)
		renamed []string
	}{
		{"visible", []string{hero}, []string{hero + " #2", hero}, []string{"renamed Aldric Vane to Aldric Vane #2"}},
		{"case-insensitive", []string{"aldric vane"}, []string{"aldric vane #2", hero}, []string{"renamed aldric vane to aldric vane #2"}},
		{"hidden", []string{"!" + hero}, []string{hero + " #2", hero}, []string{"renamed Aldric Vane to Aldric Vane #2"}},
		{"visible takes #2 from a hidden one", []string{hero, "!" + hero}, []string{hero + " #2", hero + " #3", hero},
			[]string{"renamed Aldric Vane to Aldric Vane #2", "renamed Aldric Vane #2 to Aldric Vane #3"}},
		{"no clash", []string{"Goblin", hero + " #2"}, []string{"Goblin", hero + " #2", hero}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := playing(t)
			ok(t, r, "dm", "scene.add", m{"name": "Cave", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 10, "rows": 5}})
			cave := r.State().Scenes[1].ID
			pool := map[string]string{}
			for _, n := range tc.place {
				name, hide := strings.CutPrefix(n, "!")
				if pool[name] == "" {
					ok(t, r, "dm", "npc.add", m{"name": name, "hpMax": 5})
					pool[name] = r.npcs.List()[len(pool)].ID
				}
				ok(t, r, "dm", "npc.place", m{"npcId": pool[name], "sceneId": cave})
				if hide {
					ok(t, r, "dm", "token.setHidden", m{"tokenId": lastID(r), "hidden": true})
				}
			}
			ok(t, r, "dm", "scene.select", m{"sceneId": cave})
			if got := sceneNames(r, cave); !slices.Equal(got, tc.want) {
				t.Fatalf("names %q, want %q", got, tc.want)
			}
			if got := renames(r); !slices.Equal(got, tc.renamed) {
				t.Fatalf("renames %q, want %q", got, tc.renamed)
			}
		})
	}
}

// Numbers stop at 999999999: an event that would need a bigger one fails with nothing changed.
func TestNameNumberCap(t *testing.T) {
	const top = "Goblin #999999999"
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, r *Room, sc string) (client, typ string, payload m)
	}{
		{"npc.place", func(t *testing.T, r *Room, sc string) (string, string, m) {
			ok(t, r, "dm", "npc.add", m{"name": top, "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			return "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc}
		}},
		{"token.duplicate", func(t *testing.T, r *Room, sc string) (string, string, m) {
			ok(t, r, "dm", "npc.add", m{"name": top, "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			return "dm", "token.duplicate", m{"tokenId": lastID(r)}
		}},
		{"hidden npc.place", func(t *testing.T, r *Room, sc string) (string, string, m) {
			ok(t, r, "dm", "npc.add", m{"name": top, "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5, "hidden": true})
			return "dm", "npc.place", m{"npcId": "npc2", "sceneId": sc}
		}},
		{"displacing a hidden one", func(t *testing.T, r *Room, sc string) (string, string, m) {
			ok(t, r, "dm", "npc.add", m{"name": top, "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			ok(t, r, "dm", "token.setHidden", m{"tokenId": lastID(r), "hidden": true})
			return "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc} // visible top: the hidden top needs #1000000000
		}},
		{"reveal", func(t *testing.T, r *Room, sc string) (string, string, m) {
			ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 5, "hidden": true})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			hidden := lastID(r)
			ok(t, r, "dm", "npc.add", m{"name": top, "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": sc})
			return "dm", "token.setHidden", m{"tokenId": hidden, "hidden": false}
		}},
		{"scene.select", func(t *testing.T, r *Room, sc string) (string, string, m) {
			ok(t, r, "dm", "scene.add", m{"name": "Cave", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 10, "rows": 5}})
			cave := r.State().Scenes[1].ID
			ok(t, r, "dm", "npc.add", m{"name": "Aldric Vane", "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": cave})
			ok(t, r, "dm", "npc.add", m{"name": "Aldric Vane #999999999", "hpMax": 5})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": cave})
			return "dm", "scene.select", m{"sceneId": cave}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			client, typ, payload := tc.setup(t, r, sc)
			base := "Goblin"
			if typ == "scene.select" {
				base = "Aldric Vane"
			}
			failWith(t, r, client, typ, payload, "too many copies of "+base)
			if len(r.Fresh()) != 0 || renames(r) != nil {
				t.Fatalf("log left after a failed event: %+v", r.dmLog)
			}
		})
	}
}
