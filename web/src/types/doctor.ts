// 医生

// 医生响应体
export interface Doctor {
  id: number
  real_name?: string
  id_card?: string
  hospital?: string
  certificate_img_file_id?: number
  certificate_img_url?: string
  certificate_no?: string
  ntroduction?: string
  created_at?: string
  updated_at?: string
  user_id?: number
  user?: { id: number; username: string; nickname?: string; avatar_file_url?: string }
  created_by?: number
  creator?: { id: number; username: string; nickname?: string }
  audit_status?: number  // 0-待审批 1-审批通过 2-审批拒绝
  audit_remark?: string
  audit_time?: string
  audit_by?: number
  auditor?: { id: number; username: string; nickname?: string }
}

// 创建医生请求体
export interface CreateDoctorRequest {
  user_id: number
  real_name?: string
  id_card?: string
  hospital?: string
  certificate_img_file_id?: number
  certificate_no?: string
  ntroduction?: string
}

// 更新医生请求体
export type UpdateDoctorRequest = Partial<CreateDoctorRequest>

// 医生查询参数
export interface DoctorQuery {
  page?: number
  page_size?: number
  created_by?: number
  sort_field?: string
  sort_order?: 'ascend' | 'descend'
}

// 选项项（用于下拉选择）
export interface OptionItem {
  id: number
  name: string
  count?: number
}

// 保存我的医生请求体
export interface SaveMyDoctorRequest {
  real_name?: string
  id_card?: string
  hospital?: string
  certificate_img_file_id?: number
  certificate_no?: string
  ntroduction?: string
}
