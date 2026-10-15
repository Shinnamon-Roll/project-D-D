import { useEffect, useState } from 'react'
import type { Grid, Scene } from '../protocol'
import { useStore } from '../store'
import { send, upload } from '../ws'
import { useT, type TKey } from '../i18n'
import GrantPanel from './GrantPanel'
import NotesPanel from './NotesPanel'
import NpcPanel from './NpcPanel'
import { useImg } from './MapStage'
import RulesetEditor from '../screens/RulesetEditor'

const validGrid = (g: Grid) =>
  g.cellSize >= 10 && g.cellSize <= 400 && g.cols >= 1 && g.cols <= 200 && g.rows >= 1 && g.rows <= 200

// Square cells over the whole image: cell = w/cols image px, rows follow the aspect, no offsets.
// cellSize is an int on the server (10-400 check), so it is the rounded cell; rows use the exact one.
function gridFor(cols: number, w: number, h: number): Grid {
  return { cellSize: Math.round(w / cols), offsetX: 0, offsetY: 0, cols, rows: Math.max(1, Math.round((h * cols) / w)) }
}
const colsOk = (c: number, w: number, h: number) => c >= 4 && c <= 60 && validGrid(gridFor(c, w, h))

function imageSize(url: string): Promise<{ w: number; h: number }> {
  return new Promise((res) => {
    const i = new Image()
    i.onload = () => res({ w: i.naturalWidth, h: i.naturalHeight })
    i.onerror = () => res({ w: 1000, h: 750 })
    i.src = url
  })
}

type Section = 'scenes' | 'fog' | 'npc' | 'grant' | 'notes' | 'ruleset'
const SECTIONS: [Section, TKey][] = [
  ['scenes', 'sceneMap'],
  ['fog', 'fog'],
  ['npc', 'npc'],
  ['grant', 'grant'],
  ['notes', 'story'],
  ['ruleset', 'ruleset'],
]
const SECTION_KEY = 'dm.section'
function savedSection(): Section | '' {
  try {
    const v = localStorage.getItem(SECTION_KEY)
    if (v === '') return ''
    return SECTIONS.find(([id]) => id === v)?.[0] ?? 'scenes'
  } catch {
    return 'scenes'
  }
}

interface Props {
  scene?: Scene // scene the DM is viewing
  onView: (sceneId: string) => void
  onPreview: (g: Grid | null) => void
  fogMode: boolean
  setFogMode: (on: boolean) => void
}

export default function DmPanel({ scene, onView, onPreview, fogMode, setFogMode }: Props) {
  const t = useT()
  const room = useStore((s) => s.room)!
  const [newName, setNewName] = useState('')
  const [err, setErr] = useState('')
  const [open, setOpen] = useState<Section | ''>(savedSection)

  async function addScene(file: File) {
    const url = await upload(file)
    if (!url) return
    const { w, h } = await imageSize(url)
    // New scenes start at 20 columns; else the nearest valid count below, then above.
    const cols = [...Array(17)].map((_, i) => 20 - i).concat([...Array(40)].map((_, i) => 21 + i)).find((c) => colsOk(c, w, h))
    if (cols === undefined) {
      setErr(t('gridFitError', { w, h }))
      return
    }
    setErr('')
    send('scene.add', {
      name: newName.trim() || file.name.replace(/\.[^.]+$/, '').slice(0, 32),
      imageUrl: url,
      grid: gridFor(cols, w, h),
    })
    setNewName('')
  }

  // One menu section open at a time; leaving the fog section ends fog drawing.
  const choose = (next: Section | '') => {
    setOpen(next)
    if (next !== 'fog') setFogMode(false)
    try {
      localStorage.setItem(SECTION_KEY, next)
    } catch {
      /* storage blocked: remembered for this page only */
    }
  }
  const needScene = <p className="muted">{t('uploadMapHint')}</p>

  const content = (id: Section) => {
    switch (id) {
      case 'scenes':
        return (
          <>
            {room.scenes.length > 0 && (
              <div className="row">
                <label className="grow">
                  {t('viewScene')}
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
                  {t('showToAll')}
                </button>
              </div>
            )}
            <div className="row">
              <label className="grow">
                {t('newSceneName')}
                <input data-testid="scene-name" value={newName} maxLength={32} onChange={(e) => setNewName(e.target.value)} />
              </label>
              <label>
                {t('uploadMap')}
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
            {err && (
              <p className="error" role="alert" data-testid="scene-add-error">
                {err}
              </p>
            )}
            {scene && <SceneEditor key={scene.id + JSON.stringify(scene.grid) + scene.name} scene={scene} onPreview={onPreview} />}
          </>
        )
      case 'fog':
        if (!scene) return needScene
        return (
          <>
            <button data-testid="fog-toggle" aria-pressed={fogMode} onClick={() => setFogMode(!fogMode)}>
              {t(fogMode ? 'fogOn' : 'fogOff')}
            </button>
            {fogMode && <p className="muted">{t('fogHint')}</p>}
            <ul className="chips">
              {scene.fog.map((f) => (
                <li key={f.id}>
                  {f.col},{f.row} {f.w}×{f.h}
                  <button aria-label={t('removeFog')} data-testid={`fog-remove-${f.id}`} onClick={() => send('fog.remove', { sceneId: scene.id, fogId: f.id })}>
                    ×
                  </button>
                </li>
              ))}
            </ul>
          </>
        )
      case 'npc':
        return <NpcPanel scene={scene} />
      case 'grant':
        return <GrantPanel />
      case 'notes':
        return <NotesPanel />
      case 'ruleset':
        return <RulesetEditor key={JSON.stringify(room.ruleset)} />
    }
  }

  return (
    <section className="panel dm-panel" data-testid="dm-panel">
      <h2>DM</h2>
      <ul className="dm-menu" aria-label={t('dmMenu')}>
        {SECTIONS.map(([id, label]) => (
          <li key={id}>
            <button
              className="dm-menu-head"
              data-testid={`dm-menu-${id}`}
              aria-expanded={open === id}
              aria-controls={`dm-sec-${id}`}
              onClick={() => choose(open === id ? '' : id)}
            >
              {t(label)}
            </button>
            <div id={`dm-sec-${id}`} className="dm-sec" hidden={open !== id}>
              {open === id && content(id)}
            </div>
          </li>
        ))}
      </ul>
    </section>
  )
}

function SceneEditor({ scene, onPreview }: { scene: Scene; onPreview: (g: Grid | null) => void }) {
  const t = useT()
  const [name, setName] = useState(scene.name)
  const [g, setG] = useState(scene.grid)
  const img = useImg(scene.imageUrl)
  const w = img?.width || g.cols * g.cellSize
  const h = img?.height || g.rows * g.cellSize
  useEffect(() => () => onPreview(null), [onPreview]) // drop the preview when this editor goes away
  const step = (d: number) => {
    const next = gridFor(g.cols + d, w, h)
    setG(next)
    onPreview(next)
  }
  return (
    <div className="scene-editor">
      <label>
        {t('sceneName')}
        <input data-testid="scene-edit-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
      </label>
      <div className="grid-cols" role="group" aria-label={t('grid')}>
        <span className="grid-label">{t('grid')}</span>
        <button aria-label={t('fewerCols')} data-testid="grid-cols-dec" disabled={!img || !colsOk(g.cols - 1, w, h)} onClick={() => step(-1)}>
          −
        </button>
        <output data-testid="grid-cols" aria-live="polite">
          {g.cols}
        </output>
        <button aria-label={t('moreCols')} data-testid="grid-cols-inc" disabled={!img || !colsOk(g.cols + 1, w, h)} onClick={() => step(1)}>
          +
        </button>
        <span className="muted">
          {t('colsX')} <output data-testid="grid-rows">{g.rows}</output> {t('rows')}
        </span>
      </div>
      <button
        data-testid="scene-save"
        disabled={!validGrid(g) || !name.trim()}
        onClick={() => send('scene.update', { sceneId: scene.id, name: name.trim(), grid: g })}
      >
        {t('saveScene')}
      </button>
    </div>
  )
}
