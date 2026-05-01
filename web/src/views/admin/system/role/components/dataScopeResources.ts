import type { DataScopeResource } from '@/api/role'
import type { RoleFeatureDataScopePayload } from '@/api/role'
import type { RoleFeatureDataScope } from '@/types'

export interface RolePermissionDeptOption {
  key: string | number
  title: string
  value: number
  disabled?: boolean
  children?: RolePermissionDeptOption[]
}

export type DataScopeResourceDefinition = DataScopeResource

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

export function formatDataScopeLabel(value: number): string {
  return FEATURE_SCOPE_OPTIONS.find(option => option.value === value)?.label || '-'
}

export function findDataScopeResourceDefinition(
  resources: DataScopeResourceDefinition[],
  resourceCode: string
) {
  return resources.find(resource => resource.code === resourceCode)
}

export function buildRoleFeatureDataScopeForm(
  resources: DataScopeResourceDefinition[] = [],
  scopes?: RoleFeatureDataScope[]
): RoleFeatureDataScopeFormItem[] {
  const scopeMap = new Map((scopes || []).map(scope => [scope.resource_code, scope]))
  return resources.map(resource => {
    const current = scopeMap.get(resource.code)
    return {
      resource_code: resource.code,
      data_scope: current?.data_scope ?? 0,
      dept_ids: current?.depts?.map(item => item.id) || []
    }
  })
}

export function splitKnownAndUnknownFeatureDataScopes(
  resources: DataScopeResourceDefinition[] = [],
  scopes: RoleFeatureDataScope[] = []
) {
  const resourceCodes = new Set(resources.map(resource => resource.code))

  return scopes.reduce<{
    knownScopes: RoleFeatureDataScope[]
    unknownScopes: RoleFeatureDataScopePayload[]
  }>((result, scope) => {
    if (resourceCodes.has(scope.resource_code)) {
      result.knownScopes.push(scope)
      return result
    }

    result.unknownScopes.push({
      resource_code: scope.resource_code,
      data_scope: scope.data_scope,
      dept_ids: scope.depts?.map(item => item.id) || []
    })
    return result
  }, {
    knownScopes: [],
    unknownScopes: []
  })
}

export function buildRoleFeatureDataScopePayload(
  knownItems: RoleFeatureDataScopeFormItem[] = [],
  unknownItems: RoleFeatureDataScopePayload[] = []
): RoleFeatureDataScopePayload[] {
  const knownScopePayloads = knownItems
    .filter(item => item.data_scope > 0)
    .map(item => ({
      resource_code: item.resource_code,
      data_scope: item.data_scope,
      dept_ids: item.data_scope === 2 ? [...item.dept_ids] : []
    }))

  const unknownScopePayloads = unknownItems.map(item => ({
    resource_code: item.resource_code,
    data_scope: item.data_scope,
    dept_ids: [...item.dept_ids]
  }))

  return [...knownScopePayloads, ...unknownScopePayloads]
}
