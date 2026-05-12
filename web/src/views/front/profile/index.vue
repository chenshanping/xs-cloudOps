<template>
  <div class="front-profile-page">
    <a-row :gutter="24">
      <!-- 左侧：个人信息卡片 + 身份信息 -->
      <a-col :xs="24" :md="8">
        <a-card class="profile-card">
          <div class="avatar-section">
            <AvatarUpload
              v-model:fileId="profileForm.avatar_file_id"
              v-model:url="profileForm.avatar"
              :size="120"
              tip=""
              @success="handleAvatarSuccess"
            />
            <h2 class="nickname">{{ userStore.user?.nickname || userStore.user?.username }}</h2>
            <p class="username">@{{ userStore.user?.username }}</p>
          </div>
          
          <a-divider />
          
          <div class="info-list">
            <div class="info-item">
              <MailOutlined class="info-icon" />
              <span class="info-label">邮箱</span>
              <span class="info-value">{{ userStore.user?.email || '未设置' }}</span>
            </div>
            <div class="info-item">
              <PhoneOutlined class="info-icon" />
              <span class="info-label">手机</span>
              <span class="info-value">{{ userStore.user?.phone || '未设置' }}</span>
            </div>
            <div class="info-item">
              <CalendarOutlined class="info-icon" />
              <span class="info-label">注册时间</span>
              <span class="info-value">{{ formatTime(userStore.user?.created_at) }}</span>
            </div>
          </div>
        </a-card>
        
        <!-- 身份信息卡片 -->
        <a-card 
          v-for="profile in userProfiles" 
          :key="profile.key" 
          class="identity-card" 
          :loading="profilesLoading"
          style="margin-top: 16px"
        >
          <template #title>
            <component :is="getIconComponent(profile.icon)" style="margin-right: 8px" />
            {{ profile.name }}
          </template>
          <template #extra>
            <a-space>
              <!-- 审批状态 -->
              <template v-if="profile.data?.audit_status !== undefined">
                <a-tag v-if="profile.data.audit_status === 0" color="orange">待审批</a-tag>
                <a-tag v-else-if="profile.data.audit_status === 1" color="green">已通过</a-tag>
                <a-tag v-else-if="profile.data.audit_status === 2" color="red">已拒绝</a-tag>
              </template>
              <!-- 完善状态 -->
              <a-tag v-if="profile.has_profile && profile.is_complete" color="blue">已完善</a-tag>
              <a-tag v-else-if="profile.has_profile" color="default">未完善</a-tag>
              <a-tag v-else color="default">未填写</a-tag>
            </a-space>
          </template>
          <div class="identity-content">
            <template v-if="profile.has_profile">
              <a-descriptions :column="1" size="small">
                <template v-for="field in profile.fields" :key="field.key">
                  <!-- 跳过审批状态字段（已在头部显示） -->
                  <a-descriptions-item v-if="field.key !== 'audit_status'" :label="field.label">
                    <!-- 图片类型 -->
                    <template v-if="field.type === 'image'">
                      <a-image
                        v-if="profile.data?.[field.key]"
                        :src="profile.data[field.key]"
                        :width="60"
                        style="border-radius: 4px"
                      />
                      <span v-else class="text-gray">-</span>
                    </template>
                    <!-- 文件类型 -->
                    <template v-else-if="field.type === 'file'">
                      <a v-if="profile.data?.[field.key]" :href="profile.data[field.key]" target="_blank">查看文件</a>
                      <span v-else class="text-gray">-</span>
                    </template>
                    <!-- 多图片类型 -->
                    <template v-else-if="field.type === 'images'">
                      <a-image-preview-group v-if="profile.data?.[field.key]">
                        <a-space>
                          <a-image v-for="(url, idx) in profile.data[field.key].split(',')" :key="idx" :src="url" :width="50" />
                        </a-space>
                      </a-image-preview-group>
                      <span v-else class="text-gray">-</span>
                    </template>
                    <!-- 多文件类型 -->
                    <template v-else-if="field.type === 'files'">
                      <a-space v-if="profile.data?.[field.key]" direction="vertical" size="small">
                        <a v-for="(url, idx) in profile.data[field.key].split(',')" :key="idx" :href="url" target="_blank">文件{{ idx + 1 }}</a>
                      </a-space>
                      <span v-else class="text-gray">-</span>
                    </template>
                    <!-- 审批备注 -->
                    <template v-else-if="field.key === 'audit_remark'">
                      <span :class="{ 'text-red': profile.data?.audit_status === 2 }">
                        {{ getFieldDisplayValue(profile.data, field.key) }}
                      </span>
                    </template>
                    <!-- 普通字段 -->
                    <template v-else>
                      {{ getFieldDisplayValue(profile.data, field.key) }}
                    </template>
                  </a-descriptions-item>
                </template>
              </a-descriptions>
            </template>
            <a-empty v-else description="请填写身份信息以完成认证" :image-style="{ height: '40px' }" />
            <div class="identity-actions" style="margin-top: 12px">
              <a-button 
                type="primary" 
                size="small" 
                :disabled="profile.data?.audit_status === 1" 
                @click="handleEditIdentity(profile)"
              >
                <EditOutlined /> {{ profile.has_profile ? '编辑信息' : '填写信息' }}
              </a-button>
            </div>
          </div>
        </a-card>
      </a-col>
      
      <!-- 右侧：编辑表单 -->
      <a-col :xs="24" :md="16">
        <a-card title="编辑资料" class="edit-card">
          <a-form
            :model="profileForm"
            :label-col="{ span: 4 }"
            :wrapper-col="{ span: 18 }"
            @finish="handleUpdateProfile"
          >
            <a-form-item label="昵称" name="nickname">
              <a-input v-model:value="profileForm.nickname" placeholder="请输入昵称" />
            </a-form-item>
            <a-form-item label="邮箱" name="email">
              <a-input v-model:value="profileForm.email" placeholder="请输入邮箱" />
            </a-form-item>
            <a-form-item label="手机" name="phone">
              <a-input v-model:value="profileForm.phone" placeholder="请输入手机号" />
            </a-form-item>
            <a-form-item :wrapper-col="{ offset: 4, span: 18 }">
              <a-button type="primary" html-type="submit" :loading="profileLoading">
                保存修改
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
        
        <a-card title="修改密码" class="edit-card" style="margin-top: 24px">
          <a-form
            :model="passwordForm"
            :label-col="{ span: 4 }"
            :wrapper-col="{ span: 18 }"
            @finish="handleChangePassword"
          >
            <a-form-item label="原密码" name="old_password" :rules="[{ required: true, message: '请输入原密码' }]">
              <a-input-password v-model:value="passwordForm.old_password" placeholder="请输入原密码" />
            </a-form-item>
            <a-form-item label="新密码" name="new_password" :rules="[{ required: true, message: '请输入新密码' }]">
              <a-input-password v-model:value="passwordForm.new_password" placeholder="请输入新密码（至少6位）" />
            </a-form-item>
            <a-form-item label="确认密码" name="confirm_password" :rules="[{ required: true, message: '请确认新密码' }]">
              <a-input-password v-model:value="passwordForm.confirm_password" placeholder="请再次输入新密码" />
            </a-form-item>
            <a-form-item :wrapper-col="{ offset: 4, span: 18 }">
              <a-button type="primary" html-type="submit" :loading="passwordLoading">
                修改密码
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
    
    <!-- 身份编辑抽屉（动态加载） -->
    <component
      v-if="editingProfile && identityFormComponent"
      :is="identityFormComponent"
      v-model:open="identityDrawerVisible"
      :record="editingProfile.data"
      :profileMode="true"
      @success="handleIdentitySaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, shallowRef, defineAsyncComponent, nextTick, h } from 'vue'
import { message, Modal } from 'ant-design-vue'
import * as Icons from '@ant-design/icons-vue'
import { MailOutlined, PhoneOutlined, CalendarOutlined, EditOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { getUserInfo } from '@/api/auth'
import { useUserStore } from '@/store/user'
import { updateProfile, changePassword, updateAvatar, getUserProfiles, type UserProfile } from '@/api/user'
import { formatTime } from '@/utils/format'
import AvatarUpload from '@/components/AvatarUpload.vue'
import type { FileInfo } from '@/types/file'

const userStore = useUserStore()
const profileLoading = ref(false)
const passwordLoading = ref(false)
const profilesLoading = ref(false)
const userProfiles = ref<UserProfile[]>([])

// 身份编辑相关
const identityDrawerVisible = ref(false)
const editingProfile = ref<UserProfile | null>(null)
const identityFormComponent = shallowRef<any>(null)

const profileForm = reactive({
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
  avatar_file_id: undefined as number | undefined
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const initProfileForm = () => {
  if (userStore.user) {
    profileForm.nickname = userStore.user.nickname || ''
    profileForm.email = userStore.user.email || ''
    profileForm.phone = userStore.user.phone || ''
    profileForm.avatar = userStore.user?.avatar_file_url || ''
    profileForm.avatar_file_id = userStore.user.avatar_file_id
  }
}

const handleAvatarSuccess = async (file: FileInfo) => {
  try {
    await updateAvatar(file.id)
    await userStore.getUserInfoAction()
    message.success('头像更新成功')
  } catch (error) {
    message.warning('头像已上传，请点击保存修改完成绑定')
  }
}

const handleUpdateProfile = async () => {
  profileLoading.value = true
  try {
    await updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone
    })
    message.success('更新成功')
    await userStore.getUserInfoAction()
  } finally {
    profileLoading.value = false
  }
}

const handleChangePassword = async () => {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    message.warning('两次密码不一致')
    return
  }
  if (passwordForm.new_password.length < 6) {
    message.warning('密码长度至少6位')
    return
  }
  
  passwordLoading.value = true
  try {
    await changePassword({ 
      old_password: passwordForm.old_password, 
      new_password: passwordForm.new_password 
    })
    message.success('修改成功，请重新登录')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } finally {
    passwordLoading.value = false
  }
}

// 获取用户身份信息
const fetchUserProfiles = async () => {
  profilesLoading.value = true
  try {
    const res = await getUserProfiles()
    userProfiles.value = res.data || []
  } catch (error) {
    console.error('加载身份信息失败', error)
  } finally {
    profilesLoading.value = false
  }
}

// 获取图标组件
const getIconComponent = (iconName: string) => {
  return (Icons as any)[iconName] || Icons.UserOutlined
}

// 获取字段显示值
const getFieldDisplayValue = (data: any, key: string) => {
  if (!data) return '-'
  const value = data[key]
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'boolean') return value ? '是' : '否'
  return value
}

// 蛇形命名转 PascalCase（farmer_certification -> FarmerCertification）
const toPascalCase = (str: string) => {
  return str.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join('')
}

// 编辑身份信息
const handleEditIdentity = async (profile: UserProfile) => {
  editingProfile.value = profile
  
  // 动态加载身份表单组件
  try {
    const formName = toPascalCase(profile.key)
    identityFormComponent.value = defineAsyncComponent(() => 
      import(`@/views/admin/system/${profile.key}/components/${formName}Form.vue`)
    )
    await nextTick()
    identityDrawerVisible.value = true
  } catch (error) {
    message.error('加载表单组件失败')
    console.error(error)
  }
}

// 身份信息保存成功
const handleIdentitySaved = () => {
  identityDrawerVisible.value = false
  fetchUserProfiles()
}

// 检查用户是否获得后台权限（定时轮询）
let permissionCheckTimer: ReturnType<typeof setInterval> | null = null
const permissionGranted = ref(false)

const checkPermissionChange = async () => {
  try {
    // 直接调用 API 获取最新用户信息（不走缓存）
    const res = await getUserInfo()
    const menus = res.data.menus || []
    
    // 如果有后台菜单权限，说明已被分配角色
    if (menus.length > 0 && !permissionGranted.value) {
      permissionGranted.value = true
      // 停止轮询
      if (permissionCheckTimer) {
        clearInterval(permissionCheckTimer)
        permissionCheckTimer = null
      }
      // 显示提示 - 角色变化后 Token 已失效，需要重新登录
      Modal.success({
        title: '认证已通过',
        icon: h(CheckCircleOutlined),
        content: '您的身份认证已通过审核，请重新登录以获取后台管理权限。',
        okText: '重新登录',
        async onOk() {
          // 登出并跳转到登录页
          await userStore.logoutAction()
          window.location.href = '/login'
        }
      })
    }
  } catch (error: any) {
    // 如果是 401 错误，说明 Token 已失效（角色变化导致）
    if (error?.message?.includes('401') || error?.message?.includes('Token')) {
      // 停止轮询
      if (permissionCheckTimer) {
        clearInterval(permissionCheckTimer)
        permissionCheckTimer = null
      }
      // 提示用户重新登录
      Modal.info({
        title: '权限已更新',
        content: '您的权限已更新，请重新登录以应用新权限。',
        okText: '重新登录',
        async onOk() {
          await userStore.logoutAction()
          window.location.href = '/login'
        }
      })
    }
  }
}

// 启动权限检查轮询（10秒一次）
const startPermissionCheck = () => {
  // 先检查一次
  checkPermissionChange()
  // 然后每 10 秒检查一次
  permissionCheckTimer = setInterval(checkPermissionChange, 10000)
}

onMounted(() => {
  initProfileForm()
  fetchUserProfiles()
  // 启动权限检查
  startPermissionCheck()
})

onUnmounted(() => {
  // 清理定时器
  if (permissionCheckTimer) {
    clearInterval(permissionCheckTimer)
    permissionCheckTimer = null
  }
})
</script>

<style scoped lang="less">
.front-profile-page {
  .profile-card {
    .avatar-section {
      text-align: center;
      padding: 24px 0;
      
      .nickname {
        margin: 16px 0 4px;
        font-size: 20px;
        font-weight: 600;
        color: #1a1a1a;
      }
      
      .username {
        color: #999;
        margin: 0;
      }
    }
    
    .info-list {
      .info-item {
        display: flex;
        align-items: center;
        padding: 12px 0;
        border-bottom: 1px solid #f0f0f0;
        
        &:last-child {
          border-bottom: none;
        }
        
        .info-icon {
          font-size: 16px;
          color: #1890ff;
          margin-right: 12px;
        }
        
        .info-label {
          color: #666;
          width: 60px;
        }
        
        .info-value {
          flex: 1;
          text-align: right;
          color: #1a1a1a;
        }
      }
    }
  }
  
  .edit-card {
    :deep(.ant-card-head-title) {
      font-weight: 600;
    }
  }
  
  .text-gray {
    color: #999;
  }
  
  .text-red {
    color: #ff4d4f;
  }
}

@media (max-width: 768px) {
  .front-profile-page {
    .ant-col {
      margin-bottom: 24px;
    }
  }
}
</style>
