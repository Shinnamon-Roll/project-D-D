import { useState } from 'react'
import { useStore } from '../store'
import { send, upload } from '../ws'
import GrantPanel from '../table/GrantPanel'
import RulesetEditor from './RulesetEditor'

export default function CharacterCreate() {
  const ruleset = useStore((s) => s.room!.ruleset)
  const [name, setName] = useState('')
  const [portraitUrl, setPortraitUrl] = useState('')
  const [skillIds, setSkillIds] = useState<string[]>([])
  const [stats, setStats] = useState<Record<string, number>>(() =>
    Object.fromEntries(ruleset.stats.map((s) => [s.key, s.base])),
  )

  // UX only; the server re-validates everything.
  const spent = ruleset.stats.reduce((sum, s) => sum + ((stats[s.key] ?? s.base) - s.base), 0)
  const left = ruleset.statPoints - spent
  const picks = ruleset.startingSkillPicks
  const itemName = (id: string) => ruleset.items.find((i) => i.id === id)?.name ?? id
  const canSubmit = name.trim() !== '' && left === 0 && skillIds.length === picks

  const bump = (key: string, d: number) => setStats({ ...stats, [key]: stats[key] + d })
  const toggleSkill = (id: string) =>
    setSkillIds(skillIds.includes(id) ? skillIds.filter((s) => s !== id) : [...skillIds, id])

  return (
    <section className="panel narrow">
      <h1>สร้างตัวละคร / Create character</h1>
      <label>
        ชื่อตัวละคร / Name
        <input data-testid="char-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        รูป / Portrait
        <input
          data-testid="portrait-upload"
          type="file"
          accept="image/png,image/jpeg,image/webp,image/gif"
          onChange={async (e) => {
            const f = e.target.files?.[0]
            if (f) setPortraitUrl((await upload(f)) ?? portraitUrl)
          }}
        />
      </label>
      {portraitUrl && <img className="portrait" src={portraitUrl} alt="portrait" data-testid="portrait-preview" />}

      <h2>
        Stats · แต้มเหลือ / Points left: <span data-testid="points-left">{left}</span>
      </h2>
      <ul className="stats">
        {ruleset.stats.map((s) => (
          <li key={s.key}>
            <span>{s.name}</span>
            <button aria-label={`ลด ${s.name}`} data-testid={`stat-minus-${s.key}`} disabled={stats[s.key] <= s.min} onClick={() => bump(s.key, -1)}>
              −
            </button>
            <output data-testid={`stat-value-${s.key}`}>{stats[s.key]}</output>
            <button aria-label={`เพิ่ม ${s.name}`} data-testid={`stat-plus-${s.key}`} disabled={stats[s.key] >= s.max || left <= 0} onClick={() => bump(s.key, 1)}>
              +
            </button>
          </li>
        ))}
      </ul>

      <h2>
        สกิล / Skills ({skillIds.length}/{picks})
      </h2>
      <div className="checks">
        {ruleset.skills.map((s) => (
          <label key={s.id} className="check">
            <input
              type="checkbox"
              data-testid={`skill-${s.id}`}
              checked={skillIds.includes(s.id)}
              disabled={!skillIds.includes(s.id) && skillIds.length >= picks}
              onChange={() => toggleSkill(s.id)}
            />
            {s.name}
          </label>
        ))}
      </div>

      <h2>ไอเทมเริ่มต้น / Starting items</h2>
      <ul data-testid="starting-items">
        {ruleset.startingItems.map((st) => (
          <li key={st.itemId}>
            {itemName(st.itemId)} × {st.qty}
          </li>
        ))}
      </ul>

      <button
        className="primary"
        data-testid="char-submit"
        disabled={!canSubmit}
        onClick={() => send('character.submit', { name: name.trim(), portraitUrl, stats, skillIds })}
      >
        เสร็จ / Done
      </button>
    </section>
  )
}

// DM while players create: who is done, grant panel, ruleset editor.
export function CreatingDmView() {
  const room = useStore((s) => s.room)!
  return (
    <div className="columns">
      <section className="panel">
        <h1>รอผู้เล่นสร้างตัวละคร / Waiting for characters</h1>
        <ul className="seats">
          {room.seats
            .filter((s) => s.kind === 'player')
            .map((s) => (
              <li key={s.id} data-testid={`creating-status-${s.id}`}>
                <span className="who">{s.name}</span>
                <span className={s.characterId ? 'ok' : 'muted'}>
                  {s.characterId ? 'เสร็จแล้ว / Done' : 'กำลังสร้าง / Creating…'}
                </span>
              </li>
            ))}
        </ul>
        <GrantPanel />
      </section>
      <RulesetEditor key={JSON.stringify(room.ruleset)} />
    </div>
  )
}
