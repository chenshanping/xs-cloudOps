import type { RoleFeatureDataScope } from '@/types'

export interface RolePermissionDeptOption {
  key: string | number
  title: string
  value: number
  disabled?: boolean
  children?: RolePermissionDeptOption[]
}

export interface DataScopeResourceDefinition {
  code: string
  label: string
  description: string
}

export interface RoleFeatureDataScopeFormItem {
  resource_code: string
  data_scope: number
  dept_ids: number[]
}

export const FEATURE_SCOPE_OPTIONS = [
  { value: 0, label: '继承默认数据范围' },
  { value: 1, label: '全部数据' },
  { value: 2, label: '自定义部门' },
  { value: 3, label: '本部门' },
  { value: 4, label: '本部门及下级' },
  { value: 5, label: '仅本人' }
]

export const ROLE_FEATURE_SCOPE_RESOURCES: DataScopeResourceDefinition[] = [
  {
    code: 'system:user-management',
    label: '用户管理',
    description: '控制用户列表及关联用户操作可见的数据范围。'
  },
  {
    code: 'system:dept-management',
    label: '部门管理',
    description: '控制部门树、可管理部门及部门统计可见的数据范围。'
  }
]

export function formatDataScopeLabel(value: number): string {
  return FEATURE_SCOPE_OPTIONS.find(option => option.value === value)?.label || '-'
}

export function buildRoleFeatureDataScopeForm(
  scopes?: RoleFeatureDataScope[]
): RoleFeatureDataScopeFormItem[] {
  const scopeMap = new Map((scopes || []).map(scope => [scope.resource_code, scope]))
  return ROLE_FEATURE_SCOPE_RESOURCES.map(resource => {
    const current = scopeMap.get(resource.code)
    return {
      resource_code: resource.code,
      data_scope: current?.data_scope ?? 0,
      dept_ids: current?.depts?.map(item => item.id) || []
    }
  })
}
