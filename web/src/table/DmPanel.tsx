import { useEffect, useState } from 'react'
import type { Grid, Scene } from '../protocol'
import { useStore } from '../store'
import { send, upload } from '../ws'
import GrantPanel from './GrantPanel'
import RulesetEditor from '../screens/RulesetEditor'

const validGrid = (g: Grid) =>
  g.cellSize >= 10 && g.cellSize <= 400 && g.cols >= 1 && g.cols <= 200 && g.rows >= 1 && g.rows <= 200

function imageSize(url: string): Promise<{ w: number; h: number }> {
  return new Promise((res) => {
    const i = new Image()
    i.onload = () => res({ w: i.naturalWidth, h: i.naturalHeight })
    i.onerror = () => res({ w: 1000, h: 750 })
    i.src = url
  })
}

interface Props {
  scene?: Scene // scene the DM is viewing
  onView: (sceneId: string) => void
  onPreview: (g: Grid | null) => void
  fogMode: boolean
  setFogMode: (on: boolean) => void
}

export default function DmPanel({ scene, onView, onPreview, fogMode, setFogMode }: Props) {
  const room = useStore((s) => s.room)!
  const [newName, setNewName] = useState('')
  const [npcName, setNpcName] = useState('')
  const [npcUrl, setNpcUrl] = useState('')

  async function addScene(file: File) {
    const url = await upload(file)
    if (!url) return
    const { w, h } = await imageSize(url)
    const clamp = (n: number) => Math.min(200, Math.max(1, Math.round(n / 50)))
    send('scene.add', {
      name: newName.trim() || file.name.replace(/\.[^.]+$/, '').slice(0, 32),
      imageUrl: url,
      grid: { cellSize: 50, offsetX: 0, offsetY: 0, cols: clamp(w), rows: clamp(h) },
    })
    setNewName('')
  }

  return (
    <section className="panel" data-testid="dm-panel">
      <h2>DM</h2>

      <h3>ฉาก / Scenes</h3>
      {room.scenes.length > 0 && (
        <div className="row">
          <label className="grow">
            ดู/แก้ฉาก / View scene
            <select data-testid="scene-select" value={scene?.id ?? ''} onChange={(e) => onView(e.target.value)}>
              {room.scenes.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name}
                  {s.id === room.currentSceneId ? ' ★' : ''}
                </option>
              ))}
            </select>
          </label>
          <button
            className="primary"
            data-testid="scene-activate"
            disabled={!scene || scene.id === room.currentSceneId}
            onClick={() => scene && send('scene.select', { sceneId: scene.id })}
          >
            ให้ทุกคนไปฉากนี้ / Show to all
          </button>
        </div>
      )}
      <div className="row">
        <label className="grow">
          ชื่อฉากใหม่ / New scene name
          <input data-testid="scene-name" value={newName} maxLength={32} onChange={(e) => setNewName(e.target.value)} />
        </label>
        <label>
          อัปแมพ / Upload map
          <input
            data-testid="map-upload"
            type="file"
            accept="image/png,image/jpeg,image/webp,image/gif"
            onChange={(e) => {
              const f = e.target.files?.[0]
              if (f) addScene(f)
              e.target.value = ''
            }}
          />
        </label>
      </div>
      {scene && <SceneEditor key={scene.id + JSON.stringify(scene.grid) + scene.name} scene={scene} onPreview={onPreview} />}

      {scene && (
        <>
          <h3>Fog</h3>
          <button data-testid="fog-toggle" aria-pressed={fogMode} onClick={() => setFogMode(!fogMode)}>
            {fogMode ? 'ปิดโหมด fog / Fog mode: ON' : 'เปิดโหมด fog / Fog mode: OFF'}
          </button>
          {fogMode && <p className="muted">ลากบนแมพเพื่อเพิ่ม · คลิกสี่เหลี่ยมเพื่อลบ / Drag to add, click a rect to remove</p>}
          <ul className="chips">
            {scene.fog.map((f) => (
              <li key={f.id}>
                {f.col},{f.row} {f.w}×{f.h}
                <button aria-label="ลบ fog / Remove fog" data-testid={`fog-remove-${f.id}`} onClick={() => send('fog.remove', { sceneId: scene.id, fogId: f.id })}>
                  ×
                </button>
              </li>
            ))}
          </ul>

          <h3>NPC</h3>
          <div className="row">
            <label className="grow">
              ชื่อ NPC / NPC name
              <input data-testid="npc-name" value={npcName} maxLength={32} onChange={(e) => setNpcName(e.target.value)} />
            </label>
            <label>
              รูป / Image
              <input
                data-testid="npc-upload"
                type="file"
                accept="image/png,image/jpeg,image/webp,image/gif"
                onChange={async (e) => {
                  const f = e.target.files?.[0]
                  if (f) setNpcUrl((await upload(f)) ?? '')
                }}
              />
            </label>
            <button
              data-testid="npc-add"
              disabled={!npcName.trim()}
              onClick={() => {
                send('token.addNpc', { sceneId: scene.id, name: npcName.trim(), imageUrl: npcUrl })
                setNpcName('')
              }}
            >
              + NPC
            </button>
          </div>
        </>
      )}

      <details>
        <summary>ให้ของ/สกิล/stat / Grant</summary>
        <GrantPanel />
      </details>
      <details>
        <summary>ชุดกฎ / Ruleset</summary>
        <RulesetEditor key={JSON.stringify(room.ruleset)} />
      </details>
    </section>
  )
}

function SceneEditor({ scene, onPreview }: { scene: Scene; onPreview: (g: Grid | null) => void }) {
  const [name, setName] = useState(scene.name)
  const [g, setG] = useState(scene.grid)
  useEffect(() => () => onPreview(null), [onPreview]) // drop the preview when this editor goes away
  const change = (k: keyof Grid, v: string) => {
    const next = { ...g, [k]: Number(v) }
    setG(next)
    if (validGrid(next)) onPreview(next)
  }
  const fields: [keyof Grid, string][] = [
    ['cellSize', 'ขนาดช่อง / Cell px'],
    ['offsetX', 'Offset X'],
    ['offsetY', 'Offset Y'],
    ['cols', 'คอลัมน์ / Cols'],
    ['rows', 'แถว / Rows'],
  ]
  return (
    <div className="scene-editor">
      <label>
        ชื่อฉาก / Scene name
        <input data-testid="scene-edit-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
      </label>
      <div className="stats-grid">
        {fields.map(([k, label]) => (
          <label key={k}>
            {label}
            <input type="number" data-testid={`grid-${k}`} value={g[k]} onChange={(e) => change(k, e.target.value)} />
          </label>
        ))}
      </div>
      <button
        data-testid="scene-save"
        disabled={!validGrid(g) || !name.trim()}
        onClick={() => send('scene.update', { sceneId: scene.id, name: name.trim(), grid: g })}
      >
        บันทึกฉาก / Save scene
      </button>
    </div>
  )
}
