<template>
  <div class="page-container">
    <div class="page-header">
      <h2>👁️ 视觉理解</h2>
      <p class="page-sub">上传图片并让 AI 模型理解和分析画面内容</p>
    </div>

    <div class="page-grid">
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="Qwen-VL-Max (最强视觉理解)" value="qwen-vl-max" />
            <el-option label="Qwen-VL-Plus (高效视觉理解)" value="qwen-vl-plus" />
            <el-option label="Qwen2.5-VL-72B-Instruct" value="qwen2.5-vl-72b-instruct" />
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">上传图片</label>
          <div class="upload-grid">
            <div v-for="(img, idx) in images" :key="idx" class="uploaded-thumb">
              <img :src="img.url" />
              <el-icon class="thumb-remove" @click="removeImage(idx)"><Close /></el-icon>
            </div>
            <el-upload action="" :auto-upload="false" :show-file-list="false" accept="image/*"
              multiple @change="handleImageUpload" class="upload-trigger">
              <div class="upload-placeholder">
                <el-icon><Plus /></el-icon>
                <span>添加图片</span>
              </div>
            </el-upload>
          </div>
        </div>

        <div class="field">
          <label class="field-label">对话消息</label>
          <el-input v-model="input" type="textarea" :rows="3" resize="none"
            placeholder="描述你想让 AI 分析的内容，例如：请描述这张图片中的场景，并指出主要元素" />
        </div>

        <!-- Parameters -->
        <div class="params-toggle" @click="showParams = !showParams">
          <el-icon><Setting /></el-icon>
          <span>高级参数</span>
          <el-icon class="arrow"><ArrowDown v-if="!showParams" /><ArrowUp v-else /></el-icon>
        </div>

        <transition name="param-slide">
          <div v-show="showParams">
            <div class="field">
              <div class="param-label-row">
                <span class="field-label">创意程度</span>
                <span class="param-value">{{ params.temperature.toFixed(1) }}</span>
              </div>
              <el-slider v-model="params.temperature" :min="0" :max="2" :step="0.1" :show-tooltip="false" />
            </div>

            <div class="field">
              <div class="param-label-row">
                <span class="field-label">最大回复长度</span>
                <span class="param-value">{{ params.max_tokens }} tokens</span>
              </div>
              <el-slider v-model="params.max_tokens" :min="256" :max="4096" :step="256" :show-tooltip="false" />
            </div>
          </div>
        </transition>

        <el-button @click="sendMessage" type="primary" :loading="loading" class="generate-btn" size="large"
          :disabled="images.length === 0 && !input.trim()">
          <el-icon><View /></el-icon> 开始分析
        </el-button>
      </div>

      <div class="result-card glass-card">
        <h3 class="card-title"><el-icon><ChatDotRound /></el-icon> 对话结果</h3>

        <div class="vision-messages" ref="msgRef">
          <div v-if="messages.length === 0" class="vision-welcome">
            <el-icon class="welcome-icon" :size="48"><View /></el-icon>
            <p>上传图片并输入问题，AI 将为你分析画面内容</p>
          </div>

          <div v-for="(msg, i) in messages" :key="i" class="v-msg" :class="msg.role">
            <div class="v-msg-avatar">
              <div class="v-avatar-icon" :class="msg.role">{{ msg.role === 'user' ? '👤' : '✦' }}</div>
            </div>
            <div class="v-msg-bubble">
              <div v-if="msg.images && msg.images.length > 0" class="v-msg-images">
                <img v-for="(img, j) in msg.images" :key="j" :src="img" />
              </div>
              <div class="v-msg-text" v-if="msg.role === 'assistant'" v-html="renderMarkdown(msg.content)"></div>
              <div class="v-msg-text" v-else>{{ msg.content }}</div>
            </div>
          </div>

          <div v-if="loading" class="v-msg assistant">
            <div class="v-msg-avatar"><div class="v-avatar-icon assistant">✦</div></div>
            <div class="v-msg-bubble"><div class="typing-indicator"><span></span><span></span><span></span></div></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { EditPen, Plus, Close, Setting, ArrowDown, ArrowUp, ChatDotRound, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({ html: true, linkify: true, breaks: true })

const form = ref({ model: 'qwen-vl-max' })
const images = ref([])
const input = ref('')
const messages = ref([])
const msgRef = ref(null)
const loading = ref(false)
const showParams = ref(false)

const params = ref({ temperature: 0.7, max_tokens: 2048 })

const renderMarkdown = (text) => md.render(text || '')

const handleImageUpload = (uploadFile) => {
  const file = uploadFile.raw
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => { images.value.push({ url: e.target.result, base64: e.target.result }) }
  reader.readAsDataURL(file)
}

const removeImage = (idx) => { images.value.splice(idx, 1) }

const scrollToBottom = async () => {
  await nextTick()
  if (msgRef.value) msgRef.value.scrollTop = msgRef.value.scrollHeight
}

const sendMessage = async () => {
  if (!input.value.trim() && images.value.length === 0) return
  if (loading.value) return
  loading.value = true

  const userImages = images.value.map(i => i.base64)
  const userText = input.value
  messages.value.push({ role: 'user', content: userText || '请分析图片', images: userImages })
  const assistantIdx = messages.value.length
  messages.value.push({ role: 'assistant', content: '' })
  input.value = ''
  images.value = []
  await scrollToBottom()

  try {
    // Build multimodal messages
    const historyContent = []
    // Add previous messages
    for (let i = 0; i < assistantIdx - 1; i++) {
      const m = messages.value[i]
      if (m.role === 'user') {
        const parts = []
        if (m.images && m.images.length > 0) {
          for (const img of m.images) {
            parts.push({ type: 'image_url', image_url: { url: img } })
          }
        }
        parts.push({ type: 'text', text: m.content })
        historyContent.push({ role: 'user', content: parts })
      } else {
        historyContent.push({ role: 'assistant', content: m.content })
      }
    }

    const lastMsg = messages.value[assistantIdx - 1]
    const lastParts = []
    if (lastMsg.images && lastMsg.images.length > 0) {
      for (const img of lastMsg.images) {
        lastParts.push({ type: 'image_url', image_url: { url: img } })
      }
    }
    lastParts.push({ type: 'text', text: lastMsg.content })
    historyContent.push({ role: 'user', content: lastParts })

    const res = await fetch('/api/compatible-mode/v1/chat/completions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: form.value.model,
        messages: historyContent,
        temperature: params.value.temperature,
        max_tokens: params.value.max_tokens,
      })
    })

    if (!res.ok) { const e = await res.json(); throw new Error(e.message || '请求失败') }

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
              if (delta) { messages.value[assistantIdx].content += delta; await scrollToBottom() }
            } catch {}
          }
        }
      }
      if (!messages.value[assistantIdx].content) messages.value[assistantIdx].content = '(无回复)'
    } else {
      const data = await res.json()
      messages.value[assistantIdx].content = data.choices?.[0]?.message?.content || '(无回复)'
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
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.page-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }

.glass-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; padding: 20px; backdrop-filter: blur(10px); }
.card-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }

.field { margin-bottom: 16px; }
.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.param-label-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.param-value { font-size: 13px; font-weight: 700; color: var(--gradient-start); font-variant-numeric: tabular-nums; }
.full-width { width: 100%; }

.upload-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(100px, 1fr)); gap: 8px; }
.uploaded-thumb { position: relative; aspect-ratio: 1; border-radius: 8px; overflow: hidden; }
.uploaded-thumb img { width: 100%; height: 100%; object-fit: cover; }
.thumb-remove { position: absolute; top: 4px; right: 4px; background: rgba(0,0,0,0.6); color: #fff; border-radius: 50%; width: 20px; height: 20px; font-size: 12px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.upload-trigger { cursor: pointer; }
.upload-placeholder { aspect-ratio: 1; border: 2px dashed var(--card-border); border-radius: 8px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; color: var(--text-secondary); font-size: 12px; transition: all 0.2s; }
.upload-placeholder:hover { border-color: var(--gradient-start); }

.params-toggle { display: flex; align-items: center; gap: 8px; padding: 8px 0; cursor: pointer; color: var(--text-secondary); font-size: 14px; }
.params-toggle:hover { color: var(--text-primary); }
.params-toggle .arrow { margin-left: auto; font-size: 14px; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.vision-messages { flex: 1; overflow-y: auto; max-height: calc(100vh - 200px); }
.vision-welcome { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 300px; text-align: center; color: var(--text-secondary); }
.welcome-icon { margin-bottom: 12px; color: var(--gradient-start); }

.v-msg { display: flex; gap: 10px; margin-bottom: 14px; max-width: 90%; animation: msgIn 0.3s ease; }
.v-msg.user { align-self: flex-end; flex-direction: row-reverse; margin-left: auto; }
.v-msg.assistant { align-self: flex-start; }
.v-avatar-icon { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 14px; flex-shrink: 0; }
.v-avatar-icon.user { background: linear-gradient(135deg, #3b82f6, #6366f1); }
.v-avatar-icon.assistant { background: linear-gradient(135deg, var(--gradient-start), var(--gradient-end)); }
.v-msg-bubble { padding: 10px 14px; border-radius: 12px; line-height: 1.7; font-size: 14px; word-break: break-word; }
.v-msg.user .v-msg-bubble { background: linear-gradient(135deg, rgba(99,102,241,0.2), rgba(139,92,246,0.15)); border: 1px solid rgba(99,102,241,0.2); }
.v-msg.assistant .v-msg-bubble { background: var(--card-bg); border: 1px solid var(--card-border); }

.v-msg-images { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 8px; }
.v-msg-images img { max-width: 150px; max-height: 120px; border-radius: 6px; object-fit: cover; }
.v-msg-text :deep(p) { margin: 0 0 6px; }
.v-msg-text :deep(p:last-child) { margin-bottom: 0; }
.v-msg-text :deep(code) { background: rgba(255,255,255,0.08); padding: 1px 5px; border-radius: 3px; font-size: 12px; }
.v-msg-text :deep(pre) { background: rgba(0,0,0,0.3); padding: 10px; border-radius: 6px; overflow-x: auto; margin: 6px 0; }
.v-msg-text :deep(pre code) { background: none; padding: 0; }

.typing-indicator { display: flex; gap: 4px; padding: 4px 0; }
.typing-indicator span { width: 7px; height: 7px; border-radius: 50%; background: var(--gradient-start); animation: typing 1.4s infinite; }
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes msgIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes typing { 0%, 60%, 100% { opacity: 0.3; transform: scale(0.8); } 30% { opacity: 1; transform: scale(1); } }
@media (max-width: 768px) { .page-grid { grid-template-columns: 1fr; } }
</style>
