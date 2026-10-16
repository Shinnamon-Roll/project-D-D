package room

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	"projectdnd/internal/notes"
	"projectdnd/internal/npcs"
	"projectdnd/internal/rules"
)

const (
	readLimit    = 128 << 10 // note.update: 20000 runes x 4 UTF-8 bytes + JSON escapes
	sendBuffer   = 32
	writeTimeout = 10 * time.Second
)

// buildID is random per server process; clients reload the page when it changes (protocol.md "build in state").
var buildID = func() string {
	b := make([]byte, 8)
	rand.Read(b) // never fails (crypto/rand docs)
	return hex.EncodeToString(b)
}()

// Hub owns the room and every WebSocket connection. One mutex serializes Apply and broadcasts.
// shortcut: full-state broadcast on every change, fine for <=8 clients; switch to patches if it lags.
type Hub struct {
	mu    sync.Mutex
	room  *Room
	conns map[*conn]struct{}
	log   *slog.Logger
}

type conn struct {
	clientID string // "" until hello
	send     chan []byte
	dm       bool // DM/co-DM when its last state was sent; a false->true change resends notes and npcs
}

type message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewHub logs connections and events to log; it never logs clientIds (bearer secrets), note text or NPC names.
// n holds the DM notes (see notes.Load), p the NPC pool (see npcs.Load), ro the roster (see rules.LoadRoster);
// uploads is the upload directory, used to check that NPC and roster images exist.
func NewHub(log *slog.Logger, n *notes.Store, p *npcs.Store, ro *rules.RosterStore, uploads string) *Hub {
	r := New()
	r.notes, r.npcs, r.roster, r.uploads = n, p, ro, uploads
	return &Hub{room: r, conns: map[*conn]struct{}{}, log: log}
}

// Clients returns the number of open WebSocket connections.
func (h *Hub) Clients() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.conns)
}

// seatNameLocked names the sender for logs: seat name, else seat id, else "-".
func (h *Hub) seatNameLocked(clientID string) string {
	s := h.room.seatOf(clientID)
	switch {
	case s == nil:
		return "-"
	case s.Name != "":
		return s.Name
	}
	return s.ID
}

// Summary returns a status snapshot of the room (names only, no clientIds).
func (h *Hub) Summary() Summary {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.room.Summary()
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
	h.log.Info("ws connected", "seat", "-")
	defer func() {
		h.drop(c)
		h.mu.Lock()
		name := h.seatNameLocked(c.clientID)
		h.mu.Unlock()
		h.log.Info("ws disconnected", "seat", name)
	}()

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
		h.sendStateLocked(c, h.viewsLocked(), true, true)
		return
	}
	if c.clientID == "" {
		h.errorLocked(c, "send hello first")
		return
	}
	by, status := h.seatNameLocked(c.clientID), h.room.st.Status
	before := mustJSON(h.room.PlayerState())
	if err := h.room.Apply(c.clientID, m.Type, m.Payload); err != nil {
		typ := m.Type
		if r := []rune(typ); len(r) > 40 { // unknown types are client input; keep log lines short (cut on runes, not bytes)
			typ = string(r[:40])
		}
		h.log.Warn("event rejected", "type", typ, "by", by, "err", err.Error())
		h.errorLocked(c, err.Error())
		return
	}
	attrs := []any{"by", h.seatNameLocked(c.clientID)}
	if h.room.st.Status != status {
		attrs = append(attrs, "status", h.room.st.Status)
	}
	h.log.Info("event "+m.Type, attrs...)
	v := h.viewsLocked()
	// Players get a state only when their view changed (or they sent the event), so note.*, npc.* and edits
	// of hidden tokens send them nothing.
	playerChanged := !bytes.Equal(before, v.player)
	withNotes, withNpcs := strings.HasPrefix(m.Type, "note."), strings.HasPrefix(m.Type, "npc.")
	for o := range h.conns {
		if o.clientID != "" && (o == c || playerChanged || isDM(h.room.seatOf(o.clientID))) {
			h.sendStateLocked(o, v, withNotes, withNpcs)
		}
	}
}

func (h *Hub) reply(c *conn, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.errorLocked(c, msg)
}

// views is the room JSON for DM/co-DM connections (full) and for everyone else (player: Room.PlayerState).
type views struct{ full, player json.RawMessage }

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err) // state, notes and npcs hold only plain types
	}
	return b
}

func (h *Hub) viewsLocked() views {
	return views{mustJSON(h.room.State()), mustJSON(h.room.PlayerState())}
}

// sendStateLocked sends the full room to a DM or co-DM and the player view to anyone else. The "notes" and
// "npcs" keys go only to a DM or co-DM, and only when withNotes/withNpcs (hello, note.*/npc.* events) or the
// connection has just become DM/co-DM. Without a key the client keeps its previous notes/npcs.
func (h *Hub) sendStateLocked(c *conn, v views, withNotes, withNpcs bool) {
	type you struct {
		SeatID string `json:"seatId"`
	}
	dm := isDM(h.room.seatOf(c.clientID))
	room := v.player
	var notes, npcs json.RawMessage // nil: omitempty drops the key
	if dm {
		room = v.full
		if withNotes || !c.dm {
			notes = mustJSON(h.room.notes.List())
		}
		if withNpcs || !c.dm {
			npcs = mustJSON(h.room.npcs.List())
		}
	}
	c.dm = dm
	p, _ := json.Marshal(struct {
		Room  json.RawMessage `json:"room"`
		You   you             `json:"you"`
		Build string          `json:"build"`
		Notes json.RawMessage `json:"notes,omitempty"`
		Npcs  json.RawMessage `json:"npcs,omitempty"`
	}{room, you{h.room.SeatID(c.clientID)}, buildID, notes, npcs})
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
