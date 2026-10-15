import { expect, test, type Page } from '@playwright/test'

// 200x150 solid PNG; new scenes get 20 cols, so square cells give 15 rows.
const MAP_PNG = Buffer.from(
  'iVBORw0KGgoAAAANSUhEUgAAAMgAAACWCAIAAAAUvlBOAAABIUlEQVR42u3SMQ0AAAzDsEIrlPFHMQI7d1oygijpFN5FAoyFsTAWGAtjYSwwFsbCWGAsjIWxwFgYC2OBsTAWxgJjYSyMBcbCWBgLjIWxMBYYC2NhLDAWxsJYYCyMhbHAWBgLY4GxMBbGAmNhLIwFxsJYGAuMhbEwFhgLY2EsMBbGwlhgLIyFscBYGAtjgbEwFsYCY2EsjAXGwlgYC4yFsTAWGAtjYSwwFsbCWGAsjIWxwFgYC2NhLBUwFsbCWGAsjIWxwFgYC2OBsTAWxgJjYSyMBcbCWBgLjIWxMBYYC2NhLDAWxsJYYCyMhbHAWBgLY4GxMBbGAmNhLIwFxsJYGAuMhbEwFhgLY2EsMBbGwlhgLIyFscBYGAtjgbEwFsaC2wK9d4/7bDeK+AAAAABJRU5ErkJggg==',
  'base64',
)

const tid = (page: Page, id: string) => page.getByTestId(id)
const MAP_W = 200
const MAP_H = 150

// Page coords of a cell's centre: the taller of map and grid is contain-fit and centred; the grid spans the map width.
async function cellPoint(page: Page, col: number, row: number) {
  const [cols, rows] = (await tid(page, 'map-stage').getAttribute('data-grid'))!.split('x').map(Number)
  const gridH = Math.max(MAP_H, (rows * MAP_W) / cols)
  const box = (await page.locator('[data-testid="map-stage"] .konvajs-content').boundingBox())!
  const s = Math.min(box.width / MAP_W, box.height / gridH)
  const k = (MAP_W * s) / cols
  return { x: box.x + (box.width - MAP_W * s) / 2 + (col + 0.5) * k, y: box.y + (box.height - gridH * s) / 2 + (row + 0.5) * k }
}
const noPageScroll = (page: Page) =>
  page.evaluate(() => {
    const d = document.documentElement
    return d.scrollHeight <= innerHeight && d.scrollWidth <= innerWidth
  })

// shortcut: the server has one global in-memory room, so this needs a fresh server per run (no --repeat-each, no retries); add a reset hook if more specs appear.
test('v0: DM + player full flow', async ({ browser }) => {
  const viewport = { width: 1440, height: 900 }
  const dmCtx = await browser.newContext({ viewport })
  const plCtx = await browser.newContext({ viewport })
  const dm = await dmCtx.newPage()
  const pl = await plCtx.newPage()
  const plFrames: string[] = [] // every WS message the player receives: notes must never be in them
  pl.on('websocket', (ws) => ws.on('framereceived', (f) => plFrames.push(String(f.payload))))
  await dm.goto('/')
  await pl.goto('/')
  await expect(tid(dm, 'ws-status')).toHaveClass(/\bok\b/)
  await expect(tid(pl, 'ws-status')).toHaveClass(/\bok\b/)

  // 1. DM picks the Starter Set.
  await tid(dm, 'dm-name').fill('Dee Em')
  await tid(dm, 'campaign-starter-set').click()
  await expect(tid(dm, 'my-seat')).toContainText('Dee Em')
  await expect(tid(dm, 'attribution')).toContainText('SRD')

  // 1b. Language: default Thai; EN switches labels and survives a reload; back to TH for the rest.
  await expect(tid(pl, 'lobby-heading')).toHaveText('ห้องรอ')
  await expect(tid(pl, 'lang-th')).toHaveAttribute('aria-pressed', 'true')
  await tid(pl, 'lang-en').click()
  await expect(tid(pl, 'lobby-heading')).toHaveText('Lobby')
  await expect(pl.locator('html')).toHaveAttribute('lang', 'en')
  await pl.reload()
  await expect(tid(pl, 'lobby-heading')).toHaveText('Lobby')
  await expect(tid(pl, 'lang-en')).toHaveAttribute('aria-pressed', 'true')
  await tid(pl, 'lang-th').click()
  await expect(tid(pl, 'lobby-heading')).toHaveText('ห้องรอ')
  await expect(pl.locator('html')).toHaveAttribute('lang', 'th')

  // 2. Player sits in a free player seat (add one if none) and readies.
  await expect(tid(pl, 'sit-name')).toBeVisible()
  const sitButtons = pl.locator('[data-testid^="seat-sit-"]')
  if ((await sitButtons.count()) === 0) await tid(dm, 'seat-add').click()
  await tid(pl, 'sit-name').fill('Ayla')
  await sitButtons.first().click()
  await expect(tid(pl, 'my-seat')).toContainText('Ayla')

  // Player never sees DM-only lobby controls.
  for (const id of ['seat-add', 'seat-open-codm', 'game-start', 'ruleset-editor', 'notes-open'])
    await expect(tid(pl, id)).toHaveCount(0)
  await expect(tid(dm, 'notes-open')).toBeVisible()
  await expect(pl.locator('[data-testid^="seat-remove-"]')).toHaveCount(0)
  await expect(pl.locator('[data-testid^="seat-qr-"]')).toHaveCount(0)

  // 2b. Seat QR: DM opens the QR for an empty player seat; a second player joins via its link.
  if ((await dm.locator('button[data-testid^="seat-qr-"]').count()) === 0) await tid(dm, 'seat-add').click()
  const qrBtn = dm.locator('button[data-testid^="seat-qr-"]').first()
  const qrSeatId = (await qrBtn.getAttribute('data-testid'))!.replace('seat-qr-', '')
  const seatLabel = (await tid(dm, `seat-num-${qrSeatId}`).textContent())!
  await qrBtn.click()
  await expect(dm.getByRole('dialog')).toBeVisible()
  await expect(tid(dm, 'seat-qr-label')).toHaveText(seatLabel)
  await expect(tid(dm, 'seat-qr-code')).toBeVisible()
  await expect(tid(dm, 'seat-qr-url')).toContainText(`/?seat=${qrSeatId}`)
  await expect(tid(dm, 'seat-qr-warn')).toBeVisible() // test server binds 127.0.0.1: no LAN URL
  // A click in the dialog's padding (inside its box) must not close it.
  const qrBox = (await dm.getByRole('dialog').boundingBox())!
  await dm.mouse.click(qrBox.x + 4, qrBox.y + qrBox.height / 2)
  await expect(dm.getByRole('dialog')).toBeVisible()
  await dm.keyboard.press('Escape')
  await expect(dm.getByRole('dialog')).toHaveCount(0)
  await expect(qrBtn).toBeFocused()
  await qrBtn.click() // left open: it must close itself once the QR player sits
  await expect(dm.getByRole('dialog')).toBeVisible()

  const qrCtx = await browser.newContext()
  const qr = await qrCtx.newPage()
  await qr.goto(`/?seat=${qrSeatId}`)
  await expect(tid(qr, 'qr-join')).toContainText(seatLabel.split(' / ')[0])
  await expect(tid(qr, 'sit-name')).toBeFocused()
  await qr.keyboard.type('Bram')
  await tid(qr, 'qr-sit').click()
  await expect(tid(qr, 'my-seat')).toContainText('Bram')
  await expect(tid(dm, `seat-${qrSeatId}`)).toContainText('Bram')
  await expect(dm.getByRole('dialog')).toHaveCount(0)
  await expect(tid(dm, 'lobby-heading')).toBeFocused()
  await expect(qr).toHaveURL(/\/$/) // ?seat= stripped once seated
  // Same link again now that the seat is taken: notice, normal lobby.
  const late = await browser.newContext()
  const latePage = await late.newPage()
  await latePage.goto(`/?seat=${qrSeatId}`)
  await expect(tid(latePage, 'qr-seat-gone')).toBeVisible()
  await expect(tid(latePage, 'qr-join')).toHaveCount(0)
  await late.close()
  await tid(qr, 'seat-leave').click()
  await expect(tid(qr, 'qr-join')).toHaveCount(0)
  await qrCtx.close()

  await tid(pl, 'ready-toggle').click()
  await expect(tid(pl, 'ready-toggle')).toHaveAttribute('aria-pressed', 'true')

  // 3. DM starts.
  await expect(tid(dm, 'game-start')).toBeEnabled()
  await tid(dm, 'game-start').click()

  // A seat link after start: banner, no join form.
  const lateCtx = await browser.newContext()
  const lateStart = await lateCtx.newPage()
  await lateStart.goto(`/?seat=${qrSeatId}`)
  await expect(tid(lateStart, 'qr-not-lobby')).toBeVisible()
  await expect(tid(lateStart, 'qr-join')).toHaveCount(0)
  await lateCtx.close()

  // 4. Player creates a character: spend every point, pick the required skills.
  await expect(tid(pl, 'char-name')).toBeVisible()
  await tid(pl, 'char-name').fill('Ayla')
  const plus = pl.locator('[data-testid^="stat-plus-"]')
  for (let i = 0; i < (await plus.count()); i++) {
    const b = plus.nth(i)
    while (await b.isEnabled()) await b.click()
  }
  await expect(tid(pl, 'points-left')).toHaveText('0')
  // Skill card: preset skill shows SRD examples and the d20 how-to; Esc closes and returns focus.
  const info = tid(pl, 'skill-info-stealth')
  await info.click()
  await expect(pl.getByRole('dialog')).toBeVisible()
  await expect(tid(pl, 'skill-info-examples').locator('li').first()).toBeVisible()
  await expect(tid(pl, 'skill-info-how')).toBeVisible()
  await pl.keyboard.press('Escape')
  await expect(pl.getByRole('dialog')).toHaveCount(0)
  await expect(info).toBeFocused()
  const skills = pl.locator('input[data-testid^="skill-"]:enabled:not(:checked)')
  while (await tid(pl, 'char-submit').isDisabled()) await skills.first().check()
  await tid(pl, 'char-submit').click()
  await expect(tid(pl, 'sheet')).toContainText('Ayla')

  // Player never sees the DM panel on the table.
  await expect(tid(pl, 'dm-panel')).toHaveCount(0)
  await expect(tid(pl, 'grant-panel')).toHaveCount(0)
  await expect(tid(pl, 'no-scene')).toBeVisible()

  // 5. DM uploads the PNG as a scene and shows it to everyone.
  await expect(tid(dm, 'dm-panel')).toBeVisible()
  await tid(dm, 'scene-name').fill('Test Map')
  await tid(dm, 'map-upload').setInputFiles({ name: 'map.png', mimeType: 'image/png', buffer: MAP_PNG })
  await expect(tid(dm, 'scene-activate')).toBeEnabled()
  await tid(dm, 'scene-activate').click()
  const sceneId = await tid(dm, 'map-stage').getAttribute('data-scene-id')
  await expect(tid(pl, 'map-stage')).toHaveAttribute('data-scene-id', sceneId!)

  // Table fills the viewport: no page scrollbars on desktop, map pane never scrolls.
  await expect(tid(dm, 'map-stage')).toHaveAttribute('data-grid', '20x15')
  expect(await noPageScroll(dm)).toBe(true)
  expect(await noPageScroll(pl)).toBe(true)

  // 6. Player taps a cell on the canvas; DM sees the new cell. The token row form still works too.
  const pos = pl.locator('[data-testid^="token-pos-"]')
  await expect(pos).toHaveCount(1)
  const tokenId = (await pos.getAttribute('data-testid'))!.replace('token-pos-', '')
  await expect(pos).toHaveText('(0, 0)')
  let p = await cellPoint(pl, 2, 1)
  await pl.mouse.click(p.x, p.y)
  await expect(tid(dm, `token-pos-${tokenId}`)).toHaveText('(2, 1)')
  await expect(tid(pl, `token-pos-${tokenId}`)).toHaveText('(2, 1)') // the row remounts on its own state: fill after that
  await tid(pl, `token-col-${tokenId}`).fill('3')
  await tid(pl, `token-row-${tokenId}`).fill('2')
  await tid(pl, `token-move-${tokenId}`).click()
  await expect(tid(dm, `token-pos-${tokenId}`)).toHaveText('(3, 2)')

  // 6b. DM taps the token, then a cell.
  await dm.evaluate(() => new Promise(requestAnimationFrame)) // Konva redraws its hit canvas on the next frame
  p = await cellPoint(dm, 3, 2)
  await dm.mouse.click(p.x, p.y)
  await expect(tid(dm, 'map-stage')).toHaveAttribute('data-selected-token', tokenId)
  p = await cellPoint(dm, 2, 1)
  await dm.mouse.click(p.x, p.y)
  await expect(tid(pl, `token-pos-${tokenId}`)).toHaveText('(2, 1)')

  // 6c. DM changes the column count with +/-; rows follow the image aspect; save reaches the player.
  await expect(tid(dm, 'grid-cols')).toHaveText('20')
  await expect(tid(dm, 'grid-rows')).toHaveText('15')
  await tid(dm, 'grid-cols-dec').click()
  await tid(dm, 'grid-cols-dec').click()
  await expect(tid(dm, 'grid-cols')).toHaveText('18')
  await expect(tid(dm, 'grid-rows')).toHaveText('14') // round(150 * 18 / 200) = 13.5 -> 14
  await expect(tid(dm, 'map-stage')).toHaveAttribute('data-grid', '18x14') // live preview
  await expect(tid(pl, 'map-stage')).toHaveAttribute('data-grid', '20x15')
  await tid(dm, 'scene-save').click()
  await expect(tid(pl, 'map-stage')).toHaveAttribute('data-grid', '18x14')
  await tid(dm, 'grid-cols-inc').click()
  await tid(dm, 'grid-cols-inc').click()
  await expect(tid(dm, 'grid-rows')).toHaveText('15')
  await expect(tid(dm, 'grid-cols-inc')).toBeEnabled() // 21 cols -> 10px cells, still valid
  await tid(dm, 'scene-save').click()
  await expect(tid(pl, 'map-stage')).toHaveAttribute('data-grid', '20x15')

  // 6d. Splitter: drag left widens the map; keyboard steps 2%.
  const canvas = dm.locator('[data-testid="map-stage"] .konvajs-content')
  const before = (await canvas.boundingBox())!.width
  const sp = (await tid(dm, 'splitter').boundingBox())!
  await dm.mouse.move(sp.x + sp.width / 2, sp.y + sp.height / 2)
  await dm.mouse.down()
  await dm.mouse.move(viewport.width * 0.6, sp.y + sp.height / 2, { steps: 5 })
  await dm.mouse.up()
  await expect.poll(async () => (await canvas.boundingBox())!.width).toBeLessThan(before - 50)
  const now = Number(await tid(dm, 'splitter').getAttribute('aria-valuenow'))
  await tid(dm, 'splitter').focus()
  await dm.keyboard.press('ArrowRight')
  await expect(tid(dm, 'splitter')).toHaveAttribute('aria-valuenow', String(now + 2))
  expect(await noPageScroll(dm)).toBe(true)

  // 7. DM drags a fog rect over column 3 (rows 0-2); player's view gets it.
  await tid(dm, 'dm-menu-fog').click()
  await expect(tid(dm, 'dm-menu-scenes')).toHaveAttribute('aria-expanded', 'false') // one section at a time
  await tid(dm, 'fog-toggle').click()
  await expect(tid(dm, 'fog-toggle')).toHaveAttribute('aria-pressed', 'true')
  const from = await cellPoint(dm, 3, 0)
  const to = await cellPoint(dm, 3, 2)
  await dm.mouse.move(from.x, from.y)
  await dm.mouse.down()
  await dm.mouse.move(to.x, (from.y + to.y) / 2, { steps: 3 })
  await dm.mouse.move(to.x, to.y, { steps: 3 })
  await dm.mouse.up()
  await expect(tid(dm, 'map-stage')).toHaveAttribute('data-fog-count', '1')
  await expect(tid(dm, 'dm-panel')).toContainText('3,0 1×3')
  await expect(tid(pl, 'map-stage')).toHaveAttribute('data-fog-count', '1')

  // 7b. Story notes: DM writes a note, autosave survives a reload; the player gets neither the menu nor the data.
  await tid(dm, 'dm-menu-notes').click()
  await expect(tid(dm, 'dm-menu-notes')).toHaveAttribute('aria-expanded', 'true')
  await tid(dm, 'notes-add').click()
  await expect(tid(dm, 'note-title')).toBeFocused()
  await tid(dm, 'note-title').fill('Chapter 1')
  await tid(dm, 'note-body').fill('The dragon sleeps under the lighthouse.')
  await expect(tid(dm, 'note-status')).toHaveText('บันทึกแล้ว')
  await dm.reload()
  await expect(tid(dm, 'dm-menu-notes')).toHaveAttribute('aria-expanded', 'true') // last section remembered
  const noteItem = dm.locator('[data-testid^="note-item-"]')
  await expect(noteItem).toHaveText(['Chapter 1'])
  await noteItem.click()
  await expect(tid(dm, 'note-body')).toHaveValue('The dragon sleeps under the lighthouse.')
  await tid(dm, 'note-back').click()
  await expect(noteItem).toBeFocused()
  await expect(tid(pl, 'dm-menu-notes')).toHaveCount(0)
  expect(plFrames.length).toBeGreaterThan(0)
  expect(plFrames.filter((f) => f.includes('"notes"'))).toEqual([])

  // Phone: stacked, the map pane still never scrolls and the page has no sideways scroll.
  await pl.setViewportSize({ width: 390, height: 844 })
  await expect
    .poll(() =>
      pl.evaluate(() => {
        const st = document.querySelector('[data-testid="map-stage"]')!
        return st.scrollHeight <= st.clientHeight && st.scrollWidth <= st.clientWidth && document.documentElement.scrollWidth <= innerWidth
      }),
    )
    .toBe(true)

  await expect(tid(dm, 'error-banner')).toHaveCount(0)
  await expect(tid(pl, 'error-banner')).toHaveCount(0)

  // 8. Touch phone: the same player (same clientId) taps a cell with a finger; the token moves.
  const touchCtx = await browser.newContext({ viewport: { width: 390, height: 844 }, hasTouch: true, storageState: await plCtx.storageState() })
  await plCtx.close()
  const ph = await touchCtx.newPage()
  await ph.goto('/')
  await expect(tid(ph, 'map-stage')).toHaveAttribute('data-scene-id', sceneId!)
  await expect(tid(ph, `token-pos-${tokenId}`)).toHaveText('(2, 1)')
  p = await cellPoint(ph, 5, 4)
  await ph.touchscreen.tap(p.x, p.y)
  await expect(tid(dm, `token-pos-${tokenId}`)).toHaveText('(5, 4)')
  await expect(tid(ph, `token-pos-${tokenId}`)).toHaveText('(5, 4)')
  await dmCtx.close()
  await touchCtx.close()
})
