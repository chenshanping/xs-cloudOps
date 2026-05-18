import { ref } from 'vue'
import {
  getCmdbCredentialOptions,
  getCmdbHostGroups,
  getCmdbHostTags,
  type CmdbCredentialSummary,
  type CmdbHostGroup,
  type CmdbHostTag,
  type CmdbVerifyStatus,
} from '@/api/cmdb'

export const cmdbVerifyStatusOptions = [
  { label: '待校验', value: 'pending' },
  { label: '校验成功', value: 'success' },
  { label: '校验失败', value: 'failed' },
] as const

export const cmdbAuthTypeOptions = [
  { label: '密码', value: 'password' },
  { label: '私钥', value: 'private_key' },
] as const

export const getCmdbVerifyStatusMeta = (status?: CmdbVerifyStatus | string) => {
  switch (status) {
    case 'success':
      return { text: '校验成功', color: 'success' as const, badge: 'green' as const }
    case 'failed':
      return { text: '校验失败', color: 'error' as const, badge: 'red' as const }
    default:
      return { text: '待校验', color: 'default' as const, badge: 'default' as const }
  }
}

export const getCmdbAuthTypeLabel = (value?: string) =>
  cmdbAuthTypeOptions.find(item => item.value === value)?.label || value || '-'

export const formatCmdbMemory = (memoryMB?: number) => {
  if (!memoryMB) {
    return '-'
  }
  if (memoryMB >= 1024) {
    return `${(memoryMB / 1024).toFixed(memoryMB % 1024 === 0 ? 0 : 1)} GB`
  }
  return `${memoryMB} MB`
}

export function useCmdbReferenceOptions() {
  const loading = ref(false)
  const groupOptions = ref<CmdbHostGroup[]>([])
  const tagOptions = ref<CmdbHostTag[]>([])
  const credentialOptions = ref<CmdbCredentialSummary[]>([])

  const loadReferenceOptions = async () => {
    loading.value = true
    try {
      const [groupRes, tagRes, credentialRes] = await Promise.all([
        getCmdbHostGroups(),
        getCmdbHostTags(),
        getCmdbCredentialOptions(),
      ])
      groupOptions.value = groupRes.data || []
      tagOptions.value = tagRes.data || []
      credentialOptions.value = credentialRes.data || []
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    groupOptions,
    tagOptions,
    credentialOptions,
    loadReferenceOptions,
  }
}
