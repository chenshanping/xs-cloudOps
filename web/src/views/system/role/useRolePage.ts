import { computed, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { getRoleList, getRole, createRole, updateRole, deleteRole } from '@/api/role'
import { getManageableDeptTree } from '@/api/dept'
import { useTableColumns } from '@/utils/permission'
import type { Dept, Role } from '@/types'

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
  const deptTree = ref<Dept[]>([])
  const drawerVisible = ref(false)
  const drawerTitle = ref('新增角色')
  const isEdit = ref(false)
  const currentId = ref(0)
  const currentRoleName = ref('')
  const permissionDrawerVisible = ref(false)
  const drawerInitialValue = ref<Record<string, any>>({})

  const columns = useTableColumns(
    [
      { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
      { title: '角色名称', dataIndex: 'name', key: 'name' },
      { title: '角色编码', dataIndex: 'code', key: 'code' },
      { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
      { title: '状态', key: 'status', width: 80 },
      { title: '数据范围', key: 'data_scope', width: 140 },
      { title: '备注', dataIndex: 'remark', key: 'remark' }
    ],
    { title: '操作', key: 'action', width: 200 },
    ['system:role:edit', 'system:role:delete', 'system:role:assign']
  )

  const deptSelectOptions = computed<TreeSelectOption[]>(() => buildDeptSelectOptions(deptTree.value))

  const fetchData = async () => {
    loading.value = true
    try {
      const res = await getRoleList()
      tableData.value = res.data
    } finally {
      loading.value = false
    }
  }

  const fetchDepts = async () => {
    const res = await getManageableDeptTree()
    deptTree.value = res.data.tree
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
      data_scope: 1,
      dept_ids: [],
      remark: ''
    }
    drawerVisible.value = true
  }

  const handleEdit = async (record: Role) => {
    isEdit.value = true
    drawerTitle.value = '编辑角色'
    currentId.value = record.id
    const res = await getRole(record.id)
    const role = res.data
    drawerInitialValue.value = {
      name: role.name,
      code: role.code,
      sort: role.sort,
      statusChecked: role.status === 1,
      data_scope: role.data_scope || 1,
      dept_ids: role.depts?.map(item => item.id) || [],
      remark: role.remark
    }
    drawerVisible.value = true
  }

  const handleDrawerSubmit = async (values: any) => {
    if (isEdit.value) {
      await updateRole(currentId.value, values)
      message.success('更新成功')
    } else {
      await createRole(values)
      message.success('创建成功')
    }
    drawerVisible.value = false
    fetchData()
  }

  const handleDelete = async (record: Role) => {
    await deleteRole(record.id)
    message.success('删除成功')
    fetchData()
  }

  const handleAssignPermissions = (record: Role) => {
    currentId.value = record.id
    currentRoleName.value = record.name
    permissionDrawerVisible.value = true
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
    handleAssignPermissions,
    handleDelete,
    handleDrawerSubmit,
    handleEdit,
    isEdit,
    loading,
    permissionDrawerVisible,
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
