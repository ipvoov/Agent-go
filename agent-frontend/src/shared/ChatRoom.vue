<template>
  <div style="display:flex;flex-direction:column;flex:1;min-height:0;">
    <div style="padding:12px 16px;border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;">
      <div style="font-weight:600;">{{ title }}</div>
      <div style="color:var(--muted);font-size:12px;">会话ID：{{ sessionId }}</div>
    </div>

    <div ref="scrollRef" style="flex:1; overflow:auto; padding: 12px 16px; display:flex; flex-direction:column; gap:14px;">
      <template v-for="(m, idx) in messages" :key="idx">
        <div v-if="!m.hidden" class="msg-enter row" :class="{ right: m.role === 'user' }">
          <div :style="avatarStyle(m.role)">
            <img :src="m.role === 'user' ? userAvatar : aiAvatar" alt="avatar" style="width:100%;height:100%;border-radius:50%;object-fit:cover;" />
          </div>
          <div :class="['bubble', m.role]" style="white-space:pre-wrap; word-break:break-word;">
            <template v-if="m.role === 'thinking'">
              <div style="display:flex;align-items:center;gap:6px;margin-bottom:6px;color:var(--muted);font-size:12px;">
                <span class="dot" v-if="thinkingActive && idx === currentThinkingIndex" />
                <span>{{ (thinkingActive && idx === currentThinkingIndex) ? '正在思考中…' : '思考内容' }}</span>
                <button class="btn btn--ghost" style="padding:4px 8px;font-size:12px;margin-left:8px;" @click.prevent="toggleThinking(idx)">{{ expandedThinking.has(idx) ? '收起' : '展开' }}</button>
              </div>
              <div class="thinking-text thinking-content" :class="{ 'collapsed': !expandedThinking.has(idx) }">{{ m.content }}</div>
            </template>
            <template v-else>
              {{ m.content }}
            </template>
            <div v-if="m.role === 'assistant' && lastThinkingContent" class="bubble-tools">
              <button class="link-btn" @click="showThinking=true">查看思考</button>
            </div>
          </div>
        </div>
      </template>

      <!-- 思考弹窗 -->
      <div v-if="showThinking" class="overlay" @click.self="showThinking=false">
        <div class="modal">
          <h4>思考内容</h4>
          <pre>{{ lastThinkingContent }}</pre>
          <div style="text-align:right; margin-top:8px;">
            <button class="btn" @click="showThinking=false">关闭</button>
          </div>
        </div>
      </div>
    </div>

    <form @submit.prevent="onSend" style="padding:12px; border-top:1px solid var(--border); display:flex; gap:8px;">
      <input class="input" :disabled="loading" v-model="input" :placeholder="loading ? 'AI 生成中...' : '输入你的内容，回车发送'" />
      <button class="btn" :disabled="!input.trim() || loading" type="submit">发送</button>
    </form>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { API_BASE } from '../utils/config.js'
import { generateSessionId } from '../utils/session.js'
import { decodeLatin1Utf8 } from '../utils/encoding.js'
import aiAvatar from '../assets/ai.svg'
import userAvatar from '../assets/user-cute.svg'

const props = defineProps({
  title: { type: String, default: 'Chat' },
  endpoint: { type: String, required: true }
})

const input = ref('')
const messages = ref([])
const sessionId = ref('')
const loading = ref(false)
const scrollRef = ref(null)

let es = null
let currentThinkingIndex = -1
let currentAssistantIndex = -1
const thinkingActive = ref(false)
const expandedThinking = ref(new Set())
const lastThinkingContent = ref('')
const showThinking = ref(false)

// 打字机相关
let typingTimer = null
let typingBuffer = ''
const TYPING_TICK_MS = 20
const TYPING_CHARS_PER_TICK = 4

function ensureTypingTimer() {
  if (typingTimer) return
  typingTimer = setInterval(() => {
    if (!typingBuffer || currentAssistantIndex === -1) {
      if (!loading.value && !typingBuffer) {
        clearInterval(typingTimer)
        typingTimer = null
      }
      return
    }
    const take = typingBuffer.slice(0, TYPING_CHARS_PER_TICK)
    typingBuffer = typingBuffer.slice(TYPING_CHARS_PER_TICK)
    const cur = messages.value[currentAssistantIndex]
    if (cur && cur.role === 'assistant') cur.content += take
    requestAnimationFrame(scrollToBottom)
  }, TYPING_TICK_MS)
}

function collapseThinking() {
  if (currentThinkingIndex !== -1) {
    // 直接隐藏思考气泡
    const msg = messages.value[currentThinkingIndex]
    if (msg && msg.role === 'thinking') msg.hidden = true
    const set = new Set(expandedThinking.value)
    set.delete(currentThinkingIndex)
    expandedThinking.value = set
  }
}

const bubbleStyle = (role) => ({
  display: 'flex',
  gap: '8px',
  alignItems: 'flex-start',
  maxWidth: '90%',
  marginLeft: (role === 'assistant' || role === 'thinking') ? '0' : 'auto',
  marginRight: (role === 'assistant' || role === 'thinking') ? 'auto' : '0',
  background: 'var(--panel)',
  border: '1px solid var(--border)',
  borderRadius: (role === 'assistant' || role === 'thinking') ? '12px 12px 12px 4px' : '12px 12px 4px 12px',
  padding: '10px 12px'
})

const avatarStyle = (role) => ({
  width: '32px', height: '32px', flex: '0 0 auto',
  display: 'grid', placeItems: 'center',
  borderRadius: '50%',
  background: (role === 'assistant' || role === 'thinking') ? '#e2e8f0' : 'var(--primary)',
  color: (role === 'assistant' || role === 'thinking') ? '#0f172a' : '#fff',
  fontSize: '12px', fontWeight: 700
})

const scrollToBottom = () => {
  const el = scrollRef.value
  if (!el) return
  el.scrollTop = el.scrollHeight
}

watch(messages, () => {
  // 在消息更新后保持滚动在底部
  requestAnimationFrame(scrollToBottom)
}, { deep: true })

function startSSE(query) {
  if (es) {
    es.close()
    es = null
  }
  loading.value = true

  const url = `${API_BASE}${props.endpoint}?query=${encodeURIComponent(query)}&session_id=${encodeURIComponent(sessionId.value)}`
  es = new EventSource(url)
  thinkingActive.value = false

  es.onmessage = (evt) => {
    try {
      const raw = evt.data || '{}'
      const fixed = decodeLatin1Utf8(raw)
      const data = JSON.parse(fixed)
      if (data.done) {
        loading.value = false
        es && es.close()
        es = null
        // 结束思考占位
        // 安全兜底：若还存在思考气泡，隐藏之
        if (currentThinkingIndex !== -1) {
          const msg = messages.value[currentThinkingIndex]
          if (msg && msg.role === 'thinking') msg.hidden = true
        }
        currentThinkingIndex = -1
        thinkingActive.value = false
        // 完整结束后滚到底
        requestAnimationFrame(scrollToBottom)
        return
      }
      const content = decodeLatin1Utf8(data.content || '')

      if (data.thinking === true) {
        // 流式思考内容
        if (currentThinkingIndex === -1) {
          messages.value.push({ role: 'thinking', content })
          currentThinkingIndex = messages.value.length - 1
          thinkingActive.value = true
          expandedThinking.value.add(currentThinkingIndex)
        } else {
          const cur = messages.value[currentThinkingIndex]
          if (cur && cur.role === 'thinking') cur.content += content
        }
        // 记录最新思考全文
        lastThinkingContent.value += content
        return
      }

      // 输出正式回答（非思考）
      if (thinkingActive.value) thinkingActive.value = false
      // 接收到第一段非思考内容时，自动折叠思考区
      collapseThinking()
      if (currentAssistantIndex === -1) {
        messages.value.push({ role: 'assistant', content: '' })
        currentAssistantIndex = messages.value.length - 1
      } else {
        // 已有回答气泡，继续往下打字
      }
      // 将本次片段放入缓冲，按速率输出
      typingBuffer += content
      ensureTypingTimer()

      // 一旦开始输出答案，清除思考占位（可保留历史思考文本，不删除）
      currentThinkingIndex = -1
    } catch (e) {
      // ignore parse errors
    }
  }

  es.onerror = () => {
    loading.value = false
    es && es.close()
    es = null
  }
}

function onSend() {
  const q = input.value.trim()
  if (!q || loading.value) return
  messages.value.push({ role: 'user', content: q })
  input.value = ''
  // 重置索引
  currentThinkingIndex = -1
  currentAssistantIndex = -1
  thinkingActive.value = false
  expandedThinking.value = new Set()
  // 清理打字机状态
  typingBuffer = ''
  if (typingTimer) { clearInterval(typingTimer); typingTimer = null }
  lastThinkingContent.value = ''
  startSSE(q)
}

function toggleThinking(idx) {
  const set = new Set(expandedThinking.value)
  if (set.has(idx)) set.delete(idx)
  else set.add(idx)
  expandedThinking.value = set
}

onMounted(() => {
  sessionId.value = generateSessionId()
})

onUnmounted(() => {
  es && es.close()
  es = null
})
</script>

<style scoped>
.row { display:flex; align-items:flex-end; gap:8px; }
.row.right { flex-direction: row-reverse; }
.bubble { max-width: 86%; padding: 10px 12px; border-radius: 14px; font-size: 13px; line-height: 1.6; }
.bubble.assistant, .bubble.thinking { background: var(--panel); border:1px solid var(--border); border-radius: 14px 14px 14px 6px; }
.row.right .bubble.user { background: var(--primary); color: #fff; border-radius: 14px 14px 6px 14px; }
.thinking-content {
  max-height: 220px;
  overflow: auto;
  border-top: 1px dashed var(--border);
  padding-top: 6px;
}
.thinking-content.collapsed {
  max-height: 1.6em;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

/* 思考弹窗 */
.overlay { position: fixed; inset: 0; background: rgba(0,0,0,.45); display: grid; place-items: center; }
.modal { width: min(760px, 92vw); max-height: 72vh; overflow: auto; background: var(--panel); border:1px solid var(--border); border-radius: 12px; padding: 16px; }
.modal h4 { margin: 0 0 8px 0; }
.modal pre { white-space: pre-wrap; word-break: break-word; font-family: inherit; font-size: 13px; color: var(--muted); }
.bubble-tools { margin-top: 6px; display:flex; gap:8px; }
.link-btn { background: transparent; color: var(--primary); border: none; padding: 0; cursor: pointer; font-size: 12px; }
</style>

