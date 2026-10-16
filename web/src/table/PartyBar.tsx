import { useEffect, useRef, useState, type MouseEvent } from 'react'
import type { Character } from '../protocol'
import { askConfirm, showError, useStore } from '../store'
import { send } from '../ws'
import { useL, useT } from '../i18n'
import { useClass } from '../rules'
import { Medallion } from './Medallion'
import { Meter } from './SheetPanel'
import { SheetStats, SheetTabs } from './CharSheet'
import { SwordsIcon } from './icons'

const MAX_STACK = 999 // server's largest stack: item.remove with this qty always takes the whole stack

// DM only (top of the map pane): one round avatar per hero with a tiny HP bar, plus combat and short/long rest.
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
      if ((e.target as HTMLElement).closest?.('dialog, .npc-form')) return
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
            const label = t('partyOpen', { hero: tx(c.name), player: player(c), hp: c.hp, hpMax: c.hpMax })
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
                  <Medallion heroId={c.heroId} art={c.art} />
                  <span className="mini hp" aria-hidden="true">
                    <span style={{ transform: `scaleX(${c.hpMax ? c.hp / c.hpMax : 0})` }} />
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
          onClick={() =>
            room.combat
              ? askConfirm({ title: t('combatEnd'), body: t('combatEndConfirm'), okLabel: t('combatEnd') }, () => send('combat.set', { on: false }))
              : send('combat.set', { on: true })
          }
        >
          <SwordsIcon />
          {t('combat')}
        </button>
        {(['short', 'long'] as const).map((kind) => (
          <button key={kind} className="party-rest" data-testid={`party-rest-${kind}`} disabled={room.status !== 'playing'} onClick={() => {
              const label = t(kind === 'short' ? 'restShort' : 'restLong')
              askConfirm({ title: label, body: t(kind === 'short' ? 'restShortConfirm' : 'restLongConfirm'), okLabel: label }, () => send('party.rest', { kind }))
            }}
          >
            {t(kind === 'short' ? 'restShort' : 'restLong')}
          </button>
        ))}
      </div>
      {open && <CharCard key={open.id} c={open} player={player(open)} onClose={close} />}
    </>
  )
}

function CharCard({ c, player, onClose }: { c: Character; player: string; onClose: () => void }) {
  const t = useT()
  const tx = useL()
  const cls = useClass(c.classId)
  const lastError = useStore((s) => s.lastError)
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
  const hero = tx(c.name)
  const by = (e: MouseEvent) => (e.shiftKey ? 5 : 1) // shift-click = ±5
  const adjust = (d: number) => send('character.adjust', { characterId: c.id, hp: d })

  return (
    <section id="char-card" className="panel char-card" data-testid="char-card" aria-labelledby="char-card-title">
      <header className="char-card-head">
        <Medallion heroId={c.heroId} art={c.art} size="lg" />
        <div className="char-card-who">
          <h2 id="char-card-title">{hero}</h2>
          {cls && <p className="hero-class">{tx(cls.name)}</p>}
          <p className="char-card-player">{t('playerOf', { name: player })}</p>
        </div>
        <button className="char-card-close" data-testid="char-card-close" aria-label={t('charCardClose')} onClick={onClose}>
          ×
        </button>
      </header>
      <div className="char-card-body">
        <div className="pool" role="group" aria-label={t('hpOf', { name: hero })}>
          <button data-testid="card-hp-dec" aria-label={t('hpDec', { name: hero })} title={t('shiftFive')} disabled={c.hp <= 0} onClick={(e) => adjust(-by(e))}>
            −
          </button>
          <Meter kind="hp" label={t('statHp')} value={c.hp} max={c.hpMax} testid="card-hp" />
          <button data-testid="card-hp-inc" aria-label={t('hpInc', { name: hero })} title={t('shiftFive')} disabled={c.hp >= c.hpMax} onClick={(e) => adjust(by(e))}>
            +
          </button>
        </div>
        <SheetStats c={c} prefix="card" onLevel={(delta) => send('character.level', { characterId: c.id, delta })} />
        <SheetTabs
          c={c}
          own={false}
          dm
          items={
            <>
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
                          onClick={() =>
                            askConfirm({ title: t('itemDel', { name }), body: t('itemDelConfirm', { name, qty: it.qty, hero }), okLabel: t('removeAll'), danger: true }, () => {
                              // The index is from this snapshot: if the bag changed while the popup was open, don't delete a shifted item.
                              const now = useStore.getState().room?.characters.find((x) => x.id === c.id)
                              if (JSON.stringify(now?.items) !== itemsKey) return showError(t('itemsChanged'))
                              itemSend(() => send('item.remove', { characterId: c.id, index: i, qty: MAX_STACK }))
                            })
                          }
                        >
                          {t('removeAll')}
                        </button>
                      </li>
                    )
                  })}
                </ul>
              )}
              <GiveItem busy={busy} give={(name, qty) => itemSend(() => send('item.give', { characterId: c.id, name, qty }))} />
            </>
          }
        />
      </div>
    </section>
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
