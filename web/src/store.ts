import { create } from 'zustand'
import type { RoomState, Seat, You } from './protocol'
import type { Lang } from './i18n'

interface Store {
  room: RoomState | null
  you: You
  connected: boolean
  lastError: string | null
  lang: Lang
}

function savedLang(): Lang {
  try {
    return localStorage.getItem('lang') === 'en' ? 'en' : 'th'
  } catch {
    return 'th'
  }
}

export const useStore = create<Store>(() => ({
  room: null,
  you: { seatId: '' },
  connected: false,
  lastError: null,
  lang: savedLang(),
}))

let errorTimer: ReturnType<typeof setTimeout> | undefined
export function showError(message: string) {
  useStore.setState({ lastError: message })
  clearTimeout(errorTimer)
  errorTimer = setTimeout(() => useStore.setState({ lastError: null }), 6000)
}

export const isDmSeat = (s?: Seat) => s?.kind === 'dm' || s?.kind === 'codm'

// My seat + character, derived from the latest snapshot.
export function useMe() {
  const room = useStore((s) => s.room)
  const seatId = useStore((s) => s.you.seatId)
  const seat = room?.seats.find((s) => s.id === seatId)
  const character = room?.characters.find((c) => c.seatId === seatId && seatId !== '')
  return { seat, character, isDm: isDmSeat(seat) }
}
