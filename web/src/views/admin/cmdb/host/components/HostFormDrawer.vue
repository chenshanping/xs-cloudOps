<template>
  <a-drawer
    :open="open"
    :title="title"
    width="760"
    placement="right"
    destroy-on-close
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" layout="vertical">
      <div class="drawer-section">
        <div class="drawer-section__title">台账信息</div>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="主机名称" name="name">
              <a-input v-model:value="formState.name" placeholder="请输入主机名称" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="环境标识" name="environment">
              <a-input v-model:value="formState.environment" placeholder="如 prod / test / dev" allow-clear />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="主机分组" name="group_id">
              <a-select v-model:value="formState.group_id" placeholder="请选择主机分组" allow-clear>
                <a-select-option v-for="item in groupOptions" :key="item.id" :value="item.id">
                  {{ item.name }}<span v-if="item.status !== 1">（停用）</span>
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="负责人" name="owner">
              <a-input v-model:value="formState.owner" placeholder="请输入负责人" allow-clear />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="主机标签" name="tag_ids">
          <a-select v-model:value="formState.tag_ids" mode="multiple" placeholder="请选择主机标签" allow-clear>
            <a-select-option v-for="item in tagOptions" :key="item.id" :value="item.id">
              {{ item.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </div>

      <div class="drawer-section">
        <div class="drawer-section__title">SSH 连接</div>
        <a-alert
          type="info"
          show-icon
          style="margin-bottom: 16px"
          message="保存后会自动尝试 SSH 校验；即使校验失败，主机记录也会先保存，再由校验状态反馈问题。"
        />

        <a-row :gutter="16">
          <a-col :span="16">
            <a-form-item label="SSH 地址" name="ssh_host">
              <a-input v-model:value="formState.ssh_host" placeholder="请输入 SSH 地址或域名" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="SSH 端口" name="ssh_port">
              <a-input-number v-model:value="formState.ssh_port" style="width: 100%" :min="1" :max="65535" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="SSH 凭据" name="credential_id">
          <a-select v-model:value="formState.credential_id" placeholder="请选择 SSH 凭据" allow-clear>
            <a-select-option v-for="item in credentialOptions" :key="item.id" :value="item.id">
              {{ item.name }}（{{ item.username }} / {{ getCmdbAuthTypeLabel(item.auth_type) }}）
            </a-select-option>
          </a-select>
        </a-form-item>
      </div>

      <div class="drawer-section">
        <div class="drawer-section__title">备注</div>
        <a-form-item label="补充说明" name="remark">
          <a-textarea v-model:value="formState.remark" :rows="4" :maxlength="500" show-count placeholder="请输入备注" />
        </a-form-item>
      </div>
    </a-form>

    <template #footer>
      <a-space>
        <a-button :disabled="submitting" @click="handleClose">取消</a-button>
        <a-button type="primary" :loading="submitting" :disabled="submitting" @click="handleSubmit">保存</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import type { CmdbCredentialSummary, CmdbHostGroup, CmdbHostPayload, CmdbHostTag } from '@/api/cmdb'
import { getCmdbAuthTypeLabel } from '../../shared'

interface HostFormValue {
  name: string
  group_id?: number
  tag_ids: number[]
  environment: string
  owner: string
  private_ip: string
  public_ip: string
  ssh_host: string
  ssh_port: number
  credential_id?: number
  remark: string
}

const props = defineProps<{
  open: boolean
  title: string
  groupOptions: CmdbHostGroup[]
  tagOptions: CmdbHostTag[]
  credentialOptions: CmdbCredentialSummary[]
  initialValue?: Partial<HostFormValue>
  submitting?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: CmdbHostPayload): void
}>()

const formRef = ref<FormInstance>()
const formState = reactive<HostFormValue>({
  name: '',
  group_id: undefined,
  tag_ids: [],
  environment: '',
  owner: '',
  private_ip: '',
  public_ip: '',
  ssh_host: '',
  ssh_port: 22,
  credential_id: undefined,
  remark: '',
})

const validGroupIds = computed(() => new Set(props.groupOptions.map(item => item.id)))
const validCredentialIds = computed(() => new Set(props.credentialOptions.map(item => item.id)))
const validTagIds = computed(() => new Set(props.tagOptions.map(item => item.id)))

const formRules = computed<Record<string, Rule[]>>(() => ({
  name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
  group_id: [
    { required: true, message: '请选择主机分组', trigger: 'change' },
    {
      trigger: 'change',
      validator: async (_rule, value?: number) => {
        if (!value || !validGroupIds.value.has(value)) {
          throw new Error('请选择有效的主机分组')
        }
      }
    }
  ],
  credential_id: [
    { required: true, message: '请选择 SSH 凭据', trigger: 'change' },
    {
      trigger: 'change',
      validator: async (_rule, value?: number) => {
        if (!value || !validCredentialIds.value.has(value)) {
          throw new Error('请选择有效的 SSH 凭据')
        }
      }
    }
  ],
  ssh_host: [{ required: true, message: '请输入 SSH 地址', trigger: 'blur' }],
  ssh_port: [{ type: 'number', min: 1, max: 65535, message: 'SSH 端口范围为 1-65535', trigger: 'change' }],
  tag_ids: [
    {
      trigger: 'change',
      validator: async (_rule, value?: number[]) => {
        if (!value?.length) {
          return
        }
        if (value.some(item => !validTagIds.value.has(item))) {
          throw new Error('标签选项已失效，请重新选择')
        }
      }
    }
  ],
}))

watch(
  () => [props.open, props.initialValue] as const,
  ([open]) => {
    if (!open) {
      return
    }
    Object.assign(formState, {
      name: props.initialValue?.name ?? '',
      group_id: props.initialValue?.group_id,
      tag_ids: props.initialValue?.tag_ids ? [...props.initialValue.tag_ids] : [],
      environment: props.initialValue?.environment ?? '',
      owner: props.initialValue?.owner ?? '',
      private_ip: props.initialValue?.private_ip ?? '',
      public_ip: props.initialValue?.public_ip ?? '',
      ssh_host: props.initialValue?.ssh_host ?? '',
      ssh_port: props.initialValue?.ssh_port ?? 22,
      credential_id: props.initialValue?.credential_id,
      remark: props.initialValue?.remark ?? '',
    })
  },
  { immediate: true, deep: true }
)

const handleClose = () => {
  if (props.submitting) {
    return
  }
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
    group_id: Number(formState.group_id),
    tag_ids: [...new Set(formState.tag_ids)],
    environment: formState.environment.trim(),
    owner: formState.owner.trim(),
    private_ip: formState.private_ip.trim(),
    public_ip: formState.public_ip.trim(),
    ssh_host: formState.ssh_host.trim(),
    ssh_port: formState.ssh_port || 22,
    credential_id: Number(formState.credential_id),
    remark: formState.remark.trim(),
  })
}
</script>

<style scoped>
.drawer-section + .drawer-section {
  margin-top: 16px;
}

.drawer-section {
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.drawer-section__title {
  margin-bottom: 14px;
  color: var(--app-text-strong);
  font-size: 14px;
  font-weight: 600;
}
</style>
