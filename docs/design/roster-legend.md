# Roster — วีรชนในตำนาน / Heroes of Legend (12 คลาส)

ชุดตัวละครชุดที่ 2 (`set: "legend"`): 12 ตัว 1 ตัวต่อ 1 คลาส สำหรับ feature [`../features/character-roster.md`](../features/character-roster.md).
รูปแบบเดียวกับ [`roster.md`](roster.md): Backend อ่าน JSON ท้ายไฟล์; `id` = `legend-<คลาส>`, `art` = `""` (ใช้ไอคอนคลาสใน `web/src/assets/classes/`).
ชื่อ เนื้อเรื่อง สกิล และไอเทมทั้งหมดแต่งขึ้นเอง ไม่ได้คัดลอกจากหนังสือหรือเกมใด; ชื่อคลาสเป็นคำทั่วไป แนวคลาสอ้างอิง [`../research/classes/phb-2024/`](../research/classes/phb-2024/).

## กติกาสั้นๆ สำหรับ DM
ใช้กติกาชุดเดียวกับ [roster.md — กติกาสั้นๆ สำหรับ DM](roster.md#กติกาสั้นๆ-สำหรับ-dm-ใช้ลูกเต๋าจริง) (check, โจมตี, "ถ้าโดน", ตา/รอบ, ระยะเป็นช่อง, ไอเทม).

## สูตรสมดุล stat
`total = hpMax/4 + str + end + agi + mpMax/2` → ทุกตัว = **26** (เป้า 26 ± 1).
หาร HP ด้วย 4 และมานาด้วย 2 เพื่อให้ทุกค่าอยู่สเกล 1–10 เท่ากัน. Top-2 stat (หลังแปลงสเกลแล้ว) ไม่ซ้ำกันเลยทั้ง 12 ตัว
นับแบบมีลำดับ (อันดับ 1 + อันดับ 2; ค่าเท่ากันเรียงตาม HP, STR, END, AGI, MP แบบ roster.md). คู่ stat แบบไม่สนลำดับมีแค่ 10 คู่
จึงมี 2 คู่ที่กลับลำดับกัน: Bard AGI+MP กับ Wizard MP+AGI, Druid END+MP กับ Sorcerer MP+END.

## ตารางสรุป
| # | ฮีโร่ | คลาส | role | HP | STR | END | AGI | MP | total | top-2 |
|---|---|---|---|---|---|---|---|---|---|---|
| 1 | Kaya Stormhorn (คายา สตอร์มฮอร์น) | Barbarian (นักรบเถื่อน) | damage | 36 | 8 | 5 | 2 | 4 | 26 | HP+STR |
| 2 | Odalys Reed (โอดาลิส รีด) | Bard (นักกวีเวทดนตรี) | support | 20 | 4 | 4 | 7 | 12 | 26 | AGI+MP |
| 3 | Oswin Graylight (ออสวิน เกรย์ไลต์) | Cleric (นักบวช) | healer | 20 | 6 | 5 | 2 | 16 | 26 | MP+STR |
| 4 | Tansy Hollowfern (แทนซี ฮอลโลว์เฟิร์น) | Druid (ดรูอิด) | hybrid | 20 | 3 | 7 | 5 | 12 | 26 | END+MP |
| 5 | Petra Kalvos (เพตรา คาลวอส) | Fighter (นักรบ) | tank | 24 | 8 | 7 | 3 | 4 | 26 | STR+END |
| 6 | Kiran Stillwater (คิรัน สติลวอเตอร์) | Monk (นักพรตยุทธ์) | control | 16 | 4 | 6 | 9 | 6 | 26 | AGI+END |
| 7 | Tobias Vellan (โทเบียส เวลลัน) | Paladin (อัศวินศักดิ์สิทธิ์) | hybrid | 28 | 6 | 8 | 2 | 6 | 26 | END+HP |
| 8 | Joren Vask (โยเรน วาสก์) | Ranger (พรานป่า) | damage | 28 | 4 | 4 | 8 | 6 | 26 | AGI+HP |
| 9 | Silas Wren (ไซลัส เรน) | Rogue (นักย่องเบา) | damage | 16 | 6 | 3 | 10 | 6 | 26 | AGI+STR |
| 10 | Remy Saltfire (เรมี ซอลต์ไฟร์) | Sorcerer (นักเวทโดยกำเนิด) | damage | 16 | 2 | 6 | 5 | 18 | 26 | MP+END |
| 11 | Vesna Morrow (เวสนา มอร์โรว์) | Warlock (ผู้ทำสัญญา) | hybrid | 24 | 3 | 4 | 5 | 16 | 26 | MP+HP |
| 12 | Theodric Pell (ธีโอดริก เพลล์) | Wizard (พ่อมด) | control | 16 | 2 | 4 | 6 | 20 | 26 | MP+AGI |

## 1. นักรบเถื่อน / Barbarian — `legend-barbarian`
**คายา สตอร์มฮอร์น** / **Kaya Stormhorn** · role: `damage` (ดาเมจ)

> ลูกสาวหัวหน้าเผ่าบนภูเขาหิมะ ที่ถูกเนรเทศเพราะพังประตูเมืองด้วยมือเปล่า และยังภูมิใจกับมันอยู่  
> *A snow-mountain chieftain's daughter exiled for breaking down a city gate with her bare hands, and still proud of it.*

**Stat**: HP 36 · STR 8 · END 5 · AGI 2 · มานา 4 — total 9 + 8 + 5 + 2 + 2 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `all-in-swing` | ฟันหมดหน้าตัก / All-In Swing | 1 | การโจมตีครั้งนี้ได้ +3 ในทอยโดน แต่จนถึงตาถัดไปของเรา ศัตรูที่โจมตีเราก็ได้ +3 เหมือนกัน<br>*Get +3 on this attack roll, but until your next turn, enemies attacking you also get +3.* |
| `mountain-wrath` | โทสะขุนเขา / Mountain Wrath | 2 | นาน 3 ตา ทุกครั้งที่ฟันโดนแรงขึ้น +2 และดาเมจที่เราโดนลดลงครั้งละ 2 แต้ม ระหว่างนั้นใช้สกิลอื่นไม่ได้<br>*For 3 turns your melee hits deal +2 damage and every hit on you deals 2 less. You can't use other skills meanwhile.* |
| `boulder-toss` | ขว้างหินผา / Boulder Toss | 3 | ยกหินหรือของหนักขว้างใส่ศัตรูในระยะ 4 ช่อง ถ้าโดน เสียหาย d8 และล้มลง เดินไม่ได้ในตาถัดไปของมัน<br>*Hurl a rock or something heavy at an enemy within 4 squares. On a hit it takes d8 damage and is knocked down: it can't move on its next turn.* |

**ของติดตัว**

- ค้อนศึกหินยอดเขา / Peakstone War Maul ×1 — ของประจำตัว (signature)
- เข็มขัดหนังหมีขาว / White Bear Hide Belt ×1
- เชือกปีนเขา / Climbing Rope ×1
- เนื้อแพะตากแห้ง / Dried Goat Meat ×4

## 2. นักกวีเวทดนตรี / Bard — `legend-bard`
**โอดาลิส รีด** / **Odalys Reed** · role: `support` (สนับสนุน)

> นักเขียนบทละครเร่ที่ออกเดินทางเก็บเรื่องจริง เพราะเชื่อว่าโลกแต่งพล็อตได้ดีกว่าเธอเสมอ  
> *A travelling playwright who collects true stories, because she believes the world always writes better plots than she does.*

**Stat**: HP 20 · STR 4 · END 4 · AGI 7 · มานา 12 — total 5 + 4 + 4 + 7 + 6 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `hecklers-verse` | กลอนแซว / Heckler's Verse | 1 | แต่งกลอนล้อศัตรู 1 ตัวในระยะ 6 ช่องที่ได้ยินเรา มันเสียสมาธิ ทอยครั้งถัดไปของมัน −3<br>*Mock one enemy within 6 squares that can hear you with a cutting rhyme: its next roll gets −3.* |
| `chorus-of-courage` | ประสานเสียงกล้า / Chorus of Courage | 3 | ร้องเพลงปลุกขวัญ เพื่อนทุกคนในระยะ 4 ช่องได้ +2 ในทุกการทอย จนจบรอบถัดไป<br>*Sing a rallying song: every ally within 4 squares gets +2 on all rolls until the end of the next round.* |
| `curtain-fall` | ม่านปิดฉาก / Curtain Fall | 4 | ร้องเพลงกล่อมศัตรู 1 ตัวในระยะ 5 ช่อง มันหลับไป 2 ตา หรือจนกว่าจะโดนตี (บอสอาจต้านได้ DM ตัดสิน)<br>*Sing one enemy within 5 squares to sleep: it dozes for 2 turns or until it is hit (a boss may resist, DM decides).* |

**ของติดตัว**

- ซอไม้เมเปิลแดง / Red Maple Fiddle ×1 — ของประจำตัว (signature)
- สมุดบทละครที่ยังไม่จบ / Unfinished Play Script ×1
- หน้ากากละครครึ่งหน้า / Half Stage Mask ×1
- มีดสั้นด้ามมุก / Pearl-Handled Dagger ×1
- ลูกอมน้ำผึ้งมะนาว / Honey-Lemon Drop ×4

## 3. นักบวช / Cleric — `legend-cleric`
**ออสวิน เกรย์ไลต์** / **Oswin Graylight** · role: `healer` (ฮีลเลอร์)

> อดีตทหารที่หันมาบวชหลังรอดตายจากศึกที่ไม่ควรมีใครรอด และยังถามเทพทุกคืนว่าทำไมต้องเป็นเขา  
> *A former soldier who took holy vows after surviving a battle no one should have survived, and still asks his god every night: why me?*

**Stat**: HP 20 · STR 6 · END 5 · AGI 2 · มานา 16 — total 5 + 6 + 5 + 2 + 8 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `mending-prayer` | คำสวดสมานแผล / Mending Prayer | 2 | สวดให้เพื่อน 1 คนในระยะ 4 ช่อง เพื่อนฟื้นเลือด d6 + 2<br>*Pray over one ally within 4 squares: they regain d6 + 2 HP.* |
| `blessed-ranks` | แถวทัพได้พร / Blessed Ranks | 3 | ให้พรเพื่อนสูงสุด 3 คนในระยะ 4 ช่อง ได้ +1 ในทุกการทอย นาน 3 รอบ<br>*Bless up to 3 allies within 4 squares: they get +1 on every roll for 3 rounds.* |
| `judgment-lantern` | ตะเกียงพิพากษา / Judgment Lantern | 4 | ส่องแสงศักดิ์สิทธิ์ใส่ศัตรู 1 ตัวในระยะ 6 ช่อง โดนแน่นอน เสียหาย d8 และแรงเป็น 2 เท่าถ้าเป็นซากศพเดินได้<br>*Shine holy light on one enemy within 6 squares. It hits automatically for d8 damage, doubled against undead.* |

**ของติดตัว**

- ตะเกียงเงินแห่งวิหาร / Temple Silver Lantern ×1 — ของประจำตัว (signature)
- ค้อนศึกหัวเหล็ก / Iron-Headed Warhammer ×1
- โล่ไม้ตราตะวันสีซีด / Faded Sunmark Shield ×1
- ผ้าพันแผลชุบน้ำมนต์ / Blessed Bandage ×4
- ป้ายชื่อเพื่อนร่วมรบ / Fallen Comrade's Tag ×1

## 4. ดรูอิด / Druid — `legend-druid`
**แทนซี ฮอลโลว์เฟิร์น** / **Tansy Hollowfern** · role: `hybrid` (ไฮบริด)

> หญิงชราหมอยาชายป่าที่เด็กๆ ในหมู่บ้านลือว่าเป็นแม่มด และเธอก็ไม่เคยคิดจะแก้ข่าว  
> *An old hedge-herbalist the village children swear is a witch, and she has never bothered to correct them.*

**Stat**: HP 20 · STR 3 · END 7 · AGI 5 · มานา 12 — total 5 + 3 + 7 + 5 + 6 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `herb-poultice` | ยาพอกสมุนไพร / Herb Poultice | 2 | พอกยาให้เพื่อนที่อยู่ติดกัน เพื่อนฟื้นเลือด d6 และหายจากพิษ 1 อย่าง<br>*Press a poultice on an ally next to you: they regain d6 HP and are cured of one poison.* |
| `bramble-wall` | กำแพงหนาม / Bramble Wall | 3 | เสกพุ่มหนามเป็นแนวตรงยาว 4 ช่องในระยะ 6 ช่อง ใครเดินเข้าแนวหนามต้องหยุดเดินทันทีและเสียหาย d4 อยู่นาน 3 รอบ<br>*Grow a straight line of brambles 4 squares long within 6 squares. Anyone who steps into it stops moving and takes d4 damage. It lasts 3 rounds.* |
| `old-boar-shape` | ร่างหมูป่าเฒ่า / Old Boar Shape | 5 | แปลงร่างเป็นหมูป่าตัวโตนาน 3 ตา ได้เลือดเสริม 8 แต้ม (หายเมื่อคืนร่าง) นับ STR เป็น 8 และพุ่งชนเป็นเส้นตรงได้ 4 ช่องก่อนฟัน ระหว่างเป็นหมูป่าใช้สกิลอื่นไม่ได้<br>*Become a great boar for 3 turns: gain 8 extra HP (gone when you change back), your Strength counts as 8, and you may charge 4 squares in a straight line before striking. No other skills as a boar.* |

**ของติดตัว**

- ไม้เท้าเอลเดอร์บิดเกลียว / Twisted Elder Staff ×1 — ของประจำตัว (signature)
- ตะกร้าสมุนไพร / Herb Basket ×1
- โถน้ำผึ้งป่า / Jar of Wild Honey ×2
- สร้อยลูกโอ๊ก / Acorn Necklace ×1

## 5. นักรบ / Fighter — `legend-fighter`
**เพตรา คาลวอส** / **Petra Kalvos** · role: `tank` (แทงก์)

> ลูกสาวช่างตีเหล็กที่ตั้งใจพิสูจน์ว่าดาบของพ่อดีที่สุดในแผ่นดิน ด้วยการถือมันลงสนามประลองเองทุกสนาม  
> *A blacksmith's daughter set on proving her father's blades are the finest in the land, by carrying them into every arena herself.*

**Stat**: HP 24 · STR 8 · END 7 · AGI 3 · มานา 4 — total 6 + 8 + 7 + 3 + 2 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `catch-your-breath` | ตั้งหลักหายใจ / Catch Your Breath | 1 | ฟื้นเลือดตัวเอง d6 + 2 แล้วยังเดินได้ตามปกติในตานี้<br>*Regain d6 + 2 HP yourself, and still move as normal this turn.* |
| `hold-the-line` | ยืนขวางทาง / Hold the Line | 1 | จนถึงตาถัดไปของเรา ศัตรูที่เดินเข้ามาช่องติดกับเราต้องหยุดเดินทันที และเราฟันสวนมันได้ 1 ครั้ง (ทอยโดนตามปกติ)<br>*Until your next turn, any enemy that steps next to you must stop moving, and you get one free swing at it (roll to hit as normal).* |
| `follow-through` | ฟันต่อเนื่อง / Follow-Through | 3 | ตานี้ได้ทำอีก 1 อย่างทันที เช่น โจมตีอีกครั้ง หรือเดินอีกรอบ<br>*Take one more action right away this turn, such as attacking again or moving again.* |

**ของติดตัว**

- ดาบยาวตราค้อนทั่ง / Anvil-Mark Longsword ×1 — ของประจำตัว (signature)
- โล่กลมขอบเหล็ก / Iron-Rimmed Round Shield ×1
- ค้อนตีเหล็กพกพา / Travel Smithing Hammer ×1
- เหรียญรางวัลสนามประลอง / Arena Prize Medal ×3
- ขนมปังแข็ง / Hardtack Biscuit ×4

## 6. นักพรตยุทธ์ / Monk — `legend-monk`
**คิรัน สติลวอเตอร์** / **Kiran Stillwater** · role: `control` (ควบคุม)

> ศิษย์คนสุดท้ายของวัดบนหน้าผาที่ถูกเผา ออกตามหาคัมภีร์ที่ถูกขโมย โดยพกไปแค่ชาดีกับความอดทน  
> *The last disciple of a burned cliffside monastery, tracking its stolen scrolls with nothing but good tea and patience.*

**Stat**: HP 16 · STR 4 · END 6 · AGI 9 · มานา 6 — total 4 + 4 + 6 + 9 + 3 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `hundred-palms` | ฝ่ามือร้อยครั้ง / Hundred Palms | 1 | ตบฝ่ามือรัวใส่ศัตรูที่อยู่ติดกัน ทอยโจมตี 2 ครั้ง ถ้าโดน ครั้งละ d4<br>*Rain open-palm strikes on an adjacent enemy: make 2 attack rolls, each hit deals d4 damage.* |
| `drifting-leaf` | ใบไม้ลอยลม / Drifting Leaf | 1 | ตานี้เดินได้เพิ่มอีก 4 ช่อง เดินผ่านช่องศัตรูได้ และไม่มีใครฟันสวนเราระหว่างเดิน<br>*Move 4 extra squares this turn. You can pass through enemy squares and no one gets a swing at you on the way.* |
| `still-the-pulse` | สะกดชีพจร / Still the Pulse | 3 | จิ้มจุดชีพจรศัตรูที่อยู่ติดกัน ถ้าโดน มันต้องข้ามตาถัดไปทั้งตา (บอสอาจต้านได้ DM ตัดสิน)<br>*Jab a pressure point on an adjacent enemy. On a hit it skips its next turn entirely (a boss may resist, DM decides).* |

**ของติดตัว**

- ไม้พลองไผ่ดำ / Black Bamboo Staff ×1 — ของประจำตัว (signature)
- กาน้ำชาดินเผา / Clay Teapot ×1
- ใบชาภูเขาห่อใบตอง / Mountain Tea in Leaf Wrap ×5
- ผ้าพันมือ / Hand Wraps ×1
- เศษคัมภีร์ฉีกขาด / Torn Scroll Fragment ×1

## 7. อัศวินศักดิ์สิทธิ์ / Paladin — `legend-paladin`
**โทเบียส เวลลัน** / **Tobias Vellan** · role: `hybrid` (ไฮบริด)

> ชาวนาธรรมดาที่ได้ยินเสียงเรียกในฝันคืนที่หมู่บ้านถูกเผา และยังไม่แน่ใจว่าตัวเองคู่ควรกับคำสาบานที่รับมา  
> *A plain farmer who heard a calling in his dreams the night his village burned, and still isn't sure he deserves the oath he took.*

**Stat**: HP 28 · STR 6 · END 8 · AGI 2 · มานา 6 — total 7 + 6 + 8 + 2 + 3 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `oathstrike` | ฟันตามคำสาบาน / Oathstrike | 2 | การฟันโดนครั้งถัดไปแรงขึ้น +d8 หรือ +2d8 ถ้าเป็นซากศพเดินได้หรือปีศาจ<br>*Your next melee hit deals +d8 damage, or +2d8 against undead and demons.* |
| `hearth-touch` | สัมผัสไออุ่นบ้าน / Hearth Touch | 2 | แตะเพื่อนที่อยู่ติดกัน เลือกให้เพื่อนฟื้นเลือด d6 + 3 หรือล้างอาการหวาดกลัวออก<br>*Touch an ally next to you: either they regain d6 + 3 HP, or their fear is lifted.* |
| `banner-of-home` | ธงแห่งบ้านเกิด / Banner of Home | 3 | ชูธงคำสาบาน จนจบรอบถัดไป เพื่อนทุกคนในระยะ 2 ช่องจากเรา (รวมตัวเรา) โดนตีแล้วดาเมจลดลงครั้งละ 2 แต้ม<br>*Raise your oath banner: until the end of the next round, every ally within 2 squares of you (you included) takes 2 less damage from each hit.* |

**ของติดตัว**

- ง้าวตีใหม่จากเคียวไร่ / Reforged Harvest Glaive ×1 — ของประจำตัว (signature)
- เสื้อเกราะโซ่มือสอง / Secondhand Chainmail ×1
- ผ้าพันคอปักลายรวงข้าว / Wheat-Embroidered Scarf ×1
- ขนมปังข้าวไรย์ฝีมือแม่ / Mother's Rye Loaf ×2
- กำเมล็ดข้าวจากบ้านเกิด / Handful of Home Seed Grain ×1

## 8. พรานป่า / Ranger — `legend-ranger`
**โยเรน วาสก์** / **Joren Vask** · role: `damage` (ดาเมจ)

> คนนำทางกองคาราวานที่รู้ทุกทางลัดบนภูเขา และร้อยเขี้ยวสัตว์ร้ายทุกตัวที่ล่าได้เป็นสร้อยคอ  
> *A caravan guide who knows every mountain shortcut and strings a fang from every beast he has hunted onto his necklace.*

**Stat**: HP 28 · STR 4 · END 4 · AGI 8 · มานา 6 — total 7 + 4 + 4 + 8 + 3 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `trackers-eye` | ตาพรานตามรอย / Tracker's Eye | 1 | อ่านรอยบนพื้นโดยไม่ต้องทอย DM บอกว่ามีอะไรผ่านมา กี่ตัว และไปทางไหน<br>*Read the tracks on the ground with no roll: the DM tells you what passed, how many, and which way they went.* |
| `pinning-arrow` | ธนูตรึงขา / Pinning Arrow | 2 | ยิงศัตรูในระยะ 10 ช่อง ถ้าโดน เสียหาย d6 และมันเดินไม่ได้ในตาถัดไปของมัน<br>*Shoot an enemy within 10 squares. On a hit it takes d6 damage and can't move on its next turn.* |
| `double-loose` | ยิงคู่ / Double Loose | 3 | ยิงธนู 2 ดอกใส่ศัตรู 1 หรือ 2 ตัวในระยะ 8 ช่อง ทอยโดนแยกกัน ดอกที่โดนเสียหายดอกละ d8<br>*Loose 2 arrows at one or two enemies within 8 squares. Roll each hit separately; each arrow that hits deals d8 damage.* |

**ของติดตัว**

- ธนูสั้นเขาแพะภูเขา / Ram-Horn Shortbow ×1 — ของประจำตัว (signature)
- ลูกธนูหัวเหล็ก / Iron-Tipped Arrow ×20
- สร้อยเขี้ยวสัตว์ร้าย / Fang Trophy Necklace ×1
- แผนที่ทางลัดภูเขา / Mountain Shortcut Map ×1
- เนื้อแห้งพริกไทย / Peppered Jerky ×5

## 9. นักย่องเบา / Rogue — `legend-rogue`
**ไซลัส เรน** / **Silas Wren** · role: `damage` (ดาเมจ)

> สายลับของสภาเมืองที่แฝงตัวเข้ากลุ่มนักผจญภัย และยังตัดสินใจไม่ได้ว่าจะส่งรายงานหรือเลือกเพื่อนใหม่  
> *A city council spy planted among adventurers, who still can't decide whether to file his report or side with his new friends.*

**Stat**: HP 16 · STR 6 · END 3 · AGI 10 · มานา 6 — total 4 + 6 + 3 + 10 + 3 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `gap-in-the-guard` | ช่องโหว่การ์ด / Gap in the Guard | 1 | ใช้ได้เมื่อมีเพื่อนยืนติดกับเป้าหมาย การโจมตีครั้งนี้ ถ้าโดน แรงขึ้น +d6<br>*Use when an ally stands next to your target. If this attack hits, it deals +d6 damage.* |
| `slip-away` | ถอยลับตา / Slip Away | 1 | หลังโจมตีแล้ว เดินต่อได้อีก 3 ช่อง และศัตรูฟันสวนเราไม่ได้ระหว่างถอย<br>*After you attack, move 3 more squares. Enemies can't swing at you as you slip away.* |
| `borrowed-face` | ยืมหน้าคนอื่น / Borrowed Face | 3 | ปลอมตัวเป็นคนที่เคยเห็นหน้า คนทั่วไปไม่สงสัยจนจบฉาก แต่ถ้าคนสนิทของตัวจริงมาคุยด้วย DM อาจให้ทอย AGI<br>*Disguise yourself as someone you have seen. Ordinary folk won't doubt you until the scene ends, but if someone close to the real person talks to you, the DM may ask for an AGI roll.* |

**ของติดตัว**

- มีดสั้นใบดำซ่อนแขนเสื้อ / Black Sleeve Dagger ×1 — ของประจำตัว (signature)
- แหวนตราสภาเมือง / Council Signet Ring ×1
- ชุดปลอมตัว / Disguise Kit ×1
- ขวดหมึกล่องหน / Vial of Invisible Ink ×2
- มีดขว้าง / Throwing Knife ×4

## 10. นักเวทโดยกำเนิด / Sorcerer — `legend-sorcerer`
**เรมี ซอลต์ไฟร์** / **Remy Saltfire** · role: `damage` (ดาเมจ)

> ลูกชาวประมงที่เผาเรือของพ่อตอนโกรธครั้งแรกจนต้องหนีออกจากหมู่บ้าน ตอนนี้กำลังหัดโกรธให้ถูกเวลา  
> *A fisher's son who set his father's boat ablaze the first time he lost his temper and had to flee the village. Now he is learning to get angry at the right moment.*

**Stat**: HP 16 · STR 2 · END 6 · AGI 5 · มานา 18 — total 4 + 2 + 6 + 5 + 9 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `ember-flick` | ดีดสะเก็ดไฟ / Ember Flick | 1 | ดีดสะเก็ดไฟจากปลายนิ้วใส่ศัตรูในระยะ 6 ช่อง ถ้าโดน เสียหาย d6<br>*Flick a spark from your fingertip at an enemy within 6 squares. On a hit it takes d6 damage.* |
| `bend-the-spell` | ดัดเวท / Bend the Spell | 2 | ใช้ก่อนสกิลอื่นในตาเดียวกัน เลือก 1 อย่าง: สกิลนั้นยิงได้ไกลขึ้น 2 เท่า หรือเลือกเป้าหมายเพิ่มได้อีก 1 ตัว<br>*Use before another skill in the same turn and pick one: that skill's range doubles, or it can target one extra creature.* |
| `firetide` | คลื่นเพลิง / Firetide | 5 | พ่นไฟเป็นเส้นตรงยาว 6 ช่องจากตัวเรา ทุกตัวในเส้น ทั้งศัตรูและเพื่อน เสียหาย 2d6 โดยไม่ต้องทอยโดน<br>*Breathe a straight line of fire 6 squares long from you. Everyone in the line, friend or foe, takes 2d6 damage, no hit roll needed.* |

**ของติดตัว**

- จี้เปลือกหอยเกลือ / Salt-Shell Pendant ×1 — ของประจำตัว (signature)
- ถุงมือหนังไหม้เกรียม / Scorched Leather Gloves ×1
- มีดแล่ปลา / Filleting Knife ×1
- ขวดน้ำทะเลไว้ดับใจร้อน / Flask of Calming Seawater ×2
- ปลาเค็มตากแห้ง / Salted Dried Fish ×4

## 11. ผู้ทำสัญญา / Warlock — `legend-warlock`
**เวสนา มอร์โรว์** / **Vesna Morrow** · role: `hybrid` (ไฮบริด)

> หมอยาที่ทำสัญญากับบางสิ่งจากใต้ทะเลลึกเพื่อรักษาลูกสาว และตอนนี้มันเริ่มกระซิบทวงค่าตอบแทนทุกคืน  
> *An apothecary who struck a pact with something from the deep sea to save her daughter, and now it whispers every night for its payment.*

**Stat**: HP 24 · STR 3 · END 4 · AGI 5 · มานา 16 — total 6 + 3 + 4 + 5 + 8 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `abyssal-lance` | หอกห้วงลึก / Abyssal Lance | 1 | ยิงพลังมืดสีเขียวน้ำทะเลใส่ศัตรูในระยะ 8 ช่อง ถ้าโดน เสียหาย d8<br>*Hurl a lance of sea-green darkness at an enemy within 8 squares. On a hit it takes d8 damage.* |
| `mark-of-debt` | ตราหนี้ / Mark of Debt | 2 | ประทับตราหนี้ลงบนศัตรู 1 ตัวในระยะ 6 ช่อง นาน 3 ตา ทุกครั้งที่ใครในกลุ่มโจมตีมันโดน มันเสียเพิ่มอีก 2 แต้ม<br>*Brand one enemy within 6 squares with a debt mark for 3 turns: every time anyone in the party hits it, it takes 2 extra damage.* |
| `grasp-of-the-deep` | หนวดจากห้วงลึก / Grasp of the Deep | 4 | หนวดยักษ์โผล่จากพื้นในพื้นที่ 2×2 ช่องในระยะ 6 ช่อง ศัตรูในพื้นที่เดินไม่ได้ 2 ตา และเสียหาย d4 ทุกต้นตาของมัน โดยไม่ต้องทอยโดน<br>*Giant tentacles burst from a 2×2 area within 6 squares. Enemies inside can't move for 2 turns and take d4 damage at the start of each of their turns, no hit roll needed.* |

**ของติดตัว**

- จี้ปะการังดำ / Black Coral Pendant ×1 — ของประจำตัว (signature)
- กระเป๋าหมอยา / Apothecary's Satchel ×1
- มีดผ่าแผลเงิน / Silver Lancet ×1
- ยาแก้ไข้ / Fever Draught ×3
- จดหมายจากลูกสาว / Letter from Her Daughter ×1

## 12. พ่อมด / Wizard — `legend-wizard`
**ธีโอดริก เพลล์** / **Theodric Pell** · role: `control` (ควบคุม)

> บรรณารักษ์ที่แอบเรียนเวทจากหนังสือต้องห้ามในห้องใต้ดินของหอสมุด และถูกไล่ออกตอนหนังสือเล่มหนึ่งกัดคนอ่าน  
> *A librarian who secretly learned magic from the forbidden books in the archive vault, and was dismissed when one of them bit a reader.*

**Stat**: HP 16 · STR 2 · END 4 · AGI 6 · มานา 20 — total 4 + 2 + 4 + 6 + 10 = **26**

**สกิลติดตัว**

| id | ชื่อ | มานา | ทำอะไร |
|---|---|---|---|
| `ink-dart` | ลูกดอกหมึก / Ink Dart | 1 | ยิงหมึกเวทใส่ศัตรูในระยะ 6 ช่อง ถ้าโดน เสียหาย d6 และหมึกเข้าตา ทอยโจมตีครั้งถัดไปของมัน −2<br>*Shoot magic ink at an enemy within 6 squares. On a hit it takes d6 damage and the ink blurs its eyes: its next attack roll gets −2.* |
| `paper-ward` | กระดาษคุ้มภัย / Paper Ward | 2 | หน้ากระดาษปลิวมาบังเพื่อน 1 คนในระยะ 5 ช่อง (หรือตัวเรา) ครั้งถัดไปที่เขาโดนตี ดาเมจลดลง d8 แต้ม<br>*Pages fly up to shield one ally within 5 squares (or yourself): the next hit on them deals d8 less damage.* |
| `binding-glyph` | อักขระตรึง / Binding Glyph | 5 | เขียนอักขระลงพื้นที่ 3×3 ช่องในระยะ 6 ช่อง ศัตรูที่อยู่หรือเดินเข้าไปเดินไม่ได้ 2 ตา (แต่ยังโจมตีได้) อักขระอยู่นาน 3 รอบ<br>*Inscribe a glyph over a 3×3 area within 6 squares. Enemies in it, or who step in, can't move for 2 turns (they can still attack). The glyph lasts 3 rounds.* |

**ของติดตัว**

- ตำราเวทจอมกัด / The Biting Grimoire ×1 — ของประจำตัว (signature)
- ไม้เท้าไม้มะเกลือ / Ebony Staff ×1
- แว่นขยายทองเหลือง / Brass Magnifying Glass ×1
- บัตรยืมหนังสือหมดอายุ / Expired Library Card ×1
- ชาเขียวเข้มข้น / Strong Green Tea ×2

## JSON (machine-readable)

```json
[
  {
    "id": "legend-barbarian",
    "set": "legend",
    "heroName": {
      "th": "คายา สตอร์มฮอร์น",
      "en": "Kaya Stormhorn"
    },
    "className": {
      "th": "นักรบเถื่อน",
      "en": "Barbarian"
    },
    "blurb": {
      "th": "ลูกสาวหัวหน้าเผ่าบนภูเขาหิมะ ที่ถูกเนรเทศเพราะพังประตูเมืองด้วยมือเปล่า และยังภูมิใจกับมันอยู่",
      "en": "A snow-mountain chieftain's daughter exiled for breaking down a city gate with her bare hands, and still proud of it."
    },
    "role": "damage",
    "art": "",
    "stats": {
      "hpMax": 36,
      "str": 8,
      "end": 5,
      "agi": 2,
      "mpMax": 4
    },
    "skills": [
      {
        "id": "all-in-swing",
        "name": {
          "th": "ฟันหมดหน้าตัก",
          "en": "All-In Swing"
        },
        "desc": {
          "th": "การโจมตีครั้งนี้ได้ +3 ในทอยโดน แต่จนถึงตาถัดไปของเรา ศัตรูที่โจมตีเราก็ได้ +3 เหมือนกัน",
          "en": "Get +3 on this attack roll, but until your next turn, enemies attacking you also get +3."
        },
        "cost": 1
      },
      {
        "id": "mountain-wrath",
        "name": {
          "th": "โทสะขุนเขา",
          "en": "Mountain Wrath"
        },
        "desc": {
          "th": "นาน 3 ตา ทุกครั้งที่ฟันโดนแรงขึ้น +2 และดาเมจที่เราโดนลดลงครั้งละ 2 แต้ม ระหว่างนั้นใช้สกิลอื่นไม่ได้",
          "en": "For 3 turns your melee hits deal +2 damage and every hit on you deals 2 less. You can't use other skills meanwhile."
        },
        "cost": 2
      },
      {
        "id": "boulder-toss",
        "name": {
          "th": "ขว้างหินผา",
          "en": "Boulder Toss"
        },
        "desc": {
          "th": "ยกหินหรือของหนักขว้างใส่ศัตรูในระยะ 4 ช่อง ถ้าโดน เสียหาย d8 และล้มลง เดินไม่ได้ในตาถัดไปของมัน",
          "en": "Hurl a rock or something heavy at an enemy within 4 squares. On a hit it takes d8 damage and is knocked down: it can't move on its next turn."
        },
        "cost": 3
      }
    ],
    "items": [
      {
        "name": {
          "th": "ค้อนศึกหินยอดเขา",
          "en": "Peakstone War Maul"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เข็มขัดหนังหมีขาว",
          "en": "White Bear Hide Belt"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เชือกปีนเขา",
          "en": "Climbing Rope"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เนื้อแพะตากแห้ง",
          "en": "Dried Goat Meat"
        },
        "qty": 4
      }
    ]
  },
  {
    "id": "legend-bard",
    "set": "legend",
    "heroName": {
      "th": "โอดาลิส รีด",
      "en": "Odalys Reed"
    },
    "className": {
      "th": "นักกวีเวทดนตรี",
      "en": "Bard"
    },
    "blurb": {
      "th": "นักเขียนบทละครเร่ที่ออกเดินทางเก็บเรื่องจริง เพราะเชื่อว่าโลกแต่งพล็อตได้ดีกว่าเธอเสมอ",
      "en": "A travelling playwright who collects true stories, because she believes the world always writes better plots than she does."
    },
    "role": "support",
    "art": "",
    "stats": {
      "hpMax": 20,
      "str": 4,
      "end": 4,
      "agi": 7,
      "mpMax": 12
    },
    "skills": [
      {
        "id": "hecklers-verse",
        "name": {
          "th": "กลอนแซว",
          "en": "Heckler's Verse"
        },
        "desc": {
          "th": "แต่งกลอนล้อศัตรู 1 ตัวในระยะ 6 ช่องที่ได้ยินเรา มันเสียสมาธิ ทอยครั้งถัดไปของมัน −3",
          "en": "Mock one enemy within 6 squares that can hear you with a cutting rhyme: its next roll gets −3."
        },
        "cost": 1
      },
      {
        "id": "chorus-of-courage",
        "name": {
          "th": "ประสานเสียงกล้า",
          "en": "Chorus of Courage"
        },
        "desc": {
          "th": "ร้องเพลงปลุกขวัญ เพื่อนทุกคนในระยะ 4 ช่องได้ +2 ในทุกการทอย จนจบรอบถัดไป",
          "en": "Sing a rallying song: every ally within 4 squares gets +2 on all rolls until the end of the next round."
        },
        "cost": 3
      },
      {
        "id": "curtain-fall",
        "name": {
          "th": "ม่านปิดฉาก",
          "en": "Curtain Fall"
        },
        "desc": {
          "th": "ร้องเพลงกล่อมศัตรู 1 ตัวในระยะ 5 ช่อง มันหลับไป 2 ตา หรือจนกว่าจะโดนตี (บอสอาจต้านได้ DM ตัดสิน)",
          "en": "Sing one enemy within 5 squares to sleep: it dozes for 2 turns or until it is hit (a boss may resist, DM decides)."
        },
        "cost": 4
      }
    ],
    "items": [
      {
        "name": {
          "th": "ซอไม้เมเปิลแดง",
          "en": "Red Maple Fiddle"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "สมุดบทละครที่ยังไม่จบ",
          "en": "Unfinished Play Script"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "หน้ากากละครครึ่งหน้า",
          "en": "Half Stage Mask"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มีดสั้นด้ามมุก",
          "en": "Pearl-Handled Dagger"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ลูกอมน้ำผึ้งมะนาว",
          "en": "Honey-Lemon Drop"
        },
        "qty": 4
      }
    ]
  },
  {
    "id": "legend-cleric",
    "set": "legend",
    "heroName": {
      "th": "ออสวิน เกรย์ไลต์",
      "en": "Oswin Graylight"
    },
    "className": {
      "th": "นักบวช",
      "en": "Cleric"
    },
    "blurb": {
      "th": "อดีตทหารที่หันมาบวชหลังรอดตายจากศึกที่ไม่ควรมีใครรอด และยังถามเทพทุกคืนว่าทำไมต้องเป็นเขา",
      "en": "A former soldier who took holy vows after surviving a battle no one should have survived, and still asks his god every night: why me?"
    },
    "role": "healer",
    "art": "",
    "stats": {
      "hpMax": 20,
      "str": 6,
      "end": 5,
      "agi": 2,
      "mpMax": 16
    },
    "skills": [
      {
        "id": "mending-prayer",
        "name": {
          "th": "คำสวดสมานแผล",
          "en": "Mending Prayer"
        },
        "desc": {
          "th": "สวดให้เพื่อน 1 คนในระยะ 4 ช่อง เพื่อนฟื้นเลือด d6 + 2",
          "en": "Pray over one ally within 4 squares: they regain d6 + 2 HP."
        },
        "cost": 2
      },
      {
        "id": "blessed-ranks",
        "name": {
          "th": "แถวทัพได้พร",
          "en": "Blessed Ranks"
        },
        "desc": {
          "th": "ให้พรเพื่อนสูงสุด 3 คนในระยะ 4 ช่อง ได้ +1 ในทุกการทอย นาน 3 รอบ",
          "en": "Bless up to 3 allies within 4 squares: they get +1 on every roll for 3 rounds."
        },
        "cost": 3
      },
      {
        "id": "judgment-lantern",
        "name": {
          "th": "ตะเกียงพิพากษา",
          "en": "Judgment Lantern"
        },
        "desc": {
          "th": "ส่องแสงศักดิ์สิทธิ์ใส่ศัตรู 1 ตัวในระยะ 6 ช่อง โดนแน่นอน เสียหาย d8 และแรงเป็น 2 เท่าถ้าเป็นซากศพเดินได้",
          "en": "Shine holy light on one enemy within 6 squares. It hits automatically for d8 damage, doubled against undead."
        },
        "cost": 4
      }
    ],
    "items": [
      {
        "name": {
          "th": "ตะเกียงเงินแห่งวิหาร",
          "en": "Temple Silver Lantern"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ค้อนศึกหัวเหล็ก",
          "en": "Iron-Headed Warhammer"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "โล่ไม้ตราตะวันสีซีด",
          "en": "Faded Sunmark Shield"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ผ้าพันแผลชุบน้ำมนต์",
          "en": "Blessed Bandage"
        },
        "qty": 4
      },
      {
        "name": {
          "th": "ป้ายชื่อเพื่อนร่วมรบ",
          "en": "Fallen Comrade's Tag"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-druid",
    "set": "legend",
    "heroName": {
      "th": "แทนซี ฮอลโลว์เฟิร์น",
      "en": "Tansy Hollowfern"
    },
    "className": {
      "th": "ดรูอิด",
      "en": "Druid"
    },
    "blurb": {
      "th": "หญิงชราหมอยาชายป่าที่เด็กๆ ในหมู่บ้านลือว่าเป็นแม่มด และเธอก็ไม่เคยคิดจะแก้ข่าว",
      "en": "An old hedge-herbalist the village children swear is a witch, and she has never bothered to correct them."
    },
    "role": "hybrid",
    "art": "",
    "stats": {
      "hpMax": 20,
      "str": 3,
      "end": 7,
      "agi": 5,
      "mpMax": 12
    },
    "skills": [
      {
        "id": "herb-poultice",
        "name": {
          "th": "ยาพอกสมุนไพร",
          "en": "Herb Poultice"
        },
        "desc": {
          "th": "พอกยาให้เพื่อนที่อยู่ติดกัน เพื่อนฟื้นเลือด d6 และหายจากพิษ 1 อย่าง",
          "en": "Press a poultice on an ally next to you: they regain d6 HP and are cured of one poison."
        },
        "cost": 2
      },
      {
        "id": "bramble-wall",
        "name": {
          "th": "กำแพงหนาม",
          "en": "Bramble Wall"
        },
        "desc": {
          "th": "เสกพุ่มหนามเป็นแนวตรงยาว 4 ช่องในระยะ 6 ช่อง ใครเดินเข้าแนวหนามต้องหยุดเดินทันทีและเสียหาย d4 อยู่นาน 3 รอบ",
          "en": "Grow a straight line of brambles 4 squares long within 6 squares. Anyone who steps into it stops moving and takes d4 damage. It lasts 3 rounds."
        },
        "cost": 3
      },
      {
        "id": "old-boar-shape",
        "name": {
          "th": "ร่างหมูป่าเฒ่า",
          "en": "Old Boar Shape"
        },
        "desc": {
          "th": "แปลงร่างเป็นหมูป่าตัวโตนาน 3 ตา ได้เลือดเสริม 8 แต้ม (หายเมื่อคืนร่าง) นับ STR เป็น 8 และพุ่งชนเป็นเส้นตรงได้ 4 ช่องก่อนฟัน ระหว่างเป็นหมูป่าใช้สกิลอื่นไม่ได้",
          "en": "Become a great boar for 3 turns: gain 8 extra HP (gone when you change back), your Strength counts as 8, and you may charge 4 squares in a straight line before striking. No other skills as a boar."
        },
        "cost": 5
      }
    ],
    "items": [
      {
        "name": {
          "th": "ไม้เท้าเอลเดอร์บิดเกลียว",
          "en": "Twisted Elder Staff"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ตะกร้าสมุนไพร",
          "en": "Herb Basket"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "โถน้ำผึ้งป่า",
          "en": "Jar of Wild Honey"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "สร้อยลูกโอ๊ก",
          "en": "Acorn Necklace"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-fighter",
    "set": "legend",
    "heroName": {
      "th": "เพตรา คาลวอส",
      "en": "Petra Kalvos"
    },
    "className": {
      "th": "นักรบ",
      "en": "Fighter"
    },
    "blurb": {
      "th": "ลูกสาวช่างตีเหล็กที่ตั้งใจพิสูจน์ว่าดาบของพ่อดีที่สุดในแผ่นดิน ด้วยการถือมันลงสนามประลองเองทุกสนาม",
      "en": "A blacksmith's daughter set on proving her father's blades are the finest in the land, by carrying them into every arena herself."
    },
    "role": "tank",
    "art": "",
    "stats": {
      "hpMax": 24,
      "str": 8,
      "end": 7,
      "agi": 3,
      "mpMax": 4
    },
    "skills": [
      {
        "id": "catch-your-breath",
        "name": {
          "th": "ตั้งหลักหายใจ",
          "en": "Catch Your Breath"
        },
        "desc": {
          "th": "ฟื้นเลือดตัวเอง d6 + 2 แล้วยังเดินได้ตามปกติในตานี้",
          "en": "Regain d6 + 2 HP yourself, and still move as normal this turn."
        },
        "cost": 1
      },
      {
        "id": "hold-the-line",
        "name": {
          "th": "ยืนขวางทาง",
          "en": "Hold the Line"
        },
        "desc": {
          "th": "จนถึงตาถัดไปของเรา ศัตรูที่เดินเข้ามาช่องติดกับเราต้องหยุดเดินทันที และเราฟันสวนมันได้ 1 ครั้ง (ทอยโดนตามปกติ)",
          "en": "Until your next turn, any enemy that steps next to you must stop moving, and you get one free swing at it (roll to hit as normal)."
        },
        "cost": 1
      },
      {
        "id": "follow-through",
        "name": {
          "th": "ฟันต่อเนื่อง",
          "en": "Follow-Through"
        },
        "desc": {
          "th": "ตานี้ได้ทำอีก 1 อย่างทันที เช่น โจมตีอีกครั้ง หรือเดินอีกรอบ",
          "en": "Take one more action right away this turn, such as attacking again or moving again."
        },
        "cost": 3
      }
    ],
    "items": [
      {
        "name": {
          "th": "ดาบยาวตราค้อนทั่ง",
          "en": "Anvil-Mark Longsword"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "โล่กลมขอบเหล็ก",
          "en": "Iron-Rimmed Round Shield"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ค้อนตีเหล็กพกพา",
          "en": "Travel Smithing Hammer"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เหรียญรางวัลสนามประลอง",
          "en": "Arena Prize Medal"
        },
        "qty": 3
      },
      {
        "name": {
          "th": "ขนมปังแข็ง",
          "en": "Hardtack Biscuit"
        },
        "qty": 4
      }
    ]
  },
  {
    "id": "legend-monk",
    "set": "legend",
    "heroName": {
      "th": "คิรัน สติลวอเตอร์",
      "en": "Kiran Stillwater"
    },
    "className": {
      "th": "นักพรตยุทธ์",
      "en": "Monk"
    },
    "blurb": {
      "th": "ศิษย์คนสุดท้ายของวัดบนหน้าผาที่ถูกเผา ออกตามหาคัมภีร์ที่ถูกขโมย โดยพกไปแค่ชาดีกับความอดทน",
      "en": "The last disciple of a burned cliffside monastery, tracking its stolen scrolls with nothing but good tea and patience."
    },
    "role": "control",
    "art": "",
    "stats": {
      "hpMax": 16,
      "str": 4,
      "end": 6,
      "agi": 9,
      "mpMax": 6
    },
    "skills": [
      {
        "id": "hundred-palms",
        "name": {
          "th": "ฝ่ามือร้อยครั้ง",
          "en": "Hundred Palms"
        },
        "desc": {
          "th": "ตบฝ่ามือรัวใส่ศัตรูที่อยู่ติดกัน ทอยโจมตี 2 ครั้ง ถ้าโดน ครั้งละ d4",
          "en": "Rain open-palm strikes on an adjacent enemy: make 2 attack rolls, each hit deals d4 damage."
        },
        "cost": 1
      },
      {
        "id": "drifting-leaf",
        "name": {
          "th": "ใบไม้ลอยลม",
          "en": "Drifting Leaf"
        },
        "desc": {
          "th": "ตานี้เดินได้เพิ่มอีก 4 ช่อง เดินผ่านช่องศัตรูได้ และไม่มีใครฟันสวนเราระหว่างเดิน",
          "en": "Move 4 extra squares this turn. You can pass through enemy squares and no one gets a swing at you on the way."
        },
        "cost": 1
      },
      {
        "id": "still-the-pulse",
        "name": {
          "th": "สะกดชีพจร",
          "en": "Still the Pulse"
        },
        "desc": {
          "th": "จิ้มจุดชีพจรศัตรูที่อยู่ติดกัน ถ้าโดน มันต้องข้ามตาถัดไปทั้งตา (บอสอาจต้านได้ DM ตัดสิน)",
          "en": "Jab a pressure point on an adjacent enemy. On a hit it skips its next turn entirely (a boss may resist, DM decides)."
        },
        "cost": 3
      }
    ],
    "items": [
      {
        "name": {
          "th": "ไม้พลองไผ่ดำ",
          "en": "Black Bamboo Staff"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "กาน้ำชาดินเผา",
          "en": "Clay Teapot"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ใบชาภูเขาห่อใบตอง",
          "en": "Mountain Tea in Leaf Wrap"
        },
        "qty": 5
      },
      {
        "name": {
          "th": "ผ้าพันมือ",
          "en": "Hand Wraps"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เศษคัมภีร์ฉีกขาด",
          "en": "Torn Scroll Fragment"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-paladin",
    "set": "legend",
    "heroName": {
      "th": "โทเบียส เวลลัน",
      "en": "Tobias Vellan"
    },
    "className": {
      "th": "อัศวินศักดิ์สิทธิ์",
      "en": "Paladin"
    },
    "blurb": {
      "th": "ชาวนาธรรมดาที่ได้ยินเสียงเรียกในฝันคืนที่หมู่บ้านถูกเผา และยังไม่แน่ใจว่าตัวเองคู่ควรกับคำสาบานที่รับมา",
      "en": "A plain farmer who heard a calling in his dreams the night his village burned, and still isn't sure he deserves the oath he took."
    },
    "role": "hybrid",
    "art": "",
    "stats": {
      "hpMax": 28,
      "str": 6,
      "end": 8,
      "agi": 2,
      "mpMax": 6
    },
    "skills": [
      {
        "id": "oathstrike",
        "name": {
          "th": "ฟันตามคำสาบาน",
          "en": "Oathstrike"
        },
        "desc": {
          "th": "การฟันโดนครั้งถัดไปแรงขึ้น +d8 หรือ +2d8 ถ้าเป็นซากศพเดินได้หรือปีศาจ",
          "en": "Your next melee hit deals +d8 damage, or +2d8 against undead and demons."
        },
        "cost": 2
      },
      {
        "id": "hearth-touch",
        "name": {
          "th": "สัมผัสไออุ่นบ้าน",
          "en": "Hearth Touch"
        },
        "desc": {
          "th": "แตะเพื่อนที่อยู่ติดกัน เลือกให้เพื่อนฟื้นเลือด d6 + 3 หรือล้างอาการหวาดกลัวออก",
          "en": "Touch an ally next to you: either they regain d6 + 3 HP, or their fear is lifted."
        },
        "cost": 2
      },
      {
        "id": "banner-of-home",
        "name": {
          "th": "ธงแห่งบ้านเกิด",
          "en": "Banner of Home"
        },
        "desc": {
          "th": "ชูธงคำสาบาน จนจบรอบถัดไป เพื่อนทุกคนในระยะ 2 ช่องจากเรา (รวมตัวเรา) โดนตีแล้วดาเมจลดลงครั้งละ 2 แต้ม",
          "en": "Raise your oath banner: until the end of the next round, every ally within 2 squares of you (you included) takes 2 less damage from each hit."
        },
        "cost": 3
      }
    ],
    "items": [
      {
        "name": {
          "th": "ง้าวตีใหม่จากเคียวไร่",
          "en": "Reforged Harvest Glaive"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เสื้อเกราะโซ่มือสอง",
          "en": "Secondhand Chainmail"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ผ้าพันคอปักลายรวงข้าว",
          "en": "Wheat-Embroidered Scarf"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขนมปังข้าวไรย์ฝีมือแม่",
          "en": "Mother's Rye Loaf"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "กำเมล็ดข้าวจากบ้านเกิด",
          "en": "Handful of Home Seed Grain"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-ranger",
    "set": "legend",
    "heroName": {
      "th": "โยเรน วาสก์",
      "en": "Joren Vask"
    },
    "className": {
      "th": "พรานป่า",
      "en": "Ranger"
    },
    "blurb": {
      "th": "คนนำทางกองคาราวานที่รู้ทุกทางลัดบนภูเขา และร้อยเขี้ยวสัตว์ร้ายทุกตัวที่ล่าได้เป็นสร้อยคอ",
      "en": "A caravan guide who knows every mountain shortcut and strings a fang from every beast he has hunted onto his necklace."
    },
    "role": "damage",
    "art": "",
    "stats": {
      "hpMax": 28,
      "str": 4,
      "end": 4,
      "agi": 8,
      "mpMax": 6
    },
    "skills": [
      {
        "id": "trackers-eye",
        "name": {
          "th": "ตาพรานตามรอย",
          "en": "Tracker's Eye"
        },
        "desc": {
          "th": "อ่านรอยบนพื้นโดยไม่ต้องทอย DM บอกว่ามีอะไรผ่านมา กี่ตัว และไปทางไหน",
          "en": "Read the tracks on the ground with no roll: the DM tells you what passed, how many, and which way they went."
        },
        "cost": 1
      },
      {
        "id": "pinning-arrow",
        "name": {
          "th": "ธนูตรึงขา",
          "en": "Pinning Arrow"
        },
        "desc": {
          "th": "ยิงศัตรูในระยะ 10 ช่อง ถ้าโดน เสียหาย d6 และมันเดินไม่ได้ในตาถัดไปของมัน",
          "en": "Shoot an enemy within 10 squares. On a hit it takes d6 damage and can't move on its next turn."
        },
        "cost": 2
      },
      {
        "id": "double-loose",
        "name": {
          "th": "ยิงคู่",
          "en": "Double Loose"
        },
        "desc": {
          "th": "ยิงธนู 2 ดอกใส่ศัตรู 1 หรือ 2 ตัวในระยะ 8 ช่อง ทอยโดนแยกกัน ดอกที่โดนเสียหายดอกละ d8",
          "en": "Loose 2 arrows at one or two enemies within 8 squares. Roll each hit separately; each arrow that hits deals d8 damage."
        },
        "cost": 3
      }
    ],
    "items": [
      {
        "name": {
          "th": "ธนูสั้นเขาแพะภูเขา",
          "en": "Ram-Horn Shortbow"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ลูกธนูหัวเหล็ก",
          "en": "Iron-Tipped Arrow"
        },
        "qty": 20
      },
      {
        "name": {
          "th": "สร้อยเขี้ยวสัตว์ร้าย",
          "en": "Fang Trophy Necklace"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "แผนที่ทางลัดภูเขา",
          "en": "Mountain Shortcut Map"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "เนื้อแห้งพริกไทย",
          "en": "Peppered Jerky"
        },
        "qty": 5
      }
    ]
  },
  {
    "id": "legend-rogue",
    "set": "legend",
    "heroName": {
      "th": "ไซลัส เรน",
      "en": "Silas Wren"
    },
    "className": {
      "th": "นักย่องเบา",
      "en": "Rogue"
    },
    "blurb": {
      "th": "สายลับของสภาเมืองที่แฝงตัวเข้ากลุ่มนักผจญภัย และยังตัดสินใจไม่ได้ว่าจะส่งรายงานหรือเลือกเพื่อนใหม่",
      "en": "A city council spy planted among adventurers, who still can't decide whether to file his report or side with his new friends."
    },
    "role": "damage",
    "art": "",
    "stats": {
      "hpMax": 16,
      "str": 6,
      "end": 3,
      "agi": 10,
      "mpMax": 6
    },
    "skills": [
      {
        "id": "gap-in-the-guard",
        "name": {
          "th": "ช่องโหว่การ์ด",
          "en": "Gap in the Guard"
        },
        "desc": {
          "th": "ใช้ได้เมื่อมีเพื่อนยืนติดกับเป้าหมาย การโจมตีครั้งนี้ ถ้าโดน แรงขึ้น +d6",
          "en": "Use when an ally stands next to your target. If this attack hits, it deals +d6 damage."
        },
        "cost": 1
      },
      {
        "id": "slip-away",
        "name": {
          "th": "ถอยลับตา",
          "en": "Slip Away"
        },
        "desc": {
          "th": "หลังโจมตีแล้ว เดินต่อได้อีก 3 ช่อง และศัตรูฟันสวนเราไม่ได้ระหว่างถอย",
          "en": "After you attack, move 3 more squares. Enemies can't swing at you as you slip away."
        },
        "cost": 1
      },
      {
        "id": "borrowed-face",
        "name": {
          "th": "ยืมหน้าคนอื่น",
          "en": "Borrowed Face"
        },
        "desc": {
          "th": "ปลอมตัวเป็นคนที่เคยเห็นหน้า คนทั่วไปไม่สงสัยจนจบฉาก แต่ถ้าคนสนิทของตัวจริงมาคุยด้วย DM อาจให้ทอย AGI",
          "en": "Disguise yourself as someone you have seen. Ordinary folk won't doubt you until the scene ends, but if someone close to the real person talks to you, the DM may ask for an AGI roll."
        },
        "cost": 3
      }
    ],
    "items": [
      {
        "name": {
          "th": "มีดสั้นใบดำซ่อนแขนเสื้อ",
          "en": "Black Sleeve Dagger"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "แหวนตราสภาเมือง",
          "en": "Council Signet Ring"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชุดปลอมตัว",
          "en": "Disguise Kit"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขวดหมึกล่องหน",
          "en": "Vial of Invisible Ink"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "มีดขว้าง",
          "en": "Throwing Knife"
        },
        "qty": 4
      }
    ]
  },
  {
    "id": "legend-sorcerer",
    "set": "legend",
    "heroName": {
      "th": "เรมี ซอลต์ไฟร์",
      "en": "Remy Saltfire"
    },
    "className": {
      "th": "นักเวทโดยกำเนิด",
      "en": "Sorcerer"
    },
    "blurb": {
      "th": "ลูกชาวประมงที่เผาเรือของพ่อตอนโกรธครั้งแรกจนต้องหนีออกจากหมู่บ้าน ตอนนี้กำลังหัดโกรธให้ถูกเวลา",
      "en": "A fisher's son who set his father's boat ablaze the first time he lost his temper and had to flee the village. Now he is learning to get angry at the right moment."
    },
    "role": "damage",
    "art": "",
    "stats": {
      "hpMax": 16,
      "str": 2,
      "end": 6,
      "agi": 5,
      "mpMax": 18
    },
    "skills": [
      {
        "id": "ember-flick",
        "name": {
          "th": "ดีดสะเก็ดไฟ",
          "en": "Ember Flick"
        },
        "desc": {
          "th": "ดีดสะเก็ดไฟจากปลายนิ้วใส่ศัตรูในระยะ 6 ช่อง ถ้าโดน เสียหาย d6",
          "en": "Flick a spark from your fingertip at an enemy within 6 squares. On a hit it takes d6 damage."
        },
        "cost": 1
      },
      {
        "id": "bend-the-spell",
        "name": {
          "th": "ดัดเวท",
          "en": "Bend the Spell"
        },
        "desc": {
          "th": "ใช้ก่อนสกิลอื่นในตาเดียวกัน เลือก 1 อย่าง: สกิลนั้นยิงได้ไกลขึ้น 2 เท่า หรือเลือกเป้าหมายเพิ่มได้อีก 1 ตัว",
          "en": "Use before another skill in the same turn and pick one: that skill's range doubles, or it can target one extra creature."
        },
        "cost": 2
      },
      {
        "id": "firetide",
        "name": {
          "th": "คลื่นเพลิง",
          "en": "Firetide"
        },
        "desc": {
          "th": "พ่นไฟเป็นเส้นตรงยาว 6 ช่องจากตัวเรา ทุกตัวในเส้น ทั้งศัตรูและเพื่อน เสียหาย 2d6 โดยไม่ต้องทอยโดน",
          "en": "Breathe a straight line of fire 6 squares long from you. Everyone in the line, friend or foe, takes 2d6 damage, no hit roll needed."
        },
        "cost": 5
      }
    ],
    "items": [
      {
        "name": {
          "th": "จี้เปลือกหอยเกลือ",
          "en": "Salt-Shell Pendant"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ถุงมือหนังไหม้เกรียม",
          "en": "Scorched Leather Gloves"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มีดแล่ปลา",
          "en": "Filleting Knife"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ขวดน้ำทะเลไว้ดับใจร้อน",
          "en": "Flask of Calming Seawater"
        },
        "qty": 2
      },
      {
        "name": {
          "th": "ปลาเค็มตากแห้ง",
          "en": "Salted Dried Fish"
        },
        "qty": 4
      }
    ]
  },
  {
    "id": "legend-warlock",
    "set": "legend",
    "heroName": {
      "th": "เวสนา มอร์โรว์",
      "en": "Vesna Morrow"
    },
    "className": {
      "th": "ผู้ทำสัญญา",
      "en": "Warlock"
    },
    "blurb": {
      "th": "หมอยาที่ทำสัญญากับบางสิ่งจากใต้ทะเลลึกเพื่อรักษาลูกสาว และตอนนี้มันเริ่มกระซิบทวงค่าตอบแทนทุกคืน",
      "en": "An apothecary who struck a pact with something from the deep sea to save her daughter, and now it whispers every night for its payment."
    },
    "role": "hybrid",
    "art": "",
    "stats": {
      "hpMax": 24,
      "str": 3,
      "end": 4,
      "agi": 5,
      "mpMax": 16
    },
    "skills": [
      {
        "id": "abyssal-lance",
        "name": {
          "th": "หอกห้วงลึก",
          "en": "Abyssal Lance"
        },
        "desc": {
          "th": "ยิงพลังมืดสีเขียวน้ำทะเลใส่ศัตรูในระยะ 8 ช่อง ถ้าโดน เสียหาย d8",
          "en": "Hurl a lance of sea-green darkness at an enemy within 8 squares. On a hit it takes d8 damage."
        },
        "cost": 1
      },
      {
        "id": "mark-of-debt",
        "name": {
          "th": "ตราหนี้",
          "en": "Mark of Debt"
        },
        "desc": {
          "th": "ประทับตราหนี้ลงบนศัตรู 1 ตัวในระยะ 6 ช่อง นาน 3 ตา ทุกครั้งที่ใครในกลุ่มโจมตีมันโดน มันเสียเพิ่มอีก 2 แต้ม",
          "en": "Brand one enemy within 6 squares with a debt mark for 3 turns: every time anyone in the party hits it, it takes 2 extra damage."
        },
        "cost": 2
      },
      {
        "id": "grasp-of-the-deep",
        "name": {
          "th": "หนวดจากห้วงลึก",
          "en": "Grasp of the Deep"
        },
        "desc": {
          "th": "หนวดยักษ์โผล่จากพื้นในพื้นที่ 2×2 ช่องในระยะ 6 ช่อง ศัตรูในพื้นที่เดินไม่ได้ 2 ตา และเสียหาย d4 ทุกต้นตาของมัน โดยไม่ต้องทอยโดน",
          "en": "Giant tentacles burst from a 2×2 area within 6 squares. Enemies inside can't move for 2 turns and take d4 damage at the start of each of their turns, no hit roll needed."
        },
        "cost": 4
      }
    ],
    "items": [
      {
        "name": {
          "th": "จี้ปะการังดำ",
          "en": "Black Coral Pendant"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "กระเป๋าหมอยา",
          "en": "Apothecary's Satchel"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "มีดผ่าแผลเงิน",
          "en": "Silver Lancet"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ยาแก้ไข้",
          "en": "Fever Draught"
        },
        "qty": 3
      },
      {
        "name": {
          "th": "จดหมายจากลูกสาว",
          "en": "Letter from Her Daughter"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-wizard",
    "set": "legend",
    "heroName": {
      "th": "ธีโอดริก เพลล์",
      "en": "Theodric Pell"
    },
    "className": {
      "th": "พ่อมด",
      "en": "Wizard"
    },
    "blurb": {
      "th": "บรรณารักษ์ที่แอบเรียนเวทจากหนังสือต้องห้ามในห้องใต้ดินของหอสมุด และถูกไล่ออกตอนหนังสือเล่มหนึ่งกัดคนอ่าน",
      "en": "A librarian who secretly learned magic from the forbidden books in the archive vault, and was dismissed when one of them bit a reader."
    },
    "role": "control",
    "art": "",
    "stats": {
      "hpMax": 16,
      "str": 2,
      "end": 4,
      "agi": 6,
      "mpMax": 20
    },
    "skills": [
      {
        "id": "ink-dart",
        "name": {
          "th": "ลูกดอกหมึก",
          "en": "Ink Dart"
        },
        "desc": {
          "th": "ยิงหมึกเวทใส่ศัตรูในระยะ 6 ช่อง ถ้าโดน เสียหาย d6 และหมึกเข้าตา ทอยโจมตีครั้งถัดไปของมัน −2",
          "en": "Shoot magic ink at an enemy within 6 squares. On a hit it takes d6 damage and the ink blurs its eyes: its next attack roll gets −2."
        },
        "cost": 1
      },
      {
        "id": "paper-ward",
        "name": {
          "th": "กระดาษคุ้มภัย",
          "en": "Paper Ward"
        },
        "desc": {
          "th": "หน้ากระดาษปลิวมาบังเพื่อน 1 คนในระยะ 5 ช่อง (หรือตัวเรา) ครั้งถัดไปที่เขาโดนตี ดาเมจลดลง d8 แต้ม",
          "en": "Pages fly up to shield one ally within 5 squares (or yourself): the next hit on them deals d8 less damage."
        },
        "cost": 2
      },
      {
        "id": "binding-glyph",
        "name": {
          "th": "อักขระตรึง",
          "en": "Binding Glyph"
        },
        "desc": {
          "th": "เขียนอักขระลงพื้นที่ 3×3 ช่องในระยะ 6 ช่อง ศัตรูที่อยู่หรือเดินเข้าไปเดินไม่ได้ 2 ตา (แต่ยังโจมตีได้) อักขระอยู่นาน 3 รอบ",
          "en": "Inscribe a glyph over a 3×3 area within 6 squares. Enemies in it, or who step in, can't move for 2 turns (they can still attack). The glyph lasts 3 rounds."
        },
        "cost": 5
      }
    ],
    "items": [
      {
        "name": {
          "th": "ตำราเวทจอมกัด",
          "en": "The Biting Grimoire"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ไม้เท้าไม้มะเกลือ",
          "en": "Ebony Staff"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "แว่นขยายทองเหลือง",
          "en": "Brass Magnifying Glass"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "บัตรยืมหนังสือหมดอายุ",
          "en": "Expired Library Card"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ชาเขียวเข้มข้น",
          "en": "Strong Green Tea"
        },
        "qty": 2
      }
    ]
  }
]
```
