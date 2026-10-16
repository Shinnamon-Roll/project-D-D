# Feature: ตัวละครสำเร็จรูป 10 คลาส + การ์ดเลื่อนเลือก + มานา (v0.12.0)

> v0.15.0: roster แบ่งเป็นเซ็ตแล้ว (core 10 + legend 12, `RosterEntry.set`) และ DM เลือกเซ็ตที่ใช้เล่นใน lobby — ดู [`character-sets.md`](character-sets.md). ส่วนที่เขียนว่า "10 ตัว" ด้านล่างหมายถึงเซ็ต core

ที่มา: [`../conversations/2026-10-16-character-roster.md`](../conversations/2026-10-16-character-roster.md) · ทำหลัง i18n (ข้อความ 2 ภาษา)

## ระบบใหม่ (แทนที่ white-label ruleset)
- **ลบ**: แจกแต้ม stat, เลือกสกิลตอนสร้าง, RulesetEditor แบบเดิม (เปลี่ยนเป็นตัวแก้ roster — ด้านล่าง), preset 5E (SRD 5.2.1) + attribution, event `ruleset.*`, `character.submit` แบบเดิม, คำอธิบายสกิล SRD (`web/src/content/srdSkills.ts`)
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
- **สกิล 2 ประเภท** (คำขอเพิ่ม): `skills` ของตัวละคร = **ติดตัว** (`kind: "innate"`, จาก roster 3 อัน, แก้ไม่ได้) + **ของผู้เล่น** (`kind: "custom"`, สูงสุด 3 อัน)
  - `skill.addCustom {name, desc, cost}` — เจ้าของตัวละคร (สถานะ creating/playing): name 1–40, desc ≤ 300, cost 0–10 (int); เกิน 3 อัน = error
  - `skill.editCustom {characterId, skillId, name, desc, cost}` / `skill.removeCustom {characterId, skillId}` — เจ้าของ หรือ DM/co-DM
  - `skill.use` ใช้ได้ทั้งสองประเภท (หักมานาเหมือนกัน)
- `character.adjust {characterId, hp?, mp?}` — DM/co-DM, ค่า delta (บวก/ลบ), clamp 0..max
- `party.rest {}` — DM/co-DM: mp = mpMax ทุกตัว
- `character.grant` เดิม: ให้ไอเทม (ไม่มี stat-point/skill grant แบบเดิม); DM แก้ str/end/agi/hpMax/mpMax ได้
- **ให้ของ ย้ายเข้าการ์ดโปรไฟล์** (คำขอ dm-notes): ในการ์ดข้อมูลผู้เล่นมีฟอร์ม "ให้ของ" (ชื่อของ + จำนวน, พิมพ์เอง) → `item.give {characterId, name, qty}`; ลบหัวข้อ "ให้ของ" (GrantPanel) ออกจากเมนู DM
- state: `roster` (10 ตัว, ส่งครั้งเดียวใน state ก็ได้), `characters[]` ใหม่ (`classId`, `hp`, `hpMax`, `mp`, `mpMax`, `str`, `end`, `agi`, `skills`, `items`), `feed[]` บันทึกการกระทำล่าสุด 50 รายการ `{at, text{th,en}}` เช่น "Aldric ใช้ Shield Wall (−2 มานา)", "DM: พักผ่อน — มานาเต็ม"
- ห้ามมี clientId ใน feed/log

## หน้าเลือกตัวละคร (creating, ผู้เล่น)
- การ์ดใหญ่ (~320×440) **carousel เลื่อนซ้าย/ขวา**: ปัด (drag/swipe), ปุ่ม ‹ ›, คีย์ ←/→, จุดบอกตำแหน่ง
- **Motion**: library `motion` (motion.dev, `motion/react`) — การ์ดกลางเด่น scale 1, ข้างๆ scale ~0.85 + จาง + เอียงเล็กน้อย, spring ตอนเลื่อน; animate เฉพาะ transform/opacity (กฎ performance); `prefers-reduced-motion` → เปลี่ยนแบบไม่เคลื่อนไหว
- การ์ด: รูปคลาส, ชื่อฮีโร่ + คลาส, คำโปรย, แถบ stat 5 ค่า (bar), สกิล 3 อัน (ชื่อ + มานา + ⓘ คำอธิบาย), ของติดตัว
- หลังเลือกแล้ว: ฟอร์ม **"สกิลของฉัน"** เพิ่มสกิลเอง (ชื่อ, คำอธิบาย, มานา) สูงสุด 3 อัน (ข้ามได้ เพิ่มทีหลังในหน้าเล่นได้)
- ปุ่ม **"เลือกตัวนี้"** (data-testid `pick-<classId>`); ถูกเลือกแล้ว → การ์ดเทา + ป้าย "ถูกเลือกแล้ว: <ชื่อผู้เล่น>" ปุ่มกดไม่ได้
- DM หน้าเดียวกันเห็นว่าใครเลือกอะไร (ไม่มีปุ่มเลือก)
- ทุกคนเลือกครบ → playing (เหมือนเดิม)

## หน้าเล่น
- ชีท: รูป, HP bar (current/max), มานา bar, str/end/agi, สกิลติดตัว + สกิลของผู้เล่น (แยกหัวข้อ, เพิ่ม/แก้/ลบของตัวเองได้) + ปุ่ม **"ใช้"** (ปิดเมื่อมานาไม่พอ, ⓘ ดูคำอธิบาย), ของติดตัว
- DM: HP/มานา −/+ อยู่ในการ์ดข้อมูลผู้เล่น (ด้านล่าง), ปุ่ม **"พักผ่อน"** (เติมมานาทุกคน), ให้ไอเทม
- **DM: การ์ดข้อมูลผู้เล่น** (คำขอเพิ่ม 2): แถบ **ปาร์ตี้** = รูปโปรไฟล์กลมของทุกตัวละคร (data-testid `party-<characterId>`, ปุ่มจริง, แสดง HP/มานาเป็นเส้นเล็กใต้รูป) อยู่ด้านบนฝั่งแมพ
  - คลิก → **การ์ดลอยมุมขวาล่างของจอ** (`position: fixed`, data-testid `char-card`, ปิดด้วย × / Esc / คลิกรูปเดิม): รูป, ชื่อฮีโร่ + คลาส + ชื่อผู้เล่น, stat ทั้ง 5 (HP/มานา เป็น bar current/max พร้อมปุ่ม −/+), สกิลติดตัว + สกิลของผู้เล่น (ชื่อ, มานา, คำอธิบาย), **กระเป๋า** (ชื่อ × จำนวน)
  - DM ลบของ: ปุ่ม **−1** (`item-dec-<i>`) และ **ลบทั้งกอง** (`item-del-<i>`) → event `item.remove {characterId, index, qty}` (qty ≥ 1, เกินจำนวน = ลบทั้งกอง; DM/co-DM เท่านั้น) + บันทึกลง feed
  - การ์ดอัปเดตสดตาม state; ผู้เล่นไม่เห็นแถบนี้ (เห็นชีทตัวเองตามเดิม)
- แผง **บันทึก/Feed** ล่าสุด (aria-live polite)
- token บนแมพใช้รูปคลาส (ไม่ต้องอัป portrait)

## รูป
- ตามผล research (license ฟรี, ตรงคลาส); ไฟล์ใน `web/src/assets/classes/`; credit ใน `docs/research/character-art.md` + หน้าเว็บ (footer เครดิต ถ้า license ต้องการ)

## Done
Go tests: pick (ห้ามซ้ำ, สถานะ, สิทธิ์), custom skill (limit 3, ความยาว, cost 0–10, สิทธิ์เจ้าของ/DM), skill.use (มานาพอ/ไม่พอ/ไม่ใช่เจ้าของ), adjust clamp, rest, item.remove (−1, ทั้งกอง, เกินจำนวน, สิทธิ์) · e2e: ผู้เล่นเลื่อนการ์ด (ปุ่ม + คีย์) เลือกตัว, ผู้เล่นคนที่ 2 เห็นการ์ดนั้นเทา, ใช้สกิลแล้วมานาลด + feed, DM พักผ่อน → มานาเต็ม, DM คลิกรูปปาร์ตี้ → การ์ดมุมขวาล่าง → ลบของ −1 แล้วผู้เล่นเห็นของลด · screenshot carousel 1440 + 390

## หน้า Ruleset = ตัวแก้ Roster (คำขอเพิ่ม 3)
- เมนู DM หัวข้อ **ชุดกฎ / Ruleset** (และใน lobby) → รายการ 10 ตัวละคร (รูป + ชื่อ + คลาส) → คลิกเปิดตัวแก้:
  - ชื่อฮีโร่ / ชื่อคลาส / คำโปรย (TH + EN), รูป (เลือกจาก 10 รูปคลาส หรืออัปโหลดเอง)
  - stat: hpMax 1–999, str/end/agi 0–20, mpMax 0–99
  - สกิลติดตัว 1–5 อัน (ชื่อ/คำอธิบาย TH+EN, มานา 0–10)
  - ของเริ่มต้น 0–10 ชิ้น (ชื่อ TH+EN + จำนวน)
  - ปุ่ม **รีเซ็ตตัวนี้** / **รีเซ็ตทั้งหมด** เป็นค่าเริ่มต้น (`docs/design/roster.md`)
- Events (DM/co-DM, ทุกสถานะ): `roster.update {classId, entry}` · `roster.reset {classId?}` (ไม่ใส่ = ทั้งหมด)
- มีผลกับตัวละครที่ **ถูกเลือกหลังจากนั้น** เท่านั้น (ตัวที่เลือกไปแล้ว DM แก้ผ่านการ์ดโปรไฟล์)
- เก็บถาวร `data/roster.json` (atomic write แบบ notes); ไม่มีไฟล์ = roster เริ่มต้นใน Go; ไฟล์เสีย = quarantine แล้วใช้ค่าเริ่มต้น
- ผู้เล่นได้ roster ฉบับที่ใช้จริง (ไม่มีอะไรลับ) เพื่อแสดงการ์ด

## Protocol contract (ตายตัว — backend + frontend ทำคู่ขนานตามนี้)
```ts
type L = { th: string; en: string }                       // ข้อความ 2 ภาษา
type RosterSkill = { id: string; name: L; desc: L; cost: number }          // cost 0–10
type RosterItem  = { name: L; qty: number }
type RosterEntry = {
  id: string                    // classId: knight|wizard|rogue|ranger|cleric|paladin|barbarian|bard|druid|necromancer
  heroName: L; className: L; blurb: L; role: string
  art: string                   // "" = รูปคลาสตาม id (web/src/assets/classes/<id>.svg) | "/uploads/<name>"
  stats: { hpMax: number; str: number; end: number; agi: number; mpMax: number }
  skills: RosterSkill[]         // 1–5
  items: RosterItem[]           // 0–10
}
type CharSkill = { id: string; kind: 'innate' | 'custom'; name: L; desc: L; cost: number }
// custom: name/desc ใส่ภาษาเดียวกันทั้ง th/en (ผู้เล่นพิมพ์ครั้งเดียว)
type Character = {
  id: string; seatId: string; classId: string
  name: L                       // heroName ตอนเลือก
  art: string
  hp: number; hpMax: number; mp: number; mpMax: number; str: number; end: number; agi: number
  skills: CharSkill[]           // innate (จาก roster) + custom ≤ 3
  items: { name: L; qty: number }[]
}
type FeedEntry = { at: number; text: L }   // ล่าสุด 50
// RoomState: เพิ่ม roster: RosterEntry[], feed: FeedEntry[]; characters: Character[] (แบบใหม่)
// ลบ: ruleset, attribution, characters แบบเดิม (stats map / skillIds / inventory itemId)
```
Events (client → server):
| event | payload | ใคร | สถานะ |
|---|---|---|---|
| `character.pick` | `{classId}` | ผู้เล่นที่ยังไม่มีตัว | creating |
| `skill.addCustom` | `{name, desc, cost}` (name 1–40, desc ≤300, cost 0–10) | เจ้าของ | creating, playing |
| `skill.editCustom` | `{characterId, skillId, name, desc, cost}` | เจ้าของ / DM | creating, playing |
| `skill.removeCustom` | `{characterId, skillId}` | เจ้าของ / DM | creating, playing |
| `skill.use` | `{skillId}` | เจ้าของ | playing |
| `character.adjust` | `{characterId, hp?, mp?}` (delta −999..999) | DM | creating, playing |
| `character.setStats` | `{characterId, hpMax, mpMax, str, end, agi}` | DM | creating, playing |
| `party.rest` | `{}` | DM | playing |
| `item.give` | `{characterId, name, qty}` (name 1–40 → th=en=name, qty 1–99; ชื่อซ้ำ = บวกจำนวน) | DM | creating, playing |
| `item.remove` | `{characterId, index, qty}` | DM | creating, playing |
| `roster.update` | `{classId, entry: RosterEntry}` | DM | ทุกสถานะ |
| `roster.reset` | `{classId?}` | DM | ทุกสถานะ |

ลบ events: `ruleset.update`, `ruleset.loadPreset`, `character.submit`, `character.grant` · ผู้เล่นทุกคนเลือกครบ → `playing` · DM = DM หรือ co-DM
