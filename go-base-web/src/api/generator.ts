import request from '@/utils/request'
import type { ApiResponse } from '@/types'

// 表信息
export interface TableInfo {
  table_name: string
  table_comment: string
}

// 下拉选项
export interface SelectOption {
  label: string
  value: string
}

// 开关值
export interface SwitchValue {
  active_value: string
  inactive_value: string
  active_text: string
  inactive_text: string
}

// 字段配置
export interface ColumnConfig {
  column_name: string
  field_name: string
  field_type: string
  json_name: string
  ts_type: string
  comment: string
  db_type: string
  db_length: number
  default_value: string
  is_primary_key: boolean
  is_required: boolean
  is_searchable: boolean
  search_type: string  // eq/like/gt/gte/lt/lte
  is_list_visible: boolean
  is_form_visible: boolean
  is_sortable: boolean // 是否可排序
  sort_order: string   // asc/desc
  is_unique: boolean   // 是否唯一
  form_type: string
  dict_type: string
  select_options: SelectOption[]
  switch_values: SwitchValue | null
  // belongsTo 关联配置（仅外键字段使用，如 category_id）
  related_table: string    // 关联表名（如 category）
  related_module: string   // 关联模块名（用于 API import，留空则使用表名）
  display_field: string    // 显示字段（如 name）
  use_options_api: boolean // 使用轻量options接口
  use_tree_layout: boolean // 左树右表布局
}

// 关联配置
export interface RelationConfig {
  relation_type: string
  related_table: string
  related_module: string // 关联模块名（用于 API import，留空则使用表名）
  related_model: string
  foreign_key: string
  reference_key: string
  join_table: string
  display_field: string
  comment: string
  is_required: boolean   // 是否必填
  use_options_api: boolean // 使用轻量options接口（返回id,name,count）
  use_tree_layout: boolean // 使用左树右表布局
}

// 菜单配置
export interface MenuConfig {
  parent_id: number
  menu_name: string
  menu_icon: string
  menu_sort: number
  permission: string
}

// 单个统计图表配置
export interface StatsChartConfig {
  field: string      // 分组字段（如 category_id, status）
  chart_type: string // 图表类型: pie/bar
  title: string      // 图表标题（可选，默认使用字段注释）
}

// 统计配置
export interface StatsConfig {
  enabled: boolean             // 是否启用统计
  charts: StatsChartConfig[]   // 多个分组统计图表
  time_field: string           // 时间字段（用于趋势统计，如 created_at）
}

// 生成器配置
export interface GeneratorConfig {
  id?: number  // 配置ID（编辑时传递）
  table_name: string
  module_name: string
  description: string
  author: string
  generate_backend: boolean
  generate_frontend: boolean
  generate_sql: boolean
  frontend_path: string
  // 时间字段开关（用于生成模型）
  has_created_at?: boolean
  has_updated_at?: boolean
  has_deleted_at?: boolean
  has_created_by?: boolean // 创建人字段
  // 创建人身份关联配置（显示创建者的身份信息）
  created_by_profile_table?: string // 创建者身份关联表名
  created_by_profile_field?: string // 身份表中要显示的字段
  // 数据隔离配置（仅当 has_created_by 为 true 时生效）
  data_isolation?: boolean  // 是否启用数据隔离
  admin_role_ids?: string   // 管理员角色ID列表（逗号分隔）
  // 审批功能配置
  has_audit?: boolean       // 是否启用审批功能
  // 前台接口配置
  generate_frontend_api?: boolean // 是否生成前台用户使用的接口
  // 用户关联配置
  link_to_user?: boolean      // 是否关联用户表（一对一，生成 /my 接口）
  profile_name?: string       // 身份显示名称（如"医生"、"商家"）
  profile_icon?: string       // 身份图标（ant-design图标名）
  profile_role_code?: string  // 身份限定角色编码（为空表示不限制）
  columns: ColumnConfig[]
  relations: RelationConfig[]
  menu_config: MenuConfig | null
  stats_config?: StatsConfig | null  // 统计配置
}

// 生成的文件
export interface GeneratedFile {
  path: string
  content: string
  type: string
}

// 预览结果
export interface PreviewResult {
  files: GeneratedFile[]
}

// 保存的配置信息
export interface SavedConfig {
  id: number
  table_name: string
  module_name: string
  description: string
  config_json: string
  created_at: string
  updated_at: string
}

// 获取数据库表列表
export function getTables() {
  return request.get<any, ApiResponse<TableInfo[]>>('/generator/tables')
}

// 获取表字段信息
export function getTableColumns(tableName: string) {
  return request.get<any, ApiResponse<ColumnConfig[]>>(`/generator/tables/${tableName}/columns`)
}

// 预览生成的代码
export function previewCode(config: GeneratorConfig) {
  return request.post<any, ApiResponse<PreviewResult>>('/generator/preview', config)
}

// 生成代码
export function generateCode(config: GeneratorConfig) {
  return request.post<any, ApiResponse<PreviewResult>>('/generator/generate', config)
}

// 获取已生成的模块列表
export function getGeneratedModules() {
  return request.get<any, ApiResponse<string[]>>('/generator/modules')
}

// 删除已生成的模块
export function deleteModule(moduleName: string) {
  return request.delete<any, ApiResponse>(`/generator/modules/${moduleName}`)
}

// 保存配置（有id则更新，无id则新增）
export function saveConfig(config: GeneratorConfig) {
  if (config.id) {
    // 更新
    return request.put<any, ApiResponse<SavedConfig>>(`/generator/configs/${config.id}`, config)
  }
  // 新增
  return request.post<any, ApiResponse<SavedConfig>>('/generator/configs', config)
}

// 获取已保存的配置列表
export function getSavedConfigs() {
  return request.get<any, ApiResponse<SavedConfig[]>>('/generator/configs')
}

// 获取单个配置详情
export function getSavedConfig(id: number) {
  return request.get<any, ApiResponse<SavedConfig>>(`/generator/configs/${id}`)
}

// 删除已保存的配置
export function deleteSavedConfig(id: number) {
  return request.delete<any, ApiResponse>(`/generator/configs/${id}`)
}

// 执行建表SQL
export function executeSQL(sql: string) {
  return request.post<any, ApiResponse>('/generator/execute-sql', { sql })
}
