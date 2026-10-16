# Roster — วีรชนในตำนาน / Heroes of Legend (12 คลาส)

> ไฟล์นี้เก็บแค่ชื่อ บทบาท และเนื้อเรื่อง — ข้อมูลเกม (ability, AC, skills, ของเริ่มต้น, JSON ที่ sync กับ Go) อยู่ที่ [`heroes-legend.md`](heroes-legend.md) (กติกา: [`../features/rules-v2.md`](../features/rules-v2.md))

ชุดตัวละครชุดที่ 2 (`set: "legend"`): 12 ตัว 1 ตัวต่อ 1 คลาส.
รูปแบบเดียวกับ [`roster.md`](roster.md); `id` = `legend-<คลาส>`, `art` = `""` (ใช้ไอคอนคลาสใน `web/src/assets/classes/`).
ชื่อและเนื้อเรื่องทั้งหมดแต่งขึ้นเอง ไม่ได้คัดลอกจากหนังสือหรือเกมใด; ชื่อคลาสเป็นคำทั่วไป แนวคลาสอ้างอิง [`../research/classes/phb-2024/`](../research/classes/phb-2024/).

## 1. นักรบเถื่อน / Barbarian — `legend-barbarian`
**คายา สตอร์มฮอร์น** / **Kaya Stormhorn** · role: `damage` (ดาเมจ)

> ลูกสาวหัวหน้าเผ่าบนภูเขาหิมะ ที่ถูกเนรเทศเพราะพังประตูเมืองด้วยมือเปล่า และยังภูมิใจกับมันอยู่  
> *A snow-mountain chieftain's daughter exiled for breaking down a city gate with her bare hands, and still proud of it.*

## 2. นักกวีเวทดนตรี / Bard — `legend-bard`
**โอดาลิส รีด** / **Odalys Reed** · role: `support` (สนับสนุน)

> นักเขียนบทละครเร่ที่ออกเดินทางเก็บเรื่องจริง เพราะเชื่อว่าโลกแต่งพล็อตได้ดีกว่าเธอเสมอ  
> *A travelling playwright who collects true stories, because she believes the world always writes better plots than she does.*

## 3. นักบวช / Cleric — `legend-cleric`
**ออสวิน เกรย์ไลต์** / **Oswin Graylight** · role: `healer` (ฮีลเลอร์)

> อดีตทหารที่หันมาบวชหลังรอดตายจากศึกที่ไม่ควรมีใครรอด และยังถามเทพทุกคืนว่าทำไมต้องเป็นเขา  
> *A former soldier who took holy vows after surviving a battle no one should have survived, and still asks his god every night: why me?*

## 4. ดรูอิด / Druid — `legend-druid`
**แทนซี ฮอลโลว์เฟิร์น** / **Tansy Hollowfern** · role: `hybrid` (ไฮบริด)

> หญิงชราหมอยาชายป่าที่เด็กๆ ในหมู่บ้านลือว่าเป็นแม่มด และเธอก็ไม่เคยคิดจะแก้ข่าว  
> *An old hedge-herbalist the village children swear is a witch, and she has never bothered to correct them.*

## 5. นักรบ / Fighter — `legend-fighter`
**เพตรา คาลวอส** / **Petra Kalvos** · role: `tank` (แทงก์)

> ลูกสาวช่างตีเหล็กที่ตั้งใจพิสูจน์ว่าดาบของพ่อดีที่สุดในแผ่นดิน ด้วยการถือมันลงสนามประลองเองทุกสนาม  
> *A blacksmith's daughter set on proving her father's blades are the finest in the land, by carrying them into every arena herself.*

## 6. นักพรตยุทธ์ / Monk — `legend-monk`
**คิรัน สติลวอเตอร์** / **Kiran Stillwater** · role: `control` (ควบคุม)

> ศิษย์คนสุดท้ายของวัดบนหน้าผาที่ถูกเผา ออกตามหาคัมภีร์ที่ถูกขโมย โดยพกไปแค่ชาดีกับความอดทน  
> *The last disciple of a burned cliffside monastery, tracking its stolen scrolls with nothing but good tea and patience.*

## 7. อัศวินศักดิ์สิทธิ์ / Paladin — `legend-paladin`
**โทเบียส เวลลัน** / **Tobias Vellan** · role: `hybrid` (ไฮบริด)

> ชาวนาธรรมดาที่ได้ยินเสียงเรียกในฝันคืนที่หมู่บ้านถูกเผา และยังไม่แน่ใจว่าตัวเองคู่ควรกับคำสาบานที่รับมา  
> *A plain farmer who heard a calling in his dreams the night his village burned, and still isn't sure he deserves the oath he took.*

## 8. พรานป่า / Ranger — `legend-ranger`
**โยเรน วาสก์** / **Joren Vask** · role: `damage` (ดาเมจ)

> คนนำทางกองคาราวานที่รู้ทุกทางลัดบนภูเขา และร้อยเขี้ยวสัตว์ร้ายทุกตัวที่ล่าได้เป็นสร้อยคอ  
> *A caravan guide who knows every mountain shortcut and strings a fang from every beast he has hunted onto his necklace.*

## 9. นักย่องเบา / Rogue — `legend-rogue`
**ไซลัส เรน** / **Silas Wren** · role: `damage` (ดาเมจ)

> สายลับของสภาเมืองที่แฝงตัวเข้ากลุ่มนักผจญภัย และยังตัดสินใจไม่ได้ว่าจะส่งรายงานหรือเลือกเพื่อนใหม่  
> *A city council spy planted among adventurers, who still can't decide whether to file his report or side with his new friends.*

## 10. นักเวทโดยกำเนิด / Sorcerer — `legend-sorcerer`
**เรมี ซอลต์ไฟร์** / **Remy Saltfire** · role: `damage` (ดาเมจ)

> ลูกชาวประมงที่เผาเรือของพ่อตอนโกรธครั้งแรกจนต้องหนีออกจากหมู่บ้าน ตอนนี้กำลังหัดโกรธให้ถูกเวลา  
> *A fisher's son who set his father's boat ablaze the first time he lost his temper and had to flee the village. Now he is learning to get angry at the right moment.*

## 11. ผู้ทำสัญญา / Warlock — `legend-warlock`
**เวสนา มอร์โรว์** / **Vesna Morrow** · role: `hybrid` (ไฮบริด)

> หมอยาที่ทำสัญญากับบางสิ่งจากใต้ทะเลลึกเพื่อรักษาลูกสาว และตอนนี้มันเริ่มกระซิบทวงค่าตอบแทนทุกคืน  
> *An apothecary who struck a pact with something from the deep sea to save her daughter, and now it whispers every night for its payment.*

## 12. พ่อมด / Wizard — `legend-wizard`
**ธีโอดริก เพลล์** / **Theodric Pell** · role: `control` (ควบคุม)

> บรรณารักษ์ที่แอบเรียนเวทจากหนังสือต้องห้ามในห้องใต้ดินของหอสมุด และถูกไล่ออกตอนหนังสือเล่มหนึ่งกัดคนอ่าน  
> *A librarian who secretly learned magic from the forbidden books in the archive vault, and was dismissed when one of them bit a reader.*
