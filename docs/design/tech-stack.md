# Tech stack v2: backend Go + frontend React/TS — **ยืนยันแล้ว 2026-10-16** (ยกเว้นที่ฝาก server: Q12)

วันที่ 2026-10-16 · v1 (Node/TypeScript ทั้งหมด) ถูกแทนด้วยฉบับนี้ตามที่ผู้ใช้อยากใช้ Go ([conv](../conversations/2026-10-16-tech-stack.md))
อิงจาก [`web-flow.md`](web-flow.md), [`../research/vtt-and-player-needs.md`](../research/vtt-and-player-needs.md)

## สรุป: server = Go, หน้าเว็บ = React + TypeScript, ได้ binary ไฟล์เดียว
หน้าเว็บยังต้องเป็น JavaScript/TypeScript เพราะ browser รันแค่ JS (Go → WASM ทำได้ แต่ไฟล์ใหญ่ ใช้ library canvas/React ไม่ได้ → ไม่คุ้ม)

| ชั้น | ใช้ | เพราะอะไร |
|---|---|---|
| **server** | **Go 1.26** (มีในเครื่องแล้ว) + `net/http` (standard library) | routing ของ `net/http` รองรับ method + path param แล้ว ไม่ต้องใช้ framework |
| **realtime** | **WebSocket** ด้วย `github.com/coder/websocket` | Socket.IO ไม่มีตัวทางการใน Go → ใช้ WebSocket ตรง ๆ ส่ง JSON; ไม่ใช้ WebRTC (บทเรียน Owlbear) |
| state ห้อง | struct ใน memory + `sync.Mutex` / goroutine เดียวต่อห้อง | v0 ไม่เซฟข้าม session; อัปเกรดทีหลัง = เขียน JSON ลงไฟล์ |
| ใครเป็นเจ้าของ state | **server** ตัดสินทุกอย่าง | server ตรวจสิทธิ์ → broadcast; ไม่ต้องใช้ CRDT |
| ลูกเต๋า | **ไม่มีในแอป** — ผู้เล่นทอยลูกเต๋าจริง (ตัดออก 2026-10-16, [conv](../conversations/2026-10-16-remove-dice.md)) | |
| ตรวจกฎ | Go (`internal/rules`): ข้อมูล 12 คลาส (`classes.go` + `classes_data.go` gen จาก `docs/design/classes/`), roster 2 เซ็ต (`roster.go` + `roster_default.go`/`roster_legend.go` gen จาก `heroes-core.md`/`heroes-legend.md`), สูตร rules v2 (ability mod, PB, HP, slot, counter), ขอบเขต ability/AC/skill, สิทธิ์ DM/ผู้เล่น | ตรวจที่ server เสมอ หน้าเว็บตรวจซ้ำเพื่อ UX เท่านั้น · ~~แต้ม stat / ruleset white label / preset 5E~~ แทนที่ด้วย roster 2026-10-16 ([conv](../conversations/2026-10-16-character-roster.md)) |
| เก็บรูป | โฟลเดอร์ `uploads/` + `http.MaxBytesReader` + ตรวจ `http.DetectContentType` | ไม่ต้องมี cloud storage |
| **หน้าเว็บ** | **React + Vite + TypeScript** | มาตรฐาน, library เยอะ |
| แมพ / grid / token / fog | **react-konva** (MIT) | canvas 2D ลาก token ง่าย |
| state หน้าเว็บ | **Zustand** | เบา |
| animation | **`motion`** (motion.dev, `motion/react`) | carousel การ์ดเลือกตัวละคร; animate เฉพาะ transform/opacity, รองรับ `prefers-reduced-motion` |
| QR | **`qrcode`** | QR เข้านั่งในหน้า lobby ([spec](../features/seat-qr.md)) |
| type ระหว่าง Go ↔ TS | เขียน type ข้อความ (message) ไว้ที่เดียวใน Go แล้ว gen เป็น TS ด้วย **tygo** — หรือเขียนมือถ้ายังน้อย | กันสองฝั่งไม่ตรงกัน |
| **deploy** | `go build` → **binary ไฟล์เดียว** ที่ฝังหน้าเว็บไว้ด้วย `embed` | ก๊อปไฟล์เดียวไปรันที่ไหนก็ได้ |
| ทดสอบ | `go test` (กฎ, สิทธิ์) + **Playwright** (2 แท็บ DM + ผู้เล่น) | |

## โครงโฟลเดอร์
```
go.mod
cmd/server/main.go     เปิด HTTP + WebSocket, เสิร์ฟหน้าเว็บที่ embed ไว้
internal/room/         state ห้อง, เก้าอี้, ตัวละคร, ฉาก, token, fog + จัดการ event
internal/rules/        12 คลาส (classes*.go) + roster 2 เซ็ต (roster*.go) + สูตร rules v2 + ตรวจขอบเขต
internal/server/       HTTP + WS + upload + GET /api/classes (ข้อมูลคลาส + attribution SRD)
web/                   React + Vite (package.json มี version SemVer)
web/dist/              ผล build → ถูก embed เข้า binary (gitignored)
uploads/               รูปที่อัป (gitignored)
```
Version: SemVer อันเดียวของทั้งโปรเจกต์ — เก็บใน `web/package.json` + ฝังเข้า binary ด้วย `-ldflags "-X main.version=..."` (ยืนยันกับผู้ใช้ตอน scaffold)

## Event หลัก (WebSocket JSON `{type, payload}`; server ตรวจสิทธิ์แล้ว broadcast)
| type | ใครส่งได้ |
|---|---|
| `campaign.select` | DM |
| `seat.sit`, `seat.ready` | ทุกคน (เก้าอี้ตัวเอง) |
| `seat.add`, `seat.remove`, `seat.openCoDm`, `game.start` | DM |
| `character.pick {heroId}` (ฮีโร่ของเซ็ตที่เปิด, ห้ามซ้ำ) | ผู้เล่น (ตัวเอง) |
| `slot.set`, `resource.set`, `hitDice.set` (ค่า used แบบ absolute) | เจ้าของ / DM |
| `ability.add`, `ability.update` (คาถา/ความสามารถของฉัน ≤40) | ผู้เล่น (ตัวละครตัวเอง) |
| `ability.remove` | เจ้าของ / DM |
| `character.adjust` (HP ±), `character.level` (±1), `character.setAbilities`, `character.setSkills`, `party.rest {kind}` | DM |
| `item.give`, `item.remove` | DM |
| `roster.update`, `roster.reset` (เก็บ `data/roster.json`) | DM |
| `note.*`, `npc.*`, `npc.place`, `token.setHidden`, `token.adjustHp` | DM |
| `scene.add`, `scene.select`, `fog.set` | DM |
| `token.move` | ผู้เล่น (token ตัวเอง) / DM (ทุกตัว) |

~~`ruleset.update`, `ruleset.loadPreset`, `character.submit`, `character.grant`~~ ลบใน v0.12.0 · ~~`skill.use`, `skill.addCustom/editCustom/removeCustom`, `character.setStats`~~ ลบใน v0.18.0 ([rules v2](../features/rules-v2.md)) · รายการเต็ม + เงื่อนไข: [`protocol.md`](protocol.md)

server ส่งกลับ: `state.snapshot` (ตอนเข้า/เชื่อมใหม่) + `state.patch` / event เดิมหลังตรวจผ่าน + `error`

## ฝาก server ที่ไหน
ต้องเปิดค้าง + มี disk (WebSocket + รูป)
| ทางเลือก | ข้อดี | ข้อเสีย |
|---|---|---|
| **ก. คอมตัวเอง + Cloudflare Tunnel** (แนะนำครั้งแรก) | ฟรี, ไม่เปิด port, ได้ URL https; Go binary รันได้เลย | คอมต้องเปิดตอนเล่น; ต้องติดตั้ง `cloudflared` (ยังไม่มีในเครื่อง) |
| ข. VPS เล็ก / Fly.io | เปิดตลอด; Go binary กิน RAM น้อย | เสียเงินรายเดือนเล็กน้อย (ตรวจราคาก่อน) |
| ✗ Render/Vercel free | — | หลับ / ไม่รองรับ WebSocket ค้าง / disk หาย |

## Go vs Node (ข้อเสนอ v1) — เทียบสั้น ๆ
| | Go | Node |
|---|---|---|
| ผู้ใช้ชอบ | ✅ | |
| deploy | binary ไฟล์เดียว | ต้องมี Node + node_modules |
| type ร่วมกับหน้าเว็บ | ต้อง gen (tygo) | ใช้ร่วมได้ตรง ๆ |
| realtime | WebSocket ตรง (เขียน reconnect เองฝั่งหน้าเว็บ ~20 บรรทัด) | Socket.IO มี reconnect ในตัว |
| ภาษาที่ต้องใช้ | 2 (Go + TS) | 1 |

สรุป: Go ได้เปรียบเรื่อง deploy + ผู้ใช้ถนัด; ต้นทุนเพิ่มคือ gen type + reconnect — เล็ก

## ที่ไม่เลือก
Next.js (ไม่ต้องการ SSR) · WebRTC (โดนบล็อก, ต้องมี TURN) · Yjs/CRDT (เกินจำเป็น) · Supabase/Firebase (ยังไม่ต้องเก็บถาวร) · tldraw (ไม่ใช่ open source) · Go→WASM หน้าเว็บ (ไฟล์ใหญ่, ไม่มี library canvas/UI) · framework Go (Gin/Echo/Fiber — `net/http` พอ)

## ความปลอดภัย (ขั้นต่ำ)
- server ตรวจสิทธิ์ทุก event (ผู้เล่นขยับได้แค่ token ตัวเอง, คำสั่ง DM ต้องมาจากเก้าอี้ DM)
- อัปรูป: จำกัดขนาด (เช่น 10MB), รับเฉพาะ image (`DetectContentType`), ตั้งชื่อไฟล์ใหม่เอง (กัน path traversal)
- WebSocket: ตรวจ `Origin` + จำกัดขนาดข้อความ
- ห้องไม่มีรหัส (ยอมรับแล้ว) → ใครได้ URL เข้าได้
