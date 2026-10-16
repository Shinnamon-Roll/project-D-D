import type { ReactNode } from 'react'

// Small markdown subset for DM story notes, rendered as React nodes (never innerHTML, so note text can't inject markup):
// # headings, paragraphs, - / 1. lists, > quotes, --- rules, ``` code blocks, **bold**, *italic* / _italic_, `code`.
export default function Markdown({ text }: { text: string }) {
  return <div className="md">{blocks(text)}</div>
}

const BLOCK_START = /^(#{1,4}\s|```|>|\s*[-*]\s|\s*\d+[.)]\s|(-{3,}|\*{3,})\s*$)/

function blocks(text: string): ReactNode[] {
  const lines = text.replace(/\r\n?/g, '\n').split('\n')
  const out: ReactNode[] = []
  let i = 0
  while (i < lines.length) {
    const line = lines[i]
    const key = out.length
    if (!line.trim()) {
      i++
    } else if (line.startsWith('```')) {
      const code: string[] = []
      for (i++; i < lines.length && !lines[i].startsWith('```'); i++) code.push(lines[i])
      i++ // closing fence (or end of text)
      out.push(<pre key={key}>{code.join('\n')}</pre>)
    } else if (/^#{1,4}\s/.test(line)) {
      const level = line.match(/^#+/)![0].length
      const H = (['h3', 'h4', 'h5', 'h6'] as const)[level - 1] // notes sit under the panel's h2
      out.push(<H key={key}>{inline(line.slice(level).trim())}</H>)
      i++
    } else if (/^(-{3,}|\*{3,})\s*$/.test(line)) {
      out.push(<hr key={key} />)
      i++
    } else if (line.startsWith('>')) {
      const q: string[] = []
      for (; i < lines.length && lines[i].startsWith('>'); i++) q.push(lines[i].replace(/^>\s?/, ''))
      out.push(<blockquote key={key}>{blocks(q.join('\n'))}</blockquote>)
    } else if (/^\s*[-*]\s/.test(line) || /^\s*\d+[.)]\s/.test(line)) {
      const ordered = /^\s*\d/.test(line)
      const item = ordered ? /^\s*\d+[.)]\s+/ : /^\s*[-*]\s+/
      const items: string[] = []
      for (; i < lines.length && item.test(lines[i]); i++) items.push(lines[i].replace(item, ''))
      const L = ordered ? 'ol' : 'ul'
      out.push(
        <L key={key}>
          {items.map((it, j) => (
            <li key={j}>{inline(it)}</li>
          ))}
        </L>,
      )
    } else {
      const p = [line] // the first line is never a block start here, so the loop always advances
      for (i++; i < lines.length && lines[i].trim() && !BLOCK_START.test(lines[i]); i++) p.push(lines[i])
      out.push(
        <p key={key}>
          {p.map((l, j) => (
            <span key={j}>
              {j > 0 && <br />}
              {inline(l)}
            </span>
          ))}
        </p>,
      )
    }
  }
  return out
}

const INLINE = /(`[^`]+`|\*\*[^*]+\*\*|\*[^*\s][^*]*\*|_[^_\s][^_]*_)/g

function inline(s: string): ReactNode[] {
  return s.split(INLINE).map((part, i) => {
    if (i % 2 === 0) return part
    if (part.startsWith('`')) return <code key={i}>{part.slice(1, -1)}</code>
    if (part.startsWith('**')) return <strong key={i}>{part.slice(2, -2)}</strong>
    return <em key={i}>{part.slice(1, -1)}</em>
  })
}
