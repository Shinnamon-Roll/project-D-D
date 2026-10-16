import { useRef, useState } from 'react'
import type { CharSkill, Character, L } from '../protocol'
import { send } from '../ws'
import { useL, useT } from '../i18n'
import SkillInfo from '../screens/SkillInfo'
import { classArt } from '../classArt'

const MAX_CUSTOM = 3

// Round class medallion: gold icon on oxblood, or an uploaded picture filling the circle.
export function Medallion({ classId, art, size = 'md' }: { classId: string; art?: string; size?: 'sm' | 'md' | 'lg' }) {
  const a = classArt(classId, art)
  return (
    <span className={`medallion ${size} ${a.icon ? 'icon' : 'photo'}`} aria-hidden="true">
      {a.src && <img src={a.src} alt="" draggable={false} />}
    </span>
  )
}

export function ManaPill({ cost }: { cost: number }) {
  const t = useT()
  return <span className="mana-pill">{t('manaCost', { n: cost })}</span>
}

// One skill: name, mana, ⓘ; Use when play is set (disabled without enough mana); edit/remove for custom ones.
function SkillRow(props: { s: { id: string; name: L; desc: L; cost: number }; mp?: number; onEdit?: () => void; onRemove?: () => void }) {
  const { s, mp, onEdit, onRemove } = props
  const t = useT()
  const tx = useL()
  const [info, setInfo] = useState(false)
  const name = tx(s.name)
  const short = mp !== undefined && mp < s.cost
  return (
    <li className="skill-row">
      <span className="skill-name">{name}</span>
      <ManaPill cost={s.cost} />
      <button className="info-btn" data-testid={`skill-info-${s.id}`} aria-label={t('skillInfo', { name })} onClick={() => setInfo(true)}>
        ⓘ
      </button>
      {(mp !== undefined || onEdit) && (
        <span className="skill-actions">
          {onEdit && (
            <button data-testid={`custom-skill-edit-${s.id}`} aria-label={t('customSkillEdit', { name })} onClick={onEdit}>
              {t('edit')}
            </button>
          )}
          {onRemove && (
            <button className="danger" data-testid={`custom-skill-remove-${s.id}`} aria-label={t('removeNamed', { name })} onClick={onRemove}>
              {t('remove')}
            </button>
          )}
          {mp !== undefined && (
            <button
              className="primary use-btn"
              data-testid={`use-skill-${s.id}`}
              aria-label={t('useSkillNamed', { name, n: s.cost })}
              title={short ? t('notEnoughMana') : undefined}
              disabled={short}
              onClick={() => send('skill.use', { skillId: s.id })}
            >
              {t('useSkill')}
            </button>
          )}
        </span>
      )}
      {info && <SkillInfo skill={s} onClose={() => setInfo(false)} />}
    </li>
  )
}

// Roster/innate skills (read-only). mp set = show Use buttons.
export function SkillList({ skills, mp, testid }: { skills: { id: string; name: L; desc: L; cost: number }[]; mp?: number; testid?: string }) {
  return (
    <ul className="skill-list" data-testid={testid}>
      {skills.map((s) => (
        <SkillRow key={s.id} s={s} mp={mp} />
      ))}
    </ul>
  )
}

// The player's own skills (max 3): list + one add/edit form at a time. DM may edit too (server allows it).
export function CustomSkills({ c, play, canAdd = true }: { c: Character; play?: boolean; canAdd?: boolean }) {
  const t = useT()
  const custom = c.skills.filter((s) => s.kind === 'custom')
  const [editing, setEditing] = useState<string>('') // '' closed, 'new', or a skill id
  const addRef = useRef<HTMLButtonElement>(null)
  const close = () => {
    setEditing('')
    requestAnimationFrame(() => addRef.current?.focus())
  }
  const current = custom.find((s) => s.id === editing)
  return (
    <div className="custom-skills" data-testid="custom-skills">
      {custom.length === 0 && editing !== 'new' && <p className="muted">{t('noCustomSkills')}</p>}
      <ul className="skill-list">
        {custom.map((s) =>
          s.id === editing ? (
            <li key={s.id}>
              <SkillForm c={c} skill={s} onDone={close} />
            </li>
          ) : (
            <SkillRow
              key={s.id}
              s={s}
              mp={play ? c.mp : undefined}
              onEdit={() => setEditing(s.id)}
              onRemove={() => send('skill.removeCustom', { characterId: c.id, skillId: s.id })}
            />
          ),
        )}
      </ul>
      {editing === 'new' ? (
        <SkillForm c={c} onDone={close} />
      ) : (
        canAdd && (
        <button ref={addRef} data-testid="custom-skill-add" disabled={custom.length >= MAX_CUSTOM || current !== undefined} onClick={() => setEditing('new')}>
          {custom.length >= MAX_CUSTOM ? t('customSkillFull') : t('customSkillAdd')}
        </button>
        )
      )}
    </div>
  )
}

// Limits mirror the server (name 1-40, desc <= 300, cost 0-10); UX only, the server re-checks.
function SkillForm({ c, skill, onDone }: { c: Character; skill?: CharSkill; onDone: () => void }) {
  const t = useT()
  const [name, setName] = useState(skill?.name.th ?? '')
  const [desc, setDesc] = useState(skill?.desc.th ?? '')
  const [cost, setCost] = useState(String(skill?.cost ?? 1))
  const n = Number(cost)
  const ok = name.trim() !== '' && cost !== '' && Number.isInteger(n) && n >= 0 && n <= 10
  return (
    <form
      className="skill-form"
      onSubmit={(e) => {
        e.preventDefault()
        if (!ok) return
        const p = { name: name.trim(), desc: desc.trim(), cost: n }
        const sent = skill ? send('skill.editCustom', { characterId: c.id, skillId: skill.id, ...p }) : send('skill.addCustom', p)
        if (sent) onDone()
      }}
      onKeyDown={(e) => e.key === 'Escape' && onDone()}
    >
      <label>
        {t('customSkillName')}
        <input autoFocus data-testid="custom-skill-name" value={name} maxLength={40} required onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        {t('customSkillDesc')}
        <textarea data-testid="custom-skill-desc" value={desc} maxLength={300} rows={2} onChange={(e) => setDesc(e.target.value)} />
      </label>
      <div className="row">
        <label>
          {t('customSkillCost')}
          <input data-testid="custom-skill-cost" type="number" min={0} max={10} step={1} value={cost} onChange={(e) => setCost(e.target.value)} />
        </label>
        <span className="grow" />
        <button type="button" onClick={onDone}>
          {t('cancel')}
        </button>
        <button type="submit" className="primary" data-testid="custom-skill-save" disabled={!ok}>
          {t('customSkillSave')}
        </button>
      </div>
    </form>
  )
}
