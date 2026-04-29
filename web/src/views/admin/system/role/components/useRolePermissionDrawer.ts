import { computed, ref, watch, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  getRole,
  assignMenus,
  assignApis,
  assignDataScopes,
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
  buildRoleFeatureDataScopeForm,
  formatDataScopeLabel,
  ROLE_FEATURE_SCOPE_RESOURCES,
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
  const menuTree = ref<Menu[]>([])
  const allApis = ref<Api[]>([])
  const selectedMenuKeys = ref<number[]>([])
  const checkedApiIds = ref<number[]>([])
  const defaultDataScope = ref(1)
  const featureDataScopes = ref<RoleFeatureDataScopeFormItem[]>(buildRoleFeatureDataScopeForm())
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

  const fetchMenuTree = async () => {
    const res = await getMenuTree()
    menuTree.value = res.data
    if (!selectedTopMenuId.value || !menuTree.value.some(item => item.id === selectedTopMenuId.value)) {
      selectedTopMenuId.value = menuTree.value[0]?.id ?? null
    }
  }

  const fetchAllApis = async () => {
    const res = await getAllApis()
    allApis.value = res.data.filter(api => api.need_auth)
  }

  const loadRolePermissions = async () => {
    if (!props.roleId) {
      return
    }
    const res = await getRole(props.roleId)
    defaultDataScope.value = res.data.data_scope || 1
    featureDataScopes.value = buildRoleFeatureDataScopeForm(res.data.feature_data_scopes)
    selectedMenuKeys.value = filterAssignableMenuIds(
      res.data.menus?.map((menu: Menu) => menu.id) || [],
      menuTree.value
    )
    checkedApiIds.value = res.data.apis?.map((api: Api) => api.id) || []
  }

  const validateFeatureDataScopes = () => {
    const invalidScope = featureDataScopes.value.find(item => item.data_scope === 2 && item.dept_ids.length === 0)
    if (invalidScope) {
      const resourceLabel = ROLE_FEATURE_SCOPE_RESOURCES.find(item => item.code === invalidScope.resource_code)?.label || invalidScope.resource_code
      message.warning(`请为「${resourceLabel}」选择自定义部门`)
      return false
    }
    return true
  }

  const buildFeatureDataScopePayload = (): RoleFeatureDataScopePayload[] =>
    featureDataScopes.value
      .filter(item => item.data_scope > 0)
      .map(item => ({
        resource_code: item.resource_code,
        data_scope: item.data_scope,
        dept_ids: item.data_scope === 2 ? [...item.dept_ids] : []
      }))

  const getErrorMessage = (reason: unknown, fallback: string) => {
    if (reason instanceof Error && reason.message) {
      return reason.message
    }
    return fallback
  }

  const handleSavePermissions = async () => {
    if (!validateFeatureDataScopes()) {
      return
    }

    saveLoading.value = true
    try {
      const tasks = [
        {
          label: '菜单权限',
          promise: assignMenus(props.roleId, assignableSelectedMenuKeys.value, { silent: true })
        },
        {
          label: 'API 权限',
          promise: assignApis(props.roleId, checkedApiIds.value, { silent: true })
        },
        {
          label: '数据权限',
          promise: assignDataScopes(props.roleId, buildFeatureDataScopePayload(), { silent: true })
        }
      ]

      const results = await Promise.allSettled(tasks.map(task => task.promise))
      const failedTasks = results
        .map((result, index) => ({ result, label: tasks[index].label }))
        .filter(item => item.result.status === 'rejected')

      if (!failedTasks.length) {
        message.success('权限分配成功')
        visible.value = false
        await userStore.refreshAccessAction()
        return
      }

      if (failedTasks.length === 1) {
        const failedTask = failedTasks[0]
        if (failedTask.result.status === 'rejected') {
          message.warning(`${failedTask.label}保存失败：${getErrorMessage(failedTask.result.reason, '请重试')}`)
        }
        return
      }

      message.error(`${failedTasks.map(item => item.label).join('、')}保存失败，请重试`)
    } finally {
      saveLoading.value = false
    }
  }

  watch(visible, async val => {
    if (!val) {
      featureDataScopes.value = buildRoleFeatureDataScopeForm()
      return
    }
    searchText.value = ''
    await Promise.all([fetchMenuTree(), fetchAllApis()])
    await loadRolePermissions()
  })

  return {
    checkedApiIds,
    checkedMenuKeys,
    assignableSelectedMenuKeys,
    currentSections,
    defaultDataScope,
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
    isUncategorizedChecked,
    isUncategorizedIndeterminate,
    saveLoading,
    searchText,
    selectedTopMenuId,
    topMenus,
    uncategorizedApis
  }
}
