# project-D-D

D&D board game simulator for the web — play D&D with friends in the browser, no paper.

## Status
Planning. Requirements not confirmed yet; no app code.

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
