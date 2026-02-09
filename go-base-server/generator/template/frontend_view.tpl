<template>
{{- if .HasTreeLayout}}
  <div class="{{.ModuleName}}-page tree-table-layout">
    <!-- 左侧分类树 -->
    <div class="category-tree-panel">
      <div class="tree-header">
        <span class="tree-title"><FolderOutlined /> {{range .Relations}}{{if .UseTreeLayout}}{{.Comment}}{{end}}{{end}}</span>
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
            <span class="item-count">{{"{{"}} totalCount {{"}}"}}</span>
          </div>
          <div 
            v-for="item in categoryOptions" 
            :key="item.id" 
            class="tree-item"
            :class="{ active: selectedCategoryId === item.id }"
            @click="handleSelectCategory(item.id)"
          >
            <TagOutlined />
            <span>{{"{{"}} item.name {{"}}"}}</span>
            <span class="item-count">{{"{{"}} item.count || 0 {{"}}"}}</span>
          </div>
        </div>
      </a-spin>
    </div>
    <!-- 右侧表格 -->
    <div class="table-panel">
{{- else}}
  <div class="{{.ModuleName}}-page">
{{- end}}
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
{{- if and .HasSearchFields .DataIsolation .OnlyCreatedBySearch}}
      :show-search="isAdmin"
{{- end}}
    >
{{- if .HasSearchFields}}
      <template #search>
{{- range .SearchColumns}}
{{- if eq .FormType "select"}}
        <a-form-item label="{{.Comment}}">
{{- if .DictType}}
          <a-select v-model:value="searchForm.{{.JsonName}}" placeholder="请选择" allowClear style="width: 160px">
            <a-select-option v-for="item in {{.JsonName}}DictList" :key="item.value" :value="{{if eq .FieldType "string"}}item.value{{else}}Number(item.value){{end}}">
              {{"{{"}} item.label {{"}}"}}
            </a-select-option>
          </a-select>
{{- else if .SelectOptions}}
          <a-select v-model:value="searchForm.{{.JsonName}}" placeholder="请选择" allowClear style="width: 160px">
{{- $fieldType := .FieldType}}
{{- range .SelectOptions}}
            <a-select-option {{if eq $fieldType "string"}}value="{{.Value}}"{{else}}:value="{{.Value}}"{{end}}>{{.Label}}</a-select-option>
{{- end}}
          </a-select>
{{- else}}
          <a-select v-model:value="searchForm.{{.JsonName}}" placeholder="请选择" allowClear style="width: 160px">
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">禁用</a-select-option>
          </a-select>
{{- end}}
        </a-form-item>
{{- else}}
        <a-form-item label="{{.Comment}}">
          <a-input v-model:value="searchForm.{{.JsonName}}" placeholder="请输入{{.Comment}}" allowClear style="width: 200px" />
        </a-form-item>
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
{{- if not .UseTreeLayout}}
        <a-form-item label="{{.Comment}}">
          <a-select v-model:value="searchForm.{{.ForeignKeyJson}}" placeholder="请选择{{.Comment}}" allowClear style="width: 180px" show-search :filter-option="filterOption">
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
{{- end}}
{{- if .HasCreatedBy}}
{{- if .DataIsolation}}
        <a-form-item v-if="isAdmin" label="创建人">
{{- else}}
        <a-form-item label="创建人">
{{- end}}
          <a-select v-model:value="searchForm.created_by" placeholder="请选择创建人" allowClear style="width: 180px" show-search :filter-option="filterOption">
            <a-select-option v-for="item in creatorOptions" :key="item.id" :value="item.id">
              {{"{{"}} item.name {{"}}"}}{{"{{"}} item.count !== undefined ? ` (${item.count})` : '' {{"}}"}}
            </a-select-option>
          </a-select>
        </a-form-item>
{{- end}}
      </template>
{{- end}}

      <template #toolbar>
        <a-space>
{{- if .HasMenu}}
          <a-button type="primary" @click="handleAdd" v-permission="'{{.MenuConfig.Permission}}:add'"><PlusOutlined /> 新增</a-button>
          <a-button danger :disabled="selectedRowKeys.length === 0" @click="confirmBatchDelete" v-permission="'{{.MenuConfig.Permission}}:delete'">
            <DeleteOutlined /> 批量删除 {{"{{"}} selectedRowKeys.length > 0 ? `(${selectedRowKeys.length})` : '' {{"}}"}}
          </a-button>
{{- else}}
          <a-button type="primary" @click="handleAdd"><PlusOutlined /> 新增</a-button>
          <a-button danger :disabled="selectedRowKeys.length === 0" @click="confirmBatchDelete">
            <DeleteOutlined /> 批量删除 {{"{{"}} selectedRowKeys.length > 0 ? `(${selectedRowKeys.length})` : '' {{"}}"}}
          </a-button>
{{- end}}
        </a-space>
      </template>

      <template #bodyCell="{ column, record }">
{{- range .ListColumns}}
{{- if eq .FormType "switch"}}
        <template v-if="column.key === '{{.JsonName}}'">
{{- if .SwitchValues}}
          <a-tag :color="record.{{.JsonName}} == '{{.SwitchValues.ActiveValue}}' ? 'green' : 'red'">
            {{"{{"}} record.{{.JsonName}} == '{{.SwitchValues.ActiveValue}}' ? '{{.SwitchValues.ActiveText}}' : '{{.SwitchValues.InactiveText}}' {{"}}"}}
          </a-tag>
{{- else}}
          <a-switch :checked="record.{{.JsonName}} === 1" disabled />
{{- end}}
        </template>
{{- else if eq .FormType "select"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-tag :color="{{.JsonName}}Colors[record.{{.JsonName}}]">
            {{"{{"}} {{.JsonName}}Options[record.{{.JsonName}}] || record.{{.JsonName}} {{"}}"}}
          </a-tag>
        </template>
{{- else if eq .FormType "datetime"}}
        <template v-if="column.key === '{{.JsonName}}'">{{"{{"}} formatTime(record.{{.JsonName}}) {{"}}"}}</template>
{{- else if eq .FormType "image"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-image v-if="record.{{.JsonName}}_url" :src="record.{{.JsonName}}_url" :width="40" />
        </template>
{{- else if eq .FormType "images"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-image-preview-group v-if="record.{{.JsonName}}_urls?.length">
            <a-image v-for="(url, i) in record.{{.JsonName}}_urls.slice(0, 3)" :key="i" :src="url" :width="30" style="margin-right: 4px" />
            <span v-if="record.{{.JsonName}}_urls.length > 3">+{{"{{"}} record.{{.JsonName}}_urls.length - 3 {{"}}"}}</span>
          </a-image-preview-group>
        </template>
{{- else if eq .FormType "files"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-button v-if="record.{{.JsonName}}_urls?.length" type="link" size="small" @click="handlePreviewFiles(record.{{.JsonName}}_urls)">
            {{"{{"}} record.{{.JsonName}}_urls.length {{"}}"}}个文件
          </a-button>
        </template>
{{- else if eq .FormType "file"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-button v-if="record.{{.JsonName}}_url" type="link" size="small" @click="handlePreviewFile(record.{{.JsonName}}_url)">查看文件</a-button>
        </template>
{{- else if eq .FormType "textarea"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-button type="link" size="small" @click="handleViewText(record, '{{.Comment}}', '{{.JsonName}}')">查看</a-button>
        </template>
{{- else if eq .FormType "editor"}}
        <template v-if="column.key === '{{.JsonName}}'">
          <a-button type="link" size="small" @click="handleViewContent(record, '{{.Comment}}', '{{.JsonName}}')">查看</a-button>
        </template>
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
        <template v-if="column.key === '{{.JsonName}}'">
          {{"{{"}} record.{{.JsonName}}?.{{.DisplayField}} || '-' {{"}}"}}
        </template>
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
        <template v-if="column.key === 'creator'">
{{- if .HasCreatedByProfile}}
          <a-space direction="vertical" :size="0">
            <span>{{"{{"}} record.creator?.nickname || record.creator?.username || '-' {{"}}"}}</span>
            <span v-if="record.creator_profile" style="color: #999; font-size: 12px">
              {{"{{"}} record.creator_profile.{{.CreatedByProfileField}} {{"}}"}}
            </span>
          </a-space>
{{- else}}
          {{"{{"}} record.creator?.nickname || record.creator?.username || '-' {{"}}"}}
{{- end}}
        </template>
{{- end}}
{{- if .HasAudit}}
        <template v-if="column.key === 'audit_status'">
          <a-space>
            <a-tag v-if="record.audit_status === 0" color="default">待审批</a-tag>
            <a-tag v-else-if="record.audit_status === 1" color="success">审批通过</a-tag>
            <a-tag v-else-if="record.audit_status === 2" color="error">审批拒绝</a-tag>
            <a-button v-if="record.audit_status !== 0" type="link" size="small" @click="handleViewAudit(record)">详情</a-button>
          </a-space>
        </template>
{{- end}}
{{- if .HasCreatedAt }}
        <template v-if="column.key === 'created_at'">{{"{{"}} formatTime(record.created_at) {{"}}"}}</template>
{{- end}}
        <template v-if="column.key === 'action'">
{{- if .HasMenu}}
          <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'{{.MenuConfig.Permission}}:edit'">编辑</a-button>
          <a-button type="link" size="small" @click="handleCopy(record)" v-permission="'{{.MenuConfig.Permission}}:edit'">复制</a-button>
{{- if .HasAudit}}
          <a-button v-if="record.audit_status === 0" type="link" size="small" @click="handleAudit(record)" v-permission="'{{.MenuConfig.Permission}}:audit'">审批</a-button>
{{- end}}
          <a-button type="link" size="small" danger @click="confirmDelete(record.id)" v-permission="'{{.MenuConfig.Permission}}:delete'">删除</a-button>
{{- else}}
          <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
          <a-button type="link" size="small" @click="handleCopy(record)">复制</a-button>
{{- if .HasAudit}}
          <a-button v-if="record.audit_status === 0" type="link" size="small" @click="handleAudit(record)">审批</a-button>
{{- end}}
          <a-button type="link" size="small" danger @click="confirmDelete(record.id)">删除</a-button>
{{- end}}
        </template>
      </template>
    </ProTable>
{{- if .HasTreeLayout}}
    </div>
{{- end}}

    <!-- 表单抽屉 -->
    <{{.ModelName}}Form
      v-model:open="drawerVisible"
      :record="currentRecord"
{{- if .LinkToUser}}
      :user-options="userOptions"
{{- end}}
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
      :{{.JsonName}}-options="{{.JsonName}}Options"
{{- end}}
{{- end}}
      @success="handleFormSuccess"
    />

{{- if .HasFiles}}
    <!-- 文件预览 -->
    <FilePreview
      v-model:open="previewVisible"
      :url="previewUrl"
      :name="previewName"
      :ext="previewExt"
    />

    <!-- 多文件预览列表 -->
    <a-modal v-model:open="filesModalVisible" title="文件列表" :footer="null" width="500px">
      <a-list :data-source="previewFiles" size="small">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-button type="link" @click="handlePreviewFile(item)">
              {{"{{"}} item.split('/').pop() {{"}}"}}
            </a-button>
          </a-list-item>
        </template>
      </a-list>
    </a-modal>
{{- end}}
    <!-- 文本域内容预览 -->
    <a-modal v-model:open="textModalVisible" :title="textModalTitle" :footer="null" width="600px">
      <div class="text-content">{{"{{"}} textModalContent {{"}}"}}</div>
    </a-modal>
{{- if .HasEditor}}
    <!-- 富文本内容预览 -->
    <a-modal v-model:open="contentModalVisible" :title="contentModalTitle" :footer="null" width="800px">
      <div class="rich-content" v-html="contentModalHtml"></div>
    </a-modal>
{{- end}}
{{- if .HasAudit}}
    <!-- 审批弹窗 -->
    <AuditModal
      v-model:open="auditModalVisible"
      title="审批"
      @confirm="handleAuditConfirm"
    />
    <!-- 审批详情弹窗 -->
    <a-modal v-model:open="auditDetailVisible" title="审批详情" :footer="null" width="500px">
      <a-descriptions :column="1" bordered size="small">
        <a-descriptions-item label="审批状态">
          <a-tag v-if="auditDetailRecord?.audit_status === 1" color="success">审批通过</a-tag>
          <a-tag v-else-if="auditDetailRecord?.audit_status === 2" color="error">审批拒绝</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="审批人">{{"{{"}} auditDetailRecord?.auditor?.nickname || auditDetailRecord?.auditor?.username || '-' {{"}}"}}</a-descriptions-item>
        <a-descriptions-item label="审批时间">{{"{{"}} auditDetailRecord?.audit_time ? formatTime(auditDetailRecord.audit_time) : '-' {{"}}"}}</a-descriptions-item>
        <a-descriptions-item label="审批备注">{{"{{"}} auditDetailRecord?.audit_remark || '-' {{"}}"}}</a-descriptions-item>
      </a-descriptions>
{{- if .LinkToUser}}
      <div v-if="auditDetailRecord?.audit_status === 2" style="margin-top: 16px; color: #999; font-size: 12px;">
        <InfoCircleOutlined /> 用户重新保存档案后将自动重置为待审批状态
      </div>
{{- end}}
    </a-modal>
{{- end}}
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted{{if or .HasTreeLayout (and .HasCreatedBy .DataIsolation) .HasDictSelect}}, computed{{end}}, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined, DeleteOutlined, ExclamationCircleOutlined{{if .HasTreeLayout}}, FolderOutlined, AppstoreOutlined, TagOutlined{{end}}{{if and .HasAudit .LinkToUser}}, InfoCircleOutlined{{end}} } from '@ant-design/icons-vue'
import ProTable from '@/components/ProTable.vue'
{{- if .HasFiles}}
import FilePreview from '@/components/FilePreview.vue'
{{- end}}
{{- if .HasAudit}}
import AuditModal from '@/components/AuditModal.vue'
{{- end}}
import {{.ModelName}}Form from './components/{{.ModelName}}Form.vue'
import { get{{.ModelName}}List, delete{{.ModelName}}, batchDelete{{.ModelName}}{{if .HasCreatedBy}}, get{{.ModelName}}CreatorOptions{{end}}{{if .HasAudit}}, audit{{.ModelName}}{{end}} } from '@/api/{{.ModuleName}}'
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
{{- if .UseOptionsApi}}
import { get{{.RelatedModel}}Options } from '@/api/{{.RelatedTable}}'
{{- else}}
import { get{{.RelatedModel}}List } from '@/api/{{.RelatedTable}}'
{{- end}}
{{- end}}
{{- end}}
{{- if .LinkToUser}}
import { getUserList } from '@/api/user'
{{- end}}
{{- if and .HasCreatedBy .DataIsolation}}
import { useUserStore } from '@/store/user'
{{- end}}
import { formatTime } from '@/utils/format'
import { {{.ModelName}} } from '@/types/{{.ModuleName}}'
{{- if .HasMenu}}
import { useTableColumns } from '@/utils/permission'
{{- end}}
{{- if .HasDictSelect}}
import { getDictDataByType } from '@/api/dict'
{{- end}}

const loading = ref(false)
const tableData = ref<{{.ModelName}}[]>([])
const drawerVisible = ref(false)
const currentRecord = ref<{{.ModelName}} | null>(null)
const selectedRowKeys = ref<number[]>([])
{{- if .HasAudit}}
// 审批相关
const auditModalVisible = ref(false)
const currentAuditId = ref<number>()
const auditDetailVisible = ref(false)
const auditDetailRecord = ref<{{.ModelName}} | null>(null)
{{- end}}
{{- if and .HasCreatedBy .DataIsolation}}

// 数据隔离：检查是否为管理员
const userStore = useUserStore()
const adminRoleIds = [{{.AdminRoleIds}}]
const isAdmin = computed(() => {
  const userRoleIds = userStore.user?.roles?.map((r: any) => r.id) || []
  return userRoleIds.some((id: number) => adminRoleIds.includes(id))
})
{{- end}}
{{- if .HasTreeLayout}}

// 左树右表相关
const treeLoading = ref(false)
const selectedCategoryId = ref<number | null>(null)
{{- range .Relations}}
{{- if .UseTreeLayout}}
const categoryOptions = computed(() => {{.JsonName}}Options.value)
const totalCount = computed(() => categoryOptions.value.reduce((sum: number, item: any) => sum + (item.count || 0), 0))
{{- end}}
{{- end}}

const handleSelectCategory = (id: number | null) => {
  selectedCategoryId.value = id
{{- range .Relations}}
{{- if .UseTreeLayout}}
  searchForm.{{.ForeignKeyJson}} = id || undefined
{{- end}}
{{- end}}
  pagination.current = 1
  fetchData()
}
{{- end}}
{{- if .HasFiles}}

// 文件预览
const previewVisible = ref(false)
const previewUrl = ref('')
const previewName = ref('')
const previewExt = ref('')

// 多文件预览
const filesModalVisible = ref(false)
const previewFiles = ref<string[]>([])

const handlePreviewFile = (url: string, name?: string) => {
  previewUrl.value = url
  previewName.value = name || url.split('/').pop() || 'file'
  previewExt.value = url.split('.').pop()?.toLowerCase() || ''
  previewVisible.value = true
}

const handlePreviewFiles = (urls: string[]) => {
  previewFiles.value = urls
  filesModalVisible.value = true
}
{{- end}}
// 文本域内容预览
const textModalVisible = ref(false)
const textModalTitle = ref('')
const textModalContent = ref('')

const handleViewText = (record: any, title: string, field: string) => {
  textModalTitle.value = title
  textModalContent.value = record[field] || '暂无内容'
  textModalVisible.value = true
}
{{- if .HasEditor}}

// 富文本内容预览
const contentModalVisible = ref(false)
const contentModalTitle = ref('')
const contentModalHtml = ref('')

const handleViewContent = (record: any, title: string, field: string) => {
  contentModalTitle.value = title
  contentModalHtml.value = record[field] || '<p style="color: #999">暂无内容</p>'
  contentModalVisible.value = true
}
{{- end}}

// 默认标签颜色（当字典未配置 tag_type 时作为回退）
const defaultTagColors = ['blue', 'green', 'orange', 'purple', 'cyan', 'magenta', 'gold', 'lime']

{{- range .ListColumns}}
{{- if eq .FormType "select"}}
{{- if .DictType}}
// 字典选项映射（动态获取）
const {{.JsonName}}DictList = ref<any[]>([])
const {{.JsonName}}Options = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  {{.JsonName}}DictList.value.forEach(item => { map[item.value] = item.label })
  return map
})
const {{.JsonName}}Colors = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  {{.JsonName}}DictList.value.forEach((item, i) => {
    // 优先使用字典配置的 tag_type，否则使用默认颜色
    map[item.value] = item.tag_type || defaultTagColors[i % defaultTagColors.length]
  })
  return map
})
{{- else if .SelectOptions}}
// 下拉选项映射
const {{.JsonName}}Options: Record<string, string> = {
{{- range .SelectOptions}}
  '{{.Value}}': '{{.Label}}',
{{- end}}
}
const {{.JsonName}}Colors: Record<string, string> = Object.fromEntries(
  Object.keys({{.JsonName}}Options).map((k, i) => [k, defaultTagColors[i % defaultTagColors.length]])
)
{{- else}}
const {{.JsonName}}Options: Record<string, string> = { '1': '启用', '0': '禁用' }
const {{.JsonName}}Colors: Record<string, string> = { '1': 'green', '0': 'red' }
{{- end}}
{{- end}}
{{- end}}

// 关联选项
{{- if .LinkToUser}}
const userOptions = ref<any[]>([])
{{- end}}
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
const {{.JsonName}}Options = ref<any[]>([])
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
const creatorOptions = ref<any[]>([])
{{- end}}
{{- if or .HasRelations .HasCreatedBy .LinkToUser}}

// 下拉搜索过滤
const filterOption = (input: string, option: any) => {
  return option.children?.[0]?.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0
}
{{- end}}

const searchForm = reactive({
{{- range .SearchColumns}}
  {{.JsonName}}: {{if eq .FormType "select"}}undefined as number | undefined{{else if or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64")}}undefined as number | undefined{{else if eq .SearchType "eq"}}undefined as {{.TsType}} | undefined{{else}}''{{end}},
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  {{.ForeignKeyJson}}: undefined as number | undefined,
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
  created_by: undefined as number | undefined,
{{- end}}
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

{{- if .HasMenu}}
// 基础列配置（不含操作列）
const baseColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', align: 'center', sorter: true },
{{- range .ListColumns}}
{{- if and (ne .JsonName "id") (ne .JsonName "created_at") (ne .JsonName "updated_at")}}
{{- if .IsSortable}}
  { title: '{{.Comment}}', dataIndex: '{{.JsonName}}', key: '{{.JsonName}}', align: 'center', sorter: true },
{{- else}}
  { title: '{{.Comment}}', dataIndex: '{{.JsonName}}', key: '{{.JsonName}}', align: 'center' },
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  { title: '{{.Comment}}', dataIndex: '{{.JsonName}}', key: '{{.JsonName}}', align: 'center' },
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
  { title: '创建人', dataIndex: 'creator', key: 'creator', align: 'center' },
{{- end}}
{{- if .HasAudit}}
  { title: '审批状态', dataIndex: 'audit_status', key: 'audit_status', align: 'center' },
{{- end}}
{{- if .HasCreatedAt }}
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', align: 'center', sorter: true },
{{- end}}
]

// 操作列配置
const actionColumn = { title: '操作', key: 'action', width: 200, align: 'center' }

// 根据权限动态显示操作列（有编辑或删除权限时显示）
const columns = useTableColumns(baseColumns, actionColumn, ['{{.MenuConfig.Permission}}:edit', '{{.MenuConfig.Permission}}:delete'])
{{- else}}
const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', align: 'center', sorter: true },
{{- range .ListColumns}}
{{- if and (ne .JsonName "id") (ne .JsonName "created_at") (ne .JsonName "updated_at")}}
{{- if .IsSortable}}
  { title: '{{.Comment}}', dataIndex: '{{.JsonName}}', key: '{{.JsonName}}', align: 'center', sorter: true },
{{- else}}
  { title: '{{.Comment}}', dataIndex: '{{.JsonName}}', key: '{{.JsonName}}', align: 'center' },
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  { title: '{{.Comment}}', dataIndex: '{{.JsonName}}', key: '{{.JsonName}}', align: 'center' },
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
  { title: '创建人', dataIndex: 'creator', key: 'creator', align: 'center' },
{{- end}}
{{- if .HasAudit}}
  { title: '审批状态', dataIndex: 'audit_status', key: 'audit_status', align: 'center' },
{{- end}}
{{- if .HasCreatedAt }}
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', align: 'center', sorter: true },
{{- end}}
  { title: '操作', key: 'action', width: 200, align: 'center' }
]
{{- end}}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await get{{.ModelName}}List({
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
{{- range .SearchColumns}}
  searchForm.{{.JsonName}} = {{if eq .FormType "select"}}undefined{{else if or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64")}}undefined{{else if eq .SearchType "eq"}}undefined{{else}}''{{end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
  searchForm.{{.ForeignKeyJson}} = undefined
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
  searchForm.created_by = undefined
{{- end}}
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

const handleEdit = (record: {{.ModelName}}) => {
  currentRecord.value = record
  drawerVisible.value = true
}

const handleCopy = (record: {{.ModelName}}) => {
  // 复制时不传 id，使表单识别为新增模式
  const { id, created_at, updated_at, ...copyData } = record
  currentRecord.value = copyData as {{.ModelName}}
  drawerVisible.value = true
}

const handleFormSuccess = () => {
  fetchData()
{{- range .Relations}}
{{- if and (or (eq .RelationType "belongsTo") (eq .RelationType "many2many")) .UseOptionsApi}}
  fetch{{.RelatedModel}}Options() // 刷新关联选项count
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
  fetchCreatorOptions() // 刷新创建人选项
{{- end}}
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
      await delete{{.ModelName}}(id)
      message.success('删除成功')
      fetchData()
{{- range .Relations}}
{{- if and (or (eq .RelationType "belongsTo") (eq .RelationType "many2many")) .UseOptionsApi}}
      fetch{{.RelatedModel}}Options()
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
      fetchCreatorOptions()
{{- end}}
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
      await batchDelete{{.ModelName}}(selectedRowKeys.value)
      message.success('批量删除成功')
      selectedRowKeys.value = []
      fetchData()
{{- range .Relations}}
{{- if and (or (eq .RelationType "belongsTo") (eq .RelationType "many2many")) .UseOptionsApi}}
      fetch{{.RelatedModel}}Options()
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
      fetchCreatorOptions()
{{- end}}
    }
  })
}
{{- if .HasAudit}}

// 审批
const handleAudit = (record: {{.ModelName}}) => {
  currentAuditId.value = record.id
  auditModalVisible.value = true
}

const handleAuditConfirm = async (data: { audit_status: number; audit_remark: string }) => {
  try {
    await audit{{.ModelName}}(currentAuditId.value!, data)
    message.success('审批成功')
    auditModalVisible.value = false
    fetchData()
  } catch (error: any) {
    message.error(error.response?.data?.msg || '审批失败')
  }
}

// 查看审批详情
const handleViewAudit = (record: {{.ModelName}}) => {
  auditDetailRecord.value = record
  auditDetailVisible.value = true
}
{{- end}}

{{- if or .HasRelations .HasCreatedBy .LinkToUser}}
// 获取关联选项
{{- if .LinkToUser}}
const fetchUserOptions = async () => {
  try {
    const res = await getUserList({ page: 1, page_size: 1000 })
    userOptions.value = res.data.list || []
  } catch (e) {
    console.error('获取用户选项失败', e)
  }
}
{{- end}}
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
const fetch{{.RelatedModel}}Options = async () => {
  try {
{{- if .UseOptionsApi}}
    // 使用轻量ptions接口（返回id,name,count）
{{- if $.DataIsolation}}
    // 数据隔离：非管理员统计时按当前用户过滤
    const res = await get{{.RelatedModel}}Options({
      display_field: '{{.DisplayField}}',
      count_table: '{{$.TableName}}',
      count_field: '{{.ForeignKeyJson}}'{{if $.HasDeletedAt}},
      exclude_deleted: true{{end}},
      count_created_by: isAdmin.value ? undefined : userStore.user?.id
    })
{{- else}}
    const res = await get{{.RelatedModel}}Options({
      display_field: '{{.DisplayField}}',
      count_table: '{{$.TableName}}',
      count_field: '{{.ForeignKeyJson}}'{{if $.HasDeletedAt}},
      exclude_deleted: true{{end}}
    })
{{- end}}
    {{.JsonName}}Options.value = res.data || []
{{- else}}
    const res = await get{{.RelatedModel}}List({ page: 1, page_size: 1000 })
    {{.JsonName}}Options.value = res.data.list || []
{{- end}}
  } catch (e) {
    console.error('获取{{.Comment}}选项失败', e)
  }
}
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
const fetchCreatorOptions = async () => {
  try {
    const res = await get{{.ModelName}}CreatorOptions()
    creatorOptions.value = res.data || []
  } catch (e) {
    console.error('获取创建人选项失败', e)
  }
}
{{- end}}
{{- end}}
{{- if .HasDictSelect}}

// 获取字典数据
{{- range .ListColumns}}
{{- if and (eq .FormType "select") .DictType}}
const fetch{{.FieldName}}Dict = async () => {
  try {
    const res = await getDictDataByType('{{.DictType}}')
    {{.JsonName}}DictList.value = res.data || []
  } catch (e) {
    console.error('获取{{.Comment}}字典失败', e)
  }
}
{{- end}}
{{- end}}
{{- end}}

onMounted(() => {
{{- if .LinkToUser}}
  fetchUserOptions()
{{- end}}
{{- range .Relations}}
{{- if or (eq .RelationType "belongsTo") (eq .RelationType "many2many")}}
  fetch{{.RelatedModel}}Options()
{{- end}}
{{- end}}
{{- if and .HasCreatedBy .DataIsolation}}
  // 仅管理员加载创建人选项
  if (isAdmin.value) {
    fetchCreatorOptions()
  }
{{- else if .HasCreatedBy}}
  fetchCreatorOptions()
{{- end}}
{{- if .HasDictSelect}}
  // 加载字典数据
{{- range .ListColumns}}
{{- if and (eq .FormType "select") .DictType}}
  fetch{{.FieldName}}Dict()
{{- end}}
{{- end}}
{{- end}}
  fetchData()
})
</script>

<style scoped>
.{{.ModuleName}}-page {
  padding: 0;
}
{{- if .HasTreeLayout}}
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
{{- end}}
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
{{- if .HasEditor}}
.rich-content {
  max-height: 500px;
  overflow-y: auto;
  padding: 16px;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
}
.rich-content img {
  max-width: 100%;
}
.rich-content video {
  max-width: 100%;
}
{{- end}}
</style>
