<template>
  <a-drawer
    :open="open"
    :title="`用户身份 - ${user?.username || ''}`"
    width="50%"
    placement="right"
    @close="emit('update:open', false)"
  >
    <a-spin :spinning="loading">
      <a-empty v-if="!loading && profiles.length === 0" description="未绑定任何身份" />
      <a-collapse v-else v-model:activeKey="activeKey" accordion>
        <a-collapse-panel v-for="profile in profiles" :key="profile.key" :header="profile.name">
          <template #extra>
            <a-tag v-if="profile.has_profile" :color="profile.is_complete ? 'success' : 'warning'">
              {{ profile.is_complete ? '已完善' : '未完善' }}
            </a-tag>
            <a-tag v-else color="default">未填写</a-tag>
          </template>
          <a-descriptions v-if="profile.has_profile && profile.data" :column="2" size="small" bordered>
            <a-descriptions-item v-for="field in profile.fields" :key="field.key" :label="field.label">
              <template v-if="field.type === 'image'">
                <a-image v-if="getFieldValue(profile.data, field.key)" :src="getFieldValue(profile.data, field.key)" :width="80" />
                <span v-else>-</span>
              </template>
              <template v-else-if="field.type === 'file'">
                <a v-if="getFieldValue(profile.data, field.key)" :href="getFieldValue(profile.data, field.key)" target="_blank">查看文件</a>
                <span v-else>-</span>
              </template>
              <template v-else-if="field.type === 'images'">
                <a-image-preview-group v-if="getFieldValue(profile.data, field.key)">
                  <a-space>
                    <a-image
                      v-for="(url, idx) in getFieldValue(profile.data, field.key).split(',')"
                      :key="idx"
                      :src="url"
                      :width="60"
                    />
                  </a-space>
                </a-image-preview-group>
                <span v-else>-</span>
              </template>
              <template v-else-if="field.key === 'audit_status'">
                <a-tag v-if="getFieldValue(profile.data, field.key) === 0" color="default">待审批</a-tag>
                <a-tag v-else-if="getFieldValue(profile.data, field.key) === 1" color="success">审批通过</a-tag>
                <a-tag v-else-if="getFieldValue(profile.data, field.key) === 2" color="error">审批拒绝</a-tag>
                <span v-else>-</span>
              </template>
              <template v-else>
                {{ getFieldValue(profile.data, field.key) || '-' }}
              </template>
            </a-descriptions-item>
          </a-descriptions>
          <a-empty v-else description="未填写档案信息" />
        </a-collapse-panel>
      </a-collapse>
    </a-spin>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { User } from '@/types'
import type { UserProfile } from '@/api/user'

const props = defineProps<{
  open: boolean
  loading: boolean
  user: User | null
  profiles: UserProfile[]
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const activeKey = ref<string>('')

watch(
  () => [props.open, props.profiles],
  () => {
    if (!props.open) {
      activeKey.value = ''
      return
    }
    activeKey.value = props.profiles[0]?.key || ''
  },
  { immediate: true, deep: true }
)

const getFieldValue = (data: Record<string, any> | undefined, key: string) => {
  if (!data) {
    return ''
  }
  return data[key]
}
</script>
