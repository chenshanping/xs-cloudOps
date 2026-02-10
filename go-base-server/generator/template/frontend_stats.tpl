<template>
  <div class="stats-container">
    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="{{.Description}}分布统计" size="small">
          <BaseChart
            type="pie"
            :data="groupData"
            :loading="loading"
            height="300px"
            name-field="name"
            value-field="value"
          />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="{{.Description}}数量统计" size="small">
          <BaseChart
            type="bar"
            :data="groupData"
            :loading="loading"
            height="300px"
            x-field="name"
            y-field="value"
          />
        </a-card>
      </a-col>
{{- if .StatsTimeColumn}}
      <a-col :span="24" style="margin-top: 16px">
        <a-card size="small">
          <template #title>
            <a-space>
              <span>{{.Description}}趋势统计</span>
              <a-radio-group v-model:value="trendDays" size="small" @change="fetchTrendStats">
                <a-radio-button :value="7">近7天</a-radio-button>
                <a-radio-button :value="30">近30天</a-radio-button>
                <a-radio-button :value="90">近90天</a-radio-button>
              </a-radio-group>
            </a-space>
          </template>
          <BaseChart
            type="line"
            :data="trendData"
            :loading="trendLoading"
            height="300px"
            x-field="date"
            y-field="value"
          />
        </a-card>
      </a-col>
{{- end}}
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseChart from '@/components/BaseChart.vue'
import { get{{.ModelName}}GroupStats{{if .StatsTimeColumn}}, get{{.ModelName}}TrendStats{{end}} } from '@/api/{{.ModuleName}}'
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
{{- if $.StatsGroupColumn}}
{{- if eq $.StatsGroupColumn .ForeignKeyJson}}
import { get{{.RelatedModel}}Options } from '@/api/{{.RelatedModule}}'
{{- end}}
{{- end}}
{{- end}}
{{- end}}

// 分组统计数据
const groupData = ref<{ name: string; value: number }[]>([])
const loading = ref(false)
{{- if .StatsTimeColumn}}

// 趋势统计数据
const trendData = ref<{ date: string; value: number }[]>([])
const trendLoading = ref(false)
const trendDays = ref(30)
{{- end}}

// 分组名称映射
const groupNameMap = ref<Record<string, string>>({})

// 获取分组名称映射
const fetchGroupNameMap = async () => {
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
{{- if $.StatsGroupColumn}}
{{- if eq $.StatsGroupColumn .ForeignKeyJson}}
  try {
    const res = await get{{.RelatedModel}}Options({ display_field: '{{.DisplayField}}' })
    const map: Record<string, string> = {}
    res.data?.forEach((item: any) => {
      map[item.id] = item.name
    })
    groupNameMap.value = map
  } catch (e) {
    console.error('获取分组名称失败', e)
  }
{{- end}}
{{- end}}
{{- end}}
{{- end}}
}

// 获取分组统计数据
const fetchGroupStats = async () => {
  loading.value = true
  try {
    const res = await get{{.ModelName}}GroupStats()
    const data = res.data || []
    
    // 转换数据，使用名称映射
    groupData.value = data.map((item: any) => ({
      name: groupNameMap.value[item.group_key] || item.name || `ID:${item.group_key}`,
      value: Number(item.value)
    }))
  } catch (e) {
    console.error('获取分组统计失败', e)
  } finally {
    loading.value = false
  }
}
{{- if .StatsTimeColumn}}

// 获取趋势统计数据
const fetchTrendStats = async () => {
  trendLoading.value = true
  try {
    const res = await get{{.ModelName}}TrendStats(trendDays.value)
    trendData.value = res.data || []
  } catch (e) {
    console.error('获取趋势统计失败', e)
  } finally {
    trendLoading.value = false
  }
}
{{- end}}

onMounted(async () => {
  await fetchGroupNameMap()
  fetchGroupStats()
{{- if .StatsTimeColumn}}
  fetchTrendStats()
{{- end}}
})

// 暴露刷新方法供父组件调用
defineExpose({
  refresh: () => {
    fetchGroupStats()
{{- if .StatsTimeColumn}}
    fetchTrendStats()
{{- end}}
  }
})
</script>

<style scoped>
.stats-container {
  margin-bottom: 16px;
}
</style>
