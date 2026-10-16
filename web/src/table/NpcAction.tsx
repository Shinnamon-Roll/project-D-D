import { useEffect, useRef, useState } from 'react'
import type { ClientEvents, NpcActKind, Token } from '../protocol'
import { send } from '../ws'
import { useT, type TKey } from '../i18n'
import { closeOnBackdrop } from '../screens/SeatQr'
import { SwordsIcon } from './icons'

const KINDS: [NpcActKind, TKey][] = [
  ['talk', 'npcActTalk'],
  ['attack', 'attack'],
  ['gift', 'npcActGift'],
  ['surrender', 'npcActSurrender'],
]
const intIn = (s: string, lo: number, hi: number) => s.trim() !== '' && Number.isInteger(Number(s)) && Number(s) >= lo && Number(s) <= hi

// DM/co-DM, combat on: an NPC talks to, attacks, gives an item to or surrenders to a player character.
// Native <dialog>: first control gets focus (target picker or Talk), focus trap, Esc, focus returns to whatever opened it. The server re-checks everything.
// target '' = opened from the token list: pick the player in the dialog.
export default function NpcAction({ npc, targets, target: first, onClose }: { npc: Token; targets: Token[]; target: string; onClose: () => void }) {
  const t = useT()
  const ref = useRef<HTMLDialogElement>(null)
  useEffect(() => ref.current?.showModal(), [])
  const [target, setTarget] = useState(first)
  const [kind, setKind] = useState<NpcActKind | null>(null)
  const [text, setText] = useState('')
  const [hit, setHit] = useState<boolean | null>(null)
  const [dmg, setDmg] = useState('')
  const [item, setItem] = useState('')
  const [qty, setQty] = useState('1')
  const to = targets.find((x) => x.id === target)
  const dmgOk = dmg === '' || intIn(dmg, 0, 999)
  const ok =
    !!to &&
    (kind === 'talk' ||
      kind === 'surrender' ||
      (kind === 'attack' && hit !== null && (!hit || dmgOk)) ||
      (kind === 'gift' && item.trim().length >= 1 && intIn(qty, 1, 99)))

  const submit = () => {
    if (!ok || !kind) return
    const base = { tokenId: npc.id, targetId: target, kind }
    const extra: Omit<ClientEvents['npc.act'], 'tokenId' | 'targetId' | 'kind'> =
      kind === 'talk'
        ? { text: text.trim() }
        : kind === 'attack'
          ? hit && dmg !== '' ? { hit: true, damage: Number(dmg) } : { hit: !!hit }
          : kind === 'gift'
            ? { item: { name: item.trim(), qty: Number(qty) } }
            : {}
    // shortcut: closes once sent; a server refusal (e.g. a full bag) shows in the error banner and the DM reopens.
    if (send('npc.act', { ...base, ...extra })) ref.current?.close()
  }

  return (
    <dialog
      ref={ref}
      className="panel npc-act-modal"
      aria-modal="true"
      aria-labelledby="npc-act-title"
      data-testid="npc-act-dialog"
      onClose={onClose}
      onClick={closeOnBackdrop}
    >
      <form
        onSubmit={(e) => {
          e.preventDefault()
          submit()
        }}
      >
        <h2 id="npc-act-title">
          {npc.name} → {to?.name ?? '…'}
        </h2>
        {first === '' && (
          <label className="npc-act-field">
            {t('attackTarget')}
            <select data-testid="npc-act-target" value={target} onChange={(e) => setTarget(e.target.value)}>
              <option value="">{t(targets.length ? 'attackChoose' : 'npcActNoTargets')}</option>
              {targets.map((x) => (
                <option key={x.id} value={x.id}>
                  {x.name}
                </option>
              ))}
            </select>
          </label>
        )}
        <div className="npc-act-kinds" role="group" aria-label={t('npcActWhat')}>
          {KINDS.map(([k, label]) => (
            <button
              key={k}
              type="button"
              className={`npc-kind ${k}`}
              data-testid={`npc-act-${k}`}
              aria-pressed={kind === k}
              onClick={() => setKind(k)}
            >
              {k === 'attack' && <SwordsIcon />}
              {t(label)}
            </button>
          ))}
        </div>
        {kind === 'talk' && (
          <label className="npc-act-field">
            <span>
              {t('npcActText')} <small>({t('attackDmgHint')})</small>
            </span>
            <input data-testid="npc-act-text" value={text} maxLength={200} onChange={(e) => setText(e.target.value)} />
          </label>
        )}
        {kind === 'attack' && (
          <div className="npc-act-attack">
            <div className="row" role="group" aria-label={t('npcActHitOrMiss')}>
              <button type="button" className="primary" data-testid="npc-act-hit" aria-pressed={hit === true} onClick={() => setHit(true)}>
                {t('attackHit')}
              </button>
              <button type="button" className="miss" data-testid="npc-act-miss" aria-pressed={hit === false} onClick={() => setHit(false)}>
                {t('attackMiss')}
              </button>
            </div>
            <label className="npc-act-field">
              <span>
                {t('attackDmg')} <small>({t('npcActDmgHint')})</small>
              </span>
              <input
                data-testid="npc-act-damage"
                type="number"
                inputMode="numeric"
                min={0}
                max={999}
                step={1}
                value={dmg}
                disabled={hit === false}
                aria-invalid={!dmgOk}
                onChange={(e) => setDmg(e.target.value)}
              />
            </label>
          </div>
        )}
        {kind === 'gift' && (
          <div className="row">
            <label className="npc-act-field grow">
              {t('itemName')}
              <input data-testid="npc-act-item-name" value={item} maxLength={40} onChange={(e) => setItem(e.target.value)} />
            </label>
            <label className="npc-act-field">
              {t('qty')}
              <input data-testid="npc-act-item-qty" type="number" inputMode="numeric" min={1} max={99} step={1} value={qty} aria-invalid={!intIn(qty, 1, 99)} onChange={(e) => setQty(e.target.value)} />
            </label>
          </div>
        )}
        {kind === 'surrender' && <p className="npc-act-note">{t('npcActSurrenderNote', { npc: npc.name, target: to?.name ?? '…' })}</p>}
        <div className="row npc-act-foot">
          <button type="button" data-testid="npc-act-cancel" onClick={() => ref.current?.close()}>
            {t('cancel')}
          </button>
          <button type="submit" className="primary" data-testid="npc-act-send" disabled={!ok}>
            {t('npcActSend')}
          </button>
        </div>
      </form>
    </dialog>
  )
}
