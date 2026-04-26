import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'

export type ThemeMode = 'light' | 'dark'
export type LayoutMode = 'sidebar' | 'top' | 'mixed'
export type ThemePreset =
  | 'default'
  | 'violet'
  | 'pink'
  | 'yellow'
  | 'sky-blue'
  | 'green'
  | 'zinc'
  | 'deep-green'
  | 'deep-blue'
  | 'orange'
  | 'rose'
  | 'neutral'
  | 'slate'
  | 'gray'
  | 'custom'

export interface ThemePresetOption {
  key: ThemePreset
  label: string
  color: string
  lightPrimaryColor?: string
  darkPrimaryColor?: string
}

export interface LayoutPreferences {
  mode: LayoutMode
  showHeader: boolean
  showSidebar: boolean
  showTabs: boolean
  sidebarCollapsed: boolean
  sidebarWidth: number
}

export interface ThemePreferences {
  mode: ThemeMode
  preset: ThemePreset
  primaryColor: string
  headerDark: boolean
  sidebarDark: boolean
  compactContent: boolean
}

export const THEME_PRESETS: ThemePresetOption[] = [
  { key: 'default', label: '默认', color: '#006be6' },
  { key: 'violet', label: '紫罗兰', color: '#7166f0' },
  { key: 'pink', label: '玫粉', color: '#e84a6c' },
  { key: 'yellow', label: '琥珀黄', color: '#efbd48' },
  { key: 'sky-blue', label: '天青蓝', color: '#4e69fd' },
  { key: 'green', label: '薄荷绿', color: '#0bd092' },
  {
    key: 'zinc',
    label: '锌灰',
    color: '#3f3f46',
    lightPrimaryColor: '#18181b',
    darkPrimaryColor: '#fafafa'
  },
  { key: 'deep-green', label: '深青绿', color: '#0d9496' },
  { key: 'deep-blue', label: '深海蓝', color: '#0960be' },
  { key: 'orange', label: '橙棕', color: '#c1420b' },
  { key: 'rose', label: '绛红', color: '#bb1b1b' },
  {
    key: 'neutral',
    label: '中性灰',
    color: '#404040',
    lightPrimaryColor: '#18181b',
    darkPrimaryColor: '#fafafa'
  },
  {
    key: 'slate',
    label: '石板灰',
    color: '#344256',
    lightPrimaryColor: '#18181b',
    darkPrimaryColor: '#fafafa'
  },
  {
    key: 'gray',
    label: '冷灰',
    color: '#384250',
    lightPrimaryColor: '#18181b',
    darkPrimaryColor: '#fafafa'
  },
  { key: 'custom', label: '自定义', color: '#1677ff' }
]

const STORAGE_KEY = 'go-base-ui-preferences'

const DEFAULT_LAYOUT: LayoutPreferences = {
  mode: 'sidebar',
  showHeader: true,
  showSidebar: true,
  showTabs: true,
  sidebarCollapsed: false,
  sidebarWidth: 220
}

const DEFAULT_THEME: ThemePreferences = {
  mode: 'light',
  preset: 'default',
  primaryColor: '#006be6',
  headerDark: false,
  sidebarDark: true,
  compactContent: false
}

function normalizeHexColor(color?: string) {
  return (color || '').trim().toLowerCase()
}

function getThemePreset(preset: ThemePreset) {
  return THEME_PRESETS.find((item) => item.key === preset)
}

function resolvePresetPrimaryColor(preset: ThemePreset, mode: ThemeMode) {
  const selected = getThemePreset(preset)
  if (!selected) {
    return DEFAULT_THEME.primaryColor
  }

  if (mode === 'dark') {
    return selected.darkPrimaryColor || selected.lightPrimaryColor || selected.color
  }

  return selected.lightPrimaryColor || selected.color
}

function inferThemePreset(primaryColor?: string): ThemePreset {
  const normalizedColor = normalizeHexColor(primaryColor)
  if (!normalizedColor) {
    return DEFAULT_THEME.preset
  }

  const matchedPreset = THEME_PRESETS.find((preset) => {
    return [
      preset.color,
      preset.lightPrimaryColor,
      preset.darkPrimaryColor,
    ]
      .filter(Boolean)
      .some((value) => normalizeHexColor(value) === normalizedColor)
  })

  return matchedPreset?.key || 'custom'
}

function normalizeThemePreferences(input: ThemePreferences): ThemePreferences {
  const next = {
    ...input,
    primaryColor: normalizeHexColor(input.primaryColor) || DEFAULT_THEME.primaryColor,
    preset: input.preset || inferThemePreset(input.primaryColor),
  }

  if (next.preset !== 'custom') {
    next.primaryColor = resolvePresetPrimaryColor(next.preset, next.mode)
  }

  return next
}

function loadPreferences() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return null
    }
    return JSON.parse(raw) as {
      layout?: Partial<LayoutPreferences>
      theme?: Partial<ThemePreferences>
    }
  } catch (error) {
    console.error('读取布局偏好失败', error)
    return null
  }
}

function normalizeLayoutPreferences(
  input: LayoutPreferences,
): LayoutPreferences {
  const next = { ...input }

  next.sidebarWidth = Math.min(300, Math.max(180, next.sidebarWidth))

  if (next.mode === 'top') {
    next.showHeader = true
    next.showSidebar = false
  }

  if (next.mode === 'mixed') {
    next.showHeader = true
  }

  return next
}

export const useUiStore = defineStore('ui', () => {
  const saved = loadPreferences()

  const layout = ref<LayoutPreferences>({
    ...normalizeLayoutPreferences({
      ...DEFAULT_LAYOUT,
      ...saved?.layout
    })
  })

  const theme = ref<ThemePreferences>({
    ...normalizeThemePreferences({
      ...DEFAULT_THEME,
      ...saved?.theme,
      preset: saved?.theme?.preset || inferThemePreset(saved?.theme?.primaryColor)
    })
  })

  const settingsOpen = ref(false)

  const effectiveShowSidebar = computed(() => {
    if (layout.value.mode === 'top') {
      return false
    }
    return layout.value.showSidebar
  })

  const effectiveShowHeader = computed(() => layout.value.showHeader)
  const effectiveShowTabs = computed(() => layout.value.showTabs)

  const isDark = computed(() => theme.value.mode === 'dark')
  const contentPadding = computed(() => (theme.value.compactContent ? 12 : 20))

  function persist() {
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        layout: layout.value,
        theme: theme.value
      })
    )
  }

  watch([layout, theme], persist, { deep: true })

  function updateLayout(patch: Partial<LayoutPreferences>) {
    layout.value = normalizeLayoutPreferences({
      ...layout.value,
      ...patch
    })
  }

  function updateTheme(patch: Partial<ThemePreferences>) {
    theme.value = normalizeThemePreferences({
      ...theme.value,
      ...patch
    })
  }

  function selectThemePreset(preset: ThemePreset) {
    theme.value = normalizeThemePreferences({
      ...theme.value,
      preset
    })
  }

  function setCustomPrimaryColor(color: string) {
    theme.value = normalizeThemePreferences({
      ...theme.value,
      preset: 'custom',
      primaryColor: color
    })
  }

  function resetPreferences() {
    layout.value = { ...DEFAULT_LAYOUT }
    theme.value = { ...DEFAULT_THEME }
  }

  function toggleSettings(open?: boolean) {
    settingsOpen.value = typeof open === 'boolean' ? open : !settingsOpen.value
  }

  return {
    layout,
    theme,
    settingsOpen,
    effectiveShowSidebar,
    effectiveShowHeader,
    effectiveShowTabs,
    isDark,
    contentPadding,
    updateLayout,
    updateTheme,
    selectThemePreset,
    setCustomPrimaryColor,
    toggleSettings,
    resetPreferences
  }
})
