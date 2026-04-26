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
  email: string
  phone: string
  avatar: string
  avatar_file_id?: number
  avatar_file_url?: string
  status: number
  roles: Role[]
  created_at: string
  updated_at: string
}

// 角色
export interface Role {
  id: number
  name: string
  code: string
  sort: number
  status: number
  remark: string
  menus?: Menu[]
  apis?: Api[]
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

// 慢查询日志
export interface SlowLog {
  id: number
  sql: string
  rows: number
  latency: number
  source: string
  created_at: string
}
