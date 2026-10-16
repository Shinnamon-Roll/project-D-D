# รีวิวระบบที่ยังขาด + flow แปลกๆ (v0.18.0, 2026-10-17)

บทสนทนา: [`../conversations/2026-10-17-gap-review.md`](../conversations/2026-10-17-gap-review.md)

ที่มา: agent อ่านโค้ดจริงไล่ flow ทั้งฝั่ง DM และผู้เล่น (read-only) + agent `researcher` เทียบกับกติกา SRD 5.2.1 และ VTT อื่น (Owlbear, Foundry, D&D Beyond Maps, Shard, Alchemy)
CONFIRMED = เห็นในโค้ด · SUSPECTED = น่าจะเป็น ยังไม่ได้ลองเล่นจริง

## 1. Flow ที่แปลก / ติดตาย

### สถานะห้อง (state machine)
`menu → lobby → creating → playing` เดินทางเดียว ไม่มี event ย้อนกลับ, จบเกม, หรือเริ่มใหม่ (`internal/room/room.go:438-628`) — CONFIRMED

| # | ปัญหา | ผลตอนเล่นจริง | ที่มา |
|---|---|---|---|
| F1 | **ห้องค้างที่ `creating`**: เข้า `playing` อัตโนมัติเมื่อผู้เล่น**ทุกคน**เลือกฮีโร่ ถ้ามีคนหลุดไม่กลับมา DM เตะหรือบังคับเริ่มไม่ได้ | ทั้งโต๊ะติด ต้องปิด server | `character.go:169-175` |
| F2 | **จัดการที่นั่งหลังเริ่มเกมไม่ได้**: `seat.add/remove/sit/leave/openCoDm` ใช้ได้แค่ lobby → เพื่อนมาสายเข้าไม่ได้, เตะคนไม่ได้, เปิด co-DM เพิ่มไม่ได้ | ต้องเริ่มใหม่ทั้งหมด | `room.go:444,477,493,511,722` |
| F3 | **เลือกฮีโร่ใหม่ไม่ได้** (`you already have a character`) ไม่มี `character.remove` | กดผิดแล้วแก้ไม่ได้ | `character.go:138` |
| F4 | **DM เสีย `clientId`** (เปลี่ยนเครื่อง/private window/ล้าง storage) = ไม่มีใครคุมห้องได้อีก เพราะ claim DM ทำได้แค่ใน `menu` | เกมเดินต่อไม่ได้ | `room.go:450`, `ws.ts:11-20` |
| F5 | **ไม่มีปุ่มจบ session / กลับ lobby** (web-flow.md บอกมี `playing → lobby` แต่โค้ดไม่มี) | ต้อง kill server | `web-flow.md:48` |
| F6 | **DM เตรียมแผนที่ไม่ได้ระหว่างผู้เล่นเลือกฮีโร่**: server อนุญาต `scene.add` ใน `creating` แต่หน้าจอ DM ตอน creating มีแค่ Notes + Roster | ผู้เล่นเข้าเกมแล้วเจอ "รอฉาก" ทุกครั้ง | `room.go:816`, `CharacterPick.tsx:31-64` |
| F7 | ผู้เล่นที่เลือกแล้วเห็นแค่ "รอคนอื่น" ดูชีทตัวเองไม่ได้จนทุกคนเลือกเสร็จ | รอเปล่า | `CharacterPick.tsx:68` |
| F8 | **คุยกับ NPC / ให้ของ / ยอมแพ้ ต้องเปิด combat ก่อน** — ฉากคุยนอกการต่อสู้ใช้ dialog NPC ไม่ได้ | ฉากสังคมต้องเปิด combat หลอก | `combat.go:334`, `Table.tsx:63` |
| F9 | ผู้เล่นประกาศโจมตี**ผู้เล่นด้วยกัน**ได้ (กันแค่ตัวเอง) — SUSPECTED ว่าไม่ได้ตั้งใจ | กดพลาดแล้ว DM ต้องตัดสิน | `Table.tsx:57`, `combat.go:183-197` |
| F10 | กดครั้งเดียวไม่มี confirm / undo: **พักยาว** (เติมทั้งปาร์ตี้), **จบ combat** (ลบการโจมตีที่ค้าง) | กดพลาดแก้กลับไม่ได้ | `PartyBar.tsx:75,81` |
| F11 | **ไม่รู้ว่าใครหลุด**: Seat ไม่มี online/offline, แจ้งแค่ใน DM log | DM ไม่รู้ว่ารอใคร | `room.go:84-92`, `hub.go:236-249` |
| F12 | ใช้ hit die แล้ว HP ไม่ขึ้น ต้องรอ DM บวกเอง (ผู้เล่นแก้ HP ตัวเองไม่ได้) — ตั้งใจไว้ แต่งงได้ | DM งานเยอะตอนพักสั้น | `character.go:187`, `i18n.ts:332` |
| F13 | ข้อความ error เลือกซ้ำบอก "that class is already taken" แต่จริงๆ เช็คต่อฮีโร่ (คลาสเดียวกันต่างเซ็ตเลือกได้) | ข้อความผิด | `character.go:155` |
| F14 | คนที่ไม่ได้นั่ง เปิดลิงก์ตอน `playing` จะเห็นโต๊ะ (เป็นผู้ชมโดยปริยาย) แต่ docs บอกไม่มีผู้ชม | ข้อมูลผู้เล่นรั่วให้คนนอกดู | `App.tsx:33`, `web-flow.md:21` |
| F15 | หมอก (fog) ซ่อนแค่ฝั่ง client: token ใต้หมอกยังส่งไปเบราว์เซอร์ผู้เล่น (ยอมรับไว้แล้ว) | เปิด devtools เห็น | `room.go:237-251` |

### docs ไม่ตรงโค้ด (`docs/design/web-flow.md`)
- state machine เริ่มที่ `lobby` (จริงเริ่ม `menu`), มี `playing → lobby` (ไม่มีจริง)
- ปุ่ม "เป็น DM" ใน lobby (จริงเลือกที่หน้าแคมเปญ), ปุ่ม Map ใน lobby (ไม่มี), กลับเข้าด้วยเลขห้อง (จริงใช้ `clientId` ใน localStorage, Q13 ยังเปิด)
- data model เก่า (`mapImageUrl`, fog `{x,y}`, Seat `roomId/order`), ยังเขียน "10 ฮีโร่"
- ระบบที่เอาออก (มานา, `skill.use`, `setStats`) ไม่เหลือในโค้ด/i18n แล้ว · `character-roster.md` ยังพูดถึงมานาแต่มีป้าย Superseded

## 2. การเก็บข้อมูล (persistence)
| อยู่รอด server restart | หายตอน restart / กด reload ใน TUI |
|---|---|
| `data/notes.json`, `npcs.json`, `roster.json`, ไฟล์ใน `uploads/` | สถานะห้อง, แคมเปญ, ที่นั่ง + DM, **ตัวละครทั้งหมด** (เลเวล HP slot ของ คาถา), ฉาก + grid + fog, token + HP NPC, feed, combat, DM log |

- หลัง restart ทุกคนกลับ `menu` และ**ใครกดก่อนได้เป็น DM** (ไม่มีรหัสห้อง)
- refresh เบราว์เซอร์เดิม = กลับที่นั่งเดิมได้ (`clientId`) · เปลี่ยนเครื่อง = ที่นั่งเดิมหายถาวร (นอก lobby)
- รูปแผนที่ใน `uploads/` ค้างโดยไม่มีฉากชี้ (SUSPECTED orphan)

## 3. ระบบเกมที่ยังขาด (เทียบ SRD 5.2.1 + VTT อื่น)
ขาดทั้งหมด เว้นแต่บอกไว้ · มีแล้ว: HP, เลเวล, slot, counter, hit dice, พักสั้น/ยาว, HP NPC, token ซ่อน, combat hit/miss, log

| ระบบ | เจอบ่อยแค่ไหน | VTT อื่นทำขั้นต่ำยังไง | ขั้นต่ำที่เข้ากับแอปนี้ |
|---|---|---|---|
| **Save เกม** | ทุก session ที่ 2 เป็นต้นไป (แคมเปญยาวหลายคืน) | Owlbear 2.0 เปลี่ยนเป็นห้องถาวร [12] | เขียนห้องทั้งหมดลง `data/room.json` (ใช้ atomic write แบบ notes/npcs) โหลดตอนเปิด server ⚠️ จะเขียน `clientId` ลงดิสก์ |
| **Death saves** (SRD p17) | เลเวล 1 ล้มได้ทุกไฟต์ | ช่องในชีท | HP = 0 → pip สำเร็จ 3 / ล้มเหลว 3 เจ้าของหรือ DM ติ๊ก, HP > 0 ล้างเอง, feed "stable"/"dead" |
| **Conditions** (15 อย่าง + exhaustion 1–6, p181) | Prone/Poisoned/Frightened แทบทุกไฟต์ | D&D Beyond Maps คลิกใส่ [3], Shard มี duration [9] | ไอคอน condition บน token (game-icons CC BY), exhaustion 0–6, ไม่คำนวณผลอัตโนมัติ, NPC ซ่อนกรองเหมือน HP |
| **Initiative / ลำดับเทิร์น** (p13) | ทุกไฟต์ | Owlbear: list + ปุ่ม next [4], Foundry พิมพ์ค่าเอง [7] | เฉพาะตอน combat: DM พิมพ์ค่า (ทอยเต๋าจริง) server เรียง, ตัวนับรอบ, ปุ่มถัดไป, NPC ซ่อนไม่ส่งให้ผู้เล่น ⚠️ **ขัดกับ decision เดิม** (decisions.md:29 "ไม่เพิ่ม initiative จนกว่าจะได้ลองเล่น") |
| **Temp HP** (p18) | บ่อยตั้งแต่เลเวล 2 | ช่องในชีท | `tempHp` ดาเมจหักก่อน (`attack.resolve`, `npc.act`, `character.adjust`), พักยาวล้าง |
| **Concentration** (p178, DC = max(10, ดาเมจ÷2) สูงสุด 30) | โต๊ะจริงลืมบ่อย [11] | ใส่วงสีเอง [3] | toggle "กำลัง concentrate X" + วงบน token, โดนดาเมจ → feed/DM เตือน "CON save DC N" |
| **วัดระยะ (ft)** (ช่อง 5 ft, แนวทแยงนับ 1) | ทุกไฟต์ (ยิงไกล/เดิน) | Owlbear ruler แบบ chessboard [6] | client อย่างเดียว: เลือก token แล้วชี้ช่อง แสดง "N ft" |
| **NPC stat block** | DM ต้องรู้ AC + โบนัสโจมตี | Owlbear extension [5] | เพิ่ม AC + โน้ต DM ต่อ NPC ใน pool (ห้ามใส่ stat block จาก Starter Set) |
| **Bloodied** (HP ≤ ครึ่ง, p177) | ผู้เล่นถาม "มันเจ็บยัง" ทุกไฟต์ | วงแดงทำมือ [3] | ผู้เล่นเห็นป้าย Bloodied บน NPC ที่มองเห็น ไม่เห็นตัวเลข HP |
| **Handouts** | จดหมาย/เบาะแส ไม่กี่ครั้งต่อ session | Alchemy panel [10] | DM แชร์รูปที่อัปโหลดหรือโน้ต ให้ทุกคนหรือคนเดียว |
| **เงิน** (p89) | ซื้อของ/ของตก | ช่องในชีท | ใช้ item stack ได้อยู่แล้ว (สูงสุด 999) |

**ข้ามได้สำหรับเล่นกับเพื่อน**: คลังคาถา (~339 คาถา ต้องเขียนใหม่ TH/EN, "คาถาของฉัน" ใช้แทนได้), XP (Starter Set ใช้ milestone + DM มีเลเวล ±), น้ำหนัก/encumbrance, AoE template, คลังมอนสเตอร์, ผลของ condition อัตโนมัติ, แสงเงา, แชท/เสียง (ใช้ Discord)

## 4. จัดอันดับรวม (คุ้มค่า ÷ แรง)
| # | อะไร | ทำไม | แรง |
|---|---|---|---|
| 1 | **Save/load ห้อง** | restart แล้วหายหมด เล่นหลายคืนไม่ได้ | กลาง-ต่ำ |
| 2 | **แก้ flow ติดตาย**: DM บังคับเริ่ม/เตะคนตอน creating, เพิ่ม/ลบที่นั่งหลังเริ่ม, เลือกฮีโร่ใหม่, ปุ่มจบ session | ติดแล้วต้องปิด server | กลาง |
| 3 | **Death saves** | ล้มได้ตั้งแต่ไฟต์แรก | ต่ำ |
| 4 | **ไฟ online/offline ที่ที่นั่ง** | DM รู้ว่ารอใคร | ต่ำ |
| 5 | **Temp HP** | ใช้บ่อย | ต่ำ |
| 6 | **Conditions + exhaustion** | ทุกไฟต์ | กลาง |
| 7 | **Initiative** (ต้องกลับ decision เดิม) | ทุกไฟต์ | กลาง |
| 8 | วัดระยะ ft | ทุกไฟต์ | ต่ำ (client) |
| 9 | confirm ปุ่มพักยาว/จบ combat + NPC talk นอก combat | กันพลาด / ฉากคุย | ต่ำ |
| 10 | Concentration, NPC AC, Bloodied, DM เตรียมแผนที่ตอน creating, Handouts | | ต่ำ-กลาง |
| — | ซิงก์ `web-flow.md` กับโค้ด | docs เก่า | docs |

## 5. แหล่งอ้างอิง (จาก researcher)
[1] SRD 5.2.1 — https://media.dndbeyond.com/compendium-images/srd/5.2/SRD_CC_v5.2.1.pdf ·
[2] https://www.dndbeyond.com/posts/1841-roll-for-initiative-combat-tracking-comes-to-the ·
[3] https://www.dndbeyond.com/posts/2212-condition-tracking-comes-to-the-maps-vtt ·
[4] https://extensions.owlbear.rodeo/initiative-tracker ·
[5] https://obr-extensions-guide.onrender.com/ ·
[6] https://docs.owlbear.rodeo/docs/measure/ ·
[7] https://foundryvtt.com/article/combat/ ·
[8] https://foundryvtt.com/article/measurement/ ·
[9] https://shardtabletop.com/blog/conditions-and-hero-abilities ·
[10] https://help.alchemyrpg.com/en/articles/9714185-gamemaster-orientation ·
[11] https://www.enworld.org/goto/post?id=7100706 ·
[12] https://docs.owlbear.rodeo/docs/migration/migrate-from-1-to-2/

ความมั่นใจ: กติกา SRD สูง (อ่านจาก PDF) · วิธีของ VTT อื่น กลาง (Roll20 โหลดไม่ได้ 403) · ความถี่ที่เจอ ต่ำ-กลาง (อนุมาน ไม่มีสถิติ)
