import { useCallback, useEffect, useRef, useState } from 'react'
import type { Attack, L, Token } from '../protocol'
import { useStore } from '../store'
import { send } from '../ws'
import { useL, useT } from '../i18n'
import { SwordsIcon } from './icons'

export const hhmm = (at: number) => new Date(at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
// Death lines from the server carry a skull, revive lines end "is back up"; the feed and the DM log accent both in oxblood.
// shortcut: matches the server's wording, upgrade to a "kind" field on the entry if the strings change.
export const vitalClass = (text: L) => (text.en.includes('💀') ? 'death' : text.en.endsWith(' is back up') ? 'revive' : undefined)

// Player, combat on: Attack → pick a target (map tap or list) → confirm; while pending, a bar with Cancel.
// aim: null = not targeting, '' = targeting with nothing picked, else the picked token id.
export function AttackControl({ targets, aim, setAim }: { targets: Token[]; aim: string | null; setAim: (id: string | null) => void }) {
  const t = useT()
  const myAttack = useStore((s) => s.myAttack)
  // Focus follows the control as it swaps idle / aiming / pending: if the old one held focus, the new one's first control takes it.
  const hadFocus = useRef(false)
  // Each branch has its own key so React remounts it (idle and pending are both <div>s) and the ref sees the swap.
  // stable callback: a fresh one each render would detach/reattach and yank focus off the select
  const keepFocus = useCallback((el: HTMLElement | null) => {
    if (!el) return
    if (hadFocus.current) {
      hadFocus.current = false
      el.querySelector<HTMLElement>('select, button')?.focus()
    }
    return () => {
      if (el.contains(document.activeElement)) hadFocus.current = true
    }
  }, [])

  useEffect(() => () => setAim(null), [setAim]) // combat ended mid-aim: start fresh next time
  useEffect(() => {
    if (aim === null) return
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && setAim(null)
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [aim, setAim])

  if (myAttack)
    return (
      <div key="pending" className="attack-ctl pending" data-testid="attack-pending" role="status" ref={keepFocus}>
        <p>
          <strong>{t('attackPendingOn', { name: myAttack.targetName })}</strong>
          <br />
          {t('attackPending')}
        </p>
        <button data-testid="attack-cancel" onClick={() => send('attack.cancel', {})}>
          {t('attackCancel')}
        </button>
      </div>
    )
  if (aim === null)
    return (
      <div key="idle" className="attack-ctl" ref={keepFocus}>
        <button className="primary attack-btn" data-testid="attack-btn" onClick={() => setAim('')}>
          <SwordsIcon />
          {t('attack')}
        </button>
      </div>
    )
  return (
    <form
      className="attack-ctl aiming"
      data-testid="attack-aim"
      ref={keepFocus}
      onSubmit={(e) => {
        e.preventDefault()
        if (aim && send('attack.declare', { targetTokenId: aim })) setAim(null)
      }}
    >
      <p className="attack-hint">{t(targets.length ? 'attackPick' : 'attackNoTargets')}</p>
      <label>
        {t('attackTarget')}
        <select autoFocus data-testid="attack-target-select" value={aim} onChange={(e) => setAim(e.target.value)}>
          <option value="">{t('attackChoose')}</option>
          {targets.map((x) => (
            <option key={x.id} value={x.id}>
              {x.name}
            </option>
          ))}
        </select>
      </label>
      <div className="row">
        <button type="button" onClick={() => setAim(null)}>
          {t('cancel')}
        </button>
        <button type="submit" className="primary" data-testid="attack-confirm" disabled={!aim}>
          <SwordsIcon />
          {t('attackConfirm')}
        </button>
      </div>
    </form>
  )
}

// DM/co-DM: pending attacks stacked at the top of the screen, each with Hit / Miss and optional damage.
export function AttackAlerts() {
  const t = useT()
  const attacks = useStore((s) => s.attacks)
  const host = useRef<HTMLElement>(null)
  const resolving = useRef<string | null>(null) // alert the DM just answered; once it is gone, focus moves on
  useEffect(() => {
    const el = host.current
    if (!el || resolving.current === null || attacks.some((a) => a.id === resolving.current)) return
    resolving.current = null
    const active = document.activeElement
    if (active && active !== document.body && !el.contains(active)) return // DM already moved elsewhere
    ;(el.querySelector<HTMLElement>('input:not(:disabled)') ?? el).focus()
  }, [attacks])
  return (
    <section ref={host} tabIndex={-1} className="attack-alerts" aria-label={t('attackAlerts')} aria-live="assertive" data-testid="attack-alerts">
      {attacks.map((a) => (
        <AttackAlert key={a.id} a={a} onSent={() => (resolving.current = a.id)} />
      ))}
    </section>
  )
}

function AttackAlert({ a, onSent }: { a: Attack; onSent: () => void }) {
  const t = useT()
  const tx = useL()
  const [dmg, setDmg] = useState('')
  const [sent, setSent] = useState(false) // one answer per attack; the alert leaves when the server's state arrives
  const n = Number(dmg)
  const dmgOk = dmg === '' || (Number.isInteger(n) && n >= 0 && n <= 999)
  const resolve = (hit: boolean) => {
    if (!send('attack.resolve', hit && dmg !== '' && n > 0 ? { attackId: a.id, hit, damage: n } : { attackId: a.id, hit })) return
    setSent(true)
    onSent()
  }
  return (
    <article className="attack-alert" data-testid={`attack-alert-${a.id}`}>
      <p className="attack-alert-what">
        <SwordsIcon />
        {t('attackAlert', { by: tx(a.byName), target: a.targetName })}
        <time dateTime={new Date(a.at).toISOString()}>{hhmm(a.at)}</time>
      </p>
      <div className="attack-alert-act">
        <label className="attack-dmg">
          <span>
            {t('attackDmg')} <small>({t('attackDmgHint')})</small>
          </span>
          <input
            data-testid={`attack-dmg-${a.id}`}
            type="number"
            inputMode="numeric"
            min={0}
            max={999}
            step={1}
            value={dmg}
            disabled={sent}
            aria-invalid={!dmgOk}
            onChange={(e) => setDmg(e.target.value)}
          />
        </label>
        <button className="primary" data-testid={`attack-hit-${a.id}`} disabled={sent || !dmgOk} onClick={() => resolve(true)}>
          {t('attackHit')}
        </button>
        <button className="miss" data-testid={`attack-miss-${a.id}`} disabled={sent} onClick={() => resolve(false)}>
          {t('attackMiss')}
        </button>
      </div>
    </article>
  )
}

// Everyone: big Hit!/Miss! in the middle of the screen for ~2s. Never blocks clicks.
export function CombatBanner() {
  const t = useT()
  const tx = useL()
  const banner = useStore((s) => s.banner)
  return (
    <div className="combat-banner-host" aria-live="assertive">
      {banner && (
        <div key={banner.n} className={`combat-banner ${banner.kind}`} data-testid="combat-banner" data-kind={banner.kind}>
          <strong>{t(banner.kind === 'hit' ? 'bannerHit' : 'bannerMiss')}</strong>
          <span>{tx(banner.text)}</span>
        </div>
      )}
    </div>
  )
}
