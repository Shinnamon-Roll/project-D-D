package room

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"

	"projectdnd/internal/notes"
	"projectdnd/internal/npcs"
	"projectdnd/internal/rules"
)

type wsClient struct {
	t  *testing.T
	ws *websocket.Conn
}

func dial(t *testing.T, srv *httptest.Server) *wsClient {
	t.Helper()
	ws, _, err := websocket.Dial(context.Background(), "ws"+strings.TrimPrefix(srv.URL, "http"), nil)
	if err != nil {
		t.Fatal(err)
	}
	ws.SetReadLimit(1 << 20) // the room state with every roster set is over the 32 KiB default
	t.Cleanup(func() { ws.CloseNow() })
	return &wsClient{t, ws}
}

func (c *wsClient) send(typ string, payload any) {
	c.t.Helper()
	p, _ := json.Marshal(payload)
	b, _ := json.Marshal(message{Type: typ, Payload: p})
	if err := c.ws.Write(context.Background(), websocket.MessageText, b); err != nil {
		c.t.Fatal(err)
	}
}

// next reads one message, skipping DM "log" lines, and returns its type, payload and raw text.
func (c *wsClient) next() (string, map[string]any, string) {
	c.t.Helper()
	for {
		typ, p, raw := c.nextAny()
		if typ != "log" {
			return typ, p, raw
		}
	}
}

// nextAny reads one message of any type.
func (c *wsClient) nextAny() (string, map[string]any, string) {
	c.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, b, err := c.ws.Read(ctx)
	if err != nil {
		c.t.Fatal(err)
	}
	var m struct {
		Type    string         `json:"type"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		c.t.Fatal(err)
	}
	return m.Type, m.Payload, string(b)
}

func seatID(p map[string]any) string { return p["you"].(map[string]any)["seatId"].(string) }

// lockedBuf lets the test read logs while connection goroutines still write them.
type lockedBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuf) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuf) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

func TestHub(t *testing.T) {
	var logs lockedBuf
	h := NewHub(slog.New(slog.NewTextHandler(&logs, nil)), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	const dmID, plID = "dm-secret-0001", "pl-secret-0002"

	dm, pl := dial(t, srv), dial(t, srv)
	dm.send("seat.add", m{})
	if typ, p, _ := dm.next(); typ != "error" || p["message"] != "send hello first" {
		t.Fatalf("before hello: %s %v", typ, p)
	}
	dm.send("hello", m{"clientId": "x"})
	if typ, _, _ := dm.next(); typ != "error" {
		t.Fatal("short clientId accepted")
	}
	dm.send("hello", m{"clientId": dmID})
	if typ, p, _ := dm.next(); typ != "state" || seatID(p) != "" {
		t.Fatalf("hello: %s %v", typ, p)
	}
	pl.send("hello", m{"clientId": plID})
	if _, p, _ := pl.next(); p["build"] != buildID || len(buildID) != 16 || strings.Trim(buildID, "0123456789abcdef") != "" {
		t.Fatalf("state build %v, want 16 hex %q", p["build"], buildID)
	}
	if n := h.Clients(); n != 2 {
		t.Fatalf("Clients = %d, want 2", n)
	}

	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	_, p, raw := dm.next()
	if seatID(p) == "" {
		t.Fatal("dm has no seat")
	}
	if p["build"] != buildID {
		t.Fatalf("build changed within one process: %v", p["build"])
	}
	_, p2, raw2 := pl.next()
	if seatID(p2) != "" {
		t.Fatal("player got a seat")
	}
	for _, s := range []string{raw, raw2} {
		if strings.Contains(s, dmID) || strings.Contains(s, plID) {
			t.Fatalf("clientId leaked: %s", s)
		}
	}

	// Rejected event: error only to the sender; the next message the DM sees is a state.
	pl.send("seat.add", m{})
	if typ, _, _ := pl.next(); typ != "error" {
		t.Fatal("player seat.add not rejected")
	}
	dm.send("seat.add", m{})
	if typ, _, _ := dm.next(); typ != "state" {
		t.Fatal("dm saw something other than state")
	}
	pl.next()

	if !h.Seated(dmID) || h.Seated(plID) {
		t.Fatal("Seated")
	}

	// Reconnect with the same clientId returns the same seat.
	dmSeat := seatID(p)
	dm.ws.CloseNow()
	dm2 := dial(t, srv)
	dm2.send("hello", m{"clientId": dmID})
	if _, p, _ := dm2.next(); seatID(p) != dmSeat {
		t.Fatal("reconnect lost seat")
	}
	dm2.send("hello", m{"clientId": dmID})
	if typ, _, _ := dm2.next(); typ != "error" {
		t.Fatal("second hello accepted")
	}

	out := logs.String()
	for _, want := range []string{
		"ws connected\" seat=-",
		"event campaign.select\" by=Dee status=lobby",
		"event rejected\" type=seat.add by=- err=",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("logs missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, dmID) || strings.Contains(out, plID) {
		t.Fatalf("clientId leaked into logs:\n%s", out)
	}
}

// TestHubNotes: notes reach the DM and co-DM only; players get no "notes" key and no note-only broadcasts.
func TestHubNotes(t *testing.T) {
	var logs lockedBuf
	h := NewHub(slog.New(slog.NewTextHandler(&logs, nil)), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()

	dm, co, pl := dial(t, srv), dial(t, srv), dial(t, srv)
	dm.send("hello", m{"clientId": "dm-secret-0001"})
	co.send("hello", m{"clientId": "co-secret-0003"})
	pl.send("hello", m{"clientId": "pl-secret-0002"})
	for _, c := range []*wsClient{dm, co, pl} {
		if _, p, _ := c.next(); p["notes"] != nil {
			t.Fatal("unseated connection got notes")
		}
	}
	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	_, p, _ := dm.next()
	seats := p["room"].(map[string]any)["seats"].([]any)
	plSeat := seats[1].(map[string]any)["id"].(string)
	if n, ok := p["notes"].([]any); !ok || len(n) != 0 {
		t.Fatalf("dm notes = %#v, want []", p["notes"])
	}
	co.next()
	pl.next()
	dm.send("seat.openCoDm", m{})
	_, p, _ = dm.next()
	if p["notes"] != nil {
		t.Fatal("dm got notes again on a non-note broadcast (seat.openCoDm)")
	}
	seats = p["room"].(map[string]any)["seats"].([]any)
	coSeat := seats[len(seats)-1].(map[string]any)["id"].(string)
	co.next()
	pl.next()
	co.send("seat.sit", m{"seatId": coSeat, "name": "Cora"})
	// Sitting in the co-DM seat makes that connection a DM: it gets notes once, the others do not.
	if _, p, _ = co.next(); p["notes"] == nil {
		t.Fatal("new co-DM did not get notes")
	}
	if _, p, _ = dm.next(); p["notes"] != nil {
		t.Fatal("dm got notes on co-DM sit")
	}
	pl.next()
	pl.send("seat.sit", m{"seatId": plSeat, "name": "Pat"})
	dm.next()
	if _, p, _ = co.next(); p["notes"] != nil {
		t.Fatal("co-DM got notes again on a non-note broadcast")
	}
	pl.next()

	const body = "the-dragon-is-the-mayor"
	dm.send("note.add", m{"title": "Twist"})
	dm.send("note.update", m{"id": "n1", "title": "Twist", "body": body})
	for _, c := range []*wsClient{dm, co} {
		c.next()
		if _, p, raw := c.next(); !strings.Contains(raw, body) || len(p["notes"].([]any)) != 1 {
			t.Fatalf("DM side missing note: %s", raw)
		}
	}
	pl.send("note.add", m{"title": "hack"})
	if typ, p, _ := pl.next(); typ != "error" || p["message"] != "only the DM can do that" {
		t.Fatalf("player note.add: %s %v", typ, p)
	}
	// The player's next message is the state for seat.ready: the note events above sent it nothing.
	pl.send("seat.ready", m{"ready": true})
	_, p, raw := pl.next()
	if _, has := p["notes"]; has || strings.Contains(raw, body) || strings.Contains(raw, "Twist") {
		t.Fatalf("player saw notes: %s", raw)
	}
	for _, c := range []*wsClient{dm, co} {
		if _, p, _ := c.next(); p["notes"] != nil {
			t.Fatal("non-note broadcast carried notes (client keeps its previous notes)")
		}
	}

	// A reconnect (new connection, same clientId) gets notes in its hello reply.
	dm2 := dial(t, srv)
	dm2.send("hello", m{"clientId": "dm-secret-0001"})
	if _, p, _ := dm2.next(); len(p["notes"].([]any)) != 1 {
		t.Fatalf("dm hello reply notes = %#v", p["notes"])
	}

	out := logs.String()
	if !strings.Contains(out, "event note.update\" by=Dee") {
		t.Fatalf("missing note log:\n%s", out)
	}
	if strings.Contains(out, body) {
		t.Fatalf("note body leaked into logs:\n%s", out)
	}
}

// TestHubNotesConcurrent: DM and co-DM write the same note at once; run under -race.
func TestHubNotesConcurrent(t *testing.T) {
	h := NewHub(slog.New(slog.DiscardHandler), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	dm, co := dial(t, srv), dial(t, srv)
	dm.send("hello", m{"clientId": "dm-secret-0001"})
	co.send("hello", m{"clientId": "co-secret-0003"})
	dm.next()
	co.next()
	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	dm.next()
	co.next()
	dm.send("seat.openCoDm", m{})
	_, p, _ := dm.next()
	co.next()
	seats := p["room"].(map[string]any)["seats"].([]any)
	co.send("seat.sit", m{"seatId": seats[len(seats)-1].(map[string]any)["id"], "name": "Cora"})
	dm.next()
	co.next()
	dm.send("note.add", m{"title": "Shared"})
	dm.next()
	co.next()

	const n = 10 // stays under the 32-message send buffer per connection
	var wg sync.WaitGroup
	for _, c := range []*wsClient{dm, co} {
		wg.Go(func() {
			for i := range n {
				c.send("note.update", m{"id": "n1", "title": "Shared", "body": strings.Repeat("x", i)})
			}
		})
	}
	wg.Wait()
	for _, c := range []*wsClient{dm, co} {
		for range 2 * n {
			if typ, _, _ := c.next(); typ != "state" {
				t.Fatalf("got %s", typ)
			}
		}
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if got := h.room.notes.List(); len(got) != 1 || len(got[0].Body) != n-1 {
		t.Fatalf("last write did not win: %+v", got)
	}
}

// TestHubNpcs: the pool and hidden tokens reach the DM only; players never see a hidden NPC's name or any NPC hp.
func TestHubNpcs(t *testing.T) {
	var logs lockedBuf
	h := NewHub(slog.New(slog.NewTextHandler(&logs, nil)), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	dm, pl := dial(t, srv), dial(t, srv)
	dm.send("hello", m{"clientId": "dm-secret-0001"})
	pl.send("hello", m{"clientId": "pl-secret-0002"})
	for _, c := range []*wsClient{dm, pl} {
		if _, p, _ := c.next(); p["npcs"] != nil {
			t.Fatal("unseated connection got npcs")
		}
	}
	const secret = "Secret-Hag-Zz"
	var plFrames []string
	both := func() map[string]any {
		t.Helper()
		_, p, _ := dm.next()
		_, _, raw := pl.next()
		plFrames = append(plFrames, raw)
		return p
	}
	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	p := both()
	if n, ok := p["npcs"].([]any); !ok || len(n) != 0 {
		t.Fatalf("new DM npcs = %#v, want []", p["npcs"])
	}
	plSeat := p["room"].(map[string]any)["seats"].([]any)[1].(map[string]any)["id"]
	for _, e := range []struct {
		c   *wsClient
		typ string
		p   m
	}{
		{pl, "seat.sit", m{"seatId": plSeat, "name": "Pat"}},
		{pl, "seat.ready", m{"ready": true}},
		{dm, "game.start", m{}},
		{pl, "character.pick", pick("knight")},
		{dm, "scene.add", m{"name": "Beach", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 10, "rows": 5}}},
		{dm, "scene.select", m{"sceneId": "sc4"}},
	} {
		e.c.send(e.typ, e.p)
		if p = both(); p["npcs"] != nil {
			t.Fatalf("%s: dm got npcs on a non-npc broadcast", e.typ)
		}
	}
	// probe: the player's next message must be the error for its own bad event, i.e. nothing else was queued.
	probe := func(step string) {
		t.Helper()
		pl.send("nope", m{})
		if typ, _, raw := pl.next(); typ != "error" {
			t.Fatalf("%s: player got %s: %s", step, typ, raw)
		}
	}

	dm.send("npc.add", m{"name": secret, "hpMax": 9, "hidden": true})
	if _, p, raw := dm.next(); !strings.Contains(raw, secret) || len(p["npcs"].([]any)) != 1 {
		t.Fatalf("dm npc.add: %s", raw)
	}
	dm.send("npc.place", m{"npcId": "npc1", "sceneId": "sc4"})
	_, p, raw := dm.next()
	toks := p["room"].(map[string]any)["tokens"].([]any)
	npc := toks[1].(map[string]any)
	if npc["name"] != secret || npc["hidden"] != true || npc["hp"] != 9.0 || npc["hpMax"] != 9.0 {
		t.Fatalf("dm token: %s", raw)
	}
	dm.send("token.adjustHp", m{"tokenId": npc["id"], "delta": -2})
	dm.next()
	probe("hidden npc placed")

	dm.send("token.setHidden", m{"tokenId": npc["id"], "hidden": false})
	dm.next()
	_, p, raw = pl.next()
	plFrames = append(plFrames, raw)
	toks = p["room"].(map[string]any)["tokens"].([]any)
	if len(toks) != 2 || toks[1].(map[string]any)["name"] != secret {
		t.Fatalf("revealed token missing: %s", raw)
	}
	_, hasHP := toks[1].(map[string]any)["hp"]
	if _, hasMax := toks[1].(map[string]any)["hpMax"]; hasHP || hasMax || p["npcs"] != nil {
		t.Fatalf("player got hp or the pool: %s", raw)
	}
	dm.send("token.adjustHp", m{"tokenId": npc["id"], "delta": -1})
	if _, p, _ := dm.next(); p["room"].(map[string]any)["tokens"].([]any)[1].(map[string]any)["hp"] != 6.0 {
		t.Fatal("dm hp not 6")
	}
	probe("hp change of a visible npc")

	for _, f := range plFrames[:len(plFrames)-1] {
		if strings.Contains(f, secret) {
			t.Fatalf("hidden npc name reached the player: %s", f)
		}
	}
	dm2 := dial(t, srv)
	dm2.send("hello", m{"clientId": "dm-secret-0001"})
	if _, p, _ := dm2.next(); len(p["npcs"].([]any)) != 1 {
		t.Fatalf("dm hello reply npcs = %#v", p["npcs"])
	}
	if out := logs.String(); strings.Contains(out, secret) || !strings.Contains(out, "event npc.place\" by=Dee") {
		t.Fatalf("logs:\n%s", out)
	}
}

// TestHubPickRace: two players pick the same class at once; exactly one gets it. Run under -race.
func TestHubPickRace(t *testing.T) {
	h := NewHub(slog.New(slog.DiscardHandler), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	dm, p1, p2 := dial(t, srv), dial(t, srv), dial(t, srv)
	all := []*wsClient{dm, p1, p2}
	// do sends one event that must succeed and returns the state every client then gets.
	do := func(c *wsClient, typ string, payload any) map[string]any {
		t.Helper()
		c.send(typ, payload)
		var st map[string]any
		for _, x := range all {
			typ2, p, raw := x.next()
			if typ2 != "state" {
				t.Fatalf("%s: got %s", typ, raw)
			}
			st = p
		}
		return st
	}
	for i, c := range all {
		c.send("hello", m{"clientId": []string{"dm-secret-0001", "p1-secret-0002", "p2-secret-0003"}[i]})
		c.next()
	}
	seats := func(p map[string]any) []any { return p["room"].(map[string]any)["seats"].([]any) }
	id := func(s any) any { return s.(map[string]any)["id"] }
	st := do(dm, "campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	do(p1, "seat.sit", m{"seatId": id(seats(st)[1]), "name": "Pat"})
	st = do(dm, "seat.add", m{})
	do(p2, "seat.sit", m{"seatId": id(seats(st)[2]), "name": "Sam"})
	do(p1, "seat.ready", m{"ready": true})
	do(p2, "seat.ready", m{"ready": true})
	do(dm, "game.start", m{})

	var wg sync.WaitGroup
	for _, c := range []*wsClient{p1, p2} {
		wg.Go(func() { c.send("character.pick", m{"classId": "knight"}) })
	}
	wg.Wait()
	// Each player reads until its outcome: an error (lost) or a state holding a character on its seat (won).
	won, lost := 0, 0
	for _, c := range []*wsClient{p1, p2} {
		for {
			typ, p, raw := c.next()
			if typ == "error" {
				if p["message"] != "that class is already taken" {
					t.Fatalf("loser error: %s", raw)
				}
				lost++
				break
			}
			mine := seatID(p)
			if slices.ContainsFunc(p["room"].(map[string]any)["characters"].([]any), func(x any) bool {
				return x.(map[string]any)["seatId"] == mine
			}) {
				won++
				break
			}
		}
	}
	if won != 1 || lost != 1 {
		t.Fatalf("won %d lost %d, want 1 and 1", won, lost)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if cs := h.room.st.Characters; len(cs) != 1 || cs[0].ClassID != "knight" {
		t.Fatalf("characters %+v", cs)
	}
}

type frame struct {
	typ string
	p   map[string]any
	raw string
}

// settle returns every frame c received up to the error for a probe it sends now. The hub handles one
// connection's messages in order, so all output of events processed before the probe comes first.
func (c *wsClient) settle() []frame {
	c.t.Helper()
	c.send("probe", m{})
	var out []frame
	for {
		typ, p, raw := c.nextAny()
		if typ == "error" && p["message"] == `unknown event type "probe"` {
			return out
		}
		out = append(out, frame{typ, p, raw})
	}
}

func only(fs []frame, typ string) []frame {
	var out []frame
	for _, f := range fs {
		if f.typ == typ {
			out = append(out, f)
		}
	}
	return out
}

// TestHubCombat: attacks and the DM log reach DM/co-DM connections only; players see their own myAttack;
// every connection gets the banner once; a hidden NPC's death reaches nobody but the DMs. Run under -race.
func TestHubCombat(t *testing.T) {
	h := NewHub(slog.New(slog.DiscardHandler), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	dm, p1, p2, watcher := dial(t, srv), dial(t, srv), dial(t, srv), dial(t, srv)
	all := []*wsClient{dm, p1, p2, watcher}
	secrets := []string{"dm-secret-0001", "p1-secret-0002", "p2-secret-0003", "wa-secret-0004"}
	for i, c := range all {
		c.send("hello", m{"clientId": secrets[i]})
		if _, p, _ := c.next(); p["dmLog"] != nil || p["attacks"] != nil {
			t.Fatal("unseated hello got dmLog or attacks")
		}
	}
	players := map[*wsClient][]frame{}
	// do sends one event and returns what each connection got for it (players' frames are also kept).
	do := func(c *wsClient, typ string, payload any) map[*wsClient][]frame {
		t.Helper()
		c.send(typ, payload)
		got := map[*wsClient][]frame{c: c.settle()}
		for _, o := range all {
			if o != c {
				got[o] = o.settle()
			}
		}
		for _, o := range []*wsClient{p1, p2, watcher} {
			players[o] = append(players[o], got[o]...)
		}
		for _, f := range got[c] {
			if f.typ == "error" {
				t.Fatalf("%s: %s", typ, f.raw)
			}
		}
		return got
	}

	got := do(dm, "campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	st := only(got[dm], "state")
	if len(st) != 1 || len(only(got[dm], "log")) != 0 {
		t.Fatalf("new DM: want 1 state with the full log and no log frame, got %v", got[dm])
	}
	if l, ok := st[0].p["dmLog"].([]any); !ok || len(l) != 1 || !strings.Contains(st[0].raw, "opened the room") {
		t.Fatalf("new DM dmLog: %s", st[0].raw)
	}
	seats := st[0].p["room"].(map[string]any)["seats"].([]any)
	got = do(dm, "seat.add", m{})
	if lg := only(got[dm], "log"); len(lg) != 1 || only(got[dm], "state")[0].p["dmLog"] != nil {
		t.Fatalf("dm after seat.add: %v", got[dm])
	}
	seats = only(got[dm], "state")[0].p["room"].(map[string]any)["seats"].([]any)
	do(p1, "seat.sit", m{"seatId": seats[1].(map[string]any)["id"], "name": "Pat"})
	do(p2, "seat.sit", m{"seatId": seats[2].(map[string]any)["id"], "name": "Sam"})
	do(p1, "seat.ready", m{"ready": true})
	do(p2, "seat.ready", m{"ready": true})
	do(dm, "game.start", m{})
	do(p1, "character.pick", pick("knight"))
	do(p2, "character.pick", pick("wizard"))
	got = do(dm, "scene.add", m{"name": "Beach", "imageUrl": "/uploads/a.png", "grid": m{"cellSize": 50, "cols": 10, "rows": 5}})
	sc := only(got[dm], "state")[0].p["room"].(map[string]any)["scenes"].([]any)[0].(map[string]any)["id"]
	do(dm, "scene.select", m{"sceneId": sc})
	do(dm, "npc.add", m{"name": "Wolf", "hpMax": 5})
	do(dm, "npc.add", m{"name": "Secret-Ghost", "hpMax": 5, "hidden": true})
	do(dm, "npc.place", m{"npcId": "npc1", "sceneId": sc})
	got = do(dm, "npc.place", m{"npcId": "npc2", "sceneId": sc})
	ids := map[string]string{}
	for _, x := range only(got[dm], "state")[0].p["room"].(map[string]any)["tokens"].([]any) {
		tk := x.(map[string]any)
		ids[tk["name"].(string)] = tk["id"].(string)
	}
	do(dm, "combat.set", m{"on": true})

	got = do(p1, "attack.declare", m{"targetTokenId": ids["Wolf"]})
	my := only(got[p1], "state")[0].p
	if a, ok := my["myAttack"].(map[string]any); !ok || a["targetName"] != "Wolf" || a["id"] == "" {
		t.Fatalf("p1 myAttack: %v", my["myAttack"])
	}
	if st := only(got[p2], "state"); len(st) != 1 || st[0].p["myAttack"] != nil {
		t.Fatalf("p2 after p1's declare: %v", got[p2])
	} else if _, has := st[0].p["myAttack"]; !has {
		t.Fatal("p2 state has no myAttack key")
	}
	dst := only(got[dm], "state")[0]
	atk, _ := dst.p["attacks"].([]any)
	if len(atk) != 1 || atk[0].(map[string]any)["targetName"] != "Wolf" || atk[0].(map[string]any)["byCharacterId"] == "" {
		t.Fatalf("dm attacks: %s", dst.raw)
	}
	if _, has := dst.p["myAttack"]; has {
		t.Fatal("dm got myAttack")
	}
	if lg := only(got[dm], "log"); len(lg) != 1 || lg[0].p["actor"] != "Pat" ||
		lg[0].p["text"].(map[string]any)["en"] != "Aldric Vane attacks Wolf" {
		t.Fatalf("dm log frames: %v", got[dm])
	}
	attackID := atk[0].(map[string]any)["id"]

	got = do(dm, "attack.resolve", m{"attackId": attackID, "hit": true, "damage": 3})
	for _, c := range all {
		b := only(got[c], "banner")
		if len(b) != 1 || b[0].p["kind"] != "hit" || b[0].p["text"].(map[string]any)["en"] != "Aldric Vane → Wolf" {
			t.Fatalf("banner: %v", got[c])
		}
	}
	if st := only(got[p1], "state"); len(st) != 1 || st[0].p["myAttack"] != nil {
		t.Fatalf("p1 myAttack after resolve: %v", got[p1])
	}
	if lg := only(got[dm], "log"); len(lg) != 1 || lg[0].p["actor"] != "Dee" ||
		!strings.Contains(lg[0].raw, "hit, 3 damage (HP 2/5)") {
		t.Fatalf("dm resolve log: %v", got[dm])
	}
	if a := only(got[dm], "state")[0].p["attacks"].([]any); len(a) != 0 {
		t.Fatalf("attacks after resolve: %v", a)
	}

	// Hidden NPC knocked out: players get nothing at all; the DM gets the skull in a log frame.
	got = do(dm, "token.adjustHp", m{"tokenId": ids["Secret-Ghost"], "delta": -9})
	for _, c := range []*wsClient{p1, p2, watcher} {
		if len(got[c]) != 0 {
			t.Fatalf("player got frames for a hidden NPC's hp: %v", got[c])
		}
	}
	if lg := only(got[dm], "log"); len(lg) != 2 || !strings.Contains(lg[1].raw, "💀 Secret-Ghost falls (HP 0)") {
		t.Fatalf("dm hidden death log: %v", got[dm])
	}
	dm.send("attack.resolve", m{"attackId": attackID, "hit": false}) // resolve twice
	if fs := dm.settle(); len(fs) != 1 || fs[0].typ != "error" || fs[0].p["message"] != "no such attack" {
		t.Fatalf("second resolve: %v", fs)
	}

	// A DM reconnect gets the whole log in its hello state, and no log frames for it.
	dm2 := dial(t, srv)
	dm2.send("hello", m{"clientId": secrets[0]})
	fs := dm2.settle()
	if len(fs) != 1 || fs[0].typ != "state" {
		t.Fatalf("dm reconnect: %v", fs)
	}
	l, _ := fs[0].p["dmLog"].([]any)
	if len(l) < 20 || !strings.Contains(fs[0].raw, "Secret-Ghost falls") || fs[0].p["attacks"] == nil {
		t.Fatalf("dm reconnect dmLog: %d entries", len(l))
	}
	for _, s := range secrets {
		if strings.Contains(fs[0].raw, s) {
			t.Fatal("clientId in the DM log")
		}
	}
	for c, fs := range players {
		for _, f := range fs {
			if f.typ == "log" || f.p["attacks"] != nil || f.p["dmLog"] != nil || strings.Contains(f.raw, "Secret-Ghost") ||
				strings.Contains(f.raw, "3 damage (HP") {
				t.Fatalf("player %p got DM-only data: %s", c, f.raw)
			}
			if _, has := f.p["attacks"]; has {
				t.Fatalf("player got the attacks key: %s", f.raw)
			}
			if _, has := f.p["dmLog"]; has {
				t.Fatalf("player got the dmLog key: %s", f.raw)
			}
		}
	}
}

// TestHubPresence: a seated clientId's first connection logs "joined the room" and its last disconnection
// "left the room", as log frames to DM/co-DM connections only; extra tabs and unseated clients log nothing.
func TestHubPresence(t *testing.T) {
	h := NewHub(slog.New(slog.DiscardHandler), &notes.Store{}, &npcs.Store{}, &rules.RosterStore{}, "")
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	defer srv.Close()
	// waitLive blocks until the hub counts n open connections for id (disconnects are handled asynchronously).
	waitLive := func(id string, n int) {
		t.Helper()
		for range 200 {
			h.mu.Lock()
			got := h.live[id]
			h.mu.Unlock()
			if got == n {
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
		t.Fatalf("live[%s] never reached %d", id, n)
	}
	presence := func(fs []frame) []string {
		var out []string
		for _, f := range only(fs, "log") {
			if en := f.p["text"].(map[string]any)["en"]; en == "joined the room" || en == "left the room" {
				out = append(out, f.p["actor"].(string)+" "+en.(string))
			}
		}
		return out
	}
	expect := func(c *wsClient, want ...string) {
		t.Helper()
		if got := presence(c.settle()); !slices.Equal(got, want) {
			t.Fatalf("presence frames %q, want %q", got, want)
		}
	}

	dm := dial(t, srv)
	dm.send("hello", m{"clientId": "dm-secret-0001"})
	dm.settle()
	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	dm.send("seat.add", m{})
	fs := dm.settle()
	seats := only(fs, "state")[1].p["room"].(map[string]any)["seats"].([]any)

	p1 := dial(t, srv)
	p1.send("hello", m{"clientId": "p1-secret-0002"})
	p1.send("seat.sit", m{"seatId": seats[1].(map[string]any)["id"], "name": "Pat"})
	p1.settle()
	expect(dm) // unseated at hello: nothing

	watcher := dial(t, srv) // never seated
	watcher.send("hello", m{"clientId": "wa-secret-0004"})
	watcher.settle()
	watcher.ws.CloseNow()
	waitLive("wa-secret-0004", 0)
	expect(dm)

	p1.ws.CloseNow() // last connection of Pat
	waitLive("p1-secret-0002", 0)
	expect(dm, "Pat left the room")

	a, b := dial(t, srv), dial(t, srv) // Pat back in two tabs: one join line
	a.send("hello", m{"clientId": "p1-secret-0002"})
	a.settle()
	b.send("hello", m{"clientId": "p1-secret-0002"})
	b.settle()
	expect(dm, "Pat joined the room")
	a.ws.CloseNow()
	waitLive("p1-secret-0002", 1)
	expect(dm)
	for _, f := range b.settle() {
		if f.typ == "log" {
			t.Fatalf("player got a log frame: %s", f.raw)
		}
	}

	dm2 := dial(t, srv) // DM second tab: no join line, nothing extra for the first tab
	dm2.send("hello", m{"clientId": "dm-secret-0001"})
	if fs := dm2.settle(); len(fs) != 1 || fs[0].typ != "state" {
		t.Fatalf("dm second tab: %v", fs)
	}
	expect(dm)

	b.ws.CloseNow()
	waitLive("p1-secret-0002", 0)
	expect(dm, "Pat left the room")
	expect(dm2, "Pat left the room")

	// Both lines are in the full log a new DM connection gets, stamped with Pat's seat name, never a clientId.
	dm3 := dial(t, srv)
	dm3.send("hello", m{"clientId": "dm-secret-0001"})
	st := dm3.settle()[0]
	if n := strings.Count(st.raw, `"actor":"Pat","text":{"th":"ออกจากห้อง","en":"left the room"}`); n != 2 ||
		!strings.Contains(st.raw, `"en":"joined the room"`) || strings.Contains(st.raw, "-secret-") {
		t.Fatalf("dmLog presence lines: %s", st.raw)
	}
}
