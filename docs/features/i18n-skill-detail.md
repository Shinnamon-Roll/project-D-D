# Feature: 2 ภาษา TH/EN + รายละเอียดสกิล (v0.8.0)

ที่มา: [`../conversations/2026-10-16-i18n-skill-detail.md`](../conversations/2026-10-16-i18n-skill-detail.md)
ลำดับงาน: ทำหลัง map-fit (แก้ไฟล์หน้าจอชุดเดียวกัน)

## i18n (frontend เท่านั้น ไม่เพิ่ม dependency)
- `web/src/i18n.ts`: dictionary `th` / `en` (key เดียวกัน, TypeScript บังคับว่าครบทั้งสองภาษา) + hook `useT()` + `lang` ใน Zustand
- ค่าเริ่มต้น `th`; จำใน `localStorage['lang']` (try/catch); ตั้ง `<html lang>` ตามภาษา
- ปุ่มสลับ **TH | EN** ที่ topbar (`data-testid` `lang-th` / `lang-en`, `aria-pressed`)
- ข้อความ UI ทุกหน้าจอผ่าน `t()` — เลิกแสดงแบบ "ไทย / English" คู่กัน
- ไม่แปล: ข้อมูลที่ DM/ผู้เล่นพิมพ์ (ชื่อ, ชื่อ stat/item/สกิลที่ DM สร้าง), ข้อความ attribution SRD (ต้องคงภาษาอังกฤษตามต้นฉบับ CC-BY) — แต่เพิ่มคำแปลไทยสั้นๆ ข้างใต้ได้
- error จาก server ยังเป็นภาษาอังกฤษ (ไม่อยู่ใน scope)
- e2e ใช้ `data-testid` อยู่แล้ว → ต้องไม่พังทั้งสองภาษา; เพิ่ม e2e สลับภาษา

## รายละเอียดสกิล
- Server: `SkillDef` เพิ่ม `description` (string, ≤ 500 ตัวอักษร, ไม่บังคับ) — `internal/rules` + validation + protocol.md
- RulesetEditor: ช่องคำอธิบายต่อสกิล (DM)
- Preset 5E: คำอธิบาย + ตัวอย่างการใช้ของ 18 สกิล **TH + EN** อยู่ฝั่ง frontend (`web/src/content/srdSkills.ts`) key ตาม skill id ของ preset
  - เนื้อหาเรียบเรียงจาก SRD 5.2.1 (ตาราง Skills: ability + example uses) — ห้ามใช้ข้อความจาก PHB/Basic Rules
  - ใช้เมื่อ skill id ตรง preset และ DM ไม่ได้เขียน `description` เอง (ข้อความ DM ชนะ)
- หน้าเลือกสกิล (CharacterCreate): ปุ่ม **ⓘ** ข้างแต่ละสกิล (`data-testid` `skill-info-<id>`, aria-label "รายละเอียด <ชื่อ>") → การ์ด `<dialog>` (Esc/ปุ่มปิด/คลิกนอก):
  - ชื่อสกิล, stat ที่ผูก
  - คำอธิบาย + ตัวอย่างการใช้
  - **วิธีใช้** (ลูกเต๋าจริง — แอปไม่มีระบบทอยแล้ว): เมื่อ DM ขอให้ทดสอบ → ทอย d20 จริง แล้วบวก modifier(stat) + proficiency (+2 ถ้าเลือกสกิลนี้); modifier = ⌊(ค่า stat − 10) / 2⌋; แสดงตัวอย่างจากค่า stat ปัจจุบัน เช่น "DEX 14 → +2, เลือกสกิลนี้ → ทอย d20 แล้ว **+4**"
  - ถ้า ruleset ไม่ใช่ preset 5E (ไม่มี stat ตาม key) → แสดงแค่คำอธิบาย + "ถาม DM ว่าทอยยังไง"
- (ไม่ทำรอบนี้) ปุ่ม ⓘ ในชีทตัวละครหน้าแมพ
