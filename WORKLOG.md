# WORKLOG — project-D-D

Work log by date (newest at bottom).

## 2026-10-16
- Set up repo like `../nothingPhone`: `CLAUDE.md`, `README.md`, `WORKLOG.md`, agent `researcher`, command `/new-feature`.
- Installed Graft/Trail repo-local: `trail init --no-agents --no-global` (no writes to `~/.claude`). Undo: `trail uninstall -y`.
- Next: requirement gate — what the smallest playable version is, and the stack.
- Research (agent `researcher` x2) saved to `docs/research/dnd-basics.md` and `docs/research/vtt-and-player-needs.md`; linked from `CLAUDE.md` Knowledge.
- Web flow draft (white label, lobby → character creation → map) + Salem campaign stub + docs structure (`docs/README.md`: conversations, decisions, open-questions, ideas).
- Salem research → `docs/research/salem-1692.md`; campaign design draft (proposal) → `docs/campaigns/salem/README.md`; Q12–Q17 open.
- Campaign menu (DM picks) added to flow; Salem paused. Starter Set research → `docs/research/starter-set.md`; plan → `docs/campaigns/starter-set/README.md` (DoSI, start with free intro).
- Tech stack confirmed: Go backend + React/TS frontend (`docs/design/tech-stack.md`). Open: hosting (Q12), room number (Q13).
- v0 step 1 scaffold: Go server (`cmd/server`, embed `internal/webui`), Vite React TS `web/` v0.1.0, Makefile (dev/build/test). Verified: binary serves page, WS echo, WS upgrade via Vite proxy (101).
- Agents added: `backend`, `frontend`, `auditor` (`.claude/agents/`).
- v0 built via workflow (backend + frontend agents in parallel → auditor → fixes → Playwright e2e). Audit: 3 MEDIUM fixed (fog leak in token list, dice notation raw length, scene name >32 chars), 2 LOW open. Verified: 66 Go tests (-race), vet, web build, e2e 1/1. Version 0.2.2. Spec: `docs/features/v0-table.md`.
- v0.3.0 TUI launcher (`make run`, Bubble Tea v2): welcome + start button + live slog logs; server moved to `internal/server` (ReadHeaderTimeout, graceful shutdown), `cmd/server` stays headless for Playwright. Audit: 1 MEDIUM fixed (shutdown on panic/SIGINT) + 3 LOW fixed by manager (control chars in log message, rune-safe truncation, Makefile prereq); 4 LOW open (q while starting, Serve error unread in TUI, logging under hub mutex, —). README rewritten with full usage. 84 Go tests, e2e pass.
- v0.4.0 medieval UI ("Candlelit Tavern", design brief `docs/design/visual-theme.md`): parchment panels on dark wood, Cinzel + Trirong, oxblood/gold buttons, original SVG art (crest, Stormwreck + Salem campaign cards, divider), wax-seal ready badge, candle connection status, Konva map colors. No Gemini key → SVG art instead of AI raster (ideas). Verified: web build, e2e 1/1, screenshots at 1440 + 390px.
