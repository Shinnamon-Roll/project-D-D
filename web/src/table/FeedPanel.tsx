import { useStore } from '../store'
import { useL, useT } from '../i18n'
import { hhmm, vitalClass } from './Combat'

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
            <li key={feed.length - 1 - i} className={vitalClass(e.text)}>
              <time dateTime={new Date(e.at).toISOString()}>{hhmm(e.at)}</time> {tx(e.text)}
            </li>
          ))}
      </ol>
    </section>
  )
}
