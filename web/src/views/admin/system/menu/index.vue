<template>
  <PageWrapper class="menu-page">
    <div class="menu-page__content">
      <div class="menu-summary">
      <div class="summary-item summary-item--primary">
        <div class="summary-item__icon">
          <AppstoreOutlined />
        </div>
        <div>
          <div class="summary-item__label">菜单总数</div>
          <div class="summary-item__value">{{ menuCount }}</div>
        </div>
      </div>
      <div class="summary-item">
        <div class="summary-item__icon summary-item__icon--success">
          <MenuOutlined />
        </div>
        <div>
          <div class="summary-item__label">页面菜单</div>
          <div class="summary-item__value">{{ pageMenuCount }}</div>
        </div>
      </div>
      <div class="summary-item">
        <div class="summary-item__icon summary-item__icon--warning">
          <SafetyCertificateOutlined />
        </div>
        <div>
          <div class="summary-item__label">按钮权限</div>
          <div class="summary-item__value">{{ buttonPermissionCount }}</div>
        </div>
      </div>
    </div>

      <a-card :bordered="false" class="menu-table-card">
        <div class="menu-table-card__header">
        <div>
          <h2 class="menu-table-card__title">菜单结构</h2>
          <div class="menu-table-card__subtitle">目录、页面和按钮权限共用同一棵菜单树</div>
        </div>
        <a-space>
          <a-button :disabled="!expandableMenuIds.length || isAllExpanded" @click="handleExpandAll">
            <DownOutlined /> 全部展开
          </a-button>
          <a-button :disabled="!expandedRowKeys.length" @click="handleCollapseAll">
            <UpOutlined /> 全部收起
          </a-button>
          <a-button type="primary" @click="handleAdd()" v-permission="'system:menu:add'">
            <PlusOutlined /> 新增
          </a-button>
        </a-space>
      </div>
        <a-table
          :columns="columns"
          :data-source="tableData"
          :loading="loading"
          row-key="id"
          :pagination="false"
          :expanded-row-keys="expandedRowKeys"
          :scroll="{ x: 1200 }"
          @expand="handleExpand"
        >
          <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="menu-name-cell">
              <span class="menu-name-cell__mark" :class="`menu-name-cell__mark--${record.type}`"></span>
              <div class="menu-name-cell__content">
                <div class="menu-name-cell__main">
                  <span>{{ record.name }}</span>
                  <a-tag v-if="record.hidden === 1" color="default">隐藏</a-tag>
                </div>
              </div>
            </div>
          </template>
          <template v-if="column.key === 'icon'">
            <span v-if="record.icon" class="menu-icon-box">
              <component v-if="!record.icon.startsWith('custom-')" :is="getIconComponent(record.icon)" />
              <img v-else :src="getCustomIconUrl(record.icon)" alt="" />
            </span>
            <span v-else class="menu-empty-text">-</span>
          </template>
          <template v-if="column.key === 'type'">
            <a-tag :color="getMenuTypeMeta(record.type).color">
              {{ record.type === 1 ? '目录' : record.type === 2 ? '菜单' : '按钮' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'path'">
            <a-tooltip v-if="record.path" :title="record.path">
              <a-typography-text code class="menu-ellipsis-code">{{ record.path }}</a-typography-text>
            </a-tooltip>
            <span v-else class="menu-empty-text">-</span>
          </template>
          <template v-if="column.key === 'component'">
            <a-tooltip v-if="record.component" :title="record.component">
              <a-typography-text code class="menu-ellipsis-code">{{ record.component }}</a-typography-text>
            </a-tooltip>
            <span v-else class="menu-empty-text">-</span>
          </template>
          <template v-if="column.key === 'route_mode'">
            <a-tag v-if="record.type === 2" :color="record.is_standalone === 1 ? 'purple' : 'default'">
              {{ record.is_standalone === 1 ? '独立页' : 'Layout页' }}
            </a-tag>
            <span v-else class="menu-empty-text">-</span>
          </template>
          <template v-if="column.key === 'permission'">
            <a-tooltip v-if="record.permission" :title="record.permission">
              <a-typography-text code class="menu-ellipsis-code">{{ record.permission }}</a-typography-text>
            </a-tooltip>
            <span v-else class="menu-empty-text">-</span>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'green' : 'red'">
              {{ record.status === 1 ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space :size="0">
              <a-button type="link" size="small" @click="handleAdd(record.id)" v-permission="'system:menu:add'">新增子菜单</a-button>
              <a-button type="link" size="small" @click="handleEdit(record)" v-permission="'system:menu:edit'">编辑</a-button>
              <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger v-permission="'system:menu:delete'">删除</a-button>
              </a-popconfirm>
            </a-space>
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
    >
      <a-form class="menu-drawer-form" :model="formState" layout="vertical">
        <section class="drawer-section">
          <div class="drawer-section__title">基础信息</div>
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
            <a-radio-group v-model:value="formState.type" class="menu-type-switch">
              <a-radio-button :value="1">目录</a-radio-button>
              <a-radio-button :value="2">菜单</a-radio-button>
              <a-radio-button :value="3">按钮</a-radio-button>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="菜单名称" required>
            <a-input v-model:value="formState.name" placeholder="请输入菜单名称" />
          </a-form-item>
          <a-form-item label="图标" v-if="formState.type !== 3">
            <IconSelect v-model="formState.icon" />
          </a-form-item>
        </section>

        <section class="drawer-section">
          <div class="drawer-section__title">路由与权限</div>
          <a-form-item label="路由路径" v-if="formState.type !== 3">
            <a-input v-model:value="formState.path" placeholder="例如 /system/user" />
          </a-form-item>
          <a-form-item label="组件路径" v-if="formState.type === 2">
            <a-input v-model:value="formState.component" placeholder="例如 system/user/index" />
          </a-form-item>
          <a-form-item label="页面模式" v-if="formState.type === 2">
            <a-radio-group v-model:value="formState.is_standalone" class="menu-type-switch">
              <a-radio-button :value="0">使用 Layout</a-radio-button>
              <a-radio-button :value="1">独立页面</a-radio-button>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="权限标识" v-if="formState.type !== 1">
            <a-input v-model:value="formState.permission" placeholder="例如 system:user:list" />
          </a-form-item>
        </section>

        <section class="drawer-section">
          <div class="drawer-section__title">显示设置</div>
          <div class="drawer-form-grid">
            <a-form-item label="排序">
              <a-input-number v-model:value="formState.sort" :min="0" style="width: 100%" />
            </a-form-item>
            <a-form-item label="状态">
              <a-switch v-model:checked="formState.statusChecked" />
            </a-form-item>
            <a-form-item label="隐藏" v-if="formState.type !== 3">
              <a-switch v-model:checked="formState.hiddenChecked" />
            </a-form-item>
          </div>
        </section>
      </a-form>
      <template #footer>
        <div class="drawer-footer">
          <a-button @click="drawerVisible = false">取消</a-button>
          <a-button type="primary" @click="handleModalOk">保存</a-button>
        </div>
      </template>
    </a-drawer>

    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { AppstoreOutlined, DownOutlined, MenuOutlined, PlusOutlined, SafetyCertificateOutlined, UpOutlined } from '@ant-design/icons-vue'
import * as AntIcons from '@ant-design/icons-vue'
import IconSelect from '@/components/IconSelect.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import { getMenuTree, createMenu, updateMenu, deleteMenu } from '@/api/menu'
import { useTableColumns } from '@/utils/permission'
import { useUserStore } from '@/store/user'
import type { Menu } from '@/types'

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
const expandedRowKeys = ref<number[]>([])
const drawerVisible = ref(false)
const modalTitle = ref('新增菜单')
const isEdit = ref(false)
const currentId = ref(0)
const userStore = useUserStore()


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
  hiddenChecked: false,
  is_standalone: 0
})

// 使用工具函数动态生成列配置
const columns = useTableColumns(
  [
    { title: '菜单名称', dataIndex: 'name', key: 'name', width: 260 },
    { title: '图标', dataIndex: 'icon', key: 'icon', width: 80 },
    { title: '类型', key: 'type', width: 80 },
    { title: '路由路径', dataIndex: 'path', key: 'path', ellipsis: true },
    { title: '组件路径', dataIndex: 'component', key: 'component', ellipsis: true },
    { title: '页面模式', key: 'route_mode', width: 110, align: 'center' },
    { title: '权限标识', dataIndex: 'permission', key: 'permission', ellipsis: true },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 60 },
    { title: '状态', key: 'status', width: 80 },
  ],
  { title: '操作', key: 'action', width: 280 },
  ['system:menu:add', 'system:menu:edit', 'system:menu:delete']
)

const flattenMenus = (menus: Menu[]): Menu[] =>
  menus.flatMap(menu => [menu, ...(menu.children ? flattenMenus(menu.children) : [])])

const collectExpandableMenuIds = (menus: Menu[]): number[] =>
  menus.flatMap(menu => {
    if (!menu.children?.length) {
      return []
    }
    return [menu.id, ...collectExpandableMenuIds(menu.children)]
  })

const collectMenuIds = (menu: Menu): number[] => [
  menu.id,
  ...(menu.children ? menu.children.flatMap(child => collectMenuIds(child)) : [])
]

const isTopLevelMenu = (record: Menu) =>
  tableData.value.some(menu => menu.id === record.id)

const allMenus = computed(() => flattenMenus(tableData.value))
const expandableMenuIds = computed(() => collectExpandableMenuIds(tableData.value))
const isAllExpanded = computed(() =>
  expandableMenuIds.value.length > 0 &&
  expandableMenuIds.value.every(id => expandedRowKeys.value.includes(id))
)
const menuCount = computed(() => allMenus.value.length)
const pageMenuCount = computed(() => allMenus.value.filter(menu => menu.type === 2).length)
const buttonPermissionCount = computed(() => allMenus.value.filter(menu => menu.type === 3).length)

const getMenuTypeMeta = (type: number) => {
  if (type === 1) {
    return { color: 'blue' }
  }
  if (type === 2) {
    return { color: 'green' }
  }
  return { color: 'orange' }
}
const menuTreeOptions = computed(() => {
  return [{ id: 0, name: '顶级菜单', children: tableData.value }]
})


const fetchData = async (preserveExpand = true) => {
  loading.value = true
  try {
    const res = await getMenuTree()
    tableData.value = res.data
    if (!preserveExpand) {
      expandedRowKeys.value = []
    } else {
      // 清理已不存在的 key，保留仍有效的展开状态
      const validIds = new Set(collectExpandableMenuIds(res.data))
      expandedRowKeys.value = expandedRowKeys.value.filter(id => validIds.has(id))
    }
  } finally {
    loading.value = false
  }
}


const handleExpand = (expanded: boolean, record: Menu) => {
  const keys = new Set(expandedRowKeys.value)

  if (expanded && isTopLevelMenu(record)) {
    tableData.value
      .filter(menu => menu.id !== record.id)
      .flatMap(menu => collectMenuIds(menu))
      .forEach(id => keys.delete(id))
  }

  if (expanded) {
    keys.add(record.id)
  } else {
    const idsToRemove = isTopLevelMenu(record) ? collectMenuIds(record) : [record.id]
    idsToRemove.forEach(id => keys.delete(id))
  }
  expandedRowKeys.value = Array.from(keys)
}

const handleExpandAll = () => {
  expandedRowKeys.value = expandableMenuIds.value
  message.success('已展开全部菜单')
}

const handleCollapseAll = () => {
  expandedRowKeys.value = []
  message.success('已收起全部菜单')
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
    hiddenChecked: false,
    is_standalone: 0
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
    hiddenChecked: record.hidden === 1,
    is_standalone: record.is_standalone || 0
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
    hidden: formState.hiddenChecked ? 1 : 0,
    is_standalone: formState.type === 2 ? formState.is_standalone : 0
  }
  if (isEdit.value) {
    await updateMenu(currentId.value, data)
    message.success('更新成功')
  } else {
    await createMenu(data)
    message.success('创建成功')
  }
  drawerVisible.value = false
  await Promise.all([fetchData(), userStore.refreshAccessAction()])
}

const handleDelete = async (record: Menu) => {
  await deleteMenu(record.id)
  message.success('删除成功')
  await Promise.all([fetchData(), userStore.refreshAccessAction()])
}





onMounted(() => {
  fetchData(false)
})
</script>

<style scoped>
.menu-page__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.menu-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 88px;
  padding: 16px;
  background: #fff;
  border: 1px solid #edf0f5;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgb(15 23 42 / 4%);
}

.summary-item--primary {
  border-color: #d6e4ff;
}

.summary-item__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  color: #1677ff;
  font-size: 20px;
  background: #eef5ff;
  border-radius: 8px;
}

.summary-item__icon--success {
  color: #389e0d;
  background: #f0f8e8;
}

.summary-item__icon--warning {
  color: #d46b08;
  background: #fff4e6;
}

.summary-item__icon--muted {
  color: #45556c;
  background: #f3f5f8;
}

.summary-item__label {
  color: #8c8c8c;
  font-size: 12px;
}

.summary-item__value {
  margin-top: 4px;
  color: #262626;
  font-size: 24px;
  font-weight: 650;
  line-height: 1;
}

.menu-table-card {
  border-radius: 8px;
  box-shadow: 0 4px 16px rgb(15 23 42 / 4%);
}

.menu-table-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.menu-table-card__title {
  margin: 0;
  color: #262626;
  font-size: 18px;
  font-weight: 650;
  line-height: 1.4;
}

.menu-table-card__subtitle {
  margin-top: 4px;
  color: #8c8c8c;
  font-size: 12px;
}

.menu-name-cell {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  min-width: 0;
}

.menu-name-cell__mark {
  width: 8px;
  height: 8px;
  margin-top: 7px;
  border-radius: 50%;
}

.menu-name-cell__mark--1 {
  background: #1677ff;
  box-shadow: 0 0 0 4px #e6f4ff;
}

.menu-name-cell__mark--2 {
  background: #52c41a;
  box-shadow: 0 0 0 4px #f6ffed;
}

.menu-name-cell__mark--3 {
  background: #fa8c16;
  box-shadow: 0 0 0 4px #fff7e6;
}

.menu-name-cell__main {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
  font-weight: 600;
}

.menu-icon-box {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  color: #1677ff;
  font-size: 17px;
  background: #f3f7ff;
  border: 1px solid #e6f0ff;
  border-radius: 8px;
}

.menu-icon-box img {
  width: 18px;
  height: 18px;
}

.menu-ellipsis-code {
  display: block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.menu-empty-text {
  color: #bfbfbf;
}

.menu-drawer-form {
  padding-bottom: 12px;
}

.drawer-section {
  padding: 14px;
  margin-bottom: 14px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
}

.drawer-section__title {
  margin-bottom: 12px;
  color: #262626;
  font-size: 15px;
  font-weight: 650;
}

.menu-type-switch {
  width: 100%;
}

.menu-type-switch :deep(.ant-radio-button-wrapper) {
  width: 33.333%;
  text-align: center;
}

.drawer-form-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

:deep(.ant-table-thead > tr > th) {
  color: #595959;
  font-size: 12px;
  font-weight: 650;
  background: #fafafa;
}

@media (max-width: 1180px) {
  .menu-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .menu-summary,
  .drawer-form-grid {
    grid-template-columns: 1fr;
  }

  .menu-table-card__header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
