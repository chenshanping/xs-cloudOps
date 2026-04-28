export interface AIPreferences {
  enableSearch: boolean
  enableThinking: boolean
}

export const AI_PREFERENCES_STORAGE_KEY = 'go-base-ai-preferences'

export const DEFAULT_AI_PREFERENCES: AIPreferences = {
  enableSearch: false,
  enableThinking: true,
}

interface StorageLike {
  getItem(key: string): string | null
  setItem(key: string, value: string): void
}

function resolveStorage(storage?: StorageLike | null) {
  if (storage) {
    return storage
  }

  if (typeof globalThis === 'undefined' || !('localStorage' in globalThis)) {
    return null
  }

  return globalThis.localStorage
}

function normalizeAIPreferences(input?: Partial<AIPreferences> | null): AIPreferences {
  return {
    enableSearch:
      typeof input?.enableSearch === 'boolean'
        ? input.enableSearch
        : DEFAULT_AI_PREFERENCES.enableSearch,
    enableThinking:
      typeof input?.enableThinking === 'boolean'
        ? input.enableThinking
        : DEFAULT_AI_PREFERENCES.enableThinking,
  }
}

export function loadAIPreferences(storage?: StorageLike | null): AIPreferences {
  const targetStorage = resolveStorage(storage)
  if (!targetStorage) {
    return { ...DEFAULT_AI_PREFERENCES }
  }

  try {
    const raw = targetStorage.getItem(AI_PREFERENCES_STORAGE_KEY)
    if (!raw) {
      return { ...DEFAULT_AI_PREFERENCES }
    }

    return normalizeAIPreferences(JSON.parse(raw) as Partial<AIPreferences>)
  } catch (error) {
    console.error('读取 AI 对话偏好失败', error)
    return { ...DEFAULT_AI_PREFERENCES }
  }
}

export function persistAIPreferences(
  storage: StorageLike | null | undefined,
  preferences: AIPreferences,
) {
  const targetStorage = resolveStorage(storage)
  if (!targetStorage) {
    return
  }

  try {
    targetStorage.setItem(
      AI_PREFERENCES_STORAGE_KEY,
      JSON.stringify(normalizeAIPreferences(preferences)),
    )
  } catch (error) {
    console.error('保存 AI 对话偏好失败', error)
  }
}
