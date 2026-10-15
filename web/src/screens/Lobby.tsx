import { useState } from 'react'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import RulesetEditor from './RulesetEditor'

const kindLabel = { dm: 'DM', codm: 'Co-DM', player: 'ผู้เล่น / Player' }

export default function Lobby() {
  const room = useStore((s) => s.room)!
  const { seat: me, isDm } = useMe()
  const [name, setName] = useState('')
  const trimmed = name.trim()

  const players = room.seats.filter((s) => s.kind === 'player')
  const occupiedPlayers = players.filter((s) => s.occupied)
  const canStart = occupiedPlayers.length > 0 && occupiedPlayers.every((s) => s.ready)
  const campaign = room.campaigns.find((c) => c.id === room.campaignId)

  return (
    <div className="columns">
      <section className="panel">
        <h1>ห้องรอ / Lobby</h1>
        <p className="subtitle">{campaign?.name}</p>

        {!me && (
          <label>
            ชื่อของคุณ / Your name
            <input data-testid="sit-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
          </label>
        )}

        <ul className="seats">
          {room.seats.map((s) => (
            <li key={s.id} data-testid={`seat-${s.id}`} className={s.id === me?.id ? 'mine' : ''}>
              <span className="kind">{kindLabel[s.kind]}</span>
              <span className="who">{s.occupied ? s.name : <em>ว่าง / empty</em>}</span>
              {s.kind === 'player' && s.occupied && (
                <span className={`seal ${s.ready ? 'ok' : ''}`} data-testid={`seat-ready-${s.id}`}>
                  {s.ready ? 'พร้อม / Ready' : 'ยังไม่พร้อม / Not ready'}
                </span>
              )}
              {!me && !s.occupied && s.kind !== 'dm' && (
                <button
                  data-testid={`seat-sit-${s.id}`}
                  disabled={!trimmed}
                  onClick={() => send('seat.sit', { seatId: s.id, name: trimmed })}
                >
                  นั่ง / Sit
                </button>
              )}
              {isDm && s.kind !== 'dm' && (
                <button className="danger" data-testid={`seat-remove-${s.id}`} onClick={() => send('seat.remove', { seatId: s.id })}>
                  ลบเก้าอี้ / Remove
                </button>
              )}
            </li>
          ))}
        </ul>

        <div className="row">
          {me?.kind === 'player' && (
            <button data-testid="ready-toggle" aria-pressed={me.ready} onClick={() => send('seat.ready', { ready: !me.ready })}>
              {me.ready ? 'ยังไม่พร้อม / Not ready' : 'พร้อม / Ready'}
            </button>
          )}
          {me && me.kind !== 'dm' && (
            <button data-testid="seat-leave" onClick={() => send('seat.leave', {})}>
              ลุก / Leave seat
            </button>
          )}
        </div>

        {isDm && (
          <div className="row">
            <button data-testid="seat-add" disabled={players.length >= 6} onClick={() => send('seat.add', {})}>
              + เก้าอี้ผู้เล่น / Add seat
            </button>
            <button
              data-testid="seat-open-codm"
              disabled={room.seats.some((s) => s.kind === 'codm')}
              onClick={() => send('seat.openCoDm', {})}
            >
              เปิดที่ DM ร่วม / Open co-DM seat
            </button>
            <button className="primary" data-testid="game-start" disabled={!canStart} onClick={() => send('game.start', {})}>
              เริ่มเกม / Start
            </button>
          </div>
        )}
      </section>
      {isDm && <RulesetEditor key={JSON.stringify(room.ruleset)} />}
    </div>
  )
}
