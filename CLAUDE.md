# project-D-D — D&D board game simulator (web)

Website to play D&D with friends online, no paper: shared map/board, tokens, dice,
character sheets. Single project at repo root (not a monorepo).

## Goal
Replace the paper table: a DM and friends open one link and play together in the browser.

## Confirmed requirements
- White label: DM defines stats, skills, items (no hardcoded D&D rules).
- Private use. One DM (first to click) + optional co-DM; up to 6 player seats added by DM.
- Flow: lobby + ready → character creation → map with grid, click to move. Full draft: `docs/design/web-flow.md`.
- First campaign: Salem witch hunt, witch is an NPC (`docs/campaigns/salem/`).
- Open items: `docs/open-questions.md`. Confirmed list: `docs/decisions.md`.

## Initial scope
Core only, smallest thing that plays. No extras until the core works.

## Architecture
<decide after open questions are answered>

## Working rules
- Clarify requirements before acting. Never jump into code on your own.
- Start from the smallest scope that works. No extras yet.
- Every code/behavior change bumps `version` in `package.json` (SemVer `MAJOR.MINOR.PATCH`). Docs-only changes do not.
- `README.md` is for humans; `CLAUDE.md` is instructions for Claude.
- Log work by date in `WORKLOG.md`.
- Keep everything in files under `docs/` (index: `docs/README.md`), never only in chat:
  every user request/feedback → `docs/conversations/YYYY-MM-DD-<topic>.md` (verbatim + summary);
  confirmed → `docs/decisions.md`; unclear → `docs/open-questions.md`; out of scope → `docs/ideas.md`.
- Docs for the user are written in Thai (user's language); keep technical terms in English.
- Game content: only use D&D 5e SRD (CC-BY-4.0) material. Never copy non-SRD book text,
  Basic/Free Rules text, official character sheet layout, or adventure maps/text. Recommended: SRD 5.2.1 (not confirmed yet).
- Never brand the app "D&D" / "Dungeons & Dragons" (trademark). "5E compatible" is OK.
- Salem campaign: mechanics ideas only from Salem 1692 (Facade Games). Never copy its card/rule text, art, coined terms
  ("Tryal", "Town Crier", "frame card") or use "Salem 1692" as a title. NPCs fictional; no real 1692 person as the witch.

## Knowledge (read before deciding scope or rules)
| file | covers |
|---|---|
| [`docs/research/dnd-basics.md`](docs/research/dnd-basics.md) | what D&D is, play loop, core rules (d20, combat, grid, character), SRD versions + attribution, rules APIs |
| [`docs/research/salem-1692.md`](docs/research/salem-1692.md) | Salem 1692 board game rules, 1692 history, what can't be copied, free map sources |
| [`docs/research/vtt-and-player-needs.md`](docs/research/vtt-and-player-needs.md) | existing VTTs compared, player/DM needs, MVP must/nice/skip, Owlbear Rodeo lessons, OSS libs, asset licenses |

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
