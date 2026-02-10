<template>
  <a-drawer
    v-model:open="visible"
    :title="title"
    :width="600"
    :destroy-on-close="true"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 6 }">
      <a-form-item label="产品名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入产品名称" />
      </a-form-item>
      <a-form-item label="产品数量" name="num">
        <a-input-number v-model:value="formState.num" style="width: 100%" placeholder="请输入产品数量" />
      </a-form-item>
      <a-form-item label="产品单价" name="price">
        <a-input-number v-model:value="formState.price" style="width: 100%" placeholder="请输入产品单价" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-select v-model:value="formState.status" placeholder="请选择状态">
          <a-select-option v-for="item in statusDictList" :key="item.value" :value="item.value">
            {{ item.label }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="产品类型" name="type_id">
        <a-select v-model:value="formState.type_id" placeholder="请选择产品类型" allow-clear show-search :filter-option="filterOption">
          <a-select-option v-for="item in product_typeOptions" :key="item.id" :value="item.id">
            {{ item.name }}{{ item.count !== undefined ? ` (${item.count})` : '' }}
          </a-select-option>
        </a-select>
      </a-form-item>
    </a-form>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">取消</a-button>
        <a-button type="primary" :loading="submitting" @click="handleSubmit">确定</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, watch, nextTick } from 'vue'
import { message, type FormInstance } from 'ant-design-vue'
import type { Rule } from 'ant-design-vue/es/form'
import ImageUpload from '@/components/ImageUpload.vue'
import FileUpload from '@/components/FileUpload.vue'
import { createProduct, updateProduct } from '@/api/product'
import { Product } from '@/types/product'
import { getDictDataByType } from '@/api/dict'

interface Props {
  open: boolean
  record?: Product | null
  product_typeOptions?: any[]
}

const props = withDefaults(defineProps<Props>(), {
  record: null,
  product_typeOptions: () => [],
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const visible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const title = ref('新增产品信息')
// 字典数据
const statusDictList = ref<any[]>([])

// 获取字典数据
const fetchDictData = async () => {
  try {
    const resstatus = await getDictDataByType('common_status')
    statusDictList.value = resstatus.data || []
  } catch (e) {
    console.error('Failed to fetch dict common_status:', e)
  }
}
fetchDictData()

// 图片多选映射
const imageMap = reactive<Record<string, Array<{id: number, url: string}>>>({
})

// 文件多选映射
const fileMap = reactive<Record<string, Array<{id: number, url: string, name?: string}>>>({
})
// 下拉搜索过滤
const filterOption = (input: string, option: any) => {
  return option.children?.[0]?.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

// 表单初始值
const getInitialFormState = () => ({
  name: '',
  num: 0,
  price: 0,
  status: undefined as string | undefined,
  type_id: undefined as number | undefined,
})

const formState = reactive(getInitialFormState())

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  type_id: [{ required: true, message: '请选择产品类型', trigger: 'change', type: 'number' }],
}

// 重置表单
const resetForm = () => {
  Object.assign(formState, getInitialFormState())
  // 重置映射
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 填充表单
const fillForm = (record: Product) => {
  Object.assign(formState, {
    name: record.name,
    num: record.num,
    price: record.price,
    status: record.status,
    type_id: record.type_id,
  })

  // 填充多图片/文件映射
}

// 监听 open
watch(() => props.open, (val) => {
  visible.value = val
  if (val) {
    if (props.record?.id) {
      // 编辑模式
      title.value = '编辑产品信息'
      fillForm(props.record)
    } else if (props.record) {
      // 复制模式（有数据但无 id）
      title.value = '新增产品信息'
      fillForm(props.record)
    } else {
      // 新增模式
      title.value = '新增产品信息'
      resetForm()
    }
  }
}, { immediate: true })

// 监听内部 visible
watch(visible, (val) => {
  emit('update:open', val)
})

const handleClose = () => {
  visible.value = false
  resetForm()
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  const data = {
    name: formState.name,
    num: formState.num,
    price: formState.price,
    status: formState.status,
    type_id: formState.type_id,
  }

  submitting.value = true
  try {
    if (props.record?.id) {
      await updateProduct(props.record.id, data)
      message.success('更新成功')
    } else {
      await createProduct(data)
      message.success('创建成功')
    }
    visible.value = false
    resetForm()
    emit('success')
  } catch {
    // 错误已由 request 拦截器统一处理
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.image-list { margin-top: 8px; }
.file-list { margin-top: 8px; }
</style>
