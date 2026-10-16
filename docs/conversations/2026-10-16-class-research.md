# 2026-10-16 — research คลาสจาก D&D Beyond

## คำขอผู้ใช้ (คำต่อคำ)

> scapping https://www.dndbeyond.com/classes?srsltid=AU7gw4Wia6bnHtKvGJ5L8HTFTtMHCouNDsFnMvBbRqu0I8GOEwlKzpRD คลาสจากหน้านี้ และดูดีเืล เพื่อดูว่า แต่ละคลาสมีอะไร นิสัยยังไง แล้วเก้บข้อมุลไว้ก่อนเป็นไฟล์ class รวตัวละครและข้อมูลตัวสละครทั้งหมด อย่าลืมแบ่งตามแพคด้วยนะ

## สรุป

- ดึงรายชื่อคลาสจากหน้า D&D Beyond Classes แล้วอ่านรายละเอียดแต่ละคลาส
- แต่ละคลาสเก็บ: มีอะไร (feature/กลไก), นิสัยและบุคลิก, subclass ทั้งหมด
- เก็บเป็นไฟล์ คลาสละ 1 ไฟล์ ใน `docs/research/classes/` + สารบัญที่แบ่งตามแพค (หนังสือที่มา เช่น PHB 2014, PHB 2024, Tasha's, Xanathar's)
- เป็นข้อมูล research เท่านั้น ยังไม่เอาเข้าแอป

## ข้อจำกัดที่ตั้งเอง

- เนื้อหา D&D Beyond เป็นลิขสิทธิ์ WotC: สรุปด้วยคำของเราเอง ไม่ copy ข้อความ feature คำต่อคำ (ชื่อ feature/subclass เป็นข้อเท็จจริง ใช้ได้)
- ToS ของ D&D Beyond ห้ามดึงข้อมูลอัตโนมัติจำนวนมาก: อ่านหน้า public แค่ไม่กี่หน้า และใช้ SRD 5.2.1 (CC BY 4.0) เป็นแหล่งหลักของกฎพื้นฐาน
- โรสเตอร์ 10 ตัวในแอปยังเป็นของเราเอง (original) ไม่เปลี่ยน

## ผู้ใช้ปรับ scope ระหว่างทำ (คำต่อคำ)

> this is detail of babarian https://www.dndbeyond.com/classes/2190875-barbarian scapfor basic class only

> Player’s Handbook only

## Scope ที่ใช้จริง

- เอาเฉพาะ 12 คลาสใน Player's Handbook (2024) เท่านั้น: ตัด Artificer, Blood Hunter และหนังสือเสริมอื่นออก
- เก็บเฉพาะ base class (feature ของคลาสหลัก) ไม่ลงรายละเอียด subclass: ใส่แค่ชื่อ subclass ใน PHB ไว้อ้างอิง 1 บรรทัด
- ไฟล์อยู่ที่ `docs/research/classes/phb-2024/<class>.md` (โฟลเดอร์ = แพค เผื่อเพิ่มแพคอื่นทีหลัง)
