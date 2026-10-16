# Feature: popup ยืนยันก่อนลบ / ทำสิ่งสำคัญ — v0.18.1

บทสนทนา: [`../conversations/2026-10-17-confirm-dialog.md`](../conversations/2026-10-17-confirm-dialog.md) · ที่มา: รีวิว F10 ใน [`../research/gap-review-2026-10-17.md`](../research/gap-review-2026-10-17.md)

## สรุป
- popup ยืนยันแบบธีมของแอปเอง (ไม่ใช้ `window.confirm()` ของเบราว์เซอร์) ตัวเดียวใช้ทั้งแอป, TH/EN
- frontend อย่างเดียว ไม่เปลี่ยน protocol / server

## Action ที่ต้องถาม
| action | event | ที่ | หมายเหตุ |
|---|---|---|---|
| ลบที่นั่ง | `seat.remove` | `screens/Lobby.tsx` | |
| ลบ token | `token.remove` | `screens/Table.tsx` | |
| ลบคาถา/ความสามารถของฉัน | `ability.remove` | `table/CharSheet.tsx` | |
| ลบของทั้งกอง | `item.remove` (qty = MAX_STACK) | `table/PartyBar.tsx` | ปุ่ม −1 ไม่ถาม |
| รีเซ็ตฮีโร่ตัวเดียว | `roster.reset {classId}` | `table/RosterEditor.tsx` | |
| รีเซ็ตทั้งเซ็ต | `roster.reset {set}` | `table/RosterEditor.tsx` | เดิมใช้ `confirm()` → เปลี่ยนเป็น popup ใหม่ |
| ลบโน้ต | `note.remove` | `table/NotesPanel.tsx` | เดิมใช้ `confirm()` → เปลี่ยน |
| ลบ NPC จาก pool | `npc.remove` | `table/NpcPanel.tsx` | เดิมใช้ `confirm()` → เปลี่ยน |
| พักสั้น / พักยาว | `party.rest` | `table/PartyBar.tsx` | |
| จบ combat | `combat.set {on:false}` | `table/PartyBar.tsx` | เปิด combat ไม่ถาม |
| เริ่มเกม | `game.start` | `screens/Lobby.tsx` | |
| เปิดเผย token ที่ซ่อน | `token.setHidden {hidden:false}` | `screens/Table.tsx` | ซ่อนไม่ถาม |

ไม่ถาม: ลบหมอก (ปุ่มในแผงและคลิกบนแผนที่), ปุ่มอื่นที่กดกลับได้ (เลเวล ±, HP ±, slot, เปลี่ยนฉาก, เลือกเซ็ต)

## พฤติกรรม
- หัวข้อ + ข้อความบอกว่าจะเกิดอะไร (ใส่ชื่อสิ่งที่ลบ เช่น ชื่อ token/โน้ต/ของ)
- ปุ่ม: ยกเลิก (`confirm-cancel`) + ยืนยัน (`confirm-ok`); ปุ่มยืนยันสีแดง (`danger`) สำหรับการลบ/รีเซ็ต, สีหลัก (`primary`) สำหรับอย่างอื่น
- Esc / คลิกพื้นหลัง = ยกเลิก; focus เริ่มที่ปุ่มยกเลิก; `role="alertdialog"` + `aria-modal`
- เปิดได้ทีละอัน; ถ้าสถานะห้องเปลี่ยนจนปุ่มเดิมไม่มีความหมาย server ยังเป็นคนตัดสิน (error toast เดิม)

## Test
- e2e: ทุก flow เดิมที่กดปุ่มเหล่านี้ต้องกด `confirm-ok`; เพิ่ม 1 เคส กด `confirm-cancel` แล้วของยังอยู่
