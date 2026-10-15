// Command server is the headless server (Playwright, scripts). For the TUI launcher see cmd/tabletop.
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"projectdnd/internal/server"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "listen address (0.0.0.0 = reachable from LAN; 127.0.0.1 = this machine only)")
	uploads := flag.String("uploads", "uploads", "directory for uploaded images")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := server.Run(ctx, *addr, *uploads, log); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
