# Paladin (อัศวินศักดิ์สิทธิ์)

นักรบเกราะหนักที่ได้พลังจากคำสาบาน ใช้ Charisma ร่ายเวทแบบ half caster ยืนแนวหน้า ฟันแล้วเผา spell slot เป็น Divine Smite
นอกจากตีแรงยังรักษาเพื่อนด้วย Lay On Hands และแผ่ aura ให้เพื่อนรอบตัว save ง่ายขึ้น · research: [`../../research/classes/phb-2024/paladin.md`](../../research/classes/phb-2024/paladin.md)

> กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) · คำอธิบายเขียนใหม่ด้วยคำของเราเอง

## ข้อมูลหลัก

| | |
|---|---|
| Hit Die | d10 |
| Primary ability | Strength, Charisma |
| Saving throws | Wisdom, Charisma |
| Armor / Weapons / Tools | Light, Medium, Heavy armor + Shields / Simple + Martial weapons / ไม่มี |
| Skills | เลือก 2: Athletics, Insight, Intimidation, Medicine, Persuasion, Religion |
| Spellcasting | half caster · Charisma · ใช้ Holy Symbol เป็น focus |

## Feature ตามเลเวล

| Level | Feature |
|---|---|
| 1 | Lay On Hands, Spellcasting, Weapon Mastery |
| 2 | Fighting Style, Paladin's Smite |
| 3 | Channel Divinity, Subclass feature |
| 4 | Ability Score Improvement |
| 5 | Extra Attack, Faithful Steed |
| 6 | Aura of Protection |
| 7 | Subclass feature |
| 8 | Ability Score Improvement |
| 9 | Abjure Foes |
| 10 | Aura of Courage |
| 11 | Radiant Strikes |
| 12 | Ability Score Improvement |
| 14 | Restoring Touch |
| 15 | Subclass feature |
| 16 | Ability Score Improvement |
| 18 | Aura Expansion |
| 19 | Epic Boon |
| 20 | Subclass feature |

## ตารางตามเลเวล + resources

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Channel Divinity | — | — | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 |
| Prepared spells | 2 | 3 | 4 | 5 | 6 | 6 | 7 | 7 | 9 | 9 | 10 | 10 | 11 | 11 | 12 | 12 | 14 | 14 | 15 | 15 |
| Weapon Mastery | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |

| Resource | จำนวน | คืนเมื่อ |
|---|---|---|
| Lay On Hands (pool HP) | 5 × เลเวล | Long Rest |
| Channel Divinity | ตามตาราง (เลเวล 3+) | Long Rest (Short Rest คืน 1 ครั้ง) |
| Paladin's Smite (ร่ายฟรี) | 1 (เลเวล 2+) | Long Rest |
| Faithful Steed (ร่ายฟรี) | 1 (เลเวล 5+) | Long Rest |

Spell slots: ตาราง half caster ใน [`_shared.md`](_shared.md)

## JSON (machine-readable)

```json
{
  "id": "paladin",
  "name": {"th": "อัศวินศักดิ์สิทธิ์", "en": "Paladin"},
  "hitDie": 10,
  "primary": ["str", "cha"],
  "saves": ["wis", "cha"],
  "armor": {"th": "เกราะ Light, Medium, Heavy และ Shield", "en": "Light, Medium, and Heavy armor and Shields"},
  "weapons": {"th": "อาวุธ Simple และ Martial", "en": "Simple and Martial weapons"},
  "tools": {"th": "", "en": ""},
  "skillChoices": {"count": 2, "from": ["athletics", "insight", "intimidation", "medicine", "persuasion", "religion"]},
  "caster": "half",
  "castAbility": "cha",
  "tables": [
    {"id": "channel-divinity", "name": {"th": "Channel Divinity", "en": "Channel Divinity"}, "values": ["—", "—", "2", "2", "2", "2", "2", "2", "2", "2", "3", "3", "3", "3", "3", "3", "3", "3", "3", "3"]},
    {"id": "prepared", "name": {"th": "เวทที่เตรียมได้", "en": "Prepared Spells"}, "values": ["2", "3", "4", "5", "6", "6", "7", "7", "9", "9", "10", "10", "11", "11", "12", "12", "14", "14", "15", "15"]},
    {"id": "weapon-mastery", "name": {"th": "Weapon Mastery", "en": "Weapon Mastery"}, "values": ["2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2"]}
  ],
  "resources": [
    {"id": "lay-on-hands", "name": {"th": "Lay On Hands (HP ในคลัง)", "en": "Lay On Hands (pool HP)"}, "max": null, "maxFormula": "level*5", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "channel-divinity", "name": {"th": "Channel Divinity", "en": "Channel Divinity"}, "max": [0, 0, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3], "maxFormula": "", "refill": "long", "shortRegain": 1, "shortFromLevel": 0},
    {"id": "paladins-smite", "name": {"th": "Divine Smite ฟรี", "en": "Free Divine Smite"}, "max": [0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0},
    {"id": "faithful-steed", "name": {"th": "Find Steed ฟรี", "en": "Free Find Steed"}, "max": [0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1], "maxFormula": "", "refill": "long", "shortRegain": 0, "shortFromLevel": 0}
  ],
  "features": [
    {"level": 1, "name": {"th": "Lay On Hands", "en": "Lay On Hands"}, "desc": {"th": "มีคลัง HP สำหรับรักษา = 5 × เลเวล Paladin คืนเต็มตอน Long Rest ใช้ Bonus Action แตะตัวเองหรือคนอื่นเพื่อแจก HP จากคลัง หรือจ่าย 5 แต้มเพื่อแก้ Poisoned (แต้มนั้นไม่ได้รักษา HP)", "en": "You hold a healing pool of 5 × your Paladin level HP, refilled on a Long Rest. As a Bonus Action, touch a creature (or yourself) to give it HP from the pool, or spend 5 points to end Poisoned (those points heal nothing)."}},
    {"level": 1, "name": {"th": "Spellcasting", "en": "Spellcasting"}, "desc": {"th": "ร่ายเวทของ Paladin ด้วย Charisma ใช้ Holy Symbol เป็น focus ได้ จำนวนเวทที่เตรียมดูตามตาราง ต้องเป็นเลเวลที่มี slot หลัง Long Rest สลับเวทที่เตรียมได้ 1 บท slot คืนครบตอน Long Rest", "en": "Cast Paladin spells using Charisma; a Holy Symbol works as your focus. Prepared spell count follows the table and must fit your slot levels. After a Long Rest you may swap one prepared spell. Slots return on a Long Rest."}},
    {"level": 1, "name": {"th": "Weapon Mastery", "en": "Weapon Mastery"}, "desc": {"th": "เลือกอาวุธที่ถนัด 2 ชนิด ใช้ mastery property ของอาวุธนั้นได้ เปลี่ยนชนิดได้ทุกครั้งที่จบ Long Rest", "en": "Pick two kinds of weapons you are proficient with and use their mastery properties. You can change the picks after each Long Rest."}},
    {"level": 2, "name": {"th": "Fighting Style", "en": "Fighting Style"}, "desc": {"th": "เลือก Fighting Style feat 1 อย่าง หรือเลือก Blessed Warrior แทน: ได้ cantrip ของ Cleric 2 อัน (นับเป็นเวท Paladin ใช้ Charisma) เลเวลขึ้นสลับได้ 1 อัน ให้เพิ่มเองเป็น custom ability", "en": "Take a Fighting Style feat, or Blessed Warrior instead: learn two Cleric cantrips that count as Paladin spells and use Charisma; swap one on level up. Add your choice as a custom ability."}},
    {"level": 2, "name": {"th": "Paladin's Smite", "en": "Paladin's Smite"}, "desc": {"th": "มีเวท Divine Smite เตรียมไว้เสมอ (ไม่นับโควตา) และร่ายได้ฟรีโดยไม่ใช้ slot 1 ครั้ง คืนตอน Long Rest", "en": "Divine Smite is always prepared (outside your count), and once per Long Rest you can cast it without a spell slot."}},
    {"level": 3, "name": {"th": "Channel Divinity", "en": "Channel Divinity"}, "desc": {"th": "ใช้พลังศักดิ์สิทธิ์ได้ 2 ครั้ง (3 ครั้งที่เลเวล 11) Short Rest คืน 1 ครั้ง Long Rest คืนหมด เริ่มด้วย Divine Sense: Bonus Action รู้ตำแหน่ง Celestial, Fiend, Undead และที่ที่ถูกเสกให้ศักดิ์สิทธิ์/แปดเปื้อนในระยะ 60 ft นาน 10 นาที", "en": "Channel divine power twice (three times from level 11); a Short Rest returns one use, a Long Rest all. First option, Divine Sense: Bonus Action, for 10 minutes sense Celestials, Fiends, Undead and hallowed or desecrated places within 60 ft."}},
    {"level": 3, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "เลือก Oath ของ Paladin (เช่น Oath of Devotion) ให้ DM/ผู้เล่นเพิ่มความสามารถของ subclass เองเป็น custom ability", "en": "Choose a Paladin Oath subclass (e.g. Oath of Devotion). The DM or player adds its abilities as custom abilities."}},
    {"level": 4, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 5, "name": {"th": "Extra Attack", "en": "Extra Attack"}, "desc": {"th": "เมื่อใช้ Attack action ในเทิร์นตัวเอง โจมตีได้ 2 ครั้งแทน 1 ครั้ง", "en": "When you take the Attack action on your turn, you attack twice instead of once."}},
    {"level": 5, "name": {"th": "Faithful Steed", "en": "Faithful Steed"}, "desc": {"th": "มีเวท Find Steed เตรียมไว้เสมอ และร่ายได้ฟรีโดยไม่ใช้ slot 1 ครั้ง คืนตอน Long Rest", "en": "Find Steed is always prepared, and once per Long Rest you can cast it without a spell slot."}},
    {"level": 6, "name": {"th": "Aura of Protection", "en": "Aura of Protection"}, "desc": {"th": "แผ่ aura รัศมี 10 ft รอบตัว ตัวเองและเพื่อนในนั้นบวก Charisma mod (อย่างน้อย +1) ใน saving throw ทุกครั้ง aura ดับขณะติด Incapacitated ถ้ามี Paladin หลายคน แต่ละคนรับได้ aura เดียว", "en": "You emit a 10-ft aura. You and allies inside add your Charisma modifier (min +1) to saving throws. It shuts off while you are Incapacitated. A creature benefits from only one such aura at a time."}},
    {"level": 7, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 8, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 9, "name": {"th": "Abjure Foes", "en": "Abjure Foes"}, "desc": {"th": "Magic action ใช้ Channel Divinity 1 ครั้ง เลือกศัตรูที่เห็นในระยะ 60 ft ได้จำนวน = Charisma mod (อย่างน้อย 1) ต้อง save Wisdom ไม่ผ่านติด Frightened 1 นาทีหรือจนโดนดาเมจ ระหว่างนั้นแต่ละเทิร์นทำได้แค่อย่างเดียว: เดิน, action หรือ Bonus Action", "en": "Magic action, spend one Channel Divinity: pick foes you see within 60 ft, up to your Charisma modifier (min 1). Each makes a Wisdom save or is Frightened for 1 minute or until damaged, and can only move, act, or take a Bonus Action each turn."}},
    {"level": 10, "name": {"th": "Aura of Courage", "en": "Aura of Courage"}, "desc": {"th": "ตัวเองและเพื่อนใน Aura of Protection ไม่ติด Frightened ถ้าเพื่อนที่กลัวอยู่เดินเข้ามา อาการกลัวไม่มีผลขณะอยู่ใน aura", "en": "You and allies in your Aura of Protection are immune to Frightened. An already Frightened ally ignores it while inside the aura."}},
    {"level": 11, "name": {"th": "Radiant Strikes", "en": "Radiant Strikes"}, "desc": {"th": "ทุกครั้งที่โจมตีโดนด้วยอาวุธระยะประชิดหรือ Unarmed Strike เป้าโดน Radiant เพิ่ม 1d8", "en": "Whenever you hit with a Melee weapon or an Unarmed Strike, the target takes an extra 1d8 Radiant damage."}},
    {"level": 12, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 14, "name": {"th": "Restoring Touch", "en": "Restoring Touch"}, "desc": {"th": "ตอนใช้ Lay On Hands จ่ายเพิ่ม 5 แต้มต่ออาการเพื่อปลด Blinded, Charmed, Deafened, Frightened, Paralyzed หรือ Stunned ได้ (แต้มนั้นไม่ได้รักษา HP)", "en": "When you use Lay On Hands, spend 5 points per condition to also end Blinded, Charmed, Deafened, Frightened, Paralyzed, or Stunned (those points heal nothing)."}},
    {"level": 15, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 16, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 18, "name": {"th": "Aura Expansion", "en": "Aura Expansion"}, "desc": {"th": "Aura of Protection ขยายจาก 10 ft เป็น 30 ft", "en": "Your Aura of Protection grows from 10 ft to 30 ft."}},
    {"level": 19, "name": {"th": "Epic Boon", "en": "Epic Boon"}, "desc": {"th": "ได้ Epic Boon feat 1 อย่าง (หรือ feat อื่นที่เข้าเงื่อนไข) แนะนำ Boon of Truesight ให้เพิ่มเองเป็น custom ability และถ้าเพิ่ม ability ให้ DM แก้ค่า (Epic Boon เพิ่มได้ถึง 30)", "en": "Gain an Epic Boon feat (or another feat you qualify for); Boon of Truesight is suggested. Add it as a custom ability; if it raises an ability, the DM edits it (Epic Boons allow up to 30)."}},
    {"level": 20, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถสูงสุดของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain your subclass capstone ability. The DM or player adds it as a custom ability."}}
  ]
}
```
