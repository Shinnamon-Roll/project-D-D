# Sorcerer (นักเวทโดยกำเนิด)

ผู้ที่มีเวทมนตร์ไหลอยู่ในตัวตั้งแต่เกิด ไม่ต้องเรียนจากตำรา ใช้ Charisma ร่ายเวท
จุดเด่นคือ Sorcery Points ที่เอาไปดัดแปลงเวท (Metamagic) หรือแลกกับ spell slot ได้ · research: [`../../research/classes/phb-2024/sorcerer.md`](../../research/classes/phb-2024/sorcerer.md)

> กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) · คำอธิบายเขียนใหม่ด้วยคำของเราเอง

## ข้อมูลหลัก

| | |
|---|---|
| Hit Die | d6 |
| Primary ability | Charisma |
| Saving throws | Constitution, Charisma |
| Armor / Weapons / Tools | ไม่มี / Simple weapons / ไม่มี |
| Skills | เลือก 2: Arcana, Deception, Insight, Intimidation, Persuasion, Religion |
| Spellcasting | full caster · Charisma |

## Feature ตามเลเวล

| Level | Feature |
|---|---|
| 1 | Spellcasting, Innate Sorcery |
| 2 | Font of Magic, Metamagic |
| 3 | Subclass feature |
| 4 | Ability Score Improvement |
| 5 | Sorcerous Restoration |
| 6 | Subclass feature |
| 7 | Sorcery Incarnate |
| 8 | Ability Score Improvement |
| 10 | Metamagic (+2 options) |
| 12 | Ability Score Improvement |
| 14 | Subclass feature |
| 16 | Ability Score Improvement |
| 17 | Metamagic (+2 options) |
| 18 | Subclass feature |
| 19 | Epic Boon |
| 20 | Arcane Apotheosis |

## ตารางตามเลเวล + resources

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Sorcery Points | – | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
| Cantrips | 4 | 4 | 4 | 5 | 5 | 5 | 5 | 5 | 5 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 |
| Prepared spells | 2 | 4 | 6 | 7 | 9 | 10 | 11 | 12 | 14 | 15 | 16 | 16 | 17 | 17 | 18 | 18 | 19 | 20 | 21 | 22 |
| Metamagic options | – | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 6 | 6 | 6 | 6 |

| Resource | จำนวน | คืนเมื่อ |
|---|---|---|
| Sorcery Points | ตามตาราง (เลเวล 2+) | Long Rest (Sorcerous Restoration คืนบางส่วนตอน Short Rest) |
| Innate Sorcery | 2 | Long Rest |
| Sorcerous Restoration | 1 (เลเวล 5+) | Long Rest |

Spell slots: ตาราง full caster ใน [`_shared.md`](_shared.md)

## JSON (machine-readable)

```json
{
  "id": "sorcerer",
  "name": {"th": "นักเวทโดยกำเนิด", "en": "Sorcerer"},
  "hitDie": 6,
  "primary": ["cha"],
  "saves": ["con", "cha"],
  "armor": {"th": "ไม่มี", "en": "None"},
  "weapons": {"th": "อาวุธ Simple", "en": "Simple weapons"},
  "tools": {"th": "", "en": ""},
  "skillChoices": {"count": 2, "from": ["arcana", "deception", "insight", "intimidation", "persuasion", "religion"]},
  "caster": "full",
  "castAbility": "cha",
  "tables": [
    {"id": "sorcery-points", "name": {"th": "Sorcery Points", "en": "Sorcery Points"}, "values": ["—", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"]},
    {"id": "cantrips", "name": {"th": "Cantrips ที่รู้", "en": "Cantrips"}, "values": ["4", "4", "4", "5", "5", "5", "5", "5", "5", "6", "6", "6", "6", "6", "6", "6", "6", "6", "6", "6"]},
    {"id": "prepared", "name": {"th": "เวทที่เตรียมได้", "en": "Prepared Spells"}, "values": ["2", "4", "6", "7", "9", "10", "11", "12", "14", "15", "16", "16", "17", "17", "18", "18", "19", "20", "21", "22"]},
    {"id": "metamagic", "name": {"th": "Metamagic ที่รู้", "en": "Metamagic Options"}, "values": ["—", "2", "2", "2", "2", "2", "2", "2", "2", "4", "4", "4", "4", "4", "4", "4", "6", "6", "6", "6"]}
  ],
  "resources": [
    {"id": "sorcery-points", "name": {"th": "Sorcery Points", "en": "Sorcery Points"}, "max": [0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "innate-sorcery", "name": {"th": "Innate Sorcery", "en": "Innate Sorcery"}, "max": [2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "sorcerous-restoration", "name": {"th": "Sorcerous Restoration", "en": "Sorcerous Restoration"}, "max": [0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0}
  ],
  "features": [
    {"level": 1, "name": {"th": "Spellcasting", "en": "Spellcasting"}, "desc": {"th": "ร่ายเวทของ Sorcerer ด้วย Charisma ใช้ Arcane Focus ได้ จำนวน cantrip และเวทที่เตรียมดูตามตาราง ทุกครั้งที่เลเวลขึ้นสลับ cantrip ได้ 1 และสลับเวทที่เตรียมได้ 1 บท slot คืนครบตอน Long Rest", "en": "Cast Sorcerer spells using Charisma; an Arcane Focus works as your focus. Cantrips and prepared spells follow the table. Each time you level up you may swap one cantrip and one prepared spell. Slots return on a Long Rest."}},
    {"level": 1, "name": {"th": "Innate Sorcery", "en": "Innate Sorcery"}, "desc": {"th": "ใช้ Bonus Action ปลุกพลังในตัว 1 นาที ระหว่างนั้น spell save DC ของเวท Sorcerer +1 และทอยโจมตีด้วยเวท Sorcerer ได้ Advantage ใช้ได้ 2 ครั้ง คืนตอน Long Rest", "en": "Bonus Action: awaken your inner magic for 1 minute. While active, your Sorcerer spell save DC is +1 and your Sorcerer spell attacks have Advantage. Two uses, regained on a Long Rest."}},
    {"level": 2, "name": {"th": "Font of Magic", "en": "Font of Magic"}, "desc": {"th": "ได้ Sorcery Points ตามตาราง คืนครบตอน Long Rest ทุบ spell slot เป็นแต้มได้เท่าเลเวลของ slot หรือใช้ Bonus Action แปลงแต้มเป็น slot (เลเวล 1/2/3/4/5 ใช้ 2/3/5/6/7 แต้ม ต้องเป็น Sorcerer เลเวล 2/3/5/7/9) slot ที่สร้างหายตอน Long Rest", "en": "Gain Sorcery Points per the table, refilled on a Long Rest. Burn a slot for points equal to its level, or as a Bonus Action buy a slot: level 1/2/3/4/5 costs 2/3/5/6/7 points and needs Sorcerer level 2/3/5/7/9. Created slots vanish on a Long Rest."}},
    {"level": 2, "name": {"th": "Metamagic", "en": "Metamagic"}, "desc": {"th": "เลือก Metamagic 2 แบบ (เช่น Careful, Distant, Empowered, Extended, Heightened, Quickened, Seeking, Subtle, Transmuted, Twinned) จ่าย Sorcery Points เพื่อดัดแปลงเวท ปกติใส่ได้ 1 แบบต่อเวท เลเวลขึ้นสลับได้ 1 แบบ", "en": "Learn 2 Metamagic options (e.g. Careful, Distant, Empowered, Extended, Heightened, Quickened, Seeking, Subtle, Transmuted, Twinned). Spend Sorcery Points to tweak a spell; normally one option per spell. Swap one option on level up."}},
    {"level": 3, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "เลือก subclass ของ Sorcerer (เช่น Draconic Sorcery) ให้ DM/ผู้เล่นเพิ่มความสามารถของ subclass เองเป็น custom ability", "en": "Choose a Sorcerer subclass (e.g. Draconic Sorcery). The DM or player adds its abilities as custom abilities."}},
    {"level": 4, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 5, "name": {"th": "Sorcerous Restoration", "en": "Sorcerous Restoration"}, "desc": {"th": "จบ Short Rest แล้วคืน Sorcery Points ได้ไม่เกินครึ่งหนึ่งของเลเวล Sorcerer (ปัดลง) ใช้ได้ 1 ครั้ง คืนตอน Long Rest", "en": "After a Short Rest, regain Sorcery Points up to half your Sorcerer level (round down). Once per Long Rest."}},
    {"level": 6, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 7, "name": {"th": "Sorcery Incarnate", "en": "Sorcery Incarnate"}, "desc": {"th": "ถ้า Innate Sorcery หมดแล้ว จ่าย 2 Sorcery Points เพื่อเปิดใช้อีกได้ และระหว่างที่ Innate Sorcery ทำงาน ใส่ Metamagic ได้ถึง 2 แบบต่อเวท", "en": "With no Innate Sorcery uses left, you can still activate it by paying 2 Sorcery Points. While Innate Sorcery is active, you can apply up to two Metamagic options to each spell."}},
    {"level": 8, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 10, "name": {"th": "Metamagic", "en": "Metamagic"}, "desc": {"th": "เรียน Metamagic เพิ่มอีก 2 แบบ (รวม 4)", "en": "Learn 2 more Metamagic options (4 total)."}},
    {"level": 12, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 14, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 16, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 17, "name": {"th": "Metamagic", "en": "Metamagic"}, "desc": {"th": "เรียน Metamagic เพิ่มอีก 2 แบบ (รวม 6)", "en": "Learn 2 more Metamagic options (6 total)."}},
    {"level": 18, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 19, "name": {"th": "Epic Boon", "en": "Epic Boon"}, "desc": {"th": "ได้ Epic Boon feat 1 อย่าง (หรือ feat อื่นที่เข้าเงื่อนไข) แนะนำ Boon of Dimensional Travel ให้เพิ่มเองเป็น custom ability และถ้าเพิ่ม ability ให้ DM แก้ค่า (Epic Boon เพิ่มได้ถึง 30)", "en": "Gain an Epic Boon feat (or another feat you qualify for); Boon of Dimensional Travel is suggested. Add it as a custom ability; if it raises an ability, the DM edits it (Epic Boons allow up to 30)."}},
    {"level": 20, "name": {"th": "Arcane Apotheosis", "en": "Arcane Apotheosis"}, "desc": {"th": "ระหว่างที่ Innate Sorcery ทำงาน ใช้ Metamagic 1 แบบได้ฟรีทุกเทิร์นโดยไม่เสีย Sorcery Points", "en": "While Innate Sorcery is active, once on each of your turns you can apply one Metamagic option without spending Sorcery Points."}}
  ]
}
```
