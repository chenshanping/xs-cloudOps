import { computed, onMounted, reactive, ref, watch } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  createDictData,
  createDictType,
  deleteDictData,
  deleteDictType,
  getDictDataList,
  getDictTypeList,
  updateDictData,
  updateDictType,
  type DictData,
  type DictType,
} from '@/api/dict'
import {
  canApplyDictDataResponse,
  filterDictTypes,
  reconcileSelectedType,
} from './dict-page-state'

export type StatusFilter = 'all' | 0 | 1

interface DictTypeFormValue {
  name: string
  type: string
  status: number
  remark: string
}

interface DictDataFormValue {
  label: string
  value: string
  sort: number
  tag_type: string
  status: number
  is_default: number
  remark: string
}

export function useDictPage() {
  const dictTypes = ref<DictType[]>([])
  const typeLoading = ref(false)
  const typeSearchText = ref('')
  const selectedType = ref<DictType | null>(null)
  const typePagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
    showSizeChanger: true,
    showTotal: (total: number) => `共 ${total} 条`,
  })

  const typeStatusFilter = ref<StatusFilter>('all')

  const dictDataList = ref<DictData[]>([])
  const dataLoading = ref(false)
  const dataPagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
    showSizeChanger: true,
    showTotal: (total: number) => `共 ${total} 条`,
  })
  const dataLabelSearch = ref('')
  const dataStatusFilter = ref<StatusFilter>('all')
  const selectedDataKeys = ref<number[]>([])
  const togglingStatusIds = ref<Set<number>>(new Set())
  const togglingTypeIds = ref<Set<number>>(new Set())
  const batchDeleting = ref(false)

  const typeDrawerVisible = ref(false)
  const typeDrawerTitle = ref('新增字典类型')
  const typeSubmitLoading = ref(false)
  const editingType = ref<DictType | null>(null)
  const typeDrawerInitialValue = ref<Partial<DictTypeFormValue>>({})

  const dataDrawerVisible = ref(false)
  const dataDrawerTitle = ref('新增字典数据')
  const dataSubmitLoading = ref(false)
  const editingData = ref<DictData | null>(null)
  const dataDrawerInitialValue = ref<Partial<DictDataFormValue>>({})

  const filteredDictTypes = computed(() => filterDictTypes(dictTypes.value, typeSearchText.value))
  const showSelectedOutsideFilter = computed(() => {
    if (!selectedType.value) {
      return false
    }
    return filteredDictTypes.value.every(item => item.id !== selectedType.value?.id)
  })

  let latestDataRequestId = 0

  const fetchDictTypes = async (resetPage = false) => {
    if (resetPage) {
      typePagination.current = 1
    }

    typeLoading.value = true
    try {
      const keyword = typeSearchText.value.trim()
      const res = await getDictTypeList({
        page: typePagination.current,
        page_size: typePagination.pageSize,
        name: keyword || undefined,
        type: keyword || undefined,
        status: typeStatusFilter.value === 'all' ? undefined : typeStatusFilter.value,
      })

      dictTypes.value = res.data.list || []
      typePagination.total = res.data.total
      const previousSelectedType = selectedType.value
      const reconciledSelectedType = reconcileSelectedType(dictTypes.value, previousSelectedType)

      if (!reconciledSelectedType && dictTypes.value.length > 0) {
        selectedType.value = dictTypes.value[0]
      } else {
        selectedType.value = reconciledSelectedType
      }

      if (selectedType.value?.type) {
        if (
          previousSelectedType?.id !== selectedType.value.id
          || dictDataList.value.length === 0
        ) {
          dataPagination.current = 1
          await fetchDictData(selectedType.value.type)
        }
      } else {
        dictDataList.value = []
        dataPagination.total = 0
      }
    } finally {
      typeLoading.value = false
    }
  }

  const fetchDictData = async (typeCode = selectedType.value?.type) => {
    if (!typeCode) {
      dictDataList.value = []
      dataPagination.total = 0
      return
    }

    const requestId = ++latestDataRequestId
    dataLoading.value = true

    try {
      const res = await getDictDataList({
        dict_type: typeCode,
        page: dataPagination.current,
        page_size: dataPagination.pageSize,
        label: dataLabelSearch.value.trim() || undefined,
        status: dataStatusFilter.value === 'all' ? undefined : dataStatusFilter.value,
      })

      if (requestId !== latestDataRequestId || !canApplyDictDataResponse(typeCode, selectedType.value?.type)) {
        return
      }

      dictDataList.value = res.data.list || []
      dataPagination.total = res.data.total
    } finally {
      if (requestId === latestDataRequestId) {
        dataLoading.value = false
      }
    }
  }

  const handleTypeSearch = () => {
    fetchDictTypes(true)
  }

  const handleTypePaginationChange = (page: number, pageSize: number) => {
    typePagination.current = page
    typePagination.pageSize = pageSize
    fetchDictTypes()
  }

  const handleDataPaginationChange = (page: number, pageSize: number) => {
    dataPagination.current = page
    dataPagination.pageSize = pageSize
    fetchDictData()
  }

  const handleSelectType = (record: DictType) => {
    if (selectedType.value?.id === record.id) {
      return
    }
    selectedType.value = record
    selectedDataKeys.value = []
    dataLabelSearch.value = ''
    dataStatusFilter.value = 'all'
    dataPagination.current = 1
    fetchDictData(record.type)
  }

  const handleTypeStatusFilterChange = (value: StatusFilter) => {
    typeStatusFilter.value = value
    fetchDictTypes(true)
  }

  const handleDataSearch = () => {
    dataPagination.current = 1
    fetchDictData()
  }

  const handleDataStatusFilterChange = (value: StatusFilter) => {
    dataStatusFilter.value = value
    dataPagination.current = 1
    fetchDictData()
  }

  const handleToggleTypeStatus = async (record: DictType, checked: boolean) => {
    if (togglingTypeIds.value.has(record.id)) {
      return
    }
    const nextStatus = checked ? 1 : 0
    togglingTypeIds.value.add(record.id)
    try {
      await updateDictType(record.id, { status: nextStatus })
      record.status = nextStatus
      if (selectedType.value?.id === record.id) {
        selectedType.value = { ...record, status: nextStatus }
      }
      message.success(nextStatus === 1 ? '已启用' : '已停用')
    } catch {
      // keep original status
    } finally {
      togglingTypeIds.value.delete(record.id)
    }
  }

  const handleToggleDataStatus = async (record: DictData, checked: boolean) => {
    if (togglingStatusIds.value.has(record.id)) {
      return
    }
    const nextStatus = checked ? 1 : 0
    togglingStatusIds.value.add(record.id)
    try {
      await updateDictData(record.id, { status: nextStatus })
      record.status = nextStatus
      message.success(nextStatus === 1 ? '已启用' : '已停用')
    } catch {
      // keep original status
    } finally {
      togglingStatusIds.value.delete(record.id)
    }
  }

  const handleBatchDeleteData = () => {
    if (!selectedDataKeys.value.length) {
      return
    }
    const ids = [...selectedDataKeys.value]
    Modal.confirm({
      title: `确定删除选中的 ${ids.length} 条字典数据吗？`,
      content: '删除后无法恢复，请确认操作。',
      okType: 'danger',
      okText: '删除',
      cancelText: '取消',
      onOk: async () => {
        batchDeleting.value = true
        try {
          await Promise.all(ids.map(id => deleteDictData(id)))
          message.success(`已删除 ${ids.length} 条`)
          selectedDataKeys.value = []
          await fetchDictData()
        } finally {
          batchDeleting.value = false
        }
      },
    })
  }

  const handleAddType = () => {
    editingType.value = null
    typeDrawerTitle.value = '新增字典类型'
    typeDrawerInitialValue.value = {
      name: '',
      type: '',
      status: 1,
      remark: '',
    }
    typeDrawerVisible.value = true
  }

  const handleEditType = (record: DictType) => {
    editingType.value = record
    typeDrawerTitle.value = '编辑字典类型'
    typeDrawerInitialValue.value = {
      name: record.name,
      type: record.type,
      status: record.status,
      remark: record.remark,
    }
    typeDrawerVisible.value = true
  }

  const handleTypeSubmit = async (values: DictTypeFormValue) => {
    typeSubmitLoading.value = true
    try {
      if (editingType.value) {
        await updateDictType(editingType.value.id, values)
        message.success('更新成功')
      } else {
        await createDictType(values)
        message.success('创建成功')
      }

      typeDrawerVisible.value = false
      await fetchDictTypes()
    } finally {
      typeSubmitLoading.value = false
    }
  }

  const handleDeleteType = async (record: DictType) => {
    await deleteDictType(record.id)
    message.success('删除成功')

    if (selectedType.value?.id === record.id) {
      selectedType.value = null
      dictDataList.value = []
      dataPagination.total = 0
      dataLoading.value = false
    }

    await fetchDictTypes()
  }

  const handleAddData = () => {
    if (!selectedType.value) {
      return
    }

    editingData.value = null
    dataDrawerTitle.value = '新增字典数据'
      dataDrawerInitialValue.value = {
      label: '',
      value: '',
      sort: 0,
      tag_type: '',
      status: 1,
      is_default: 0,
      remark: '',
    }
    dataDrawerVisible.value = true
  }

  const handleEditData = (record: DictData) => {
    editingData.value = record
    dataDrawerTitle.value = '编辑字典数据'
    dataDrawerInitialValue.value = {
      label: record.label,
      value: record.value,
      sort: record.sort,
      tag_type: record.tag_type || '',
      status: record.status,
      is_default: record.is_default,
      remark: record.remark,
    }
    dataDrawerVisible.value = true
  }

  const handleDataSubmit = async (values: DictDataFormValue) => {
    if (!selectedType.value) {
      return
    }

    dataSubmitLoading.value = true
    try {
      if (editingData.value) {
        await updateDictData(editingData.value.id, {
          label: values.label,
          value: values.value,
          sort: values.sort,
          tag_type: values.tag_type,
          status: values.status,
          is_default: values.is_default,
          remark: values.remark,
        })
        message.success('更新成功')
      } else {
        await createDictData({
          ...values,
          dict_type: selectedType.value.type,
        })
        message.success('创建成功')
      }

      dataDrawerVisible.value = false
      await fetchDictData(selectedType.value.type)
    } finally {
      dataSubmitLoading.value = false
    }
  }

  const handleDeleteData = async (record: DictData) => {
    await deleteDictData(record.id)
    message.success('删除成功')
    selectedDataKeys.value = selectedDataKeys.value.filter(id => id !== record.id)
    await fetchDictData()
  }

  const handleCopy = async (text: string, label: string) => {
    try {
      if (!navigator?.clipboard?.writeText) {
        throw new Error('clipboard unavailable')
      }
      await navigator.clipboard.writeText(text)
      message.success(`${label}已复制`)
    } catch {
      message.error(`复制${label}失败`)
    }
  }

  watch(
    typeSearchText,
    value => {
      if (!value.trim()) {
        fetchDictTypes(true)
      }
    },
  )

  onMounted(() => {
    fetchDictTypes()
  })

  return {
    dataDrawerInitialValue,
    dataDrawerTitle,
    dataDrawerVisible,
    dataLoading,
    dataPagination,
    dataSubmitLoading,
    dictDataList,
    dictTypes,
    editingData,
    editingType,
    fetchDictData,
    fetchDictTypes,
    filteredDictTypes,
    handleAddData,
    handleAddType,
    handleCopy,
    handleDataPaginationChange,
    handleDataSubmit,
    handleDeleteData,
    handleDeleteType,
    handleEditData,
    handleEditType,
    handleSelectType,
    handleTypePaginationChange,
    handleTypeSearch,
    handleTypeSubmit,
    handleTypeStatusFilterChange,
    handleDataSearch,
    handleDataStatusFilterChange,
    handleToggleTypeStatus,
    handleToggleDataStatus,
    handleBatchDeleteData,
    dataLabelSearch,
    dataStatusFilter,
    typeStatusFilter,
    selectedDataKeys,
    togglingStatusIds,
    togglingTypeIds,
    batchDeleting,
    selectedType,
    showSelectedOutsideFilter,
    typeDrawerInitialValue,
    typeDrawerTitle,
    typeDrawerVisible,
    typeLoading,
    typePagination,
    typeSearchText,
    typeSubmitLoading,
  }
}
