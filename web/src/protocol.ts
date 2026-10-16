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

// Text in both UI languages; game content shows text[lang].
export interface L {
  th: string
  en: string
}

export interface RosterSkill {
  id: string
  name: L
  desc: L
  cost: number // mana 0-10
}

export interface RosterItem {
  name: L
  qty: number
}

export interface CharStats {
  hpMax: number
  str: number
  end: number
  agi: number
  mpMax: number
}

// A character set: a fixed group of roster entries the DM switches on in the lobby.
export interface CharSet {
  id: string // "core" | "legend"
  name: L
}

// One pre-made character (server-owned, DM-editable).
export interface RosterEntry {
  id: string // classId: core keeps its 10 ids (knight, wizard, ...); legend = "legend-<class>"
  set: string // CharSet id; server-owned, a client-sent value is ignored
  heroName: L
  className: L
  blurb: L
  role: string
  art: string // "" = class art by id (src/assets/classes/<id>.svg) | "/uploads/<name>"
  stats: CharStats
  skills: RosterSkill[] // 1-5
  items: RosterItem[] // 0-10
}

// innate = from the roster, fixed; custom = player-made (name/desc same text in th and en), at most 3.
export interface CharSkill {
  id: string
  kind: 'innate' | 'custom'
  name: L
  desc: L
  cost: number
}

export interface Character {
  id: string
  seatId: string
  classId: string
  name: L // heroName at pick time
  art: string
  hp: number
  hpMax: number
  mp: number
  mpMax: number
  str: number
  end: number
  agi: number
  skills: CharSkill[]
  items: RosterItem[]
}

export interface FeedEntry {
  at: number // unix ms
  text: L
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
  npcId?: string // set on tokens placed from the NPC pool
  hidden?: boolean // players never receive hidden tokens
  hp?: number // NPC only, DM/co-DM only (players never get hp/hpMax)
  hpMax?: number
}

// DM's NPC pool entry (white label). Sent to DM/co-DM only.
export interface Npc {
  id: string
  name: string // 1-40
  imageUrl: string // "/uploads/<file>" or "" = default silhouette
  hpMax: number // 1-999
  hidden: boolean // default for tokens placed from it
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
  sets: CharSet[] // every set, fixed
  activeSets: string[] // sets players pick from, in `sets` order; default ["core"]
  roster: RosterEntry[] // DM/co-DM: every set; players: active sets only
  characters: Character[]
  feed: FeedEntry[] // latest 50
  scenes: Scene[]
  currentSceneId: string
  tokens: Token[]
  combat: boolean // combat mode, set by DM/co-DM
}

// A declared attack waiting for the DM's hit/miss. DM/co-DM only.
export interface Attack {
  id: string
  byCharacterId: string
  byName: L // attacker's hero name
  targetTokenId: string
  targetName: string
  at: number // unix ms
}

// One line of the DM log ("who did what"). DM/co-DM only.
export interface LogEntry {
  at: number // unix ms
  actor: string
  text: L
}

export interface You {
  seatId: string // "" when not seated
}

// Server → client
export type ServerMessage =
  // build: random per server process. notes/npcs: DM/co-DM only, sent on connect and on their own events; no key = unchanged. Players never get them.
  // attacks: DM/co-DM only, every state. myAttack: a player's own pending attack. dmLog: DM/co-DM only, full list on
  // hello / becoming DM; no key = keep the previous list (later lines come as "log" messages).
  | {
      type: 'state'
      payload: {
        room: RoomState
        you: You
        build?: string
        notes?: Note[] | null
        npcs?: Npc[] | null
        attacks?: Attack[] | null
        myAttack?: { id: string; targetName: string } | null
        dmLog?: LogEntry[] | null
      }
    }
  | { type: 'error'; payload: { message: string } }
  // Not in state: shown once, to every connection, after attack.resolve.
  | { type: 'banner'; payload: { kind: 'hit' | 'miss'; text: L } }
  // Appended to the DM log; DM/co-DM connections only.
  | { type: 'log'; payload: LogEntry }

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
  'sets.set': { sets: string[] } // DM/co-DM, lobby only, at least one
  'game.start': Empty
  'character.pick': { classId: string }
  'skill.addCustom': { name: string; desc: string; cost: number }
  'skill.editCustom': { characterId: string; skillId: string; name: string; desc: string; cost: number }
  'skill.removeCustom': { characterId: string; skillId: string }
  'skill.use': { skillId: string }
  'character.adjust': { characterId: string; hp?: number; mp?: number } // deltas
  'character.setStats': { characterId: string } & CharStats
  'party.rest': Empty
  'item.give': { characterId: string; name: string; qty: number }
  'item.remove': { characterId: string; index: number; qty: number }
  'roster.update': { classId: string; entry: RosterEntry }
  'roster.reset': { classId?: string } | { set: string } // {} = every entry of every set
  'scene.add': { name: string; imageUrl: string; grid: Grid }
  'scene.update': { sceneId: string; name: string; grid: Grid }
  'scene.select': { sceneId: string }
  'npc.add': { name: string; imageUrl: string; hpMax: number; hidden: boolean }
  'npc.update': { id: string; name: string; imageUrl: string; hpMax: number; hidden: boolean }
  'npc.remove': { id: string }
  'npc.place': { npcId: string; sceneId: string }
  'token.setHidden': { tokenId: string; hidden: boolean }
  'token.adjustHp': { tokenId: string; delta: number }
  'token.move': { tokenId: string; col: number; row: number }
  'token.remove': { tokenId: string }
  'fog.add': { sceneId: string; col: number; row: number; w: number; h: number }
  'fog.remove': { sceneId: string; fogId: string }
  'note.add': { title: string }
  'note.update': { id: string; title: string; body: string }
  'note.remove': { id: string }
  'note.move': { id: string; index: number }
  'combat.set': { on: boolean }
  'attack.declare': { targetTokenId: string }
  'attack.cancel': Empty
  'attack.resolve': { attackId: string; hit: boolean; damage?: number }
}
