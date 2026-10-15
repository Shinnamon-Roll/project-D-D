import { useState } from 'react'
import type { Ruleset } from '../protocol'
import { useStore } from '../store'
import { send } from '../ws'

// Parent keys this by the server ruleset, so the draft resets whenever the server copy changes.
export default function RulesetEditor() {
  const [r, setR] = useState<Ruleset>(() => structuredClone(useStore.getState().room!.ruleset))
  const set = (patch: Partial<Ruleset>) => setR({ ...r, ...patch })
  function edit<K extends 'stats' | 'skills' | 'items' | 'startingItems'>(list: K, i: number, patch: Partial<Ruleset[K][number]>) {
    set({ [list]: r[list].map((x, j) => (j === i ? { ...x, ...patch } : x)) })
  }
  const drop = (list: 'stats' | 'skills' | 'items' | 'startingItems', i: number) =>
    set({ [list]: r[list].filter((_, j) => j !== i) })
  const num = (v: string) => (v === '' ? 0 : Number(v))

  return (
    <section className="panel" data-testid="ruleset-editor">
      <h2>ชุดกฎ / Ruleset</h2>
      <div className="row">
        <label>
          ชื่อ / Name
          <input data-testid="ruleset-name" value={r.name} onChange={(e) => set({ name: e.target.value })} />
        </label>
        <button data-testid="ruleset-load-preset" onClick={() => send('ruleset.loadPreset', { presetId: 'dnd5e-srd' })}>
          โหลด preset 5E (SRD 5.2.1) / Load 5E preset
        </button>
      </div>

      <h3>Stats</h3>
      <table>
        <thead>
          <tr>
            <th>key</th><th>ชื่อ / name</th><th>min</th><th>max</th><th>base</th><th />
          </tr>
        </thead>
        <tbody>
          {r.stats.map((s, i) => (
            <tr key={i}>
              <td><input aria-label="stat key" value={s.key} onChange={(e) => edit('stats', i, { key: e.target.value })} /></td>
              <td><input aria-label="stat name" value={s.name} onChange={(e) => edit('stats', i, { name: e.target.value })} /></td>
              <td><input aria-label="min" type="number" value={s.min} onChange={(e) => edit('stats', i, { min: num(e.target.value) })} /></td>
              <td><input aria-label="max" type="number" value={s.max} onChange={(e) => edit('stats', i, { max: num(e.target.value) })} /></td>
              <td><input aria-label="base" type="number" value={s.base} onChange={(e) => edit('stats', i, { base: num(e.target.value) })} /></td>
              <td><button aria-label={`ลบ ${s.name}`} onClick={() => drop('stats', i)}>×</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      <button data-testid="ruleset-add-stat" onClick={() => set({ stats: [...r.stats, { key: '', name: '', min: 0, max: 10, base: 0 }] })}>
        + stat
      </button>
      <label>
        แต้ม stat ตอนสร้าง / Stat points
        <input data-testid="ruleset-stat-points" type="number" min={0} value={r.statPoints} onChange={(e) => set({ statPoints: num(e.target.value) })} />
      </label>

      <h3>Skills</h3>
      <table>
        <tbody>
          {r.skills.map((s, i) => (
            <tr key={i}>
              <td><input aria-label="skill id" value={s.id} onChange={(e) => edit('skills', i, { id: e.target.value })} /></td>
              <td><input aria-label="skill name" value={s.name} onChange={(e) => edit('skills', i, { name: e.target.value })} /></td>
              <td>
                <select aria-label="skill stat" value={s.statKey} onChange={(e) => edit('skills', i, { statKey: e.target.value })}>
                  <option value="">—</option>
                  {r.stats.map((st) => <option key={st.key} value={st.key}>{st.name || st.key}</option>)}
                </select>
              </td>
              <td><button aria-label={`ลบ ${s.name}`} onClick={() => drop('skills', i)}>×</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      <button data-testid="ruleset-add-skill" onClick={() => set({ skills: [...r.skills, { id: '', name: '', statKey: '' }] })}>
        + skill
      </button>
      <label>
        เลือกสกิลตอนเริ่ม / Starting skill picks
        <input data-testid="ruleset-skill-picks" type="number" min={0} value={r.startingSkillPicks} onChange={(e) => set({ startingSkillPicks: num(e.target.value) })} />
      </label>

      <h3>Items</h3>
      <table>
        <tbody>
          {r.items.map((it, i) => (
            <tr key={i}>
              <td><input aria-label="item id" value={it.id} onChange={(e) => edit('items', i, { id: e.target.value })} /></td>
              <td><input aria-label="item name" value={it.name} onChange={(e) => edit('items', i, { name: e.target.value })} /></td>
              <td><button aria-label={`ลบ ${it.name}`} onClick={() => drop('items', i)}>×</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      <button data-testid="ruleset-add-item" onClick={() => set({ items: [...r.items, { id: '', name: '' }] })}>
        + item
      </button>

      <h3>ไอเทมเริ่มต้น / Starting items</h3>
      <table>
        <tbody>
          {r.startingItems.map((st, i) => (
            <tr key={i}>
              <td>
                <select aria-label="starting item" value={st.itemId} onChange={(e) => edit('startingItems', i, { itemId: e.target.value })}>
                  <option value="">—</option>
                  {r.items.map((it) => <option key={it.id} value={it.id}>{it.name || it.id}</option>)}
                </select>
              </td>
              <td><input aria-label="qty" type="number" min={1} value={st.qty} onChange={(e) => edit('startingItems', i, { qty: num(e.target.value) })} /></td>
              <td><button aria-label="ลบ" onClick={() => drop('startingItems', i)}>×</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      <button data-testid="ruleset-add-starting-item" onClick={() => set({ startingItems: [...r.startingItems, { itemId: r.items[0]?.id ?? '', qty: 1 }] })}>
        + starting item
      </button>

      <div className="row">
        <button className="primary" data-testid="ruleset-save" onClick={() => send('ruleset.update', { ruleset: r })}>
          บันทึก / Save ruleset
        </button>
      </div>
      {r.attribution && <p className="attribution">{r.attribution}</p>}
    </section>
  )
}
