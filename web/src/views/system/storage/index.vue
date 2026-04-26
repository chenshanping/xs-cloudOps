<template>
  <div class="storage-page">
    <ProTable
      title="存储配置"
      :columns="columns"
      :data-source="storageList"
      :loading="loading"
    >
      <template #toolbar>
        <a-button type="primary" @click="handleAdd">
          <PlusOutlined /> 新增
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          <a-tag :color="getTypeColor(record.type)">
            {{ getTypeLabel(record.type) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'is_default'">
          <a-tag v-if="record.is_default === 1" color="green">默认</a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-badge :status="record.status === 1 ? 'success' : 'error'" :text="record.status === 1 ? '启用' : '禁用'" />
        </template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="link" size="small" @click="handleSetDefault(record)" :disabled="record.is_default === 1">
            设为默认
          </a-button>
          <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record)">
            <a-button type="link" size="small" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </ProTable>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? '编辑存储配置' : '新增存储配置'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      width="600px"
    >
      <a-form :model="formState" :label-col="{ span: 5 }">
        <a-form-item label="名称" required>
          <a-input v-model:value="formState.name" placeholder="请输入存储名称" />
        </a-form-item>
        <a-form-item label="类型" required>
          <a-select v-model:value="formState.type" placeholder="请选择存储类型" @change="handleTypeChange">
            <a-select-option v-for="item in storageTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <!-- 本地存储配置 -->
        <template v-if="formState.type === 'local'">
<!--          <a-form-item label="存储路径" required>-->
<!--            <a-input v-model:value="configForm.base_path" placeholder="如：./uploads" />-->
<!--          </a-form-item>-->
<!--          <a-form-item label="访问URL" required>-->
<!--            <a-input v-model:value="configForm.base_url" placeholder="如：http://localhost:8080/uploads" />-->
<!--          </a-form-item>-->
        </template>

        <!-- 阿里云OSS配置 -->
        <template v-if="formState.type === 'aliyun'">
          <a-form-item label="Endpoint" required>
            <a-input v-model:value="configForm.endpoint" placeholder="如：oss-cn-hangzhou.aliyuncs.com" />
          </a-form-item>
          <a-form-item label="AccessKey ID" required>
            <a-input v-model:value="configForm.access_key_id" />
          </a-form-item>
          <a-form-item label="AccessKey Secret" required>
            <a-input-password v-model:value="configForm.access_key_secret" />
          </a-form-item>
          <a-form-item label="Bucket" required>
            <a-input v-model:value="configForm.bucket_name" />
          </a-form-item>
          <a-form-item label="Region">
            <a-input v-model:value="configForm.region" placeholder="如：cn-hangzhou" />
          </a-form-item>
        </template>

        <!-- 腾讯云COS配置 -->
        <template v-if="formState.type === 'tencent'">
          <a-form-item label="Region" required>
            <a-input v-model:value="configForm.region" placeholder="如：ap-guangzhou" />
          </a-form-item>
          <a-form-item label="SecretId" required>
            <a-input v-model:value="configForm.secret_id" />
          </a-form-item>
          <a-form-item label="SecretKey" required>
            <a-input-password v-model:value="configForm.secret_key" />
          </a-form-item>
          <a-form-item label="Bucket" required>
            <a-input v-model:value="configForm.bucket" placeholder="如：mybucket-1234567890" />
          </a-form-item>
          <a-form-item label="AppID">
            <a-input v-model:value="configForm.app_id" />
          </a-form-item>
        </template>

        <!-- MinIO配置 -->
        <template v-if="formState.type === 'minio'">
          <a-form-item label="Endpoint" required>
            <a-input v-model:value="configForm.endpoint" placeholder="如：127.0.0.1:9000" />
          </a-form-item>
          <a-form-item label="AccessKey ID" required>
            <a-input v-model:value="configForm.access_key_id" />
          </a-form-item>
          <a-form-item label="SecretAccessKey" required>
            <a-input-password v-model:value="configForm.secret_access_key" />
          </a-form-item>
          <a-form-item label="Bucket" required>
            <a-input v-model:value="configForm.bucket_name" />
          </a-form-item>
          <a-form-item label="使用SSL">
            <a-switch v-model:checked="configForm.use_ssl" />
          </a-form-item>
        </template>

        <a-form-item label="设为默认">
          <a-switch v-model:checked="formState.is_default_checked" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="formState.remark" :rows="2" />
        </a-form-item>
      </a-form>

      <template #footer>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="default" :loading="testing" @click="handleTest">测试连接</a-button>
        <a-button type="primary" :loading="submitting" @click="handleSubmit">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import ProTable from '@/components/ProTable.vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import type { Storage, StorageType } from '@/types/storage'
import { storageTypeOptions } from '@/types/storage'
import { getStorageList, createStorage, updateStorage, deleteStorage, setDefaultStorage, testStorage } from '@/api/storage'

const loading = ref(false)
const submitting = ref(false)
const testing = ref(false)
const modalVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(0)
const storageList = ref<Storage[]>([])

const formState = reactive({
  name: '',
  type: '' as StorageType | '',
  remark: '',
  is_default_checked: false,
})

const configForm = reactive<Record<string, any>>({})

const columns = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '类型', key: 'type' },
  { title: '默认', key: 'is_default', width: 80 },
  { title: '状态', key: 'status', width: 80 },
  { title: '备注', dataIndex: 'remark', key: 'remark', ellipsis: true },
  { title: '操作', key: 'action', width: 200 },
]

const getTypeLabel = (type: string) => {
  const item = storageTypeOptions.find(t => t.value === type)
  return item?.label || type
}

const getTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    local: 'default',
    aliyun: 'orange',
    tencent: 'blue',
    minio: 'green',
  }
  return colors[type] || 'default'
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getStorageList()
    storageList.value = res.data
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(formState, { name: '', type: '', remark: '', is_default_checked: false })
  Object.keys(configForm).forEach(key => delete configForm[key])
  modalVisible.value = true
}

const handleEdit = (record: Storage) => {
  isEdit.value = true
  currentId.value = record.id
  Object.assign(formState, {
    name: record.name,
    type: record.type,
    remark: record.remark,
    is_default_checked: record.is_default === 1,
  })
  // 解析配置
  try {
    const config = JSON.parse(record.config)
    Object.keys(configForm).forEach(key => delete configForm[key])
    Object.assign(configForm, config)
  } catch {
    // ignore
  }
  modalVisible.value = true
}

const handleTypeChange = () => {
  // 清空配置
  Object.keys(configForm).forEach(key => delete configForm[key])
}

const handleTest = async () => {
  if (!formState.type) {
    message.error('请选择存储类型')
    return
  }
  testing.value = true
  try {
    await testStorage({
      name: formState.name,
      type: formState.type as StorageType,
      config: JSON.stringify(configForm),
    })
    message.success('连接测试成功')
  } catch {
    // 错误已在拦截器中处理
  } finally {
    testing.value = false
  }
}

const handleSubmit = async () => {
  if (!formState.name || !formState.type) {
    message.error('请填写必填项')
    return
  }
  submitting.value = true
  try {
    const data = {
      name: formState.name,
      type: formState.type as StorageType,
      config: JSON.stringify(configForm),
      is_default: formState.is_default_checked ? 1 : 0,
      remark: formState.remark,
      status: 1,
    }
    if (isEdit.value) {
      await updateStorage(currentId.value, data)
      message.success('更新成功')
    } else {
      await createStorage(data)
      message.success('创建成功')
    }
    modalVisible.value = false
    fetchList()
  } finally {
    submitting.value = false
  }
}

const handleSetDefault = async (record: Storage) => {
  await setDefaultStorage(record.id)
  message.success('设置成功')
  fetchList()
}

const handleDelete = async (record: Storage) => {
  await deleteStorage(record.id)
  message.success('删除成功')
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.storage-page {
  height: 100%;
}
</style>
