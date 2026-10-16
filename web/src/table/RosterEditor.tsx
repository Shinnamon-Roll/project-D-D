import { useEffect, useRef, useState } from 'react'
import type { AbilityId, Class, L, RosterEntry } from '../protocol'
import { useStore } from '../store'
import { send, upload } from '../ws'
import { useL, useT, type TKey } from '../i18n'
import { classArt } from '../classArt'
import { closeOnBackdrop } from '../screens/SeatQr'
import { Medallion } from './Medallion'
import { abilityKey, skillKey } from './StatBlocks'
import { ABILITIES, useClass, useSkills } from '../rules'

// Limits mirror internal/rules/roster.go + classes.go (UX only, the server re-checks everything).
const NAME = 40
const DESC = 300
const MAX_ITEMS = 10
const SCORE = [1, 30] as const // abilities and AC
const ROLES = ['tank', 'damage', 'healer', 'support', 'control', 'hybrid']

const intIn = (n: number, lo: number, hi: number) => Number.isInteger(n) && n >= lo && n <= hi
// Mirrors rules.CleanText: 1..max runes after trim, no control chars (\n only when multiline).
const textOk = (s: string, max: number, multiline = false) => {
  const r = [...s.trim()]
  return r.length >= 1 && r.length <= max && !r.some((c) => /\p{Cc}/u.test(c) && !(multiline && c === '\n'))
}
const lOk = (l: L, max: number, multiline = false) => textOk(l.th, max, multiline) && textOk(l.en, max, multiline)
const blank = (): L => ({ th: '', en: '' })

// Skills and expertise need no check here: the checkboxes only ever make known, unique ids with expertise ⊆ skills.
function valid(e: RosterEntry) {
  return (
    lOk(e.heroName, NAME) &&
    lOk(e.blurb, DESC, true) &&
    ROLES.includes(e.role) &&
    ABILITIES.every((k) => intIn(e.abilities[k], ...SCORE)) &&
    intIn(e.ac, ...SCORE) &&
    e.items.length <= MAX_ITEMS &&
    e.items.every((it) => lOk(it.name, NAME) && intIn(it.qty, 1, 99))
  )
}

// art is "" (this class's own icon) or an upload, so another class's icon goes up as a PNG,
// inset like the medallion draws icons.
async function iconPng(id: string): Promise<File | null> {
  const img = new Image()
  img.src = classArt(id).src
  await img.decode()
  const c = document.createElement('canvas')
  c.width = c.height = 512
  c.getContext('2d')!.drawImage(img, 77, 77, 358, 358)
  const blob = await new Promise<Blob | null>((r) => c.toBlob(r, 'image/png'))
  return blob && new File([blob], `${id}.png`, { type: 'image/png' })
}

const numIn = (v: string) => (v === '' ? NaN : Number(v)) // empty stays invalid instead of turning into 0
const numOut = (n: number) => (Number.isNaN(n) ? '' : n)

// DM: the pre-made heroes, one tab per set (DM/co-DM get every set). Edits apply to heroes picked after saving.
export default function RosterEditor() {
  const t = useT()
  const tx = useL()
  const roster = useStore((s) => s.room!.roster)
  const sets = useStore((s) => s.room!.sets)
  const classes = useStore((s) => s.classes?.classes) ?? []
  const [tab, setTab] = useState(sets[0]?.id ?? '')
  const [openId, setOpenId] = useState('')
  const backFrom = useRef('')
  const list = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (openId || !backFrom.current) return
    list.current?.querySelector<HTMLElement>(`[data-testid="roster-item-${CSS.escape(backFrom.current)}"]`)?.focus()
    backFrom.current = ''
  }, [openId])

  const entry = roster.find((r) => r.id === openId)
  if (entry)
    return (
      <EntryEditor
        key={entry.id}
        entry={entry}
        onBack={() => {
          backFrom.current = entry.id
          setOpenId('')
        }}
      />
    )

  return (
    <div className="roster" data-testid="roster-editor" ref={list}>
      <p className="muted">{t('rosterLead')}</p>
      <div className="set-tabs" role="group" aria-label={t('sets')}>
        {sets.map((x) => (
          <button key={x.id} data-testid={`roster-tab-${x.id}`} aria-pressed={x.id === tab} onClick={() => setTab(x.id)}>
            {tx(x.name)}
          </button>
        ))}
      </div>
      <ul className="roster-list">
        {roster.filter((r) => r.set === tab).map((r) => (
          <li key={r.id}>
            <button data-testid={`roster-item-${r.id}`} onClick={() => setOpenId(r.id)}>
              <Medallion heroId={r.id} art={r.art} size="sm" />
              <span className="roster-name">{tx(r.heroName)}</span>
              <span className="roster-class">{tx(classes.find((c) => c.id === r.classId)?.name)}</span>
            </button>
          </li>
        ))}
      </ul>
      <button
        className="danger"
        data-testid="roster-reset-set"
        onClick={() => confirm(t('rosterResetSetConfirm', { name: tx(sets.find((x) => x.id === tab)?.name) })) && send('roster.reset', { set: tab })}
      >
        {t('rosterResetSet')}
      </button>
    </div>
  )
}

function LField(props: { label: string; value: L; max: number; multiline?: boolean; testid: string; onChange: (l: L) => void }) {
  const { label, value, max, multiline, testid, onChange } = props
  return (
    <fieldset className="l-pair">
      <legend>{label}</legend>
      {(['th', 'en'] as const).map((k) => {
        const common = {
          'data-testid': `${testid}-${k}`,
          lang: k,
          value: value[k],
          maxLength: max,
          'aria-invalid': !textOk(value[k], max),
        }
        return (
          <label key={k}>
            {k.toUpperCase()}
            {multiline ? (
              <textarea rows={2} {...common} onChange={(e) => onChange({ ...value, [k]: e.target.value })} />
            ) : (
              <input {...common} onChange={(e) => onChange({ ...value, [k]: e.target.value })} />
            )}
          </label>
        )
      })}
    </fieldset>
  )
}

function EntryEditor({ entry, onBack }: { entry: RosterEntry; onBack: () => void }) {
  const t = useT()
  const tx = useL()
  const roster = useStore((s) => s.room!.roster)
  const cls = useClass(entry.classId)
  const skillDefs = useSkills()
  const [d, setD] = useState(entry)
  const [busy, setBusy] = useState(false)
  const json = JSON.stringify(entry)
  useEffect(() => setD(entry), [json]) // saved, reset or changed by the co-DM: show the server's copy
  const dirty = JSON.stringify(d) !== json
  const ok = valid(d)

  const set = (p: Partial<RosterEntry>) => setD((x) => ({ ...x, ...p }))
  const setAbility = (k: AbilityId, n: number) => setD((x) => ({ ...x, abilities: { ...x.abilities, [k]: n } }))
  // Lists stay in sheet order (re-ticking a skill is not a change); unticking a skill drops its expertise too.
  const pick = (list: string[], id: string, on: boolean) => skillDefs.map((s) => s.id).filter((k) => (k === id ? on : list.includes(k)))
  const toggleSkill = (id: string, on: boolean) =>
    setD((x) => ({ ...x, skills: pick(x.skills, id, on), expertise: on ? x.expertise : x.expertise.filter((k) => k !== id) }))
  const toggleExpertise = (id: string, on: boolean) => setD((x) => ({ ...x, expertise: pick(x.expertise, id, on) }))
  const setItem = (i: number, p: Partial<RosterEntry['items'][number]>) =>
    setD((x) => ({ ...x, items: x.items.map((s, j) => (j === i ? { ...s, ...p } : s)) }))

  const applyArt = async (file: File | null) => {
    if (!file) return
    setBusy(true)
    const url = await upload(file)
    setBusy(false)
    if (url) set({ art: url })
  }
  const ownIcon = classArt(d.id).src
  // one button per distinct icon: "legend-bard" draws the same icon as "bard"
  const iconOpts = roster.filter((r, i) => roster.findIndex((x) => classArt(x.id).src === classArt(r.id).src) === i)
  const pickIcon = async (id: string) => {
    if (classArt(id).src === ownIcon) return set({ art: '' })
    setBusy(true)
    const f = await iconPng(id).catch(() => null)
    setBusy(false)
    await applyArt(f)
  }

  return (
    <form
      className="roster-form"
      data-testid="roster-form"
      onSubmit={(e) => {
        e.preventDefault()
        if (ok && dirty && !busy) send('roster.update', { classId: entry.id, entry: d })
      }}
    >
      <div className="row roster-tools">
        <button type="button" data-testid="roster-back" onClick={onBack}>
          {t('rosterBack')}
        </button>
        <span className="note-status" role="status" data-testid="roster-status">
          {dirty ? t(ok ? 'rosterUnsaved' : 'rosterInvalid') : t('noteSaved')}
        </span>
      </div>
      <header className="roster-head">
        <Medallion heroId={d.id} art={d.art} size="lg" />
        <div>
          <h3>{tx(d.heroName) || tx(entry.heroName)}</h3>
          <p className="roster-class">{tx(cls?.name)}</p>
        </div>
      </header>

      <LField label={t('heroName')} value={d.heroName} max={NAME} testid="roster-hero" onChange={(heroName) => set({ heroName })} />
      <LField label={t('blurb')} value={d.blurb} max={DESC} multiline testid="roster-blurb" onChange={(blurb) => set({ blurb })} />
      <label className="roster-role">
        {t('role')}
        <select data-testid="roster-role" value={d.role} onChange={(e) => set({ role: e.target.value })}>
          {ROLES.map((r) => (
            <option key={r} value={r}>
              {t(`role_${r}` as TKey)}
            </option>
          ))}
        </select>
      </label>

      <fieldset className="art-pick">
        <legend>{t('image')}</legend>
        <div className="art-grid">
          {iconOpts.map((r) => (
            <button
              type="button"
              key={r.id}
              className="art-opt"
              data-testid={`roster-art-${r.id}`}
              aria-pressed={d.art === '' && classArt(r.id).src === ownIcon}
              aria-label={t('artUse', { name: tx(r.heroName) })}
              title={tx(r.heroName)}
              disabled={busy}
              onClick={() => pickIcon(r.id)}
            >
              <img src={classArt(r.id).src} alt="" />
            </button>
          ))}
        </div>
        <label>
          {t('artUpload')}
          <input
            data-testid="roster-art-upload"
            type="file"
            accept="image/png,image/jpeg,image/webp,image/gif"
            disabled={busy}
            onChange={(e) => {
              const f = e.target.files?.[0] ?? null
              e.target.value = ''
              applyArt(f)
            }}
          />
        </label>
      </fieldset>

      <fieldset>
        <legend>{t('abilityScores')}</legend>
        <div className="stats-inputs">
          {ABILITIES.map((k) => (
            <label key={k}>
              {t(abilityKey(k))} ({k.toUpperCase()})
              <input
                data-testid={`roster-ab-${k}`}
                type="number"
                min={SCORE[0]}
                max={SCORE[1]}
                step={1}
                value={numOut(d.abilities[k])}
                aria-invalid={!intIn(d.abilities[k], ...SCORE)}
                onChange={(e) => setAbility(k, numIn(e.target.value))}
              />
            </label>
          ))}
          <label>
            {t('rosterAc')}
            <input
              data-testid="roster-ac"
              type="number"
              min={SCORE[0]}
              max={SCORE[1]}
              step={1}
              value={numOut(d.ac)}
              aria-invalid={!intIn(d.ac, ...SCORE)}
              onChange={(e) => set({ ac: numIn(e.target.value) })}
            />
          </label>
        </div>
      </fieldset>

      <fieldset>
        <legend>{t('skillsProf')}</legend>
        <p className="muted roster-hint">{t('rosterSkillsLead')}</p>
        <ul className="roster-skills">
          {skillDefs.map((s) => {
            const name = t(skillKey(s.id))
            const on = d.skills.includes(s.id)
            return (
              <li key={s.id} className={cls?.skillChoices.from.includes(s.id) ? 'class-choice' : undefined}>
                <label>
                  <input type="checkbox" data-testid={`roster-skill-${s.id}`} checked={on} onChange={(e) => toggleSkill(s.id, e.target.checked)} />
                  {name} <small>{s.ability.toUpperCase()}</small>
                </label>
                <label className="exp-toggle">
                  <input
                    type="checkbox"
                    data-testid={`roster-exp-${s.id}`}
                    aria-label={t('expertiseOf', { name })}
                    checked={d.expertise.includes(s.id)}
                    disabled={!on}
                    onChange={(e) => toggleExpertise(s.id, e.target.checked)}
                  />
                  <span aria-hidden="true">{t('expertise')}</span>
                </label>
              </li>
            )
          })}
        </ul>
      </fieldset>

      {cls && <ClassSummary cls={cls} />}

      <h3>{t('startItems')}</h3>
      <ol className="roster-rows">
        {d.items.map((it, i) => (
          <li key={i} className="roster-item-row" data-testid={`roster-gear-${i}`}>
            <LField label={t('itemN', { n: i + 1 })} value={it.name} max={NAME} testid={`roster-gear-${i}-name`} onChange={(name) => setItem(i, { name })} />
            <label>
              {t('qty')}
              <input
                data-testid={`roster-gear-${i}-qty`}
                type="number"
                min={1}
                max={99}
                step={1}
                value={numOut(it.qty)}
                aria-invalid={!intIn(it.qty, 1, 99)}
                onChange={(e) => setItem(i, { qty: numIn(e.target.value) })}
              />
            </label>
            <button
              type="button"
              className="danger"
              data-testid={`roster-gear-${i}-remove`}
              aria-label={t('removeNamed', { name: tx(it.name) || t('itemN', { n: i + 1 }) })}
              onClick={() => set({ items: d.items.filter((_, j) => j !== i) })}
            >
              {t('remove')}
            </button>
          </li>
        ))}
      </ol>
      <button
        type="button"
        data-testid="roster-gear-add"
        disabled={d.items.length >= MAX_ITEMS}
        onClick={() => set({ items: [...d.items, { name: blank(), qty: 1 }] })}
      >
        {t('addItem')}
      </button>

      <div className="row roster-actions">
        <button type="button" className="danger" data-testid="roster-reset" onClick={() => send('roster.reset', { classId: entry.id })}>
          {t('rosterReset')}
        </button>
        <span className="grow" />
        <button type="submit" className="primary" data-testid="roster-save" disabled={!ok || !dirty || busy}>
          {t('rosterSave')}
        </button>
      </div>
    </form>
  )
}

// Class data is fixed (SRD-based); shown so the DM sees what the hero's class already gives.
function ClassSummary({ cls }: { cls: Class }) {
  const t = useT()
  const tx = useL()
  const ab = (id: string) => t(abilityKey(id))
  const rows: [TKey, string][] = [
    ['hitDie', `d${cls.hitDie}`],
    ['classSaves', cls.saves.map(ab).join(', ')],
    ['classArmor', tx(cls.armor)],
    ['classWeapons', tx(cls.weapons)],
    ['classTools', tx(cls.tools)],
    ['classCaster', t(`caster_${cls.caster}` as TKey, { ab: cls.castAbility ? ab(cls.castAbility) : '' })],
  ]
  return (
    <section className="class-summary" data-testid="roster-class-info" aria-labelledby="class-info-title">
      <h3 id="class-info-title">
        {t('classInfo')}: {tx(cls.name)}
      </h3>
      <dl>
        {rows.map(([k, v]) => (
          <div key={k}>
            <dt>{t(k)}</dt>
            <dd>{v || '—'}</dd>
          </div>
        ))}
      </dl>
      <p className="muted">
        {t('classSkillChoices', { n: cls.skillChoices.count, list: cls.skillChoices.from.map((id) => t(skillKey(id))).join(', ') })}
      </p>
    </section>
  )
}

// Lobby / character pick: the same editor in a modal (like the story notes).
export function RosterButton() {
  const t = useT()
  const [open, setOpen] = useState(false)
  return (
    <>
      <button data-testid="roster-open" onClick={() => setOpen(true)}>
        {t('ruleset')}
      </button>
      {open && <RosterDialog onClose={() => setOpen(false)} />}
    </>
  )
}

function RosterDialog({ onClose }: { onClose: () => void }) {
  const t = useT()
  const ref = useRef<HTMLDialogElement>(null)
  useEffect(() => ref.current?.showModal(), [])
  return (
    <dialog ref={ref} className="panel notes-modal roster-modal" aria-labelledby="roster-title" data-testid="roster-modal" onClose={onClose} onClick={closeOnBackdrop}>
      <h2 id="roster-title">{t('ruleset')}</h2>
      <RosterEditor />
      <button className="primary" data-testid="roster-close" onClick={() => ref.current?.close()}>
        {t('close')}
      </button>
    </dialog>
  )
}
