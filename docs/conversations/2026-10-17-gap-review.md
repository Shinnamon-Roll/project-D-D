# 2026-10-17 — รีวิวระบบที่ยังขาด + flow แปลกๆ

## คำขอผู้ใช้ (คำต่อคำ)

> ลอง research ดูหน่อยว่ายังขาดระบบอะไรอีกบ้าง และflow การทำงานของเกมแปลกๆรึเปล่า รีวิวหน่อย

## สรุป

- รัน 2 agent คู่กัน: (1) อ่านโค้ดไล่ flow DM + ผู้เล่น (read-only) (2) `researcher` เทียบ SRD 5.2.1 + VTT อื่น
- ผลเต็ม: [`../research/gap-review-2026-10-17.md`](../research/gap-review-2026-10-17.md)
- ยังไม่ได้เขียนโค้ด ไม่ bump version (docs อย่างเดียว)

## ข้อค้นพบหลัก

- ไม่มี save: server restart / reload ใน TUI = ตัวละคร ฉาก token หายหมด เหลือแค่ notes, NPC pool, roster
- flow ติดตาย: ผู้เล่นหลุดตอนเลือกฮีโร่ = ห้องค้าง `creating`; หลังเริ่มเกมเพิ่ม/ลบที่นั่งไม่ได้; เลือกฮีโร่ใหม่ไม่ได้; ไม่มีปุ่มจบ session; DM เสีย `clientId` = คุมห้องไม่ได้
- ระบบเกมที่ขาด (เรียงความคุ้ม): death saves, temp HP, conditions + exhaustion, initiative (ขัด decision เดิม), วัดระยะ ft, concentration, NPC AC, Bloodied, handouts
- `docs/design/web-flow.md` ไม่ตรงกับโค้ดหลายจุด

## รอผู้ใช้ตัดสิน

Q18–Q20 ใน [`../open-questions.md`](../open-questions.md)
