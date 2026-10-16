# Warlock (ผู้ทำสัญญา)

ผู้ได้พลังเวทจากสัญญากับสิ่งเหนือธรรมชาติ (patron) ใช้ Charisma ร่ายเวท
ใช้ Pact Magic: slot น้อยแต่เป็นเลเวลสูงเท่ากันทุกช่องและคืนตอน Short Rest ปรับตัวละครด้วย Eldritch Invocations · research: [`../../research/classes/phb-2024/warlock.md`](../../research/classes/phb-2024/warlock.md)

> กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) · คำอธิบายเขียนใหม่ด้วยคำของเราเอง

## ข้อมูลหลัก

| | |
|---|---|
| Hit Die | d8 |
| Primary ability | Charisma |
| Saving throws | Wisdom, Charisma |
| Armor / Weapons / Tools | Light armor / Simple weapons / ไม่มี |
| Skills | เลือก 2: Arcana, Deception, History, Intimidation, Investigation, Nature, Religion |
| Spellcasting | pact (Pact Magic) · Charisma |

## Feature ตามเลเวล

| Level | Feature |
|---|---|
| 1 | Eldritch Invocations, Pact Magic |
| 2 | Magical Cunning |
| 3 | Subclass feature |
| 4 | Ability Score Improvement |
| 6 | Subclass feature |
| 8 | Ability Score Improvement |
| 9 | Contact Patron |
| 10 | Subclass feature |
| 11 | Mystic Arcanum (level 6 spell) |
| 12 | Ability Score Improvement |
| 13 | Mystic Arcanum (level 7 spell) |
| 14 | Subclass feature |
| 15 | Mystic Arcanum (level 8 spell) |
| 16 | Ability Score Improvement |
| 17 | Mystic Arcanum (level 9 spell) |
| 19 | Epic Boon |
| 20 | Eldritch Master |

## ตารางตามเลเวล + resources

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Invocations | 1 | 3 | 3 | 3 | 5 | 5 | 6 | 6 | 7 | 7 | 7 | 8 | 8 | 8 | 9 | 9 | 9 | 10 | 10 | 10 |
| Cantrips | 2 | 2 | 2 | 3 | 3 | 3 | 3 | 3 | 3 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 |
| Prepared spells | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 10 | 11 | 11 | 12 | 12 | 13 | 13 | 14 | 14 | 15 | 15 |
| Pact slots | 1 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 3 | 3 | 3 | 3 | 3 | 3 | 4 | 4 | 4 | 4 |
| Slot level | 1 | 1 | 2 | 2 | 3 | 3 | 4 | 4 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 |

| Resource | จำนวน | คืนเมื่อ |
|---|---|---|
| Magical Cunning | 1 (เลเวล 2+) | Long Rest |
| Contact Patron | 1 (เลเวล 9+) | Long Rest |
| Mystic Arcanum ระดับ 6 / 7 / 8 / 9 | 1 ต่อบท (เลเวล 11 / 13 / 15 / 17+) | Long Rest |

Pact slots: ตาราง pact ใน [`_shared.md`](_shared.md) (คืนตอน Short Rest)

## JSON (machine-readable)

```json
{
  "id": "warlock",
  "name": {"th": "ผู้ทำสัญญา", "en": "Warlock"},
  "hitDie": 8,
  "primary": ["cha"],
  "saves": ["wis", "cha"],
  "armor": {"th": "เกราะเบา (Light)", "en": "Light armor"},
  "weapons": {"th": "อาวุธ Simple", "en": "Simple weapons"},
  "tools": {"th": "", "en": ""},
  "skillChoices": {"count": 2, "from": ["arcana", "deception", "history", "intimidation", "investigation", "nature", "religion"]},
  "caster": "pact",
  "castAbility": "cha",
  "tables": [
    {"id": "invocations", "name": {"th": "Eldritch Invocations", "en": "Eldritch Invocations"}, "values": ["1", "3", "3", "3", "5", "5", "6", "6", "7", "7", "7", "8", "8", "8", "9", "9", "9", "10", "10", "10"]},
    {"id": "cantrips", "name": {"th": "Cantrips ที่รู้", "en": "Cantrips"}, "values": ["2", "2", "2", "3", "3", "3", "3", "3", "3", "4", "4", "4", "4", "4", "4", "4", "4", "4", "4", "4"]},
    {"id": "prepared", "name": {"th": "เวทที่เตรียมได้", "en": "Prepared Spells"}, "values": ["2", "3", "4", "5", "6", "7", "8", "9", "10", "10", "11", "11", "12", "12", "13", "13", "14", "14", "15", "15"]}
  ],
  "resources": [
    {"id": "magical-cunning", "name": {"th": "Magical Cunning", "en": "Magical Cunning"}, "max": [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "contact-patron", "name": {"th": "Contact Patron", "en": "Contact Patron"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "mystic-arcanum-6", "name": {"th": "Mystic Arcanum (ระดับ 6)", "en": "Mystic Arcanum (level 6)"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "mystic-arcanum-7", "name": {"th": "Mystic Arcanum (ระดับ 7)", "en": "Mystic Arcanum (level 7)"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "mystic-arcanum-8", "name": {"th": "Mystic Arcanum (ระดับ 8)", "en": "Mystic Arcanum (level 8)"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "mystic-arcanum-9", "name": {"th": "Mystic Arcanum (ระดับ 9)", "en": "Mystic Arcanum (level 9)"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0}
  ],
  "features": [
    {"level": 1, "name": {"th": "Eldritch Invocations", "en": "Eldritch Invocations"}, "desc": {"th": "เลือก invocation 1 อย่าง (เช่น Pact of the Blade, Pact of the Chain, Pact of the Tome, Agonizing Blast) จำนวนเพิ่มตามตาราง บางอย่างมีเงื่อนไขเลเวล เลเวลขึ้นสลับได้ 1 อย่าง ให้เพิ่มแต่ละอย่างเองเป็น custom ability", "en": "Pick 1 invocation (e.g. Pact of the Blade, Pact of the Chain, Pact of the Tome, Agonizing Blast); more per the table, some need a minimum level. Swap one on level up. Add each as a custom ability."}},
    {"level": 1, "name": {"th": "Pact Magic", "en": "Pact Magic"}, "desc": {"th": "ร่ายเวทของ Warlock ด้วย Charisma ใช้ Arcane Focus ได้ pact slot ทุกช่องเป็นเลเวลเดียวกัน (ดูตาราง) ร่ายเวทเลเวลต่ำด้วย slot นี้จะกลายเป็นเลเวลของ slot คืนครบตอน Short หรือ Long Rest เลเวลขึ้นสลับ cantrip และเวทได้อย่างละ 1", "en": "Cast Warlock spells with Charisma; an Arcane Focus works. All pact slots share one level (see table) and lower spells cast with them are upcast to it. Slots return on a Short or Long Rest. Swap one cantrip and one spell on level up."}},
    {"level": 2, "name": {"th": "Magical Cunning", "en": "Magical Cunning"}, "desc": {"th": "ทำพิธี 1 นาที แล้วคืน pact slot ที่ใช้ไปได้ไม่เกินครึ่งหนึ่งของจำนวนสูงสุด (ปัดขึ้น) ใช้ได้ 1 ครั้ง คืนตอน Long Rest", "en": "Perform a 1-minute rite to regain expended pact slots, up to half your maximum (round up). Once per Long Rest."}},
    {"level": 3, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "เลือก patron (subclass เช่น Fiend Patron) ให้ DM/ผู้เล่นเพิ่มความสามารถของ subclass เองเป็น custom ability", "en": "Choose a patron subclass (e.g. Fiend Patron). The DM or player adds its abilities as custom abilities."}},
    {"level": 4, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 6, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 8, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 9, "name": {"th": "Contact Patron", "en": "Contact Patron"}, "desc": {"th": "มีเวท Contact Other Plane เตรียมไว้เสมอ ร่ายเพื่อคุยกับ patron ได้โดยไม่ใช้ slot และผ่าน save ของเวทนี้อัตโนมัติ ใช้แบบนี้ได้ 1 ครั้ง คืนตอน Long Rest", "en": "Contact Other Plane is always prepared. Once per Long Rest you can cast it without a slot to reach your patron, and you automatically pass its saving throw."}},
    {"level": 10, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 11, "name": {"th": "Mystic Arcanum (ระดับ 6)", "en": "Mystic Arcanum (level 6 spell)"}, "desc": {"th": "เลือกเวท Warlock ระดับ 6 หนึ่งบท ร่ายได้ 1 ครั้งโดยไม่ใช้ slot คืนตอน Long Rest เลเวลขึ้นเปลี่ยนเป็นเวทระดับเดียวกันบทอื่นได้", "en": "Pick one level 6 Warlock spell. Cast it once without a slot; regained on a Long Rest. On level up you may swap it for another spell of that level."}},
    {"level": 12, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 13, "name": {"th": "Mystic Arcanum (ระดับ 7)", "en": "Mystic Arcanum (level 7 spell)"}, "desc": {"th": "ได้ arcanum อีกบท เป็นเวท Warlock ระดับ 7 ร่ายได้ 1 ครั้งโดยไม่ใช้ slot คืนตอน Long Rest", "en": "Gain another arcanum: a level 7 Warlock spell cast once without a slot, regained on a Long Rest."}},
    {"level": 14, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 15, "name": {"th": "Mystic Arcanum (ระดับ 8)", "en": "Mystic Arcanum (level 8 spell)"}, "desc": {"th": "ได้ arcanum อีกบท เป็นเวท Warlock ระดับ 8 ร่ายได้ 1 ครั้งโดยไม่ใช้ slot คืนตอน Long Rest", "en": "Gain another arcanum: a level 8 Warlock spell cast once without a slot, regained on a Long Rest."}},
    {"level": 16, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 17, "name": {"th": "Mystic Arcanum (ระดับ 9)", "en": "Mystic Arcanum (level 9 spell)"}, "desc": {"th": "ได้ arcanum อีกบท เป็นเวท Warlock ระดับ 9 ร่ายได้ 1 ครั้งโดยไม่ใช้ slot คืนตอน Long Rest", "en": "Gain another arcanum: a level 9 Warlock spell cast once without a slot, regained on a Long Rest."}},
    {"level": 19, "name": {"th": "Epic Boon", "en": "Epic Boon"}, "desc": {"th": "ได้ Epic Boon feat 1 อย่าง (หรือ feat อื่นที่เข้าเงื่อนไข) แนะนำ Boon of Fate ให้เพิ่มเองเป็น custom ability และถ้าเพิ่ม ability ให้ DM แก้ค่า (Epic Boon เพิ่มได้ถึง 30)", "en": "Gain an Epic Boon feat (or another feat you qualify for); Boon of Fate is suggested. Add it as a custom ability; if it raises an ability, the DM edits it (Epic Boons allow up to 30)."}},
    {"level": 20, "name": {"th": "Eldritch Master", "en": "Eldritch Master"}, "desc": {"th": "เมื่อใช้ Magical Cunning จะคืน pact slot ที่ใช้ไปได้ครบทุกช่อง", "en": "Magical Cunning now restores all of your expended pact slots."}}
  ]
}
```
