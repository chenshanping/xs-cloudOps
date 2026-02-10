<template>
  <a-drawer
    v-model:open="visible"
    :title="title"
    :width="600"
    :destroy-on-close="true"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 6 }">
      <a-form-item label="产品类型名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入产品类型名称" />
      </a-form-item>
      <a-form-item label="类型图标" name="icon">
        <a-input v-model:value="formState.icon" placeholder="请输入类型图标" />
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
import { createProductType, updateProductType } from '@/api/productType'
import { ProductType } from '@/types/productType'

interface Props {
  open: boolean
  record?: ProductType | null
}

const props = withDefaults(defineProps<Props>(), {
  record: null,
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const visible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const title = ref('新增产品类型')

// 图片多选映射
const imageMap = reactive<Record<string, Array<{id: number, url: string}>>>({
})

// 文件多选映射
const fileMap = reactive<Record<string, Array<{id: number, url: string, name?: string}>>>({
})

// 表单初始值
const getInitialFormState = () => ({
  name: '',
  icon: '',
})

const formState = reactive(getInitialFormState())

// 表单验证规则
const formRules: Record<string, Rule[]> = {
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
const fillForm = (record: ProductType) => {
  Object.assign(formState, {
    name: record.name,
    icon: record.icon,
  })

  // 填充多图片/文件映射
}

// 监听 open
watch(() => props.open, (val) => {
  visible.value = val
  if (val) {
    if (props.record?.id) {
      // 编辑模式
      title.value = '编辑产品类型'
      fillForm(props.record)
    } else if (props.record) {
      // 复制模式（有数据但无 id）
      title.value = '新增产品类型'
      fillForm(props.record)
    } else {
      // 新增模式
      title.value = '新增产品类型'
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
    icon: formState.icon,
  }

  submitting.value = true
  try {
    if (props.record?.id) {
      await updateProductType(props.record.id, data)
      message.success('更新成功')
    } else {
      await createProductType(data)
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
