import test from 'node:test'
import assert from 'node:assert/strict'
import {
  mergeImportedModels,
  normalizeAIConfig,
} from '../dist-tests-ai-config-state/src/views/admin/system/config/components/ai-config/state.js'

test('normalizes empty AI config values into a stable shape', () => {
  assert.deepEqual(
    normalizeAIConfig(),
    {
      default_provider: '',
      providers: [],
    },
  )
})

test('mergeImportedModels appends only new model ids', () => {
  const merged = mergeImportedModels(
    [
      { id: 'qwen-plus', name: '千问 Plus', description: '本地说明' },
    ],
    [
      { id: 'qwen-plus', owned_by: 'provider-a' },
      { id: 'deepseek-v3', owned_by: 'provider-b' },
    ],
  )

  assert.equal(merged.importedCount, 1)
  assert.equal(merged.skippedCount, 1)
  assert.deepEqual(
    merged.models.map(item => item.id),
    ['qwen-plus', 'deepseek-v3'],
  )
})

test('mergeImportedModels preserves local name and description for existing models', () => {
  const merged = mergeImportedModels(
    [
      { id: 'gpt-4o', name: '自定义显示名', description: '手工备注' },
    ],
    [
      { id: 'gpt-4o', owned_by: 'openai' },
      { id: 'gpt-4.1' },
    ],
  )

  assert.deepEqual(merged.models[0], {
    id: 'gpt-4o',
    name: '自定义显示名',
    description: '手工备注',
  })
  assert.deepEqual(merged.models[1], {
    id: 'gpt-4.1',
    name: 'gpt-4.1',
    description: '',
  })
})
