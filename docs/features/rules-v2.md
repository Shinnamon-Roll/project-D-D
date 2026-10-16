# Feature: ระบบเกม v2 ตามข้อมูลคลาสอ้างอิง (6 ability, spell slot, เลเวล 1–20)

บทสนทนา: [`../conversations/2026-10-16-rules-v2.md`](../conversations/2026-10-16-rules-v2.md) · แทนที่ระบบ stat/มานาใน [`character-roster.md`](character-roster.md)

## ที่มาของข้อมูล + ลิขสิทธิ์
- กลไกและชื่อ feature ของ 12 คลาสอ้างอิง **SRD 5.2.1** (Wizards of the Coast, **CC BY 4.0**) ซึ่งมีครบทั้ง 12 คลาส (base class) — ใช้ได้ถูกกฎหมายเมื่อใส่ attribution
- คำอธิบาย feature ในแอปเขียนใหม่ด้วยคำของเราเอง สั้นๆ TH/EN (ไม่ copy ข้อความ SRD ยาวๆ, ไม่ใช้ข้อความจาก PHB หรือ D&D Beyond)
- ใส่ attribution ตามที่ CC BY กำหนด ในหน้าเลือกตัวละคร (ข้างเครดิตรูป) และ README:
  "This work includes material from the System Reference Document 5.2.1 ("SRD 5.2.1") by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd. The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International License, available at https://creativecommons.org/licenses/by/4.0/legalcode."
- ยังห้ามแบรนด์แอปว่า "D&D" และห้ามใช้ชื่อ "Player's Handbook" ใน UI
- อ้างอิงงาน research: [`../research/classes/phb-2024/`](../research/classes/phb-2024/)

## ข้อมูลคลาส (ใหม่) — `docs/design/classes.md` + JSON → Go (`internal/rules/classes.go`, มี sync test แบบ roster)
ต่อคลาส (12 คลาส: barbarian bard cleric druid fighter monk paladin ranger rogue sorcerer warlock wizard):
- `id`, `name{th,en}`, `hitDie` (6/8/10/12), `primary` (ability ids), `saves` (2 ability)
- `armor`, `weapons` (ข้อความสั้น TH/EN), `skillChoices` (รายการ skill ที่เลือกได้ + จำนวน)
- `caster`: `none | full | half | pact` + `castAbility`
- `tables`: ค่าตามเลเวล 1–20 ที่แสดงอย่างเดียว เช่น Rage Damage, Sneak Attack dice, Martial Arts die, Cantrips, Prepared spells, Invocations
- `resources`: counter ที่ใช้/คืนได้ `{id, name{th,en}, max: ตารางตามเลเวล | "abilityMod:<ab>" | "level×N", refill: "short"|"long"}` เช่น Rage, Bardic Inspiration, Channel Divinity, Wild Shape, Second Wind, Action Surge, Focus Points, Lay on Hands (pool), Sorcery Points, Favored Enemy, Arcane Recovery
- `features`: `{level, name{th,en}, desc{th,en}}` ทุก feature ของ base class เลเวล 1–20 (เลเวลที่เป็น subclass feature = แถว "Subclass feature" ให้เพิ่มเองผ่าน custom feature; ASI = แถว "Ability Score Improvement")
- ตารางกลาง: proficiency bonus ตามเลเวล, spell slot full/half/pact ตามเลเวล

## กติกาตัวละคร
- **6 ability**: STR DEX CON INT WIS CHA (1–30) → modifier = floor((score−10)/2)
- **Proficiency bonus** ตามเลเวล (+2 … +6)
- **Saving throws**: modifier + PB ถ้า proficient (จากคลาส)
- **18 skills** (Athletics=STR; Acrobatics, Sleight of Hand, Stealth=DEX; Arcana, History, Investigation, Nature, Religion=INT; Animal Handling, Insight, Medicine, Perception, Survival=WIS; Deception, Intimidation, Performance, Persuasion=CHA): modifier + PB ถ้า proficient, + 2×PB ถ้า expertise
- **HP**: เลเวล 1 = hit die เต็ม + CON mod; ต่อเลเวล = ค่าเฉลี่ยปัดขึ้น (d6=4, d8=5, d10=6, d12=7) + CON mod (ขั้นต่ำ 1)
- **AC**: ค่าตัวเลขต่อฮีโร่ (คำนวณจากของเริ่มต้นไว้ให้แล้ว), DM แก้ได้
- **Initiative** = DEX mod (แสดงอย่างเดียว) · **Passive Perception** = 10 + Perception
- **Spell slots**: ตามตาราง full/half/pact ตามเลเวลคลาส; ใช้ = กดช่อง slot (pip) ให้ทึบ, คืนตาม rest (pact = short rest)
- **คาถา/ความสามารถของฉัน**: ผู้เล่นเพิ่มรายการเองได้ (ชื่อ, เลเวลคาถา 0–9 หรือ "feature", คำอธิบาย ≤ 300) สูงสุด 40 — ไม่มีคลังคาถาในแอป (เล่นตามหนังสือ/การ์ดของโต๊ะ)
- **Hit dice**: จำนวน = เลเวล, ใช้ตอน short rest (ทอยเต๋าจริง แล้ว DM/ผู้เล่นใส่ HP ที่ฟื้น)
- **Rest** (DM กดให้ทั้งปาร์ตี้): **Short** = คืน counter แบบ short + pact slots · **Long** = HP เต็ม, slot ทั้งหมด, counter ทุกตัว, hit dice ทั้งหมด
- **เลเวล**: ค่าเริ่ม 1; DM กด เลเวล +/− ในการ์ดโปรไฟล์ → hpMax, PB, slots, counters, features ปรับตามตาราง (ASI = DM แก้ ability เอง)
- ลูกเต๋ายังเป็นของจริง: แอปแสดง modifier ที่ต้องบวก ไม่ทอยให้

## ฮีโร่ (สร้างข้อมูลใหม่)
- **วีรชนในตำนาน** 12 ตัว: เก็บชื่อ/เรื่อง/รูปเดิม แต่ข้อมูลเกมใหม่: `classId`, ability scores (standard array 15/14/13/12/10/8 + โบนัส +2/+1 ใส่ไว้แล้ว), skill proficiencies (ตามตัวเลือกคลาส, rogue/bard มี expertise), AC, ของเริ่มต้น (ตามคลาส แต่งชื่อเองได้)
- **ฮีโร่ชุดแรก** 10 ตัว: แปลงเป็นระบบเดียวกัน — knight → Fighter, necromancer → Wizard, ที่เหลือคลาสตรงชื่อ
- สกิลประจำตัวแบบเดิม (มานา) เลิกใช้ — ความสามารถมาจาก feature ของคลาส
- ข้อมูลใน `docs/design/roster.md` และ `docs/design/roster-legend.md` เขียนใหม่ทั้ง JSON (sync test เดิม)

## เอาออก
- มานา (`mpMax`, ปุ่มใช้สกิลหักมานา `skill.use`), stat STR/END/AGI แบบเดิม, สูตรสมดุล 26
- `party.rest` แบบ "มานาเต็ม" → แทนด้วย short/long rest
- `data/roster.json` รูปแบบเก่า: โหลดไม่ผ่าน → ย้ายไป `.bad-*` แล้วใช้ค่าเริ่มต้น (กลไก quarantine เดิม) + log เตือน

## Events (ร่าง — backend ยืนยันใน protocol.md)
- `character.pick {classId}` — เหมือนเดิม (heroId)
- `character.level {characterId, delta}` — DM, ±1, 1..20
- `character.setAbilities {characterId, abilities, ac}` — DM (แทน `character.setStats`)
- `character.adjust {characterId, hp?}` — เดิม (HP)
- `slot.toggle {characterId, level, used}` / `resource.adjust {characterId, resourceId, delta}` / `hitDice.adjust {characterId, delta}` — เจ้าของตัวละคร หรือ DM
- `party.rest {kind: "short"|"long"}` — DM
- `ability.add/update/remove` (custom spells/features ของผู้เล่น) — แทน `skill.addCustom`
- RoomState: ข้อมูลคลาส (คงที่) อาจส่งแค่ตอน hello ถ้าใหญ่เกิน (backend วัดขนาดก่อนตัดสิน)

## UI
- **การ์ดเลือกตัวละคร**: ชื่อ/คลาส/รูป/เรื่อง + 6 ability (score + mod), HP เลเวล 1, AC, Hit die, saves, skills เด่น, feature เลเวล 1, ของ
- **ชีทผู้เล่น** (แถบข้อมูล): หัว (ชื่อ, คลาส, เลเวล, HP, AC, PB, Initiative, Passive Perception) · แท็บ **Abilities** (6 ช่อง + saves + 18 skills) · **Resources** (spell slot pips ต่อเลเวล, counters +/−, hit dice) · **Features** (ตามเลเวล ≤ ปัจจุบัน, ดูรายละเอียดได้) · **คาถา/ความสามารถของฉัน** · **ของ**
- **การ์ดโปรไฟล์ DM**: เหมือนชีท + ปุ่ม เลเวล +/−, แก้ ability/AC, HP ±, ให้/ลบของ
- **หน้าเลือกตัวละครของ DM** (grid): ใช้การ์ดแบบใหม่
- **Ruleset editor**: แก้ฮีโร่ (ชื่อ, เรื่อง, รูป, ability, AC, skill proficiencies, ของ); ข้อมูลคลาสอ่านอย่างเดียว
- ธีมเดิม + กฎ perf เดิม, i18n th/en ครบ

## ลำดับงาน (commit ทีละขั้น)
1. ข้อมูลคลาส 12 คลาส (`classes.md` + Go + sync test) — content agents แบ่งคลาสกันเขียน
2. ข้อมูลฮีโร่ใหม่ 22 ตัว (roster.md + roster-legend.md)
3. Backend: model ตัวละคร v2 + events + tests
4. Frontend: การ์ด, ชีท, การ์ด DM, editor, i18n, e2e
5. Audit → แก้ → verify → commit (version MINOR 0.17.0)

## นอก scope (ไว้ใน ideas)
- Subclass, species, background, feats, คลังคาถา, death saves, conditions, multiclass, XP
- ทอยเต๋าในแอป (ยังเป็นเต๋าจริง)
