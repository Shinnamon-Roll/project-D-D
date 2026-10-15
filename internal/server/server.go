// Package server wires the WebSocket hub, uploads and the embedded web app into one HTTP server.
package server

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"projectdnd/internal/room"
	"projectdnd/internal/webui"
)

const shutdownTimeout = 5 * time.Second

// Server is a running HTTP server. Create it with Start.
type Server struct {
	hub  *room.Hub
	http *http.Server
	ln   net.Listener
	log  *slog.Logger
	done chan error // result of Serve
}

// Start binds addr and serves in the background. Binding happens before Start returns,
// so a busy port is reported here, not later.
func Start(addr, uploads string, log *slog.Logger) (*Server, error) {
	dist, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(uploads, 0o755); err != nil {
		return nil, err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	hub := room.NewHub(log)
	base, cancel := context.WithCancel(context.Background())
	s := &Server{
		hub: hub,
		http: &http.Server{
			Handler:           newMux(hub, uploads, dist, log),
			ReadHeaderTimeout: 10 * time.Second,
			ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
			// WebSocket connections are hijacked, so Shutdown does not wait for them;
			// cancelling their base context makes their read loops return.
			BaseContext: func(net.Listener) context.Context { return base },
		},
		ln:   ln,
		log:  log,
		done: make(chan error, 1),
	}
	s.http.RegisterOnShutdown(cancel)
	go func() { s.done <- s.http.Serve(ln) }()
	log.Info("server started", "addr", s.Addr())
	return s, nil
}

// Addr is the bound address (useful with port 0).
func (s *Server) Addr() string { return s.ln.Addr().String() }

// Clients is the number of open WebSocket connections.
func (s *Server) Clients() int { return s.hub.Clients() }

// Shutdown stops accepting connections and waits for in-flight requests until ctx ends.
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.http.Shutdown(ctx)
	s.log.Info("server stopped")
	return err
}

// Run serves addr until ctx is cancelled, then shuts down gracefully.
func Run(ctx context.Context, addr, uploads string, log *slog.Logger) error {
	s, err := Start(addr, uploads, log)
	if err != nil {
		return err
	}
	select {
	case err := <-s.done:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
	}
	sctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.Shutdown(sctx)
}

func newMux(hub *room.Hub, uploads string, dist fs.FS, log *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.ServeWS)
	mux.HandleFunc("POST /api/upload", handleUpload(uploads, hub.Seated, log))
	mux.HandleFunc("GET /uploads/{name}", handleUploads(uploads))
	mux.Handle("/", http.FileServerFS(dist))
	return mux
}
