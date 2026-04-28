import test from 'node:test'
import assert from 'node:assert/strict'
import {
  DEFAULT_AI_PREFERENCES,
  loadAIPreferences,
  persistAIPreferences,
} from '../dist-tests-ai-preferences/src/store/ai-preferences.js'

function createStorage(initial = {}) {
  const store = new Map(Object.entries(initial))
  return {
    getItem(key) {
      return store.has(key) ? store.get(key) : null
    },
    setItem(key, value) {
      store.set(key, String(value))
    },
  }
}

test('loadAIPreferences falls back to defaults when storage is empty', () => {
  const prefs = loadAIPreferences(createStorage())
  assert.deepEqual(prefs, DEFAULT_AI_PREFERENCES)
})

test('loadAIPreferences restores persisted booleans', () => {
  const prefs = loadAIPreferences(
    createStorage({
      'go-base-ai-preferences': JSON.stringify({
        enableSearch: true,
        enableThinking: false,
      }),
    }),
  )

  assert.deepEqual(prefs, {
    enableSearch: true,
    enableThinking: false,
  })
})

test('loadAIPreferences ignores malformed or non-boolean values', () => {
  const originalConsoleError = console.error
  console.error = () => {}
  try {
    const malformed = loadAIPreferences(
      createStorage({
        'go-base-ai-preferences': '{bad json',
      }),
    )
    assert.deepEqual(malformed, DEFAULT_AI_PREFERENCES)
  } finally {
    console.error = originalConsoleError
  }

  const invalidShape = loadAIPreferences(
    createStorage({
      'go-base-ai-preferences': JSON.stringify({
        enableSearch: 'true',
        enableThinking: null,
      }),
    }),
  )
  assert.deepEqual(invalidShape, DEFAULT_AI_PREFERENCES)
})

test('persistAIPreferences writes both toggles as JSON', () => {
  const storage = createStorage()
  persistAIPreferences(storage, {
    enableSearch: true,
    enableThinking: false,
  })

  assert.equal(
    storage.getItem('go-base-ai-preferences'),
    JSON.stringify({
      enableSearch: true,
      enableThinking: false,
    }),
  )
})
