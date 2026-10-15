import { useState } from 'react'
import { useStore } from '../store'
import { send } from '../ws'

export default function CampaignMenu() {
  const campaigns = useStore((s) => s.room?.campaigns ?? [])
  const [name, setName] = useState('')
  const trimmed = name.trim()

  return (
    <section className="panel narrow">
      <h1>เลือกแคมเปญ / Choose a campaign</h1>
      <p>คนที่เลือกแคมเปญจะเป็น DM / Whoever picks a campaign becomes the DM.</p>
      <label>
        ชื่อ DM / DM name
        <input data-testid="dm-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
      </label>
      <div className="stack">
        {campaigns.map((c) =>
          c.status === 'available' ? (
            <button
              key={c.id}
              className="big"
              data-testid={`campaign-${c.id}`}
              disabled={!trimmed}
              onClick={() => send('campaign.select', { campaignId: c.id, name: trimmed })}
            >
              {c.name}
            </button>
          ) : (
            <button key={c.id} className="big" data-testid={`campaign-${c.id}`} disabled>
              {c.name} — Coming soon
            </button>
          ),
        )}
      </div>
    </section>
  )
}
