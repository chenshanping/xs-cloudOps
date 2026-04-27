<template>
  <div style="border: 1px solid #ccc">
    <Toolbar
      style="border-bottom: 1px solid #ccc"
      :editor="editorRef"
      :defaultConfig="toolbarConfig"
      mode="default"
    />
    <Editor
      style="height: 300px; overflow-y: hidden"
      v-model="valueHtml"
      :defaultConfig="editorConfig"
      mode="default"
      @onCreated="handleCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, shallowRef, watch, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'
import '@wangeditor/editor/dist/css/style.css'
import { calculateMD5, multipartUpload } from '@/utils/upload'

interface Props {
  modelValue?: string
  placeholder?: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入内容...',
  disabled: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef<IDomEditor>()

// 内容 HTML
const valueHtml = ref(props.modelValue)

// 工具栏配置
const toolbarConfig: Partial<IToolbarConfig> = {
  excludeKeys: []
}

// 自定义上传图片
async function customUploadImage(file: File, insertFn: (url: string, alt?: string, href?: string) => void) {
  try {
    const md5 = await calculateMD5(file)
    const result = await multipartUpload(file, md5)
    insertFn(result.url, result.name || '', result.url)
  } catch (error) {
    message.error('图片上传失败: ' + (error as Error).message)
  }
}

// 自定义上传视频
async function customUploadVideo(file: File, insertFn: (url: string, poster?: string) => void) {
  try {
    const md5 = await calculateMD5(file)
    const result = await multipartUpload(file, md5)
    insertFn(result.url, '')
  } catch (error) {
    message.error('视频上传失败: ' + (error as Error).message)
  }
}

// 编辑器配置
const editorConfig: Partial<IEditorConfig> = {
  placeholder: props.placeholder,
  readOnly: props.disabled,
  MENU_CONF: {
    uploadImage: {
      maxFileSize: 10 * 1024 * 1024, // 10M
      allowedFileTypes: ['image/*'],
      customUpload: customUploadImage
    },
    uploadVideo: {
      maxFileSize: 100 * 1024 * 1024, // 100M
      allowedFileTypes: ['video/*'],
      customUpload: customUploadVideo
    }
  }
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

const handleCreated = (editor: IDomEditor) => {
  editorRef.value = editor
}

// 监听 props 变化
watch(() => props.modelValue, (val) => {
  if (val !== valueHtml.value) {
    valueHtml.value = val
  }
})

// 监听内容变化，同步到父组件
watch(valueHtml, (val) => {
  emit('update:modelValue', val)
})

// 监听禁用状态
watch(() => props.disabled, (val) => {
  if (editorRef.value) {
    if (val) {
      editorRef.value.disable()
    } else {
      editorRef.value.enable()
    }
  }
})
</script>

<style scoped>
:deep(.w-e-text-container) {
  background-color: #fff;
}
</style>
