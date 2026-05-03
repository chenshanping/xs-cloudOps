import type { Menu } from '@/types'

export function normalizePath(path?: string) {
  if (!path) {
    return ''
  }
  return path.startsWith('/') ? path : `/${path}`
}

export function filterEnabledMenus(menus: Menu[] = []): Menu[] {
  return menus
    .filter((menu) => menu.status === 1)
    .map((menu) => ({
      ...menu,
      path: normalizePath(menu.path),
      children: filterEnabledMenus(menu.children || []),
    }))
}

export function filterVisibleMenus(menus: Menu[] = []): Menu[] {
  return menus
    .filter((menu) => menu.status === 1 && menu.hidden !== 1)
    .map((menu) => ({
      ...menu,
      path: normalizePath(menu.path),
      children: filterVisibleMenus(menu.children || []),
    }))
}

export function findMenuTrail(menus: Menu[], path: string): Menu[] {
  for (const menu of menus) {
    const menuPath = normalizePath(menu.path)
    if (menuPath === path) {
      return [menu]
    }
    if (menu.children?.length) {
      const childTrail = findMenuTrail(menu.children, path)
      if (childTrail.length) {
        return [menu, ...childTrail]
      }
    }
  }
  return []
}

export function firstNavigablePath(menu?: Menu | null): string {
  if (!menu) {
    return ''
  }
  if (menu.type === 2 && menu.path) {
    return normalizePath(menu.path)
  }
  for (const child of menu.children || []) {
    const target = firstNavigablePath(child)
    if (target) {
      return target
    }
  }
  return ''
}

export function getMixedSidebarMenus(menus: Menu[], path: string): Menu[] {
  const trail = findMenuTrail(menus, path)
  const activeTop = trail[0]

  if (!activeTop) {
    return []
  }

  if (activeTop.type === 1 && activeTop.children?.length) {
    return activeTop.children
  }

  return []
}

export function getBreadcrumbs(menus: Menu[], path: string) {
  const trail = findMenuTrail(menus, path)
  return trail.map((menu) => ({
    path: menu.type === 2 ? normalizePath(menu.path) : undefined,
    title: menu.name,
  }))
}
