<template>
  <a-drawer
    v-model:open="visible"
    title="富文本 Demo"
    width="760"
    placement="right"
    :destroy-on-close="true"
    @close="handleCancel"
  >
    <div class="rich-text-demo-drawer" @paste.stop>
      <a-form layout="vertical">
        <a-form-item label="标题">
          <a-input
            v-model:value="formState.title"
            :maxlength="80"
            show-count
            placeholder="请输入标题"
          />
        </a-form-item>

        <a-form-item label="内容">
          <RichTextEditor
            v-if="visible"
            v-model="formState.content"
            placeholder="请输入富文本内容"
          />
        </a-form-item>
      </a-form>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <a-button @click="handleCancel">取消</a-button>
        <a-button type="primary" @click="handleSubmit">保存</a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { message } from 'ant-design-vue'
import RichTextEditor from '@/components/RichTextEditor.vue'

interface RichTextDemoValue {
  title: string
  content: string
}

const props = withDefaults(defineProps<{
  open: boolean
  initialValue?: RichTextDemoValue
}>(), {
  initialValue: () => ({
    title: '',
    content: '',
  }),
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: RichTextDemoValue): void
}>()

const visible = computed({
  get: () => props.open,
  set: (value: boolean) => emit('update:open', value),
})

const formState = reactive<RichTextDemoValue>({
  title: '',
  content: '',
})

const resetForm = () => {
  formState.title = props.initialValue.title || ''
  formState.content = props.initialValue.content || ''
}

watch(
  () => props.open,
  (open) => {
    if (open) {
      resetForm()
    }
  },
  { immediate: true },
)

const handleCancel = () => {
  visible.value = false
}

const handleSubmit = () => {
  if (!formState.title.trim()) {
    message.warning('请输入标题')
    return
  }

  emit('submit', {
    title: formState.title.trim(),
    content: formState.content,
  })
  message.success('内容已保存')
  visible.value = false
}
</script>

<style scoped>
.rich-text-demo-drawer {
  padding-bottom: 16px;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
