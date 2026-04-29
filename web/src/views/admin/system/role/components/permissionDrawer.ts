import type { Api, Menu } from '@/types'

export interface PermissionMenuItem {
  menu: Menu
  level: number
}

export interface PermissionPageSection {
  id: string
  topMenuId: number
  topMenuName: string
  pageMenu: Menu
  menuItems: PermissionMenuItem[]
  apis: Api[]
}

export interface PermissionViewModel {
  topMenus: Menu[]
  sectionsByTopMenu: Record<number, PermissionPageSection[]>
  uncategorizedApis: Api[]
}

export interface UncategorizedApiGroup {
  id: string
  label: string
  apis: Api[]
}

const API_SECTION_OVERRIDES: Record<string, string> = {
  '/api/v1/logs/route-groups': '操作日志'
}

const PATH_STOP_WORDS = new Set(['api', 'v1', 'system', 'admin'])
const API_ACTION_WORDS = new Set([
  'list',
  'detail',
  'create',
  'update',
  'delete',
  'batch',
  'test',
  'send',
  'fetch',
  'save',
  'get',
  'put',
  'post'
])

const SECTION_ALIAS_MAP: Record<string, string[]> = {
  参数配置: ['配置管理', '系统配置', 'config', 'configs', 'setting', 'settings'],
  AI配置: ['ai配置', 'ai管理', 'provider', 'model'],
  用户管理: ['users', 'user'],
  角色管理: ['roles', 'role'],
  部门管理: ['depts', 'dept', 'department'],
  菜单管理: ['menus', 'menu'],
  API管理: ['apis', 'api'],
  文件管理: ['files', 'file', 'storage'],
  操作日志: ['operationlog', 'operation', 'audit'],
  登录日志: ['loginlog', 'login']
}

export function normalizeText(value: string | undefined | null): string {
  return (value || '')
    .toLowerCase()
    .replace(/[\s/_:-]+/g, '')
    .trim()
}

export function extractPathSegments(value: string | undefined | null): string[] {
  return (value || '')
    .toLowerCase()
    .split('/')
    .map(segment => segment.trim())
    .filter(segment => segment && !PATH_STOP_WORDS.has(segment))
}

function buildSectionKeywords(section: PermissionPageSection): string[] {
  const values = new Set<string>()
  const pageName = section.pageMenu.name
  const pageText = normalizeText(pageName)
  const permissionText = normalizeText(section.pageMenu.permission)
  const pathSegments = extractPathSegments(section.pageMenu.path)

  if (pageText) {
    values.add(pageText)
  }
  if (permissionText) {
    values.add(permissionText)
  }
  pathSegments.forEach(segment => {
    values.add(segment)
  })

  const aliases = SECTION_ALIAS_MAP[pageName] || []
  aliases.forEach(alias => {
    const normalized = normalizeText(alias)
    if (normalized) {
      values.add(normalized)
    }
  })

  return Array.from(values)
}

export function collectMenuIds(menu: Menu): number[] {
  const ids = [menu.id]
  menu.children?.forEach(child => {
    ids.push(...collectMenuIds(child))
  })
  return ids
}

export function collectAssignableMenuIds(menu: Menu): number[] {
  const ids: number[] = []
  if (menu.type !== 1) {
    ids.push(menu.id)
  }
  menu.children?.forEach(child => {
    ids.push(...collectAssignableMenuIds(child))
  })
  return ids
}

export function filterAssignableMenuIds(selectedIds: number[], menuTree: Menu[]): number[] {
  const menuMap = new Map<number, Menu>()
  const walk = (menus: Menu[]) => {
    menus.forEach(menu => {
      menuMap.set(menu.id, menu)
      if (menu.children?.length) {
        walk(menu.children)
      }
    })
  }
  walk(menuTree)

  return selectedIds.filter(id => {
    const menu = menuMap.get(id)
    return menu ? menu.type !== 1 : false
  })
}

export function flattenMenuItems(menu: Menu, level = 0): PermissionMenuItem[] {
  const items: PermissionMenuItem[] = [{ menu, level }]
  menu.children?.forEach(child => {
    items.push(...flattenMenuItems(child, level + 1))
  })
  return items
}

export function findMenuPath(
  menus: Menu[],
  targetId: number,
  ancestors: number[] = []
): number[] | null {
  for (const menu of menus) {
    const currentPath = [...ancestors, menu.id]
    if (menu.id === targetId) {
      return currentPath
    }
    if (menu.children?.length) {
      const childPath = findMenuPath(menu.children, targetId, currentPath)
      if (childPath) {
        return childPath
      }
    }
  }
  return null
}

export function normalizeMenuSelection(selectedIds: number[], menuTree: Menu[]): number[] {
  const normalized = new Set(selectedIds)
  selectedIds.forEach(id => {
    const path = findMenuPath(menuTree, id)
    path?.forEach(pathId => normalized.add(pathId))
  })
  return Array.from(normalized)
}

function scoreApiToSection(api: Api, section: PermissionPageSection): number {
  const groupText = normalizeText(api.group)
  const pageName = normalizeText(section.pageMenu.name)
  const descriptionText = normalizeText(api.description)
  const pathSegments = extractPathSegments(api.path)
  const sectionKeywords = buildSectionKeywords(section)
  const matchedPathSegments = pathSegments.filter(segment => !API_ACTION_WORDS.has(segment))

  if (!groupText && !descriptionText && pathSegments.length === 0) {
    return 0
  }
  if (groupText && groupText === pageName) {
    return 100
  }
  if (groupText && (groupText.includes(pageName) || pageName.includes(groupText))) {
    return 80
  }

  if (sectionKeywords.some(keyword => keyword && groupText.includes(keyword))) {
    return 88
  }

  if (sectionKeywords.some(keyword => keyword && descriptionText.includes(keyword))) {
    return 72
  }

  if (matchedPathSegments.some(segment => sectionKeywords.includes(segment))) {
    return 68
  }

  if (pageName && descriptionText.includes(pageName)) {
    return 60
  }

  return 0
}

function getSectionRoots(topMenu: Menu): Menu[] {
  const children = topMenu.children || []
  const roots = children.filter(child => child.type !== 3)
  if (roots.length > 0) {
    return roots
  }
  return topMenu.type !== 3 ? [topMenu] : []
}

export function buildPermissionViewModel(menuTree: Menu[], allApis: Api[]): PermissionViewModel {
  const topMenus = menuTree
  const sectionsByTopMenu: Record<number, PermissionPageSection[]> = {}
  const allSections: PermissionPageSection[] = []

  topMenus.forEach(topMenu => {
    const sections = getSectionRoots(topMenu).map(root => {
      const section: PermissionPageSection = {
        id: `section-${topMenu.id}-${root.id}`,
        topMenuId: topMenu.id,
        topMenuName: topMenu.name,
        pageMenu: root,
        menuItems: flattenMenuItems(root),
        apis: []
      }
      allSections.push(section)
      return section
    })
    sectionsByTopMenu[topMenu.id] = sections
  })

  const uncategorizedApis: Api[] = []
  allApis.forEach(api => {
    const overridePageName = API_SECTION_OVERRIDES[api.path]
    if (overridePageName) {
      const matchedSection = allSections.find(section => section.pageMenu.name === overridePageName)
      if (matchedSection) {
        matchedSection.apis.push(api)
        return
      }
    }

    let bestSection: PermissionPageSection | null = null
    let bestScore = 0
    allSections.forEach(section => {
      const score = scoreApiToSection(api, section)
      if (score > bestScore) {
        bestScore = score
        bestSection = section
      }
    })
    if (bestSection && bestScore >= 68) {
      bestSection.apis.push(api)
      return
    }
    uncategorizedApis.push(api)
  })

  return {
    topMenus,
    sectionsByTopMenu,
    uncategorizedApis
  }
}

export function matchesMenuKeyword(menu: Menu, keyword: string): boolean {
  if (!keyword) {
    return true
  }
  const normalizedKeyword = normalizeText(keyword)
  return [
    menu.name,
    menu.permission,
    menu.path
  ].some(value => normalizeText(value).includes(normalizedKeyword))
}

export function matchesApiKeyword(api: Api, keyword: string): boolean {
  if (!keyword) {
    return true
  }
  const normalizedKeyword = normalizeText(keyword)
  return [
    api.group,
    api.path,
    api.description,
    api.method
  ].some(value => normalizeText(value).includes(normalizedKeyword))
}

export function groupUncategorizedApis(apis: Api[]): UncategorizedApiGroup[] {
  const buckets = new Map<string, UncategorizedApiGroup>()

  apis.forEach(api => {
    const rawGroup = (api.group || '').trim()
    const label = rawGroup || '其他系统接口'
    const id = `system-${normalizeText(label) || 'other'}`

    if (!buckets.has(id)) {
      buckets.set(id, {
        id,
        label,
        apis: []
      })
    }

    buckets.get(id)!.apis.push(api)
  })

  return Array.from(buckets.values()).sort((left, right) => left.label.localeCompare(right.label, 'zh-CN'))
}
