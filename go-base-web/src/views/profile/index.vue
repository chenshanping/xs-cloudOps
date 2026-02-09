<template>
  <div class="profile-page">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card title="个人信息">
          <div class="avatar-section">
            <AvatarUpload
              v-model:fileId="profileForm.avatar_file_id"
              v-model:url="profileForm.avatar"
              :size="100"
              tip=""
              @success="handleAvatarSuccess"
            />
            <h3>{{ userStore.user?.nickname || userStore.user?.username }}</h3>
            <p class="username">@{{ userStore.user?.username }}</p>
          </div>
          <a-descriptions :column="1" size="small">
            <a-descriptions-item label="邮箱">{{ userStore.user?.email || '-' }}</a-descriptions-item>
            <a-descriptions-item label="手机">{{ userStore.user?.phone || '-' }}</a-descriptions-item>
            <a-descriptions-item label="角色">
              <a-tag v-for="role in userStore.user?.roles" :key="role.id" color="blue">{{ role.name }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="创建时间">{{ formatTime(userStore.user?.created_at) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
        <!-- 身份信息卡片 -->
       <div  v-if="userStore.user?.roles[0].name!='超级管理员' && userStore.user?.roles[0].name!='系统管理员'">
         <a-card  v-for="profile in userProfiles" :key="profile.key" class="identity-card" :loading="profilesLoading">
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
              <a-tag v-if="profile.is_complete" color="blue">已完善</a-tag>
              <a-tag v-else color="default">未完善</a-tag>
            </a-space>
          </template>
          <div class="identity-content">
            <a-descriptions :column="1" size="small" :labelStyle="{ width: '200px' }">
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
            <div class="identity-actions">
              <a-button type="primary" size="small" :disabled="profile.data?.audit_status==1" @click="handleEditIdentity(profile)">
                <EditOutlined /> 编辑信息 
              </a-button>
            </div>
          </div>
        </a-card>
       </div>
      </a-col>
      <a-col :span="16">
        <a-card title="编辑资料" style="margin-bottom: 16px">
          <a-form :model="profileForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 16 }">
            <a-form-item label="昵称"><a-input v-model:value="profileForm.nickname" /></a-form-item>
            <a-form-item label="邮箱"><a-input v-model:value="profileForm.email" /></a-form-item>
            <a-form-item label="手机"><a-input v-model:value="profileForm.phone" /></a-form-item>
            <a-form-item :wrapper-col="{ offset: 4 }">
              <a-button type="primary" :loading="profileLoading" @click="handleUpdateProfile">保存修改</a-button>
            </a-form-item>
          </a-form>
        </a-card>
        <a-card title="修改密码">
          <a-form :model="passwordForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 16 }">
            <a-form-item label="原密码"><a-input-password v-model:value="passwordForm.old_password" /></a-form-item>
            <a-form-item label="新密码"><a-input-password v-model:value="passwordForm.new_password" /></a-form-item>
            <a-form-item label="确认密码"><a-input-password v-model:value="passwordForm.confirm_password" /></a-form-item>
            <a-form-item :wrapper-col="{ offset: 4 }">
              <a-button type="primary" :loading="passwordLoading" @click="handleChangePassword">修改密码</a-button>
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
import { ref, reactive, onMounted, shallowRef, defineAsyncComponent, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import * as Icons from '@ant-design/icons-vue'
import { EditOutlined } from '@ant-design/icons-vue'
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

// 获取图标组件
const getIconComponent = (iconName: string) => {
  return (Icons as any)[iconName] || Icons.UserOutlined
}

// 获取字段显示值
const getFieldDisplayValue = (data: any, key: string) => {
  if (!data) return '-'
  const value = data[key]
  if (value === null || value === undefined || value === '') return '-'
  // 布尔值转中文
  if (typeof value === 'boolean') return value ? '是' : '否'
  // 数字 0 也显示
  return value
}


// 加载用户身份
const loadUserProfiles = async () => {
  profilesLoading.value = true
  try {
    const res = await getUserProfiles()
    userProfiles.value = res.data || []
    
  } catch (error) {
    console.error('获取用户身份失败', error)
  } finally {
    profilesLoading.value = false
  }
}

// 蛇形命名转 PascalCase（farmer_certification -> FarmerCertification）
const toPascalCase = (str: string) => {
  return str.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join('')
}

// 编辑身份
const handleEditIdentity = async (profile: UserProfile) => {
  editingProfile.value = profile
  // 动态加载对应的表单组件
  try {
    const moduleName = profile.key
    const formName = toPascalCase(moduleName)
    identityFormComponent.value = defineAsyncComponent(
      () => import(`@/views/system/${moduleName}/components/${formName}Form.vue`)
    )
    // 等待组件加载后打开抽屉
    await nextTick()
    identityDrawerVisible.value = true
  } catch (e) {
    message.error('加载表单组件失败')
    identityFormComponent.value = null
  }
}

// 身份信息保存成功回调
const handleIdentitySaved = () => {
  loadUserProfiles()
}

// 头像上传成功后自动保存
const handleAvatarSuccess = async (file: FileInfo) => {
  try {
    await updateAvatar(file.id)
    await userStore.getUserInfoAction()
    message.success('头像更新成功')
  } catch (error) {
    // 上传成功但绑定失败，提示用户
    message.warning('头像已上传，请点击保存修改完成绑定')
  }
}

const handleUpdateProfile = async () => {
  profileLoading.value = true
  try {
    await updateProfile(profileForm)
    message.success('更新成功')
    await userStore.getUserInfoAction()
  } finally {
    profileLoading.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordForm.old_password || !passwordForm.new_password) {
    message.warning('请填写完整')
    return
  }
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
    await changePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
    message.success('修改成功，请重新登录')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  initProfileForm()
  loadUserProfiles()
})
</script>

<style scoped>
.avatar-section {
  text-align: center;
  padding: 20px 0;
}
.avatar-section h3 {
  margin: 16px 0 4px;
}
.avatar-section .username {
  color: #999;
  margin: 0;
}
.identity-card {
  margin-top: 16px;
}
.identity-card :deep(.ant-card-body) {
  padding: 16px;
}
.identity-content {
  position: relative;
}
.identity-actions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}
.text-gray {
  color: #999;
}
.text-red {
  color: #ff4d4f;
}
</style>
