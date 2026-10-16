package rules

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCampaigns(t *testing.T) {
	c := Campaigns()
	if len(c) != 2 || c[0].ID != "starter-set" || c[0].Status != "available" || c[1].ID != "salem" || c[1].Status != "coming_soon" {
		t.Fatalf("campaigns: %+v", c)
	}
}

func TestSets(t *testing.T) {
	want := []Set{{"core", L{"ฮีโร่ชุดแรก", "First Company"}}, {"legend", L{"วีรชนในตำนาน", "Heroes of Legend"}}}
	if got := Sets(); !reflect.DeepEqual(got, want) {
		t.Fatalf("sets %+v", got)
	}
	for _, tt := range []struct {
		id string
		ok bool
	}{{"core", true}, {"legend", true}, {"", false}, {"Core", false}, {"pirate", false}} {
		if IsSet(tt.id) != tt.ok {
			t.Errorf("IsSet(%q) != %v", tt.id, tt.ok)
		}
	}
}

// docRoster reads the JSON block at the end of a docs/design roster file.
func docRoster(t *testing.T, name string) []RosterEntry {
	t.Helper()
	b, err := os.ReadFile("../../docs/design/" + name)
	if err != nil {
		t.Fatal(err)
	}
	_, block, _ := strings.Cut(string(b), "## JSON (machine-readable)\n\n```json\n")
	block, _, _ = strings.Cut(block, "```")
	var doc []RosterEntry
	if err := json.Unmarshal([]byte(block), &doc); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return doc
}

func TestDefaultRosterMatchesDoc(t *testing.T) {
	core, legend := docRoster(t, "roster.md"), docRoster(t, "roster-legend.md")
	got := DefaultRoster()
	if !reflect.DeepEqual(got, append(core, legend...)) {
		t.Fatal("internal/rules/roster_default.go + roster_legend.go differ from the JSON blocks in docs/design/roster.md + roster-legend.md")
	}
	if len(core) != 10 || len(legend) != 12 {
		t.Fatalf("%d core, %d legend", len(core), len(legend))
	}
	ids := map[string]bool{}
	for i, e := range got {
		if want := map[bool]string{true: "core", false: "legend"}[i < len(core)]; e.Set != want {
			t.Fatalf("%s: set %q, want %q", e.ID, e.Set, want)
		}
		if ids[e.ID] {
			t.Fatalf("duplicate id %s", e.ID)
		}
		ids[e.ID] = true
		if !IsClass(e.ID) {
			t.Fatalf("%s: not a class", e.ID)
		}
		c, err := cleanEntry(e)
		if err != nil || !reflect.DeepEqual(c, e) {
			t.Fatalf("%s: default not clean: %v", e.ID, err)
		}
	}
	got[0].Skills[0].Cost = 99 // copies are independent
	if DefaultRoster()[0].Skills[0].Cost == 99 || defaults[0].Skills[0].Cost == 99 {
		t.Fatal("DefaultRoster shares slices")
	}
}

func TestCleanEntry(t *testing.T) {
	long := strings.Repeat("ก", MaxName+1)
	tests := []struct {
		name   string
		mutate func(e *RosterEntry)
		ok     bool
	}{
		{"default", func(e *RosterEntry) {}, true},
		{"unknown class", func(e *RosterEntry) { e.ID = "pirate" }, false},
		{"client set ignored", func(e *RosterEntry) { e.Set = "legend" }, true},
		{"legend entry", func(e *RosterEntry) { *e = DefaultRoster()[len(DefaultRoster())-1] }, true},
		{"hero th empty", func(e *RosterEntry) { e.HeroName.TH = "  " }, false},
		{"hero en too long", func(e *RosterEntry) { e.HeroName.EN = long }, false},
		{"hero 40 runes", func(e *RosterEntry) { e.HeroName.TH = strings.Repeat("ก", MaxName) }, true},
		{"hero newline", func(e *RosterEntry) { e.HeroName.EN = "a\nb" }, false},
		{"class empty", func(e *RosterEntry) { e.ClassName.EN = "" }, false},
		{"blurb newline ok", func(e *RosterEntry) { e.Blurb.EN = "a\nb" }, true},
		{"blurb control", func(e *RosterEntry) { e.Blurb.EN = "a\x00" }, false},
		{"blurb too long", func(e *RosterEntry) { e.Blurb.TH = strings.Repeat("x", MaxDesc+1) }, false},
		{"bad role", func(e *RosterEntry) { e.Role = "boss" }, false},
		{"art upload ok", func(e *RosterEntry) { e.Art = "/uploads/" + strings.Repeat("a", 32) + ".png" }, true},
		{"art other url", func(e *RosterEntry) { e.Art = "https://x/y.png" }, false},
		{"art traversal", func(e *RosterEntry) { e.Art = "/uploads/../x.png" }, false},
		{"hpMax 0", func(e *RosterEntry) { e.Stats.HPMax = 0 }, false},
		{"hpMax 999", func(e *RosterEntry) { e.Stats.HPMax = 999 }, true},
		{"hpMax 1000", func(e *RosterEntry) { e.Stats.HPMax = 1000 }, false},
		{"str 21", func(e *RosterEntry) { e.Stats.Str = 21 }, false},
		{"agi -1", func(e *RosterEntry) { e.Stats.Agi = -1 }, false},
		{"end 0", func(e *RosterEntry) { e.Stats.End = 0 }, true},
		{"mpMax 100", func(e *RosterEntry) { e.Stats.MPMax = 100 }, false},
		{"mpMax 0", func(e *RosterEntry) { e.Stats.MPMax = 0 }, true},
		{"no skills", func(e *RosterEntry) { e.Skills = nil }, false},
		{"6 skills", func(e *RosterEntry) {
			for i, id := range []string{"a", "b", "c"} {
				e.Skills = append(e.Skills, e.Skills[i])
				e.Skills[len(e.Skills)-1].ID = id
			}
		}, false},
		{"dup skill id", func(e *RosterEntry) { e.Skills[1].ID = e.Skills[0].ID }, false},
		{"bad skill id", func(e *RosterEntry) { e.Skills[0].ID = "Bad Id" }, false},
		{"skill name empty", func(e *RosterEntry) { e.Skills[0].Name.TH = "" }, false},
		{"skill desc empty", func(e *RosterEntry) { e.Skills[0].Desc.EN = "" }, false},
		{"cost 10", func(e *RosterEntry) { e.Skills[0].Cost = 10 }, true},
		{"cost 11", func(e *RosterEntry) { e.Skills[0].Cost = 11 }, false},
		{"cost -1", func(e *RosterEntry) { e.Skills[0].Cost = -1 }, false},
		{"no items", func(e *RosterEntry) { e.Items = nil }, true},
		{"11 items", func(e *RosterEntry) {
			for len(e.Items) < MaxItems+1 {
				e.Items = append(e.Items, e.Items[0])
			}
		}, false},
		{"item qty 0", func(e *RosterEntry) { e.Items[0].Qty = 0 }, false},
		{"item qty 100", func(e *RosterEntry) { e.Items[0].Qty = 100 }, false},
		{"item name empty", func(e *RosterEntry) { e.Items[0].Name.EN = " " }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := DefaultRoster()[0]
			tt.mutate(&e)
			if _, err := cleanEntry(e); (err == nil) != tt.ok {
				t.Fatalf("ok=%v, err=%v", tt.ok, err)
			}
		})
	}
	e := DefaultRoster()[1]
	e.HeroName.EN, e.Skills[0].Name.TH = "  Soraya  ", " x "
	c, err := cleanEntry(e)
	if err != nil || c.HeroName.EN != "Soraya" || c.Skills[0].Name.TH != "x" || c.Items == nil {
		t.Fatalf("trim: %+v %v", c, err)
	}
	for _, tt := range []struct{ id, sent, want string }{
		{"knight", "legend", "core"},
		{"knight", "", "core"},
		{"legend-barbarian", "core", "legend"},
		{"legend-barbarian", "pirate", "legend"},
	} {
		e, ok := (&RosterStore{}).Get(tt.id)
		if !ok {
			t.Fatalf("no %s", tt.id)
		}
		e.Set = tt.sent
		if c, err := cleanEntry(e); err != nil || c.Set != tt.want {
			t.Errorf("%s with set %q: got %q %v, want %q", tt.id, tt.sent, c.Set, err, tt.want)
		}
	}
}

func TestRosterResetSet(t *testing.T) {
	var s RosterStore
	k, b := DefaultRoster()[0], DefaultRoster()[10] // knight (core), legend-barbarian
	k.Stats.Str, b.Stats.Str = 20, 20
	if b.ID != "legend-barbarian" {
		t.Fatalf("index 10 is %s", b.ID)
	}
	for _, e := range []RosterEntry{k, b} {
		if err := s.Update(e); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.ResetSet("pirate"); !errors.Is(err, ErrNoSet) {
		t.Fatalf("unknown set: %v", err)
	}
	if err := s.ResetSet(""); !errors.Is(err, ErrNoSet) {
		t.Fatalf("empty set: %v", err)
	}
	if err := s.ResetSet("legend"); err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Get("legend-barbarian"); got.Stats.Str == 20 {
		t.Fatal("legend edit kept")
	}
	if got, _ := s.Get("knight"); got.Stats.Str != 20 {
		t.Fatal("core edit lost by resetting legend")
	}
	if err := s.ResetSet("core"); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s.List(), DefaultRoster()) {
		t.Fatal("reset core")
	}
}

func TestRosterStore(t *testing.T) {
	dir := t.TempDir()
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, nil))
	s, err := LoadRoster(dir, "", log)
	if err != nil || !reflect.DeepEqual(s.List(), DefaultRoster()) {
		t.Fatalf("missing file: %v", err)
	}
	var zero RosterStore // in-memory store starts at the defaults too
	if !reflect.DeepEqual(zero.List(), DefaultRoster()) {
		t.Fatal("zero store")
	}

	e := DefaultRoster()[2]
	e.HeroName.EN = " Zed "
	if err := s.Update(e); err != nil {
		t.Fatal(err)
	}
	e.ID = "pirate"
	if err := s.Update(e); !errors.Is(err, ErrNoClass) {
		t.Fatalf("unknown class: %v", err)
	}
	if err := s.Reset("pirate"); !errors.Is(err, ErrNoClass) {
		t.Fatalf("reset unknown: %v", err)
	}
	reload := func() *RosterStore {
		t.Helper()
		s2, err := LoadRoster(dir, "", log)
		if err != nil {
			t.Fatal(err)
		}
		return s2
	}
	if got, _ := reload().Get("rogue"); got.HeroName.EN != "Zed" {
		t.Fatalf("persisted %q", got.HeroName.EN)
	}
	if fi, err := os.Stat(filepath.Join(dir, rosterFile)); err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("file mode: %v %v", fi, err)
	}
	k := DefaultRoster()[0]
	k.Stats.Str = 20
	if err := s.Update(k); err != nil {
		t.Fatal(err)
	}
	if err := s.Reset("rogue"); err != nil {
		t.Fatal(err)
	}
	r2 := reload()
	if got, _ := r2.Get("rogue"); got.HeroName.EN != "Zarif Kesh" {
		t.Fatalf("reset rogue: %q", got.HeroName.EN)
	}
	if got, _ := r2.Get("knight"); got.Stats.Str != 20 {
		t.Fatal("knight edit lost by resetting rogue")
	}
	if err := s.Reset(""); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(reload().List(), DefaultRoster()) {
		t.Fatal("reset all")
	}

	// A failed save changes nothing.
	s.path = filepath.Join(dir, "missing-dir", rosterFile)
	if err := s.Update(k); err == nil || err.Error() != "could not save roster" {
		t.Fatalf("save failure: %v", err)
	}
	if got, _ := s.Get("knight"); got.Stats.Str == 20 {
		t.Fatal("failed save changed the roster")
	}
}

func TestLoadRosterBadFiles(t *testing.T) {
	log := slog.New(slog.DiscardHandler)
	t.Run("corrupt json is quarantined", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, rosterFile), []byte("{nope"), 0o600)
		s, err := LoadRoster(dir, "", log)
		if err != nil || !reflect.DeepEqual(s.List(), DefaultRoster()) {
			t.Fatalf("err %v", err)
		}
		if m, _ := filepath.Glob(filepath.Join(dir, rosterFile+".bad-*")); len(m) != 1 {
			t.Fatalf("quarantine files %v", m)
		}
	})
	t.Run("invalid and duplicate entries fall back to defaults", func(t *testing.T) {
		dir := t.TempDir()
		good, bad, dup := DefaultRoster()[0], DefaultRoster()[1], DefaultRoster()[0]
		good.Stats.Str, bad.Stats.HPMax, dup.Stats.Str = 20, 0, 1
		b, _ := json.Marshal(rosterFileData{Edits: []RosterEntry{good, bad, dup, {ID: "pirate"}}})
		os.WriteFile(filepath.Join(dir, rosterFile), b, 0o600)
		s, err := LoadRoster(dir, "", log)
		if err != nil {
			t.Fatal(err)
		}
		if k, _ := s.Get("knight"); k.Stats.Str != 20 {
			t.Fatal("valid edit not loaded")
		}
		if w, _ := s.Get("wizard"); w.Stats.HPMax != DefaultRoster()[1].Stats.HPMax {
			t.Fatal("invalid edit loaded")
		}
		if len(s.List()) != len(DefaultRoster()) {
			t.Fatal("roster size")
		}
	})
	t.Run("old file without set still loads", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, rosterFile), []byte(`{"edits":[`+strings.Replace(mustMarshal(t, func() RosterEntry {
			e := DefaultRoster()[0]
			e.Stats.Str = 20
			return e
		}()), `"set":"core",`, "", 1)+`]}`), 0o600)
		s, err := LoadRoster(dir, "", log)
		if err != nil {
			t.Fatal(err)
		}
		if k, _ := s.Get("knight"); k.Stats.Str != 20 || k.Set != "core" {
			t.Fatalf("old edit: str %d set %q", k.Stats.Str, k.Set)
		}
		if err := s.ResetSet("core"); err != nil {
			t.Fatal(err)
		}
		if k, _ := s.Get("knight"); k.Stats.Str == 20 {
			t.Fatal("old edit not reset by its set")
		}
	})
	t.Run("art whose upload is gone falls back to the class art", func(t *testing.T) {
		const kept, gone = "/uploads/0123456789abcdef0123456789abcdef.png", "/uploads/fedcba9876543210fedcba9876543210.webp"
		dir, uploads := t.TempDir(), t.TempDir()
		os.WriteFile(filepath.Join(uploads, "0123456789abcdef0123456789abcdef.png"), []byte("x"), 0o600)
		os.Mkdir(filepath.Join(uploads, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.gif"), 0o700) // a dir, not a file
		k, w, r := DefaultRoster()[0], DefaultRoster()[1], DefaultRoster()[2]
		k.Art, w.Art, r.Art = kept, gone, "/uploads/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.gif"
		b, _ := json.Marshal(rosterFileData{Edits: []RosterEntry{k, w, r}})
		os.WriteFile(filepath.Join(dir, rosterFile), b, 0o600)
		for _, tc := range []struct {
			name, uploads, class, want string
		}{
			{"file present", uploads, "knight", kept},
			{"file gone", uploads, "wizard", ""},
			{"dir, not a file", uploads, "rogue", ""},
			{"no uploads dir", "", "knight", ""},
		} {
			s, err := LoadRoster(dir, tc.uploads, log)
			if err != nil {
				t.Fatal(err)
			}
			if e, _ := s.Get(tc.class); e.Art != tc.want {
				t.Errorf("%s: art %q, want %q", tc.name, e.Art, tc.want)
			}
		}
	})
}

func mustMarshal(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"set":"core",`) {
		t.Fatalf("no set key in %s", b)
	}
	return string(b)
}
