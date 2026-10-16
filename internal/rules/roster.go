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

type RosterItem struct {
	Name L   `json:"name"`
	Qty  int `json:"qty"`
}

// RosterEntry is one pickable hero (docs/design/heroes-core.md, heroes-legend.md). ID is unique across every set
// (the legend set uses "legend-<class>"). Its class data (name, HP, saves, features) comes from Class ClassID.
type RosterEntry struct {
	ID        string       `json:"id"`
	Set       string       `json:"set"`     // id from Sets(); always the built-in entry's set, never client-chosen
	ClassID   string       `json:"classId"` // Class id; always the built-in entry's class, never client-chosen
	HeroName  L            `json:"heroName"`
	Blurb     L            `json:"blurb"`
	Role      string       `json:"role"`
	Art       string       `json:"art"` // "" = class art by id (frontend asset) | "/uploads/<generated name>"
	Abilities Abilities    `json:"abilities"`
	AC        int          `json:"ac"`
	Skills    []string     `json:"skills"`    // proficient skill ids
	Expertise []string     `json:"expertise"` // subset of Skills
	Items     []RosterItem `json:"items"`
}

const (
	MaxName    = 40 // hero/item names
	MaxDesc    = 300
	MaxQty     = 99
	MaxItems   = 10
	rosterFile = "roster.json"
)

var (
	roles       = []string{"tank", "damage", "healer", "support", "control", "hybrid"}
	artName     = regexp.MustCompile(`^/uploads/[0-9a-f]{32}\.(png|jpg|gif|webp)$`) // same names as room.UploadName
	defaults    = DefaultRoster()                                                   // never mutated
	ErrNoClass  = errors.New("no such class")
	ErrNoSet    = errors.New("no such set")
	errNotSaved = errors.New("could not save roster")
)

// Set is one character set; RosterEntry.Set holds its ID.
type Set struct {
	ID   string `json:"id"`
	Name L      `json:"name"`
}

// Sets returns every character set in display order (fixed; docs/features/character-sets.md).
func Sets() []Set {
	return []Set{
		{ID: "core", Name: L{"ฮีโร่ชุดแรก", "First Company"}},
		{ID: "legend", Name: L{"วีรชนในตำนาน", "Heroes of Legend"}},
	}
}

// IsSet reports whether id is one of Sets().
func IsSet(id string) bool {
	return slices.ContainsFunc(Sets(), func(s Set) bool { return s.ID == id })
}

// DefaultRoster returns a fresh copy of every built-in character: the core set, then the legend set.
func DefaultRoster() []RosterEntry {
	return append(coreRoster(), legendRoster()...)
}

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

// CheckAbilities validates the six scores and AC (roster entries and the DM's character.setAbilities).
func CheckAbilities(a Abilities, ac int) error {
	if !a.valid() {
		return fmt.Errorf("abilities must be %d-%d", MinAbility, MaxAbility)
	}
	if ac < MinAC || ac > MaxAC {
		return fmt.Errorf("ac must be %d-%d", MinAC, MaxAC)
	}
	return nil
}

// CleanSkills checks the proficient skills and expertise: known skill ids, no duplicates, expertise ⊆ skills.
// The DM may give any skill, not only the class's choices (roster.update, character.setSkills). It returns fresh
// non-nil slices.
func CleanSkills(skills, expertise []string) ([]string, []string, error) {
	for i, k := range skills {
		if _, ok := SkillAbility(k); !ok {
			return nil, nil, errors.New("unknown skill id")
		}
		if slices.Contains(skills[:i], k) {
			return nil, nil, errors.New("duplicate skill")
		}
	}
	for i, k := range expertise {
		if !slices.Contains(skills, k) {
			return nil, nil, errors.New("expertise must be a proficient skill")
		}
		if slices.Contains(expertise[:i], k) {
			return nil, nil, errors.New("duplicate expertise")
		}
	}
	return append([]string{}, skills...), append([]string{}, expertise...), nil
}

// cleanEntry trims every text and checks the limits in docs/features/character-roster.md. It returns fresh slices.
// It checks only the art format; the room also checks that an uploaded file exists. Set and ClassID are forced to
// the built-in entry's, so a client can never move a hero to another set or class.
func cleanEntry(e RosterEntry) (RosterEntry, error) {
	i := slices.IndexFunc(defaults, func(d RosterEntry) bool { return d.ID == e.ID })
	if i < 0 {
		return e, ErrNoClass
	}
	e.Set, e.ClassID = defaults[i].Set, defaults[i].ClassID
	var ok bool
	if e.HeroName, ok = cleanL(e.HeroName, MaxName, false); !ok {
		return e, fmt.Errorf("heroName must be 1-%d characters (th and en)", MaxName)
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
	if err := CheckAbilities(e.Abilities, e.AC); err != nil {
		return e, err
	}
	sk, ex, err := CleanSkills(e.Skills, e.Expertise)
	if err != nil {
		return e, err
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
	e.Skills, e.Expertise, e.Items = sk, ex, items
	return e, nil
}

// IsClass reports whether id is a built-in roster entry id (any set).
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
// renamed to a unique roster.json.bad-<random> and the defaults are used (a pre-rules-v2 file does not decode, so
// it lands here too); an entry that fails validation is
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
		log.Warn("roster file is corrupt or an old format, using the default roster", "err", bad.Err.Error(), "moved_to", bad.MovedTo)
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

// List returns the roster in DefaultRoster order (core, then legend). Callers must not mutate it.
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

// Reset restores class id to its default; id "" resets every class (every set).
func (s *RosterStore) Reset(id string) error {
	if id == "" {
		return s.commit(nil)
	}
	if !IsClass(id) {
		return ErrNoClass
	}
	return s.commit(slices.DeleteFunc(slices.Clone(s.edits), func(x RosterEntry) bool { return x.ID == id }))
}

// ResetSet restores every entry of set id to its default; edits of other sets stay.
func (s *RosterStore) ResetSet(id string) error {
	if !IsSet(id) {
		return ErrNoSet
	}
	return s.commit(slices.DeleteFunc(slices.Clone(s.edits), func(x RosterEntry) bool { return x.Set == id }))
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
