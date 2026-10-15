# project-D-D — D&D board game simulator (web)

Website to play D&D with friends online, no paper: shared map/board, tokens, dice,
character sheets. Single project at repo root (not a monorepo).

## Goal
Replace the paper table: a DM and friends open one link and play together in the browser.

## Confirmed requirements
- [ ] Not confirmed yet — run the pipeline gate below before writing any app code

## Initial scope
Core only, smallest thing that plays. No extras until the core works.

## Architecture
<decide after research + requirement gate>

## Working rules
- Clarify requirements before acting. Never jump into code on your own.
- Start from the smallest scope that works. No extras yet.
- Every code/behavior change bumps `version` in `package.json` (SemVer `MAJOR.MINOR.PATCH`). Docs-only changes do not.
- `README.md` is for humans; `CLAUDE.md` is instructions for Claude.
- Log work by date in `WORKLOG.md`.
- Game content: only use D&D 5e SRD (CC-BY-4.0) material. Never copy non-SRD book text.

## Pipeline for a new feature (checklist)
Two human gates. Do not skip them; this cannot run autonomously.
1. **Research** — unknown rules, libraries, realtime/hosting options: agent `researcher` (report + citations).
2. **GATE: confirm requirements** — agree smallest scope + architecture with the user. No code before this.
3. **Scaffold** — `/new-feature <name>` writes the spec stub and stops for the gate.
4. **Code** — smallest scope that works.
5. **Test** — run it in the browser, ideally with 2 tabs as 2 players.
6. **Review** — `ecc:typescript-reviewer` / `ecc:react-reviewer`; `ecc:security-reviewer` for anything multiplayer/network.
7. **Record** — update `WORKLOG.md` and the requirement checklist above.

## Slash commands (`.claude/commands/`)
| command | does |
|---|---|
| `/new-feature <name>` | write `docs/features/<name>.md` stub, then stop for the requirement gate |

## Agents (`.claude/agents/`)
| Agent | use when |
|---|---|
| `researcher` | need outside facts before deciding: D&D 5e SRD rules, web libs, realtime sync, hosting. Cited report |

## Tooling
- **Graft/Trail** (`@trailhq/trail-code`, bin `trail` at `~/.local/bin`): code graph over MCP to cut tokens when navigating.
  Repo-local only (`.claude/` + `.mcp.json`, cache `graft/` gitignored). Remove: `trail uninstall -y`.
  MCP loads when Claude Code launches from this folder; restart the session after install.
