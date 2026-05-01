import test from 'node:test'
import assert from 'node:assert/strict'
import {
  filterModelsByCapabilityAndKeyword,
  groupRemoteModelsByOwner,
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
    is_thinking: false,
    support_vision: false,
    support_tools: false,
    search_strategy: 'none',
    support_embedding: false,
    support_rerank: false,
    is_free: false,
    temperature: null,
    context_window: null,
  })
  assert.deepEqual(merged.models[1], {
    id: 'gpt-4.1',
    name: 'gpt-4.1',
    description: '',
    is_thinking: false,
    support_vision: false,
    support_tools: false,
    search_strategy: 'none',
    support_embedding: false,
    support_rerank: false,
    is_free: false,
    temperature: null,
    context_window: null,
  })
})

test('filterModelsByCapabilityAndKeyword applies capability and keyword together', () => {
  const models = [
    { id: 'qwen-max', name: 'Qwen Max', description: '标准模型', support_tools: true },
    { id: 'kimi-k2.6', name: 'Kimi K2.6', description: '联网模型', search_strategy: 'builtin' },
    { id: 'qwen-vl-max', name: 'Qwen VL', description: '视觉模型', support_vision: true },
  ]

  const filtered = filterModelsByCapabilityAndKeyword(models, 'search', 'kimi')

  assert.deepEqual(
    filtered.map(item => item.id),
    ['kimi-k2.6'],
  )
})

test('groupRemoteModelsByOwner falls back to 其他 when owned_by is empty', () => {
  const groups = groupRemoteModelsByOwner([
    { id: 'qwen-max', name: 'Qwen Max', owned_by: 'qwen' },
    { id: 'kimi-k2.6', name: 'Kimi K2.6', owned_by: '' },
    { id: 'qwen-vl-max', name: 'Qwen VL', owned_by: 'qwen' },
  ])

  assert.deepEqual(
    groups.map(group => ({ name: group.name, ids: group.items.map(item => item.id) })),
    [
      { name: 'qwen', ids: ['qwen-max', 'qwen-vl-max'] },
      { name: '其他', ids: ['kimi-k2.6'] },
    ],
  )
})
