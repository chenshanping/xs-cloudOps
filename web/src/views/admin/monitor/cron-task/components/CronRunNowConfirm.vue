<template></template>

<script setup lang="ts">
import { Modal, message } from 'ant-design-vue'
import { runCronTaskNow, type CronTask } from '@/api/cron'

const emit = defineEmits<{
  (e: 'success', logId: number, task: CronTask): void
}>()

const confirmRun = (task: CronTask) => {
  Modal.confirm({
    title: '立即执行定时任务',
    content: `确认立即执行「${task.name}」吗？任务将在后台异步执行。`,
    okText: '立即执行',
    cancelText: '取消',
    async onOk() {
      const res = await runCronTaskNow(task.id)
      message.success('任务已提交执行')
      emit('success', res.data.log_id, task)
    },
  })
}

defineExpose({ confirmRun })
</script>
