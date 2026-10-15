import { useEffect, useState } from 'react'

// Step 1 smoke screen: proves the browser ↔ Go WebSocket round trip.
export default function App() {
  const [status, setStatus] = useState('connecting…')

  useEffect(() => {
    const ws = new WebSocket(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`)
    ws.onopen = () => ws.send(JSON.stringify({ type: 'hello' }))
    ws.onmessage = (e) => setStatus(`connected · server says ${JSON.parse(e.data).type}`)
    ws.onclose = () => setStatus('disconnected')
    return () => ws.close()
  }, [])

  return (
    <main style={{ fontFamily: 'system-ui', padding: 24 }}>
      <h1>Tabletop v{__APP_VERSION__}</h1>
      <p data-testid="ws-status">{status}</p>
    </main>
  )
}
