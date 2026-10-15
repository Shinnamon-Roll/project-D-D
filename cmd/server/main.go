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
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
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
