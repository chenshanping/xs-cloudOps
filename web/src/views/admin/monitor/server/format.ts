export function formatBytes(value?: number) {
  const size = Number(value || 0)
  if (size <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const index = Math.min(Math.floor(Math.log(size) / Math.log(1024)), units.length - 1)
  return `${(size / Math.pow(1024, index)).toFixed(index === 0 ? 0 : 2)} ${units[index]}`
}

export function formatDuration(seconds?: number) {
  const total = Math.max(0, Math.floor(Number(seconds || 0)))
  const days = Math.floor(total / 86400)
  const hours = Math.floor((total % 86400) / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  if (days > 0) return `${days} 天 ${hours} 小时`
  if (hours > 0) return `${hours} 小时 ${minutes} 分钟`
  if (minutes > 0) return `${minutes} 分钟`
  return `${total} 秒`
}

export function formatPercent(value?: number) {
  const percent = Number(value || 0)
  return Math.max(0, Math.min(100, percent))
}

export function formatDateTime(value?: string) {
  if (!value) return '-'
  return value.replace('T', ' ').replace(/\.\d+Z?$/, '').replace(/Z$/, '')
}

export function healthColor(reachable?: boolean) {
  return reachable ? 'success' : 'error'
}

export function progressStatus(percent?: number) {
  const value = Number(percent || 0)
  if (value >= 90) return 'exception'
  if (value >= 75) return 'active'
  return 'normal'
}
