<template>
  <div class="page-shell chat-page">
    <header class="page-header">
      <div><div class="page-eyebrow">智能对话</div><h1 class="page-title">随时可用的思考搭档</h1><p class="page-subtitle">讨论想法、整理资料、分析问题，支持通义千问与 DeepSeek 系列模型。</p></div>
      <v-btn v-if="messages.length" variant="tonal" rounded="xl" prepend-icon="mdi-broom" @click="clearChat">清空对话</v-btn>
    </header>

    <v-card rounded="xl" class="surface-card chat-shell">
      <div class="chat-toolbar">
        <div class="model-badge"><div class="model-badge__icon"><v-icon size="20">mdi-creation-outline</v-icon></div><div><span>智能对话</span><strong>选择合适的模型开始创作</strong></div></div>
        <div class="chat-controls">
          <v-select
            v-model="model"
            :items="chatModels"
            item-title="name"
            item-value="id"
            label="对话模型"
            aria-label="对话模型"
            density="compact"
            variant="outlined"
            hide-details
            class="chat-model-select"
          />
          <v-switch v-model="useStream" label="流式" aria-label="流式显示回答" color="primary" density="compact" hide-details inset />
          <v-chip size="small" variant="tonal" color="success"><span class="status-dot mr-2" />在线</v-chip>
        </div>
      </div>

      <div ref="chatRef" class="message-area">
        <div v-if="!messages.length && !streaming" class="welcome-state">
          <div class="welcome-state__icon"><v-icon size="36">mdi-chat-processing-outline</v-icon></div>
          <h2>你好，今天想聊什么？</h2>
          <p>你可以直接提问，也可以从下面选一个开始。</p>
          <div class="suggestion-grid">
            <button v-for="item in suggestions" :key="item.text" @click="input = item.text"><v-icon size="18">{{ item.icon }}</v-icon><span>{{ item.text }}</span><v-icon size="15">mdi-arrow-right</v-icon></button>
          </div>
        </div>

        <div v-for="(msg,index) in messages" :key="index" class="message" :class="`message--${msg.role}`">
          <div class="message__avatar"><v-icon size="18">{{ msg.role === 'user' ? 'mdi-account-outline' : 'mdi-creation-outline' }}</v-icon></div>
          <div class="message__body">
            <div class="message__name">{{ msg.role === 'user' ? '你' : currentModelName }}</div>
            <div v-if="msg.role === 'assistant'" class="markdown-body" v-html="renderMarkdown(msg.content)" />
            <div v-else class="message__text">{{ msg.content }}</div>
          </div>
        </div>
        <div v-if="streaming" class="message message--assistant">
          <div class="message__avatar"><v-icon size="18">mdi-creation-outline</v-icon></div>
          <div class="message__body"><div class="message__name">{{ currentModelName }}</div><div v-if="streamContent" class="markdown-body" v-html="renderMarkdown(streamContent)" /><div v-else class="typing"><i /><i /><i /></div></div>
        </div>
      </div>

      <div class="composer-wrap">
        <div class="composer">
          <textarea v-model="input" rows="1" placeholder="输入你的问题，Enter 发送，Shift + Enter 换行" @keydown.enter.exact.prevent="send" />
          <div class="composer__footer"><span>{{ input.length }}/8000</span><v-btn icon="mdi-arrow-up" color="primary" size="small" :loading="streaming" :disabled="!input.trim() || streaming" @click="send" /></div>
        </div>
        <div class="composer-tip">AI 生成的内容可能存在错误，请核对重要信息。</div>
      </div>
    </v-card>
  </div>
</template>

<script setup>
import { computed, nextTick, ref } from 'vue'
import { marked } from 'marked'
import api, { chatStream } from '../api'

const chatModels = [
  { id:'qwen-plus', name:'通义千问 Plus' }, { id:'qwen-max', name:'通义千问 Max' }, { id:'qwen-turbo', name:'通义千问 Turbo' }, { id:'deepseek-r1', name:'DeepSeek R1' }, { id:'deepseek-v3', name:'DeepSeek V3' },
]
const suggestions = [
  { icon:'mdi-lightbulb-outline', text:'帮我构思一个适合短视频的创意脚本' },
  { icon:'mdi-text-box-edit-outline', text:'把一段文字改得更自然、更有说服力' },
  { icon:'mdi-chart-timeline-variant', text:'帮我拆解一个复杂问题并制定行动计划' },
  { icon:'mdi-image-edit-outline', text:'为产品海报写一组高质量图片提示词' },
]
const model = ref('qwen-plus'), input = ref(''), messages = ref([]), streaming = ref(false), streamContent = ref(''), useStream = ref(true), chatRef = ref(null)
const currentModelName = computed(() => chatModels.find(item => item.id === model.value)?.name || model.value)
function renderMarkdown(text) { return marked(text || '', { breaks:true }) }
function clearChat() { messages.value = []; streamContent.value = '' }
async function scrollBottom() { await nextTick(); chatRef.value?.scrollTo({ top:chatRef.value.scrollHeight, behavior:'smooth' }) }
async function send() {
  if (!input.value.trim() || streaming.value) return
  const content = input.value.trim().slice(0,8000)
  input.value = ''
  messages.value.push({ role:'user', content })
  await scrollBottom()
  if (useStream.value) await sendStreaming(); else await sendNormal()
}
async function sendNormal() {
  streaming.value = true
  try {
    const response = await api.post('/chat/completions', { model:model.value, messages:messages.value.slice(-12) })
    messages.value.push({ role:'assistant', content:response.data.choices?.[0]?.message?.content || '没有收到有效回答。' })
  } catch (error) { messages.value.push({ role:'assistant', content:`请求失败：${error.response?.data?.message || error.message}` }) }
  streaming.value = false
  await scrollBottom()
}
async function sendStreaming() {
  streaming.value = true; streamContent.value = ''
  try {
    const response = await chatStream({ model:model.value, messages:messages.value.slice(-12), stream:true })
    if (!response.ok) throw new Error(`服务返回 ${response.status}`)
    const reader = response.body.getReader(), decoder = new TextDecoder(); let buffer = ''
    while (true) {
      const { done, value } = await reader.read(); if (done) break
      buffer += decoder.decode(value, { stream:true }); const lines = buffer.split('\n'); buffer = lines.pop()
      for (const line of lines) if (line.startsWith('data: ') && line !== 'data: [DONE]') try { const delta = JSON.parse(line.slice(6)).choices?.[0]?.delta?.content; if (delta) { streamContent.value += delta; await scrollBottom() } } catch {}
    }
    messages.value.push({ role:'assistant', content:streamContent.value || '没有收到有效回答。' })
  } catch (error) { messages.value.push({ role:'assistant', content:`请求失败：${error.message}` }) }
  streaming.value = false; streamContent.value = ''
  await scrollBottom()
}
</script>

<style scoped>
.chat-page { width:min(1160px,100%); }.chat-shell { height:calc(100vh - 190px); min-height:620px; display:flex; flex-direction:column; overflow:hidden; }.chat-toolbar { min-height:82px; flex:0 0 82px; display:flex; align-items:center; justify-content:space-between; gap:20px; padding:0 22px; border-bottom:1px solid rgba(var(--v-border-color),.1); }.model-badge { display:flex; align-items:center; gap:11px; min-width:0; }.model-badge__icon { width:38px; height:38px; flex:0 0 38px; display:grid; place-items:center; border-radius:13px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.1); }.model-badge div:last-child { display:flex; min-width:0; flex-direction:column; }.model-badge span { color:rgb(var(--v-theme-on-surface-variant)); font-size:10px; }.model-badge strong { margin-top:2px; overflow:hidden; font-size:13px; text-overflow:ellipsis; white-space:nowrap; }.chat-controls { display:flex; align-items:center; gap:12px; }.chat-model-select { width:230px; }.chat-controls :deep(.v-switch .v-label) { font-size:12px; }.message-area { flex:1; overflow-y:auto; padding:28px clamp(20px,6vw,70px); scroll-behavior:smooth; }.welcome-state { min-height:100%; display:flex; flex-direction:column; align-items:center; justify-content:center; text-align:center; }.welcome-state__icon { width:68px; height:68px; display:grid; place-items:center; margin-bottom:20px; border-radius:22px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.1); }.welcome-state h2 { margin:0; font-size:27px; }.welcome-state p { margin:8px 0 24px; color:rgb(var(--v-theme-on-surface-variant)); font-size:13px; }.suggestion-grid { width:min(680px,100%); display:grid; grid-template-columns:1fr 1fr; gap:10px; }.suggestion-grid button { display:flex; align-items:center; gap:10px; padding:14px; border:1px solid rgba(var(--v-border-color),.12); border-radius:14px; color:rgb(var(--v-theme-on-surface)); background:rgba(var(--v-theme-surface),.6); font:inherit; font-size:12px; text-align:left; cursor:pointer; }.suggestion-grid button:hover { border-color:rgba(var(--v-theme-primary),.35); background:rgba(var(--v-theme-primary),.05); }.suggestion-grid button span { flex:1; }.message { display:flex; gap:13px; max-width:850px; margin-bottom:26px; }.message--user { margin-left:auto; flex-direction:row-reverse; }.message__avatar { width:34px; height:34px; flex:0 0 34px; display:grid; place-items:center; border-radius:11px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.1); }.message--user .message__avatar { color:white; background:rgb(var(--v-theme-primary)); }.message__body { min-width:0; }.message__name { margin:0 0 7px; color:rgb(var(--v-theme-on-surface-variant)); font-size:11px; }.message--user .message__name { text-align:right; }.message__text,.markdown-body { padding:13px 16px; border-radius:5px 17px 17px 17px; background:rgb(var(--v-theme-surface-variant)); font-size:14px; line-height:1.75; }.message--user .message__text { color:white; background:rgb(var(--v-theme-primary)); border-radius:17px 5px 17px 17px; }.markdown-body :deep(p) { margin:0 0 10px; }.markdown-body :deep(p:last-child) { margin:0; }.markdown-body :deep(pre) { overflow:auto; padding:12px; border-radius:10px; background:rgba(0,0,0,.15); }.typing { display:flex; gap:5px; padding:17px; border-radius:5px 17px 17px; background:rgb(var(--v-theme-surface-variant)); }.typing i { width:6px; height:6px; border-radius:50%; background:rgb(var(--v-theme-primary)); animation:pulse 1s infinite alternate; }.typing i:nth-child(2) { animation-delay:.2s; }.typing i:nth-child(3) { animation-delay:.4s; }@keyframes pulse { to { opacity:.2; transform:translateY(-3px); } }.composer-wrap { flex:0 0 auto; padding:14px clamp(16px,5vw,58px) 17px; border-top:1px solid rgba(var(--v-border-color),.08); }.composer { padding:12px 12px 8px 16px; border:1px solid rgba(var(--v-border-color),.18); border-radius:18px; background:rgb(var(--v-theme-surface)); box-shadow:0 10px 28px rgba(15,23,42,.06); }.composer:focus-within { border-color:rgba(var(--v-theme-primary),.6); box-shadow:0 0 0 3px rgba(var(--v-theme-primary),.09); }.composer textarea { width:100%; max-height:120px; resize:none; outline:0; border:0; color:inherit; background:transparent; font:inherit; font-size:14px; line-height:1.6; }.composer__footer { display:flex; align-items:center; justify-content:flex-end; gap:10px; }.composer__footer span { color:rgb(var(--v-theme-on-surface-variant)); font-size:10px; }.composer-tip { margin-top:8px; color:rgb(var(--v-theme-on-surface-variant)); font-size:10px; text-align:center; }
@media (max-width:700px) { .chat-shell { height:calc(100vh - 170px); min-height:580px; }.suggestion-grid { grid-template-columns:1fr; }.message-area { padding-inline:16px; }.chat-toolbar { min-height:78px; flex-basis:78px; padding-inline:15px; }.model-badge { display:none; }.chat-controls { width:100%; }.chat-model-select { width:auto; flex:1; }.chat-controls .v-chip { display:none; } }
</style>
