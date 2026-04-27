<template>
  <a-drawer
    :open="uiStore.settingsOpen"
    title="布局设置"
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
import { useUiStore, type LayoutMode } from '@/store/ui'

const uiStore = useUiStore()

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

.drawer-footer {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
</style>
