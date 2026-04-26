<template>
  <a-drawer
    :open="uiStore.settingsOpen"
    title="布局与主题"
    placement="right"
    width="360"
    @close="uiStore.toggleSettings(false)"
  >
    <div class="settings-section">
      <div class="section-title">布局模式</div>
      <a-radio-group
        :value="uiStore.layout.mode"
        button-style="solid"
        class="full-width"
        @update:value="handleLayoutModeChange"
      >
        <a-radio-button value="sidebar">侧边布局</a-radio-button>
        <a-radio-button value="top">顶部布局</a-radio-button>
        <a-radio-button value="mixed">混合布局</a-radio-button>
      </a-radio-group>
    </div>

    <div class="settings-section">
      <div class="section-title">显示控制</div>
      <div class="setting-item">
        <span>显示顶栏</span>
        <a-switch
          :checked="uiStore.layout.showHeader"
          :disabled="uiStore.layout.mode !== 'sidebar'"
          @update:checked="(checked) => uiStore.updateLayout({ showHeader: checked })"
        />
      </div>
      <div class="setting-item">
        <span>显示侧栏</span>
        <a-switch
          :checked="uiStore.layout.showSidebar"
          :disabled="uiStore.layout.mode === 'top'"
          @update:checked="(checked) => uiStore.updateLayout({ showSidebar: checked })"
        />
      </div>
      <div class="setting-item">
        <span>显示标签栏</span>
        <a-switch
          :checked="uiStore.layout.showTabs"
          @update:checked="(checked) => uiStore.updateLayout({ showTabs: checked })"
        />
      </div>
      <div class="setting-item">
        <span>紧凑内容区</span>
        <a-switch
          :checked="uiStore.theme.compactContent"
          @update:checked="(checked) => uiStore.updateTheme({ compactContent: checked })"
        />
      </div>
    </div>

    <div class="settings-section">
      <div class="section-title">主题</div>
      <a-radio-group
        :value="uiStore.theme.mode"
        button-style="solid"
        class="full-width"
        @update:value="(value) => uiStore.updateTheme({ mode: value as ThemeMode })"
      >
        <a-radio-button value="light">浅色</a-radio-button>
        <a-radio-button value="dark">深色</a-radio-button>
      </a-radio-group>
      <div class="preset-label">主题色</div>
      <div class="preset-grid">
        <button
          v-for="preset in themePresets"
          :key="preset.key"
          :class="['preset-card', { 'preset-card-active': uiStore.theme.preset === preset.key }]"
          type="button"
          @click="handlePresetSelect(preset.key)"
        >
          <template v-if="preset.key !== 'custom'">
            <span :style="{ background: preset.color }" class="preset-swatch"></span>
          </template>
          <template v-else>
            <span :style="{ background: uiStore.theme.primaryColor }" class="preset-swatch preset-swatch-custom">
              <BgColorsOutlined />
            </span>
            <input
              ref="customColorInput"
              :value="uiStore.theme.primaryColor"
              class="preset-color-input"
              type="color"
              @click.stop
              @input="handleColorInput"
            />
          </template>
          <span class="preset-name">{{ preset.label }}</span>
          <CheckOutlined v-if="uiStore.theme.preset === preset.key" class="preset-check" />
        </button>
      </div>
      <div class="setting-item">
        <span>深色顶栏</span>
        <a-switch
          :checked="uiStore.theme.headerDark"
          :disabled="uiStore.isDark"
          @update:checked="(checked) => uiStore.updateTheme({ headerDark: checked })"
        />
      </div>
      <div class="setting-item">
        <span>深色侧栏</span>
        <a-switch
          :checked="uiStore.theme.sidebarDark"
          :disabled="uiStore.isDark || uiStore.layout.mode === 'top'"
          @update:checked="(checked) => uiStore.updateTheme({ sidebarDark: checked })"
        />
      </div>
    </div>

    <div class="settings-section">
      <div class="section-title">侧栏</div>
      <div class="setting-item">
        <span>默认折叠</span>
        <a-switch
          :checked="uiStore.layout.sidebarCollapsed"
          :disabled="!uiStore.layout.showSidebar || uiStore.layout.mode === 'top'"
          @update:checked="(checked) => uiStore.updateLayout({ sidebarCollapsed: checked })"
        />
      </div>
      <div class="slider-row">
        <span>宽度</span>
        <a-slider
          :value="uiStore.layout.sidebarWidth"
          :min="180"
          :max="300"
          :disabled="!uiStore.layout.showSidebar || uiStore.layout.mode === 'top'"
          @update:value="(value) => uiStore.updateLayout({ sidebarWidth: Number(value) })"
        />
      </div>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <a-button @click="uiStore.resetPreferences()">恢复默认</a-button>
        <a-button type="primary" @click="uiStore.toggleSettings(false)">完成</a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { BgColorsOutlined, CheckOutlined } from '@ant-design/icons-vue'
import { THEME_PRESETS, useUiStore, type ThemeMode, type LayoutMode, type ThemePreset } from '@/store/ui'

const uiStore = useUiStore()
const themePresets = THEME_PRESETS
const customColorInput = ref<HTMLInputElement | null>(null)

const handleColorInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  uiStore.setCustomPrimaryColor(target.value)
}

const handlePresetSelect = (preset: ThemePreset) => {
  if (preset === 'custom') {
    uiStore.selectThemePreset('custom')
    customColorInput.value?.click()
    return
  }

  uiStore.selectThemePreset(preset)
}

const handleLayoutModeChange = (value: LayoutMode) => {
  if (value === 'top') {
    uiStore.updateLayout({
      mode: value,
      showSidebar: false
    })
    return
  }

  uiStore.updateLayout({
    mode: value,
    showSidebar: true
  })
}
</script>

<style scoped>
.settings-section + .settings-section {
  margin-top: 28px;
}

.section-title {
  margin-bottom: 12px;
  color: #111827;
  font-size: 14px;
  font-weight: 600;
}

.full-width {
  width: 100%;
}

.preset-label {
  margin: 18px 0 10px;
  color: #4b5563;
  font-size: 13px;
}

.setting-item,
.slider-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 10px 0;
}

.slider-row :deep(.ant-slider) {
  flex: 1;
}

.preset-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.preset-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px 8px 10px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #ffffff;
  cursor: pointer;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.preset-card:hover {
  border-color: rgba(0, 107, 230, 0.28);
  transform: translateY(-1px);
}

.preset-card-active {
  border-color: var(--app-primary-color, #006be6);
  box-shadow: 0 0 0 3px var(--app-primary-color-soft, rgba(0, 107, 230, 0.12));
}

.preset-swatch {
  width: 26px;
  height: 26px;
  border-radius: 8px;
  box-shadow: inset 0 0 0 1px rgba(15, 23, 42, 0.08);
}

.preset-swatch-custom {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 14px;
}

.preset-name {
  color: #374151;
  font-size: 12px;
  line-height: 1.2;
}

.preset-check {
  position: absolute;
  top: 8px;
  right: 8px;
  color: var(--app-primary-color, #006be6);
  font-size: 12px;
}

.preset-color-input {
  position: absolute;
  inset: 0;
  opacity: 0;
  pointer-events: none;
}

.drawer-footer {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
</style>
