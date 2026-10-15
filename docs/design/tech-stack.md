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
| ทอยเต๋า | **เขียน parser เองใน Go** (`NdM`, `+/-K`, `kh/kl` สำหรับ advantage) + `crypto/rand` | ~100 บรรทัด + test; ไม่ต้องพึ่ง library; ทอยบน server = โกงไม่ได้ |
| ตรวจกฎ | Go (`internal/rules`): แต้ม stat, จำนวนสกิล, สิทธิ์ DM/ผู้เล่น | ตรวจที่ server เสมอ หน้าเว็บตรวจซ้ำเพื่อ UX เท่านั้น |
| เก็บรูป | โฟลเดอร์ `uploads/` + `http.MaxBytesReader` + ตรวจ `http.DetectContentType` | ไม่ต้องมี cloud storage |
| **หน้าเว็บ** | **React + Vite + TypeScript** | มาตรฐาน, library เยอะ |
| แมพ / grid / token / fog | **react-konva** (MIT) | canvas 2D ลาก token ง่าย |
| state หน้าเว็บ | **Zustand** | เบา |
| type ระหว่าง Go ↔ TS | เขียน type ข้อความ (message) ไว้ที่เดียวใน Go แล้ว gen เป็น TS ด้วย **tygo** — หรือเขียนมือถ้ายังน้อย | กันสองฝั่งไม่ตรงกัน |
| **deploy** | `go build` → **binary ไฟล์เดียว** ที่ฝังหน้าเว็บไว้ด้วย `embed` | ก๊อปไฟล์เดียวไปรันที่ไหนก็ได้ |
| ทดสอบ | `go test` (เต๋า, กฎ, สิทธิ์) + **Playwright** (2 แท็บ DM + ผู้เล่น) | |

## โครงโฟลเดอร์
```
go.mod
cmd/server/main.go     เปิด HTTP + WebSocket, เสิร์ฟหน้าเว็บที่ embed ไว้
internal/room/         state ห้อง, เก้าอี้, ตัวละคร, ฉาก, token, fog + จัดการ event
internal/rules/        ruleset, preset D&D 5e (SRD), ตรวจแต้ม stat/สกิล
internal/dice/         parser + ทอยเต๋า (มี _test.go)
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
| `ruleset.update`, `ruleset.loadPreset` | DM |
| `character.submit` | ผู้เล่น (ตัวเอง; server ตรวจแต้ม stat/จำนวนสกิล) |
| `character.grant` (ไอเทม/สกิล/stat) | DM |
| `scene.add`, `scene.select`, `fog.set` | DM |
| `token.move` | ผู้เล่น (token ตัวเอง) / DM (ทุกตัว) |
| `dice.roll` | ทุกคน |

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
| เต๋า | เขียนเอง (~100 บรรทัด) | มี library |
| ภาษาที่ต้องใช้ | 2 (Go + TS) | 1 |

สรุป: Go ได้เปรียบเรื่อง deploy + ผู้ใช้ถนัด; ต้นทุนเพิ่มคือ gen type + reconnect + parser เต๋า — เล็ก

## ที่ไม่เลือก
Next.js (ไม่ต้องการ SSR) · WebRTC (โดนบล็อก, ต้องมี TURN) · Yjs/CRDT (เกินจำเป็น) · Supabase/Firebase (ยังไม่ต้องเก็บถาวร) · tldraw (ไม่ใช่ open source) · Go→WASM หน้าเว็บ (ไฟล์ใหญ่, ไม่มี library canvas/UI) · framework Go (Gin/Echo/Fiber — `net/http` พอ)

## ความปลอดภัย (ขั้นต่ำ)
- server ตรวจสิทธิ์ทุก event (ผู้เล่นขยับได้แค่ token ตัวเอง, คำสั่ง DM ต้องมาจากเก้าอี้ DM)
- อัปรูป: จำกัดขนาด (เช่น 10MB), รับเฉพาะ image (`DetectContentType`), ตั้งชื่อไฟล์ใหม่เอง (กัน path traversal)
- WebSocket: ตรวจ `Origin` + จำกัดขนาดข้อความ
- ห้องไม่มีรหัส (ยอมรับแล้ว) → ใครได้ URL เข้าได้
