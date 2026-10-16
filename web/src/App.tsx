import { useEffect } from 'react'
import { useMe, useStore } from './store'
import CampaignMenu from './screens/CampaignMenu'
import Lobby from './screens/Lobby'
import CharacterPick from './screens/CharacterPick'
import Table from './screens/Table'
import crest from './assets/crest.svg'
import { setLang, useT, type Lang } from './i18n'
import { CombatBanner } from './table/Combat'
import ConfirmHost from './screens/ConfirmDialog'

export default function App() {
  const room = useStore((s) => s.room)
  const connected = useStore((s) => s.connected)
  const lastError = useStore((s) => s.lastError)
  const lang = useStore((s) => s.lang)
  const t = useT()
  const { seat } = useMe()
  const linkedSeat = new URLSearchParams(location.search).get('seat')

  // A seated browser ignores ?seat= from a QR link; drop it so a later leave doesn't re-offer it.
  useEffect(() => {
    if (seat && linkedSeat) history.replaceState(null, '', location.pathname + location.hash)
  }, [seat, linkedSeat])
  useEffect(() => {
    document.documentElement.lang = lang
  }, [lang])

  let screen
  if (!room) screen = <p>{t('connecting')}</p>
  else if (room.status === 'menu') screen = <CampaignMenu />
  else if (room.status === 'lobby') screen = <Lobby />
  else if (room.status === 'creating') screen = <CharacterPick />
  else screen = <Table />

  return (
    <div className="app">
      <header className="topbar">
        <span className="brand">
          <img src={crest} alt="" />
          <strong>Tabletop</strong> <small>v{__APP_VERSION__}</small>
        </span>
        {seat && (
          <span className="me" data-testid="my-seat">
            {seat.name} · {t(seat.kind === 'player' ? 'kindPlayer' : seat.kind === 'dm' ? 'kindDm' : 'kindCodm')}
          </span>
        )}
        <span className="lang" role="group" aria-label={t('language')}>
          {(['th', 'en'] as Lang[]).map((l) => (
            <button key={l} lang={l} aria-label={l === 'th' ? 'ไทย' : 'English'} data-testid={`lang-${l}`} aria-pressed={lang === l} onClick={() => setLang(l)}>
              {l.toUpperCase()}
            </button>
          ))}
        </span>
        <span className={`status ${connected ? 'ok' : 'bad'}`} data-testid="ws-status">
          <span className="flame" aria-hidden="true" />
          {t(connected ? 'connected' : 'disconnected')}
        </span>
      </header>
      {lastError && (
        <div className="error" role="alert" data-testid="error-banner">
          {lastError}
          <button aria-label={t('dismiss')} onClick={() => useStore.setState({ lastError: null })}>
            ×
          </button>
        </div>
      )}
      {linkedSeat && !seat && room && room.status !== 'lobby' && (
        <p className="notice banner" role="status" data-testid="qr-not-lobby">
          {t('qrNotLobby')}
        </p>
      )}
      <main>{screen}</main>
      <CombatBanner />
      <ConfirmHost />
    </div>
  )
}
