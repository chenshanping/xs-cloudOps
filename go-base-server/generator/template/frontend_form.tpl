<template>
  <a-drawer
    v-model:open="visible"
    :title="title"
    :width="{{if .HasEditor}}900{{else}}600{{end}}"
    :destroy-on-close="true"
    @close="handleClose"
  >
{{- if .HasEditor}}
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="basic" tab="基础信息">
        <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 4 }">
{{- if .LinkToUser}}
          <a-form-item v-if="!profileMode" label="关联用户" name="user_id">
            <a-select v-model:value="formState.user_id" placeholder="请选择关联用户" allow-clear show-search :filter-option="filterOption" :disabled="!!props.record">
              <a-select-option v-for="item in userOptions" :key="item.id" :value="item.id">
                {{ "{{" }} item.nickname || item.username {{ "}}" }}
              </a-select-option>
            </a-select>
          </a-form-item>
{{- end}}
{{- range .FormColumns}}
{{- if ne .FormType "editor"}}
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
{{- if .DictType}}
            <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
              <a-select-option v-for="item in {{.JsonName}}DictList" :key="item.value" :value="{{if eq .FieldType "string"}}item.value{{else}}Number(item.value){{end}}">
                {{"{{"}} item.label {{"}}"}}
              </a-select-option>
            </a-select>
{{- else if .SelectOptions}}
            <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
{{- $fieldType := .FieldType}}
{{- range .SelectOptions}}
              <a-select-option {{if eq $fieldType "string"}}value="{{.Value}}"{{else}}:value="{{.Value}}"{{end}}>{{.Label}}</a-select-option>
{{- end}}
            </a-select>
{{- else}}
            <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
              <a-select-option :value="1">启用</a-select-option>
              <a-select-option :value="0">禁用</a-select-option>
            </a-select>
{{- end}}
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
{{- else if eq .FormType "images"}}
          <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_ids">
            <FileUpload accept="image/*" :multiple="true" @success="(f) => addFileToMap('{{.JsonName}}', f, 'image')" />
            <div class="image-list" v-if="imageMap['{{.JsonName}}']?.length">
              <div v-for="(f, i) in imageMap['{{.JsonName}}']" :key="f.id" style="display: inline-block; margin-right: 8px; position: relative;">
                <a-image :src="f.url" :width="60" />
                <a-button type="link" size="small" danger style="position: absolute; top: -8px; right: -8px;" @click="removeFromMap('{{.JsonName}}', i, 'image')">x</a-button>
              </div>
            </div>
          </a-form-item>
{{- else if or (eq .FormType "file") (eq .FormType "upload")}}
          <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_id">
            <FileUpload :multiple="false" @success="(f) => { formState.{{.JsonName}}_file_id = f.id; formState.{{.JsonName}}_url = f.url; formState.{{.JsonName}}_name = f.name }" />
            <div v-if="formState.{{.JsonName}}_url" style="margin-top: 8px;">
              <a-button type="link" size="small" @click="handlePreview(formState.{{.JsonName}}_url, formState.{{.JsonName}}_name)">预览文件</a-button>
            </div>
            <div style="color: #999; font-size: 12px;">再次上传将覆盖当前文件</div>
          </a-form-item>
{{- else if eq .FormType "files"}}
          <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_ids">
            <FileUpload :multiple="true" @success="(f) => addFileToMap('{{.JsonName}}', f, 'file')" />
            <div class="file-list" v-if="fileMap['{{.JsonName}}']?.length">
              <div v-for="(f, i) in fileMap['{{.JsonName}}']" :key="f.id" style="display: flex; align-items: center; margin-top: 4px;">
                <a-button type="link" size="small" @click="handlePreview(f.url, f.name)" style="flex: 1; text-align: left; overflow: hidden; text-overflow: ellipsis;">{{"{{"}} f.name || f.url {{"}}"}}</a-button>
                <a-button type="link" size="small" danger @click="removeFromMap('{{.JsonName}}', i, 'file')">x</a-button>
              </div>
            </div>
          </a-form-item>
{{- else}}
          <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
            <a-input v-model:value="formState.{{.JsonName}}" placeholder="请输入{{.Comment}}" />
          </a-form-item>
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
          <a-form-item label="{{.Comment}}" name="{{.ForeignKeyJson}}">
            <a-select v-model:value="formState.{{.ForeignKeyJson}}" placeholder="请选择{{.Comment}}" allow-clear show-search :filter-option="filterOption">
              <a-select-option v-for="item in {{.JsonName}}Options" :key="item.id" :value="item.id">
{{- if .UseOptionsApi}}
                {{"{{"}} item.name {{"}}"}}{{"{{"}} item.count !== undefined ? ` (${item.count})` : '' {{"}}"}}
{{- else}}
                {{"{{"}} item.{{.DisplayField}} {{"}}"}}
{{- end}}
              </a-select-option>
            </a-select>
          </a-form-item>
{{- else if eq .RelationType "many2many"}}
          <a-form-item label="{{.Comment}}" name="{{.JsonName}}_ids">
            <a-select v-model:value="formState.{{.JsonName}}_ids" mode="multiple" placeholder="请选择{{.Comment}}" allow-clear show-search :filter-option="filterOption">
              <a-select-option v-for="item in {{.JsonName}}Options" :key="item.id" :value="item.id">
{{- if .UseOptionsApi}}
                {{"{{"}} item.name {{"}}"}}{{"{{"}} item.count !== undefined ? ` (${item.count})` : '' {{"}}"}}
{{- else}}
                {{"{{"}} item.{{.DisplayField}} {{"}}"}}
{{- end}}
              </a-select-option>
            </a-select>
          </a-form-item>
{{- end}}
{{- end}}
        </a-form>
      </a-tab-pane>
{{- range .FormColumns}}
{{- if eq .FormType "editor"}}
      <a-tab-pane key="{{.JsonName}}" tab="{{.Comment}}">
        <RichTextEditor v-model="formState.{{.JsonName}}" placeholder="请输入{{.Comment}}" />
      </a-tab-pane>
{{- end}}
{{- end}}
    </a-tabs>
{{- else}}
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 6 }">
{{- if .LinkToUser}}
      <a-form-item v-if="!profileMode" label="关联用户" name="user_id">
        <a-select v-model:value="formState.user_id" placeholder="请选择关联用户" allow-clear show-search :filter-option="filterOption" :disabled="!!props.record">
          <a-select-option v-for="item in userOptions" :key="item.id" :value="item.id">
            {{ "{{" }} item.nickname || item.username {{ "}}" }}
          </a-select-option>
        </a-select>
      </a-form-item>
{{- end}}
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
{{- if .DictType}}
        <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
          <a-select-option v-for="item in {{.JsonName}}DictList" :key="item.value" :value="{{if eq .FieldType "string"}}item.value{{else}}Number(item.value){{end}}">
            {{"{{"}} item.label {{"}}"}}
          </a-select-option>
        </a-select>
{{- else if .SelectOptions}}
        <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
{{- $fieldType := .FieldType}}
{{- range .SelectOptions}}
          <a-select-option {{if eq $fieldType "string"}}value="{{.Value}}"{{else}}:value="{{.Value}}"{{end}}>{{.Label}}</a-select-option>
{{- end}}
        </a-select>
{{- else}}
        <a-select v-model:value="formState.{{.JsonName}}" placeholder="请选择{{.Comment}}">
          <a-select-option :value="1">启用</a-select-option>
          <a-select-option :value="0">禁用</a-select-option>
        </a-select>
{{- end}}
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
{{- else if eq .FormType "images"}}
      <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_ids">
        <FileUpload accept="image/*" :multiple="true" @success="(f) => addFileToMap('{{.JsonName}}', f, 'image')" />
        <div class="image-list" v-if="imageMap['{{.JsonName}}']?.length">
          <div v-for="(f, i) in imageMap['{{.JsonName}}']" :key="f.id" style="display: inline-block; margin-right: 8px; position: relative;">
            <a-image :src="f.url" :width="60" />
            <a-button type="link" size="small" danger style="position: absolute; top: -8px; right: -8px;" @click="removeFromMap('{{.JsonName}}', i, 'image')">x</a-button>
          </div>
        </div>
      </a-form-item>
{{- else if or (eq .FormType "file") (eq .FormType "upload")}}
      <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_id">
        <FileUpload :multiple="false" @success="(f) => { formState.{{.JsonName}}_file_id = f.id; formState.{{.JsonName}}_url = f.url; formState.{{.JsonName}}_name = f.name }" />
        <div v-if="formState.{{.JsonName}}_url" style="margin-top: 8px;">
          <a-button type="link" size="small" @click="handlePreview(formState.{{.JsonName}}_url, formState.{{.JsonName}}_name)">预览文件</a-button>
        </div>
        <div style="color: #999; font-size: 12px;">再次上传将覆盖当前文件</div>
      </a-form-item>
{{- else if eq .FormType "files"}}
      <a-form-item label="{{.Comment}}" name="{{.JsonName}}_file_ids">
        <FileUpload :multiple="true" @success="(f) => addFileToMap('{{.JsonName}}', f, 'file')" />
        <div class="file-list" v-if="fileMap['{{.JsonName}}']?.length">
          <div v-for="(f, i) in fileMap['{{.JsonName}}']" :key="f.id" style="display: flex; align-items: center; margin-top: 4px;">
            <a-button type="link" size="small" @click="handlePreview(f.url, f.name)" style="flex: 1; text-align: left; overflow: hidden; text-overflow: ellipsis;">{{"{{"}} f.name || f.url {{"}}"}}</a-button>
            <a-button type="link" size="small" danger @click="removeFromMap('{{.JsonName}}', i, 'file')">x</a-button>
          </div>
        </div>
      </a-form-item>
{{- else}}
      <a-form-item label="{{.Comment}}" name="{{.JsonName}}">
        <a-input v-model:value="formState.{{.JsonName}}" placeholder="请输入{{.Comment}}" />
      </a-form-item>
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
      <a-form-item label="{{.Comment}}" name="{{.ForeignKeyJson}}">
        <a-select v-model:value="formState.{{.ForeignKeyJson}}" placeholder="请选择{{.Comment}}" allow-clear show-search :filter-option="filterOption">
          <a-select-option v-for="item in {{.JsonName}}Options" :key="item.id" :value="item.id">
{{- if .UseOptionsApi}}
            {{"{{"}} item.name {{"}}"}}{{"{{"}} item.count !== undefined ? ` (${item.count})` : '' {{"}}"}}
{{- else}}
            {{"{{"}} item.{{.DisplayField}} {{"}}"}}
{{- end}}
          </a-select-option>
        </a-select>
      </a-form-item>
{{- else if eq .RelationType "many2many"}}
      <a-form-item label="{{.Comment}}" name="{{.JsonName}}_ids">
        <a-select v-model:value="formState.{{.JsonName}}_ids" mode="multiple" placeholder="请选择{{.Comment}}" allow-clear show-search :filter-option="filterOption">
          <a-select-option v-for="item in {{.JsonName}}Options" :key="item.id" :value="item.id">
{{- if .UseOptionsApi}}
            {{"{{"}} item.name {{"}}"}}{{"{{"}} item.count !== undefined ? ` (${item.count})` : '' {{"}}"}}
{{- else}}
            {{"{{"}} item.{{.DisplayField}} {{"}}"}}
{{- end}}
          </a-select-option>
        </a-select>
      </a-form-item>
{{- end}}
{{- end}}
    </a-form>
{{- end}}

    <template #footer>
      <a-space>
        <a-button @click="handleClose">取消</a-button>
        <a-button type="primary" :loading="submitting" @click="handleSubmit">确定</a-button>
      </a-space>
    </template>

{{- if .HasFiles}}
    <!-- 文件预览 -->
    <FilePreview
      v-model:open="previewVisible"
      :url="previewUrl"
      :name="previewName"
      :ext="previewExt"
    />
{{- end}}
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, watch, nextTick{{if .LinkToUser}}, computed{{end}} } from 'vue'
import { message, type FormInstance } from 'ant-design-vue'
import type { Rule } from 'ant-design-vue/es/form'
import ImageUpload from '@/components/ImageUpload.vue'
import FileUpload from '@/components/FileUpload.vue'
{{- if .HasFiles}}
import FilePreview from '@/components/FilePreview.vue'
{{- end}}
{{- if .HasEditor}}
import RichTextEditor from '@/components/RichTextEditor.vue'
{{- end}}
import { create{{.ModelName}}, update{{.ModelName}}{{if .LinkToUser}}, saveMy{{.ModelName}}{{end}} } from '@/api/{{.ModuleName}}'
import { {{.ModelName}} } from '@/types/{{.ModuleName}}'
{{- if .HasDictSelect}}
import { getDictDataByType } from '@/api/dict'
{{- end}}

interface Props {
  open: boolean
  record?: {{.ModelName}} | null
{{- if .LinkToUser}}
  userOptions?: any[]
  profileMode?: boolean  // 个人中心模式，不显示用户选择器，调用 saveMy API
{{- end}}
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
  {{.JsonName}}Options?: any[]
{{- end}}
{{- end}}
}

const props = withDefaults(defineProps<Props>(), {
  record: null,
{{- if .LinkToUser}}
  userOptions: () => [],
  profileMode: false,
{{- end}}
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
  {{.JsonName}}Options: () => [],
{{- end}}
{{- end}}
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const visible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()
{{- if .HasEditor}}
const activeTab = ref('basic')
{{- end}}
{{- if .HasFiles}}

// 文件预览
const previewVisible = ref(false)
const previewUrl = ref('')
const previewName = ref('')
const previewExt = ref('')

const handlePreview = (url: string, name?: string) => {
  previewUrl.value = url
  previewName.value = name || url.split('/').pop() || 'file'
  previewExt.value = url.split('.').pop()?.toLowerCase() || ''
  previewVisible.value = true
}
{{- end}}

const title = ref('新增{{.Description}}')

{{- if .HasDictSelect}}
// 字典数据
{{- range .FormColumns}}
{{- if and (eq .FormType "select") .DictType}}
const {{.JsonName}}DictList = ref<any[]>([])
{{- end}}
{{- end}}

// 获取字典数据
const fetchDictData = async () => {
{{- range .FormColumns}}
{{- if and (eq .FormType "select") .DictType}}
  try {
    const res{{.JsonName}} = await getDictDataByType('{{.DictType}}')
    {{.JsonName}}DictList.value = res{{.JsonName}}.data || []
  } catch (e) {
    console.error('Failed to fetch dict {{.DictType}}:', e)
  }
{{- end}}
{{- end}}
}
fetchDictData()
{{- end}}

// 图片多选映射
const imageMap = reactive<Record<string, Array<{id: number, url: string}>>>({
{{- range .FormColumns}}
{{- if eq .FormType "images"}}
  '{{.JsonName}}': [],
{{- end}}
{{- end}}
})

// 文件多选映射
const fileMap = reactive<Record<string, Array<{id: number, url: string, name?: string}>>>({
{{- range .FormColumns}}
{{- if eq .FormType "files"}}
  '{{.JsonName}}': [],
{{- end}}
{{- end}}
})

{{- if or .HasRelations .LinkToUser}}
// 下拉搜索过滤
const filterOption = (input: string, option: any) => {
  return option.children?.[0]?.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0
}
{{- end}}

// 表单初始值
const getInitialFormState = () => ({
{{- if .LinkToUser}}
  user_id: undefined as number | undefined,
{{- end}}
{{- range .FormColumns}}
{{- if eq .FormType "switch"}}
  {{.JsonName}}Checked: true,
{{- else if eq .FormType "number"}}
  {{.JsonName}}: 0,
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id: 0,
  {{.JsonName}}_url: '',
  {{.JsonName}}_name: '',
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
  {{.JsonName}}_file_ids: '',
{{- else if eq .FormType "select"}}
  {{.JsonName}}: undefined as {{if or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint")}}number{{else}}string{{end}} | undefined,
{{- else}}
  {{.JsonName}}: '',
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  {{.ForeignKeyJson}}: undefined as number | undefined,
{{- else if eq .RelationType "many2many"}}
  {{.JsonName}}_ids: [] as number[],
{{- end}}
{{- end}}
})

const formState = reactive(getInitialFormState())

// 表单验证规则
{{- if .LinkToUser}}
const formRules = computed<Record<string, Rule[]>>(() => ({
  // 个人中心模式不需要选择用户
  ...(props.profileMode ? {} : { user_id: [{ required: true, message: '请选择关联用户', trigger: 'change', type: 'number' }] }),
{{- range .FormColumns}}
{{- if .IsRequired}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id: [{ required: true, message: '请上传{{.Comment}}', trigger: 'change', type: 'number', min: 1 }],
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
  {{.JsonName}}_file_ids: [{ required: true, message: '请上传{{.Comment}}', trigger: 'change' }],
{{- else if eq .FormType "number"}}
  {{.JsonName}}: [{ required: true, message: '请输入{{.Comment}}', trigger: 'blur', type: 'number' }],
{{- else if eq .FormType "select"}}
  {{.JsonName}}: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change' }],
{{- else if or (eq .FormType "datetime") (eq .FormType "date")}}
  {{.JsonName}}: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change' }],
{{- else if ne .FormType "switch"}}
  {{.JsonName}}: [{ required: true, message: '请输入{{.Comment}}', trigger: 'blur' }],
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if and .IsRequired (eq .RelationType "belongsTo")}}
  {{.ForeignKeyJson}}: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change', type: 'number' }],
{{- else if and .IsRequired (eq .RelationType "many2many")}}
  {{.JsonName}}_ids: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change', type: 'array' }],
{{- end}}
{{- end}}
}))
{{- else}}
const formRules: Record<string, Rule[]> = {
{{- range .FormColumns}}
{{- if .IsRequired}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
  {{.JsonName}}_file_id: [{ required: true, message: '请上传{{.Comment}}', trigger: 'change', type: 'number', min: 1 }],
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
  {{.JsonName}}_file_ids: [{ required: true, message: '请上传{{.Comment}}', trigger: 'change' }],
{{- else if eq .FormType "number"}}
  {{.JsonName}}: [{ required: true, message: '请输入{{.Comment}}', trigger: 'blur', type: 'number' }],
{{- else if eq .FormType "select"}}
  {{.JsonName}}: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change' }],
{{- else if or (eq .FormType "datetime") (eq .FormType "date")}}
  {{.JsonName}}: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change' }],
{{- else if ne .FormType "switch"}}
  {{.JsonName}}: [{ required: true, message: '请输入{{.Comment}}', trigger: 'blur' }],
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if and .IsRequired (eq .RelationType "belongsTo")}}
  {{.ForeignKeyJson}}: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change', type: 'number' }],
{{- else if and .IsRequired (eq .RelationType "many2many")}}
  {{.JsonName}}_ids: [{ required: true, message: '请选择{{.Comment}}', trigger: 'change', type: 'array' }],
{{- end}}
{{- end}}
}
{{- end}}

{{- if .HasMultiFiles}}
// 添加文件到映射
const addFileToMap = (field: string, file: {id: number, url: string, name?: string}, type: 'image' | 'file') => {
  if (type === 'image') {
    imageMap[field].push(file)
    ;(formState as any)[field + '_file_ids'] = imageMap[field].map(f => f.id).join(',')
  } else {
    fileMap[field].push(file)
    ;(formState as any)[field + '_file_ids'] = fileMap[field].map(f => f.id).join(',')
  }
}

// 从映射中移除
const removeFromMap = (field: string, index: number, type: 'image' | 'file') => {
  if (type === 'image') {
    imageMap[field].splice(index, 1)
    ;(formState as any)[field + '_file_ids'] = imageMap[field].map(f => f.id).join(',')
  } else {
    fileMap[field].splice(index, 1)
    ;(formState as any)[field + '_file_ids'] = fileMap[field].map(f => f.id).join(',')
  }
}
{{- end}}

// 重置表单
const resetForm = () => {
  Object.assign(formState, getInitialFormState())
  // 重置映射
{{- range .FormColumns}}
{{- if eq .FormType "images"}}
  imageMap['{{.JsonName}}'] = []
{{- else if eq .FormType "files"}}
  fileMap['{{.JsonName}}'] = []
{{- end}}
{{- end}}
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 填充表单
const fillForm = (record: {{.ModelName}}) => {
  Object.assign(formState, {
{{- if .LinkToUser}}
    user_id: record.user_id,
{{- end}}
{{- range .FormColumns}}
{{- if eq .FormType "switch"}}
{{- if .SwitchValues}}
    {{.JsonName}}Checked: record.{{.JsonName}} == '{{.SwitchValues.ActiveValue}}',
{{- else}}
    {{.JsonName}}Checked: record.{{.JsonName}} === 1,
{{- end}}
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
    {{.JsonName}}_file_id: record.{{.JsonName}}_file_id || 0,
    {{.JsonName}}_url: record.{{.JsonName}}_url || '',
    {{.JsonName}}_name: '',
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
    {{.JsonName}}_file_ids: record.{{.JsonName}}_file_ids || '',
{{- else}}
    {{.JsonName}}: record.{{.JsonName}},
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
    {{.ForeignKeyJson}}: record.{{.ForeignKeyJson}},
{{- else if eq .RelationType "many2many"}}
    {{.JsonName}}_ids: record.{{.JsonName}}?.map(r => r.id) || [],
{{- end}}
{{- end}}
  })

  // 填充多图片/文件映射
{{- range .FormColumns}}
{{- if eq .FormType "images"}}
  if (record.{{.JsonName}}_file_ids && record.{{.JsonName}}_urls) {
    const ids = record.{{.JsonName}}_file_ids.split(',').map(Number)
    imageMap['{{.JsonName}}'] = ids.map((id, i) => ({ id, url: record.{{.JsonName}}_urls[i] || '' }))
  } else {
    imageMap['{{.JsonName}}'] = []
  }
{{- else if eq .FormType "files"}}
  if (record.{{.JsonName}}_file_ids && record.{{.JsonName}}_urls) {
    const ids = record.{{.JsonName}}_file_ids.split(',').map(Number)
    const names = record.{{.JsonName}}_names || []
    fileMap['{{.JsonName}}'] = ids.map((id, i) => ({ id, url: record.{{.JsonName}}_urls[i] || '', name: names[i] || '' }))
  } else {
    fileMap['{{.JsonName}}'] = []
  }
{{- end}}
{{- end}}
}

// 监听 open
watch(() => props.open, (val) => {
  visible.value = val
  if (val) {
    if (props.record?.id) {
      // 编辑模式
      title.value = '编辑{{.Description}}'
      fillForm(props.record)
    } else if (props.record) {
      // 复制模式（有数据但无 id）
      title.value = '新增{{.Description}}'
      fillForm(props.record)
    } else {
      // 新增模式
      title.value = '新增{{.Description}}'
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
{{- if .LinkToUser}}
    user_id: formState.user_id,
{{- end}}
{{- range .FormColumns}}
{{- if eq .FormType "switch"}}
{{- if .SwitchValues}}
    {{.JsonName}}: formState.{{.JsonName}}Checked ? '{{.SwitchValues.ActiveValue}}' : '{{.SwitchValues.InactiveValue}}',
{{- else}}
    {{.JsonName}}: formState.{{.JsonName}}Checked ? 1 : 0,
{{- end}}
{{- else if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
    {{.JsonName}}_file_id: formState.{{.JsonName}}_file_id,
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
    {{.JsonName}}_file_ids: formState.{{.JsonName}}_file_ids || '',
{{- else}}
    {{.JsonName}}: formState.{{.JsonName}},
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
    {{.ForeignKeyJson}}: formState.{{.ForeignKeyJson}},
{{- else if eq .RelationType "many2many"}}
    {{.JsonName}}_ids: formState.{{.JsonName}}_ids,
{{- end}}
{{- end}}
  }

  submitting.value = true
  try {
{{- if .LinkToUser}}
    if (props.profileMode) {
      // 个人中心模式，调用 saveMy API
      await saveMy{{.ModelName}}(data)
      message.success('保存成功')
    } else if (props.record?.id) {
{{- else}}
    if (props.record?.id) {
{{- end}}
      await update{{.ModelName}}(props.record.id, data)
      message.success('更新成功')
    } else {
      await create{{.ModelName}}(data)
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
