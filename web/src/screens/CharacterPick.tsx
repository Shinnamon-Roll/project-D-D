import { useRef, useState, type KeyboardEvent } from 'react'
import { MotionConfig, motion, type PanInfo } from 'motion/react'
import type { RosterEntry } from '../protocol'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import { useL, useT } from '../i18n'
import { NotesButton } from '../table/NotesPanel'
import { RosterButton } from '../table/RosterEditor'
import { Medallion } from '../table/Medallion'
import { AbilityGrid, FeatureList, SkillChips } from '../table/StatBlocks'
import { ART_CREDIT, ART_LICENSE_URL } from '../classArt'
import { featuresAt, hpLevel1, profBonus, signed, useClass } from '../rules'

const STEP = 292 // px between card centres; neighbours peek in from the sides
const SPRING = { type: 'spring', stiffness: 260, damping: 30, mass: 0.9 } as const

// creating: players pick one of the roster cards; the DM watches who took what.
export default function CharacterPick() {
  const t = useT()
  const tx = useL()
  const room = useStore((s) => s.room)!
  const { seat, character, isDm } = useMe()
  // DM gets every set (for the Ruleset editor); the pick screen shows only the sets in play.
  const roster = room.roster.filter((r) => room.activeSets.includes(r.set))

  // heroId -> player name, from the live snapshot
  const takenBy = new Map<string, string>()
  for (const c of room.characters) takenBy.set(c.heroId, room.seats.find((s) => s.id === c.seatId)?.name ?? '?')

  // DM: no carousel, one full card per player seat showing the hero they took (or that they're still choosing)
  if (isDm) {
    const setName = (id: string) => (room.activeSets.length > 1 ? tx(room.sets.find((x) => x.id === id)?.name) : '')
    return (
      <div className="pick">
        <h1 className="pick-title">{t('dmPickTitle')}</h1>
        <div className="pick-status-tools">
          <NotesButton />
          <RosterButton />
        </div>
        <ul className="dm-pick-grid">
          {room.seats
            .filter((s) => s.kind === 'player')
            .map((s) => {
              const c = room.characters.find((x) => x.seatId === s.id)
              const r = c && room.roster.find((x) => x.id === c.heroId)
              return (
                <li key={s.id} className={`hero-card${r ? '' : ' waiting'}`} data-testid={`creating-status-${s.id}`}>
                  <p className="dm-pick-who">{s.name}</p>
                  <div className="hero-inner">
                    {r ? (
                      <HeroCard r={r} setName={setName(r.set)} />
                    ) : (
                      <div className="hero-wait">
                        <span className="medallion lg empty" aria-hidden="true" />
                        <p className="muted">{t('charCreating')}</p>
                      </div>
                    )}
                  </div>
                </li>
              )
            })}
        </ul>
        <ArtCredit />
      </div>
    )
  }
  if (!seat) return <p data-testid="waiting">{t('waitingOthers')}</p>
  if (character) return <p data-testid="waiting">{t('waitingOthers')}</p>
  return (
    <div className="pick">
      <h1 className="pick-title">{t('pickTitle')}</h1>
      <p className="pick-lead">{t('pickLead', { n: roster.length })}</p>
      <Carousel roster={roster} takenBy={takenBy} canPick />
      <ArtCredit />
    </div>
  )
}

// Art credit, plus the SRD 5.2.1 CC BY 4.0 attribution once the class data (the SRD material) has loaded.
function ArtCredit() {
  const t = useT()
  const srd = useStore((s) => s.classes?.attribution)
  return (
    <>
      <p className="art-credit" data-testid="art-credit">
        {t('artCredit')} {ART_CREDIT}{' '}
        <a href={ART_LICENSE_URL} target="_blank" rel="noreferrer">
          {ART_LICENSE_URL.replace('https://', '')}
        </a>
      </p>
      {srd && (
        <p className="art-credit srd-credit" data-testid="srd-credit" lang="en">
          {srd}
        </p>
      )}
    </>
  )
}

function Carousel({ roster, takenBy, canPick }: { roster: RosterEntry[]; takenBy: Map<string, string>; canPick?: boolean }) {
  const t = useT()
  const tx = useL()
  const room = useStore((s) => s.room)!
  // set badge only when heroes from more than one set share the carousel
  const setName = (id: string) => (room.activeSets.length > 1 ? tx(room.sets.find((x) => x.id === id)?.name) : '')
  const [index, setIndex] = useState(0)
  const box = useRef<HTMLElement>(null)
  const n = roster.length
  const go = (i: number) => setIndex(Math.max(0, Math.min(n - 1, i)))

  const onKey = (e: KeyboardEvent) => {
    const next = { ArrowLeft: index - 1, ArrowRight: index + 1, Home: 0, End: n - 1 }[e.key]
    if (next === undefined || (e.target as HTMLElement).closest('dialog')) return
    e.preventDefault()
    go(next)
    // a focused button in the card that slides away goes inert; keep focus in the carousel
    if ((e.target as HTMLElement).closest('.hero-card')) box.current?.focus()
  }
  const onDragEnd = (_: unknown, info: PanInfo) => {
    const swipe = info.offset.x + info.velocity.x * 0.2
    if (swipe < -STEP / 4) go(index + 1)
    else if (swipe > STEP / 4) go(index - 1)
  }
  const cur = roster[index]

  return (
    // reducedMotion="user": transforms jump instead of sliding when the OS asks for less motion
    <MotionConfig reducedMotion="user" transition={SPRING}>
      <section ref={box} tabIndex={-1} className="carousel" data-testid="carousel" aria-roledescription="carousel" aria-label={t('pickTitle')} onKeyDown={onKey}>
        <p className="sr-only" aria-live="polite">
          {cur && t('cardOf', { n: index + 1, total: n, name: tx(cur.heroName) })}
        </p>
        <motion.div className="carousel-track" drag="x" dragSnapToOrigin dragElastic={0.6} onDragEnd={onDragEnd}>
          {roster.map((r, i) => {
            const off = i - index
            const far = Math.abs(off) > 2
            const owner = takenBy.get(r.id)
            return (
              <motion.article
                key={r.id}
                className={`hero-card${owner ? ' taken' : ''}`}
                data-testid={`card-${r.id}`}
                data-current={off === 0}
                aria-roledescription="slide"
                aria-label={t('cardOf', { n: i + 1, total: n, name: tx(r.heroName) })}
                // tapping an off-centre card brings it to the centre
                onClick={off !== 0 ? () => go(i) : undefined}
                initial={false}
                animate={{
                  x: off * STEP,
                  scale: off === 0 ? 1 : 0.85,
                  opacity: [1, 0.6, 0.25][Math.abs(off)] ?? 0,
                  rotateY: Math.max(-1, Math.min(1, off)) * -10,
                }}
                // only the centre card and its neighbours get their own compositor layer
                style={{ zIndex: n - Math.abs(off), pointerEvents: far ? 'none' : undefined, willChange: Math.abs(off) <= 1 ? 'transform' : 'auto' }}
              >
                {/* off-centre cards stay visible but out of the tab order and the a11y tree */}
                <div className="hero-inner" inert={off !== 0}>
                  <HeroCard r={r} owner={owner} canPick={canPick} setName={setName(r.set)} />
                </div>
              </motion.article>
            )
          })}
        </motion.div>
        <div className="carousel-nav">
          <button className="nav-arrow" data-testid="carousel-prev" aria-label={t('prevCard')} disabled={index === 0} onClick={() => go(index - 1)}>
            ‹
          </button>
          <div className="dots" role="group" aria-label={t('pickTitle')}>
            {roster.map((r, i) => (
              <button
                key={r.id}
                className={`dot${takenBy.has(r.id) ? ' taken' : ''}`}
                aria-label={t('goToCard', { name: tx(r.heroName) })}
                aria-current={i === index ? 'true' : undefined}
                onClick={() => go(i)}
              />
            ))}
          </div>
          <button className="nav-arrow" data-testid="carousel-next" aria-label={t('nextCard')} disabled={index === n - 1} onClick={() => go(index + 1)}>
            ›
          </button>
        </div>
      </section>
    </MotionConfig>
  )
}

// One hero at level 1: who they are, the six abilities with saves, vitals, proficient skills, level-1 features, gear.
function HeroCard(props: { r: RosterEntry; owner?: string; canPick?: boolean; setName: string }) {
  const { r, owner, canPick, setName } = props
  const t = useT()
  const tx = useL()
  const cls = useClass(r.classId)
  const pb = profBonus(1)
  const vitals: [string, string, string, string][] = [
    ['hp', t('hpL1'), t('hpL1Hint'), cls ? String(hpLevel1(cls, r.abilities.con)) : '—'],
    ['ac', t('ac'), t('acHint'), String(r.ac)],
    ['hd', t('hitDie'), t('hitDieHint'), cls ? `d${cls.hitDie}` : '—'],
    ['pb', t('profBonus'), t('profBonusHint'), signed(pb)],
  ]
  return (
    <>
      <header className="hero-head">
        <Medallion heroId={r.id} art={r.art} size="lg" />
        <div>
          <h2>{tx(r.heroName)}</h2>
          <p className="hero-class">{tx(cls?.name)}</p>
          {setName && <p className="set-badge">{setName}</p>}
        </div>
      </header>
      <p className="hero-blurb">{tx(r.blurb)}</p>
      <dl className="hero-vitals">
        {vitals.map(([k, label, hint, v]) => (
          <div key={k} title={hint}>
            <dt>{label}</dt>
            <dd data-testid={`hero-${k}`}>{v}</dd>
          </div>
        ))}
      </dl>
      <AbilityGrid abilities={r.abilities} saves={cls?.saves ?? []} pb={pb} />
      {r.skills.length > 0 && (
        <>
          <h3>{t('skillsProf')}</h3>
          <SkillChips abilities={r.abilities} skills={r.skills} expertise={r.expertise} pb={pb} />
        </>
      )}
      <h3>{t('featuresL1')}</h3>
      {cls ? <FeatureList features={featuresAt(cls, 1)} /> : <p className="hero-loading muted">{t('classesLoading')}</p>}
      {r.items.length > 0 && (
        <>
          <h3>{t('items')}</h3>
          <ul className="hero-items">
            {r.items.map((it, i) => (
              <li key={i}>
                {tx(it.name)}
                {it.qty > 1 && <span className="qty"> ×{it.qty}</span>}
              </li>
            ))}
          </ul>
        </>
      )}
      <footer className="hero-foot">
        {owner && (
          <p className="taken-tag" data-testid={`taken-${r.id}`}>
            {t('takenBy', { name: owner })}
          </p>
        )}
        {canPick && (
          <button className="primary" data-testid={`pick-${r.id}`} disabled={!!owner} onClick={() => send('character.pick', { heroId: r.id })}>
            {t('pickThis')}
          </button>
        )}
      </footer>
    </>
  )
}
