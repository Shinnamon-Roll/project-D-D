---
title: Repo setup mirroring nothingPhone + repo-local Graft
author: "-Shinnamon-"
date: 2026-10-15
branch: main
cost: { minutes: 3, tokens: 72621 }
touches:
  - README.md
  - CLAUDE.md
  - WORKLOG.md
---

## Decided
- `../working` does not exist; used `../nothingPhone` as the template (only sibling repo with Graft).
- Single project at repo root, not an `apps/` monorepo: one website.
- Kept the human gates: no app code until requirements + stack are confirmed. `/new-feature` writes `docs/features/<name>.md` and stops.
- Only D&D 5e SRD (CC-BY-4.0) content in the game.
- Docs in English (user wrote English); nothingPhone uses Thai.

## Tried and ruled out
- `trail init` without flags: would write `~/.claude/settings.json`, `~/.claude.json`, `~/.claude/helpers/` and affect all repos. Used `--no-agents --no-global`.

## Watch out
- `trail init --dry-run --no-global` still lists the global writes, but the real run with `--no-global` wrote none (verified with diff).
- Graft MCP loads only after restarting Claude Code from this folder.
