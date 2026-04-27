export interface AIModel {
  id: string
  name: string
  description: string
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

export interface RemoteProviderModel {
  id: string
  object?: string
  created?: number
  owned_by?: string
}

export interface ImportMergeResult {
  models: AIModel[]
  importedCount: number
  skippedCount: number
}

export function createEmptyModel(): AIModel {
  return {
    id: '',
    name: '',
    description: '',
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

export function normalizeAIConfig(input?: Partial<AIConfigState> | null): AIConfigState {
  const providers = Array.isArray(input?.providers)
    ? input!.providers.map(provider => ({
      name: provider?.name ?? '',
      api_key: provider?.api_key ?? '',
      base_url: provider?.base_url ?? '',
      models: Array.isArray(provider?.models)
        ? provider.models.map(model => ({
          id: model?.id ?? '',
          name: model?.name ?? '',
          description: model?.description ?? '',
        }))
        : [],
    }))
    : []

  return {
    default_provider: input?.default_provider ?? '',
    providers,
  }
}

export function mergeImportedModels(existingModels: AIModel[], importedModels: RemoteProviderModel[]): ImportMergeResult {
  const merged = existingModels.map(model => ({
    id: model.id ?? '',
    name: model.name ?? '',
    description: model.description ?? '',
  }))
  const existingIDs = new Set(merged.map(model => model.id))

  let importedCount = 0
  let skippedCount = 0

  for (const imported of importedModels) {
    const id = imported.id?.trim() ?? ''
    if (!id || existingIDs.has(id)) {
      skippedCount += 1
      continue
    }
    merged.push({
      id,
      name: id,
      description: '',
    })
    existingIDs.add(id)
    importedCount += 1
  }

  return {
    models: merged,
    importedCount,
    skippedCount,
  }
}
