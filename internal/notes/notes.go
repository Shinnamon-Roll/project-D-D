// Package notes is the DM's story notes, persisted to <dir>/notes.json (docs/features/dm-notes.md).
// Permissions live in room.Apply; this package only validates and stores.
package notes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxTitle = 80
	MaxBody  = 20000
	MaxNotes = 200
	fileName = "notes.json"
)

type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	UpdatedAt int64  `json:"updatedAt"` // unix milliseconds
}

// Store is not safe for concurrent use; the room Hub serializes access.
// The zero value is an in-memory store (no file), used by tests.
type Store struct {
	path  string // "" = memory only
	log   *slog.Logger
	seq   int
	notes []Note // replaced, never mutated in place: List callers may hold the old slice
}

type file struct {
	Seq   int    `json:"seq"`
	Notes []Note `json:"notes"`
}

var (
	errTitle   = fmt.Errorf("note title must be 1-%d characters", MaxTitle)
	errBody    = fmt.Errorf("note body must be at most %d characters", MaxBody)
	errFull    = fmt.Errorf("at most %d notes", MaxNotes)
	errNoNote  = errors.New("no such note")
	errIndex   = errors.New("note index out of range")
	errNotSave = errors.New("could not save notes")
)

// Load reads <dir>/notes.json. A missing file is an empty store. A file that is not valid JSON is
// renamed to a unique notes.json.bad-<random> (never overwritten) and the store starts empty.
func Load(dir string, log *slog.Logger) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	s := &Store{path: filepath.Join(dir, fileName), log: log, notes: []Note{}}
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err // unreadable: refuse to start rather than overwrite it later
	}
	var f file
	if err := json.Unmarshal(b, &f); err != nil {
		bad, rerr := quarantine(s.path)
		if rerr != nil {
			return nil, fmt.Errorf("notes file is corrupt and could not be moved aside: %w", rerr)
		}
		log.Error("notes file is corrupt, starting empty", "err", err.Error(), "moved_to", bad)
		return s, nil
	}
	if f.Notes != nil {
		s.notes = f.Notes
	}
	s.seq = f.Seq
	for _, n := range s.notes { // a hand-edited or stale seq must never reissue an existing id
		if num, ok := strings.CutPrefix(n.ID, "n"); ok {
			if k, err := strconv.Atoi(num); err == nil && k > s.seq {
				s.seq = k
			}
		}
	}
	return s, nil
}

// quarantine moves path to a fresh unique path.bad-<random> and returns that name.
// CreateTemp reserves the name, so two quarantines never collide; the rename then replaces the empty file.
func quarantine(path string) (string, error) {
	f, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".bad-*")
	if err != nil {
		return "", err
	}
	f.Close()
	if err := os.Rename(path, f.Name()); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// List returns the notes in order. Callers must not mutate it.
func (s *Store) List() []Note {
	if s.notes == nil {
		return []Note{}
	}
	return s.notes
}

func cleanTitle(t string) (string, error) {
	t = strings.TrimSpace(t)
	if n := utf8.RuneCountInString(t); n < 1 || n > MaxTitle {
		return "", errTitle
	}
	return t, nil
}

func (s *Store) index(id string) int {
	return slices.IndexFunc(s.notes, func(n Note) bool { return n.ID == id })
}

func now() int64 { return time.Now().UnixMilli() }

// Add appends a note with an empty body.
func (s *Store) Add(title string) error {
	title, err := cleanTitle(title)
	if err != nil {
		return err
	}
	if len(s.notes) >= MaxNotes {
		return errFull
	}
	n := Note{ID: fmt.Sprintf("n%d", s.seq+1), Title: title, UpdatedAt: now()}
	return s.commit(s.seq+1, append(slices.Clone(s.notes), n))
}

// Update replaces a note's title and body (last write wins).
func (s *Store) Update(id, title, body string) error {
	title, err := cleanTitle(title)
	if err != nil {
		return err
	}
	if utf8.RuneCountInString(body) > MaxBody {
		return errBody
	}
	i := s.index(id)
	if i < 0 {
		return errNoNote
	}
	next := slices.Clone(s.notes)
	next[i] = Note{ID: id, Title: title, Body: body, UpdatedAt: now()}
	return s.commit(s.seq, next)
}

func (s *Store) Remove(id string) error {
	i := s.index(id)
	if i < 0 {
		return errNoNote
	}
	return s.commit(s.seq, slices.Delete(slices.Clone(s.notes), i, i+1))
}

// Move puts the note at position index of the resulting list (0 ≤ index < len).
func (s *Store) Move(id string, index int) error {
	i := s.index(id)
	if i < 0 {
		return errNoNote
	}
	if index < 0 || index >= len(s.notes) {
		return errIndex
	}
	n := s.notes[i]
	next := slices.Delete(slices.Clone(s.notes), i, i+1)
	return s.commit(s.seq, slices.Insert(next, index, n))
}

// commit saves next to disk first and only then swaps it in, so a failed write changes nothing.
func (s *Store) commit(seq int, next []Note) error {
	if s.path != "" {
		if err := s.save(file{Seq: seq, Notes: next}); err != nil {
			s.log.Error("notes save failed", "err", err.Error())
			return errNotSave
		}
	}
	s.seq, s.notes = seq, next
	return nil
}

// save writes atomically: temp file in the same dir (0600) + fsync + rename, then fsync the dir.
func (s *Store) save(f file) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, fileName+".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // no-op after a successful rename
	if _, err := tmp.Write(b); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), s.path); err != nil {
		return err
	}
	if d, err := os.Open(dir); err == nil {
		d.Sync() // makes the rename durable; best effort (not supported everywhere)
		d.Close()
	}
	return nil
}
