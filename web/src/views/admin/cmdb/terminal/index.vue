<template>
  <div class="cmdb-terminal-standalone">
    <PageWrapper
      class="cmdb-terminal-page"
      title="SSH 在线终端"
      description="通过主机已绑定 SSH 凭据直接进入当前主机终端，并保留完整会话审计。"
    >
      <div class="terminal-page">
      <div class="host-overview">
        <div class="host-overview__main">
          <div class="host-overview__title">
            <CodeOutlined />
            <div>
              <strong>{{ hostDetail?.name || '主机终端' }}</strong>
              <span>{{ hostDetail ? `${hostDetail.ssh_host}:${hostDetail.ssh_port}` : '正在加载主机信息' }}</span>
            </div>
          </div>
          <div class="host-overview__meta">
            <a-tag v-if="hostDetail?.group?.name" color="blue">{{ hostDetail.group.name }}</a-tag>
            <a-tag v-if="hostDetail?.environment">{{ hostDetail.environment }}</a-tag>
            <a-tag v-if="hostDetail?.credential_summary?.name" color="purple">{{ hostDetail.credential_summary.name }}</a-tag>
          </div>
        </div>
        <a-space wrap>
          <a-button
            type="primary"
            :loading="connecting"
            :disabled="!hostId || connecting"
            v-permission="'cmdb:terminal:connect'"
            @click="handleReconnect"
          >
            重新连接
          </a-button>
          <a-button
            :disabled="!currentSession || currentSession.status !== 'active'"
            v-permission="'cmdb:terminal:disconnect'"
            @click="handleDisconnect"
          >
            断开连接
          </a-button>
          <a-button
            :disabled="!currentSession"
            v-permission="'cmdb:terminal:audit'"
            @click="handleOpenLogs"
          >
            查看日志
          </a-button>
          <a-button
            danger
            :disabled="!canForceDisconnectCurrent"
            v-permission="'cmdb:terminal:force_disconnect'"
            @click="handleForceDisconnect"
          >
            强制断开
          </a-button>
        </a-space>
      </div>

      <div class="terminal-layout">
        <div class="terminal-workbench">
          <div class="terminal-surface" @click="focusTerminal">
            <div ref="terminalRef" class="terminal-surface__canvas" />
            <div v-if="!isTerminalOpen" class="terminal-placeholder">
              <CodeOutlined />
              <strong>{{ placeholderTitle }}</strong>
              <span>{{ placeholderDescription }}</span>
            </div>
          </div>
        </div>

        <div class="meta-stack">
          <div class="meta-card">
            <span>当前状态</span>
            <strong>{{ currentSession ? getSessionStatusMeta(currentSession.status).text : '未连接' }}</strong>
            <p>{{ currentSession?.disconnect_reason || statusDescription }}</p>
          </div>
          <div class="meta-card">
            <span>主机指纹</span>
            <strong>{{ currentSession?.host_key_fingerprint || '-' }}</strong>
            <p>首次连接记录指纹，后续连接严格校验。</p>
          </div>
          <div class="meta-card">
            <span>会话信息</span>
            <strong>{{ currentSession?.username_snapshot || hostDetail?.credential_summary?.username || '-' }}</strong>
            <p>
              开始：{{ formatTime(currentSession?.start_time || currentSession?.created_at) }}
              <br />
              最近活动：{{ formatTime(currentSession?.last_activity_at) }}
            </p>
          </div>
          <div class="meta-card">
            <span>主机信息</span>
            <strong>{{ hostDetail?.platform_version || hostDetail?.platform || hostDetail?.os || '-' }}</strong>
            <p>
              主机名：{{ hostDetail?.hostname || '-' }}
              <br />
              内核：{{ hostDetail?.kernel_version || '-' }}
            </p>
          </div>
        </div>
      </div>

        <TerminalLogDrawer
          v-model:open="logDrawerVisible"
          :loading="logLoading"
          :session="currentSession"
          :logs="logList"
          :pagination="logPagination"
          :stream-type="logFilter.streamType"
          @change="handleLogChange"
          @refresh="fetchLogs"
        />
      </div>
    </PageWrapper>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { CodeOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { createVNode } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import PageWrapper from '@/components/page/PageWrapper.vue'
import {
  createCmdbTerminalSession,
  disconnectCmdbTerminalSession,
  forceDisconnectCmdbTerminalSession,
  getCmdbHost,
  getCmdbTerminalSession,
  getCmdbTerminalSessionLogs,
  type CmdbHostItem,
  type CmdbTerminalConnectPayload,
  type CmdbTerminalLogItem,
  type CmdbTerminalLogStreamType,
  type CmdbTerminalSessionItem,
  type CmdbTerminalSessionStatus,
} from '@/api/cmdb'
import { formatTime } from '@/utils/format'
import { usePermission } from '@/utils/permission'
import TerminalLogDrawer from './components/TerminalLogDrawer.vue'

interface TerminalServerMessage {
  type?: string
  data?: string
  status?: string
  message?: string
  session_id?: number
}

interface TerminalClientMessage {
  type: 'input' | 'resize' | 'disconnect'
  data?: string
  cols?: number
  rows?: number
}

const route = useRoute()
const { hasPermission } = usePermission()

const hostId = computed(() => Number(route.params.hostId || 0))
const terminalRef = ref<HTMLElement>()
const hostDetail = ref<CmdbHostItem | null>(null)
const currentSession = ref<CmdbTerminalSessionItem | null>(null)
const connecting = ref(false)
const logDrawerVisible = ref(false)
const logLoading = ref(false)
const logList = ref<CmdbTerminalLogItem[]>([])

const logPagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
})

const logFilter = reactive<{
  streamType?: CmdbTerminalLogStreamType
}>({
  streamType: undefined,
})

let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let socket: WebSocket | null = null
let resizeTimer: number | null = null
let lastSessionID = 0

const isTerminalOpen = ref(false)
const canForceDisconnectCurrent = computed(
  () =>
    !!currentSession.value &&
    currentSession.value.status === 'active' &&
    hasPermission('cmdb:terminal:force_disconnect')
)

const statusDescription = computed(() => {
  if (connecting.value) {
    return '系统正在校验指纹、建立 SSH 会话，请稍候。'
  }
  if (hostDetail.value?.credential_summary?.name) {
    return `当前主机将使用凭据「${hostDetail.value.credential_summary.name}」连接。`
  }
  return '当前无活动终端连接。'
})

const placeholderTitle = computed(() => {
  if (!hostId.value) {
    return '主机参数无效'
  }
  if (connecting.value) {
    return '正在建立终端连接'
  }
  if (currentSession.value?.status === 'failed') {
    return '终端连接失败'
  }
  if (currentSession.value?.status === 'closed') {
    return '终端会话已结束'
  }
  return '准备连接当前主机'
})

const placeholderDescription = computed(() => {
  if (!hostId.value) {
    return '缺少有效主机 ID，无法创建终端会话。'
  }
  if (currentSession.value?.disconnect_reason) {
    return currentSession.value.disconnect_reason
  }
  if (connecting.value) {
    return '终端页打开后会自动连接当前主机，连接失败会保留在本页并显示明确错误。'
  }
  return '点击“重新连接”可再次创建新会话，不会复用历史会话。'
})

const getSessionStatusMeta = (status?: CmdbTerminalSessionStatus) => {
  switch (status) {
    case 'prepared':
      return { text: '准备中', badge: 'processing' as const }
    case 'active':
      return { text: '在线中', badge: 'green' as const }
    case 'failed':
      return { text: '连接失败', badge: 'red' as const }
    case 'closed':
      return { text: '已关闭', badge: 'default' as const }
    default:
      return { text: '未连接', badge: 'default' as const }
  }
}

const ensureTerminal = async () => {
  if (!terminalRef.value || terminal) {
    return
  }
  terminal = new Terminal({
    cursorBlink: true,
    fontSize: 13,
    fontFamily: 'Consolas, Monaco, monospace',
    theme: {
      background: '#0f172a',
      foreground: '#e2e8f0',
      cursor: '#1677ff',
      selectionBackground: 'rgba(22, 119, 255, 0.28)',
    },
    convertEol: true,
    disableStdin: false,
  })
  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.open(terminalRef.value)
  fitAddon.fit()
  terminal.focus()
  terminal.onData(data => {
    sendMessage({ type: 'input', data })
  })
}

const clearTerminal = () => {
  terminal?.clear()
  terminal?.write('\x1bc')
}

const focusTerminal = () => {
  terminal?.focus()
}

const writeTerminalLine = (content: string) => {
  terminal?.write(content)
}

const sendMessage = (payload: TerminalClientMessage) => {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    return
  }
  socket.send(JSON.stringify(payload))
}

const syncTerminalResize = () => {
  if (!terminal || !fitAddon || !socket || socket.readyState !== WebSocket.OPEN) {
    return
  }
  fitAddon.fit()
  sendMessage({
    type: 'resize',
    cols: terminal.cols,
    rows: terminal.rows,
  })
}

const debounceResize = () => {
  if (resizeTimer) {
    window.clearTimeout(resizeTimer)
  }
  resizeTimer = window.setTimeout(() => {
    syncTerminalResize()
  }, 120)
}

const closeSocket = () => {
  if (socket) {
    try {
      socket.close()
    } catch {
      // noop
    }
  }
  socket = null
  isTerminalOpen.value = false
}

const fetchHost = async () => {
  if (!hostId.value) {
    throw new Error('主机参数无效')
  }
  const res = await getCmdbHost(hostId.value)
  hostDetail.value = res.data
}

const fetchSessionDetail = async (sessionID: number) => {
  const res = await getCmdbTerminalSession(sessionID)
  currentSession.value = res.data
  return res.data
}

const fetchLogs = async () => {
  if (!currentSession.value) {
    logList.value = []
    logPagination.total = 0
    return
  }
  logLoading.value = true
  try {
    const res = await getCmdbTerminalSessionLogs(currentSession.value.id, {
      page: logPagination.current,
      page_size: logPagination.pageSize,
      stream_type: logFilter.streamType,
    })
    logList.value = res.data.list || []
    logPagination.total = res.data.total || 0
  } finally {
    logLoading.value = false
  }
}

const connectWebSocket = async (payload: CmdbTerminalConnectPayload) => {
  await ensureTerminal()
  clearTerminal()
  const baseProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsURL = payload.ws_url.startsWith('ws')
    ? payload.ws_url
    : `${baseProtocol}//${window.location.host}${payload.ws_url.startsWith('/') ? payload.ws_url : `/${payload.ws_url}`}`

  return new Promise<void>((resolve, reject) => {
    const ws = new WebSocket(wsURL)
    socket = ws
    isTerminalOpen.value = false
    let settled = false

    ws.onopen = () => {
      isTerminalOpen.value = true
      syncTerminalResize()
      focusTerminal()
    }

    ws.onmessage = async event => {
      const data = typeof event.data === 'string' ? event.data : ''
      let parsed: TerminalServerMessage | null = null
      try {
        parsed = JSON.parse(data)
      } catch {
        writeTerminalLine(data)
        return
      }
      if (!parsed) {
        return
      }
      if (parsed.type === 'output') {
        writeTerminalLine(parsed.data || '')
        return
      }
      if (parsed.type === 'status') {
        if (parsed.message) {
          writeTerminalLine(`\r\n[系统提示] ${parsed.message}\r\n`)
        }
        if (parsed.status === 'connected' && payload.session_id) {
          lastSessionID = payload.session_id
          await fetchSessionDetail(payload.session_id)
          if (!settled) {
            settled = true
            resolve()
          }
        }
        if (parsed.status === 'error' && !settled) {
          settled = true
          reject(new Error(parsed.message || '终端连接失败'))
        }
      }
    }

    ws.onerror = () => {
      isTerminalOpen.value = false
      if (!settled) {
        settled = true
        reject(new Error('终端 WebSocket 连接失败'))
      }
    }

    ws.onclose = async () => {
      socket = null
      isTerminalOpen.value = false
      connecting.value = false
      if (lastSessionID) {
        try {
          await fetchSessionDetail(lastSessionID)
        } catch {
          // noop
        }
      }
      if (!settled) {
        settled = true
        reject(new Error('终端连接已关闭'))
      }
    }
  })
}

const createAndConnect = async () => {
  if (!hostId.value) {
    throw new Error('主机参数无效')
  }
  const res = await createCmdbTerminalSession({ host_id: hostId.value })
  await connectWebSocket(res.data)
}

const handleReconnect = async () => {
  connecting.value = true
  closeSocket()
  currentSession.value = null
  try {
    await fetchHost()
    await createAndConnect()
    message.success('终端连接成功')
  } catch (error: any) {
    closeSocket()
    if (error?.message === '主机指纹异常') {
      message.error('主机指纹异常，已拒绝连接，请联系管理员确认指纹变更')
    } else if (error?.message) {
      message.error(error.message)
    }
  } finally {
    connecting.value = false
  }
}

const handleDisconnect = async () => {
  if (!currentSession.value) {
    return
  }
  await disconnectCmdbTerminalSession(currentSession.value.id)
  sendMessage({ type: 'disconnect' })
  closeSocket()
  await fetchSessionDetail(currentSession.value.id)
  message.success('终端会话已断开')
}

const handleForceDisconnect = () => {
  if (!currentSession.value) {
    return
  }
  Modal.confirm({
    title: '确认强制断开',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要强制断开当前会话 #${currentSession.value.id} 吗？此操作会立即结束目标用户的 SSH 连接。`,
    okText: '强制断开',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      if (!currentSession.value) {
        return
      }
      await forceDisconnectCmdbTerminalSession(currentSession.value.id)
      closeSocket()
      await fetchSessionDetail(currentSession.value.id)
      message.success('已强制断开终端会话')
    },
  })
}

const handleOpenLogs = async () => {
  if (!currentSession.value) {
    message.warning('当前没有可查看日志的会话')
    return
  }
  logDrawerVisible.value = true
  logPagination.current = 1
  await fetchLogs()
}

const handleLogChange = async (payload: { page: number; pageSize: number; streamType?: CmdbTerminalLogStreamType }) => {
  logPagination.current = payload.page
  logPagination.pageSize = payload.pageSize
  logFilter.streamType = payload.streamType
  await fetchLogs()
}

onMounted(async () => {
  await nextTick()
  await ensureTerminal()
  window.addEventListener('resize', debounceResize)
  await handleReconnect()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', debounceResize)
  if (resizeTimer) {
    window.clearTimeout(resizeTimer)
  }
  closeSocket()
  terminal?.dispose()
  terminal = null
  fitAddon = null
})
</script>

<style scoped>
.cmdb-terminal-standalone {
  min-height: 100vh;
  padding: 20px;
  background:
    radial-gradient(circle at top left, rgba(22, 119, 255, 0.12), transparent 32%),
    linear-gradient(180deg, #f5f7fa 0%, #eef2f6 100%);
}

.cmdb-terminal-page {
  max-width: 1600px;
  height: 100%;
  margin: 0 auto;
}

.terminal-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.host-overview {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.host-overview__main {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

.host-overview__title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.host-overview__title :deep(.anticon) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  color: #1677ff;
  background: rgba(22, 119, 255, 0.1);
  border-radius: 10px;
}

.host-overview__title strong {
  display: block;
  color: var(--app-text-strong);
  font-size: 18px;
}

.host-overview__title span {
  color: var(--app-text-muted);
  font-size: 12px;
}

.host-overview__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.terminal-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 16px;
  min-height: 0;
}

.terminal-workbench,
.meta-stack {
  min-height: 0;
}

.terminal-surface {
  position: relative;
  min-height: 620px;
  padding: 12px;
  overflow: hidden;
  background: #0f172a;
  border: 1px solid rgba(148, 163, 184, 0.16);
  border-radius: 10px;
}

.terminal-surface__canvas {
  width: 100%;
  height: 100%;
}

.terminal-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: rgba(226, 232, 240, 0.82);
  text-align: center;
  background: rgba(15, 23, 42, 0.82);
}

.terminal-placeholder :deep(.anticon) {
  font-size: 28px;
}

.terminal-placeholder strong {
  font-size: 16px;
}

.terminal-placeholder span {
  width: min(420px, 80%);
  font-size: 12px;
  line-height: 1.7;
}

.meta-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.meta-card {
  padding: 14px;
  background: var(--app-surface-color);
  border: 1px solid var(--app-border-color);
  border-radius: 8px;
}

.meta-card span {
  color: var(--app-text-muted);
  font-size: 12px;
}

.meta-card strong {
  display: block;
  margin-top: 8px;
  color: var(--app-text-strong);
  font-size: 14px;
  word-break: break-word;
}

.meta-card p {
  margin: 8px 0 0;
  color: var(--app-text-muted);
  font-size: 12px;
  line-height: 1.7;
  word-break: break-word;
}

@media (max-width: 1200px) {
  .terminal-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .host-overview {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
