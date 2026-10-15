# project-D-D — Tabletop

เว็บโต๊ะเล่นเกมเล่าเรื่อง (แนว D&D 5e) สำหรับเล่นกับเพื่อนผ่าน browser ไม่ต้องใช้กระดาษ
แนว white label: DM กำหนด stat / skill / ไอเทมเองได้ มี preset กฎ D&D 5e (SRD 5.2.1) ให้โหลด

**เวอร์ชัน:** 0.5.0 · **สถานะ:** เล่นได้บนเครื่องตัวเอง (localhost) — เมนูแคมเปญ, lobby, ตั้งกฎ, สร้างตัวละคร, แมพ + grid + token + fog, ทอยเต๋า, TUI launcher, หน้าตาโทน medieval, QR เข้านั่ง

---

## 1. ติดตั้งครั้งแรก
ต้องมี **Go 1.26+** และ **Node 24+**
```sh
git clone <repo> && cd project-D-D
cd web && npm install && cd ..
```

## 2. วิธีรัน
| คำสั่ง | ใช้เมื่อ | เปิดที่ |
|---|---|---|
| `make run` | **เล่นจริง** — เปิด TUI: หน้ายินดีต้อนรับ → กดปุ่มเริ่ม → ดู log สด | http://127.0.0.1:8080 (เครื่องเรา) · ลิงก์ LAN ขึ้นใน TUI |
| `make dev` | **พัฒนา** — Go server + Vite (แก้หน้าเว็บแล้วรีเฟรชเอง) | http://localhost:5173 |
| `make build` | สร้าง binary: `bin/tabletop` (TUI) + `bin/dnd` (server เปล่า ไม่มี TUI) | — |
| `./bin/tabletop` | รัน TUI จาก binary ที่ build แล้ว | http://127.0.0.1:8080 |
| `./bin/dnd` | รัน server แบบไม่มี TUI (log ออก terminal) | http://127.0.0.1:8080 |

**ค่าเริ่มต้นเปิดให้ Wi-Fi เดียวกันเข้าได้** (`-addr 0.0.0.0:8080`): เพื่อนเปิดลิงก์ LAN ที่ TUI แสดง หรือสแกน QR จาก DM

> ⚠️ **ความปลอดภัย:** ห้องไม่มีรหัสผ่าน ใครอยู่ Wi-Fi เดียวกันก็เปิดหน้าเว็บได้ กดเป็น DM ได้ถ้ากดก่อน และนั่งเก้าอี้ว่างได้
> ใช้ Wi-Fi สาธารณะ (ร้านกาแฟ ฯลฯ) ให้รัน `./bin/tabletop -addr 127.0.0.1:8080` (เครื่องเราเท่านั้น) · macOS อาจถามเรื่อง firewall ครั้งแรก ให้กด Allow

### ปุ่มใน TUI (`make run`)
| หน้า | ปุ่ม | ทำอะไร |
|---|---|---|
| ยินดีต้อนรับ | `enter` / `space` / คลิกปุ่ม | เริ่ม server (ถ้า port ไม่ว่าง จะขึ้น error + ปุ่ม "ลองใหม่") |
| ยินดีต้อนรับ | `q` | ออก |
| log | `o` | เปิด browser |
| log | `↑` `↓` `pgup` `pgdn` | เลื่อน log |
| log | `f` | เปิด/ปิดการเลื่อนตาม log ล่าสุด |
| log | `c` | ล้าง log |
| log | `q` / `ctrl+c` | หยุด server อย่างปลอดภัยแล้วออก |

## 3. วิธีเล่น
1. **DM** เปิดเว็บ → ใส่ชื่อ → เลือกแคมเปญ **Starter Set** (Salem = "Coming soon") → คนที่เลือกคือ DM
2. **Lobby** — DM กด **+ เก้าอี้ผู้เล่น** (สูงสุด 6 คน) / เปิดที่ DM ร่วม; เพื่อนเปิดลิงก์เดียวกัน → เลือกเก้าอี้ → ใส่ชื่อ → กด **พร้อม / Ready**
   - **QR:** DM กด **QR / สแกนเข้านั่ง** ที่เก้าอี้ว่าง → QR ขึ้นกลางจอพร้อม "ที่นั่ง N" → เพื่อนสแกนด้วยกล้องมือถือ → ใส่ชื่อ → กด **นั่ง / Sit**
3. (ถ้าต้องการ) DM แก้กฎในส่วน **ชุดกฎ / Ruleset**: stat, แต้มเริ่มต้น, skill, ไอเทม — หรือกด **โหลด preset 5E (SRD 5.2.1)**
4. DM กด **เริ่มเกม / Start** (กดได้เมื่อทุกคนที่นั่งพร้อมแล้ว)
5. **สร้างตัวละคร** — ทุกคน: ชื่อ, อัปรูป, กระจายแต้ม stat ให้หมด, เลือก skill ตามจำนวน → ส่ง; DM ให้ไอเทม/skill/แก้ stat ได้ตลอด
   (ใช้ตัวละครสำเร็จรูปจาก Starter Set: กรอกตามชีท แล้วให้ DM แก้ตัวเลขที่เหลือ เช่น HP/AC)
6. **แมพ** — DM อัปรูปแมพ → ตั้งขนาดช่อง grid ให้ตรงรูป → เลือกฉาก (ทุกคนย้ายไปฉากนั้น)
   - ผู้เล่น: คลิก token ตัวเอง → คลิกช่องปลายทาง
   - DM: ขยับได้ทุก token, เพิ่ม token NPC, โหมด fog (ลากสี่เหลี่ยม = ซ่อน, คลิก = เปิด)
7. **ทอยเต๋า** — พิมพ์ในช่องเต๋า ทุกคนเห็นผลใน log: `1d20+3`, `2d20kh1+5` (advantage), `2d20kl1` (disadvantage), `4d6kl3`, `d100`

เสียงคุยใช้ Discord · แมพ/เนื้อเรื่อง Starter Set ต้องมาจากของที่ DM มีเอง (แอปไม่รวมเนื้อหาลิขสิทธิ์) — ดู `docs/campaigns/starter-set/README.md`

## 4. ทดสอบ
```sh
make test                    # Go unit tests (dice, rules, ห้อง/สิทธิ์, upload, server, TUI)
cd web && npm run test:e2e   # Playwright: DM + ผู้เล่น เล่นครบ flow (build + รัน bin/dnd ที่ port 18080 เอง)
cd web && npm run build      # ตรวจ type + build หน้าเว็บ
```

## 5. ปัญหาที่เจอบ่อย
| อาการ | สาเหตุ / วิธีแก้ |
|---|---|
| สแกน QR แล้วมือถือเข้าไม่ได้ | มือถือต้องอยู่ Wi-Fi เดียวกับเครื่อง · ถ้า QR ขึ้นคำเตือน = server รันแบบ 127.0.0.1 → รันใหม่ด้วย `-addr 0.0.0.0:8080` · มีหลาย IP ให้กดเลือกที่อยู่อื่นในหน้าต่าง QR |
| TUI ขึ้น error ตอนกดเริ่ม | port 8080 ถูกใช้อยู่ → ปิดโปรแกรมนั้น หรือรัน `./bin/tabletop -addr 127.0.0.1:8081` |
| ปิด server แล้วเกมหาย | ตั้งใจ: v0 เก็บ state ใน memory · รูปที่อัปยังอยู่ใน `uploads/` |
| รีเฟรชหน้าแล้วกลับมาเก้าอี้เดิม | ตั้งใจ: browser จำ id ไว้ (เปิดหน้าต่างส่วนตัว = คนใหม่) |
| `make build` error เรื่อง embed | รัน `cd web && npm install` ก่อน แล้ว `make build` ใหม่ |
| อัปรูปไม่ได้ | ต้องนั่งเก้าอี้ก่อน, ไฟล์ ≤ 10MB, png/jpeg/gif/webp เท่านั้น |

## โครงโปรเจกต์
```
cmd/tabletop/         TUI launcher (Bubble Tea) → bin/tabletop
cmd/server/           server เปล่า (ไม่มี TUI) → bin/dnd
internal/server/      HTTP + WebSocket + upload
internal/room/        state ห้อง + ตรวจสิทธิ์ทุก event + hub
internal/rules/       ruleset, ตรวจตัวละคร, preset D&D 5e
internal/dice/        ตัวทอยเต๋า
internal/webui/       ฝังหน้าเว็บที่ build แล้วเข้า binary
web/                  React + Vite + TypeScript (+ e2e/ Playwright)
uploads/              รูปที่อัป (ไม่ commit)
docs/                 เอกสารทั้งหมด — สารบัญ docs/README.md
CLAUDE.md             กติกา + pipeline สำหรับ Claude Code
WORKLOG.md            บันทึกงานตามวันที่
.claude/agents/       researcher, backend, frontend, auditor
```

## กติกา
- แก้โค้ด → bump `version` ใน `web/package.json` (SemVer)
- เนื้อหากฎใช้เฉพาะ SRD 5.2.1 (CC-BY-4.0) พร้อม attribution; ไม่ใช้ชื่อ "D&D" เป็นชื่อแอป
- เครื่องมือ: **Graft/Trail** (`trail`) code graph ผ่าน MCP (cache `graft/` ไม่ commit)
