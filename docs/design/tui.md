# TUI launcher — design (v0.3)

`make run` → builds the web app, then runs `cmd/tabletop` (Bubble Tea). `cmd/server` stays as the headless server (Playwright, scripts).

## Screen 1 — Welcome
```
╭──────────────────────────────────────────────╮
│                                              │
│        T A B L E T O P   (ascii logo)        │
│      โต๊ะเล่นเกมเล่าเรื่องของแก๊งเพื่อน         │
│                                              │
│            ╭──────────────────╮              │
│            │  ▶  เริ่มเซิร์ฟเวอร์  │  ← focused   │
│            ╰──────────────────╯              │
│                                              │
│     v0.3.0 · http://127.0.0.1:8080           │
╰──────────────────────────────────────────────╯
        enter/space เริ่ม · q ออก
```
- Card and button centered in the terminal; re-centers on resize (`tea.WindowSizeMsg`).
- Button: Enter / Space / mouse click starts. Pressed state shows a spinner "กำลังเริ่ม…".
- If start fails (e.g. port in use) → red error line under the button; button becomes "ลองใหม่".

## Screen 2 — Running / logs
```
 ● RUNNING  http://127.0.0.1:8080   uptime 00:03:12   clients 3   v0.3.0
──────────────────────────────────────────────────────────────────────
 12:00:01 INFO  server started addr=127.0.0.1:8080
 12:00:09 INFO  ws connected seat=-
 12:00:12 INFO  event campaign.select by=Mira status=lobby
 12:00:20 WARN  event rejected type=token.move by=Bo err="you can only move your own token"
 ...                                         (scrollable viewport, auto-follow)
──────────────────────────────────────────────────────────────────────
 o เปิดเบราว์เซอร์ · ↑↓ pgup/pgdn เลื่อน · f ตามล่าสุด on · c ล้าง · q หยุด+ออก
```
- Logs come from the server's `slog` logger through a handler that sends each record to the TUI as a message (also kept in a ring buffer of the last 1000 lines).
- Colors by level: DEBUG dim, INFO default, WARN yellow, ERROR red; time dim.
- `q` / `ctrl+c` → graceful shutdown (`http.Server.Shutdown` with timeout), then exit.
- Never log clientIds (bearer secrets): log seat names / seat ids only.

## v0.6 dashboard (แทน Screen 2) — ที่มา [`../conversations/2026-10-16-tui-dashboard.md`](../conversations/2026-10-16-tui-dashboard.md)
แนวภาพ: พื้นดำ, ขอบเส้นบาง (rounded/normal border สีเทาเข้ม), หัวแผงตัวพิมพ์ใหญ่สีเทา + ตัวเลขใหญ่, accent เขียวมิ้นต์ (`#5EEAD4`/`#34D399`), bar `█░` และ sparkline `▁▂▃▄▅▆▇█`
```
╭ SERVER ───────────╮╭ LOGS ─────────────────────────── 412 · 3 warn ╮╭ CONNECTIONS ──────╮
│ ● RUNNING  v0.6.0 ││ 12:00:01 INFO  server started               ││ 3 online          │
│ local 127.0.0.1:… ││ 12:00:09 INFO  ws connected seat=Mira       ││ ▁▂▂▃▅▅▆▆▅▃ 25s    │
│ LAN  192.168.1.5:…││ 12:00:12 INFO  event campaign.select …      │╰───────────────────╯
│ uptime 00:03:12   ││ ...                                          │╭ ACTIVITY ─────────╮
╰───────────────────╯│                                              ││ events 128 · rej 2│
╭ ROOM ─────────────╮│                                              ││ ▁▁▃▇▅▂▁ /10s      │
│ LOBBY             ││                                              ││ dice 14 · last    │
│ menu▸LOBBY▸create…││                                              ││ Mira 1d20+3 = 17  │
│ Starter Set …     ││                                              │╰───────────────────╯
│ DM  Dee Em        ││                                              │╭ SYSTEM ───────────╮
│ seats ███░░░ 3/6  ││                                              ││ heap 6.1 MB       │
│ ready ██░░░░ 2/3  ││                                              ││ goroutines 21     │
│ chars 0 · scenes 1││                                              │╰───────────────────╯
╰───────────────────╯╰──────────────────────────────────────────────╯
 o เปิดเว็บ · ↑↓ pgup/pgdn เลื่อน · f ตามล่าสุด · c ล้าง · q หยุด+ออก
```
- **ซ้าย SERVER**: สถานะ, version, ลิงก์ local + LAN ทุกอัน, uptime
- **ซ้าย ROOM** (จาก snapshot ห้อง อ่านแบบ thread-safe ทุก tick): status + ขั้น menu▸lobby▸creating▸playing (ขั้นปัจจุบันเน้นสี), ชื่อแคมเปญ, ชื่อ DM (+ co-DM), เก้าอี้ผู้เล่นที่นั่ง/ทั้งหมด (bar), พร้อม/นั่ง (bar), จำนวนตัวละคร, จำนวนฉาก + ฉากปัจจุบัน
- **กลาง LOGS**: log เดิม (สีตาม level, scroll, follow) — หัวแผงบอกจำนวนบรรทัด + จำนวน warn/error
- **ขวา CONNECTIONS**: จำนวนคนต่อ + sparkline (sample ทุก 1s) แสดงเท่าที่พอดีความกว้างแผง (≈25 วินาทีในแผง 32 คอลัมน์)
- **ขวา ACTIVITY**: นับ event ที่สำเร็จ/ถูกปฏิเสธ (จาก log record), sparkline event ต่อ 10 วินาที, จำนวนทอยเต๋า + ทอยล่าสุด (ชื่อผู้ทอย, notation, total)
- **ขวา SYSTEM**: heap (runtime.MemStats), goroutines
- Responsive: กว้าง ≥ 130 → 3 คอลัมน์; 90–129 → ซ่อนคอลัมน์ขวา (ย้ายตัวเลขสำคัญไปหัวแผง log); < 90 → แถบสถานะบรรทัดเดียวด้านบน + log เต็ม
- ปุ่มเดิมทั้งหมด; ห้ามแสดง clientId
- mouse on for ⟳ RELOAD; select log text with shift+drag; wheel scrolls log
- reload ปิดอยู่ (หา repo root ไม่เจอ) → แผง/แถบสถานะแสดงแค่ "⟳ RELOAD off"; เหตุผลเต็มขึ้นใน log เมื่อกด `r`

## v0.6 reload — กด `r` / คลิก **⟳ RELOAD** (ที่แผง SERVER)
1. ถ้ามี client ต่ออยู่ → ขึ้น "ห้องจะรีเซ็ต กด r อีกครั้งเพื่อยืนยัน" (หมดเวลา 5 วิ)
2. รัน build ใน repo root (โฟลเดอร์ที่มี `Makefile`): web build + `go build` → `bin/tabletop`; output ไหลเข้าแผง LOGS; แผง SERVER แสดง "BUILDING…"
3. build fail → แสดง error, server เดิมรันต่อ
4. build ผ่าน → `Shutdown` server → ออกจาก Bubble Tea (คืน terminal) → `syscall.Exec` binary ใหม่ด้วย args เดิม + `-autostart`
5. `-autostart` → เริ่ม server ทันทีไม่ผ่านหน้ายินดีต้อนรับ
6. ไม่เจอ `Makefile`/repo → ปุ่ม reload ปิดใช้งาน พร้อมเหตุผล
Browser: state มี `build` (random ต่อ process) → client เห็นค่าเปลี่ยนหลังต่อใหม่ → `location.reload()`
