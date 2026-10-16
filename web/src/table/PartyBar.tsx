import { useEffect, useRef, useState, type MouseEvent } from 'react'
import type { CharStats, Character } from '../protocol'
import { useStore } from '../store'
import { send } from '../ws'
import { useL, useT, type TKey } from '../i18n'
import { CustomSkills, Medallion, SkillList } from './Skills'
import { Meter } from './SheetPanel'
import { SwordsIcon } from './icons'

const MAX_STACK = 999 // server's largest stack: item.remove with this qty always takes the whole stack
const STAT_FORM: [keyof CharStats, TKey, number, number][] = [
  ['hpMax', 'statHpMax', 1, 999],
  ['mpMax', 'statMpMax', 0, 99],
  ['str', 'statStr', 0, 20],
  ['end', 'statEnd', 0, 20],
  ['agi', 'statAgi', 0, 20],
]

// DM only (top of the map pane): one round avatar per hero with tiny HP/mana bars, plus Rest.
// An avatar toggles that hero's card, fixed bottom-right; ×, Esc or the same avatar closes it.
export default function PartyBar() {
  const t = useT()
  const tx = useL()
  const room = useStore((s) => s.room)!
  const [openId, setOpenId] = useState('')
  const bar = useRef<HTMLDivElement>(null)
  const open = room.characters.find((c) => c.id === openId)
  const player = (c: Character) => room.seats.find((s) => s.id === c.seatId)?.name ?? ''

  const close = () => {
    const id = openId
    setOpenId('')
    bar.current?.querySelector<HTMLElement>(`[data-testid="party-${CSS.escape(id)}"]`)?.focus()
  }

  // Esc anywhere closes the card, except inside a dialog or a form that uses Esc to close itself.
  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key !== 'Escape' || e.defaultPrevented) return
      if ((e.target as HTMLElement).closest?.('dialog, .skill-form, .stats-form, .npc-form')) return
      close()
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  })

  if (room.characters.length === 0) return null
  return (
    <>
      <div className="party-bar" data-testid="party-bar" ref={bar}>
        <ul className="party" aria-label={t('party')}>
          {room.characters.map((c) => {
            const label = t('partyOpen', { hero: tx(c.name), player: player(c), hp: c.hp, hpMax: c.hpMax, mp: c.mp, mpMax: c.mpMax })
            return (
              <li key={c.id}>
                <button
                  className={`party-avatar${c.hp === 0 ? ' down' : ''}`}
                  data-testid={`party-${c.id}`}
                  aria-expanded={c.id === openId}
                  aria-controls="char-card"
                  aria-label={label}
                  title={label}
                  onClick={() => (c.id === openId ? close() : setOpenId(c.id))}
                >
                  <Medallion classId={c.classId} art={c.art} />
                  <span className="mini hp" aria-hidden="true">
                    <span style={{ transform: `scaleX(${c.hpMax ? c.hp / c.hpMax : 0})` }} />
                  </span>
                  <span className="mini mp" aria-hidden="true">
                    <span style={{ transform: `scaleX(${c.mpMax ? c.mp / c.mpMax : 0})` }} />
                  </span>
                </button>
              </li>
            )
          })}
        </ul>
        <button
          className="combat-toggle"
          data-testid="combat-toggle"
          aria-pressed={room.combat}
          disabled={room.status !== 'playing'}
          onClick={() => send('combat.set', { on: !room.combat })}
        >
          <SwordsIcon />
          {t('combat')}
        </button>
        <button className="party-rest" data-testid="party-rest" disabled={room.status !== 'playing'} onClick={() => send('party.rest', {})}>
          {t('partyRest')}
        </button>
      </div>
      {open && <CharCard key={open.id} c={open} player={player(open)} onClose={close} />}
    </>
  )
}

function CharCard({ c, player, onClose }: { c: Character; player: string; onClose: () => void }) {
  const t = useT()
  const tx = useL()
  const entry = useStore((s) => s.room!.roster.find((r) => r.id === c.classId))
  const lastError = useStore((s) => s.lastError)
  const [editStats, setEditStats] = useState(false)
  // Bag indices come from this snapshot: after an item.remove/give, lock the bag until new items
  // (or an error) arrive, so a second click never targets a shifted index.
  const itemsKey = JSON.stringify(c.items)
  const [sent, setSent] = useState<{ items: string; err: string | null } | null>(null)
  const busy = sent !== null && sent.items === itemsKey && sent.err === lastError
  const itemSend = (sendIt: () => boolean) => {
    if (busy || !sendIt()) return false
    setSent({ items: itemsKey, err: lastError })
    return true
  }
  const editRef = useRef<HTMLButtonElement>(null)
  const hero = tx(c.name)
  const by = (e: MouseEvent) => (e.shiftKey ? 5 : 1) // shift-click = ±5
  const adjust = (kind: 'hp' | 'mp', d: number) =>
    send('character.adjust', kind === 'hp' ? { characterId: c.id, hp: d } : { characterId: c.id, mp: d })
  const doneStats = () => {
    setEditStats(false)
    requestAnimationFrame(() => editRef.current?.focus())
  }

  const pool = (kind: 'hp' | 'mp', value: number, max: number) => (
    <div className="pool" role="group" aria-label={t(kind === 'hp' ? 'hpOf' : 'mpOf', { name: hero })}>
      <button
        data-testid={`card-${kind}-dec`}
        aria-label={t(kind === 'hp' ? 'hpDec' : 'mpDec', { name: hero })}
        title={t('shiftFive')}
        disabled={value <= 0}
        onClick={(e) => adjust(kind, -by(e))}
      >
        −
      </button>
      <Meter kind={kind} label={t(kind === 'hp' ? 'statHp' : 'statMp')} value={value} max={max} testid={`card-${kind}`} />
      <button
        data-testid={`card-${kind}-inc`}
        aria-label={t(kind === 'hp' ? 'hpInc' : 'mpInc', { name: hero })}
        title={t('shiftFive')}
        disabled={value >= max}
        onClick={(e) => adjust(kind, by(e))}
      >
        +
      </button>
    </div>
  )

  return (
    <section id="char-card" className="panel char-card" data-testid="char-card" aria-labelledby="char-card-title">
      <header className="char-card-head">
        <Medallion classId={c.classId} art={c.art} size="lg" />
        <div className="char-card-who">
          <h2 id="char-card-title">{hero}</h2>
          {entry && <p className="hero-class">{tx(entry.className)}</p>}
          <p className="char-card-player">{t('playerOf', { name: player })}</p>
        </div>
        <button className="char-card-close" data-testid="char-card-close" aria-label={t('charCardClose')} onClick={onClose}>
          ×
        </button>
      </header>
      <div className="char-card-body">
        {pool('hp', c.hp, c.hpMax)}
        {pool('mp', c.mp, c.mpMax)}
        {editStats ? (
          <StatsForm c={c} onDone={doneStats} />
        ) : (
          <div className="card-stats">
            <dl className="stats-grid">
              {STAT_FORM.slice(2).map(([k, label]) => (
                <div key={k}>
                  <dt>{t(label)}</dt>
                  <dd data-testid={`card-${k}`}>{c[k]}</dd>
                </div>
              ))}
            </dl>
            <button ref={editRef} data-testid="card-stats-edit" onClick={() => setEditStats(true)}>
              {t('editStats')}
            </button>
          </div>
        )}

        <h3>{t('innateSkills')}</h3>
        <SkillList skills={c.skills.filter((s) => s.kind === 'innate')} />
        <h3>{t('playerSkills')}</h3>
        <CustomSkills c={c} canAdd={false} />

        <h3>{t('bag')}</h3>
        {c.items.length === 0 ? (
          <p className="muted">{t('bagEmpty')}</p>
        ) : (
          <ul className="bag" data-testid="card-items">
            {c.items.map((it, i) => {
              const name = tx(it.name)
              return (
                <li key={i}>
                  <span className="bag-name">{name}</span>
                  <span className="qty" data-testid={`item-qty-${i}`}>
                    ×{it.qty}
                  </span>
                  <button data-testid={`item-dec-${i}`} aria-label={t('itemDec', { name })} disabled={busy} onClick={() => itemSend(() => send('item.remove', { characterId: c.id, index: i, qty: 1 }))}>
                    −1
                  </button>
                  <button
                    className="danger"
                    data-testid={`item-del-${i}`}
                    aria-label={t('itemDel', { name })}
                    disabled={busy}
                    onClick={() => itemSend(() => send('item.remove', { characterId: c.id, index: i, qty: MAX_STACK }))}
                  >
                    {t('removeAll')}
                  </button>
                </li>
              )
            })}
          </ul>
        )}
        <GiveItem busy={busy} give={(name, qty) => itemSend(() => send('item.give', { characterId: c.id, name, qty }))} />
      </div>
    </section>
  )
}

// Limits mirror character.setStats on the server (UX only).
function StatsForm({ c, onDone }: { c: Character; onDone: () => void }) {
  const t = useT()
  const [v, setV] = useState<Record<keyof CharStats, string>>({
    hpMax: String(c.hpMax),
    mpMax: String(c.mpMax),
    str: String(c.str),
    end: String(c.end),
    agi: String(c.agi),
  })
  const num = (k: keyof CharStats) => Number(v[k])
  const bad = (k: keyof CharStats, lo: number, hi: number) => v[k] === '' || !Number.isInteger(num(k)) || num(k) < lo || num(k) > hi
  const ok = STAT_FORM.every(([k, , lo, hi]) => !bad(k, lo, hi))
  return (
    <form
      className="stats-form"
      onSubmit={(e) => {
        e.preventDefault()
        if (ok && send('character.setStats', { characterId: c.id, hpMax: num('hpMax'), mpMax: num('mpMax'), str: num('str'), end: num('end'), agi: num('agi') }))
          onDone()
      }}
      onKeyDown={(e) => e.key === 'Escape' && onDone()}
    >
      <div className="stats-inputs">
        {STAT_FORM.map(([k, label, lo, hi], i) => (
          <label key={k}>
            {t(label)}
            <input
              autoFocus={i === 0}
              data-testid={`card-stat-${k.toLowerCase()}`}
              type="number"
              min={lo}
              max={hi}
              step={1}
              value={v[k]}
              aria-invalid={bad(k, lo, hi)}
              onChange={(e) => setV({ ...v, [k]: e.target.value })}
            />
          </label>
        ))}
      </div>
      <div className="row">
        <span className="grow" />
        <button type="button" onClick={onDone}>
          {t('cancel')}
        </button>
        <button type="submit" className="primary" data-testid="card-stats-save" disabled={!ok}>
          {t('save')}
        </button>
      </div>
    </form>
  )
}

// item.give: name 1-40, qty 1-99; the same name adds to the existing stack (server side).
function GiveItem({ busy, give }: { busy: boolean; give: (name: string, qty: number) => boolean }) {
  const t = useT()
  const [name, setName] = useState('')
  const [qty, setQty] = useState('1')
  const nameRef = useRef<HTMLInputElement>(null)
  const n = Number(qty)
  const ok = name.trim() !== '' && Number.isInteger(n) && n >= 1 && n <= 99
  return (
    <form
      className="give-item"
      onSubmit={(e) => {
        e.preventDefault()
        if (!ok || !give(name.trim(), n)) return
        setName('')
        setQty('1')
        nameRef.current?.focus()
      }}
    >
      <h3>{t('giveItem')}</h3>
      <div className="row">
        <label className="grow">
          {t('itemName')}
          <input ref={nameRef} data-testid="give-item-name" value={name} maxLength={40} onChange={(e) => setName(e.target.value)} />
        </label>
        <label>
          {t('qty')}
          <input data-testid="give-item-qty" type="number" min={1} max={99} step={1} value={qty} onChange={(e) => setQty(e.target.value)} />
        </label>
        <button type="submit" className="primary" data-testid="give-item-submit" disabled={!ok || busy}>
          {t('give')}
        </button>
      </div>
    </form>
  )
}
