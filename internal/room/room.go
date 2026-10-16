// Package room holds the single server-authoritative game room. Every client event goes through
// Room.Apply, which checks the sender's role and the room status before changing state.
// Event table: docs/design/protocol.md.
package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"projectdnd/internal/notes"
	"projectdnd/internal/npcs"
	"projectdnd/internal/rules"
)

const (
	StatusMenu     = "menu"
	StatusLobby    = "lobby"
	StatusCreating = "creating"
	StatusPlaying  = "playing"

	KindDM     = "dm"
	KindCoDM   = "codm"
	KindPlayer = "player"

	MaxPlayerSeats = 6
	MaxPerCell     = 4 // visible tokens per grid cell; hidden tokens never count
	maxName        = 32
)

type Seat struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Occupied    bool   `json:"occupied"`
	Ready       bool   `json:"ready"`
	CharacterID string `json:"characterId"`
	client      string // bearer secret; unexported so it is never serialized
}

// Character is a player's pick from the roster; character.go holds its events.
type Character struct {
	ID      string             `json:"id"`
	SeatID  string             `json:"seatId"`
	ClassID string             `json:"classId"`
	Name    rules.L            `json:"name"` // heroName when picked
	Art     string             `json:"art"`
	HP      int                `json:"hp"`
	HPMax   int                `json:"hpMax"`
	MP      int                `json:"mp"`
	MPMax   int                `json:"mpMax"`
	Str     int                `json:"str"`
	End     int                `json:"end"`
	Agi     int                `json:"agi"`
	Skills  []CharSkill        `json:"skills"`
	Items   []rules.RosterItem `json:"items"`
}

type CharSkill struct {
	ID   string  `json:"id"`
	Kind string  `json:"kind"` // "innate" (from the roster) | "custom" (the player's own)
	Name rules.L `json:"name"`
	Desc rules.L `json:"desc"`
	Cost int     `json:"cost"`
}

// FeedEntry is one line of the action log every client sees. Text holds names only, never clientIds.
type FeedEntry struct {
	At   int64   `json:"at"` // unix milliseconds
	Text rules.L `json:"text"`
}

type Grid struct {
	CellSize int `json:"cellSize"`
	OffsetX  int `json:"offsetX"`
	OffsetY  int `json:"offsetY"`
	Cols     int `json:"cols"`
	Rows     int `json:"rows"`
}

type Fog struct {
	ID  string `json:"id"`
	Col int    `json:"col"`
	Row int    `json:"row"`
	W   int    `json:"w"`
	H   int    `json:"h"`
}

type Scene struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
	Grid     Grid   `json:"grid"`
	Fog      []Fog  `json:"fog"`
}

type Token struct {
	ID          string `json:"id"`
	SceneID     string `json:"sceneId"`
	CharacterID string `json:"characterId"` // "" = NPC
	NpcID       string `json:"npcId"`       // pool NPC it was placed from; "" for characters
	Name        string `json:"name"`
	ImageURL    string `json:"imageUrl"`
	Col         int    `json:"col"`
	Row         int    `json:"row"`
	Hidden      bool   `json:"hidden"`          // NPC only; players never receive hidden tokens
	HP          *int   `json:"hp,omitempty"`    // NPC only; stripped from the players' view
	HPMax       *int   `json:"hpMax,omitempty"` // NPC only; stripped from the players' view
}

// State is the RoomState sent to every client.
type State struct {
	Status         string              `json:"status"`
	Campaigns      []rules.Campaign    `json:"campaigns"`
	CampaignID     string              `json:"campaignId"`
	Seats          []Seat              `json:"seats"`
	Sets           []rules.Set         `json:"sets"`       // every character set (fixed)
	ActiveSets     []string            `json:"activeSets"` // sets players pick from, in Sets order; replaced, never mutated
	Roster         []rules.RosterEntry `json:"roster"`     // filled by State() from the roster store
	Characters     []Character         `json:"characters"`
	Scenes         []Scene             `json:"scenes"`
	CurrentSceneID string              `json:"currentSceneId"`
	Tokens         []Token             `json:"tokens"`
	Feed           []FeedEntry         `json:"feed"` // newest last, at most maxFeed
	Combat         bool                `json:"combat"`
}

// Room is not safe for concurrent use; the Hub serializes access.
type Room struct {
	st      State
	seq     int
	notes   *notes.Store // DM-only; never part of State (players must not see it)
	npcs    *npcs.Store  // DM-only pool; never part of State
	roster  *rules.RosterStore
	uploads string // directory of uploaded files; "" = no uploads exist

	attacks    []Attack         // pending, DM-only; never part of State
	dmLog      []LogEntry       // DM-only, oldest first, at most maxLog; never part of State
	fresh      []LogEntry       // log lines of the event being applied; Apply keeps them only on success
	banner     *Banner          // set by attack.resolve for the event being applied
	noteLogged map[string]int64 // note id -> unix ms of its last "edited note" log line
	now        func() time.Time
}

func New() *Room {
	r := &Room{st: State{
		Status:     StatusMenu,
		Campaigns:  rules.Campaigns(),
		Sets:       rules.Sets(),
		ActiveSets: []string{"core"},
		Seats:      []Seat{},
		Characters: []Character{},
		Scenes:     []Scene{},
		Tokens:     []Token{},
		Feed:       []FeedEntry{},
	}, notes: &notes.Store{}, npcs: &npcs.Store{}, roster: &rules.RosterStore{}, now: time.Now}
	return r
}

// State returns the full state (DM/co-DM view) for serialization. Callers must not mutate it.
func (r *Room) State() State {
	st := r.st
	st.Roster = r.roster.List()
	return st
}

// PlayerState is the view for players and unseated connections: no hidden tokens, no NPC hp/hpMax/npcId, and
// only the roster entries of the active sets (the DM view keeps every set for the Ruleset editor).
func (r *Room) PlayerState() State {
	st := r.State()
	st.Roster = slices.DeleteFunc(slices.Clone(st.Roster), func(e rules.RosterEntry) bool {
		return !slices.Contains(r.st.ActiveSets, e.Set)
	})
	st.Tokens = make([]Token, 0, len(r.st.Tokens))
	for _, t := range r.st.Tokens {
		if t.Hidden {
			continue
		}
		t.HP, t.HPMax, t.NpcID = nil, nil, ""
		st.Tokens = append(st.Tokens, t)
	}
	return st
}

// SeatID returns the seat bound to clientID, or "".
func (r *Room) SeatID(clientID string) string {
	if s := r.seatOf(clientID); s != nil {
		return s.ID
	}
	return ""
}

func (r *Room) id(prefix string) string {
	r.seq++
	return fmt.Sprintf("%s%d", prefix, r.seq)
}

func (r *Room) seatOf(clientID string) *Seat {
	if clientID == "" {
		return nil
	}
	for i := range r.st.Seats {
		if r.st.Seats[i].Occupied && r.st.Seats[i].client == clientID {
			return &r.st.Seats[i]
		}
	}
	return nil
}

func (r *Room) seat(id string) *Seat {
	for i := range r.st.Seats {
		if r.st.Seats[i].ID == id {
			return &r.st.Seats[i]
		}
	}
	return nil
}

func (r *Room) character(id string) *Character {
	for i := range r.st.Characters {
		if r.st.Characters[i].ID == id {
			return &r.st.Characters[i]
		}
	}
	return nil
}

func (r *Room) scene(id string) *Scene {
	for i := range r.st.Scenes {
		if r.st.Scenes[i].ID == id {
			return &r.st.Scenes[i]
		}
	}
	return nil
}

func (r *Room) token(id string) *Token {
	for i := range r.st.Tokens {
		if r.st.Tokens[i].ID == id {
			return &r.st.Tokens[i]
		}
	}
	return nil
}

func isDM(s *Seat) bool { return s != nil && (s.Kind == KindDM || s.Kind == KindCoDM) }

func cleanName(s string) (string, error) {
	s = strings.TrimSpace(s)
	if n := utf8.RuneCountInString(s); n < 1 || n > maxName {
		return "", fmt.Errorf("name must be 1-%d characters", maxName)
	}
	return s, nil
}

// uploadURL reports whether s points at a file under /uploads/ (as returned by POST /api/upload).
func uploadURL(s string) bool {
	name, ok := strings.CutPrefix(s, "/uploads/")
	return ok && name != "" && len(name) <= 64 && !strings.ContainsAny(name, `/\?#`) && !strings.Contains(name, "..")
}

// uploadName matches only names the upload handler generates, so no user input reaches a filesystem path.
var uploadName = regexp.MustCompile(`^[0-9a-f]{32}\.(png|jpg|gif|webp)$`)

// UploadName reports whether name is one the upload handler could have generated.
func UploadName(name string) bool { return uploadName.MatchString(name) }

// uploadedFile reports whether s is /uploads/<generated name> and that file exists.
func (r *Room) uploadedFile(s string) bool {
	name, ok := strings.CutPrefix(s, "/uploads/")
	if !ok || !UploadName(name) || r.uploads == "" {
		return false
	}
	fi, err := os.Stat(filepath.Join(r.uploads, name))
	return err == nil && fi.Mode().IsRegular()
}

func gridOK(g Grid) error {
	if g.CellSize < 10 || g.CellSize > 400 {
		return errors.New("grid cellSize must be 10-400")
	}
	if g.Cols < 1 || g.Cols > 200 || g.Rows < 1 || g.Rows > 200 {
		return errors.New("grid cols and rows must be 1-200")
	}
	return nil
}

func decode[T any](raw json.RawMessage) (T, error) {
	var v T
	if len(raw) == 0 || string(raw) == "null" {
		return v, nil
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return v, errors.New("malformed payload")
	}
	return v, nil
}

var (
	errNotSeated = errors.New("you are not seated")
	errNotDM     = errors.New("only the DM can do that")
	errCellFull  = fmt.Errorf("that cell is full (max %d)", MaxPerCell)
)

// visibleAt counts the visible tokens on one cell of a scene, not counting token skip.
// Hidden tokens never count, so a full-cell error never hints that a hidden token exists.
func (r *Room) visibleAt(sceneID string, col, row int, skip string) int {
	n := 0
	for _, t := range r.st.Tokens {
		if t.SceneID == sceneID && t.Col == col && t.Row == row && !t.Hidden && t.ID != skip {
			n++
		}
	}
	return n
}

func (r *Room) needStatus(ok ...string) error {
	if slices.Contains(ok, r.st.Status) {
		return nil
	}
	return fmt.Errorf("not allowed while the room is %s", r.st.Status)
}

func (r *Room) needDM(s *Seat, status ...string) error {
	if err := r.needStatus(status...); err != nil {
		return err
	}
	if !isDM(s) {
		return errNotDM
	}
	return nil
}

var inGame = []string{StatusCreating, StatusPlaying}

// Apply validates and applies one client event. On error the state is unchanged.
// On success the event's DM log lines are in Fresh() (already appended to DMLog) and a hit/miss in Banner().
func (r *Room) Apply(clientID, typ string, payload json.RawMessage) error {
	r.fresh, r.banner = nil, nil
	actor := r.actorName(clientID) // before the event: seat.leave clears the name
	if err := r.apply(clientID, typ, payload); err != nil {
		r.fresh, r.banner = nil, nil
		return err
	}
	if actor == "" { // campaign.select, seat.sit: named only after the event
		actor = r.actorName(clientID)
	}
	for i := range r.fresh {
		r.fresh[i].Actor = actor
	}
	r.keepLog(r.fresh...)
	return nil
}

// Fresh returns the DM log lines added by the last successful Apply. Callers must not mutate it.
func (r *Room) Fresh() []LogEntry { return r.fresh }

// Banner returns the hit/miss banner of the last successful Apply, or nil.
func (r *Room) Banner() *Banner { return r.banner }

func (r *Room) actorName(clientID string) string {
	if s := r.seatOf(clientID); s != nil {
		return s.Name
	}
	return ""
}

func (r *Room) apply(clientID, typ string, payload json.RawMessage) error {
	s := r.seatOf(clientID)
	switch typ {
	case "campaign.select":
		return r.campaignSelect(clientID, s, payload)
	case "seat.sit":
		return r.seatSit(clientID, s, payload)
	case "seat.leave":
		if err := r.needStatus(StatusLobby); err != nil {
			return err
		}
		if s == nil {
			return errNotSeated
		}
		if s.Kind == KindDM {
			return errors.New("the DM cannot leave the DM seat")
		}
		*s = Seat{ID: s.ID, Kind: s.Kind}
		r.addLog("ลุกจากที่นั่ง", "left their seat")
		return nil
	case "seat.ready":
		p, err := decode[struct {
			Ready bool `json:"ready"`
		}](payload)
		if err != nil {
			return err
		}
		if err := r.needStatus(StatusLobby); err != nil {
			return err
		}
		if s == nil || s.Kind != KindPlayer {
			return errors.New("only a seated player can mark ready")
		}
		s.Ready = p.Ready
		if p.Ready {
			r.addLog("พร้อมแล้ว", "is ready")
		} else {
			r.addLog("ยังไม่พร้อม", "is not ready")
		}
		return nil
	case "seat.add":
		if err := r.needDM(s, StatusLobby); err != nil {
			return err
		}
		if r.count(KindPlayer) >= MaxPlayerSeats {
			return fmt.Errorf("at most %d player seats", MaxPlayerSeats)
		}
		r.st.Seats = append(r.st.Seats, Seat{ID: r.id("s"), Kind: KindPlayer})
		r.addLog("เพิ่มที่นั่งผู้เล่น", "added a player seat")
		return nil
	case "seat.remove":
		p, err := decode[struct {
			SeatID string `json:"seatId"`
		}](payload)
		if err != nil {
			return err
		}
		if err := r.needDM(s, StatusLobby); err != nil {
			return err
		}
		t := r.seat(p.SeatID)
		if t == nil {
			return errors.New("no such seat")
		}
		if t.Kind == KindDM {
			return errors.New("the DM seat cannot be removed")
		}
		if t.Occupied {
			r.addLog("ลบที่นั่งของ "+t.Name, "removed "+t.Name+"'s seat")
		} else {
			r.addLog("ลบที่นั่งว่าง", "removed an empty seat")
		}
		r.st.Seats = slices.DeleteFunc(r.st.Seats, func(x Seat) bool { return x.ID == p.SeatID })
		return nil
	case "seat.openCoDm":
		if err := r.needDM(s, StatusLobby); err != nil {
			return err
		}
		if r.count(KindCoDM) >= 1 {
			return errors.New("there is already a co-DM seat")
		}
		r.st.Seats = append(r.st.Seats, Seat{ID: r.id("s"), Kind: KindCoDM})
		r.addLog("เปิดที่นั่ง co-DM", "opened a co-DM seat")
		return nil
	case "sets.set":
		return r.setsSet(s, payload)
	case "game.start":
		return r.gameStart(s)
	case "character.pick":
		return r.characterPick(s, payload)
	case "skill.addCustom", "skill.editCustom", "skill.removeCustom":
		return r.customSkill(s, typ, payload)
	case "skill.use":
		return r.skillUse(s, payload)
	case "character.adjust":
		return r.characterAdjust(s, payload)
	case "character.setStats":
		return r.characterSetStats(s, payload)
	case "party.rest":
		return r.partyRest(s)
	case "item.give":
		return r.itemGive(s, payload)
	case "item.remove":
		return r.itemRemove(s, payload)
	case "roster.update", "roster.reset":
		if !isDM(s) { // any room status, like notes
			return errNotDM
		}
		return r.rosterEdit(typ, payload)
	case "scene.add", "scene.update":
		return r.sceneUpsert(s, typ, payload)
	case "scene.select":
		return r.sceneSelect(s, payload)
	case "npc.place":
		return r.npcPlace(s, payload)
	case "token.setHidden", "token.adjustHp":
		return r.tokenNpcEdit(s, typ, payload)
	case "token.move":
		return r.tokenMove(s, payload)
	case "token.remove":
		p, err := decode[struct {
			TokenID string `json:"tokenId"`
		}](payload)
		if err != nil {
			return err
		}
		if err := r.needDM(s, StatusPlaying); err != nil {
			return err
		}
		t := r.token(p.TokenID)
		if t == nil {
			return errors.New("no such token")
		}
		r.addLog("ลบ "+t.Name+" ออกจากแมพ", "removed "+t.Name+" from the map")
		r.attacks = slices.DeleteFunc(r.attacks, func(a Attack) bool { return a.TargetTokenID == p.TokenID })
		r.st.Tokens = slices.DeleteFunc(r.st.Tokens, func(t Token) bool { return t.ID == p.TokenID })
		return nil
	case "fog.add":
		return r.fogAdd(s, payload)
	case "fog.remove":
		p, err := decode[struct {
			SceneID string `json:"sceneId"`
			FogID   string `json:"fogId"`
		}](payload)
		if err != nil {
			return err
		}
		if err := r.needDM(s, StatusPlaying); err != nil {
			return err
		}
		sc := r.scene(p.SceneID)
		if sc == nil {
			return errors.New("no such scene")
		}
		n := len(sc.Fog)
		sc.Fog = slices.DeleteFunc(sc.Fog, func(f Fog) bool { return f.ID == p.FogID })
		if len(sc.Fog) == n {
			return errors.New("no such fog")
		}
		r.addLog("ลบหมอกใน "+sc.Name, "removed fog in "+sc.Name)
		return nil
	case "note.add", "note.update", "note.remove", "note.move":
		if !isDM(s) { // any room status: the DM prepares the story before the game starts
			return errNotDM
		}
		return r.note(typ, payload)
	case "npc.add", "npc.update", "npc.remove":
		if !isDM(s) { // any room status, like notes
			return errNotDM
		}
		return r.npc(typ, payload)
	case "combat.set":
		return r.combatSet(s, payload)
	case "attack.declare":
		return r.attackDeclare(s, payload)
	case "attack.cancel":
		return r.attackCancel(s)
	case "attack.resolve":
		return r.attackResolve(s, payload)
	case "npc.act":
		return r.npcAct(s, payload)
	}
	return fmt.Errorf("unknown event type %q", typ)
}

// Summary is a status snapshot for the TUI dashboard. It holds names only, never clientIds.
type Summary struct {
	Status, Campaign string
	DM, CoDM         string // "" when the seat is missing or empty
	Players, Seats   int    // occupied / total player seats
	Ready            int    // ready players
	Characters       int
	Scenes           int
	Scene            string // current scene name
}

// Summary returns a snapshot of the room for status displays.
func (r *Room) Summary() Summary {
	s := Summary{Status: r.st.Status, Characters: len(r.st.Characters), Scenes: len(r.st.Scenes)}
	for _, c := range r.st.Campaigns {
		if c.ID == r.st.CampaignID {
			s.Campaign = c.Name
		}
	}
	for _, x := range r.st.Seats {
		switch {
		case x.Kind == KindDM && x.Occupied:
			s.DM = x.Name
		case x.Kind == KindCoDM && x.Occupied:
			s.CoDM = x.Name
		case x.Kind == KindPlayer:
			s.Seats++
			if x.Occupied {
				s.Players++
				if x.Ready {
					s.Ready++
				}
			}
		}
	}
	if sc := r.scene(r.st.CurrentSceneID); sc != nil {
		s.Scene = sc.Name
	}
	return s
}

func (r *Room) count(kind string) int {
	n := 0
	for _, s := range r.st.Seats {
		if s.Kind == kind {
			n++
		}
	}
	return n
}

func (r *Room) campaignSelect(clientID string, s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		CampaignID string `json:"campaignId"`
		Name       string `json:"name"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(StatusMenu); err != nil {
		return err
	}
	if s != nil {
		return errors.New("you are already seated")
	}
	i := slices.IndexFunc(r.st.Campaigns, func(c rules.Campaign) bool { return c.ID == p.CampaignID })
	if i < 0 || r.st.Campaigns[i].Status != "available" {
		return errors.New("campaign not available")
	}
	name, err := cleanName(p.Name)
	if err != nil {
		return err
	}
	r.st.CampaignID = p.CampaignID
	r.st.Seats = []Seat{
		{ID: r.id("s"), Kind: KindDM, Name: name, Occupied: true, client: clientID},
		{ID: r.id("s"), Kind: KindPlayer},
	}
	r.st.Status = StatusLobby
	r.addLog("เปิดห้อง: "+r.st.Campaigns[i].Name, "opened the room: "+r.st.Campaigns[i].Name)
	return nil
}

func (r *Room) seatSit(clientID string, s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		SeatID string `json:"seatId"`
		Name   string `json:"name"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(StatusLobby); err != nil {
		return err
	}
	if s != nil {
		return errors.New("you are already seated")
	}
	t := r.seat(p.SeatID)
	if t == nil {
		return errors.New("no such seat")
	}
	if t.Occupied {
		return errors.New("seat is taken")
	}
	if t.Kind != KindPlayer && t.Kind != KindCoDM {
		return errors.New("you cannot sit there")
	}
	name, err := cleanName(p.Name)
	if err != nil {
		return err
	}
	*t = Seat{ID: t.ID, Kind: t.Kind, Name: name, Occupied: true, client: clientID}
	if t.Kind == KindCoDM {
		r.addLog("นั่งที่นั่ง co-DM", "took the co-DM seat")
	} else {
		r.addLog("นั่งที่นั่งผู้เล่น", "took a player seat")
	}
	return nil
}

// setsSet is sets.set: the DM/co-DM picks which character sets players choose from, lobby only.
func (r *Room) setsSet(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		Sets []string `json:"sets"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusLobby); err != nil {
		return err
	}
	if len(p.Sets) == 0 {
		return errors.New("pick at least one set")
	}
	for _, id := range p.Sets {
		if !rules.IsSet(id) {
			return rules.ErrNoSet
		}
	}
	var active, th, en []string // Sets order, duplicates dropped
	for _, x := range r.st.Sets {
		if slices.Contains(p.Sets, x.ID) {
			active, th, en = append(active, x.ID), append(th, x.Name.TH), append(en, x.Name.EN)
		}
	}
	if slices.Equal(active, r.st.ActiveSets) {
		return nil // no change: no log line
	}
	r.st.ActiveSets = active
	r.addLog("เลือกเซ็ตตัวละคร: "+strings.Join(th, ", "), "set the character sets: "+strings.Join(en, ", "))
	return nil
}

func (r *Room) gameStart(s *Seat) error {
	if err := r.needDM(s, StatusLobby); err != nil {
		return err
	}
	seated := 0
	for _, x := range r.st.Seats {
		if x.Kind == KindPlayer && x.Occupied {
			if !x.Ready {
				return fmt.Errorf("%s is not ready", x.Name)
			}
			seated++
		}
	}
	if seated == 0 {
		return errors.New("need at least one seated player")
	}
	r.st.Seats = slices.DeleteFunc(r.st.Seats, func(x Seat) bool { return x.Kind == KindPlayer && !x.Occupied })
	r.st.Status = StatusCreating
	r.addLog("เริ่มเกม — สร้างตัวละคร", "started the game — character creation")
	return nil
}

func (r *Room) sceneUpsert(s *Seat, typ string, payload json.RawMessage) error {
	p, err := decode[struct {
		SceneID  string `json:"sceneId"`
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl"`
		Grid     Grid   `json:"grid"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, inGame...); err != nil {
		return err
	}
	name, err := cleanName(p.Name)
	if err != nil {
		return err
	}
	if err := gridOK(p.Grid); err != nil {
		return err
	}
	if typ == "scene.add" {
		if !uploadURL(p.ImageURL) {
			return errors.New("imageUrl must start with /uploads/")
		}
		r.st.Scenes = append(r.st.Scenes, Scene{ID: r.id("sc"), Name: name, ImageURL: p.ImageURL, Grid: p.Grid, Fog: []Fog{}})
		r.addLog(fmt.Sprintf("เพิ่มฉาก %s (กริด %d×%d)", name, p.Grid.Cols, p.Grid.Rows), fmt.Sprintf("added scene %s (grid %d×%d)", name, p.Grid.Cols, p.Grid.Rows))
		return nil
	}
	sc := r.scene(p.SceneID)
	if sc == nil {
		return errors.New("no such scene")
	}
	visible := 0
	for _, t := range r.st.Tokens {
		if t.SceneID == sc.ID && !t.Hidden {
			visible++
		}
	}
	if need := (visible + MaxPerCell - 1) / MaxPerCell; p.Grid.Cols*p.Grid.Rows < need {
		return fmt.Errorf("grid too small: %d visible tokens need at least %d cells", visible, need)
	}
	sc.Name, sc.Grid = name, p.Grid
	// Keep tokens on the board when the grid shrinks. A clamped visible token that would overfill its cell
	// goes to the nearest cell with room. Tokens not yet clamped sit outside the new grid, so they never count.
	for i := range r.st.Tokens {
		t := &r.st.Tokens[i]
		if t.SceneID != sc.ID || (t.Col < p.Grid.Cols && t.Row < p.Grid.Rows) {
			continue
		}
		t.Col, t.Row = min(t.Col, p.Grid.Cols-1), min(t.Row, p.Grid.Rows-1)
		if !t.Hidden && r.visibleAt(sc.ID, t.Col, t.Row, t.ID) >= MaxPerCell {
			// Always finds a cell: the grid holds at least ceil(visible/4) cells (checked above).
			best, c0, r0 := -1, t.Col, t.Row
			for c := range p.Grid.Cols * p.Grid.Rows {
				col, row := c%p.Grid.Cols, c/p.Grid.Cols
				d := max(abs(col-c0), abs(row-r0))
				if (best < 0 || d < best) && r.visibleAt(sc.ID, col, row, t.ID) < MaxPerCell {
					best, t.Col, t.Row = d, col, row
				}
			}
		}
	}
	// Clip fog to the grid; drop rects left fully outside. New slice: State() shares the old one.
	fog := make([]Fog, 0, len(sc.Fog))
	for _, f := range sc.Fog {
		if f.Col < p.Grid.Cols && f.Row < p.Grid.Rows {
			f.W, f.H = min(f.W, p.Grid.Cols-f.Col), min(f.H, p.Grid.Rows-f.Row)
			fog = append(fog, f)
		}
	}
	sc.Fog = fog
	r.addLog(fmt.Sprintf("แก้ฉาก %s (กริด %d×%d)", name, p.Grid.Cols, p.Grid.Rows), fmt.Sprintf("edited scene %s (grid %d×%d)", name, p.Grid.Cols, p.Grid.Rows))
	return nil
}

func (r *Room) sceneSelect(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		SceneID string `json:"sceneId"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	sc := r.scene(p.SceneID)
	if sc == nil {
		return errors.New("no such scene")
	}
	// One token per cell, row-major from 0,0, skipping full cells; wrap around when the grid runs out.
	n, cells, next := len(r.st.Tokens), sc.Grid.Cols*sc.Grid.Rows, 0
	for _, c := range r.st.Characters {
		if slices.ContainsFunc(r.st.Tokens, func(t Token) bool { return t.SceneID == sc.ID && t.CharacterID == c.ID }) {
			continue
		}
		cell := -1
		for k := range cells {
			if i := (next + k) % cells; r.visibleAt(sc.ID, i%sc.Grid.Cols, i/sc.Grid.Cols, "") < MaxPerCell {
				cell = i
				break
			}
		}
		if cell < 0 {
			r.st.Tokens = r.st.Tokens[:n] // state unchanged on error
			return errors.New("scene has no free cell")
		}
		r.st.Tokens = append(r.st.Tokens, Token{
			ID: r.id("t"), SceneID: sc.ID, CharacterID: c.ID, Name: c.Name.EN, ImageURL: c.Art,
			Col: cell % sc.Grid.Cols, Row: cell / sc.Grid.Cols,
		})
		next = cell + 1
	}
	r.st.CurrentSceneID = sc.ID
	r.addLog("เปลี่ยนไปฉาก "+sc.Name, "switched to scene "+sc.Name)
	return nil
}

func abs(x int) int { return max(x, -x) }

// cellText formats a grid cell for log lines: "(col,row)", 0-based like the protocol.
func cellText(col, row int) string { return fmt.Sprintf("(%d,%d)", col, row) }

func (r *Room) tokenMove(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		TokenID string `json:"tokenId"`
		Col     int    `json:"col"`
		Row     int    `json:"row"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needStatus(StatusPlaying); err != nil {
		return err
	}
	if s == nil {
		return errNotSeated
	}
	t := r.token(p.TokenID)
	// Same error for a missing and a foreign token: players must not learn that a hidden token id exists.
	if !isDM(s) && (t == nil || s.Kind != KindPlayer || t.CharacterID == "" || t.CharacterID != s.CharacterID) {
		return errors.New("you can only move your own token")
	}
	if t == nil {
		return errors.New("no such token")
	}
	sc := r.scene(t.SceneID)
	if sc == nil || p.Col < 0 || p.Row < 0 || p.Col >= sc.Grid.Cols || p.Row >= sc.Grid.Rows {
		return errors.New("target cell is outside the grid")
	}
	if !t.Hidden && r.visibleAt(t.SceneID, p.Col, p.Row, t.ID) >= MaxPerCell {
		return errCellFull
	}
	from := cellText(t.Col, t.Row)
	t.Col, t.Row = p.Col, p.Row
	r.moveLog(s.ID+"/"+t.ID, t.Name, from, cellText(p.Col, p.Row))
	return nil
}

func (r *Room) fogAdd(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		SceneID string `json:"sceneId"`
		Col     int    `json:"col"`
		Row     int    `json:"row"`
		W       int    `json:"w"`
		H       int    `json:"h"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	sc := r.scene(p.SceneID)
	if sc == nil {
		return errors.New("no such scene")
	}
	if p.W < 1 || p.H < 1 || p.Col < 0 || p.Row < 0 || p.W > sc.Grid.Cols-p.Col || p.H > sc.Grid.Rows-p.Row {
		return errors.New("fog rectangle must be inside the grid")
	}
	sc.Fog = append(sc.Fog, Fog{ID: r.id("f"), Col: p.Col, Row: p.Row, W: p.W, H: p.H})
	r.addLog(fmt.Sprintf("เพิ่มหมอก %d×%d ที่ %s ใน %s", p.W, p.H, cellText(p.Col, p.Row), sc.Name),
		fmt.Sprintf("added fog %d×%d at %s in %s", p.W, p.H, cellText(p.Col, p.Row), sc.Name))
	return nil
}

func (r *Room) note(typ string, payload json.RawMessage) error {
	p, err := decode[struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
		Index int    `json:"index"`
	}](payload)
	if err != nil {
		return err
	}
	title := "" // a removed note's title, read before it goes
	if i := slices.IndexFunc(r.notes.List(), func(n notes.Note) bool { return n.ID == p.ID }); i >= 0 {
		title = r.notes.List()[i].Title
	}
	switch typ {
	case "note.add":
		err = r.notes.Add(p.Title)
	case "note.update":
		err = r.notes.Update(p.ID, p.Title, p.Body)
	case "note.remove":
		err = r.notes.Remove(p.ID)
	default:
		err = r.notes.Move(p.ID, p.Index)
	}
	if err != nil {
		return err
	}
	l := r.notes.List()
	if typ == "note.add" {
		p.ID = l[len(l)-1].ID
	}
	if i := slices.IndexFunc(l, func(n notes.Note) bool { return n.ID == p.ID }); i >= 0 {
		title = l[i].Title
	}
	r.noteLog(p.ID, title)
	return nil
}

var errNpcImage = errors.New("imageUrl must be empty or an uploaded image")

func (r *Room) npc(typ string, payload json.RawMessage) error {
	p, err := decode[npcs.NPC](payload)
	if err != nil {
		return err
	}
	if typ == "npc.remove" {
		n, _ := r.npcs.Get(p.ID)
		if err := r.npcs.Remove(p.ID); err != nil {
			return err
		}
		r.addLog("ลบ NPC "+n.Name+" ออกจากคลัง", "removed NPC "+n.Name+" from the pool")
		return nil
	}
	if p.ImageURL != "" && !r.uploadedFile(p.ImageURL) {
		return errNpcImage
	}
	name := strings.TrimSpace(p.Name)
	if typ == "npc.add" {
		if err := r.npcs.Add(p); err != nil {
			return err
		}
		r.addLog("เพิ่ม NPC "+name+" ในคลัง", "added NPC "+name+" to the pool")
		return nil
	}
	if err := r.npcs.Update(p); err != nil {
		return err
	}
	r.addLog("แก้ NPC "+name, "edited NPC "+name)
	return nil
}

// npcPlace puts a copy of a pool NPC on the first least-crowded cell (row-major) of a scene.
func (r *Room) npcPlace(s *Seat, payload json.RawMessage) error {
	p, err := decode[struct {
		NpcID   string `json:"npcId"`
		SceneID string `json:"sceneId"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	n, ok := r.npcs.Get(p.NpcID)
	if !ok {
		return errors.New("no such npc")
	}
	sc := r.scene(p.SceneID)
	if sc == nil {
		return errors.New("no such scene")
	}
	// A visible token ignores hidden ones for its cell and "#n" suffix, so neither hints that hidden tokens exist.
	// It may stack on a hidden token (the DM can move it); a name may repeat after a reveal.
	// A hidden one counts every token for its cell but never fills a cell, so it always fits.
	used := map[[2]int]int{}
	names := map[string]bool{}
	for _, t := range r.st.Tokens {
		if t.SceneID == sc.ID && (n.Hidden || !t.Hidden) {
			used[[2]int{t.Col, t.Row}]++
			names[t.Name] = true
		}
	}
	cell := 0
	for i := range sc.Grid.Cols * sc.Grid.Rows {
		if used[[2]int{i % sc.Grid.Cols, i / sc.Grid.Cols}] < used[[2]int{cell % sc.Grid.Cols, cell / sc.Grid.Cols}] {
			cell = i
		}
	}
	if !n.Hidden && used[[2]int{cell % sc.Grid.Cols, cell / sc.Grid.Cols}] >= MaxPerCell {
		return errors.New("scene has no free cell")
	}
	name := n.Name // second copy in the scene: "<name> #2", then #3, ...
	for k := 2; names[name]; k++ {
		name = n.Name + " #" + strconv.Itoa(k)
	}
	hp, hpMax := n.HPMax, n.HPMax
	r.st.Tokens = append(r.st.Tokens, Token{
		ID: r.id("t"), SceneID: sc.ID, NpcID: n.ID, Name: name, ImageURL: n.ImageURL,
		Col: cell % sc.Grid.Cols, Row: cell / sc.Grid.Cols, Hidden: n.Hidden, HP: &hp, HPMax: &hpMax,
	})
	th, en := fmt.Sprintf("วาง %s ในฉาก %s", name, sc.Name), fmt.Sprintf("placed %s in %s", name, sc.Name)
	if n.Hidden {
		th, en = th+" (ซ่อน)", en+" (hidden)"
	}
	r.addLog(th, en)
	return nil
}

// tokenNpcEdit handles token.setHidden and token.adjustHp (DM, NPC tokens only).
func (r *Room) tokenNpcEdit(s *Seat, typ string, payload json.RawMessage) error {
	p, err := decode[struct {
		TokenID string `json:"tokenId"`
		Hidden  bool   `json:"hidden"`
		Delta   int    `json:"delta"`
	}](payload)
	if err != nil {
		return err
	}
	if err := r.needDM(s, StatusPlaying); err != nil {
		return err
	}
	t := r.token(p.TokenID)
	if t == nil {
		return errors.New("no such token")
	}
	if t.CharacterID != "" || t.HP == nil || t.HPMax == nil {
		return errors.New("only NPC tokens from the pool")
	}
	if typ == "token.setHidden" {
		if t.Hidden && !p.Hidden && r.visibleAt(t.SceneID, t.Col, t.Row, t.ID) >= MaxPerCell {
			return errCellFull // the DM moves it first
		}
		t.Hidden = p.Hidden
		if p.Hidden {
			r.addLog("ซ่อน "+t.Name, "hid "+t.Name)
		} else {
			r.addLog("เปิดเผย "+t.Name, "revealed "+t.Name)
		}
		return nil
	}
	if p.Delta < -npcs.MaxHP || p.Delta > npcs.MaxHP {
		return fmt.Errorf("delta must be -%d to %d", npcs.MaxHP, npcs.MaxHP)
	}
	old := *t.HP
	*t.HP = min(max(old+p.Delta, 0), *t.HPMax)
	if d := *t.HP - old; d != 0 {
		r.addLog(fmt.Sprintf("%s เลือด %s (%d/%d)", t.Name, signed(d), *t.HP, *t.HPMax), fmt.Sprintf("%s HP %s (%d/%d)", t.Name, signed(d), *t.HP, *t.HPMax))
	}
	r.downed(rules.L{TH: t.Name, EN: t.Name}, old, *t.HP, false, t.Hidden)
	return nil
}
