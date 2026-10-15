# v0-table — first playable version

Project rules and pipeline live in the root `CLAUDE.md`. Flow: `docs/design/web-flow.md` · protocol: `docs/design/protocol.md`.

## Goal
Play the Starter Set (Intro to Stormwreck Isle) online with friends: campaign menu → lobby → character creation → map.

## Confirmed requirements
- [x] Campaign menu (Starter Set playable, Salem "Coming soon"); first to pick = DM
- [x] Lobby: DM adds/removes player seats (≤6), opens 1 co-DM seat; players sit, ready; DM starts
- [x] White-label ruleset editor + "D&D 5e (SRD 5.2.1)" preset with attribution
- [x] Character creation: name, portrait upload, starting skills, stat points, starting items; DM grants anytime
- [x] Scenes: DM uploads map, calibrates grid, selects scene; tokens click-to-move; DM fog rectangles
- [ ] ~~Dice notation + shared log~~ — removed 2026-10-16, physical dice ([conv](../conversations/2026-10-16-remove-dice.md))
- [x] Rejoin by `clientId` in browser storage (Q13 room number skipped)

## Design
Go server-authoritative (`internal/room` `Apply` checks role + status per event), full-state broadcast per connection,
React + react-konva client. Details: `docs/design/tech-stack.md`.

## How to test
- `make test` — Go unit tests (rules, room permissions/flow, upload, hub)
- `cd web && npm run test:e2e` — Playwright: DM + player contexts run the whole flow
- Manual: `make dev` → open http://localhost:5173 in a normal and a private window (DM + player)

## Known limits (accepted)
- State lives in memory: restarting the server clears the game (uploads stay on disk).
- Fog rects and fogged tokens reach player browsers; the UI hides them, devtools could reveal them.
- Bound to 127.0.0.1 (localhost only). LAN: `./bin/dnd -addr 0.0.0.0:8080`.
- Audit LOW (not fixed): no HTTP server timeouts (matters once exposed beyond localhost); test gaps for player-moves-other-player-token / co-DM moves.
