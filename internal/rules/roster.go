package rules

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"projectdnd/internal/jsonfile"
)

// L is game text in both UI languages.
type L struct {
	TH string `json:"th"`
	EN string `json:"en"`
}

type Stats struct {
	HPMax int `json:"hpMax"`
	Str   int `json:"str"`
	End   int `json:"end"`
	Agi   int `json:"agi"`
	MPMax int `json:"mpMax"`
}

type RosterSkill struct {
	ID   string `json:"id"`
	Name L      `json:"name"`
	Desc L      `json:"desc"`
	Cost int    `json:"cost"`
}

type RosterItem struct {
	Name L   `json:"name"`
	Qty  int `json:"qty"`
}

// RosterEntry is one pickable character; ID is its class id.
type RosterEntry struct {
	ID        string        `json:"id"`
	HeroName  L             `json:"heroName"`
	ClassName L             `json:"className"`
	Blurb     L             `json:"blurb"`
	Role      string        `json:"role"`
	Art       string        `json:"art"` // "" = class art by id (frontend asset) | "/uploads/<generated name>"
	Stats     Stats         `json:"stats"`
	Skills    []RosterSkill `json:"skills"`
	Items     []RosterItem  `json:"items"`
}

const (
	MaxName      = 40 // hero/class/skill/item names
	MaxDesc      = 300
	MaxCost      = 10
	MaxQty       = 99
	MaxSkills    = 5
	MaxItems     = 10
	rosterFile   = "roster.json"
	maxStatValue = 20
)

var (
	roles       = []string{"tank", "damage", "healer", "support", "control", "hybrid"}
	skillID     = regexp.MustCompile(`^[a-z0-9-]{1,40}$`)
	artName     = regexp.MustCompile(`^/uploads/[0-9a-f]{32}\.(png|jpg|gif|webp)$`) // same names as room.UploadName
	defaults    = DefaultRoster()                                                   // never mutated
	ErrNoClass  = errors.New("no such class")
	errNotSaved = errors.New("could not save roster")
)

// CleanText trims s and checks it has min..max runes and no control characters (multiline allows '\n').
func CleanText(s string, min, max int, multiline bool) (string, bool) {
	s = strings.TrimSpace(s)
	if n := utf8.RuneCountInString(s); n < min || n > max {
		return s, false
	}
	for _, c := range s {
		if unicode.IsControl(c) && !(multiline && c == '\n') {
			return s, false
		}
	}
	return s, true
}

func cleanL(l L, max int, multiline bool) (L, bool) {
	th, a := CleanText(l.TH, 1, max, multiline)
	en, b := CleanText(l.EN, 1, max, multiline)
	return L{th, en}, a && b
}

// Check validates stat limits (roster entries and the DM's character.setStats).
func (s Stats) Check() error {
	switch {
	case s.HPMax < 1 || s.HPMax > 999:
		return errors.New("hpMax must be 1-999")
	case s.MPMax < 0 || s.MPMax > 99:
		return errors.New("mpMax must be 0-99")
	}
	for _, v := range []int{s.Str, s.End, s.Agi} {
		if v < 0 || v > maxStatValue {
			return fmt.Errorf("str, end and agi must be 0-%d", maxStatValue)
		}
	}
	return nil
}

// cleanEntry trims every text and checks the limits in docs/features/character-roster.md. It returns fresh slices.
// It checks only the art format; the room also checks that an uploaded file exists.
func cleanEntry(e RosterEntry) (RosterEntry, error) {
	if !IsClass(e.ID) {
		return e, ErrNoClass
	}
	var ok bool
	if e.HeroName, ok = cleanL(e.HeroName, MaxName, false); !ok {
		return e, fmt.Errorf("heroName must be 1-%d characters (th and en)", MaxName)
	}
	if e.ClassName, ok = cleanL(e.ClassName, MaxName, false); !ok {
		return e, fmt.Errorf("className must be 1-%d characters (th and en)", MaxName)
	}
	if e.Blurb, ok = cleanL(e.Blurb, MaxDesc, true); !ok {
		return e, fmt.Errorf("blurb must be 1-%d characters (th and en)", MaxDesc)
	}
	if !slices.Contains(roles, e.Role) {
		return e, errors.New("role must be one of " + strings.Join(roles, ", "))
	}
	if e.Art != "" && !artName.MatchString(e.Art) {
		return e, errors.New("art must be empty or an uploaded image")
	}
	if err := e.Stats.Check(); err != nil {
		return e, err
	}
	if len(e.Skills) < 1 || len(e.Skills) > MaxSkills {
		return e, fmt.Errorf("a class needs 1-%d skills", MaxSkills)
	}
	skills := make([]RosterSkill, len(e.Skills))
	for i, k := range e.Skills {
		if !skillID.MatchString(k.ID) {
			return e, errors.New("skill id must be 1-40 characters a-z, 0-9 or -")
		}
		if slices.ContainsFunc(skills[:i], func(x RosterSkill) bool { return x.ID == k.ID }) {
			return e, errors.New("duplicate skill id")
		}
		if k.Name, ok = cleanL(k.Name, MaxName, false); !ok {
			return e, fmt.Errorf("skill name must be 1-%d characters (th and en)", MaxName)
		}
		if k.Desc, ok = cleanL(k.Desc, MaxDesc, true); !ok {
			return e, fmt.Errorf("skill desc must be 1-%d characters (th and en)", MaxDesc)
		}
		if k.Cost < 0 || k.Cost > MaxCost {
			return e, fmt.Errorf("skill cost must be 0-%d", MaxCost)
		}
		skills[i] = k
	}
	if len(e.Items) > MaxItems {
		return e, fmt.Errorf("a class starts with at most %d items", MaxItems)
	}
	items := make([]RosterItem, len(e.Items))
	for i, it := range e.Items {
		if it.Name, ok = cleanL(it.Name, MaxName, false); !ok {
			return e, fmt.Errorf("item name must be 1-%d characters (th and en)", MaxName)
		}
		if it.Qty < 1 || it.Qty > MaxQty {
			return e, fmt.Errorf("item qty must be 1-%d", MaxQty)
		}
		items[i] = it
	}
	e.Skills, e.Items = skills, items
	return e, nil
}

// IsClass reports whether id is one of the 10 class ids.
func IsClass(id string) bool {
	return slices.ContainsFunc(defaults, func(e RosterEntry) bool { return e.ID == id })
}

// RosterStore is the roster the DM may edit, persisted to <dir>/roster.json. The file holds only the edited
// entries, so unedited classes follow the built-in defaults. Not safe for concurrent use; the room Hub
// serializes access. The zero value is the default roster in memory (no file), used by tests.
type RosterStore struct {
	path  string // "" = memory only
	log   *slog.Logger
	edits []RosterEntry // replaced, never mutated in place
	list  []RosterEntry // defaults with edits applied; nil = defaults
}

type rosterFileData struct {
	Edits []RosterEntry `json:"edits"`
}

// LoadRoster reads <dir>/roster.json. A missing file is the default roster. A file that is not valid JSON is
// renamed to a unique roster.json.bad-<random> and the defaults are used; an entry that fails validation is
// dropped (that class falls back to its default) and logged. An art whose file is gone from the uploads dir
// becomes "" (the class art); with uploads "" every uploaded art is treated as gone.
func LoadRoster(dir, uploads string, log *slog.Logger) (*RosterStore, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	s := &RosterStore{path: filepath.Join(dir, rosterFile), log: log}
	var f rosterFileData
	var bad *jsonfile.CorruptError
	switch err := jsonfile.Load(s.path, &f); {
	case errors.As(err, &bad):
		log.Error("roster file is corrupt, using the default roster", "err", bad.Err.Error(), "moved_to", bad.MovedTo)
		return s, nil
	case err != nil:
		return nil, err // unreadable: refuse to start rather than overwrite it later
	}
	var edits []RosterEntry
	for _, e := range f.Edits {
		c, err := cleanEntry(e)
		if err == nil && slices.ContainsFunc(edits, func(x RosterEntry) bool { return x.ID == c.ID }) {
			err = errors.New("duplicate class")
		}
		if err != nil {
			log.Warn("roster entry ignored, using its default", "class", c.ID, "err", err.Error())
			continue
		}
		if c.Art != "" && !uploadExists(uploads, c.Art) {
			log.Warn("roster art file is missing, using the class art", "class", c.ID, "art", c.Art)
			c.Art = ""
		}
		edits = append(edits, c)
	}
	s.edits, s.list = edits, merge(edits)
	return s, nil
}

// uploadExists reports whether art (already matched by artName, so a plain file name) is a regular file in uploads.
func uploadExists(uploads, art string) bool {
	if uploads == "" {
		return false
	}
	fi, err := os.Stat(filepath.Join(uploads, strings.TrimPrefix(art, "/uploads/")))
	return err == nil && fi.Mode().IsRegular()
}

func merge(edits []RosterEntry) []RosterEntry {
	if len(edits) == 0 {
		return nil
	}
	list := slices.Clone(defaults)
	for _, e := range edits {
		list[slices.IndexFunc(list, func(x RosterEntry) bool { return x.ID == e.ID })] = e
	}
	return list
}

// List returns the roster in class order. Callers must not mutate it.
func (s *RosterStore) List() []RosterEntry {
	if s.list == nil {
		return defaults
	}
	return s.list
}

// Get returns the entry for class id. Its slices are shared: callers must copy before changing them.
func (s *RosterStore) Get(id string) (RosterEntry, bool) {
	l := s.List()
	if i := slices.IndexFunc(l, func(e RosterEntry) bool { return e.ID == id }); i >= 0 {
		return l[i], true
	}
	return RosterEntry{}, false
}

// Update replaces the entry for class e.ID.
func (s *RosterStore) Update(e RosterEntry) error {
	e, err := cleanEntry(e)
	if err != nil {
		return err
	}
	edits := slices.DeleteFunc(slices.Clone(s.edits), func(x RosterEntry) bool { return x.ID == e.ID })
	return s.commit(append(edits, e))
}

// Reset restores class id to its default; id "" resets every class.
func (s *RosterStore) Reset(id string) error {
	if id == "" {
		return s.commit(nil)
	}
	if !IsClass(id) {
		return ErrNoClass
	}
	return s.commit(slices.DeleteFunc(slices.Clone(s.edits), func(x RosterEntry) bool { return x.ID == id }))
}

// commit saves edits to disk first and only then swaps them in, so a failed write changes nothing.
func (s *RosterStore) commit(edits []RosterEntry) error {
	if edits == nil {
		edits = []RosterEntry{}
	}
	if s.path != "" {
		if err := jsonfile.Save(s.path, rosterFileData{Edits: edits}); err != nil {
			s.log.Error("roster save failed", "err", err.Error())
			return errNotSaved
		}
	}
	s.edits, s.list = edits, merge(edits)
	return nil
}
