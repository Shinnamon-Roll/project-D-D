.PHONY: dev server web build test

# Run Go server (:8080) and Vite dev server (:5173) together; open http://localhost:5173
dev:
	$(MAKE) -j2 server web

server: internal/webui/dist/.keep
	go run ./cmd/server

web:
	cd web && npm run dev

build:
	cd web && npm run build
	go build -o bin/dnd ./cmd/server

# Not ./... : web/node_modules may contain stray .go files.
test:
	go test ./cmd/... ./internal/...

# go:embed needs at least one file in dist before the first frontend build.
internal/webui/dist/.keep:
	mkdir -p internal/webui/dist && touch $@
