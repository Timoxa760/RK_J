export interface FnsQrFields {
  fn: string
  fd: string
  fp: string
}

/** Парсит QR-строку ФНС: fn=…&i=…&fp=… (или fd=). */
export function parseFnsQr(qr: string): FnsQrFields | null {
  const trimmed = qr.trim()
  if (!trimmed) return null

  const params: Record<string, string> = {}
  for (const part of trimmed.split('&')) {
    const eq = part.indexOf('=')
    if (eq <= 0) continue
    const key = part.slice(0, eq).trim()
    const value = part.slice(eq + 1).trim()
    if (key && value) params[key] = decodeURIComponent(value)
  }

  const fn = params.fn
  const fd = params.i ?? params.fd
  const fp = params.fp

  if (fn && fd && fp) {
    return { fn, fd, fp }
  }

  return null
}
