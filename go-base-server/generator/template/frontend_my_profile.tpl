<template>
  <div class="my-{{.ModuleName}}-page">
    <a-card :title="pageTitle">
      <template #extra>
        <a-button type="link" @click="router.push('/profile')">
          <ArrowLeftOutlined /> 返回个人中心
        </a-button>
      </template>

      <a-spin :spinning="loading">
        <a-empty v-if="!hasProfile && !loading" description="您还没有完善{{.Description}}信息">
          <a-button type="primary" @click="handleEdit">立即完善</a-button>
        </a-empty>

        <template v-if="hasProfile">
          <!-- 信息展示 -->
          <a-descriptions :column="2" bordered v-if="!editing && profileData">
{{- range .FormColumns}}
{{- if eq .FormType "image"}}
            <a-descriptions-item label="{{.Comment}}">
              <a-image v-if="profileData.{{.JsonName}}_url" :src="profileData.{{.JsonName}}_url" :width="80" />
              <span v-else>-</span>
            </a-descriptions-item>
{{- else if eq .FormType "switch"}}
            <a-descriptions-item label="{{.Comment}}">
{{- if .SwitchValues}}
              <a-tag :color="profileData.{{.JsonName}} == '{{.SwitchValues.ActiveValue}}' ? 'green' : 'red'">
                {{"{{"}} profileData.{{.JsonName}} == '{{.SwitchValues.ActiveValue}}' ? '{{.SwitchValues.ActiveText}}' : '{{.SwitchValues.InactiveText}}' {{"}}"}}
              </a-tag>
{{- else}}
              <a-tag :color="profileData.{{.JsonName}} === 1 ? 'green' : 'red'">
                {{"{{"}} profileData.{{.JsonName}} === 1 ? '是' : '否' {{"}}"}}
              </a-tag>
{{- end}}
            </a-descriptions-item>
{{- else if eq .FormType "select"}}
            <a-descriptions-item label="{{.Comment}}">
              <a-tag>{{"{{"}} {{.JsonName}}Labels[profileData.{{.JsonName}}] || profileData.{{.JsonName}} || '-' {{"}}"}}</a-tag>
            </a-descriptions-item>
{{- else if eq .FormType "textarea"}}
            <a-descriptions-item label="{{.Comment}}" :span="2">
              <div style="white-space: pre-wrap">{{"{{"}} profileData.{{.JsonName}} || '-' {{"}}"}}</div>
            </a-descriptions-item>
{{- else if eq .FormType "editor"}}
            <a-descriptions-item label="{{.Comment}}" :span="2">
              <div v-html="profileData.{{.JsonName}} || '-'"></div>
            </a-descriptions-item>
{{- else}}
            <a-descriptions-item label="{{.Comment}}">{{"{{"}} profileData.{{.JsonName}} || '-' {{"}}"}}</a-descriptions-item>
{{- end}}
{{- end}}
          </a-descriptions>

          <div style="margin-top: 16px; text-align: center" v-if="!editing">
            <a-button type="primary" @click="handleEdit">编辑信息</a-button>
          </div>

          <!-- 编辑表单 -->
          <a-form v-if="editing" ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 4 }" :wrapper-col="{ span: 16 }">
{{- range .FormColumns}}
{{- if eq .FormType "input"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-input v-model:value="formState.{{.JsonName}}" placeholder="请输入{{.Comment}}" />
            </a-form-item>
{{- else if eq .FormType "textarea"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-textarea v-model:value="formState.{{.JsonName}}" :rows="4" placeholder="请输入{{.Comment}}" />
            </a-form-item>
{{- else if eq .FormType "number"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-input-number v-model:value="formState.{{.JsonName}}" style="width: 100%" placeholder="请输入{{.Comment}}" />
            </a-form-item>
{{- else if eq .FormType "select"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
{{- if .SelectOptions}}
{{- $fieldType := .FieldType}}
{{- range .SelectOptions}}
                <a-select-option {{if eq $fieldType "string"}}value="{{.Value}}"{{else}}:value="{{.Value}}"{{end}}>{{.Label}}</a-select-option>
{{- end}}
{{- else}}
                <a-select-option :value="1">启用</a-select-option>
                <a-select-option :value="0">禁用</a-select-option>
{{- end}}
              </a-select>
            </a-form-item>
{{- else if eq .FormType "switch"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
{{- if .SwitchValues}}
              <a-switch v-model:checked="formState.{{.JsonName}}Checked" checked-children="{{.SwitchValues.ActiveText}}" un-checked-children="{{.SwitchValues.InactiveText}}" />
{{- else}}
              <a-switch v-model:checked="formState.{{.JsonName}}Checked" />
{{- end}}
            </a-form-item>
{{- else if eq .FormType "datetime"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-date-picker v-model:value="formState.{{.JsonName}}" show-time style="width: 100%" placeholder="请选择{{.Comment}}" />
            </a-form-item>
{{- else if eq .FormType "date"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-date-picker v-model:value="formState.{{.JsonName}}" style="width: 100%" placeholder="请选择{{.Comment}}" />
            </a-form-item>
{{- else if eq .FormType "image"}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_id">
              <ImageUpload v-model:fileId="formState.{{.JsonName}}_file_id" v-model:url="formState.{{.JsonName}}_url" />
            </a-form-item>
{{- else}}
            <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
              <a-input v-model:value="formState.{{.JsonName}}" placeholder="请输入{{.Comment}}" />
            </a-form-item>
{{- end}}
{{- end}}

            <a-form-item :wrapper-col="{ offset: 4 }">
              <a-space>
                <a-button type="primary" :loading="submitting" @click="handleSubmit">保存</a-button>
                <a-button @click="handleCancel">取消</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </template>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message, type FormInstance } from 'ant-design-vue'
import type { Rule } from 'ant-design-vue/es/form'
import { useRouter } from 'vue-router'
import { ArrowLeftOutlined } from '@ant-design/icons-vue'
import ImageUpload from '@/components/ImageUpload.vue'
import { getMy{{.ModelName}}, saveMy{{.ModelName}} } from '@/api/{{.ModuleName}}'
import { {{.ModelName}} } from '@/types/{{.ModuleName}}'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const editing = ref(false)
const formRef = ref<FormInstance>()
const profileData = ref<{{.ModelName}} | null>(null)

const hasProfile = computed(() => profileData.value !== null)
const pageTitle = computed(() => editing.value ? '编辑{{.Description}}信息' : '我的{{.Description}}信息')

{{- range .FormColumns}}
{{- if eq .FormType "select"}}
{{- if .SelectOptions}}
// {{.Comment}}选项映射
const {{.JsonName}}Labels: Record<string, string> = {
{{- range .SelectOptions}}
  '{{.Value}}': '{{.Label}}',
{{- end}}
}
{{- else}}
const {{.JsonName}}Labels: Record<string, string> = { '1': '启用', '0': '禁用' }
{{- end}}
{{- end}}
{{- end}}

// 表单初始值
const getInitialFormState = () => ({
{{- range .FormColumns}}
{{- if eq .FormType "switch"}}
  {{.JsonName}}Checked: true,
{{- else if eq .FormType "number"}}
  {{.JsonName}}: 0,
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id: 0,
  {{.JsonName}}_url: '',
{{- else}}
  {{.JsonName}}: '',
{{- end}}
{{- end}}
})

const formState = reactive(getInitialFormState())

// 表单验证规则
const formRules: Record<string, Rule[]> = {
{{- range .FormColumns}}
{{- if .IsRequired}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id: [{ required: true, message: '请上传{{.Comment}}', trigger: 'change', type: 'number', min: 1 }],
{{- else if eq .FormType "number"}}
  {{.JsonName}}: [{ required: true, message: '请输入{{.Comment}}', trigger: 'blur', type: 'number' }],
{{- else if ne .FormType "switch"}}
  {{.JsonName}}: [{ required: true, message: '请输入{{.Comment}}', trigger: 'blur' }],
{{- end}}
{{- end}}
{{- end}}
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getMy{{.ModelName}}()
    profileData.value = res.data
    if (res.data) {
      fillForm(res.data)
    }
  } catch (e) {
    console.error('获取{{.Description}}信息失败', e)
  } finally {
    loading.value = false
  }
}

// 填充表单
const fillForm = (data: {{.ModelName}}) => {
  Object.assign(formState, {
{{- range .FormColumns}}
{{- if eq .FormType "switch"}}
{{- if .SwitchValues}}
    {{.JsonName}}Checked: data.{{.JsonName}} == '{{.SwitchValues.ActiveValue}}',
{{- else}}
    {{.JsonName}}Checked: data.{{.JsonName}} === 1,
{{- end}}
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
    {{.JsonName}}_file_id: data.{{.JsonName}}_file_id || 0,
    {{.JsonName}}_url: data.{{.JsonName}}_url || '',
{{- else}}
    {{.JsonName}}: data.{{.JsonName}},
{{- end}}
{{- end}}
  })
}

// 重置表单
const resetForm = () => {
  Object.assign(formState, getInitialFormState())
}

// 编辑
const handleEdit = () => {
  if (profileData.value) {
    fillForm(profileData.value)
  } else {
    resetForm()
  }
  editing.value = true
}

// 取消
const handleCancel = () => {
  editing.value = false
  if (profileData.value) {
    fillForm(profileData.value)
  }
}

// 提交
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  const data = {
{{- range .FormColumns}}
{{- if eq .FormType "switch"}}
{{- if .SwitchValues}}
    {{.JsonName}}: formState.{{.JsonName}}Checked ? '{{.SwitchValues.ActiveValue}}' : '{{.SwitchValues.InactiveValue}}',
{{- else}}
    {{.JsonName}}: formState.{{.JsonName}}Checked ? 1 : 0,
{{- end}}
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
    {{.JsonName}}_file_id: formState.{{.JsonName}}_file_id,
{{- else}}
    {{.JsonName}}: formState.{{.JsonName}},
{{- end}}
{{- end}}
  }

  submitting.value = true
  try {
    await saveMy{{.ModelName}}(data)
    message.success('保存成功')
    editing.value = false
    fetchData()
  } finally {
    submitting.value = false
  }
}

onMounted(() => fetchData())
</script>

<style scoped>
.my-{{.ModuleName}}-page {
  padding: 0;
}
</style>
