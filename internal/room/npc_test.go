package room

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// withUpload points the room at a temp uploads dir holding one generated-name image; returns its URL.
func withUpload(t *testing.T, r *Room) string {
	t.Helper()
	r.uploads = t.TempDir()
	name := strings.Repeat("a", 32) + ".png"
	if err := os.WriteFile(filepath.Join(r.uploads, name), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(r.uploads, strings.Repeat("d", 32)+".png"), 0o700); err != nil {
		t.Fatal(err)
	}
	return "/uploads/" + name
}

func TestNpcPoolPermissions(t *testing.T) {
	r, _ := lobby(t)
	ok(t, r, "dm", "seat.openCoDm", m{})
	ok(t, r, "co", "seat.sit", m{"seatId": r.State().Seats[len(r.State().Seats)-1].ID, "name": "Cora"})

	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 7, "hidden": true}) // lobby: any status
	ok(t, r, "co", "npc.update", m{"id": "npc1", "name": "Goblin", "hpMax": 8, "hidden": true})
	for _, who := range []string{"p1", "stranger"} {
		for _, e := range []struct {
			typ string
			p   m
		}{
			{"npc.add", m{"name": "x", "hpMax": 1}},
			{"npc.update", m{"id": "npc1", "name": "x", "hpMax": 1}},
			{"npc.remove", m{"id": "npc1"}},
		} {
			if err := apply(r, who, e.typ, e.p); err != errNotDM {
				t.Fatalf("%s %s: %v", who, e.typ, err)
			}
		}
	}
	if l := r.npcs.List(); len(l) != 1 || l[0].HPMax != 8 {
		t.Fatalf("pool changed by a non-DM: %+v", l)
	}
	if b, _ := json.Marshal(r.State()); strings.Contains(string(b), "Goblin") {
		t.Fatalf("pool leaked into RoomState: %s", b)
	}
	fail(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": "sc1"}) // not playing
	ok(t, r, "co", "npc.remove", m{"id": "npc1"})
	if len(r.npcs.List()) != 0 {
		t.Fatal("not removed")
	}
}

func TestNpcImageURL(t *testing.T) {
	r, _ := lobby(t)
	good := withUpload(t, r)
	for i, tc := range []struct {
		url string
		ok  bool
	}{
		{"", true},
		{good, true},
		{"/uploads/" + strings.Repeat("b", 32) + ".png", false}, // generated name, no file
		{"/uploads/" + strings.Repeat("d", 32) + ".png", false}, // a directory
		{"/uploads/a.png", false},                               // not a generated name
		{"/uploads/../" + strings.Repeat("a", 32) + ".png", false},
		{"/uploads/" + strings.Repeat("a", 32) + ".png?x", false},
		{"https://evil.example/" + strings.Repeat("a", 32) + ".png", false},
		{strings.Repeat("a", 32) + ".png", false},
	} {
		for _, typ := range []string{"npc.add", "npc.update"} {
			name := "Gob" // npc.update of npc1 keeps its own name; every npc.add needs a new one
			if typ == "npc.add" {
				name = "Gob " + strconv.Itoa(i)
			}
			err := apply(r, "dm", typ, m{"id": "npc1", "name": name, "hpMax": 3, "imageUrl": tc.url})
			if (err == nil) != tc.ok {
				t.Errorf("%s %q: err %v, want ok=%v", typ, tc.url, err, tc.ok)
			}
		}
	}
	if err := apply(New(), "dm", "npc.add", m{"name": "x", "hpMax": 1, "imageUrl": good}); err == nil {
		t.Fatal("image accepted by a room without an uploads dir")
	}
}

func TestNpcPlace(t *testing.T) {
	r, sc := playing(t) // Mira at 0,0 on a 10x5 grid
	img := withUpload(t, r)
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 7, "hidden": true, "imageUrl": img})
	ok(t, r, "dm", "npc.add", m{"name": "Wolf", "hpMax": 11})
	for range 3 {
		ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	}
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": sc})
	toks := r.State().Tokens[1:]
	for i, want := range []struct {
		name     string
		col, hp  int
		hidden   bool
		npc, img string
	}{
		{"Goblin", 1, 7, true, "npc1", img},
		{"Goblin #2", 2, 7, true, "npc1", img},
		{"Goblin #3", 3, 7, true, "npc1", img},
		{"Wolf", 1, 11, false, "npc2", ""}, // visible: stacks on the hidden Goblin
	} {
		tk := toks[i]
		if tk.Name != want.name || tk.Col != want.col || tk.Row != 0 || *tk.HP != want.hp || *tk.HPMax != want.hp ||
			tk.Hidden != want.hidden || tk.NpcID != want.npc || tk.ImageURL != want.img || tk.CharacterID != "" {
			t.Fatalf("token %d = %+v, want %+v", i, tk, want)
		}
	}

	// A freed name is reused; a pool edit or removal leaves placed tokens alone.
	ok(t, r, "dm", "token.remove", m{"tokenId": toks[1].ID})
	ok(t, r, "dm", "npc.update", m{"id": "npc1", "name": "Orc", "hpMax": 20})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	last := r.State().Tokens[len(r.State().Tokens)-1]
	if last.Name != "Orc" || last.Col != 2 || *last.HP != 20 || last.Hidden {
		t.Fatalf("after update: %+v", last)
	}
	ok(t, r, "dm", "npc.remove", m{"id": "npc1"})
	if r.State().Tokens[1].Name != "Goblin" {
		t.Fatal("npc.remove touched a placed token")
	}
	fail(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})

	// A scene whose cells all hold MaxPerCell visible tokens is rejected.
	ok(t, r, "dm", "scene.add", m{"name": "Tiny", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 1, "rows": 2}})
	tiny := r.State().Scenes[1].ID
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": tiny})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": tiny})
	if tk := r.State().Tokens[len(r.State().Tokens)-1]; tk.Name != "Wolf #2" || tk.Col != 0 || tk.Row != 1 {
		t.Fatalf("second row: %+v", tk) // names are counted per scene
	}
	for range 2*MaxPerCell - 2 { // then the least-crowded cell, until both hold MaxPerCell
		ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": tiny})
	}
	fail(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": tiny})
}

// A visible NPC ignores hidden tokens for its cell and number (the clashing hidden one is renamed); a hidden one
// counts every token.
func TestNpcPlaceIgnoresHidden(t *testing.T) {
	for _, tc := range []struct {
		name          string
		first, second bool // hidden flag of the first and second placed Goblin
		wantName      string
		wantCol       int
		firstName     string // the first Goblin's name afterwards
	}{
		{"hidden then visible", true, false, "Goblin", 1, "Goblin #2"},
		{"hidden then hidden", true, true, "Goblin #2", 2, "Goblin"},
		{"visible then hidden", false, true, "Goblin #2", 2, "Goblin"},
		{"visible then visible", false, false, "Goblin #2", 2, "Goblin"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 7, "hidden": tc.first})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			ok(t, r, "dm", "npc.update", m{"id": "npc1", "name": "Goblin", "hpMax": 7, "hidden": tc.second})
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
			tk := r.State().Tokens[2]
			if tk.Name != tc.wantName || tk.Col != tc.wantCol || tk.Row != 0 {
				t.Fatalf("second token = %+v, want %s at col %d", tk, tc.wantName, tc.wantCol)
			}
			if n := r.State().Tokens[1].Name; n != tc.firstName {
				t.Fatalf("first token %q, want %q", n, tc.firstName)
			}
		})
	}
}

func TestTokenNpcEdit(t *testing.T) {
	r, sc := playing(t)
	ok(t, r, "dm", "npc.add", m{"name": "Goblin", "hpMax": 10})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	pc, npc := r.State().Tokens[0].ID, r.State().Tokens[1].ID

	for _, tc := range []struct {
		delta, want int
	}{{-3, 7}, {-20, 0}, {-1, 0}, {5, 5}, {999, 10}, {0, 10}} {
		ok(t, r, "dm", "token.adjustHp", m{"tokenId": npc, "delta": tc.delta})
		if hp := *r.State().Tokens[1].HP; hp != tc.want {
			t.Fatalf("delta %d: hp %d, want %d", tc.delta, hp, tc.want)
		}
	}
	ok(t, r, "dm", "token.setHidden", m{"tokenId": npc, "hidden": true})
	if !r.State().Tokens[1].Hidden {
		t.Fatal("not hidden")
	}
	for _, e := range []struct {
		who, typ string
		p        m
	}{
		{"dm", "token.adjustHp", m{"tokenId": npc, "delta": 1000}},
		{"dm", "token.adjustHp", m{"tokenId": npc, "delta": -1000}},
		{"dm", "token.adjustHp", m{"tokenId": npc, "delta": "1"}},
		{"dm", "token.adjustHp", m{"tokenId": pc, "delta": 1}},
		{"dm", "token.setHidden", m{"tokenId": pc, "hidden": true}},
		{"dm", "token.setHidden", m{"tokenId": "nope", "hidden": true}},
		{"p1", "token.setHidden", m{"tokenId": npc, "hidden": false}},
		{"p1", "token.adjustHp", m{"tokenId": npc, "delta": -1}},
		{"stranger", "token.adjustHp", m{"tokenId": npc, "delta": -1}},
	} {
		fail(t, r, e.who, e.typ, e.p)
	}

	// Players get one error for a hidden token and a missing one: no existence oracle.
	errHidden := apply(r, "p1", "token.move", m{"tokenId": npc, "col": 1, "row": 1})
	errMissing := apply(r, "p1", "token.move", m{"tokenId": "t999", "col": 1, "row": 1})
	if errHidden == nil || errMissing == nil || errHidden.Error() != errMissing.Error() {
		t.Fatalf("hidden %v vs missing %v", errHidden, errMissing)
	}
	ok(t, r, "dm", "token.move", m{"tokenId": npc, "col": 3, "row": 3})
}

func TestPlayerState(t *testing.T) {
	r, sc := playing(t)
	ok(t, r, "dm", "npc.add", m{"name": "Hidden-Hag", "hpMax": 9, "hidden": true})
	ok(t, r, "dm", "npc.add", m{"name": "Open-Ogre", "hpMax": 30})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": sc})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": sc})

	full, _ := json.Marshal(r.State().Tokens)
	pl, _ := json.Marshal(r.PlayerState().Tokens)
	for _, s := range []string{"Hidden-Hag", `"hp":9`, `"hpMax":30`} {
		if !strings.Contains(string(full), s) {
			t.Errorf("full state missing %s: %s", s, full)
		}
	}
	if b, _ := json.Marshal(r.PlayerState()); strings.Contains(string(b), "Hidden-Hag") || strings.Contains(string(b), "npcs") {
		t.Fatalf("player state leaks the pool or a hidden name: %s", b)
	}
	for _, s := range []string{`"hp"`, `"hpMax"`} {
		if strings.Contains(string(pl), s) {
			t.Errorf("player state has %s: %s", s, pl)
		}
	}
	if strings.Contains(string(pl), `"npcId":"npc`) {
		t.Errorf("player state has a pool npcId: %s", pl)
	}
	if r.State().Tokens[2].NpcID == "" {
		t.Fatal("PlayerState stripped npcId from the full state")
	}
	if !strings.Contains(string(pl), "Open-Ogre") || len(r.PlayerState().Tokens) != 2 {
		t.Fatalf("player state lost visible tokens: %s", pl)
	}
	if tk := r.State().Tokens[0]; tk.HP != nil || tk.HPMax != nil {
		t.Fatalf("character token has hp: %+v", tk)
	}
	if r.State().Tokens[2].HP == nil {
		t.Fatal("PlayerState stripped hp from the full state")
	}
}
