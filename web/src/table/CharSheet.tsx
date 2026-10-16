import { useEffect, useRef, useState, type KeyboardEvent, type MouseEvent, type ReactNode } from 'react'
import type { Abilities, AbilityId, Character, Class, CustomAbility, Feature, Resource } from '../protocol'
import { send } from '../ws'
import { useL, useT, type TKey } from '../i18n'
import { ABILITIES, abilityMod, featuresAt, signed, skillBonus, useClass, useSkills } from '../rules'
import { AbilityGrid, SkillList, abilityKey, skillKey } from './StatBlocks'

// The rules v2 sheet, shared by the player's own sheet and the DM's hero card. Every number the server owns
// (maxima, PB) comes from the Character; the client only adds modifiers for display. The server checks every event.

const by = (e: MouseEvent) => (e.shiftKey ? 5 : 1) // shift-click = ×5
const CAB_NAME = 40 // mirrors the server's custom ability limits (UX only)
const CAB_DESC = 300
const CAB_MAX = 40

// Header numbers: level, AC, proficiency bonus, initiative (DEX modifier), passive Perception (10 + Perception).
// onLevel (DM) turns the level cell into a −/+ stepper.
export function SheetStats({ c, prefix, onLevel }: { c: Character; prefix: 'sheet' | 'card'; onLevel?: (d: 1 | -1) => void }) {
  const t = useT()
  const tx = useL()
  const per = skillBonus(c.abilities, 'wis', 'perception', c.pb, c.skills, c.expertise)
  const stats: [string, string, string, string][] = [
    ['ac', t('ac'), t('acHint'), String(c.ac)],
    ['pb', t('profBonus'), t('profBonusHint'), signed(c.pb)],
    ['init', t('initiative'), t('initiativeHint'), signed(abilityMod(c.abilities.dex))],
    ['passive', t('passive'), t('passiveHint'), String(10 + per)],
  ]
  const lvl = (
    <span className="lvl-num" data-testid={`${prefix}-level`}>
      {c.level}
    </span>
  )
  return (
    <dl className={`hero-vitals sheet-stats${onLevel ? ' with-step' : ''}`}>
      <div>
        <dt>{t('level')}</dt>
        <dd>
          {onLevel ? (
            <span className="lvl-step" role="group" aria-label={t('level')}>
              <button type="button" data-testid={`${prefix}-level-dec`} aria-label={t('levelDec', { name: tx(c.name) })} disabled={c.level <= 1} onClick={() => onLevel(-1)}>
                −
              </button>
              {lvl}
              <button type="button" data-testid={`${prefix}-level-inc`} aria-label={t('levelInc', { name: tx(c.name) })} disabled={c.level >= 20} onClick={() => onLevel(1)}>
                +
              </button>
            </span>
          ) : (
            lvl
          )}
        </dd>
      </div>
      {stats.map(([k, label, hint, v]) => (
        <div key={k} title={hint}>
          <dt>{label}</dt>
          <dd data-testid={`${prefix}-${k}`}>{v}</dd>
        </div>
      ))}
    </dl>
  )
}

const TABS = ['abilities', 'resources', 'features', 'mine', 'items'] as const
type Tab = (typeof TABS)[number]
const tabKey = (k: Tab) => `tab_${k}` as TKey

// own: the viewer's own character (custom abilities editable). dm: DM/co-DM controls (abilities/AC edit, remove).
// Inactive panels stay mounted but hidden, so the tab keeps its scroll and drafts.
export function SheetTabs({ c, own, dm, items }: { c: Character; own: boolean; dm: boolean; items: ReactNode }) {
  const t = useT()
  const cls = useClass(c.classId)
  const [tab, setTab] = useState<Tab>('abilities')
  const base = `${dm ? 'card' : 'sheet'}-${c.id}`
  const onKey = (e: KeyboardEvent) => {
    const i = TABS.indexOf(tab)
    const next = { ArrowLeft: i - 1, ArrowRight: i + 1, Home: 0, End: TABS.length - 1 }[e.key]
    if (next === undefined) return
    e.preventDefault()
    const k = TABS[(next + TABS.length) % TABS.length]
    setTab(k)
    document.getElementById(`${base}-tab-${k}`)?.focus()
  }
  const panel = (k: Tab) => {
    switch (k) {
      case 'abilities':
        return (
          <>
            {dm && (
              <div className="sheet-dm-tools">
                <AbilityEdit c={c} />
                <ProfEdit c={c} cls={cls} />
              </div>
            )}
            <AbilityGrid abilities={c.abilities} saves={c.saves} pb={c.pb} />
            <h3>{t('skills')}</h3>
            <SkillList abilities={c.abilities} skills={c.skills} expertise={c.expertise} pb={c.pb} />
          </>
        )
      case 'resources':
        return <Resources c={c} cls={cls} canEdit={own || dm} />
      case 'features':
        return cls ? <Features c={c} cls={cls} /> : <p className="muted">{t('classesLoading')}</p>
      case 'mine':
        return <MyAbilities c={c} own={own} dm={dm} />
      case 'items':
        return items
    }
  }
  return (
    <div className="sheet-tabs">
      <div className="tablist" role="tablist" aria-label={t('sheetTabs')} onKeyDown={onKey}>
        {TABS.map((k) => (
          <button
            key={k}
            type="button"
            role="tab"
            id={`${base}-tab-${k}`}
            aria-selected={tab === k}
            aria-controls={`${base}-panel-${k}`}
            tabIndex={tab === k ? 0 : -1}
            data-testid={`sheet-tab-${k}`}
            onClick={() => setTab(k)}
          >
            {t(tabKey(k))}
          </button>
        ))}
      </div>
      {TABS.map((k) => (
        <div key={k} className="tabpanel" role="tabpanel" id={`${base}-panel-${k}`} aria-labelledby={`${base}-tab-${k}`} hidden={tab !== k} data-testid={`sheet-panel-${k}`}>
          {panel(k)}
        </div>
      ))}
    </div>
  )
}

// DM: the six scores and AC in one form (character.setAbilities). Opening it copies the server's current values.
function AbilityEdit({ c }: { c: Character }) {
  const t = useT()
  const [open, setOpen] = useState(false)
  const [ab, setAb] = useState<Record<string, string>>({})
  const [ac, setAc] = useState('')
  const first = useRef<HTMLInputElement>(null)
  const toggle = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    if (open) first.current?.focus()
  }, [open])
  const ok = (v: string) => /^\d+$/.test(v) && Number(v) >= 1 && Number(v) <= 30
  const valid = ABILITIES.every((k) => ok(ab[k])) && ok(ac)
  const close = () => {
    setOpen(false)
    requestAnimationFrame(() => toggle.current?.focus())
  }
  if (!open)
    return (
      <button
        type="button"
        ref={toggle}
        className="ab-edit-open"
        data-testid="card-ab-edit"
        onClick={() => {
          setAb(Object.fromEntries(ABILITIES.map((k) => [k, String(c.abilities[k])])))
          setAc(String(c.ac))
          setOpen(true)
        }}
      >
        {t('abEdit')}
      </button>
    )
  return (
    <form
      className="ab-edit"
      onSubmit={(e) => {
        e.preventDefault()
        if (!valid) return
        const abilities = Object.fromEntries(ABILITIES.map((k) => [k, Number(ab[k])])) as Abilities
        if (send('character.setAbilities', { characterId: c.id, abilities, ac: Number(ac) })) close()
      }}
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          e.preventDefault() // keeps the card open
          close()
        }
      }}
    >
      <div className="ab-edit-grid">
        {ABILITIES.map((k, i) => (
          <label key={k}>
            {t(abilityKey(k))}
            <input
              ref={i === 0 ? first : undefined}
              type="number"
              min={1}
              max={30}
              step={1}
              data-testid={`card-edit-${k}`}
              aria-invalid={!ok(ab[k])}
              value={ab[k]}
              onChange={(e) => setAb({ ...ab, [k]: e.target.value })}
            />
          </label>
        ))}
        <label>
          {t('ac')}
          <input type="number" min={1} max={30} step={1} data-testid="card-edit-ac" aria-invalid={!ok(ac)} value={ac} onChange={(e) => setAc(e.target.value)} />
        </label>
      </div>
      <div className="row">
        <button type="submit" className="primary" data-testid="card-ab-save" disabled={!valid}>
          {t('abEditSave')}
        </button>
        <button type="button" data-testid="card-ab-cancel" onClick={close}>
          {t('cancel')}
        </button>
      </div>
    </form>
  )
}

// DM: skill/expertise and extra save proficiencies a class feature grants after the pick (character.setSkills).
// Class saves show ticked and locked; the server stores only the extras. Opening it copies the server's lists.
function ProfEdit({ c, cls }: { c: Character; cls?: Class }) {
  const t = useT()
  const skillDefs = useSkills()
  const [open, setOpen] = useState(false)
  const [skills, setSkills] = useState<string[]>([])
  const [exp, setExp] = useState<string[]>([])
  const [saves, setSaves] = useState<AbilityId[]>([])
  const first = useRef<HTMLInputElement>(null)
  const toggle = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    if (open) first.current?.focus()
  }, [open])
  const classSaves = cls?.saves ?? []
  // sheet order, like the server stores them
  const pick = (list: string[], id: string, on: boolean) => skillDefs.map((s) => s.id).filter((k) => (k === id ? on : list.includes(k)))
  const close = () => {
    setOpen(false)
    requestAnimationFrame(() => toggle.current?.focus())
  }
  if (!open)
    return (
      <button
        type="button"
        ref={toggle}
        className="ab-edit-open"
        data-testid="card-prof-edit"
        disabled={skillDefs.length === 0}
        onClick={() => {
          setSkills(c.skills)
          setExp(c.expertise)
          setSaves(c.extraSaves)
          setOpen(true)
        }}
      >
        {t('profEdit')}
      </button>
    )
  return (
    <form
      className="ab-edit prof-edit"
      onSubmit={(e) => {
        e.preventDefault()
        // a class save ticked as extra is dropped: the class already gives it
        const extra = ABILITIES.filter((k) => saves.includes(k) && !classSaves.includes(k))
        if (send('character.setSkills', { characterId: c.id, skills, expertise: exp, saves: extra })) close()
      }}
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          e.preventDefault() // keeps the card open
          close()
        }
      }}
    >
      <p className="muted roster-hint">{t('profEditLead')}</p>
      <fieldset>
        <legend>{t('skillsProf')}</legend>
        <ul className="roster-skills">
          {skillDefs.map((s, i) => {
            const name = t(skillKey(s.id))
            const on = skills.includes(s.id)
            return (
              <li key={s.id}>
                <label>
                  <input
                    ref={i === 0 ? first : undefined}
                    type="checkbox"
                    data-testid={`card-prof-skill-${s.id}`}
                    checked={on}
                    onChange={(e) => {
                      setSkills(pick(skills, s.id, e.target.checked))
                      if (!e.target.checked) setExp(exp.filter((k) => k !== s.id))
                    }}
                  />
                  {name} <small>{s.ability.toUpperCase()}</small>
                </label>
                <label className="exp-toggle">
                  <input
                    type="checkbox"
                    data-testid={`card-prof-exp-${s.id}`}
                    aria-label={t('expertiseOf', { name })}
                    checked={exp.includes(s.id)}
                    disabled={!on}
                    onChange={(e) => setExp(pick(exp, s.id, e.target.checked))}
                  />
                  <span aria-hidden="true">{t('expertise')}</span>
                </label>
              </li>
            )
          })}
        </ul>
      </fieldset>
      <fieldset>
        <legend>{t('classSaves')}</legend>
        <ul className="roster-skills">
          {ABILITIES.map((k) => {
            const locked = classSaves.includes(k)
            return (
              <li key={k}>
                <label>
                  <input
                    type="checkbox"
                    data-testid={`card-prof-save-${k}`}
                    aria-label={t('saveOf', { name: t(abilityKey(k)) })}
                    checked={locked || saves.includes(k)}
                    disabled={locked}
                    onChange={(e) => setSaves(e.target.checked ? [...saves, k] : saves.filter((x) => x !== k))}
                  />
                  {t(abilityKey(k))} <small>{k.toUpperCase()}</small>
                </label>
                {locked && <small className="exp-toggle">{t('classSaveLocked')}</small>}
              </li>
            )
          })}
        </ul>
      </fieldset>
      <div className="row">
        <button type="submit" className="primary" data-testid="card-prof-save">
          {t('profEditSave')}
        </button>
        <button type="button" data-testid="card-prof-cancel" onClick={close}>
          {t('cancel')}
        </button>
      </div>
    </form>
  )
}

// Spell slot pips: filled = spent. Tapping an open pip spends up to it; tapping a filled one gives back from it on.
function Pips({ id, label, max, used, disabled, onSet }: { id: string; label: string; max: number; used: number; disabled: boolean; onSet: (used: number) => void }) {
  const t = useT()
  return (
    <div className="slot-row" role="group" aria-label={label} data-testid={`slot-${id}`}>
      <span className="slot-label">{label}</span>
      <span className="pips">
        {Array.from({ length: max }, (_, j) => {
          const spent = j < used
          return (
            <button
              key={j}
              type="button"
              className={`pip${spent ? ' spent' : ''}`}
              aria-pressed={spent}
              aria-label={t('pipN', { n: j + 1 })}
              title={t('pipHint')}
              disabled={disabled}
              data-testid={`slot-${id}-${j}`}
              onClick={() => onSet(spent ? j : j + 1)}
            />
          )
        })}
      </span>
      <output className="left" data-testid={`slot-${id}-left`}>
        {max - used}/{max}
      </output>
    </div>
  )
}

// A use counter: − spends one, + gives one back (shift ×5). Shows what is left.
function Counter(props: { id: string; label: string; note?: string; max: number; used: number; disabled: boolean; onSet: (used: number) => void }) {
  const { id, label, note, max, used, disabled, onSet } = props
  const t = useT()
  return (
    <li className="counter" data-testid={id}>
      <span className="counter-name">
        {label}
        {note && <small>{note}</small>}
      </span>
      <span className="counter-ctl" role="group" aria-label={label}>
        <button type="button" data-testid={`${id}-use`} aria-label={t('useOne', { name: label })} title={t('shiftFive')} disabled={disabled || used >= max} onClick={(e) => onSet(Math.min(max, used + by(e)))}>
          −
        </button>
        <output className="left" data-testid={`${id}-left`}>
          {max - used}/{max}
        </output>
        <button type="button" data-testid={`${id}-restore`} aria-label={t('restoreOne', { name: label })} title={t('shiftFive')} disabled={disabled || used <= 0} onClick={(e) => onSet(Math.max(0, used - by(e)))}>
          +
        </button>
      </span>
    </li>
  )
}

// Which rest gives a resource back, in the server's terms (refill, shortRegain, shortFromLevel).
function refillNote(r: Resource, level: number): [TKey, Record<string, number>?] {
  if (r.refill === 'short' || (r.shortFromLevel > 0 && level >= r.shortFromLevel)) return ['refillShort']
  if (r.shortRegain > 0) return ['refillRegain', { n: r.shortRegain }]
  return ['refillLong']
}

function Resources({ c, cls, canEdit }: { c: Character; cls?: Class; canEdit: boolean }) {
  const t = useT()
  const tx = useL()
  const slots = c.slotsMax.map((max, i) => ({ level: i + 1, max, used: c.slotsUsed[i] ?? 0 })).filter((s) => s.max > 0)
  const res = (cls?.resources ?? []).filter((r) => (c.resourcesMax[r.id] ?? 0) > 0)
  const rows = (cls?.tables ?? []).map((tb) => [tx(tb.name), tb.values[c.level - 1]] as const).filter(([, v]) => v && v !== '—')
  return (
    <>
      {(slots.length > 0 || c.pact.count > 0) && (
        <>
          <h3>{t('spellSlots')}</h3>
          {slots.map((s) => (
            <Pips
              key={s.level}
              id={String(s.level)}
              label={t('slotLevel', { n: s.level })}
              max={s.max}
              used={s.used}
              disabled={!canEdit}
              onSet={(used) => send('slot.set', { characterId: c.id, level: s.level, used })}
            />
          ))}
          {c.pact.count > 0 && (
            <Pips
              id="pact"
              label={t('pactSlots', { n: c.pact.level })}
              max={c.pact.count}
              used={c.pactUsed}
              disabled={!canEdit}
              onSet={(used) => send('slot.set', { characterId: c.id, level: 'pact', used })}
            />
          )}
        </>
      )}
      <h3>{t('classRes')}</h3>
      <ul className="counters">
        {res.map((r) => {
          const [k, v] = refillNote(r, c.level)
          return (
            <Counter
              key={r.id}
              id={`res-${r.id}`}
              label={tx(r.name)}
              note={t(k, v)}
              max={c.resourcesMax[r.id]}
              used={c.resourcesUsed[r.id] ?? 0}
              disabled={!canEdit}
              onSet={(used) => send('resource.set', { characterId: c.id, resourceId: r.id, used })}
            />
          )
        })}
        <Counter
          id="hd"
          label={`${t('hitDice')}${cls ? ` d${cls.hitDie}` : ''}`}
          note={t('hitDiceNote')}
          max={c.hitDiceMax}
          used={c.hitDiceUsed}
          disabled={!canEdit}
          onSet={(used) => send('hitDice.set', { characterId: c.id, used })}
        />
      </ul>
      {rows.length > 0 && (
        <>
          <h3>{t('classTable')}</h3>
          <dl className="class-table" data-testid="class-table">
            {rows.map(([name, v]) => (
              <div key={name}>
                <dt>{name}</dt>
                <dd>{v}</dd>
              </div>
            ))}
          </dl>
        </>
      )}
    </>
  )
}

function FeatureRows({ list, testid }: { list: Feature[]; testid: string }) {
  const t = useT()
  const tx = useL()
  return list.map((f, i) => (
    <li key={i} data-testid={`${testid}-${i}`}>
      <details>
        <summary>
          {tx(f.name)} <small>{t('lvShort', { n: f.level })}</small>
        </summary>
        <p>{tx(f.desc)}</p>
      </details>
    </li>
  ))
}

// Class features up to the current level (tap to read), then the next level's, muted.
function Features({ c, cls }: { c: Character; cls: Class }) {
  const t = useT()
  const next = cls.features.filter((f) => f.level === c.level + 1)
  return (
    <>
      <h3>{t('featuresNow', { n: c.level })}</h3>
      <ul className="feat-list" data-testid="feat-now">
        <FeatureRows list={featuresAt(cls, c.level)} testid="feat" />
      </ul>
      {next.length > 0 && (
        <>
          <h3>{t('featuresNext', { n: c.level + 1 })}</h3>
          <ul className="feat-list next" data-testid="feat-next">
            <FeatureRows list={next} testid="feat-next" />
          </ul>
        </>
      )}
    </>
  )
}

const order = (a: CustomAbility) => (a.kind === 'feature' ? 10 : a.level) // spells by level, then features

// Spells and features the player typed in from the table's own book or cards (no spell library in the app).
function MyAbilities({ c, own, dm }: { c: Character; own: boolean; dm: boolean }) {
  const t = useT()
  const [editId, setEditId] = useState('')
  const list = [...c.abilitiesCustom].sort((a, b) => order(a) - order(b) || a.name.localeCompare(b.name))
  const meta = (a: CustomAbility) => (a.kind === 'feature' ? t('kind_feature') : a.level === 0 ? t('cantrip') : t('slotLevel', { n: a.level }))
  return (
    <>
      {list.length === 0 ? (
        <p className="muted">{t('mineEmpty')}</p>
      ) : (
        <ul className="feat-list cab-list" data-testid="cab-list">
          {list.map((a) =>
            a.id === editId ? (
              <li key={a.id}>
                <AbilityForm c={c} init={a} onDone={() => setEditId('')} />
              </li>
            ) : (
              <li key={a.id} data-testid={`cab-${a.id}`}>
                <details>
                  <summary>
                    {a.name} <small>{meta(a)}</small>
                  </summary>
                  {a.desc && <p>{a.desc}</p>}
                </details>
                {(own || dm) && (
                  <span className="cab-tools">
                    {own && (
                      <button type="button" data-testid={`cab-edit-${a.id}`} aria-label={t('cabEditOf', { name: a.name })} onClick={() => setEditId(a.id)}>
                        {t('edit')}
                      </button>
                    )}
                    <button
                      type="button"
                      className="danger"
                      data-testid={`cab-remove-${a.id}`}
                      aria-label={t('removeNamed', { name: a.name })}
                      onClick={() => send('ability.remove', { characterId: c.id, abilityId: a.id })}
                    >
                      {t('remove')}
                    </button>
                  </span>
                )}
              </li>
            ),
          )}
        </ul>
      )}
      {own && !editId && (c.abilitiesCustom.length < CAB_MAX ? <AbilityForm c={c} onDone={() => {}} /> : <p className="muted">{t('cabFull')}</p>)}
    </>
  )
}

// Add (no init) or edit one custom ability. A feature has no spell level (the server stores 0).
function AbilityForm({ c, init, onDone }: { c: Character; init?: CustomAbility; onDone: () => void }) {
  const t = useT()
  const [name, setName] = useState(init?.name ?? '')
  const [kind, setKind] = useState<CustomAbility['kind']>(init?.kind ?? 'spell')
  const [level, setLevel] = useState(init?.level ?? 1)
  const [desc, setDesc] = useState(init?.desc ?? '')
  const nameRef = useRef<HTMLInputElement>(null)
  const n = [...name.trim()].length
  const ok = n >= 1 && n <= CAB_NAME && [...desc].length <= CAB_DESC
  const base = `cab-${init?.id ?? 'new'}-${c.id}`
  return (
    <form
      className="cab-form"
      onSubmit={(e) => {
        e.preventDefault()
        if (!ok) return
        const body = { name: name.trim(), kind, level: kind === 'feature' ? 0 : level, desc: desc.trim() }
        if (!(init ? send('ability.update', { abilityId: init.id, ...body }) : send('ability.add', body))) return
        if (!init) {
          setName('')
          setDesc('')
          nameRef.current?.focus()
        }
        onDone()
      }}
      onKeyDown={(e) => {
        if (init && e.key === 'Escape') {
          e.preventDefault()
          onDone()
        }
      }}
    >
      {!init && <h3>{t('cabAdd')}</h3>}
      <div className="row">
        <label className="grow">
          {t('cabName')}
          <input ref={nameRef} autoFocus={!!init} data-testid="cab-name" maxLength={CAB_NAME} value={name} onChange={(e) => setName(e.target.value)} />
        </label>
        <label>
          {t('cabKind')}
          <select data-testid="cab-kind" value={kind} onChange={(e) => setKind(e.target.value as CustomAbility['kind'])}>
            <option value="spell">{t('kind_spell')}</option>
            <option value="feature">{t('kind_feature')}</option>
          </select>
        </label>
        <label>
          {t('cabLevel')}
          <select data-testid="cab-level" value={level} disabled={kind === 'feature'} onChange={(e) => setLevel(Number(e.target.value))}>
            {Array.from({ length: 10 }, (_, i) => (
              <option key={i} value={i}>
                {i === 0 ? t('cantrip') : t('slotLevel', { n: i })}
              </option>
            ))}
          </select>
        </label>
      </div>
      <label htmlFor={`${base}-desc`}>{t('cabDesc')}</label>
      <textarea id={`${base}-desc`} data-testid="cab-desc" rows={2} maxLength={CAB_DESC} value={desc} onChange={(e) => setDesc(e.target.value)} />
      <div className="row cab-foot">
        {init && (
          <button type="button" data-testid="cab-cancel" onClick={onDone}>
            {t('cancel')}
          </button>
        )}
        <button type="submit" className="primary" data-testid="cab-submit" disabled={!ok}>
          {init ? t('save') : t('cabSubmit')}
        </button>
      </div>
    </form>
  )
}
