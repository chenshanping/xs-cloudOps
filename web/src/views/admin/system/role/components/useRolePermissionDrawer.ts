import { computed, ref, watch, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import { getRole, assignMenus, assignApis } from '@/api/role'
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

export interface RolePermissionDrawerProps {
  roleId: number
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
    selectedMenuKeys.value = filterAssignableMenuIds(
      res.data.menus?.map((menu: Menu) => menu.id) || [],
      menuTree.value
    )
    checkedApiIds.value = res.data.apis?.map((api: Api) => api.id) || []
  }

  const getErrorMessage = (reason: unknown, fallback: string) => {
    if (reason instanceof Error && reason.message) {
      return reason.message
    }
    return fallback
  }

  const handleSavePermissions = async () => {
    saveLoading.value = true
    try {
      const results = await Promise.allSettled([
        assignMenus(props.roleId, assignableSelectedMenuKeys.value, { silent: true }),
        assignApis(props.roleId, checkedApiIds.value, { silent: true })
      ])

      const [menuResult, apiResult] = results
      const menuFailed = menuResult.status === 'rejected'
      const apiFailed = apiResult.status === 'rejected'

      if (!menuFailed && !apiFailed) {
        message.success('权限分配成功')
        visible.value = false
        await userStore.refreshAccessAction()
        return
      }

      if (menuFailed && apiFailed) {
        message.error('菜单权限和 API 权限保存失败，请重试')
        return
      }

      if (menuFailed) {
        message.warning(`菜单权限保存失败：${getErrorMessage(menuResult.reason, '请重试')}`)
        return
      }

      message.warning(`API 权限保存失败：${getErrorMessage(apiResult.reason, '请重试')}`)
    } finally {
      saveLoading.value = false
    }
  }

  watch(visible, async val => {
    if (!val) {
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
