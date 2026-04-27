<template>
  <div class="dashboard">
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card>
          <a-statistic title="用户总数" :value="112893" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="角色总数" :value="12" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="菜单总数" :value="28" />
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="API总数" :value="56" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 图表区域 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="12">
        <a-card title="用户角色占比" :loading="roleStatsLoading">
          <BaseChart
            type="pie"
            :data="roleStatsData"
            :loading="roleStatsLoading"
            height="350px"
          />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="用户状态分布" :loading="statusStatsLoading">
          <BaseChart
            type="pie"
            :data="statusStatsData"
            :loading="statusStatsLoading"
            height="350px"
            :colors="['#52c41a', '#ff4d4f']"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="24">
        <a-card title="用户注册趋势（近30天）" :loading="trendLoading">
          <BaseChart
            type="line"
            :data="trendData"
            :loading="trendLoading"
            height="350px"
            x-field="date"
            y-field="count"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-card title="欢迎使用" style="margin-top: 16px">
      <p>Go RBAC Admin 是一个基于 Go + Gin + Vue3 + TypeScript 的后台权限管理系统</p>
      <p>当前用户: {{ userStore.user?.nickname || userStore.user?.username }}</p>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import BaseChart from '@/components/BaseChart.vue'
import { getUserRoleStats, getUserStatusStats, getUserRegisterTrend, type ChartItem, type TrendItem } from '@/api/echart'

const userStore = useUserStore()

// 用户角色占比数据
const roleStatsData = ref<ChartItem[]>([])
const roleStatsLoading = ref(false)

// 用户状态统计数据
const statusStatsData = ref<ChartItem[]>([])
const statusStatsLoading = ref(false)

// 用户注册趋势数据
const trendData = ref<TrendItem[]>([])
const trendLoading = ref(false)

// 加载用户角色占比
async function loadRoleStats() {
  roleStatsLoading.value = true
  try {
    const res = await getUserRoleStats()
    roleStatsData.value = res.data || []
  } catch (error) {
    console.error('加载用户角色统计失败:', error)
  } finally {
    roleStatsLoading.value = false
  }
}

// 加载用户状态统计
async function loadStatusStats() {
  statusStatsLoading.value = true
  try {
    const res = await getUserStatusStats()
    statusStatsData.value = res.data || []
  } catch (error) {
    console.error('加载用户状态统计失败:', error)
  } finally {
    statusStatsLoading.value = false
  }
}

// 加载用户注册趋势
async function loadTrend() {
  trendLoading.value = true
  try {
    const res = await getUserRegisterTrend()
    trendData.value = res.data || []
  } catch (error) {
    console.error('加载用户注册趋势失败:', error)
  } finally {
    trendLoading.value = false
  }
}

onMounted(() => {
  loadRoleStats()
  loadStatusStats()
  loadTrend()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}
</style>
