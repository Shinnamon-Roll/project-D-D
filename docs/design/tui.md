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
