---
title: Rules v2 data pipeline and class data delivery
author: "-Shinnamon-"
date: 2026-10-16
branch: main
cost: { minutes: 236, tokens: 1136157 }
touches:
  - docs/design/protocol.md
  - CLAUDE.md
  - README.md
  - WORKLOG.md
  - docs/README.md
  - docs/design/classes/cleric.md
  - docs/design/classes/paladin.md
  - docs/design/classes/ranger.md
  - docs/design/classes/sorcerer.md
  - docs/design/roster-legend.md
---

## Decided
- Rules v2 class data lives in docs/design/classes/<id>.md JSON blocks (+ _shared.md tables) and heroes in docs/design/heroes-core.md / heroes-legend.md; Go data is generated (`go test ./internal/rules -run MatchDoc -update`) and sync tests keep docs and Go equal.
- Class data (~136KB) is static and served once via GET /api/classes (ETag), never inside RoomState broadcasts.
- SRD 5.2.1 (CC BY 4.0) feature names are allowed with the exact attribution on the pick screen + README; descriptions stay in our own words.

## Tried and ruled out
- Putting class tables in RoomState: 136KB per broadcast.

## Watch out
- Edit class/hero data in docs, then regenerate; editing the generated Go files directly breaks the sync tests.
- Hidden NPC tokens: visible names are numbered from visible tokens only; never let a visible #n reveal a hidden token.
