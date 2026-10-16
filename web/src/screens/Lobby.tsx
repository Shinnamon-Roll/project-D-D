import { useEffect, useRef, useState } from 'react'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import { useT } from '../i18n'
import SeatQr from './SeatQr'
import { NotesButton } from '../table/NotesPanel'
import { RosterButton } from '../table/RosterEditor'

const kindLabel = { dm: 'kindDm', codm: 'kindCodm', player: 'kindPlayer' } as const

export default function Lobby() {
  const room = useStore((s) => s.room)!
  const t = useT()
  const { seat: me, isDm } = useMe()
  const [name, setName] = useState('')
  const [qrFor, setQrFor] = useState('')
  const headingRef = useRef<HTMLHeadingElement>(null)
  const trimmed = name.trim()

  const players = room.seats.filter((s) => s.kind === 'player')
  const seatNum = (id: string) => players.findIndex((s) => s.id === id) + 1
  // ?seat=<id> from a scanned QR; App strips it once this browser is seated.
  const linked = me ? null : new URLSearchParams(location.search).get('seat')
  const target = linked ? players.find((s) => s.id === linked && !s.occupied) : undefined
  const occupiedPlayers = players.filter((s) => s.occupied)
  const canStart = occupiedPlayers.length > 0 && occupiedPlayers.every((s) => s.ready)
  const campaign = room.campaigns.find((c) => c.id === room.campaignId)
  // The QR stays open only while its seat exists and is empty.
  const qrOpen = isDm && players.some((s) => s.id === qrFor && !s.occupied)

  // Seat taken or removed under the open QR: the unmounted dialog can't return focus, so land on the heading.
  useEffect(() => {
    if (qrFor && !qrOpen) {
      setQrFor('')
      headingRef.current?.focus()
    }
  }, [qrFor, qrOpen])

  return (
    <div className="columns">
      <section className="panel">
        <h1 ref={headingRef} tabIndex={-1} data-testid="lobby-heading">
          {t('lobby')}
        </h1>
        <p className="subtitle">{campaign?.name}</p>

        {target ? (
          <form
            className="qr-join"
            data-testid="qr-join"
            onSubmit={(e) => {
              e.preventDefault()
              if (trimmed) send('seat.sit', { seatId: target.id, name: trimmed })
            }}
          >
            <strong>
              {t('joiningSeat', { n: seatNum(target.id) })}
            </strong>
            <label>
              {t('yourName')}
              <input autoFocus data-testid="sit-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
            </label>
            <button type="submit" className="primary" data-testid="qr-sit" disabled={!trimmed}>
              {t('sit')}
            </button>
          </form>
        ) : (
          !me && (
            <>
              {linked && (
                <p className="notice" role="status" data-testid="qr-seat-gone">
                  {t('seatGone')}
                </p>
              )}
              <label>
                {t('yourName')}
                <input data-testid="sit-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
              </label>
            </>
          )
        )}

        <ul className="seats">
          {room.seats.map((s) => (
            <li key={s.id} data-testid={`seat-${s.id}`} className={s.id === me?.id ? 'mine' : ''}>
              <span className="kind" data-testid={s.kind === 'player' ? `seat-num-${s.id}` : undefined}>
                {s.kind === 'player' ? t('seatN', { n: seatNum(s.id) }) : t(kindLabel[s.kind])}
              </span>
              <span className="who">{s.occupied ? s.name : <em>{t('empty')}</em>}</span>
              {s.kind === 'player' && s.occupied && (
                <span className={`seal ${s.ready ? 'ok' : ''}`} data-testid={`seat-ready-${s.id}`}>
                  {t(s.ready ? 'ready' : 'notReady')}
                </span>
              )}
              {!me && !s.occupied && s.kind !== 'dm' && (
                <button
                  data-testid={`seat-sit-${s.id}`}
                  disabled={!trimmed}
                  onClick={() => send('seat.sit', { seatId: s.id, name: trimmed })}
                >
                  {t('sit')}
                </button>
              )}
              {isDm && s.kind === 'player' && !s.occupied && (
                <button data-testid={`seat-qr-${s.id}`} onClick={() => setQrFor(s.id)}>
                  {t('seatQr')}
                </button>
              )}
              {isDm && s.kind !== 'dm' && (
                <button className="danger" data-testid={`seat-remove-${s.id}`} onClick={() => send('seat.remove', { seatId: s.id })}>
                  {t('removeSeat')}
                </button>
              )}
            </li>
          ))}
        </ul>

        <div className="row">
          {me?.kind === 'player' && (
            <button data-testid="ready-toggle" aria-pressed={me.ready} onClick={() => send('seat.ready', { ready: !me.ready })}>
              {t(me.ready ? 'notReady' : 'ready')}
            </button>
          )}
          {me && me.kind !== 'dm' && (
            <button data-testid="seat-leave" onClick={() => send('seat.leave', {})}>
              {t('leaveSeat')}
            </button>
          )}
        </div>

        {isDm && (
          <div className="row">
            <button data-testid="seat-add" disabled={players.length >= 6} onClick={() => send('seat.add', {})}>
              {t('addSeat')}
            </button>
            <button
              data-testid="seat-open-codm"
              disabled={room.seats.some((s) => s.kind === 'codm')}
              onClick={() => send('seat.openCoDm', {})}
            >
              {t('openCodm')}
            </button>
            <button className="primary" data-testid="game-start" disabled={!canStart} onClick={() => send('game.start', {})}>
              {t('start')}
            </button>
            <NotesButton />
            <RosterButton />
          </div>
        )}
      </section>
      {qrOpen && <SeatQr seatId={qrFor} num={seatNum(qrFor)} onClose={() => setQrFor('')} />}
    </div>
  )
}
