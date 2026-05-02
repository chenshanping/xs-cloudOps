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

const DATA_SCOPE_INHERIT = 0
const DATA_SCOPE_ALL = 1
const DATA_SCOPE_CUSTOM = 2
const DATA_SCOPE_DEPT = 3
const DATA_SCOPE_DEPT_AND_CHILDREN = 4
const DATA_SCOPE_SELF = 5
const OWNER_FIELD_DEPT_ID = 'dept_id'
const OWNER_FIELD_CREATED_BY = 'created_by'

export const FEATURE_SCOPE_OPTIONS = [
  { value: DATA_SCOPE_INHERIT, label: '继承默认数据范围' },
  { value: DATA_SCOPE_ALL, label: '全部数据' },
  { value: DATA_SCOPE_CUSTOM, label: '自定义部门' },
  { value: DATA_SCOPE_DEPT, label: '本部门' },
  { value: DATA_SCOPE_DEPT_AND_CHILDREN, label: '本部门及下级' },
  { value: DATA_SCOPE_SELF, label: '仅本人' }
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
    const nextDataScope = normalizeSupportedDataScope(resource, current?.data_scope ?? DATA_SCOPE_INHERIT)
    return {
      resource_code: resource.code,
      data_scope: nextDataScope,
      dept_ids: nextDataScope === DATA_SCOPE_CUSTOM ? current?.depts?.map(item => item.id) || [] : []
    }
  })
}

export function resourceSupportsDeptScope(resource?: DataScopeResourceDefinition) {
  return Boolean(resource?.owner_fields?.includes(OWNER_FIELD_DEPT_ID))
}

export function resourceSupportsSelfScope(resource?: DataScopeResourceDefinition) {
  return Boolean(resource?.owner_fields?.includes(OWNER_FIELD_CREATED_BY))
}

export function getSupportedFeatureScopeOptions(resource?: DataScopeResourceDefinition) {
  return FEATURE_SCOPE_OPTIONS.filter(option => {
    switch (option.value) {
      case DATA_SCOPE_CUSTOM:
      case DATA_SCOPE_DEPT:
      case DATA_SCOPE_DEPT_AND_CHILDREN:
        return resourceSupportsDeptScope(resource)
      case DATA_SCOPE_SELF:
        return resourceSupportsSelfScope(resource)
      default:
        return true
    }
  })
}

export function normalizeSupportedDataScope(
  resource: DataScopeResourceDefinition | undefined,
  dataScope: number
) {
  const supportedValues = new Set(getSupportedFeatureScopeOptions(resource).map(option => option.value))
  return supportedValues.has(dataScope) ? dataScope : DATA_SCOPE_INHERIT
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
    .filter(item => item.data_scope > DATA_SCOPE_INHERIT)
    .map(item => ({
      resource_code: item.resource_code,
      data_scope: item.data_scope,
      dept_ids: item.data_scope === DATA_SCOPE_CUSTOM ? [...item.dept_ids] : []
    }))

  const unknownScopePayloads = unknownItems.map(item => ({
    resource_code: item.resource_code,
    data_scope: item.data_scope,
    dept_ids: [...item.dept_ids]
  }))

  return [...knownScopePayloads, ...unknownScopePayloads]
}
