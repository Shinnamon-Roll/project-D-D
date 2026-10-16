// Package npcs is the DM's NPC pool, persisted to <dir>/npcs.json (docs/features/npc-pool.md).
// Permissions and the imageUrl check live in room.Apply; this package only validates and stores.
package npcs

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"projectdnd/internal/jsonfile"
)

const (
	MaxName  = 40
	MaxHP    = 999
	MaxNPCs  = 100
	fileName = "npcs.json"
	idPrefix = "npc"
)

type NPC struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"` // "" = default silhouette
	HPMax    int    `json:"hpMax"`
	Hidden   bool   `json:"hidden"` // default for new tokens placed from this NPC
}

// Store is not safe for concurrent use; the room Hub serializes access.
// The zero value is an in-memory store (no file), used by tests.
type Store struct {
	path string // "" = memory only
	log  *slog.Logger
	seq  int
	npcs []NPC // replaced, never mutated in place: List callers may hold the old slice
}

type file struct {
	Seq  int   `json:"seq"`
	NPCs []NPC `json:"npcs"`
}

var (
	errName    = fmt.Errorf("npc name must be 1-%d characters", MaxName)
	errHP      = fmt.Errorf("npc hpMax must be 1-%d", MaxHP)
	errFull    = fmt.Errorf("at most %d npcs", MaxNPCs)
	errNoNPC   = errors.New("no such npc")
	errDupName = errors.New("that name is already in the pool")
	errNotSave = errors.New("could not save npcs")
)

// Load reads <dir>/npcs.json. A missing file is an empty store. A file that is not valid JSON is
// renamed to a unique npcs.json.bad-<random> (never overwritten) and the store starts empty.
func Load(dir string, log *slog.Logger) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	s := &Store{path: filepath.Join(dir, fileName), log: log, npcs: []NPC{}}
	var f file
	var bad *jsonfile.CorruptError
	switch err := jsonfile.Load(s.path, &f); {
	case errors.As(err, &bad):
		log.Error("npcs file is corrupt, starting empty", "err", bad.Err.Error(), "moved_to", bad.MovedTo)
		return s, nil
	case err != nil:
		return nil, err // unreadable: refuse to start rather than overwrite it later
	}
	if f.NPCs != nil {
		s.npcs = f.NPCs
	}
	s.seq = f.Seq
	for _, n := range s.npcs { // a hand-edited or stale seq must never reissue an existing id
		if num, ok := strings.CutPrefix(n.ID, idPrefix); ok {
			if k, err := strconv.Atoi(num); err == nil && k > s.seq {
				s.seq = k
			}
		}
	}
	return s, nil
}

// List returns the pool in order. Callers must not mutate it.
func (s *Store) List() []NPC {
	if s.npcs == nil {
		return []NPC{}
	}
	return s.npcs
}

// Get returns the NPC with id.
func (s *Store) Get(id string) (NPC, bool) {
	if i := s.index(id); i >= 0 {
		return s.npcs[i], true
	}
	return NPC{}, false
}

func (s *Store) index(id string) int {
	return slices.IndexFunc(s.npcs, func(n NPC) bool { return n.ID == id })
}

// clean trims the name and checks the limits. imageUrl is checked by the caller (room).
func clean(n NPC) (NPC, error) {
	n.Name = strings.TrimSpace(n.Name)
	if c := utf8.RuneCountInString(n.Name); c < 1 || c > MaxName {
		return n, errName
	}
	if n.HPMax < 1 || n.HPMax > MaxHP {
		return n, errHP
	}
	return n, nil
}

// nameTaken reports whether an NPC other than self already has name (trimmed, case-insensitive).
// Duplicates in an old npcs.json load as they are; only adds and edits are checked.
func (s *Store) nameTaken(name, self string) bool {
	return slices.ContainsFunc(s.npcs, func(n NPC) bool { return n.ID != self && strings.EqualFold(strings.TrimSpace(n.Name), name) })
}

// Add appends n with a new id (n.ID is ignored).
func (s *Store) Add(n NPC) error {
	n, err := clean(n)
	if err != nil {
		return err
	}
	if len(s.npcs) >= MaxNPCs {
		return errFull
	}
	if s.nameTaken(n.Name, "") {
		return errDupName
	}
	n.ID = idPrefix + strconv.Itoa(s.seq+1)
	return s.commit(s.seq+1, append(slices.Clone(s.npcs), n))
}

// Update replaces the NPC with n.ID. Tokens already placed from it keep their own copy.
func (s *Store) Update(n NPC) error {
	n, err := clean(n)
	if err != nil {
		return err
	}
	i := s.index(n.ID)
	if i < 0 {
		return errNoNPC
	}
	if s.nameTaken(n.Name, n.ID) {
		return errDupName
	}
	next := slices.Clone(s.npcs)
	next[i] = n
	return s.commit(s.seq, next)
}

// Remove deletes the NPC; tokens placed from it stay on the map.
func (s *Store) Remove(id string) error {
	i := s.index(id)
	if i < 0 {
		return errNoNPC
	}
	return s.commit(s.seq, slices.Delete(slices.Clone(s.npcs), i, i+1))
}

// commit saves next to disk first and only then swaps it in, so a failed write changes nothing.
func (s *Store) commit(seq int, next []NPC) error {
	if s.path != "" {
		if err := jsonfile.Save(s.path, file{Seq: seq, NPCs: next}); err != nil {
			s.log.Error("npcs save failed", "err", err.Error())
			return errNotSave
		}
	}
	s.seq, s.npcs = seq, next
	return nil
}
