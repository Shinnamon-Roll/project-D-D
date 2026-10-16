# docs — สารบัญเอกสารทั้งหมด

ทุกการคุย แผน flow ไอเดีย และการตัดสินใจ เก็บเป็นไฟล์ที่นี่ ไม่ปล่อยไว้แค่ในแชท

| โฟลเดอร์/ไฟล์ | เก็บอะไร | อัปเดตเมื่อ |
|---|---|---|
| [`conversations/`](conversations/) | บันทึกการคุยแต่ละรอบ: คำขอผู้ใช้ (คำต่อคำ) + สรุป + สิ่งที่ตัดสินใจ | ทุกครั้งที่ผู้ใช้ให้ requirement/feedback ใหม่ |
| [`design/`](design/) | flow การทำงานของเว็บ, หน้าจอ, data model, tech stack | เมื่อ requirement เปลี่ยน |
| [`campaigns/`](campaigns/) | เนื้อเรื่อง/ฉาก/NPC ของแต่ละแคมเปญ (เช่น Salem) | เมื่อออกแบบเนื้อเรื่อง |
| [`research/`](research/) | ความรู้พื้นฐานจาก agent `researcher` (มี citation) | เมื่อ research เรื่องใหม่ |
| [`features/`](features/) | spec ทีละฟีเจอร์ (จาก `/new-feature`) | ก่อนเขียนโค้ดแต่ละฟีเจอร์ |
| [`decisions.md`](decisions.md) | การตัดสินใจที่ยืนยันแล้ว (วันที่ + เหตุผล) | เมื่อผู้ใช้ยืนยันอะไร |
| [`open-questions.md`](open-questions.md) | คำถามที่รอผู้ใช้ตอบ | เพิ่มเมื่อเจอจุดไม่ชัด, ย้ายไป decisions เมื่อได้คำตอบ |
| [`ideas.md`](ideas.md) | ไอเดียที่ยังไม่อยู่ใน scope (backlog) | เมื่อมีไอเดียใหม่ |

ชื่อไฟล์บทสนทนา: `conversations/YYYY-MM-DD-<หัวข้อ>.md`

ไฟล์หลักของระบบเกม v2 (v0.18.0):
- spec [`features/rules-v2.md`](features/rules-v2.md) — กติกา 6 ability, เลเวล 1–20, spell slot, counter, rest, events
- ข้อมูลคลาส 12 คลาส [`design/classes/`](design/classes/) (`<id>.md` + ตารางกลาง `_shared.md`, SRD 5.2.1 CC BY 4.0)
- ข้อมูลเกมฮีโร่ [`design/heroes-core.md`](design/heroes-core.md) (core 10) · [`design/heroes-legend.md`](design/heroes-legend.md) (legend 12) — ชื่อ/เรื่องอยู่ใน `design/roster*.md`
- research คลาส [`research/classes/`](research/classes/)

ล่าสุด: รีวิวระบบที่ขาด + flow [`research/gap-review-2026-10-17.md`](research/gap-review-2026-10-17.md) · spec [`features/rules-v2.md`](features/rules-v2.md) · spec [`features/character-sets.md`](features/character-sets.md) · [`conversations/2026-10-16-tui-dashboard.md`](conversations/2026-10-16-tui-dashboard.md) · [`conversations/2026-10-16-map-fit.md`](conversations/2026-10-16-map-fit.md) · spec [`features/seat-qr.md`](features/seat-qr.md) · [`features/map-fit.md`](features/map-fit.md) · ธีม [`design/visual-theme.md`](design/visual-theme.md)
