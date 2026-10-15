# VTT ที่มีอยู่ + ผู้เล่นต้องการอะไรจากระบบ

ผลวิจัยวันที่ 2026-10-16 (agent `researcher`) — ราคา/เวอร์ชัน/ดาว ณ วันนั้น
แหล่งอ้างอิงเลข `[n]` อยู่ท้ายไฟล์ · ความเห็นจาก Reddit = ความเห็นรายบุคคล ไม่ใช่ผลสำรวจ

## สรุปสั้น
- กลุ่มเพื่อนเล็ก ๆ เลือก VTT จาก **ความง่ายในการเริ่ม** มากกว่าจำนวนฟีเจอร์: "เปิดลิงก์ ลากแมพ ขยับ token ทอยเต๋า" ชนะ automation [23][25][26][30]
- โมเดลที่ใกล้โปรเจกต์นี้สุด = **Owlbear Rodeo**: ฐานเบา (แมพ, token, fog, วาด, ไม้บรรทัด) ผู้เล่นไม่ต้องสมัคร ที่เหลือเป็น extension [2][3][4]
- เสียงคุยใช้ Discord, character sheet ใช้ D&D Beyond อยู่แล้ว → ไม่ต้องทำเองใน MVP [25][29]

## 1. เทียบฟีเจอร์ VTT
Y = มีในตัว, Ext = extension, N = ไม่มี, ? = ยังไม่ยืนยัน

| ฟีเจอร์ | Roll20 | Foundry VTT | Fantasy Grounds | Owlbear Rodeo 2.x | Tableplop |
|---|---|---|---|---|---|
| แมพ + grid | Y [16] | Y [12] | Y [19] | Y สี่เหลี่ยม/hex, canvas ไม่จำกัด [5] | Y [21][22] |
| Token | Y | Y | Y | Y, 6 ชนิดภาพ = layer [5] | Y [22] |
| Fog of war | Y [16] | Y [12] | Y (LOS) [19] | Y, วาด shape + Fill/Cut [5] | Y [21] |
| Dynamic lighting | Plus ขึ้นไป $5.99/เดือน [15] | Y [12] | Y [19] | Y ทุก tier + ext `dynamic-fog` [1][10] | ? |
| ทอยเต๋า | Y | Y | Y 3D [19] | Ext (3D dice ทางการ) [4] | Y, log, ทอยลับ, adv/dis [22] |
| Character sheet | Y 250+ ระบบ [16] | Y ตามระบบ [12] | Y | N [25] | Y template เอง [21][22] |
| Rules automation | macro; API เฉพาะ Pro [15] | หนัก ผ่าน module [12] | หนัก [19] | N | บางส่วน [22] |
| เสียง/วิดีโอ | มีในตัว [18] ไม่มีบน tablet [17] | มีในตัว [12] | ไม่มี [19] | N (ย้ายไป Kenku FM บน Discord) [4] | ? |
| Initiative tracker | Turn Tracker [18] (เรื่อง tier ขัดกัน) | Combat tracker | Y | Ext ทางการ [4][10] | Y [22] |
| วัดระยะ/AoE | Y | Y | Y | ไม้บรรทัด sync + shape [4][5] | Y [22] |
| มือถือ | มีแอป ไม่มี A/V [17] | ไม่รองรับ iPad ใช้ไม่ได้ [12] | desktop เท่านั้น [19] | responsive มือถือได้ [3] | ? |
| ราคา | ฟรี; $5.99 / $10.99 / $14.99 ต่อเดือน [15] | $50 ครั้งเดียว (GM) [11] | ฟรีตั้งแต่ 2025-11-08, content เสียเงิน [19][20] | ฟรี (200MB, 2 ห้อง); $3.99 / $7.99 [1] | ฟรี (Patreon) [21] |
| Self-host | N | ได้ (port-forward) [13] | GM รันแอป [19] | N (1.0 source เปิด non-commercial) [9] | N |
| บัญชี | ต้องสมัคร [16] | 1 license/host; ผู้เล่นเข้า browser [12] | ต้อง [19] | GM ต้อง; ผู้เล่นไม่ต้อง [2][6] | ต้อง (ยังไม่ยืนยัน) |

Roll20 ฟรี: พื้นที่ 100MB, อัปได้ไฟล์ละ 20MB [15]

## 2. สองแนวคิด: เบา vs automation หนัก
**สายเบา (Owlbear)**
- ชอบ: "ไม่ต้องให้ผู้เล่นสมัคร แค่ส่งลิงก์", เล่นบนมือถือ + Discord ได้ [30]; โหลดเร็ว ไม่ล้น [25]; "แค่อยากวางแมพแล้วลากของไปมา" [23]
- บ่น: ไม่มี character sheet [25]; initiative เทอะทะ → DM ไปใช้ D&D Beyond แทน [25][29]; 1.0 เก็บข้อมูลใน browser แล้วหาย [30][4]

**สาย automation (Foundry, Roll20)**
- ชอบ: automation ช้าตอนแรก แต่เร็วขึ้นภายหลัง [26]; import walls/lighting ไว [29]
- บ่น: GM ใช้เวลา "เขียนโค้ดเกม" นาน [29]; self-host ต้อง port-forward [13][29]; อัปเดตทำ module พัง ทั้งเกมล่ม [24][28]; UI Roll20 "รก" [27]; A/V ในตัวไม่เสถียร → ย้ายไป Discord [29]

## 3. ปัญหาที่เจอบ่อย (เรียงตามความถี่ — เป็นการสรุปเอง)
1. **ติดตั้ง/พังกลางเซสชัน** — port-forward, module พัง, A/V หลุด [13][24][28][29]
2. **เรียนรู้ยาก, UI รก** [26][27][29]
3. **เตรียมเกมหนัก** — วาด walls, lighting, ตั้ง automation [29]
4. **ค่าใช้จ่าย** — ฟีเจอร์ล็อกหลัง subscription, content เสียเงิน [15][19][11]
5. **ข้อมูลหาย** เพราะเก็บแค่ใน browser (Safari ลบหลัง 7 วัน) [4][9]
6. **อุปกรณ์** — Foundry ไม่รองรับ iPad/มือถือ แต่ผู้เล่นใช้มือถือจริง [12][17][30]
7. visual/automation มากไปทำลายจินตนาการ (เสียงส่วนน้อย) [24][26]

## 4. ผู้เล่นต้องการ vs DM ต้องการ
**ผู้เล่น**
- เข้าห้องจากลิงก์ ไม่ต้องสมัคร [2][30]
- ใช้ได้บนมือถือ/โน้ตบุ๊ก [3][30]
- เห็นแมพ + ขยับ token ตัวเอง
- ทอยเต๋าที่ทุกคนเห็น [4]
- เห็นตำแหน่งศัตรู, วัดระยะ
- (เสริม) เห็น AC/HP ของตัวเอง [21]

**DM**
- อัปแมพเร็ว + จัด grid ให้ตรง [4][5]
- ซ่อน/เปิดพื้นที่ (fog) [5]
- token ที่ GM เห็นคนเดียว [5]
- copy token มอนสเตอร์ + เลขลำดับอัตโนมัติ [3]
- ติดตาม initiative + HP [25][29]
- กำหนดสิทธิ์ว่าผู้เล่นขยับ/แก้อะไรได้ [6]
- เซฟข้าม session [4]
- สลับฉาก/แมพ [5]

## 5. ลำดับความสำคัญสำหรับ MVP ("แทนโต๊ะกระดาษ" เล่นเซสชันแรก)
| ระดับ | ฟีเจอร์ |
|---|---|
| **ต้องมี** | ลิงก์ห้อง ผู้เล่นไม่ต้องสมัคร · role GM/ผู้เล่น · อัปรูปแมพ + grid สี่เหลี่ยม · ลาก token sync สด · fog สี่เหลี่ยมให้ GM ซ่อน/เปิด · ทอยเต๋าแบบ notation (`2d6+3`) + log ทุกคนเห็น · state เซฟข้าม session |
| **ควรมี** | ไม้บรรทัด · ชื่อ + HP บน token · ลิสต์ initiative · token ซ่อน · หลายฉาก · ping ชี้ตำแหน่ง · layout มือถือ |
| **ตัดออกจาก MVP** | character sheet + rules automation · dynamic lighting/walls · เสียง/วิดีโอในตัว (ใช้ Discord) · เพลง · compendium · extension API · เต๋า 3D physics |

## 6. Owlbear Rodeo เจาะลึก (ต้นแบบที่ใกล้สุด)
- จุดยืน: "ประสบการณ์ battle map ที่ดี", เครื่องมือเลียนแบบโต๊ะจริง [3]; ฐานเบา ที่เหลือเป็น extension [4]
- **ห้อง**: GM สร้างห้อง → URL ถาวร (ชื่อห้องอยู่ใน URL เปลี่ยนไม่ได้) → ผู้เล่นเปิดลิงก์ "ขอเข้า" → GM อนุมัติ [2] (1.0 ใช้รหัสผ่าน [4])
- ผู้เล่นไม่มีบัญชีใช้ได้ทุกเครื่องมือ ยกเว้นเซฟ/อัปรูป [2][3]
- Role 2 แบบ GM / Player + สิทธิ์ create/edit/delete ต่อเครื่องมือทั้งห้อง [6]
- **1.0 → 2.0**: IndexedDB → cloud; เกมหายหลัง 24 ชม. → ห้องถาวร; แมพเดียว → หลายฉาก canvas ไม่จำกัด; fog เป็น shape + Fill/Cut; ไม้บรรทัด sync; เต๋า + initiative ย้ายไปเป็น extension [4][5]
- **Tech 1.0** [9]: React + Konva, เต๋า 3D Babylon.js, เก็บ Dexie/IndexedDB, socket.io + msgpack, WebRTC (simple-peer) ส่งรูป, Express
- **บทเรียนที่ผู้สร้างเขียนไว้** [9]:
  - React context เป็น "ไอเดียแย่" → ใช้ Zustand
  - ไม่จัดการ conflict ทำให้ 2 คนแก้พร้อมกัน "แมพพัง" → 2.0 ใช้ CRDT
  - เก็บใน IndexedDB browser ไม่เชื่อถือได้
  - WebRTC "ใช้ทีไรเสียใจทุกที" — VPN/ISP บล็อก, TURN แพง
- **Tech 2.0**: AWS serverless หนึ่งเซิร์ฟเวอร์ต่อเกม, cold start ~10–15 วิ [7]; extension ใช้ `@owlbear-rodeo/sdk` (MIT) [10]
- ตั้งใจไม่ทำ: rules engine, character sheet, compendium, เสียง/วิดีโอ [3][4][25]

## 7. Open-source ที่ควรศึกษา
**VTT**
| Repo | License | ดาว | สถานะ | หมายเหตุ |
|---|---|---|---|---|
| owlbear-rodeo/owlbear-rodeo-legacy | ไม่ใช่ OSS (ส่วนตัว non-commercial) | 194 | push ล่าสุด 2023-12-05 | อ่านศึกษาได้ **ห้าม copy โค้ด** [9] |
| samcf/ogres | AGPL-3.0, Clojure | 191 | 2026-10-12 | ไม่ต้องสมัคร, initiative, มือถือ — แนวคิดใกล้สุด [34] |
| Kruptein/PlanarAlly | MIT, TypeScript | 559 | v2026.2 (2026-09-13) | self-host, dynamic lighting [35] |
| dungeon-revealer | ISC, TypeScript | 472 | v1.17.1 (2022) เงียบ | เปิดแมพ, token, dice chat [36] |
| RPTools/maptool | AGPL-3.0, Java | 931 | 1.18.7 (2026-10-09) | desktop |

**Library**
| แพ็กเกจ | เวอร์ชัน (วันที่) | License | ใช้ทำ |
|---|---|---|---|
| konva / react-konva | 10.9.0 / 19.3.2 (2026-10-15) | MIT | canvas แมพ + token |
| pixi.js | 8.22.0 (2026-10-01) | MIT | WebGL ถ้าต้องเร็ว |
| fabric | 7.4.0 (2026-05-18) | MIT | canvas |
| tldraw | 5.5.2 | **ไม่ใช่ open source** (ต้อง license key) [37] | — |
| @dice-roller/rpg-dice-roller | 5.5.1 (2025-02-08) | MIT | parse `2d20kh1+5` |
| @3d-dice/dice-box | 1.1.4 (2024-08-24) | MIT | เต๋า 3D (เก่า) |
| yjs | 13.6.33 (2026-09-23) | MIT | CRDT sync |
| socket.io | 4.8.4 (2026-09-25) | MIT | WebSocket |
| colyseus | 0.18.9 (2026-09-30) | MIT | game server multiplayer |
| partykit | 0.0.115 (2025-05-21) | MIT | room server (Cloudflare) |

## 8. Asset + ลิขสิทธิ์
- **SRD 5.1 และ 5.2.1** (2025-05-01) = CC-BY-4.0 ต้องใส่ attribution ตามข้อความใน SRD [38]
- **รวมในโปรเจกต์ได้**: Kenney (CC0 ไม่ต้องให้เครดิต) [41]; game-icons.net (CC BY 3.0 ~4,180 ไอคอน ให้เครดิต) [40] — เหมาะกับ token/condition
- **ใช้ส่วนตัวเท่านั้น ห้าม bundle**: 2-Minute Tabletop (CC BY-NC 4.0) [42]; Forgotten Adventures (ยกเว้น software) [43]; Tabletop Audio (CC BY-NC-ND, ลิงก์ออกไปปลอดภัยสุด) [39]
- **รูปที่ผู้ใช้อัปเอง**: เว็บสาธารณะในสหรัฐต้องมี DMCA agent ($6, ต่อทุก 3 ปี) [45] — แอปส่วนตัวกลุ่มเพื่อนความเสี่ยงต่ำ แต่ควรมีปุ่มลบรูป + ข้อตกลงว่าผู้อัปต้องมีสิทธิ์ในรูป; Owlbear ให้อัปได้เฉพาะคนมีบัญชี [2]

## ความมั่นใจ
- สูง: ราคา/tier (หน้าทางการ), Owlbear 1.0/2.0, เวอร์ชัน library (npm/GitHub API), license Kenney/game-icons/Tabletop Audio/SRD
- กลาง: แถว chat/handout/compendium ของ Roll20/Foundry/FG, Owlbear 2.3, ความเห็น Reddit, ลำดับ MVP (สรุปเอง), license 2-Minute Tabletop/Forgotten Adventures
- ต่ำ: แถว Tableplop (มือถือ, lighting, บัญชี)

## ยังไม่ได้ตรวจ
- Fantasy Grounds/Tableplop เรื่อง chat, handout, มือถือ, A/V
- Token-Stamp, Gridzzly, Homebrewery, GM Binder (เครื่องมือเตรียมเกม ไม่ใช่ฟีเจอร์บนโต๊ะ)
- ข้อความ attribution ที่แน่นอนใน SRD 5.2.1 — อ่าน PDF ก่อนใส่ข้อความกฎใด ๆ

## แหล่งอ้างอิง
[1] Owlbear Rodeo Pricing — https://www.owlbear.rodeo/pricing
[2] Owlbear docs: Rooms — https://docs.owlbear.rodeo/docs/rooms/
[3] Owlbear docs: Getting Started — https://docs.owlbear.rodeo/docs/getting-started/
[4] Owlbear docs: Migrating 1.0 to 2.0 — https://docs.owlbear.rodeo/docs/migration/migrate-from-1-to-2/
[5] Owlbear docs: Scenes / Images / Fog / Drawing — https://docs.owlbear.rodeo/docs/scenes/
[6] Owlbear 2.0 Dev Log 7 — https://blog.owlbear.rodeo/owlbear-rodeo-dev-log-7/
[7] Owlbear 2.0 Dev Log 4 — https://blog.owlbear.rodeo/owlbear-rodeo-2-0-dev-log-4/
[8] Owlbear 2.3 Release Notes — https://blog.owlbear.rodeo/owlbear-rodeo-2-3-release-notes/
[9] owlbear-rodeo-legacy README — https://github.com/owlbear-rodeo/owlbear-rodeo-legacy
[10] owlbear-rodeo GitHub org — https://github.com/owlbear-rodeo
[11] Foundry Purchase — https://foundryvtt.com/purchase/
[12] Foundry FAQ — https://foundryvtt.com/article/faq/
[13] Foundry Minimum Requirements — https://foundryvtt.com/article/requirements/
[14] Foundry Licence — https://foundryvtt.com/article/license/
[15] Roll20 Feature Breakdown — https://help.roll20.net/hc/en-us/articles/360037774633-Feature-Breakdown
[16] Roll20 — https://roll20.net/
[17] Roll20 Table on Tablet / Mobile FAQ — https://help.roll20.net/hc/en-us/articles/360037774593-Table-on-Tablet ; https://help.roll20.net/hc/en-us/articles/4411213438231-Roll20-Mobile-App-FAQ
[18] Roll20 Turn Tracker / A/V — https://help.roll20.net/hc/articles/360039178634 ; https://wiki.roll20.net/Init ; https://help.roll20.net/hc/articles/360046621893
[19] Fantasy Grounds — https://www.fantasygrounds.com/
[20] FG free-to-play — https://www.gamingonlinux.com/2025/11/fantasy-grounds-virtual-tabletop-vtt-is-now-free-to-play ; https://www.wargamer.com/dnd/fantasy-grounds-free
[21] Tableplop — https://new.tableplop.com/
[22] Tableplop Docs — https://tableplop.notion.site/Tableplop-Docs-a441ac6f11f342b4b0f76e30a309a5b1
[23] r/DMAcademy (2021-06-22) — https://www.reddit.com/r/DMAcademy/comments/o5yfe9/
[24] r/DMAcademy (2024-07-29) — https://www.reddit.com/r/DMAcademy/comments/1eevpdy/
[25] r/DMAcademy (2026-08-02) — https://www.reddit.com/r/DMAcademy/comments/1vds3a3/
[26] r/rpg (2024-08-08) — https://www.reddit.com/r/rpg/comments/1ems8pn/
[27] r/DnD (2022-07-12) — https://www.reddit.com/r/DnD/comments/vx9j95/
[28] r/FoundryVTT (2026-09-18) — https://www.reddit.com/r/FoundryVTT/comments/1wk6be3/
[29] r/DMAcademy (2024-06-26) — https://www.reddit.com/r/DMAcademy/comments/1dowmi3/
[30] r/rpg (2021-04-05) — https://www.reddit.com/r/rpg/comments/mk9ypm/
[31] Wikipedia: Foundry VTT — https://en.wikipedia.org/wiki/Foundry_VTT
[32] GitHub Search API — https://api.github.com/search/repositories
[33] npm registry — https://www.npmjs.com/
[34] samcf/ogres — https://github.com/samcf/ogres
[35] PlanarAlly — https://github.com/Kruptein/PlanarAlly
[36] dungeon-revealer — https://github.com/dungeon-revealer/dungeon-revealer
[37] tldraw licence — https://tldraw.dev/community/license
[38] D&D Beyond SRD — https://www.dndbeyond.com/srd
[39] Tabletop Audio About — https://tabletopaudio.com/about.html
[40] game-icons.net About — https://game-icons.net/about.html
[41] Kenney Support — https://kenney.nl/support
[42] 2-Minute Tabletop / Roll20 Commercial Use — https://2minutetabletop.com/treasure-map-assets ; https://help.roll20.net/hc/en-us/articles/36165245156247-Commercial-Use-Policy
[43] Forgotten Adventures — https://new.forgotten-adventures.net/about/ ; https://www.patreon.com/posts/21728541
[44] Dungeon Scrawl — https://dungeonscrawl.com/
[45] US Copyright Office DMCA Directory — https://www.copyright.gov/dmca-directory/ ; https://www.copyright.gov/dmca-directory/faq.html
