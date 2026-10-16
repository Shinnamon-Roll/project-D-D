# คลาส — research แบ่งตามแพค

ข้อมูล research จากหน้า [D&D Beyond Classes](https://www.dndbeyond.com/classes) ใช้ประกอบการออกแบบ ยังไม่ได้เอาเข้าแอป
(โรสเตอร์ 10 ตัวในแอปเป็นของเราเอง ดู `docs/design/roster.md`)

- ทุกไฟล์สรุปด้วยคำของเราเอง ไม่ copy ข้อความ feature/flavor คำต่อคำ (เนื้อหาเป็นลิขสิทธิ์ WotC)
- เก็บเฉพาะ base class ไม่ลงรายละเอียด subclass (มีแค่ชื่อ subclass ไว้อ้างอิง)
- บทสนทนา: [`../../conversations/2026-10-16-class-research.md`](../../conversations/2026-10-16-class-research.md)

## แพค: Player's Handbook (2024) — [`phb-2024/`](phb-2024/)

| คลาส | ชื่อไทย | Hit Die | Primary ability | บทบาท | ความยาก |
|---|---|---|---|---|---|
| [Barbarian](phb-2024/barbarian.md) | นักรบเถื่อน | d12 | Strength | tank + damage (melee) | ง่าย |
| [Bard](phb-2024/bard.md) | นักกวีเวทดนตรี | d8 | Charisma | support + face + control | กลาง |
| [Cleric](phb-2024/cleric.md) | นักบวช | d8 | Wisdom | support (healer) + damage/tank รอง | กลาง |
| [Druid](phb-2024/druid.md) | ดรูอิด | d8 | Wisdom | support / control / utility / scout | กลาง-ยาก |
| [Fighter](phb-2024/fighter.md) | นักรบ | d10 | Strength หรือ Dexterity | tank / damage | ง่าย |
| [Monk](phb-2024/monk.md) | นักพรตยุทธ์ | d8 | Dexterity + Wisdom | damage / scout / control | กลาง |
| [Paladin](phb-2024/paladin.md) | อัศวินศักดิ์สิทธิ์ | d10 | Strength + Charisma | tank / damage (burst) / support / face | ง่าย-กลาง |
| [Ranger](phb-2024/ranger.md) | พรานป่า | d10 | Dexterity + Wisdom | damage (ranged) / scout / utility | กลาง |
| [Rogue](phb-2024/rogue.md) | นักย่องเบา | d8 | Dexterity | damage (burst) / scout / utility / face | ง่าย-กลาง |
| [Sorcerer](phb-2024/sorcerer.md) | นักเวทโดยกำเนิด | d6 | Charisma | damage / control / face | กลาง |
| [Warlock](phb-2024/warlock.md) | ผู้ทำสัญญา | d8 | Charisma | damage / face / utility | กลาง |
| [Wizard](phb-2024/wizard.md) | พ่อมด/นักเวทศึกษา | d6 | Intelligence | control / damage (area) / utility | ยาก |

## จุดที่ยังไม่ได้ตรวจกับต้นทาง

- ชื่อ subclass ส่วนใหญ่: หน้า D&D Beyond โชว์ subclass เต็มแค่ตัวเดียวต่อคลาส ชื่อที่เหลือมาจากความรู้ของ agent หรือเว็บอื่น (dexerto, rpgbot)
- HP เลเวล 1 / ต่อเลเวล: หน้าไม่ได้เขียนไว้ คำนวณจากกฎมาตรฐาน (ลูกเต๋าเต็ม + CON ที่เลเวล 1, ค่าเฉลี่ย + CON ต่อเลเวล)
- Wizard spellbook: จำนวนขั้นต่ำที่คำนวณเอง (6 + 2 ต่อเลเวล)

## แพคอื่น

ยังไม่ได้ทำ (ผู้ใช้สั่ง Player's Handbook เท่านั้น) เช่น Artificer (Eberron/Tasha's), Blood Hunter (Critical Role)
