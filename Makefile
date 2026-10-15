.PHONY: dev server web build test run tabletop

# Run Go server (:8080) and Vite dev server (:5173) together; open http://localhost:5173
dev:
	$(MAKE) -j2 server web

server: internal/webui/dist/.keep
	go run ./cmd/server

web:
	cd web && npm run dev

VERSION = $(shell node -p "require('./web/package.json').version")

# TUI launcher: welcome screen → start button → live server logs (http://127.0.0.1:8080).
run:
	cd web && npm run build
	$(MAKE) tabletop
	./bin/tabletop

tabletop: internal/webui/dist/.keep
	go build -ldflags "-X main.version=$(VERSION)" -o bin/tabletop ./cmd/tabletop

build:
	cd web && npm run build
	go build -o bin/dnd ./cmd/server
	$(MAKE) tabletop

# Not ./... : web/node_modules may contain stray .go files.
test:
	go test ./cmd/... ./internal/...

# go:embed needs at least one file in dist before the first frontend build.
internal/webui/dist/.keep:
	mkdir -p internal/webui/dist && touch $@
