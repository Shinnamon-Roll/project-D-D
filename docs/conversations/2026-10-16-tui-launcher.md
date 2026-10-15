# 2026-10-16 — TUI launcher (Bubble Tea)

## คำขอ (คำต่อคำ)
> use bubble tea TUI เพื่อ run พร้อมพัน โดยสร้างหน้าแรกเมื่อ make run จะโชวืหน้า ยินดีต้อนรับ จานั้นมีปุ่มกลางจอกดเพื่อเริ่ม แล้วก้จะโชว์ log ฝาก Design TUI ให้หด้วย

## ตีความ
- `make run` → เปิด TUI (Bubble Tea) แทนการรัน server เปล่า ๆ
- หน้าแรก: ยินดีต้อนรับ + ปุ่มกลางจอ → กดแล้วเริ่ม server (หน้าเว็บฝังอยู่ใน binary → frontend + backend รันพร้อมกัน)
- หลังเริ่ม: หน้า log สด
- ออกแบบ TUI ให้เอง (ดู `docs/design/tui.md`)
