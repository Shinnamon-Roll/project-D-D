import type { ClientEvents, ServerMessage } from './protocol'
import { showError, useStore } from './store'
import { translate } from './i18n'

// crypto.randomUUID only exists in secure contexts (https/localhost); plain-http LAN needs the fallback.
function newId(): string {
  if (crypto.randomUUID) return crypto.randomUUID()
  return Array.from(crypto.getRandomValues(new Uint8Array(16)), (b) => b.toString(16).padStart(2, '0')).join('')
}

function loadClientId(): string {
  try {
    let id = localStorage.getItem('clientId')
    if (!id) {
      id = newId()
      localStorage.setItem('clientId', id)
    }
    return id
  } catch {
    return newId() // storage blocked: works, but no rejoin after reload
  }
}

export const clientId = loadClientId()

let socket: WebSocket | null = null
let delay = 500
let firstBuild = '' // server build seen first in this page load

// Server restarted with a new build (TUI reload) → reload once to fetch the new web bundle.
function checkBuild(build?: string) {
  if (!build) return // old server without the field
  if (!firstBuild) {
    firstBuild = build
    return
  }
  if (build === firstBuild) return
  try {
    if (sessionStorage.getItem('reloadedFor') === build) return // already reloaded for this build: no loop
    sessionStorage.setItem('reloadedFor', build)
  } catch {
    // storage blocked: reload anyway; the new page's first build is this one, so it settles
  }
  location.reload()
}

export function connect() {
  const ws = new WebSocket(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`)
  socket = ws
  ws.onopen = () => {
    delay = 500
    ws.send(JSON.stringify({ type: 'hello', payload: { clientId } }))
    useStore.setState({ connected: true })
  }
  ws.onmessage = (e) => {
    let msg: ServerMessage
    try {
      msg = JSON.parse(e.data)
    } catch {
      return
    }
    if (msg.type === 'state') {
      checkBuild(msg.payload.build)
      useStore.setState({ room: msg.payload.room, you: msg.payload.you })
    }
    else if (msg.type === 'error') showError(msg.payload.message)
  }
  ws.onclose = () => {
    useStore.setState({ connected: false })
    setTimeout(connect, delay)
    delay = Math.min(delay * 2, 10000)
  }
}

export function send<T extends keyof ClientEvents>(type: T, payload: ClientEvents[T]) {
  if (socket?.readyState !== WebSocket.OPEN) {
    showError(translate(useStore.getState().lang, 'notConnected'))
    return
  }
  socket.send(JSON.stringify({ type, payload }))
}

// POST /api/upload → "/uploads/<name>". Errors surface in the banner and return null.
export async function upload(file: File): Promise<string | null> {
  const body = new FormData()
  body.append('file', file)
  try {
    const r = await fetch('/api/upload', { method: 'POST', headers: { 'X-Client-ID': clientId }, body })
    const j = await r.json().catch(() => ({}))
    if (!r.ok || typeof j.url !== 'string') throw new Error(j.error || `upload failed (${r.status})`)
    return j.url
  } catch (err) {
    showError(translate(useStore.getState().lang, 'uploadFailed', { error: (err as Error).message }))
    return null
  }
}
