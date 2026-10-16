# Roster — ตัวละครสำเร็จรูป 10 คลาส

ข้อมูลตัวละคร 10 ตัว (1 ตัวต่อ 1 คลาส) สำหรับ feature [`../features/character-roster.md`](../features/character-roster.md).
เซ็ต `core` (ฮีโร่ชุดแรก / First Company) — ทุกตัวมี `"set": "core"`; อีกเซ็ตอยู่ใน [`roster-legend.md`](roster-legend.md).
Backend แปลง JSON ท้ายไฟล์เป็น Go (`internal/rules/roster_default.go`); `art` = `id` (รูปใน `web/src/assets/classes/`).
ชื่อ เนื้อเรื่อง สกิล และไอเทมทั้งหมดแต่งขึ้นเอง ไม่ได้คัดลอกจากหนังสือหรือเกมใด

## กติกาสั้นๆ สำหรับ DM (ใช้ลูกเต๋าจริง)
1. **Check**: ทอย d20 + stat ที่เกี่ยวข้อง แล้วเทียบกับ DC ที่ DM ตั้ง (ง่าย 10 · ปานกลาง 15 · ยาก 20). STR = ยก ผลัก ฟัน · END = ทนพิษ ทนหนาว วิ่งไกล · AGI = หลบ แอบ ยิง ปีน
2. **โจมตี**: ทอย d20 + STR (ระยะประชิด) หรือ + AGI (ระยะไกล/อาวุธเบา) ให้ได้ ≥ ค่าป้องกันศัตรูที่ DM กำหนด (ปกติ 12). โดนแล้วทอยดาเมจ: อาวุธเบา d6, อาวุธใหญ่ d8
3. **HP (เลือด)**: ลดตามดาเมจที่โดน; HP 0 = หมดสติ ล้มลง ยังไม่ตาย จนกว่าจะมีคนรักษา. DM ปรับ HP ด้วยปุ่ม −/+ ในการ์ดตัวละคร
4. **มานา**: กด "ใช้" สกิล ระบบหักมานาให้เอง มานาไม่พอกดไม่ได้. DM กด "พักผ่อน" = มานาเต็มทุกคน (HP ไม่ฟื้นเอง)
5. **สกิล**: ทำตามคำอธิบายได้ทันที. สกิลที่เขียนว่า "ถ้าโดน" ต้องทอยโจมตีตามข้อ 2 ก่อน; ที่เขียนว่า "โดนแน่นอน/ไม่ต้องทอยโดน" ไม่ต้องทอย
6. **ตา (turn)** = ช่วงที่ตัวละคร 1 ตัวได้เล่น (เดิน + ทำ 1 อย่าง). **รอบ (round)** = ทุกตัวเล่นครบคนละ 1 ตา. ระยะนับเป็นช่องบน grid
7. **ไอเทม**: ของกิน/ยา/ของใช้ DM ตัดสินผลเอง (เช่น ยาสมุนไพร ฟื้น d4 HP) แล้วลดจำนวนในกระเป๋า. ชิ้นแรกของแต่ละตัว = ของประจำตัว (signature)
8. DM คือผู้ตัดสินสุดท้ายเสมอ ถ้าคำอธิบายไม่ครอบคลุม ให้ตัดสินตามความสนุกของโต๊ะ

## สูตรสมดุล stat
`total = hpMax/4 + str + end + agi + mpMax/2` → ทุกตัว = **26** (เป้า 26 ± 1).
หาร HP ด้วย 4 และมานาด้วย 2 เพื่อให้ทุกค่าอยู่สเกล 1–10 เท่ากัน. Top-2 stat (หลังแปลงสเกลแล้ว) ไม่ซ้ำกันเลยทั้ง 10 ตัว.

## ตารางสรุป
| # | ฮีโร่ | คลาส | role | HP | STR | END | AGI | MP | total | top-2 |
|---|---|---|---|---|---|---|---|---|---|---|
| 1 | Aldric Vane (อัลดริก เวน) | Knight (อัศวิน) | tank | 36 | 5 | 8 | 2 | 4 | 26 | HP+END |
| 2 | Soraya Quill (โซรายา ควิลล์) | Wizard (จอมเวท) | control | 16 | 2 | 3 | 7 | 20 | 26 | MP+AGI |
| 3 | Zarif Kesh (ซารีฟ เคช) | Rogue (โจร) | damage | 18 | 6 | 3 | 9 | 7 | 26 | AGI+STR |
| 4 | Sigrun Ashvale (ซีกรุน แอชเวล) | Ranger (พรานป่า) | damage | 22 | 4 | 6 | 8 | 5 | 26 | AGI+END |
| 5 | Imara Solenne (อิมารา โซเลนน์) | Cleric (นักบวช) | healer | 28 | 4 | 5 | 2 | 16 | 26 | MP+HP |
| 6 | Brenna Dawnmere (เบรนนา ดอว์นเมียร์) | Paladin (พาลาดิน) | hybrid | 24 | 7 | 7 | 2 | 8 | 26 | STR+END |
| 7 | Batu Redmane (บาตู เรดเมน) | Barbarian (นักรบเถื่อน) | damage | 32 | 9 | 4 | 2 | 6 | 26 | STR+HP |
| 8 | Cosimo Lark (โคซิโม ลาร์ก) | Bard (กวีพเนจร) | support | 24 | 3 | 4 | 8 | 10 | 26 | AGI+HP |
| 9 | Malai of the Deepwood (มาลัย แห่งพงไพร) | Druid (ดรูอิด) | hybrid | 20 | 7 | 4 | 3 | 14 | 26 | STR+MP |
| 10 | Ivo Grimsel (อีโว กริมเซล) | Necromancer (เนโครแมนเซอร์) | control | 16 | 3 | 7 | 3 | 18 | 26 | MP+END |

## 1. อัศวิน / Knight — `knight`
**อัลดริก เวน** / **Aldric Vane** · role: `tank` (แทงก์)

> อัศวินที่เคยพลาดจนเสียเพื่อนไป จึงสาบานว่าจะไม่ยอมให้ใครล้มก่อนตัวเองอีก  
> *A knight who once lost a comrade and swore no friend would ever fall before him again.*

**Stat**: HP 36 · STR 5 · END 8 · AGI 2 · มานา 4 — total 9 + 5 + 8 + 2 + 2 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `stand-guard` | ยืนบังให้ / Stand Guard | 1 | เลือกเพื่อน 1 คนที่อยู่ติดกัน จนถึงตาถัดไปของเรา การโจมตีที่เล็งเพื่อนคนนั้นจะมาโดนเราแทน<br>*Pick an ally next to you. Until your next turn, any attack aimed at them hits you instead.* |
| `shield-bash` | กระแทกโล่ / Shield Bash | 2 | เอาโล่กระแทกศัตรูที่อยู่ติดกัน ถ้าโดน ศัตรูกระเด็นถอยไป 2 ช่อง และเดินไม่ได้ในตาถัดไปของมัน<br>*Slam an adjacent enemy with your shield. On a hit it is knocked back 2 squares and can't move on its next turn.* |
| `last-stand` | ยืนหยัดครั้งสุดท้าย / Last Stand | 4 | จนจบรอบถัดไป เลือดของเราจะไม่ลดต่ำกว่า 1 ไม่ว่าจะโดนหนักแค่ไหน<br>*Until the end of the next round, your HP can't drop below 1, no matter how hard you are hit.* |

**ของติดตัว**

- โล่ตราสิงห์ตระกูลเวน / Vane Family Lion Shield ×1 — ของประจำตัว (signature)
- ดาบยาวอัศวิน / Knight's Longsword ×1
- หินลับดาบกับน้ำมันขัด / Whetstone and Blade Oil ×1
- ม้วนผ้าพันแผลสนาม / Field Bandage Roll ×3

## 2. จอมเวท / Wizard — `wizard`
**โซรายา ควิลล์** / **Soraya Quill** · role: `control` (ควบคุม)

> นักปราชญ์จากหอดูดาวกลางทะเลทราย ผู้อ่านดวงดาวได้คล่องเหมือนอ่านแผนที่สนามรบ  
> *A scholar from the desert observatories who reads the stars like a battle map.*

**Stat**: HP 16 · STR 2 · END 3 · AGI 7 · มานา 20 — total 4 + 2 + 3 + 7 + 10 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `arcane-bolt` | ลูกศรเวท / Arcane Bolt | 1 | ยิงลำแสงเวทใส่ศัตรู 1 ตัวในระยะ 6 ช่อง ถ้าโดน เสียหาย d6<br>*Fire a bolt of light at one enemy up to 6 squares away. On a hit it takes d6 damage.* |
| `frozen-moment` | หยุดกาล / Frozen Moment | 3 | ตรึงเวลาศัตรู 1 ตัวในระยะ 6 ช่อง มันต้องข้ามตาถัดไปทั้งตา (บอสอาจต้านได้ DM ตัดสิน)<br>*Freeze one enemy within 6 squares in time: it skips its next turn entirely (a boss may resist, DM decides).* |
| `starfall` | ดาวตก / Starfall | 5 | เรียกดาวตกลงพื้นที่ 3×3 ช่องในระยะ 8 ช่อง ทุกตัวในพื้นที่ ทั้งศัตรูและเพื่อน เสียหาย 2d6 โดยไม่ต้องทอยโดน<br>*Call falling stars onto a 3×3 area within 8 squares. Everyone inside, friend or foe, takes 2d6 damage, no hit roll needed.* |

**ของติดตัว**

- ไม้เท้าแก้วดารา / Starglass Staff ×1 — ของประจำตัว (signature)
- ตำราเวทปกหนังสีคราม / Indigo-Bound Spellbook ×1
- ขวดหมึกกับปากกาขนนกฮูก / Ink Pot and Owl-Feather Quill ×1
- ขวดน้ำยาจิตใส / Clear-Mind Tonic ×2

## 3. โจร / Rogue — `rogue`
**ซารีฟ เคช** / **Zarif Kesh** · role: `damage` (ดาเมจ)

> หัวขโมยเจ้าเสน่ห์จากเมืองคลอง ที่ติดหนี้คนไว้ครึ่งเมือง  
> *A charming thief from the canal city who owes money to half the town.*

**Stat**: HP 18 · STR 6 · END 3 · AGI 9 · มานา 7 — total 4.5 + 6 + 3 + 9 + 3.5 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `quick-fingers` | มือไว / Quick Fingers | 1 | สะเดาะกุญแจ หรือขโมยของชิ้นเล็ก 1 ชิ้นจากคนที่อยู่ติดกันได้เลยโดยไม่ต้องทอย ยกเว้นเขากำลังจ้องเราอยู่<br>*Pick a lock, or steal one small item from someone next to you, with no roll, unless they are watching you closely.* |
| `shadow-step` | หลบเข้าเงา / Shadow Step | 2 | กลืนหายไปในเงา ศัตรูเล็งเราไม่ได้จนกว่าเราจะโจมตีหรือจบตาถัดไป การโจมตีแรกจากที่ซ่อนแรงขึ้น +d6<br>*Melt into the shadows: enemies can't target you until you attack or your next turn ends. Your first attack from hiding deals +d6 damage.* |
| `venom-edge` | คมพิษ / Venom Edge | 4 | ทาพิษลงใบมีด ศัตรูตัวถัดไปที่เราฟันโดนจะเสียเลือด 2 แต้มทุกต้นตาของมัน นาน 3 ตา<br>*Coat your blade in poison. The next enemy you hit loses 2 HP at the start of each of its turns for 3 turns.* |

**ของติดตัว**

- มีดสั้นคู่เงินจันทรา / Twin Moonsilver Daggers ×1 — ของประจำตัว (signature)
- ชุดสะเดาะกุญแจ / Lockpick Roll ×1
- ลูกระเบิดควัน / Smoke Pellet ×3
- ตะขอเกี่ยวกับเชือกไหม / Grappling Hook and Silk Rope ×1
- ลูกเต๋าถ่วงน้ำหนัก / Loaded Dice ×2

## 4. พรานป่า / Ranger — `ranger`
**ซีกรุน แอชเวล** / **Sigrun Ashvale** · role: `damage` (ดาเมจ)

> นักล่าเงียบขรึมจากป่าสนทางเหนือ ที่ไว้ใจเหยี่ยวคู่ใจมากกว่าคน  
> *A quiet huntress from the northern pines who trusts her falcon more than people.*

**Stat**: HP 22 · STR 4 · END 6 · AGI 8 · มานา 5 — total 5.5 + 4 + 6 + 8 + 2.5 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `mark-quarry` | หมายตาเหยื่อ / Mark the Quarry | 1 | หมายตาศัตรู 1 ตัวที่มองเห็น จนกว่ามันจะล้มหรือเราพักผ่อน การโจมตีของเราใส่มันแรงขึ้น +2<br>*Mark one enemy you can see. Until it falls or you rest, your attacks against it deal +2 damage.* |
| `falcon-scout` | เหยี่ยวสอดแนม / Falcon Scout | 2 | ปล่อยเหยี่ยวบินล่วงหน้าไกลสูงสุด 20 ช่อง DM จะบอกว่ามันเห็นอะไร เช่น ศัตรู กับดัก หรือทางออก<br>*Send your falcon up to 20 squares ahead. The DM tells you what it sees: enemies, traps, exits.* |
| `deadeye-shot` | เล็งสังหาร / Deadeye Shot | 4 | ใช้ทั้งตายืนเล็งโดยไม่เดิน ตาถัดไปลูกธนูจะโดนแน่นอน เสียหาย 2d8<br>*Spend your whole turn aiming without moving. On your next turn the arrow hits automatically for 2d8 damage.* |

**ของติดตัว**

- ธนูยาวไม้แอช / Ashwood Longbow ×1 — ของประจำตัว (signature)
- ลูกธนูขนนกสีเทา / Grey-Fletched Arrow ×20
- นกหวีดกระดูกเรียกเหยี่ยว / Bone Falcon Whistle ×1
- เนื้อกวางรมควัน / Smoked Venison Strip ×5

## 5. นักบวช / Cleric — `cleric`
**อิมารา โซเลนน์** / **Imara Solenne** · role: `healer` (ฮีลเลอร์)

> หมอแห่งวิหารสุริยันที่หัวเราะเสียงดัง สวดเสียงเบา และไม่เคยทิ้งใครไว้ข้างหลัง  
> *A sun-temple healer who laughs loudly, prays quietly, and never leaves anyone behind.*

**Stat**: HP 28 · STR 4 · END 5 · AGI 2 · มานา 16 — total 7 + 4 + 5 + 2 + 8 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `purify` | ชำระล้าง / Purify | 1 | แตะตัวเพื่อน 1 คน เพื่อล้างพิษ คำสาป หรืออาการหวาดกลัวออก 1 อย่าง<br>*Touch an ally to remove one poison, curse, or fear effect.* |
| `healing-touch` | สัมผัสเยียวยา / Healing Touch | 2 | วางมือบนเพื่อนที่อยู่ติดกัน เพื่อนฟื้นเลือด d6 + 2<br>*Lay a hand on an ally next to you: they regain d6 + 2 HP.* |
| `dawn-return` | คืนชีพยามรุ่ง / Return at Dawn | 5 | ปลุกเพื่อนที่อยู่ติดกันและเลือดเหลือ 0 ให้ลุกขึ้นมาพร้อมเลือดครึ่งหนึ่งของเลือดสูงสุด<br>*Bring an adjacent ally at 0 HP back to their feet with half their max HP.* |

**ของติดตัว**

- จี้สุริยันศักดิ์สิทธิ์ / Sunburst Holy Pendant ×1 — ของประจำตัว (signature)
- กระบองเหล็ก / Iron Mace ×1
- ถุงสมุนไพรห้ามเลือด / Styptic Herb Pouch ×4
- ขวดน้ำมนต์ / Flask of Holy Water ×2
- สร้อยลูกประคำ / Prayer Beads ×1

## 6. พาลาดิน / Paladin — `paladin`
**เบรนนา ดอว์นเมียร์** / **Brenna Dawnmere** · role: `hybrid` (ไฮบริด)

> พาลาดินพูดตรงใจใหญ่ ผู้สาบานจะปกป้องคนอ่อนแอ และถือคำสาบานนั้นแบบตรงตัวทุกคำ  
> *A blunt, big-hearted paladin sworn to protect the weak, who takes that oath word for word.*

**Stat**: HP 24 · STR 7 · END 7 · AGI 2 · มานา 8 — total 6 + 7 + 7 + 2 + 4 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `life-bond` | สายใยชีวิต / Life Bond | 1 | ถ่ายเลือดของตัวเองให้เพื่อนที่เราแตะอยู่ ได้สูงสุด 5 แต้ม<br>*Give up to 5 of your own HP to an ally you are touching.* |
| `radiant-smite` | ฟาดแสงศักดิ์สิทธิ์ / Radiant Smite | 2 | การฟันโดนครั้งถัดไปแรงขึ้น +d8 และแรงเป็น 2 เท่าถ้าเป็นซากศพเดินได้หรือปีศาจ<br>*Your next melee hit deals +d8 damage, doubled against undead and demons.* |
| `holy-aegis` | โล่แสงศักดิ์สิทธิ์ / Holy Aegis | 4 | สร้างโล่แสงคลุมเพื่อน 1 คนในระยะ 4 ช่อง ทุกครั้งที่เพื่อนโดนตี ความเสียหายลดลง 5 แต้ม จนถึงตาถัดไปของเรา<br>*Wrap one ally within 4 squares in a shield of light: every hit on them deals 5 less damage until your next turn.* |

**ของติดตัว**

- ดาบใหญ่แห่งคำปฏิญาณ / Oathbound Greatsword ×1 — ของประจำตัว (signature)
- เกราะอกขลิบเงิน / Silver-Trimmed Breastplate ×1
- ม้วนคำปฏิญาณผนึกครั่ง / Wax-Sealed Oath Scroll ×1
- ชาดอกตะวัน / Sunpetal Tea ×2

## 7. นักรบเถื่อน / Barbarian — `barbarian`
**บาตู เรดเมน** / **Batu Redmane** · role: `damage` (ดาเมจ)

> นักรบจากทุ่งหญ้ากว้างที่ออกจากเผ่ามาตามหาศึกที่คู่ควรจะถูกแต่งเป็นเพลง  
> *A steppe warrior who left his clan to find a fight worthy of a song.*

**Stat**: HP 32 · STR 9 · END 4 · AGI 2 · มานา 6 — total 8 + 9 + 4 + 2 + 3 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `war-cry` | คำรามศึก / War Cry | 1 | คำรามจนขนลุก ศัตรูอ่อนแอในระยะ 3 ช่องจะหนีหรือลังเลไป 1 ตา (DM ตัดสินว่าตัวไหนนับว่าอ่อนแอ)<br>*Let out a terrifying roar: weak enemies within 3 squares flee or hesitate for 1 turn (the DM decides who counts as weak).* |
| `blood-rage` | โหมโทสะ / Blood Rage | 2 | นาน 3 ตา ทุกครั้งที่ฟันโดนแรงขึ้น +3 แต่ทอยที่ใช้ความว่องไว (AGI) −2<br>*For 3 turns your melee hits deal +3 damage, but your Agility rolls get −2.* |
| `unstoppable-charge` | พุ่งชนไม่ยั้ง / Unstoppable Charge | 4 | พุ่งเป็นเส้นตรงไกลสูงสุด 6 ช่อง ไม่มีอะไรขวางได้ แล้วฟันศัตรูที่ปลายทาง ถ้าโดนแรงขึ้น +d8<br>*Charge up to 6 squares in a straight line, nothing can block you, then strike the enemy at the end. On a hit it deals +d8 damage.* |

**ของติดตัว**

- ขวานยักษ์ตระกูลเรดเมน / Redmane Greataxe ×1 — ของประจำตัว (signature)
- เสื้อคลุมหนังหมาป่า / Wolf-Pelt Cloak ×1
- ขวานเล็กสำหรับขว้าง / Throwing Hatchet ×3
- ถุงนมม้าหมัก / Skin of Fermented Mare's Milk ×1

## 8. กวีพเนจร / Bard — `bard`
**โคซิโม ลาร์ก** / **Cosimo Lark** · role: `support` (สนับสนุน)

> นักแต่งเพลงพเนจรที่สาบานว่าทุกเรื่องที่เล่าเป็นเรื่องจริง โดยเฉพาะเรื่องของตัวเอง  
> *A wandering songwriter who swears every tale he tells is true, especially the ones about himself.*

**Stat**: HP 24 · STR 3 · END 4 · AGI 8 · มานา 10 — total 6 + 3 + 4 + 8 + 5 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `silver-tongue` | ปากหวาน / Silver Tongue | 1 | ทอยครั้งถัดไปเพื่อเกลี้ยกล่อม เอาใจ หรือหลอกใครสักคน ได้ +3<br>*Get +3 on your next roll to persuade, charm, or fool someone.* |
| `rousing-tune` | เพลงปลุกใจ / Rousing Tune | 2 | บรรเลงเพลงฮึกเหิม เพื่อนที่ได้ยินสูงสุด 2 คนได้ +2 ในการทอยครั้งถัดไป<br>*Play a stirring tune: up to 2 allies who can hear it get +2 on their next roll.* |
| `encore` | อีกสักรอบ! / Encore! | 5 | เลือกเพื่อน 1 คนที่เพิ่งเล่นจบตา เพื่อนได้เล่นอีก 1 ตาเต็มๆ ทันที<br>*Pick an ally who just finished their turn: they immediately take another full turn.* |

**ของติดตัว**

- พิณลูทสายเงิน / Silver-Stringed Lute ×1 — ของประจำตัว (signature)
- ดาบเรเพียร์ / Rapier ×1
- สำรับไพ่ทำเครื่องหมาย / Marked Deck of Cards ×1
- หมวกขนนกสุดหรู / Plumed Hat ×1
- ขวดไวน์น้ำผึ้ง / Bottle of Honeyed Wine ×2

## 9. ดรูอิด / Druid — `druid`
**มาลัย แห่งพงไพร** / **Malai of the Deepwood** · role: `hybrid` (ไฮบริด)

> ผู้พิทักษ์ป่าลึกที่พูดเสียงนุ่ม แต่กลายร่างเป็นหมีได้เมื่อจำเป็น  
> *A guardian of the deep jungle who speaks softly and becomes a bear when she must.*

**Stat**: HP 20 · STR 7 · END 4 · AGI 3 · มานา 14 — total 5 + 7 + 4 + 3 + 7 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `thorn-lash` | แส้หนาม / Thorn Lash | 1 | ฟาดแส้เถาหนามใส่ศัตรูในระยะ 3 ช่อง ถ้าโดน เสียหาย d4 และถูกดึงเข้ามาหาเรา 2 ช่อง<br>*Lash an enemy within 3 squares with a thorny vine. On a hit it takes d4 damage and is pulled 2 squares toward you.* |
| `entangling-roots` | รากไม้พันธนาการ / Entangling Roots | 2 | รากไม้พุ่งขึ้นจากพื้นรัดศัตรู 1 ตัวในระยะ 5 ช่อง มันเดินไม่ได้ 2 ตา (แต่ยังโจมตีได้)<br>*Roots burst up around one enemy within 5 squares: it can't move for 2 turns (it can still attack).* |
| `bear-form` | ร่างหมีป่า / Bear Form | 4 | แปลงร่างเป็นหมีป่าตัวใหญ่นาน 3 ตา ได้เลือดเสริม 10 แต้ม (หายเมื่อคืนร่าง) และนับ STR เป็น 10 ระหว่างเป็นหมีใช้สกิลอื่นไม่ได้<br>*Become a great bear for 3 turns: gain 10 extra HP (gone when you change back) and your Strength counts as 10. You can't use other skills as a bear.* |

**ของติดตัว**

- ไม้เท้าเถาวัลย์มีชีวิต / Living Vine Staff ×1 — ของประจำตัว (signature)
- ถุงเมล็ดพันธุ์ป่า / Pouch of Wild Seeds ×1
- เครื่องรางด้วงหยก / Jade Beetle Charm ×1
- มัดตะไคร่จันทร์ / Bundle of Moonmoss ×3

## 10. เนโครแมนเซอร์ / Necromancer — `necromancer`
**อีโว กริมเซล** / **Ivo Grimsel** · role: `control` (ควบคุม)

> ลูกชายคนเฝ้าสุสาน ที่ฟังเสียงคนตายเพราะคนเป็นไม่เคยฟังเขา  
> *A gravekeeper's son who listens to the dead because the living never listened to him.*

**Stat**: HP 16 · STR 3 · END 7 · AGI 3 · มานา 18 — total 4 + 3 + 7 + 3 + 9 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `grave-whisper` | กระซิบจากหลุม / Grave Whisper | 1 | ถามศพ 1 คำถาม ศพจะตอบสั้นๆ ตามที่มันรู้ตอนยังมีชีวิต<br>*Ask a corpse one question. It answers in a few words with what it knew in life.* |
| `life-siphon` | ดูดชีวิต / Life Siphon | 2 | ดูดพลังชีวิตศัตรูในระยะ 4 ช่อง ถ้าโดน มันเสียหาย d6 และเราฟื้นเลือดเท่ากับที่ดูดได้<br>*Drain an enemy within 4 squares. On a hit it takes d6 damage and you regain that much HP.* |
| `raise-fallen` | ปลุกซากผู้ล่วงลับ / Raise the Fallen | 5 | ปลุกศพในระยะ 3 ช่องเป็นโครงกระดูกรับใช้ (เลือด 6 ตีแรง d4) ทำตามคำสั่งเราจนกว่าจะถูกทำลายหรือจบฉาก<br>*Raise a corpse within 3 squares as a skeleton servant (6 HP, hits for d4) that obeys you until destroyed or the scene ends.* |

**ของติดตัว**

- ตะเกียงกะโหลก / Skull Lantern ×1 — ของประจำตัว (signature)
- เคียวด้ามกระดูก / Bone-Handled Sickle ×1
- โถดินจากหลุมศพ / Jar of Grave Dirt ×1
- สมุดรายชื่อผู้วายชนม์ / Ledger of the Departed ×1
- เสื้อคลุมขนอีกา / Crow-Feather Cloak ×1

## JSON (machine-readable)

```json
[
  {
    "id": "knight",
    "set": "core",
    "heroName": {
      "th": "อัลดริก เวน",
      "en": "Aldric Vane"
    },
    "className": {
      "th": "อัศวิน",
      "en": "Knight"
    },
    "blurb": {
      "th": "อัศวินที่เคยพลาดจนเสียเพื่อนไป จึงสาบานว่าจะไม่ยอมให้ใครล้มก่อนตัวเองอีก",
      "en": "A knight who once lost a comrade and swore no friend would ever fall before him again."
    },
    "role": "tank",
    "stats": {
      "hpMax": 36,
      "str": 5,
      "end": 8,
      "agi": 2,
      "mpMax": 4
    },
    "skills": [
      {
        "id": "stand-guard",
        "name": {
          "th": "ยืนบังให้",
          "en": "Stand Guard"
        },
        "desc": {
          "th": "เลือกเพื่อน 1 คนที่อยู่ติดกัน จนถึงตาถัดไปของเรา การโจมตีที่เล็งเพื่อนคนนั้นจะมาโดนเราแทน",
          "en": "Pick an ally next to you. Until your next turn, any attack aimed at them hits you instead."
        },
        "cost": 1
      },
      {
        "id": "shield-bash",
        "name": {
          "th": "กระแทกโล่",
          "en": "Shield Bash"
        },
        "desc": {
          "th": "เอาโล่กระแทกศัตรูที่อยู่ติดกัน ถ้าโดน ศัตรูกระเด็นถอยไป 2 ช่อง และเดินไม่ได้ในตาถัดไปของมัน",
          "en": "Slam an adjacent enemy with your shield. On a hit it is knocked back 2 squares and can't move on its next turn."
        },
        "cost": 2
      },
      {
        "id": "last-stand",
        "name": {
          "th": "ยืนหยัดครั้งสุดท้าย",
          "en": "Last Stand"
        },
        "desc": {
          "th": "จนจบรอบถัดไป เลือดของเราจะไม่ลดต่ำกว่า 1 ไม่ว่าจะโดนหนักแค่ไหน",
          "en": "Until the end of the next round, your HP can't drop below 1, no matter how hard you are hit."
        },
        "cost": 4
      }
    ],
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
          "th": "ดาบยาวอัศวิน",
          "en": "Knight's Longsword"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "หินลับดาบกับน้ำมันขัด",
          "en": "Whetstone and Blade Oil"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ม้วนผ้าพันแผลสนาม",
          "en": "Field Bandage Roll"
        },
        "qty": 3
      }
    ]
  },
  {
    "id": "wizard",
    "set": "core",
    "heroName": {
      "th": "โซรายา ควิลล์",
      "en": "Soraya Quill"
    },
    "className": {
      "th": "จอมเวท",
      "en": "Wizard"
    },
    "blurb": {
      "th": "นักปราชญ์จากหอดูดาวกลางทะเลทราย ผู้อ่านดวงดาวได้คล่องเหมือนอ่านแผนที่สนามรบ",
      "en": "A scholar from the desert observatories who reads the stars like a battle map."
    },
    "role": "control",
    "stats": {
      "hpMax": 16,
      "str": 2,
      "end": 3,
      "agi": 7,
      "mpMax": 20
    },
    "skills": [
      {
        "id": "arcane-bolt",
        "name": {
          "th": "ลูกศรเวท",
          "en": "Arcane Bolt"
        },
        "desc": {
          "th": "ยิงลำแสงเวทใส่ศัตรู 1 ตัวในระยะ 6 ช่อง ถ้าโดน เสียหาย d6",
          "en": "Fire a bolt of light at one enemy up to 6 squares away. On a hit it takes d6 damage."
        },
        "cost": 1
      },
      {
        "id": "frozen-moment",
        "name": {
          "th": "หยุดกาล",
          "en": "Frozen Moment"
        },
        "desc": {
          "th": "ตรึงเวลาศัตรู 1 ตัวในระยะ 6 ช่อง มันต้องข้ามตาถัดไปทั้งตา (บอสอาจต้านได้ DM ตัดสิน)",
          "en": "Freeze one enemy within 6 squares in time: it skips its next turn entirely (a boss may resist, DM decides)."
        },
        "cost": 3
      },
      {
        "id": "starfall",
        "name": {
          "th": "ดาวตก",
          "en": "Starfall"
        },
        "desc": {
          "th": "เรียกดาวตกลงพื้นที่ 3×3 ช่องในระยะ 8 ช่อง ทุกตัวในพื้นที่ ทั้งศัตรูและเพื่อน เสียหาย 2d6 โดยไม่ต้องทอยโดน",
          "en": "Call falling stars onto a 3×3 area within 8 squares. Everyone inside, friend or foe, takes 2d6 damage, no hit roll needed."
        },
        "cost": 5
      }
    ],
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
          "th": "ขวดหมึกกับปากกาขนนกฮูก",
          "en": "Ink Pot and Owl-Feather Quill"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขวดน้ำยาจิตใส",
          "en": "Clear-Mind Tonic"
        },
        "qty": 2
      }
    ]
  },
  {
    "id": "rogue",
    "set": "core",
    "heroName": {
      "th": "ซารีฟ เคช",
      "en": "Zarif Kesh"
    },
    "className": {
      "th": "โจร",
      "en": "Rogue"
    },
    "blurb": {
      "th": "หัวขโมยเจ้าเสน่ห์จากเมืองคลอง ที่ติดหนี้คนไว้ครึ่งเมือง",
      "en": "A charming thief from the canal city who owes money to half the town."
    },
    "role": "damage",
    "stats": {
      "hpMax": 18,
      "str": 6,
      "end": 3,
      "agi": 9,
      "mpMax": 7
    },
    "skills": [
      {
        "id": "quick-fingers",
        "name": {
          "th": "มือไว",
          "en": "Quick Fingers"
        },
        "desc": {
          "th": "สะเดาะกุญแจ หรือขโมยของชิ้นเล็ก 1 ชิ้นจากคนที่อยู่ติดกันได้เลยโดยไม่ต้องทอย ยกเว้นเขากำลังจ้องเราอยู่",
          "en": "Pick a lock, or steal one small item from someone next to you, with no roll, unless they are watching you closely."
        },
        "cost": 1
      },
      {
        "id": "shadow-step",
        "name": {
          "th": "หลบเข้าเงา",
          "en": "Shadow Step"
        },
        "desc": {
          "th": "กลืนหายไปในเงา ศัตรูเล็งเราไม่ได้จนกว่าเราจะโจมตีหรือจบตาถัดไป การโจมตีแรกจากที่ซ่อนแรงขึ้น +d6",
          "en": "Melt into the shadows: enemies can't target you until you attack or your next turn ends. Your first attack from hiding deals +d6 damage."
        },
        "cost": 2
      },
      {
        "id": "venom-edge",
        "name": {
          "th": "คมพิษ",
          "en": "Venom Edge"
        },
        "desc": {
          "th": "ทาพิษลงใบมีด ศัตรูตัวถัดไปที่เราฟันโดนจะเสียเลือด 2 แต้มทุกต้นตาของมัน นาน 3 ตา",
          "en": "Coat your blade in poison. The next enemy you hit loses 2 HP at the start of each of its turns for 3 turns."
        },
        "cost": 4
      }
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
          "th": "ชุดสะเดาะกุญแจ",
          "en": "Lockpick Roll"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ลูกระเบิดควัน",
          "en": "Smoke Pellet"
        },
        "qty": 3
      },
      {
        "name": {
          "th": "ตะขอเกี่ยวกับเชือกไหม",
          "en": "Grappling Hook and Silk Rope"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ลูกเต๋าถ่วงน้ำหนัก",
          "en": "Loaded Dice"
        },
        "qty": 2
      }
    ]
  },
  {
    "id": "ranger",
    "set": "core",
    "heroName": {
      "th": "ซีกรุน แอชเวล",
      "en": "Sigrun Ashvale"
    },
    "className": {
      "th": "พรานป่า",
      "en": "Ranger"
    },
    "blurb": {
      "th": "นักล่าเงียบขรึมจากป่าสนทางเหนือ ที่ไว้ใจเหยี่ยวคู่ใจมากกว่าคน",
      "en": "A quiet huntress from the northern pines who trusts her falcon more than people."
    },
    "role": "damage",
    "stats": {
      "hpMax": 22,
      "str": 4,
      "end": 6,
      "agi": 8,
      "mpMax": 5
    },
    "skills": [
      {
        "id": "mark-quarry",
        "name": {
          "th": "หมายตาเหยื่อ",
          "en": "Mark the Quarry"
        },
        "desc": {
          "th": "หมายตาศัตรู 1 ตัวที่มองเห็น จนกว่ามันจะล้มหรือเราพักผ่อน การโจมตีของเราใส่มันแรงขึ้น +2",
          "en": "Mark one enemy you can see. Until it falls or you rest, your attacks against it deal +2 damage."
        },
        "cost": 1
      },
      {
        "id": "falcon-scout",
        "name": {
          "th": "เหยี่ยวสอดแนม",
          "en": "Falcon Scout"
        },
        "desc": {
          "th": "ปล่อยเหยี่ยวบินล่วงหน้าไกลสูงสุด 20 ช่อง DM จะบอกว่ามันเห็นอะไร เช่น ศัตรู กับดัก หรือทางออก",
          "en": "Send your falcon up to 20 squares ahead. The DM tells you what it sees: enemies, traps, exits."
        },
        "cost": 2
      },
      {
        "id": "deadeye-shot",
        "name": {
          "th": "เล็งสังหาร",
          "en": "Deadeye Shot"
        },
        "desc": {
          "th": "ใช้ทั้งตายืนเล็งโดยไม่เดิน ตาถัดไปลูกธนูจะโดนแน่นอน เสียหาย 2d8",
          "en": "Spend your whole turn aiming without moving. On your next turn the arrow hits automatically for 2d8 damage."
        },
        "cost": 4
      }
    ],
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
          "th": "ลูกธนูขนนกสีเทา",
          "en": "Grey-Fletched Arrow"
        },
        "qty": 20
      },
      {
        "name": {
          "th": "นกหวีดกระดูกเรียกเหยี่ยว",
          "en": "Bone Falcon Whistle"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เนื้อกวางรมควัน",
          "en": "Smoked Venison Strip"
        },
        "qty": 5
      }
    ]
  },
  {
    "id": "cleric",
    "set": "core",
    "heroName": {
      "th": "อิมารา โซเลนน์",
      "en": "Imara Solenne"
    },
    "className": {
      "th": "นักบวช",
      "en": "Cleric"
    },
    "blurb": {
      "th": "หมอแห่งวิหารสุริยันที่หัวเราะเสียงดัง สวดเสียงเบา และไม่เคยทิ้งใครไว้ข้างหลัง",
      "en": "A sun-temple healer who laughs loudly, prays quietly, and never leaves anyone behind."
    },
    "role": "healer",
    "stats": {
      "hpMax": 28,
      "str": 4,
      "end": 5,
      "agi": 2,
      "mpMax": 16
    },
    "skills": [
      {
        "id": "purify",
        "name": {
          "th": "ชำระล้าง",
          "en": "Purify"
        },
        "desc": {
          "th": "แตะตัวเพื่อน 1 คน เพื่อล้างพิษ คำสาป หรืออาการหวาดกลัวออก 1 อย่าง",
          "en": "Touch an ally to remove one poison, curse, or fear effect."
        },
        "cost": 1
      },
      {
        "id": "healing-touch",
        "name": {
          "th": "สัมผัสเยียวยา",
          "en": "Healing Touch"
        },
        "desc": {
          "th": "วางมือบนเพื่อนที่อยู่ติดกัน เพื่อนฟื้นเลือด d6 + 2",
          "en": "Lay a hand on an ally next to you: they regain d6 + 2 HP."
        },
        "cost": 2
      },
      {
        "id": "dawn-return",
        "name": {
          "th": "คืนชีพยามรุ่ง",
          "en": "Return at Dawn"
        },
        "desc": {
          "th": "ปลุกเพื่อนที่อยู่ติดกันและเลือดเหลือ 0 ให้ลุกขึ้นมาพร้อมเลือดครึ่งหนึ่งของเลือดสูงสุด",
          "en": "Bring an adjacent ally at 0 HP back to their feet with half their max HP."
        },
        "cost": 5
      }
    ],
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
          "th": "กระบองเหล็ก",
          "en": "Iron Mace"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ถุงสมุนไพรห้ามเลือด",
          "en": "Styptic Herb Pouch"
        },
        "qty": 4
      },
      {
        "name": {
          "th": "ขวดน้ำมนต์",
          "en": "Flask of Holy Water"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "สร้อยลูกประคำ",
          "en": "Prayer Beads"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "paladin",
    "set": "core",
    "heroName": {
      "th": "เบรนนา ดอว์นเมียร์",
      "en": "Brenna Dawnmere"
    },
    "className": {
      "th": "พาลาดิน",
      "en": "Paladin"
    },
    "blurb": {
      "th": "พาลาดินพูดตรงใจใหญ่ ผู้สาบานจะปกป้องคนอ่อนแอ และถือคำสาบานนั้นแบบตรงตัวทุกคำ",
      "en": "A blunt, big-hearted paladin sworn to protect the weak, who takes that oath word for word."
    },
    "role": "hybrid",
    "stats": {
      "hpMax": 24,
      "str": 7,
      "end": 7,
      "agi": 2,
      "mpMax": 8
    },
    "skills": [
      {
        "id": "life-bond",
        "name": {
          "th": "สายใยชีวิต",
          "en": "Life Bond"
        },
        "desc": {
          "th": "ถ่ายเลือดของตัวเองให้เพื่อนที่เราแตะอยู่ ได้สูงสุด 5 แต้ม",
          "en": "Give up to 5 of your own HP to an ally you are touching."
        },
        "cost": 1
      },
      {
        "id": "radiant-smite",
        "name": {
          "th": "ฟาดแสงศักดิ์สิทธิ์",
          "en": "Radiant Smite"
        },
        "desc": {
          "th": "การฟันโดนครั้งถัดไปแรงขึ้น +d8 และแรงเป็น 2 เท่าถ้าเป็นซากศพเดินได้หรือปีศาจ",
          "en": "Your next melee hit deals +d8 damage, doubled against undead and demons."
        },
        "cost": 2
      },
      {
        "id": "holy-aegis",
        "name": {
          "th": "โล่แสงศักดิ์สิทธิ์",
          "en": "Holy Aegis"
        },
        "desc": {
          "th": "สร้างโล่แสงคลุมเพื่อน 1 คนในระยะ 4 ช่อง ทุกครั้งที่เพื่อนโดนตี ความเสียหายลดลง 5 แต้ม จนถึงตาถัดไปของเรา",
          "en": "Wrap one ally within 4 squares in a shield of light: every hit on them deals 5 less damage until your next turn."
        },
        "cost": 4
      }
    ],
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
          "th": "เกราะอกขลิบเงิน",
          "en": "Silver-Trimmed Breastplate"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ม้วนคำปฏิญาณผนึกครั่ง",
          "en": "Wax-Sealed Oath Scroll"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชาดอกตะวัน",
          "en": "Sunpetal Tea"
        },
        "qty": 2
      }
    ]
  },
  {
    "id": "barbarian",
    "set": "core",
    "heroName": {
      "th": "บาตู เรดเมน",
      "en": "Batu Redmane"
    },
    "className": {
      "th": "นักรบเถื่อน",
      "en": "Barbarian"
    },
    "blurb": {
      "th": "นักรบจากทุ่งหญ้ากว้างที่ออกจากเผ่ามาตามหาศึกที่คู่ควรจะถูกแต่งเป็นเพลง",
      "en": "A steppe warrior who left his clan to find a fight worthy of a song."
    },
    "role": "damage",
    "stats": {
      "hpMax": 32,
      "str": 9,
      "end": 4,
      "agi": 2,
      "mpMax": 6
    },
    "skills": [
      {
        "id": "war-cry",
        "name": {
          "th": "คำรามศึก",
          "en": "War Cry"
        },
        "desc": {
          "th": "คำรามจนขนลุก ศัตรูอ่อนแอในระยะ 3 ช่องจะหนีหรือลังเลไป 1 ตา (DM ตัดสินว่าตัวไหนนับว่าอ่อนแอ)",
          "en": "Let out a terrifying roar: weak enemies within 3 squares flee or hesitate for 1 turn (the DM decides who counts as weak)."
        },
        "cost": 1
      },
      {
        "id": "blood-rage",
        "name": {
          "th": "โหมโทสะ",
          "en": "Blood Rage"
        },
        "desc": {
          "th": "นาน 3 ตา ทุกครั้งที่ฟันโดนแรงขึ้น +3 แต่ทอยที่ใช้ความว่องไว (AGI) −2",
          "en": "For 3 turns your melee hits deal +3 damage, but your Agility rolls get −2."
        },
        "cost": 2
      },
      {
        "id": "unstoppable-charge",
        "name": {
          "th": "พุ่งชนไม่ยั้ง",
          "en": "Unstoppable Charge"
        },
        "desc": {
          "th": "พุ่งเป็นเส้นตรงไกลสูงสุด 6 ช่อง ไม่มีอะไรขวางได้ แล้วฟันศัตรูที่ปลายทาง ถ้าโดนแรงขึ้น +d8",
          "en": "Charge up to 6 squares in a straight line, nothing can block you, then strike the enemy at the end. On a hit it deals +d8 damage."
        },
        "cost": 4
      }
    ],
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
          "th": "เสื้อคลุมหนังหมาป่า",
          "en": "Wolf-Pelt Cloak"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขวานเล็กสำหรับขว้าง",
          "en": "Throwing Hatchet"
        },
        "qty": 3
      },
      {
        "name": {
          "th": "ถุงนมม้าหมัก",
          "en": "Skin of Fermented Mare's Milk"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "bard",
    "set": "core",
    "heroName": {
      "th": "โคซิโม ลาร์ก",
      "en": "Cosimo Lark"
    },
    "className": {
      "th": "กวีพเนจร",
      "en": "Bard"
    },
    "blurb": {
      "th": "นักแต่งเพลงพเนจรที่สาบานว่าทุกเรื่องที่เล่าเป็นเรื่องจริง โดยเฉพาะเรื่องของตัวเอง",
      "en": "A wandering songwriter who swears every tale he tells is true, especially the ones about himself."
    },
    "role": "support",
    "stats": {
      "hpMax": 24,
      "str": 3,
      "end": 4,
      "agi": 8,
      "mpMax": 10
    },
    "skills": [
      {
        "id": "silver-tongue",
        "name": {
          "th": "ปากหวาน",
          "en": "Silver Tongue"
        },
        "desc": {
          "th": "ทอยครั้งถัดไปเพื่อเกลี้ยกล่อม เอาใจ หรือหลอกใครสักคน ได้ +3",
          "en": "Get +3 on your next roll to persuade, charm, or fool someone."
        },
        "cost": 1
      },
      {
        "id": "rousing-tune",
        "name": {
          "th": "เพลงปลุกใจ",
          "en": "Rousing Tune"
        },
        "desc": {
          "th": "บรรเลงเพลงฮึกเหิม เพื่อนที่ได้ยินสูงสุด 2 คนได้ +2 ในการทอยครั้งถัดไป",
          "en": "Play a stirring tune: up to 2 allies who can hear it get +2 on their next roll."
        },
        "cost": 2
      },
      {
        "id": "encore",
        "name": {
          "th": "อีกสักรอบ!",
          "en": "Encore!"
        },
        "desc": {
          "th": "เลือกเพื่อน 1 คนที่เพิ่งเล่นจบตา เพื่อนได้เล่นอีก 1 ตาเต็มๆ ทันที",
          "en": "Pick an ally who just finished their turn: they immediately take another full turn."
        },
        "cost": 5
      }
    ],
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
          "th": "ดาบเรเพียร์",
          "en": "Rapier"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "สำรับไพ่ทำเครื่องหมาย",
          "en": "Marked Deck of Cards"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "หมวกขนนกสุดหรู",
          "en": "Plumed Hat"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขวดไวน์น้ำผึ้ง",
          "en": "Bottle of Honeyed Wine"
        },
        "qty": 2
      }
    ]
  },
  {
    "id": "druid",
    "set": "core",
    "heroName": {
      "th": "มาลัย แห่งพงไพร",
      "en": "Malai of the Deepwood"
    },
    "className": {
      "th": "ดรูอิด",
      "en": "Druid"
    },
    "blurb": {
      "th": "ผู้พิทักษ์ป่าลึกที่พูดเสียงนุ่ม แต่กลายร่างเป็นหมีได้เมื่อจำเป็น",
      "en": "A guardian of the deep jungle who speaks softly and becomes a bear when she must."
    },
    "role": "hybrid",
    "stats": {
      "hpMax": 20,
      "str": 7,
      "end": 4,
      "agi": 3,
      "mpMax": 14
    },
    "skills": [
      {
        "id": "thorn-lash",
        "name": {
          "th": "แส้หนาม",
          "en": "Thorn Lash"
        },
        "desc": {
          "th": "ฟาดแส้เถาหนามใส่ศัตรูในระยะ 3 ช่อง ถ้าโดน เสียหาย d4 และถูกดึงเข้ามาหาเรา 2 ช่อง",
          "en": "Lash an enemy within 3 squares with a thorny vine. On a hit it takes d4 damage and is pulled 2 squares toward you."
        },
        "cost": 1
      },
      {
        "id": "entangling-roots",
        "name": {
          "th": "รากไม้พันธนาการ",
          "en": "Entangling Roots"
        },
        "desc": {
          "th": "รากไม้พุ่งขึ้นจากพื้นรัดศัตรู 1 ตัวในระยะ 5 ช่อง มันเดินไม่ได้ 2 ตา (แต่ยังโจมตีได้)",
          "en": "Roots burst up around one enemy within 5 squares: it can't move for 2 turns (it can still attack)."
        },
        "cost": 2
      },
      {
        "id": "bear-form",
        "name": {
          "th": "ร่างหมีป่า",
          "en": "Bear Form"
        },
        "desc": {
          "th": "แปลงร่างเป็นหมีป่าตัวใหญ่นาน 3 ตา ได้เลือดเสริม 10 แต้ม (หายเมื่อคืนร่าง) และนับ STR เป็น 10 ระหว่างเป็นหมีใช้สกิลอื่นไม่ได้",
          "en": "Become a great bear for 3 turns: gain 10 extra HP (gone when you change back) and your Strength counts as 10. You can't use other skills as a bear."
        },
        "cost": 4
      }
    ],
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
          "th": "ถุงเมล็ดพันธุ์ป่า",
          "en": "Pouch of Wild Seeds"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เครื่องรางด้วงหยก",
          "en": "Jade Beetle Charm"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มัดตะไคร่จันทร์",
          "en": "Bundle of Moonmoss"
        },
        "qty": 3
      }
    ]
  },
  {
    "id": "necromancer",
    "set": "core",
    "heroName": {
      "th": "อีโว กริมเซล",
      "en": "Ivo Grimsel"
    },
    "className": {
      "th": "เนโครแมนเซอร์",
      "en": "Necromancer"
    },
    "blurb": {
      "th": "ลูกชายคนเฝ้าสุสาน ที่ฟังเสียงคนตายเพราะคนเป็นไม่เคยฟังเขา",
      "en": "A gravekeeper's son who listens to the dead because the living never listened to him."
    },
    "role": "control",
    "stats": {
      "hpMax": 16,
      "str": 3,
      "end": 7,
      "agi": 3,
      "mpMax": 18
    },
    "skills": [
      {
        "id": "grave-whisper",
        "name": {
          "th": "กระซิบจากหลุม",
          "en": "Grave Whisper"
        },
        "desc": {
          "th": "ถามศพ 1 คำถาม ศพจะตอบสั้นๆ ตามที่มันรู้ตอนยังมีชีวิต",
          "en": "Ask a corpse one question. It answers in a few words with what it knew in life."
        },
        "cost": 1
      },
      {
        "id": "life-siphon",
        "name": {
          "th": "ดูดชีวิต",
          "en": "Life Siphon"
        },
        "desc": {
          "th": "ดูดพลังชีวิตศัตรูในระยะ 4 ช่อง ถ้าโดน มันเสียหาย d6 และเราฟื้นเลือดเท่ากับที่ดูดได้",
          "en": "Drain an enemy within 4 squares. On a hit it takes d6 damage and you regain that much HP."
        },
        "cost": 2
      },
      {
        "id": "raise-fallen",
        "name": {
          "th": "ปลุกซากผู้ล่วงลับ",
          "en": "Raise the Fallen"
        },
        "desc": {
          "th": "ปลุกศพในระยะ 3 ช่องเป็นโครงกระดูกรับใช้ (เลือด 6 ตีแรง d4) ทำตามคำสั่งเราจนกว่าจะถูกทำลายหรือจบฉาก",
          "en": "Raise a corpse within 3 squares as a skeleton servant (6 HP, hits for d4) that obeys you until destroyed or the scene ends."
        },
        "cost": 5
      }
    ],
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
          "th": "เคียวด้ามกระดูก",
          "en": "Bone-Handled Sickle"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "โถดินจากหลุมศพ",
          "en": "Jar of Grave Dirt"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "สมุดรายชื่อผู้วายชนม์",
          "en": "Ledger of the Departed"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เสื้อคลุมขนอีกา",
          "en": "Crow-Feather Cloak"
        },
        "qty": 1
      }
    ]
  }
]
```
