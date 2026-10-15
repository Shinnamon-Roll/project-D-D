import { useState } from 'react'
import { useStore } from '../store'
import { send, upload } from '../ws'
import GrantPanel from '../table/GrantPanel'
import RulesetEditor from './RulesetEditor'
import SkillInfo from './SkillInfo'
import { useT } from '../i18n'

export default function CharacterCreate() {
  const t = useT()
  const ruleset = useStore((s) => s.room!.ruleset)
  const [name, setName] = useState('')
  const [infoId, setInfoId] = useState('')
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
  const infoSkill = ruleset.skills.find((s) => s.id === infoId)
  const canSubmit = name.trim() !== '' && left === 0 && skillIds.length === picks

  const bump = (key: string, d: number) => setStats({ ...stats, [key]: stats[key] + d })
  const toggleSkill = (id: string) =>
    setSkillIds(skillIds.includes(id) ? skillIds.filter((s) => s !== id) : [...skillIds, id])

  return (
    <section className="panel narrow">
      <h1>{t('createCharacter')}</h1>
      <label>
        {t('charName')}
        <input data-testid="char-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        {t('portrait')}
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
      {portraitUrl && <img className="portrait" src={portraitUrl} alt={t('portraitAlt')} data-testid="portrait-preview" />}

      <h2>
        {t('pointsLeft')} <span data-testid="points-left">{left}</span>
      </h2>
      <ul className="stats">
        {ruleset.stats.map((s) => (
          <li key={s.key}>
            <span>{s.name}</span>
            <button aria-label={t('decStat', { name: s.name })} data-testid={`stat-minus-${s.key}`} disabled={stats[s.key] <= s.min} onClick={() => bump(s.key, -1)}>
              −
            </button>
            <output data-testid={`stat-value-${s.key}`}>{stats[s.key]}</output>
            <button aria-label={t('incStat', { name: s.name })} data-testid={`stat-plus-${s.key}`} disabled={stats[s.key] >= s.max || left <= 0} onClick={() => bump(s.key, 1)}>
              +
            </button>
          </li>
        ))}
      </ul>

      <h2>{t('skillsCount', { n: skillIds.length, max: picks })}</h2>
      <div className="checks">
        {ruleset.skills.map((s) => (
          <div key={s.id} className="skill-row">
            <label className="check">
              <input
                type="checkbox"
                data-testid={`skill-${s.id}`}
                checked={skillIds.includes(s.id)}
                disabled={!skillIds.includes(s.id) && skillIds.length >= picks}
                onChange={() => toggleSkill(s.id)}
              />
              {s.name}
            </label>
            <button
              className="info-btn"
              data-testid={`skill-info-${s.id}`}
              aria-label={t('skillInfo', { name: s.name })}
              onClick={() => setInfoId(s.id)}
            >
              ⓘ
            </button>
          </div>
        ))}
      </div>
      {infoSkill && (
        <SkillInfo skill={infoSkill} stats={stats} picked={skillIds.includes(infoSkill.id)} onClose={() => setInfoId('')} />
      )}

      <h2>{t('startingItems')}</h2>
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
        {t('done')}
      </button>
    </section>
  )
}

// DM while players create: who is done, grant panel, ruleset editor.
export function CreatingDmView() {
  const t = useT()
  const room = useStore((s) => s.room)!
  return (
    <div className="columns">
      <section className="panel">
        <h1>{t('waitingCharacters')}</h1>
        <ul className="seats">
          {room.seats
            .filter((s) => s.kind === 'player')
            .map((s) => (
              <li key={s.id} data-testid={`creating-status-${s.id}`}>
                <span className="who">{s.name}</span>
                <span className={s.characterId ? 'ok' : 'muted'}>
                  {t(s.characterId ? 'charDone' : 'charCreating')}
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
