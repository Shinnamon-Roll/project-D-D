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
- Reconnect: same `clientId` → same seat (no room number in v0).

## RoomState
```jsonc
{
  "status": "menu" | "lobby" | "creating" | "playing",
  "campaigns": [{ "id": "starter-set", "name": "Starter Set: Dragons of Stormwreck Isle", "status": "available" },
                { "id": "salem", "name": "Salem", "status": "coming_soon" }],
  "campaignId": "",                       // "" in menu
  "seats": [{ "id": "s1", "kind": "dm" | "codm" | "player", "name": "", "occupied": false, "ready": false, "characterId": "" }],
  "ruleset": {
    "name": "D&D 5e (SRD 5.2.1)",
    "attribution": "This work includes material from the System Reference Document 5.2.1 ...",  // "" for custom
    "stats":  [{ "key": "str", "name": "STR", "min": 8, "max": 15, "base": 8 }],
    "statPoints": 27,                     // points a player adds on top of base at creation
    "skills": [{ "id": "athletics", "name": "Athletics", "statKey": "str", "description": "" }],   // statKey may be ""; description optional (preset: "")
    "startingSkillPicks": 2,
    "items":  [{ "id": "rope", "name": "Rope" }],
    "startingItems": [{ "itemId": "rope", "qty": 1 }]
  },
  "characters": [{ "id": "c1", "seatId": "s2", "name": "Mira", "portraitUrl": "/uploads/x.png",
                   "stats": { "str": 10 }, "skillIds": ["athletics"], "inventory": [{ "itemId": "rope", "qty": 1 }] }],
  "scenes": [{ "id": "sc1", "name": "Beach", "imageUrl": "/uploads/y.png",
               "grid": { "cellSize": 50, "offsetX": 0, "offsetY": 0, "cols": 20, "rows": 15 },
               "fog": [{ "id": "f1", "col": 0, "row": 0, "w": 3, "h": 2 }] }],
  "currentSceneId": "",
  "tokens": [{ "id": "t1", "sceneId": "sc1", "characterId": "c1", "name": "Mira", "imageUrl": "/uploads/x.png", "col": 0, "row": 0 }]
                                          // characterId "" = NPC token
}
```

## Events (client → server)
Roles: **DM** = seat kind `dm` or `codm` (co-DM has every DM right except removing the `dm` seat). **Player** = seat kind `player`.
| type | payload | who | when | effect / validation |
|---|---|---|---|---|
| `hello` | `{clientId}` | anyone | any | bind connection; must be first |
| `campaign.select` | `{campaignId, name}` | anyone not seated | `menu` | campaign must be `available`; sender sits in a new `dm` seat; ruleset = campaign preset; status → `lobby`; 1 empty `player` seat created |
| `seat.sit` | `{seatId, name}` | not seated | `lobby` | seat empty, kind `player` or `codm`; name 1–32 chars |
| `seat.leave` | `{}` | seated, not `dm` | `lobby` | frees seat, clears ready |
| `seat.ready` | `{ready}` | player | `lobby` | |
| `seat.add` | `{}` | DM | `lobby` | max 6 `player` seats |
| `seat.remove` | `{seatId}` | DM | `lobby` | not the `dm` seat; kicks occupant |
| `seat.openCoDm` | `{}` | DM | `lobby` | max 1 `codm` seat |
| `game.start` | `{}` | DM | `lobby` | ≥1 occupied player seat, all occupied player seats ready; empty player seats are removed; → `creating` |
| `ruleset.update` | `{ruleset}` | DM | not `menu` | unique keys/ids; `min ≤ base ≤ max`; counts ≥ 0; startingItems reference items |
| `ruleset.loadPreset` | `{presetId: "dnd5e-srd"}` | DM | not `menu` | replace ruleset |
| `character.submit` | `{name, portraitUrl, stats, skillIds}` | player without character | `creating` | every stat key present, `min ≤ v ≤ max`, `Σ(v − base) == statPoints`; exactly `startingSkillPicks` distinct known skills; inventory = startingItems; when every player has a character → `playing` |
| `character.grant` | `{characterId, setStats?, addSkillIds?, removeSkillIds?, addItems?, removeItems?}` | DM | `creating`/`playing` | DM bypasses min/max; skills/items must exist in ruleset; qty > 0; item qty never below 0 (removed at 0) |
| `scene.add` | `{name, imageUrl, grid}` | DM | `creating`/`playing` | imageUrl must start with `/uploads/`; grid: cellSize 10–400, cols/rows 1–200 |
| `scene.update` | `{sceneId, name, grid}` | DM | `creating`/`playing` | same grid limits (cellSize is an int: send `round(imgW/cols)`; offsetX/offsetY 0). Grid shrink: that scene's tokens clamp to the nearest edge cell `[0,cols-1]×[0,rows-1]`; fog rects clip to the grid, rects fully outside are dropped |
| `scene.select` | `{sceneId}` | DM | `playing` | sets `currentSceneId`; adds a token for every character without one in that scene, at row 0, col 0..n |
| `token.addNpc` | `{sceneId, name, imageUrl}` | DM | `playing` | placed at col 0, row 0; imageUrl "" or `/uploads/...` |
| `token.move` | `{tokenId, col, row}` | DM any token; player only own character's token | `playing` | 0 ≤ col < cols, 0 ≤ row < rows of the token's scene |
| `token.remove` | `{tokenId}` | DM | `playing` | |
| `fog.add` | `{sceneId, col, row, w, h}` | DM | `playing` | w,h ≥ 1, rect inside grid |
| `fog.remove` | `{sceneId, fogId}` | DM | `playing` | |
| `note.add` | `{title}` | DM | any | appends `{id, title, body: "", updatedAt}`; title trimmed, 1–80 chars; max 200 notes |
| `note.update` | `{id, title, body}` | DM | any | title 1–80 (trimmed), body ≤ 20000 chars (not trimmed); sets `updatedAt`; last write wins |
| `note.remove` | `{id}` | DM | any | |
| `note.move` | `{id, index}` | DM | any | `index` = position in the resulting list, `0 ≤ index < len(notes)` |

Unknown `type`, malformed payload, or failed check → `error` to sender.

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
- Ruleset: names/keys/ids 1–64 chars, attribution ≤1000, skill `statKey` "" or an existing stat, `startingSkillPicks` ≤ number of skills, `statPoints` ≤ Σ(max − base), startingItems qty ≥ 1 and unique. Skill `description` optional: trimmed, ≤500 chars (runes), newlines allowed, other control characters rejected; preset skills send `""` (frontend shows SRD text when the skill id matches the preset and description is empty).
- Names (DM, seat, character, scene, NPC): 1–32 chars.
- `hello.clientId`: 8–64 chars `[A-Za-z0-9_-]`; any event before `hello`, or a second `hello` → error.
- Image URLs: exactly `/uploads/<name>` (no extra `/`, `..`, `?`, `#`; name ≤64).
- `scene.update` clamps tokens and clips/drops fog in a shrunken grid; `scene.select` wraps auto-placed tokens to the next row.
- `character.grant`: `setStats` keys and `removeSkillIds` must exist.
- A co-DM may remove their own seat; an empty co-DM seat survives `game.start`.
- Preset HP and AC: min = max = base = 10 (DM sets them via grant).
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
