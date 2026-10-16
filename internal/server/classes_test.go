package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"projectdnd/internal/rules"
)

func TestClassesEndpoint(t *testing.T) {
	h := handleClasses()
	res := httptest.NewRecorder()
	h(res, httptest.NewRequest(http.MethodGet, "/api/classes", nil))
	if res.Code != http.StatusOK || res.Header().Get("Content-Type") != "application/json" ||
		res.Header().Get("Cache-Control") != "no-cache" || res.Header().Get("ETag") == "" {
		t.Fatalf("status %d headers %v", res.Code, res.Header())
	}
	var d classesDoc
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&d); err != nil {
		t.Fatal(err)
	}
	if len(d.Classes) != 12 || len(d.Skills) != 18 || d.Attribution != SRDAttribution {
		t.Fatalf("classes %d skills %d attribution %q", len(d.Classes), len(d.Skills), d.Attribution)
	}
	for _, tc := range []struct {
		level, pb int
		full      [9]int
		half      [9]int
		pact      rules.PactSlots
	}{
		{1, 2, [9]int{2}, [9]int{2}, rules.PactSlots{Count: 1, Level: 1}},
		{5, 3, [9]int{4, 3, 2}, [9]int{4, 2}, rules.PactSlots{Count: 2, Level: 3}},
		{20, 6, [9]int{4, 3, 3, 3, 3, 2, 2, 1, 1}, [9]int{4, 3, 3, 3, 2}, rules.PactSlots{Count: 4, Level: 5}},
	} {
		i := tc.level - 1
		if len(d.Proficiency) != 20 || d.Proficiency[i] != tc.pb || d.FullSlots[i] != tc.full || d.HalfSlots[i] != tc.half || d.PactSlots[i] != tc.pact {
			t.Fatalf("level %d: pb %d full %v half %v pact %v", tc.level, d.Proficiency[i], d.FullSlots[i], d.HalfSlots[i], d.PactSlots[i])
		}
	}

	// Revalidation with the ETag: 304, no body.
	req := httptest.NewRequest(http.MethodGet, "/api/classes", nil)
	req.Header.Set("If-None-Match", res.Header().Get("ETag"))
	again := httptest.NewRecorder()
	h(again, req)
	if again.Code != http.StatusNotModified || again.Body.Len() != 0 {
		t.Fatalf("revalidate: %d, %d bytes", again.Code, again.Body.Len())
	}
}

func TestClassesRoute(t *testing.T) {
	s, err := Start("127.0.0.1:0", t.TempDir(), t.TempDir(), slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown(t.Context())
	res, err := http.Get("http://" + s.Addr() + "/api/classes")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("status %d type %q", res.StatusCode, res.Header.Get("Content-Type"))
	}
}
