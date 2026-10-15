# Feature: แมพพอดีกรอบ + แบ่งจอ + จิ้มเดิน (v0.7.0)

ที่มา: [`../conversations/2026-10-16-map-fit.md`](../conversations/2026-10-16-map-fit.md)

## Layout หน้า Table
- เต็มความสูงจอใต้ topbar (`100dvh - topbar`); **ทั้งหน้าไม่ scroll** บน desktop
- ซ้าย = **แมพ** (ยืดหยุ่น) · กลาง = **แถบลาก** (`role="separator"`, `aria-orientation="vertical"`, `aria-valuenow` %, ลากด้วย pointer, ปุ่ม ←/→ ทีละ 2%, Home/End) · ขวา = **ข้อมูล** (ชีท, เต๋า, DM panel, รายการ token) scroll ภายในตัวเอง
- ขวากว้าง 22–60% ของจอ (ขั้นต่ำ 280px); จำค่าใน `localStorage` (try/catch)
- จอ < 900px: ซ้อนแนวตั้ง แมพบน (สูง ~55dvh) ข้อมูลล่าง ไม่มีแถบลาก

## แมพ
- Konva Stage ขนาด = กรอบแมพ (ResizeObserver) — เปลี่ยนตามการลาก/ย่อหน้าต่าง
- รูป = พื้นหลังเฉยๆ วาดแบบ **contain** (ทั้งรูปอยู่ในกรอบ ไม่ตัด ไม่ยืด) จัดกลาง
- **กริดของระบบคลุมทั้งรูป**: `cols` จาก DM, ช่องจัตุรัส `cell = imgW / cols` (หน่วย px ของรูป), `rows = max(1, round(imgH / cell))`, offset = 0
  - ส่ง `grid` เดิมใน protocol: `{cellSize: imgW/cols, offsetX: 0, offsetY: 0, cols, rows}` → server ไม่ต้องเปลี่ยนรูปแบบข้อมูล
- token = วงกลมกลางช่อง รัศมี ~0.42 ช่อง ชื่อใต้ token ย่อตามขนาดช่อง

## เดิน (จิ้ม/คลิก)
- **ผู้เล่น**: จิ้มช่องไหนก็ได้ → token ของตัวเองในฉากนี้เดินไป (`token.move`) · ไม่มี token ในฉาก = ไม่ทำอะไร
- **DM**: จิ้ม token → เลือก (วงทอง) → จิ้มช่อง → เดิน · จิ้ม token เดิม = ยกเลิกเลือก
- โหมด fog (DM): ลากสี่เหลี่ยม = ซ่อน, จิ้ม fog = เปิด (เหมือนเดิม)
- ใช้ pointer events (เมาส์ + นิ้ว); ไม่ zoom/pan ใน v0.7
- รายการ token (ช่องตัวเลข + ปุ่ม "ย้าย") ยังอยู่ในฝั่งข้อมูล เป็นทางเลือกสำหรับคีย์บอร์ด/screen reader

## DM ตั้งกริด
- ในแผง DM: **"กริด / Grid"** ปุ่ม `−` / `+` (data-testid `grid-cols-dec` / `grid-cols-inc`) + ตัวเลขคอลัมน์ (`grid-cols`), แถวแสดงอัตโนมัติ (`grid-rows`), preview สดบนแมพ, ปุ่มบันทึก
- คอลัมน์ 4–60; ฉากใหม่เริ่ม 20 คอลัมน์
- ลบช่อง cellSize/offsetX/offsetY เดิมออกจาก UI

## Server (internal/room)
- `scene.update` ที่กริดเล็กลง: token ที่อยู่นอกกริด → ดึงเข้าช่องขอบที่ใกล้สุด (clamp), fog → ตัดให้อยู่ในกริด, fog ที่หลุดทั้งหมดทิ้ง
- `cols`/`rows` จำกัด 1–100 (ตรวจของเดิม)

## Done
Go test clamp/clip · e2e: ผู้เล่นจิ้มช่องบน canvas แล้ว token ย้าย (DM เห็นช่องใหม่) · DM กด `+`/`−` แล้ว rows เปลี่ยน + บันทึก · แมพไม่มี scrollbar ที่ 1440x900 และ 390x844 · ลากแถบแบ่งแล้วแมพย่อ/ขยายตาม
