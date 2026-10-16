# Monk (นักพรตยุทธ์)

Monk คือนักสู้มือเปล่าที่ไม่ใส่เกราะ เร็ว คล่อง และตีได้หลายครั้งต่อเทิร์นด้วย Martial Arts  
ใช้ Focus Points เปิดท่าพิเศษ เช่น ตีรัว หลบ วิ่งไว และทำให้ศัตรู Stunned วิ่งเข้าออกแนวหน้าตลอดเวลา

Research ประกอบ: [`../../research/classes/phb-2024/monk.md`](../../research/classes/phb-2024/monk.md) · กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) คำอธิบายเขียนเองด้วยคำของเรา

## ข้อมูลหลัก
| | |
|---|---|
| Hit Die | d8 |
| Primary ability | DEX, WIS |
| Saving throws | STR, DEX |
| Armor | ไม่มี |
| Weapons | อาวุธ Simple และ Martial ที่มี Light |
| Tools | Artisan's Tools หรือ Musical Instrument 1 ชนิด |
| Skills | เลือก 2 จาก acrobatics, athletics, history, insight, religion, stealth |
| Spellcasting | ไม่มี |

## Feature ตามเลเวล
| Level | Feature | สรุป |
|---|---|---|
| 1 | Martial Arts | ไม่ใส่เกราะ/โล่ และใช้มือเปล่าหรือ Monk weapon (อาวุธ Simple ระยะประชิด, Martial ระยะประชิดที่มี Light): ตีมือเปล่าเพิ่มด้วย Bonus Action, ทอย Martial Arts die แทนดาเมจปกติ, ใช้ Dex แทน Str ตีและ DC Grapple/Shove ได้ |
| 1 | Unarmored Defense | ไม่ใส่เกราะและไม่ถือโล่ AC ฐาน = 10 + Dex mod + Wis mod (DM ใส่ค่า AC ในการ์ด) |
| 2 | Monk's Focus | ได้ Focus Points ตามตาราง คืนหมดเมื่อพัก จ่าย 1 แต้ม: Flurry of Blows (ตีมือเปล่า 2 ครั้ง), Patient Defense (Disengage+Dodge), Step of the Wind (Disengage+Dash, กระโดดไกล 2 เท่า) ทั้งหมดเป็น Bonus Action เซฟ DC = 8 + Wis mod + PB |
| 2 | Unarmored Movement | ไม่ใส่เกราะและไม่ถือโล่ Speed เพิ่มตามตาราง (เริ่ม +10 ft) |
| 2 | Uncanny Metabolism | ตอนทอย Initiative คืน Focus Points ทั้งหมด และฟื้น HP = เลเวล Monk + ทอย Martial Arts die ใช้ได้ 1 ครั้งต่อ Long Rest |
| 3 | Deflect Attacks | Reaction เมื่อโดนโจมตีที่มีดาเมจ Bludgeoning/Piercing/Slashing ลดดาเมจ 1d10 + Dex mod + เลเวล Monk ถ้าลดเหลือ 0 จ่าย 1 Focus ส่งแรงใส่เป้า (เซฟ Dex) ดาเมจ Martial Arts die 2 ลูก + Dex mod |
| 3 | Subclass feature | เลือก subclass ของ Monk (เช่น Warrior of the Open Hand) แล้วเพิ่ม feature ของ subclass เองเป็น custom ability |
| 4 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 4 | Slow Fall | Reaction ตอนตก ลดดาเมจจากการตก 5 × เลเวล Monk |
| 5 | Extra Attack | Attack action โจมตีได้ 2 ครั้ง |
| 5 | Stunning Strike | เทิร์นละครั้ง ตีโดนด้วย Monk weapon หรือมือเปล่า จ่าย 1 Focus ให้เป้าเซฟ Con พลาด = Stunned ถึงต้นเทิร์นหน้าของเรา ผ่าน = Speed ลดครึ่งและการโจมตีครั้งถัดไปใส่มันได้ Advantage |
| 6 | Empowered Strikes | ดาเมจจาก Unarmed Strike เลือกเป็น Force หรือชนิดปกติได้ |
| 6 | Subclass feature | ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability |
| 7 | Evasion | เซฟ Dex ที่ผ่านแล้วโดนครึ่ง: ผ่าน = ไม่โดนเลย พลาด = โดนครึ่ง ใช้ไม่ได้ตอน Incapacitated |
| 8 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 9 | Acrobatic Movement | ไม่ใส่เกราะและไม่ถือโล่ เดินบนผนังแนวตั้งและบนผิวของเหลวได้ระหว่างเทิร์นโดยไม่ตก |
| 10 | Heightened Focus | Flurry of Blows ตี 3 ครั้ง, Patient Defense แบบจ่ายแต้มได้ temp HP = Martial Arts die 2 ลูก, Step of the Wind แบบจ่ายแต้มพาเพื่อนขนาด Large หรือเล็กกว่าที่อยู่ติดตัวไปด้วย |
| 10 | Self-Restoration | จบเทิร์นแต่ละครั้ง ปลด Charmed, Frightened หรือ Poisoned ออกจากตัวเองได้ 1 อย่าง และอดอาหาร/น้ำไม่ทำให้ได้ Exhaustion |
| 11 | Subclass feature | ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability |
| 12 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 13 | Deflect Energy | Deflect Attacks ใช้กับการโจมตีดาเมจทุกชนิดได้ |
| 14 | Disciplined Survivor | proficient ทุก saving throw (DM ติ๊กเพิ่มเอง) และเซฟพลาดจ่าย 1 Focus ทอยใหม่ได้ ต้องใช้ผลใหม่ |
| 15 | Perfect Focus | ทอย Initiative โดยไม่ใช้ Uncanny Metabolism ถ้า Focus เหลือ 3 หรือน้อยกว่า เติมกลับเป็น 4 |
| 16 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 17 | Subclass feature | ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability |
| 18 | Superior Defense | ต้นเทิร์น จ่าย 3 Focus ได้ Resistance ทุกดาเมจยกเว้น Force นาน 1 นาที หรือจนกว่าจะ Incapacitated |
| 19 | Epic Boon | ได้ Epic Boon feat 1 อัน หรือ feat อื่นที่เข้าเงื่อนไข (แนะนำ Boon of Irresistible Offense) ให้เพิ่มเองเป็น custom ability |
| 20 | Body and Mind | Dex และ Wis +4 (สูงสุด 25) — DM แก้ ability ให้ |

## ค่าตามเลเวล (แสดงอย่างเดียว)
| ตาราง | เลเวล: ค่า |
|---|---|
| Martial Arts die | 1–4: 1d6 · 5–10: 1d8 · 11–16: 1d10 · 17–20: 1d12 |
| Unarmored Movement | 1: — · 2–5: +10 ft · 6–9: +15 ft · 10–13: +20 ft · 14–17: +25 ft · 18–20: +30 ft |

## Resources (counter ใช้/คืนได้)
| Resource | จำนวนตามเลเวล | คืนเมื่อ |
|---|---|---|
| Focus Points | 2: 2 · 3: 3 · 4: 4 · 5: 5 · 6: 6 · 7: 7 · 8: 8 · 9: 9 · 10: 10 · 11: 11 · 12: 12 · 13: 13 · 14: 14 · 15: 15 · 16: 16 · 17: 17 · 18: 18 · 19: 19 · 20: 20 | Short Rest (คืนหมด) |
| Uncanny Metabolism | 2–20: 1 | Long Rest (คืนหมด) |

## JSON (machine-readable)

```json
{
  "id": "monk",
  "name": {
    "th": "นักพรตยุทธ์",
    "en": "Monk"
  },
  "hitDie": 8,
  "primary": [
    "dex",
    "wis"
  ],
  "saves": [
    "str",
    "dex"
  ],
  "armor": {
    "th": "ไม่มี",
    "en": "None"
  },
  "weapons": {
    "th": "อาวุธ Simple และ Martial ที่มี Light",
    "en": "Simple weapons, Martial weapons with Light"
  },
  "tools": {
    "th": "Artisan's Tools หรือ Musical Instrument 1 ชนิด",
    "en": "One Artisan's Tools or Musical Instrument"
  },
  "skillChoices": {
    "count": 2,
    "from": [
      "acrobatics",
      "athletics",
      "history",
      "insight",
      "religion",
      "stealth"
    ]
  },
  "caster": "none",
  "castAbility": "",
  "tables": [
    {
      "id": "martial-arts",
      "name": {
        "th": "Martial Arts die",
        "en": "Martial Arts die"
      },
      "values": [
        "1d6",
        "1d6",
        "1d6",
        "1d6",
        "1d8",
        "1d8",
        "1d8",
        "1d8",
        "1d8",
        "1d8",
        "1d10",
        "1d10",
        "1d10",
        "1d10",
        "1d10",
        "1d10",
        "1d12",
        "1d12",
        "1d12",
        "1d12"
      ]
    },
    {
      "id": "unarmored-movement",
      "name": {
        "th": "Unarmored Movement",
        "en": "Unarmored Movement"
      },
      "values": [
        "—",
        "+10 ft",
        "+10 ft",
        "+10 ft",
        "+10 ft",
        "+15 ft",
        "+15 ft",
        "+15 ft",
        "+15 ft",
        "+20 ft",
        "+20 ft",
        "+20 ft",
        "+20 ft",
        "+25 ft",
        "+25 ft",
        "+25 ft",
        "+25 ft",
        "+30 ft",
        "+30 ft",
        "+30 ft"
      ]
    }
  ],
  "resources": [
    {
      "id": "focus-points",
      "name": {
        "th": "Focus Points",
        "en": "Focus Points"
      },
      "max": [
        0,
        2,
        3,
        4,
        5,
        6,
        7,
        8,
        9,
        10,
        11,
        12,
        13,
        14,
        15,
        16,
        17,
        18,
        19,
        20
      ],
      "maxFormula": "",
      "refill": "short",
      "shortRegain": 0,
      "shortFromLevel": 0
    },
    {
      "id": "uncanny-metabolism",
      "name": {
        "th": "Uncanny Metabolism",
        "en": "Uncanny Metabolism"
      },
      "max": [
        0,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1
      ],
      "maxFormula": "",
      "refill": "long",
      "shortRegain": 0,
      "shortFromLevel": 0
    }
  ],
  "features": [
    {
      "level": 1,
      "name": {
        "th": "Martial Arts",
        "en": "Martial Arts"
      },
      "desc": {
        "th": "ไม่ใส่เกราะ/โล่ และใช้มือเปล่าหรือ Monk weapon (อาวุธ Simple ระยะประชิด, Martial ระยะประชิดที่มี Light): ตีมือเปล่าเพิ่มด้วย Bonus Action, ทอย Martial Arts die แทนดาเมจปกติ, ใช้ Dex แทน Str ตีและ DC Grapple/Shove ได้",
        "en": "Unarmored, no Shield, using Unarmed Strikes or Monk weapons (Simple melee, Light Martial melee): Bonus Action Unarmed Strike, roll the Martial Arts die for damage, use Dex for attacks, damage and Grapple/Shove DC."
      }
    },
    {
      "level": 1,
      "name": {
        "th": "Unarmored Defense",
        "en": "Unarmored Defense"
      },
      "desc": {
        "th": "ไม่ใส่เกราะและไม่ถือโล่ AC ฐาน = 10 + Dex mod + Wis mod (DM ใส่ค่า AC ในการ์ด)",
        "en": "Without armor or a Shield, your base AC is 10 + Dex mod + Wis mod (the DM sets AC on the card)."
      }
    },
    {
      "level": 2,
      "name": {
        "th": "Monk's Focus",
        "en": "Monk's Focus"
      },
      "desc": {
        "th": "ได้ Focus Points ตามตาราง คืนหมดเมื่อพัก จ่าย 1 แต้ม: Flurry of Blows (ตีมือเปล่า 2 ครั้ง), Patient Defense (Disengage+Dodge), Step of the Wind (Disengage+Dash, กระโดดไกล 2 เท่า) ทั้งหมดเป็น Bonus Action เซฟ DC = 8 + Wis mod + PB",
        "en": "Gain Focus Points (per the table), all back on any rest. 1 point, as a Bonus Action: Flurry of Blows (two Unarmed Strikes), Patient Defense (Disengage + Dodge) or Step of the Wind (Disengage + Dash, double jump). Save DC 8 + Wis mod + PB."
      }
    },
    {
      "level": 2,
      "name": {
        "th": "Unarmored Movement",
        "en": "Unarmored Movement"
      },
      "desc": {
        "th": "ไม่ใส่เกราะและไม่ถือโล่ Speed เพิ่มตามตาราง (เริ่ม +10 ft)",
        "en": "Without armor or a Shield, your Speed rises per the table (starts at +10 ft)."
      }
    },
    {
      "level": 2,
      "name": {
        "th": "Uncanny Metabolism",
        "en": "Uncanny Metabolism"
      },
      "desc": {
        "th": "ตอนทอย Initiative คืน Focus Points ทั้งหมด และฟื้น HP = เลเวล Monk + ทอย Martial Arts die ใช้ได้ 1 ครั้งต่อ Long Rest",
        "en": "When you roll Initiative, regain all Focus Points and heal Monk level + a Martial Arts die roll. Once per Long Rest."
      }
    },
    {
      "level": 3,
      "name": {
        "th": "Deflect Attacks",
        "en": "Deflect Attacks"
      },
      "desc": {
        "th": "Reaction เมื่อโดนโจมตีที่มีดาเมจ Bludgeoning/Piercing/Slashing ลดดาเมจ 1d10 + Dex mod + เลเวล Monk ถ้าลดเหลือ 0 จ่าย 1 Focus ส่งแรงใส่เป้า (เซฟ Dex) ดาเมจ Martial Arts die 2 ลูก + Dex mod",
        "en": "Reaction when hit by Bludgeoning, Piercing or Slashing damage: reduce it by 1d10 + Dex mod + Monk level. At 0, spend 1 Focus to redirect it (Dex save) for two Martial Arts dice + Dex mod."
      }
    },
    {
      "level": 3,
      "name": {
        "th": "ฟีเจอร์ subclass",
        "en": "Subclass feature"
      },
      "desc": {
        "th": "เลือก subclass ของ Monk (เช่น Warrior of the Open Hand) แล้วเพิ่ม feature ของ subclass เองเป็น custom ability",
        "en": "Choose a Monk subclass (e.g. Warrior of the Open Hand) and add its features yourself as custom abilities."
      }
    },
    {
      "level": 4,
      "name": {
        "th": "Ability Score Improvement",
        "en": "Ability Score Improvement"
      },
      "desc": {
        "th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์",
        "en": "+2 to one ability or +1 to two (max 20), or take a feat instead. The DM edits the abilities on the profile card."
      }
    },
    {
      "level": 4,
      "name": {
        "th": "Slow Fall",
        "en": "Slow Fall"
      },
      "desc": {
        "th": "Reaction ตอนตก ลดดาเมจจากการตก 5 × เลเวล Monk",
        "en": "Reaction when you fall: reduce falling damage by 5 × Monk level."
      }
    },
    {
      "level": 5,
      "name": {
        "th": "Extra Attack",
        "en": "Extra Attack"
      },
      "desc": {
        "th": "Attack action โจมตีได้ 2 ครั้ง",
        "en": "Attack twice when you take the Attack action."
      }
    },
    {
      "level": 5,
      "name": {
        "th": "Stunning Strike",
        "en": "Stunning Strike"
      },
      "desc": {
        "th": "เทิร์นละครั้ง ตีโดนด้วย Monk weapon หรือมือเปล่า จ่าย 1 Focus ให้เป้าเซฟ Con พลาด = Stunned ถึงต้นเทิร์นหน้าของเรา ผ่าน = Speed ลดครึ่งและการโจมตีครั้งถัดไปใส่มันได้ Advantage",
        "en": "Once per turn on a Monk weapon or Unarmed Strike hit, spend 1 Focus: Con save or Stunned until your next turn starts. On a save, its Speed halves and the next attack on it has Advantage."
      }
    },
    {
      "level": 6,
      "name": {
        "th": "Empowered Strikes",
        "en": "Empowered Strikes"
      },
      "desc": {
        "th": "ดาเมจจาก Unarmed Strike เลือกเป็น Force หรือชนิดปกติได้",
        "en": "Your Unarmed Strikes can deal Force damage or their normal type."
      }
    },
    {
      "level": 6,
      "name": {
        "th": "ฟีเจอร์ subclass",
        "en": "Subclass feature"
      },
      "desc": {
        "th": "ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability",
        "en": "Gain your subclass feature for this level; add it yourself as a custom ability."
      }
    },
    {
      "level": 7,
      "name": {
        "th": "Evasion",
        "en": "Evasion"
      },
      "desc": {
        "th": "เซฟ Dex ที่ผ่านแล้วโดนครึ่ง: ผ่าน = ไม่โดนเลย พลาด = โดนครึ่ง ใช้ไม่ได้ตอน Incapacitated",
        "en": "On a Dex save for half damage, take none on a success and half on a failure. Not while Incapacitated."
      }
    },
    {
      "level": 8,
      "name": {
        "th": "Ability Score Improvement",
        "en": "Ability Score Improvement"
      },
      "desc": {
        "th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์",
        "en": "+2 to one ability or +1 to two (max 20), or take a feat instead. The DM edits the abilities on the profile card."
      }
    },
    {
      "level": 9,
      "name": {
        "th": "Acrobatic Movement",
        "en": "Acrobatic Movement"
      },
      "desc": {
        "th": "ไม่ใส่เกราะและไม่ถือโล่ เดินบนผนังแนวตั้งและบนผิวของเหลวได้ระหว่างเทิร์นโดยไม่ตก",
        "en": "Without armor or a Shield, move along vertical surfaces and across liquids on your turn without falling."
      }
    },
    {
      "level": 10,
      "name": {
        "th": "Heightened Focus",
        "en": "Heightened Focus"
      },
      "desc": {
        "th": "Flurry of Blows ตี 3 ครั้ง, Patient Defense แบบจ่ายแต้มได้ temp HP = Martial Arts die 2 ลูก, Step of the Wind แบบจ่ายแต้มพาเพื่อนขนาด Large หรือเล็กกว่าที่อยู่ติดตัวไปด้วย",
        "en": "Flurry of Blows makes three strikes; paid Patient Defense grants temp HP of two Martial Arts dice; paid Step of the Wind carries an adjacent willing Large or smaller ally."
      }
    },
    {
      "level": 10,
      "name": {
        "th": "Self-Restoration",
        "en": "Self-Restoration"
      },
      "desc": {
        "th": "จบเทิร์นแต่ละครั้ง ปลด Charmed, Frightened หรือ Poisoned ออกจากตัวเองได้ 1 อย่าง และอดอาหาร/น้ำไม่ทำให้ได้ Exhaustion",
        "en": "At the end of each turn, end one of Charmed, Frightened or Poisoned on yourself. Going without food or drink gives no Exhaustion."
      }
    },
    {
      "level": 11,
      "name": {
        "th": "ฟีเจอร์ subclass",
        "en": "Subclass feature"
      },
      "desc": {
        "th": "ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability",
        "en": "Gain your subclass feature for this level; add it yourself as a custom ability."
      }
    },
    {
      "level": 12,
      "name": {
        "th": "Ability Score Improvement",
        "en": "Ability Score Improvement"
      },
      "desc": {
        "th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์",
        "en": "+2 to one ability or +1 to two (max 20), or take a feat instead. The DM edits the abilities on the profile card."
      }
    },
    {
      "level": 13,
      "name": {
        "th": "Deflect Energy",
        "en": "Deflect Energy"
      },
      "desc": {
        "th": "Deflect Attacks ใช้กับการโจมตีดาเมจทุกชนิดได้",
        "en": "Deflect Attacks now works against any damage type."
      }
    },
    {
      "level": 14,
      "name": {
        "th": "Disciplined Survivor",
        "en": "Disciplined Survivor"
      },
      "desc": {
        "th": "proficient ทุก saving throw (DM ติ๊กเพิ่มเอง) และเซฟพลาดจ่าย 1 Focus ทอยใหม่ได้ ต้องใช้ผลใหม่",
        "en": "Gain proficiency in all saving throws (DM marks them). After a failed save, spend 1 Focus to reroll; use the new roll."
      }
    },
    {
      "level": 15,
      "name": {
        "th": "Perfect Focus",
        "en": "Perfect Focus"
      },
      "desc": {
        "th": "ทอย Initiative โดยไม่ใช้ Uncanny Metabolism ถ้า Focus เหลือ 3 หรือน้อยกว่า เติมกลับเป็น 4",
        "en": "When you roll Initiative without Uncanny Metabolism and have 3 or fewer Focus Points, refill to 4."
      }
    },
    {
      "level": 16,
      "name": {
        "th": "Ability Score Improvement",
        "en": "Ability Score Improvement"
      },
      "desc": {
        "th": "เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์",
        "en": "+2 to one ability or +1 to two (max 20), or take a feat instead. The DM edits the abilities on the profile card."
      }
    },
    {
      "level": 17,
      "name": {
        "th": "ฟีเจอร์ subclass",
        "en": "Subclass feature"
      },
      "desc": {
        "th": "ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability",
        "en": "Gain your subclass feature for this level; add it yourself as a custom ability."
      }
    },
    {
      "level": 18,
      "name": {
        "th": "Superior Defense",
        "en": "Superior Defense"
      },
      "desc": {
        "th": "ต้นเทิร์น จ่าย 3 Focus ได้ Resistance ทุกดาเมจยกเว้น Force นาน 1 นาที หรือจนกว่าจะ Incapacitated",
        "en": "At the start of your turn, spend 3 Focus for Resistance to all damage except Force for 1 minute or until Incapacitated."
      }
    },
    {
      "level": 19,
      "name": {
        "th": "Epic Boon",
        "en": "Epic Boon"
      },
      "desc": {
        "th": "ได้ Epic Boon feat 1 อัน หรือ feat อื่นที่เข้าเงื่อนไข (แนะนำ Boon of Irresistible Offense) ให้เพิ่มเองเป็น custom ability",
        "en": "Gain an Epic Boon feat or another feat you qualify for (Boon of Irresistible Offense suggested). Add it as a custom ability."
      }
    },
    {
      "level": 20,
      "name": {
        "th": "Body and Mind",
        "en": "Body and Mind"
      },
      "desc": {
        "th": "Dex และ Wis +4 (สูงสุด 25) — DM แก้ ability ให้",
        "en": "Dex and Wis rise by 4 (max 25). The DM edits the abilities."
      }
    }
  ]
}
```
