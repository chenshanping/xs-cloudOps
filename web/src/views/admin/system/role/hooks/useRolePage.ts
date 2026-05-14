import { computed, onMounted, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  getRoleList,
  getRole,
  createRole,
  updateRole,
  deleteRole,
  type RoleUpsertPayload
} from '@/api/role'
import { getManageableDeptTree } from '@/api/dept'
import { getUserList } from '@/api/user'
import { useTableColumns } from '@/utils/permission'
import type { Dept, Role, User } from '@/types'

interface TreeSelectOption {
  key: string | number
  title: string
  value: number
  disabled?: boolean
  children?: TreeSelectOption[]
}

export function useRolePage() {
  const loading = ref(false)
  const tableData = ref<Role[]>([])
  const searchKeyword = ref('')
  const deptTree = ref<Dept[]>([])
  const drawerVisible = ref(false)
  const drawerTitle = ref('新增角色')
  const isEdit = ref(false)
  const currentId = ref(0)
  const currentRoleName = ref('')
  const menuPermissionDrawerVisible = ref(false)
  const apiPermissionDrawerVisible = ref(false)
  const dataScopeDrawerVisible = ref(false)
  const roleUsersDrawerVisible = ref(false)
  const roleUsersLoading = ref(false)
  const roleUsers = ref<User[]>([])
  const roleUsersRoleId = ref(0)
  const roleUsersRoleName = ref('')
  const statusLoadingMap = ref<Record<number, boolean>>({})
  const superAdminLoadingMap = ref<Record<number, boolean>>({})
  const drawerInitialValue = ref<Record<string, any>>({})
  const currentRoleScopeState = ref({
    data_scope: 1,
    dept_ids: [] as number[]
  })
  const roleUsersPagination = ref({
    current: 1,
    pageSize: 10,
    total: 0
  })

  const columns = useTableColumns(
    [
      { title: '角色名称', dataIndex: 'name', key: 'name', width: 220 },
      { title: '角色编码', dataIndex: 'code', key: 'code', width: 180 },
      { title: '超管', key: 'is_super_admin', width: 90 },
      { title: '关联用户', key: 'users', width: 160 },
      { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
      { title: '状态', key: 'status', width: 80 },
      { title: '默认数据范围', key: 'data_scope', width: 150 },
    ],
    { title: '操作', key: 'action', width: 260, fixed: 'right' },
    ['system:role:edit', 'system:role:delete', 'system:role:assign', 'system:role:dataScope']
  )

  const deptSelectOptions = computed<TreeSelectOption[]>(() => buildDeptSelectOptions(deptTree.value))
  const filteredTableData = computed(() => {
    const keyword = searchKeyword.value.trim().toLowerCase()
    if (!keyword) {
      return tableData.value
    }

    return tableData.value.filter(role =>
      role.name.toLowerCase().includes(keyword)
      || role.code.toLowerCase().includes(keyword)
    )
  })

  const fetchData = async () => {
    loading.value = true
    try {
      const res = await getRoleList()
      tableData.value = res.data
    } finally {
      loading.value = false
    }
  }

  const showDeptTreeErrorModal = (content: string) => {
    Modal.error({
      title: '部门树加载失败',
      content,
      okText: '知道了'
    })
  }

  const fetchDepts = async () => {
    try {
      const res = await getManageableDeptTree({ silent: true })
      const tree = Array.isArray(res.data?.tree) ? res.data.tree : []
      if (!Array.isArray(res.data?.tree)) {
        showDeptTreeErrorModal('可管理部门树数据异常，请刷新页面后重试。')
      }
      deptTree.value = tree
    } catch (error) {
      console.error('获取可管理部门树失败:', error)
      deptTree.value = []
      showDeptTreeErrorModal('获取可管理部门树失败，请稍后重试。')
    }
  }

  const handleAdd = () => {
    isEdit.value = false
    drawerTitle.value = '新增角色'
    currentId.value = 0
    drawerInitialValue.value = {
      name: '',
      code: '',
      sort: 0,
      statusChecked: true,
      is_super_admin: false,
      remark: ''
    }
    currentRoleScopeState.value = {
      data_scope: 1,
      dept_ids: []
    }
    drawerVisible.value = true
  }

  const handleEdit = async (record: Role) => {
    isEdit.value = true
    drawerTitle.value = '编辑角色'
    currentId.value = record.id
    const res = await getRole(record.id)
    const role = res.data
    currentRoleScopeState.value = {
      data_scope: role.data_scope || 1,
      dept_ids: role.depts?.map(item => item.id) || []
    }
    drawerInitialValue.value = {
      name: role.name,
      code: role.code,
      sort: role.sort,
      statusChecked: role.status === 1,
      is_super_admin: role.is_super_admin,
      remark: role.remark
    }
    drawerVisible.value = true
  }

  const handleDrawerSubmit = async (values: any) => {
    const payload: RoleUpsertPayload = {
      ...values,
      data_scope: isEdit.value ? currentRoleScopeState.value.data_scope : 1,
      dept_ids: isEdit.value ? [...currentRoleScopeState.value.dept_ids] : [],
    }
    if (isEdit.value) {
      await updateRole(currentId.value, payload)
      message.success('更新成功')
    } else {
      await createRole(payload)
      message.success('创建成功')
    }
    drawerVisible.value = false
    fetchData()
  }

  const handleDelete = async (record: Role) => {
    if (getRoleUserCount(record) > 0) {
      message.warning('当前角色下存在用户，无法删除')
      return
    }
    await deleteRole(record.id)
    message.success('删除成功')
    fetchData()
  }

  const getRoleUsers = (record: Role) => record.users || []
  const getRoleUserCount = (record: Role) => record.user_count ?? getRoleUsers(record).length

  const handleAssignMenuPermissions = (record: Role) => {
    currentId.value = record.id
    currentRoleName.value = record.name
    menuPermissionDrawerVisible.value = true
  }

  const handleAssignApiPermissions = (record: Role) => {
    currentId.value = record.id
    currentRoleName.value = record.name
    apiPermissionDrawerVisible.value = true
  }

  const handleAssignDataScope = (record: Role) => {
    currentId.value = record.id
    currentRoleName.value = record.name
    dataScopeDrawerVisible.value = true
  }

  const buildRolePayload = (record: Role, overrides: Partial<RoleUpsertPayload> = {}): RoleUpsertPayload => ({
    name: record.name,
    code: record.code,
    sort: record.sort,
    status: record.status,
    is_super_admin: record.is_super_admin,
    data_scope: record.data_scope || 1,
    dept_ids: record.depts?.map(item => item.id) || [],
    remark: record.remark || '',
    ...overrides
  })

  const fetchRoleUsers = async () => {
    if (!roleUsersRoleId.value) {
      roleUsers.value = []
      roleUsersPagination.value.total = 0
      return
    }

    roleUsersLoading.value = true
    try {
      const res = await getUserList({
        role_id: roleUsersRoleId.value,
        page: roleUsersPagination.value.current,
        page_size: roleUsersPagination.value.pageSize
      })
      roleUsers.value = res.data.list
      roleUsersPagination.value.total = res.data.total
    } finally {
      roleUsersLoading.value = false
    }
  }

  const handleViewRoleUsers = async (record: Role) => {
    roleUsersRoleId.value = record.id
    roleUsersRoleName.value = record.name
    roleUsersPagination.value.current = 1
    roleUsersDrawerVisible.value = true
    await fetchRoleUsers()
  }

  const handleRoleUsersPageChange = async (page: number) => {
    roleUsersPagination.value.current = page
    await fetchRoleUsers()
  }

  const handleRoleUsersDrawerOpenChange = (open: boolean) => {
    roleUsersDrawerVisible.value = open
    if (!open) {
      roleUsers.value = []
      roleUsersRoleId.value = 0
      roleUsersRoleName.value = ''
      roleUsersPagination.value.current = 1
      roleUsersPagination.value.total = 0
    }
  }

  const setRowLoading = (
    mapRef: typeof statusLoadingMap,
    roleId: number,
    loading: boolean
  ) => {
    mapRef.value = {
      ...mapRef.value,
      [roleId]: loading
    }
  }

  const isRoleStatusLoading = (roleId: number) => !!statusLoadingMap.value[roleId]

  const isRoleSuperAdminLoading = (roleId: number) => !!superAdminLoadingMap.value[roleId]

  const isRoleRowUpdating = (roleId: number) =>
    isRoleStatusLoading(roleId) || isRoleSuperAdminLoading(roleId)

  const handleStatusToggle = async (record: Role, checked: boolean) => {
    if (isRoleRowUpdating(record.id)) {
      return
    }

    const previous = record.status
    record.status = checked ? 1 : 0
    setRowLoading(statusLoadingMap, record.id, true)

    try {
      await updateRole(record.id, buildRolePayload(record, { status: checked ? 1 : 0 }))
      message.success('修改成功')
      await fetchData()
    } catch {
      record.status = previous
    } finally {
      setRowLoading(statusLoadingMap, record.id, false)
    }
  }

  const updateSuperAdminStatus = async (record: Role, checked: boolean) => {
    if (isRoleRowUpdating(record.id)) {
      return
    }

    const previous = record.is_super_admin
    record.is_super_admin = checked
    setRowLoading(superAdminLoadingMap, record.id, true)

    try {
      await updateRole(record.id, buildRolePayload(record, { is_super_admin: checked }))
      message.success('修改成功')
      await fetchData()
    } catch {
      record.is_super_admin = previous
    } finally {
      setRowLoading(superAdminLoadingMap, record.id, false)
    }
  }

  const handleSuperAdminToggle = (record: Role, checked: boolean) => {
    Modal.confirm({
      title: checked ? '确认设为超管角色' : '确认取消超管角色',
      content: checked
        ? `确定将角色「${record.name}」设为超管角色吗？`
        : `确定取消角色「${record.name}」的超管角色吗？`,
      okText: '确定',
      cancelText: '取消',
      async onOk() {
        await updateSuperAdminStatus(record, checked)
      }
    })
  }

  const formatDataScope = (value: number) => {
    switch (value) {
      case 1:
        return '全部数据'
      case 2:
        return '自定义部门'
      case 3:
        return '本部门'
      case 4:
        return '本部门及下级'
      case 5:
        return '仅本人'
      default:
        return '-'
    }
  }

  onMounted(async () => {
    await Promise.all([fetchData(), fetchDepts()])
  })

  return {
    columns,
    currentId,
    currentRoleName,
    deptSelectOptions,
    drawerInitialValue,
    drawerTitle,
    drawerVisible,
    formatDataScope,
    handleAdd,
    handleAssignApiPermissions,
    handleAssignDataScope,
    handleAssignMenuPermissions,
    handleDelete,
    handleDrawerSubmit,
    handleEdit,
    handleRoleUsersDrawerOpenChange,
    handleRoleUsersPageChange,
    handleStatusToggle,
    handleSuperAdminToggle,
    handleViewRoleUsers,
    getRoleUserCount,
    isEdit,
    isRoleRowUpdating,
    isRoleStatusLoading,
    isRoleSuperAdminLoading,
    loading,
    apiPermissionDrawerVisible,
    dataScopeDrawerVisible,
    filteredTableData,
    menuPermissionDrawerVisible,
    roleUsers,
    roleUsersDrawerVisible,
    roleUsersLoading,
    roleUsersPagination,
    roleUsersRoleName,
    searchKeyword,
    tableData
  }
}

const buildDeptSelectOptions = (depts: Dept[]): TreeSelectOption[] =>
  depts.map(dept => ({
    key: `role-dept-${dept.id}`,
    title: dept.name,
    value: dept.id,
    disabled: dept.manageable === false,
    children: dept.children ? buildDeptSelectOptions(dept.children) : undefined
  }))
