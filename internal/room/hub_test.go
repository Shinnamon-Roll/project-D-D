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
	h := NewHub(slog.New(slog.NewTextHandler(&logs, nil)))
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
