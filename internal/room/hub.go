package room

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	readLimit    = 64 << 10
	sendBuffer   = 32
	writeTimeout = 10 * time.Second
)

// Hub owns the room and every WebSocket connection. One mutex serializes Apply and broadcasts.
// shortcut: full-state broadcast on every change, fine for <=8 clients; switch to patches if it lags.
type Hub struct {
	mu    sync.Mutex
	room  *Room
	conns map[*conn]struct{}
}

type conn struct {
	clientID string // "" until hello
	send     chan []byte
}

type message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewHub() *Hub {
	return &Hub{room: New(), conns: map[*conn]struct{}{}}
}

// Seated reports whether clientID holds a seat (used by the upload endpoint).
func (h *Hub) Seated(clientID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.room.SeatID(clientID) != ""
}

// validClientID accepts a browser-generated UUID-like token.
func validClientID(id string) bool {
	if len(id) < 8 || len(id) > 64 {
		return false
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

// ServeWS handles GET /ws. The first message must be hello{clientId}.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil) // default same-origin check
	if err != nil {
		return
	}
	defer ws.CloseNow()
	ws.SetReadLimit(readLimit)
	ctx := r.Context()

	c := &conn{send: make(chan []byte, sendBuffer)}
	h.mu.Lock()
	h.conns[c] = struct{}{}
	h.mu.Unlock()
	defer h.drop(c)

	go func() {
		for b := range c.send {
			wctx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := ws.Write(wctx, websocket.MessageText, b)
			cancel()
			if err != nil {
				ws.CloseNow()
				return
			}
		}
		ws.CloseNow() // dropped as a slow consumer
	}()

	for {
		_, data, err := ws.Read(ctx)
		if err != nil {
			return
		}
		var m message
		if err := json.Unmarshal(data, &m); err != nil {
			h.reply(c, "malformed message")
			continue
		}
		h.handle(c, m)
	}
}

func (h *Hub) handle(c *conn, m message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m.Type == "hello" {
		if c.clientID != "" {
			h.errorLocked(c, "already said hello")
			return
		}
		var p struct {
			ClientID string `json:"clientId"`
		}
		if json.Unmarshal(m.Payload, &p) != nil || !validClientID(p.ClientID) {
			h.errorLocked(c, "hello needs a clientId")
			return
		}
		c.clientID = p.ClientID
		h.sendStateLocked(c, h.roomJSONLocked())
		return
	}
	if c.clientID == "" {
		h.errorLocked(c, "send hello first")
		return
	}
	if err := h.room.Apply(c.clientID, m.Type, m.Payload); err != nil {
		h.errorLocked(c, err.Error())
		return
	}
	room := h.roomJSONLocked()
	for o := range h.conns {
		if o.clientID != "" {
			h.sendStateLocked(o, room)
		}
	}
}

func (h *Hub) reply(c *conn, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.errorLocked(c, msg)
}

func (h *Hub) roomJSONLocked() json.RawMessage {
	b, err := json.Marshal(h.room.State())
	if err != nil {
		panic(err) // State holds only plain types
	}
	return b
}

func (h *Hub) sendStateLocked(c *conn, room json.RawMessage) {
	type you struct {
		SeatID string `json:"seatId"`
	}
	p, _ := json.Marshal(struct {
		Room json.RawMessage `json:"room"`
		You  you             `json:"you"`
	}{room, you{h.room.SeatID(c.clientID)}})
	h.enqueueLocked(c, "state", p)
}

func (h *Hub) errorLocked(c *conn, msg string) {
	p, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{msg})
	h.enqueueLocked(c, "error", p)
}

func (h *Hub) enqueueLocked(c *conn, typ string, payload json.RawMessage) {
	if _, ok := h.conns[c]; !ok {
		return
	}
	b, _ := json.Marshal(message{Type: typ, Payload: payload})
	select {
	case c.send <- b:
	default: // client too slow: drop it; it reconnects and gets a fresh state
		delete(h.conns, c)
		close(c.send)
	}
}

// drop removes a closed connection. The seat stays bound to the clientId for reconnects.
func (h *Hub) drop(c *conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.conns[c]; ok {
		delete(h.conns, c)
		close(c.send)
	}
}
