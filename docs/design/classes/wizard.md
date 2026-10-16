# Wizard (พ่อมด/นักเวทศึกษา)

นักวิชาการที่เรียนเวทจากตำรา ใช้ Intelligence ร่ายเวท และจดเวทไว้ใน spellbook
รายชื่อเวทกว้างที่สุด เปลี่ยนเวทที่เตรียมได้ทุก Long Rest แต่ตัวบางมาก · research: [`../../research/classes/phb-2024/wizard.md`](../../research/classes/phb-2024/wizard.md)

> กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) · คำอธิบายเขียนใหม่ด้วยคำของเราเอง

## ข้อมูลหลัก

| | |
|---|---|
| Hit Die | d6 |
| Primary ability | Intelligence |
| Saving throws | Intelligence, Wisdom |
| Armor / Weapons / Tools | ไม่มี / Simple weapons / ไม่มี |
| Skills | เลือก 2: Arcana, History, Insight, Investigation, Medicine, Nature, Religion |
| Spellcasting | full caster · Intelligence (เตรียมเวทจาก spellbook) |

## Feature ตามเลเวล

| Level | Feature |
|---|---|
| 1 | Spellcasting, Ritual Adept, Arcane Recovery |
| 2 | Scholar |
| 3 | Subclass feature |
| 4 | Ability Score Improvement |
| 5 | Memorize Spell |
| 6 | Subclass feature |
| 8 | Ability Score Improvement |
| 10 | Subclass feature |
| 12 | Ability Score Improvement |
| 14 | Subclass feature |
| 16 | Ability Score Improvement |
| 18 | Spell Mastery |
| 19 | Epic Boon |
| 20 | Signature Spells |

## ตารางตามเลเวล + resources

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Cantrips | 3 | 3 | 3 | 4 | 4 | 4 | 4 | 4 | 4 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 |
| Prepared spells | 4 | 5 | 6 | 7 | 9 | 10 | 11 | 12 | 14 | 15 | 16 | 16 | 17 | 18 | 19 | 21 | 22 | 23 | 24 | 25 |
| เวทใน spellbook (ขั้นต่ำ) | 6 | 8 | 10 | 12 | 14 | 16 | 18 | 20 | 22 | 24 | 26 | 28 | 30 | 32 | 34 | 36 | 38 | 40 | 42 | 44 |

| Resource | จำนวน | คืนเมื่อ |
|---|---|---|
| Arcane Recovery | 1 | Long Rest |
| Signature Spells | 2 (เลเวล 20) | Short หรือ Long Rest |

Spell slots: ตาราง full caster ใน [`_shared.md`](_shared.md)

## JSON (machine-readable)

```json
{
  "id": "wizard",
  "name": {"th": "พ่อมด/นักเวทศึกษา", "en": "Wizard"},
  "hitDie": 6,
  "primary": ["int"],
  "saves": ["int", "wis"],
  "armor": {"th": "ไม่มี", "en": "None"},
  "weapons": {"th": "อาวุธ Simple", "en": "Simple weapons"},
  "tools": {"th": "", "en": ""},
  "skillChoices": {"count": 2, "from": ["arcana", "history", "insight", "investigation", "medicine", "nature", "religion"]},
  "caster": "full",
  "castAbility": "int",
  "tables": [
    {"id": "cantrips", "name": {"th": "Cantrips ที่รู้", "en": "Cantrips"}, "values": ["3", "3", "3", "4", "4", "4", "4", "4", "4", "5", "5", "5", "5", "5", "5", "5", "5", "5", "5", "5"]},
    {"id": "prepared", "name": {"th": "เวทที่เตรียมได้", "en": "Prepared Spells"}, "values": ["4", "5", "6", "7", "9", "10", "11", "12", "14", "15", "16", "16", "17", "18", "19", "21", "22", "23", "24", "25"]},
    {"id": "spellbook", "name": {"th": "เวทใน spellbook (ขั้นต่ำ)", "en": "Spellbook Spells (min)"}, "values": ["6", "8", "10", "12", "14", "16", "18", "20", "22", "24", "26", "28", "30", "32", "34", "36", "38", "40", "42", "44"]}
  ],
  "resources": [
    {"id": "arcane-recovery", "name": {"th": "Arcane Recovery", "en": "Arcane Recovery"}, "max": [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "signature-spells", "name": {"th": "Signature Spells", "en": "Signature Spells"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2], "maxFormula": "", "refill": "short", "shortRegain": 0, "shortFromLevel": 0}
  ],
  "features": [
    {"level": 1, "name": {"th": "Spellcasting", "en": "Spellcasting"}, "desc": {"th": "ร่ายเวท Wizard ด้วย Intelligence ใช้ Arcane Focus หรือ spellbook ได้ spellbook เริ่มมีเวทระดับ 1 หกบท เลเวลขึ้นเพิ่มได้ 2 บท เตรียมเวทจากหนังสือตามตาราง เปลี่ยนเวทที่เตรียมและ cantrip 1 บทได้ทุก Long Rest คัดลอกเวทที่เจอลงหนังสือได้ (2 ชม. + 50 GP ต่อระดับเวท)", "en": "Cast Wizard spells with Intelligence via an Arcane Focus or your spellbook. The book starts with six level 1 spells, +2 per level. Prepare from it per the table; each Long Rest re-pick prepared spells and swap one cantrip. Copy found spells in (2 hours + 50 GP per spell level)."}},
    {"level": 1, "name": {"th": "Ritual Adept", "en": "Ritual Adept"}, "desc": {"th": "ร่ายเวทที่มีแท็ก Ritual จาก spellbook แบบ Ritual ได้แม้ไม่ได้เตรียมไว้ แต่ต้องเปิดอ่านจากหนังสือ", "en": "Cast any Ritual-tagged spell in your spellbook as a Ritual even if it is not prepared, as long as you read it from the book."}},
    {"level": 1, "name": {"th": "Arcane Recovery", "en": "Arcane Recovery"}, "desc": {"th": "จบ Short Rest แล้วเลือกคืน spell slot ที่ใช้ไป รวมเลเวลไม่เกินครึ่งหนึ่งของเลเวล Wizard (ปัดขึ้น) และห้ามเป็น slot ระดับ 6 ขึ้นไป ใช้ได้ 1 ครั้ง คืนตอน Long Rest", "en": "After a Short Rest, recover expended slots whose levels add up to half your Wizard level (round up), none of level 6 or higher. Once per Long Rest."}},
    {"level": 2, "name": {"th": "Scholar", "en": "Scholar"}, "desc": {"th": "เลือก skill ที่มี proficiency อยู่แล้ว 1 อย่างจาก Arcana, History, Investigation, Medicine, Nature, Religion ได้ Expertise ใน skill นั้น", "en": "Pick one skill you are already proficient in from Arcana, History, Investigation, Medicine, Nature, or Religion; you gain Expertise in it."}},
    {"level": 3, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "เลือก subclass ของ Wizard (เช่น Evoker) ให้ DM/ผู้เล่นเพิ่มความสามารถของ subclass เองเป็น custom ability", "en": "Choose a Wizard subclass (e.g. Evoker). The DM or player adds its abilities as custom abilities."}},
    {"level": 4, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 5, "name": {"th": "Memorize Spell", "en": "Memorize Spell"}, "desc": {"th": "ทุกครั้งที่จบ Short Rest สลับเวทระดับ 1+ ที่เตรียมไว้ 1 บทกับเวทอื่นใน spellbook ได้", "en": "Whenever you finish a Short Rest, you can swap one prepared level 1+ spell for another spell from your spellbook."}},
    {"level": 6, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 8, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 10, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 12, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 14, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 16, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 18, "name": {"th": "Spell Mastery", "en": "Spell Mastery"}, "desc": {"th": "เลือกเวทระดับ 1 และระดับ 2 อย่างละบทใน spellbook (ร่ายด้วย action) เตรียมไว้เสมอ ร่ายที่ระดับต่ำสุดได้ไม่จำกัดโดยไม่ใช้ slot เปลี่ยนบทได้ทุก Long Rest", "en": "Choose one level 1 and one level 2 spell in your book with a casting time of an action. They are always prepared and you can cast them at their lowest level without a slot, at will. Swap them on a Long Rest."}},
    {"level": 19, "name": {"th": "Epic Boon", "en": "Epic Boon"}, "desc": {"th": "ได้ Epic Boon feat 1 อย่าง (หรือ feat อื่นที่เข้าเงื่อนไข) แนะนำ Boon of Spell Recall ให้เพิ่มเองเป็น custom ability และถ้าเพิ่ม ability ให้ DM แก้ค่า (Epic Boon เพิ่มได้ถึง 30)", "en": "Gain an Epic Boon feat (or another feat you qualify for); Boon of Spell Recall is suggested. Add it as a custom ability; if it raises an ability, the DM edits it (Epic Boons allow up to 30)."}},
    {"level": 20, "name": {"th": "Signature Spells", "en": "Signature Spells"}, "desc": {"th": "เลือกเวทระดับ 3 สองบทใน spellbook เตรียมไว้เสมอ ร่ายที่ระดับ 3 ได้บทละ 1 ครั้งโดยไม่ใช้ slot คืนตอน Short หรือ Long Rest", "en": "Pick two level 3 spells from your book as signature spells. They are always prepared, and each can be cast once at level 3 without a slot, regained on a Short or Long Rest."}}
  ]
}
```
