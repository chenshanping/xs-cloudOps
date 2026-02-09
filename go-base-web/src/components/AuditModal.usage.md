# 审批弹窗组件使用说明

## 组件位置
`src/components/AuditModal.vue`

## 在生成的页面中使用审批功能

### 1. 导入组件和API

在生成的 index.vue 文件中添加导入：

```vue
<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import AuditModal from '@/components/AuditModal.vue'
import { audit{{ModelName}} } from '@/api/{{module_name}}'

// 审批弹窗状态
const auditModalVisible = ref(false)
const currentAuditId = ref<number>()

// 打开审批弹窗
const handleAudit = (record: any) => {
  currentAuditId.value = record.id
  auditModalVisible.value = true
}

// 确认审批
const handleAuditConfirm = async (data: { audit_status: number; audit_remark: string }) => {
  try {
    await audit{{ModelName}}(currentAuditId.value!, data)
    message.success('审批成功')
    auditModalVisible.value = false
    fetchData() // 刷新列表
  } catch (error: any) {
    message.error(error.response?.data?.msg || '审批失败')
  }
}
</script>
```

### 2. 在模板中添加组件

```vue
<template>
  <!-- 表格操作列中添加审批按钮 -->
  <a-table>
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'action'">
        <a-space>
          <!-- 其他按钮... -->
          <a-button 
            v-if="record.audit_status === 0"
            type="link" 
            size="small" 
            @click="handleAudit(record)"
          >
            审批
          </a-button>
        </a-space>
      </template>
    </template>
  </a-table>

  <!-- 审批弹窗 -->
  <AuditModal
    v-model:open="auditModalVisible"
    title="审批"
    @confirm="handleAuditConfirm"
  />
</template>
```

### 3. 添加审批状态列

在表格列定义中添加审批状态列：

```typescript
const columns = [
  // ... 其他列
  {
    title: '审批状态',
    dataIndex: 'audit_status',
    key: 'audit_status',
    width: 100
  }
]
```

在模板中自定义渲染审批状态：

```vue
<template #bodyCell="{ column, record }">
  <template v-if="column.key === 'audit_status'">
    <a-tag v-if="record.audit_status === 0" color="default">待审批</a-tag>
    <a-tag v-else-if="record.audit_status === 1" color="success">审批通过</a-tag>
    <a-tag v-else-if="record.audit_status === 2" color="error">审批拒绝</a-tag>
  </template>
</template>
```

### 4. 在详情/表单中显示审批信息

```vue
<template>
  <!-- 审批信息展示 -->
  <a-descriptions v-if="detailData?.audit_status" title="审批信息" bordered>
    <a-descriptions-item label="审批状态">
      <a-tag v-if="detailData.audit_status === 0" color="default">待审批</a-tag>
      <a-tag v-else-if="detailData.audit_status === 1" color="success">审批通过</a-tag>
      <a-tag v-else-if="detailData.audit_status === 2" color="error">审批拒绝</a-tag>
    </a-descriptions-item>
    <a-descriptions-item label="审批备注">
      {{ detailData.audit_remark || '-' }}
    </a-descriptions-item>
    <a-descriptions-item label="审批时间">
      {{ detailData.audit_time || '-' }}
    </a-descriptions-item>
    <a-descriptions-item label="审批人">
      {{ detailData.auditor?.nickname || detailData.auditor?.username || '-' }}
    </a-descriptions-item>
  </a-descriptions>
</template>
```

## 完整示例

参考下面的完整示例代码片段：

```vue
<template>
  <div class="page-container">
    <a-card>
      <!-- 搜索表单 -->
      <a-form>
        <!-- ... -->
      </a-form>

      <!-- 数据表格 -->
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
      >
        <template #bodyCell="{ column, record }">
          <!-- 审批状态列 -->
          <template v-if="column.key === 'audit_status'">
            <a-tag v-if="record.audit_status === 0" color="default">待审批</a-tag>
            <a-tag v-else-if="record.audit_status === 1" color="success">审批通过</a-tag>
            <a-tag v-else-if="record.audit_status === 2" color="error">审批拒绝</a-tag>
          </template>

          <!-- 操作列 -->
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="handleView(record)">查看</a-button>
              <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-button 
                v-if="record.audit_status === 0"
                type="link" 
                size="small" 
                @click="handleAudit(record)"
              >
                审批
              </a-button>
              <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record.id)">
                <a-button type="link" size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 审批弹窗 -->
    <AuditModal
      v-model:open="auditModalVisible"
      title="审批"
      @confirm="handleAuditConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import AuditModal from '@/components/AuditModal.vue'
import { get{{ModelName}}List, audit{{ModelName}} } from '@/api/{{module_name}}'

// 表格列定义
const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  // ... 其他列
  { title: '审批状态', dataIndex: 'audit_status', key: 'audit_status', width: 100 },
  { title: '操作', key: 'action', width: 200, fixed: 'right' }
]

const dataSource = ref([])
const loading = ref(false)

// 审批相关
const auditModalVisible = ref(false)
const currentAuditId = ref<number>()

const handleAudit = (record: any) => {
  currentAuditId.value = record.id
  auditModalVisible.value = true
}

const handleAuditConfirm = async (data: { audit_status: number; audit_remark: string }) => {
  try {
    await audit{{ModelName}}(currentAuditId.value!, data)
    message.success('审批成功')
    auditModalVisible.value = false
    fetchData()
  } catch (error: any) {
    message.error(error.response?.data?.msg || '审批失败')
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await get{{ModelName}}List({ page: 1, page_size: 10 })
    dataSource.value = res.data.list
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
```

## 审批状态说明

- **0** - 待审批（灰色 default） - 用户提交表单后的默认状态
- **1** - 审批通过（绿色 success）
- **2** - 审批拒绝（红色 error）

## 审批流程

1. 用户创建/提交表单 → `audit_status` 自动设置为 0（待审批）
2. 审批人点击“审批”按钮
3. 在弹窗中选择“通过”或“拒绝”，并填写审批备注
4. 确认后，后端更新审批状态、审批人、审批时间和审批备注

## 注意事项

1. 审批按钮只在 `audit_status === 0` 时显示
2. 已审批的数据不可再次审批（后端会校验）
3. 审批弹窗只有“通过”和“拒绝”两个选项，不需要选择“待审批”
4. 审批备注为必填项，最多500字
5. 审批后会自动记录审批人、审批时间和审批备注
