# D&D พื้นฐาน: เกมคืออะไร เล่นยังไง กฎที่ระบบต้องรู้

ผลวิจัยวันที่ 2026-10-16 (agent `researcher`) — อ่านจาก SRD PDF โดยตรง
`[1 pN]` = SRD 5.2.1 หน้า N · แหล่งอ้างอิงท้ายไฟล์ · "(อนุมานเอง)" = สูตรที่ตรวจกับตารางใน SRD แล้วแต่ SRD ไม่ได้เขียนเป็นสูตร

## สรุปสั้น
- **D&D** = เกมเล่าเรื่องร่วมกัน คนหนึ่งเป็น **DM/GM** (บรรยายโลก, เล่น NPC, ตัดสิน) ที่เหลือเล่นตัวละครคนละตัว เหมาะ 4–6 คน + DM [1 p5][6]
- **ลูปหลัก**: DM บรรยาย → ผู้เล่นบอกว่าทำอะไร → DM ตัดสิน (ทอย d20 เฉพาะเมื่อผลไม่แน่นอน) → DM เล่าผล [1 p5–6]
- **กลไกหลัก** = **D20 Test**: d20 + modifier + proficiency ≥ DC/AC ผ่าน [1 p6–7]
- **Grid เป็นทางเลือก**: เล่นแบบไม่มีแมพ ("theater of the mind") ก็ได้ ใช้ grid ตอนต่อสู้ใหญ่ [1 p14][10][11]
- **ลิขสิทธิ์**: ใช้ได้เฉพาะ **SRD (CC-BY-4.0)** แนะนำ **SRD 5.2.1 (กฎ 2024)** · Basic/Free Rules และ adventure ทางการ **ไม่ใช่** open content · ห้ามตั้งชื่อแอปว่า "D&D"/"Dungeons & Dragons" (trademark)

## 1. เกมนี้คืออะไร
- **บทบาท**: GM บอกว่าตัวละครอยู่ที่ไหน มีอะไรรอบตัว และเล่นเป็น NPC ทุกตัว [1 p5, p10]; DM เป็น "คนเล่าเรื่องหลักและกรรมการ" มีอำนาจสุดท้าย [7]
- **โครงสร้าง**:
  - **Session**: เล่นหนึ่งครั้ง ประกอบด้วย encounter หลายอัน [6]
  - **Adventure**: เนื้อเรื่องหนึ่งชุด มี hook, plot, แมพ, encounter [6]
  - **Campaign**: หลาย adventure ต่อกัน; หรือ **one-shot** จบในครั้งเดียว [8]
- **Tier ตามเลเวล**: 1–4 ภัยระดับท้องถิ่น · 5–10 · 11–16 · 17–20 ภัยระดับโลก [1 p24]
- **3 เสาหลัก** [1 p5]:
  - **Social**: คุยกับ NPC ท่าที Friendly / Indifferent / Hostile ใช้ roleplay + action Influence [1 p10]
  - **Exploration**: แสง (Bright / Dim / Darkness), พื้นที่มองเห็นยาก, ซ่อนตัว, กับดัก, เดินทาง [1 p11–12]
  - **Combat**: ต่อสู้แบบผลัดตา (ข้อ 4)
- **Theater of the mind vs grid**: SRD 5.2.1 มีหมวด "Playing on a Grid" เป็นทางเลือก [1 p14]; SRD 5.1 ไม่มีกฎ grid เลย; D&D Beyond แนะนำผสม — สู้เล็กเล่าเอา สู้ใหญ่ใช้แมพ [10][11]

## 2. ลูปการเล่น
1. "GM บรรยายฉาก"
2. "ผู้เล่นบอกว่าตัวละครทำอะไร"
3. "GM เล่าผล" [1 p5]

- ทอยเต๋าเฉพาะ "เมื่อผลไม่แน่นอนและน่าสนใจต่อเรื่อง" GM ตั้ง DC แล้วบอกผล [1 p6]
- **ข้อยกเว้นชนะกฎทั่วไป**: ความสามารถ/เวทที่ขัดกฎทั่วไป ให้ยึดความสามารถนั้น [1 p5] — ข้อจำกัดสำคัญถ้าจะทำ rules engine

## 3. กลไกหลักที่ระบบต้องจำลอง
### ลูกเต๋า
- ชุดมาตรฐาน: d4, d6, d8, d10, d12, d20 [6]
- Notation: `3d8+5` = ทอย d8 สามลูก รวม บวก 5 [7]
- d100 = d10 สองลูก (หลักสิบ + หลักหน่วย, 00+0 = 100); d3 = d6 หาร 2 ปัดขึ้น; d2 = คี่/คู่ [7]

### D20 Test [1 p6–7]
- 3 แบบ: **ability check**, **saving throw**, **attack roll**
- รวม = d20 + ability modifier + proficiency bonus (ถ้าถนัด) + โบนัสอื่น → ผ่านเมื่อ ≥ เป้า (DC หรือ AC)
- DC ทั่วไป: ง่ายมาก 5 · ง่าย 10 · กลาง 15 · ยาก 20 · ยากมาก 25 · แทบเป็นไปไม่ได้ 30
- Attack: ใช้ STR (ระยะประชิด), DEX (ยิงไกล), ค่าร่ายเวท (spell attack); อาวุธ Finesse เลือก STR หรือ DEX ได้
- **เฉพาะ attack**: ทอยได้ 20 เต็ม = โดนเสมอ (critical), ได้ 1 = พลาดเสมอ · check/save ไม่มีผ่าน/ตกอัตโนมัติ
- AC พื้นฐาน = 10 + DEX modifier

### Advantage / Disadvantage [1 p7–8]
- ทอย 2d20 เอาค่าสูง (Advantage) / ต่ำ (Disadvantage)
- หลายแหล่งก็ยังเป็น 2 ลูก ไม่ซ้อนกัน; มีทั้งสองอย่าง (กี่อันก็ตาม) = หักล้าง ทอยลูกเดียวปกติ
- **Heroic Inspiration**: ทอยใหม่ได้หนึ่งลูก ถือได้ทีละ 1

### Ability score [1 p5–6]
- 6 ค่า: **STR, DEX, CON, INT, WIS, CHA** ปกติ 1–20 (มอนสเตอร์ถึง 30)
- **Modifier = floor((score − 10) / 2)** (อนุมานเอง ตรงทุกแถวในตาราง) เช่น 1 → −5, 10–11 → +0, 30 → +10
- ปัดลงเสมอ เว้นแต่กฎบอกเป็นอย่างอื่น

### Proficiency bonus [1 p9, p23]
- เลเวล/CR ≤4 → +2 · 5–8 → +3 · 9–12 → +4 · 13–16 → +5 · 17–20 → +6 (CR 21–30 ขึ้นถึง +9)
- **สูตร = 2 + floor((level − 1) / 4)** (อนุมานเอง)
- บวกได้ครั้งเดียวต่อการทอย; **Expertise** = คูณ 2

### Skill 18 อย่าง [1 p9]
| Ability | Skill |
|---|---|
| STR | Athletics |
| DEX | Acrobatics, Sleight of Hand, Stealth |
| INT | Arcana, History, Investigation, Nature, Religion |
| WIS | Animal Handling, Insight, Medicine, Perception, Survival |
| CHA | Deception, Intimidation, Performance, Persuasion |

### ค่าที่คำนวณได้ [1 p22–23]
- Passive Perception = 10 + WIS (Perception) modifier (+5 ถ้า Advantage, −5 ถ้า Disadvantage)
- Initiative แบบคงที่ = 10 + DEX modifier
- Spell save DC = 8 + spellcasting modifier + proficiency
- Spell attack bonus = spellcasting modifier + proficiency

### HP และ damage [1 p16–18]
- HP ลดถึง 0 ไม่ติดลบ; **Bloodied** = HP ≤ ครึ่ง
- Damage = ลูกเต๋า + modifier เดียวกับที่ใช้โจมตี (ขั้นต่ำ 0)
- เวทวงกว้างทอย damage ครั้งเดียวใช้กับทุกเป้า; save ผ่านมักได้ครึ่ง (ปัดลง)
- ลำดับคำนวณ: โบนัส/ลบ → **Resistance** (หาร 2) → **Vulnerability** (คูณ 2); ไม่ซ้อน; **Immunity** = ไม่โดน
- ฮีลไม่เกิน max HP
- **Temp HP** โดนหักก่อน ไม่ซ้อน (เลือกค่าสูงกว่า) หายหลัง Long Rest
- โจมตีประชิดให้ล้มที่ 1 HP + Unconscious แทนตายได้

### HP 0 / Death save [1 p17–18]
- มอนสเตอร์ตายที่ 0 HP (GM เลือกเป็นอื่นได้)
- Damage ที่เหลือหลังถึง 0 ≥ max HP = **ตายทันที**
- ไม่งั้น **Unconscious** + ทอย death save ทุกต้นเทิร์น:
  - d20 ≥ 10 = ผ่าน · ผ่าน 3 = Stable · ตก 3 = ตาย (ไม่ต้องติดกัน)
  - ได้ 1 = ตก 2 ครั้ง · ได้ 20 = ฟื้น 1 HP
- โดน damage ตอน 0 HP = ตก 1 (critical = ตก 2)
- ช่วยให้ Stable: action Help + WIS (Medicine) DC 10; Stable แล้วฟื้น 1 HP หลัง 1d4 ชั่วโมง

### Condition 15 อย่าง [1 glossary]
Blinded, Charmed, Deafened, Exhaustion, Frightened, Grappled, Incapacitated, Invisible, Paralyzed, Petrified, Poisoned, Prone, Restrained, Stunned, Unconscious

ที่มีผลกับการต่อสู้บ่อย:
- **Paralyzed / Unconscious**: โดนตีจากระยะ 5 ฟุต = critical; STR/DEX save ตกอัตโนมัติ
- **Prone**: ศัตรูใน 5 ฟุตได้ Advantage, ไกลกว่านั้นได้ Disadvantage; ลุกใช้ครึ่ง Speed
- **Exhaustion (2024)**: ต่อระดับ −2 ทุก D20 Test และ −5 ฟุต Speed; ระดับ 6 = ตาย; Long Rest ลด 1 ระดับ
- Exhaustion (2014) ต่างกัน เป็นตาราง 6 ขั้น [25]

### Concentration [1 glossary]
- โดน damage → CON save DC = max(10, floor(damage / 2)) สูงสุด 30
- หลุดเมื่อร่ายเวท concentration อันใหม่, Incapacitated, หรือตาย

### พัก [1 glossary]
- **Short Rest**: 1 ชม. ใช้ Hit Dice ได้ (ทอย + CON mod ขั้นต่ำ 1 ต่อลูก)
- **Long Rest**: 8 ชม. (นอน ≥6) ฟื้น HP + Hit Dice ทั้งหมด ลด Exhaustion 1; ต้องมี ≥1 HP; ได้ 1 ครั้งต่อ 16 ชม.
- ถูกขัดเมื่อ: ทอย initiative, ร่ายเวท (ที่ไม่ใช่ cantrip), โดน damage; Long Rest ขัดเพิ่มถ้าออกแรง 1 ชม.

## 4. การต่อสู้ [1 p13–16]
### ลำดับ
- GM จัดตำแหน่ง → ทุกคนทอย **Initiative** (DEX check) → เล่นตามลำดับ; 1 round ≈ 6 วินาที
- มอนสเตอร์ชนิดเดียวกันทอยครั้งเดียวทั้งกลุ่ม
- เสมอ: ผู้เล่นจัดลำดับกันเอง, GM จัดมอนสเตอร์, GM ตัดสินผู้เล่น vs มอนสเตอร์
- **Surprise (2024)**: Disadvantage ตอนทอย initiative (2014: เทิร์นแรกขยับ/ทำอะไรไม่ได้) [2]

### ในหนึ่งเทิร์น
| อย่าง | จำนวน |
|---|---|
| Movement | ถึงค่า Speed แบ่งก่อน/หลัง action ได้ |
| Action | 1 |
| Bonus Action | ≤1 เฉพาะเมื่อมีความสามารถให้ |
| Reaction | 1 ต่อ round (ใช้ได้นอกเทิร์นตัวเอง) |
| หยิบจับของ | ฟรี 1 ครั้ง (ครั้งที่ 2 ใช้ action Utilize) |
| พูดสั้น ๆ, ล้มตัวลง | ฟรี |

### Action
- **2024**: Attack, Dash, Disengage, Dodge, Help, Hide, Influence, Magic, Ready, Search, Study, Utilize [1 p9]
- 2014: Attack, Cast a Spell, Dash, Disengage, Dodge, Help, Hide, Ready, Search, Use an Object [2]
- Hide (2024): DEX (Stealth) DC 15 ขณะมองไม่เห็นหรืออยู่หลังกำบัง ≥3/4 → ได้ condition Invisible

### Grid [1 p14]
- 1 ช่อง = **5 ฟุต**; Speed เป็นช่อง = Speed / 5
- เดินช่องติดกันใช้ 1 ช่อง **รวมแนวทแยง** (ไม่มีกฎ 5/10/5)
- Difficult terrain = 2 ช่อง; เดินทแยงตัดมุมกำแพงไม่ได้
- วัดระยะ: นับจากช่องข้างตัวหนึ่งไปถึงอีกตัว ทางสั้นสุด

### ขนาด [1 p14–15]
| Size | พื้นที่ |
|---|---|
| Tiny | 4 ตัวต่อช่อง |
| Small, Medium | 1 ช่อง |
| Large | 2×2 |
| Huge | 3×3 |
| Gargantuan | 4×4 |

- ผ่านช่องของพวกเดียวกัน, ตัว Incapacitated, Tiny, หรือตัวที่ขนาดต่างกัน ≥2 ระดับได้; ช่องอื่นนับเป็น difficult terrain
- หยุดในช่องที่มีคนอยู่ไม่ได้

### โจมตี [1 p15–16]
- มองไม่เห็นเป้า = Disadvantage; เป้ามองไม่เห็นเรา = Advantage
- ยิงไกล: ระยะปกติ/ระยะไกล (ไกลเกินปกติ = Disadvantage); มีศัตรูใน 5 ฟุตที่มองเห็นเรา = Disadvantage
- Reach ปกติ 5 ฟุต
- **Opportunity Attack**: ศัตรูเดินออกจาก reach → ใช้ Reaction ตีประชิด 1 ครั้งก่อนออก; เลี่ยงได้ด้วย Disengage, teleport, หรือถูกผลักโดยไม่ได้เดินเอง
- **Critical**: ทอยลูกเต๋า damage ทุกลูก 2 รอบ (รวม Sneak Attack ฯลฯ) แล้วบวก modifier ครั้งเดียว

### Cover [1 p15]
- ครึ่ง: +2 AC และ DEX save · 3/4: +5 · ทั้งหมด: เล็งตรงไม่ได้ · ใช้อันที่ดีสุดอันเดียว

### พื้นที่ผลของเวท (AoE) [1 glossary]
- 2024 มี 6 รูป: Cone, Cube, Cylinder, Emanation, Line, Sphere (2014 ไม่มี Emanation)
- Cone กว้างเท่าระยะจากจุดเริ่ม; Emanation เคลื่อนตามตัวผู้ใช้
- **ช่องว่าง**: SRD ทั้งสองฉบับไม่บอกว่า AoE ครอบช่อง grid ไหน (อยู่ใน DMG ซึ่งไม่ใช่ SRD) → **โปรเจกต์ต้องเขียนกฎแปลง AoE → ช่องเอง ห้าม copy DMG**

## 5. ตัวละคร [1 p19–25][9]
### ข้อมูลบน character sheet (สร้าง layout เอง ห้าม copy ชีททางการ — ชีทให้ถ่ายเอกสารใช้ส่วนตัวเท่านั้น)
- ตัวตน: ชื่อ, background, class, subclass, species, level, XP
- ป้องกัน: AC, shield · HP ปัจจุบัน/max/temp · Hit Dice max/ใช้ไป · death save ผ่าน/ตก
- Proficiency bonus · 6 ability (ค่า + mod) พร้อม save และ skill
- Heroic Inspiration, initiative, speed, size, passive Perception
- ตารางโจมตี: ชื่อ, โบนัสโจมตี/DC, damage + ชนิด, หมายเหตุ
- Class features, species traits, feats
- ความถนัด: เกราะ (Light/Medium/Heavy/Shield), อาวุธ, เครื่องมือ
- Spellcasting: ability, mod, save DC, attack bonus · spell slot ระดับ 1–9 (ทั้งหมด/ใช้ไป) · ลิสต์เวท
- รูปลักษณ์, ประวัติ, นิสัย, alignment, ภาษา, ของ, magic item attunement
- เงิน: CP, SP, EP, GP, PP

### ขั้นตอนสร้างตัวละคร
1. เลือก **class**
2. **Origin**: background, species, ภาษา (Common + 2)
3. **Ability score** — เลือก 1 วิธี:
   - Standard array: 15, 14, 13, 12, 10, 8
   - ทอย 4d6 ตัดลูกต่ำสุด × 6
   - Point buy 27 แต้ม (8 = 0, 9 = 1, 10 = 2, 11 = 3, 12 = 4, 13 = 5, 14 = 7, 15 = 9)
   - แล้ว background บวก +2/+1 หรือ +1/+1/+1 (ไม่เกิน 20)
4. Alignment
5. ใส่ตัวเลขที่คำนวณ

### HP เลเวล 1
| HP | Class | Hit Die | HP/เลเวลแบบคงที่ |
|---|---|---|---|
| 12 + CON | Barbarian | d12 | 7 + CON |
| 10 + CON | Fighter, Paladin, Ranger | d10 | 6 + CON |
| 8 + CON | Bard, Cleric, Druid, Monk, Rogue, Warlock | d8 | 5 + CON |
| 6 + CON | Sorcerer, Wizard | d6 | 4 + CON |

### XP ที่ต้องใช้ เลเวล 1–20
0; 300; 900; 2,700; 6,500; 14,000; 23,000; 34,000; 48,000; 64,000; 85,000; 100,000; 120,000; 140,000; 165,000; 195,000; 225,000; 265,000; 305,000; 355,000

### Spell slot (full caster)
| Level | Slot ต่อระดับเวท 1/2/3/… |
|---|---|
| 1 | 2 |
| 2 | 3 |
| 3 | 4/2 |
| 4 | 4/3 |
| 5 | 4/3/2 |
| 7 | 4/3/3/1 |
| 9 | 4/3/3/3/1 |
| 11 | 4/3/3/3/2/1 |
| 13 | เพิ่ม slot ระดับ 7 |
| 15 | เพิ่ม slot ระดับ 8 |
| 17 | 4/3/3/3/2/1/1/1/1 |
| 20 | 4/3/3/3/3/2/2/1/1 |

### เนื้อหาที่มีใน SRD
- **SRD 5.2.1**:
  - 12 class มี subclass ละ 1: Berserker, Lore, Life, Land, Champion, Open Hand, Devotion, Hunter, Thief, Draconic, Fiend, Evoker
  - 4 background: Acolyte, Criminal, Sage, Soldier
  - 9 species: Dragonborn, Dwarf, Elf, Gnome, Goliath, Halfling, Human, Orc, Tiefling
  - feats (API นับ 17), อุปกรณ์, magic item, มอนสเตอร์ A–Z
- **SRD 5.1**: 12 class (subclass ละ 1), 9 race (มี Half-Elf, Half-Orc), background 1 (Acolyte), feat 1 (Grappler) [2][13]

## 6. ฉบับกฎ + ลิขสิทธิ์
| ฉบับ | กฎ | License | วันที่ |
|---|---|---|---|
| SRD 5.1 | 2014 | CC-BY-4.0 (ตั้งแต่ ม.ค. 2023) + OGL 1.0a | — |
| SRD 5.2.0 | 2024 | CC-BY-4.0 | 2025-04-22 |
| SRD 5.2.1 | 2024 | CC-BY-4.0 | 2025-05-01 (อังกฤษ) |

- 5.2.1 แก้ errata + เพิ่ม magic item 15 ชิ้นที่ตกหล่น [3]; ใช้ 5.1 + 5.2 ร่วมกันได้ถ้าให้เครดิตถูก [3]
- SRD ในอนาคตจะเป็น Creative Commons "อย่างถาวร เพิกถอนไม่ได้" [4]

**ข้อความ attribution ที่ต้องใส่ (SRD 5.2.1 คำต่อคำ)** [1 p1]:

> This work includes material from the System Reference Document 5.2.1 ("SRD 5.2.1") by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd. The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International License, available at https://creativecommons.org/licenses/by/4.0/legalcode.

- ห้ามใส่เครดิต Wizards อย่างอื่นเพิ่ม; เขียนว่า "compatible with fifth edition" / "5E compatible" ได้ [1]

**สิ่งที่ไม่ใช่ open content (ห้ามใส่ในแอป)**:
- Basic Rules / Free Rules "ใช้สร้าง content ไม่ได้" [4][5]
- ชื่อ setting/lore เช่น Phandelver, Stormwreck [4]
- เครื่องหมายการค้า: CC-BY ไม่ครอบคลุม trademark → ห้ามตั้งชื่อแอปว่า "D&D" / "Dungeons & Dragons" (ความมั่นใจกลาง — จากความรู้เรื่อง license ไม่ได้ fetch ซ้ำ)

**ข้อมูลกฎแบบให้โปรแกรมอ่าน**:
| แหล่ง | รายละเอียด |
|---|---|
| **dnd5eapi.co** (5e-bits/5e-srd-api) | โค้ด MIT; endpoint `/api/2024` และ `/api/2014`; ไม่ต้อง auth; self-host ได้ (Node 24 + pnpm); 2024 มี 12 class, 9 species, 4 background, 17 feat, 339 spell, 341 monster, 15 condition, 182 equipment; ข้อมูลยังต้องให้เครดิต SRD [12][13] |
| **Open5e v2** | กรอง `document__key=srd-2024` (331 creature, 339 spell); **ระวัง**: มี content ของเจ้าอื่น (OGL เท่านั้น เช่น Kobold Press) ปนอยู่ ต้องกรอง document key ทุกครั้ง; license โค้ดยังไม่ชัด [14][15][16] |

**แนะนำ: ใช้ SRD 5.2.1 (กฎ 2024)** — เป็นคำแนะนำ ยังไม่ได้ยืนยันกับผู้ใช้
- CC อย่างเดียว ไม่มีภาระ OGL; เป็นฉบับปัจจุบัน; มีกฎ grid ฉบับเดียว; API ทั้งสองรองรับ
- ข้อเสีย: เพื่อนที่เคยเล่นกฎ 2014 จะเจอความต่าง: surprise, exhaustion, race → species, ที่มาของโบนัส ability

## 7. Adventure สำหรับมือใหม่
| Adventure | เลเวล | สถานะ | มีอะไร |
|---|---|---|---|
| Intro to Stormwreck Isle | ไม่ระบุ | ฟรีใน D&D Beyond (ต้องมีบัญชี) | 2 encounter นำเข้า Stormwreck [19] |
| Dragons of Stormwreck Isle (Starter Set) | 1–3 | เสียเงิน, ออก 2022-10-04 | adventure 48 หน้า, กฎ 32 หน้า, ตัวละครสำเร็จรูป 5, เต๋า 6 [18][22] |
| Lost Mine of Phandelver | 1–5 | หมดช่วงรับฟรี 2023-10-06 แล้ว | 64 หน้า, แมพในเล่ม, ตัวละครสำเร็จรูป 5 [20][21] |
| The Delian Tomb | 1 | ไม่ใช่ของ WotC; จาก Matt Colville "Running the Game" #1 (2016) | goblin, ~5 ห้อง, เล่นจบ ~3–4 ชม.; ลิขสิทธิ์ PDF แฟนทำไม่ชัด [23] |
| Heroes of the Borderlands (starter กฎ 2024, 2025-09-16) | ไม่ระบุ | เสียเงิน | poster map 20+, token 200+, การ์ด 200+, เต๋า 11 [24] |

- โต๊ะต้องมี: แมพ, stat block มอนสเตอร์, ตัวละครสำเร็จรูป, handout, initiative tracker
- **ไม่มีอันไหนเป็น SRD** → ห้ามใส่ข้อความ/แมพ/ตัวละครของ adventure ในแอป ให้ DM อัปแมพ/โน้ตเอง + ใช้มอนสเตอร์จาก SRD

## ผลต่อการออกแบบระบบ (สรุปเอง)
- กลไกที่โค้ดได้ง่ายและควรเริ่ม: ทอยเต๋า notation + adv/dis, modifier/proficiency, initiative order, HP/temp HP, death save, condition 15 อย่าง
- ยาก/ควรเลื่อน: rules engine เต็ม (ข้อยกเว้นชนะกฎทั่วไป), AoE → ช่อง grid (ต้องคิดกฎเอง), spell slot ทุก class
- ทุกกฎตัดสินโดย DM ได้เสมอ → ระบบควรช่วยคำนวณ แต่ให้ DM แก้ทับได้

## ความมั่นใจ
| เรื่อง | ระดับ | เหตุผล |
|---|---|---|
| ตัวเลขกฎข้อ 1–5 | สูง | อ่านจาก SRD 5.2.1 PDF โดยตรง |
| สูตร modifier/proficiency | สูง | อนุมานเอง ตรวจกับทุกแถวในตารางแล้ว |
| ฉบับ SRD, วันที่, ข้อความ attribution | สูง | หน้า SRD ของ D&D Beyond + PDF |
| จำนวนใน API | สูง | เรียกจริง 2026-10-16 (จะเปลี่ยนตามเวลา) |
| license โค้ด Open5e | ต่ำ | GitHub ระบุ NOASSERTION |
| SRD ไม่มีกฎ AoE บน grid | กลาง-สูง | grep ข้อความ ไม่ได้อ่านทุกหน้า |
| Delian Tomb | กลาง-ต่ำ | แหล่งรอง/แฟน |
| trademark clause ของ CC | กลาง | จากความรู้ ไม่ได้ fetch ซ้ำ |

## ตรวจด้วยมือ
```sh
# attribution
curl -sO https://media.dndbeyond.com/compendium-images/srd/5.2/SRD_CC_v5.2.1.pdf && pdftotext SRD_CC_v5.2.1.pdf - | head -15
# API
curl -s https://www.dnd5eapi.co/api/2024/monsters | jq .count          # ~341
curl -s https://www.dnd5eapi.co/api/2024/conditions | jq '.results[].name'
curl -s 'https://api.open5e.com/v2/creatures/?document__key=srd-2024&limit=1' | jq .count   # 331
```
```js
// สูตร
[1,3,8,10,11,15,20,30].map(s=>Math.floor((s-10)/2))  // [-5,-4,-1,0,0,2,5,10]
[1,4,5,9,13,17,20].map(l=>2+Math.floor((l-1)/4))     // [2,2,3,4,5,6,6]
```

## ยังไม่ได้ตรวจ
- Quick Reference / Player Reference Sheet ทางการ — ไม่พบฉบับ open license; ทำ quick reference เองจาก SRD
- ตาราง feature ของแต่ละ class และลิสต์เวทเต็ม (มีแค่ตาราง slot ของ full caster)

## แหล่งอ้างอิง
[1] SRD 5.2.1 (PDF) — https://media.dndbeyond.com/compendium-images/srd/5.2/SRD_CC_v5.2.1.pdf
[2] SRD 5.1 CC (PDF) — https://media.dndbeyond.com/compendium-images/srd/5.1/SRD_CC_v5.1.pdf
[3] D&D Beyond SRD — https://www.dndbeyond.com/srd
[4] D&D Beyond, "You Can Now Publish Your Own Creations Using the New Core Rules" — https://www.dndbeyond.com/posts/1949-you-can-now-publish-your-own-creations-using-the
[5] D&D Beyond, "2024 D&D Free Rules" — https://www.dndbeyond.com/posts/1804-start-playing-today-with-the-2024-d-d-free-rules
[6] Free Rules 2024, The Basics — https://www.dndbeyond.com/sources/dnd/br-2024/the-basics
[7] Basic Rules 2014, Introduction — https://www.dndbeyond.com/sources/dnd/basic-rules-2014/introduction
[8] D&D Beyond, How to Play D&D — https://www.dndbeyond.com/en/how-to-play-dnd
[9] Character sheet 2024 (PDF) — https://media.dndbeyond.com/compendium-images/br/ph/character-sheet.pdf
[10] D&D Beyond, Tactical Combat — https://www.dndbeyond.com/posts/784-new-players-guide-how-to-try-tactical-combat
[11] D&D Beyond, DM Guide to Miniatures — https://www.dndbeyond.com/posts/178-a-new-dungeon-masters-guide-to-miniatures
[12] 5e-bits/5e-srd-api — https://github.com/5e-bits/5e-srd-api
[13] dnd5eapi.co — https://www.dnd5eapi.co/api/2024 ; https://www.dnd5eapi.co/api/2014
[14] Open5e API docs — https://open5e.com/api-docs
[15] Open5e documents — https://api.open5e.com/v2/documents/
[16] open5e/open5e-api — https://github.com/open5e/open5e-api
[17] CC BY 4.0 legal code — https://creativecommons.org/licenses/by/4.0/legalcode
[18] WPN, Dragons of Stormwreck Isle — https://wpn.wizards.com/products/dragons-of-stormwreck-isle
[19] D&D Beyond, Intro to Stormwreck Isle — https://www.dndbeyond.com/posts/1582-intro-to-stormwreck-isle-your-free-introduction-to
[20] ComicBook.com, LMoP legacy — https://comicbook.com/gaming/news/dungeons-dragons-lost-mine-phandelver-legacy-content-dnd-beyond/
[21] d20diaries, Starter Set review — https://d20diaries.com/2019/02/26/dd-starter-set/
[22] D&D Beyond, Choosing Your First Adventure — https://www.dndbeyond.com/posts/721-new-players-guide-choosing-your-first-adventure
[23] Delian Tomb — https://dungeonsanddragonsfan.com/best-free-dnd-one-shots/ ; https://mcdm.fandom.com/wiki/Running_the_Game?oldid=1021
[24] WPN, Heroes of the Borderlands — https://wpn.wizards.com/en/products/starter-set-heroes-of-the-borderlands
[25] dnd5eapi.co, 2014 Exhaustion — https://www.dnd5eapi.co/api/2014/conditions/exhaustion
