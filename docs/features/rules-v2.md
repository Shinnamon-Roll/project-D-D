# Feature: ระบบเกม v2 ตามข้อมูลคลาสอ้างอิง (6 ability, spell slot, เลเวล 1–20) — v0.18.0

บทสนทนา: [`../conversations/2026-10-16-rules-v2.md`](../conversations/2026-10-16-rules-v2.md) · แทนที่ระบบ stat/มานาใน [`character-roster.md`](character-roster.md) · protocol เต็ม: [`../design/protocol.md`](../design/protocol.md)

## ที่มาของข้อมูล + ลิขสิทธิ์
- กลไกและชื่อ feature ของ 12 คลาสอ้างอิง **SRD 5.2.1** (Wizards of the Coast, **CC BY 4.0**) ซึ่งมีครบทั้ง 12 คลาส (base class) — ใช้ได้ถูกกฎหมายเมื่อใส่ attribution
- คำอธิบาย feature ในแอปเขียนใหม่ด้วยคำของเราเอง สั้นๆ TH/EN (ไม่ copy ข้อความ SRD ยาวๆ, ไม่ใช้ข้อความจาก PHB, D&D Beyond หรือ Basic Rules)
- attribution ตามที่ CC BY กำหนด แสดงในหน้าเลือกตัวละคร (ข้างเครดิตรูป, ส่งมากับ `GET /api/classes` ช่อง `attribution`) และใน README — ข้อความตรงตัว:
  "This work includes material from the System Reference Document 5.2.1 ("SRD 5.2.1") by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd. The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International License, available at https://creativecommons.org/licenses/by/4.0/legalcode."
- ยังห้ามแบรนด์แอปว่า "D&D" และห้ามใช้ชื่อ "Player's Handbook" ใน UI
- อ้างอิงงาน research: [`../research/classes/phb-2024/`](../research/classes/phb-2024/)

## ไฟล์ข้อมูล
| ข้อมูล | ไฟล์ docs (ต้นฉบับ, มี JSON block) | Go (gen จาก docs) |
|---|---|---|
| 12 คลาส | [`../design/classes/<id>.md`](../design/classes/) (barbarian bard cleric druid fighter monk paladin ranger rogue sorcerer warlock wizard) | `internal/rules/classes_data.go` (type ใน `classes.go`) |
| ตารางกลาง (proficiency bonus, slot full/half/pact) | [`../design/classes/_shared.md`](../design/classes/_shared.md) | `internal/rules/classes.go` (เขียนมือ, test เทียบกับ docs) |
| ฮีโร่เซ็ต core 10 ตัว | [`../design/heroes-core.md`](../design/heroes-core.md) | `internal/rules/roster_default.go` |
| ฮีโร่เซ็ต legend 12 ตัว | [`../design/heroes-legend.md`](../design/heroes-legend.md) | `internal/rules/roster_legend.go` |
| ชื่อ + เนื้อเรื่องฮีโร่ (อ่านอย่างเดียว) | [`../design/roster.md`](../design/roster.md), [`../design/roster-legend.md`](../design/roster-legend.md) | — |

- sync test เทียบ Go กับ docs; แก้ docs แล้ว gen ไฟล์ `*_data.go`/`roster_*.go` ใหม่: `go test ./internal/rules -run MatchDoc -update`
- ข้อมูลคลาส (~136 KB) **ไม่อยู่ใน RoomState**: client โหลดครั้งเดียวจาก `GET /api/classes` (`{classes, skills, proficiency, fullSlots, halfSlots, pactSlots, attribution}`, ETag + `no-cache`)

ต่อคลาส: `id`, `name{th,en}`, `hitDie` (6/8/10/12), `primary`, `saves` (2 ability), `armor`/`weapons`/`tools` (TH/EN), `skillChoices {count, from}`,
`caster` (`none|full|half|pact`) + `castAbility`, `tables` (ค่าตามเลเวล 1–20 แสดงอย่างเดียว เช่น Rage Damage, Sneak Attack),
`resources` (counter: `max` ตารางตามเลเวล หรือ `maxFormula` `abilityMod:<ab>` / `level*<n>`, `refill` `short|long`, `shortRegain`, `shortFromLevel`, `fromLevel`),
`features` (`{level, name, desc}` ทุก feature ของ base class 1–20; subclass feature = แถว "Subclass feature" ให้ผู้เล่นเพิ่มเองในคาถาของฉัน; ASI = แถว "Ability Score Improvement")

## กติกาตัวละคร
- **6 ability**: STR DEX CON INT WIS CHA (1–30) → modifier = floor((score−10)/2)
- **Proficiency bonus** ตามเลเวล (+2 … +6)
- **Saving throws**: modifier + PB ถ้า proficient (save ของคลาส ∪ `extraSaves` ที่ DM เพิ่มจาก feature ผ่าน `character.setSkills`)
- **18 skills**: modifier + PB ถ้า proficient, + 2×PB ถ้า expertise
- **HP**: เลเวล 1 = hit die เต็ม + CON mod; ต่อเลเวล = ค่าเฉลี่ยปัดขึ้น (d6=4, d8=5, d10=6, d12=7) + CON mod (ขั้นต่ำเลเวลละ 1)
- **AC**: ค่าตัวเลขต่อฮีโร่ (คำนวณจากของเริ่มต้นไว้แล้ว), DM แก้ได้
- **Initiative** = DEX mod · **Passive Perception** = 10 + Perception (แสดงอย่างเดียว)
- **Spell slots**: ตาราง full/half/pact ตามเลเวล; ผู้เล่นกด pip ให้ทึบ = ใช้, กดช่องทึบ = คืน
- **Counters**: ตาม `resources` ของคลาส (+/−); max เปลี่ยนตามเลเวล/ability
- **คาถา/ความสามารถของฉัน**: ผู้เล่นเพิ่มเอง (ชื่อ 1–40, ชนิด spell ขั้น 0–9 หรือ feature, คำอธิบาย ≤ 300) สูงสุด 40 — ไม่มีคลังคาถาในแอป
- **Hit dice**: จำนวน = เลเวล; ผู้เล่นติ๊กว่าใช้ไปกี่ลูก, ทอยเต๋าจริง แล้ว **DM ใส่ HP ที่ฟื้นเอง** (HP ± ในการ์ดโปรไฟล์)
- **Rest** (DM กดให้ทั้งปาร์ตี้, สถานะ `playing`): **Short** = counter แบบ `short` (หรือถึง `shortFromLevel`) เต็ม, counter อื่นคืน `shortRegain`, pact slot เต็ม · **Long** = HP, slot, counter, hit dice เต็มทั้งหมด
- **เลเวล**: เริ่ม 1; DM กด −/+ ทีละ 1 (1..20) → hpMax, PB, slot, counter, hit dice, feature ปรับตามตาราง; HP ปัจจุบันเลื่อนตามส่วนต่าง hpMax (ลดเลเวล = **ดาเมจที่โดนยังอยู่** และ HP ถึง 0 = หมดสติได้); used ที่เกิน max ใหม่ถูกตัด · ASI = DM แก้ ability เอง
- ลูกเต๋ายังเป็นของจริง: แอปแสดง modifier ที่ต้องบวก ไม่ทอยให้

## ฮีโร่
- 22 ตัว (core 10 + legend 12): `id` เดิม, `classId` (knight → fighter, necromancer → wizard, ที่เหลือตรงชื่อ), ability (standard array 15/14/13/12/10/8 + โบนัส +2/+1), skill proficiency (ตามตัวเลือกคลาส, rogue มี expertise 2), AC, ของเริ่มต้นตามคลาส (ชื่อแต่งเอง)
- สกิลประจำตัวแบบมานาเลิกใช้ — ความสามารถมาจาก feature ของคลาส

## เอาออก
- มานา (`mpMax`, `skill.use`), stat STR/END/AGI แบบเดิม, สูตรสมดุล 26, `skill.addCustom/editCustom/removeCustom`, `character.setStats`
- `party.rest` แบบ "มานาเต็ม" → แทนด้วย short/long
- `data/roster.json` รูปแบบเก่า: decode ไม่ผ่าน → ย้ายเป็น `roster.json.bad-<random>` แล้วใช้ค่าเริ่มต้น + log เตือน

## Events (ตามโค้ด — รายละเอียดใน protocol.md)
ค่า used ส่งเป็น **absolute** (ไม่ใช่ delta) เพื่อกดซ้ำแล้วไม่เพี้ยน
- `character.pick {heroId}` — ผู้เล่น, สถานะ `creating`, ฮีโร่ของเซ็ตที่เปิด, ห้ามซ้ำ (`classId` ชื่อเดิมยังรับได้)
- `character.adjust {characterId, hp}` — DM, delta HP (±999), clamp 0..hpMax
- `character.level {characterId, delta}` — DM, delta = ±1, ผลต้องอยู่ 1..20
- `character.setAbilities {characterId, abilities, ac}` — DM, ability 1–30, AC 1–30
- `character.setSkills {characterId, skills, expertise, saves}` — DM, แทนทั้ง 3 รายการ (saves = save เพิ่มจาก feature)
- `slot.set {characterId, level: 1-9 | "pact", used}` / `resource.set {characterId, resourceId, used}` / `hitDice.set {characterId, used}` — เจ้าของ หรือ DM, `0 ≤ used ≤ max`
- `party.rest {kind: "short"|"long"}` — DM
- `ability.add {name, kind, level, desc}` / `ability.update {abilityId, …}` — เจ้าของ (ตัวละครตัวเอง) · `ability.remove {characterId, abilityId}` — เจ้าของ หรือ DM
- HTTP: `GET /api/classes` — ข้อมูลคลาส + ตารางกลาง + attribution

## UI
- **การ์ดเลือกตัวละคร**: ชื่อ/คลาส/รูป/เรื่อง + 6 ability (score + mod), HP เลเวล 1, AC, hit die, saves, skill ที่ถนัด, feature เลเวล 1 (ⓘ), ของ · footer: เครดิตรูป + attribution SRD
- **ชีทผู้เล่น**: หัว (ชื่อ, คลาส, เลเวล, HP, AC, PB, Initiative, Passive Perception) · แท็บ **Abilities** (ค่าสถานะ: 6 ability + saves + 18 skills) · **Resources** (ทรัพยากร: slot pips, counters, hit dice) · **Features** (ความสามารถถึงเลเวลปัจจุบัน + เลเวลถัดไป, ตารางคลาส) · **คาถาของฉัน** · **Items** (ของ)
- **การ์ดโปรไฟล์ DM** (แถบปาร์ตี้): ชีทเดียวกัน + เลเวล −/+, แก้ ability/AC, แก้ skill/expertise/save เพิ่ม, HP −/+, ให้/ลบของ · ปุ่ม **พักสั้น / พักยาว** ในแถบปาร์ตี้
- **Ruleset editor**: แก้ฮีโร่ (ชื่อ, เรื่อง, role, รูป, ability, AC, skill + expertise, ของ); ข้อมูลคลาสอ่านอย่างเดียว
- ธีมเดิม + กฎ perf เดิม, i18n th/en ครบ

## นอก scope (อยู่ใน [`../ideas.md`](../ideas.md))
- Subclass, species, background, feats, คลังคาถา, death saves, conditions, multiclass, XP
- การแปลงทรัพยากร (เช่น Font of Magic, Wild Resurgence), การคืน counter ตอนทอย initiative
- ทอยเต๋าในแอป (ยังเป็นเต๋าจริง)
