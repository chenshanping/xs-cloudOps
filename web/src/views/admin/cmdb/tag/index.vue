<template>
  <PageWrapper class="cmdb-tag-page">
    <div class="cmdb-tag-page__content">
      <div class="summary-grid">
        <div class="summary-item">
          <TagsOutlined />
          <span>标签总数</span>
          <strong>{{ tableData.length }}</strong>
        </div>
        <div class="summary-item success">
          <BgColorsOutlined />
          <span>已配置颜色</span>
          <strong>{{ colorConfiguredCount }}</strong>
        </div>
      </div>

      <ProTable
        title="主机标签"
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        row-key="id"
        :scroll="{ x: 860 }"
        @search="handleSearch"
        @reset="handleReset"
      >
        <template #search>
          <a-form-item label="标签名称">
            <a-input v-model:value="searchForm.name" placeholder="请输入标签名称" allow-clear />
          </a-form-item>
        </template>

        <template #toolbar>
          <a-button type="primary" v-permission="'cmdb:tag:create'" @click="handleAdd">
            <PlusOutlined /> 新增标签
          </a-button>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-space :size="8">
              <a-tag :color="record.color || '#1677ff'">{{ record.name }}</a-tag>
              <span class="tag-code">{{ record.color || '-' }}</span>
            </a-space>
          </template>
          <template v-else-if="column.key === 'color'">
            <div class="color-display">
              <span class="color-display__block" :style="{ background: record.color || '#d9d9d9' }"></span>
              <span>{{ record.color || '-' }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space :size="0">
              <a-button type="link" size="small" v-permission="'cmdb:tag:update'" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" size="small" danger v-permission="'cmdb:tag:delete'" @click="handleDelete(record)">删除</a-button>
            </a-space>
          </template>
        </template>
      </ProTable>

      <TagFormDrawer
        v-model:open="drawerVisible"
        :title="drawerTitle"
        :initial-value="drawerInitialValue"
        @submit="handleSubmit"
      />
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { BgColorsOutlined, ExclamationCircleOutlined, PlusOutlined, TagsOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import ProTable from '@/components/ProTable.vue'
import { createCmdbHostTag, deleteCmdbHostTag, getCmdbHostTags, updateCmdbHostTag, type CmdbHostTag } from '@/api/cmdb'
import { useTableColumns } from '@/utils/permission'
import TagFormDrawer from './components/TagFormDrawer.vue'

interface TagDrawerValue {
  name: string
  color: string
  remark: string
}

const loading = ref(false)
const tableData = ref<CmdbHostTag[]>([])
const drawerVisible = ref(false)
const drawerTitle = ref('新增标签')
const editingId = ref<number>()
const drawerInitialValue = ref<Partial<TagDrawerValue>>({})

const searchForm = reactive({
  name: '',
})

const columns = useTableColumns(
  [
    { title: '标签名称', key: 'name', width: 220 },
    { title: '颜色值', key: 'color', width: 180 },
    { title: '备注', dataIndex: 'remark', key: 'remark', ellipsis: true },
  ],
  { title: '操作', key: 'action', width: 150, fixed: 'right', align: 'center' },
  ['cmdb:tag:update', 'cmdb:tag:delete']
)

const colorConfiguredCount = computed(() => tableData.value.filter(item => !!item.color).length)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getCmdbHostTags({ name: searchForm.name || undefined })
    tableData.value = res.data || []
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchData()
}

const handleReset = () => {
  searchForm.name = ''
  fetchData()
}

const handleAdd = () => {
  editingId.value = undefined
  drawerTitle.value = '新增标签'
  drawerInitialValue.value = { name: '', color: '#1677ff', remark: '' }
  drawerVisible.value = true
}

const handleEdit = (record: CmdbHostTag) => {
  editingId.value = record.id
  drawerTitle.value = '编辑标签'
  drawerInitialValue.value = {
    name: record.name,
    color: record.color || '#1677ff',
    remark: record.remark,
  }
  drawerVisible.value = true
}

const handleSubmit = async (values: { name: string; color: string; remark: string }) => {
  try {
    if (editingId.value) {
      await updateCmdbHostTag(editingId.value, values)
      message.success('标签更新成功')
    } else {
      await createCmdbHostTag(values)
      message.success('标签创建成功')
    }
    drawerVisible.value = false
    await fetchData()
  } catch {
    // handled by interceptor
  }
}

const handleDelete = (record: CmdbHostTag) => {
  Modal.confirm({
    title: '确认删除标签',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除标签「${record.name}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteCmdbHostTag(record.id)
      message.success('标签删除成功')
      await fetchData()
    },
  })
}

fetchData()
</script>

<style scoped>
.cmdb-tag-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.summary-item {
  display: grid;
  grid-template-columns: 34px 1fr auto;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  color: var(--app-text-strong);
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.summary-item :deep(.anticon) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  color: #1677ff;
  background: rgba(22, 119, 255, 0.1);
  border-radius: 8px;
}

.summary-item.success :deep(.anticon) {
  color: #389e0d;
  background: rgba(82, 196, 26, 0.12);
}

.summary-item span {
  color: var(--app-text-muted);
  font-size: 13px;
}

.summary-item strong {
  font-size: 22px;
  font-weight: 650;
}

.tag-code {
  color: var(--app-text-muted);
  font-size: 12px;
}

.color-display {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.color-display__block {
  width: 18px;
  height: 18px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
}

@media (max-width: 960px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
