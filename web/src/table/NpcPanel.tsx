import { useEffect, useRef, useState } from 'react'
import type { Npc, Scene } from '../protocol'
import { useStore } from '../store'
import { send, upload } from '../ws'
import { useT } from '../i18n'
import { EyeIcon, EyeOffIcon } from './icons'
import silhouette from '../assets/npc-silhouette.svg'

const MAX_NPCS = 100

// DM/co-DM NPC pool. The server sends `npcs` only to DM seats; players never get the key.
export default function NpcPanel({ scene }: { scene?: Scene }) {
  const t = useT()
  const npcs = useStore((s) => s.npcs)
  const [editing, setEditing] = useState<Npc | 'new' | null>(null)
  const backTo = useRef('') // testid to focus once the form closes
  const list = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (editing || !backTo.current) return
    list.current?.querySelector<HTMLElement>(`[data-testid="${CSS.escape(backTo.current)}"]`)?.focus()
    backTo.current = ''
  }, [editing])

  const close = (focusId: string) => {
    backTo.current = focusId
    setEditing(null)
  }

  if (editing)
    return (
      <NpcForm
        key={editing === 'new' ? 'new' : editing.id}
        npc={editing === 'new' ? undefined : editing}
        onDone={() => close(editing === 'new' ? 'npc-add' : `npc-edit-${editing.id}`)}
      />
    )

  return (
    <div className="npc-pool" data-testid="npc-pool" ref={list}>
      <button className="primary" data-testid="npc-add" disabled={npcs.length >= MAX_NPCS} onClick={() => setEditing('new')}>
        {t('npcAdd')}
      </button>
      {npcs.length === 0 ? (
        <p className="muted">{t('npcNone')}</p>
      ) : (
        <ul className="npc-list">
          {npcs.map((n) => (
            <li key={n.id} data-testid={`npc-item-${n.id}`}>
              <img className="npc-thumb" src={n.imageUrl || silhouette} alt="" />
              <span className="npc-name">{n.name}</span>
              <span className="npc-meta">
                HP {n.hpMax}
                <button
                  className="eye-toggle"
                  data-testid={`npc-hidden-${n.id}`}
                  aria-pressed={n.hidden}
                  aria-label={t(n.hidden ? 'npcIsHidden' : 'npcIsShown', { name: n.name })}
                  title={t(n.hidden ? 'npcIsHidden' : 'npcIsShown', { name: n.name })}
                  onClick={() => send('npc.update', { id: n.id, name: n.name, imageUrl: n.imageUrl, hpMax: n.hpMax, hidden: !n.hidden })}
                >
                  {n.hidden ? <EyeOffIcon /> : <EyeIcon />}
                </button>
              </span>
              <span className="npc-actions">
                <button
                  data-testid={`npc-place-${n.id}`}
                  disabled={!scene}
                  title={scene ? undefined : t('uploadMapHint')}
                  aria-label={t('npcPlaceNamed', { name: n.name })}
                  onClick={() => scene && send('npc.place', { npcId: n.id, sceneId: scene.id })}
                >
                  {t('npcPlace')}
                </button>
                <button data-testid={`npc-edit-${n.id}`} aria-label={t('npcEditNamed', { name: n.name })} onClick={() => setEditing(n)}>
                  {t('npcEdit')}
                </button>
                <button
                  className="danger"
                  data-testid={`npc-delete-${n.id}`}
                  aria-label={t('removeNamed', { name: n.name })}
                  onClick={() => confirm(t('npcDeleteConfirm', { name: n.name })) && send('npc.remove', { id: n.id })}
                >
                  {t('remove')}
                </button>
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

function NpcForm({ npc, onDone }: { npc?: Npc; onDone: () => void }) {
  const t = useT()
  const npcs = useStore((s) => s.npcs)
  const [name, setName] = useState(npc?.name ?? '')
  const [imageUrl, setImageUrl] = useState(npc?.imageUrl ?? '')
  const [hp, setHp] = useState(String(npc?.hpMax ?? 10))
  const [hidden, setHidden] = useState(npc?.hidden ?? true)
  const [uploading, setUploading] = useState(false)

  const hpMax = Number(hp)
  // UX only, the server checks too: same name as another pool NPC (trimmed, case-insensitive). Editing may keep its own name.
  const taken = npcs.some((n) => n.id !== npc?.id && n.name.trim().toLowerCase() === name.trim().toLowerCase())
  const nameOk = name.trim().length >= 1 && name.trim().length <= 40 && !taken
  const hpOk = Number.isInteger(hpMax) && hpMax >= 1 && hpMax <= 999

  const save = () => {
    const body = { name: name.trim(), imageUrl, hpMax, hidden }
    if (npc ? send('npc.update', { id: npc.id, ...body }) : send('npc.add', body)) onDone()
  }

  return (
    <form
      className="npc-form"
      data-testid="npc-form"
      onSubmit={(e) => {
        e.preventDefault()
        if (nameOk && hpOk && !uploading) save()
      }}
      onKeyDown={(e) => e.key === 'Escape' && onDone()}
    >
      <h3>{t(npc ? 'npcEditTitle' : 'npcNew')}</h3>
      <div className="npc-portrait">
        <img className="npc-thumb big" src={imageUrl || silhouette} alt={t('npcPreviewAlt')} data-testid="npc-form-preview" />
        <label>
          {t('image')}
          <input
            data-testid="npc-form-image"
            type="file"
            accept="image/png,image/jpeg,image/webp,image/gif"
            onChange={async (e) => {
              const f = e.target.files?.[0]
              e.target.value = ''
              if (!f) return
              setUploading(true)
              const url = await upload(f)
              setUploading(false)
              if (url) setImageUrl(url)
            }}
          />
          <small className="muted">{t('npcNoImage')}</small>
        </label>
      </div>
      <label>
        {t('npcName')}
        <input
          autoFocus
          data-testid="npc-form-name"
          maxLength={40}
          value={name}
          aria-invalid={!nameOk}
          aria-describedby={taken ? 'npc-name-taken' : undefined}
          onChange={(e) => setName(e.target.value)}
        />
      </label>
      {taken && (
        <small className="name-taken" id="npc-name-taken" role="alert" data-testid="npc-form-name-taken">
          {t('npcNameTaken')}
        </small>
      )}
      <label>
        {t('npcHp')}
        <input data-testid="npc-form-hp" type="number" min={1} max={999} step={1} value={hp} aria-invalid={!hpOk} onChange={(e) => setHp(e.target.value)} />
      </label>
      <label className="switch">
        <input type="checkbox" role="switch" data-testid="npc-form-hidden" checked={hidden} onChange={(e) => setHidden(e.target.checked)} />
        <span className="switch-track" aria-hidden="true" />
        {t('npcHidden')}
      </label>
      <div className="row">
        <button className="primary" type="submit" data-testid="npc-save" disabled={!nameOk || !hpOk || uploading}>
          {t('npcSave')}
        </button>
        <button type="button" data-testid="npc-cancel" onClick={onDone}>
          {t('cancel')}
        </button>
      </div>
    </form>
  )
}
