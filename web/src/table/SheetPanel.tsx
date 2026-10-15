import type { Character } from '../protocol'
import { useStore } from '../store'

export default function SheetPanel({ character: c }: { character: Character }) {
  const ruleset = useStore((s) => s.room!.ruleset)
  const name = (list: { id: string; name: string }[], id: string) => list.find((x) => x.id === id)?.name ?? id

  return (
    <section className="panel" data-testid="sheet">
      <h2 className="row">
        {c.portraitUrl && <img className="avatar" src={c.portraitUrl} alt="" />}
        {c.name}
      </h2>
      <dl className="stats-grid">
        {ruleset.stats.map((s) => (
          <div key={s.key}>
            <dt>{s.name}</dt>
            <dd data-testid={`sheet-stat-${s.key}`}>{c.stats[s.key] ?? '—'}</dd>
          </div>
        ))}
      </dl>
      <h3>สกิล / Skills</h3>
      <ul className="chips">
        {c.skillIds.map((id) => <li key={id}>{name(ruleset.skills, id)}</li>)}
      </ul>
      <h3>ไอเทม / Items</h3>
      <ul data-testid="sheet-items">
        {c.inventory.map((it) => (
          <li key={it.itemId}>
            {name(ruleset.items, it.itemId)} × {it.qty}
          </li>
        ))}
      </ul>
    </section>
  )
}
