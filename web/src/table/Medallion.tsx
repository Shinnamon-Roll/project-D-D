import { classArt } from '../classArt'

// Round class medallion: gold icon on oxblood, or an uploaded picture filling the circle. heroId picks the icon.
export function Medallion({ heroId, art, size = 'md' }: { heroId: string; art?: string; size?: 'sm' | 'md' | 'lg' }) {
  const a = classArt(heroId, art)
  return (
    <span className={`medallion ${size} ${a.icon ? 'icon' : 'photo'}`} aria-hidden="true">
      {a.src && <img src={a.src} alt="" draggable={false} />}
    </span>
  )
}
