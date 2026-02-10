<template>
  <div class="product-page">
    <div class="tree-table-layout">
      <!-- 左侧分组树 -->
      <div class="category-tree-panel">
      <div class="tree-header">
        <span class="tree-title"><FolderOutlined /> 产品分类</span>
      </div>
      <a-spin :spinning="treeLoading">
        <div class="tree-content">
          <div 
            class="tree-item" 
            :class="{ active: !selectedCategoryId }" 
            @click="handleSelectCategory(null)"
          >
            <AppstoreOutlined />
            <span>全部</span>
            <span class="item-count">{{ totalCount }}</span>
          </div>
          <div 
            v-for="item in categoryOptions" 
            :key="item.id" 
            class="tree-item"
            :class="{ active: selectedCategoryId === item.id }"
            @click="handleSelectCategory(item.id)"
          >
            <TagOutlined />
            <span>{{ item.name }}</span>
            <span class="item-count">{{ item.count || 0 }}</span>
          </div>
        </div>
        </a-spin>
      </div>
      <!-- 右侧表格 -->
      <div class="table-panel">
    <ProTable
      :columns="columns"
      :data-source="tableData"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      :row-selection="{ selectedRowKeys, onChange: onSelectChange }"
      @update:pagination="val => Object.assign(pagination, val)"
       :scroll="{ x: 2000,y: 400 }"
      @change="handleTableChange"
      @search="handleSearch"
      @reset="handleReset"
    >

      <template #toolbar>
        <a-space>
          <a-button type="primary" @click="handleAdd" v-permission="'product:add'"><PlusOutlined /> 新增</a-button>
          <a-button danger :disabled="selectedRowKeys.length === 0" @click="confirmBatchDelete" v-permission="'product:delete'">
            <DeleteOutlined /> 批量删除 {{ selectedRowKeys.length > 0 ? `(${selectedRowKeys.length})` : '' }}
          </a-button>
        </a-space>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColors[record.status]">
            {{ statusOptions[record.status] || record.status }}
          </a-tag>
        </template>
        <template v-if="column.key === 'product_type'">
          {{ record.product_type?.name || '-' }}
        </template>
        <template v-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'product:edit'">编辑</a-button>
          <a-button type="link" size="small" @click="handleCopy(record)" v-permission="'product:edit'">复制</a-button>
          <a-button type="link" size="small" danger @click="confirmDelete(record.id)" v-permission="'product:delete'">删除</a-button>
        </template>
      </template>
    </ProTable>
      </div>
    </div>

    <!-- 表单抽屉 -->
    <ProductForm
      v-model:open="drawerVisible"
      :record="currentRecord"
      :product_type-options="product_typeOptions"
      @success="handleFormSuccess"
    />
    <!-- 文本域内容预览 -->
    <a-modal v-model:open="textModalVisible" :title="textModalTitle" :footer="null" width="600px">
      <div class="text-content">{{ textModalContent }}</div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined, DeleteOutlined, ExclamationCircleOutlined, FolderOutlined, AppstoreOutlined, TagOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import ProductForm from './components/ProductForm.vue'
import { getProductList, deleteProduct, batchDeleteProduct } from '@/api/product'
import { getProductTypeOptions } from '@/api/productType'
import { formatTime } from '@/utils/format'
import { Product } from '@/types/product'
import { useTableColumns } from '@/utils/permission'
import { getDictDataByType } from '@/api/dict'

const loading = ref(false)
const tableData = ref<Product[]>([])
const drawerVisible = ref(false)
const currentRecord = ref<Product | null>(null)
const selectedRowKeys = ref<number[]>([])

// 左树右表相关
const treeLoading = ref(false)
const selectedCategoryId = ref<number | null>(null)
const categoryOptions = computed(() => product_typeOptions.value)
const totalCount = computed(() => categoryOptions.value.reduce((sum: number, item: any) => sum + (item.count || 0), 0))

const handleSelectCategory = (id: number | null) => {
  selectedCategoryId.value = id
  searchForm.type_id = id || undefined
  pagination.current = 1
  fetchData()
}
// 文本域内容预览
const textModalVisible = ref(false)
const textModalTitle = ref('')
const textModalContent = ref('')

const handleViewText = (record: any, title: string, field: string) => {
  textModalTitle.value = title
  textModalContent.value = record[field] || '暂无内容'
  textModalVisible.value = true
}

// 默认标签颜色（当字典未配置 tag_type 时作为回退）
const defaultTagColors = ['blue', 'green', 'orange', 'purple', 'cyan', 'magenta', 'gold', 'lime']
// 字典选项映射（动态获取）
const statusDictList = ref<any[]>([])
const statusOptions = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  statusDictList.value.forEach(item => { map[item.value] = item.label })
  return map
})
const statusColors = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  statusDictList.value.forEach((item, i) => {
    // 优先使用字典配置的 tag_type，否则使用默认颜色
    map[item.value] = item.tag_type || defaultTagColors[i % defaultTagColors.length]
  })
  return map
})

// 关联选项
const product_typeOptions = ref<any[]>([])

// 下拉搜索过滤
const filterOption = (input: string, option: any) => {
  return option.children?.[0]?.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const searchForm = reactive({
  type_id: undefined as number | undefined,
})

// 排序参数
const sortInfo = reactive({
  field: '',
  order: '' as '' | 'ascend' | 'descend'
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})
// 基础列配置（不含操作列）
const baseColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', align: 'center', sorter: true },
  { title: '产品名称', dataIndex: 'name', key: 'name', align: 'center' },
  { title: '产品数量', dataIndex: 'num', key: 'num', align: 'center' },
  { title: '产品单价', dataIndex: 'price', key: 'price', align: 'center' },
  { title: '状态', dataIndex: 'status', key: 'status', align: 'center' },
  { title: '产品分类', dataIndex: 'product_type', key: 'product_type', align: 'center' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', align: 'center', sorter: true },
]

// 操作列配置
const actionColumn = { title: '操作', key: 'action', width: 200, align: 'center' }

// 根据权限动态显示操作列（有编辑或删除权限时显示）
const columns = useTableColumns(baseColumns, actionColumn, ['product:edit', 'product:delete'])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getProductList({
      page: pagination.current,
      page_size: pagination.pageSize,
      ...searchForm,
      sort_field: sortInfo.field || undefined,
      sort_order: sortInfo.order || undefined
    })
    tableData.value = res.data.list
    pagination.total = res.data.total
  } catch {
    // 错误已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchForm.type_id = undefined
  sortInfo.field = ''
  sortInfo.order = ''
  handleSearch()
}

const handleTableChange = (pag: any, _filters: any, sorter: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  // 处理排序
  if (sorter && sorter.field) {
    sortInfo.field = sorter.field
    sortInfo.order = sorter.order || ''
  } else {
    sortInfo.field = ''
    sortInfo.order = ''
  }
  fetchData()
}

const handleAdd = () => {
  // 左树右表：新增时传递当前选中的分类
  if (selectedCategoryId.value) {
    currentRecord.value = {
      type_id: selectedCategoryId.value
    } as any
  } else {
    currentRecord.value = null
  }
  drawerVisible.value = true
}

const handleEdit = (record: Product) => {
  currentRecord.value = record
  drawerVisible.value = true
}

const handleCopy = (record: Product) => {
  // 复制时不传 id，使表单识别为新增模式
  const { id, created_at, updated_at, ...copyData } = record
  currentRecord.value = copyData as Product
  drawerVisible.value = true
}

const handleFormSuccess = () => {
  fetchData()
  fetchProductTypeOptions() // 刷新关联选项count
}

// 确认删除
const confirmDelete = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确定要删除该条数据吗？',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteProduct(id)
      message.success('删除成功')
      fetchData()
      fetchProductTypeOptions()
    }
  })
}

// 行选择变化
const onSelectChange = (keys: number[]) => {
  selectedRowKeys.value = keys
}

// 确认批量删除
const confirmBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的数据')
    return
  }
  Modal.confirm({
    title: '确认批量删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 条数据吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await batchDeleteProduct(selectedRowKeys.value)
      message.success('批量删除成功')
      selectedRowKeys.value = []
      fetchData()
      fetchProductTypeOptions()
    }
  })
}
// 获取关联选项
const fetchProductTypeOptions = async () => {
  try {
    // 使用轻量ptions接口（返回id,name,count）
    const res = await getProductTypeOptions({
      display_field: 'name',
      count_table: 'product',
      count_field: 'type_id'
    })
    product_typeOptions.value = res.data || []
  } catch (e) {
    console.error('获取产品分类选项失败', e)
  }
}

// 获取字典数据
const fetchStatusDict = async () => {
  try {
    const res = await getDictDataByType('common_status')
    statusDictList.value = res.data || []
  } catch (e) {
    console.error('获取状态字典失败', e)
  }
}

onMounted(() => {
  fetchProductTypeOptions()
  // 加载字典数据
  fetchStatusDict()
  fetchData()
})
</script>

<style scoped>
.product-page {
  padding: 0;
}
.tree-table-layout {
  display: flex;
  height: calc(100vh - 120px);
  gap: 16px;
}
.category-tree-panel {
  width: 240px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.category-tree-panel .tree-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}
.category-tree-panel .tree-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}
.category-tree-panel .tree-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}
.category-tree-panel .tree-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  margin: 2px 8px;
  border-radius: 6px;
  cursor: pointer;
  color: #666;
  transition: all 0.2s;
}
.category-tree-panel .tree-item:hover {
  background: #f5f5f5;
  color: #333;
}
.category-tree-panel .tree-item.active {
  background: #e6f7ff;
  color: #1890ff;
  font-weight: 500;
}
.category-tree-panel .tree-item.active .item-count {
  background: #1890ff;
  color: #fff;
}
.category-tree-panel .item-count {
  margin-left: auto;
  background: #f0f0f0;
  color: #999;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
  min-width: 24px;
  text-align: center;
}
.table-panel {
  flex: 1;
  min-width: 0;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}
.text-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  color: #333;
}
</style>
