<template>
  <div class="user-page">
    <ProTable
      :columns="columns"
      :data-source="tableData"
      :loading="loading"
      :pagination="pagination"
      @search="handleSearch"
      @reset="handleReset"
      @change="handleTableChange"
    >
      <!-- 搜索区域 -->
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
        <a-form-item label="角色">
          <a-select v-model:value="searchForm.roleId" placeholder="请选择角色" allowClear style="width: 150px">
            <a-select-option v-for="role in roleList" :key="role.id" :value="role.id">
              {{ role.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </template>

      <!-- 工具栏 -->
      <template #toolbar>
        <a-button type="primary" @click="handleAdd" v-permission="'system:user:add'">
          <PlusOutlined /> 新增
        </a-button>
      </template>

      <!-- 表格单元格 -->
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'avatar'">
          <a-avatar :size="40" :src="record.avatar_file_url">
            <template #icon><UserOutlined /></template>
          </a-avatar>
        </template>
        <template v-if="column.key === 'status'">
          <a-switch
            :checked="record.status === 1"
            @change="(checked: boolean) => handleStatusChange(record, checked)"
          />
        </template>
        <template v-if="column.key === 'roles'">
          <a-tag v-for="role in record.roles" :key="role.id" color="blue">
            {{ role.name }}
          </a-tag>
        </template>
        <template v-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
        <template v-if="column.key === 'action'">
          <a-space :size="0">
            <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:user:edit'">编辑</a-button>
            <a-button type="link" size="small" @click="handleViewProfiles(record)">身份</a-button>
            <a-dropdown>
              <a-button type="link" size="small">更多 <DownOutlined /></a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="resetPwd" v-permission="'system:user:resetPwd'" @click="handleResetPwd(record)">重置密码</a-menu-item>
                  <a-menu-item key="delete" v-permission="'system:user:delete'" @click="confirmDelete(record)">
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

    <!-- 用户身份弹窗 -->
    <a-modal
      v-model:open="profilesVisible"
      :title="`用户身份 - ${profilesUser?.username || ''}`"
      :footer="null"
      width="50%"
    >
      <a-spin :spinning="profilesLoading">
        <a-empty v-if="!profilesLoading && userProfiles.length === 0" description="未绑定任何身份" />
        <a-collapse v-else v-model:activeKey="profilesActiveKey" accordion>
          <a-collapse-panel v-for="profile in userProfiles" :key="profile.key" :header="profile.name">
            <template #extra>
              <a-tag v-if="profile.has_profile" :color="profile.is_complete ? 'success' : 'warning'">
                {{ profile.is_complete ? '已完善' : '未完善' }}
              </a-tag>
              <a-tag v-else color="default">未填写</a-tag>
            </template>
            <a-descriptions v-if="profile.has_profile && profile.data" :column="2" size="small" bordered>
              <a-descriptions-item v-for="field in profile.fields" :key="field.key" :label="field.label">
                <!-- 图片类型 -->
                <template v-if="field.type === 'image'">
                  <a-image v-if="getFieldValue(profile.data, field.key)" :src="getFieldValue(profile.data, field.key)" :width="80" />
                  <span v-else>-</span>
                </template>
                <!-- 文件类型 -->
                <template v-else-if="field.type === 'file'">
                  <a v-if="getFieldValue(profile.data, field.key)" :href="getFieldValue(profile.data, field.key)" target="_blank">查看文件</a>
                  <span v-else>-</span>
                </template>
                <!-- 多图片类型 -->
                <template v-else-if="field.type === 'images'">
                  <a-image-preview-group v-if="getFieldValue(profile.data, field.key)">
                    <a-space>
                      <a-image v-for="(url, idx) in getFieldValue(profile.data, field.key).split(',')" :key="idx" :src="url" :width="60" />
                    </a-space>
                  </a-image-preview-group>
                  <span v-else>-</span>
                </template>
                <!-- 审批状态特殊处理 -->
                <template v-else-if="field.key === 'audit_status'">
                  <a-tag v-if="getFieldValue(profile.data, field.key) === 0" color="default">待审批</a-tag>
                  <a-tag v-else-if="getFieldValue(profile.data, field.key) === 1" color="success">审批通过</a-tag>
                  <a-tag v-else-if="getFieldValue(profile.data, field.key) === 2" color="error">审批拒绝</a-tag>
                  <span v-else>-</span>
                </template>
                <!-- 默认文本 -->
                <template v-else>
                  {{ getFieldValue(profile.data, field.key) || '-' }}
                </template>
              </a-descriptions-item>
            </a-descriptions>
            <a-empty v-else description="未填写档案信息" />
          </a-collapse-panel>
        </a-collapse>
      </a-spin>
    </a-modal>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 5 }">
        <a-form-item label="头像" name="avatar_file_id">
          <AvatarUpload
            v-model:fileId="formState.avatar_file_id"
            :url="formState.avatar_file_url"
            :size="80"
            tip=""
          />
        </a-form-item>
        <a-form-item label="用户名" name="username">
          <a-input v-model:value="formState.username" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="密码" name="password" v-if="!isEdit">
          <a-input-password v-model:value="formState.password" />
        </a-form-item>
        <a-form-item label="昵称">
          <a-input v-model:value="formState.nickname" />
        </a-form-item>
        <a-form-item label="邮箱">
          <a-input v-model:value="formState.email" />
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model:value="formState.phone" />
        </a-form-item>
        <a-form-item label="角色">
          <a-select v-model:value="formState.role_ids" mode="multiple" placeholder="请选择角色">
            <a-select-option v-for="role in roleList" :key="role.id" :value="role.id">
              {{ role.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="formState.statusChecked" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message, Modal, type FormInstance } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { PlusOutlined, UserOutlined, DownOutlined } from '@ant-design/icons-vue'
import AvatarUpload from '@/components/AvatarUpload.vue'
import ProTable from '@/components/ProTable.vue'
import { getUserList, createUser, updateUser, deleteUser, updateUserStatus, resetPassword, forceUserOffline, getUserProfilesById, type UserProfile } from '@/api/user'
import { getRoleList } from '@/api/role'
import { formatTime } from '@/utils/format'
import { useTableColumns } from '@/utils/permission'
import type { User, Role } from '@/types'
import type { Rule } from 'ant-design-vue/es/form'
const formRef = ref<FormInstance>()
const loading = ref(false)
const tableData = ref<User[]>([])
const roleList = ref<Role[]>([])
const modalVisible = ref(false)
const modalTitle = ref('新增用户')
const isEdit = ref(false)
const currentId = ref(0)

// 用户身份弹窗
const profilesVisible = ref(false)
const profilesLoading = ref(false)
const profilesUser = ref<User | null>(null)
const userProfiles = ref<UserProfile[]>([])
const profilesActiveKey = ref<string>('')

const searchForm = reactive({
  username: '',
  status: undefined as number | undefined,
  roleId: undefined as number | undefined
})
// 表单验证规则
const formRules: Record<string, Rule[]> = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  avatar_file_id: [{ required: true, message: '请输入头像', trigger: 'blur' }],
}
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formState = reactive({
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  role_ids: [1] as number[],
  statusChecked: true,
  avatar_file_id: undefined as number | undefined,
  avatar_file_url: ''
})

// 使用工具函数动态生成列配置
const columns = useTableColumns(
  [
    { title: '头像', key: 'avatar', width: 80 },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '状态', key: 'status' },
    { title: '角色', key: 'roles' },
    { title: '创建时间', key: 'created_at' },
  ],
  { title: '操作', key: 'action', width: 200,fixed: 'right' },
  ['system:user:edit', 'system:user:delete', 'system:user:resetPwd']
)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getUserList({
      page: pagination.current,
      page_size: pagination.pageSize,
      username: searchForm.username,
      status: searchForm.status,
      role_id: searchForm.roleId,
    })
    tableData.value = res.data.list
    pagination.total = res.data.total
  } catch {
    // 错误已由 request 拦截器处理
  } finally {
    loading.value = false
  }
}

const fetchRoles = async () => {
  try {
    const res = await getRoleList()
    roleList.value = res.data
  } catch {
    // 错误已由 request 拦截器处理
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.status = undefined
  searchForm.roleId = undefined
  pagination.current = 1
  fetchData()
}

const handleTableChange = (pag: any) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchData()
}

const handleAdd = () => {
  isEdit.value = false
  modalTitle.value = '新增用户'
  Object.assign(formState, {
    username: '', password: '123456', nickname: '', email: '', phone: '', role_ids: [2], statusChecked: true,
    avatar_file_id: undefined, avatar_file_url: ''
  })
  modalVisible.value = true
}

const handleEdit = (record: User) => {
  isEdit.value = true
  modalTitle.value = '编辑用户'
  currentId.value = record.id
  Object.assign(formState, {
    username: record.username,
    nickname: record.nickname,
    email: record.email,
    phone: record.phone,
    role_ids: record.roles?.map(r => r.id) || [],
    statusChecked: record.status === 1,
    avatar_file_id: record.avatar_file_id,
    avatar_file_url: record.avatar_file_url || ''
  })
  modalVisible.value = true
}

const handleModalOk = async () => {
  // 提交时验证表单
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  
  const data: any = {
    username: formState.username,
    password: formState.password,
    nickname: formState.nickname,
    email: formState.email,
    phone: formState.phone,
    role_ids: formState.role_ids,
    status: formState.statusChecked ? 1 : 0
  }
  if (formState.avatar_file_id) {
    data.avatar_file_id = formState.avatar_file_id
  }
  
  try {
    if (isEdit.value) {
      await updateUser(currentId.value, data)
      message.success('更新成功')
    } else {
      await createUser(data)
      message.success('创建成功')
    }
    modalVisible.value = false
    // 清空表单数据
    formRef.value?.resetFields()
    Object.assign(formState, {
      username: '',
      password: '',
      nickname: '',
      email: '',
      phone: '',
      role_ids: [],
      statusChecked: true,
      avatar_file_id: undefined,
      avatar_file_url: ''
    })
    fetchData()
  } catch {
    // 错误已由 request 拦截器处理，这里只需捕获防止冒泡到 ErrorBoundary
  }
}

const handleModalCancel = () => {
  modalVisible.value = false
  // 清空表单数据和验证状态
  formRef.value?.resetFields()
  Object.assign(formState, {
    username: '',
    password: '',
    nickname: '',
    email: '',
    phone: '',
    role_ids: [],
    statusChecked: true,
    avatar_file_id: undefined,
    avatar_file_url: ''
  })
}

const handleStatusChange = async (record: User, checked: boolean) => {
  try {
    await updateUserStatus(record.id, checked ? 1 : 0)
    message.success('修改成功')
    fetchData()
  } catch {
    // 错误已由 request 拦截器处理
  }
}

const handleResetPwd = async (record: User) => {
  try {
    await resetPassword(record.id, '123456')
    message.success('密码已重置为 123456')
  } catch {
    // 错误已由 request 拦截器处理
  }
}

// 确认删除用户
const confirmDelete = (record: User) => {
  Modal.confirm({
    title: '确认删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除用户「${record.username}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteUser(record.id)
        message.success('删除成功')
        fetchData()
      } catch {
        // 错误已由 request 拦截器处理
      }
    }
  })
}

// 确认强制下线
const confirmForceOffline = (record: User) => {
  Modal.confirm({
    title: '确认强制下线',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要强制用户「${record.username}」下线吗？`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await forceUserOffline(record.id)
        message.success('已强制该用户下线')
      } catch {
        // 错误已由 request 拦截器处理
      }
    }
  })
}

// 查看用户身份
const handleViewProfiles = async (record: User) => {
  profilesUser.value = record
  profilesVisible.value = true
  profilesLoading.value = true
  profilesActiveKey.value = ''
  try {
    const res = await getUserProfilesById(record.id)
    userProfiles.value = res.data || []
    if (userProfiles.value.length > 0) {
      profilesActiveKey.value = userProfiles.value[0].key
    }
  } catch {
    userProfiles.value = []
  } finally {
    profilesLoading.value = false
  }
}

// 获取字段值
const getFieldValue = (data: any, key: string) => {
  if (!data) return ''
  return data[key]
}

onMounted(() => {
  fetchData()
  fetchRoles()
})
</script>

<style scoped>
/* 样式由 ProTable 组件管理 */
</style>
