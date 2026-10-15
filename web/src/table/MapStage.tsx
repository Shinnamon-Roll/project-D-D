import { useEffect, useRef, useState } from 'react'
import { Circle, Group, Image as KImage, Layer, Line, Rect, Stage, Text } from 'react-konva'
import type Konva from 'konva'
import type { Grid, Scene, Token } from '../protocol'
import { send } from '../ws'

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

export default function MapStage({ scene, grid, tokens, isDm, myCharacterId, fogMode }: Props) {
  const map = useImg(scene.imageUrl)
  const [selected, setSelected] = useState('')
  const [drag, setDrag] = useState<{ a: Cell; b: Cell } | null>(null)
  const dragging = useRef(false)
  const { cellSize: cs, offsetX: ox, offsetY: oy, cols, rows } = grid
  const width = Math.max(map?.width ?? 0, ox + cols * cs)
  const height = Math.max(map?.height ?? 0, oy + rows * cs)

  const cellAt = (e: Konva.KonvaEventObject<PointerEvent>): Cell | null => {
    const p = e.target.getStage()?.getPointerPosition()
    if (!p) return null
    const col = Math.floor((p.x - ox) / cs)
    const row = Math.floor((p.y - oy) / cs)
    return col >= 0 && col < cols && row >= 0 && row < rows ? { col, row } : null
  }
  const canMove = (t: Token) => isDm || (t.characterId !== '' && t.characterId === myCharacterId)

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
  const onClick = (e: Konva.KonvaEventObject<PointerEvent>) => {
    if (fogMode || !selected) return
    const c = cellAt(e)
    if (c) send('token.move', { tokenId: selected, ...c })
    setSelected('')
  }

  const lines = []
  for (let c = 0; c <= cols; c++) lines.push([ox + c * cs, oy, ox + c * cs, oy + rows * cs])
  for (let r = 0; r <= rows; r++) lines.push([ox, oy + r * cs, ox + cols * cs, oy + r * cs])
  const dragRect = drag && {
    x: ox + Math.min(drag.a.col, drag.b.col) * cs,
    y: oy + Math.min(drag.a.row, drag.b.row) * cs,
    width: (Math.abs(drag.a.col - drag.b.col) + 1) * cs,
    height: (Math.abs(drag.a.row - drag.b.row) + 1) * cs,
  }

  return (
    <div
      className={`stage ${fogMode ? 'fog-mode' : ''}`}
      data-testid="map-stage"
      data-scene-id={scene.id}
      data-fog-count={scene.fog.length}
      data-selected-token={selected}
      aria-label={`แมพ / Map: ${scene.name}`}
    >
      <Stage width={width} height={height} onPointerDown={onDown} onPointerMove={onMove} onPointerUp={onUp} onPointerClick={onClick}>
        <Layer listening={false}>
          {map && <KImage image={map} />}
          {lines.map((p, i) => (
            <Line key={i} points={p} stroke="rgba(255,255,255,0.45)" strokeWidth={1} />
          ))}
        </Layer>
        <Layer listening={!fogMode}>
          {tokens.map((t) => (
            <TokenNode
              key={t.id}
              t={t}
              grid={grid}
              selected={t.id === selected}
              onSelect={canMove(t) ? () => setSelected(t.id === selected ? '' : t.id) : undefined}
            />
          ))}
        </Layer>
        <Layer listening={isDm && fogMode}>
          {scene.fog.map((f) => (
            <Rect
              key={f.id}
              x={ox + f.col * cs}
              y={oy + f.row * cs}
              width={f.w * cs}
              height={f.h * cs}
              fill="#0b0b0f"
              opacity={isDm ? 0.5 : 1}
              stroke={isDm ? '#ef4444' : undefined}
              onPointerDown={(e) => (e.cancelBubble = true)}
              onPointerClick={(e) => {
                e.cancelBubble = true
                send('fog.remove', { sceneId: scene.id, fogId: f.id })
              }}
            />
          ))}
          {dragRect && <Rect {...dragRect} fill="rgba(239,68,68,0.25)" stroke="#ef4444" dash={[6, 4]} listening={false} />}
        </Layer>
      </Stage>
    </div>
  )
}

function TokenNode({ t, grid, selected, onSelect }: { t: Token; grid: Grid; selected: boolean; onSelect?: () => void }) {
  const img = useImg(t.imageUrl)
  const cs = grid.cellSize
  const r = cs * 0.42
  return (
    <Group
      x={grid.offsetX + t.col * cs + cs / 2}
      y={grid.offsetY + t.row * cs + cs / 2}
      onPointerClick={
        onSelect &&
        ((e) => {
          e.cancelBubble = true
          onSelect()
        })
      }
    >
      <Circle radius={r} fill="#4b5563" stroke={selected ? '#f59e0b' : '#fff'} strokeWidth={selected ? 4 : 2} />
      {img ? (
        <Group clipFunc={(ctx) => ctx.arc(0, 0, r, 0, Math.PI * 2)}>
          <KImage image={img} x={-r} y={-r} width={r * 2} height={r * 2} />
        </Group>
      ) : (
        <Text text={t.name.slice(0, 1).toUpperCase()} x={-r} y={-r} width={r * 2} height={r * 2} align="center" verticalAlign="middle" fontSize={r} fill="#fff" />
      )}
      <Text text={t.name} x={-cs} y={r + 2} width={cs * 2} align="center" fontSize={12} fill="#fff" shadowColor="#000" shadowBlur={3} shadowOpacity={1} />
    </Group>
  )
}
