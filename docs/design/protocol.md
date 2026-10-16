# Protocol v0 — contract between Go server and React client

Source of truth for both sides. Go structs in `internal/` and `web/src/protocol.ts` must match this file
(JSON field names exactly as written). Change this file first, then both sides.

## Transport
- WebSocket `GET /ws`. Every message is JSON `{"type": string, "payload": object}`.
- Client's **first** message: `hello {clientId}` — `clientId` is a random UUID the browser keeps in `localStorage`.
  It is a bearer secret: the server **never** sends any clientId to other clients.
- Server → client:
  - `state {room: RoomState, you: {seatId: string}}` — sent to each connection after `hello` and after **every** change
    (`you.seatId` = "" when not seated; per-connection, so each client knows its own seat).
  - `error {message: string}` — only to the sender of a rejected event. State is unchanged.
  - DM/co-DM connections sometimes also get `notes: Note[]` at the top level of `state` (next to `room`, `you`,
    `build`); player and unseated connections **never** get the `notes` key. A `state` without the `notes` key means
    "keep the previous notes". See "DM notes" below.
  - Same rule for `npcs: Npc[]` (the NPC pool, see "NPC pool" below). **`room` is filtered per connection**: DM/co-DM
    get every token; players and unseated connections get no `hidden: true` tokens and no `hp`/`hpMax`/`npcId` on any token.
  - Combat (v0.14.0, see "Combat mode, attacks and DM log" below): DM/co-DM `state` also has `attacks: Attack[]` (every
    state) and sometimes `dmLog: LogEntry[]`; every other connection's `state` has `myAttack: {id, targetName} | null`.
  - `banner {kind: "hit" | "miss" | "talk" | "gift" | "surrender", text: L}` — ephemeral, to **every** connection once
    after `attack.resolve` or `npc.act`; not in state.
  - `log LogEntry` — one new DM log line, DM/co-DM connections only, sent after the `state` of the event that made it.
- Reconnect: same `clientId` → same seat (no room number in v0).

## RoomState
```jsonc
{
  "status": "menu" | "lobby" | "creating" | "playing",
  "campaigns": [{ "id": "starter-set", "name": "Starter Set: Dragons of Stormwreck Isle", "status": "available" },
                { "id": "salem", "name": "Salem", "status": "coming_soon" }],
  "campaignId": "",                       // "" in menu
  "seats": [{ "id": "s1", "kind": "dm" | "codm" | "player", "name": "", "occupied": false, "ready": false, "characterId": "" }],
  "sets": [{ "id": "core", "name": { "th": "ฮีโร่ชุดแรก", "en": "First Company" } },
           { "id": "legend", "name": { "th": "วีรชนในตำนาน", "en": "Heroes of Legend" } }], // fixed, every connection
  "activeSets": ["core"],                 // sets players pick from, in "sets" order; default ["core"]
  "roster": [ /* RosterEntry, see "Character roster" below; DM/co-DM: every set (22), players: active sets only */ ],
  "characters": [ /* Character, see "Character roster" below */ ],
  "scenes": [{ "id": "sc1", "name": "Beach", "imageUrl": "/uploads/y.png",
               "grid": { "cellSize": 50, "offsetX": 0, "offsetY": 0, "cols": 20, "rows": 15 },
               "fog": [{ "id": "f1", "col": 0, "row": 0, "w": 3, "h": 2 }] }],
  "currentSceneId": "",
  "tokens": [{ "id": "t1", "sceneId": "sc1", "characterId": "c1", "npcId": "", "name": "Mira", "imageUrl": "/uploads/x.png",
               "col": 0, "row": 0, "hidden": false },
             { "id": "t9", "sceneId": "sc1", "characterId": "", "npcId": "npc1", "name": "Goblin #2", "imageUrl": "",
               "col": 1, "row": 0, "hidden": true, "hp": 7, "hpMax": 7 }]
                                          // characterId "" = NPC token; hp/hpMax only on NPC tokens and only for DM/co-DM;
                                          // hidden tokens only for DM/co-DM
                                          // character tokens: name = heroName.en, imageUrl = character art ("" = class art)
  "feed": [{ "at": 1760000000000, "text": { "th": "...", "en": "..." } }],  // newest last, at most 50
  "combat": false                         // combat mode (combat.set); same for every connection
}
```

## Events (client → server)
Roles: **DM** = seat kind `dm` or `codm` (co-DM has every DM right except removing the `dm` seat). **Player** = seat kind `player`.
| type | payload | who | when | effect / validation |
|---|---|---|---|---|
| `hello` | `{clientId}` | anyone | any | bind connection; must be first |
| `campaign.select` | `{campaignId, name}` | anyone not seated | `menu` | campaign must be `available`; sender sits in a new `dm` seat; status → `lobby`; 1 empty `player` seat created |
| `seat.sit` | `{seatId, name}` | not seated | `lobby` | seat empty, kind `player` or `codm`; name 1–32 chars |
| `seat.leave` | `{}` | seated, not `dm` | `lobby` | frees seat, clears ready |
| `seat.ready` | `{ready}` | player | `lobby` | |
| `seat.add` | `{}` | DM | `lobby` | max 6 `player` seats |
| `seat.remove` | `{seatId}` | DM | `lobby` | not the `dm` seat; kicks occupant |
| `seat.openCoDm` | `{}` | DM | `lobby` | max 1 `codm` seat |
| `sets.set` | `{sets: string[]}` | DM | `lobby` | ids from `RoomState.sets`; empty/missing → `pick at least one set`, unknown id → `no such set`; duplicates dropped, stored in `sets` order; DM log `"set the character sets: First Company, Heroes of Legend"` |
| `game.start` | `{}` | DM | `lobby` | ≥1 occupied player seat, all occupied player seats ready; empty player seats are removed; → `creating` |
| `character.pick` | `{heroId}` (pre-v2 key `classId` still accepted when `heroId` is empty) | player without character | `creating` | entry in the roster, its `set` in `activeSets` (else `that set is not in play`), and not taken by another character (by hero id: the same class from another set is fine); creates the Character at level 1 from the current roster entry (see "Rules v2"); when every occupied player seat has a character → `playing` |
| `character.adjust` | `{characterId, hp}` | DM | `creating`/`playing` | hp delta int −999..999 (required); clamps to `0..hpMax`; knocked-out / back-up lines |
| `character.level` | `{characterId, delta}` | DM | `creating`/`playing` | delta `1` or `-1`; level stays 1–20; recomputes hpMax/maxima (see "Rules v2") |
| `character.setAbilities` | `{characterId, abilities, ac}` | DM | `creating`/`playing` | six scores 1–30, ac 1–30; recomputes hpMax (CON) and `abilityMod` resources |
| `character.setSkills` | `{characterId, skills, expertise, saves}` | DM | `creating`/`playing` | replaces all three (missing key = `[]`): skills/expertise as `roster.update` (known skill ids, no duplicates, expertise ⊆ skills); saves = **extra** save proficiencies from class features (Expertise, Deft Explorer, Scholar, Disciplined Survivor, Slippery Mind): ability ids `str`…`cha`, duplicates dropped, stored as `extraSaves` in sheet order; recomputes `saves`. DM log only |
| `slot.set` | `{characterId, level, used}` | owner or DM | `creating`/`playing` | level int 1–9 → `slotsUsed[level-1]`, or `"pact"` → `pactUsed`; `0 ≤ used ≤ max` (max 0 → `no such spell slot`) |
| `resource.set` | `{characterId, resourceId, used}` | owner or DM | `creating`/`playing` | a resource of the character's class with `resourcesMax[id] > 0` (else `no such resource`); `0 ≤ used ≤ max` |
| `hitDice.set` | `{characterId, used}` | owner or DM | `creating`/`playing` | `0 ≤ used ≤ hitDiceMax` (= level) |
| `party.rest` | `{kind}` | DM | `playing` | `"short"` or `"long"`, every character; see "Rules v2" |
| `ability.add` | `{name, kind, level, desc}` | owner (player with a character; acts on the sender's own character) | `creating`/`playing` | name trimmed 1–40 (one line), desc trimmed ≤ 300 (newlines ok), kind `spell`/`feature`, spell level int 0–9 (feature: level ignored, stored 0); at most 40 |
| `ability.update` | `{abilityId, name, kind, level, desc}` | owner (own character only; `characterId` ignored) | `creating`/`playing` | same limits; keeps its id |
| `ability.remove` | `{characterId, abilityId}` | owner or DM | `creating`/`playing` | |
| `item.give` | `{characterId, name, qty}` | DM | `creating`/`playing` | name trimmed 1–40, qty 1–99; a stack whose `name.th` or `name.en` equals name gets `+qty` (stack ≤ 999); else a new stack `{name: {th: name, en: name}, qty}` (≤ 50 stacks) |
| `item.remove` | `{characterId, index, qty}` | DM | `creating`/`playing` | index into `items`; qty ≥ 1; qty ≥ stack qty removes the stack |
| `roster.update` | `{classId, entry: RosterEntry}` | DM | any | classId is the key (`entry.id` is ignored); limits below; affects later picks only; persisted |
| `roster.reset` | `{classId?}` or `{set}` | DM | any | `{classId}` = that entry; `{}`/`""` = every entry of every set; `{set}` = only that set's entries (unknown → `no such set`; both keys → `send classId or set, not both`); persisted |
| `scene.add` | `{name, imageUrl, grid}` | DM | `creating`/`playing` | imageUrl must start with `/uploads/`; grid: cellSize 10–400, cols/rows 1–200 |
| `scene.update` | `{sceneId, name, grid}` | DM | `creating`/`playing` | same grid limits (cellSize is an int: send `round(imgW/cols)`; offsetX/offsetY 0). A grid with fewer than ceil(visibleTokensInScene/4) cells (hidden tokens never count) → `grid too small: <n> visible tokens need at least <m> cells`, nothing changes. Grid shrink: that scene's tokens clamp to the nearest edge cell `[0,cols-1]×[0,rows-1]`; a clamped visible token whose cell already holds 4 visible tokens goes to the nearest cell (Chebyshev distance, then row-major) with room; fog rects clip to the grid, rects fully outside are dropped |
| `scene.select` | `{sceneId}` | DM | `playing` | sets `currentSceneId`; adds a token for every character without one in that scene, one per cell row-major from 0,0, skipping cells with 4 visible tokens and wrapping around; no cell with room → `scene has no free cell` (name = `heroName.en`, imageUrl = character `art`) |
| `npc.place` | `{npcId, sceneId}` | DM | `playing` | copies the pool NPC into a new token on the first least-crowded cell (row-major) of the scene; `hp = hpMax`, `hidden`, `imageUrl` from the pool; name by "Token names" below; a visible NPC ignores hidden tokens for its cell (may stack on one); a hidden NPC counts every token for its cell and always fits; visible NPC with every cell holding 4 visible → `scene has no free cell`. Replaces `token.addNpc` (removed v0.11.0) |
| `token.move` | `{tokenId, col, row}` | DM any token; player only own character's token | `playing` | 0 ≤ col < cols, 0 ≤ row < rows of the token's scene. A player gets `you can only move your own token` for a missing token too (no hidden-token oracle). Target cell already holding 4 visible tokens (the mover not counted; hidden tokens never count) → `that cell is full (max 4)`; a hidden token may move anywhere |
| `token.setHidden` | `{tokenId, hidden}` | DM | `playing` | NPC tokens placed from the pool only. Reveal (`hidden: false`) on a cell already holding 4 visible tokens → `that cell is full (max 4)` (move it first); hiding is always allowed and keeps the name; a reveal may rename (see "Token names") |
| `token.duplicate` | `{tokenId}` | DM | `playing` | NPC tokens placed from the pool only (character → `only NPC tokens from the pool`, missing → `no such token`). New token: same `imageUrl`, `hpMax`, `hidden`, `npcId`; `hp = hpMax`; name by "Token names"; cell = nearest to the original (Chebyshev distance, then row-major; the original's own cell first) holding < 4 tokens — visible copy counts visible tokens only, hidden copy counts all and falls back to the original's cell; visible copy with no such cell → `scene has no free cell`. Line `"added Goblin #6"` / `"เพิ่ม Goblin #6"`: visible = feed + DM log; hidden = DM log only, `"added Goblin #6 (hidden)"` / `"เพิ่ม Goblin #6 (ซ่อน)"`. No combat mode needed |
| `token.adjustHp` | `{tokenId, delta}` | DM | `playing` | NPC tokens placed from the pool only; `-999 ≤ delta ≤ 999` (int); `hp` clamps to `0..hpMax` |
| `token.remove` | `{tokenId}` | DM | `playing` | |
| `fog.add` | `{sceneId, col, row, w, h}` | DM | `playing` | w,h ≥ 1, rect inside grid |
| `fog.remove` | `{sceneId, fogId}` | DM | `playing` | |
| `note.add` | `{title}` | DM | any | appends `{id, title, body: "", updatedAt}`; title trimmed, 1–80 chars; max 200 notes |
| `note.update` | `{id, title, body}` | DM | any | title 1–80 (trimmed), body ≤ 20000 chars (not trimmed); sets `updatedAt`; last write wins |
| `note.remove` | `{id}` | DM | any | |
| `note.move` | `{id, index}` | DM | any | `index` = position in the resulting list, `0 ≤ index < len(notes)` |
| `npc.add` | `{name, imageUrl, hpMax, hidden}` | DM | any | name trimmed 1–40 chars and not already in the pool (trimmed, case-insensitive) → else `that name is already in the pool`; `hpMax` 1–999; `imageUrl` "" or an existing `/uploads/<generated name>`; max 100 NPCs |
| `npc.update` | `{id, name, imageUrl, hpMax, hidden}` | DM | any | same checks (its own current name is allowed); tokens already placed keep their copy. An old `npcs.json` with duplicate names loads unchanged |
| `npc.remove` | `{id}` | DM | any | tokens already placed stay |
| `combat.set` | `{on}` | DM | `playing` | sets `combat`; same value again = no change, no line; `on: false` cancels every pending attack |
| `attack.declare` | `{targetTokenId}` | player with a character | `playing`, `combat` | at most 1 pending attack per character; attacker needs a token in `currentSceneId`; target exists and is not hidden (a hidden token gets the same `no such target` as a missing one), is not the attacker's own token, and is in the attacker's scene |
| `attack.cancel` | `{}` | player with a pending attack | `playing` | removes it |
| `attack.resolve` | `{attackId, hit, damage?}` | DM | `playing` | damage int 0–999, only with `hit: true` (miss + damage > 0 → error); hit + damage: NPC token `hp` or the target character's `hp` −= damage, clamp 0; the attack is removed (second resolve → `no such attack`); banner to everyone |
| `npc.act` | `{tokenId, targetId, kind, text?, hit?, damage?, item?: {name, qty}}` | DM | `playing`, `combat` | see "NPC actions" below |

Unknown `type`, malformed payload, or failed check → `error` to sender. Removed in v0.12.0 (now unknown types):
`ruleset.update`, `ruleset.loadPreset`, `character.submit`, `character.grant`. Removed in v0.18.0 (rules v2, now
`unknown event type "<type>"`): `skill.addCustom`, `skill.editCustom`, `skill.removeCustom`, `skill.use`,
`character.setStats`; `character.adjust` has no `mp` and `party.rest` needs `kind`.

"Owner or DM" = a DM/co-DM, or the player seated on that character. A player naming another or a missing
character gets `you can only change your own character` (never a hint which ids exist).

## NPC actions (`docs/features/npc-actions.md`)
- `npc.act`: the DM makes a visible NPC token act on a player character. `kind`: `"talk" | "attack" | "gift" | "surrender"`.
  - `tokenId`: an NPC token (`characterId` ""), in `currentSceneId`; hidden → `reveal the NPC first`.
  - `targetId`: a player-character token in the NPC's scene (missing, NPC or hidden token → `no such target`).
  - talk: `text` optional, trimmed, one line, ≤ 200 chars. attack: `hit` bool, `damage` int 0–999 (miss + damage > 0
    → error); hit + damage: target character `hp` −= damage, clamp 0 (+ knocked-out line). gift: `item.name` 1–40,
    `item.qty` 1–99, added to the target's items exactly like `item.give` (same stack merge, limits and errors).
    surrender: no fields. Fields of other kinds are ignored.
- Lines (feed + DM log + banner all use the full token name, `" #n"` included, e.g. `"Wolf #2 attacks Mira — hit! −5"`):
  talk `"Wolf talks with Mira: “Gold, now!”"` / `"Wolf คุยกับ Mira: “Gold, now!”"` (no text: `"Wolf talks with Mira"`);
  hit `"Wolf attacks Mira — hit! −5"` / `"Wolf โจมตี Mira — โดน! −5"` (damage 0: no number); miss
  `"Wolf attacks Mira — miss"` / `"Wolf โจมตี Mira — ไม่โดน"`; gift `"Wolf gives Mira Rusty Key ×1"` /
  `"Wolf มอบ Rusty Key ×1 ให้ Mira"`; surrender `"Wolf surrenders to Mira"` / `"Wolf ยอมแพ้ต่อ Mira"`.
- `banner`: every successful `npc.act`, to every connection. `kind` = `hit` / `miss` (attack) or the action kind.
  `text`: hit/miss `"Wolf → Mira"` (as `attack.resolve`); talk/gift/surrender = the feed line.
- Errors: `combat mode is off`, `no such NPC token`, `reveal the NPC first`, `the NPC is not in this scene`,
  `no such target`, `target is not in this scene`, `kind must be talk, attack, gift or surrender`,
  `talk text must be one line of at most 200 characters`, `damage must be 0-999`, `a miss deals no damage`,
  `item name must be 1-40 characters`, `item qty must be 1-99`, `an item stack holds at most 999`,
  `at most 50 different items`, plus the shared ones.

## Combat mode, attacks and DM log (v0.14.0, `docs/features/combat-log.md`)
```ts
type Attack   = { id: string; byCharacterId: string; byName: L; targetTokenId: string; targetName: string; at: number }
type LogEntry = { at: number; actor: string; text: L }   // at = unix ms; actor = sender's seat name, never a clientId
```
- Pending attacks are **not** in `RoomState`. DM/co-DM `state` carries `attacks: Attack[]` (`[]` when none) on every
  state; player and unseated `state` carries `myAttack: {id, targetName}` (own pending attack) or `null` — never `attacks`.
  Players also get a `state` whenever the pending-attack list changes (so their `myAttack` is current).
- `byName` = attacker hero name; `targetName` = target token name (character tokens: `heroName.en`). Attack ids `a<seq>`.
- `token.remove` drops pending attacks on that token. `combat.set {on: false}` drops all.
- Feed after `attack.resolve`: miss `"Aldric Vane → Goblin: miss"` / `"… : ไม่โดน"`; hit on an NPC `"… : hit!"` / `"… : โดน!"`
  (never numbers); hit on a character with damage `"Aldric Vane → Mira: hit! −3 (HP 9/12)"` / `"… : โดน! −3 (เลือด 9/12)"`.
  The DM log gets the NPC numbers: `"Aldric Vane → Goblin: hit, 3 damage (HP 4/7)"`; a target without HP:
  `"Aldric Vane → Goblin: hit, 3 damage (no HP)"` / `"… : โดน ดาเมจ 3 (ไม่มีเลือด)"`.
- Target **hidden** at resolve time (hidden after the declare): every result line goes to the DM log only — no feed
  line and **no `banner`** (players must not learn the result).
- `banner.text` = `"<attacker> → <target>"` (th: hero/character names in Thai, en in English).
- Death / revive (attack.resolve, character.adjust, token.adjustHp), when hp goes from > 0 to 0 or from 0 to > 0:
  character `"💀 <name> is knocked out (HP 0)"` / `"💀 <name> หมดสติ (HP 0)"`; NPC `"💀 <name> falls (HP 0)"` /
  `"💀 <name> ล้มลง (HP 0)"`; back up `"<name> is back up"` / `"<name> ฟื้นขึ้น"`. Feed + DM log, except for a **hidden**
  NPC token: DM log only.
- `combat.set` feed: `"Combat started"` / `"เข้าสู่โหมดต่อสู้"` (no ⚔: fonts render it as ×), `"Combat ended"` / `"จบการต่อสู้"`.
- DM log: in memory only (lost on restart), newest 300. Every successful event adds its line(s): every feed line, plus
  seat sit/leave/ready, seat add/remove, co-DM seat, sets.set, game start, custom abilities, character.setAbilities, character.setSkills,
  slot/resource/hit-dice changes, roster edit/reset, scene add/edit/select,
  NPC pool add/edit/remove, npc.place, hide/reveal, NPC hp, token move (`"moved Wolf from (0,0) to (4,3)"`, 0-based
  col,row), token remove, fog add/remove, attack declare/cancel/resolve, combat on/off (+ cancelled count). Notes only
  as `"edited note <title>"`, at most once per note per 60 s (never the body). Failed events add nothing.
- Token moves: moves of the **same token by the same seat** within 3 s of that line's first move share one line —
  the server rewrites the stored line in place (first `from`, latest `to`, `at` = latest move) and sends **no** `log`
  frame for it, so the client stays append-only and spamming moves cannot push older lines out of the 300. The live
  list may show the first destination until the next full `dmLog` (hello / becoming DM); the map shows the real cell.
- Presence: a **seated** client's first open connection (by clientId) adds `"joined the room"` / `"เข้าห้อง"`, its
  last closed connection `"left the room"` / `"ออกจากห้อง"` (actor = seat name). Extra tabs and unseated clients add
  nothing. Sent as a `log` frame to the other DM/co-DM connections (a joining DM gets it inside its hello `dmLog`).
- `state.dmLog` (full list, oldest first) goes to a DM/co-DM connection only in the reply to `hello` and in the first
  `state` after it becomes DM/co-DM; **no key = keep the previous list**. Each new line then comes as a `log` message
  after the event's `state`; a connection that just got the full `dmLog` does not get those lines again.
- Errors: `combat mode is off`, `you have no character`, `you already have a pending attack`,
  `your token is not in this scene`, `no such target`, `you cannot attack yourself`, `target is not in this scene`,
  `you have no pending attack`, `no such attack`, `damage must be 0-999`, `a miss deals no damage`,
  plus the shared `only the DM can do that`, `not allowed while the room is <status>`, `malformed payload`.

## Upload
`POST /api/upload`, header `X-Client-ID: <clientId>` of a **seated** client, multipart field `file`.
- ≤ 10 MB (`413`), content sniffed with `http.DetectContentType` ∈ png/jpeg/gif/webp (`415`), not seated (`401`).
- Saved as `uploads/<random hex>.<ext>`; response `200 {"url": "/uploads/<name>"}`; errors `{"error": "..."}`.
- `GET /uploads/<name>` serves files (no directory listing).

## Server bind
Default `-addr 0.0.0.0:8080` (reachable from other devices on the LAN; use `-addr 127.0.0.1:8080` for this machine only).
`GET /api/lan` → `{"urls": [...]}`: private IPv4 only (192.168/16, 10/8, 172.16/12; no link-local, CGNAT/Tailscale or public),
preferred first: 192.168 → 10 → 172.16 (Docker/VM bridges usually last). Loopback bind → `[]`.

## Server limits beyond the table (implemented v0.2)
- Names (DM, seat, scene): 1–32 chars. Pool NPC names: 1–40.
- `hello.clientId`: 8–64 chars `[A-Za-z0-9_-]`; any event before `hello`, or a second `hello` → error.
- Image URLs: exactly `/uploads/<name>` (no extra `/`, `..`, `?`, `#`; name ≤64).
- `scene.update` clamps tokens and clips/drops fog in a shrunken grid; `scene.select` wraps auto-placed tokens to the next row.
- At most 4 visible tokens per cell (`MaxPerCell`, v0.13.0, `docs/features/shared-cell.md`); hidden tokens never count, so the full-cell error never reveals one.
- A co-DM may remove their own seat; an empty co-DM seat survives `game.start`.
- Hub drops a client with >32 queued messages (it reconnects and gets fresh state).

## `GET /api/lan`
No auth (same as the page). Returns `200 {"urls": ["http://<ipv4>:<port>"]}` — every non-loopback IPv4 of the host with the listening port; `[]` when the server listens on loopback only. Used by the seat QR (`docs/features/seat-qr.md`).

## `build` in `state`
The `state` payload has a top-level `"build": "<16 hex>"`, random per server process. The client remembers the first value; when a later `state` carries a different value (server restarted, e.g. TUI reload), it calls `location.reload()` to load the new web build.

## DM notes (v0.10.0, `docs/features/dm-notes.md`)
- `Note = {"id": "n1", "title": "...", "body": "...", "updatedAt": 1760000000000}` (`updatedAt` = unix **milliseconds**).
  Ids are `n<seq>`, never reused (seq is persisted; on load it is raised to at least the highest `n<k>` id).
- Not part of `RoomState`. Sent only in `state.notes` to DM/co-DM connections (`[]` when empty), and only in:
  (a) the `state` reply to `hello`; (b) the broadcast after a `note.*` event; (c) the first `state` after a connection
  becomes DM/co-DM (e.g. `campaign.select`, sitting in the co-DM seat). Every other `state` omits the `notes` key:
  **`state` without `notes` = keep the previous notes**. A note event's broadcast goes to DM/co-DM connections only;
  players get nothing for it and never get the `notes` key.
- Errors: `only the DM can do that`, `note title must be 1-80 characters`, `note body must be at most 20000 characters`,
  `at most 200 notes`, `no such note`, `note index out of range`, `could not save notes` (disk write failed; nothing changed),
  `malformed payload`.
- Lengths count characters (runes). WebSocket read limit is 128 KB per message (fits a 20000-char body).
- Persisted to `<-data>/notes.json` (flag `-data`, default `data`, gitignored) on every change: temp file in the
  same dir + fsync + rename; dir `0700`, file `0600`. Loaded at server start. A file that is not valid JSON is renamed
  to a unique `notes.json.bad-<random>` and the server starts with no notes (error logged); an unreadable file stops startup.
- Logs name the event and sender only (`event note.update by=<name>`), never titles or bodies.

## Token names (`docs/features/npc-names.md`)
- Token names are unique per scene (case-insensitive); other scenes are independent. Character tokens are never
  renamed but count as taken.
- Base = name minus a trailing `" #<n>"` (`n` 1–999999999, no leading zero); a name without one is number 1.
  Name from source `src` (pool name for `npc.place`, the original token's name for `token.duplicate`):
  `n = max(src's number, highest number of that base among the counted tokens + 1)`; `n = 1` → the bare base, else
  `"<base> #<n>"`. E.g. `Goblin` → `Goblin #2`, `#3`; highest `Goblin #5` → `Goblin #6`; 5 goblins, duplicate `#2` →
  `Goblin #6`; pool `Goblin #5` in an empty scene → `Goblin #5`, then `#6`.
- Visible new token: counts visible tokens only. Hidden new token: counts every token (scene-wide next).
- A visible name that a hidden token already holds: the hidden token is renamed to its base's scene-wide next.
- Reveal (`token.setHidden false`): the token takes the visible next of its base (starting at its pool NPC's own number
  when that NPC still has the same base) if its name differs, then any clashing hidden token is renamed; pending
  attacks on it take the new `targetName`. Players never see a hidden rename (their view is unchanged, `myAttack` too).
- Every automatic rename adds a DM-log-only line `"renamed Goblin #2 to Goblin #3"` / `"เปลี่ยนชื่อ Goblin #2 เป็น Goblin #3"`.

## NPC pool (v0.11.0, `docs/features/npc-pool.md`)
- `Npc = {"id": "npc1", "name": "Goblin", "imageUrl": "", "hpMax": 7, "hidden": true}` — `hidden` is the default for
  new tokens placed from it; `imageUrl` "" = default silhouette. Ids are `npc<seq>`, never reused (same rule as notes).
- Not part of `RoomState`. Sent only in `state.npcs` to DM/co-DM connections (`[]` when empty), only in:
  (a) the `state` reply to `hello`; (b) the broadcast after an `npc.*` event (incl. `npc.place`); (c) the first `state`
  after a connection becomes DM/co-DM. **`state` without `npcs` = keep the previous pool.** Players never get the key.
- Players (and unseated connections) get a `state` only when their filtered view changed or they sent the event:
  `note.*`, `npc.add/update/remove`, placing a hidden NPC, and hp changes send them nothing.
- `imageUrl` check: exactly `/uploads/<32 hex>.<png|jpg|gif|webp>` (a name the upload endpoint generates) and the
  file exists in the uploads dir.
- Errors: `only the DM can do that`, `npc name must be 1-40 characters`, `that name is already in the pool`, `npc hpMax must be 1-999`, `at most 100 npcs`,
  `no such npc`, `imageUrl must be empty or an uploaded image`, `could not save npcs`, `scene has no free cell`,
  `only NPC tokens from the pool`, `delta must be -999 to 999`, `no such token`, `malformed payload`.
- Persisted to `<-data>/npcs.json` exactly like notes (atomic write, dir `0700`, file `0600`, corrupt file → `npcs.json.bad-<random>`).
- Logs name the event and sender only (`event npc.place by=<name>`), never NPC or token names.
- Token ids share the room counter (`t<seq>`): accepted: id gaps may hint that hidden tokens exist; names/hp/position never leak.

## Character roster (v0.12.0, `docs/features/character-roster.md`; sets v0.15.0, `docs/features/character-sets.md`)
```ts
type L = { th: string; en: string }
type Abilities  = { str: number; dex: number; con: number; int: number; wis: number; cha: number }  // 1-30
type RosterItem = { name: L; qty: number }
type RosterEntry = {     // a hero
  id: string            // unique across sets. core: knight|wizard|rogue|ranger|cleric|paladin|barbarian|bard|druid|necromancer
                        // legend: legend-<barbarian|bard|cleric|druid|fighter|monk|paladin|ranger|rogue|sorcerer|warlock|wizard>
  set: 'core' | 'legend' // always the built-in entry's set; a client-sent set is ignored
  classId: string       // Class id from /api/classes; always the built-in entry's class, a client-sent one is ignored
  heroName: L; blurb: L
  role: string          // tank|damage|healer|support|control|hybrid
  art: string           // "" = class art by hero id (web/src/assets/classes/<id>.svg) | "/uploads/<generated name>"
  abilities: Abilities; ac: number
  skills: string[]      // proficient skill ids (18 known ids, no duplicates)
  expertise: string[]   // subset of skills
  items: RosterItem[]   // 0-10
}
type CustomAbility = { id: string; name: string; kind: 'spell' | 'feature'; level: number; desc: string }
                        // player-typed, one language; id "ab<n>"; level = spell level 0-9, 0 for a feature
type Character = {
  id: string; seatId: string
  classId: string       // Class id (fighter, wizard, ...)
  heroId: string        // RosterEntry id picked (knight, legend-bard, ...)
  name: L               // heroName at pick time
  art: string           // entry art at pick time
  level: number         // 1-20
  hp: number; hpMax: number
  abilities: Abilities; ac: number
  skills: string[]; expertise: string[]
  extraSaves: string[]  // ability ids the DM added with character.setSkills (sheet order); [] at pick
  slotsUsed: number[]   // 9 entries, index 0 = 1st-level slots (full/half casters)
  pactUsed: number      // pact casters (Warlock)
  resourcesUsed: Record<string, number>   // every class resource id, 0 when unused
  hitDiceUsed: number
  abilitiesCustom: CustomAbility[]        // ≤ 40
  items: RosterItem[]
  // derived by the server from class + level + abilities (clients never compute these):
  pb: number            // proficiency bonus
  saves: string[]       // proficient saves = class.saves ∪ extraSaves, sheet order (str dex con int wis cha)
  slotsMax: number[]    // 9 entries
  pact: { count: number; level: number }  // {0,0} unless a pact caster
  resourcesMax: Record<string, number>    // every class resource id; 0 = not gained at this level
  hitDiceMax: number    // = level (die size = class hitDie)
}
type FeedEntry = { at: number; text: L }   // at = unix ms
```
- `roster` is in every `RoomState` (nothing secret): DM/co-DM get every set (Ruleset editor), players and unseated
  connections only entries whose `set` is in `activeSets`. Default = the JSON blocks in `docs/design/roster.md` (core)
  then `docs/design/roster-legend.md` (legend) (Go: `internal/rules/roster_default.go` + `roster_legend.go`, kept in
  sync by a test). A roster edit of an inactive set sends players nothing.
- `roster.update` limits: heroName th+en trimmed 1–40, blurb th+en trimmed 1–300 (newlines ok), role in the
  enum, abilities 1–30 and ac 1–30 (as `character.setAbilities`), skills = known skill ids without duplicates (any
  skill, not only the class's choices), expertise ⊆ skills without duplicates; items 0–10 with name th+en 1–40, qty
  1–99; `art` "" or an existing `/uploads/<32 hex>.<png|jpg|gif|webp>`. `set` and `classId` are always the built-in ones.
  Text with control characters (other than newline in blurb/desc) is rejected.
- Persisted to `<-data>/roster.json` as `{"edits": RosterEntry[]}` (edited classes only; unedited classes follow the
  built-in defaults; one file for every set, keyed by entry id), atomic write like notes. Missing file = defaults;
  invalid JSON or a pre-rules-v2 file → `roster.json.bad-<random>` + defaults (logged as a warning); an entry that
  fails validation on load is dropped (that class uses its default) and logged.
- `feed` keeps the newest 50 entries. Lines are written by the server, never by clients, and hold seat and hero names only
  (no clientIds). Events that add a line: `character.pick` ("Pat picked Aldric Vane (Fighter)" — class name),
  `character.adjust` ("DM: Aldric Vane HP −5 (7/12)" — the change after clamping; no change = no line),
  `character.level`, `party.rest` (lines in "Rules v2"), `item.give`, `item.remove`
  ("DM: removed Field Bandages ×2 from Aldric Vane" — the qty actually removed). setAbilities, setSkills, slot, resource,
  hit-dice, custom-ability and roster events add a DM log line only.
- Errors: `only a seated player can pick a character`, `you already have a character`, `no such class`,
  `that set is not in play`, `pick at least one set`, `no such set`, `send classId or set, not both`,
  `that class is already taken`, `you have no character`, `no such character`, `no such item`,
  `hp delta is required`, `delta must be -999 to 999`, `item name must be 1-40 characters`, `item qty must be 1-99`,
  `an item stack holds at most 999`, `at most 50 different items`, `qty must be at least 1`,
  `heroName must be 1-40 characters (th and en)`, `blurb must be 1-300 characters (th and en)`,
  `role must be one of tank, damage, healer, support, control, hybrid`, `art must be empty or an uploaded image`,
  `abilities must be 1-30`, `ac must be 1-30`, `unknown skill id`, `duplicate skill`,
  `expertise must be a proficient skill`, `duplicate expertise`,
  `a class starts with at most 10 items`, `item name must be 1-40 characters (th and en)`, `could not save roster`,
  the "Rules v2" errors,
  plus the shared `only the DM can do that`, `not allowed while the room is <status>`, `malformed payload`.

## Rules v2 (v0.18.0, `docs/features/rules-v2.md`)
- Class data (12 classes, SRD 5.2.1 based) is static and **not** in `RoomState`: `GET /api/classes` (no auth, like
  the page). `Cache-Control: no-cache` + `ETag` (revalidate → `304`); about 140 KB uncompressed. Body:
  ```ts
  type Class = {
    id: string; name: L; hitDie: 6 | 8 | 10 | 12; primary: string[]; saves: string[]   // ability ids
    armor: L; weapons: L; tools: L                                                    // tools {"",""} = none
    skillChoices: { count: number; from: string[] }
    caster: 'none' | 'full' | 'half' | 'pact'; castAbility: string                    // "" when none
    tables: { id: string; name: L; values: string[] }[]                               // 20 values, "—" = none yet
    resources: { id: string; name: L; max: number[] | null; maxFormula: string;       // "abilityMod:<ab>" | "level*<n>" | ""
                 refill: 'short' | 'long'; shortRegain: number; shortFromLevel: number; fromLevel: number }[]
    features: { level: number; name: L; desc: L }[]                                   // sorted by level
  }
  { "classes": Class[],                      // 12
    "skills": [{ "id": "athletics", "ability": "str" }, ...],   // 18, sheet order
    "proficiency": number[20],               // index = level - 1
    "fullSlots": number[20][9], "halfSlots": number[20][9],
    "pactSlots": [{ "count": 1, "level": 1 }, ...],             // 20
    "attribution": "This work includes material from the System Reference Document 5.2.1 (\"SRD 5.2.1\") by ..." }
  ```
  `attribution` is the exact CC BY 4.0 credit; show it on the pick screen (next to the art credit).
- Client-side display only (simple, from the class data + the Character): ability mod `floor((score−10)/2)`, save =
  mod + `pb` if in `character.saves` (not `class.saves`: it misses the extra saves), skill = mod + `pb` if in `skills`, + 2×`pb` if in `expertise`, initiative = DEX mod,
  passive Perception = 10 + Perception, features with `level ≤ character.level`. Maxima and hpMax come from the server.
- `character.pick`: `level 1`, `hpMax = hitDie + CON mod` (≥ 1), `hp = hpMax`, abilities/ac/skills/expertise/items
  copied from the hero, every used counter 0, `abilitiesCustom []`.
- hpMax at level n = `max(1, hitDie + CON mod) + (n−1) × max(1, hitDie/2 + 1 + CON mod)`. `character.level` and
  `character.setAbilities` recompute it and keep the damage taken: `hp += newHpMax − oldHpMax`, clamped `0..hpMax`
  (crossing 0 adds the knocked-out / back-up line). Every used counter clamps down to its new max.
- `party.rest` short: a resource with `refill: "short"`, or with `shortFromLevel > 0` and `level ≥ shortFromLevel`,
  resets to 0 used; any other resource gives back `shortRegain` uses (never below 0); `pactUsed = 0`; hp, spell slots
  and hit dice unchanged (players spend hit dice with `hitDice.set` and the DM adds the healed HP with
  `character.adjust` — dice stay physical). Long: `hp = hpMax`, all slots, pact, resources and hit dice 0 used.
- Feed (+ DM log): `"Ayla reaches level 4"` / `"Ayla ขึ้นเลเวล 4"`, `"Ayla drops to level 3"` / `"Ayla ลดเป็นเลเวล 3"`,
  `"DM: short rest — short-rest uses and pact slots back"` / `"DM: พักสั้น — คืนความสามารถแบบพักสั้นและ pact slot"`,
  `"DM: long rest — HP, spell slots, uses and hit dice refilled"` / `"DM: พักยาว — เลือด ช่องคาถา ความสามารถ และ hit dice เต็ม"`.
- DM log only: `"set Ayla: STR 17, DEX 10, CON 8, INT 8, WIS 13, CHA 12, AC 20 (HP 7/9)"`,
  `"set skills of Ayla: athletics, stealth; expertise: stealth; extra saves: wis, cha"` (empty list = `—`; th
  `"ตั้ง skill ของ Ayla: …; expertise: …; save เพิ่ม: …"`),
  `"Ayla: level 1 spell slots used 2/4"`, `"Ayla: pact slots used 1/2"`, `"Ayla: Second Wind used 1/2"` (resource
  name), `"Ayla: hit dice used 1/4"`, `"added spell Fireball to Ayla"` / `"added feature Tough to Ayla"`,
  `"edited Fireball of Ayla"`, `"removed Fireball from Ayla"` (th lines alongside).
- Errors: `delta must be 1 or -1`, `level must be 1-20`, `abilities must be 1-30`, `ac must be 1-30`,
  `you can only change your own character`, `level must be 1-9 or "pact"`, `no such spell slot`, `no such resource`,
  `used must be 0-<max>`, `kind must be short or long`, `you have no character`, `no such ability`,
  `ability name must be 1-40 characters`, `ability desc must be at most 300 characters`, `kind must be spell or feature`,
  `spell level must be 0-9`, `at most 40 custom abilities`, `unknown skill id`, `duplicate skill`,
  `expertise must be a proficient skill`, `duplicate expertise`, `unknown ability id`, plus the shared ones.
