<template>
  <a-drawer
    :open="open"
    title="管理平台模型"
    width="920"
    placement="right"
    @close="handleClose"
  >
    <div class="drawer-toolbar">
      <div>
        <div class="provider-title">{{ providerName || '未命名平台' }}</div>
        <div class="provider-meta">{{ providerBaseUrl || '请先填写 Base URL' }}</div>
      </div>
      <a-space>
        <a-input-search
          v-model:value="keyword"
          placeholder="搜索模型 ID / 所属平台"
          allow-clear
          class="search-input"
        />
        <a-button type="primary" :loading="loading" @click="handleFetch">
          获取模型列表
        </a-button>
      </a-space>
    </div>

    <a-alert
      class="drawer-alert"
      type="info"
      show-icon
      message="导入只更新当前编辑态；仍需点击页面底部“保存配置”才会真正落库。"
    />

    <a-table
      :data-source="filteredModels"
      :columns="columns"
      :pagination="{ pageSize: 10, hideOnSinglePage: true }"
      :loading="loading"
      row-key="id"
      size="small"
      :row-selection="rowSelection"
      :locale="tableLocale"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag v-if="isImported(record.id)" color="success">已导入</a-tag>
          <a-tag v-else-if="selectedRowKeys.includes(record.id)" color="processing">待导入</a-tag>
          <a-tag v-else>未导入</a-tag>
        </template>
        <template v-else-if="column.key === 'created'">
          {{ formatCreated(record.created) }}
        </template>
      </template>
    </a-table>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">关闭</a-button>
        <a-button type="primary" :disabled="selectedRowKeys.length === 0" @click="handleImport">
          导入已选
        </a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import { fetchAIProviderModels } from '@/api/ai'
import type { AIModel, RemoteProviderModel } from './state'

const props = defineProps<{
  open: boolean
  providerName: string
  apiKey: string
  providerBaseUrl: string
  localModels: AIModel[]
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'import', value: RemoteProviderModel[]): void
}>()

const columns = [
  { title: '模型 ID', dataIndex: 'id', key: 'id', ellipsis: true },
  { title: '对象类型', dataIndex: 'object', key: 'object', width: 120 },
  { title: '所属平台', dataIndex: 'owned_by', key: 'owned_by', width: 180, ellipsis: true },
  { title: '创建时间', dataIndex: 'created', key: 'created', width: 180 },
  { title: '状态', key: 'status', width: 120 },
]

const loading = ref(false)
const keyword = ref('')
const remoteModels = ref<RemoteProviderModel[]>([])
const selectedRowKeys = ref<string[]>([])

const localModelIDs = computed(() => new Set(props.localModels.map(model => model.id)))

const filteredModels = computed(() => {
  const search = keyword.value.trim().toLowerCase()
  if (!search) {
    return remoteModels.value
  }
  return remoteModels.value.filter(model =>
    model.id.toLowerCase().includes(search)
    || (model.owned_by ?? '').toLowerCase().includes(search)
    || (model.object ?? '').toLowerCase().includes(search),
  )
})

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys: Array<string | number>) => {
    selectedRowKeys.value = keys.map(key => String(key))
  },
  getCheckboxProps: (record: RemoteProviderModel) => ({
    disabled: isImported(record.id),
  }),
}))

const tableLocale = computed(() => ({
  emptyText: loading.value ? '正在获取模型列表...' : '点击“获取模型列表”查看平台返回的模型',
}))

const isImported = (id: string) => localModelIDs.value.has(id)

const formatCreated = (created?: number) => {
  if (!created) {
    return '-'
  }
  return dayjs.unix(created).format('YYYY-MM-DD HH:mm:ss')
}

const handleClose = () => {
  emit('update:open', false)
}

const handleFetch = async () => {
  if (!props.apiKey.trim()) {
    message.warning('请先填写当前平台的 API Key')
    return
  }
  if (!props.providerBaseUrl.trim()) {
    message.warning('请先填写当前平台的 Base URL')
    return
  }

  loading.value = true
  try {
    const res = await fetchAIProviderModels({
      api_key: props.apiKey,
      base_url: props.providerBaseUrl,
    })
    remoteModels.value = res.data.models ?? []
    selectedRowKeys.value = []
    message.success(`已获取 ${remoteModels.value.length} 个平台模型`)
  } catch (error: any) {
    message.error(error.message || '获取模型列表失败')
  } finally {
    loading.value = false
  }
}

const handleImport = () => {
  const selectedModels = remoteModels.value.filter(model => selectedRowKeys.value.includes(model.id))
  if (selectedModels.length === 0) {
    message.warning('请先选择要导入的模型')
    return
  }

  emit('import', selectedModels)
  selectedRowKeys.value = []
}

watch(
  () => [props.open, props.providerName, props.providerBaseUrl],
  ([open]) => {
    if (!open) {
      return
    }
    keyword.value = ''
    remoteModels.value = []
    selectedRowKeys.value = []
  },
)
</script>

<style scoped>
.drawer-toolbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.provider-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f1f1f;
}

.provider-meta {
  margin-top: 4px;
  color: #8c8c8c;
  word-break: break-all;
}

.search-input {
  width: 260px;
}

.drawer-alert {
  margin-bottom: 16px;
}
</style>
