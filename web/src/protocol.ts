// Hand-written mirror of docs/design/protocol.md (v0). JSON field names must match exactly.

export type RoomStatus = 'menu' | 'lobby' | 'creating' | 'playing'

export interface Campaign {
  id: string
  name: string
  status: 'available' | 'coming_soon'
}

export type SeatKind = 'dm' | 'codm' | 'player'

export interface Seat {
  id: string
  kind: SeatKind
  name: string
  occupied: boolean
  ready: boolean
  characterId: string
}

export interface StatDef {
  key: string
  name: string
  min: number
  max: number
  base: number
}

export interface SkillDef {
  id: string
  name: string
  statKey: string // may be ""
  description?: string // DM text, <= 500 chars
}

export interface ItemDef {
  id: string
  name: string
}

export interface ItemStack {
  itemId: string
  qty: number
}

export interface Ruleset {
  name: string
  attribution: string // "" for custom
  stats: StatDef[]
  statPoints: number
  skills: SkillDef[]
  startingSkillPicks: number
  items: ItemDef[]
  startingItems: ItemStack[]
}

export interface Character {
  id: string
  seatId: string
  name: string
  portraitUrl: string
  stats: Record<string, number>
  skillIds: string[]
  inventory: ItemStack[]
}

export interface Grid {
  cellSize: number
  offsetX: number
  offsetY: number
  cols: number
  rows: number
}

export interface FogRect {
  id: string
  col: number
  row: number
  w: number
  h: number
}

export interface Scene {
  id: string
  name: string
  imageUrl: string
  grid: Grid
  fog: FogRect[]
}

export interface Token {
  id: string
  sceneId: string
  characterId: string // "" = NPC token
  name: string
  imageUrl: string
  col: number
  row: number
}

export interface Note {
  id: string
  title: string // 1-80
  body: string // <= 20000
  updatedAt: number // unix ms
}

export interface RoomState {
  status: RoomStatus
  campaigns: Campaign[]
  campaignId: string
  seats: Seat[]
  ruleset: Ruleset
  characters: Character[]
  scenes: Scene[]
  currentSceneId: string
  tokens: Token[]
}

export interface You {
  seatId: string // "" when not seated
}

// Server → client
export type ServerMessage =
  // build: random per server process. notes: DM/co-DM only, on connect and note.* events; no key = unchanged. Players never get it.
  | { type: 'state'; payload: { room: RoomState; you: You; build?: string; notes?: Note[] | null } }
  | { type: 'error'; payload: { message: string } }

type Empty = Record<string, never>

// Client → server: type → payload
export interface ClientEvents {
  hello: { clientId: string }
  'campaign.select': { campaignId: string; name: string }
  'seat.sit': { seatId: string; name: string }
  'seat.leave': Empty
  'seat.ready': { ready: boolean }
  'seat.add': Empty
  'seat.remove': { seatId: string }
  'seat.openCoDm': Empty
  'game.start': Empty
  'ruleset.update': { ruleset: Ruleset }
  'ruleset.loadPreset': { presetId: 'dnd5e-srd' }
  'character.submit': { name: string; portraitUrl: string; stats: Record<string, number>; skillIds: string[] }
  'character.grant': {
    characterId: string
    setStats?: Record<string, number>
    addSkillIds?: string[]
    removeSkillIds?: string[]
    addItems?: ItemStack[]
    removeItems?: ItemStack[]
  }
  'scene.add': { name: string; imageUrl: string; grid: Grid }
  'scene.update': { sceneId: string; name: string; grid: Grid }
  'scene.select': { sceneId: string }
  'token.addNpc': { sceneId: string; name: string; imageUrl: string }
  'token.move': { tokenId: string; col: number; row: number }
  'token.remove': { tokenId: string }
  'fog.add': { sceneId: string; col: number; row: number; w: number; h: number }
  'fog.remove': { sceneId: string; fogId: string }
  'note.add': { title: string }
  'note.update': { id: string; title: string; body: string }
  'note.remove': { id: string }
  'note.move': { id: string; index: number }
}
