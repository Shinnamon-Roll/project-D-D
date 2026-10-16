# Feature: NPC Pool ของ DM (v0.11.0)

ที่มา: [`../conversations/2026-10-16-npc-pool.md`](../conversations/2026-10-16-npc-pool.md) · ทำหลัง dm-notes (ใช้ที่เก็บไฟล์ `data/` เดียวกัน) และก่อน roster

## Pool (white label)
- เมนู DM หัวข้อ **NPC** เปลี่ยนเป็น **NPC Pool**: ปุ่ม **+** → ฟอร์ม: **รูป** (อัปโหลด, ไม่ใส่ = รูป silhouette เริ่มต้น), **ชื่อ** (1–40), **ปกปิด** (สวิตช์ ค่าเริ่มต้นตอนวาง), **เลือด** (HP สูงสุด 1–999) — แค่นี้
- รายการ pool: รูปเล็ก + ชื่อ + HP + ไอคอนปกปิด; แก้ไข / ลบ / **วางบนแมพ**
- เก็บถาวร `data/npcs.json` (atomic write แบบเดียวกับ notes), สูงสุด 100 ตัว; ส่งเฉพาะ DM/co-DM
- Events (DM/co-DM, ทุกสถานะห้อง): `npc.add {name, imageUrl, hpMax, hidden}` · `npc.update {id, name, imageUrl, hpMax, hidden}` · `npc.remove {id}` (token ที่วางไว้แล้วอยู่ต่อ)
- `imageUrl` ต้องเป็น `/uploads/<ชื่อไฟล์>` ที่ server สร้าง หรือว่าง

## วางบนแมพ
- `npc.place {npcId, sceneId}` → token ใหม่ในฉาก (ช่องว่างแรก), `hp = hpMax`, `hidden` = ค่าจาก pool, รูปจาก pool; วางตัวเดียวกันซ้ำได้ (ชื่อ token ต่อท้าย #2, #3)
  - กฎชื่อ token ไม่ซ้ำ + ปุ่ม Duplicate: ดู [npc-names.md](../features/npc-names.md)
- แทนที่ event เดิม `token.addNpc` (ชื่ออย่างเดียว)
- DM ต่อ token NPC: **ซ่อน/เปิดเผย** (`token.setHidden {tokenId, hidden}`), **HP −/+** (`token.adjustHp {tokenId, delta}` clamp 0..hpMax), ลบ (เดิม) — อยู่ในรายการ token ฝั่งข้อมูล
- บนแมพ: token NPC ใช้รูป; HP 0 = token จาง + กากบาท; DM เห็น token ที่ซ่อนแบบโปร่ง + ไอคอนตาปิด

## ความลับ (server เป็นคนกรอง)
- state ที่ส่งให้ **ผู้เล่น**: ไม่มี token ที่ `hidden: true`, token NPC ไม่มี field `hp`/`hpMax`, ไม่มี `npcs` (pool)
- DM/co-DM ได้ครบ
- log ไม่มีชื่อ NPC ที่ซ่อน (log แค่ id)

## Done
Go tests: สิทธิ์, limits, imageUrl ต้องเป็น uploads, persist/reload, ผู้เล่นไม่ได้รับ hidden token / hp / pool, place ซ้ำ + ชื่อ #n, adjustHp clamp · e2e: DM + NPC (ใส่รูป, ชื่อ, ปกปิด, HP) → วาง → ผู้เล่นไม่เห็น (และ WS frame ไม่มีชื่อ NPC นั้น) → DM เปิดเผย → ผู้เล่นเห็น → DM −HP

## แก้ไข (คำขอเพิ่ม): ไอคอนตาเป็น toggle
- ปุ่มไอคอนตาอย่างเดียว (ไม่มีข้อความ, `aria-label` บอกสถานะ + `aria-pressed`): ตาเปิด = มองเห็น, ตาปิด = ซ่อน; คลิกสลับ
  - แถว token NPC (`token-hide-<tokenId>`) → `token.setHidden`
  - แถว NPC Pool (`npc-hidden-<id>`) → `npc.update` เปลี่ยนค่าเริ่มต้น hidden
- บนแมพ **ไม่มีไอคอนตา** — token ที่ซ่อนแสดงแบบโปร่ง (~45%) สำหรับ DM เท่านั้น
