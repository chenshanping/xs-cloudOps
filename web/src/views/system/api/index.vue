<template>
  <div class="api-page">
    <a-row :gutter="16">
      <!-- 左侧分组列表 -->
      <a-col :span="4">
        <a-card title="API分组" size="small">
          <a-menu
            v-model:selectedKeys="selectedGroup"
            mode="inline"
            @click="handleGroupClick"
          >
            <a-menu-item key="">
              <template #icon><AppstoreOutlined /></template>
              全部
              <a-badge :count="totalApiCount" :overflow-count="999" class="group-badge" />
            </a-menu-item>
            <a-menu-item v-for="item in groupList" :key="item.group">
              <template #icon><FolderOutlined /></template>
              {{ item.group }}
              <a-badge :count="item.api_count" :overflow-count="999" class="group-badge" />            
            </a-menu-item>
          </a-menu>
        </a-card>
      </a-col>

      <!-- 右侧表格 -->
      <a-col :span="20">
        <ProTable
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          :pagination="pagination"
          @search="handleSearch"
          @reset="handleReset"
          @change="handleTableChange"
        >
          <!-- 搜索区域 -->
          <template #search>
            <a-form-item label="API路径">
              <a-input v-model:value="searchForm.path" placeholder="请输入API路径" allowClear />
            </a-form-item>
            <a-form-item label="请求方法">
              <a-select v-model:value="searchForm.method" placeholder="请选择" allowClear style="width: 120px">
                <a-select-option value="GET">GET</a-select-option>
                <a-select-option value="POST">POST</a-select-option>
                <a-select-option value="PUT">PUT</a-select-option>
                <a-select-option value="DELETE">DELETE</a-select-option>
              </a-select>
            </a-form-item>
          </template>

          <!-- 工具栏 -->
          <template #toolbar>
            <a-space>
              <a-button @click="handleSync" :loading="syncLoading" v-permission="'system:api:sync'">
                <SyncOutlined /> 同步API
              </a-button>
              <a-button type="primary" @click="handleAdd" v-permission="'system:api:add'">
                <PlusOutlined /> 新增
              </a-button>
            </a-space>
          </template>

          <!-- 表格单元格 -->
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'method'">
              <a-tag :color="getMethodColor(record.method)">{{ record.method }}</a-tag>
            </template>
            <template v-if="column.key === 'need_auth'">
              <a-tag :color="record.need_auth ? 'orange' : 'green'">{{ record.need_auth ? '需要' : '公开' }}</a-tag>
            </template>
            <template v-if="column.key === 'params'">
              <a-tag 
                v-if="parseParams(record.request_params).length > 0" 
                color="blue" 
                class="params-tag"
                @click="showParams(record)"
              >
                {{ parseParams(record.request_params).length }} 个参数
              </a-tag>
              <span v-else class="text-gray">-</span>
            </template>
            <template v-if="column.key === 'action'">
              <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:api:edit'">编辑</a-button>
              <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger v-permission="'system:api:delete'">删除</a-button>
              </a-popconfirm>
            </template>
          </template>

        </ProTable>
      </a-col>
    </a-row>

    <!-- 新增/编辑弹窗 -->
    <a-modal v-model:open="modalVisible" :title="modalTitle" @ok="handleModalOk">
      <a-form :model="formState" :label-col="{ span: 5 }">
        <a-form-item label="API路径" required><a-input v-model:value="formState.path" /></a-form-item>
        <a-form-item label="请求方法" required>
          <a-select v-model:value="formState.method">
            <a-select-option value="GET">GET</a-select-option>
            <a-select-option value="POST">POST</a-select-option>
            <a-select-option value="PUT">PUT</a-select-option>
            <a-select-option value="DELETE">DELETE</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分组"><a-input v-model:value="formState.group" /></a-form-item>
        <a-form-item label="描述"><a-input v-model:value="formState.description" /></a-form-item>
      </a-form>
    </a-modal>

    <!-- 参数详情弹窗 -->
    <a-modal v-model:open="paramsModalVisible" title="API参数详情" :footer="null" width="700px">
      <a-descriptions :column="2" bordered size="small" class="mb-4">
        <a-descriptions-item label="API路径">{{ currentApi?.path }}</a-descriptions-item>
        <a-descriptions-item label="请求方法">
          <a-tag :color="getMethodColor(currentApi?.method || '')">{{ currentApi?.method }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">{{ currentApi?.description || '-' }}</a-descriptions-item>
      </a-descriptions>
      
      <h4>请求参数</h4>
      <a-table 
        :columns="paramsColumns" 
        :data-source="currentParams" 
        :pagination="false"
        size="small"
        bordered
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'required'">
            <a-tag :color="record.required ? 'red' : 'default'">{{ record.required ? '是' : '否' }}</a-tag>
          </template>
          <template v-if="column.key === 'in'">
            <a-tag>{{ record.in }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, SyncOutlined, AppstoreOutlined, FolderOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import { getApiList, getApiGroups, createApi, updateApi, deleteApi, syncApis, type ApiGroupStats } from '@/api/api'
import { useTableColumns } from '@/utils/permission'
import type { Api, ApiFieldInfo } from '@/types'

const loading = ref(false)
const syncLoading = ref(false)
const tableData = ref<Api[]>([])
const groupList = ref<ApiGroupStats[]>([])
const selectedGroup = ref<string[]>([''])
const modalVisible = ref(false)
const paramsModalVisible = ref(false)
const modalTitle = ref('新增API')
const isEdit = ref(false)
const currentId = ref(0)
const currentApi = ref<Api | null>(null)
const currentParams = ref<ApiFieldInfo[]>([])
const searchForm = reactive({ path: '', method: undefined as string | undefined })
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
const formState = reactive({ path: '', method: 'GET', group: '', description: '' })

// 参数表格列
const paramsColumns = [
  { title: '参数名', dataIndex: 'name', key: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', key: 'type', width: 120 },
  { title: '位置', key: 'in', width: 80 },
  { title: '必填', key: 'required', width: 80 },
  { title: '描述', dataIndex: 'description', key: 'description' },
]

// 使用工具函数动态生成列配置
const columns = useTableColumns(
  [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'API路径', dataIndex: 'path', key: 'path', ellipsis: true },
    { title: '方法', key: 'method', width: 80 },
    { title: '认证', key: 'need_auth', width: 80 },
    { title: '参数', key: 'params', width: 100 },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  ],
  { title: '操作', key: 'action', width: 120 },
  ['system:api:edit', 'system:api:delete']
)

const getMethodColor = (m: string) => ({ GET: 'green', POST: 'blue', PUT: 'orange', DELETE: 'red' }[m] || 'default')

// 解析参数 JSON
const parseParams = (paramsStr?: string): ApiFieldInfo[] => {
  if (!paramsStr) return []
  try {
    return JSON.parse(paramsStr)
  } catch {
    return []
  }
}

// 显示参数弹窗
const showParams = (record: Api) => {
  currentApi.value = record
  currentParams.value = parseParams(record.request_params)
  paramsModalVisible.value = true
}

// 计算总 API 数量
const totalApiCount = computed(() => {
  return groupList.value.reduce((sum, item) => sum + item.api_count, 0)
})

const fetchGroups = async () => {
  try {
    const res = await getApiGroups()
    groupList.value = res.data || []
  } catch { /* ignore */ }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getApiList({
      page: pagination.current,
      page_size: pagination.pageSize,
      path: searchForm.path,
      method: searchForm.method,
      group: selectedGroup.value[0] || undefined
    })
    tableData.value = res.data.list
    pagination.total = res.data.total
  } finally { loading.value = false }
}

const handleGroupClick = ({ key }: { key: string }) => {
  selectedGroup.value = [key]
  pagination.current = 1
  fetchData()
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchForm.path = ''
  searchForm.method = undefined
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pag: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchData()
}
const handleAdd = () => { isEdit.value = false; modalTitle.value = '新增API'; Object.assign(formState, { path: '', method: 'GET', group: '', description: '' }); modalVisible.value = true }
const handleEdit = (record: Api) => { isEdit.value = true; modalTitle.value = '编辑API'; currentId.value = record.id; Object.assign(formState, record); modalVisible.value = true }
const handleModalOk = async () => {
  if (isEdit.value) { await updateApi(currentId.value, formState); message.success('更新成功') }
  else { await createApi(formState); message.success('创建成功') }
  modalVisible.value = false; fetchData()
}
const handleDelete = async (record: Api) => { await deleteApi(record.id); message.success('删除成功'); fetchData() }
const handleSync = async () => {
  syncLoading.value = true
  try {
    const res = await syncApis()
    message.success(`同步完成，新增 ${res.data.added} 条，更新 ${res.data.updated} 条，删除 ${res.data.deleted} 条`)
    fetchGroups()
    fetchData()
  } finally { syncLoading.value = false }
}
onMounted(() => { fetchGroups(); fetchData() })
</script>

<style scoped>
.api-page {
  height: 100%;
}
.api-page :deep(.ant-menu) {
  border-right: none;
}
.api-page :deep(.ant-card-body) {
  padding: 0;
}
.text-gray {
  color: #999;
}
.params-tag {
  cursor: pointer;
}
.params-tag:hover {
  opacity: 0.8;
}
.mb-4 {
  margin-bottom: 16px;
}
.group-badge {
  margin-left: auto;
}
.group-badge :deep(.ant-badge-count) {
  background: #e6f7ff;
  color: #1890ff;
  box-shadow: none;
}
.api-page :deep(.ant-menu-item-selected) .group-badge :deep(.ant-badge-count) {
  background: #1890ff;
  color: #fff;
}
</style>
