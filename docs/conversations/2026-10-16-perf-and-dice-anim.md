# 2026-10-16 — หน้าสร้างตัวละครกระตุก + research animation ลูกเต๋า

## คำขอ (คำต่อคำ)
> เช็คหน้าเว็ตอัพตัวละครของทางฝั่ง player หน่อย กระตุกหนักมากๆ ทั้งที่ไม่มีอะไรเลย ลองเช็ค flow การทำงาน มันไม่ไ้เป็น live session อ่ะ ไม่ควรกระตุกขนาดนั้น

> research การทำ animation ลูกเต๋าหน่อย animation การทอยลุกเต๋างี้ มีมั้ย ถ้าไม่มีเอาเป้นวงล้อก็ได้เพื่อน

## สรุป
- กระตุก: สมมติฐาน = ต้นทุนการวาด (paint) จากธีม medieval — พื้นผิว SVG `feTurbulence` แบบ runtime + `background-attachment: fixed` + box-shadow หลายชั้นบนแผงยาว
  → agent `frontend` วัดก่อน/หลัง (เลื่อน, กดปุ่ม, พิมพ์, จำนวนข้อความ WS ตอนว่าง, CPU ช้าลง 4 เท่า) แล้วแก้
- animation เต๋า: agent `researcher` เทียบ library 3D (ต้องหยุดที่ผลจาก server), CSS 3D, วงล้อ → `docs/research/dice-animation.md`
  ข้อบังคับ: server เป็นคนทอย ทุกคนต้องเห็น animation จบที่ผลเดียวกัน, เบา, ข้ามได้, เคารพ reduced-motion
