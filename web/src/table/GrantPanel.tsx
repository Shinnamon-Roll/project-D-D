import { useState } from 'react'
import type { ClientEvents } from '../protocol'
import { useStore } from '../store'
import { send } from '../ws'

// DM: give/take items and skills, overwrite stats (server lets the DM bypass min/max).
export default function GrantPanel() {
  const room = useStore((s) => s.room)!
  const { ruleset, characters } = room
  const [charId, setCharId] = useState('')
  const [skillId, setSkillId] = useState('')
  const [itemId, setItemId] = useState('')
  const [qty, setQty] = useState(1)
  const [statEdits, setStatEdits] = useState<Record<string, number>>({})

  const c = characters.find((x) => x.id === charId) ?? characters[0]
  if (!c) return <p className="muted">ยังไม่มีตัวละคร / No characters yet</p>
  const name = (list: { id: string; name: string }[], id: string) => list.find((x) => x.id === id)?.name ?? id
  const grant = (p: Omit<ClientEvents['character.grant'], 'characterId'>) =>
    send('character.grant', { characterId: c.id, ...p })

  return (
    <div className="grant" data-testid="grant-panel">
      <h2>ให้ของ / Grant</h2>
      <label>
        ตัวละคร / Character
        <select data-testid="grant-character" value={c.id} onChange={(e) => { setCharId(e.target.value); setStatEdits({}) }}>
          {characters.map((x) => <option key={x.id} value={x.id}>{x.name}</option>)}
        </select>
      </label>

      <div className="stats-grid">
        {ruleset.stats.map((s) => (
          <label key={s.key}>
            {s.name}
            <input
              type="number"
              data-testid={`grant-stat-${s.key}`}
              value={statEdits[s.key] ?? c.stats[s.key] ?? s.base}
              onChange={(e) => setStatEdits({ ...statEdits, [s.key]: Number(e.target.value) })}
            />
          </label>
        ))}
      </div>
      <button
        data-testid="grant-set-stats"
        disabled={Object.keys(statEdits).length === 0}
        onClick={() => { grant({ setStats: statEdits }); setStatEdits({}) }}
      >
        ตั้งค่า stat / Set stats
      </button>

      <div className="row">
        <select aria-label="skill" data-testid="grant-skill" value={skillId} onChange={(e) => setSkillId(e.target.value)}>
          <option value="">— skill —</option>
          {ruleset.skills.filter((s) => !c.skillIds.includes(s.id)).map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
        </select>
        <button data-testid="grant-add-skill" disabled={!skillId} onClick={() => { grant({ addSkillIds: [skillId] }); setSkillId('') }}>
          + skill
        </button>
      </div>
      <ul className="chips">
        {c.skillIds.map((id) => (
          <li key={id}>
            {name(ruleset.skills, id)}
            <button aria-label={`เอาสกิล ${name(ruleset.skills, id)} ออก`} data-testid={`grant-remove-skill-${id}`} onClick={() => grant({ removeSkillIds: [id] })}>
              ×
            </button>
          </li>
        ))}
      </ul>

      <div className="row">
        <select aria-label="item" data-testid="grant-item" value={itemId} onChange={(e) => setItemId(e.target.value)}>
          <option value="">— item —</option>
          {ruleset.items.map((it) => <option key={it.id} value={it.id}>{it.name}</option>)}
        </select>
        <input aria-label="qty" data-testid="grant-qty" type="number" min={1} value={qty} onChange={(e) => setQty(Number(e.target.value))} />
        <button data-testid="grant-add-item" disabled={!itemId || qty < 1} onClick={() => grant({ addItems: [{ itemId, qty }] })}>
          + item
        </button>
        <button data-testid="grant-remove-item" disabled={!itemId || qty < 1} onClick={() => grant({ removeItems: [{ itemId, qty }] })}>
          − item
        </button>
      </div>
      <ul className="chips">
        {c.inventory.map((it) => (
          <li key={it.itemId}>{name(ruleset.items, it.itemId)} × {it.qty}</li>
        ))}
      </ul>
    </div>
  )
}
