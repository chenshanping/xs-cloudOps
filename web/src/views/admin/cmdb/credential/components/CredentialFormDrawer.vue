<template>
  <a-drawer
    :open="open"
    :title="title"
    width="620"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <template v-if="mode === 'view'">
      <div class="detail-layout">
        <div class="detail-hero">
          <div>
            <div class="detail-hero__eyebrow">SSH 凭据</div>
            <div class="detail-hero__title">{{ detail?.name || '-' }}</div>
            <div class="detail-hero__desc">凭据明文字段不会在查看场景返回，编辑时需重新填写认证材料。</div>
          </div>
          <a-space wrap :size="[8, 8]">
            <a-tag color="blue">{{ getCmdbAuthTypeLabel(detail?.auth_type) }}</a-tag>
            <a-tag>{{ detail?.username || '-' }}</a-tag>
            <a-tag color="processing">绑定主机 {{ detail?.bind_count ?? 0 }}</a-tag>
          </a-space>
        </div>

        <div class="detail-section">
          <div class="detail-section__title">凭据摘要</div>
          <a-descriptions :column="2" bordered size="small">
            <a-descriptions-item label="凭据名称">{{ detail?.name || '-' }}</a-descriptions-item>
            <a-descriptions-item label="认证方式">{{ getCmdbAuthTypeLabel(detail?.auth_type) }}</a-descriptions-item>
            <a-descriptions-item label="登录用户名">{{ detail?.username || '-' }}</a-descriptions-item>
            <a-descriptions-item label="绑定主机数">{{ detail?.bind_count ?? 0 }}</a-descriptions-item>
            <a-descriptions-item label="备注" :span="2">{{ detail?.remark || '-' }}</a-descriptions-item>
          </a-descriptions>
        </div>
      </div>
    </template>

    <a-form v-else ref="formRef" :model="formState" :rules="formRules" layout="vertical">
      <div class="drawer-section">
        <div class="drawer-section__title">基础信息</div>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="凭据名称" name="name">
              <a-input v-model:value="formState.name" placeholder="请输入凭据名称" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="认证方式" name="auth_type">
              <a-segmented v-model:value="formState.auth_type" :options="authTypeSegmentOptions" block />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="登录用户名" name="username">
          <a-input v-model:value="formState.username" placeholder="请输入登录用户名" allow-clear />
        </a-form-item>

        <a-form-item v-if="formState.auth_type === 'password'" label="登录密码" name="password">
          <a-input-password v-model:value="formState.password" placeholder="请输入登录密码" />
        </a-form-item>

        <template v-else>
          <a-form-item label="私钥内容" name="private_key">
            <a-textarea v-model:value="formState.private_key" :rows="8" placeholder="请输入私钥内容" />
          </a-form-item>
          <a-form-item label="私钥口令" name="passphrase">
            <a-input-password v-model:value="formState.passphrase" placeholder="如私钥未加密可留空" />
          </a-form-item>
        </template>
      </div>

      <div class="drawer-section">
        <div class="drawer-section__title">使用提示</div>
        <div class="tip-panel">
          <div class="tip-panel__title">编辑说明</div>
          <div class="tip-panel__content">为避免敏感信息回显，编辑已有凭据时需要重新填写密码或私钥内容。</div>
        </div>
        <a-form-item label="备注" name="remark" style="margin-top: 16px">
          <a-textarea v-model:value="formState.remark" :rows="4" :maxlength="500" show-count placeholder="请输入备注" />
        </a-form-item>
      </div>
    </a-form>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">{{ mode === 'view' ? '关闭' : '取消' }}</a-button>
        <a-button v-if="mode !== 'view'" type="primary" @click="handleSubmit">保存</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import type { CmdbCredentialPayload, CmdbCredentialSummary } from '@/api/cmdb'
import { cmdbAuthTypeOptions, getCmdbAuthTypeLabel } from '../../shared'

type DrawerMode = 'create' | 'edit' | 'view'

interface CredentialFormValue {
  name: string
  auth_type: 'password' | 'private_key'
  username: string
  password: string
  private_key: string
  passphrase: string
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  mode: DrawerMode
  initialValue?: Partial<CredentialFormValue>
  detail?: CmdbCredentialSummary | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: CmdbCredentialPayload): void
}>()

const formRef = ref<FormInstance>()
const formState = reactive<CredentialFormValue>({
  name: '',
  auth_type: 'password',
  username: '',
  password: '',
  private_key: '',
  passphrase: '',
  remark: '',
})

const authTypeSegmentOptions = cmdbAuthTypeOptions.map(item => ({ label: item.label, value: item.value }))

const formRules = computed<Record<string, Rule[]>>(() => ({
  name: [{ required: true, message: '请输入凭据名称', trigger: 'blur' }],
  username: [{ required: true, message: '请输入登录用户名', trigger: 'blur' }],
  password: formState.auth_type === 'password'
    ? [{ required: true, message: '请输入登录密码', trigger: 'blur' }]
    : [],
  private_key: formState.auth_type === 'private_key'
    ? [{ required: true, message: '请输入私钥内容', trigger: 'blur' }]
    : [],
}))

watch(
  () => [props.open, props.initialValue, props.mode] as const,
  ([open]) => {
    if (!open || props.mode === 'view') {
      return
    }
    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      auth_type: props.initialValue?.auth_type ?? 'password',
      username: props.initialValue?.username ?? '',
      password: '',
      private_key: '',
      passphrase: '',
      remark: props.initialValue?.remark ?? '',
    })
  },
  { immediate: true, deep: true }
)

watch(
  () => formState.auth_type,
  (value) => {
    if (value === 'password') {
      formState.private_key = ''
      formState.passphrase = ''
    } else {
      formState.password = ''
    }
  }
)

const handleClose = () => {
  emit('update:open', false)
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  emit('submit', {
    name: formState.name.trim(),
    auth_type: formState.auth_type,
    username: formState.username.trim(),
    password: formState.password,
    private_key: formState.private_key,
    passphrase: formState.passphrase,
    remark: formState.remark.trim(),
  })
}
</script>

<style scoped>
.drawer-section + .drawer-section {
  margin-top: 16px;
}

.drawer-section,
.detail-section {
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.drawer-section__title,
.detail-section__title {
  margin-bottom: 14px;
  color: var(--app-text-strong);
  font-size: 14px;
  font-weight: 600;
}

.tip-panel {
  padding: 14px;
  background: rgba(22, 119, 255, 0.06);
  border: 1px solid rgba(22, 119, 255, 0.14);
  border-radius: 8px;
}

.tip-panel__title {
  color: #1677ff;
  font-size: 13px;
  font-weight: 600;
}

.tip-panel__content {
  margin-top: 6px;
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.6;
}

.detail-layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(22, 119, 255, 0.08), rgba(22, 119, 255, 0.02));
  border: 1px solid rgba(22, 119, 255, 0.16);
  border-radius: 10px;
}

.detail-hero__eyebrow {
  color: #1677ff;
  font-size: 12px;
  font-weight: 600;
}

.detail-hero__title {
  margin-top: 6px;
  color: var(--app-text-strong);
  font-size: 18px;
  font-weight: 700;
}

.detail-hero__desc {
  margin-top: 8px;
  color: var(--app-text-muted);
  font-size: 13px;
  line-height: 1.6;
}

@media (max-width: 900px) {
  .detail-hero {
    flex-direction: column;
  }
}
</style>
