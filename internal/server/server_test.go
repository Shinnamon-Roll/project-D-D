package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestRunBindError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), ln.Addr().String(), t.TempDir(), slog.New(slog.DiscardHandler))
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Run on a taken port returned nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not report the bind error")
	}
}

func TestStartServesAndShutsDown(t *testing.T) {
	s, err := Start("127.0.0.1:0", t.TempDir(), slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.Get("http://" + s.Addr() + "/uploads/nope")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d", res.StatusCode)
	}
	if s.Clients() != 0 {
		t.Fatal("clients")
	}
	if err := s.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := http.Get("http://" + s.Addr() + "/"); err == nil {
		t.Fatal("still serving after Shutdown")
	}
}

func TestRunStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, "127.0.0.1:0", t.TempDir(), slog.New(slog.DiscardHandler)) }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(shutdownTimeout + time.Second):
		t.Fatal("Run did not return after cancel")
	}
}
