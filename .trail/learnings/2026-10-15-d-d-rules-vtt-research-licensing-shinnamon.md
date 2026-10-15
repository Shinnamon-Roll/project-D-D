---
title: "D&D rules + VTT research: licensing and MVP scope"
author: "-Shinnamon-"
date: 2026-10-15
branch: main
cost: { minutes: 30, tokens: 110823 }
touches:
  - CLAUDE.md
  - README.md
  - WORKLOG.md
  - ".playwright-mcp/console-2026-10-15T17-32-37-871Z.log"
  - ".playwright-mcp/console-2026-10-15T17-33-41-224Z.log"
  - ".playwright-mcp/console-2026-10-15T17-34-03-148Z.log"
  - ".playwright-mcp/console-2026-10-15T17-34-30-951Z.log"
  - ".playwright-mcp/console-2026-10-15T17-34-57-156Z.log"
  - ".playwright-mcp/console-2026-10-15T17-35-04-416Z.log"
  - ".playwright-mcp/console-2026-10-15T17-35-22-051Z.log"
---

## Decided
- Research lives in docs/research/dnd-basics.md and docs/research/vtt-and-player-needs.md (Thai, cited), linked from CLAUDE.md Knowledge.
- Recommended rules base: SRD 5.2.1 (2024 rules, CC-BY-4.0 only). Not yet confirmed by the user.
- MVP shape from research: Owlbear-Rodeo-like. Room link with no player sign-up, map + square grid, live token sync, GM rectangle fog, shared dice log, persistence. Skip sheets, rules automation, lighting, voice (friends use Discord).

## Tried and ruled out
- Bundling starter adventures (LMoP, Stormwreck, Delian Tomb), the official character sheet layout, or Basic/Free Rules text: not open content.
- tldraw: not open source (licence key). Owlbear legacy source: no OSS licence, study only.
- WebRTC for sync: Owlbear authors regretted it (VPN/ISP blocks, TURN cost).

## Watch out
- Open5e mixes OGL-only third-party content; always filter document__key=srd-2024.
- No SRD rule maps AoE shapes to grid squares (that is DMG); write our own.
- Don't brand the app "D&D" (trademark); "5E compatible" is OK. SRD attribution text is quoted in dnd-basics.md section 6.
