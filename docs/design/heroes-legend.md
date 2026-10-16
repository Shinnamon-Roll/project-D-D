# วีรชนในตำนาน / Heroes of Legend — ข้อมูลเกม v2 (12 ตัว)

ข้อมูลฮีโร่ 12 ตัวของเซ็ต `legend` (วีรชนในตำนาน) แปลงจาก [`roster-legend.md`](roster-legend.md) เป็นระบบ v2.
id = `legend-<คลาส>`, `classId` = คลาสหลัง `legend-`. อีกเซ็ตอยู่ใน [`heroes-core.md`](heroes-core.md).
กติกา v2 ([`../features/rules-v2.md`](../features/rules-v2.md)): 6 ability, AC, skill proficiency ตามคลาส ([`classes/`](classes/)).
Ability = standard array 15/14/13/12/10/8 จัดให้เข้ากับคลาสและเรื่องของฮีโร่ แล้วบวกโบนัส +2/+1 (สูงสุด 17 ที่เลเวล 1, ผลรวม 75).
AC คำนวณจากเกราะ/โล่ในของเริ่มต้น (ตัวเลขเกราะตาม SRD 5.2.1, CC BY 4.0); HP เลเวล 1 = hit die เต็ม + CON mod. DM แก้ได้ในแอป.
Expertise: Rogue ได้ 2 skill ตั้งแต่เลเวล 1; Bard/Ranger ได้ทีหลัง (DM/ผู้เล่นเพิ่มเอง).
ชื่อ เรื่อง (`heroName`, `blurb`, `role`) คงไว้ตามไฟล์เดิม; ของชิ้นแรก = ของประจำตัวจากไฟล์เดิม. ชื่อ เนื้อเรื่อง และไอเทมแต่งขึ้นเอง

## ตารางสรุป
abilities = STR / DEX / CON / INT / WIS / CHA

| # | ฮีโร่ | คลาส | abilities | AC | HP L1 | skills |
|---|---|---|---|---|---|---|
| 1 | Kaya Stormhorn (คายา สตอร์มฮอร์น) | Barbarian (นักรบเถื่อน) | 16 / 12 / 16 / 8 / 13 / 10 | 14 | 15 | athletics, intimidation |
| 2 | Odalys Reed (โอดาลิส รีด) | Bard (นักกวีเวทดนตรี) | 8 / 14 / 12 / 14 / 10 / 17 | 13 | 9 | history, insight, performance |
| 3 | Oswin Graylight (ออสวิน เกรย์ไลต์) | Cleric (นักบวช) | 15 / 10 / 13 / 8 / 17 / 12 | 16 | 9 | insight, religion |
| 4 | Tansy Hollowfern (แทนซี ฮอลโลว์เฟิร์น) | Druid (ดรูอิด) | 8 / 12 / 13 / 15 / 17 / 10 | 12 | 9 | medicine, nature |
| 5 | Petra Kalvos (เพตรา คาลวอส) | Fighter (นักรบ) | 16 / 12 / 16 / 10 / 8 / 13 | 18 | 13 | athletics, persuasion |
| 6 | Kiran Stillwater (คิรัน สติลวอเตอร์) | Monk (นักพรตยุทธ์) | 10 / 17 / 13 / 12 / 15 / 8 | 15 | 9 | history, insight |
| 7 | Tobias Vellan (โทเบียส เวลลัน) | Paladin (อัศวินศักดิ์สิทธิ์) | 16 / 10 / 16 / 8 / 12 / 13 | 16 | 13 | insight, religion |
| 8 | Joren Vask (โยเรน วาสก์) | Ranger (พรานป่า) | 12 / 17 / 14 / 10 / 14 / 8 | 15 | 12 | athletics, nature, survival |
| 9 | Silas Wren (ไซลัส เรน) | Rogue (นักย่องเบา) | 8 / 17 / 10 / 14 / 13 / 13 | 14 | 8 | deception, insight, investigation, stealth |
| 10 | Remy Saltfire (เรมี ซอลต์ไฟร์) | Sorcerer (นักเวทโดยกำเนิด) | 12 / 13 / 15 / 8 / 10 / 17 | 11 | 8 | insight, intimidation |
| 11 | Vesna Morrow (เวสนา มอร์โรว์) | Warlock (ผู้ทำสัญญา) | 8 / 12 / 14 / 14 / 10 / 17 | 12 | 10 | arcana, nature |
| 12 | Theodric Pell (ธีโอดริก เพลล์) | Wizard (พ่อมด/นักเวทศึกษา) | 8 / 13 / 12 / 17 / 15 / 10 | 11 | 7 | arcana, investigation |

## 1. นักรบเถื่อน / Barbarian — `legend-barbarian`
**คายา สตอร์มฮอร์น** / **Kaya Stormhorn** · คลาส `barbarian` · role: `damage` (ดาเมจ)

> ลูกสาวหัวหน้าเผ่าบนภูเขาหิมะ ที่ถูกเนรเทศเพราะพังประตูเมืองด้วยมือเปล่า และยังภูมิใจกับมันอยู่  
> *A snow-mountain chieftain's daughter exiled for breaking down a city gate with her bare hands, and still proud of it.*

- **Abilities**: STR 16 (+3) DEX 12 (+1) CON 16 (+3) INT 8 (-1) WIS 13 (+1) CHA 10 (+0)
  - standard array + โบนัส: STR 15+1, CON 14+2, WIS 13, DEX 12, CHA 10, INT 8
- **AC 14**: Unarmored Defense 10 + DEX + CON
- **HP เลเวล 1**: 15 (d12 เต็ม 12 + CON +3)
- **Skills**: athletics, intimidation
- **ของ**: ค้อนศึกหินยอดเขา / Peakstone War Maul · ขวานมือ / Handaxe ×4 · เสื้อคลุมขนหมี / Bearskin Cloak · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 2. นักกวีเวทดนตรี / Bard — `legend-bard`
**โอดาลิส รีด** / **Odalys Reed** · คลาส `bard` · role: `support` (สนับสนุน)

> นักเขียนบทละครเร่ที่ออกเดินทางเก็บเรื่องจริง เพราะเชื่อว่าโลกแต่งพล็อตได้ดีกว่าเธอเสมอ  
> *A travelling playwright who collects true stories, because she believes the world always writes better plots than she does.*

- **Abilities**: STR 8 (-1) DEX 14 (+2) CON 12 (+1) INT 14 (+2) WIS 10 (+0) CHA 17 (+3)
  - standard array + โบนัส: CHA 15+2, DEX 14, INT 13+1, CON 12, WIS 10, STR 8
- **AC 13**: Leather Armor 11 + DEX +2
- **HP เลเวล 1**: 9 (d8 เต็ม 8 + CON +1)
- **Skills**: history, insight, performance
- **ของ**: ซอไม้เมเปิลแดง / Red Maple Fiddle · เกราะหนัง / Leather Armor · มีดสั้น / Dagger ×2 · สมุดบทละคร / Playwright's Notebook · ชุดนักแสดง / Entertainer's Pack (ชิ้นแรก = ของประจำตัว)

## 3. นักบวช / Cleric — `legend-cleric`
**ออสวิน เกรย์ไลต์** / **Oswin Graylight** · คลาส `cleric` · role: `healer` (ฮีลเลอร์)

> อดีตทหารที่หันมาบวชหลังรอดตายจากศึกที่ไม่ควรมีใครรอด และยังถามเทพทุกคืนว่าทำไมต้องเป็นเขา  
> *A former soldier who took holy vows after surviving a battle no one should have survived, and still asks his god every night: why me?*

- **Abilities**: STR 15 (+2) DEX 10 (+0) CON 13 (+1) INT 8 (-1) WIS 17 (+3) CHA 12 (+1)
  - standard array + โบนัส: WIS 15+2, STR 14+1, CON 13, CHA 12, DEX 10, INT 8
- **AC 16**: Scale Mail 14 + DEX +0 + Shield 2
- **HP เลเวล 1**: 9 (d8 เต็ม 8 + CON +1)
- **Skills**: insight, religion
- **ของ**: ตะเกียงเงินแห่งวิหาร / Temple Silver Lantern · ชุดเกราะเกล็ด / Scale Mail · โล่ / Shield · กระบอง / Mace · สัญลักษณ์ศักดิ์สิทธิ์ / Holy Symbol · ชุดนักบวช / Priest's Pack (ชิ้นแรก = ของประจำตัว)

## 4. ดรูอิด / Druid — `legend-druid`
**แทนซี ฮอลโลว์เฟิร์น** / **Tansy Hollowfern** · คลาส `druid` · role: `hybrid` (ไฮบริด)

> หญิงชราหมอยาชายป่าที่เด็กๆ ในหมู่บ้านลือว่าเป็นแม่มด และเธอก็ไม่เคยคิดจะแก้ข่าว  
> *An old hedge-herbalist the village children swear is a witch, and she has never bothered to correct them.*

- **Abilities**: STR 8 (-1) DEX 12 (+1) CON 13 (+1) INT 15 (+2) WIS 17 (+3) CHA 10 (+0)
  - standard array + โบนัส: WIS 15+2, INT 14+1, CON 13, DEX 12, CHA 10, STR 8
- **AC 12**: Leather Armor 11 + DEX +1
- **HP เลเวล 1**: 9 (d8 เต็ม 8 + CON +1)
- **Skills**: medicine, nature
- **ของ**: ไม้เท้าเอลเดอร์บิดเกลียว / Twisted Elder Staff · เกราะหนัง / Leather Armor · เคียว / Sickle · ชุดสมุนไพร / Herbalism Kit · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 5. นักรบ / Fighter — `legend-fighter`
**เพตรา คาลวอส** / **Petra Kalvos** · คลาส `fighter` · role: `tank` (แทงก์)

> ลูกสาวช่างตีเหล็กที่ตั้งใจพิสูจน์ว่าดาบของพ่อดีที่สุดในแผ่นดิน ด้วยการถือมันลงสนามประลองเองทุกสนาม  
> *A blacksmith's daughter set on proving her father's blades are the finest in the land, by carrying them into every arena herself.*

- **Abilities**: STR 16 (+3) DEX 12 (+1) CON 16 (+3) INT 10 (+0) WIS 8 (-1) CHA 13 (+1)
  - standard array + โบนัส: STR 15+1, CON 14+2, CHA 13, DEX 12, INT 10, WIS 8
- **AC 18**: Chain Mail 16 + Shield 2
- **HP เลเวล 1**: 13 (d10 เต็ม 10 + CON +3)
- **Skills**: athletics, persuasion
- **ของ**: ดาบยาวตราค้อนทั่ง / Anvil-Mark Longsword · ชุดเกราะโซ่ / Chain Mail · โล่ / Shield · เครื่องมือช่างตีเหล็ก / Smith's Tools · ชุดสำรวจดันเจียน / Dungeoneer's Pack (ชิ้นแรก = ของประจำตัว)

## 6. นักพรตยุทธ์ / Monk — `legend-monk`
**คิรัน สติลวอเตอร์** / **Kiran Stillwater** · คลาส `monk` · role: `control` (ควบคุม)

> ศิษย์คนสุดท้ายของวัดบนหน้าผาที่ถูกเผา ออกตามหาคัมภีร์ที่ถูกขโมย โดยพกไปแค่ชาดีกับความอดทน  
> *The last disciple of a burned cliffside monastery, tracking its stolen scrolls with nothing but good tea and patience.*

- **Abilities**: STR 10 (+0) DEX 17 (+3) CON 13 (+1) INT 12 (+1) WIS 15 (+2) CHA 8 (-1)
  - standard array + โบนัส: DEX 15+2, WIS 14+1, CON 13, INT 12, STR 10, CHA 8
- **AC 15**: Unarmored Defense 10 + DEX + WIS
- **HP เลเวล 1**: 9 (d8 เต็ม 8 + CON +1)
- **Skills**: history, insight
- **ของ**: ไม้พลองไผ่ดำ / Black Bamboo Staff · มีดสั้น / Dagger ×5 · ชุดชงชา / Tea Set · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 7. อัศวินศักดิ์สิทธิ์ / Paladin — `legend-paladin`
**โทเบียส เวลลัน** / **Tobias Vellan** · คลาส `paladin` · role: `hybrid` (ไฮบริด)

> ชาวนาธรรมดาที่ได้ยินเสียงเรียกในฝันคืนที่หมู่บ้านถูกเผา และยังไม่แน่ใจว่าตัวเองคู่ควรกับคำสาบานที่รับมา  
> *A plain farmer who heard a calling in his dreams the night his village burned, and still isn't sure he deserves the oath he took.*

- **Abilities**: STR 16 (+3) DEX 10 (+0) CON 16 (+3) INT 8 (-1) WIS 12 (+1) CHA 13 (+1)
  - standard array + โบนัส: STR 15+1, CON 14+2, CHA 13, WIS 12, DEX 10, INT 8
- **AC 16**: Chain Mail 16
- **HP เลเวล 1**: 13 (d10 เต็ม 10 + CON +3)
- **Skills**: insight, religion
- **ของ**: ง้าวตีใหม่จากเคียวไร่ / Reforged Harvest Glaive · ชุดเกราะโซ่ / Chain Mail · หอกซัด / Javelin ×5 · สัญลักษณ์ศักดิ์สิทธิ์ / Holy Symbol · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 8. พรานป่า / Ranger — `legend-ranger`
**โยเรน วาสก์** / **Joren Vask** · คลาส `ranger` · role: `damage` (ดาเมจ)

> คนนำทางกองคาราวานที่รู้ทุกทางลัดบนภูเขา และร้อยเขี้ยวสัตว์ร้ายทุกตัวที่ล่าได้เป็นสร้อยคอ  
> *A caravan guide who knows every mountain shortcut and strings a fang from every beast he has hunted onto his necklace.*

- **Abilities**: STR 12 (+1) DEX 17 (+3) CON 14 (+2) INT 10 (+0) WIS 14 (+2) CHA 8 (-1)
  - standard array + โบนัส: DEX 15+2, CON 14, WIS 13+1, STR 12, INT 10, CHA 8
- **AC 15**: Studded Leather Armor 12 + DEX +3
- **HP เลเวล 1**: 12 (d10 เต็ม 10 + CON +2)
- **Skills**: athletics, nature, survival
- **ของ**: ธนูสั้นเขาแพะภูเขา / Ram-Horn Shortbow · เกราะหนังหมุด / Studded Leather Armor · ลูกธนู / Arrows ×20 · ดาบโค้ง / Scimitar · สร้อยเขี้ยวสัตว์ / Fang Necklace · ชุดเดินทาง / Explorer's Pack (ชิ้นแรก = ของประจำตัว)

## 9. นักย่องเบา / Rogue — `legend-rogue`
**ไซลัส เรน** / **Silas Wren** · คลาส `rogue` · role: `damage` (ดาเมจ)

> สายลับของสภาเมืองที่แฝงตัวเข้ากลุ่มนักผจญภัย และยังตัดสินใจไม่ได้ว่าจะส่งรายงานหรือเลือกเพื่อนใหม่  
> *A city council spy planted among adventurers, who still can't decide whether to file his report or side with his new friends.*

- **Abilities**: STR 8 (-1) DEX 17 (+3) CON 10 (+0) INT 14 (+2) WIS 13 (+1) CHA 13 (+1)
  - standard array + โบนัส: DEX 15+2, INT 14, CHA 13, WIS 12+1, CON 10, STR 8
- **AC 14**: Leather Armor 11 + DEX +3
- **HP เลเวล 1**: 8 (d8 เต็ม 8 + CON +0)
- **Skills**: deception, insight, investigation, stealth · **Expertise**: deception, investigation
- **ของ**: มีดสั้นใบดำซ่อนแขนเสื้อ / Black Sleeve Dagger · เกราะหนัง / Leather Armor · ดาบสั้น / Shortsword · เครื่องมือโจร / Thieves' Tools · ชุดปลอมตัว / Disguise Kit · ชุดย่องเบา / Burglar's Pack (ชิ้นแรก = ของประจำตัว)

## 10. นักเวทโดยกำเนิด / Sorcerer — `legend-sorcerer`
**เรมี ซอลต์ไฟร์** / **Remy Saltfire** · คลาส `sorcerer` · role: `damage` (ดาเมจ)

> ลูกชาวประมงที่เผาเรือของพ่อตอนโกรธครั้งแรกจนต้องหนีออกจากหมู่บ้าน ตอนนี้กำลังหัดโกรธให้ถูกเวลา  
> *A fisher's son who set his father's boat ablaze the first time he lost his temper and had to flee the village. Now he is learning to get angry at the right moment.*

- **Abilities**: STR 12 (+1) DEX 13 (+1) CON 15 (+2) INT 8 (-1) WIS 10 (+0) CHA 17 (+3)
  - standard array + โบนัส: CHA 15+2, CON 14+1, DEX 13, STR 12, WIS 10, INT 8
- **AC 11**: ไม่ใส่เกราะ 10 + DEX
- **HP เลเวล 1**: 8 (d6 เต็ม 6 + CON +2)
- **Skills**: insight, intimidation
- **ของ**: จี้เปลือกหอยเกลือ / Salt-Shell Pendant · หอก / Spear · มีดสั้น / Dagger ×2 · อุปกรณ์ตกปลา / Fishing Tackle · ชุดสำรวจดันเจียน / Dungeoneer's Pack (ชิ้นแรก = ของประจำตัว)

## 11. ผู้ทำสัญญา / Warlock — `legend-warlock`
**เวสนา มอร์โรว์** / **Vesna Morrow** · คลาส `warlock` · role: `hybrid` (ไฮบริด)

> หมอยาที่ทำสัญญากับบางสิ่งจากใต้ทะเลลึกเพื่อรักษาลูกสาว และตอนนี้มันเริ่มกระซิบทวงค่าตอบแทนทุกคืน  
> *An apothecary who struck a pact with something from the deep sea to save her daughter, and now it whispers every night for its payment.*

- **Abilities**: STR 8 (-1) DEX 12 (+1) CON 14 (+2) INT 14 (+2) WIS 10 (+0) CHA 17 (+3)
  - standard array + โบนัส: CHA 15+2, CON 14, INT 13+1, DEX 12, WIS 10, STR 8
- **AC 12**: Leather Armor 11 + DEX +1
- **HP เลเวล 1**: 10 (d8 เต็ม 8 + CON +2)
- **Skills**: arcana, nature
- **ของ**: จี้ปะการังดำ / Black Coral Pendant · เกราะหนัง / Leather Armor · เคียว / Sickle · มีดสั้น / Dagger ×2 · ชุดสมุนไพร / Herbalism Kit · ชุดนักปราชญ์ / Scholar's Pack (ชิ้นแรก = ของประจำตัว)

## 12. พ่อมด/นักเวทศึกษา / Wizard — `legend-wizard`
**ธีโอดริก เพลล์** / **Theodric Pell** · คลาส `wizard` · role: `control` (ควบคุม)

> บรรณารักษ์ที่แอบเรียนเวทจากหนังสือต้องห้ามในห้องใต้ดินของหอสมุด และถูกไล่ออกตอนหนังสือเล่มหนึ่งกัดคนอ่าน  
> *A librarian who secretly learned magic from the forbidden books in the archive vault, and was dismissed when one of them bit a reader.*

- **Abilities**: STR 8 (-1) DEX 13 (+1) CON 12 (+1) INT 17 (+3) WIS 15 (+2) CHA 10 (+0)
  - standard array + โบนัส: INT 15+2, WIS 14+1, DEX 13, CON 12, CHA 10, STR 8
- **AC 11**: ไม่ใส่เกราะ 10 + DEX
- **HP เลเวล 1**: 7 (d6 เต็ม 6 + CON +1)
- **Skills**: arcana, investigation
- **ของ**: ตำราเวทจอมกัด / The Biting Grimoire · ไม้พลอง / Quarterstaff · มีดสั้น / Dagger ×2 · แว่นอ่านหนังสือ / Reading Spectacles · ชุดนักปราชญ์ / Scholar's Pack (ชิ้นแรก = ของประจำตัว)

## JSON (machine-readable)

```json
[
  {
    "id": "legend-barbarian",
    "set": "legend",
    "classId": "barbarian",
    "heroName": {
      "th": "คายา สตอร์มฮอร์น",
      "en": "Kaya Stormhorn"
    },
    "blurb": {
      "th": "ลูกสาวหัวหน้าเผ่าบนภูเขาหิมะ ที่ถูกเนรเทศเพราะพังประตูเมืองด้วยมือเปล่า และยังภูมิใจกับมันอยู่",
      "en": "A snow-mountain chieftain's daughter exiled for breaking down a city gate with her bare hands, and still proud of it."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 16,
      "dex": 12,
      "con": 16,
      "int": 8,
      "wis": 13,
      "cha": 10
    },
    "ac": 14,
    "skills": [
      "athletics",
      "intimidation"
    ],
    "expertise": [],
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
          "th": "ขวานมือ",
          "en": "Handaxe"
        },
        "qty": 4
      },
      {
        "name": {
          "th": "เสื้อคลุมขนหมี",
          "en": "Bearskin Cloak"
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
    "id": "legend-bard",
    "set": "legend",
    "classId": "bard",
    "heroName": {
      "th": "โอดาลิส รีด",
      "en": "Odalys Reed"
    },
    "blurb": {
      "th": "นักเขียนบทละครเร่ที่ออกเดินทางเก็บเรื่องจริง เพราะเชื่อว่าโลกแต่งพล็อตได้ดีกว่าเธอเสมอ",
      "en": "A travelling playwright who collects true stories, because she believes the world always writes better plots than she does."
    },
    "role": "support",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 14,
      "con": 12,
      "int": 14,
      "wis": 10,
      "cha": 17
    },
    "ac": 13,
    "skills": [
      "history",
      "insight",
      "performance"
    ],
    "expertise": [],
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
          "th": "สมุดบทละคร",
          "en": "Playwright's Notebook"
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
    "id": "legend-cleric",
    "set": "legend",
    "classId": "cleric",
    "heroName": {
      "th": "ออสวิน เกรย์ไลต์",
      "en": "Oswin Graylight"
    },
    "blurb": {
      "th": "อดีตทหารที่หันมาบวชหลังรอดตายจากศึกที่ไม่ควรมีใครรอด และยังถามเทพทุกคืนว่าทำไมต้องเป็นเขา",
      "en": "A former soldier who took holy vows after surviving a battle no one should have survived, and still asks his god every night: why me?"
    },
    "role": "healer",
    "art": "",
    "abilities": {
      "str": 15,
      "dex": 10,
      "con": 13,
      "int": 8,
      "wis": 17,
      "cha": 12
    },
    "ac": 16,
    "skills": [
      "insight",
      "religion"
    ],
    "expertise": [],
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
          "th": "ชุดเกราะเกล็ด",
          "en": "Scale Mail"
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
    "id": "legend-druid",
    "set": "legend",
    "classId": "druid",
    "heroName": {
      "th": "แทนซี ฮอลโลว์เฟิร์น",
      "en": "Tansy Hollowfern"
    },
    "blurb": {
      "th": "หญิงชราหมอยาชายป่าที่เด็กๆ ในหมู่บ้านลือว่าเป็นแม่มด และเธอก็ไม่เคยคิดจะแก้ข่าว",
      "en": "An old hedge-herbalist the village children swear is a witch, and she has never bothered to correct them."
    },
    "role": "hybrid",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 12,
      "con": 13,
      "int": 15,
      "wis": 17,
      "cha": 10
    },
    "ac": 12,
    "skills": [
      "medicine",
      "nature"
    ],
    "expertise": [],
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
          "th": "เกราะหนัง",
          "en": "Leather Armor"
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
    "id": "legend-fighter",
    "set": "legend",
    "classId": "fighter",
    "heroName": {
      "th": "เพตรา คาลวอส",
      "en": "Petra Kalvos"
    },
    "blurb": {
      "th": "ลูกสาวช่างตีเหล็กที่ตั้งใจพิสูจน์ว่าดาบของพ่อดีที่สุดในแผ่นดิน ด้วยการถือมันลงสนามประลองเองทุกสนาม",
      "en": "A blacksmith's daughter set on proving her father's blades are the finest in the land, by carrying them into every arena herself."
    },
    "role": "tank",
    "art": "",
    "abilities": {
      "str": 16,
      "dex": 12,
      "con": 16,
      "int": 10,
      "wis": 8,
      "cha": 13
    },
    "ac": 18,
    "skills": [
      "athletics",
      "persuasion"
    ],
    "expertise": [],
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
          "th": "ชุดเกราะโซ่",
          "en": "Chain Mail"
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
          "th": "เครื่องมือช่างตีเหล็ก",
          "en": "Smith's Tools"
        },
        "qty": 1
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
    "id": "legend-monk",
    "set": "legend",
    "classId": "monk",
    "heroName": {
      "th": "คิรัน สติลวอเตอร์",
      "en": "Kiran Stillwater"
    },
    "blurb": {
      "th": "ศิษย์คนสุดท้ายของวัดบนหน้าผาที่ถูกเผา ออกตามหาคัมภีร์ที่ถูกขโมย โดยพกไปแค่ชาดีกับความอดทน",
      "en": "The last disciple of a burned cliffside monastery, tracking its stolen scrolls with nothing but good tea and patience."
    },
    "role": "control",
    "art": "",
    "abilities": {
      "str": 10,
      "dex": 17,
      "con": 13,
      "int": 12,
      "wis": 15,
      "cha": 8
    },
    "ac": 15,
    "skills": [
      "history",
      "insight"
    ],
    "expertise": [],
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
          "th": "มีดสั้น",
          "en": "Dagger"
        },
        "qty": 5
      },
      {
        "name": {
          "th": "ชุดชงชา",
          "en": "Tea Set"
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
    "id": "legend-paladin",
    "set": "legend",
    "classId": "paladin",
    "heroName": {
      "th": "โทเบียส เวลลัน",
      "en": "Tobias Vellan"
    },
    "blurb": {
      "th": "ชาวนาธรรมดาที่ได้ยินเสียงเรียกในฝันคืนที่หมู่บ้านถูกเผา และยังไม่แน่ใจว่าตัวเองคู่ควรกับคำสาบานที่รับมา",
      "en": "A plain farmer who heard a calling in his dreams the night his village burned, and still isn't sure he deserves the oath he took."
    },
    "role": "hybrid",
    "art": "",
    "abilities": {
      "str": 16,
      "dex": 10,
      "con": 16,
      "int": 8,
      "wis": 12,
      "cha": 13
    },
    "ac": 16,
    "skills": [
      "insight",
      "religion"
    ],
    "expertise": [],
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
          "th": "ชุดเดินทาง",
          "en": "Explorer's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-ranger",
    "set": "legend",
    "classId": "ranger",
    "heroName": {
      "th": "โยเรน วาสก์",
      "en": "Joren Vask"
    },
    "blurb": {
      "th": "คนนำทางกองคาราวานที่รู้ทุกทางลัดบนภูเขา และร้อยเขี้ยวสัตว์ร้ายทุกตัวที่ล่าได้เป็นสร้อยคอ",
      "en": "A caravan guide who knows every mountain shortcut and strings a fang from every beast he has hunted onto his necklace."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 12,
      "dex": 17,
      "con": 14,
      "int": 10,
      "wis": 14,
      "cha": 8
    },
    "ac": 15,
    "skills": [
      "athletics",
      "nature",
      "survival"
    ],
    "expertise": [],
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
          "th": "ดาบโค้ง",
          "en": "Scimitar"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "สร้อยเขี้ยวสัตว์",
          "en": "Fang Necklace"
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
    "id": "legend-rogue",
    "set": "legend",
    "classId": "rogue",
    "heroName": {
      "th": "ไซลัส เรน",
      "en": "Silas Wren"
    },
    "blurb": {
      "th": "สายลับของสภาเมืองที่แฝงตัวเข้ากลุ่มนักผจญภัย และยังตัดสินใจไม่ได้ว่าจะส่งรายงานหรือเลือกเพื่อนใหม่",
      "en": "A city council spy planted among adventurers, who still can't decide whether to file his report or side with his new friends."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 17,
      "con": 10,
      "int": 14,
      "wis": 13,
      "cha": 13
    },
    "ac": 14,
    "skills": [
      "deception",
      "insight",
      "investigation",
      "stealth"
    ],
    "expertise": [
      "deception",
      "investigation"
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
          "th": "เกราะหนัง",
          "en": "Leather Armor"
        },
        "qty": 1
      },
      {
        "name": {
          "th": "ดาบสั้น",
          "en": "Shortsword"
        },
        "qty": 1
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
          "th": "ชุดปลอมตัว",
          "en": "Disguise Kit"
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
    "id": "legend-sorcerer",
    "set": "legend",
    "classId": "sorcerer",
    "heroName": {
      "th": "เรมี ซอลต์ไฟร์",
      "en": "Remy Saltfire"
    },
    "blurb": {
      "th": "ลูกชาวประมงที่เผาเรือของพ่อตอนโกรธครั้งแรกจนต้องหนีออกจากหมู่บ้าน ตอนนี้กำลังหัดโกรธให้ถูกเวลา",
      "en": "A fisher's son who set his father's boat ablaze the first time he lost his temper and had to flee the village. Now he is learning to get angry at the right moment."
    },
    "role": "damage",
    "art": "",
    "abilities": {
      "str": 12,
      "dex": 13,
      "con": 15,
      "int": 8,
      "wis": 10,
      "cha": 17
    },
    "ac": 11,
    "skills": [
      "insight",
      "intimidation"
    ],
    "expertise": [],
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
          "th": "หอก",
          "en": "Spear"
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
          "th": "อุปกรณ์ตกปลา",
          "en": "Fishing Tackle"
        },
        "qty": 1
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
    "id": "legend-warlock",
    "set": "legend",
    "classId": "warlock",
    "heroName": {
      "th": "เวสนา มอร์โรว์",
      "en": "Vesna Morrow"
    },
    "blurb": {
      "th": "หมอยาที่ทำสัญญากับบางสิ่งจากใต้ทะเลลึกเพื่อรักษาลูกสาว และตอนนี้มันเริ่มกระซิบทวงค่าตอบแทนทุกคืน",
      "en": "An apothecary who struck a pact with something from the deep sea to save her daughter, and now it whispers every night for its payment."
    },
    "role": "hybrid",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 12,
      "con": 14,
      "int": 14,
      "wis": 10,
      "cha": 17
    },
    "ac": 12,
    "skills": [
      "arcana",
      "nature"
    ],
    "expertise": [],
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
          "th": "เกราะหนัง",
          "en": "Leather Armor"
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
          "th": "มีดสั้น",
          "en": "Dagger"
        },
        "qty": 2
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
          "th": "ชุดนักปราชญ์",
          "en": "Scholar's Pack"
        },
        "qty": 1
      }
    ]
  },
  {
    "id": "legend-wizard",
    "set": "legend",
    "classId": "wizard",
    "heroName": {
      "th": "ธีโอดริก เพลล์",
      "en": "Theodric Pell"
    },
    "blurb": {
      "th": "บรรณารักษ์ที่แอบเรียนเวทจากหนังสือต้องห้ามในห้องใต้ดินของหอสมุด และถูกไล่ออกตอนหนังสือเล่มหนึ่งกัดคนอ่าน",
      "en": "A librarian who secretly learned magic from the forbidden books in the archive vault, and was dismissed when one of them bit a reader."
    },
    "role": "control",
    "art": "",
    "abilities": {
      "str": 8,
      "dex": 13,
      "con": 12,
      "int": 17,
      "wis": 15,
      "cha": 10
    },
    "ac": 11,
    "skills": [
      "arcana",
      "investigation"
    ],
    "expertise": [],
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
          "th": "ไม้พลอง",
          "en": "Quarterstaff"
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
          "th": "แว่นอ่านหนังสือ",
          "en": "Reading Spectacles"
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
