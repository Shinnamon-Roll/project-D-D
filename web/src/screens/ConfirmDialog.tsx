import { useEffect, useRef } from 'react'
import { useStore } from '../store'
import { useT } from '../i18n'
import { closeOnBackdrop } from './SeatQr'

// The one confirm popup, mounted once in App; opened with askConfirm().
// Native <dialog>: focus trap, Esc = cancel, focus returns to the button that opened it.
export default function ConfirmHost() {
  const ask = useStore((s) => s.confirm)
  return ask ? <ConfirmDialog key={ask.n} /> : null
}

function ConfirmDialog() {
  const t = useT()
  const ask = useStore((s) => s.confirm)!
  const ref = useRef<HTMLDialogElement>(null)
  const cancelRef = useRef<HTMLButtonElement>(null)
  const ok = useRef(false)
  useEffect(() => {
    ref.current?.showModal()
    cancelRef.current?.focus()
  }, [])

  // Every way out (OK, cancel, Esc, backdrop) goes through close(), so focus is restored first.
  const onClose = () => {
    useStore.setState({ confirm: null })
    if (ok.current) ask.onOk()
  }

  return (
    <dialog
      ref={ref}
      role="alertdialog"
      className="panel skill-modal confirm-modal"
      aria-modal="true"
      aria-labelledby="confirm-title"
      aria-describedby="confirm-body"
      data-testid="confirm-modal"
      onClose={onClose}
      onClick={closeOnBackdrop}
      onPointerDown={(e) => e.stopPropagation()}
    >
      <h2 id="confirm-title">{ask.title}</h2>
      <p id="confirm-body" data-testid="confirm-body">
        {ask.body}
      </p>
      <div className="row confirm-actions">
        <button ref={cancelRef} data-testid="confirm-cancel" onClick={() => ref.current?.close()}>
          {t('cancel')}
        </button>
        <button
          className={ask.danger ? 'danger' : 'primary'}
          data-testid="confirm-ok"
          onClick={() => {
            ok.current = true
            ref.current?.close()
          }}
        >
          {ask.okLabel}
        </button>
      </div>
    </dialog>
  )
}
