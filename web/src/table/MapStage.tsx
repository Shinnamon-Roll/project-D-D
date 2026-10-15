import { useEffect, useRef, useState, type CSSProperties } from 'react'
import { Circle, Group, Image as KImage, Layer, Line, Rect, Stage, Text } from 'react-konva'
import type Konva from 'konva'
import type { Grid, Scene, Token } from '../protocol'
import { send } from '../ws'
import { useT } from '../i18n'

export function useImg(url: string) {
  const [img, setImg] = useState<HTMLImageElement | null>(null)
  useEffect(() => {
    if (!url) return setImg(null)
    const i = new window.Image()
    i.onload = () => setImg(i)
    i.src = url
    return () => {
      i.onload = null
    }
  }, [url])
  return img
}

interface Props {
  scene: Scene
  grid: Grid // scene.grid, or the DM's unsaved preview
  tokens: Token[]
  isDm: boolean
  myCharacterId: string
  fogMode: boolean
}

type Cell = { col: number; row: number }
// Where the grid sits on the stage: top-left (x, y), drawn cell size k, drawn width w and height h of image-or-grid, whichever is taller.
type Fit = { x: number; y: number; k: number; w: number; h: number }

// Konva can't read CSS vars; mirrors the tokens in docs/design/visual-theme.md.
const C = {
  grid: 'rgba(43,29,18,.45)',
  parchment: '#efe2c3',
  gold: '#c9a227',
  oxblood: '#7a1f1f',
  oxbloodHi: '#9b2d2d',
  iron: '#4a4f55',
  fog: '#120d09',
  font: 'Trirong, Georgia, serif',
}

export default function MapStage({ scene, grid, tokens, isDm, myCharacterId, fogMode }: Props) {
  const t = useT()
  const map = useImg(scene.imageUrl)
  const [selected, setSelected] = useState('')
  const [drag, setDrag] = useState<{ a: Cell; b: Cell } | null>(null)
  const dragging = useRef(false)
  const box = useRef<HTMLDivElement>(null)
  const [size, setSize] = useState({ w: 0, h: 0 })
  useEffect(() => {
    const ro = new ResizeObserver(([e]) => setSize({ w: Math.floor(e.contentRect.width), h: Math.floor(e.contentRect.height) }))
    ro.observe(box.current!)
    return () => ro.disconnect()
  }, [])

  const { cols, rows } = grid
  // Grid spans the image width (cols) at square cells; rounded rows may run past the image bottom,
  // so contain-fit the taller of image and grid (image drawn at the top) to keep every row on screen.
  const imgW = map?.width || cols * grid.cellSize
  const imgH = map?.height || rows * grid.cellSize
  const gridH = Math.max(imgH, (rows * imgW) / cols)
  const s = Math.min(size.w / imgW, size.h / gridH) || 0
  const f: Fit = { w: imgW * s, h: gridH * s, x: (size.w - imgW * s) / 2, y: (size.h - gridH * s) / 2, k: (imgW * s) / cols }
  const clip = { clipX: f.x, clipY: f.y, clipWidth: f.w, clipHeight: f.h }

  const cellAt = (e: Konva.KonvaEventObject<PointerEvent>): Cell | null => {
    const p = e.target.getStage()?.getPointerPosition()
    if (!p || !f.k) return null
    const col = Math.floor((p.x - f.x) / f.k)
    const row = Math.floor((p.y - f.y) / f.k)
    return col >= 0 && col < cols && row >= 0 && row < rows ? { col, row } : null
  }
  const mine = tokens.find((t) => t.characterId !== '' && t.characterId === myCharacterId)

  const onDown = (e: Konva.KonvaEventObject<PointerEvent>) => {
    if (!fogMode) return
    const c = cellAt(e)
    if (c) {
      dragging.current = true
      setDrag({ a: c, b: c })
    }
  }
  const onMove = (e: Konva.KonvaEventObject<PointerEvent>) => {
    if (!dragging.current) return
    const c = cellAt(e)
    if (c) setDrag((d) => d && { a: d.a, b: c })
  }
  const onUp = () => {
    if (!dragging.current || !drag) return
    dragging.current = false
    const { a, b } = drag
    send('fog.add', {
      sceneId: scene.id,
      col: Math.min(a.col, b.col),
      row: Math.min(a.row, b.row),
      w: Math.abs(a.col - b.col) + 1,
      h: Math.abs(a.row - b.row) + 1,
    })
    setDrag(null)
  }
  // Tap a cell: the DM moves the selected token there; a player moves their own token.
  const onTap = (e: Konva.KonvaEventObject<PointerEvent>) => {
    if (fogMode) return
    const id = isDm ? selected : mine?.id
    const c = cellAt(e)
    const t = tokens.find((t) => t.id === id)
    if (t && c && (c.col !== t.col || c.row !== t.row)) send('token.move', { tokenId: t.id, ...c })
    setSelected('')
  }

  const lines = []
  for (let c = 0; c <= cols; c++) lines.push([f.x + c * f.k, f.y, f.x + c * f.k, f.y + rows * f.k])
  for (let r = 0; r <= rows; r++) lines.push([f.x, f.y + r * f.k, f.x + f.w, f.y + r * f.k])
  const dragRect = drag && {
    x: f.x + Math.min(drag.a.col, drag.b.col) * f.k,
    y: f.y + Math.min(drag.a.row, drag.b.row) * f.k,
    width: (Math.abs(drag.a.col - drag.b.col) + 1) * f.k,
    height: (Math.abs(drag.a.row - drag.b.row) + 1) * f.k,
  }

  return (
    <div
      ref={box}
      className={`stage ${fogMode ? 'fog-mode' : ''}`}
      data-testid="map-stage"
      data-scene-id={scene.id}
      data-fog-count={scene.fog.length}
      data-grid={`${cols}x${rows}`}
      data-selected-token={selected}
      style={{ '--map-aspect': imgW / gridH } as CSSProperties}
      aria-label={t('map', { name: scene.name })}
    >
      <Stage width={size.w} height={size.h} onPointerDown={onDown} onPointerMove={onMove} onPointerUp={onUp} onPointerClick={onTap}>
        <Layer listening={false} {...clip}>
          {map && <KImage image={map} x={f.x} y={f.y} width={f.w} height={imgH * s} />}
          {lines.map((p, i) => (
            <Line key={i} points={p} stroke={C.grid} strokeWidth={1} />
          ))}
        </Layer>
        <Layer listening={!fogMode}>
          {tokens.map((t) => (
            <TokenNode
              key={t.id}
              t={t}
              f={f}
              selected={t.id === selected}
              // Only the DM selects; a player's tap on any token falls through to the cell under it.
              onSelect={isDm ? () => setSelected(t.id === selected ? '' : t.id) : undefined}
            />
          ))}
        </Layer>
        <Layer listening={isDm && fogMode} {...clip}>
          {scene.fog.map((r) => (
            <Rect
              key={r.id}
              x={f.x + r.col * f.k}
              y={f.y + r.row * f.k}
              width={r.w * f.k}
              height={r.h * f.k}
              fill={C.fog}
              opacity={isDm ? 0.5 : 1}
              stroke={isDm ? C.oxblood : undefined}
              strokeWidth={2}
              onPointerDown={(e) => (e.cancelBubble = true)}
              onPointerClick={(e) => {
                e.cancelBubble = true
                send('fog.remove', { sceneId: scene.id, fogId: r.id })
              }}
            />
          ))}
          {dragRect && <Rect {...dragRect} fill="rgba(122,31,31,0.3)" stroke={C.oxbloodHi} dash={[6, 4]} listening={false} />}
        </Layer>
      </Stage>
    </div>
  )
}

function TokenNode({ t, f, selected, onSelect }: { t: Token; f: Fit; selected: boolean; onSelect?: () => void }) {
  const img = useImg(t.imageUrl)
  const k = f.k
  const r = k * 0.42
  const label = Math.min(14, k * 0.3)
  // Keep the name box inside the image so edge tokens' names are not cut by the canvas edge.
  const left = -(t.col + 0.5) * k
  const right = f.w + left - 3 * k
  const lx = Math.min(Math.max(-1.5 * k, left), right)
  const align = lx > -1.5 * k ? 'left' : lx < -1.5 * k ? 'right' : 'center'
  return (
    <Group
      x={f.x + t.col * k + k / 2}
      y={f.y + t.row * k + k / 2}
      onPointerClick={
        onSelect &&
        ((e) => {
          e.cancelBubble = true
          onSelect()
        })
      }
    >
      <Circle
        radius={r}
        fill={t.characterId ? C.oxblood : C.iron}
        stroke={selected ? C.gold : C.parchment}
        strokeWidth={selected ? Math.max(2, k * 0.08) : Math.max(1, k * 0.04)}
        shadowColor={selected ? C.gold : '#000'}
        shadowBlur={selected ? 12 : 4}
        shadowOpacity={0.8}
      />
      {img ? (
        <Group clipFunc={(ctx) => ctx.arc(0, 0, r, 0, Math.PI * 2)}>
          <KImage image={img} x={-r} y={-r} width={r * 2} height={r * 2} />
        </Group>
      ) : (
        <Text text={t.name.slice(0, 1).toUpperCase()} x={-r} y={-r} width={r * 2} height={r * 2} align="center" verticalAlign="middle" fontSize={r} fontFamily={C.font} fill={C.parchment} />
      )}
      {/* name hides when cells are too small to read it */}
      {label >= 8 && (
        <Text text={t.name} x={lx} y={r + 2} width={k * 3} align={align} fontSize={label} fontFamily={C.font} fill={C.parchment} shadowColor="#000" shadowBlur={3} shadowOpacity={1} listening={false} wrap="none" ellipsis />
      )}
    </Group>
  )
}
