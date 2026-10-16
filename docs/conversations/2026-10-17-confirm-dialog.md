# 2026-10-17 — popup ยืนยันก่อนลบ / ทำสิ่งสำคัญ

## คำขอผู้ใช้ (คำต่อคำ)

> if delete or do something that important create confirmmation popup to make sure

## คำถาม + คำตอบ

- รูปแบบ popup → **popup ในแอปตามธีม** (ไม่ใช้ `confirm()` ของเบราว์เซอร์); 3 จุดที่ใช้ `confirm()` อยู่แล้วเปลี่ยนมาใช้ตัวใหม่
- นอกจากการลบ ให้ถามด้วย → **พักสั้น/ยาว + จบ combat**, **เริ่มเกม**, **เปิดเผย token ที่ซ่อน** (ไม่เลือก: ลบหมอก)

## สรุป

- spec: [`../features/confirm-dialog.md`](../features/confirm-dialog.md)
- frontend อย่างเดียว, v0.18.1
