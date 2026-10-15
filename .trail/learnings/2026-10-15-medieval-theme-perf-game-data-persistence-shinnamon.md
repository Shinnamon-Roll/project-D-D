---
title: "Medieval theme perf, game data persistence and per-connection secrets"
author: "-Shinnamon-"
date: 2026-10-15
branch: main
cost: { minutes: 132, tokens: 478678 }
touches:
  - docs/conversations/2026-10-16-character-roster.md
  - docs/decisions.md
  - docs/features/character-roster.md
  - docs/conversations/2026-10-16-npc-pool.md
  - docs/design/roster.md
  - docs/features/npc-pool.md
  - docs/research/character-art.md
  - docs/research/img/character-art-contact-sheet.png
  - web/src/assets/classes/barbarian.svg
  - web/src/assets/classes/bard.svg
---

## Decided
- UI theme textures are pre-rendered WebP tiles (web/src/assets/tex-*.webp); the wood background is a fixed `body::before` layer.
- In-app dice were removed (physical dice). The `dice.roll` event is rejected as unknown.
- DM-only data (notes, and next the NPC pool and hidden NPC tokens/HP) is filtered per connection in hub.sendStateLocked.
  - Notes are sent only on hello, on note.* events, or when a connection becomes DM.
  - A state message without a `notes` key means "keep the previous notes".
- DM data that must survive TUI reload is persisted under `-data` (default `data/`, gitignored), via internal/notes: atomic temp+fsync+rename, 0600 file, corrupt-file quarantine via CreateTemp.
- Playwright runs the server with a temp `-data` dir so saved notes don't leak between runs.

## Tried and ruled out
- Runtime SVG feTurbulence textures plus `background-attachment: fixed`: 3.4s of raster per 3s of scroll on the character screen.
- Full-notes payload in every state broadcast: up to MBs per token move.

## Watch out
- Never reintroduce runtime SVG filters or fixed backgrounds (docs/design/visual-theme.md "Performance rules").
- Any new secret field must be stripped for players in sendStateLocked and in every broadcast path, and must never be logged.
- Parallel agents editing the same files: run sequential workflows for room/DmPanel/e2e changes.
- e2e v0.spec.ts has had timing flakes (Konva hit map redraw, token row re-render). Wait for state before acting.
