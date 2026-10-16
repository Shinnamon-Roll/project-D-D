import { useEffect, useRef, useState } from 'react'
import type { CharStats, L, RosterEntry } from '../protocol'
import { useStore } from '../store'
import { send, upload } from '../ws'
import { useL, useT, type TKey } from '../i18n'
import { classArt } from '../classArt'
import { closeOnBackdrop } from '../screens/SeatQr'
import { Medallion } from './Skills'

// Limits mirror internal/rules/roster.go (UX only, the server re-checks everything).
const NAME = 40
const DESC = 300
const MAX_SKILLS = 5
const MAX_ITEMS = 10
const STATS: [keyof CharStats, TKey, number, number][] = [
  ['hpMax', 'statHpMax', 1, 999],
  ['str', 'statStr', 0, 20],
  ['end', 'statEnd', 0, 20],
  ['agi', 'statAgi', 0, 20],
  ['mpMax', 'statMpMax', 0, 99],
]

const intIn = (n: number, lo: number, hi: number) => Number.isInteger(n) && n >= lo && n <= hi
// Mirrors rules.CleanText: 1..max runes after trim, no control chars (\n only when multiline).
const textOk = (s: string, max: number, multiline = false) => {
  const r = [...s.trim()]
  return r.length >= 1 && r.length <= max && !r.some((c) => /\p{Cc}/u.test(c) && !(multiline && c === '\n'))
}
const lOk = (l: L, max: number, multiline = false) => textOk(l.th, max, multiline) && textOk(l.en, max, multiline)
const blank = (): L => ({ th: '', en: '' })

function valid(e: RosterEntry) {
  return (
    lOk(e.heroName, NAME) &&
    lOk(e.className, NAME) &&
    lOk(e.blurb, DESC, true) &&
    STATS.every(([k, , lo, hi]) => intIn(e.stats[k], lo, hi)) &&
    e.skills.length >= 1 &&
    e.skills.length <= MAX_SKILLS &&
    e.skills.every((s) => lOk(s.name, NAME) && lOk(s.desc, DESC, true) && intIn(s.cost, 0, 10)) &&
    e.items.length <= MAX_ITEMS &&
    e.items.every((it) => lOk(it.name, NAME) && intIn(it.qty, 1, 99))
  )
}

// Skill ids: ^[a-z0-9-]{1,40}$, unique in the entry.
function newSkillId(skills: { id: string }[]) {
  let n = skills.length + 1
  while (skills.some((s) => s.id === `skill-${n}`)) n++
  return `skill-${n}`
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

// DM: the 10 pre-made heroes. Edits apply to heroes picked after saving.
export default function RosterEditor() {
  const t = useT()
  const tx = useL()
  const roster = useStore((s) => s.room!.roster)
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
      <ul className="roster-list">
        {roster.map((r) => (
          <li key={r.id}>
            <button data-testid={`roster-item-${r.id}`} onClick={() => setOpenId(r.id)}>
              <Medallion classId={r.id} art={r.art} size="sm" />
              <span className="roster-name">{tx(r.heroName)}</span>
              <span className="roster-class">{tx(r.className)}</span>
            </button>
          </li>
        ))}
      </ul>
      <button className="danger" data-testid="roster-reset-all" onClick={() => confirm(t('rosterResetAllConfirm')) && send('roster.reset', {})}>
        {t('rosterResetAll')}
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
  const [d, setD] = useState(entry)
  const [busy, setBusy] = useState(false)
  const json = JSON.stringify(entry)
  useEffect(() => setD(entry), [json]) // saved, reset or changed by the co-DM: show the server's copy
  const dirty = JSON.stringify(d) !== json
  const ok = valid(d)

  const set = (p: Partial<RosterEntry>) => setD((x) => ({ ...x, ...p }))
  const setStat = (k: keyof CharStats, n: number) => setD((x) => ({ ...x, stats: { ...x.stats, [k]: n } }))
  const setSkill = (i: number, p: Partial<RosterEntry['skills'][number]>) =>
    setD((x) => ({ ...x, skills: x.skills.map((s, j) => (j === i ? { ...s, ...p } : s)) }))
  const setItem = (i: number, p: Partial<RosterEntry['items'][number]>) =>
    setD((x) => ({ ...x, items: x.items.map((s, j) => (j === i ? { ...s, ...p } : s)) }))

  const applyArt = async (file: File | null) => {
    if (!file) return
    setBusy(true)
    const url = await upload(file)
    setBusy(false)
    if (url) set({ art: url })
  }
  const pickIcon = async (id: string) => {
    if (id === d.id) return set({ art: '' })
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
        <Medallion classId={d.id} art={d.art} size="lg" />
        <h3>{tx(d.heroName) || tx(entry.heroName)}</h3>
      </header>

      <LField label={t('heroName')} value={d.heroName} max={NAME} testid="roster-hero" onChange={(heroName) => set({ heroName })} />
      <LField label={t('className')} value={d.className} max={NAME} testid="roster-class" onChange={(className) => set({ className })} />
      <LField label={t('blurb')} value={d.blurb} max={DESC} multiline testid="roster-blurb" onChange={(blurb) => set({ blurb })} />

      <fieldset className="art-pick">
        <legend>{t('image')}</legend>
        <div className="art-grid">
          {roster.map((r) => (
            <button
              type="button"
              key={r.id}
              className="art-opt"
              data-testid={`roster-art-${r.id}`}
              aria-pressed={d.art === '' && r.id === d.id}
              aria-label={t('artUse', { name: tx(r.className) })}
              title={tx(r.className)}
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
        <legend>{t('stats')}</legend>
        <div className="stats-inputs">
          {STATS.map(([k, label, lo, hi]) => (
            <label key={k}>
              {t(label)}
              <input
                data-testid={`roster-${k.toLowerCase()}`}
                type="number"
                min={lo}
                max={hi}
                step={1}
                value={numOut(d.stats[k])}
                aria-invalid={!intIn(d.stats[k], lo, hi)}
                onChange={(e) => setStat(k, numIn(e.target.value))}
              />
            </label>
          ))}
        </div>
      </fieldset>

      <h3>{t('innateSkills')}</h3>
      <ol className="roster-rows">
        {d.skills.map((s, i) => (
          <li key={s.id} data-testid={`roster-skill-${i}`}>
            <LField label={t('skillN', { n: i + 1 })} value={s.name} max={NAME} testid={`roster-skill-${i}-name`} onChange={(name) => setSkill(i, { name })} />
            <LField label={t('skillDesc')} value={s.desc} max={DESC} multiline testid={`roster-skill-${i}-desc`} onChange={(desc) => setSkill(i, { desc })} />
            <div className="row">
              <label>
                {t('customSkillCost')}
                <input
                  data-testid={`roster-skill-${i}-cost`}
                  type="number"
                  min={0}
                  max={10}
                  step={1}
                  value={numOut(s.cost)}
                  aria-invalid={!intIn(s.cost, 0, 10)}
                  onChange={(e) => setSkill(i, { cost: numIn(e.target.value) })}
                />
              </label>
              <span className="grow" />
              <button
                type="button"
                className="danger"
                data-testid={`roster-skill-${i}-remove`}
                aria-label={t('removeNamed', { name: tx(s.name) || t('skillN', { n: i + 1 }) })}
                disabled={d.skills.length <= 1}
                onClick={() => set({ skills: d.skills.filter((_, j) => j !== i) })}
              >
                {t('remove')}
              </button>
            </div>
          </li>
        ))}
      </ol>
      <button
        type="button"
        data-testid="roster-skill-add"
        disabled={d.skills.length >= MAX_SKILLS}
        onClick={() => set({ skills: [...d.skills, { id: newSkillId(d.skills), name: blank(), desc: blank(), cost: 1 }] })}
      >
        {t('customSkillAdd')}
      </button>

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
