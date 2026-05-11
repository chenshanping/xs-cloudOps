<template>
  <div class="ai-history-page">
    <!-- 左侧：用户列表 -->
    <a-card :bordered="false" class="left-pane" body-style="padding: 0; display: flex; flex-direction: column; height: 100%;">
      <div class="filter-bar">
        <a-input-search
          v-model:value="userKeyword"
          placeholder="搜索用户名/昵称"
          allow-clear
          @search="handleUserSearch"
          @change="handleUserKeywordChange"
        />
      </div>

      <div class="user-list">
        <a-empty v-if="!loadingUsers && users.length === 0" description="暂无活跃用户" />

        <div
          v-for="user in users"
          :key="user.id"
          :class="['user-item', { active: selectedUser?.id === user.id }]"
          @click="handleSelectUser(user)"
        >
          <div class="user-name">
            <span class="user-display">{{ user.nickname || user.username }}</span>
            <span v-if="user.nickname" class="user-username">@{{ user.username }}</span>
          </div>
          <div class="user-meta">
            <a-tag color="blue">{{ user.conversation_count }} 个对话</a-tag>
          </div>
          <div class="user-active">最近活动：{{ user.last_active_at }}</div>
        </div>

        <div v-if="loadingUsers" class="load-more-area">
          <a-spin />
        </div>
      </div>

      <div v-if="userTotal > 0" class="user-pagination">
        <a-pagination
          v-model:current="userPage"
          v-model:page-size="userPageSize"
          :total="userTotal"
          :page-size-options="['20', '50', '100']"
          show-size-changer
          size="small"
          simple
          @change="fetchUsers"
        />
      </div>
    </a-card>

    <!-- 右侧：对话列表 -->
    <a-card :bordered="false" class="right-pane" body-style="padding: 0; display: flex; flex-direction: column; height: 100%;">
      <div v-if="!selectedUser" class="empty-right">
        <a-empty description="请从左侧选择一个用户查看其对话" />
      </div>

      <template v-else>
        <div class="convs-header">
          <div class="convs-title">
            {{ selectedUser.nickname || selectedUser.username }} 的对话记录
          </div>
          <div class="convs-toolbar">
            <a-input-search
              v-model:value="conversationKeyword"
              placeholder="按对话标题搜索"
              allow-clear
              style="width: 220px"
              @search="handleConversationSearch"
            />
          </div>
        </div>

        <div class="conversation-list" @scroll="handleConversationScroll">
          <a-empty v-if="!loadingConversations && conversations.length === 0" description="该用户暂无对话" />

          <div
            v-for="conv in conversations"
            :key="conv.id"
            class="conversation-item"
            @click="handleOpenDetail(conv)"
          >
            <div class="conv-header">
              <span class="conv-title" :title="conv.title">{{ conv.title }}</span>
              <a-space :size="6">
                <a-button
                  v-if="userStore.hasPermission('ai:history:view')"
                  type="link"
                  size="small"
                  @click.stop="handleOpenDetail(conv)"
                >
                  查看详情
                </a-button>
                <a-popconfirm
                  v-if="userStore.hasPermission('ai:history:delete')"
                  title="确定删除该对话及其所有消息吗？"
                  ok-text="删除"
                  ok-type="danger"
                  cancel-text="取消"
                  @confirm.stop="handleDeleteConversation(conv)"
                >
                  <a-button type="link" size="small" danger @click.stop>
                    <DeleteOutlined />
                  </a-button>
                </a-popconfirm>
              </a-space>
            </div>
            <div class="conv-meta">
              <a-tag>{{ conv.model || '默认模型' }}</a-tag>
              <span class="meta-text">{{ conv.message_count }} 条消息</span>
              <span class="meta-text">最近更新 {{ conv.updated_at }}</span>
            </div>
          </div>

          <div v-if="loadingConversations" class="load-more-area">
            <a-spin />
          </div>
          <div v-else-if="conversationsHasMore" class="load-more-area">
            <a-button size="small" @click="loadMoreConversations">加载更多</a-button>
          </div>
          <div v-else-if="conversations.length > 0" class="load-more-area">
            <span class="no-more">— 已加载全部 —</span>
          </div>
        </div>
      </template>
    </a-card>

    <ConversationDetailDrawer
      v-model:open="drawerOpen"
      :conversation="drawerConversation"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { DeleteOutlined } from '@ant-design/icons-vue'
import {
  getAdminAIUsers,
  getAdminConversations,
  deleteAdminConversation,
  type AdminAIUserItem,
  type AdminConversationItem,
} from '@/api/ai'
import { useUserStore } from '@/store/user'
import ConversationDetailDrawer from './components/ConversationDetailDrawer.vue'

const userStore = useUserStore()

const CONV_PAGE_SIZE = 20

// ===== 用户列表（左侧） =====
const userKeyword = ref<string>('')
const users = ref<AdminAIUserItem[]>([])
const userTotal = ref<number>(0)
const userPage = ref<number>(1)
const userPageSize = ref<number>(50)
const loadingUsers = ref<boolean>(false)
const selectedUser = ref<AdminAIUserItem | null>(null)

let userSearchDebounce: ReturnType<typeof setTimeout> | null = null

const fetchUsers = async () => {
  if (loadingUsers.value) {
    return
  }
  loadingUsers.value = true
  try {
    const res = await getAdminAIUsers({
      page: userPage.value,
      page_size: userPageSize.value,
      keyword: userKeyword.value?.trim() || undefined,
    })
    users.value = res.data.list || []
    userTotal.value = res.data.total || 0
  } catch (err) {
    console.error('加载用户列表失败', err)
  } finally {
    loadingUsers.value = false
  }
}

const handleUserSearch = () => {
  userPage.value = 1
  fetchUsers()
}

const handleUserKeywordChange = () => {
  if (userSearchDebounce) {
    clearTimeout(userSearchDebounce)
  }
  userSearchDebounce = setTimeout(() => {
    userPage.value = 1
    fetchUsers()
  }, 300)
}

const handleSelectUser = (user: AdminAIUserItem) => {
  selectedUser.value = user
  conversationKeyword.value = ''
  fetchConversations(true)
}

// ===== 对话列表（右侧） =====
const conversationKeyword = ref<string>('')
const conversations = ref<AdminConversationItem[]>([])
const conversationsCursor = ref<number>(0)
const conversationsHasMore = ref<boolean>(false)
const loadingConversations = ref<boolean>(false)

const fetchConversations = async (reset = true) => {
  if (!selectedUser.value || loadingConversations.value) {
    return
  }
  loadingConversations.value = true
  try {
    const res = await getAdminConversations({
      cursor: reset ? 0 : conversationsCursor.value,
      limit: CONV_PAGE_SIZE,
      user_id: selectedUser.value.id,
      keyword: conversationKeyword.value?.trim() || undefined,
    })
    const data = res.data
    if (reset) {
      conversations.value = data.list || []
    } else {
      conversations.value = [...conversations.value, ...(data.list || [])]
    }
    conversationsCursor.value = data.next_cursor || 0
    conversationsHasMore.value = !!data.has_more
  } catch (err) {
    console.error('加载对话列表失败', err)
  } finally {
    loadingConversations.value = false
  }
}

const loadMoreConversations = () => {
  if (!conversationsHasMore.value) {
    return
  }
  fetchConversations(false)
}

const handleConversationSearch = () => {
  fetchConversations(true)
}

const handleConversationScroll = (e: Event) => {
  const target = e.target as HTMLElement
  if (
    !loadingConversations.value &&
    conversationsHasMore.value &&
    target.scrollTop + target.clientHeight >= target.scrollHeight - 24
  ) {
    loadMoreConversations()
  }
}

// ===== 详情 Drawer =====
const drawerOpen = ref<boolean>(false)
const drawerConversation = ref<AdminConversationItem | null>(null)

const handleOpenDetail = (conv: AdminConversationItem) => {
  if (!userStore.hasPermission('ai:history:view')) {
    message.warning('无权查看对话详情')
    return
  }
  drawerConversation.value = conv
  drawerOpen.value = true
}

const handleDeleteConversation = async (conv: AdminConversationItem) => {
  try {
    await deleteAdminConversation(conv.id)
    message.success('删除成功')
    if (drawerConversation.value?.id === conv.id) {
      drawerOpen.value = false
      drawerConversation.value = null
    }
    // 刷新对话列表 + 用户列表（对话数会变）
    fetchConversations(true)
    fetchUsers()
  } catch {
    // 错误已由 request 拦截器处理
  }
}

watch(userPageSize, () => {
  userPage.value = 1
  fetchUsers()
})

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="less">
.ai-history-page {
  display: flex;
  gap: 12px;
  height: calc(100vh - 140px);
  min-height: 480px;
}

.left-pane {
  width: 320px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}

.right-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.filter-bar {
  padding: 12px;
  border-bottom: 1px solid var(--ant-color-border-secondary, #f0f0f0);
}

.user-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.user-item {
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 4px;
}

.user-item:hover {
  background-color: var(--ant-color-fill-tertiary, #f5f5f5);
}

.user-item.active {
  background-color: var(--ant-color-primary-bg, #e6f4ff);
}

.user-name {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-wrap: wrap;
}

.user-display {
  font-weight: 500;
}

.user-username {
  font-size: 12px;
  color: var(--ant-color-text-secondary, #00000073);
}

.user-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
}

.user-active {
  margin-top: 4px;
  font-size: 12px;
  color: var(--ant-color-text-tertiary, #00000040);
}

.user-pagination {
  padding: 8px 12px;
  border-top: 1px solid var(--ant-color-border-secondary, #f0f0f0);
  display: flex;
  justify-content: center;
}

.empty-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.convs-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--ant-color-border-secondary, #f0f0f0);
  flex-wrap: wrap;
}

.convs-title {
  font-size: 16px;
  font-weight: 600;
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.conversation-item {
  padding: 12px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 8px;
  border: 1px solid var(--ant-color-border-secondary, #f0f0f0);
}

.conversation-item:hover {
  background-color: var(--ant-color-fill-tertiary, #fafafa);
  border-color: var(--ant-color-primary-border, #91caff);
}

.conv-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.conv-title {
  flex: 1;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conv-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 6px;
}

.meta-text {
  font-size: 12px;
  color: var(--ant-color-text-secondary, #00000073);
}

.load-more-area {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 0;
}

.no-more {
  font-size: 12px;
  color: var(--ant-color-text-tertiary, #00000040);
}
</style>
