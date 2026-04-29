<template>
  <a-drawer
    :open="open"
    :title="`关联用户 - ${roleName}`"
    width="720"
    placement="right"
    destroyOnClose
    @close="emit('update:open', false)"
  >
    <a-table
      :columns="columns"
      :data-source="users"
      :loading="loading"
      :pagination="false"
      row-key="id"
      size="small"
      bordered
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'dept'">
          <span v-if="record.dept?.name">{{ record.dept.name }}</span>
          <a-tag v-else color="error">未绑定部门</a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </a-tag>
        </template>
      </template>
    </a-table>

    <div v-if="pagination.total > 0" class="drawer-pagination">
      <a-pagination
        :current="pagination.current"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        :show-size-changer="false"
        :show-total="total => `共 ${total} 条`"
        @change="handlePageChange"
      />
    </div>
  </a-drawer>
</template>

<script setup lang="ts">
import type { User } from '@/types'

interface DrawerPagination {
  current: number
  pageSize: number
  total: number
}

interface Props {
  open: boolean
  roleName: string
  loading: boolean
  users: User[]
  pagination: DrawerPagination
}

defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'page-change': [page: number]
}>()

const columns = [
  { title: '用户名', dataIndex: 'username', key: 'username', width: 140 },
  { title: '昵称', dataIndex: 'nickname', key: 'nickname', width: 140 },
  { title: '所属部门', key: 'dept' },
  { title: '状态', key: 'status', width: 100 }
]

const handlePageChange = (page: number) => {
  emit('page-change', page)
}
</script>

<style scoped>
.drawer-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
