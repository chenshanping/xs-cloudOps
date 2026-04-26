export function withAlpha(hex: string, alpha: number) {
  const normalized = hex.replace('#', '')
  const safeHex = normalized.length === 3
    ? normalized
        .split('')
        .map((char) => `${char}${char}`)
        .join('')
    : normalized

  if (safeHex.length !== 6) {
    return hex
  }

  const r = Number.parseInt(safeHex.slice(0, 2), 16)
  const g = Number.parseInt(safeHex.slice(2, 4), 16)
  const b = Number.parseInt(safeHex.slice(4, 6), 16)

  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}
