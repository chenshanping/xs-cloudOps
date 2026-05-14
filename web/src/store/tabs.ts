import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import {
  closeAllTabs,
  closeLeftTabs,
  closeOtherTabs,
  closeRightTabs,
  closeTab,
  createInitialTabsState,
  type TabItem,
  type TabsStateSnapshot,
  upsertTab,
} from './tabs-state'

const STORAGE_KEY = 'go-base-tabs'

function loadTabsState(): TabsStateSnapshot {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return createInitialTabsState([])
    }
    const parsed = JSON.parse(raw) as Partial<TabsStateSnapshot>
    const sanitizedTabs = (parsed.tabs || []).map((tab) => ({ ...tab, affix: false }))
    return createInitialTabsState(sanitizedTabs, parsed.activeKey)
  } catch (error) {
    console.error('读取标签页状态失败', error)
    return createInitialTabsState([])
  }
}

export const useTabsStore = defineStore('tabs', () => {
  const saved = loadTabsState()
  const tabs = ref<TabItem[]>(saved.tabs)
  const activeKey = ref(saved.activeKey)

  const cachedViews = computed(() => {
    return Array.from(
      new Set(
        tabs.value
          .map((tab) => tab.name)
          .filter((name): name is string => Boolean(name)),
      ),
    )
  })

  function applySnapshot(snapshot: TabsStateSnapshot) {
    tabs.value = snapshot.tabs
    activeKey.value = snapshot.activeKey
  }

  function persist() {
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        tabs: tabs.value,
        activeKey: activeKey.value,
      }),
    )
  }

  watch([tabs, activeKey], persist, { deep: true })

  function openTab(tab: TabItem) {
    applySnapshot(upsertTab({ tabs: tabs.value, activeKey: activeKey.value }, tab))
  }

  function setActiveKey(key: string) {
    activeKey.value = key
  }

  function removeTab(targetKey: string) {
    const previousActiveKey = activeKey.value
    const snapshot = closeTab({ tabs: tabs.value, activeKey: activeKey.value }, targetKey)
    applySnapshot(snapshot)
    return previousActiveKey === targetKey ? snapshot.activeKey : null
  }

  function removeLeftTabs(targetKey: string) {
    applySnapshot(closeLeftTabs({ tabs: tabs.value, activeKey: activeKey.value }, targetKey))
  }

  function removeRightTabs(targetKey: string) {
    applySnapshot(closeRightTabs({ tabs: tabs.value, activeKey: activeKey.value }, targetKey))
  }

  function removeOtherTabs(targetKey: string) {
    applySnapshot(closeOtherTabs({ tabs: tabs.value, activeKey: activeKey.value }, targetKey))
  }

  function removeAllTabs() {
    applySnapshot(closeAllTabs())
  }

  return {
    tabs,
    activeKey,
    cachedViews,
    openTab,
    setActiveKey,
    removeTab,
    removeLeftTabs,
    removeRightTabs,
    removeOtherTabs,
    removeAllTabs,
  }
})
