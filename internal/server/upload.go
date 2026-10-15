package server

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const maxUpload = 10 << 20

var imageExt = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// uploadName matches only names this server generates, so no user input reaches the filesystem path.
var uploadName = regexp.MustCompile(`^[0-9a-f]{32}\.(png|jpg|gif|webp)$`)

func jsonError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// handleUpload serves POST /api/upload: a seated client sends multipart field "file" (image, <= 10 MB).
func handleUpload(dir string, seated func(clientID string) bool, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if id := r.Header.Get("X-Client-ID"); id == "" || !seated(id) {
			jsonError(w, http.StatusUnauthorized, "not seated")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxUpload+64<<10) // room for multipart headers
		mr, err := r.MultipartReader()
		if err != nil {
			jsonError(w, http.StatusBadRequest, "expected multipart/form-data")
			return
		}
		var part io.Reader
		for {
			p, err := mr.NextPart()
			if err != nil {
				if tooLarge(err) {
					jsonError(w, http.StatusRequestEntityTooLarge, "file larger than 10 MB")
				} else {
					jsonError(w, http.StatusBadRequest, "missing field \"file\"")
				}
				return
			}
			if p.FormName() == "file" {
				part = p
				break
			}
		}
		head := make([]byte, 512)
		n, err := io.ReadFull(part, head)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			if tooLarge(err) {
				jsonError(w, http.StatusRequestEntityTooLarge, "file larger than 10 MB")
			} else {
				jsonError(w, http.StatusBadRequest, "empty file")
			}
			return
		}
		head = head[:n]
		ext, ok := imageExt[http.DetectContentType(head)]
		if !ok {
			jsonError(w, http.StatusUnsupportedMediaType, "only png, jpeg, gif or webp images")
			return
		}
		name := randomHex() + ext
		path := filepath.Join(dir, name)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "cannot save file")
			return
		}
		written, err := io.Copy(f, io.LimitReader(io.MultiReader(bytes.NewReader(head), part), maxUpload+1))
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil || written > maxUpload {
			os.Remove(path)
			if written > maxUpload || tooLarge(err) {
				jsonError(w, http.StatusRequestEntityTooLarge, "file larger than 10 MB")
			} else {
				jsonError(w, http.StatusBadRequest, "upload failed")
			}
			return
		}
		log.Info("upload saved", "size", written, "type", http.DetectContentType(head))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"url": "/uploads/" + name})
	}
}

func tooLarge(err error) bool {
	var mbe *http.MaxBytesError
	return errors.As(err, &mbe)
}

func randomHex() string {
	b := make([]byte, 16)
	rand.Read(b) // never fails (crypto/rand panics internally on failure)
	return hex.EncodeToString(b)
}

// handleUploads serves GET /uploads/{name}: generated names only, no directory listing.
func handleUploads(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if !uploadName.MatchString(name) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		http.ServeFile(w, r, filepath.Join(dir, name))
	}
}
