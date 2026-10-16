import { useStore } from '../store'
import { useL, useT } from '../i18n'

const hhmm = (at: number) => new Date(at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })

// Latest actions (skill use, rest, items…), newest first. Polite live region: new lines are read out.
export default function FeedPanel() {
  const t = useT()
  const tx = useL()
  const feed = useStore((s) => s.room!.feed)
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
            <li key={feed.length - 1 - i}>
              <time dateTime={new Date(e.at).toISOString()}>{hhmm(e.at)}</time> {tx(e.text)}
            </li>
          ))}
      </ol>
    </section>
  )
}
