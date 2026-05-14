<template>
  <PageWrapper class="dict-page">
    <AdminSplitLayout class="dict-page__layout" :aside-width="420" :content-min-width="940">
      <template #aside>
        <a-card class="dict-panel dict-panel--types">
          <template #title>
            <div class="dict-panel__title">
              <span>字典类型</span>
              <a-tag color="blue">{{ typePagination.total }}</a-tag>
            </div>
          </template>
          <template #extra>
            <a-space :size="4">
              <a-tooltip title="刷新列表">
                <a-button type="text" @click="fetchDictTypes()">
                  <ReloadOutlined />
                </a-button>
              </a-tooltip>
              <a-button type="primary" @click="handleAddType" v-permission="'system:dict:add'">
                <PlusOutlined />
                新增
              </a-button>
            </a-space>
          </template>

          <a-input-search
            v-model:value="typeSearchText"
            allow-clear
            class="dict-panel__search"
            placeholder="搜索字典名称或类型编码"
            @search="handleTypeSearch"
          />

          <a-segmented
            :value="typeStatusFilter"
            :options="statusFilterOptions"
            block
            class="dict-panel__segmented"
            @change="handleTypeStatusFilterChange"
          />

          <a-alert
            v-if="showSelectedOutsideFilter && selectedType"
            class="dict-panel__hint"
            type="info"
            show-icon
            :message="`当前搜索结果未包含已选字典：${selectedType.name}`"
          />

          <div class="dict-type-list">
            <a-skeleton v-if="typeLoading && !filteredDictTypes.length" active :paragraph="{ rows: 6 }" />
            <a-list
              v-else-if="filteredDictTypes.length"
              :data-source="filteredDictTypes"
              :loading="typeLoading"
              item-layout="horizontal"
            >
              <template #renderItem="{ item }">
                <a-list-item
                  :class="['dict-type-item', { 'dict-type-item--active': selectedType?.id === item.id }]"
                  role="button"
                  tabindex="0"
                  @click="handleSelectType(item)"
                  @keydown.enter.prevent="handleSelectType(item)"
                  @keydown.space.prevent="handleSelectType(item)"
                >
                  <template #actions>
                    <a-tooltip :title="item.status === 1 ? '点击停用' : '点击启用'">
                      <a-switch
                        size="small"
                        :checked="item.status === 1"
                        :loading="togglingTypeIds.has(item.id)"
                        @click.stop
                        @change="checked => handleToggleTypeStatus(item, Boolean(checked))"
                        v-permission="'system:dict:edit'"
                      />
                    </a-tooltip>
                    <a-tooltip title="复制编码">
                      <a-button type="text" size="small" @click.stop="handleCopy(item.type, '字典类型编码')">
                        <CopyOutlined />
                      </a-button>
                    </a-tooltip>
                    <a-tooltip title="编辑">
                      <a-button type="text" size="small" @click.stop="handleEditType(item)" v-permission="'system:dict:edit'">
                        <EditOutlined />
                      </a-button>
                    </a-tooltip>
                    <a-popconfirm
                      title="确定删除此字典类型及其所有字典数据吗？"
                      @confirm="handleDeleteType(item)"
                    >
                      <a-button type="text" size="small" danger @click.stop v-permission="'system:dict:delete'">
                        <DeleteOutlined />
                      </a-button>
                    </a-popconfirm>
                  </template>

                  <a-list-item-meta>
                    <template #title>
                      <div class="dict-type-item__title">
                        <span
                          class="dict-type-item__dot"
                          :class="{ 'dict-type-item__dot--off': item.status !== 1 }"
                        />
                        <span class="dict-type-item__name">{{ item.name }}</span>
                      </div>
                    </template>
                    <template #description>
                      <div class="dict-type-item__description">
                        <a-typography-text code>{{ item.type }}</a-typography-text>
                        <span class="dict-type-item__remark">{{ item.remark || '暂无备注' }}</span>
                      </div>
                    </template>
                  </a-list-item-meta>
                </a-list-item>
              </template>
            </a-list>

            <a-empty v-else-if="!typeLoading" :image="false" description="当前条件下没有匹配的字典类型" />
          </div>

          <div class="dict-panel__pagination">
            <a-pagination
              size="small"
              :current="typePagination.current"
              :page-size="typePagination.pageSize"
              :total="typePagination.total"
              :show-size-changer="typePagination.showSizeChanger"
              :show-total="typePagination.showTotal"
              @change="handleTypePaginationChange"
              @showSizeChange="handleTypePaginationChange"
            />
          </div>
        </a-card>
      </template>

      <div class="dict-page__main">
        <a-card class="dict-panel dict-panel--data">
          <template #title>
            <div class="dict-data-header">
              <div class="dict-data-header__main">
                {{ selectedType ? selectedType.name : '字典数据' }}
              </div>
              <div v-if="selectedType" class="dict-data-header__meta">
                <a-typography-text code>{{ selectedType.type }}</a-typography-text>
                <a-button type="link" size="small" @click="handleCopy(selectedType.type, '字典类型编码')">
                  复制编码
                </a-button>
                <a-tag :color="selectedType.status === 1 ? 'success' : 'default'">
                  {{ selectedType.status === 1 ? '正常' : '停用' }}
                </a-tag>
              </div>
            </div>
          </template>
          <template #extra>
            <a-space>
              <a-button :disabled="!selectedType" @click="fetchDictData()">
                <ReloadOutlined />
                刷新
              </a-button>
              <a-button type="primary" :disabled="!selectedType" @click="handleAddData" v-permission="'system:dict:add'">
                <PlusOutlined />
                新增字典数据
              </a-button>
            </a-space>
          </template>

          <template v-if="selectedType">
            <div class="dict-data-summary">
              <span class="dict-data-summary__label">备注</span>
              <span class="dict-data-summary__text">{{ selectedType.remark || '暂无备注' }}</span>
            </div>

            <div class="dict-data-toolbar">
              <a-input-search
                v-model:value="dataLabelSearch"
                allow-clear
                class="dict-data-toolbar__search"
                placeholder="搜索字典标签"
                @search="handleDataSearch"
              />
              <a-segmented
                :value="dataStatusFilter"
                :options="statusFilterOptions"
                @change="handleDataStatusFilterChange"
              />
              <a-popconfirm
                v-if="selectedDataKeys.length"
                :title="`确定删除选中的 ${selectedDataKeys.length} 条字典数据吗？`"
                @confirm="handleBatchDeleteData"
              >
                <a-button danger :loading="batchDeleting" v-permission="'system:dict:delete'">
                  <DeleteOutlined />
                  批量删除 ({{ selectedDataKeys.length }})
                </a-button>
              </a-popconfirm>
            </div>

            <a-table
              :columns="dataColumns"
              :data-source="dictDataList"
              :loading="dataLoading"
              :pagination="dataPagination"
              row-key="id"
              :row-selection="{ selectedRowKeys: selectedDataKeys, onChange: keys => (selectedDataKeys = keys as number[]) }"
              @change="pagination => handleDataPaginationChange(pagination.current, pagination.pageSize)"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'label'">
                  <a-tag :color="record.tag_type || 'default'">{{ record.label }}</a-tag>
                </template>
                <template v-if="column.key === 'value'">
                  <div class="dict-value-cell">
                    <a-typography-text>{{ record.value }}</a-typography-text>
                    <a-tooltip title="复制键值">
                      <a-button type="text" size="small" @click="handleCopy(record.value, '字典键值')">
                        <CopyOutlined />
                      </a-button>
                    </a-tooltip>
                  </div>
                </template>
                <template v-if="column.key === 'status'">
                  <a-switch
                    size="small"
                    :checked="record.status === 1"
                    :loading="togglingStatusIds.has(record.id)"
                    checked-children="开"
                    un-checked-children="关"
                    @change="checked => handleToggleDataStatus(record, Boolean(checked))"
                  />
                </template>
                <template v-if="column.key === 'is_default'">
                  <a-tag :color="record.is_default === 1 ? 'blue' : 'default'">
                    {{ record.is_default === 1 ? '默认' : '否' }}
                  </a-tag>
                </template>
                <template v-if="column.key === 'remark'">
                  <span>{{ record.remark || '-' }}</span>
                </template>
                <template v-if="column.key === 'action'">
                  <a-space :size="0">
                    <a-button type="link" size="small" @click="handleEditData(record)" v-permission="'system:dict:edit'">编辑</a-button>
                    <a-popconfirm title="确定删除此字典数据吗？" @confirm="handleDeleteData(record)">
                      <a-button type="link" size="small" danger v-permission="'system:dict:delete'">删除</a-button>
                    </a-popconfirm>
                  </a-space>
                </template>
              </template>
            </a-table>
          </template>

        <a-empty
          v-else
          class="dict-panel__empty"
          description="请先从左侧选择一个字典类型，再管理对应的字典数据"
        />
        </a-card>
      </div>
    </AdminSplitLayout>

    <DictTypeDrawer
      v-model:open="typeDrawerVisible"
      :title="typeDrawerTitle"
      :is-edit="Boolean(editingType)"
      :submitting="typeSubmitLoading"
      :initial-value="typeDrawerInitialValue"
      @submit="handleTypeSubmit"
    />

    <DictDataDrawer
      v-model:open="dataDrawerVisible"
      :title="dataDrawerTitle"
      :submitting="dataSubmitLoading"
      :current-type-name="selectedType?.name || '未选择'"
      :current-type-code="selectedType?.type || ''"
      :initial-value="dataDrawerInitialValue"
      @submit="handleDataSubmit"
    />
  </PageWrapper>
</template>

<script setup lang="ts">
import { CopyOutlined, DeleteOutlined, EditOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import AdminSplitLayout from '@/components/AdminSplitLayout.vue'
import PageWrapper from '@/components/page/PageWrapper.vue'
import DictDataDrawer from './components/DictDataDrawer.vue'
import DictTypeDrawer from './components/DictTypeDrawer.vue'
import { useDictPage } from './useDictPage'

const statusFilterOptions = [
  { label: '全部', value: 'all' },
  { label: '正常', value: 1 },
  { label: '停用', value: 0 },
]

const dataColumns = [
  { title: '字典标签', key: 'label', width: 160 },
  { title: '字典键值', key: 'value', width: 180 },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
  { title: '状态', key: 'status', width: 90 },
  { title: '默认', key: 'is_default', width: 90 },
  { title: '备注', key: 'remark' },
  { title: '操作', key: 'action', width: 140 },
]

const {
  batchDeleting,
  dataDrawerInitialValue,
  dataDrawerTitle,
  dataDrawerVisible,
  dataLabelSearch,
  dataLoading,
  dataPagination,
  dataStatusFilter,
  dataSubmitLoading,
  dictDataList,
  editingType,
  fetchDictData,
  fetchDictTypes,
  filteredDictTypes,
  handleAddData,
  handleAddType,
  handleBatchDeleteData,
  handleCopy,
  handleDataPaginationChange,
  handleDataSearch,
  handleDataStatusFilterChange,
  handleDataSubmit,
  handleDeleteData,
  handleDeleteType,
  handleEditData,
  handleEditType,
  handleSelectType,
  handleToggleDataStatus,
  handleToggleTypeStatus,
  handleTypePaginationChange,
  handleTypeSearch,
  handleTypeStatusFilterChange,
  handleTypeSubmit,
  selectedDataKeys,
  selectedType,
  showSelectedOutsideFilter,
  togglingStatusIds,
  togglingTypeIds,
  typeDrawerInitialValue,
  typeDrawerTitle,
  typeDrawerVisible,
  typeLoading,
  typePagination,
  typeSearchText,
  typeStatusFilter,
  typeSubmitLoading,
} = useDictPage()
</script>

<style scoped>
.dict-page__layout,
.dict-page__main {
  min-width: 0;
}

.dict-page {
  color: var(--app-text-color);
}

.dict-panel {
  min-height: 680px;
}

.dict-panel__title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dict-panel__search {
  margin-bottom: 12px;
}

.dict-panel__segmented {
  margin-bottom: 12px;
}

.dict-panel__hint {
  margin-bottom: 12px;
}

.dict-type-item__dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--app-success-color, #52c41a);
  box-shadow: 0 0 0 2px rgba(82, 196, 26, 0.16);
  flex-shrink: 0;
}

.dict-type-item__dot--off {
  background: var(--app-text-muted, #bfbfbf);
  box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.06);
}

.dict-data-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.dict-data-toolbar__search {
  width: 240px;
  max-width: 100%;
}

.dict-type-list {
  min-height: 520px;
}

.dict-type-item {
  padding: 12px 8px;
  cursor: pointer;
  border-radius: 10px;
  transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease;
  border: 1px solid transparent;
  color: var(--app-text-color);
}

.dict-type-item:hover {
  background: var(--app-hover-bg);
}

.dict-type-item--active {
  background: var(--app-primary-color-soft);
  border-color: var(--app-primary-color);
  box-shadow: inset 3px 0 0 var(--app-primary-color);
}

.dict-type-item__title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.dict-type-item__name {
  font-weight: 600;
  color: var(--app-text-strong);
}

.dict-type-item__description {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dict-type-item__remark {
  color: var(--app-text-secondary);
  line-height: 1.5;
}

.dict-panel__pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.dict-panel__empty {
  min-height: 560px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dict-data-header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dict-data-header__main {
  font-size: 16px;
  font-weight: 600;
  color: var(--app-text-strong);
}

.dict-data-header__meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.dict-data-summary {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 0 0 16px;
}

.dict-data-summary__label {
  color: var(--app-text-secondary);
  white-space: nowrap;
}

.dict-data-summary__text {
  color: var(--app-text-color);
  line-height: 1.5;
}

.dict-value-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.dict-page :deep(.ant-card) {
  border-color: var(--app-border-color);
  background: var(--app-surface-color);
}

.dict-page :deep(.ant-card-head) {
  border-bottom-color: var(--app-border-color);
}

.dict-page :deep(.ant-card-head-title),
.dict-page :deep(.ant-tag),
.dict-page :deep(.ant-list-item-action),
.dict-page :deep(.ant-typography),
.dict-page :deep(.ant-empty-description) {
  color: var(--app-text-color);
}

.dict-page :deep(.ant-alert) {
  border-color: var(--app-border-color);
}

.dict-page :deep(.ant-alert-message),
.dict-page :deep(.ant-alert-description) {
  color: var(--app-text-color);
}

.dict-page :deep(.ant-input-search .ant-input),
.dict-page :deep(.ant-input-search .ant-input-group-addon),
.dict-page :deep(.ant-input-affix-wrapper),
.dict-page :deep(.ant-input),
.dict-page :deep(.ant-pagination .ant-select-selector) {
  background: var(--app-surface-soft);
  border-color: var(--app-border-color);
  color: var(--app-text-color);
}

.dict-page :deep(.ant-input::placeholder),
.dict-page :deep(.ant-input-affix-wrapper input::placeholder),
.dict-page :deep(.ant-empty-description),
.dict-page :deep(.ant-pagination),
.dict-page :deep(.ant-typography-secondary) {
  color: var(--app-text-muted);
}

.dict-page :deep(.ant-list-item) {
  border-bottom-color: transparent;
}

.dict-page :deep(.ant-table) {
  background: transparent;
  color: var(--app-text-color);
}

.dict-page :deep(.ant-table-container) {
  border: 1px solid var(--app-border-color);
  border-radius: 10px;
  overflow: hidden;
}

.dict-page :deep(.ant-table-thead > tr > th) {
  background: var(--app-surface-soft);
  color: var(--app-text-strong);
  border-bottom-color: var(--app-border-color);
}

.dict-page :deep(.ant-table-tbody > tr > td) {
  background: var(--app-surface-color);
  border-bottom-color: var(--app-border-color);
  color: var(--app-text-color);
}

.dict-page :deep(.ant-table-tbody > tr:hover > td) {
  background: var(--app-hover-bg);
}

.dict-page :deep(.ant-table-placeholder > td) {
  background: var(--app-surface-color);
}

.dict-page :deep(.ant-btn-link) {
  color: var(--app-primary-color);
}

@media (max-width: 1199px) {
  .dict-panel {
    min-height: auto;
  }

  .dict-type-list {
    min-height: auto;
  }

  .dict-panel__empty {
    min-height: 240px;
  }
}
</style>
