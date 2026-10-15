import { useEffect, useRef, useState, type CSSProperties } from 'react'
import type { Grid, Token } from '../protocol'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import MapStage from '../table/MapStage'
import DmPanel from '../table/DmPanel'
import SheetPanel from '../table/SheetPanel'

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
  const room = useStore((s) => s.room)!
  const { character, isDm } = useMe()
  const [viewId, setViewId] = useState('') // DM may view/calibrate a scene that is not active yet
  const [preview, setPreview] = useState<Grid | null>(null)
  const [fogMode, setFogMode] = useState(false)

  const sceneId = isDm ? viewId || room.currentSceneId || room.scenes.at(-1)?.id : room.currentSceneId
  const scene = room.scenes.find((s) => s.id === sceneId)
  // Players never see tokens under fog (their own token stays visible); the DM sees all.
  const fogged = (t: Token) => scene?.fog.some((f) => t.col >= f.col && t.col < f.col + f.w && t.row >= f.row && t.row < f.row + f.h)
  const tokens = room.tokens.filter(
    (t) => t.sceneId === scene?.id && (isDm || (t.characterId !== '' && t.characterId === character?.id) || !fogged(t)),
  )
  const canMove = (t: Token) => isDm || (t.characterId !== '' && t.characterId === character?.id)

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
    <div className="table" ref={tableRef} style={{ '--map-pct': `${mapPct}%` } as CSSProperties}>
      <div className="map-area">
        {scene ? (
          <>
            <h2>
              {scene.name}
              {isDm && scene.id !== room.currentSceneId && <small className="muted"> (ผู้เล่นยังไม่เห็น / not shown to players)</small>}
            </h2>
            <MapStage
              scene={scene}
              grid={isDm && preview ? preview : scene.grid}
              tokens={tokens}
              isDm={isDm}
              myCharacterId={character?.id ?? ''}
              fogMode={isDm && fogMode}
            />
          </>
        ) : (
          <p data-testid="no-scene">{isDm ? 'อัปแมพเพื่อสร้างฉาก / Upload a map to create a scene' : 'รอ DM เลือกฉาก / Waiting for the DM to pick a scene'}</p>
        )}
      </div>
      <div
        className="splitter"
        role="separator"
        tabIndex={0}
        aria-label="ขนาดแมพ / Map width"
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
        {character && <SheetPanel character={character} />}
        {isDm && (
          <DmPanel scene={scene} onView={setViewId} onPreview={setPreview} fogMode={fogMode} setFogMode={setFogMode} />
        )}
        {scene && (
          <ul className="token-list" data-testid="token-list" aria-label="Token">
            {tokens.map((t) => (
              <TokenRow key={t.id + t.col + ',' + t.row} t={t} movable={canMove(t)} removable={isDm && t.characterId === ''} />
            ))}
          </ul>
        )}
      </aside>
    </div>
  )
}

// Keyboard/screen-reader alternative to clicking on the canvas; also shows each token's cell.
function TokenRow({ t, movable, removable }: { t: Token; movable: boolean; removable: boolean }) {
  const [col, setCol] = useState(t.col)
  const [row, setRow] = useState(t.row)
  return (
    <li data-testid={`token-${t.id}`}>
      <span>{t.name}</span>{' '}
      <span data-testid={`token-pos-${t.id}`}>
        ({t.col}, {t.row})
      </span>
      {movable && (
        <form
          className="row inline"
          onSubmit={(e) => {
            e.preventDefault()
            send('token.move', { tokenId: t.id, col, row })
          }}
        >
          <input aria-label={`${t.name} col`} type="number" min={0} value={col} onChange={(e) => setCol(Number(e.target.value))} />
          <input aria-label={`${t.name} row`} type="number" min={0} value={row} onChange={(e) => setRow(Number(e.target.value))} />
          <button data-testid={`token-move-${t.id}`} type="submit">
            ย้าย / Move
          </button>
        </form>
      )}
      {removable && (
        <button className="danger" data-testid={`token-remove-${t.id}`} onClick={() => send('token.remove', { tokenId: t.id })}>
          ลบ / Remove
        </button>
      )}
    </li>
  )
}
