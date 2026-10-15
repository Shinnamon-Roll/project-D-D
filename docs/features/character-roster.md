# Feature: ตัวละครสำเร็จรูป 10 คลาส + การ์ดเลื่อนเลือก + มานา (v0.10.0)

ที่มา: [`../conversations/2026-10-16-character-roster.md`](../conversations/2026-10-16-character-roster.md) · ทำหลัง i18n (ข้อความ 2 ภาษา)

## ระบบใหม่ (แทนที่ white-label ruleset)
- **ลบ**: แจกแต้ม stat, เลือกสกิลตอนสร้าง, RulesetEditor, preset 5E (SRD 5.2.1) + attribution, event `ruleset.*`, `character.submit` แบบเดิม, คำอธิบายสกิล SRD (`web/src/content/srdSkills.ts`)
- **Stat** (ทุกตัวมีชุดเดียวกัน ค่าต่างกัน):
  | key | ไทย | EN |
  |---|---|---|
  | `hp` (มี current/max) | เลือด | HP |
  | `str` | ความแข็งแกร่ง | Strength |
  | `end` | ความทนทาน | Endurance |
  | `agi` | ความว่องไว | Agility |
  | `mp` (มี current/max) | มานา | Mana |
- **Roster 10 ตัว** อยู่ใน Go `internal/rules/roster.go` (server เป็นเจ้าของข้อมูล): `{id (class), heroName{th,en}, className{th,en}, blurb{th,en}, art, stats{hpMax,str,end,agi,mpMax}, skills[3]{id, name{th,en}, desc{th,en}, cost}, items[]{name{th,en}, qty}}`
  - เน้นหลากหลาย: บทบาทต่างกัน (แทงก์, ดาเมจ, ฮีล, สนับสนุน, ควบคุม), stat กระจายไม่ซ้ำแพทเทิร์น, ของติดตัวไม่ซ้ำกันเลย, สกิลค่ามานา 1–5
  - รวมแต้ม stat ใกล้เคียงกัน (สมดุลคร่าวๆ) — ตัวมานาเยอะ HP น้อย และกลับกัน
  - 10 คลาสยึดตามรูปที่ researcher หาได้ (`docs/research/character-art.md`)

## Protocol
- `character.pick {classId}` — ผู้เล่นที่ยังไม่มีตัวละคร, สถานะ `creating`; คลาสต้องยังว่าง (ห้ามซ้ำ) → สร้าง Character จาก roster (hp=hpMax, mp=mpMax, items ตาม roster)
- `skill.use {skillId}` — เจ้าของตัวละคร, สถานะ `playing`; mp ≥ cost → mp −= cost; ไม่พอ → error
- `character.adjust {characterId, hp?, mp?}` — DM/co-DM, ค่า delta (บวก/ลบ), clamp 0..max
- `party.rest {}` — DM/co-DM: mp = mpMax ทุกตัว
- `character.grant` เดิม: ให้ไอเทม (ไม่มี stat-point/skill grant แบบเดิม); DM แก้ str/end/agi/hpMax/mpMax ได้
- state: `roster` (10 ตัว, ส่งครั้งเดียวใน state ก็ได้), `characters[]` ใหม่ (`classId`, `hp`, `hpMax`, `mp`, `mpMax`, `str`, `end`, `agi`, `skills`, `items`), `feed[]` บันทึกการกระทำล่าสุด 50 รายการ `{at, text{th,en}}` เช่น "Aldric ใช้ Shield Wall (−2 มานา)", "DM: พักผ่อน — มานาเต็ม"
- ห้ามมี clientId ใน feed/log

## หน้าเลือกตัวละคร (creating, ผู้เล่น)
- การ์ดใหญ่ (~320×440) **carousel เลื่อนซ้าย/ขวา**: ปัด (drag/swipe), ปุ่ม ‹ ›, คีย์ ←/→, จุดบอกตำแหน่ง
- **Motion**: library `motion` (motion.dev, `motion/react`) — การ์ดกลางเด่น scale 1, ข้างๆ scale ~0.85 + จาง + เอียงเล็กน้อย, spring ตอนเลื่อน; animate เฉพาะ transform/opacity (กฎ performance); `prefers-reduced-motion` → เปลี่ยนแบบไม่เคลื่อนไหว
- การ์ด: รูปคลาส, ชื่อฮีโร่ + คลาส, คำโปรย, แถบ stat 5 ค่า (bar), สกิล 3 อัน (ชื่อ + มานา + ⓘ คำอธิบาย), ของติดตัว
- ปุ่ม **"เลือกตัวนี้"** (data-testid `pick-<classId>`); ถูกเลือกแล้ว → การ์ดเทา + ป้าย "ถูกเลือกแล้ว: <ชื่อผู้เล่น>" ปุ่มกดไม่ได้
- DM หน้าเดียวกันเห็นว่าใครเลือกอะไร (ไม่มีปุ่มเลือก)
- ทุกคนเลือกครบ → playing (เหมือนเดิม)

## หน้าเล่น
- ชีท: รูป, HP bar (current/max), มานา bar, str/end/agi, สกิล + ปุ่ม **"ใช้"** (ปิดเมื่อมานาไม่พอ, ⓘ ดูคำอธิบาย), ของติดตัว
- DM: แต่ละตัวละคร HP/มานา −/+ , ปุ่ม **"พักผ่อน"** (เติมมานาทุกคน), ให้ไอเทม
- แผง **บันทึก/Feed** ล่าสุด (aria-live polite)
- token บนแมพใช้รูปคลาส (ไม่ต้องอัป portrait)

## รูป
- ตามผล research (license ฟรี, ตรงคลาส); ไฟล์ใน `web/src/assets/classes/`; credit ใน `docs/research/character-art.md` + หน้าเว็บ (footer เครดิต ถ้า license ต้องการ)

## Done
Go tests: pick (ห้ามซ้ำ, สถานะ, สิทธิ์), skill.use (มานาพอ/ไม่พอ/ไม่ใช่เจ้าของ), adjust clamp, rest · e2e: ผู้เล่นเลื่อนการ์ด (ปุ่ม + คีย์) เลือกตัว, ผู้เล่นคนที่ 2 เห็นการ์ดนั้นเทา, ใช้สกิลแล้วมานาลด + feed, DM พักผ่อน → มานาเต็ม · screenshot carousel 1440 + 390
