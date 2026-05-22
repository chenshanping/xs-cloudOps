// 响应结构
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应
export interface PageResponse<T = any> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 用户
export interface User {
  id: number
  username: string
  nickname: string
  gender: number
  email: string
  phone: string
  avatar_file_id?: number
  avatar_file_url?: string
  status: number
  dept_id: number
  dept?: Dept
  roles: Role[]
  created_at: string
  updated_at: string
}

// 部门
export interface Dept {
  id: number
  parent_id: number
  ancestors: string
  name: string
  sort: number
  status: number
  remark: string
  direct_user_count?: number
  total_user_count?: number
  has_children?: boolean
  bindable?: boolean
  manageable?: boolean
  children?: Dept[]
  created_at: string
  updated_at: string
}

export interface ManageableDeptTreeData {
  tree: Dept[]
  unassigned_user_count: number
  default_avatar_url?: string
}

// 角色
export interface RoleUserSummary {
  id?: number
  username: string
  nickname: string
  status?: number
}

export interface RoleFeatureDataScope {
  id?: number
  role_id?: number
  resource_code: string
  data_scope: number
  depts?: Dept[]
}

export interface Role {
  id: number
  name: string
  code: string
  sort: number
  status: number
  is_super_admin: boolean
  data_scope: number
  remark: string
  user_count?: number
  menus?: Menu[]
  apis?: Api[]
  depts?: Dept[]
  feature_data_scopes?: RoleFeatureDataScope[]
  users?: RoleUserSummary[]
  created_at: string
  updated_at: string
}

// 菜单
export interface Menu {
  id: number
  parent_id: number
  name: string
  path: string
  component: string
  icon: string
  sort: number
  type: number
  permission: string
  status: number
  hidden: number
  is_standalone: number
  apis?: Api[]
  children?: Menu[]
  created_at: string
  updated_at: string
}

// API 字段信息
export interface ApiFieldInfo {
  name: string
  type: string
  description: string
  required: boolean
  in: string // query, body, path
}

// API
export interface Api {
  id: number
  path: string
  method: string
  group: string
  description: string
  request_params?: string // JSON 字符串
  response_params?: string // JSON 字符串
  need_auth: boolean
  created_at: string
  updated_at: string
}

// 操作日志
export interface OperationLog {
  id: number
  user_id: number
  username: string
  ip: string
  method: string
  path: string
  request: string
  response: string
  status: number
  latency: number
  user_agent: string
  created_at: string
}

// 登录日志
export interface LoginLog {
  id: number
  user_id: number
  username: string
  ip: string
  location: string
  browser: string
  os: string
  status: number
  msg: string
  created_at: string
}

// 登录表单
export interface LoginForm {
  username: string
  password: string
}
