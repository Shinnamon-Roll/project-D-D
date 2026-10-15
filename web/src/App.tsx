import { useEffect } from 'react'
import { useMe, useStore } from './store'
import CampaignMenu from './screens/CampaignMenu'
import Lobby from './screens/Lobby'
import CharacterCreate, { CreatingDmView } from './screens/CharacterCreate'
import Table from './screens/Table'
import crest from './assets/crest.svg'

export default function App() {
  const room = useStore((s) => s.room)
  const connected = useStore((s) => s.connected)
  const lastError = useStore((s) => s.lastError)
  const { seat, character, isDm } = useMe()
  const linkedSeat = new URLSearchParams(location.search).get('seat')

  // A seated browser ignores ?seat= from a QR link; drop it so a later leave doesn't re-offer it.
  useEffect(() => {
    if (seat && linkedSeat) history.replaceState(null, '', location.pathname + location.hash)
  }, [seat, linkedSeat])

  let screen
  if (!room) screen = <p>กำลังเชื่อมต่อ… / Connecting…</p>
  else if (room.status === 'menu') screen = <CampaignMenu />
  else if (room.status === 'lobby') screen = <Lobby />
  else if (room.status === 'creating') {
    if (isDm) screen = <CreatingDmView />
    else if (seat && !character) screen = <CharacterCreate key={JSON.stringify(room.ruleset)} />
    else screen = <p data-testid="waiting">รอผู้เล่นคนอื่นสร้างตัวละคร… / Waiting for other players…</p>
  } else screen = <Table />

  return (
    <div className="app">
      <header className="topbar">
        <span className="brand">
          <img src={crest} alt="" />
          <strong>Tabletop</strong> <small>v{__APP_VERSION__}</small>
        </span>
        {seat && (
          <span className="me" data-testid="my-seat">
            {seat.name} · {seat.kind === 'player' ? 'ผู้เล่น / Player' : seat.kind === 'dm' ? 'DM' : 'Co-DM'}
          </span>
        )}
        <span className={`status ${connected ? 'ok' : 'bad'}`} data-testid="ws-status">
          <span className="flame" aria-hidden="true" />
          {connected ? 'connected' : 'disconnected'}
        </span>
      </header>
      {lastError && (
        <div className="error" role="alert" data-testid="error-banner">
          {lastError}
          <button aria-label="ปิด / Dismiss" onClick={() => useStore.setState({ lastError: null })}>
            ×
          </button>
        </div>
      )}
      {linkedSeat && !seat && room && room.status !== 'lobby' && (
        <p className="notice banner" role="status" data-testid="qr-not-lobby">
          ลิงก์ที่นั่งใช้ได้เฉพาะตอนห้องรอ / Seat links only work in the lobby
        </p>
      )}
      <main>{screen}</main>
      {room?.ruleset.attribution && (
        <footer className="attribution" data-testid="attribution">
          {room.ruleset.attribution}
        </footer>
      )}
    </div>
  )
}
