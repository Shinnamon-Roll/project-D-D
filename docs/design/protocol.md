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
    "skills": [{ "id": "athletics", "name": "Athletics", "statKey": "str" }],   // statKey may be ""
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
  "tokens": [{ "id": "t1", "sceneId": "sc1", "characterId": "c1", "name": "Mira", "imageUrl": "/uploads/x.png", "col": 0, "row": 0 }],
                                          // characterId "" = NPC token
  "dice": [{ "id": "d1", "by": "Mira", "notation": "1d20+3", "rolls": [14], "total": 17, "at": "2026-10-16T12:00:00Z" }]
                                          // newest last, max 200
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
| `scene.update` | `{sceneId, name, grid}` | DM | `creating`/`playing` | same grid limits |
| `scene.select` | `{sceneId}` | DM | `playing` | sets `currentSceneId`; adds a token for every character without one in that scene, at row 0, col 0..n |
| `token.addNpc` | `{sceneId, name, imageUrl}` | DM | `playing` | placed at col 0, row 0; imageUrl "" or `/uploads/...` |
| `token.move` | `{tokenId, col, row}` | DM any token; player only own character's token | `playing` | 0 ≤ col < cols, 0 ≤ row < rows of the token's scene |
| `token.remove` | `{tokenId}` | DM | `playing` | |
| `fog.add` | `{sceneId, col, row, w, h}` | DM | `playing` | w,h ≥ 1, rect inside grid |
| `fog.remove` | `{sceneId, fogId}` | DM | `playing` | |
| `dice.roll` | `{notation}` | seated | not `menu` | grammar below; result appended to `dice` |

Unknown `type`, malformed payload, or failed check → `error` to sender.

## Dice notation
Terms joined by `+`/`-`: `NdM` (1 ≤ N ≤ 100, 2 ≤ M ≤ 1000; `dM` = `1dM`), optional `khK`/`klK` (keep highest/lowest K, 1 ≤ K ≤ N),
or an integer constant. Examples: `1d20+3`, `2d20kh1+5`, `4d6kl3`, `d100`, `3d8-1`. Whitespace ignored, case-insensitive.
`rolls` = every die rolled (in order), `total` = kept dice + constants. `crypto/rand`.

## Upload
`POST /api/upload`, header `X-Client-ID: <clientId>` of a **seated** client, multipart field `file`.
- ≤ 10 MB (`413`), content sniffed with `http.DetectContentType` ∈ png/jpeg/gif/webp (`415`), not seated (`401`).
- Saved as `uploads/<random hex>.<ext>`; response `200 {"url": "/uploads/<name>"}`; errors `{"error": "..."}`.
- `GET /uploads/<name>` serves files (no directory listing).

## Server bind
Default `-addr 127.0.0.1:8080`.
