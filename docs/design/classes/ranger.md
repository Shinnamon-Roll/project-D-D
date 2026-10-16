# Ranger (พรานป่า)

นักล่าแห่งป่าเถื่อน ผสมฝีมืออาวุธ (มักเป็นธนู) กับเวทธรรมชาติแบบ half caster ใช้ Wisdom
หัวใจของคลาสคือ Hunter's Mark: ปักเป้าแล้วตามเก็บด้วยดาเมจเพิ่ม พร้อมทักษะสำรวจและสอดแนม · research: [`../../research/classes/phb-2024/ranger.md`](../../research/classes/phb-2024/ranger.md)

> กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) · คำอธิบายเขียนใหม่ด้วยคำของเราเอง

## ข้อมูลหลัก

| | |
|---|---|
| Hit Die | d10 |
| Primary ability | Dexterity, Wisdom |
| Saving throws | Strength, Dexterity |
| Armor / Weapons / Tools | Light, Medium armor + Shields / Simple + Martial weapons / ไม่มี |
| Skills | เลือก 3: Animal Handling, Athletics, Insight, Investigation, Nature, Perception, Stealth, Survival |
| Spellcasting | half caster · Wisdom · ใช้ Druidic Focus เป็น focus |

## Feature ตามเลเวล

| Level | Feature |
|---|---|
| 1 | Spellcasting, Favored Enemy, Weapon Mastery |
| 2 | Deft Explorer, Fighting Style |
| 3 | Subclass feature |
| 4 | Ability Score Improvement |
| 5 | Extra Attack |
| 6 | Roving |
| 7 | Subclass feature |
| 8 | Ability Score Improvement |
| 9 | Expertise |
| 10 | Tireless |
| 11 | Subclass feature |
| 12 | Ability Score Improvement |
| 13 | Relentless Hunter |
| 14 | Nature's Veil |
| 15 | Subclass feature |
| 16 | Ability Score Improvement |
| 17 | Precise Hunter |
| 18 | Feral Senses |
| 19 | Epic Boon |
| 20 | Foe Slayer |

## ตารางตามเลเวล + resources

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Favored Enemy | 2 | 2 | 2 | 2 | 3 | 3 | 3 | 3 | 4 | 4 | 4 | 4 | 5 | 5 | 5 | 5 | 6 | 6 | 6 | 6 |
| Prepared spells | 2 | 3 | 4 | 5 | 6 | 6 | 7 | 7 | 9 | 9 | 10 | 10 | 11 | 11 | 12 | 12 | 14 | 14 | 15 | 15 |
| Weapon Mastery | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |

| Resource | จำนวน | คืนเมื่อ |
|---|---|---|
| Favored Enemy (Hunter's Mark ฟรี) | ตามตาราง | Long Rest |
| Tireless (Temporary HP) | Wisdom mod, อย่างน้อย 1 (เลเวล 10+) | Long Rest — ยังไม่อยู่ใน JSON: schema กำหนดเลเวลเริ่มของสูตรไม่ได้ จดเอง |
| Nature's Veil | Wisdom mod, อย่างน้อย 1 (เลเวล 14+) | Long Rest — ยังไม่อยู่ใน JSON (เหตุผลเดียวกัน) |

Spell slots: ตาราง half caster ใน [`_shared.md`](_shared.md)

## JSON (machine-readable)

```json
{
  "id": "ranger",
  "name": {"th": "พรานป่า", "en": "Ranger"},
  "hitDie": 10,
  "primary": ["dex", "wis"],
  "saves": ["str", "dex"],
  "armor": {"th": "เกราะ Light, Medium และ Shield", "en": "Light and Medium armor and Shields"},
  "weapons": {"th": "อาวุธ Simple และ Martial", "en": "Simple and Martial weapons"},
  "tools": {"th": "", "en": ""},
  "skillChoices": {"count": 3, "from": ["animal-handling", "athletics", "insight", "investigation", "nature", "perception", "stealth", "survival"]},
  "caster": "half",
  "castAbility": "wis",
  "tables": [
    {"id": "favored-enemy", "name": {"th": "Favored Enemy", "en": "Favored Enemy"}, "values": ["2", "2", "2", "2", "3", "3", "3", "3", "4", "4", "4", "4", "5", "5", "5", "5", "6", "6", "6", "6"]},
    {"id": "prepared", "name": {"th": "เวทที่เตรียมได้", "en": "Prepared Spells"}, "values": ["2", "3", "4", "5", "6", "6", "7", "7", "9", "9", "10", "10", "11", "11", "12", "12", "14", "14", "15", "15"]},
    {"id": "weapon-mastery", "name": {"th": "Weapon Mastery", "en": "Weapon Mastery"}, "values": ["2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2"]}
  ],
  "resources": [
    {"id": "favored-enemy", "name": {"th": "Hunter's Mark ฟรี", "en": "Free Hunter's Mark"}, "max": [2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0}
  ],
  "features": [
    {"level": 1, "name": {"th": "Spellcasting", "en": "Spellcasting"}, "desc": {"th": "ร่ายเวทของ Ranger ด้วย Wisdom ใช้ Druidic Focus เป็น focus ได้ จำนวนเวทที่เตรียมดูตามตาราง ต้องเป็นเลเวลที่มี slot หลัง Long Rest สลับเวทที่เตรียมได้ 1 บท slot คืนครบตอน Long Rest", "en": "Cast Ranger spells using Wisdom; a Druidic Focus works as your focus. Prepared spell count follows the table and must fit your slot levels. After a Long Rest you may swap one prepared spell. Slots return on a Long Rest."}},
    {"level": 1, "name": {"th": "Favored Enemy", "en": "Favored Enemy"}, "desc": {"th": "มีเวท Hunter's Mark เตรียมไว้เสมอ (ไม่นับโควตา) ร่ายฟรีโดยไม่ใช้ slot ได้ตามจำนวนในตาราง (2 ถึง 6 ครั้ง) คืนหมดตอน Long Rest", "en": "Hunter's Mark is always prepared (outside your count). You can cast it without a slot a number of times shown in the table (2 up to 6), all regained on a Long Rest."}},
    {"level": 1, "name": {"th": "Weapon Mastery", "en": "Weapon Mastery"}, "desc": {"th": "เลือกอาวุธที่ถนัด 2 ชนิด ใช้ mastery property ของอาวุธนั้นได้ เปลี่ยนชนิดได้ทุกครั้งที่จบ Long Rest", "en": "Pick two kinds of weapons you are proficient with and use their mastery properties. You can change the picks after each Long Rest."}},
    {"level": 2, "name": {"th": "Deft Explorer", "en": "Deft Explorer"}, "desc": {"th": "ได้ Expertise ใน skill ที่ถนัดอยู่แล้ว 1 อย่าง และรู้ภาษาเพิ่ม 2 ภาษา ให้ DM/ผู้เล่นติ๊ก expertise ในชีทเอง", "en": "Gain Expertise in one skill you are already proficient in, and learn two extra languages. Mark the expertise on the sheet."}},
    {"level": 2, "name": {"th": "Fighting Style", "en": "Fighting Style"}, "desc": {"th": "เลือก Fighting Style feat 1 อย่าง หรือเลือก Druidic Warrior แทน: ได้ cantrip ของ Druid 2 อัน (นับเป็นเวท Ranger ใช้ Wisdom) เลเวลขึ้นสลับได้ 1 อัน ให้เพิ่มเองเป็น custom ability", "en": "Take a Fighting Style feat, or Druidic Warrior instead: learn two Druid cantrips that count as Ranger spells and use Wisdom; swap one on level up. Add your choice as a custom ability."}},
    {"level": 3, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "เลือก subclass ของ Ranger (เช่น Hunter) ให้ DM/ผู้เล่นเพิ่มความสามารถของ subclass เองเป็น custom ability", "en": "Choose a Ranger subclass (e.g. Hunter). The DM or player adds its abilities as custom abilities."}},
    {"level": 4, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 5, "name": {"th": "Extra Attack", "en": "Extra Attack"}, "desc": {"th": "เมื่อใช้ Attack action ในเทิร์นตัวเอง โจมตีได้ 2 ครั้งแทน 1 ครั้ง", "en": "When you take the Attack action on your turn, you attack twice instead of once."}},
    {"level": 6, "name": {"th": "Roving", "en": "Roving"}, "desc": {"th": "Speed +10 ft ถ้าไม่ได้ใส่ Heavy armor และได้ Climb Speed กับ Swim Speed เท่ากับ Speed ปกติ", "en": "Your Speed rises by 10 ft while not in Heavy armor, and you gain a Climb Speed and a Swim Speed equal to your Speed."}},
    {"level": 7, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 8, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 9, "name": {"th": "Expertise", "en": "Expertise"}, "desc": {"th": "ได้ Expertise เพิ่มอีก 2 skill ที่ถนัดแต่ยังไม่มี Expertise ให้ติ๊กในชีทเอง", "en": "Gain Expertise in two more proficient skills that lack it. Mark them on the sheet."}},
    {"level": 10, "name": {"th": "Tireless", "en": "Tireless"}, "desc": {"th": "Magic action ให้ Temporary HP ตัวเอง 1d8 + Wisdom mod (อย่างน้อย 1) ใช้ได้ Wisdom mod ครั้ง (อย่างน้อย 1) คืนตอน Long Rest (จดจำนวนเอง) และทุกครั้งที่จบ Short Rest ระดับ Exhaustion ลด 1", "en": "Magic action: gain 1d8 + Wisdom modifier (min 1) Temporary HP, usable Wisdom modifier times (min 1) per Long Rest (track by hand). Each Short Rest also lowers your Exhaustion by 1."}},
    {"level": 11, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 12, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 13, "name": {"th": "Relentless Hunter", "en": "Relentless Hunter"}, "desc": {"th": "โดนดาเมจแล้ว Concentration ของ Hunter's Mark ไม่หลุด", "en": "Taking damage can no longer break your Concentration on Hunter's Mark."}},
    {"level": 14, "name": {"th": "Nature's Veil", "en": "Nature's Veil"}, "desc": {"th": "Bonus Action ทำให้ตัวเองติด Invisible จนจบเทิร์นถัดไป ใช้ได้ Wisdom mod ครั้ง (อย่างน้อย 1) คืนตอน Long Rest (จดจำนวนเอง)", "en": "Bonus Action: become Invisible until the end of your next turn. Usable Wisdom modifier times (min 1), regained on a Long Rest (track by hand)."}},
    {"level": 15, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 16, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 17, "name": {"th": "Precise Hunter", "en": "Precise Hunter"}, "desc": {"th": "โจมตีเป้าที่ติด Hunter's Mark ของเราอยู่ได้ Advantage", "en": "You have Advantage on attack rolls against the creature marked by your Hunter's Mark."}},
    {"level": 18, "name": {"th": "Feral Senses", "en": "Feral Senses"}, "desc": {"th": "ได้ Blindsight ระยะ 30 ft", "en": "You gain Blindsight out to 30 ft."}},
    {"level": 19, "name": {"th": "Epic Boon", "en": "Epic Boon"}, "desc": {"th": "ได้ Epic Boon feat 1 อย่าง (หรือ feat อื่นที่เข้าเงื่อนไข) แนะนำ Boon of Dimensional Travel ให้เพิ่มเองเป็น custom ability และถ้าเพิ่ม ability ให้ DM แก้ค่า (Epic Boon เพิ่มได้ถึง 30)", "en": "Gain an Epic Boon feat (or another feat you qualify for); Boon of Dimensional Travel is suggested. Add it as a custom ability; if it raises an ability, the DM edits it (Epic Boons allow up to 30)."}},
    {"level": 20, "name": {"th": "Foe Slayer", "en": "Foe Slayer"}, "desc": {"th": "ลูกเต๋าดาเมจของ Hunter's Mark เปลี่ยนจาก d6 เป็น d10", "en": "Your Hunter's Mark damage die becomes a d10 instead of a d6."}}
  ]
}
```
