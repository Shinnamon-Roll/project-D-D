# Visual theme — "Candlelit Tavern" (medieval)

**สถานะ:** ใช้งาน (v0.4.0) · ที่มา: [`../conversations/2026-10-16-medieval-ui.md`](../conversations/2026-10-16-medieval-ui.md)

## แนวคิด
โต๊ะไม้ในโรงเตี๊ยมยามค่ำ: พื้นหลังไม้/หนังสีเข้ม, เอกสารทุกแผ่นเป็นกระดาษ parchment, หัวข้อเป็นอักษรโรมันจารึก,
ปุ่มหลักเป็นแดงเลือดหมูขอบทอง, สถานะ "พร้อม" เป็นตราประทับขี้ผึ้ง

## Tokens (CSS custom properties ใน `web/src/styles.css`)
| token | ค่า | ใช้กับ |
|---|---|---|
| `--wood` | `#1c140f` | พื้นหลังหน้า |
| `--wood-2` | `#2a1e16` | topbar, กรอบแมพ |
| `--parchment` | `#efe2c3` | แผง (panel) |
| `--parchment-2` | `#e2cfa3` | input, แถวสลับ |
| `--parchment-edge` | `#b8975e` | เส้นขอบ, เส้นคั่น |
| `--ink` | `#2b1d12` | ตัวหนังสือบนกระดาษ |
| `--ink-muted` | `#6b5640` | ข้อความรอง |
| `--oxblood` | `#7a1f1f` | ปุ่มหลัก, ตราขี้ผึ้ง |
| `--oxblood-hi` | `#9b2d2d` | hover |
| `--gold` | `#c9a227` | ขอบปุ่มหลัก, focus ring, token ที่เลือก |
| `--gold-dim` | `#8a6d1d` | ornament |
| `--iron` | `#4a4f55` | token NPC, ปุ่มรอง |
| `--moss` | `#4f6b3a` | สถานะ ok |
| `--blood` | `#a12a2a` | error/danger |
| `--mana` | `#2f4f7a` | แถบมานา (สีเย็นสีเดียวในพาเลต: น้ำหมึกน้ำเงิน) — ใช้เป็นต้นของ gradient |
| `--mana-hi` | `#4d77ad` | ปลาย gradient แถบมานา, เส้นมานาเล็กใต้รูปปาร์ตี้ |

Contrast: `--ink` บน `--parchment` ≈ 12:1, ข้อความ `--parchment` บน `--oxblood` ≈ 7:1 (ผ่าน WCAG AA)

## ฟอนต์
- หัวข้อ: **Cinzel** (Latin) → fallback อัตโนมัติเป็น Trirong สำหรับอักษรไทย
- เนื้อความ: **Trirong** (serif ไทยแบบคลาสสิก รองรับอังกฤษ)
- โหลดจาก Google Fonts ใน `web/index.html`; ไม่มีเน็ต → `Georgia, serif`

## ภาพประกอบ (SVG วาดเอง ไม่มีลิขสิทธิ์ภายนอก — `web/src/assets/`)
| ไฟล์ | ใช้ที่ |
|---|---|
| `crest.svg` | โลโก้ topbar + หน้าเมนู: โล่ + เต๋า 20 หน้า + ดาบไขว้ |
| `stormwreck.svg` | การ์ดแคมเปญ Starter Set: ทะเลพายุ เกาะ ประภาคาร (ภาพต้นฉบับ ไม่อิงอาร์ตจากเล่ม) |
| `salem.svg` | การ์ด Salem (ล็อก): หมู่บ้านใต้พระจันทร์ เทียน + ตราขี้ผึ้ง "Coming soon" |
| `divider.svg` | เส้นคั่นลายเถาใต้หัวข้อ |
| `tex-noise.webp` `tex-stains.webp` `tex-grain.webp` | ผิวกระดาษ/ไม้ — pre-render จาก SVG `feTurbulence` เป็น WebP แล้ว (ห้ามกลับไปใช้ filter แบบ runtime: ทำให้ repaint หนักมาก ดู perf ใน WORKLOG 0.6.1) |

## แมพ (Konva)
grid สีซีเปีย `rgba(43,29,18,.45)`, token ผู้เล่นขอบ `--parchment`, เลือกแล้วขอบ `--gold`, token NPC พื้น `--iron`,
fog ผู้เล่นเห็นเป็นหมอกดำอมน้ำตาล `#120d09`, DM เห็นโปร่ง + ขอบ `--oxblood`

## ห้าม
อาร์ต/โลโก้ WotC, ชื่อ "D&D" เป็นแบรนด์, ฟอนต์/ลายจากชีทตัวละครทางการ

## Performance rules
- ห้ามใช้ SVG filter (`feTurbulence` ฯลฯ) เป็น background แบบ runtime → ใช้ภาพ raster ที่ render ไว้แล้ว
- ห้าม `background-attachment: fixed` → ใช้ `body::before { position: fixed }` แยก layer
