export interface TabItem {
  key: string
  path: string
  fullPath: string
  title: string
  name?: string
  affix: boolean
}

export interface TabsStateSnapshot {
  tabs: TabItem[]
  activeKey: string
}

export const HOME_TAB: TabItem = {
  key: '/dashboard',
  path: '/dashboard',
  fullPath: '/dashboard',
  title: '首页',
  name: 'Dashboard',
  affix: true,
}

function uniqueTabs(tabs: TabItem[]) {
  const seen = new Set<string>()
  const result: TabItem[] = []

  for (const tab of tabs) {
    if (!tab.key || seen.has(tab.key)) {
      continue
    }
    seen.add(tab.key)
    result.push(tab)
  }

  return result
}

function ensureHomeTab(tabs: TabItem[]) {
  const merged = uniqueTabs([HOME_TAB, ...tabs])
  const homeIndex = merged.findIndex((tab) => tab.key === HOME_TAB.key)

  if (homeIndex > 0) {
    const [home] = merged.splice(homeIndex, 1)
    merged.unshift(home)
  }

  return merged
}

function resolveActiveKey(tabs: TabItem[], preferredKey?: string) {
  if (preferredKey && tabs.some((tab) => tab.key === preferredKey)) {
    return preferredKey
  }
  return tabs[tabs.length - 1]?.key || HOME_TAB.key
}

export function createInitialTabsState(tabs: TabItem[], activeKey?: string): TabsStateSnapshot {
  const normalizedTabs = ensureHomeTab(tabs)
  return {
    tabs: normalizedTabs,
    activeKey: resolveActiveKey(normalizedTabs, activeKey),
  }
}

export function upsertTab(state: TabsStateSnapshot, tab: TabItem): TabsStateSnapshot {
  const nextTabs = [...state.tabs]
  const index = nextTabs.findIndex((item) => item.key === tab.key)

  if (index >= 0) {
    nextTabs[index] = { ...nextTabs[index], ...tab }
  } else {
    nextTabs.push(tab)
  }

  return createInitialTabsState(nextTabs, tab.key)
}

export function closeTab(state: TabsStateSnapshot, targetKey: string): TabsStateSnapshot {
  const currentIndex = state.tabs.findIndex((tab) => tab.key === targetKey)
  if (currentIndex === -1) {
    return state
  }

  const target = state.tabs[currentIndex]
  if (target.affix) {
    return state
  }

  const nextTabs = state.tabs.filter((tab) => tab.key !== targetKey)
  const fallbackKey = state.activeKey === targetKey
    ? nextTabs[currentIndex]?.key || nextTabs[currentIndex - 1]?.key || HOME_TAB.key
    : state.activeKey

  return createInitialTabsState(nextTabs, fallbackKey)
}

export function closeLeftTabs(state: TabsStateSnapshot, targetKey: string): TabsStateSnapshot {
  const currentIndex = state.tabs.findIndex((tab) => tab.key === targetKey)
  if (currentIndex === -1) {
    return state
  }

  const nextTabs = state.tabs.filter((tab, index) => tab.affix || index >= currentIndex)
  return createInitialTabsState(nextTabs, resolveActiveKey(nextTabs, targetKey))
}

export function closeRightTabs(state: TabsStateSnapshot, targetKey: string): TabsStateSnapshot {
  const currentIndex = state.tabs.findIndex((tab) => tab.key === targetKey)
  if (currentIndex === -1) {
    return state
  }

  const nextTabs = state.tabs.filter((tab, index) => tab.affix || index <= currentIndex)
  return createInitialTabsState(nextTabs, resolveActiveKey(nextTabs, targetKey))
}

export function closeOtherTabs(state: TabsStateSnapshot, targetKey: string): TabsStateSnapshot {
  const nextTabs = state.tabs.filter((tab) => tab.affix || tab.key === targetKey)
  return createInitialTabsState(nextTabs, resolveActiveKey(nextTabs, targetKey))
}

export function closeAllTabs(): TabsStateSnapshot {
  return createInitialTabsState([HOME_TAB], HOME_TAB.key)
}
