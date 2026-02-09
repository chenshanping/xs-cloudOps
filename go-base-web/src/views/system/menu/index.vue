<template>
  <div class="menu-page">
    <a-card>
      <template #extra>
        <a-button type="primary" @click="handleAdd()" v-permission="'system:menu:add'">
          <PlusOutlined /> 新增
        </a-button>
      </template>
      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        row-key="id"
        :pagination="false"
        default-expand-all-rows
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'icon'">
          <template v-if="record.icon">
            <component v-if="!record.icon.startsWith('custom-')" :is="getIconComponent(record.icon)" style="font-size: 18px" />
            <img v-else :src="getCustomIconUrl(record.icon)" style="width: 1em; height: 1em" alt="" />
<!--            <svg v-else style="width: 18px; height: 18px">-->
<!--              <use :xlink:href="getCustomIconUrl(record.icon)"></use>-->
<!--            </svg>-->

          </template>
            <span v-else>-</span>
          </template>
          <template v-if="column.key === 'type'">
            <a-tag :color="record.type === 1 ? 'blue' : record.type === 2 ? 'green' : 'orange'">
              {{ record.type === 1 ? '目录' : record.type === 2 ? '菜单' : '按鑢' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'green' : 'red'">
              {{ record.status === 1 ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="handleAdd(record.id)" v-permission="'system:menu:add'">新增子菜单</a-button>
            <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:menu:edit'">编辑</a-button>
            <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record)">
              <a-button type="link" size="small" danger v-permission="'system:menu:delete'">删除</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑抽屉 -->
    <a-drawer
      v-model:open="drawerVisible"
      :title="modalTitle"
      width="600"
      placement="right"
      @ok="handleModalOk"
    >
      <a-form :model="formState" :label-col="{ span: 6 }">
        <a-form-item label="上级菜单">
          <a-tree-select
            v-model:value="formState.parent_id"
            :tree-data="menuTreeOptions"
            :field-names="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择上级菜单"
            allow-clear
            tree-default-expand-all
          />
        </a-form-item>
        <a-form-item label="菜单类型" required>
          <a-radio-group v-model:value="formState.type">
            <a-radio :value="1">目录</a-radio>
            <a-radio :value="2">菜单</a-radio>
            <a-radio :value="3">按顲</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="菜单名称" required>
          <a-input v-model:value="formState.name" />
        </a-form-item>
        <a-form-item label="路由路径" v-if="formState.type !== 3">
          <a-input v-model:value="formState.path" />
        </a-form-item>
        <a-form-item label="组件路径" v-if="formState.type === 2">
          <a-input v-model:value="formState.component" />
        </a-form-item>
        <a-form-item label="权限标识" v-if="formState.type !== 1">
          <a-input v-model:value="formState.permission" />
        </a-form-item>
        <a-form-item label="图标" v-if="formState.type !== 3">
          <IconSelect v-model="formState.icon" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="formState.sort" :min="0" />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch v-model:checked="formState.statusChecked" />
        </a-form-item>
        <a-form-item label="隐藏" v-if="formState.type !== 3">
          <a-switch v-model:checked="formState.hiddenChecked" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="drawerVisible = false">取消</a-button>
        <a-button type="primary" @click="handleModalOk">保存</a-button>
      </template>
    </a-drawer>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import * as AntIcons from '@ant-design/icons-vue'
import IconSelect from '@/components/IconSelect.vue'
import { getMenuTree, createMenu, updateMenu, deleteMenu } from '@/api/menu'
import { getAllApis } from '@/api/api'
import { useTableColumns } from '@/utils/permission'
import type { Menu, Api } from '@/types'

const getIconComponent = (iconName?: string) => {
  if (!iconName) return 'DatabaseOutlined'
  
  let name = iconName
  if (name.startsWith('official-')) {
    // 官方图标已经包含完整的组件名称
    name = name.replace('official-', '')
  } else if (name.startsWith('custom-')) {
    return null
  }
  // 其他情况，假设是官方图标的完整组件名称
  return (AntIcons as any)[name] || 'DatabaseOutlined'
}

// 预加载所有自定义图标
const iconModules = import.meta.glob('@/assets/icons/*.svg', { eager: true, query: '?url', import: 'default' })

const getCustomIconPath = (iconName?: string) => {
  if (!iconName) return ''
  const name = iconName.replace('custom-', '')
  return `/src/assets/icons/${name}.svg#icon`
}

const getCustomIconUrl = (iconName?: string) => {
  if (!iconName) return ''
  const name = iconName.replace('custom-', '')
  const key = `/src/assets/icons/${name}.svg`
  return (iconModules[key] as string) || ''
}

const loading = ref(false)
const tableData = ref<Menu[]>([])
const drawerVisible = ref(false)
const modalTitle = ref('新增菜单')
const isEdit = ref(false)
const currentId = ref(0)


const formState = reactive({
  parent_id: 0 as number | undefined,
  name: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  type: 2,
  permission: '',
  statusChecked: true,
  hiddenChecked: false
})

// 使用工具函数动态生成列配置
const columns = useTableColumns(
  [
    { title: '菜单名称', dataIndex: 'name', key: 'name', width: 180 },
    { title: '图标', dataIndex: 'icon', key: 'icon', width: 80 },
    { title: '类型', key: 'type', width: 80 },
    { title: '路由路径', dataIndex: 'path', key: 'path' },
    { title: '组件路径', dataIndex: 'component', key: 'component' },
    { title: '权限标识', dataIndex: 'permission', key: 'permission' },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 60 },
    { title: '状态', key: 'status', width: 80 },
  ],
  { title: '操作', key: 'action', width: 280 },
  ['system:menu:add', 'system:menu:edit', 'system:menu:delete']
)
















const menuTreeOptions = computed(() => {
  return [{ id: 0, name: '顶级菜单', children: tableData.value }]
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getMenuTree()
    tableData.value = res.data
  } finally {
    loading.value = false
  }
}

const handleAdd = (parentId?: number) => {
  isEdit.value = false
  modalTitle.value = '新增菜单'
  Object.assign(formState, {
    parent_id: parentId || 0,
    name: '',
    path: '',
    component: '',
    icon: '',
    sort: 0,
    type: 2,
    permission: '',
    statusChecked: true,
    hiddenChecked: false
  })
  drawerVisible.value = true
}

const handleEdit = (record: Menu) => {
  isEdit.value = true
  modalTitle.value = '编辑菜单'
  currentId.value = record.id
  // 处理图标前缀
  let iconValue = record.icon || ''
  if (iconValue && !iconValue.startsWith('official-') && !iconValue.startsWith('custom-')) {
    // 如果没有前缀，假设是官方图标
    iconValue = `official-${iconValue}`
  }
  Object.assign(formState, {
    parent_id: record.parent_id || undefined,
    name: record.name,
    path: record.path,
    component: record.component,
    icon: iconValue,
    sort: record.sort,
    type: record.type,
    permission: record.permission,
    statusChecked: record.status === 1,
    hiddenChecked: record.hidden === 1
  })
  drawerVisible.value = true
}

const handleModalOk = async () => {
  // 处理图标前缀，移除前缀后再保存
  let iconValue = formState.icon || ''
  if (iconValue.startsWith('official-')) {
    iconValue = iconValue.replace('official-', '')
  }
  
  const data = {
    parent_id: formState.parent_id || 0,
    name: formState.name,
    path: formState.path,
    component: formState.component,
    icon: iconValue,
    sort: formState.sort,
    type: formState.type,
    permission: formState.permission,
    status: formState.statusChecked ? 1 : 0,
    hidden: formState.hiddenChecked ? 1 : 0
  }
  if (isEdit.value) {
    await updateMenu(currentId.value, data)
    message.success('更新成功')
  } else {
    await createMenu(data)
    message.success('创建成功')
  }
  drawerVisible.value = false
  fetchData()
}

const handleDelete = async (record: Menu) => {
  await deleteMenu(record.id)
  message.success('删除成功')
  fetchData()
}




onMounted(() => {
  fetchData()
})
</script>
