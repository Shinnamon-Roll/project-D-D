package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"projectdnd/internal/webui"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	flag.Parse()

	dist, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWS)
	mux.Handle("/", http.FileServerFS(dist))

	log.Printf("listening on http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

type message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
}

// handleWS echoes every message back; replaced by the room hub in step 3.
func handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer c.CloseNow()
	c.SetReadLimit(64 << 10)

	ctx := context.Background()
	for {
		var m message
		if err := wsjson.Read(ctx, c, &m); err != nil {
			return
		}
		if err := wsjson.Write(ctx, c, message{Type: "echo", Payload: m}); err != nil {
			return
		}
	}
}
