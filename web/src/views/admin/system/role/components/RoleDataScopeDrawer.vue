<template>
  <a-drawer
    v-model:open="visible"
    :width="drawerWidth"
    placement="right"
    :mask-closable="false"
    destroyOnClose
    :class="['permission-drawer', { 'permission-drawer--dark': uiStore.isDark }]"
  >
    <template #title>
      <RoleDrawerContextHeader
        title="数据权限"
        :role-name="roleName"
      />
    </template>

    <a-spin :spinning="loading">
      <div v-if="!loading && featureDataScopes.length === 0" class="empty-tip">
        <a-empty description="暂无可配置的数据权限资源" />
      </div>

      <RolePermissionDataScopePanel
        v-else
        v-model:model-value="featureDataScopes"
        :resource-definitions="resourceDefinitions"
        :dept-options="deptOptions"
        :default-data-scope-label="formatDataScopeLabel(defaultDataScope)"
      />
    </a-spin>

    <template #footer>
      <div class="drawer-footer">
        <a-button @click="visible = false">取消</a-button>
        <a-button type="primary" :loading="saveLoading" :disabled="loading" @click="handleSave">
          保存
        </a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import RoleDrawerContextHeader from './RoleDrawerContextHeader.vue'
import RolePermissionDataScopePanel from './RolePermissionDataScopePanel.vue'
import { useResponsiveDrawerWidth } from './useResponsiveDrawerWidth'
import {
  assignDataScopes,
  getDataScopeResources,
  getRole,
  type DataScopeResource,
  type RoleFeatureDataScopePayload
} from '@/api/role'
import { useUiStore } from '@/store/ui'
import {
  buildRoleFeatureDataScopePayload,
  buildRoleFeatureDataScopeForm,
  findDataScopeResourceDefinition,
  formatDataScopeLabel,
  splitKnownAndUnknownFeatureDataScopes,
  type RoleFeatureDataScopeFormItem,
  type RolePermissionDeptOption
} from './dataScopeResources'

interface Props {
  roleId: number
  roleName: string
  deptOptions: RolePermissionDeptOption[]
}

const props = defineProps<Props>()
const visible = defineModel<boolean>('open', { default: false })
const uiStore = useUiStore()
const { drawerWidth } = useResponsiveDrawerWidth(960, 820)

const loading = ref(false)
const saveLoading = ref(false)
const requestToken = ref(0)
const resourceDefinitions = ref<DataScopeResource[]>([])
const featureDataScopes = ref<RoleFeatureDataScopeFormItem[]>([])
const unknownFeatureDataScopes = ref<RoleFeatureDataScopePayload[]>([])
const defaultDataScope = ref(1)

const resetState = () => {
  resourceDefinitions.value = []
  featureDataScopes.value = []
  unknownFeatureDataScopes.value = []
  defaultDataScope.value = 1
}

const fetchResources = async (token: number) => {
  const res = await getDataScopeResources()
  if (token !== requestToken.value) return
  resourceDefinitions.value = res.data || []
}

const loadRoleData = async (token: number) => {
  if (!props.roleId) return
  const res = await getRole(props.roleId)
  if (token !== requestToken.value) return
  const { knownScopes, unknownScopes } = splitKnownAndUnknownFeatureDataScopes(
    resourceDefinitions.value,
    res.data.feature_data_scopes || []
  )
  defaultDataScope.value = res.data.data_scope || 1
  featureDataScopes.value = buildRoleFeatureDataScopeForm(resourceDefinitions.value, knownScopes)
  unknownFeatureDataScopes.value = unknownScopes
}

const validate = () => {
  const invalid = featureDataScopes.value.find(item => item.data_scope === 2 && item.dept_ids.length === 0)
  if (invalid) {
    const label = findDataScopeResourceDefinition(resourceDefinitions.value, invalid.resource_code)?.label
      || invalid.resource_code
    message.warning(`请为「${label}」选择自定义部门`)
    return false
  }
  return true
}

const handleSave = async () => {
  if (!validate()) return
  saveLoading.value = true
  try {
    await assignDataScopes(
      props.roleId,
      buildRoleFeatureDataScopePayload(featureDataScopes.value, unknownFeatureDataScopes.value),
      { silent: true } as any
    )
    message.success('数据权限保存成功')
    visible.value = false
  } catch (error) {
    const msg = error instanceof Error && error.message ? error.message : '请重试'
    message.warning(`数据权限保存失败：${msg}`)
  } finally {
    saveLoading.value = false
  }
}

watch([visible, () => props.roleId], async ([isVisible]) => {
  if (!isVisible) {
    requestToken.value += 1
    loading.value = false
    resetState()
    return
  }
  const token = ++requestToken.value
  loading.value = true
  resetState()
  try {
    await fetchResources(token)
    if (token !== requestToken.value) return
    await loadRoleData(token)
  } finally {
    if (token === requestToken.value) {
      loading.value = false
    }
  }
})
</script>

<style scoped>
@import './rolePermissionDrawerShared.css';

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.empty-tip {
  padding: 48px 0;
}

</style>
