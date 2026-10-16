// Small stroke icons drawn for this project; decorative (the button carries the name).
const svg = { width: 16, height: 16, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', strokeWidth: 2, strokeLinecap: 'round', strokeLinejoin: 'round', 'aria-hidden': true } as const

export const EYE_PATH = 'M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12z'
export const EYE_OFF_SLASH = 'M3 3l18 18'

export function EyeIcon() {
  return (
    <svg {...svg}>
      <path d={EYE_PATH} />
      <circle cx="12" cy="12" r="3" />
    </svg>
  )
}

export function EyeOffIcon() {
  return (
    <svg {...svg}>
      <path d={EYE_PATH} />
      <circle cx="12" cy="12" r="3" />
      <path d={EYE_OFF_SLASH} />
    </svg>
  )
}

// Crossed swords: combat mode and attacks. The ⚔ glyph renders like a close "×" in the UI fonts.
export function SwordsIcon() {
  return (
    <svg {...svg} className="ico">
      <path d="M4 4l12 12M14 19l5-5M17 17l3 3M20 4L8 16M5 14l5 5M7 17l-3 3" />
    </svg>
  )
}
