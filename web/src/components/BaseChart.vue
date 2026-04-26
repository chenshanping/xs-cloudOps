<template>
  <div ref="chartContainer" class="base-chart" :style="containerStyle">
    <v-chart
      v-if="!loading"
      ref="chartRef"
      :option="mergedOption"
      :autoresize="true"
      :loading="loading"
      :loading-options="loadingOptions"
      @click="handleClick"
    />
    <a-spin v-else class="chart-loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, PropType } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import {
  PieChart,
  BarChart,
  LineChart,
  RadarChart,
  ScatterChart,
  GaugeChart
} from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  ToolboxComponent
} from 'echarts/components'
import type { EChartsOption, ComposeOption } from 'echarts/core'
import type {
  PieSeriesOption,
  BarSeriesOption,
  LineSeriesOption,
  RadarSeriesOption,
  ScatterSeriesOption,
  GaugeSeriesOption
} from 'echarts/charts'
import type {
  TitleComponentOption,
  TooltipComponentOption,
  LegendComponentOption,
  GridComponentOption,
  DatasetComponentOption,
  ToolboxComponentOption
} from 'echarts/components'

// 注册 ECharts 组件
use([
  CanvasRenderer,
  PieChart,
  BarChart,
  LineChart,
  RadarChart,
  ScatterChart,
  GaugeChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  ToolboxComponent
])

// 定义组合类型
type ECOption = ComposeOption<
  | PieSeriesOption
  | BarSeriesOption
  | LineSeriesOption
  | RadarSeriesOption
  | ScatterSeriesOption
  | GaugeSeriesOption
  | TitleComponentOption
  | TooltipComponentOption
  | LegendComponentOption
  | GridComponentOption
  | DatasetComponentOption
  | ToolboxComponentOption
>

// 图表类型枚举
export type ChartType = 'pie' | 'bar' | 'line' | 'radar' | 'scatter' | 'gauge'

// Props 定义
const props = defineProps({
  // 图表类型
  type: {
    type: String as PropType<ChartType>,
    default: 'bar'
  },
  // 图表数据
  data: {
    type: Array as PropType<any[]>,
    default: () => []
  },
  // 图表标题
  title: {
    type: String,
    default: ''
  },
  // 自定义配置（会与默认配置合并）
  option: {
    type: Object as PropType<EChartsOption>,
    default: () => ({})
  },
  // 宽度
  width: {
    type: String,
    default: '100%'
  },
  // 高度
  height: {
    type: String,
    default: '400px'
  },
  // 加载状态
  loading: {
    type: Boolean,
    default: false
  },
  // 主题
  theme: {
    type: String,
    default: ''
  },
  // X轴字段名（用于 bar/line）
  xField: {
    type: String,
    default: 'name'
  },
  // Y轴字段名（用于 bar/line）
  yField: {
    type: String,
    default: 'value'
  },
  // 名称字段（用于 pie）
  nameField: {
    type: String,
    default: 'name'
  },
  // 值字段（用于 pie）
  valueField: {
    type: String,
    default: 'value'
  },
  // 颜色列表
  colors: {
    type: Array as PropType<string[]>,
    default: () => ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc']
  }
})

// Emits
const emit = defineEmits(['click', 'legendselectchanged'])

// Refs
const chartContainer = ref<HTMLElement>()
const chartRef = ref<InstanceType<typeof VChart>>()

// 加载配置
const loadingOptions = {
  text: '加载中...',
  color: '#1890ff',
  maskColor: 'rgba(255, 255, 255, 0.8)'
}

// 容器样式
const containerStyle = computed(() => ({
  width: props.width,
  height: props.height
}))

// 根据图表类型生成默认配置
const defaultOption = computed<ECOption>(() => {
  const baseOption: ECOption = {
    color: props.colors,
    title: props.title ? {
      text: props.title,
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold'
      }
    } : undefined,
    tooltip: {
      trigger: props.type === 'pie' ? 'item' : 'axis',
      formatter: props.type === 'pie' ? '{b}: {c} ({d}%)' : undefined
    },
    legend: {
      show: true,
      bottom: 10
    }
  }

  // 根据类型生成特定配置
  switch (props.type) {
    case 'pie':
      return {
        ...baseOption,
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['50%', '50%'],
          avoidLabelOverlap: true,
          itemStyle: {
            borderRadius: 4,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}: {d}%'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold'
            },
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          data: props.data.map(item => ({
            name: item[props.nameField],
            value: item[props.valueField]
          }))
        }]
      }

    case 'bar':
      return {
        ...baseOption,
        grid: {
          left: '3%',
          right: '4%',
          bottom: '15%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: props.data.map(item => item[props.xField]),
          axisLabel: {
            rotate: props.data.length > 6 ? 45 : 0
          }
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          type: 'bar',
          data: props.data.map(item => item[props.yField]),
          barWidth: '50%',
          itemStyle: {
            borderRadius: [4, 4, 0, 0]
          }
        }]
      }

    case 'line':
      return {
        ...baseOption,
        grid: {
          left: '3%',
          right: '4%',
          bottom: '15%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: props.data.map(item => item[props.xField]),
          boundaryGap: false
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          type: 'line',
          data: props.data.map(item => item[props.yField]),
          smooth: true,
          areaStyle: {
            opacity: 0.3
          }
        }]
      }

    case 'gauge':
      return {
        ...baseOption,
        series: [{
          type: 'gauge',
          progress: {
            show: true,
            width: 18
          },
          axisLine: {
            lineStyle: {
              width: 18
            }
          },
          axisTick: {
            show: false
          },
          splitLine: {
            length: 15,
            lineStyle: {
              width: 2,
              color: '#999'
            }
          },
          axisLabel: {
            distance: 25,
            color: '#999',
            fontSize: 12
          },
          anchor: {
            show: true,
            showAbove: true,
            size: 25,
            itemStyle: {
              borderWidth: 10
            }
          },
          title: {
            show: false
          },
          detail: {
            valueAnimation: true,
            fontSize: 28,
            offsetCenter: [0, '70%'],
            formatter: '{value}%'
          },
          data: props.data
        }]
      }

    case 'radar':
      return {
        ...baseOption,
        radar: {
          indicator: props.data.map(item => ({
            name: item[props.nameField],
            max: item.max || 100
          }))
        },
        series: [{
          type: 'radar',
          data: [{
            value: props.data.map(item => item[props.valueField]),
            name: props.title || '数据'
          }],
          areaStyle: {
            opacity: 0.3
          }
        }]
      }

    case 'scatter':
      return {
        ...baseOption,
        grid: {
          left: '3%',
          right: '4%',
          bottom: '15%',
          containLabel: true
        },
        xAxis: {
          type: 'value'
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          type: 'scatter',
          symbolSize: 10,
          data: props.data.map(item => [item[props.xField], item[props.yField]])
        }]
      }

    default:
      return baseOption
  }
})

// 合并用户配置
const mergedOption = computed<ECOption>(() => {
  return deepMerge(defaultOption.value, props.option) as ECOption
})

// 深度合并对象
function deepMerge(target: any, source: any): any {
  if (!source) return target
  const result = { ...target }
  for (const key in source) {
    if (source[key] !== null && typeof source[key] === 'object' && !Array.isArray(source[key])) {
      result[key] = deepMerge(target[key] || {}, source[key])
    } else {
      result[key] = source[key]
    }
  }
  return result
}

// 点击事件
function handleClick(params: any) {
  emit('click', params)
}

// 暴露方法
defineExpose({
  // 获取 ECharts 实例
  getChart: () => chartRef.value,
  // 刷新图表
  refresh: () => {
    chartRef.value?.setOption(mergedOption.value, true)
  },
  // 重置大小
  resize: () => {
    chartRef.value?.resize()
  }
})
</script>

<style scoped>
.base-chart {
  position: relative;
  min-height: 200px;
}

.base-chart :deep(.vue-echarts) {
  width: 100% !important;
  height: 100% !important;
}

.chart-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
