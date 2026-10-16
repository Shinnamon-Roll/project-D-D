import { create } from 'zustand'
import type { Attack, BannerKind, L, LogEntry, Note, Npc, RoomState, Seat, You } from './protocol'
import type { Lang } from './i18n'

interface Store {
  room: RoomState | null
  you: You
  notes: Note[] // DM story notes; always empty for players
  npcs: Npc[] // DM NPC pool; always empty for players
  attacks: Attack[] // pending attacks; DM/co-DM only
  myAttack: { id: string; targetName: string } | null // a player's own pending attack
  dmLog: LogEntry[] // oldest first, at most LOG_CAP; DM/co-DM only
  banner: { kind: BannerKind; text: L; n: number } | null // big hit/miss/NPC-action banner; n remounts it for a repeat
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
  notes: [],
  npcs: [],
  attacks: [],
  myAttack: null,
  dmLog: [],
  banner: null,
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

export const LOG_CAP = 300

let bannerN = 0
let bannerTimer: ReturnType<typeof setTimeout> | undefined
export function showBanner(kind: BannerKind, text: L) {
  useStore.setState({ banner: { kind, text, n: ++bannerN } })
  clearTimeout(bannerTimer)
  bannerTimer = setTimeout(() => useStore.setState({ banner: null }), 2200) // CSS animation runs 2s
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
