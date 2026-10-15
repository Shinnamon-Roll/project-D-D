package jsonfile

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type doc struct {
	N int `json:"n"`
}

func TestLoad(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content *string // nil = no file
		want    int
		corrupt bool
	}{
		{"missing file", nil, -1, false},
		{"valid", ptr(`{"n":3}`), 3, false},
		{"truncated", ptr(`{"n":`), 0, true},
		{"wrong type", ptr(`{"n":"x"}`), 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "f.json")
			if tc.content != nil {
				if err := os.WriteFile(path, []byte(*tc.content), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			v := doc{N: -1}
			err := Load(path, &v)
			var bad *CorruptError
			if tc.corrupt {
				if !errors.As(err, &bad) {
					t.Fatalf("err %v, want CorruptError", err)
				}
				if b, _ := os.ReadFile(bad.MovedTo); string(b) != *tc.content {
					t.Fatalf("quarantined content %q", b)
				}
				if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("corrupt file left in place")
				}
				return
			}
			if err != nil || v.N != tc.want {
				t.Fatalf("got %+v %v, want %d", v, err, tc.want)
			}
		})
	}
}

func ptr(s string) *string { return &s }

func TestSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.json")
	for _, n := range []int{1, 2} { // second save replaces the first
		if err := Save(path, doc{N: n}); err != nil {
			t.Fatal(err)
		}
	}
	var v doc
	if err := Load(path, &v); err != nil || v.N != 2 {
		t.Fatalf("got %+v %v", v, err)
	}
	if fi, err := os.Stat(path); err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode %v %v", fi.Mode().Perm(), err)
	}
	if m, _ := filepath.Glob(filepath.Join(dir, "*.tmp-*")); len(m) != 0 {
		t.Fatalf("temp files left: %v", m)
	}
	if err := Save(filepath.Join(dir, "missing", "f.json"), doc{}); err == nil {
		t.Fatal("save into a missing dir succeeded")
	}
}
