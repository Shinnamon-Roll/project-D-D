# 2026-10-16 — หน้าเลือกแคมเปญ + Starter Set ก่อน, พัก Salem

## คำขอ (คำต่อคำ)
> เวฟเอาไว้กา่อนเรื่อง salem หน้าแรกเลย สามารถเลือกเกมได้ ว่าจะเล่นแคมเปยอะไร อันแรก ขอเป้น starterset เอา starterset ก่อน เนื่องจากว่า อยากลองเล่นกันก่อน salem เขียนเป็น menu ไว้้ แบบ comming อะไรงี้ ส่วน starter ser มีเป็นบอรืดเกมอยุ่จริงๆ ลองหาข้อมูลหน่อย เอามาทั้งหมดเลย ทั้งเนื้อเรื่องทั้ังเซ้ตติ้ง รวมถึงคุ่มือด้วย

## สรุป
1. **พัก Salem** — เก็บร่างไว้ (`campaigns/salem/`) ยังไม่ทำต่อ
2. **หน้าแรก = เมนูเลือกแคมเปญ**
3. แคมเปญแรก = **D&D Starter Set** (อยากลองเล่นก่อน)
4. **Salem** แสดงในเมนูเป็น "Coming soon" กดไม่ได้
5. Research Starter Set ตัวจริง: เนื้อเรื่อง, setting, คู่มือ

## ข้อจำกัดที่แจ้งผู้ใช้
Starter Set เป็นลิขสิทธิ์ Wizards of the Coast (ไม่ใช่ SRD) → **ห้าม copy เนื้อเรื่อง/คู่มือ/แมพเต็มเข้า repo หรือแอป**
- เก็บได้: สรุปข้อเท็จจริง (โครงเรื่อง, setting, บท, สถานที่, NPC, มอนสเตอร์, ของในกล่อง)
- DM อ่านจากเล่มจริง (ซื้อ) หรือ "Intro to Stormwreck Isle" (ฟรีบน D&D Beyond)
- แอปให้ DM อัปแมพจากเล่มตัวเองไว้เล่นส่วนตัวได้ แต่ไม่ bundle
- คู่มือกฎ: ทำเองจาก SRD 5.2.1

## งานที่สั่ง
- agent `researcher` → `docs/research/starter-set.md` + `docs/campaigns/starter-set/README.md`

## คำตอบรอบถัดไป (คำต่อคำ)
> DM เลือก ไม่มีกล่อง starter set เลือกมาซักเซ็๖ อยากลองเล่น

- Q18 → **DM เลือกแคมเปญ**
- Q19 → ไม่มีกล่อง ให้เลือกให้ → เลือก **Dragons of Stormwreck Isle** เริ่มจาก **Intro to Stormwreck Isle** (ฟรีบน D&D Beyond, 2 encounter)
  - เหตุผล: ลองเล่นได้ฟรีไม่ต้องซื้อ, ต่อเข้าเรื่องหลักได้ ถ้าชอบค่อยซื้อฉบับเต็ม
  - ข้อเสีย: กฎ 2014 vs preset 2024 → รอ research สรุปการแปลง

## ผล research
- ข้อเท็จจริง: [`../research/starter-set.md`](../research/starter-set.md) — researcher แนะนำ DoSI ตรงกับที่เลือก
- แผนเล่นในแอป: [`../campaigns/starter-set/README.md`](../campaigns/starter-set/README.md)
- คำถามใหม่ Q20–Q21
