export type ConfigTabKey = 'basic' | 'file' | 'login' | 'email' | 'ai'

export type ConfigDirtyState = Partial<Record<ConfigTabKey, boolean>>

export interface ConfigTabGuardHandle {
  isDirty: () => boolean
  save: () => Promise<boolean>
  discardChanges: () => void
  closeTransientUi: () => void
}

export function createSnapshot(value: unknown): string {
  const snapshot = JSON.stringify(value === undefined ? null : value)
  return snapshot ?? 'null'
}

export function cloneFromSnapshot<T>(snapshot: string): T {
  return JSON.parse(snapshot) as T
}

export function isSnapshotDirty(snapshot: string, value: unknown): boolean {
  return snapshot !== createSnapshot(value)
}

export function hasDirtyConfigTabs(state: ConfigDirtyState): boolean {
  return Object.values(state).some(Boolean)
}
