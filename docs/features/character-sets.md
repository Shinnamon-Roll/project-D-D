# Feature: เซ็ตตัวละคร (character sets)

บทสนทนา: [`../conversations/2026-10-16-character-sets.md`](../conversations/2026-10-16-character-sets.md) · ต่อยอดจาก [`character-roster.md`](character-roster.md)

## เป้าหมาย
DM เลือกได้ว่าเกมนี้ใช้ฮีโร่จากเซ็ตไหนบ้าง (เลือกได้หลายเซ็ต) ระหว่างรอผู้เล่นใน lobby

## เซ็ต
| id | ชื่อ (TH / EN) | ฮีโร่ | ข้อมูล |
|---|---|---|---|
| `core` | ฮีโร่ชุดแรก / First Company | 10 | [`../design/roster.md`](../design/roster.md) |
| `legend` | วีรชนในตำนาน / Heroes of Legend | 12 (1 ต่อคลาส: barbarian bard cleric druid fighter monk paladin ranger rogue sorcerer warlock wizard) | [`../design/roster-legend.md`](../design/roster-legend.md) |

- id ฮีโร่ไม่ซ้ำข้ามเซ็ต: เซ็ต legend ใช้ `legend-<class>`
- เนื้อหาเซ็ต legend แต่งเองทั้งหมด (อ้างแค่แนวคิดคลาสจาก research) ไม่ใช้ชื่อ feature ทางการ ไม่มีคำว่า "Player's Handbook"/"D&D" ใน UI
- รูป: icon ตามคลาส (`legend-barbarian` ใช้ `barbarian.svg`) เพิ่ม icon ใหม่ fighter/monk/sorcerer/warlock จาก game-icons.net (CC BY 3.0)

## กติกา
1. ค่าเริ่มต้นห้องใหม่ = เปิดเฉพาะ `core`
2. DM หรือ co-DM ติ๊กเลือกเซ็ตได้เฉพาะสถานะ `lobby` ต้องเปิดอย่างน้อย 1 เซ็ต
3. หลังกดเริ่มเกม เซ็ตล็อก ผู้เล่นเลือกได้เฉพาะฮีโร่ในเซ็ตที่เปิด
4. ฮีโร่ 1 ตัวเลือกได้คนเดียว แต่คลาสเดียวกันจากคนละเซ็ตเลือกพร้อมกันได้
5. หน้า Ruleset: แท็บละเซ็ต แก้ฮีโร่ในแท็บนั้น ปุ่ม reset = รีเซ็ตเฉพาะเซ็ตนั้น (ยังเก็บใน `data/roster.json` ไฟล์เดียว key ตาม id)

## Protocol (contract)
- `RosterEntry.set: "core" | "legend"` (server บังคับให้ตรงกับเซ็ตเดิมของ id นั้น)
- `RoomState.sets: {id: string, name: {th, en}}[]` — รายชื่อเซ็ตทั้งหมด (คงที่)
- `RoomState.activeSets: string[]` — เซ็ตที่เปิด เรียงตามลำดับใน `sets` ค่าเริ่ม `["core"]`
- `RoomState.roster`: DM/co-DM ได้ครบทุกเซ็ต (ใช้แก้ใน Ruleset) · ผู้เล่นได้เฉพาะเซ็ตที่เปิด
- event ใหม่ `sets.set {sets: string[]}` — DM/co-DM, สถานะ lobby เท่านั้น; ว่าง/มี id ไม่รู้จัก = error; ตัดตัวซ้ำ; เขียน DM Log 1 บรรทัด
- `roster.reset {set: "legend"}` — รีเซ็ตเฉพาะเซ็ต; `{classId}` และ `{}` ทำงานเหมือนเดิม
- `character.pick {classId}` — ปฏิเสธถ้าฮีโร่อยู่ในเซ็ตที่ไม่ได้เปิด ("that set is not in play")

## UI
- Lobby (DM): กล่องเลือกเซ็ตแบบ checkbox หลายอัน (ชื่อ + จำนวนฮีโร่) อันสุดท้ายที่ติ๊กอยู่กดเอาออกไม่ได้
- Lobby (ผู้เล่น): เห็นชื่อเซ็ตที่เปิด (อ่านอย่างเดียว)
- เลือกตัวละคร: carousel โชว์เฉพาะฮีโร่ในเซ็ตที่เปิด; เปิดมากกว่า 1 เซ็ต มีป้ายชื่อเซ็ตบนการ์ด
- Ruleset: แท็บเซ็ตด้านบนรายการ

## นอก scope
- DM สร้างเซ็ตใหม่เอง / ลบเซ็ต
- เปลี่ยนเซ็ตหลังเริ่มเกม
