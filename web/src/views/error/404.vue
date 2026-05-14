<template>
  <div class="not-found">
    <a-result status="404" title="404" sub-title="抱歉，您访问的页面不存在">
      <template #extra>
        <a-button type="primary" @click="router.push(backendHomePath)">返回首页</a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { filterVisibleMenus, firstNavigableMenuPath } from '@/layouts/components/layout-menu'

const router = useRouter()
const userStore = useUserStore()
const backendHomePath = computed(() => firstNavigableMenuPath(filterVisibleMenus(userStore.menus || [])) || '/no-permission')
</script>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
</style>
