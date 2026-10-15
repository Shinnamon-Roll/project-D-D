import { useEffect, useRef, useState } from 'react'
import type { Note } from '../protocol'
import { useStore } from '../store'
import { send } from '../ws'
import { useT } from '../i18n'
import { closeOnBackdrop } from '../screens/SeatQr'

const MAX_NOTES = 200

// DM/co-DM story notes. The server sends `notes` only to DM seats; players never get the key.
export default function NotesPanel() {
  const t = useT()
  const notes = useStore((s) => s.notes)
  const [openId, setOpenId] = useState('')
  // "add" pressed: ids that existed then + the title sent; open the first new note with that title.
  const awaitNew = useRef<{ ids: Set<string>; title: string } | null>(null)
  const backFrom = useRef('') // note just closed: focus its list button
  const list = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const w = awaitNew.current
    const n = w && notes.find((n) => !w.ids.has(n.id) && n.title === w.title)
    if (n) {
      awaitNew.current = null
      setOpenId(n.id)
    }
  }, [notes])

  useEffect(() => {
    if (openId || !backFrom.current) return
    list.current?.querySelector<HTMLElement>(`[data-testid="note-item-${CSS.escape(backFrom.current)}"]`)?.focus()
    backFrom.current = ''
  }, [openId])

  const i = notes.findIndex((n) => n.id === openId)
  if (i >= 0)
    return (
      <NoteEditor
        key={openId}
        note={notes[i]}
        index={i}
        count={notes.length}
        onBack={() => {
          backFrom.current = openId
          setOpenId('')
        }}
      />
    )

  return (
    <div className="notes" data-testid="notes" ref={list}>
      <button
        className="primary"
        data-testid="notes-add"
        disabled={notes.length >= MAX_NOTES}
        onClick={() => {
          const title = t('newNoteTitle')
          if (send('note.add', { title })) awaitNew.current = { ids: new Set(notes.map((n) => n.id)), title }
        }}
      >
        {t('notesAdd')}
      </button>
      {notes.length === 0 ? (
        <p className="muted">{t('noNotes')}</p>
      ) : (
        <ul className="note-list">
          {notes.map((n) => (
            <li key={n.id}>
              <button data-testid={`note-item-${n.id}`} onClick={() => setOpenId(n.id)}>
                {n.title}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

type Draft = { title: string; body: string }

function NoteEditor({ note, index, count, onBack }: { note: Note; index: number; count: number; onBack: () => void }) {
  const t = useT()
  const [title, setTitle] = useState(note.title)
  const [body, setBody] = useState(note.body)
  const [status, setStatus] = useState<'' | 'saving' | 'saved' | 'offline'>('')
  const connected = useStore((s) => s.connected)
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined)
  const pending = useRef<Draft | null>(null) // typed, not sent yet
  const sent = useRef<Draft | null>(null) // sent, server echo not seen yet
  const latest = useRef(note)

  // Server change: while the DM has unsaved typing keep their text; else the first change is our echo, later ones adopt.
  useEffect(() => {
    latest.current = note
  })
  useEffect(() => {
    if (pending.current) return
    if (sent.current) {
      sent.current = null
      setStatus('saved')
      return
    }
    setTitle(note.title)
    setBody(note.body)
  }, [note.title, note.body, note.updatedAt])

  const flush = () => {
    clearTimeout(timer.current)
    const p = pending.current
    if (!p) return
    pending.current = null
    // An empty title is not valid on the server: keep the stored one so the body still saves.
    const d = { title: p.title.trim() ? p.title : latest.current.title, body: p.body }
    if (d.title === latest.current.title && d.body === latest.current.body) {
      setStatus('saved') // no change: the server would not echo anything
      return
    }
    if (!send('note.update', { id: note.id, ...d })) {
      pending.current = p // offline: keep it, the effect below retries on reconnect
      setStatus('offline')
      return
    }
    sent.current = d
  }

  // Dropped before the echo: the server may not have it, so queue it again. Reconnected: send what is queued.
  // A resend of an already-saved draft is harmless: flush skips it if the state already matches, else the server echoes.
  useEffect(() => {
    if (!connected) {
      if (sent.current) pending.current ??= sent.current
      sent.current = null
      if (pending.current) setStatus('offline')
    } else if (pending.current) flush()
  }, [connected])

  // Leaving the editor or the page with typing in the 600ms window: save it now.
  useEffect(() => {
    window.addEventListener('pagehide', flush)
    return () => {
      window.removeEventListener('pagehide', flush)
      if (useStore.getState().notes.some((n) => n.id === note.id)) flush() // deleted (by anyone): nothing to save into
    }
  }, []) // flush only reads refs and the fixed note id

  const edit = (d: Draft) => {
    setTitle(d.title)
    setBody(d.body)
    pending.current = d
    setStatus('saving')
    clearTimeout(timer.current)
    timer.current = setTimeout(flush, 600)
  }

  const statusText = !title.trim() ? t('noteNeedTitle') : status && t(status === 'saving' ? 'noteSaving' : status === 'offline' ? 'noteOffline' : 'noteSaved')

  return (
    <div className="note-editor" data-testid="note-editor">
      <div className="row note-tools">
        <button data-testid="note-back" onClick={onBack}>
          {t('noteBack')}
        </button>
        <span className="note-status" role="status" data-testid="note-status">
          {statusText}
        </span>
        <button aria-label={t('noteUp')} title={t('noteUp')} data-testid="note-up" disabled={index === 0} onClick={() => send('note.move', { id: note.id, index: index - 1 })}>
          ↑
        </button>
        <button aria-label={t('noteDown')} title={t('noteDown')} data-testid="note-down" disabled={index >= count - 1} onClick={() => send('note.move', { id: note.id, index: index + 1 })}>
          ↓
        </button>
        <button
          className="danger"
          data-testid="note-delete"
          onClick={() => {
            if (!confirm(t('noteDeleteConfirm', { title: note.title }))) return
            clearTimeout(timer.current)
            pending.current = null // never save into a removed note
            send('note.remove', { id: note.id })
          }}
        >
          {t('noteDelete')}
        </button>
      </div>
      <label>
        {t('noteTitle')}
        <input autoFocus data-testid="note-title" maxLength={80} value={title} aria-invalid={!title.trim()} onChange={(e) => edit({ title: e.target.value, body })} />
      </label>
      <label className="note-body">
        {t('noteBody')}
        <textarea data-testid="note-body" maxLength={20000} value={body} onChange={(e) => edit({ title, body: e.target.value })} />
      </label>
    </div>
  )
}

// Lobby / character creation: the same notes in a modal.
export function NotesButton() {
  const t = useT()
  const [open, setOpen] = useState(false)
  return (
    <>
      <button data-testid="notes-open" onClick={() => setOpen(true)}>
        {t('storyNotes')}
      </button>
      {open && <NotesDialog onClose={() => setOpen(false)} />}
    </>
  )
}

function NotesDialog({ onClose }: { onClose: () => void }) {
  const t = useT()
  const ref = useRef<HTMLDialogElement>(null)
  useEffect(() => ref.current?.showModal(), [])
  return (
    <dialog ref={ref} className="panel notes-modal" aria-labelledby="notes-title" data-testid="notes-modal" onClose={onClose} onClick={closeOnBackdrop}>
      <h2 id="notes-title">{t('storyNotes')}</h2>
      <NotesPanel />
      <button className="primary" data-testid="notes-close" onClick={() => ref.current?.close()}>
        {t('close')}
      </button>
    </dialog>
  )
}
