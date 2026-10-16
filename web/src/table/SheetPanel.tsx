import type { ReactNode } from 'react'
import type { Character } from '../protocol'
import { useL, useT } from '../i18n'
import { useClass } from '../rules'
import { Medallion } from './Medallion'
import { SheetStats, SheetTabs } from './CharSheet'

export function Meter({ kind, label, value, max, testid = `sheet-${kind}` }: { kind: 'hp'; label: string; value: number; max: number; testid?: string }) {
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

// The player's own sheet: header (who, HP, AC, PB, initiative, passive Perception) and the rules v2 tabs.
export default function SheetPanel({ character: c, children }: { character: Character; children?: ReactNode }) {
  const t = useT()
  const tx = useL()
  const cls = useClass(c.classId)

  return (
    <section className="panel sheet" data-testid="sheet">
      <header className="sheet-head">
        <Medallion heroId={c.heroId} art={c.art} />
        <div>
          <h2>{tx(c.name)}</h2>
          {cls && <p className="hero-class">{tx(cls.name)}</p>}
        </div>
      </header>
      {children}
      <Meter kind="hp" label={t('statHp')} value={c.hp} max={c.hpMax} />
      <SheetStats c={c} prefix="sheet" />
      <SheetTabs
        c={c}
        own
        dm={false}
        items={
          <ul className="hero-items" data-testid="sheet-items">
            {c.items.map((it, i) => (
              <li key={i}>
                {tx(it.name)}
                {it.qty > 1 && <span className="qty"> ×{it.qty}</span>}
              </li>
            ))}
          </ul>
        }
      />
    </section>
  )
}
