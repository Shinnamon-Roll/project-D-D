import type { FeedEntry, LogEntry } from '../protocol'
import { useMe, useStore } from '../store'
import { useL, useT } from '../i18n'
import { hhmm, vitalClass } from './Combat'

// Latest actions (skill use, rest, items…), newest first. Polite live region: new lines are read out.
// The DM gets the full DM log here instead (every action, with who did it, incl. DM-only lines).
export default function FeedPanel() {
  const t = useT()
  const tx = useL()
  const { isDm } = useMe()
  const pub = useStore((s) => s.room!.feed)
  const dmLog = useStore((s) => s.dmLog)
  const feed: (FeedEntry & Partial<Pick<LogEntry, 'actor'>>)[] = isDm ? dmLog : pub
  return (
    <section className="panel feed-panel">
      <h2>{t('feed')}</h2>
      <ol className="feed" data-testid="feed" aria-live="polite">
        {feed.length === 0 && <li className="muted">{t('feedEmpty')}</li>}
        {feed
          .slice()
          .reverse()
          .map((e, i) => (
            // position from the oldest end stays put as new lines arrive on top
            <li key={feed.length - 1 - i} className={vitalClass(e.text)}>
              <time dateTime={new Date(e.at).toISOString()}>{hhmm(e.at)}</time>
              {e.actor && e.actor !== '-' && <span className="log-actor">{e.actor}</span>} {tx(e.text)}
            </li>
          ))}
      </ol>
    </section>
  )
}
