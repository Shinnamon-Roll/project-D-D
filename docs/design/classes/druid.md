# Druid (ดรูอิด)

Druid คือผู้ร่ายเวทสายธรรมชาติแบบ full caster ใช้ Wisdom ทั้งรักษา ควบคุมพื้นที่ และเรียกพลังธาตุ  
จุดเด่นคือ Wild Shape แปลงร่างเป็นสัตว์ไปสอดแนม เดินทาง หรือเข้าแนวหน้า สลับกับร่างคนที่ร่ายเวท

Research ประกอบ: [`../../research/classes/phb-2024/druid.md`](../../research/classes/phb-2024/druid.md) · กลไกและชื่อ feature อ้างอิง SRD 5.2.1 (CC BY 4.0) คำอธิบายเขียนเองด้วยคำของเรา

## ข้อมูลหลัก
| | |
|---|---|
| Hit Die | d8 |
| Primary ability | WIS |
| Saving throws | INT, WIS |
| Armor | เกราะเบา, โล่ |
| Weapons | อาวุธ Simple |
| Tools | Herbalism Kit |
| Skills | เลือก 2 จาก animal-handling, arcana, insight, medicine, nature, perception, religion, survival |
| Spellcasting | full caster · WIS |

## Feature ตามเลเวล
| Level | Feature | สรุป |
|---|---|---|
| 1 | Spellcasting | ร่ายเวทสาย Druid ด้วย Wisdom รู้ cantrip และเตรียมเวทได้ตามตาราง เปลี่ยนรายการเวทที่เตรียมได้หลัง Long Rest และใช้ Druidic Focus เป็น focus ได้ (spell slot ตามตาราง full caster) |
| 1 | Druidic | รู้ภาษาลับของดรูอิด ใช้ทิ้งข้อความซ่อนที่ดรูอิดด้วยกันเห็นทันที และมีเวท Speak with Animals เตรียมไว้เสมอ (ไม่นับโควตา) |
| 1 | Primal Order | เลือก 1 บทบาท: Magician (ได้ cantrip เพิ่ม 1 และบวก Wis mod ขั้นต่ำ +1 ให้เช็ก Arcana/Nature) หรือ Warden (ใช้อาวุธ Martial และเกราะกลางได้) |
| 2 | Wild Shape | Bonus Action แปลงเป็นร่างสัตว์ที่รู้ ได้ temp HP เท่าเลเวล Druid อยู่ได้ครึ่งเลเวลชั่วโมง ร่ายเวทไม่ได้ระหว่างแปลง (สมาธิเดิมไม่หลุด) ออกจากร่างด้วย Bonus Action จำนวนร่างและ CR ดูตาราง |
| 2 | Wild Companion | Magic action ใช้ spell slot หรือ Wild Shape 1 ครั้ง ร่าย Find Familiar โดยไม่ต้องใช้วัสดุ สัตว์รับใช้เป็น Fey และหายไปเมื่อ Long Rest |
| 3 | Subclass feature | เลือก subclass ของ Druid (เช่น Circle of the Land) แล้วเพิ่ม feature ของ subclass เองเป็น custom ability |
| 4 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 5 | Wild Resurgence | เทิร์นละครั้ง ถ้า Wild Shape หมด ใช้ spell slot แลกคืน 1 ครั้งได้ และแลก Wild Shape 1 ครั้งเป็น spell slot ระดับ 1 ได้ 1 ครั้งต่อ Long Rest |
| 6 | Subclass feature | ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability |
| 7 | Elemental Fury | เลือก 1: Potent Spellcasting (บวก Wis mod ให้ดาเมจ cantrip ของ Druid) หรือ Primal Strike (เทิร์นละครั้ง ตีโดนด้วยอาวุธหรือร่างสัตว์ เพิ่มดาเมจ 1d8 Cold, Fire, Lightning หรือ Thunder) |
| 8 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 10 | Subclass feature | ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability |
| 12 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 14 | Subclass feature | ได้ feature ของ subclass ที่เลือกในเลเวลนี้ ให้เพิ่มเองเป็น custom ability |
| 15 | Improved Elemental Fury | ตัวเลือกจาก Elemental Fury แรงขึ้น: Potent Spellcasting ทำให้ cantrip ที่ระยะ 10 ft ขึ้นไปไกลขึ้นอีก 300 ft, Primal Strike เพิ่มเป็น 2d8 |
| 16 | Ability Score Improvement | เพิ่ม ability 1 ค่า +2 หรือ 2 ค่า +1 หรือเลือก feat แทน (ability สูงสุด 20) — DM เป็นคนแก้ ability ให้ในการ์ดโปรไฟล์ |
| 18 | Beast Spells | ร่ายเวทได้ระหว่างอยู่ในร่าง Wild Shape ยกเว้นเวทที่ต้องใช้วัสดุที่มีราคาหรือวัสดุที่ถูกใช้หมดไป |
| 19 | Epic Boon | ได้ Epic Boon feat 1 อัน หรือ feat อื่นที่เข้าเงื่อนไข (แนะนำ Boon of Dimensional Travel) ให้เพิ่มเองเป็น custom ability |
| 20 | Archdruid | ทอย Initiative ตอน Wild Shape หมดจะคืน 1 ครั้ง, 1 ครั้งต่อ Long Rest แลก Wild Shape หลายครั้งเป็น spell slot เดียว (1 ครั้ง = 2 ระดับ) และแก่ช้าลงสิบเท่า |

## ค่าตามเลเวล (แสดงอย่างเดียว)
| ตาราง | เลเวล: ค่า |
|---|---|
| Cantrips | 1–3: 2 · 4–9: 3 · 10–20: 4 |
| เวทที่เตรียมได้ | 1: 4 · 2: 5 · 3: 6 · 4: 7 · 5: 9 · 6: 10 · 7: 11 · 8: 12 · 9: 14 · 10: 15 · 11–12: 16 · 13–14: 17 · 15–16: 18 · 17: 19 · 18: 20 · 19: 21 · 20: 22 |
| ร่างสัตว์ที่รู้ (Known Forms) | 1: — · 2–3: 4 · 4–7: 6 · 8–20: 8 |
| CR สูงสุดของร่างสัตว์ | 1: — · 2–3: 1/4 · 4–7: 1/2 · 8–20: 1 (บินได้ / Fly) |

## Resources (counter ใช้/คืนได้)
| Resource | จำนวนตามเลเวล | คืนเมื่อ |
|---|---|---|
| Wild Shape | 2–5: 2 · 6–16: 3 · 17–20: 4 | Long Rest (คืนหมด), Short Rest คืน 1 |
| Wild Resurgence slot | 5–20: 1 | Long Rest (คืนหมด) |
| Nature Magician | 20: 1 | Long Rest (คืนหมด) |

## JSON (machine-readable)

```json
{
  "id": "druid",
  "name": {
    "th": "ดรูอิด",
    "en": "Druid"
  },
  "hitDie": 8,
  "primary": [
    "wis"
  ],
  "saves": [
    "int",
    "wis"
  ],
  "armor": {
    "th": "เกราะเบา, โล่",
    "en": "Light armor, Shields"
  },
  "weapons": {
    "th": "อาวุธ Simple",
    "en": "Simple weapons"
  },
  "tools": {
    "th": "Herbalism Kit",
    "en": "Herbalism Kit"
  },
  "skillChoices": {
    "count": 2,
    "from": [
      "animal-handling",
      "arcana",
      "insight",
      "medicine",
      "nature",
      "perception",
      "religion",
      "survival"
    ]
  },
  "caster": "full",
  "castAbility": "wis",
  "tables": [
    {
      "id": "cantrips",
      "name": {
        "th": "Cantrips",
        "en": "Cantrips"
      },
      "values": [
        "2",
        "2",
        "2",
        "3",
        "3",
        "3",
        "3",
        "3",
        "3",
        "4",
        "4",
        "4",
        "4",
        "4",
        "4",
        "4",
        "4",
        "4",
        "4",
        "4"
      ]
    },
    {
      "id": "prepared",
      "name": {
        "th": "เวทที่เตรียมได้",
        "en": "Prepared Spells"
      },
      "values": [
        "4",
        "5",
        "6",
        "7",
        "9",
        "10",
        "11",
        "12",
        "14",
        "15",
        "16",
        "16",
        "17",
        "17",
        "18",
        "18",
        "19",
        "20",
        "21",
        "22"
      ]
    },
    {
      "id": "beast-forms",
      "name": {
        "th": "ร่างสัตว์ที่รู้ (Known Forms)",
        "en": "Known Forms"
      },
      "values": [
        "—",
        "4",
        "4",
        "6",
        "6",
        "6",
        "6",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8",
        "8"
      ]
    },
    {
      "id": "beast-max-cr",
      "name": {
        "th": "CR สูงสุดของร่างสัตว์",
        "en": "Max Beast CR"
      },
      "values": [
        "—",
        "1/4",
        "1/4",
        "1/2",
        "1/2",
        "1/2",
        "1/2",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)",
        "1 (บินได้ / Fly)"
      ]
    }
  ],
  "resources": [
    {
      "id": "wild-shape",
      "name": {
        "th": "Wild Shape",
        "en": "Wild Shape"
      },
      "max": [
        0,
        2,
        2,
        2,
        2,
        3,
        3,
        3,
        3,
        3,
        3,
        3,
        3,
        3,
        3,
        3,
        4,
        4,
        4,
        4
      ],
      "maxFormula": "",
      "refill": "long",
      "shortRegain": 1,
      "shortFromLevel": 0
    },
    {
      "id": "wild-resurgence-slot",
      "name": {
        "th": "Wild Resurgence (แลก slot ระดับ 1)",
        "en": "Wild Resurgence slot"
      },
      "max": [
        0,
        0,
        0,
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
        1
      ],
      "maxFormula": "",
      "refill": "long",
      "shortRegain": 0,
      "shortFromLevel": 0
    },
    {
      "id": "nature-magician",
      "name": {
        "th": "Nature Magician",
        "en": "Nature Magician"
      },
      "max": [
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
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
        "th": "Spellcasting",
        "en": "Spellcasting"
      },
      "desc": {
        "th": "ร่ายเวทสาย Druid ด้วย Wisdom รู้ cantrip และเตรียมเวทได้ตามตาราง เปลี่ยนรายการเวทที่เตรียมได้หลัง Long Rest และใช้ Druidic Focus เป็น focus ได้ (spell slot ตามตาราง full caster)",
        "en": "Cast Druid spells with Wisdom. Know cantrips and prepare spells per the table, change prepared spells after a Long Rest, and use a Druidic Focus. Spell slots follow the full-caster table."
      }
    },
    {
      "level": 1,
      "name": {
        "th": "Druidic",
        "en": "Druidic"
      },
      "desc": {
        "th": "รู้ภาษาลับของดรูอิด ใช้ทิ้งข้อความซ่อนที่ดรูอิดด้วยกันเห็นทันที และมีเวท Speak with Animals เตรียมไว้เสมอ (ไม่นับโควตา)",
        "en": "You know Druidic, a secret language for hidden messages other druids spot at once, and always have Speak with Animals prepared (free)."
      }
    },
    {
      "level": 1,
      "name": {
        "th": "Primal Order",
        "en": "Primal Order"
      },
      "desc": {
        "th": "เลือก 1 บทบาท: Magician (ได้ cantrip เพิ่ม 1 และบวก Wis mod ขั้นต่ำ +1 ให้เช็ก Arcana/Nature) หรือ Warden (ใช้อาวุธ Martial และเกราะกลางได้)",
        "en": "Pick a role: Magician (one extra cantrip, add Wis mod, min +1, to Arcana and Nature checks) or Warden (Martial weapons and Medium armor training)."
      }
    },
    {
      "level": 2,
      "name": {
        "th": "Wild Shape",
        "en": "Wild Shape"
      },
      "desc": {
        "th": "Bonus Action แปลงเป็นร่างสัตว์ที่รู้ ได้ temp HP เท่าเลเวล Druid อยู่ได้ครึ่งเลเวลชั่วโมง ร่ายเวทไม่ได้ระหว่างแปลง (สมาธิเดิมไม่หลุด) ออกจากร่างด้วย Bonus Action จำนวนร่างและ CR ดูตาราง",
        "en": "Bonus Action: take a known Beast form, gain temp HP equal to Druid level, last hours equal to half your level. No casting while shifted (Concentration holds). Leave as a Bonus Action. Forms and CR per the table."
      }
    },
    {
      "level": 2,
      "name": {
        "th": "Wild Companion",
        "en": "Wild Companion"
      },
      "desc": {
        "th": "Magic action ใช้ spell slot หรือ Wild Shape 1 ครั้ง ร่าย Find Familiar โดยไม่ต้องใช้วัสดุ สัตว์รับใช้เป็น Fey และหายไปเมื่อ Long Rest",
        "en": "Magic action: spend a spell slot or a Wild Shape use to cast Find Familiar with no materials. The familiar is Fey and leaves on a Long Rest."
      }
    },
    {
      "level": 3,
      "name": {
        "th": "ฟีเจอร์ subclass",
        "en": "Subclass feature"
      },
      "desc": {
        "th": "เลือก subclass ของ Druid (เช่น Circle of the Land) แล้วเพิ่ม feature ของ subclass เองเป็น custom ability",
        "en": "Choose a Druid subclass (e.g. Circle of the Land) and add its features yourself as custom abilities."
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
      "level": 5,
      "name": {
        "th": "Wild Resurgence",
        "en": "Wild Resurgence"
      },
      "desc": {
        "th": "เทิร์นละครั้ง ถ้า Wild Shape หมด ใช้ spell slot แลกคืน 1 ครั้งได้ และแลก Wild Shape 1 ครั้งเป็น spell slot ระดับ 1 ได้ 1 ครั้งต่อ Long Rest",
        "en": "Once per turn, if out of Wild Shape, spend a spell slot to regain one use. Once per Long Rest, spend a Wild Shape use to gain a level 1 spell slot."
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
        "th": "Elemental Fury",
        "en": "Elemental Fury"
      },
      "desc": {
        "th": "เลือก 1: Potent Spellcasting (บวก Wis mod ให้ดาเมจ cantrip ของ Druid) หรือ Primal Strike (เทิร์นละครั้ง ตีโดนด้วยอาวุธหรือร่างสัตว์ เพิ่มดาเมจ 1d8 Cold, Fire, Lightning หรือ Thunder)",
        "en": "Choose one: Potent Spellcasting (add Wis mod to Druid cantrip damage) or Primal Strike (once per turn, a weapon or Beast-form hit deals +1d8 Cold, Fire, Lightning or Thunder)."
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
      "level": 10,
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
      "level": 14,
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
      "level": 15,
      "name": {
        "th": "Improved Elemental Fury",
        "en": "Improved Elemental Fury"
      },
      "desc": {
        "th": "ตัวเลือกจาก Elemental Fury แรงขึ้น: Potent Spellcasting ทำให้ cantrip ที่ระยะ 10 ft ขึ้นไปไกลขึ้นอีก 300 ft, Primal Strike เพิ่มเป็น 2d8",
        "en": "Your Elemental Fury pick improves: Potent Spellcasting adds 300 ft to cantrips with 10+ ft range; Primal Strike becomes 2d8."
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
      "level": 18,
      "name": {
        "th": "Beast Spells",
        "en": "Beast Spells"
      },
      "desc": {
        "th": "ร่ายเวทได้ระหว่างอยู่ในร่าง Wild Shape ยกเว้นเวทที่ต้องใช้วัสดุที่มีราคาหรือวัสดุที่ถูกใช้หมดไป",
        "en": "You can cast spells in Beast form, except spells with a costly or consumed Material component."
      }
    },
    {
      "level": 19,
      "name": {
        "th": "Epic Boon",
        "en": "Epic Boon"
      },
      "desc": {
        "th": "ได้ Epic Boon feat 1 อัน หรือ feat อื่นที่เข้าเงื่อนไข (แนะนำ Boon of Dimensional Travel) ให้เพิ่มเองเป็น custom ability",
        "en": "Gain an Epic Boon feat or another feat you qualify for (Boon of Dimensional Travel suggested). Add it as a custom ability."
      }
    },
    {
      "level": 20,
      "name": {
        "th": "Archdruid",
        "en": "Archdruid"
      },
      "desc": {
        "th": "ทอย Initiative ตอน Wild Shape หมดจะคืน 1 ครั้ง, 1 ครั้งต่อ Long Rest แลก Wild Shape หลายครั้งเป็น spell slot เดียว (1 ครั้ง = 2 ระดับ) และแก่ช้าลงสิบเท่า",
        "en": "Regain a Wild Shape use when you roll Initiative with none left. Once per Long Rest, turn Wild Shape uses into one spell slot (2 levels per use). You age ten times slower."
      }
    }
  ]
}
```
