package room

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"projectdnd/internal/notes"
	"projectdnd/internal/npcs"
	"projectdnd/internal/rules"
)

// rosterSets returns the set of each roster entry, in order, deduplicated.
func rosterSets(ro []rules.RosterEntry) []string {
	var out []string
	for _, e := range ro {
		if !slices.Contains(out, e.Set) {
			out = append(out, e.Set)
		}
	}
	return out
}

func TestSetsSet(t *testing.T) {
	r, _ := lobby(t)
	ok(t, r, "dm", "seat.openCoDm", m{})
	ok(t, r, "co", "seat.sit", m{"seatId": r.State().Seats[2].ID, "name": "Cora"})
	if got := r.State().ActiveSets; !slices.Equal(got, []string{"core"}) {
		t.Fatalf("default active sets %v", got)
	}
	for _, tt := range []struct {
		name, client string
		payload      any
		want         string
	}{
		{"player", "p1", m{"sets": []string{"legend"}}, "only the DM can do that"},
		{"unseated", "stranger", m{"sets": []string{"legend"}}, "only the DM can do that"},
		{"empty list", "dm", m{"sets": []string{}}, "pick at least one set"},
		{"no payload", "dm", nil, "pick at least one set"},
		{"null sets", "dm", m{"sets": nil}, "pick at least one set"},
		{"unknown", "dm", m{"sets": []string{"core", "pirate"}}, "no such set"},
		{"empty id", "dm", m{"sets": []string{""}}, "no such set"},
		{"not a list", "dm", `{"sets":"core"}`, "malformed payload"},
	} {
		t.Run(tt.name, func(t *testing.T) { failWith(t, r, tt.client, "sets.set", tt.payload, tt.want) })
	}

	for _, tt := range []struct {
		client string
		sets   []string
		want   []string
		log    string
		actor  string
	}{
		{"dm", []string{"legend", "core", "legend"}, []string{"core", "legend"}, "set the character sets: First Company, Heroes of Legend", "Dee"},
		{"co", []string{"legend"}, []string{"legend"}, "set the character sets: Heroes of Legend", "Cora"},
		{"dm", []string{"core", "core"}, []string{"core"}, "set the character sets: First Company", "Dee"},
	} {
		ok(t, r, tt.client, "sets.set", m{"sets": tt.sets})
		if got := r.State().ActiveSets; !slices.Equal(got, tt.want) {
			t.Fatalf("%v: active %v, want %v", tt.sets, got, tt.want)
		}
		if l := lastLog(r); l.Text.EN != tt.log || l.Actor != tt.actor || l.Text.TH == "" {
			t.Fatalf("%v: log %+v", tt.sets, l)
		}
	}

	// A state handed out earlier keeps its own slice.
	before := r.State().ActiveSets
	ok(t, r, "dm", "sets.set", m{"sets": []string{"legend"}})
	if !slices.Equal(before, []string{"core"}) {
		t.Fatalf("old state changed: %v", before)
	}

	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "game.start", m{})
	failWith(t, r, "dm", "sets.set", m{"sets": []string{"core"}}, "not allowed while the room is creating")
	failWith(t, r, "p1", "sets.set", m{"sets": []string{"core"}}, "not allowed while the room is creating")
	if s := New(); apply(s, "dm", "sets.set", m{"sets": []string{"core"}}) == nil {
		t.Fatal("sets.set accepted in the menu")
	}
}

func TestPlayerStateRosterBySet(t *testing.T) {
	r, _ := lobby(t)
	all := len(rules.DefaultRoster())
	for _, tt := range []struct {
		sets []string
		n    int
	}{
		{[]string{"core"}, 10},
		{[]string{"legend"}, 12},
		{[]string{"core", "legend"}, all},
	} {
		ok(t, r, "dm", "sets.set", m{"sets": tt.sets})
		if got := r.PlayerState().Roster; len(got) != tt.n || !slices.Equal(rosterSets(got), tt.sets) {
			t.Fatalf("%v: player roster %d entries, sets %v", tt.sets, len(got), rosterSets(got))
		}
		if got := r.State().Roster; len(got) != all {
			t.Fatalf("%v: DM roster %d entries, want %d", tt.sets, len(got), all)
		}
	}
	if len(r.roster.List()) != all {
		t.Fatal("PlayerState changed the store's list")
	}
}

func TestCharacterPickBySet(t *testing.T) {
	r := creating(t) // only core is active
	failWith(t, r, "p1", "character.pick", pick("legend-barbarian"), "that set is not in play")

	r, _ = lobby(t)
	ok(t, r, "dm", "seat.add", m{})
	ok(t, r, "p2", "seat.sit", m{"seatId": r.State().Seats[2].ID, "name": "Sam"})
	ok(t, r, "dm", "sets.set", m{"sets": []string{"legend"}})
	ok(t, r, "p1", "seat.ready", m{"ready": true})
	ok(t, r, "p2", "seat.ready", m{"ready": true})
	ok(t, r, "dm", "game.start", m{})
	failWith(t, r, "p1", "character.pick", pick("barbarian"), "that set is not in play")

	r, _ = lobby(t)
	ok(t, r, "dm", "seat.add", m{})
	ok(t, r, "p2", "seat.sit", m{"seatId": r.State().Seats[2].ID, "name": "Sam"})
	ok(t, r, "dm", "seat.add", m{})
	ok(t, r, "p3", "seat.sit", m{"seatId": r.State().Seats[3].ID, "name": "Lu"})
	ok(t, r, "dm", "sets.set", m{"sets": []string{"core", "legend"}})
	for _, p := range []string{"p1", "p2", "p3"} {
		ok(t, r, p, "seat.ready", m{"ready": true})
	}
	ok(t, r, "dm", "game.start", m{})
	ok(t, r, "p1", "character.pick", pick("legend-barbarian"))
	failWith(t, r, "p2", "character.pick", pick("legend-barbarian"), "that class is already taken")
	ok(t, r, "p2", "character.pick", pick("barbarian")) // same class, other set
	cs := r.State().Characters
	if len(cs) != 2 || cs[0].HeroID != "legend-barbarian" || cs[1].HeroID != "barbarian" || cs[1].ClassID != "barbarian" || cs[0].Name.EN == cs[1].Name.EN {
		t.Fatalf("characters %+v", cs)
	}
	ok(t, r, "p3", "character.pick", pick("legend-wizard"))
	if r.State().Status != StatusPlaying {
		t.Fatal("not playing after every pick")
	}
}

// TestHubSets: the DM connection gets every set's roster, players only the active sets; a roster edit of an
// inactive set sends players nothing.
func TestHubSets(t *testing.T) {
	h := NewHub(slog.New(slog.DiscardHandler), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	dm, pl := dial(t, srv), dial(t, srv)
	dm.send("hello", m{"clientId": "dm-secret-0001"})
	dm.next()
	pl.send("hello", m{"clientId": "pl-secret-0002"})
	pl.next()

	roomOf := func(p map[string]any) (sets []string, active []any, n int) {
		t.Helper()
		room := p["room"].(map[string]any)
		for _, e := range room["roster"].([]any) {
			if s := e.(map[string]any)["set"].(string); !slices.Contains(sets, s) {
				sets = append(sets, s)
			}
		}
		if len(room["sets"].([]any)) != 2 {
			t.Fatalf("sets %v", room["sets"])
		}
		return sets, room["activeSets"].([]any), len(room["roster"].([]any))
	}
	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	_, p, _ := dm.next()
	if sets, active, n := roomOf(p); n != 22 || !slices.Equal(sets, []string{"core", "legend"}) || len(active) != 1 || active[0] != "core" {
		t.Fatalf("dm view: %v %v %d", sets, active, n)
	}
	_, p, _ = pl.next()
	if sets, _, n := roomOf(p); n != 10 || !slices.Equal(sets, []string{"core"}) {
		t.Fatalf("player view: %v %d", sets, n)
	}

	e := rules.DefaultRoster()[10]
	e.HeroName.EN = "Edited Legend"
	dm.send("roster.update", m{"classId": e.ID, "entry": e})
	dm.next()
	dm.send("sets.set", m{"sets": []string{"legend"}})
	dm.next()
	_, p, raw := pl.next() // the roster edit sent the player nothing: this is the sets.set state
	if sets, active, n := roomOf(p); n != 12 || !slices.Equal(sets, []string{"legend"}) || len(active) != 1 || active[0] != "legend" {
		t.Fatalf("player view after sets.set: %v %v %d", sets, active, n)
	}
	if !slices.ContainsFunc(p["room"].(map[string]any)["roster"].([]any), func(x any) bool {
		return x.(map[string]any)["heroName"].(map[string]any)["en"] == "Edited Legend"
	}) {
		t.Fatalf("player view lacks the edit: %.200s", raw)
	}
}
