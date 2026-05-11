<template>
  <a-drawer
    :open="open"
    title="对话详情"
    width="780"
    placement="right"
    @close="handleClose"
  >
    <template v-if="conversation">
      <a-descriptions :column="2" size="small" bordered class="conv-info">
        <a-descriptions-item label="对话标题" :span="2">{{ conversation.title }}</a-descriptions-item>
        <a-descriptions-item label="所属用户">
          {{ conversation.nickname || conversation.username || '#' + conversation.user_id }}
        </a-descriptions-item>
        <a-descriptions-item label="使用模型">{{ conversation.model || '默认模型' }}</a-descriptions-item>
        <a-descriptions-item label="消息总数">{{ conversation.message_count }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ conversation.created_at }}</a-descriptions-item>
        <a-descriptions-item label="最近更新" :span="2">{{ conversation.updated_at }}</a-descriptions-item>
      </a-descriptions>
    </template>

    <a-divider class="qa-divider">询问记录</a-divider>

    <div v-if="!loading && pairs.length === 0" class="empty-area">
      <a-empty description="该对话暂无消息" />
    </div>

    <div class="qa-list">
      <div v-for="(pair, idx) in pairs" :key="pair.key" class="qa-pair">
        <div class="qa-index">#{{ idx + 1 }}</div>

        <!-- 提问 -->
        <div v-if="pair.question" class="qa-block role-user">
          <div class="qa-block-header">
            <a-tag color="blue">提问</a-tag>
            <span class="qa-time">{{ pair.question.created_at }}</span>
          </div>
          <div class="qa-content">{{ pair.question.content }}</div>
        </div>
        <div v-else class="qa-block role-user empty">
          <a-tag color="default">无提问记录</a-tag>
        </div>

        <!-- 答案（默认折叠） -->
        <div v-if="pair.answer" class="qa-block role-assistant">
          <div class="qa-block-header">
            <a-tag color="green">答案</a-tag>
            <span class="qa-time">{{ pair.answer.created_at }}</span>
            <a-button type="link" size="small" @click="togglePair(pair.key)">
              {{ expandedPairs.has(pair.key) ? '收起答案' : '展开答案' }}
            </a-button>
          </div>
          <div v-if="expandedPairs.has(pair.key)" class="qa-content">{{ pair.answer.content }}</div>
          <div v-else class="qa-content collapsed">{{ truncateAnswer(pair.answer.content) }}</div>

          <!-- 思考过程也默认折叠 -->
          <div v-if="pair.answer.reasoning_content" class="qa-reasoning">
            <a-collapse ghost>
              <a-collapse-panel key="reasoning" header="思考过程">
                <pre class="reasoning-pre">{{ pair.answer.reasoning_content }}</pre>
              </a-collapse-panel>
            </a-collapse>
          </div>
        </div>
        <div v-else class="qa-block role-assistant empty">
          <a-tag color="default">未生成回答</a-tag>
        </div>

        <!-- 系统消息（如果有） -->
        <div v-if="pair.system" class="qa-block role-system">
          <div class="qa-block-header">
            <a-tag color="orange">系统</a-tag>
            <span class="qa-time">{{ pair.system.created_at }}</span>
          </div>
          <div class="qa-content">{{ pair.system.content }}</div>
        </div>
      </div>

      <div v-if="loading" class="load-more-area">
        <a-spin />
      </div>
      <div v-else-if="hasMore" class="load-more-area">
        <a-button size="small" @click="handleLoadMore">加载更多记录</a-button>
      </div>
      <div v-else-if="pairs.length > 0" class="load-more-area">
        <span class="no-more">— 已加载全部 —</span>
      </div>
    </div>
  </a-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getAdminMessages, type AdminConversationItem, type AIMessage } from '@/api/ai'

const props = defineProps<{
  open: boolean
  conversation: AdminConversationItem | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const PAGE_SIZE = 20
const ANSWER_PREVIEW_LIMIT = 120

const messages = ref<AIMessage[]>([])
const cursor = ref<number>(0)
const hasMore = ref<boolean>(false)
const loading = ref<boolean>(false)
const expandedPairs = ref<Set<string>>(new Set())

interface QAPair {
  key: string
  question: AIMessage | null
  answer: AIMessage | null
  system: AIMessage | null
}

const pairs = computed<QAPair[]>(() => {
  const result: QAPair[] = []
  let pendingUser: AIMessage | null = null

  const flushUserOnly = () => {
    if (pendingUser) {
      result.push({
        key: `pair-${pendingUser.id}`,
        question: pendingUser,
        answer: null,
        system: null,
      })
      pendingUser = null
    }
  }

  for (const msg of messages.value) {
    if (msg.role === 'user') {
      flushUserOnly()
      pendingUser = msg
    } else if (msg.role === 'assistant') {
      result.push({
        key: pendingUser ? `pair-${pendingUser.id}-${msg.id}` : `pair-orphan-${msg.id}`,
        question: pendingUser,
        answer: msg,
        system: null,
      })
      pendingUser = null
    } else {
      flushUserOnly()
      result.push({
        key: `system-${msg.id}`,
        question: null,
        answer: null,
        system: msg,
      })
    }
  }
  flushUserOnly()
  return result
})

const truncateAnswer = (content: string) => {
  if (!content) {
    return ''
  }
  if (content.length <= ANSWER_PREVIEW_LIMIT) {
    return content
  }
  return content.slice(0, ANSWER_PREVIEW_LIMIT) + '...'
}

const togglePair = (key: string) => {
  if (expandedPairs.value.has(key)) {
    expandedPairs.value.delete(key)
  } else {
    expandedPairs.value.add(key)
  }
  // 触发响应式更新
  expandedPairs.value = new Set(expandedPairs.value)
}

const fetchMessages = async (reset = true) => {
  if (!props.conversation || loading.value) {
    return
  }
  loading.value = true
  try {
    const params = {
      cursor: reset ? 0 : cursor.value,
      limit: PAGE_SIZE,
    }
    const res = await getAdminMessages(props.conversation.id, params)
    const data = res.data
    if (reset) {
      messages.value = data.list || []
      expandedPairs.value = new Set()
    } else {
      messages.value = [...messages.value, ...(data.list || [])]
    }
    cursor.value = data.next_cursor || 0
    hasMore.value = !!data.has_more
  } catch (err) {
    console.error('加载消息失败', err)
  } finally {
    loading.value = false
  }
}

const handleLoadMore = () => {
  if (!hasMore.value) {
    return
  }
  fetchMessages(false)
}

const handleClose = () => {
  emit('update:open', false)
}

watch(
  () => [props.open, props.conversation?.id],
  ([nextOpen]) => {
    if (nextOpen && props.conversation) {
      messages.value = []
      cursor.value = 0
      hasMore.value = false
      fetchMessages(true)
    }
  },
  { immediate: true }
)
</script>

<style scoped lang="less">
.conv-info {
  margin-bottom: 12px;
}

.qa-divider {
  margin: 16px 0 12px;
  font-weight: 600;
}

.empty-area {
  padding: 32px 0;
}

.qa-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.qa-pair {
  position: relative;
  padding: 12px 14px;
  border: 1px solid var(--ant-color-border-secondary, #f0f0f0);
  border-radius: 8px;
  background-color: var(--ant-color-bg-container, #fff);
}

.qa-index {
  position: absolute;
  top: 12px;
  right: 14px;
  font-size: 12px;
  color: var(--ant-color-text-tertiary, #00000040);
}

.qa-block {
  margin-top: 8px;
}

.qa-block:first-child {
  margin-top: 0;
}

.qa-block.empty {
  padding: 4px 0;
}

.qa-block-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.qa-time {
  font-size: 12px;
  color: var(--ant-color-text-tertiary, #00000073);
}

.qa-content {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 14px;
  line-height: 1.6;
  padding: 6px 8px;
  border-radius: 4px;
}

.role-user .qa-content {
  background-color: var(--ant-color-primary-bg, #e6f4ff);
}

.role-assistant .qa-content {
  background-color: var(--ant-color-success-bg, #f6ffed);
}

.role-system .qa-content {
  background-color: var(--ant-color-warning-bg, #fffbe6);
}

.qa-content.collapsed {
  color: var(--ant-color-text-secondary, #00000073);
  font-style: italic;
}

.qa-reasoning {
  margin-top: 6px;
}

.reasoning-pre {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  color: var(--ant-color-text-secondary, #00000073);
  margin: 0;
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
