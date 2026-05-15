export type ScheduleMode = 'simple' | 'advanced'
export type SimpleScheduleType = 'daily' | 'weekly' | 'monthly'

export interface SimpleScheduleState {
  type: SimpleScheduleType
  time: string
  weekday: number
  monthDay: number
}

const weekdayTextMap: Record<number, string> = {
  0: '周日',
  1: '周一',
  2: '周二',
  3: '周三',
  4: '周四',
  5: '周五',
  6: '周六',
}

export const defaultSimpleSchedule = (): SimpleScheduleState => ({
  type: 'daily',
  time: '02:00',
  weekday: 1,
  monthDay: 1,
})

export function buildCronExprFromSimple(schedule: SimpleScheduleState): string {
  const { hour, minute } = parseTimeParts(schedule.time)
  if (schedule.type === 'weekly') {
    return `${minute} ${hour} * * ${normalizeWeekday(schedule.weekday)}`
  }
  if (schedule.type === 'monthly') {
    return `${minute} ${hour} ${normalizeMonthDay(schedule.monthDay)} * *`
  }
  return `${minute} ${hour} * * *`
}

export function describeSimpleSchedule(schedule: SimpleScheduleState): string {
  if (schedule.type === 'weekly') {
    return `每${weekdayTextMap[normalizeWeekday(schedule.weekday)]} ${schedule.time} 执行`
  }
  if (schedule.type === 'monthly') {
    return `每月 ${normalizeMonthDay(schedule.monthDay)} 日 ${schedule.time} 执行`
  }
  return `每天 ${schedule.time} 执行`
}

export function parseSimpleCronExpr(expr?: string | null): SimpleScheduleState | null {
  const fields = String(expr || '').trim().split(/\s+/)
  if (fields.length !== 5) {
    return null
  }

  const [minuteText, hourText, domText, monthText, dowText] = fields
  if (!isSimpleNumber(minuteText, 0, 59) || !isSimpleNumber(hourText, 0, 23) || monthText !== '*') {
    return null
  }

  const time = `${hourText.padStart(2, '0')}:${minuteText.padStart(2, '0')}`
  if (domText === '*' && dowText === '*') {
    return { type: 'daily', time, weekday: 1, monthDay: 1 }
  }
  if (domText === '*' && isSimpleNumber(dowText, 0, 7)) {
    return {
      type: 'weekly',
      time,
      weekday: normalizeWeekday(Number(dowText)),
      monthDay: 1,
    }
  }
  if (dowText === '*' && isSimpleNumber(domText, 1, 31)) {
    return {
      type: 'monthly',
      time,
      weekday: 1,
      monthDay: Number(domText),
    }
  }
  return null
}

export function describeCronExpr(expr?: string | null): string {
  const parsed = parseSimpleCronExpr(expr)
  if (parsed) {
    return describeSimpleSchedule(parsed)
  }
  return '当前表达式较复杂，请在高级模式维护'
}

function parseTimeParts(time: string) {
  const match = /^(\d{1,2}):(\d{2})$/.exec(time || '')
  if (!match) {
    return { hour: 2, minute: 0 }
  }
  const hour = clampNumber(Number(match[1]), 0, 23)
  const minute = clampNumber(Number(match[2]), 0, 59)
  return { hour, minute }
}

function isSimpleNumber(value: string, min: number, max: number) {
  if (!/^\d+$/.test(value)) {
    return false
  }
  const numeric = Number(value)
  return numeric >= min && numeric <= max
}

function normalizeWeekday(value: number) {
  return value === 7 ? 0 : clampNumber(value, 0, 6)
}

function normalizeMonthDay(value: number) {
  return clampNumber(value, 1, 31)
}

function clampNumber(value: number, min: number, max: number) {
  if (Number.isNaN(value)) {
    return min
  }
  return Math.min(Math.max(value, min), max)
}
