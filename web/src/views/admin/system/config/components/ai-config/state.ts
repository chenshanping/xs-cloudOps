export type AIModelCapabilityKey =
  | 'all'
  | 'reasoning'
  | 'vision'
  | 'search'
  | 'free'
  | 'embedding'
  | 'rerank'
  | 'tool'

export type AISearchStrategy = 'none' | 'builtin' | 'tool'

export interface AIModel {
  id: string
  name: string
  description: string
  is_thinking: boolean
  support_vision: boolean
  support_tools: boolean
  search_strategy: AISearchStrategy
  support_embedding: boolean
  support_rerank: boolean
  is_free: boolean
  temperature: number | null
  context_window: number | null
}

export interface AIProvider {
  name: string
  api_key: string
  base_url: string
  models: AIModel[]
}

export interface AIConfigState {
  default_provider: string
  providers: AIProvider[]
}

export interface RemoteProviderModel extends AIModel {
  object?: string
  created?: number
  owned_by?: string
  tags?: string[]
}

export interface ImportMergeResult {
  models: AIModel[]
  importedCount: number
  skippedCount: number
}

export interface RemoteModelGroup {
  name: string
  items: RemoteProviderModel[]
}

export const capabilityTabOptions: Array<{ label: string; value: AIModelCapabilityKey }> = [
  { label: '全部', value: 'all' },
  { label: '推理', value: 'reasoning' },
  { label: '视觉', value: 'vision' },
  { label: '联网', value: 'search' },
  { label: '免费', value: 'free' },
  { label: '嵌入', value: 'embedding' },
  { label: '重排', value: 'rerank' },
  { label: '工具', value: 'tool' },
]

export const capabilityTagMetaMap: Record<Exclude<AIModelCapabilityKey, 'all'>, { color: string; label: string }> = {
  reasoning: { color: 'gold', label: '推理' },
  vision: { color: 'purple', label: '视觉' },
  search: { color: 'blue', label: '联网' },
  free: { color: 'green', label: '免费' },
  embedding: { color: 'cyan', label: '嵌入' },
  rerank: { color: 'volcano', label: '重排' },
  tool: { color: 'geekblue', label: '工具' },
}

export const searchStrategyOptions: Array<{ label: string; value: AISearchStrategy }> = [
  { label: '不支持', value: 'none' },
  { label: '内置联网', value: 'builtin' },
  { label: '工具联网', value: 'tool' },
]

const capabilityTagAliasMap: Record<string, Exclude<AIModelCapabilityKey, 'all'>> = {
  reasoning: 'reasoning',
  inference: 'reasoning',
  chat: 'reasoning',
  llm: 'reasoning',
  text: 'reasoning',
  推理: 'reasoning',
  vision: 'vision',
  visual: 'vision',
  image: 'vision',
  multimodal: 'vision',
  multi_modal: 'vision',
  视觉: 'vision',
  search: 'search',
  websearch: 'search',
  web_search: 'search',
  online: 'search',
  联网: 'search',
  free: 'free',
  免费: 'free',
  embedding: 'embedding',
  embeddings: 'embedding',
  embed: 'embedding',
  嵌入: 'embedding',
  rerank: 'rerank',
  reranker: 'rerank',
  重排: 'rerank',
  tool: 'tool',
  tools: 'tool',
  function: 'tool',
  工具: 'tool',
}

export function createEmptyModel(): AIModel {
  return {
    id: '',
    name: '',
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
  }
}

export function createEmptyProvider(): AIProvider {
  return {
    name: '',
    api_key: '',
    base_url: '',
    models: [],
  }
}

export function normalizeSearchStrategy(value?: unknown): AISearchStrategy {
  const normalized = String(value ?? '').trim().toLowerCase()
  if (normalized === 'builtin' || normalized === 'tool') {
    return normalized
  }
  return 'none'
}

function normalizeBooleanFlag(value?: unknown): boolean {
  if (typeof value === 'boolean') {
    return value
  }
  if (typeof value === 'number') {
    return value !== 0
  }
  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase()
    return ['1', 'true', 'yes', 'y', 'on'].includes(normalized)
  }
  return false
}

function normalizeNumberValue(value?: unknown): number | null {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value
  }
  if (typeof value === 'string' && value.trim()) {
    const parsed = Number(value)
    if (Number.isFinite(parsed)) {
      return parsed
    }
  }
  return null
}

function normalizeCapabilityTags(value?: unknown): Array<Exclude<AIModelCapabilityKey, 'all'>> {
  if (!Array.isArray(value)) {
    return []
  }
  const seen = new Set<Exclude<AIModelCapabilityKey, 'all'>>()
  const normalized: Array<Exclude<AIModelCapabilityKey, 'all'>> = []
  for (const item of value) {
    if (typeof item !== 'string') {
      continue
    }
    const key = item.trim().toLowerCase().replace(/[\s-]+/g, '_')
    const mapped = capabilityTagAliasMap[key] ?? capabilityTagAliasMap[item.trim()]
    if (!mapped || seen.has(mapped)) {
      continue
    }
    seen.add(mapped)
    normalized.push(mapped)
  }
  return normalized
}

function applyLegacyTags(model: AIModel, tags: Array<Exclude<AIModelCapabilityKey, 'all'>>) {
  if (!model.is_thinking && tags.includes('reasoning')) {
    model.is_thinking = true
  }
  if (!model.support_vision && tags.includes('vision')) {
    model.support_vision = true
  }
  if (!model.support_tools && tags.includes('tool')) {
    model.support_tools = true
  }
  if (model.search_strategy === 'none' && tags.includes('search')) {
    model.search_strategy = 'builtin'
  }
  if (!model.support_embedding && tags.includes('embedding')) {
    model.support_embedding = true
  }
  if (!model.support_rerank && tags.includes('rerank')) {
    model.support_rerank = true
  }
  if (!model.is_free && tags.includes('free')) {
    model.is_free = true
  }
  if (model.search_strategy === 'tool') {
    model.support_tools = true
  }
}

export function normalizeModel(input?: Partial<AIModel & { tags?: string[] }> | null): AIModel {
  const model: AIModel = {
    id: String(input?.id ?? '').trim(),
    name: String(input?.name ?? '').trim(),
    description: String(input?.description ?? '').trim(),
    is_thinking: normalizeBooleanFlag(input?.is_thinking),
    support_vision: normalizeBooleanFlag(input?.support_vision),
    support_tools: normalizeBooleanFlag(input?.support_tools),
    search_strategy: normalizeSearchStrategy(input?.search_strategy),
    support_embedding: normalizeBooleanFlag(input?.support_embedding),
    support_rerank: normalizeBooleanFlag(input?.support_rerank),
    is_free: normalizeBooleanFlag(input?.is_free),
    temperature: normalizeNumberValue(input?.temperature),
    context_window: (() => {
      const value = normalizeNumberValue(input?.context_window)
      return value === null ? null : Math.max(0, Math.trunc(value))
    })(),
  }

  applyLegacyTags(model, normalizeCapabilityTags((input as { tags?: string[] } | null | undefined)?.tags))
  if (!model.name) {
    model.name = model.id
  }
  return model
}

export function normalizeRemoteProviderModel(input?: Partial<RemoteProviderModel> | null): RemoteProviderModel {
  const model = normalizeModel(input)
  return {
    ...model,
    object: String(input?.object ?? '').trim(),
    created: typeof input?.created === 'number' && Number.isFinite(input.created) ? input.created : undefined,
    owned_by: String(input?.owned_by ?? '').trim(),
    tags: Array.isArray(input?.tags) ? input?.tags.filter(tag => typeof tag === 'string') : undefined,
  }
}

export function getModelCapabilityTags(model?: Partial<AIModel> | null): Array<Exclude<AIModelCapabilityKey, 'all'>> {
  if (!model) {
    return []
  }
  const tags: Array<Exclude<AIModelCapabilityKey, 'all'>> = []
  if (model.is_thinking) {
    tags.push('reasoning')
  }
  if (model.support_vision) {
    tags.push('vision')
  }
  if (normalizeSearchStrategy(model.search_strategy) === 'builtin') {
    tags.push('search')
  }
  if (model.is_free) {
    tags.push('free')
  }
  if (model.support_embedding) {
    tags.push('embedding')
  }
  if (model.support_rerank) {
    tags.push('rerank')
  }
  if (model.support_tools) {
    tags.push('tool')
  }
  return tags
}

export function matchesModelCapability(model: AIModel, capability: AIModelCapabilityKey): boolean {
  if (capability === 'all') {
    return true
  }
  return getModelCapabilityTags(model).includes(capability)
}

export function filterModelsByCapabilityAndKeyword<T extends AIModel>(
  models: T[],
  capability: AIModelCapabilityKey,
  keyword?: string,
): T[] {
  const search = String(keyword ?? '').trim().toLowerCase()
  return models.filter((model) => {
    if (!matchesModelCapability(model, capability)) {
      return false
    }
    if (!search) {
      return true
    }
    return [model.id, model.name, model.description]
      .some(value => String(value ?? '').toLowerCase().includes(search))
  })
}

export function groupRemoteModelsByOwner(models: RemoteProviderModel[]): RemoteModelGroup[] {
  const groups = new Map<string, RemoteProviderModel[]>()

  for (const model of models) {
    const groupName = String(model.owned_by ?? '').trim() || '其他'
    if (!groups.has(groupName)) {
      groups.set(groupName, [])
    }
    groups.get(groupName)?.push(model)
  }

  return Array.from(groups.entries()).map(([name, items]) => ({
    name,
    items,
  }))
}

export function serializeModel(model: AIModel) {
  return {
    ...model,
    search_strategy: normalizeSearchStrategy(model.search_strategy),
    tags: getModelCapabilityTags(model),
  }
}

export function normalizeAIConfig(input?: Partial<AIConfigState> | null): AIConfigState {
  const providers = Array.isArray(input?.providers)
    ? input.providers.map(provider => ({
      name: String(provider?.name ?? ''),
      api_key: String(provider?.api_key ?? ''),
      base_url: String(provider?.base_url ?? ''),
      models: Array.isArray(provider?.models)
        ? provider.models.map(model => normalizeModel(model))
        : [],
    }))
    : []

  return {
    default_provider: String(input?.default_provider ?? ''),
    providers,
  }
}

export function mergeImportedModels(existingModels: AIModel[], importedModels: RemoteProviderModel[]): ImportMergeResult {
  const merged = existingModels.map(model => normalizeModel(model))
  const existingIDs = new Set(merged.map(model => model.id).filter(Boolean))

  let importedCount = 0
  let skippedCount = 0

  for (const imported of importedModels) {
    const normalized = normalizeModel(imported)
    if (!normalized.id || existingIDs.has(normalized.id)) {
      skippedCount += 1
      continue
    }
    merged.push(normalized)
    existingIDs.add(normalized.id)
    importedCount += 1
  }

  return {
    models: merged,
    importedCount,
    skippedCount,
  }
}

export function formatSearchStrategyLabel(value?: AISearchStrategy): string {
  return searchStrategyOptions.find(option => option.value === normalizeSearchStrategy(value))?.label ?? '不支持'
}
