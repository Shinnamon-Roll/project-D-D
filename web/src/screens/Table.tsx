import { useState } from 'react'
import type { Grid, Token } from '../protocol'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import MapStage from '../table/MapStage'
import DmPanel from '../table/DmPanel'
import DicePanel from '../table/DicePanel'
import SheetPanel from '../table/SheetPanel'

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

  return (
    <div className="table">
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
            <ul className="token-list" data-testid="token-list">
              {tokens.map((t) => (
                <TokenRow key={t.id + t.col + ',' + t.row} t={t} movable={canMove(t)} removable={isDm && t.characterId === ''} />
              ))}
            </ul>
          </>
        ) : (
          <p data-testid="no-scene">{isDm ? 'อัปแมพเพื่อสร้างฉาก / Upload a map to create a scene' : 'รอ DM เลือกฉาก / Waiting for the DM to pick a scene'}</p>
        )}
      </div>
      <aside className="side">
        {character && <SheetPanel character={character} />}
        <DicePanel />
        {isDm && (
          <DmPanel scene={scene} onView={setViewId} onPreview={setPreview} fogMode={fogMode} setFogMode={setFogMode} />
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
