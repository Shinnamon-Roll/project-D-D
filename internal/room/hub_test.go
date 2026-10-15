package room

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"

	"projectdnd/internal/notes"
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

// next reads one message and returns its type, payload and raw text.
func (c *wsClient) next() (string, map[string]any, string) {
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
	h := NewHub(slog.New(slog.NewTextHandler(&logs, nil)), &notes.Store{})
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
	h := NewHub(slog.New(slog.NewTextHandler(&logs, nil)), &notes.Store{})
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
	h := NewHub(slog.New(slog.DiscardHandler), &notes.Store{})
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
