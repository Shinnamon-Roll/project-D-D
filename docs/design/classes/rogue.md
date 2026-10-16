# Rogue (นักย่องเบา)

ผู้เชี่ยวชาญทักษะและการลอบโจมตี ไม่มีเวทใน base class แต่มี skill และ Expertise มากที่สุด
ในการต่อสู้หาจังหวะ Sneak Attack 1 ครั้งต่อเทิร์นให้ดาเมจก้อนใหญ่ แล้วใช้ Cunning Action ซ่อนหรือถอยออก · research: [`../../research/classes/phb-2024/rogue.md`](../../research/classes/phb-2024/rogue.md)

> กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) · คำอธิบายเขียนใหม่ด้วยคำของเราเอง

## ข้อมูลหลัก

| | |
|---|---|
| Hit Die | d8 |
| Primary ability | Dexterity |
| Saving throws | Dexterity, Intelligence (เลเวล 15 ได้ Wisdom, Charisma เพิ่ม) |
| Armor / Weapons / Tools | Light armor / Simple + Martial ที่มี Finesse หรือ Light / Thieves' Tools |
| Skills | เลือก 4: Acrobatics, Athletics, Deception, Insight, Intimidation, Investigation, Perception, Persuasion, Sleight of Hand, Stealth |
| Spellcasting | ไม่มี |

## Feature ตามเลเวล

| Level | Feature |
|---|---|
| 1 | Expertise, Sneak Attack, Thieves' Cant, Weapon Mastery |
| 2 | Cunning Action |
| 3 | Subclass feature, Steady Aim |
| 4 | Ability Score Improvement |
| 5 | Cunning Strike, Uncanny Dodge |
| 6 | Expertise |
| 7 | Evasion, Reliable Talent |
| 8 | Ability Score Improvement |
| 9 | Subclass feature |
| 10 | Ability Score Improvement |
| 11 | Improved Cunning Strike |
| 12 | Ability Score Improvement |
| 13 | Subclass feature |
| 14 | Devious Strikes |
| 15 | Slippery Mind |
| 16 | Ability Score Improvement |
| 17 | Subclass feature |
| 18 | Elusive |
| 19 | Epic Boon |
| 20 | Stroke of Luck |

## ตารางตามเลเวล + resources

| Level | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Sneak Attack | 1d6 | 1d6 | 2d6 | 2d6 | 3d6 | 3d6 | 4d6 | 4d6 | 5d6 | 5d6 | 6d6 | 6d6 | 7d6 | 7d6 | 8d6 | 8d6 | 9d6 | 9d6 | 10d6 | 10d6 |
| Weapon Mastery | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |

| Resource | จำนวน | คืนเมื่อ |
|---|---|---|
| Stroke of Luck | 1 (เลเวล 20) | Short Rest หรือ Long Rest |

## JSON (machine-readable)

```json
{
  "id": "rogue",
  "name": {"th": "นักย่องเบา", "en": "Rogue"},
  "hitDie": 8,
  "primary": ["dex"],
  "saves": ["dex", "int"],
  "armor": {"th": "เกราะ Light", "en": "Light armor"},
  "weapons": {"th": "อาวุธ Simple และ Martial ที่มี Finesse หรือ Light", "en": "Simple weapons and Martial weapons with Finesse or Light"},
  "tools": {"th": "Thieves' Tools", "en": "Thieves' Tools"},
  "skillChoices": {"count": 4, "from": ["acrobatics", "athletics", "deception", "insight", "intimidation", "investigation", "perception", "persuasion", "sleight-of-hand", "stealth"]},
  "caster": "none",
  "castAbility": "",
  "tables": [
    {"id": "sneak-attack", "name": {"th": "Sneak Attack", "en": "Sneak Attack"}, "values": ["1d6", "1d6", "2d6", "2d6", "3d6", "3d6", "4d6", "4d6", "5d6", "5d6", "6d6", "6d6", "7d6", "7d6", "8d6", "8d6", "9d6", "9d6", "10d6", "10d6"]},
    {"id": "weapon-mastery", "name": {"th": "Weapon Mastery", "en": "Weapon Mastery"}, "values": ["2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2", "2"]}
  ],
  "resources": [
    {"id": "stroke-of-luck", "name": {"th": "Stroke of Luck", "en": "Stroke of Luck"}, "max": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1], "maxFormula": "", "refill": "short", "shortRegain": 0, "shortFromLevel": 0}
  ],
  "features": [
    {"level": 1, "name": {"th": "Expertise", "en": "Expertise"}, "desc": {"th": "ได้ Expertise ใน skill ที่ถนัด 2 อย่าง (แนะนำ Sleight of Hand กับ Stealth ถ้ามี) ให้ติ๊กในชีทเอง", "en": "Gain Expertise in two of your proficient skills (Sleight of Hand and Stealth suggested if you have them). Mark them on the sheet."}},
    {"level": 1, "name": {"th": "Sneak Attack", "en": "Sneak Attack"}, "desc": {"th": "1 ครั้งต่อเทิร์น เมื่อโจมตีโดนด้วยอาวุธ Finesse หรือ Ranged และมี Advantage ใส่ดาเมจเพิ่มตามตาราง (1d6 ถึง 10d6) ไม่ต้องมี Advantage ถ้ามีเพื่อนที่ไม่ Incapacitated ยืนห่างเป้าไม่เกิน 5 ft และเราไม่ติด Disadvantage", "en": "Once per turn, when you hit with a Finesse or Ranged weapon and have Advantage, add the table's extra damage (1d6 up to 10d6). You skip the Advantage need if an ally that is not Incapacitated is within 5 ft of the target and you lack Disadvantage."}},
    {"level": 1, "name": {"th": "Thieves' Cant", "en": "Thieves' Cant"}, "desc": {"th": "รู้ภาษาลับ Thieves' Cant และภาษาอื่นอีก 1 ภาษาที่เลือกเอง", "en": "You know Thieves' Cant plus one other language of your choice."}},
    {"level": 1, "name": {"th": "Weapon Mastery", "en": "Weapon Mastery"}, "desc": {"th": "เลือกอาวุธที่ถนัด 2 ชนิด ใช้ mastery property ของอาวุธนั้นได้ เปลี่ยนชนิดได้ทุกครั้งที่จบ Long Rest", "en": "Pick two kinds of weapons you are proficient with and use their mastery properties. You can change the picks after each Long Rest."}},
    {"level": 2, "name": {"th": "Cunning Action", "en": "Cunning Action"}, "desc": {"th": "ในเทิร์นตัวเอง ใช้ Dash, Disengage หรือ Hide เป็น Bonus Action ได้", "en": "On your turn, you can take Dash, Disengage, or Hide as a Bonus Action."}},
    {"level": 3, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "เลือก subclass ของ Rogue (เช่น Thief) ให้ DM/ผู้เล่นเพิ่มความสามารถของ subclass เองเป็น custom ability", "en": "Choose a Rogue subclass (e.g. Thief). The DM or player adds its abilities as custom abilities."}},
    {"level": 3, "name": {"th": "Steady Aim", "en": "Steady Aim"}, "desc": {"th": "Bonus Action ให้การโจมตีครั้งถัดไปในเทิร์นนี้ได้ Advantage ใช้ได้เฉพาะถ้ายังไม่ขยับในเทิร์นนี้ และหลังใช้ Speed เป็น 0 จนจบเทิร์น", "en": "Bonus Action: gain Advantage on your next attack this turn. Only if you have not moved this turn, and your Speed becomes 0 until the turn ends."}},
    {"level": 4, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 5, "name": {"th": "Cunning Strike", "en": "Cunning Strike"}, "desc": {"th": "ตอนลง Sneak Attack ตัดลูกเต๋าออกเพื่อใส่ผลพิเศษ 1 อย่าง: Poison (1d6, save CON ไม่ผ่านติด Poisoned ต้องมี Poisoner's Kit), Trip (1d6, เป้า Large ลงไป save DEX ไม่ผ่านล้ม Prone), Withdraw (1d6, เดินครึ่ง Speed ไม่โดน Opportunity Attack) DC = 8 + DEX mod + PB", "en": "When you deal Sneak Attack, drop dice to add one effect: Poison (1d6, CON save or Poisoned; needs a Poisoner's Kit), Trip (1d6, Large or smaller target, DEX save or Prone), Withdraw (1d6, move half Speed without Opportunity Attacks). DC = 8 + DEX mod + PB."}},
    {"level": 5, "name": {"th": "Uncanny Dodge", "en": "Uncanny Dodge"}, "desc": {"th": "เมื่อโดนโจมตีจากศัตรูที่มองเห็น ใช้ Reaction ลดดาเมจครั้งนั้นเหลือครึ่ง (ปัดลง)", "en": "When an attacker you can see hits you, use your Reaction to halve that attack's damage (round down)."}},
    {"level": 6, "name": {"th": "Expertise", "en": "Expertise"}, "desc": {"th": "ได้ Expertise เพิ่มอีก 2 skill ที่ถนัด ให้ติ๊กในชีทเอง", "en": "Gain Expertise in two more of your proficient skills. Mark them on the sheet."}},
    {"level": 7, "name": {"th": "Evasion", "en": "Evasion"}, "desc": {"th": "effect ที่ให้ save DEX เพื่อรับดาเมจครึ่ง: ผ่าน = ไม่โดนเลย ไม่ผ่าน = โดนครึ่ง ใช้ไม่ได้ขณะติด Incapacitated", "en": "Against effects where a DEX save halves damage, you take none on a success and half on a failure. Not while Incapacitated."}},
    {"level": 7, "name": {"th": "Reliable Talent", "en": "Reliable Talent"}, "desc": {"th": "ทอย ability check ที่ใช้ skill หรือ tool ที่ถนัด ถ้า d20 ออก 9 หรือต่ำกว่า นับเป็น 10", "en": "On an ability check using a skill or tool you are proficient with, treat a d20 roll of 9 or lower as a 10."}},
    {"level": 8, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 9, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 10, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 11, "name": {"th": "Improved Cunning Strike", "en": "Improved Cunning Strike"}, "desc": {"th": "ใส่ Cunning Strike ได้ถึง 2 ผลต่อ Sneak Attack ครั้งเดียว โดยจ่ายลูกเต๋าของทั้งสองผล", "en": "You can add up to two Cunning Strike effects to one Sneak Attack, paying the dice cost of each."}},
    {"level": 12, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 13, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 14, "name": {"th": "Devious Strikes", "en": "Devious Strikes"}, "desc": {"th": "Cunning Strike เพิ่ม 3 แบบ: Daze (2d6, save CON ไม่ผ่านเทิร์นถัดไปทำได้อย่างเดียว เดิน/action/Bonus Action), Knock Out (6d6, save CON ไม่ผ่าน Unconscious 1 นาทีหรือจนโดนดาเมจ save ซ้ำทุกจบเทิร์น), Obscure (3d6, save DEX ไม่ผ่าน Blinded จนจบเทิร์นถัดไป)", "en": "New Cunning Strike options: Daze (2d6, CON save or next turn only move, action, or Bonus Action), Knock Out (6d6, CON save or Unconscious 1 minute or until damaged; repeat save each turn end), Obscure (3d6, DEX save or Blinded until end of its next turn)."}},
    {"level": 15, "name": {"th": "Slippery Mind", "en": "Slippery Mind"}, "desc": {"th": "ได้ proficiency ใน Wisdom และ Charisma saving throw เพิ่ม (รวมเป็น 4 save) ให้ DM ติ๊กในชีทเอง", "en": "Gain proficiency in Wisdom and Charisma saving throws (four saves total). The DM marks them on the sheet."}},
    {"level": 16, "name": {"th": "Ability Score Improvement", "en": "Ability Score Improvement"}, "desc": {"th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 (สูงสุด 20) หรือเลือก feat แทน ให้ DM แก้ค่า ability ในการ์ด", "en": "Raise one ability by 2 or two abilities by 1 (max 20), or take a feat instead. The DM edits the abilities on the card."}},
    {"level": 17, "name": {"th": "ฟีเจอร์ subclass", "en": "Subclass feature"}, "desc": {"th": "ได้ความสามารถใหม่ของ subclass ให้ DM/ผู้เล่นเพิ่มเองเป็น custom ability", "en": "Gain a new subclass ability. The DM or player adds it as a custom ability."}},
    {"level": 18, "name": {"th": "Elusive", "en": "Elusive"}, "desc": {"th": "ไม่มีการโจมตีใดได้ Advantage ใส่เรา เว้นแต่เราติด Incapacitated", "en": "No attack roll against you can have Advantage unless you are Incapacitated."}},
    {"level": 19, "name": {"th": "Epic Boon", "en": "Epic Boon"}, "desc": {"th": "ได้ Epic Boon feat 1 อย่าง (หรือ feat อื่นที่เข้าเงื่อนไข) แนะนำ Boon of the Night Spirit ให้เพิ่มเองเป็น custom ability และถ้าเพิ่ม ability ให้ DM แก้ค่า (Epic Boon เพิ่มได้ถึง 30)", "en": "Gain an Epic Boon feat (or another feat you qualify for); Boon of the Night Spirit is suggested. Add it as a custom ability; if it raises an ability, the DM edits it (Epic Boons allow up to 30)."}},
    {"level": 20, "name": {"th": "Stroke of Luck", "en": "Stroke of Luck"}, "desc": {"th": "ถ้าทอย D20 Test พลาด เปลี่ยนผลเป็น 20 ได้ ใช้ได้ 1 ครั้ง คืนตอนจบ Short Rest หรือ Long Rest", "en": "When you fail a D20 Test, you can turn the roll into a 20. Once per Short or Long Rest."}}
  ]
}
```
