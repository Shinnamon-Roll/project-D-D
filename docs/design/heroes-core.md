# ฮีโร่ชุดแรก / First Company — ข้อมูลเกม v2 (10 ตัว)

ข้อมูลฮีโร่ 10 ตัวของเซ็ต `core` (ฮีโร่ชุดแรก) แปลงจาก [`roster.md`](roster.md) เป็นระบบ v2. id เดิม;
knight → คลาส `fighter`, necromancer → คลาส `wizard`, ที่เหลือคลาสตรงชื่อ. อีกเซ็ตอยู่ใน [`heroes-legend.md`](heroes-legend.md).
กติกา v2 ([`../features/rules-v2.md`](../features/rules-v2.md)): 6 ability, AC, skill proficiency ตามคลาส ([`classes/`](classes/)).
Ability = standard array 15/14/13/12/10/8 จัดให้เข้ากับคลาสและเรื่องของฮีโร่ แล้วบวกโบนัส +2/+1 (สูงสุด 17 ที่เลเวล 1, ผลรวม 75).
AC คำนวณจากเกราะ/โล่ในของเริ่มต้น (ตัวเลขเกราะตาม SRD 5.2.1, CC BY 4.0); HP เลเวล 1 = hit die เต็ม + CON mod. DM แก้ได้ในแอป.
Expertise: Rogue ได้ 2 skill ตั้งแต่เลเวล 1; Bard/Ranger ได้ทีหลัง (DM/ผู้เล่นเพิ่มเอง).
ชื่อ เรื่อง (`heroName`, `blurb`, `role`) คงไว้ตามไฟล์เดิม; ของชิ้นแรก = ของประจำตัวจากไฟล์เดิม. ชื่อ เนื้อเรื่อง และไอเทมแต่งขึ้นเอง

## ตารางสรุป
abilities = STR / DEX / CON / INT / WIS / CHA

| # | ฮีโร่ | คลาส | abilities | AC | HP L1 | skills |
|---|---|---|---|---|---|---|
| 1 | Aldric Vane (อัลดริก เวน) | Fighter (นักรบ) | 17 / 10 / 15 / 8 / 13 / 12 | 18 | 12 | athletics, perception |
| 2 | Soraya Quill (โซรายา ควิลล์) | Wizard (พ่อมด/นักเวทศึกษา) | 8 / 14 / 13 / 17 / 13 / 10 | 12 | 7 | arcana, history |
| 3 | Zarif Kesh (ซารีฟ เคช) | Rogue (นักย่องเบา) | 8 / 17 / 13 / 12 / 10 / 15 | 14 | 9 | sleight-of-hand, stealth, deception, persuasion |
| 4 | Sigrun Ashvale (ซีกรุน แอชเวล) | Ranger (พรานป่า) | 12 / 17 / 13 / 10 / 15 / 8 | 15 | 11 | animal-handling, perception, survival |
| 5 | Imara Solenne (อิมารา โซเลนน์) | Cleric (นักบวช) | 13 / 10 / 15 / 8 / 17 / 12 | 15 | 10 | medicine, persuasion |
| 6 | Brenna Dawnmere (เบรนนา ดอว์นเมียร์) | Paladin (อัศวินศักดิ์สิทธิ์) | 17 / 10 / 14 / 8 / 12 / 14 | 16 | 12 | athletics, intimidation |
| 7 | Batu Redmane (บาตู เรดเมน) | Barbarian (นักรบเถื่อน) | 17 / 13 / 15 / 8 / 12 / 10 | 13 | 14 | athletics, survival |
| 8 | Cosimo Lark (โคซิโม ลาร์ก) | Bard (นักกวีเวทดนตรี) | 8 / 15 / 13 / 12 / 10 / 17 | 13 | 9 | deception, performance, persuasion |
| 9 | Malai of the Deepwood (มาลัย แห่งพงไพร) | Druid (ดรูอิด) | 14 / 12 / 14 / 10 / 17 / 8 | 14 | 10 | animal-handling, nature |
| 10 | Ivo Grimsel (อีโว กริมเซล) | Wizard (พ่อมด/นักเวทศึกษา) | 8 / 12 / 15 / 17 / 13 / 10 | 11 | 8 | medicine, religion |

## 1. นักรบ / Fighter — `knight`
**อัลดริก เวน** / **Aldric Vane** · คลาส `fighter` · role: `tank` (แทงก์)

> อัศวินที่เคยพลาดจนเสียเพื่อนไป จึงสาบานว่าจะไม่ยอมให้ใครล้มก่อนตัวเองอีก  
> *A knight who once lost a comrade and swore no friend would ever fall before him again.*

- **Abilities**: STR 17 (+3) DEX 10 (+0) CON 15 (+2) INT 8 (-1) WIS 13 (+1) CHA 12 (+1)
  - standard array + โบนัส: STR 15+2, CON 14+1, WIS 13, CHA 12, DEX 10, INT 8
- **AC 18**: Chain Mail 16 + Shield 2
- **HP เลเวล 1**: 12 (d10 เต็ม 10 + CON +2)
- **Skills**: athletics, perception
- **ของ**: โล่ตราสิงห์ตระกูลเวน / Vane Family Lion Shield · ชุดเกราะโซ่ / Chain Mail · ดาบยาว / Longsword · หอกซัด / Javelin ×8 · ผ้าพันแผลสนาม / Field Bandages ×3 · ชุดสำรวจดันเจียน / Dungeoneer's Pack (ชิ้นแรก = ของประจำตัว)

## 2. พ่อมด/นักเวทศึกษา / Wizard — `wizard`
**โซรายา ควิลล์** / **Soraya Quill** · คลาส `wizard` · role: `control` (ควบคุม)

> นักปราชญ์จากหอดูดาวกลางทะเลทราย ผู้อ่านดวงดาวได้คล่องเหมือนอ่านแผนที่สนามรบ  
> *A scholar from the desert observatories who reads the stars like a battle map.*

- **Abilities**: STR 8 (-1) DEX 14 (+2) CON 13 (+1) INT 17 (+3) WIS 13 (+1) CHA 10 (+0)
  - standard array + โบนัส: INT 15+2, DEX 14, CON 13, WIS 12+1, CHA 10, STR 8
- **AC 12**: ไม่ใส่เกราะ 10 + DEX
- **HP เลเวล 1**: 7 (d6 เต็ม 6 + CON +1)
- **Skills**: arcana, history
- **ของ**: ไม้เท้าแก้วดารา / Starglass Staff · ตำราเวทปกหนังสีคราม / Indigo-Bound Spellbook · มีดสั้น / Dagger ×2 · แผนที่ดาว / Star Chart · ชุดนักปราชญ์ / Scholar's Pack (ชิ้นแรก = ของประจำตัว)

## 3. นักย่องเบา / Rogue — `rogue`
**ซารีฟ เคช** / **Zarif Kesh** · คลาส `rogue` · role: `damage` (ดาเมจ)

> หัวขโมยเจ้าเสน่ห์จากเมืองคลอง ที่ติดหนี้คนไว้ครึ่งเมือง  
> *A charming thief from the canal city who owes money to half the town.*

- **Abilities**: STR 8 (-1) DEX 17 (+3) CON 13 (+1) INT 12 (+1) WIS 10 (+0) CHA 15 (+2)
  - standard array + โบนัส: DEX 15+2, CHA 14+1, CON 13, INT 12, WIS 10, STR 8
- **AC 14**: Leather Armor 11 + DEX +3
- **HP เลเวล 1**: 9 (d8 เต็ม 8 + CON +1)
- **Skills**: sleight-of-hand, stealth, deception, persuasion · **Expertise**: sleight-of-hand, deception
- **ของ**: มีดสั้นคู่เงินจันทรา / Twin Moonsilver Daggers · เกราะหนัง / Leather Armor · ธนูสั้น / Shortbow · ลูกธนู / Arrows ×20 · เครื่องมือโจร / Thieves' Tools · ชุดย่องเบา / Burglar's Pack (ชิ้นแรก = ของประจำตัว)

## 4. พรานป่า / Ranger — `ranger`
**ซีกรุน แอชเวล** / **Sigrun Ashvale** · คลาส `ranger` · role: `damage` (ดาเมจ)

> นักล่าเงียบขรึมจากป่าสนทางเหนือ ที่ไว้ใจเหยี่ยวคู่ใจมากกว่าคน  
> *A quiet huntress from the northern pines who trusts her falcon more than people.*

- **Abilities**: STR 12 (+1) DEX 17 (+3) CON 13 (+1) INT 10 (+0) WIS 15 (+2) CHA 8 (-1)
  - standard array + โบนัส: DEX 15+2, WIS 14+1, CON 13, STR 12, INT 10, CHA 8
- **AC 15**: Studded Leather Armor 12 + DEX +3
- **HP เลเวล 1**: 11 (d10 เต็ม 10 + CON +1)
- **Skills**: animal-handling, perception, survival
- **ของ**: ธนูยาวไม้แอช / Ashwood Longbow · เกราะหนังหมุด / Studded Leather Armor · ลูกธนู / Arrows ×20 · ดาบสั้น / Shortsword ×2 · ถุงมือหนังเหยี่ยว / Falconer's Glove · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 5. นักบวช / Cleric — `cleric`
**อิมารา โซเลนน์** / **Imara Solenne** · คลาส `cleric` · role: `healer` (ฮีลเลอร์)

> หมอแห่งวิหารสุริยันที่หัวเราะเสียงดัง สวดเสียงเบา และไม่เคยทิ้งใครไว้ข้างหลัง  
> *A sun-temple healer who laughs loudly, prays quietly, and never leaves anyone behind.*

- **Abilities**: STR 13 (+1) DEX 10 (+0) CON 15 (+2) INT 8 (-1) WIS 17 (+3) CHA 12 (+1)
  - standard array + โบนัส: WIS 15+2, CON 14+1, STR 13, CHA 12, DEX 10, INT 8
- **AC 15**: Chain Shirt 13 + DEX +0 + Shield 2
- **HP เลเวล 1**: 10 (d8 เต็ม 8 + CON +2)
- **Skills**: medicine, persuasion
- **ของ**: จี้สุริยันศักดิ์สิทธิ์ / Sunburst Holy Pendant · เสื้อเกราะโซ่ / Chain Shirt · โล่ / Shield · กระบอง / Mace · ชุดผ้าพันแผล / Healer's Kit · ชุดนักบวช / Priest's Pack (ชิ้นแรก = ของประจำตัว)

## 6. อัศวินศักดิ์สิทธิ์ / Paladin — `paladin`
**เบรนนา ดอว์นเมียร์** / **Brenna Dawnmere** · คลาส `paladin` · role: `hybrid` (ไฮบริด)

> พาลาดินพูดตรงใจใหญ่ ผู้สาบานจะปกป้องคนอ่อนแอ และถือคำสาบานนั้นแบบตรงตัวทุกคำ  
> *A blunt, big-hearted paladin sworn to protect the weak, who takes that oath word for word.*

- **Abilities**: STR 17 (+3) DEX 10 (+0) CON 14 (+2) INT 8 (-1) WIS 12 (+1) CHA 14 (+2)
  - standard array + โบนัส: STR 15+2, CHA 14, CON 13+1, WIS 12, DEX 10, INT 8
- **AC 16**: Chain Mail 16
- **HP เลเวล 1**: 12 (d10 เต็ม 10 + CON +2)
- **Skills**: athletics, intimidation
- **ของ**: ดาบใหญ่แห่งคำปฏิญาณ / Oathbound Greatsword · ชุดเกราะโซ่ / Chain Mail · หอกซัด / Javelin ×5 · สัญลักษณ์ศักดิ์สิทธิ์ / Holy Symbol · ชุดนักบวช / Priest's Pack (ชิ้นแรก = ของประจำตัว)

## 7. นักรบเถื่อน / Barbarian — `barbarian`
**บาตู เรดเมน** / **Batu Redmane** · คลาส `barbarian` · role: `damage` (ดาเมจ)

> นักรบจากทุ่งหญ้ากว้างที่ออกจากเผ่ามาตามหาศึกที่คู่ควรจะถูกแต่งเป็นเพลง  
> *A steppe warrior who left his clan to find a fight worthy of a song.*

- **Abilities**: STR 17 (+3) DEX 13 (+1) CON 15 (+2) INT 8 (-1) WIS 12 (+1) CHA 10 (+0)
  - standard array + โบนัส: STR 15+2, CON 14+1, DEX 13, WIS 12, CHA 10, INT 8
- **AC 13**: Unarmored Defense 10 + DEX + CON
- **HP เลเวล 1**: 14 (d12 เต็ม 12 + CON +2)
- **Skills**: athletics, survival
- **ของ**: ขวานยักษ์ตระกูลเรดเมน / Redmane Greataxe · ขวานมือ / Handaxe ×4 · แตรเขาสัตว์ / Steppe Horn · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 8. นักกวีเวทดนตรี / Bard — `bard`
**โคซิโม ลาร์ก** / **Cosimo Lark** · คลาส `bard` · role: `support` (สนับสนุน)

> นักแต่งเพลงพเนจรที่สาบานว่าทุกเรื่องที่เล่าเป็นเรื่องจริง โดยเฉพาะเรื่องของตัวเอง  
> *A wandering songwriter who swears every tale he tells is true, especially the ones about himself.*

- **Abilities**: STR 8 (-1) DEX 15 (+2) CON 13 (+1) INT 12 (+1) WIS 10 (+0) CHA 17 (+3)
  - standard array + โบนัส: CHA 15+2, DEX 14+1, CON 13, INT 12, WIS 10, STR 8
- **AC 13**: Leather Armor 11 + DEX +2
- **HP เลเวล 1**: 9 (d8 เต็ม 8 + CON +1)
- **Skills**: deception, performance, persuasion
- **ของ**: พิณลูทสายเงิน / Silver-Stringed Lute · เกราะหนัง / Leather Armor · มีดสั้น / Dagger ×2 · สมุดเพลง / Songbook · ชุดนักแสดง / Entertainer's Pack (ชิ้นแรก = ของประจำตัว)

## 9. ดรูอิด / Druid — `druid`
**มาลัย แห่งพงไพร** / **Malai of the Deepwood** · คลาส `druid` · role: `hybrid` (ไฮบริด)

> ผู้พิทักษ์ป่าลึกที่พูดเสียงนุ่ม แต่กลายร่างเป็นหมีได้เมื่อจำเป็น  
> *A guardian of the deep jungle who speaks softly and becomes a bear when she must.*

- **Abilities**: STR 14 (+2) DEX 12 (+1) CON 14 (+2) INT 10 (+0) WIS 17 (+3) CHA 8 (-1)
  - standard array + โบนัส: WIS 15+2, CON 14, STR 13+1, DEX 12, INT 10, CHA 8
- **AC 14**: Leather Armor 11 + DEX +1 + Shield 2
- **HP เลเวล 1**: 10 (d8 เต็ม 8 + CON +2)
- **Skills**: animal-handling, nature
- **ของ**: ไม้เท้าเถาวัลย์มีชีวิต / Living Vine Staff · เกราะหนัง / Leather Armor · โล่ไม้ / Wooden Shield · เคียว / Sickle · ชุดสมุนไพร / Herbalism Kit · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 10. พ่อมด/นักเวทศึกษา / Wizard — `necromancer`
**อีโว กริมเซล** / **Ivo Grimsel** · คลาส `wizard` · role: `control` (ควบคุม)

> ลูกชายคนเฝ้าสุสาน ที่ฟังเสียงคนตายเพราะคนเป็นไม่เคยฟังเขา  
> *A gravekeeper's son who listens to the dead because the living never listened to him.*

- **Abilities**: STR 8 (-1) DEX 12 (+1) CON 15 (+2) INT 17 (+3) WIS 13 (+1) CHA 10 (+0)
  - standard array + โบนัส: INT 15+2, CON 14+1, WIS 13, DEX 12, CHA 10, STR 8
- **AC 11**: ไม่ใส่เกราะ 10 + DEX
- **HP เลเวล 1**: 8 (d6 เต็ม 6 + CON +2)
- **Skills**: medicine, religion
- **ของ**: ตะเกียงกะโหลก / Skull Lantern · ตำราเวทปกดำ / Black-Bound Spellbook · มีดสั้น / Dagger ×2 · พลั่วขุดหลุม / Gravedigger's Spade · ชุดนักปราชญ์ / Scholar's Pack (ชิ้นแรก = ของประจำตัว)

## JSON (machine-readable)

```json
[
  {
    "id": "knight",
    "set": "core",
    "classId": "fighter",
    "heroName": {
      "th": "อัลดริก เวน",
      "en": "Aldric Vane"
    },
    "blurb": {
      "th": "อัศวินที่เคยพลาดจนเสียเพื่อนไป จึงสาบานว่าจะไม่ยอมให้ใครล้มก่อนตัวเองอีก",
      "en": "A knight who once lost a comrade and swore no friend would ever fall before him again."
    },
    "role": "tank",
    "art": "",
    "abilities": {
      "str": 17,
      "dex": 10,
      "con": 15,
      "int": 8,
      "wis": 13,
      "cha": 12
    },
    "ac": 18,
    "skills": [
      "athletics",
      "perception"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "โล่ตราสิงห์ตระกูลเวน",
          "en": "Vane Family Lion Shield"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดเกราะโซ่",
          "en": "Chain Mail"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ดาบยาว",
          "en": "Longsword"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "หอกซัด",
          "en": "Javelin"
        },
        "qty": 8
      },
      {
        "name": {
          "th": "ผ้าพันแผลสนาม",
          "en": "Field Bandages"
        },
        "qty": 3
      },
      {
        "name": {
          "th": "ชุดสำรวจดันเจียน",
          "en": "Dungeoneer's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "wizard",
    "set": "core",
    "classId": "wizard",
    "heroName": {
      "th": "โซรายา ควิลล์",
      "en": "Soraya Quill"
    },
    "blurb": {
      "th": "นักปราชญ์จากหอดูดาวกลางทะเลทราย ผู้อ่านดวงดาวได้คล่องเหมือนอ่านแผนที่สนามรบ",
      "en": "A scholar from the desert observatories who reads the stars like a battle map."
    },
    "role": "control",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 14,
      "con": 13,
      "int": 17,
      "wis": 13,
      "cha": 10
    },
    "ac": 12,
    "skills": [
      "arcana",
      "history"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "ไม้เท้าแก้วดารา",
          "en": "Starglass Staff"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ตำราเวทปกหนังสีคราม",
          "en": "Indigo-Bound Spellbook"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มีดสั้น",
          "en": "Dagger"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "แผนที่ดาว",
          "en": "Star Chart"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดนักปราชญ์",
          "en": "Scholar's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "rogue",
    "set": "core",
    "classId": "rogue",
    "heroName": {
      "th": "ซารีฟ เคช",
      "en": "Zarif Kesh"
    },
    "blurb": {
      "th": "หัวขโมยเจ้าเสน่ห์จากเมืองคลอง ที่ติดหนี้คนไว้ครึ่งเมือง",
      "en": "A charming thief from the canal city who owes money to half the town."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 17,
      "con": 13,
      "int": 12,
      "wis": 10,
      "cha": 15
    },
    "ac": 14,
    "skills": [
      "sleight-of-hand",
      "stealth",
      "deception",
      "persuasion"
    ],
    "expertise": [
      "sleight-of-hand",
      "deception"
    ],
    "items": [
      {
        "name": {
          "th": "มีดสั้นคู่เงินจันทรา",
          "en": "Twin Moonsilver Daggers"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เกราะหนัง",
          "en": "Leather Armor"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ธนูสั้น",
          "en": "Shortbow"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ลูกธนู",
          "en": "Arrows"
        },
        "qty": 20
      },
      {
        "name": {
          "th": "เครื่องมือโจร",
          "en": "Thieves' Tools"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดย่องเบา",
          "en": "Burglar's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "ranger",
    "set": "core",
    "classId": "ranger",
    "heroName": {
      "th": "ซีกรุน แอชเวล",
      "en": "Sigrun Ashvale"
    },
    "blurb": {
      "th": "นักล่าเงียบขรึมจากป่าสนทางเหนือ ที่ไว้ใจเหยี่ยวคู่ใจมากกว่าคน",
      "en": "A quiet huntress from the northern pines who trusts her falcon more than people."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 12,
      "dex": 17,
      "con": 13,
      "int": 10,
      "wis": 15,
      "cha": 8
    },
    "ac": 15,
    "skills": [
      "animal-handling",
      "perception",
      "survival"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "ธนูยาวไม้แอช",
          "en": "Ashwood Longbow"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เกราะหนังหมุด",
          "en": "Studded Leather Armor"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ลูกธนู",
          "en": "Arrows"
        },
        "qty": 20
      },
      {
        "name": {
          "th": "ดาบสั้น",
          "en": "Shortsword"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "ถุงมือหนังเหยี่ยว",
          "en": "Falconer's Glove"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดเดินทาง",
          "en": "Explorer's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "cleric",
    "set": "core",
    "classId": "cleric",
    "heroName": {
      "th": "อิมารา โซเลนน์",
      "en": "Imara Solenne"
    },
    "blurb": {
      "th": "หมอแห่งวิหารสุริยันที่หัวเราะเสียงดัง สวดเสียงเบา และไม่เคยทิ้งใครไว้ข้างหลัง",
      "en": "A sun-temple healer who laughs loudly, prays quietly, and never leaves anyone behind."
    },
    "role": "healer",
    "art": "",
    "abilities": {
      "str": 13,
      "dex": 10,
      "con": 15,
      "int": 8,
      "wis": 17,
      "cha": 12
    },
    "ac": 15,
    "skills": [
      "medicine",
      "persuasion"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "จี้สุริยันศักดิ์สิทธิ์",
          "en": "Sunburst Holy Pendant"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เสื้อเกราะโซ่",
          "en": "Chain Shirt"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "โล่",
          "en": "Shield"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "กระบอง",
          "en": "Mace"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดผ้าพันแผล",
          "en": "Healer's Kit"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดนักบวช",
          "en": "Priest's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "paladin",
    "set": "core",
    "classId": "paladin",
    "heroName": {
      "th": "เบรนนา ดอว์นเมียร์",
      "en": "Brenna Dawnmere"
    },
    "blurb": {
      "th": "พาลาดินพูดตรงใจใหญ่ ผู้สาบานจะปกป้องคนอ่อนแอ และถือคำสาบานนั้นแบบตรงตัวทุกคำ",
      "en": "A blunt, big-hearted paladin sworn to protect the weak, who takes that oath word for word."
    },
    "role": "hybrid",
    "art": "",
    "abilities": {
      "str": 17,
      "dex": 10,
      "con": 14,
      "int": 8,
      "wis": 12,
      "cha": 14
    },
    "ac": 16,
    "skills": [
      "athletics",
      "intimidation"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "ดาบใหญ่แห่งคำปฏิญาณ",
          "en": "Oathbound Greatsword"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดเกราะโซ่",
          "en": "Chain Mail"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "หอกซัด",
          "en": "Javelin"
        },
        "qty": 5
      },
      {
        "name": {
          "th": "สัญลักษณ์ศักดิ์สิทธิ์",
          "en": "Holy Symbol"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดนักบวช",
          "en": "Priest's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "barbarian",
    "set": "core",
    "classId": "barbarian",
    "heroName": {
      "th": "บาตู เรดเมน",
      "en": "Batu Redmane"
    },
    "blurb": {
      "th": "นักรบจากทุ่งหญ้ากว้างที่ออกจากเผ่ามาตามหาศึกที่คู่ควรจะถูกแต่งเป็นเพลง",
      "en": "A steppe warrior who left his clan to find a fight worthy of a song."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 17,
      "dex": 13,
      "con": 15,
      "int": 8,
      "wis": 12,
      "cha": 10
    },
    "ac": 13,
    "skills": [
      "athletics",
      "survival"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "ขวานยักษ์ตระกูลเรดเมน",
          "en": "Redmane Greataxe"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขวานมือ",
          "en": "Handaxe"
        },
        "qty": 4
      },
      {
        "name": {
          "th": "แตรเขาสัตว์",
          "en": "Steppe Horn"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดเดินทาง",
          "en": "Explorer's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "bard",
    "set": "core",
    "classId": "bard",
    "heroName": {
      "th": "โคซิโม ลาร์ก",
      "en": "Cosimo Lark"
    },
    "blurb": {
      "th": "นักแต่งเพลงพเนจรที่สาบานว่าทุกเรื่องที่เล่าเป็นเรื่องจริง โดยเฉพาะเรื่องของตัวเอง",
      "en": "A wandering songwriter who swears every tale he tells is true, especially the ones about himself."
    },
    "role": "support",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 15,
      "con": 13,
      "int": 12,
      "wis": 10,
      "cha": 17
    },
    "ac": 13,
    "skills": [
      "deception",
      "performance",
      "persuasion"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "พิณลูทสายเงิน",
          "en": "Silver-Stringed Lute"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เกราะหนัง",
          "en": "Leather Armor"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มีดสั้น",
          "en": "Dagger"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "สมุดเพลง",
          "en": "Songbook"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดนักแสดง",
          "en": "Entertainer's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "druid",
    "set": "core",
    "classId": "druid",
    "heroName": {
      "th": "มาลัย แห่งพงไพร",
      "en": "Malai of the Deepwood"
    },
    "blurb": {
      "th": "ผู้พิทักษ์ป่าลึกที่พูดเสียงนุ่ม แต่กลายร่างเป็นหมีได้เมื่อจำเป็น",
      "en": "A guardian of the deep jungle who speaks softly and becomes a bear when she must."
    },
    "role": "hybrid",
    "art": "",
    "abilities": {
      "str": 14,
      "dex": 12,
      "con": 14,
      "int": 10,
      "wis": 17,
      "cha": 8
    },
    "ac": 14,
    "skills": [
      "animal-handling",
      "nature"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "ไม้เท้าเถาวัลย์มีชีวิต",
          "en": "Living Vine Staff"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เกราะหนัง",
          "en": "Leather Armor"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "โล่ไม้",
          "en": "Wooden Shield"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เคียว",
          "en": "Sickle"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดสมุนไพร",
          "en": "Herbalism Kit"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดเดินทาง",
          "en": "Explorer's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "necromancer",
    "set": "core",
    "classId": "wizard",
    "heroName": {
      "th": "อีโว กริมเซล",
      "en": "Ivo Grimsel"
    },
    "blurb": {
      "th": "ลูกชายคนเฝ้าสุสาน ที่ฟังเสียงคนตายเพราะคนเป็นไม่เคยฟังเขา",
      "en": "A gravekeeper's son who listens to the dead because the living never listened to him."
    },
    "role": "control",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 12,
      "con": 15,
      "int": 17,
      "wis": 13,
      "cha": 10
    },
    "ac": 11,
    "skills": [
      "medicine",
      "religion"
    ],
    "expertise": [],
    "items": [
      {
        "name": {
          "th": "ตะเกียงกะโหลก",
          "en": "Skull Lantern"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ตำราเวทปกดำ",
          "en": "Black-Bound Spellbook"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มีดสั้น",
          "en": "Dagger"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "พลั่วขุดหลุม",
          "en": "Gravedigger's Spade"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดนักปราชญ์",
          "en": "Scholar's Pack"
        },
        "qty": 1
      }
    ]
  }
]
```
