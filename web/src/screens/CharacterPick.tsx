import { useRef, useState, type KeyboardEvent } from 'react'
import { MotionConfig, motion, type PanInfo } from 'motion/react'
import type { CharStats, RosterEntry } from '../protocol'
import { useMe, useStore } from '../store'
import { send } from '../ws'
import { useL, useT, type TKey } from '../i18n'
import { NotesButton } from '../table/NotesPanel'
import { RosterButton } from '../table/RosterEditor'
import { CustomSkills, Medallion, SkillList } from '../table/Skills'
import { ART_CREDIT, ART_LICENSE_URL } from '../classArt'

const STEP = 292 // px between card centres; neighbours peek in from the sides
const SPRING = { type: 'spring', stiffness: 260, damping: 30, mass: 0.9 } as const
const STATS: [keyof CharStats, TKey][] = [
  ['hpMax', 'statHp'],
  ['str', 'statStr'],
  ['end', 'statEnd'],
  ['agi', 'statAgi'],
  ['mpMax', 'statMp'],
]

// creating: players pick one of the roster cards, then make their own skills; the DM watches who took what.
export default function CharacterPick() {
  const t = useT()
  const tx = useL()
  const room = useStore((s) => s.room)!
  const { seat, character, isDm } = useMe()
  const [skillsDone, setSkillsDone] = useState(false)
  // DM gets every set (for the Ruleset editor); the pick screen shows only the sets in play.
  const roster = room.roster.filter((r) => room.activeSets.includes(r.set))

  // classId -> player name, from the live snapshot
  const takenBy = new Map<string, string>()
  for (const c of room.characters) takenBy.set(c.classId, room.seats.find((s) => s.id === c.seatId)?.name ?? '?')

  if (isDm) {
    return (
      <div className="pick">
        <h1 className="pick-title">{t('dmPickTitle')}</h1>
        <ul className="seats pick-status panel">
          {room.seats
            .filter((s) => s.kind === 'player')
            .map((s) => {
              const c = room.characters.find((x) => x.seatId === s.id)
              return (
                <li key={s.id} data-testid={`creating-status-${s.id}`}>
                  {c ? <Medallion classId={c.classId} art={c.art} size="sm" /> : <span className="medallion sm empty" aria-hidden="true" />}
                  <span className="who">{s.name}</span>
                  <span className={c ? 'ok' : 'muted'}>{c ? tx(c.name) : t('charCreating')}</span>
                </li>
              )
            })}
          <li className="pick-status-tools">
            <NotesButton />
            <RosterButton />
          </li>
        </ul>
        <Carousel roster={roster} takenBy={takenBy} />
        <ArtCredit />
      </div>
    )
  }
  if (!seat) return <p data-testid="waiting">{t('waitingOthers')}</p>
  if (character && !skillsDone) {
    return (
      <section className="panel narrow my-skills">
        <div className="my-hero">
          <Medallion classId={character.classId} art={character.art} size="lg" />
          <div>
            <h1>{tx(character.name)}</h1>
            <p className="muted">{tx(room.roster.find((r) => r.id === character.classId)?.className)}</p>
          </div>
        </div>
        <h2>{t('mySkills')}</h2>
        <p className="muted">{t('mySkillsLead')}</p>
        <CustomSkills c={character} />
        <button className="primary" data-testid="skills-done" onClick={() => setSkillsDone(true)}>
          {t('skillsDone')}
        </button>
      </section>
    )
  }
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

function ArtCredit() {
  const t = useT()
  return (
    <p className="art-credit" data-testid="art-credit">
      {t('artCredit')} {ART_CREDIT}{' '}
      <a href={ART_LICENSE_URL} target="_blank" rel="noreferrer">
        {ART_LICENSE_URL.replace('https://', '')}
      </a>
    </p>
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
  // Bar scale: each stat against the highest value in the roster, so cards compare at a glance.
  const top = Object.fromEntries(STATS.map(([k]) => [k, Math.max(1, ...roster.map((r) => r.stats[k]))])) as Record<keyof CharStats, number>

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
                  <HeroCard r={r} top={top} owner={owner} canPick={canPick} setName={setName(r.set)} />
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

function HeroCard(props: { r: RosterEntry; top: Record<keyof CharStats, number>; owner?: string; canPick?: boolean; setName: string }) {
  const { r, top, owner, canPick, setName } = props
  const t = useT()
  const tx = useL()
  return (
    <>
      <header className="hero-head">
        <Medallion classId={r.id} art={r.art} size="lg" />
        <div>
          <h2>{tx(r.heroName)}</h2>
          <p className="hero-class">{tx(r.className)}</p>
          {setName && <p className="set-badge">{setName}</p>}
        </div>
      </header>
      <p className="hero-blurb">{tx(r.blurb)}</p>
      <dl className="hero-stats">
        {STATS.map(([k, label]) => (
          <div key={k} className={`stat-${k}`}>
            <dt>{t(label)}</dt>
            <dd>
              <span className="bar" aria-hidden="true">
                <span style={{ width: `${(r.stats[k] / top[k]) * 100}%` }} />
              </span>
              <span className="num">{r.stats[k]}</span>
            </dd>
          </div>
        ))}
      </dl>
      <h3>{t('innateSkills')}</h3>
      <SkillList skills={r.skills} />
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
          <button className="primary" data-testid={`pick-${r.id}`} disabled={!!owner} onClick={() => send('character.pick', { classId: r.id })}>
            {t('pickThis')}
          </button>
        )}
      </footer>
    </>
  )
}
