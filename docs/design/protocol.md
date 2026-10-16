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
- Reconnect: same `clientId` → same seat (no room number in v0).

## RoomState
```jsonc
{
  "status": "menu" | "lobby" | "creating" | "playing",
  "campaigns": [{ "id": "starter-set", "name": "Starter Set: Dragons of Stormwreck Isle", "status": "available" },
                { "id": "salem", "name": "Salem", "status": "coming_soon" }],
  "campaignId": "",                       // "" in menu
  "seats": [{ "id": "s1", "kind": "dm" | "codm" | "player", "name": "", "occupied": false, "ready": false, "characterId": "" }],
  "roster": [ /* RosterEntry x 10, see "Character roster" below; same for every connection */ ],
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
  "feed": [{ "at": 1760000000000, "text": { "th": "...", "en": "..." } }]   // newest last, at most 50
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
| `game.start` | `{}` | DM | `lobby` | ≥1 occupied player seat, all occupied player seats ready; empty player seats are removed; → `creating` |
| `character.pick` | `{classId}` | player without character | `creating` | class in the roster and not taken by another character; creates the Character from the current roster entry (`hp = hpMax`, `mp = mpMax`, innate skills, items copied); when every occupied player seat has a character → `playing` |
| `skill.addCustom` | `{name, desc, cost}` | owner (player with a character) | `creating`/`playing` | name trimmed 1–40, desc trimmed ≤ 300 (newlines ok), cost int 0–10; at most 3 custom skills; `th = en =` the typed text |
| `skill.editCustom` | `{characterId, skillId, name, desc, cost}` | owner or DM | `creating`/`playing` | same limits; `kind: "custom"` only |
| `skill.removeCustom` | `{characterId, skillId}` | owner or DM | `creating`/`playing` | `kind: "custom"` only |
| `skill.use` | `{skillId}` | owner | `playing` | skill of the sender's own character (innate or custom); `mp ≥ cost` → `mp −= cost` |
| `character.adjust` | `{characterId, hp?, mp?}` | DM | `creating`/`playing` | deltas, int −999..999, at least one; result clamps to `0..max` |
| `character.setStats` | `{characterId, hpMax, mpMax, str, end, agi}` | DM | `creating`/`playing` | hpMax 1–999, mpMax 0–99, str/end/agi 0–20; `hp`/`mp` clamp down to the new max |
| `party.rest` | `{}` | DM | `playing` | every character `mp = mpMax` (hp unchanged) |
| `item.give` | `{characterId, name, qty}` | DM | `creating`/`playing` | name trimmed 1–40, qty 1–99; a stack whose `name.th` or `name.en` equals name gets `+qty` (stack ≤ 999); else a new stack `{name: {th: name, en: name}, qty}` (≤ 50 stacks) |
| `item.remove` | `{characterId, index, qty}` | DM | `creating`/`playing` | index into `items`; qty ≥ 1; qty ≥ stack qty removes the stack |
| `roster.update` | `{classId, entry: RosterEntry}` | DM | any | classId is the key (`entry.id` is ignored); limits below; affects later picks only; persisted |
| `roster.reset` | `{classId?}` | DM | any | `""`/missing = every class back to its default; persisted |
| `scene.add` | `{name, imageUrl, grid}` | DM | `creating`/`playing` | imageUrl must start with `/uploads/`; grid: cellSize 10–400, cols/rows 1–200 |
| `scene.update` | `{sceneId, name, grid}` | DM | `creating`/`playing` | same grid limits (cellSize is an int: send `round(imgW/cols)`; offsetX/offsetY 0). Grid shrink: that scene's tokens clamp to the nearest edge cell `[0,cols-1]×[0,rows-1]`; fog rects clip to the grid, rects fully outside are dropped |
| `scene.select` | `{sceneId}` | DM | `playing` | sets `currentSceneId`; adds a token for every character without one in that scene, at row 0, col 0..n (name = `heroName.en`, imageUrl = character `art`) |
| `npc.place` | `{npcId, sceneId}` | DM | `playing` | copies the pool NPC into a new token on the first free cell (row-major) of the scene; `hp = hpMax`, `hidden`, `imageUrl` from the pool; name `<name>`, then `<name> #2`, `#3`… (first name not used by a token in that scene); a visible NPC ignores hidden tokens for both the cell and the name (may stack on one; a name may repeat after a reveal); full scene → error. Replaces `token.addNpc` (removed v0.11.0) |
| `token.move` | `{tokenId, col, row}` | DM any token; player only own character's token | `playing` | 0 ≤ col < cols, 0 ≤ row < rows of the token's scene. A player gets `you can only move your own token` for a missing token too (no hidden-token oracle) |
| `token.setHidden` | `{tokenId, hidden}` | DM | `playing` | NPC tokens placed from the pool only |
| `token.adjustHp` | `{tokenId, delta}` | DM | `playing` | NPC tokens placed from the pool only; `-999 ≤ delta ≤ 999` (int); `hp` clamps to `0..hpMax` |
| `token.remove` | `{tokenId}` | DM | `playing` | |
| `fog.add` | `{sceneId, col, row, w, h}` | DM | `playing` | w,h ≥ 1, rect inside grid |
| `fog.remove` | `{sceneId, fogId}` | DM | `playing` | |
| `note.add` | `{title}` | DM | any | appends `{id, title, body: "", updatedAt}`; title trimmed, 1–80 chars; max 200 notes |
| `note.update` | `{id, title, body}` | DM | any | title 1–80 (trimmed), body ≤ 20000 chars (not trimmed); sets `updatedAt`; last write wins |
| `note.remove` | `{id}` | DM | any | |
| `note.move` | `{id, index}` | DM | any | `index` = position in the resulting list, `0 ≤ index < len(notes)` |
| `npc.add` | `{name, imageUrl, hpMax, hidden}` | DM | any | name trimmed 1–40 chars; `hpMax` 1–999; `imageUrl` "" or an existing `/uploads/<generated name>`; max 100 NPCs |
| `npc.update` | `{id, name, imageUrl, hpMax, hidden}` | DM | any | same checks; tokens already placed keep their copy |
| `npc.remove` | `{id}` | DM | any | tokens already placed stay |

Unknown `type`, malformed payload, or failed check → `error` to sender. Removed in v0.12.0 (now unknown types):
`ruleset.update`, `ruleset.loadPreset`, `character.submit`, `character.grant`.

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
- Errors: `only the DM can do that`, `npc name must be 1-40 characters`, `npc hpMax must be 1-999`, `at most 100 npcs`,
  `no such npc`, `imageUrl must be empty or an uploaded image`, `could not save npcs`, `scene has no free cell`,
  `only NPC tokens from the pool`, `delta must be -999 to 999`, `no such token`, `malformed payload`.
- Persisted to `<-data>/npcs.json` exactly like notes (atomic write, dir `0700`, file `0600`, corrupt file → `npcs.json.bad-<random>`).
- Logs name the event and sender only (`event npc.place by=<name>`), never NPC or token names.
- Token ids share the room counter (`t<seq>`): accepted: id gaps may hint that hidden tokens exist; names/hp/position never leak.

## Character roster (v0.12.0, `docs/features/character-roster.md`)
```ts
type L = { th: string; en: string }
type RosterSkill = { id: string; name: L; desc: L; cost: number }
type RosterItem  = { name: L; qty: number }
type RosterEntry = {
  id: string            // classId: knight|wizard|rogue|ranger|cleric|paladin|barbarian|bard|druid|necromancer
  heroName: L; className: L; blurb: L
  role: string          // tank|damage|healer|support|control|hybrid
  art: string           // "" = class art by id (web/src/assets/classes/<id>.svg) | "/uploads/<generated name>"
  stats: { hpMax: number; str: number; end: number; agi: number; mpMax: number }
  skills: RosterSkill[] // 1-5
  items: RosterItem[]   // 0-10
}
type CharSkill = { id: string; kind: 'innate' | 'custom'; name: L; desc: L; cost: number }
type Character = {
  id: string; seatId: string; classId: string
  name: L               // heroName at pick time
  art: string           // entry art at pick time
  hp: number; hpMax: number; mp: number; mpMax: number; str: number; end: number; agi: number
  skills: CharSkill[]   // innate (roster order) then custom (≤ 3, ids "cs_<n>")
  items: { name: L; qty: number }[]
}
type FeedEntry = { at: number; text: L }   // at = unix ms
```
- `roster` is in every `RoomState` (DM and players alike; nothing secret). Default = the JSON block in
  `docs/design/roster.md` (Go: `internal/rules/roster_default.go`, kept in sync by a test).
- `roster.update` limits: heroName/className th+en trimmed 1–40, blurb th+en trimmed 1–300 (newlines ok), role in the
  enum, stats as `character.setStats`, skills 1–5 with unique id `[a-z0-9-]{1,40}`, name th+en 1–40, desc th+en 1–300,
  cost 0–10; items 0–10 with name th+en 1–40, qty 1–99; `art` "" or an existing `/uploads/<32 hex>.<png|jpg|gif|webp>`.
  Text with control characters (other than newline in blurb/desc) is rejected.
- Persisted to `<-data>/roster.json` as `{"edits": RosterEntry[]}` (edited classes only; unedited classes follow the
  built-in defaults), atomic write like notes. Missing file = defaults; invalid JSON → `roster.json.bad-<random>` +
  defaults; an entry that fails validation on load is dropped (that class uses its default) and logged.
- `feed` keeps the newest 50 entries. Lines are written by the server, never by clients, and hold seat and hero names only
  (no clientIds). Events that add a line: `character.pick` ("Pat picked Aldric Vane (Knight)"), `skill.use`
  ("Aldric Vane used Last Stand (−4 MP)"), `character.adjust` ("DM: Aldric Vane HP −5 (31/36), MP +2 (4/4)" — the
  change after clamping), `character.setStats`, `party.rest` ("DM: rest — mana refilled"), `item.give`, `item.remove`
  ("DM: removed Field Bandage Roll ×2 from Aldric Vane" — the qty actually removed). Custom-skill and roster events add none.
- Errors: `only a seated player can pick a character`, `you already have a character`, `no such class`,
  `that class is already taken`, `you have no character`, `no such character`, `no such skill`, `no such item`,
  `only custom skills can be changed`, `you can only change your own character's skills`, `at most 3 custom skills`,
  `skill name must be 1-40 characters`, `skill desc must be at most 300 characters`, `skill cost must be 0-10`,
  `not enough mana`, `hp or mp delta is required`, `delta must be -999 to 999`, `hpMax must be 1-999`, `mpMax must be 0-99`,
  `str, end and agi must be 0-20`, `item name must be 1-40 characters`, `item qty must be 1-99`,
  `an item stack holds at most 999`, `at most 50 different items`, `qty must be at least 1`,
  `heroName must be 1-40 characters (th and en)`, `className must be 1-40 characters (th and en)`,
  `blurb must be 1-300 characters (th and en)`, `role must be one of tank, damage, healer, support, control, hybrid`,
  `art must be empty or an uploaded image`, `a class needs 1-5 skills`, `skill id must be 1-40 characters a-z, 0-9 or -`,
  `duplicate skill id`, `skill name must be 1-40 characters (th and en)`, `skill desc must be 1-300 characters (th and en)`,
  `a class starts with at most 10 items`, `item name must be 1-40 characters (th and en)`, `could not save roster`,
  plus the shared `only the DM can do that`, `not allowed while the room is <status>`, `malformed payload`.
