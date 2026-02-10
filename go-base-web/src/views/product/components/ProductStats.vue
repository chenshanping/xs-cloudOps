<template>
  <div class="stats-container">
    <a-row :gutter="16">
      <a-col :span="12" >
        <a-card title="产品类型" size="small">
          <a-spin :spinning="chartLoading['type_id']">
            <template v-if="chartData['type_id']?.length > 0">
              <BaseChart
                type="pie"
                :data="chartData['type_id']"
                :loading="false"
                height="300px"
                name-field="name"
                value-field="value"
              />
            </template>
            <a-empty v-else description="暂无数据" style="height: 300px; display: flex; flex-direction: column; justify-content: center" />
          </a-spin>
        </a-card>
      </a-col>
      <a-col :span="12" >
        <a-card title="产品状态" size="small">
          <a-spin :spinning="chartLoading['status']">
            <template v-if="chartData['status']?.length > 0">
              <BaseChart
                type="bar"
                :data="chartData['status']"
                :loading="false"
                height="300px"
                x-field="name"
                y-field="value"
              />
            </template>
            <a-empty v-else description="暂无数据" style="height: 300px; display: flex; flex-direction: column; justify-content: center" />
          </a-spin>
        </a-card>
      </a-col>
      <a-col :span="24" style="margin-top: 16px">
        <a-card size="small">
          <template #title>
            <a-space>
              <span>产品信息趋势统计</span>
              <a-radio-group v-model:value="trendDays" size="small" @change="fetchTrendStats">
                <a-radio-button :value="7">近7天</a-radio-button>
                <a-radio-button :value="30">近30天</a-radio-button>
                <a-radio-button :value="90">近90天</a-radio-button>
              </a-radio-group>
            </a-space>
          </template>
          <a-spin :spinning="trendLoading">
            <template v-if="trendData?.length > 0">
              <BaseChart
                type="line"
                :data="trendData"
                :loading="false"
                height="300px"
                x-field="date"
                y-field="value"
              />
            </template>
            <a-empty v-else description="暂无数据" style="height: 300px; display: flex; flex-direction: column; justify-content: center" />
          </a-spin>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import BaseChart from '@/components/BaseChart.vue'
import { getProductStatsTypeId, getProductStatsStatus, getProductTrendStats } from '@/api/product'
import { getProductTypeOptions } from '@/api/productType'
import { getDictDataByType } from '@/api/dict'

// 各分组图表数据
const chartData = reactive<Record<string, { name: string; value: number }[]>>({
  'type_id': [],
  'status': [],
})
const chartLoading = reactive<Record<string, boolean>>({
  'type_id': false,
  'status': false,
})

// 趋势统计数据
const trendData = ref<{ date: string; value: number }[]>([])
const trendLoading = ref(false)
const trendDays = ref(30)

// 分组名称映射（用于外键字段和字典字段）
const nameMap = reactive<Record<string, Record<string, string>>>({})

// 获取分组名称映射
const fetchNameMaps = async () => {
  // 获取产品类型关联表名称
  try {
    const res = await getProductTypeOptions({ display_field: 'name' })
    const map: Record<string, string> = {}
    res.data?.forEach((item: any) => {
      map[item.id] = item.name
    })
    nameMap['type_id'] = map
  } catch (e) {
    console.error('获取type_id分组名称失败', e)
  }
  // 获取产品状态字典名称
  try {
    const res = await getDictDataByType('common_status')
    const map: Record<string, string> = {}
    res.data?.forEach((item: any) => {
      map[item.value] = item.label
    })
    nameMap['status'] = map
  } catch (e) {
    console.error('获取status字典名称失败', e)
  }
}

// 获取产品类型统计数据
const fetchStatsTypeId = async () => {
  chartLoading['type_id'] = true
  try {
    const res = await getProductStatsTypeId()
    const data = res.data || []
    const map = nameMap['type_id'] || {}
    chartData['type_id'] = data.map((item: any) => ({
      name: map[item.group_key] || item.name || `ID:${item.group_key}`,
      value: Number(item.value)
    }))
  } catch (e) {
    console.error('获取产品类型统计失败', e)
  } finally {
    chartLoading['type_id'] = false
  }
}

// 获取产品状态统计数据
const fetchStatsStatus = async () => {
  chartLoading['status'] = true
  try {
    const res = await getProductStatsStatus()
    const data = res.data || []
    const map = nameMap['status'] || {}
    chartData['status'] = data.map((item: any) => ({
      name: map[item.group_key] || item.name || `ID:${item.group_key}`,
      value: Number(item.value)
    }))
  } catch (e) {
    console.error('获取产品状态统计失败', e)
  } finally {
    chartLoading['status'] = false
  }
}

// 获取趋势统计数据
const fetchTrendStats = async () => {
  trendLoading.value = true
  try {
    const res = await getProductTrendStats(trendDays.value)
    // 格式化日期显示
    trendData.value = (res.data || []).map((item: any) => ({
      ...item,
      date: item.date ? item.date.substring(5, 10).replace('-', '/') : item.date
    }))
  } catch (e) {
    console.error('获取趋势统计失败', e)
  } finally {
    trendLoading.value = false
  }
}

// 暴露刷新方法供父组件调用
defineExpose({
  refresh: async () => {
    await fetchNameMaps()
    fetchStatsTypeId()
    fetchStatsStatus()
    fetchTrendStats()
  }
})
</script>

<style scoped>
.stats-container {
  margin-bottom: 16px;
}
</style>
