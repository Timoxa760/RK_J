/** Дата формы (YYYY-MM-DD) → ISO 8601 для receipt/manual. */
export function toReceiptIsoDate(dateStr: string): string {
  if (dateStr.includes('T')) return dateStr
  const d = new Date(`${dateStr}T12:00:00`)
  return Number.isNaN(d.getTime()) ? new Date().toISOString() : d.toISOString()
}
