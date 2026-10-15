import { expect, test, type Page } from '@playwright/test'

// 200x150 solid PNG; with the default 50px grid the scene is 4 cols x 3 rows.
const MAP_PNG = Buffer.from(
  'iVBORw0KGgoAAAANSUhEUgAAAMgAAACWCAIAAAAUvlBOAAABIUlEQVR42u3SMQ0AAAzDsEIrlPFHMQI7d1oygijpFN5FAoyFsTAWGAtjYSwwFsbCWGAsjIWxwFgYC2OBsTAWxgJjYSyMBcbCWBgLjIWxMBYYC2NhLDAWxsJYYCyMhbHAWBgLY4GxMBbGAmNhLIwFxsJYGAuMhbEwFhgLY2EsMBbGwlhgLIyFscBYGAtjgbEwFsYCY2EsjAXGwlgYC4yFsTAWGAtjYSwwFsbCWGAsjIWxwFgYC2NhLBUwFsbCWGAsjIWxwFgYC2OBsTAWxgJjYSyMBcbCWBgLjIWxMBYYC2NhLDAWxsJYYCyMhbHAWBgLY4GxMBbGAmNhLIwFxsJYGAuMhbEwFhgLY2EsMBbGwlhgLIyFscBYGAtjgbEwFsaC2wK9d4/7bDeK+AAAAABJRU5ErkJggg==',
  'base64',
)

const tid = (page: Page, id: string) => page.getByTestId(id)

// shortcut: the server has one global in-memory room, so this needs a fresh server per run (no --repeat-each, no retries); add a reset hook if more specs appear.
test('v0: DM + player full flow', async ({ browser }) => {
  const dmCtx = await browser.newContext()
  const plCtx = await browser.newContext()
  const dm = await dmCtx.newPage()
  const pl = await plCtx.newPage()
  await dm.goto('/')
  await pl.goto('/')
  await expect(tid(dm, 'ws-status')).toHaveText('connected')
  await expect(tid(pl, 'ws-status')).toHaveText('connected')

  // 1. DM picks the Starter Set.
  await tid(dm, 'dm-name').fill('Dee Em')
  await tid(dm, 'campaign-starter-set').click()
  await expect(tid(dm, 'my-seat')).toContainText('Dee Em')
  await expect(tid(dm, 'attribution')).toContainText('SRD')

  // 2. Player sits in a free player seat (add one if none) and readies.
  await expect(tid(pl, 'sit-name')).toBeVisible()
  const sitButtons = pl.locator('[data-testid^="seat-sit-"]')
  if ((await sitButtons.count()) === 0) await tid(dm, 'seat-add').click()
  await tid(pl, 'sit-name').fill('Ayla')
  await sitButtons.first().click()
  await expect(tid(pl, 'my-seat')).toContainText('Ayla')

  // Player never sees DM-only lobby controls.
  for (const id of ['seat-add', 'seat-open-codm', 'game-start', 'ruleset-editor'])
    await expect(tid(pl, id)).toHaveCount(0)
  await expect(pl.locator('[data-testid^="seat-remove-"]')).toHaveCount(0)

  await tid(pl, 'ready-toggle').click()
  await expect(tid(pl, 'ready-toggle')).toHaveAttribute('aria-pressed', 'true')

  // 3. DM starts.
  await expect(tid(dm, 'game-start')).toBeEnabled()
  await tid(dm, 'game-start').click()

  // 4. Player creates a character: spend every point, pick the required skills.
  await expect(tid(pl, 'char-name')).toBeVisible()
  await tid(pl, 'char-name').fill('Ayla')
  const plus = pl.locator('[data-testid^="stat-plus-"]')
  for (let i = 0; i < (await plus.count()); i++) {
    const b = plus.nth(i)
    while (await b.isEnabled()) await b.click()
  }
  await expect(tid(pl, 'points-left')).toHaveText('0')
  const skills = pl.locator('[data-testid^="skill-"]:enabled:not(:checked)')
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

  // 6. Player moves their token; DM sees the new cell.
  const pos = pl.locator('[data-testid^="token-pos-"]')
  await expect(pos).toHaveCount(1)
  const tokenId = (await pos.getAttribute('data-testid'))!.replace('token-pos-', '')
  await expect(pos).toHaveText('(0, 0)')
  await pl.getByLabel('Ayla col').fill('2')
  await pl.getByLabel('Ayla row').fill('1')
  await tid(pl, `token-move-${tokenId}`).click()
  await expect(tid(dm, `token-pos-${tokenId}`)).toHaveText('(2, 1)')

  // 7. Player rolls 1d20+3; DM's log shows it.
  await tid(pl, 'dice-input').fill('1d20+3')
  await tid(pl, 'dice-roll').click()
  const entry = tid(dm, 'dice-log').locator('li').first()
  await expect(entry).toContainText('Ayla')
  await expect(entry).toContainText('1d20+3')
  const total = Number(await entry.getByTestId('dice-total').textContent())
  expect(total).toBeGreaterThanOrEqual(4)
  expect(total).toBeLessThanOrEqual(23)

  // 8. DM drags a fog rect over column 3 (rows 0-2); player's view gets it.
  await tid(dm, 'fog-toggle').click()
  await expect(tid(dm, 'fog-toggle')).toHaveAttribute('aria-pressed', 'true')
  const content = dm.locator('[data-testid="map-stage"] .konvajs-content')
  await content.scrollIntoViewIfNeeded() // mouse events only land inside the viewport
  const box = (await content.boundingBox())!
  await dm.mouse.move(box.x + 175, box.y + 25)
  await dm.mouse.down()
  await dm.mouse.move(box.x + 175, box.y + 75, { steps: 3 })
  await dm.mouse.move(box.x + 175, box.y + 125, { steps: 3 })
  await dm.mouse.up()
  await expect(tid(dm, 'map-stage')).toHaveAttribute('data-fog-count', '1')
  await expect(tid(dm, 'dm-panel')).toContainText('3,0 1×3')
  await expect(tid(pl, 'map-stage')).toHaveAttribute('data-fog-count', '1')

  await expect(tid(dm, 'error-banner')).toHaveCount(0)
  await expect(tid(pl, 'error-banner')).toHaveCount(0)
  await dmCtx.close()
  await plCtx.close()
})
