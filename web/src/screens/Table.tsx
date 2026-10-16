import { useEffect, useRef, useState, type CSSProperties } from 'react'
import type { Grid, Token } from '../protocol'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import { useL, useT } from '../i18n'
import MapStage from '../table/MapStage'
import DmPanel from '../table/DmPanel'
import SheetPanel from '../table/SheetPanel'
import FeedPanel from '../table/FeedPanel'
import PartyBar from '../table/PartyBar'
import { classArt } from '../classArt'
import { EyeIcon, EyeOffIcon, SwordsIcon } from '../table/icons'
import { AttackAlerts, AttackControl } from '../table/Combat'
import NpcAction from '../table/NpcAction'

const PCT_KEY = 'table.mapPct'
const MAP_MIN = 40
const MAP_MAX = 78
const INFO_MIN = 280 // px; mirrors the clamp() in .table's grid-template-columns
// Highest map % that still leaves the info pane its 280px on a table w px wide.
const maxPct = (w: number) => (w ? Math.max(MAP_MIN, Math.min(MAP_MAX, Math.floor(100 - (INFO_MIN * 100) / w))) : MAP_MAX)
const clampPct = (n: number, w = 0) => Math.round(Math.min(maxPct(w), Math.max(MAP_MIN, n)))
const savePct = (n: number) => {
  try {
    localStorage.setItem(PCT_KEY, String(n))
  } catch {
    /* storage blocked: keep the in-memory value */
  }
}

export default function Table() {
  const t = useT()
  const room = useStore((s) => s.room)!
  const { character, isDm } = useMe()
  const [viewId, setViewId] = useState('') // DM may view/calibrate a scene that is not active yet
  const [preview, setPreview] = useState<Grid | null>(null)
  const [fogMode, setFogMode] = useState(false)
  const myAttack = useStore((s) => s.myAttack)
  const [aimWant, setAim] = useState<string | null>(null) // attack targeting: null = off, '' = nothing picked yet

  const sceneId = isDm ? viewId || room.currentSceneId || room.scenes.at(-1)?.id : room.currentSceneId
  const scene = room.scenes.find((s) => s.id === sceneId)
  // Players never see tokens under fog (their own token stays visible); the DM sees all.
  const fogged = (t: Token) => scene?.fog.some((f) => t.col >= f.col && t.col < f.col + f.w && t.row >= f.row && t.row < f.row + f.h)
  const tx = useL()
  const heroName = (t: Token) => room.characters.find((c) => c.id === t.characterId)?.name
  const tokens = room.tokens
    .filter((t) => t.sceneId === scene?.id && (isDm || (t.characterId !== '' && t.characterId === character?.id) || !fogged(t)))
    .map((t) => (heroName(t) ? { ...t, name: tx(heroName(t)) } : t)) // hero names follow the UI language
  // Character tokens show their class art (or upload) unless the token has its own image.
  const charArt = Object.fromEntries(room.characters.map((c) => [c.id, classArt(c.classId, c.art)]))
  const canMove = (t: Token) => isDm || (t.characterId !== '' && t.characterId === character?.id)
  // Attack targets: what this player can see in the scene, minus their own token. The server re-checks.
  const targets = isDm || !character ? [] : tokens.filter((t) => t.characterId !== character.id)
  const aiming = room.combat && !myAttack && character ? aimWant : null
  const aim = aiming && !targets.some((t) => t.id === aiming) ? '' : aiming
  // DM NPC action dialog: target '' = pick in the dialog. Shown only while combat is on and the NPC token is still here.
  const [act, setAct] = useState<{ npc: string; target: string } | null>(null)
  // only in the scene being played: the server refuses NPC actions anywhere else, and hidden tokens can't take part
  const canAct = isDm && room.combat && scene?.id === room.currentSceneId
  const actNpc = canAct ? tokens.find((t) => t.id === act?.npc && t.characterId === '' && !t.hidden) : undefined
  const players = tokens.filter((t) => t.characterId !== '' && !t.hidden)
  useEffect(() => {
    if (!room.combat) setAct(null) // combat ended: a reopened combat starts with no dialog
  }, [room.combat])

  // Splitter: map pane share of the table width (info pane = rest, 22–60%, never under 280px).
  const tableRef = useRef<HTMLDivElement>(null)
  const [tableW, setTableW] = useState(0)
  const [wantPct, setWantPct] = useState(() => {
    try {
      return clampPct(Number(localStorage.getItem(PCT_KEY)) || 72)
    } catch {
      return 72
    }
  })
  const mapPct = clampPct(wantPct, tableW) // what the layout really shows, also on narrow windows
  useEffect(() => {
    const ro = new ResizeObserver(([e]) => setTableW(e.contentRect.width))
    ro.observe(tableRef.current!)
    return () => ro.disconnect()
  }, [])
  // Drag: at most one layout update per frame; the pending pointer x waits in a ref.
  const dragX = useRef<number | null>(null)
  const frame = useRef(0)
  const pctAt = (x: number) => {
    const box = tableRef.current!.getBoundingClientRect()
    return clampPct(((x - box.left) / box.width) * 100, box.width)
  }
  useEffect(() => () => cancelAnimationFrame(frame.current), [])

  return (
    <div className={`table${room.combat ? ' in-combat' : ''}`} ref={tableRef} style={{ '--map-pct': `${mapPct}%` } as CSSProperties}>
      {isDm && <AttackAlerts />}
      <div className="map-area">
        {isDm && <PartyBar />}
        {(scene || room.combat) && (
          <div className="map-head">
            {scene && (
              <h2>
                {scene.name}
                {isDm && scene.id !== room.currentSceneId && <small className="muted"> {t('notShown')}</small>}
              </h2>
            )}
            {room.combat && (
              <span className="combat-badge" data-testid="combat-indicator" role="status">
                <SwordsIcon />
                {t('combatOn')}
              </span>
            )}
          </div>
        )}
        {scene ? (
          <>
            <MapStage
              scene={scene}
              grid={isDm && preview ? preview : scene.grid}
              tokens={tokens}
              isDm={isDm}
              myCharacterId={character?.id ?? ''}
              fogMode={isDm && fogMode}
              charArt={charArt}
              aim={aim}
              onAim={setAim}
              onNpcAct={canAct ? (npc, target) => setAct({ npc, target }) : undefined}
            />
          </>
        ) : (
          <p data-testid="no-scene">{t(isDm ? 'uploadMapHint' : 'waitingScene')}</p>
        )}
      </div>
      <div
        className="splitter"
        role="separator"
        tabIndex={0}
        aria-label={t('mapWidth')}
        aria-orientation="vertical"
        aria-valuenow={mapPct}
        aria-valuemin={MAP_MIN}
        aria-valuemax={maxPct(tableW)}
        data-testid="splitter"
        onPointerDown={(e) => {
          e.currentTarget.setPointerCapture(e.pointerId)
          e.preventDefault()
        }}
        onPointerMove={(e) => {
          if (!e.currentTarget.hasPointerCapture(e.pointerId)) return
          dragX.current = e.clientX
          frame.current ||= requestAnimationFrame(() => {
            frame.current = 0
            if (dragX.current !== null) setWantPct(pctAt(dragX.current))
          })
        }}
        onLostPointerCapture={() => {
          cancelAnimationFrame(frame.current)
          frame.current = 0
          if (dragX.current === null) return
          const n = pctAt(dragX.current)
          dragX.current = null
          setWantPct(n)
          savePct(n)
        }}
        onKeyDown={(e) => {
          const next = { ArrowLeft: mapPct - 2, ArrowRight: mapPct + 2, Home: MAP_MIN, End: MAP_MAX }[e.key]
          if (next === undefined) return
          e.preventDefault()
          const n = clampPct(next, tableW)
          setWantPct(n)
          savePct(n)
        }}
      />
      <aside className="side">
        {character && (
          <SheetPanel character={character}>
            {room.combat && <AttackControl targets={targets} aim={aim} setAim={setAim} />}
          </SheetPanel>
        )}
        {isDm && (
          <DmPanel scene={scene} onView={setViewId} onPreview={setPreview} fogMode={fogMode} setFogMode={setFogMode} />
        )}
        {scene && (
          <ul className="token-list" data-testid="token-list" aria-label={t('tokens')}>
            {tokens.map((t) => (
              <TokenRow key={t.id} t={t} dmNpc={isDm && t.characterId === ''} onAct={canAct && !t.hidden ? () => setAct({ npc: t.id, target: '' }) : undefined} />
            ))}
          </ul>
        )}
        <FeedPanel />
      </aside>
      {actNpc && <NpcAction key={`${act!.npc}:${act!.target}`} npc={actNpc} targets={players} target={act!.target} onClose={() => setAct(null)} />}
    </div>
  )
}

// Each token's cell, plus the DM's NPC controls. Moving is by tapping the map only.
function TokenRow({ t, dmNpc, onAct }: { t: Token; dmNpc: boolean; onAct?: () => void }) {
  const tr = useT()
  return (
    <li data-testid={`token-${t.id}`} className={t.hidden ? 'is-hidden' : undefined}>
      <span>{t.name}</span>{' '}
      <span data-testid={`token-pos-${t.id}`}>
        ({t.col}, {t.row})
      </span>
      {dmNpc && (
        <button
          className="eye-toggle"
          data-testid={`token-hide-${t.id}`}
          aria-pressed={!!t.hidden}
          aria-label={tr(t.hidden ? 'tokenIsHidden' : 'tokenIsShown', { name: t.name })}
          title={tr(t.hidden ? 'tokenIsHidden' : 'tokenIsShown', { name: t.name })}
          onClick={() => send('token.setHidden', { tokenId: t.id, hidden: !t.hidden })}
        >
          {t.hidden ? <EyeOffIcon /> : <EyeIcon />}
        </button>
      )}
      {dmNpc && t.hp !== undefined && t.hpMax !== undefined && (
        <span className="hp-ctl" role="group" aria-label={tr('hpOf', { name: t.name })}>
          <button aria-label={tr('hpDec', { name: t.name })} data-testid={`token-hp-dec-${t.id}`} disabled={t.hp <= 0} onClick={() => send('token.adjustHp', { tokenId: t.id, delta: -1 })}>
            −
          </button>
          <output data-testid={`token-hp-${t.id}`} className={t.hp === 0 ? 'down' : undefined}>
            {t.hp}/{t.hpMax}
          </output>
          <button aria-label={tr('hpInc', { name: t.name })} data-testid={`token-hp-inc-${t.id}`} disabled={t.hp >= t.hpMax} onClick={() => send('token.adjustHp', { tokenId: t.id, delta: 1 })}>
            +
          </button>
        </span>
      )}
      {dmNpc && onAct && (
        <button className="npc-act-open" data-testid={`npc-act-open-${t.id}`} aria-label={tr('npcActOf', { name: t.name })} onClick={onAct}>
          <SwordsIcon />
          {tr('npcAct')}
        </button>
      )}
      {dmNpc && (
        <button className="danger" data-testid={`token-remove-${t.id}`} onClick={() => send('token.remove', { tokenId: t.id })}>
          {tr('remove')}
        </button>
      )}
    </li>
  )
}
