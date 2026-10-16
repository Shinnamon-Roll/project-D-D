// Hand-written mirror of docs/design/protocol.md (rules v2). JSON field names must match exactly.

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

export interface RosterItem {
  name: L
  qty: number
}

export type AbilityId = 'str' | 'dex' | 'con' | 'int' | 'wis' | 'cha'
export type Abilities = Record<AbilityId, number> // scores 1-30

// Class data (internal/rules/classes.go), from GET /api/classes. Static: never in RoomState.
export interface LevelTable {
  id: string
  name: L
  values: string[] // 20 values, index = level - 1; "—" = none yet
}

export interface Resource {
  id: string
  name: L
  max: number[] | null // per level; null when maxFormula is used
  maxFormula: string // "abilityMod:<ab>" | "level*<n>" | ""
  refill: 'short' | 'long'
  shortRegain: number
  shortFromLevel: number
  fromLevel: number // 0 = from level 1
}

export interface Feature {
  level: number
  name: L
  desc: L
}

export interface Class {
  id: string
  name: L
  hitDie: number
  primary: AbilityId[]
  saves: AbilityId[]
  armor: L
  weapons: L
  tools: L
  skillChoices: { count: number; from: string[] }
  caster: 'none' | 'full' | 'half' | 'pact'
  castAbility: AbilityId | ''
  tables: LevelTable[]
  resources: Resource[]
  features: Feature[] // sorted by level
}

export interface SkillDef {
  id: string
  ability: AbilityId
}

export interface PactSlots {
  count: number
  level: number
}

// GET /api/classes body. Tables are indexed by level - 1.
export interface ClassesDoc {
  classes: Class[]
  skills: SkillDef[] // 18, sheet order
  proficiency: number[]
  fullSlots: number[][]
  halfSlots: number[][]
  pactSlots: PactSlots[]
  attribution: string // SRD 5.2.1 CC BY 4.0 credit, shown on the pick screen
}

// A character set: a fixed group of roster entries the DM switches on in the lobby.
export interface CharSet {
  id: string // "core" | "legend"
  name: L
}

// One pre-made hero (server-owned, DM-editable). Class data comes from Class classId.
export interface RosterEntry {
  id: string // hero id: core keeps its 10 ids (knight, wizard, ...); legend = "legend-<class>"
  set: string // CharSet id; server-owned, a client-sent value is ignored
  classId: string // Class id; server-owned, a client-sent value is ignored
  heroName: L
  blurb: L
  role: string // tank | damage | healer | support | control | hybrid
  art: string // "" = class art by id (src/assets/classes/<id>.svg) | "/uploads/<name>"
  abilities: Abilities
  ac: number // 1-30
  skills: string[] // proficient skill ids
  expertise: string[] // subset of skills
  items: RosterItem[] // 0-10
}

// A spell or feature the player typed in (no spell library in the app).
export interface CustomAbility {
  id: string
  name: string // 1-40
  kind: 'spell' | 'feature'
  level: number // spell level 0-9; always 0 for a feature
  desc: string // <= 300
}

export interface Character {
  id: string
  seatId: string
  classId: string // Class id (fighter, wizard, ...)
  heroId: string // roster entry picked (knight, legend-bard, ...)
  name: L // heroName at pick time
  art: string
  level: number // 1-20
  hp: number
  hpMax: number
  abilities: Abilities
  ac: number
  skills: string[]
  expertise: string[]
  extraSaves: AbilityId[] // save proficiencies beyond the class's, set by character.setSkills
  slotsUsed: number[] // 9, index 0 = 1st-level slots; full/half casters
  pactUsed: number
  resourcesUsed: Record<string, number>
  hitDiceUsed: number
  abilitiesCustom: CustomAbility[]
  items: RosterItem[]
  // computed by the server from class, level and abilities; clients never send these
  pb: number
  saves: AbilityId[] // class saves + extraSaves, sheet order
  slotsMax: number[] // 9
  pact: PactSlots
  resourcesMax: Record<string, number> // 0 = not gained yet
  hitDiceMax: number
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

// hit/miss: text "<attacker> → <target>"; talk/gift/surrender: text is the whole feed line.
export type BannerKind = 'hit' | 'miss' | 'talk' | 'gift' | 'surrender'

export type NpcActKind = 'talk' | 'attack' | 'gift' | 'surrender'

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
  // Not in state: shown once, to every connection, after attack.resolve or npc.act.
  | { type: 'banner'; payload: { kind: BannerKind; text: L } }
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
  'character.pick': { heroId: string }
  'character.adjust': { characterId: string; hp: number } // delta
  'character.level': { characterId: string; delta: 1 | -1 } // DM
  'character.setAbilities': { characterId: string; abilities: Abilities; ac: number } // DM
  // DM/co-DM, creating/playing: replaces all three lists. saves = extra save proficiencies only (not the class's).
  'character.setSkills': { characterId: string; skills: string[]; expertise: string[]; saves: AbilityId[] }
  'slot.set': { characterId: string; level: number | 'pact'; used: number } // owner or DM
  'resource.set': { characterId: string; resourceId: string; used: number } // owner or DM
  'hitDice.set': { characterId: string; used: number } // owner or DM
  'party.rest': { kind: 'short' | 'long' } // DM, playing only
  'ability.add': { name: string; kind: CustomAbility['kind']; level: number; desc: string } // own character
  'ability.update': { abilityId: string; name: string; kind: CustomAbility['kind']; level: number; desc: string }
  'ability.remove': { characterId: string; abilityId: string } // owner or DM
  'item.give': { characterId: string; name: string; qty: number }
  'item.remove': { characterId: string; index: number; qty: number }
  'roster.update': { classId: string; entry: RosterEntry } // classId = hero id
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
  'token.duplicate': { tokenId: string } // DM/co-DM: copy of an NPC token on the nearest free cell, named by the server
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
  // DM/co-DM, combat on: a visible NPC token acts on a player-character token in the same scene. Fields of other kinds are ignored.
  'npc.act': { tokenId: string; targetId: string; kind: NpcActKind; text?: string; hit?: boolean; damage?: number; item?: { name: string; qty: number } }
}
