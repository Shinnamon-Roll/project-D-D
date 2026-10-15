import { useState } from 'react'
import { useStore } from '../store'
import { send } from '../ws'
import crest from '../assets/crest.svg'
import stormwreck from '../assets/stormwreck.svg'
import salem from '../assets/salem.svg'

// Original art only (never WotC covers); unknown campaigns fall back to the crest.
const art: Record<string, string> = { 'starter-set': stormwreck, salem }

export default function CampaignMenu() {
  const campaigns = useStore((s) => s.room?.campaigns ?? [])
  const [name, setName] = useState('')
  const trimmed = name.trim()

  return (
    <section className="panel narrow menu">
      <img className="crest-hero" src={crest} alt="" />
      <h1>เลือกแคมเปญ / Choose a campaign</h1>
      <p className="lead">คนที่เลือกแคมเปญจะเป็น DM / Whoever picks a campaign becomes the DM.</p>
      <label>
        ชื่อ DM / DM name
        <input data-testid="dm-name" value={name} maxLength={32} onChange={(e) => setName(e.target.value)} />
      </label>
      <div className="cards">
        {campaigns.map((c) =>
          c.status === 'available' ? (
            <button
              key={c.id}
              className="card"
              data-testid={`campaign-${c.id}`}
              disabled={!trimmed}
              onClick={() => send('campaign.select', { campaignId: c.id, name: trimmed })}
            >
              <img src={art[c.id] ?? crest} alt="" />
              <span className="title">{c.name}</span>
              <span className="cta">{trimmed ? 'เริ่มการผจญภัย / Begin' : 'ใส่ชื่อ DM ก่อน / Enter DM name first'}</span>
            </button>
          ) : (
            <button key={c.id} className="card locked" data-testid={`campaign-${c.id}`} disabled>
              <img src={art[c.id] ?? crest} alt="" />
              <span className="title">{c.name}</span>
              <span className="cta">ยังไม่เปิด / Locked</span>
              <span className="wax">Coming soon</span>
            </button>
          ),
        )}
      </div>
    </section>
  )
}
