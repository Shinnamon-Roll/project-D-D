package room

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestHub(t *testing.T) {
	h := NewHub()
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
	pl.next()

	dm.send("campaign.select", m{"campaignId": "starter-set", "name": "Dee"})
	_, p, raw := dm.next()
	if seatID(p) == "" {
		t.Fatal("dm has no seat")
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
}
