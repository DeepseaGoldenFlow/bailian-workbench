<template>
  <div class="pa-8" style="max-width:1000px;margin:0 auto">
    <div class="d-flex align-center mb-8">
      <v-avatar color="blue" size="48" class="mr-4 elevation-4"><v-icon size="28">mdi-chat-outline</v-icon></v-avatar>
      <div><h1 class="text-h4 font-weight-bold">Chat</h1><p class="text-body-1 text-medium-emphasis mt-1">AI conversation with streaming response</p></div>
    </div>
    <v-card rounded="xl" variant="outlined" class="mb-4" max-height="55vh" style="overflow-y:auto">
      <div ref="chatRef" class="pa-4">
        <div v-for="(msg, i) in messages" :key="i" :class="['d-flex mb-4', msg.role === 'user' ? 'justify-end' : '']">
          <v-card :color="msg.role === 'user' ? 'primary' : 'surface-variant'" :variant="msg.role === 'user' ? 'flat' : 'tonal'" max-width="80%" rounded="xl" class="pa-3">
            <div v-if="msg.role === 'assistant'" v-html="renderMarkdown(msg.content)" class="text-body-2 text-high-emphasis" />
            <div v-else class="text-body-2" style="color:white">{{ msg.content }}</div>
          </v-card>
        </div>
        <div v-if="streaming" class="d-flex mb-4">
          <v-card color="surface-variant" variant="tonal" max-width="80%" rounded="xl" class="pa-3">
            <div v-html="renderMarkdown(streamContent)" class="text-body-2 text-high-emphasis" /><v-progress-linear indeterminate color="primary" class="mt-2" style="height:2px" />
          </v-card>
        </div>
        <div v-if="messages.length === 0 && !streaming" class="text-center py-16 text-medium-emphasis">
          <v-icon size="48" class="mb-2">mdi-chat-plus-outline</v-icon><p>Start a conversation</p>
        </div>
      </div>
    </v-card>
    <v-card rounded="xl" class="pa-4" elevation="2" style="background:linear-gradient(180deg, rgba(25,118,210,0.06) 0%, transparent 40%)">
      <v-row dense align="center">
        <v-col cols="12" md="3"><v-select v-model="model" :items="chatModels" item-title="name" item-value="id" label="Model" variant="outlined" density="compact" hide-details color="blue" /></v-col>
        <v-col cols="12" md><v-text-field v-model="input" placeholder="Type a message..." variant="outlined" density="compact" hide-details @keydown.enter.exact.prevent="send" /></v-col>
        <v-col cols="auto" class="d-flex ga-2">
          <v-btn color="primary" rounded="lg" @click="send" :loading="streaming" :disabled="!input.trim()" variant="elevated"><v-icon start>mdi-send</v-icon> Send</v-btn>
          <v-switch v-model="useStream" label="Stream" density="compact" hide-details color="primary" />
        </v-col>
      </v-row>
    </v-card>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { marked } from 'marked'
import api from '../api'

const chatModels = [
  { id: 'qwen-plus', name: 'Qwen Plus' }, { id: 'qwen-max', name: 'Qwen Max' },
  { id: 'qwen-turbo', name: 'Qwen Turbo' }, { id: 'deepseek-r1', name: 'DeepSeek R1' },
  { id: 'deepseek-v3', name: 'DeepSeek V3' },
]
const model = ref('qwen-plus'); const input = ref(''); const messages = ref([])
const streaming = ref(false); const streamContent = ref(''); const useStream = ref(true)
const chatRef = ref(null)

function renderMarkdown(text) { return marked(text || '', { breaks: true }) }

async function send() {
  if (!input.value.trim()) return
  const msg = input.value.trim(); input.value = ''
  messages.value.push({ role: 'user', content: msg })
  useStream.value ? await sendStream(msg) : await sendNormal(msg)
}

async function sendNormal(msg) {
  try {
    const r = await api.post('/chat/completions', { model: model.value, messages: [...messages.value.slice(-10)] })
    messages.value.push({ role: 'assistant', content: r.data.choices?.[0]?.message?.content || '' })
    await nextTick(); chatRef.value?.scrollTo(0, chatRef.value.scrollHeight)
  } catch (e) { messages.value.push({ role: 'assistant', content: 'Error: ' + (e.response?.data?.message || e.message) }) }
}

async function sendStream(msg) {
  streaming.value = true; streamContent.value = ''
  try {
    const resp = await fetch('/api/chat/completions', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ model: model.value, messages: [...messages.value.slice(-10)], stream: true }) })
    const reader = resp.body.getReader(); const decoder = new TextDecoder(); let buf = ''
    while (true) { const { done, value } = await reader.read(); if (done) break; buf += decoder.decode(value, { stream: true }); const NL = String.fromCharCode(10); const lines = buf.split(NL); buf = lines.pop(); for (const line of lines) { if (line.startsWith('data: ') && line !== 'data: [DONE]') { try { const d = JSON.parse(line.slice(6)); const c = d.choices?.[0]?.delta?.content; if (c) streamContent.value += c; await nextTick(); chatRef.value?.scrollTo(0, chatRef.value.scrollHeight) } catch(_) {} } } }
    messages.value.push({ role: 'assistant', content: streamContent.value })
  } catch (e) { messages.value.push({ role: 'assistant', content: 'Error: ' + e.message }) }
  streaming.value = false; streamContent.value = ''
}
</script>