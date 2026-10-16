import { useEffect, useRef } from 'react'
import type { L } from '../protocol'
import { useL, useT } from '../i18n'
import { closeOnBackdrop } from './SeatQr'

// Skill detail: name, mana cost, what it does. Native <dialog>: focus trap, Esc, focus returns to the ⓘ button.
export default function SkillInfo({ skill, onClose }: { skill: { name: L; desc: L; cost: number }; onClose: () => void }) {
  const t = useT()
  const tx = useL()
  const ref = useRef<HTMLDialogElement>(null)
  useEffect(() => ref.current?.showModal(), [])

  return (
    <dialog
      ref={ref}
      className="panel skill-modal"
      aria-modal="true"
      aria-labelledby="skill-info-title"
      data-testid="skill-info-modal"
      onClose={onClose}
      onClick={closeOnBackdrop}
      onPointerDown={(e) => e.stopPropagation()} // opened from a draggable carousel card: don't start a drag
    >
      <h2 id="skill-info-title" data-testid="skill-info-name">
        {tx(skill.name)}
      </h2>
      <p className="skill-stat">
        <span className="mana-pill">{t('manaCost', { n: skill.cost })}</span>
      </p>
      <p data-testid="skill-info-desc">{tx(skill.desc)}</p>
      <button className="primary" data-testid="skill-info-close" onClick={() => ref.current?.close()}>
        {t('close')}
      </button>
    </dialog>
  )
}
