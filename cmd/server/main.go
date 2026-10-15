package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"

	"projectdnd/internal/room"
	"projectdnd/internal/webui"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	uploads := flag.String("uploads", "uploads", "directory for uploaded images")
	flag.Parse()

	dist, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*uploads, 0o755); err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, newMux(room.NewHub(), *uploads, dist)))
}

func newMux(hub *room.Hub, uploads string, dist fs.FS) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.ServeWS)
	mux.HandleFunc("POST /api/upload", handleUpload(uploads, hub.Seated))
	mux.HandleFunc("GET /uploads/{name}", handleUploads(uploads))
	mux.Handle("/", http.FileServerFS(dist))
	return mux
}
