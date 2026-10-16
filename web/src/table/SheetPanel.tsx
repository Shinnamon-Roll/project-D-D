import type { ReactNode } from 'react'
import type { Character } from '../protocol'
import { useStore } from '../store'
import { useL, useT, type TKey } from '../i18n'
import { CustomSkills, Medallion, SkillList } from './Skills'

export function Meter({ kind, label, value, max, testid = `sheet-${kind}` }: { kind: 'hp' | 'mp'; label: string; value: number; max: number; testid?: string }) {
  return (
    <div className={`meter ${kind}`}>
      <span className="meter-label">{label}</span>
      <span className="meter-track" role="meter" aria-label={label} aria-valuemin={0} aria-valuemax={max} aria-valuenow={value}>
        <span className="meter-fill" style={{ transform: `scaleX(${max ? value / max : 0})` }} />
      </span>
      <output className="meter-num" data-testid={testid}>
        {value}/{max}
      </output>
    </div>
  )
}

// The player's own sheet: HP/mana, stats, innate + own skills with Use, gear.
export default function SheetPanel({ character: c, children }: { character: Character; children?: ReactNode }) {
  const t = useT()
  const tx = useL()
  const entry = useStore((s) => s.room!.roster.find((r) => r.id === c.classId))
  const innate = c.skills.filter((s) => s.kind === 'innate')
  const stats: [TKey, number][] = [
    ['statStr', c.str],
    ['statEnd', c.end],
    ['statAgi', c.agi],
  ]

  return (
    <section className="panel sheet" data-testid="sheet">
      <header className="sheet-head">
        <Medallion classId={c.classId} art={c.art} />
        <div>
          <h2>{tx(c.name)}</h2>
          {entry && <p className="hero-class">{tx(entry.className)}</p>}
        </div>
      </header>
      {children}
      <Meter kind="hp" label={t('statHp')} value={c.hp} max={c.hpMax} />
      <Meter kind="mp" label={t('statMp')} value={c.mp} max={c.mpMax} />
      <dl className="stats-grid">
        {stats.map(([k, v]) => (
          <div key={k}>
            <dt>{t(k)}</dt>
            <dd>{v}</dd>
          </div>
        ))}
      </dl>
      <h3>{t('innateSkills')}</h3>
      <SkillList skills={innate} mp={c.mp} testid="innate-skills" />
      <h3>{t('mySkills')}</h3>
      <CustomSkills c={c} play />
      <h3>{t('items')}</h3>
      <ul className="hero-items" data-testid="sheet-items">
        {c.items.map((it, i) => (
          <li key={i}>
            {tx(it.name)}
            {it.qty > 1 && <span className="qty"> ×{it.qty}</span>}
          </li>
        ))}
      </ul>
    </section>
  )
}
