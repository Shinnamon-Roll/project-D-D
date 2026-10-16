# Feature: NPC ทำ action กับผู้เล่น (โหมดต่อสู้)

บทสนทนา: [`../conversations/2026-10-16-npc-actions.md`](../conversations/2026-10-16-npc-actions.md) · ต่อยอดจาก [`combat-log.md`](combat-log.md)

## Flow (DM)
1. โหมดต่อสู้เปิดอยู่ (`room.combat`)
2. DM จิ้ม token NPC/มอนสเตอร์บนแมพ (เลือกเหมือนเดิม)
3. DM จิ้ม token ตัวละครผู้เล่น → ป๊อปอัป "<NPC> → <ผู้เล่น>" มี 4 ปุ่ม: **คุย / โจมตี / มอบของ / ยอมแพ้**
   - ปิดโหมดต่อสู้: จิ้มแบบนี้ = เลือก token ผู้เล่นแทน (เหมือนเดิม)
   - Esc / ปุ่มยกเลิก = ปิดป๊อปอัป ไม่ส่งอะไร
4. ช่องในป๊อปอัป
   - คุย: ข้อความไม่บังคับ (≤ 200 ตัวอักษร)
   - โจมตี: ปุ่ม โดน / ไม่โดน + ดาเมจ (เลขเต็ม 0–999, ใช้เมื่อโดน)
   - มอบของ: ชื่อของ (1–40) + จำนวน (1–99) → เข้ากระเป๋าผู้เล่นจริง (กติกาเดียวกับ DM ให้ของในการ์ดโปรไฟล์)
   - ยอมแพ้: ไม่มีช่อง
5. ส่งแล้ว → ป้ายใหญ่กลางจอทุกคน + บรรทัดในบันทึกการเล่น (ทุกคน) + DM log

## ข้อความ (TH / EN) — ชื่อ NPC = ชื่อ token เต็ม รวมเลข `#n` ทั้ง feed, ป้าย และ DM log (เช่น "Goblin #2 โจมตี Ayla — โดน! −5") ตาม [`npc-names.md`](npc-names.md)
| action | TH | EN | ป้าย |
|---|---|---|---|
| talk | Goblin คุยกับ Ayla: “ส่งทองมา!” (ไม่มีข้อความ = Goblin คุยกับ Ayla) | Goblin talks with Ayla: “…” | talk |
| attack โดน | Goblin โจมตี Ayla — โดน! −5 | Goblin attacks Ayla — hit! −5 | hit |
| attack ไม่โดน | Goblin โจมตี Ayla — ไม่โดน | Goblin attacks Ayla — miss | miss |
| gift | Goblin มอบ Rusty Key ×1 ให้ Ayla | Goblin gives Ayla Rusty Key ×1 | gift |
| surrender | Goblin ยอมแพ้ต่อ Ayla | Goblin surrenders to Ayla | surrender |

โดนแล้ว HP ผู้เล่นถึง 0 = บรรทัด 💀 ล้มลง ตาม helper เดิม (`downed`)

## Protocol (contract)
- event ใหม่ `npc.act {tokenId, targetId, kind, text?, hit?, damage?, item?: {name, qty}}`
  - `kind`: `"talk" | "attack" | "gift" | "surrender"`
  - DM/co-DM เท่านั้น, สถานะ `playing`, โหมดต่อสู้เปิด (ไม่งั้น `combat mode is off`)
  - `tokenId` = token NPC (ไม่ใช่ตัวละครผู้เล่น) ที่มองเห็นอยู่ (ซ่อน = error "reveal the NPC first")
  - `targetId` = token ตัวละครผู้เล่น ในฉากเดียวกัน
  - talk: `text` ≤ 200 (CleanText, ว่างได้) · attack: `hit` bool, `damage` 0–999 (ใช้เมื่อ hit) · gift: `item.name` 1–40, `item.qty` 1–99, กระเป๋าเต็ม = error เดิม · surrender: ไม่มี field
- `Banner.kind` เพิ่ม `"talk" | "gift" | "surrender"` (เดิม `"hit" | "miss"`) ส่งให้ทุก connection ครั้งเดียว (`{type:"banner"}`)
- ผลกระทบ: attack โดน = ลด HP ตัวละครเป้า (clamp 0) ผ่าน path เดียวกับ `attack.resolve` · gift = เพิ่มของในกระเป๋า (path เดียวกับ `item.give`)

## UI
- ป๊อปอัปแบบ dialog (focus อยู่ในป๊อปอัป, Esc ปิด, คืน focus) ตามสไตล์ DM alert ของการโจมตีเดิม
- ป้าย: hit/miss แบบเดิม · talk/gift/surrender ใช้ป้ายขนาดเท่ากัน สีต่างกัน (talk = หมึก/กระดาษ, gift = ทอง, surrender = ขาว/เทา) + reduced-motion
- ปุ่มใน token list (DM, โหมดต่อสู้): แถว NPC มีปุ่ม "Action" ให้เลือกเป้าจากรายชื่อผู้เล่นได้ด้วย (ทางเลือกสำหรับ keyboard/จอเล็ก)

## นอก scope
- NPC ทำ action กับ NPC ด้วยกัน
- ผู้เล่นตอบกลับ/ปฏิเสธของ
- action นอกโหมดต่อสู้
