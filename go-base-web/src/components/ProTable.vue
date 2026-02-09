<template>
  <div class="pro-table">
    <!-- 搜索区域 -->
    <a-card v-if="$slots.search && showSearch" class="search-card">
      <a-form layout="inline" @keyup.enter="handleSearch">
        <slot name="search" />
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">
              <SearchOutlined /> 搜索
            </a-button>
            <a-button @click="handleReset">
              <ReloadOutlined /> 重置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 表格区域 -->
    <a-card :class="{ 'table-card': $slots.search }">
      <template #title v-if="title">{{ title }}</template>
      <template #extra>
        <slot name="toolbar" />
      </template>
      
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="paginationConfig"
        :row-key="rowKey"
        :row-selection="rowSelection"
        :scroll="scroll"
        @change="handleTableChange"
      >
        <!-- 透传所有插槽 -->
        <template #bodyCell="scope">
          <slot name="bodyCell" v-bind="scope" />
        </template>
        <template #headerCell="scope" v-if="$slots.headerCell">
          <slot name="headerCell" v-bind="scope" />
        </template>
        <template #expandedRowRender="scope" v-if="$slots.expandedRowRender">
          <slot name="expandedRowRender" v-bind="scope" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import type { TableProps } from 'ant-design-vue'

export interface ProTableProps {
  title?: string
  columns: any[]
  dataSource: any[]
  loading?: boolean
  rowKey?: string | ((record: any) => string)
  rowSelection?: TableProps['rowSelection']
  scroll?: TableProps['scroll']
  // 分页相关
  pagination?: {
    current: number
    pageSize: number
    total: number
  }
  pageSizeOptions?: string[]
  showSizeChanger?: boolean
  showQuickJumper?: boolean
  showTotal?: boolean
  // 搜索区域控制
  showSearch?: boolean
}

const props = withDefaults(defineProps<ProTableProps>(), {
  rowKey: 'id',
  loading: false,
  pageSizeOptions: () => ['10', '20', '50', '100'],
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: true,
  showSearch: true,
})

const emit = defineEmits<{
  (e: 'search'): void
  (e: 'reset'): void
  (e: 'change', pagination: any, filters: any, sorter: any): void
  (e: 'update:pagination', pagination: any): void
}>()

// 分页配置
const paginationConfig = computed(() => {
  if (!props.pagination) return false
  
  return {
    current: props.pagination.current,
    pageSize: props.pagination.pageSize,
    total: props.pagination.total,
    showSizeChanger: props.showSizeChanger,
    showQuickJumper: props.showQuickJumper,
    pageSizeOptions: props.pageSizeOptions,
    showTotal: props.showTotal 
      ? (total: number, range: [number, number]) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
      : undefined,
  }
})

// 搜索
const handleSearch = () => {
  emit('search')
}

// 重置
const handleReset = () => {
  emit('reset')
}

// 表格变化（分页、排序、筛选）
const handleTableChange = (pagination: any, filters: any, sorter: any) => {
  if (props.pagination) {
    emit('update:pagination', {
      current: pagination.current,
      pageSize: pagination.pageSize,
      total: props.pagination.total
    })
  }
  emit('change', pagination, filters, sorter)
}
</script>

<style scoped>
.pro-table .search-card {
  margin-bottom: 16px;
}

.pro-table .table-card {
  margin-top: 0;
}

.pro-table :deep(.ant-form-item) {
  margin-bottom: 16px;
}

.pro-table :deep(.ant-card-body) {
  padding: 16px 24px;
}
</style>
