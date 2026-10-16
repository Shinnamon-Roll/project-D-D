import { useState } from 'react'
import type { Abilities, AbilityId, Feature } from '../protocol'
import { useL, useT, type TKey } from '../i18n'
import { ABILITIES, abilityMod, saveBonus, signed, skillBonus, useSkills } from '../rules'
import InfoDialog from '../screens/InfoDialog'

export const abilityKey = (id: string) => `ab_${id}` as TKey
export const skillKey = (id: string) => `skill_${id}` as TKey

// Six ability cells: modifier (what you add to the die), score, and the saving throw; proficient saves are marked.
export function AbilityGrid({ abilities, saves, pb }: { abilities: Abilities; saves: AbilityId[]; pb: number }) {
  const t = useT()
  return (
    <dl className="ability-grid" aria-label={t('abilityScores')}>
      {ABILITIES.map((id) => {
        const label = t(abilityKey(id))
        const abbr = id.toUpperCase()
        const prof = saves.includes(id)
        const save = signed(saveBonus(abilities, id, pb, saves))
        return (
          <div key={id} className={`ab${prof ? ' prof' : ''}`} data-testid={`ab-${id}`}>
            <dt>
              {label}
              {label !== abbr && <small className="ab-abbr">{abbr}</small>}
            </dt>
            <dd className="ab-nums">
              <span className="ab-mod">{signed(abilityMod(abilities[id]))}</span>
              <span className="ab-score" data-testid={`ab-${id}-score`}>
                {abilities[id]}
              </span>
            </dd>
            <dd className="ab-save" data-testid={`ab-${id}-save`} title={prof ? t('saveProf') : undefined}>
              {t('saveShort', { n: save })}
              {prof && <span className="sr-only">, {t('saveProf')}</span>}
            </dd>
          </div>
        )
      })}
    </dl>
  )
}

// Proficient skills with the bonus to add, in sheet order; expertise doubles the proficiency bonus.
export function SkillChips({ abilities, skills, expertise, pb }: { abilities: Abilities; skills: string[]; expertise: string[]; pb: number }) {
  const t = useT()
  const defs = useSkills().filter((s) => skills.includes(s.id))
  return (
    <ul className="skill-chips" data-testid="skill-chips">
      {defs.map((s) => {
        const exp = expertise.includes(s.id)
        return (
          <li key={s.id} className={exp ? 'exp' : undefined} data-testid={`skill-${s.id}`} title={exp ? t('expertiseMark') : undefined}>
            {t(skillKey(s.id))} <b>{signed(skillBonus(abilities, s.ability, s.id, pb, skills, expertise))}</b>
            {exp && <span className="sr-only">, {t('expertiseMark')}</span>}
          </li>
        )
      })}
    </ul>
  )
}

// All 18 skills in sheet order with the bonus to add; ● proficient, ◆ expertise (double proficiency bonus).
export function SkillList({ abilities, skills, expertise, pb }: { abilities: Abilities; skills: string[]; expertise: string[]; pb: number }) {
  const t = useT()
  return (
    <ul className="skill-list" aria-label={t('skills')}>
      {useSkills().map((s) => {
        const mark = expertise.includes(s.id) ? 'exp' : skills.includes(s.id) ? 'prof' : ''
        return (
          <li key={s.id} className={mark || undefined} data-testid={`skill-${s.id}`}>
            <span className="sk-mark" aria-hidden="true">
              {mark === 'exp' ? '◆' : mark === 'prof' ? '●' : '○'}
            </span>
            <span className="sk-name">
              {t(skillKey(s.id))} <small>{s.ability.toUpperCase()}</small>
              {mark && <span className="sr-only">, {t(mark === 'exp' ? 'expertiseMark' : 'profMark')}</span>}
            </span>
            <b data-testid={`skill-${s.id}-bonus`}>{signed(skillBonus(abilities, s.ability, s.id, pb, skills, expertise))}</b>
          </li>
        )
      })}
    </ul>
  )
}

// Class features as chips; each chip opens the description.
export function FeatureList({ features }: { features: Feature[] }) {
  const t = useT()
  const tx = useL()
  const [open, setOpen] = useState<Feature | null>(null)
  return (
    <>
      <ul className="feature-chips">
        {features.map((f, i) => {
          const name = tx(f.name)
          return (
            <li key={i}>
              <button type="button" data-testid={`feature-info-${i}`} aria-label={t('infoOf', { name })} onClick={() => setOpen(f)}>
                {name} <span aria-hidden="true">ⓘ</span>
              </button>
            </li>
          )
        })}
      </ul>
      {open && <InfoDialog name={open.name} desc={open.desc} meta={t('featureLevel', { n: open.level })} onClose={() => setOpen(null)} />}
    </>
  )
}
