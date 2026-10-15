package notes

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
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
		{"empty title", errTitle, func() error { return s.Add("   ") }},
		{"81 runes", errTitle, func() error { return s.Add(strings.Repeat("ก", 81)) }},
		{"80 runes", nil, func() error { return s.Add(strings.Repeat("ก", 80)) }},
		{"body 20000", nil, func() error { return s.Update("n1", "t", strings.Repeat("ข", MaxBody)) }},
		{"body 20001", errBody, func() error { return s.Update("n1", "t", strings.Repeat("ข", MaxBody+1)) }},
		{"update empty title", errTitle, func() error { return s.Update("n1", "", "") }},
		{"update unknown", errNoNote, func() error { return s.Update("n9", "t", "") }},
		{"remove unknown", errNoNote, func() error { return s.Remove("n9") }},
		{"move unknown", errNoNote, func() error { return s.Move("n9", 0) }},
		{"move -1", errIndex, func() error { return s.Move("n1", -1) }},
		{"move past end", errIndex, func() error { return s.Move("n1", 1) }},
		{"move to self", nil, func() error { return s.Move("n1", 0) }},
	} {
		if err := tc.do(); err != tc.err {
			t.Errorf("%s: err %v, want %v", tc.name, err, tc.err)
		}
	}
	if n := s.List()[0]; n.Title != "t" || n.UpdatedAt == 0 {
		t.Fatalf("update not applied: %+v", n)
	}

	for len(s.List()) < MaxNotes {
		if err := s.Add("x"); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Add("one too many"); err != errFull {
		t.Fatalf("201st note: %v", err)
	}
}

func TestMoveRemoveOrder(t *testing.T) {
	s := &Store{}
	for _, title := range []string{"a", "b", "c", "d"} {
		if err := s.Add(title); err != nil {
			t.Fatal(err)
		}
	}
	held := s.List()
	steps := []struct {
		do   func() error
		want string
	}{
		{func() error { return s.Move("n1", 3) }, "n2,n3,n4,n1"},
		{func() error { return s.Move("n4", 0) }, "n4,n2,n3,n1"},
		{func() error { return s.Remove("n2") }, "n4,n3,n1"},
		{func() error { return s.Add("e") }, "n4,n3,n1,n5"}, // ids never reused
	}
	for i, st := range steps {
		if err := st.do(); err != nil {
			t.Fatal(i, err)
		}
		if got := ids(s); got != st.want {
			t.Fatalf("step %d: %s, want %s", i, got, st.want)
		}
	}
	if got := held[0].ID + held[3].ID; got != "n1n4" {
		t.Fatalf("an old List() slice was mutated: %s", got)
	}
}

func TestPersistReload(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	log := slog.New(slog.DiscardHandler)
	s, err := Load(dir, log)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.List()) != 0 {
		t.Fatal("fresh store not empty")
	}
	for _, step := range []func() error{
		func() error { return s.Add("Chapter 1") },
		func() error { return s.Add("Chapter 2") },
		func() error { return s.Update("n1", "Chapter 1", "มังกร\nline 2") },
		func() error { return s.Move("n2", 0) },
	} {
		if err := step(); err != nil {
			t.Fatal(err)
		}
	}
	if fi, err := os.Stat(dir); err != nil || fi.Mode().Perm() != 0o700 {
		t.Fatalf("dir mode %v %v", fi.Mode().Perm(), err)
	}
	if fi, err := os.Stat(filepath.Join(dir, fileName)); err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("file mode %v %v", fi.Mode().Perm(), err)
	}
	if m, _ := filepath.Glob(filepath.Join(dir, "*.tmp-*")); len(m) != 0 {
		t.Fatalf("temp files left: %v", m)
	}

	r, err := Load(dir, log)
	if err != nil {
		t.Fatal(err)
	}
	if ids(r) != "n2,n1" || r.List()[1].Body != "มังกร\nline 2" {
		t.Fatalf("reload: %+v", r.List())
	}
	if err := r.Add("Chapter 3"); err != nil || ids(r) != "n2,n1,n3" {
		t.Fatalf("seq not restored: %s %v", ids(r), err)
	}
}

func TestCorruptFileQuarantined(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, []byte(`{"notes": [`), 0o600); err != nil {
		t.Fatal(err)
	}
	var logs bytes.Buffer
	s, err := Load(dir, slog.New(slog.NewTextHandler(&logs, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if len(s.List()) != 0 {
		t.Fatal("corrupt store not empty")
	}
	bad, _ := filepath.Glob(path + ".bad-*")
	if len(bad) != 1 {
		t.Fatalf("quarantine files: %v", bad)
	}
	if b, _ := os.ReadFile(bad[0]); string(b) != `{"notes": [` {
		t.Fatalf("quarantined content changed: %q", b)
	}
	if !strings.Contains(logs.String(), "level=ERROR") {
		t.Fatalf("no error logged: %s", logs.String())
	}
	if err := s.Add("fresh"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(bad[0]); string(b) != `{"notes": [` {
		t.Fatal("quarantined file overwritten")
	}
}

// TestCorruptFileQuarantinedTwice: a second corrupt file gets its own quarantine name (no overwrite).
func TestCorruptFileQuarantinedTwice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, fileName)
	log := slog.New(slog.DiscardHandler)
	for _, junk := range []string{"first", "second"} {
		if err := os.WriteFile(path, []byte(junk), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(dir, log); err != nil {
			t.Fatal(err)
		}
	}
	bad, _ := filepath.Glob(path + ".bad-*")
	var got []string
	for _, b := range bad {
		c, _ := os.ReadFile(b)
		got = append(got, string(c))
	}
	slices.Sort(got)
	if strings.Join(got, ",") != "first,second" {
		t.Fatalf("quarantine files: %v %v", bad, got)
	}
}

func TestSeqAtLeastHighestID(t *testing.T) {
	for _, tc := range []struct {
		name, file, next string
	}{
		{"stale seq", `{"seq":1,"notes":[{"id":"n1","title":"a"},{"id":"n7","title":"b"}]}`, "n8"},
		{"missing seq", `{"notes":[{"id":"n3","title":"a"}]}`, "n4"},
		{"seq ahead", `{"seq":9,"notes":[{"id":"n2","title":"a"}]}`, "n10"},
		{"odd ids ignored", `{"seq":2,"notes":[{"id":"x99","title":"a"},{"id":"n","title":"b"},{"id":"n-5","title":"c"}]}`, "n3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, fileName), []byte(tc.file), 0o600); err != nil {
				t.Fatal(err)
			}
			s, err := Load(dir, slog.New(slog.DiscardHandler))
			if err != nil {
				t.Fatal(err)
			}
			if err := s.Add("new"); err != nil {
				t.Fatal(err)
			}
			if l := s.List(); l[len(l)-1].ID != tc.next {
				t.Fatalf("new id = %s, want %s", l[len(l)-1].ID, tc.next)
			}
		})
	}
}

func TestFailedSaveChangesNothing(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Add("kept"); err != nil {
		t.Fatal(err)
	}
	s.path = filepath.Join(dir, "missing-dir", fileName) // CreateTemp fails
	if err := s.Add("lost"); err != errNotSave {
		t.Fatalf("err %v", err)
	}
	if ids(s) != "n1" || s.seq != 1 {
		t.Fatalf("state changed on failed save: %s seq=%d", ids(s), s.seq)
	}
}
