import { computed, ref, watch, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  getDataScopeResources,
  getRole,
  saveRolePermissions,
  type DataScopeResource,
  type RoleFeatureDataScopePayload
} from '@/api/role'
import { getMenuTree } from '@/api/menu'
import { getAllApis } from '@/api/api'
import { useUserStore } from '@/store/user'
import type { Api, Menu } from '@/types'
import {
  buildPermissionViewModel,
  collectAssignableMenuIds,
  filterAssignableMenuIds,
  matchesApiKeyword,
  matchesMenuKeyword,
  normalizeMenuSelection,
  type PermissionMenuItem,
  type PermissionPageSection
} from './permissionDrawer'
import {
  buildRoleFeatureDataScopePayload,
  buildRoleFeatureDataScopeForm,
  findDataScopeResourceDefinition,
  formatDataScopeLabel,
  splitKnownAndUnknownFeatureDataScopes,
  type RoleFeatureDataScopeFormItem,
  type RolePermissionDeptOption
} from './dataScopeResources'

export interface RolePermissionDrawerProps {
  roleId: number
  deptOptions: RolePermissionDeptOption[]
}

export interface FilteredPermissionSection {
  id: string
  raw: PermissionPageSection
  visibleMenuItems: PermissionMenuItem[]
  visibleApis: Api[]
}

export function useRolePermissionDrawer(
  props: RolePermissionDrawerProps,
  visible: Ref<boolean>
) {
  const userStore = useUserStore()

  const saveLoading = ref(false)
  const permissionLoading = ref(false)
  const permissionRequestToken = ref(0)
  const menuTree = ref<Menu[]>([])
  const allApis = ref<Api[]>([])
  const resourceDefinitions = ref<DataScopeResource[]>([])
  const selectedMenuKeys = ref<number[]>([])
  const checkedApiIds = ref<number[]>([])
  const defaultDataScope = ref(1)
  const featureDataScopes = ref<RoleFeatureDataScopeFormItem[]>([])
  const unknownFeatureDataScopes = ref<RoleFeatureDataScopePayload[]>([])
  const searchText = ref('')
  const selectedTopMenuId = ref<number | null>(null)

  const checkedMenuKeys = computed(() =>
    normalizeMenuSelection(selectedMenuKeys.value, menuTree.value)
  )

  const assignableSelectedMenuKeys = computed(() =>
    filterAssignableMenuIds(selectedMenuKeys.value, menuTree.value)
  )

  const permissionViewModel = computed(() =>
    buildPermissionViewModel(menuTree.value, allApis.value)
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
      result[key] = Array.from(value).sort((left, right) => left.localeCompare(right, 'zh-CN'))
    })
    return result
  })

  const inheritedApiIds = computed(() =>
    Object.keys(inheritedApiSourceMap.value).map(id => Number(id))
  )

  const effectiveApiIds = computed(() =>
    Array.from(new Set([...checkedApiIds.value, ...inheritedApiIds.value]))
  )

  const topMenus = computed(() => permissionViewModel.value.topMenus)

  const currentSections = computed(() => {
    if (selectedTopMenuId.value == null) {
      return []
    }
    return permissionViewModel.value.sectionsByTopMenu[selectedTopMenuId.value] || []
  })

  const uncategorizedApis = computed(() => permissionViewModel.value.uncategorizedApis)

  function buildVisibleMenuItems(section: PermissionPageSection) {
    const keyword = searchText.value.trim()
    if (!keyword) {
      return section.menuItems
    }
    if (matchesMenuKeyword(section.pageMenu, keyword)) {
      return section.menuItems
    }
    return section.menuItems.filter(item => matchesMenuKeyword(item.menu, keyword))
  }

  const filteredSections = computed<FilteredPermissionSection[]>(() => {
    const keyword = searchText.value.trim()
    return currentSections.value
      .map(section => {
        const visibleMenuItems = buildVisibleMenuItems(section)
        const visibleApis = !keyword || matchesMenuKeyword(section.pageMenu, keyword)
          ? section.apis
          : section.apis.filter(api => matchesApiKeyword(api, keyword))

        return {
          id: section.id,
          raw: section,
          visibleMenuItems,
          visibleApis
        }
      })
      .filter(section => {
        if (!keyword) {
          return true
        }
        return (
          matchesMenuKeyword(section.raw.pageMenu, keyword) ||
          section.visibleMenuItems.length > 0 ||
          section.visibleApis.length > 0
        )
      })
  })

  const filteredUncategorizedApis = computed(() => {
    const keyword = searchText.value.trim()
    if (!keyword) {
      return uncategorizedApis.value
    }
    return uncategorizedApis.value.filter(api => matchesApiKeyword(api, keyword))
  })

  const isUncategorizedChecked = computed(() => {
    const ids = uncategorizedApis.value.map(api => api.id)
    return ids.length > 0 && ids.every(id => checkedApiIds.value.includes(id))
  })

  const isUncategorizedIndeterminate = computed(() => {
    const ids = uncategorizedApis.value.map(api => api.id)
    const checkedCount = ids.filter(id => checkedApiIds.value.includes(id)).length
    return checkedCount > 0 && checkedCount < ids.length
  })

  const addUniqueIds = (source: number[], incoming: number[]) =>
    Array.from(new Set([...source, ...incoming]))

  const removeIds = (source: number[], removing: number[]) =>
    source.filter(id => !removing.includes(id))

  const handleMenuToggle = (menu: Menu, checked: boolean) => {
    const ids = menu.type === 1 ? collectAssignableMenuIds(menu) : [menu.id]
    selectedMenuKeys.value = checked
      ? addUniqueIds(selectedMenuKeys.value, ids)
      : removeIds(selectedMenuKeys.value, ids)
  }

  const handleSectionMenusToggle = (section: PermissionPageSection, checked: boolean) => {
    const ids = section.menuItems
      .filter(item => item.menu.type !== 1)
      .map(item => item.menu.id)
    selectedMenuKeys.value = checked
      ? addUniqueIds(selectedMenuKeys.value, ids)
      : removeIds(selectedMenuKeys.value, ids)
  }

  const handleSectionKeepPageOnly = (section: PermissionPageSection) => {
    const pageId = section.pageMenu.type !== 1 ? section.pageMenu.id : null
    const childIds = collectAssignableMenuIds(section.pageMenu).filter(id => id !== pageId)
    const baseIds = removeIds(selectedMenuKeys.value, childIds)
    selectedMenuKeys.value = pageId == null
      ? baseIds
      : addUniqueIds(baseIds, [pageId])
  }

  const handleSectionSelectChildPermissions = (section: PermissionPageSection) => {
    const ids = collectAssignableMenuIds(section.pageMenu)
    selectedMenuKeys.value = addUniqueIds(selectedMenuKeys.value, ids)
  }

  const handleSectionClearChildPermissions = (section: PermissionPageSection) => {
    const pageId = section.pageMenu.type !== 1 ? section.pageMenu.id : null
    const childIds = collectAssignableMenuIds(section.pageMenu).filter(id => id !== pageId)
    selectedMenuKeys.value = removeIds(selectedMenuKeys.value, childIds)
  }

  const handleApiToggle = (apiId: number, checked: boolean) => {
    checkedApiIds.value = checked
      ? addUniqueIds(checkedApiIds.value, [apiId])
      : removeIds(checkedApiIds.value, [apiId])
  }

  const handleSectionApisToggle = (section: PermissionPageSection, checked: boolean) => {
    const ids = section.apis.map(api => api.id)
    checkedApiIds.value = checked
      ? addUniqueIds(checkedApiIds.value, ids)
      : removeIds(checkedApiIds.value, ids)
  }

  const handleUncategorizedToggle = (checked: boolean) => {
    const ids = uncategorizedApis.value.map(api => api.id)
    checkedApiIds.value = checked
      ? addUniqueIds(checkedApiIds.value, ids)
      : removeIds(checkedApiIds.value, ids)
  }

  const resetPermissionState = () => {
    menuTree.value = []
    allApis.value = []
    resourceDefinitions.value = []
    selectedMenuKeys.value = []
    checkedApiIds.value = []
    defaultDataScope.value = 1
    featureDataScopes.value = []
    unknownFeatureDataScopes.value = []
    selectedTopMenuId.value = null
  }

  const fetchMenuTree = async (requestToken: number) => {
    const res = await getMenuTree()
    if (requestToken !== permissionRequestToken.value) {
      return
    }
    menuTree.value = res.data
    if (!selectedTopMenuId.value || !menuTree.value.some(item => item.id === selectedTopMenuId.value)) {
      selectedTopMenuId.value = menuTree.value[0]?.id ?? null
    }
  }

  const fetchAllApis = async (requestToken: number) => {
    const res = await getAllApis()
    if (requestToken !== permissionRequestToken.value) {
      return
    }
    allApis.value = res.data.filter(api => api.need_auth)
  }

  const fetchDataScopeResources = async (requestToken: number) => {
    const res = await getDataScopeResources()
    if (requestToken !== permissionRequestToken.value) {
      return
    }
    resourceDefinitions.value = res.data || []
  }

  const loadRolePermissions = async (requestToken: number) => {
    if (!props.roleId) {
      return
    }
    const res = await getRole(props.roleId)
    if (requestToken !== permissionRequestToken.value) {
      return
    }
    const scopeResources = resourceDefinitions.value
    const {
      knownScopes,
      unknownScopes
    } = splitKnownAndUnknownFeatureDataScopes(scopeResources, res.data.feature_data_scopes || [])

    defaultDataScope.value = res.data.data_scope || 1
    featureDataScopes.value = buildRoleFeatureDataScopeForm(
      scopeResources,
      knownScopes
    )
    unknownFeatureDataScopes.value = unknownScopes
    selectedMenuKeys.value = filterAssignableMenuIds(
      res.data.menus?.map((menu: Menu) => menu.id) || [],
      menuTree.value
    )
    checkedApiIds.value = res.data.apis?.map((api: Api) => api.id) || []
  }

  const validateFeatureDataScopes = () => {
    const invalidScope = featureDataScopes.value.find(item => item.data_scope === 2 && item.dept_ids.length === 0)
    if (invalidScope) {
      const resourceLabel = findDataScopeResourceDefinition(
        resourceDefinitions.value,
        invalidScope.resource_code
      )?.label || invalidScope.resource_code
      message.warning(`请为「${resourceLabel}」选择自定义部门`)
      return false
    }
    return true
  }

  const getErrorMessage = (reason: unknown, fallback: string) => {
    if (reason instanceof Error && reason.message) {
      return reason.message
    }
    return fallback
  }

  const handleSavePermissions = async () => {
    if (permissionLoading.value) {
      return
    }
    if (!validateFeatureDataScopes()) {
      return
    }

    saveLoading.value = true
    try {
      await saveRolePermissions(props.roleId, {
        menu_ids: assignableSelectedMenuKeys.value,
        direct_api_ids: checkedApiIds.value,
        scopes: buildRoleFeatureDataScopePayload(
          featureDataScopes.value,
          unknownFeatureDataScopes.value
        )
      }, { silent: true })
      message.success('权限分配成功')
      visible.value = false
      await userStore.refreshAccessAction()
    } catch (error) {
      message.warning(`权限保存失败：${getErrorMessage(error, '请重试')}`)
    } finally {
      saveLoading.value = false
    }
  }

  watch([visible, () => props.roleId], async ([isVisible]) => {
    if (!isVisible) {
      permissionRequestToken.value += 1
      permissionLoading.value = false
      resetPermissionState()
      return
    }

    const requestToken = ++permissionRequestToken.value
    permissionLoading.value = true
    searchText.value = ''
    resetPermissionState()

    try {
      await Promise.all([
        fetchMenuTree(requestToken),
        fetchAllApis(requestToken),
        fetchDataScopeResources(requestToken)
      ])
      if (requestToken !== permissionRequestToken.value) {
        return
      }
      await loadRolePermissions(requestToken)
    } finally {
      if (requestToken !== permissionRequestToken.value) {
        return
      }
      permissionLoading.value = false
    }
  })

  return {
    checkedApiIds,
    checkedMenuKeys,
    assignableSelectedMenuKeys,
    currentSections,
    defaultDataScope,
    effectiveApiIds,
    featureDataScopes,
    formatDataScopeLabel,
    filteredSections,
    filteredUncategorizedApis,
    handleApiToggle,
    handleMenuToggle,
    handleSavePermissions,
    handleSectionClearChildPermissions,
    handleSectionKeepPageOnly,
    handleSectionApisToggle,
    handleSectionMenusToggle,
    handleSectionSelectChildPermissions,
    handleUncategorizedToggle,
    inheritedApiIds,
    inheritedApiSourceMap,
    isUncategorizedChecked,
    isUncategorizedIndeterminate,
    permissionLoading,
    resourceDefinitions,
    saveLoading,
    searchText,
    selectedTopMenuId,
    topMenus,
    uncategorizedApis
  }
}
