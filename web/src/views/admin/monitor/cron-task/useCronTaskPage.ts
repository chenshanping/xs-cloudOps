import { computed, reactive, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  createCronTask,
  deleteCronTask,
  disableCronTask,
  enableCronTask,
  getCronRegistry,
  getCronTaskList,
  updateCronTask,
  type CronTask,
  type CronTaskPayload,
  type RegisteredCronTask,
} from '@/api/cron'

export function useCronTaskPage() {
  const loading = ref(false)
  const submitting = ref(false)
  const tableData = ref<CronTask[]>([])
  const registry = ref<RegisteredCronTask[]>([])
  const drawerVisible = ref(false)
  const isEdit = ref(false)
  const currentRecord = ref<CronTask | null>(null)
  const pagination = reactive({ current: 1, pageSize: 10, total: 0 })
  const searchForm = reactive({
    code: '',
    task_code: undefined as string | undefined,
    name: '',
    status: undefined as number | undefined,
  })

  const enabledCount = computed(() => tableData.value.filter((item) => item.status === 1).length)
  const failedCount = computed(() => tableData.value.filter((item) => item.last_status === 'failure').length)
  const runningCount = computed(() => tableData.value.filter((item) => item.last_status === 'running').length)

  const fetchRegistry = async () => {
    const res = await getCronRegistry()
    registry.value = res.data || []
  }

  const fetchData = async () => {
    loading.value = true
    try {
      const res = await getCronTaskList({
        page: pagination.current,
        page_size: pagination.pageSize,
        code: searchForm.code || undefined,
        task_code: searchForm.task_code || undefined,
        name: searchForm.name || undefined,
        status: searchForm.status,
      })
      tableData.value = res.data.list || []
      pagination.total = res.data.total
    } finally {
      loading.value = false
    }
  }

  const handleSearch = () => {
    pagination.current = 1
    fetchData()
  }

  const handleReset = () => {
    searchForm.code = ''
    searchForm.task_code = undefined
    searchForm.name = ''
    searchForm.status = undefined
    handleSearch()
  }

  const handleTableChange = (pag: any) => {
    pagination.current = pag.current
    pagination.pageSize = pag.pageSize
    fetchData()
  }

  const handleAdd = () => {
    isEdit.value = false
    currentRecord.value = null
    drawerVisible.value = true
  }

  const handleEdit = (record: CronTask) => {
    isEdit.value = true
    currentRecord.value = record
    drawerVisible.value = true
  }

  const handleSubmit = async (payload: CronTaskPayload) => {
    submitting.value = true
    try {
      if (isEdit.value && currentRecord.value) {
        await updateCronTask(currentRecord.value.id, payload)
        message.success('更新成功')
      } else {
        await createCronTask(payload)
        message.success('创建成功')
      }
      drawerVisible.value = false
      fetchData()
    } finally {
      submitting.value = false
    }
  }

  const handleDelete = (record: CronTask) => {
    Modal.confirm({
      title: '删除定时任务',
      content: `确认删除「${record.name}」吗？`,
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      async onOk() {
        await deleteCronTask(record.id)
        message.success('删除成功')
        fetchData()
      },
    })
  }

  const handleEnable = async (record: CronTask) => {
    await enableCronTask(record.id)
    message.success('启用成功')
    fetchData()
  }

  const handleDisable = async (record: CronTask) => {
    await disableCronTask(record.id)
    message.success('停用成功')
    fetchData()
  }

  return {
    loading,
    submitting,
    tableData,
    registry,
    drawerVisible,
    isEdit,
    currentRecord,
    pagination,
    searchForm,
    enabledCount,
    failedCount,
    runningCount,
    fetchRegistry,
    fetchData,
    handleSearch,
    handleReset,
    handleTableChange,
    handleAdd,
    handleEdit,
    handleSubmit,
    handleDelete,
    handleEnable,
    handleDisable,
  }
}
