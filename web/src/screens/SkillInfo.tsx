import { useEffect, useRef } from 'react'
import type { SkillDef } from '../protocol'
import { useStore } from '../store'
import { useT } from '../i18n'
import { SRD_SKILLS } from '../content/srdSkills'
import { closeOnBackdrop } from './SeatQr'

const signed = (n: number) => (n < 0 ? `−${-n}` : `+${n}`)

// Skill card for a physical d20: what the skill is, and what to add to the roll with the current stats.
export default function SkillInfo(props: { skill: SkillDef; stats: Record<string, number>; picked: boolean; onClose: () => void }) {
  const { skill, stats, picked, onClose } = props
  const t = useT()
  const lang = useStore((s) => s.lang)
  const ruleset = useStore((s) => s.room!.ruleset)
  const ref = useRef<HTMLDialogElement>(null)
  useEffect(() => ref.current?.showModal(), [])

  const srd = SRD_SKILLS[skill.id]?.[lang]
  const own = skill.description?.trim()
  const desc = own || srd?.desc
  const examples = own ? [] : (srd?.examples ?? []) // DM text wins; SRD examples could contradict it
  const stat = ruleset.stats.find((s) => s.key === skill.statKey)
  const value = stat ? (stats[stat.key] ?? stat.base) : undefined
  const mod = value === undefined ? 0 : Math.floor((value - 10) / 2)

  return (
    <dialog
      ref={ref}
      className="panel skill-modal"
      aria-modal="true"
      aria-labelledby="skill-info-title"
      data-testid="skill-info-modal"
      onClose={onClose}
      onClick={closeOnBackdrop}
    >
      <h2 id="skill-info-title" data-testid="skill-info-name">
        {srd?.name ?? skill.name}
      </h2>
      {stat && <p className="skill-stat">{t('skillUses', { stat: stat.name || stat.key })}</p>}
      <p data-testid="skill-info-desc">{desc || t('skillAskDm')}</p>
      {examples.length > 0 && (
        <>
          <h3>{t('skillExamples')}</h3>
          <ul data-testid="skill-info-examples">
            {examples.map((x) => (
              <li key={x}>{x}</li>
            ))}
          </ul>
        </>
      )}
      {stat && value !== undefined ? (
        <div className="skill-how" data-testid="skill-info-how">
          <h3>{t('skillHowTo')}</h3>
          <p>{t('skillHowIntro')}</p>
          <p className="muted">
            {t('skillMod', { stat: stat.name || stat.key, value, mod: signed(mod) })} · {t('skillFormula')}
          </p>
          <ul>
            <li className={picked ? 'current' : ''}>{t('skillIfPicked', { total: signed(mod + 2) })}</li>
            <li className={picked ? '' : 'current'}>{t('skillIfNot', { total: signed(mod) })}</li>
          </ul>
        </div>
      ) : (
        desc && <p className="muted">{t('skillAskDm')}</p>
      )}
      <button className="primary" data-testid="skill-info-close" onClick={() => ref.current?.close()}>
        {t('close')}
      </button>
    </dialog>
  )
}
