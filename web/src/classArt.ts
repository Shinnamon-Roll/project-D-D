// Class icons (game-icons.net, CC BY 3.0 — credit shown on the pick screen) recoloured gold as data URLs,
// so the same file works in <img> and as a Konva image on the map.
const raw = import.meta.glob('./assets/classes/*.svg', { query: '?raw', import: 'default', eager: true }) as Record<string, string>

const GOLD = '#c9a227' // --gold; Konva can't read CSS vars
const icons: Record<string, string> = {}
for (const [path, svg] of Object.entries(raw)) {
  const id = path.slice(path.lastIndexOf('/') + 1, -4)
  // width/height: canvas needs an intrinsic size to rasterise an SVG image
  icons[id] = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(svg.replace('<svg ', `<svg width="512" height="512" fill="${GOLD}" `))
}

// art "" = the class icon; else an uploaded image URL. icon=false means a photo-like upload (cover, no inset).
export function classArt(classId: string, art = ''): { src: string; icon: boolean } {
  if (art) return { src: art, icon: false }
  return { src: icons[classId] ?? '', icon: true }
}

export const ART_CREDIT =
  'Class icons by Lorc, Delapouite, DarkZaitzev, Cathelineau and Skoll — game-icons.net — CC BY 3.0. Recolored.'
export const ART_LICENSE_URL = 'https://creativecommons.org/licenses/by/3.0/'
