# 2026-10-16 — ตอบ Q20/Q21 + ขอข้อเสนอ tech stack

## คำขอ (คำต่อคำ)
> เอาตามลิสเก่าเลย อยากรุ้ก่อนว่าเกมมันเป็นยังไง techstack เป้นยังไง ไหนลองเสนอมาหน่อย

## ตีความ
- "เอาตามลิสเก่า / อยากรู้ก่อนว่าเกมเป็นยังไง" → **ไม่เพิ่มฟีเจอร์เข้า v0** (Q21: initiative + HP มอนสเตอร์ลับ ยังไม่ทำ) และ **ใช้ตัวละครสำเร็จรูปตามชีทเดิม** (Q20) — ลองเล่นก่อนแล้วค่อยเพิ่ม
- ขอข้อเสนอ tech stack → [`../design/tech-stack.md`](../design/tech-stack.md)

## รอบถัดไป (คำต่อคำ)
> how about go lang cuz I love go hahaha

- เปลี่ยนข้อเสนอเป็น **server = Go 1.26** (มีในเครื่อง), หน้าเว็บยัง React + TypeScript (browser รันได้แค่ JS)
- WebSocket ตรง (`coder/websocket`) แทน Socket.IO, parser เต๋าเขียนเอง, binary ไฟล์เดียว embed หน้าเว็บ
- รายละเอียด: [`../design/tech-stack.md`](../design/tech-stack.md) (v2)

## ยืนยัน (คำต่อคำ)
> frontend is good but I need backend in go

→ **ยืนยัน stack v2**: backend Go, frontend React + Vite + TypeScript
