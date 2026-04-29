<template>
  <div class="user-page">
    <div class="page-layout">
      <div class="left-tree">
        <div class="tree-header">
          <span>部门用户</span>
          <a-space :size="4">
            <a-button type="link" size="small" @click="expandAllTree">全部展开</a-button>
            <a-button type="link" size="small" @click="collapseAllTree">全部收缩</a-button>
          </a-space>
        </div>
        <a-spin :spinning="deptLoading">
          <a-tree
            v-if="deptTreeNodes.length"
            block-node
            :tree-data="deptTreeNodes"
            :selected-keys="selectedTreeKeys"
            :expanded-keys="expandedTreeKeys"
            @select="handleDeptSelect"
            @expand="handleTreeExpand"
          />
          <a-empty v-else description="暂无部门" :image="simpleImage" />
        </a-spin>
      </div>
      <div class="right-content">
        <ProTable
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          :pagination="pagination"
          row-key="id"
          :row-selection="{ selectedRowKeys, onChange: onSelectChange }"
          @search="handleSearch"
          @reset="handleReset"
          @change="handleTableChange"
        >
          <template #search>
            <a-form-item label="用户名">
              <a-input v-model:value="searchForm.username" placeholder="请输入用户名" allowClear />
            </a-form-item>
            <a-form-item label="状态">
              <a-select v-model:value="searchForm.status" placeholder="请选择" allowClear style="width: 120px">
                <a-select-option :value="1">启用</a-select-option>
                <a-select-option :value="0">禁用</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="性别">
              <a-select v-model:value="searchForm.gender" placeholder="请选择性别" allowClear style="width: 120px">
                <a-select-option v-for="option in genderOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="角色">
              <a-select v-model:value="searchForm.roleId" placeholder="请选择角色" allowClear style="width: 150px">
                <a-select-option v-for="role in roleList" :key="role.id" :value="role.id">
                  {{ role.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </template>

          <template #toolbar>
            <a-space>
              <a-button type="primary" @click="handleAdd" v-permission="'system:user:add'">
                <PlusOutlined /> 新增
              </a-button>
              <a-button
                type="primary"
                :disabled="selectedRowKeys.length === 0 || hasRestrictedManagedSelection"
                @click="handleBatchStatusChange(1)"
                v-permission="'system:user:batchEnable'"
              >
                批量启用
                <span v-if="selectedRowKeys.length > 0">({{ selectedRowKeys.length }})</span>
              </a-button>
              <a-button
                danger
                :disabled="selectedRowKeys.length === 0 || hasRestrictedManagedSelection"
                @click="handleBatchStatusChange(0)"
                v-permission="'system:user:batchDisable'"
              >
                批量禁用
                <span v-if="selectedRowKeys.length > 0">({{ selectedRowKeys.length }})</span>
              </a-button>
              <a-button
                :disabled="selectedRowKeys.length === 0 || hasRestrictedManagedSelection"
                @click="handleBatchResetPwd"
                v-permission="'system:user:batchResetPwd'"
              >
                批量重置密码
                <span v-if="selectedRowKeys.length > 0">({{ selectedRowKeys.length }})</span>
              </a-button>
              <a-button
                type="primary"
                danger
                :disabled="selectedRowKeys.length === 0 || hasRestrictedManagedSelection"
                @click="handleBatchDelete"
                v-permission="'system:user:delete'"
              >
                <DeleteOutlined /> 批量删除
                <span v-if="selectedRowKeys.length > 0">({{ selectedRowKeys.length }})</span>
              </a-button>
            </a-space>
          </template>

          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'avatar'">
              <a-avatar :size="40" :src="record.avatar_file_url">
                <template #icon><UserOutlined /></template>
              </a-avatar>
            </template>
            <template v-if="column.key === 'status'">
              <a-switch
                :checked="record.status === 1"
                :disabled="!canMutateManagedRecord(record)"
                @change="(checked: boolean) => handleStatusChange(record, checked)"
              />
            </template>
            <template v-if="column.key === 'gender'">
              <a-tag :color="getGenderOption(record.gender)?.tag_type || 'default'">
                {{ getGenderOption(record.gender)?.label || '-' }}
              </a-tag>
            </template>
            <template v-if="column.key === 'roles'">
              <a-tag v-for="role in record.roles" :key="role.id" color="blue">
                {{ role.name }}
              </a-tag>
            </template>
            <template v-if="column.key === 'dept'">
              <span v-if="record.dept?.name">{{ record.dept.name }}</span>
              <a-tag v-else color="error">未绑定部门</a-tag>
            </template>
            <template v-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
            <template v-if="column.key === 'action'">
              <a-space :size="0">
                <a-button
                  v-if="canMutateManagedRecord(record)"
                  type="link"
                  size="small"
                  @click="handleEdit(record)"
                  v-permission="'system:user:edit'"
                >编辑</a-button>
                <a-button v-if="showProfileButton" type="link" size="small" @click="handleViewProfiles(record)">身份</a-button>
                <a-dropdown v-if="canMutateManagedRecord(record)">
                  <a-button type="link" size="small">更多 <DownOutlined /></a-button>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item key="resetPwd" v-permission="'system:user:resetPwd'" @click="handleResetPwd(record)">重置密码</a-menu-item>
                      <a-menu-item
                        v-if="canDeleteRecord(record)"
                        key="delete"
                        v-permission="'system:user:delete'"
                        @click="confirmDelete(record)"
                      >
                        <span style="color: #ff4d4f">删除</span>
                      </a-menu-item>
                      <a-menu-item key="offline" v-permission="'system:user:forceOffline'" @click="confirmForceOffline(record)">
                        <span style="color: #ff4d4f">强制下线</span>
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </a-space>
            </template>
          </template>
        </ProTable>
      </div>
    </div>

    <UserFormDrawer
      v-model:open="drawerVisible"
      :title="drawerTitle"
      :is-edit="isEdit"
      :role-options="roleList"
      :gender-options="genderOptions"
      :dept-options="deptSelectOptions"
      :initial-value="drawerInitialValue"
      @submit="handleDrawerSubmit"
    />

    <UserProfilesDrawer
      v-model:open="profilesVisible"
      :loading="profilesLoading"
      :user="profilesUser"
      :profiles="userProfiles"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { h } from 'vue'
import { Empty, message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, PlusOutlined, UserOutlined, DownOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import ProTable from '@/components/ProTable.vue'
import UserFormDrawer from './components/UserFormDrawer.vue'
import UserProfilesDrawer from './components/UserProfilesDrawer.vue'
import { getDictDataByType } from '@/api/dict'
import {
  getUserList,
  createUser,
  updateUser,
  deleteUser,
  batchDeleteUsers,
  updateUserStatus,
  batchUpdateUserStatus,
  batchResetPassword,
  resetPassword,
  forceUserOffline,
  getUserProfilesById,
  type UserProfile
} from '@/api/user'
import { getRoleList } from '@/api/role'
import { getManageableDeptTree } from '@/api/dept'
import { formatTime } from '@/utils/format'
import { useTableColumns } from '@/utils/permission'
import { useConfigStore } from '@/store/config'
import { useUserStore } from '@/store/user'
import type { Dept, Role, User } from '@/types'
import { normalizeGenderDictOptions, resolveGenderOption, type GenderOption } from './user-gender'

interface DeptTreeNode {
  key: string
  title: any
  deptId?: number
  selectableType: 'all' | 'unassigned' | 'dept'
  children?: DeptTreeNode[]
}

interface TreeSelectOption {
  key: string
  title: string
  value: number
  disabled?: boolean
  selectable?: boolean
  isLeaf?: boolean
  children?: TreeSelectOption[]
}

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const configStore = useConfigStore()
const userStore = useUserStore()
const loading = ref(false)
const deptLoading = ref(false)
const tableData = ref<User[]>([])
const roleList = ref<Role[]>([])
const genderOptions = ref<GenderOption[]>([])
const deptTree = ref<Dept[]>([])
const deptSelectTree = ref<Dept[]>([])
const unassignedUserCount = ref(0)
const defaultUserAvatarUrl = ref('')
const selectedTreeKey = ref<string>('all')
const expandedTreeKeys = ref<string[]>([])
const treeInitialized = ref(false)
const drawerVisible = ref(false)
const drawerTitle = ref('新增用户')
const isEdit = ref(false)
const currentId = ref(0)
const selectedRowKeys = ref<number[]>([])
const profilesVisible = ref(false)
const profilesLoading = ref(false)
const profilesUser = ref<User | null>(null)
const userProfiles = ref<UserProfile[]>([])

const showProfileButton = computed(() => {
  const value = configStore.get('user_profile_button_visible')
  return value === 'true' || value === '1'
})

const selectedUsers = computed(() =>
  tableData.value.filter(item => selectedRowKeys.value.includes(item.id))
)

const currentUserId = computed(() => userStore.user?.id ?? 0)
const hasRestrictedManagedSelection = computed(() =>
  selectedUsers.value.some(user => isRestrictedManagedRecord(user))
)

const searchForm = reactive({
  username: '',
  status: undefined as number | undefined,
  gender: undefined as number | undefined,
  roleId: undefined as number | undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const selectedTreeKeys = computed(() => [selectedTreeKey.value])

const columns = useTableColumns(
  [
    { title: '头像', key: 'avatar', width: 80 },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    { title: '性别', key: 'gender', width: 90 },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '所属部门', key: 'dept' },
    { title: '状态', key: 'status' },
    { title: '角色', key: 'roles' },
    { title: '创建时间', key: 'created_at' }
  ],
  { title: '操作', key: 'action', width: 200, fixed: 'right' },
  ['system:user:edit', 'system:user:delete', 'system:user:resetPwd']
)

const deptTreeNodes = computed<DeptTreeNode[]>(() => {
  const deptChildren = buildDeptTreeNodes(deptTree.value)
  return [
    {
      key: 'all',
      title: `全部部门 (${getDeptTreeTotalUsers(deptTree.value) + unassignedUserCount.value})`,
      selectableType: 'all',
      children: [
        ...deptChildren,
        {
          key: 'unassigned',
          title: h('span', { class: 'unassigned-tree-node' }, `未绑定部门 (${unassignedUserCount.value})`),
          selectableType: 'unassigned'
        }
      ]
    }
  ]
})

const deptSelectOptions = computed<TreeSelectOption[]>(() => buildDeptSelectOptions(deptSelectTree.value))

const drawerInitialValue = ref<Record<string, any>>({})

const fetchData = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: pagination.current,
      page_size: pagination.pageSize,
      username: searchForm.username,
      status: searchForm.status,
      gender: searchForm.gender,
      role_id: searchForm.roleId
    }

    if (selectedTreeKey.value === 'unassigned') {
      params.unassigned_dept = true
    } else {
      const deptId = getSelectedDeptId()
      if (deptId) {
        params.dept_id = deptId
      }
    }

    const res = await getUserList(params)
    tableData.value = res.data.list
    pagination.total = res.data.total
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  const res = await getRoleList()
  roleList.value = res.data
}

const fetchGenderOptions = async () => {
  const res = await getDictDataByType('sys_gender')
  genderOptions.value = normalizeGenderDictOptions(res.data || [])
}

const showDeptTreeErrorModal = (content: string) => {
  Modal.error({
    title: '部门树加载失败',
    content,
    okText: '知道了'
  })
}

const fetchDeptTree = async () => {
  deptLoading.value = true
  try {
    const [userTreeRes, deptSelectRes] = await Promise.all([
      getManageableDeptTree('system:user-management', { silent: true }),
      getManageableDeptTree('system:user-management', { silent: true })
    ])

    const tree = Array.isArray(userTreeRes.data?.tree) ? userTreeRes.data.tree : []
    const unassignedCount = typeof userTreeRes.data?.unassigned_user_count === 'number' ? userTreeRes.data.unassigned_user_count : 0
    const defaultAvatarUrl = typeof userTreeRes.data?.default_avatar_url === 'string' ? userTreeRes.data.default_avatar_url : ''
    const selectTree = Array.isArray(deptSelectRes.data?.tree) ? deptSelectRes.data.tree : []

    if (!Array.isArray(userTreeRes.data?.tree)) {
      showDeptTreeErrorModal('可管理部门树数据异常，请刷新页面后重试。')
    }

    deptTree.value = tree
    deptSelectTree.value = selectTree
    unassignedUserCount.value = unassignedCount
    defaultUserAvatarUrl.value = defaultAvatarUrl
    const allKeys = collectTreeKeys(tree)
    expandedTreeKeys.value = treeInitialized.value
      ? expandedTreeKeys.value.filter(key => allKeys.includes(key))
      : allKeys
    treeInitialized.value = true
  } catch (error) {
    console.error('获取可管理部门树失败:', error)
    deptTree.value = []
    deptSelectTree.value = []
    unassignedUserCount.value = 0
    defaultUserAvatarUrl.value = ''
    showDeptTreeErrorModal('获取可管理部门树失败，请稍后重试。')
  } finally {
    deptLoading.value = false
  }
}

const handleSearch = () => {
  clearSelectedRows()
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  clearSelectedRows()
  searchForm.username = ''
  searchForm.status = undefined
  searchForm.gender = undefined
  searchForm.roleId = undefined
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pag: any) => {
  clearSelectedRows()
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchData()
}

const handleDeptSelect = (keys: string[]) => {
  clearSelectedRows()
  selectedTreeKey.value = keys[0] ?? 'all'
  pagination.current = 1
  fetchData()
}

const handleTreeExpand = (keys: string[]) => {
  expandedTreeKeys.value = keys
}

const handleAdd = () => {
  const defaultDeptId = getDefaultDeptIdFromSelection()
  const defaultDept = defaultDeptId ? findDeptById(deptTree.value, defaultDeptId) : null
  const defaultAvatarUrl = defaultUserAvatarUrl.value || configStore.get('register_logo')
  isEdit.value = false
  currentId.value = 0
  drawerTitle.value = '新增用户'
  drawerInitialValue.value = {
    username: '',
    password: '123456',
    nickname: '',
    gender: 0,
    email: '',
    phone: '',
    dept_id: defaultDeptId,
    dept_label: defaultDept?.name,
    role_ids: [2],
    statusChecked: true,
    avatar_file_id: undefined,
    avatar_file_url: defaultAvatarUrl || ''
  }
  drawerVisible.value = true
}

const handleEdit = (record: User) => {
  if (!canMutateManagedRecord(record)) {
    message.warning(isCurrentUserRecord(record) ? '不能在用户管理中修改当前登录账号' : '受保护管理员账号不允许编辑')
    return
  }
  isEdit.value = true
  currentId.value = record.id
  drawerTitle.value = '编辑用户'
  drawerInitialValue.value = {
    username: record.username,
    nickname: record.nickname,
    gender: record.gender ?? 0,
    email: record.email,
    phone: record.phone,
    dept_id: record.dept_id || undefined,
    dept_label: record.dept?.name,
    role_ids: record.roles?.map(r => r.id) || [],
    statusChecked: record.status === 1,
    avatar_file_id: record.avatar_file_id,
    avatar_file_url: record.avatar_file_url || ''
  }
  drawerVisible.value = true
}

const handleDrawerSubmit = async (values: any) => {
  try {
    if (isEdit.value) {
      await updateUser(currentId.value, values)
      message.success('更新成功')
    } else {
      await createUser(values)
      message.success('创建成功')
    }
    drawerVisible.value = false
    await Promise.all([fetchDeptTree(), fetchData()])
  } catch {
    // handled by interceptor
  }
}

const handleStatusChange = async (record: User, checked: boolean) => {
  if (!canMutateManagedRecord(record)) {
    message.warning(isCurrentUserRecord(record) ? '不能修改当前登录账号状态' : '受保护管理员账号不允许修改状态')
    return
  }
  await updateUserStatus(record.id, checked ? 1 : 0)
  message.success('修改成功')
  fetchData()
}

const handleResetPwd = async (record: User) => {
  if (!canMutateManagedRecord(record)) {
    message.warning(isCurrentUserRecord(record) ? '不能在用户管理中重置当前登录账号密码' : '受保护管理员账号不允许重置密码')
    return
  }
  Modal.confirm({
    title: '确认重置密码',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要将用户「${record.username}」的密码重置为系统设置中的用户默认密码吗？`,
    okText: '确认重置',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await resetPassword(record.id)
      message.success('密码已重置为系统默认密码')
    }
  })
}

const confirmDelete = (record: User) => {
  if (!canDeleteRecord(record)) {
    message.warning(record.id === currentUserId.value ? '不能删除当前登录账号' : '受保护管理员账号不允许删除')
    return
  }
  Modal.confirm({
    title: '确认删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除用户「${record.username}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteUser(record.id)
      message.success('删除成功')
      await Promise.all([fetchDeptTree(), fetchData()])
    }
  })
}

const onSelectChange = (keys: number[]) => {
  selectedRowKeys.value = keys
}

const clearSelectedRows = () => {
  selectedRowKeys.value = []
}

const isCurrentUserRecord = (record: User) => record.id === currentUserId.value

const isProtectedDeleteRecord = (record: User) => {
  if (record.id === 1 || record.username === 'admin') {
    return true
  }
  return (record.roles || []).some(role => role.id === 1 || role.code === 'admin' || role.code === 'super_admin')
}

const isRestrictedManagedRecord = (record: User) =>
  isCurrentUserRecord(record) || isProtectedDeleteRecord(record)

const canMutateManagedRecord = (record: User) =>
  !isRestrictedManagedRecord(record)

const canDeleteRecord = (record: User) =>
  canMutateManagedRecord(record)

const getBatchStatusTargetIds = (status: number) => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要修改状态的用户')
    return null
  }

  if (selectedUsers.value.some(user => isProtectedDeleteRecord(user))) {
    message.warning('当前选择包含受保护管理员账号，无法批量修改状态')
    return null
  }

  if (selectedUsers.value.some(user => isCurrentUserRecord(user))) {
    message.warning('不能批量修改当前登录账号状态')
    return null
  }

  return [...selectedRowKeys.value]
}

const handleBatchStatusChange = (status: number) => {
  const targetIds = getBatchStatusTargetIds(status)
  if (!targetIds || targetIds.length === 0) {
    return
  }

  const actionText = status === 1 ? '启用' : '禁用'
  Modal.confirm({
    title: `确认批量${actionText}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要批量${actionText}选中的 ${targetIds.length} 个用户吗？`,
    okText: `批量${actionText}`,
    okType: status === 0 ? 'danger' : 'primary',
    cancelText: '取消',
    async onOk() {
      await batchUpdateUserStatus(targetIds, status)
      message.success(`批量${actionText}成功`)
      clearSelectedRows()
      await Promise.all([fetchDeptTree(), fetchData()])
    }
  })
}

const handleBatchResetPwd = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要重置密码的用户')
    return
  }
  if (selectedUsers.value.some(user => isCurrentUserRecord(user))) {
    message.warning('不能批量重置当前登录账号密码')
    return
  }
  if (selectedUsers.value.some(user => isProtectedDeleteRecord(user))) {
    message.warning('当前选择包含受保护管理员账号，无法批量重置密码')
    return
  }

  Modal.confirm({
    title: '确认批量重置密码',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要将选中的 ${selectedRowKeys.value.length} 个用户密码重置为系统设置中的用户默认密码吗？`,
    okText: '确认重置',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await batchResetPassword(selectedRowKeys.value)
      message.success('批量重置密码成功')
      clearSelectedRows()
      await fetchData()
    }
  })
}

const handleBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的用户')
    return
  }
  if (selectedUsers.value.some(user => isCurrentUserRecord(user))) {
    message.warning('不能批量删除当前登录账号')
    return
  }
  if (selectedUsers.value.some(user => isProtectedDeleteRecord(user))) {
    message.warning('当前选择包含受保护管理员账号，无法批量删除')
    return
  }
  Modal.confirm({
    title: '确认批量删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个用户吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      const res = await batchDeleteUsers(selectedRowKeys.value)
      if (res.data?.failed_count > 0) {
        message.warning(`成功删除 ${res.data.success_count} 个，失败 ${res.data.failed_count} 个`)
      } else {
        message.success('批量删除成功')
      }
      clearSelectedRows()
      await Promise.all([fetchDeptTree(), fetchData()])
    }
  })
}

const confirmForceOffline = (record: User) => {
  if (!canMutateManagedRecord(record)) {
    message.warning(isCurrentUserRecord(record) ? '不能强制下线自己' : '受保护管理员账号不允许强制下线')
    return
  }
  Modal.confirm({
    title: '确认强制下线',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要强制用户「${record.username}」下线吗？`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await forceUserOffline(record.id)
      message.success('已强制该用户下线')
    }
  })
}

const handleViewProfiles = async (record: User) => {
  profilesUser.value = record
  profilesVisible.value = true
  profilesLoading.value = true
  try {
    const res = await getUserProfilesById(record.id)
    userProfiles.value = res.data || []
  } catch {
    userProfiles.value = []
  } finally {
    profilesLoading.value = false
  }
}

const getSelectedDeptId = () => {
  if (selectedTreeKey.value === 'all' || selectedTreeKey.value === 'unassigned') {
    return undefined
  }
  const prefix = 'dept-'
  if (!selectedTreeKey.value.startsWith(prefix)) {
    return undefined
  }
  const raw = selectedTreeKey.value.slice(prefix.length)
  const deptId = Number(raw)
  return Number.isFinite(deptId) ? deptId : undefined
}

const findDeptById = (depts: Dept[], targetId: number): Dept | null => {
  for (const dept of depts) {
    if (dept.id === targetId) {
      return dept
    }
    const found = dept.children ? findDeptById(dept.children, targetId) : null
    if (found) {
      return found
    }
  }
  return null
}

const getDefaultDeptIdFromSelection = () => {
  const deptId = getSelectedDeptId()
  if (!deptId) {
    return undefined
  }
  const dept = findDeptById(deptSelectTree.value, deptId)
  if (!dept || !dept.bindable) {
    return undefined
  }
  return dept.id
}

const getDeptDisplayCount = (dept: Dept) =>
  dept.has_children ? dept.total_user_count || 0 : dept.direct_user_count || 0

const getDeptTreeTotalUsers = (depts: Dept[]): number =>
  depts.reduce((sum, dept) => sum + (dept.total_user_count || dept.direct_user_count || 0), 0)

const buildDeptTreeNodes = (depts: Dept[]): DeptTreeNode[] =>
  depts.map(dept => ({
    key: `dept-${dept.id}`,
    title: `${dept.name} (${getDeptDisplayCount(dept)})`,
    deptId: dept.id,
    selectableType: 'dept',
    children: dept.children ? buildDeptTreeNodes(dept.children) : undefined
  }))

const collectTreeKeys = (depts: Dept[]) => {
  const keys = ['all', 'unassigned']

  const walk = (items: Dept[]) => {
    items.forEach((dept) => {
      keys.push(`dept-${dept.id}`)
      if (dept.children?.length) {
        walk(dept.children)
      }
    })
  }

  walk(depts)
  return keys
}

const expandAllTree = () => {
  expandedTreeKeys.value = collectTreeKeys(deptTree.value)
}

const collapseAllTree = () => {
  expandedTreeKeys.value = ['all']
}

const buildDeptSelectOptions = (depts: Dept[]): TreeSelectOption[] =>
  depts.map(dept => ({
    key: `dept-option-${dept.id}`,
    title: `${dept.name} (${getDeptDisplayCount(dept)})`,
    value: dept.id,
    disabled: dept.manageable === false || dept.bindable !== true,
    selectable: dept.bindable === true,
    isLeaf: !dept.has_children,
    children: dept.children ? buildDeptSelectOptions(dept.children) : undefined
  }))

const getGenderOption = (value: number) => resolveGenderOption(genderOptions.value, value)

onMounted(async () => {
  await Promise.all([configStore.loadConfigs(), fetchRoles(), fetchDeptTree(), fetchGenderOptions()])
  await fetchData()
})
</script>

<style scoped>
.page-layout {
  display: flex;
  gap: 16px;
}

.left-tree {
  width: 260px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 4px;
  padding: 12px;
  border: 1px solid #f0f0f0;
}

.tree-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 500;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.left-tree :deep(.ant-tree) {
  max-height: 620px;
  overflow-y: auto;
}

.right-content {
  flex: 1;
  min-width: 0;
}

.left-tree :deep(.unassigned-tree-node) {
  color: #cf1322;
  font-weight: 500;
}
</style>
