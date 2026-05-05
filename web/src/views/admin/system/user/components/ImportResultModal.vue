<template>
  <a-modal
    :open="open"
    title="导入结果"
    :footer="null"
    width="720"
    @cancel="handleClose"
  >
    <a-result
      v-if="result"
      :status="result.failed_count === 0 ? 'success' : (result.success_count === 0 ? 'error' : 'warning')"
      :title="resultTitle"
      :sub-title="resultSubTitle"
    >
      <template #extra>
        <a-button type="primary" @click="handleClose">确定</a-button>
      </template>
    </a-result>

    <a-table
      v-if="result && result.errors && result.errors.length > 0"
      :columns="errorColumns"
      :data-source="result.errors"
      :pagination="{ pageSize: 10, showSizeChanger: true, showTotal: (total: number) => `共 ${total} 条错误` }"
      row-key="(record: any, index: number) => index"
      size="small"
      bordered
      style="margin-top: -16px"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'row'">
          <a-tag v-if="record.row > 0" color="blue">第{{ record.row }}行</a-tag>
          <a-tag v-else color="default">-</a-tag>
        </template>
        <template v-if="column.key === 'message'">
          <span style="color: #cf1322">{{ record.message }}</span>
        </template>
      </template>
    </a-table>
  </a-modal>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ImportResult } from '@/api/user'

const props = defineProps<{
  open: boolean
  result: ImportResult | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const resultTitle = computed(() => {
  if (!props.result) return ''
  if (props.result.failed_count === 0) return '导入成功'
  if (props.result.success_count === 0) return '导入失败'
  return '部分导入成功'
})

const resultSubTitle = computed(() => {
  if (!props.result) return ''
  return `共 ${props.result.total_count} 条数据，成功 ${props.result.success_count} 条，失败 ${props.result.failed_count} 条`
})

const errorColumns = [
  { title: '行号', key: 'row', dataIndex: 'row', width: 80 },
  { title: '列名', dataIndex: 'column', key: 'column', width: 100 },
  { title: '原始值', dataIndex: 'value', key: 'value', width: 120, ellipsis: true },
  { title: '错误信息', dataIndex: 'message', key: 'message' }
]

const handleClose = () => {
  emit('update:open', false)
}
</script>
