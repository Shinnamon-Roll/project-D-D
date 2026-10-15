// Package jsonfile loads and atomically saves one JSON file (notes.json, npcs.json).
package jsonfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CorruptError means the file was not valid JSON for v and was moved aside to MovedTo.
type CorruptError struct {
	MovedTo string
	Err     error // the decode error
}

func (e *CorruptError) Error() string {
	return "corrupt file moved to " + e.MovedTo + ": " + e.Err.Error()
}

// Load decodes path into v. A missing file is not an error and leaves v unchanged. A file that does not decode
// is renamed to a unique path.bad-<random> (never overwritten) and Load returns a *CorruptError; the caller
// should then start from an empty value (v may be partly filled). Any other read error is returned as is.
func Load(path string, v any) error {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if derr := json.Unmarshal(b, v); derr != nil {
		bad, err := quarantine(path)
		if err != nil {
			return fmt.Errorf("%s is corrupt and could not be moved aside: %w", filepath.Base(path), err)
		}
		return &CorruptError{MovedTo: bad, Err: derr}
	}
	return nil
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

// Save writes v atomically: temp file in the same dir (0600) + fsync + rename, then fsync the dir.
func Save(path string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
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
	if err := os.Rename(tmp.Name(), path); err != nil {
		return err
	}
	if d, err := os.Open(dir); err == nil {
		d.Sync() // makes the rename durable; best effort (not supported everywhere)
		d.Close()
	}
	return nil
}
