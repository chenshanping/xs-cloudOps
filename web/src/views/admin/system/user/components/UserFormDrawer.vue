<template>
  <a-drawer
    :open="open"
    :title="title"
    width="560"
    placement="right"
    @close="handleClose"
  >
    <a-form ref="formRef" :model="formState" :rules="formRules" :label-col="{ span: 5 }">
      <a-alert
        v-if="isLegacyDeptBinding"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
        :message="legacyDeptMessage"
      />
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
      <a-form-item v-if="!isEdit" label="密码" name="password">
        <a-input-password v-model:value="formState.password" />
      </a-form-item>
      <a-form-item label="昵称">
        <a-input v-model:value="formState.nickname" />
      </a-form-item>
      <a-form-item label="性别">
        <a-select v-model:value="formState.gender" placeholder="请选择性别">
          <a-select-option v-for="option in genderOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="邮箱" name="email">
        <a-input v-model:value="formState.email" />
      </a-form-item>
      <a-form-item label="手机号" name="phone">
        <a-input v-model:value="formState.phone" />
      </a-form-item>
      <a-form-item label="所属部门" name="deptSelection">
        <a-tree-select
          v-model:value="formState.deptSelection"
          :tree-data="deptOptions"
          :field-names="{ label: 'title', value: 'value', children: 'children' }"
          placeholder="请选择所属部门"
          tree-default-expand-all
          tree-node-filter-prop="title"
          label-in-value
          show-search
        />
      </a-form-item>
      <a-form-item label="角色">
        <a-select v-model:value="formState.role_ids" mode="multiple" placeholder="请选择角色">
          <a-select-option v-for="role in roleOptions" :key="role.id" :value="role.id">
            {{ role.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="状态">
        <a-switch v-model:checked="formState.statusChecked" />
      </a-form-item>
    </a-form>

    <template #footer>
      <a-space>
        <a-button @click="handleClose">取消</a-button>
        <a-button type="primary" @click="handleSubmit">保存</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'
import { message } from 'ant-design-vue'
import AvatarUpload from '@/components/AvatarUpload.vue'
import type { Role } from '@/types'
import type { GenderOption } from '../user-gender'

interface TreeSelectOption {
  key: string | number
  title: string
  value: number
  disabled?: boolean
  selectable?: boolean
  isLeaf?: boolean
  children?: TreeSelectOption[]
}

interface TreeSelectValue {
  value: number
  label?: string
}

interface UserFormValue {
  username: string
  password: string
  nickname: string
  gender: number
  email: string
  phone: string
  deptSelection?: TreeSelectValue
  dept_id?: number
  dept_label?: string
  role_ids: number[]
  statusChecked: boolean
  avatar_file_id?: number
  avatar_file_url: string
}

const props = defineProps<{
  open: boolean
  title: string
  isEdit: boolean
  roleOptions: Role[]
  genderOptions: GenderOption[]
  deptOptions: TreeSelectOption[]
  initialValue?: Partial<UserFormValue>
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', value: {
    username: string
    password: string
    nickname: string
    gender: number
    email: string
    phone: string
    dept_id?: number
    role_ids: number[]
    status: number
    avatar_file_id?: number
  }): void
}>()

const formRef = ref<FormInstance>()
const mainlandPhonePattern = /^1[3-9]\d{9}$/
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const formState = reactive<UserFormValue>({
  username: '',
  password: '',
  nickname: '',
  gender: 0,
  email: '',
  phone: '',
  deptSelection: undefined,
  dept_id: undefined,
  role_ids: [2],
  statusChecked: true,
  avatar_file_id: undefined,
  avatar_file_url: ''
})

const initialDeptOption = computed(() => findDeptOption(props.deptOptions, props.initialValue?.dept_id))

const isLegacyDeptBinding = computed(
  () => props.isEdit && !!props.initialValue?.dept_id && initialDeptOption.value?.isLeaf === false
)

const legacyDeptMessage = computed(() => {
  const deptName = props.initialValue?.dept_label || initialDeptOption.value?.title || '当前部门'
  return `该用户当前属于历史绑定数据，${deptName} 已存在下级部门。可先保存非部门字段；如需调整部门，请改到具体叶子部门。`
})

const findDeptOption = (options: TreeSelectOption[], targetValue?: number): TreeSelectOption | undefined => {
  if (!targetValue) {
    return undefined
  }
  for (const option of options) {
    if (option.value === targetValue) {
      return option
    }
    const childOption = option.children ? findDeptOption(option.children, targetValue) : undefined
    if (childOption) {
      return childOption
    }
  }
  return undefined
}

const canKeepLegacyDeptSelection = (value?: TreeSelectValue) =>
  isLegacyDeptBinding.value && value?.value === props.initialValue?.dept_id

const normalizeOptionalText = (value?: string) => value?.trim() ?? ''

const validateOptionalEmail = async (_rule: Rule, value?: string) => {
  const normalized = normalizeOptionalText(value)
  if (!normalized) {
    return
  }
  if (!emailPattern.test(normalized)) {
    throw new Error('请输入正确的邮箱格式')
  }
}

const validateOptionalPhone = async (_rule: Rule, value?: string) => {
  const normalized = normalizeOptionalText(value)
  if (!normalized) {
    return
  }
  if (!mainlandPhonePattern.test(normalized)) {
    throw new Error('请输入正确的手机号格式')
  }
}

const formRules = computed<Record<string, Rule[]>>(() => ({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: props.isEdit ? [] : [{ required: true, message: '请输入密码', trigger: 'blur' }],
  avatar_file_id: [
    {
      trigger: 'change',
      validator: async () => {
        if (!formState.avatar_file_id && !normalizeOptionalText(formState.avatar_file_url)) {
          throw new Error('请上传头像')
        }
      }
    }
  ],
  email: [{ trigger: 'blur', validator: validateOptionalEmail }],
  phone: [{ trigger: 'blur', validator: validateOptionalPhone }],
  deptSelection: [
    { required: true, message: '请选择所属部门', trigger: 'change' },
    {
      trigger: 'change',
      validator: async (_rule, value?: TreeSelectValue) => {
        if (canKeepLegacyDeptSelection(value)) {
          return
        }
        const selected = findDeptOption(props.deptOptions, value?.value)
        if (!selected || selected.selectable !== true || selected.disabled) {
          throw new Error('只能绑定叶子部门')
        }
      }
    }
  ]
}))

const findDeptOptionLabel = (options: TreeSelectOption[], targetValue?: number): string | undefined => {
  if (!targetValue) {
    return undefined
  }
  for (const option of options) {
    if (option.value === targetValue) {
      return option.title
    }
    const childLabel = option.children ? findDeptOptionLabel(option.children, targetValue) : undefined
    if (childLabel) {
      return childLabel
    }
  }
  return undefined
}

watch(
  () => [props.open, props.initialValue, props.isEdit, props.deptOptions],
  () => {
    if (!props.open) {
      return
    }
    const deptId = props.initialValue?.dept_id
    const deptLabel = props.initialValue?.dept_label || findDeptOptionLabel(props.deptOptions, deptId)
    Object.assign(formState, {
      username: props.initialValue?.username ?? '',
      password: props.initialValue?.password ?? '',
      nickname: props.initialValue?.nickname ?? '',
      gender: props.initialValue?.gender ?? 0,
      email: props.initialValue?.email ?? '',
      phone: props.initialValue?.phone ?? '',
      deptSelection: deptId ? { value: deptId, label: deptLabel ?? String(deptId) } : undefined,
      dept_id: deptId,
      dept_label: deptLabel,
      role_ids: props.initialValue?.role_ids ? [...props.initialValue.role_ids] : [2],
      statusChecked: props.initialValue?.statusChecked ?? true,
      avatar_file_id: props.initialValue?.avatar_file_id,
      avatar_file_url: props.initialValue?.avatar_file_url ?? ''
    })
  },
  { immediate: true, deep: true }
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

  const selectedDept = findDeptOption(props.deptOptions, formState.deptSelection?.value)
  if (!canKeepLegacyDeptSelection(formState.deptSelection) && (!selectedDept || selectedDept.selectable !== true || selectedDept.disabled)) {
    message.warning('只能绑定叶子部门')
    return
  }

  emit('submit', {
    username: formState.username,
    password: formState.password,
    nickname: formState.nickname,
    gender: formState.gender,
    email: normalizeOptionalText(formState.email),
    phone: normalizeOptionalText(formState.phone),
    dept_id: formState.deptSelection?.value,
    role_ids: [...formState.role_ids],
    status: formState.statusChecked ? 1 : 0,
    avatar_file_id: formState.avatar_file_id
  })
}
</script>
