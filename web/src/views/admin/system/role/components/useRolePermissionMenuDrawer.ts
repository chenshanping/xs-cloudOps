import { computed, ref, watch, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import { assignMenus, getRole } from '@/api/role'
import { getMenuTree } from '@/api/menu'
import { useUserStore } from '@/store/user'
import type { Menu } from '@/types'
import {
  buildPermissionViewModel,
  collectAssignableMenuIds,
  filterAssignableMenuIds,
  matchesMenuKeyword,
  normalizeMenuSelection,
  type PermissionMenuItem,
  type PermissionPageSection
} from './permissionDrawer'

export interface RolePermissionMenuDrawerProps {
  roleId: number
}

export interface FilteredMenuPermissionSection {
  id: string
  raw: PermissionPageSection
  visibleMenuItems: PermissionMenuItem[]
}

/**
 * 菜单权限抽屉 hook：只关心菜单/按钮 (`menu_ids`)，不涉及高级 API 直授。
 * 保存时调用 `PUT /roles/:id/menus`。
 */
export function useRolePermissionMenuDrawer(
  props: RolePermissionMenuDrawerProps,
  visible: Ref<boolean>
) {
  const userStore = useUserStore()

  const saveLoading = ref(false)
  const permissionLoading = ref(false)
  const permissionRequestToken = ref(0)
  const menuTree = ref<Menu[]>([])
  const selectedMenuKeys = ref<number[]>([])
  const searchText = ref('')
  const selectedTopMenuId = ref<number | null>(null)

  const checkedMenuKeys = computed(() =>
    normalizeMenuSelection(selectedMenuKeys.value, menuTree.value)
  )

  const assignableSelectedMenuKeys = computed(() =>
    filterAssignableMenuIds(selectedMenuKeys.value, menuTree.value)
  )

  const permissionViewModel = computed(() =>
    buildPermissionViewModel(menuTree.value, [])
  )

  const topMenus = computed(() => permissionViewModel.value.topMenus)

  const currentSections = computed(() => {
    if (selectedTopMenuId.value == null) {
      return []
    }
    return permissionViewModel.value.sectionsByTopMenu[selectedTopMenuId.value] || []
  })

  const filteredSections = computed<FilteredMenuPermissionSection[]>(() => {
    const keyword = searchText.value.trim()
    return currentSections.value
      .map(section => {
        let visibleMenuItems: PermissionMenuItem[]
        if (!keyword || matchesMenuKeyword(section.pageMenu, keyword)) {
          visibleMenuItems = section.menuItems
        } else {
          visibleMenuItems = section.menuItems.filter(item => matchesMenuKeyword(item.menu, keyword))
        }
        return {
          id: section.id,
          raw: section,
          visibleMenuItems
        }
      })
      .filter(section => {
        if (!keyword) return true
        return matchesMenuKeyword(section.raw.pageMenu, keyword)
          || section.visibleMenuItems.length > 0
      })
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

  const resetState = () => {
    menuTree.value = []
    selectedMenuKeys.value = []
    selectedTopMenuId.value = null
  }

  const fetchMenuTree = async (token: number) => {
    const res = await getMenuTree()
    if (token !== permissionRequestToken.value) return
    menuTree.value = res.data
    if (!selectedTopMenuId.value || !menuTree.value.some(item => item.id === selectedTopMenuId.value)) {
      selectedTopMenuId.value = menuTree.value[0]?.id ?? null
    }
  }

  const loadRolePermissions = async (token: number) => {
    if (!props.roleId) return
    const res = await getRole(props.roleId)
    if (token !== permissionRequestToken.value) return
    selectedMenuKeys.value = filterAssignableMenuIds(
      res.data.menus?.map((menu: Menu) => menu.id) || [],
      menuTree.value
    )
  }

  const getErrorMessage = (reason: unknown, fallback: string) => {
    if (reason instanceof Error && reason.message) {
      return reason.message
    }
    return fallback
  }

  const handleSavePermissions = async () => {
    if (permissionLoading.value) return
    saveLoading.value = true
    try {
      await assignMenus(props.roleId, assignableSelectedMenuKeys.value, { silent: true })
      message.success('菜单权限保存成功')
      visible.value = false
      await userStore.refreshAccessAction()
    } catch (error) {
      message.warning(`菜单权限保存失败：${getErrorMessage(error, '请重试')}`)
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
      await fetchMenuTree(token)
      if (token !== permissionRequestToken.value) return
      await loadRolePermissions(token)
    } finally {
      if (token !== permissionRequestToken.value) return
      permissionLoading.value = false
    }
  })

  return {
    assignableSelectedMenuKeys,
    checkedMenuKeys,
    filteredSections,
    handleMenuToggle,
    handleSavePermissions,
    handleSectionClearChildPermissions,
    handleSectionKeepPageOnly,
    handleSectionMenusToggle,
    handleSectionSelectChildPermissions,
    permissionLoading,
    saveLoading,
    searchText,
    selectedTopMenuId,
    topMenus
  }
}
