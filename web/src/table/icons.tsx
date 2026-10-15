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
