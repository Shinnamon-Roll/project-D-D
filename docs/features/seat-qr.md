# Feature: QR code เข้านั่งที่นั่ง (v0.5.0)

ที่มา: [`../conversations/2026-10-16-seat-qr.md`](../conversations/2026-10-16-seat-qr.md)

## Scope
- Lobby: DM / co-DM เห็นปุ่ม **"QR / สแกนเข้านั่ง"** บนเก้าอี้ผู้เล่นที่ว่าง
- กดแล้ว modal กลางจอ: QR ใหญ่ + label **"ที่นั่ง N / Seat N"** (N = ลำดับเก้าอี้ผู้เล่น เริ่ม 1) + URL ตัวหนังสือ + ปุ่มปิด (Esc ปิดได้)
- URL ใน QR: `http://<LAN IP>:<port>/?seat=<seatId>`
- เปิด URL นั้น (lobby): หน้าเลือกเก้าอี้นั้นไว้แล้ว + ช่องชื่อ focus → กด **"นั่ง / Sit"** → ส่ง `seat.sit` เดิม; นั่งสำเร็จแล้วลบ `?seat=` ออกจาก URL
  - เก้าอี้ถูกนั่งแล้ว/ไม่มีแล้ว/ห้องไม่ได้อยู่ lobby → ข้อความแจ้ง แล้วใช้หน้า lobby ปกติ
  - browser นี้นั่งอยู่แล้ว → ไม่สนใจ `?seat=`
- Network: `bin/tabletop` และ `bin/dnd` ค่าเริ่มต้น `-addr 0.0.0.0:8080`
- `GET /api/lan` → `200 {"urls": ["http://192.168.1.5:8080"]}`: IPv4 ที่ไม่ใช่ loopback ของเครื่อง + port ที่ฟังอยู่;
  ถ้าฟังแค่ loopback (เช่น `-addr 127.0.0.1:8080`) → `{"urls": []}` → modal ใช้ `location.origin` + เตือนว่ามือถือเข้าไม่ได้
- TUI: แสดงลิงก์ LAN ด้วย (นอกจาก localhost); `o` ยังเปิด `http://127.0.0.1:<port>`

## ไม่ทำ
สแกน QR ด้วยกล้องในเว็บ (ใช้แอปกล้องมือถือ), QR สำหรับเก้าอี้ DM ร่วม, รหัสผ่าน/token ในลิงก์

## Done
Go test ของ `/api/lan` · e2e: ผู้เล่นเปิด `/?seat=<id>` → ใส่ชื่อ → นั่งเก้าอี้นั้น · DM เปิด modal เห็น label + `<img>`/`<svg>` QR
