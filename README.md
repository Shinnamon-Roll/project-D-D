# project-D-D

D&D board game simulator for the web — play D&D with friends in the browser, no paper.

## Status
v0.2 playable locally: campaign menu, lobby, ruleset editor, character creation, map with grid/tokens/fog, dice.

## Run
Requires Go 1.26+ and Node 24+.
```sh
cd web && npm install && cd ..
make dev            # Go :8080 + Vite :5173 → open http://localhost:5173
make build          # → bin/dnd (single binary, frontend embedded)
./bin/dnd           # http://127.0.0.1:8080  (LAN: ./bin/dnd -addr 0.0.0.0:8080)
make test           # Go tests
cd web && npm run test:e2e   # Playwright (DM + player)
```
Game state is in memory: restarting the server clears it. Uploaded images stay in `uploads/`.

## Structure
```
CLAUDE.md             rules + pipeline (instructions for Claude Code)
README.md             this file (for humans)
WORKLOG.md            work log by date
docs/README.md        index of all docs (conversations, design, campaigns, decisions, questions, ideas)
docs/research/        background knowledge: D&D rules basics, VTTs + player needs
docs/features/        one spec per feature (from /new-feature)
.claude/agents/       researcher (cited research)
.claude/commands/     /new-feature
.trail/learnings/     session learnings (committed)
graft/                local code-map cache (gitignored)
```

## Rules
- Every code change bumps the SemVer version.
- Only D&D 5e SRD (CC-BY-4.0) content.

## Tools
- **Graft/Trail** (`trail` CLI): code graph over MCP to cut tokens when navigating (repo-local, cache `graft/` not committed).
