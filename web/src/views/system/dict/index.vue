<template>
  <div class="dict-page">
    <a-row :gutter="16">
      <!-- 左侧：字典类型列表 -->
      <a-col :span="8">
        <a-card title="字典类型">
          <template #extra>
            <a-button type="primary" size="small" @click="handleAddType">
              <PlusOutlined /> 新增
            </a-button>
          </template>

          <!-- 搜索 -->
          <a-input-search
            v-model:value="typeSearchText"
            placeholder="搜索字典名称/类型"
            style="margin-bottom: 16px"
            @search="fetchDictTypes"
            allow-clear
          />

          <!-- 字典类型表格 -->
          <a-table
            :columns="typeColumns"
            :data-source="dictTypes"
            :loading="typeLoading"
            :pagination="typePagination"
            row-key="id"
            size="small"
            :row-class-name="(record: DictType) => selectedType?.id === record.id ? 'selected-row' : ''"
            :custom-row="(record: DictType) => ({ onClick: () => handleSelectType(record) })"
            @change="handleTypeTableChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === 1 ? 'success' : 'error'">
                  {{ record.status === 1 ? '正常' : '停用' }}
                </a-tag>
              </template>
              <template v-if="column.key === 'action'">
                <a-space>
                  <a-button type="link" size="small" @click.stop="handleEditType(record)">编辑</a-button>
                  <a-popconfirm title="确定删除此字典类型及其所有数据?" @confirm="handleDeleteType(record.id)">
                    <a-button type="link" size="small" danger @click.stop>删除</a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 右侧：字典数据列表 -->
      <a-col :span="16">
        <a-card :title="selectedType ? `字典数据 - ${selectedType.name}` : '字典数据'">
          <template #extra>
            <a-button type="primary" size="small" @click="handleAddData" :disabled="!selectedType">
              <PlusOutlined /> 新增
            </a-button>
          </template>

          <a-empty v-if="!selectedType" description="请选择左侧字典类型" />

          <template v-else>
            <!-- 字典数据表格 -->
            <a-table
              :columns="dataColumns"
              :data-source="dictDataList"
              :loading="dataLoading"
              :pagination="dataPagination"
              row-key="id"
              size="small"
              @change="handleDataTableChange"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'label'">
                  <a-tag v-if="record.tag_type" :color="record.tag_type">{{ record.label }}</a-tag>
                  <span v-else>{{ record.label }}</span>
                </template>
                <template v-if="column.key === 'status'">
                  <a-tag :color="record.status === 1 ? 'success' : 'error'">
                    {{ record.status === 1 ? '正常' : '停用' }}
                  </a-tag>
                </template>
                <template v-if="column.key === 'is_default'">
                  <a-tag v-if="record.is_default === 1" color="blue">是</a-tag>
                  <span v-else>否</span>
                </template>
                <template v-if="column.key === 'action'">
                  <a-space>
                    <a-button type="link" size="small" @click="handleEditData(record)">编辑</a-button>
                    <a-popconfirm title="确定删除此字典数据?" @confirm="handleDeleteData(record.id)">
                      <a-button type="link" size="small" danger>删除</a-button>
                    </a-popconfirm>
                  </a-space>
                </template>
              </template>
            </a-table>
          </template>
        </a-card>
      </a-col>
    </a-row>

    <!-- 字典类型表单弹窗 -->
    <a-modal
      v-model:open="typeModalVisible"
      :title="typeModalTitle"
      @ok="handleTypeSubmit"
      :confirm-loading="typeSubmitLoading"
    >
      <a-form :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }">
        <a-form-item label="字典名称" required>
          <a-input v-model:value="typeForm.name" placeholder="请输入字典名称" />
        </a-form-item>
        <a-form-item label="字典类型" required>
          <a-input v-model:value="typeForm.type" placeholder="请输入字典类型（英文）" :disabled="!!editingType" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="typeForm.statusBool" checked-children="正常" un-checked-children="停用" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="typeForm.remark" :rows="3" placeholder="请输入备注" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 字典数据表单弹窗 -->
    <a-modal
      v-model:open="dataModalVisible"
      :title="dataModalTitle"
      @ok="handleDataSubmit"
      :confirm-loading="dataSubmitLoading"
    >
      <a-form :label-col="{ span: 5 }" :wrapper-col="{ span: 18 }">
        <a-form-item label="字典标签" required>
          <a-input v-model:value="dataForm.label" placeholder="请输入字典标签（显示名称）" />
        </a-form-item>
        <a-form-item label="字典键值" required>
          <a-input v-model:value="dataForm.value" placeholder="请输入字典键值" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="dataForm.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="标签颜色">
          <a-select v-model:value="dataForm.tag_type" placeholder="请选择标签颜色" allow-clear>
            <a-select-option value="success">
              <a-tag color="success">成功/绿色</a-tag>
            </a-select-option>
            <a-select-option value="processing">
              <a-tag color="processing">处理中/蓝色</a-tag>
            </a-select-option>
            <a-select-option value="warning">
              <a-tag color="warning">警告/橙色</a-tag>
            </a-select-option>
            <a-select-option value="error">
              <a-tag color="error">错误/红色</a-tag>
            </a-select-option>
            <a-select-option value="default">
              <a-tag color="default">默认/灰色</a-tag>
            </a-select-option>
            <a-select-option value="pink">
              <a-tag color="pink">粉色</a-tag>
            </a-select-option>
            <a-select-option value="purple">
              <a-tag color="purple">紫色</a-tag>
            </a-select-option>
            <a-select-option value="cyan">
              <a-tag color="cyan">青色</a-tag>
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="dataForm.statusBool" checked-children="正常" un-checked-children="停用" />
        </a-form-item>
        <a-form-item label="是否默认">
          <a-switch v-model:checked="dataForm.isDefaultBool" checked-children="是" un-checked-children="否" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="dataForm.remark" :rows="2" placeholder="请输入备注" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import {
  getDictTypeList,
  createDictType,
  updateDictType,
  deleteDictType,
  getDictDataList,
  createDictData,
  updateDictData,
  deleteDictData,
  type DictType,
  type DictData
} from '@/api/dict'

// ==================== 字典类型 ====================
const dictTypes = ref<DictType[]>([])
const typeLoading = ref(false)
const typeSearchText = ref('')
const selectedType = ref<DictType | null>(null)
const typePagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const typeColumns = [
  { title: '字典名称', dataIndex: 'name', key: 'name', ellipsis: true },
  { title: '字典类型', dataIndex: 'type', key: 'type', ellipsis: true },
  { title: '状态', key: 'status', width: 70 },
  { title: '操作', key: 'action', width: 100 }
]

// 字典类型表单
const typeModalVisible = ref(false)
const typeModalTitle = ref('新增字典类型')
const typeSubmitLoading = ref(false)
const editingType = ref<DictType | null>(null)
const typeForm = reactive({
  name: '',
  type: '',
  statusBool: true,
  remark: ''
})

// 获取字典类型列表
const fetchDictTypes = async () => {
  typeLoading.value = true
  try {
    const res = await getDictTypeList({
      page: typePagination.current,
      page_size: typePagination.pageSize,
      name: typeSearchText.value || undefined,
      type: typeSearchText.value || undefined
    })
    dictTypes.value = res.data.list || []
    typePagination.total = res.data.total
  } finally {
    typeLoading.value = false
  }
}

// 字典类型表格分页变化
const handleTypeTableChange = (pagination: any) => {
  typePagination.current = pagination.current
  typePagination.pageSize = pagination.pageSize
  fetchDictTypes()
}

// 选择字典类型
const handleSelectType = (record: DictType) => {
  selectedType.value = record
  dataPagination.current = 1
  fetchDictData()
}

// 新增字典类型
const handleAddType = () => {
  editingType.value = null
  typeModalTitle.value = '新增字典类型'
  Object.assign(typeForm, { name: '', type: '', statusBool: true, remark: '' })
  typeModalVisible.value = true
}

// 编辑字典类型
const handleEditType = (record: DictType) => {
  editingType.value = record
  typeModalTitle.value = '编辑字典类型'
  Object.assign(typeForm, {
    name: record.name,
    type: record.type,
    statusBool: record.status === 1,
    remark: record.remark
  })
  typeModalVisible.value = true
}

// 提交字典类型
const handleTypeSubmit = async () => {
  if (!typeForm.name || !typeForm.type) {
    message.warning('请填写字典名称和类型')
    return
  }
  typeSubmitLoading.value = true
  try {
    const data = {
      name: typeForm.name,
      type: typeForm.type,
      status: typeForm.statusBool ? 1 : 0,
      remark: typeForm.remark
    }
    if (editingType.value) {
      await updateDictType(editingType.value.id, data)
      message.success('更新成功')
    } else {
      await createDictType(data)
      message.success('创建成功')
    }
    typeModalVisible.value = false
    fetchDictTypes()
  } finally {
    typeSubmitLoading.value = false
  }
}

// 删除字典类型
const handleDeleteType = async (id: number) => {
  await deleteDictType(id)
  message.success('删除成功')
  if (selectedType.value?.id === id) {
    selectedType.value = null
    dictDataList.value = []
  }
  fetchDictTypes()
}

// ==================== 字典数据 ====================
const dictDataList = ref<DictData[]>([])
const dataLoading = ref(false)
const dataPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`
})

const dataColumns = [
  { title: '字典标签', key: 'label', ellipsis: true },
  { title: '字典键值', dataIndex: 'value', key: 'value', width: 100 },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 60 },
  { title: '状态', key: 'status', width: 70 },
  { title: '默认', key: 'is_default', width: 60 },
  { title: '操作', key: 'action', width: 100 }
]

// 字典数据表单
const dataModalVisible = ref(false)
const dataModalTitle = ref('新增字典数据')
const dataSubmitLoading = ref(false)
const editingData = ref<DictData | null>(null)
const dataForm = reactive({
  label: '',
  value: '',
  sort: 0,
  tag_type: '',
  statusBool: true,
  isDefaultBool: false,
  remark: ''
})

// 获取字典数据列表
const fetchDictData = async () => {
  if (!selectedType.value) return
  dataLoading.value = true
  try {
    const res = await getDictDataList({
      dict_type: selectedType.value.type,
      page: dataPagination.current,
      page_size: dataPagination.pageSize
    })
    dictDataList.value = res.data.list || []
    dataPagination.total = res.data.total
  } finally {
    dataLoading.value = false
  }
}

// 字典数据表格分页变化
const handleDataTableChange = (pagination: any) => {
  dataPagination.current = pagination.current
  dataPagination.pageSize = pagination.pageSize
  fetchDictData()
}

// 新增字典数据
const handleAddData = () => {
  if (!selectedType.value) return
  editingData.value = null
  dataModalTitle.value = '新增字典数据'
  Object.assign(dataForm, {
    label: '', value: '', sort: 0, tag_type: '',
    statusBool: true, isDefaultBool: false, remark: ''
  })
  dataModalVisible.value = true
}

// 编辑字典数据
const handleEditData = (record: DictData) => {
  editingData.value = record
  dataModalTitle.value = '编辑字典数据'
  Object.assign(dataForm, {
    label: record.label,
    value: record.value,
    sort: record.sort,
    tag_type: record.tag_type || '',
    statusBool: record.status === 1,
    isDefaultBool: record.is_default === 1,
    remark: record.remark
  })
  dataModalVisible.value = true
}

// 提交字典数据
const handleDataSubmit = async () => {
  if (!dataForm.label || !dataForm.value) {
    message.warning('请填写字典标签和键值')
    return
  }
  dataSubmitLoading.value = true
  try {
    const data = {
      dict_type: selectedType.value!.type,
      label: dataForm.label,
      value: dataForm.value,
      sort: dataForm.sort,
      tag_type: dataForm.tag_type || '',
      status: dataForm.statusBool ? 1 : 0,
      is_default: dataForm.isDefaultBool ? 1 : 0,
      remark: dataForm.remark
    }
    if (editingData.value) {
      await updateDictData(editingData.value.id, data)
      message.success('更新成功')
    } else {
      await createDictData(data)
      message.success('创建成功')
    }
    dataModalVisible.value = false
    fetchDictData()
  } finally {
    dataSubmitLoading.value = false
  }
}

// 删除字典数据
const handleDeleteData = async (id: number) => {
  await deleteDictData(id)
  message.success('删除成功')
  fetchDictData()
}

onMounted(() => {
  fetchDictTypes()
})
</script>

<style scoped>
.dict-page {
  padding: 16px;
}
:deep(.selected-row) {
  background-color: #e6f7ff;
}
:deep(.ant-table-tbody > tr) {
  cursor: pointer;
}
</style>
