package npcs

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func ids(s *Store) string {
	var b []string
	for _, n := range s.List() {
		b = append(b, n.ID)
	}
	return strings.Join(b, ",")
}

func TestLimits(t *testing.T) {
	s := &Store{}
	for _, tc := range []struct {
		name string
		err  error
		do   func() error
	}{
		{"empty name", errName, func() error { return s.Add(NPC{Name: "  ", HPMax: 1}) }},
		{"41 runes", errName, func() error { return s.Add(NPC{Name: strings.Repeat("ก", 41), HPMax: 1}) }},
		{"40 runes", nil, func() error { return s.Add(NPC{Name: strings.Repeat("ก", 40), HPMax: 1}) }},
		{"hp 0", errHP, func() error { return s.Add(NPC{Name: "a", HPMax: 0}) }},
		{"hp -1", errHP, func() error { return s.Add(NPC{Name: "a", HPMax: -1}) }},
		{"hp 1000", errHP, func() error { return s.Add(NPC{Name: "a", HPMax: 1000}) }},
		{"hp 999", nil, func() error { return s.Add(NPC{Name: "a", HPMax: 999}) }},
		{"update", nil, func() error { return s.Update(NPC{ID: "npc1", Name: " Gob ", HPMax: 5, Hidden: true}) }},
		{"update bad hp", errHP, func() error { return s.Update(NPC{ID: "npc1", Name: "Gob", HPMax: 0}) }},
		{"update unknown", errNoNPC, func() error { return s.Update(NPC{ID: "npc9", Name: "x", HPMax: 1}) }},
		{"remove unknown", errNoNPC, func() error { return s.Remove("npc9") }},
	} {
		if err := tc.do(); err != tc.err {
			t.Errorf("%s: err %v, want %v", tc.name, err, tc.err)
		}
	}
	if n, _ := s.Get("npc1"); n != (NPC{ID: "npc1", Name: "Gob", HPMax: 5, Hidden: true}) {
		t.Fatalf("update not applied: %+v", n)
	}
	for len(s.List()) < MaxNPCs {
		if err := s.Add(NPC{Name: "x" + strconv.Itoa(len(s.List())), HPMax: 1}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Add(NPC{Name: "one too many", HPMax: 1}); err != errFull {
		t.Fatalf("101st npc: %v", err)
	}
}

func TestUniqueName(t *testing.T) {
	s := &Store{}
	for _, n := range []string{"Goblin", "Wolf"} {
		if err := s.Add(NPC{Name: n, HPMax: 1}); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		name string
		err  error
		do   func() error
	}{
		{"add same", errDupName, func() error { return s.Add(NPC{Name: "Goblin", HPMax: 1}) }},
		{"add other case + spaces", errDupName, func() error { return s.Add(NPC{Name: "  gOBLIN ", HPMax: 1}) }},
		{"add numbered is another name", nil, func() error { return s.Add(NPC{Name: "Goblin #5", HPMax: 1}) }},
		{"update to another's name", errDupName, func() error { return s.Update(NPC{ID: "npc2", Name: "goblin", HPMax: 1}) }},
		{"update keeps own name", nil, func() error { return s.Update(NPC{ID: "npc1", Name: "Goblin", HPMax: 3}) }},
		{"update own name other case", nil, func() error { return s.Update(NPC{ID: "npc1", Name: "GOBLIN", HPMax: 3}) }},
		{"update to a free name", nil, func() error { return s.Update(NPC{ID: "npc2", Name: "Wolf Alpha", HPMax: 1}) }},
		{"add the freed name", nil, func() error { return s.Add(NPC{Name: "Wolf", HPMax: 1}) }},
	} {
		if err := tc.do(); err != tc.err {
			t.Errorf("%s: err %v, want %v", tc.name, err, tc.err)
		}
	}
	if len(s.List()) != 4 {
		t.Fatalf("pool %+v", s.List())
	}
}

// An old npcs.json may hold duplicate names: it loads unchanged; only new adds and edits must be unique.
func TestLoadKeepsOldDuplicates(t *testing.T) {
	dir := t.TempDir()
	old := `{"seq":2,"npcs":[{"id":"npc1","name":"Goblin","hpMax":1},{"id":"npc2","name":"goblin","hpMax":2}]}`
	if err := os.WriteFile(filepath.Join(dir, fileName), []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := Load(dir, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if l := s.List(); len(l) != 2 || l[0].Name != "Goblin" || l[1].Name != "goblin" {
		t.Fatalf("loaded %+v", l)
	}
	if err := s.Add(NPC{Name: "GOBLIN", HPMax: 1}); err != errDupName {
		t.Fatalf("add: %v", err)
	}
	if err := s.Update(NPC{ID: "npc2", Name: "Goblin Boss", HPMax: 2}); err != nil {
		t.Fatalf("rename a duplicate: %v", err)
	}
}

func TestPersistReload(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	log := slog.New(slog.DiscardHandler)
	s, err := Load(dir, log)
	if err != nil {
		t.Fatal(err)
	}
	held := s.List()
	for _, step := range []func() error{
		func() error { return s.Add(NPC{Name: "Goblin", HPMax: 7, Hidden: true}) },
		func() error { return s.Add(NPC{Name: "Wolf", HPMax: 11}) },
		func() error { return s.Add(NPC{Name: "Gone", HPMax: 1}) },
		func() error { return s.Remove("npc3") },
	} {
		if err := step(); err != nil {
			t.Fatal(err)
		}
	}
	if len(held) != 0 {
		t.Fatal("an old List() slice was mutated")
	}
	if fi, err := os.Stat(filepath.Join(dir, fileName)); err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("file mode %v %v", fi.Mode().Perm(), err)
	}
	r, err := Load(dir, log)
	if err != nil {
		t.Fatal(err)
	}
	if ids(r) != "npc1,npc2" || !r.List()[0].Hidden || r.List()[1].HPMax != 11 {
		t.Fatalf("reload: %+v", r.List())
	}
	if err := r.Add(NPC{Name: "New", HPMax: 1}); err != nil || ids(r) != "npc1,npc2,npc4" {
		t.Fatalf("seq not restored (ids never reused): %s %v", ids(r), err)
	}
}

func TestSeqAtLeastHighestID(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, fileName), []byte(`{"seq":1,"npcs":[{"id":"npc7","name":"a","hpMax":1}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := Load(dir, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Add(NPC{Name: "b", HPMax: 1}); err != nil || ids(s) != "npc7,npc8" {
		t.Fatalf("%s %v", ids(s), err)
	}
}

func TestCorruptFileQuarantined(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, []byte(`{"npcs": [`), 0o600); err != nil {
		t.Fatal(err)
	}
	var logs bytes.Buffer
	s, err := Load(dir, slog.New(slog.NewTextHandler(&logs, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if len(s.List()) != 0 || !strings.Contains(logs.String(), "level=ERROR") {
		t.Fatalf("store %v, logs %s", s.List(), logs.String())
	}
	if bad, _ := filepath.Glob(path + ".bad-*"); len(bad) != 1 {
		t.Fatalf("quarantine files: %v", bad)
	}
}

func TestFailedSaveChangesNothing(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Add(NPC{Name: "kept", HPMax: 1}); err != nil {
		t.Fatal(err)
	}
	s.path = filepath.Join(dir, "missing-dir", fileName) // CreateTemp fails
	if err := s.Add(NPC{Name: "lost", HPMax: 1}); err != errNotSave {
		t.Fatalf("err %v", err)
	}
	if ids(s) != "npc1" || s.seq != 1 {
		t.Fatalf("state changed on failed save: %s seq=%d", ids(s), s.seq)
	}
}
