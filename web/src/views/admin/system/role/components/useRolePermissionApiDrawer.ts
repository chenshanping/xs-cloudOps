import { computed, ref, watch, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import { assignApis, getRole } from '@/api/role'
import { getAllApis } from '@/api/api'
import { getMenuTree } from '@/api/menu'
import { useUserStore } from '@/store/user'
import type { Api, Menu } from '@/types'
import { filterAssignableMenuIds, matchesApiKeyword, normalizeText } from './permissionDrawer'

export interface RolePermissionApiDrawerProps {
  roleId: number
}

export interface ApiPermissionGroup {
  id: number
  label: string
  apis: Api[]
}

export function useRolePermissionApiDrawer(
  props: RolePermissionApiDrawerProps,
  visible: Ref<boolean>
) {
  const userStore = useUserStore()
  const canUseAdvancedApiPermissions = computed(() => userStore.permissions.includes('*'))

  const saveLoading = ref(false)
  const permissionLoading = ref(false)
  const permissionRequestToken = ref(0)
  const menuTree = ref<Menu[]>([])
  const allApis = ref<Api[]>([])
  const selectedMenuKeys = ref<number[]>([])
  const checkedApiIds = ref<number[]>([])
  const searchText = ref('')
  const selectedGroupId = ref<number | null>(null)

  const assignableSelectedMenuKeys = computed(() =>
    filterAssignableMenuIds(selectedMenuKeys.value, menuTree.value)
  )

  const inheritedApiSourceMap = computed<Record<number, string[]>>(() => {
    const selectedIds = new Set(assignableSelectedMenuKeys.value)
    const sources = new Map<number, Set<string>>()

    const walkMenus = (menus: Menu[]) => {
      menus.forEach(menu => {
        if (selectedIds.has(menu.id)) {
          menu.apis?.forEach(api => {
            const apiSources = sources.get(api.id) || new Set<string>()
            apiSources.add(menu.name)
            sources.set(api.id, apiSources)
          })
        }
        if (menu.children?.length) {
          walkMenus(menu.children)
        }
      })
    }

    walkMenus(menuTree.value)

    const result: Record<number, string[]> = {}
    sources.forEach((value, key) => {
      result[key] = Array.from(value).sort((a, b) => a.localeCompare(b, 'zh-CN'))
    })
    return result
  })

  const inheritedApiIds = computed(() =>
    Object.keys(inheritedApiSourceMap.value).map(id => Number(id))
  )

  const effectiveApiIds = computed(() =>
    Array.from(new Set([...checkedApiIds.value, ...inheritedApiIds.value]))
  )

  const apiGroups = computed<ApiPermissionGroup[]>(() => {
    const keyword = searchText.value.trim()
    const groups = new Map<string, Api[]>()

    allApis.value.forEach(api => {
      const groupName = (api.group || '').trim() || '未分组'
      if (keyword) {
        const apiMatch = matchesApiKeyword(api, keyword)
        const groupMatch = normalizeText(groupName).includes(normalizeText(keyword))
        if (!apiMatch && !groupMatch) {
          return
        }
      }

      const current = groups.get(groupName) || []
      current.push(api)
      groups.set(groupName, current)
    })

    return Array.from(groups.entries())
      .sort((left, right) => left[0].localeCompare(right[0], 'zh-CN'))
      .map(([label, apis], index) => ({
        id: index + 1,
        label,
        apis
      }))
  })

  const activeGroup = computed(() => {
    if (!apiGroups.value.length) {
      return null
    }
    return apiGroups.value.find(group => group.id === selectedGroupId.value) || apiGroups.value[0]
  })

  const addUniqueIds = (source: number[], incoming: number[]) =>
    Array.from(new Set([...source, ...incoming]))

  const removeIds = (source: number[], removing: number[]) =>
    source.filter(id => !removing.includes(id))

  const handleApiToggle = (apiId: number, checked: boolean) => {
    if (!canUseAdvancedApiPermissions.value) return
    checkedApiIds.value = checked
      ? addUniqueIds(checkedApiIds.value, [apiId])
      : removeIds(checkedApiIds.value, [apiId])
  }

  const resetState = () => {
    menuTree.value = []
    allApis.value = []
    selectedMenuKeys.value = []
    checkedApiIds.value = []
    selectedGroupId.value = null
  }

  const fetchMenuTree = async (token: number) => {
    const res = await getMenuTree()
    if (token !== permissionRequestToken.value) return
    menuTree.value = res.data
  }

  const fetchAllApis = async (token: number) => {
    if (!canUseAdvancedApiPermissions.value) {
      allApis.value = []
      return
    }
    const res = await getAllApis()
    if (token !== permissionRequestToken.value) return
    allApis.value = res.data.filter(api => api.need_auth)
  }

  const loadRolePermissions = async (token: number) => {
    if (!props.roleId) return
    const res = await getRole(props.roleId)
    if (token !== permissionRequestToken.value) return
    selectedMenuKeys.value = filterAssignableMenuIds(
      res.data.menus?.map((menu: Menu) => menu.id) || [],
      menuTree.value
    )
    checkedApiIds.value = canUseAdvancedApiPermissions.value
      ? res.data.apis?.map((api: Api) => api.id) || []
      : []
  }

  const getErrorMessage = (reason: unknown, fallback: string) => {
    if (reason instanceof Error && reason.message) {
      return reason.message
    }
    return fallback
  }

  const handleSavePermissions = async () => {
    if (permissionLoading.value || !canUseAdvancedApiPermissions.value) return
    saveLoading.value = true
    try {
      await assignApis(props.roleId, checkedApiIds.value, { silent: true })
      message.success('API 权限保存成功')
      visible.value = false
      await userStore.refreshAccessAction()
    } catch (error) {
      message.warning(`API 权限保存失败：${getErrorMessage(error, '请重试')}`)
    } finally {
      saveLoading.value = false
    }
  }

  watch([visible, () => props.roleId], async ([isVisible]) => {
    if (!isVisible) {
      permissionRequestToken.value += 1
      permissionLoading.value = false
      resetState()
      return
    }

    const token = ++permissionRequestToken.value
    permissionLoading.value = true
    searchText.value = ''
    resetState()

    try {
      await Promise.all([fetchMenuTree(token), fetchAllApis(token)])
      if (token !== permissionRequestToken.value) return
      await loadRolePermissions(token)
      selectedGroupId.value = apiGroups.value[0]?.id ?? null
    } finally {
      if (token !== permissionRequestToken.value) return
      permissionLoading.value = false
    }
  })

  watch(apiGroups, (groups) => {
    if (!groups.length) {
      selectedGroupId.value = null
      return
    }
    if (!groups.some(group => group.id === selectedGroupId.value)) {
      selectedGroupId.value = groups[0].id
    }
  }, { immediate: true })

  return {
    activeGroup,
    apiGroups,
    canUseAdvancedApiPermissions,
    checkedApiIds,
    effectiveApiIds,
    handleApiToggle,
    handleSavePermissions,
    inheritedApiIds,
    inheritedApiSourceMap,
    permissionLoading,
    saveLoading,
    searchText,
    selectedGroupId
  }
}
