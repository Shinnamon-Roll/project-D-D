package room

import "testing"

const cellFull = "that cell is full (max 4)"

// withPool adds a visible pool NPC (npc1) and a hidden one (npc2).
func withPool(t *testing.T, r *Room) {
	t.Helper()
	ok(t, r, "dm", "npc.add", m{"name": "Wolf", "hpMax": 5})
	ok(t, r, "dm", "npc.add", m{"name": "Ghost", "hpMax": 5, "hidden": true})
}

// put places pool NPC npcID on scene sc, moves it to col,row and returns its token id.
func put(t *testing.T, r *Room, sc, npcID string, col, row int) string {
	t.Helper()
	ok(t, r, "dm", "npc.place", m{"npcId": npcID, "sceneId": sc})
	id := r.State().Tokens[len(r.State().Tokens)-1].ID
	ok(t, r, "dm", "token.move", m{"tokenId": id, "col": col, "row": row})
	return id
}

func TestTokenMoveCellLimit(t *testing.T) {
	for _, tc := range []struct {
		name            string
		visible, hidden int    // tokens already on 2,2 (besides the mover)
		mover           string // "pc" (Mira, moved by p1), "npc" or "ghost" (moved by the DM)
		inCell          bool   // mover already stands on 2,2
		wantErr         bool
	}{
		{"fourth visible fits", 3, 0, "pc", false, false},
		{"fifth visible rejected", 4, 0, "pc", false, true},
		{"fifth visible rejected for the DM too", 4, 0, "npc", false, true},
		{"hidden tokens do not count", 3, 3, "pc", false, false},
		{"full cell with hidden tokens: same error", 4, 2, "pc", false, true},
		{"mover is not counted", 3, 0, "pc", true, false},
		{"hidden mover ignores the limit", 4, 0, "ghost", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			withPool(t, r)
			for range tc.visible {
				put(t, r, sc, "npc1", 2, 2)
			}
			for range tc.hidden {
				put(t, r, sc, "npc2", 2, 2)
			}
			who, id := "p1", r.State().Tokens[0].ID
			switch tc.mover {
			case "npc":
				who, id = "dm", put(t, r, sc, "npc1", 5, 0)
			case "ghost":
				who, id = "dm", put(t, r, sc, "npc2", 5, 0)
			}
			if tc.inCell {
				r.token(id).Col, r.token(id).Row = 2, 2 // 4 visible already there, so token.move would refuse
			}
			move := m{"tokenId": id, "col": 2, "row": 2}
			if tc.wantErr {
				failWith(t, r, who, "token.move", move, cellFull)
				return
			}
			ok(t, r, who, "token.move", move)
			if tk := r.token(id); tk.Col != 2 || tk.Row != 2 {
				t.Fatalf("token at %d,%d", tk.Col, tk.Row)
			}
		})
	}
}

func TestSetHiddenRevealCellLimit(t *testing.T) {
	for _, tc := range []struct {
		name    string
		visible int // visible tokens on 2,2 besides the hidden one
		wantErr bool
	}{
		{"reveal into a cell with room", 3, false},
		{"reveal into a full cell", 4, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			withPool(t, r)
			var last string
			for range tc.visible {
				last = put(t, r, sc, "npc1", 2, 2)
			}
			ghost := put(t, r, sc, "npc2", 2, 2)
			reveal := m{"tokenId": ghost, "hidden": false}
			if !tc.wantErr {
				ok(t, r, "dm", "token.setHidden", reveal)
				return
			}
			failWith(t, r, "dm", "token.setHidden", reveal, cellFull)
			ok(t, r, "dm", "token.setHidden", m{"tokenId": last, "hidden": true}) // hiding is always allowed
			ok(t, r, "dm", "token.setHidden", reveal)
		})
	}
}

func TestNpcPlaceSkipsFullCells(t *testing.T) {
	for _, tc := range []struct {
		name            string
		visible, hidden int // tokens already on the 1x1 scene
		npc             string
		wantErr         bool
	}{
		{"visible fits", 3, 2, "npc1", false},
		{"visible into a full cell", 4, 0, "npc1", true},
		{"hidden into a full cell", 4, 0, "npc2", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := playing(t)
			withPool(t, r)
			ok(t, r, "dm", "scene.add", m{"name": "Cell", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 1, "rows": 1}})
			sc := r.State().Scenes[1].ID
			for range tc.visible {
				put(t, r, sc, "npc1", 0, 0)
			}
			for range tc.hidden {
				put(t, r, sc, "npc2", 0, 0)
			}
			p := m{"npcId": tc.npc, "sceneId": sc}
			if tc.wantErr {
				failWith(t, r, "dm", "npc.place", p, "scene has no free cell")
				return
			}
			ok(t, r, "dm", "npc.place", p)
		})
	}
}

func TestSceneSelectSkipsFullCells(t *testing.T) {
	type pos struct{ col, row int }
	for _, tc := range []struct {
		name    string
		fill    []pos  // one visible NPC per entry on the 2x1 scene, before select
		hidden  int    // hidden NPCs on 0,0
		want    pos    // where Mira lands
		wantErr string // non-empty: select refused
	}{
		{"empty", nil, 0, pos{0, 0}, ""},
		{"first cell full", []pos{{0, 0}, {0, 0}, {0, 0}, {0, 0}}, 0, pos{1, 0}, ""},
		{"hidden do not fill", nil, 4, pos{0, 0}, ""},
		{"every cell full", []pos{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {1, 0}, {1, 0}, {1, 0}, {1, 0}}, 0, pos{}, "scene has no free cell"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := playing(t)
			withPool(t, r)
			ok(t, r, "dm", "scene.add", m{"name": "Hall", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 2, "rows": 1}})
			sc := r.State().Scenes[1].ID
			for _, f := range tc.fill {
				put(t, r, sc, "npc1", f.col, f.row)
			}
			for range tc.hidden {
				put(t, r, sc, "npc2", 0, 0)
			}
			if tc.wantErr != "" {
				failWith(t, r, "dm", "scene.select", m{"sceneId": sc}, tc.wantErr)
				return
			}
			ok(t, r, "dm", "scene.select", m{"sceneId": sc})
			tk := r.State().Tokens[len(r.State().Tokens)-1]
			if tk.CharacterID == "" || tk.SceneID != sc || tk.Col != tc.want.col || tk.Row != tc.want.row {
				t.Fatalf("mira = %+v, want %+v", tk, tc.want)
			}
		})
	}
}

func TestSceneShrinkRelocatesFromFullCell(t *testing.T) {
	type pos struct{ col, row int }
	for _, tc := range []struct {
		name    string
		visible int  // visible NPCs on 2,2 (stay in bounds)
		hidden  bool // the clamped token is hidden
		want    pos  // where the token from 9,2 ends up after shrinking to 3x5
	}{
		{"room left: plain clamp", 3, false, pos{2, 2}},
		{"full: nearest cell with room", 4, false, pos{1, 1}},
		{"hidden: plain clamp into a full cell", 4, true, pos{2, 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, sc := playing(t)
			withPool(t, r)
			stay := make([]string, 0, tc.visible)
			for range tc.visible {
				stay = append(stay, put(t, r, sc, "npc1", 2, 2))
			}
			npc := "npc1"
			if tc.hidden {
				npc = "npc2"
			}
			id := put(t, r, sc, npc, 9, 2)
			ok(t, r, "dm", "scene.update", m{"sceneId": sc, "name": "Beach", "grid": m{"cellSize": 50, "cols": 3, "rows": 5}})
			if tk := r.token(id); tk.Col != tc.want.col || tk.Row != tc.want.row {
				t.Fatalf("clamped token at %d,%d, want %+v", tk.Col, tk.Row, tc.want)
			}
			for _, s := range stay {
				if tk := r.token(s); tk.Col != 2 || tk.Row != 2 {
					t.Fatalf("in-bounds token moved to %d,%d", tk.Col, tk.Row)
				}
			}
			if n := r.visibleAt(sc, 2, 2, ""); n > MaxPerCell {
				t.Fatalf("%d visible tokens on 2,2", n)
			}
		})
	}
}
