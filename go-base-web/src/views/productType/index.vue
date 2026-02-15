<template>
  <div class="productType-page">
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
          <a-button type="primary" @click="handleAdd" v-permission="'product_type:add'"><PlusOutlined /> 新增</a-button>
          <a-button danger :disabled="selectedRowKeys.length === 0" @click="confirmBatchDelete" v-permission="'product_type:delete'">
            <DeleteOutlined /> 批量删除 {{ selectedRowKeys.length > 0 ? `(${selectedRowKeys.length})` : '' }}
          </a-button>
          <a-button @click="handleExport" v-permission="'product_type:export'"><DownloadOutlined /> 导出</a-button>
          <a-upload
            :show-upload-list="false"
            :before-upload="handleImport"
            accept=".xlsx,.xls"
            v-permission="'product_type:import'"
          >
            <a-button><UploadOutlined /> 导入</a-button>
          </a-upload>
          <a-button @click="handleDownloadTemplate" v-permission="'product_type:import'"><FileExcelOutlined /> 下载模板</a-button>
        </a-space>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColors[record.status]">
            {{ statusOptions[record.status] || record.status }}
          </a-tag>
        </template>
        <template v-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'product_type:edit'">编辑</a-button>
          <a-button type="link" size="small" @click="handleCopy(record)" v-permission="'product_type:edit'">复制</a-button>
          <a-button type="link" size="small" danger @click="confirmDelete(record.id)" v-permission="'product_type:delete'">删除</a-button>
        </template>
      </template>
    </ProTable>

    <!-- 表单抽屉 -->
    <ProductTypeForm
      v-model:open="drawerVisible"
      :record="currentRecord"
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
import { PlusOutlined, DeleteOutlined, ExclamationCircleOutlined, DownloadOutlined, UploadOutlined, FileExcelOutlined } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
import ProductTypeForm from './components/ProductTypeForm.vue'
import { getProductTypeList, deleteProductType, batchDeleteProductType, exportProductType, importProductType, downloadTemplateProductType } from '@/api/productType'
import { formatTime } from '@/utils/format'
import { ProductType } from '@/types/productType'
import { useTableColumns } from '@/utils/permission'
import { getDictDataByType } from '@/api/dict'

const loading = ref(false)
const tableData = ref<ProductType[]>([])
const drawerVisible = ref(false)
const currentRecord = ref<ProductType | null>(null)
const selectedRowKeys = ref<number[]>([])
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

const searchForm = reactive({
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
  { title: '产品类型名称', dataIndex: 'name', key: 'name', align: 'center' },
  { title: '类型图标', dataIndex: 'icon', key: 'icon', align: 'center' },
  { title: '', dataIndex: 'status', key: 'status', align: 'center' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', align: 'center', sorter: true },
]

// 操作列配置
const actionColumn = { title: '操作', key: 'action', width: 200, align: 'center' }

// 根据权限动态显示操作列（有编辑或删除权限时显示）
const columns = useTableColumns(baseColumns, actionColumn, ['product_type:edit', 'product_type:delete'])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getProductTypeList({
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
  currentRecord.value = null
  drawerVisible.value = true
}

const handleEdit = (record: ProductType) => {
  currentRecord.value = record
  drawerVisible.value = true
}

const handleCopy = (record: ProductType) => {
  // 复制时不传 id，使表单识别为新增模式
  const { id, created_at, updated_at, ...copyData } = record
  currentRecord.value = copyData as ProductType
  drawerVisible.value = true
}

const handleFormSuccess = () => {
  fetchData()
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
      await deleteProductType(id)
      message.success('删除成功')
      fetchData()
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
      await batchDeleteProductType(selectedRowKeys.value)
      message.success('批量删除成功')
      selectedRowKeys.value = []
      fetchData()
    }
  })
}

// 导出数据
const handleExport = async () => {
  try {
    loading.value = true
    const res = await exportProductType(searchForm)
    const blob = new Blob([res as any], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `productType_${new Date().getTime()}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    message.success('导出成功')
  } catch (error) {
    message.error('导出失败')
  } finally {
    loading.value = false
  }
}

// 导入数据
const handleImport = async (file: File) => {
  try {
    loading.value = true
    const res = await importProductType(file)
    if (res.data.fail_count > 0) {
      Modal.info({
        title: '导入完成',
        width: 600,
        content: createVNode('div', {}, [
          createVNode('p', {}, `成功: ${res.data.success_count} 条，失败: ${res.data.fail_count} 条`),
          res.data.errors.length > 0 && createVNode('div', { style: 'max-height: 300px; overflow-y: auto; margin-top: 10px;' }, [
            createVNode('p', { style: 'font-weight: bold; color: #ff4d4f;' }, '错误详情:'),
            ...res.data.errors.map((err: string) => createVNode('p', { style: 'color: #666; font-size: 12px;' }, err))
          ])
        ])
      })
    } else {
      message.success(`导入成功 ${res.data.success_count} 条数据`)
    }
    fetchData()
  } catch (error: any) {
    message.error(error.message || '导入失败')
  } finally {
    loading.value = false
  }
  return false // 阻止默认上传行为
}

// 下载导入模板
const handleDownloadTemplate = async () => {
  try {
    const res = await downloadTemplateProductType()
    const blob = new Blob([res as any], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'productType_template.xlsx'
    link.click()
    window.URL.revokeObjectURL(url)
    message.success('模板下载成功')
  } catch (error) {
    message.error('模板下载失败')
  }
}

// 获取字典数据
const fetchStatusDict = async () => {
  try {
    const res = await getDictDataByType('common_status')
    statusDictList.value = res.data || []
  } catch (e) {
    console.error('获取字典失败', e)
  }
}

onMounted(() => {
  // 加载字典数据
  fetchStatusDict()
  fetchData()
})
</script>

<style scoped>
.productType-page {
  padding: 0;
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
