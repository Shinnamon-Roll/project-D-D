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
- v0.5.0 seat QR join (spec `docs/features/seat-qr.md`): DM/co-DM "QR / สแกนเข้านั่ง" on empty player seats → centered modal with "ที่นั่ง N" label; `/?seat=<id>` preselects seat + name → `seat.sit`. Default bind now `0.0.0.0:8080` (user accepted LAN exposure, no password); `GET /api/lan` lists private IPv4 URLs (192.168 first); TUI shows LAN links. Dep: `qrcode`. Workflow: backend + frontend parallel → auditor (1 HIGH version, 2 MEDIUM IP choice + README security note, 6 LOW) → fixes.
- v0.6.0 TUI dashboard (SERVER/ROOM · LOGS · CONNECTIONS/ACTIVITY/SYSTEM, mint on dark, responsive 130/90 cols) + reload (`r` / click ⟳ RELOAD: npm build + make tabletop → graceful shutdown → re-exec with `-autostart`; confirm when clients connected; build fail keeps server). State carries `build` id; browser reloads once when it changes. Room summary accessor (names only). Audit: 2 MEDIUM (version, worklog) + LOW fixed (process-group kill, repo root from executable + go.mod check, short disabled label). 127 Go tests -race.
- v0.6.1 perf: player character screen jank = paint, not JS (0 WS msgs idle). Runtime SVG feTurbulence textures + `background-attachment: fixed` re-rasterized on every scroll/keystroke. Pre-rendered to `web/src/assets/tex-*.webp`, fixed bg moved to `body::before`. Scroll raster 3368ms → 17ms (3s, headless Chromium, CPU 4x); typing raster 191 → 31ms; look unchanged (<1/255 diff).
