<template>
  <a-drawer
    v-model:open="visible"
    :title="`分配权限 - ${roleName}`"
    width="900"
    placement="right"
    class="permission-drawer"
  >
    <a-tabs v-model:activeKey="activeTabKey">
      <!-- 角色菜单 -->
      <a-tab-pane key="menu" tab="角色菜单">
        <div class="permission-header">
          <a-checkbox
            :checked="isAllMenuChecked"
            :indeterminate="isMenuIndeterminateAll"
            @change="handleCheckAllMenus"
          >
            全选
          </a-checkbox>
          <span class="selected-count">已选 {{ checkedMenuKeys.length }} 个菜单</span>
          <a-input-search
            v-model:value="menuSearchText"
            placeholder="搜索菜单名称"
            style="width: 200px; margin-left: auto"
            allow-clear
          />
        </div>
        <div class="menu-layout">
          <!-- 左侧一级菜单列表 -->
          <div class="menu-group-list">
            <div class="group-list-header">一级菜单</div>
            <div class="group-list-content">
              <div
                v-for="menu in menuTree"
                :key="menu.id"
                :class="['group-item', { active: selectedMenuId === menu.id }]"
                @click="selectedMenuId = menu.id"
              >
                <a-checkbox
                  :checked="isMenuChecked(menu.id)"
                  :indeterminate="isMenuIndeterminate(menu)"
                  @change="handleMenuChange(menu, $event)"
                  @click.stop
                />
                <span class="group-item-name">{{ menu.name }}</span>
                <a-badge
                  :count="getMenuSelectedCount(menu)"
                  :number-style="{ backgroundColor: '#52c41a' }"
                  :show-zero="false"
                />
                <span class="group-item-total">({{ getMenuTotalCount(menu) }})</span>
              </div>
            </div>
          </div>
          <!-- 右侧子菜单列表 -->
          <div class="menu-detail-list">
            <div class="menu-detail-header">
              <span>{{ selectedMenu?.name || '请选择菜单' }} 子菜单</span>
              <a-checkbox
                v-if="selectedMenu?.children?.length"
                :checked="isMenuChecked(selectedMenu.id)"
                :indeterminate="isMenuIndeterminate(selectedMenu)"
                @change="handleMenuChange(selectedMenu, $event)"
              >
                全选当前菜单
              </a-checkbox>
            </div>
            <div class="menu-detail-content">
              <template v-if="filteredChildMenus.length">
                <div
                  v-for="child in filteredChildMenus"
                  :key="child.id"
                  class="menu-child-block"
                >
                  <div class="menu-child-header">
                    <a-checkbox
                      :checked="isMenuChecked(child.id)"
                      :indeterminate="isMenuIndeterminate(child)"
                      @change="handleMenuChange(child, $event)"
                    >
                      <span class="menu-child-name">
                        <a-tag v-if="child.type === 2" color="blue" size="small">菜单</a-tag>
                        <a-tag v-else-if="child.type === 3" color="orange" size="small">按钮</a-tag>
                        {{ child.name }}
                      </span>
                    </a-checkbox>
                    <span v-if="child.permission" class="permission-code">{{ child.permission }}</span>
                  </div>
                  <!-- 按钮权限 -->
                  <div v-if="child.children?.length" class="menu-child-items">
                    <div
                      v-for="subChild in child.children"
                      :key="subChild.id"
                      class="btn-permission-item"
                    >
                      <a-checkbox
                        :checked="isMenuChecked(subChild.id)"
                        @change="handleMenuChange(subChild, $event)"
                      >
                        <span class="btn-permission-content">
                          <a-tag color="orange" size="small">按钮</a-tag>
                          <span class="btn-name">{{ subChild.name }}</span>
                        </span>
                      </a-checkbox>
                      <span v-if="subChild.permission" class="permission-code">{{ subChild.permission }}</span>
                    </div>
                  </div>
                </div>
              </template>
              <a-empty v-else-if="selectedMenu" description="暂无子菜单" />
              <a-empty v-else description="请在左侧选择一级菜单" />
            </div>
          </div>
        </div>
      </a-tab-pane>

      <!-- 角色API -->
      <a-tab-pane key="api" tab="角色API">
        <div class="permission-header">
          <a-checkbox
            :checked="isAllApiChecked"
            :indeterminate="isApiIndeterminateAll"
            @change="handleCheckAllApis"
          >
            全选
          </a-checkbox>
          <span class="selected-count">已选 {{ checkedApiIds.length }} 个 API</span>
          <a-input-search
            v-model:value="apiSearchText"
            placeholder="搜索接口路径或描述"
            style="width: 200px; margin-left: auto"
            allow-clear
          />
        </div>
        <div class="api-layout">
          <!-- 左侧分组列表 -->
          <div class="api-group-list">
            <div class="group-list-header">API 分组</div>
            <div class="group-list-content">
              <div
                v-for="group in apiGroups"
                :key="group.name"
                :class="['group-item', { active: selectedGroup === group.name }]"
                @click="selectedGroup = group.name"
              >
                <a-checkbox
                  :checked="isGroupChecked(group.name)"
                  :indeterminate="isGroupIndeterminate(group.name)"
                  @change="handleGroupChange(group.name, $event)"
                  @click.stop
                />
                <span class="group-item-name">{{ group.name || '未分组' }}</span>
                <a-badge
                  :count="getGroupSelectedCount(group.name)"
                  :number-style="{ backgroundColor: '#52c41a' }"
                  :show-zero="false"
                />
                <span class="group-item-total">({{ group.apis.length }})</span>
              </div>
            </div>
          </div>
          <!-- 右侧接口列表 -->
          <div class="api-detail-list">
            <div class="api-detail-header">
              <span>{{ selectedGroup || '未分组' }} 接口列表</span>
              <a-checkbox
                v-if="selectedGroupApis.length"
                :checked="isGroupChecked(selectedGroup)"
                :indeterminate="isGroupIndeterminate(selectedGroup)"
                @change="handleGroupChange(selectedGroup, $event)"
              >
                全选当前分组
              </a-checkbox>
            </div>
            <div class="api-detail-content">
              <template v-if="filteredGroupApis.length">
                <div
                  v-for="api in filteredGroupApis"
                  :key="api.id"
                  class="api-item"
                >
                  <a-checkbox
                    :checked="checkedApiIds.includes(api.id)"
                    @change="handleApiChange(api.id, $event)"
                  >
                    <div class="api-item-content">
                      <a-tag :color="getMethodColor(api.method)" size="small">{{ api.method }}</a-tag>
                      <span class="api-path">{{ api.path }}</span>
                      <span class="api-desc">{{ api.description }}</span>
                    </div>
                  </a-checkbox>
                </div>
              </template>
              <a-empty v-else description="暂无接口" />
            </div>
          </div>
        </div>
      </a-tab-pane>
    </a-tabs>
    <template #footer>
      <div style="display: flex; justify-content: flex-end; gap: 8px">
        <a-button @click="visible = false">取消</a-button>
        <a-button type="primary" :loading="saveLoading" @click="handleSavePermissions">保存</a-button>
      </div>
    </template>
  </a-drawer>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { getRole, assignMenus, assignApis } from '@/api/role'
import { getMenuTree } from '@/api/menu'
import { getAllApis } from '@/api/api'
import { useUserStore } from '@/store/user'
import type { Menu, Api } from '@/types'

const userStore = useUserStore()

interface Props {
  roleId: number
  roleName: string
}

const props = defineProps<Props>()
const visible = defineModel<boolean>('open', { default: false })

const saveLoading = ref(false)
const activeTabKey = ref('menu')
const menuTree = ref<Menu[]>([])
const allApis = ref<Api[]>([])
const checkedMenuKeys = ref<number[]>([])
const checkedApiIds = ref<number[]>([])
const selectedGroup = ref('')
const apiSearchText = ref('')
const selectedMenuId = ref<number | null>(null)
const menuSearchText = ref('')

// 按分组组织API
interface ApiGroup {
  name: string
  apis: Api[]
}

const apiGroups = computed<ApiGroup[]>(() => {
  const groups: Record<string, Api[]> = {}
  allApis.value.forEach(api => {
    const groupName = api.group || ''
    if (!groups[groupName]) {
      groups[groupName] = []
    }
    groups[groupName].push(api)
  })
  return Object.entries(groups).map(([name, apis]) => ({ name, apis }))
})

// 当前选中分组的接口
const selectedGroupApis = computed(() => {
  const group = apiGroups.value.find(g => g.name === selectedGroup.value)
  return group?.apis || []
})

// 当前选中的一级菜单
const selectedMenu = computed(() => {
  return menuTree.value.find(m => m.id === selectedMenuId.value)
})

// 过滤后的子菜单列表
const filteredChildMenus = computed(() => {
  const children = selectedMenu.value?.children || []
  if (!menuSearchText.value) return children
  const keyword = menuSearchText.value.toLowerCase()
  return children.filter(child => {
    // 匹配二级菜单名称
    if (child.name.toLowerCase().includes(keyword)) return true
    // 匹配三级菜单名称
    if (child.children?.some(sub => sub.name.toLowerCase().includes(keyword))) return true
    return false
  })
})

// 获取一级菜单下已选数量
const getMenuSelectedCount = (menu: Menu) => {
  let count = 0
  if (checkedMenuKeys.value.includes(menu.id)) count++
  if (menu.children?.length) {
    menu.children.forEach(child => {
      if (checkedMenuKeys.value.includes(child.id)) count++
      if (child.children?.length) {
        child.children.forEach(sub => {
          if (checkedMenuKeys.value.includes(sub.id)) count++
        })
      }
    })
  }
  return count
}

// 获取一级菜单下总数量
const getMenuTotalCount = (menu: Menu) => {
  let count = 1 // 包含自己
  if (menu.children?.length) {
    menu.children.forEach(child => {
      count++
      if (child.children?.length) {
        count += child.children.length
      }
    })
  }
  return count
}

// 过滤后的接口列表
const filteredGroupApis = computed(() => {
  if (!apiSearchText.value) return selectedGroupApis.value
  const keyword = apiSearchText.value.toLowerCase()
  return selectedGroupApis.value.filter(api => 
    api.path.toLowerCase().includes(keyword) || 
    (api.description && api.description.toLowerCase().includes(keyword))
  )
})

// 获取分组已选数量
const getGroupSelectedCount = (groupName: string) => {
  const group = apiGroups.value.find(g => g.name === groupName)
  if (!group) return 0
  return group.apis.filter(api => checkedApiIds.value.includes(api.id)).length
}

// 获取所有菜单ID
const getAllMenuIds = (menus: Menu[]): number[] => {
  const ids: number[] = []
  const traverse = (items: Menu[]) => {
    items.forEach(item => {
      ids.push(item.id)
      if (item.children?.length) traverse(item.children)
    })
  }
  traverse(menus)
  return ids
}

// 菜单全选相关
const isAllMenuChecked = computed(() => {
  const allMenuIds = getAllMenuIds(menuTree.value)
  return allMenuIds.length > 0 && allMenuIds.every(id => checkedMenuKeys.value.includes(id))
})

const isMenuIndeterminateAll = computed(() => {
  const allMenuIds = getAllMenuIds(menuTree.value)
  const checked = checkedMenuKeys.value.filter(id => allMenuIds.includes(id)).length
  return checked > 0 && checked < allMenuIds.length
})

// API全选相关
const isAllApiChecked = computed(() => {
  return allApis.value.length > 0 && allApis.value.every(api => checkedApiIds.value.includes(api.id))
})

const isApiIndeterminateAll = computed(() => {
  const checked = checkedApiIds.value.length
  return checked > 0 && checked < allApis.value.length
})

const isMenuChecked = (menuId: number) => {
  return checkedMenuKeys.value.includes(menuId)
}

const isMenuIndeterminate = (menu: Menu) => {
  if (!menu.children?.length) return false
  const childIds = getAllMenuIds(menu.children)
  const checked = childIds.filter(id => checkedMenuKeys.value.includes(id)).length
  return checked > 0 && checked < childIds.length
}

// API分组相关
const isGroupChecked = (groupName: string) => {
  const group = apiGroups.value.find(g => g.name === groupName)
  if (!group || !group.apis.length) return false
  return group.apis.every(api => checkedApiIds.value.includes(api.id))
}

const isGroupIndeterminate = (groupName: string) => {
  const group = apiGroups.value.find(g => g.name === groupName)
  if (!group) return false
  const checked = group.apis.filter(api => checkedApiIds.value.includes(api.id)).length
  return checked > 0 && checked < group.apis.length
}

const getMethodColor = (method: string) => {
  const colors: Record<string, string> = {
    'GET': 'green',
    'POST': 'blue',
    'PUT': 'orange',
    'DELETE': 'red',
    'PATCH': 'purple'
  }
  return colors[method.toUpperCase()] || 'default'
}

const handleCheckAllMenus = (e: any) => {
  if (e.target.checked) {
    checkedMenuKeys.value = getAllMenuIds(menuTree.value)
  } else {
    checkedMenuKeys.value = []
  }
}

const handleCheckAllApis = (e: any) => {
  if (e.target.checked) {
    checkedApiIds.value = allApis.value.map(api => api.id)
  } else {
    checkedApiIds.value = []
  }
}

const handleMenuChange = (menu: Menu, e: any) => {
  const menuIds = [menu.id]
  // 收集所有子菜单ID
  if (menu.children?.length) {
    menuIds.push(...getAllMenuIds(menu.children))
  }
  
  if (e.target.checked) {
    checkedMenuKeys.value = [...new Set([...checkedMenuKeys.value, ...menuIds])]
  } else {
    checkedMenuKeys.value = checkedMenuKeys.value.filter(id => !menuIds.includes(id))
  }
}

const handleGroupChange = (groupName: string, e: any) => {
  const group = apiGroups.value.find(g => g.name === groupName)
  if (!group) return
  const groupApiIds = group.apis.map(api => api.id)
  
  if (e.target.checked) {
    checkedApiIds.value = [...new Set([...checkedApiIds.value, ...groupApiIds])]
  } else {
    checkedApiIds.value = checkedApiIds.value.filter(id => !groupApiIds.includes(id))
  }
}

const handleApiChange = (apiId: number, e: any) => {
  if (e.target.checked) {
    checkedApiIds.value = [...checkedApiIds.value, apiId]
  } else {
    checkedApiIds.value = checkedApiIds.value.filter(id => id !== apiId)
  }
}

const fetchMenuTree = async () => {
  const res = await getMenuTree()
  menuTree.value = res.data
  // 默认选中第一个一级菜单
  if (menuTree.value.length > 0) {
    selectedMenuId.value = menuTree.value[0].id
  }
}

const fetchAllApis = async () => {
  const res = await getAllApis()
  allApis.value = res.data
  // 默认选中第一个分组
  if (apiGroups.value.length > 0) {
    selectedGroup.value = apiGroups.value[0].name
  }
}

const loadRolePermissions = async () => {
  if (!props.roleId) return
  const res = await getRole(props.roleId)
  checkedMenuKeys.value = res.data.menus?.map((m: Menu) => m.id) || []
  checkedApiIds.value = res.data.apis?.map((a: Api) => a.id) || []
}

const handleSavePermissions = async () => {
  saveLoading.value = true
  try {
    await Promise.all([
      assignMenus(props.roleId, checkedMenuKeys.value),
      assignApis(props.roleId, checkedApiIds.value)
    ])
    message.success('权限分配成功')
    visible.value = false
    // 刷新当前用户的菜单和权限
    userStore.getUserInfoAction()
  } finally {
    saveLoading.value = false
  }
}

// 监听打开状态，加载数据
watch(visible, async (val) => {
  if (val) {
    activeTabKey.value = 'menu'
    apiSearchText.value = ''
    menuSearchText.value = ''
    selectedMenuId.value = null
    await Promise.all([fetchMenuTree(), fetchAllApis()])
    await loadRolePermissions()
  }
})
</script>

<style scoped>
.permission-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.selected-count {
  color: #666;
  font-size: 13px;
}

/* 通用左右分栏布局 */
.menu-layout,
.api-layout {
  display: flex;
  height: calc(100vh - 280px);
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  overflow: hidden;
}

.menu-group-list,
.api-group-list {
  width: 240px;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
}

.group-list-header {
  padding: 12px 16px;
  font-weight: 500;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.group-list-content {
  flex: 1;
  overflow-y: auto;
}

.group-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  border-bottom: 1px solid #f5f5f5;
  transition: background 0.2s;
}

.group-item:hover {
  background: #f5f5f5;
}

.group-item.active {
  background: var(--app-primary-color-soft, rgba(0, 107, 230, 0.12));
  border-left: 3px solid var(--app-primary-color, #006be6);
}

.group-item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-item-total {
  color: #999;
  font-size: 12px;
}

/* 菜单右侧详情 */
.menu-detail-list,
.api-detail-list {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.menu-detail-header,
.api-detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  font-weight: 500;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.menu-detail-content,
.api-detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

/* 菜单子项块 */
.menu-child-block {
  margin-bottom: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  overflow: hidden;
}

.menu-child-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: #fafafa;
  font-weight: 500;
}

.menu-child-name {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.menu-child-name :deep(.ant-tag) {
  margin: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
}

.permission-code {
  font-size: 12px;
  color: #999;
  font-family: monospace;
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
}

.menu-child-items {
  padding: 8px 12px;
  background: #fff;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.btn-permission-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-permission-item:hover {
  background: #f5f5f5;
}

.btn-permission-content {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-permission-content :deep(.ant-tag) {
  margin: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
}

.btn-name {
  font-size: 13px;
}

.api-detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.api-item {
  padding: 8px 12px;
  border-radius: 4px;
  transition: background 0.2s;
}

.api-item:hover {
  background: #f5f5f5;
}

.api-item-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.api-item :deep(.ant-tag) {
  margin-right: 0;
  font-size: 10px;
  padding: 0 4px;
  line-height: 16px;
  min-width: 50px;
  text-align: center;
}

.api-path {
  font-size: 13px;
  color: #666;
  font-family: monospace;
}

.api-desc {
  font-size: 13px;
  color: #333;
}
</style>
