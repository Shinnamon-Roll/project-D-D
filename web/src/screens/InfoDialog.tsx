import { useEffect, useRef } from 'react'
import type { L } from '../protocol'
import { useL, useT } from '../i18n'
import { closeOnBackdrop } from './SeatQr'

// Detail of a class feature (or any named rule text): name, an optional line under it, the description.
// Native <dialog>: focus trap, Esc, focus returns to the ⓘ button.
export default function InfoDialog({ name, desc, meta, onClose }: { name: L; desc: L; meta?: string; onClose: () => void }) {
  const t = useT()
  const tx = useL()
  const ref = useRef<HTMLDialogElement>(null)
  useEffect(() => ref.current?.showModal(), [])

  return (
    <dialog
      ref={ref}
      className="panel skill-modal"
      aria-modal="true"
      aria-labelledby="info-title"
      data-testid="info-modal"
      onClose={onClose}
      onClick={closeOnBackdrop}
      onPointerDown={(e) => e.stopPropagation()} // opened from a draggable carousel card: don't start a drag
    >
      <h2 id="info-title" data-testid="info-name">
        {tx(name)}
      </h2>
      {meta && <p className="skill-stat">{meta}</p>}
      <p data-testid="info-desc">{tx(desc)}</p>
      <button className="primary" data-testid="info-close" onClick={() => ref.current?.close()}>
        {t('close')}
      </button>
    </dialog>
  )
}
