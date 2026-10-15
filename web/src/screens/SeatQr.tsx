import { useEffect, useRef, useState, type MouseEvent } from 'react'
import QRCode from 'qrcode'
import { useT } from '../i18n'

// Modal with a QR that opens /?seat=<seatId>. Native <dialog> gives focus trap, Esc and focus return.
// The dialog's own padding is also e.currentTarget; only clicks outside its box hit the backdrop.
export function closeOnBackdrop(e: MouseEvent<HTMLDialogElement>) {
  if (e.target !== e.currentTarget) return
  const r = e.currentTarget.getBoundingClientRect()
  if (e.clientX < r.left || e.clientX > r.right || e.clientY < r.top || e.clientY > r.bottom) e.currentTarget.close()
}

export default function SeatQr({ seatId, num, onClose }: { seatId: string; num: number; onClose: () => void }) {
  const t = useT()
  const ref = useRef<HTMLDialogElement>(null)
  const [urls, setUrls] = useState<string[] | null>(null)
  const [img, setImg] = useState('')
  const [pick, setPick] = useState(0)

  useEffect(() => {
    ref.current?.showModal()
    let live = true
    fetch('/api/lan')
      .then((r) => (r.ok ? r.json() : { urls: [] }))
      .catch(() => ({ urls: [] }))
      .then((j) => live && setUrls(Array.isArray(j.urls) ? j.urls : []))
    return () => {
      live = false
    }
  }, [])

  const url = urls ? `${urls[pick] ?? location.origin}/?seat=${encodeURIComponent(seatId)}` : ''

  useEffect(() => {
    if (!url) return
    let live = true
    QRCode.toDataURL(url, { width: 512, margin: 2, color: { dark: '#2b1d12', light: '#fffdf5' } })
      .then((d) => live && setImg(d))
      .catch(() => live && setImg(''))
    return () => {
      live = false
    }
  }, [url])

  return (
    <dialog
      ref={ref}
      className="panel qr-modal"
      aria-modal="true"
      aria-labelledby="seat-qr-title"
      data-testid="seat-qr-modal"
      onClose={onClose}
      onClick={closeOnBackdrop}
    >
      <h2 id="seat-qr-title" data-testid="seat-qr-label">
        {t('seatN', { n: num })}
      </h2>
      {img ? (
        <img className="qr" src={img} alt={t('qrAlt', { n: num })} data-testid="seat-qr-code" />
      ) : (
        <p className="qr qr-wait">{t('generating')}</p>
      )}
      <p className="qr-hint">{t('qrHint')}</p>
      {urls && urls.length > 1 && (
        <div className="row qr-addrs" role="group" aria-label={t('lanAddress')}>
          {urls.map((u, i) => (
            <button key={u} data-testid={`seat-qr-addr-${i}`} aria-pressed={i === pick} onClick={() => setPick(i)}>
              {u.replace(/^https?:\/\//, '')}
            </button>
          ))}
        </div>
      )}
      {url && (
        <p className="qr-url" data-testid="seat-qr-url">
          {url}
        </p>
      )}
      {urls?.length === 0 && (
        <p className="qr-warn" role="note" data-testid="seat-qr-warn">
          {t('qrWarn')} <code>-addr 0.0.0.0:8080</code>
        </p>
      )}
      <button className="primary" data-testid="seat-qr-close" onClick={() => ref.current?.close()}>
        {t('close')}
      </button>
    </dialog>
  )
}
