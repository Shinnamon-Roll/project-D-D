---
description: write a feature spec stub under docs/features/, then stop to confirm requirements before code
argument-hint: <feature-name-kebab-case>
allowed-tools: Bash, Read, Write, Edit
---

Scaffold the spec for feature `$ARGUMENTS` following this repo's rules.

**Do not write app code in this step** — this command writes docs and sets the gate only.

1. Check the name: kebab-case, and `docs/features/$ARGUMENTS.md` must not exist yet. If empty or taken, ask again.

2. Create `docs/features/$ARGUMENTS.md` from this template:
   ```md
   # $ARGUMENTS

   Project rules and pipeline live in the root `CLAUDE.md`.

   ## Goal
   <summarize after the requirement talk>

   ## Confirmed requirements
   - [ ] not confirmed yet

   ## Initial scope
   Core only. No extras.

   ## Design
   <decide after research>

   ## How to test
   <browser steps, e.g. 2 tabs as DM + player>
   ```

3. Append to `WORKLOG.md` under today's date: `- Started feature \`$ARGUMENTS\` (spec stub).`

4. **STOP — requirement gate**: per repo rules ("clarify requirements before acting, never jump into code"),
   ask the user before going on:
   - What does this feature do? What is the smallest version that is playable?
   - Who sees/controls it: DM only, every player, or both? Must it sync live between players?
   - Does anything need research first (rules, library, realtime)? If yes, use agent `researcher`.

   Write code only after the user confirms requirements + design.
