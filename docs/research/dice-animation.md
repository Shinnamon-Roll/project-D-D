# Research: animation ทอยลูกเต๋า (2026-10-16)

> **ไม่ใช้แล้ว** — ผู้ใช้ตัดสินใจใช้ลูกเต๋าจริง ตัดระบบเต๋าออกจากแอป ([conv](../conversations/2026-10-16-remove-dice.md)) · เก็บไว้อ้างอิง

ที่มา: [`../conversations/2026-10-16-perf-and-dice-anim.md`](../conversations/2026-10-16-perf-and-dice-anim.md) · agent `researcher` · ขนาด bundle วัดจริงด้วย esbuild + `gzip -9`, วันที่/license จาก `npm view`

## สรุปสั้น
- **มี 3D ที่บังคับผลได้**: `@3d-dice/dice-box-threejs` — `Box.roll("6d6@4,4,4,4,4,4")` หยุดที่หน้าที่กำหนด [1] แต่ ~150 KB gzip, release ล่าสุด 2022, physics รันบน main thread, ไม่มี d3/d7/d1000
- **แนะนำหลัก: ทำ overlay 2D เอง** (DOM + CSS transform / Web Animations API): รูปทรงเต๋าตามจำนวนหน้า หมุน+สั่น ~0.9s ข้างในมี "วงล้อ/รีล" ตัวเลขเลื่อนแล้วหยุดที่ค่าจาก server — 0 dependency (~2–4 KB ประมาณ), รองรับ d2–d1000, ไม่แตะ Konva จึงไม่ทำแมพกระตุก
- **ทางเลือกเสริม: โหมด 3D แบบเปิดเอง** — โหลด `dice-box-threejs` แบบ lazy เฉพาะเมื่อผู้ใช้เปิด, ใช้เฉพาะเต๋ามาตรฐาน ≤ ~10 ลูก

## เทียบ
| ตัวเลือก | บังคับผลได้ | ชนิดเต๋า | ขนาด gzip | license | release ล่าสุด | สรุป |
|---|---|---|---|---|---|---|
| `@3d-dice/dice-box` (Babylon + Ammo) | **ไม่ได้** | d4–d20, d100 | ~450–610 KB+ | MIT | 1.1.4 (2024-08) | ไม่เหมาะ |
| `@3d-dice/dice-box-threejs` | **ได้** `NdM@v1,v2` | d4–d20, d100 | ~145–156 KB | MIT | 0.0.12 (2022-10) | ใช้เป็นโหมด 3D ได้ |
| `@lambersond/3d-dice-core` (fork) | ได้ | 4–100 | 2 KB + 177 KB lazy | MIT | 0.4.0 (2026-06) | ใหม่เกินไป รอดูก่อน |
| `react-3d-dice` | ได้ (`results`) | d4–d20 | ~140 KB | MIT | 0.1.0 (2026-02) | ใหม่เกินไป |
| three + cannon-es เขียนเอง | ทำเอง | ทุกแบบ | ≥160 KB | MIT | — | งานเยอะเกิน |
| Owlbear Rodeo dice | ไม่ (จำลองเหมือนกันทุกเครื่อง) | — | — | **GPL-3.0** | — | license ไม่ได้ |
| `react-dice-roll` / `react-dice-complete` | ได้ | **d6 เท่านั้น** | เล็ก | MIT/ISC | 2021 / 2026 | ไม่พอ |
| CSS 3D ล้วน | ทำเอง | d6 ง่าย, d20 ยาก | 1–3 KB | ของเรา | — | d6 เท่านั้น |
| Lottie | ไม่ (เล่นคลิปแล้วซ้อนตัวเลข) | ทุกแบบ | 46–76 KB + ไฟล์ animation | MIT | 2025-05 | ไม่คุ้ม |
| วงล้อ `react-custom-roulette` | ได้ (`prizeNumber`) | d4–d20 ok, d100 แน่น, d1000 ไม่ได้ | ~131 KB unpacked | MIT | 2023 (peer React 18) | ไม่ขยาย |
| **overlay 2D เขียนเอง** | **ได้เสมอ** | **d2–d1000** | **~2–4 KB (ประมาณ)** | ของเรา | — | **แนะนำ** |

## แผน overlay 2D (ร่าง)
- `web/src/table/DiceOverlay.tsx` mount ครั้งเดียวในหน้า Table; `position: fixed; inset: 0; pointer-events: none` (ขณะเล่นให้คลิก/แตะ/Esc เพื่อข้าม)
- trigger: เมื่อมี roll id ใหม่ใน `room.dice` — จำ id ล่าสุดตอน mount เพื่อ **ไม่เล่นย้อนประวัติ 200 รายการ** ตอนเข้าห้อง/ต่อใหม่
- ~1.0–1.2s แล้ว fade 0.3s; roll ใหม่มาระหว่างเล่น → จบอันเก่าทันที (ไม่ต่อคิว); แสดงสูงสุด 12 ลูก + "+N"
- kh/kl: ลูกที่ไม่นับ จาง + ขีดฆ่า แล้วโชว์ผลรวม
- `prefers-reduced-motion` → โชว์ผลทันที ไม่หมุน
- **ต้องแก้ protocol ก่อน**: `DiceRoll.rolls[]` ไม่บอกว่าแต่ละลูกกี่หน้า/ถูกนับไหม → เพิ่ม `groups: [{sides, rolls, kept[]}]` (เพิ่ม field ไม่ทำของเดิมพัง, server มีข้อมูลอยู่แล้ว) — ไม่งั้น client ต้อง parse notation ซ้ำกับ server

## ความเสี่ยง
1. protocol ไม่มี sides/kept → แก้ก่อนทำ
2. ย้อนเล่นประวัติตอนต่อใหม่ → จำ id ตอน mount
3. overlay บังแมพบนมือถือ → `pointer-events: none` เมื่อจบ + ซ่อนเอง
4. โหมด 3D: จำลองแบบ synchronous บน main thread (ยังไม่วัด), lib ไม่อัปเดตตั้งแต่ 2022, d100/mixed notation อ่านจาก source ไม่ได้ทดสอบ, WebGL กินแบต
5. log แสดงผลรวมทันทีก่อน animation จบ (สปอยล์) → ตัดสินใจว่าจะหน่วงไหม

## ความมั่นใจ
- dice-box-threejs บังคับผลด้วย `@`: สูง (README + อ่านโค้ด `simulateThrow` → `swapDiceFace`)
- dice-box (Babylon) ไม่มี forced result: สูง
- overlay 2D ลื่นบนมือถือ: กลาง (ประเมินจากเทคนิค ยังไม่ได้สร้าง/วัด)

## แหล่งที่มา
[1] https://github.com/3d-dice/dice-box-threejs · [2] https://codesandbox.io/s/dice-box-threejs-j79h35 · [3] https://www.npmjs.com/package/@3d-dice/dice-box-threejs (และ npm ของแต่ละแพ็กเกจ) · [4] https://www.npmjs.com/package/@lambersond/3d-dice-core · [5] https://fantasticdice.games/docs/usage/config · [6] https://fantasticdice.games/blog · [7] https://github.com/3d-dice/dice-box · [8] https://h.daily-dev-tips.com/creating-a-3d-dice-in-css.md · [9] https://codepen.io/Hewitt/pen/ymJoRr · [10] https://codepen.io/rjpvroegop/pen/BabBQKQ · [11] https://www.npmjs.com/package/react-custom-roulette · [12] https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-reduced-motion · [13] https://github.com/owlbear-rodeo/dice
