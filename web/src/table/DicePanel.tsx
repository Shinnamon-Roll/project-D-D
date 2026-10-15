import { useState } from 'react'
import { useStore } from '../store'
import { send } from '../ws'

export default function DicePanel() {
  const dice = useStore((s) => s.room!.dice)
  const [notation, setNotation] = useState('')

  return (
    <section className="panel">
      <h2>ทอยเต๋า / Dice</h2>
      <form
        className="row"
        onSubmit={(e) => {
          e.preventDefault()
          if (notation.trim()) send('dice.roll', { notation: notation.trim() })
        }}
      >
        <label className="grow">
          <span className="sr-only">Notation</span>
          <input data-testid="dice-input" placeholder="1d20+3" value={notation} maxLength={64} onChange={(e) => setNotation(e.target.value)} />
        </label>
        <button className="primary" data-testid="dice-roll" type="submit">
          ทอย / Roll
        </button>
      </form>
      <ol className="dice-log" data-testid="dice-log" aria-live="polite">
        {[...dice].reverse().map((d) => (
          <li key={d.id} data-testid={`dice-entry-${d.id}`}>
            <strong>{d.by}</strong> {d.notation} = <b data-testid="dice-total">{d.total}</b>{' '}
            <small>[{d.rolls.join(', ')}]</small>
          </li>
        ))}
      </ol>
    </section>
  )
}
