import type { StorageType } from '../../../../types/storage'

type StorageConfigDraft = Record<string, any>
type StorageConfigValues = Partial<Record<string, string>>

const STORAGE_TYPES: StorageType[] = ['local', 'aliyun', 'tencent', 'minio']

export const STORAGE_CONFIG_KEY_MAP: Record<StorageType, string> = {
  local: 'storage_local_config',
  aliyun: 'storage_aliyun_config',
  tencent: 'storage_tencent_config',
  minio: 'storage_minio_config',
}

export interface StorageDraftsState {
  activeType: StorageType
  drafts: Record<StorageType, StorageConfigDraft>
}

export function getDefaultStorageConfig(type: StorageType): StorageConfigDraft {
  switch (type) {
    case 'aliyun':
      return {
        endpoint: '',
        access_key_id: '',
        access_key_secret: '',
        bucket_name: '',
        region: '',
      }
    case 'tencent':
      return {
        region: '',
        secret_id: '',
        secret_key: '',
        bucket: '',
        app_id: '',
      }
    case 'minio':
      return {
        endpoint: '',
        access_key_id: '',
        secret_access_key: '',
        bucket_name: '',
        use_ssl: false,
      }
    case 'local':
    default:
      return {
        base_path: 'uploads',
        base_url: '/api/v1/upload',
      }
  }
}

function normalizeStorageType(value?: string): StorageType {
  if (value && STORAGE_TYPES.includes(value as StorageType)) {
    return value as StorageType
  }
  return 'local'
}

export function parseStorageConfig(type: StorageType, rawValue?: string): StorageConfigDraft {
  let parsed: StorageConfigDraft = {}
  if (rawValue) {
    try {
      parsed = JSON.parse(rawValue)
    } catch {
      parsed = {}
    }
  }

  return {
    ...getDefaultStorageConfig(type),
    ...parsed,
  }
}

export function createStorageDraftsState(values: StorageConfigValues): StorageDraftsState {
  const drafts = {} as Record<StorageType, StorageConfigDraft>
  for (const type of STORAGE_TYPES) {
    drafts[type] = parseStorageConfig(type, values[STORAGE_CONFIG_KEY_MAP[type]])
  }

  return {
    activeType: normalizeStorageType(values.storage_type),
    drafts,
  }
}

export function getActiveStorageConfig(state: StorageDraftsState): StorageConfigDraft {
  return state.drafts[state.activeType]
}

export function changeActiveStorageType(state: StorageDraftsState, storageType: StorageType): StorageDraftsState {
  return {
    ...state,
    activeType: normalizeStorageType(storageType),
  }
}

export function getStorageBucketLabel(type: StorageType, configJSON?: string): string {
  let config: Record<string, any> = {}
  if (configJSON) {
    try {
      config = JSON.parse(configJSON)
    } catch {
      config = {}
    }
  }

  switch (type) {
    case 'aliyun':
      return config.bucket_name || ''
    case 'tencent':
      return config.bucket || ''
    case 'minio':
      return config.bucket_name || ''
    case 'local':
      return config.base_path || 'uploads'
    default:
      return ''
  }
}

export function buildStorageSavePayload(state: StorageDraftsState): Record<string, string> {
  const activeType = state.activeType
  return {
    storage_type: activeType,
    [STORAGE_CONFIG_KEY_MAP[activeType]]: JSON.stringify(getActiveStorageConfig(state)),
  }
}
