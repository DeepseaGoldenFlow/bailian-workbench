<template>
  <div class="chat-container">
    <!-- Header -->
    <div class="chat-header">
      <div class="header-left">
        <h2>✦ AI 对话</h2>
        <span class="header-sub">与千问大模型进行智能对话</span>
      </div>
      <div class="header-right">
        <el-select v-model="selectedModel" class="model-select" placeholder="选择模型">
          <el-option label="Qwen 3.6 Plus (最新旗舰)" value="qwen3.6-plus" />
          <el-option label="Qwen 3 Max (最强推理)" value="qwen3-max" />
        </el-select>
        <el-button @click="clearChat" :icon="Delete" circle title="清空对话" />
      </div>
    </div>

    <div class="chat-body">
      <!-- Parameters Panel -->
      <div class="params-panel" :class="{ open: showParams }">
        <div class="params-header" @click="showParams = !showParams">
          <el-icon><Setting /></el-icon>
          <span>参数配置</span>
          <el-icon class="arrow"><ArrowDown v-if="!showParams" /><ArrowUp v-else /></el-icon>
        </div>
        <transition name="param-slide">
          <div v-show="showParams" class="params-content">
            <!-- 创意程度 -->
            <div class="param-group">
              <div class="param-label-row">
                <span class="param-label">创意程度</span>
                <span class="param-value">{{ params.temperature.toFixed(1) }}</span>
              </div>
              <el-slider v-model="params.temperature" :min="0" :max="2" :step="0.1" :show-tooltip="false" />
              <p class="param-desc">较低时回复更精准稳定，较高时更有创造力和多样性</p>
            </div>

            <!-- 回复多样性 -->
            <div class="param-group">
              <div class="param-label-row">
                <span class="param-label">回复多样性</span>
                <span class="param-value">{{ params.top_p.toFixed(1) }}</span>
              </div>
              <el-slider v-model="params.top_p" :min="0.1" :max="1" :step="0.05" :show-tooltip="false" />
              <p class="param-desc">控制生成文本的随机性，值越高输出越多样化</p>
            </div>

            <!-- 最大回复长度 -->
            <div class="param-group">
              <div class="param-label-row">
                <span class="param-label">最大回复长度</span>
                <span class="param-value">{{ params.max_tokens }} tokens</span>
              </div>
              <el-slider v-model="params.max_tokens" :min="256" :max="8192" :step="256" :show-tooltip="false" />
              <p class="param-desc">限制模型生成回复的最大 token 数量</p>
            </div>

            <!-- 惩罚系数 -->
            <div class="param-group">
              <div class="param-label-row">
                <span class="param-label">重复惩罚</span>
                <span class="param-value">{{ params.repetition_penalty.toFixed(1) }}</span>
              </div>
              <el-slider v-model="params.repetition_penalty" :min="1" :max="2" :step="0.1" :show-tooltip="false" />
              <p class="param-desc">降低模型重复生成相同内容的倾向，值越高惩罚越强</p>
            </div>

            <!-- 联网搜索 -->
            <div class="param-group">
              <div class="param-label-row">
                <span class="param-label">联网搜索</span>
                <el-switch v-model="params.enable_search" active-text="开启" inactive-text="关闭" />
              </div>
              <p class="param-desc">开启后模型可以联网获取最新信息来回答你的问题</p>
            </div>
          </div>
        </transition>
      </div>

      <!-- Messages Area -->
      <div class="messages-area" ref="msgRef">
        <div v-if="messages.length === 0" class="welcome-screen">
          <div class="welcome-icon">✦</div>
          <h3>开始与 AI 对话</h3>
          <p>选择左侧模型并输入你的问题</p>
          <div class="suggestions">
            <div class="suggestion-card" v-for="s in suggestions" :key="s" @click="input = s">
              <span>{{ s }}</span>
            </div>
          </div>
        </div>

        <div v-for="(msg, i) in messages" :key="i" class="message-row" :class="msg.role">
          <div class="message-avatar">
            <div class="avatar-icon" :class="msg.role">{{ msg.role === 'user' ? '👤' : '✦' }}</div>
          </div>
          <div class="message-bubble">
            <div class="message-content" v-if="msg.role === 'assistant'" v-html="renderMarkdown(msg.content)"></div>
            <div class="message-content" v-else>{{ msg.content }}</div>
            <div class="message-actions" v-if="msg.role === 'assistant' && msg.content">
              <el-button link size="small" @click="copyMessage(msg.content)"><el-icon><CopyDocument /></el-icon> 复制</el-button>
            </div>
          </div>
        </div>

        <div v-if="loading" class="message-row assistant">
          <div class="message-avatar"><div class="avatar-icon assistant">✦</div></div>
          <div class="message-bubble"><div class="typing-indicator"><span></span><span></span><span></span></div></div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="input-area">
        <div class="input-wrapper">
          <el-input v-model="input" type="textarea" :autosize="{ minRows: 1, maxRows: 6 }" resize="none"
            placeholder="输入消息，按 Enter 发送，Shift+Enter 换行..." @keydown="handleKeydown" />
          <el-button @click="sendMessage" :loading="loading" type="primary" :icon="Promotion" class="send-btn"
            :disabled="!input.trim()">
            发送
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { Delete, Setting, ArrowDown, ArrowUp, CopyDocument, Promotion } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({ html: true, linkify: true, breaks: true })

const input = ref('')
const messages = ref([])
const msgRef = ref(null)
const loading = ref(false)
const showParams = ref(false)
const selectedModel = ref('qwen-turbo')

const params = ref({
  temperature: 0.7,
  top_p: 0.8,
  max_tokens: 2048,
  repetition_penalty: 1.1,
  enable_search: false,
})

const suggestions = [
  '帮我写一段 Python 代码实现快速排序',
  '解释一下量子计算的基本原理',
  '用通俗语言介绍大语言模型是什么',
  '帮我写一首关于秋天的现代诗',
]

const renderMarkdown = (text) => md.render(text || '')

const scrollToBottom = async () => {
  await nextTick()
  if (msgRef.value) msgRef.value.scrollTop = msgRef.value.scrollHeight
}

const handleKeydown = (e) => {
  if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage() }
}

const clearChat = () => { messages.value = []; ElMessage.success('已清空对话') }

const copyMessage = async (text) => {
  try { await navigator.clipboard.writeText(text); ElMessage.success('已复制到剪贴板') }
  catch { ElMessage.error('复制失败') }
}

const sendMessage = async () => {
  if (!input.value.trim() || loading.value) return
  loading.value = true

  const userMsg = input.value
  messages.value.push({ role: 'user', content: userMsg })
  const assistantIdx = messages.value.length
  messages.value.push({ role: 'assistant', content: '' })
  input.value = ''
  await scrollToBottom()

  try {
    const historyMessages = messages.value.slice(0, -1).map(m => ({ role: m.role, content: m.content }))
    const body = {
      model: selectedModel.value,
      messages: historyMessages,
      temperature: params.value.temperature,
      top_p: params.value.top_p,
      max_tokens: params.value.max_tokens,
    }
    if (params.value.enable_search) body.enable_search = true
    if (params.value.repetition_penalty !== 1.1) body.repetition_penalty = params.value.repetition_penalty

    const res = await fetch('/api/compatible-mode/v1/chat/completions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })

    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.message || `HTTP ${res.status}`)
    }

    // Streaming
    if (res.headers.get('content-type')?.includes('text/event-stream')) {
      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''
        for (const line of lines) {
          if (line.startsWith('data: ') && line !== 'data: [DONE]') {
            try {
              const json = JSON.parse(line.slice(6))
              const delta = json.choices?.[0]?.delta?.content
              if (delta) {
                messages.value[assistantIdx].content += delta
                await scrollToBottom()
              }
            } catch {}
          }
        }
      }
      // Fallback if streaming produced nothing
      if (!messages.value[assistantIdx].content) {
        messages.value[assistantIdx].content = '(无回复)'
      }
    } else {
      // Non-streaming fallback
      const data = await res.json()
      const content = data.choices?.[0]?.message?.content || '(无回复)'
      messages.value[assistantIdx].content = content
    }
  } catch (e) {
    messages.value[assistantIdx].content = '❌ ' + e.message
  } finally {
    loading.value = false
    await scrollToBottom()
  }
}
</script>

<style scoped>
.chat-container { height: calc(100vh - 48px); display: flex; flex-direction: column; max-width: 1000px; margin: 0 auto; width: 100%; }

.chat-header { display: flex; align-items: center; justify-content: space-between; padding: 8px 0 16px; flex-shrink: 0; }
.header-left h2 { font-size: 20px; font-weight: 700; background: linear-gradient(135deg, var(--gradient-start), var(--gradient-end)); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.header-sub { font-size: 13px; color: var(--text-secondary); }
.header-right { display: flex; align-items: center; gap: 12px; }
.model-select { width: 180px; }

.chat-body { flex: 1; display: flex; flex-direction: column; overflow: hidden; background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; backdrop-filter: blur(10px); }

/* Params Panel */
.params-panel { border-bottom: 1px solid var(--card-border); flex-shrink: 0; }
.params-header { display: flex; align-items: center; gap: 8px; padding: 12px 16px; cursor: pointer; color: var(--text-secondary); font-size: 14px; font-weight: 500; transition: all 0.2s; }
.params-header:hover { color: var(--text-primary); }
.params-header .arrow { margin-left: auto; font-size: 14px; }
.params-content { padding: 0 16px 12px; }

.param-group { padding: 10px 0; border-bottom: 1px solid var(--card-border); }
.param-group:last-child { border-bottom: none; }
.param-label-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.param-label { font-size: 13px; font-weight: 600; color: var(--text-primary); }
.param-value { font-size: 13px; font-weight: 700; color: var(--gradient-start); font-variant-numeric: tabular-nums; }
.param-desc { font-size: 11px; color: var(--text-secondary); margin-top: 4px; line-height: 1.5; }

/* Messages */
.messages-area { flex: 1; overflow-y: auto; padding: 16px; }

.welcome-screen { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; text-align: center; }
.welcome-icon { font-size: 48px; background: linear-gradient(135deg, var(--gradient-start), var(--gradient-end)); -webkit-background-clip: text; -webkit-text-fill-color: transparent; margin-bottom: 16px; }
.welcome-screen h3 { font-size: 18px; margin-bottom: 8px; }
.welcome-screen p { font-size: 13px; color: var(--text-secondary); margin-bottom: 24px; }
.suggestions { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; width: 100%; max-width: 500px; }
.suggestion-card { padding: 12px 16px; background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 10px; font-size: 13px; color: var(--text-secondary); cursor: pointer; transition: all 0.2s; text-align: left; }
.suggestion-card:hover { border-color: var(--gradient-start); color: var(--text-primary); transform: translateY(-1px); }

.message-row { display: flex; gap: 12px; margin-bottom: 16px; max-width: 85%; animation: msgIn 0.3s ease; }
.message-row.user { align-self: flex-end; flex-direction: row-reverse; margin-left: auto; }
.message-row.assistant { align-self: flex-start; }

.avatar-icon { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0; }
.avatar-icon.user { background: linear-gradient(135deg, #3b82f6, #6366f1); }
.avatar-icon.assistant { background: linear-gradient(135deg, var(--gradient-start), var(--gradient-end)); }

.message-bubble { padding: 12px 16px; border-radius: 14px; line-height: 1.7; font-size: 14px; word-break: break-word; }
.message-row.user .message-bubble { background: linear-gradient(135deg, rgba(99,102,241,0.2), rgba(139,92,246,0.15)); border: 1px solid rgba(99,102,241,0.2); }
.message-row.assistant .message-bubble { background: var(--card-bg); border: 1px solid var(--card-border); }

.message-content :deep(p) { margin: 0 0 8px; }
.message-content :deep(p:last-child) { margin-bottom: 0; }
.message-content :deep(code) { background: rgba(255,255,255,0.08); padding: 2px 6px; border-radius: 4px; font-size: 13px; font-family: 'Fira Code', monospace; }
.message-content :deep(pre) { background: rgba(0,0,0,0.3); padding: 12px; border-radius: 8px; overflow-x: auto; margin: 8px 0; }
.message-content :deep(pre code) { background: none; padding: 0; }
.message-content :deep(ul), .message-content :deep(ol) { padding-left: 20px; margin: 8px 0; }
.message-content :deep(blockquote) { border-left: 3px solid var(--gradient-start); padding-left: 12px; color: var(--text-secondary); margin: 8px 0; }
.message-content :deep(table) { border-collapse: collapse; margin: 8px 0; width: 100%; }
.message-content :deep(th), .message-content :deep(td) { border: 1px solid var(--card-border); padding: 6px 10px; text-align: left; }

.message-actions { display: flex; gap: 8px; margin-top: 8px; }

/* Typing Indicator */
.typing-indicator { display: flex; gap: 4px; padding: 4px 0; }
.typing-indicator span { width: 8px; height: 8px; border-radius: 50%; background: var(--gradient-start); animation: typing 1.4s infinite; }
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

/* Input Area */
.input-area { padding: 12px 16px; border-top: 1px solid var(--card-border); flex-shrink: 0; }
.input-wrapper { display: flex; gap: 10px; align-items: flex-end; }
.input-wrapper :deep(.el-textarea__inner) { border-radius: 12px; padding: 10px 14px; font-size: 14px; }
.send-btn { min-height: 40px; border-radius: 10px; padding: 0 20px; }

@keyframes msgIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes typing { 0%, 60%, 100% { opacity: 0.3; transform: scale(0.8); } 30% { opacity: 1; transform: scale(1); } }
</style>
