package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func pngBytes(t *testing.T) []byte {
	var b bytes.Buffer
	if err := png.Encode(&b, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func multipartBody(field string, data []byte) (*bytes.Buffer, string) {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, _ := mw.CreateFormFile(field, "../../etc/evil.png")
	fw.Write(data)
	mw.Close()
	return &b, mw.FormDataContentType()
}

func TestUpload(t *testing.T) {
	dir := t.TempDir()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/upload", handleUpload(dir, func(id string) bool { return id == "seated" }))
	mux.HandleFunc("GET /uploads/{name}", handleUploads(dir))
	mux.Handle("/", http.NotFoundHandler())

	post := func(client, field string, data []byte) *httptest.ResponseRecorder {
		body, ct := multipartBody(field, data)
		req := httptest.NewRequest("POST", "/api/upload", body)
		req.Header.Set("Content-Type", ct)
		if client != "" {
			req.Header.Set("X-Client-ID", client)
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}

	big := append(pngBytes(t), make([]byte, maxUpload)...)
	tests := []struct {
		name, client, field string
		data                []byte
		code                int
	}{
		{"no header", "", "file", pngBytes(t), 401},
		{"not seated", "stranger", "file", pngBytes(t), 401},
		{"text", "seated", "file", []byte("hello <script>"), 415},
		{"html", "seated", "file", []byte("<!DOCTYPE html><html>"), 415},
		{"wrong field", "seated", "image", pngBytes(t), 400},
		{"too big", "seated", "file", big, 413},
		{"png", "seated", "file", pngBytes(t), 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := post(tt.client, tt.field, tt.data)
			if rec.Code != tt.code {
				t.Fatalf("code %d, want %d: %s", rec.Code, tt.code, rec.Body)
			}
			if tt.code != 200 {
				return
			}
			var res struct{ URL string }
			json.Unmarshal(rec.Body.Bytes(), &res)
			name := filepath.Base(res.URL)
			if !uploadName.MatchString(name) || res.URL != "/uploads/"+name {
				t.Fatalf("url %q", res.URL)
			}
			get := httptest.NewRecorder()
			mux.ServeHTTP(get, httptest.NewRequest("GET", res.URL, nil))
			if get.Code != 200 || !bytes.Equal(get.Body.Bytes(), pngBytes(t)) || get.Header().Get("Content-Type") != "image/png" {
				t.Fatalf("get %d %q", get.Code, get.Header().Get("Content-Type"))
			}
		})
	}
	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Fatalf("%d files in uploads, want 1 (rejected uploads must not stay)", len(files))
	}

	for _, p := range []string{"/uploads/", "/uploads/x.png", "/uploads/..%2fgo.mod"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", p, nil))
		if rec.Code != 404 {
			t.Errorf("GET %s = %d, want 404", p, rec.Code)
		}
	}
}
