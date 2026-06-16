<template>
  <v-container fluid class="pa-6" style="max-width:900px">
    <h1 class="text-h4 font-weight-bold mb-1">Chat</h1>
    <p class="text-body-1 text-medium-emphasis mb-6">AI conversation via Bailian models</p>

    <v-card variant="outlined" rounded="xl" class="mb-4" max-height="60vh" style="overflow-y:auto">
      <div ref="chatRef" class="pa-4">
        <div v-for="(msg, i) in messages" :key="i" :class="['d-flex mb-4', msg.role === 'user' ? 'justify-end' : '']">
          <v-card :color="msg.role === 'user' ? 'primary' : 'surface-variant'" :variant="msg.role === 'user' ? 'flat' : 'tonal'" max-width="80%" rounded="xl" class="pa-3">
            <div v-if="msg.role === 'assistant'" v-html="renderMarkdown(msg.content)" class="text-body-2" />
            <div v-else class="text-body-2">{{ msg.content }}</div>
          </v-card>
        </div>
        <div v-if="streaming" class="d-flex mb-4">
          <v-card color="surface-variant" variant="tonal" max-width="80%" rounded="xl" class="pa-3">
            <div v-html="renderMarkdown(streamContent)" class="text-body-2" />
            <v-progress-linear indeterminate color="primary" class="mt-2" style="height:2px" />
          </v-card>
        </div>
      </div>
    </v-card>

    <v-card variant="outlined" rounded="xl" class="pa-4">
      <v-select v-model="model" :items="models" label="Model" variant="outlined" density="compact" hide-details class="mb-3" style="max-width:300px" />
      <v-textarea v-model="input" label="Message" variant="outlined" rows="2" density="compact" placeholder="Type your message..." hide-details @keydown.enter.exact.prevent="send" />
      <div class="d-flex align-center ga-2 mt-3">
        <v-btn color="primary" rounded="lg" @click="send" :loading="streaming" :disabled="!input.trim()">Send</v-btn>
        <v-switch v-model="useStream" label="Stream" density="compact" hide-details color="primary" />
        <v-spacer />
        <v-btn variant="text" color="error" size="small" @click="messages = []; streamContent = ''">Clear</v-btn>
      </div>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { marked } from 'marked'
import api from '../api'

const models = [{ id: 'qwen-plus', name: 'Qwen Plus' }, { id: 'qwen-max', name: 'Qwen Max' }, { id: 'qwen-turbo', name: 'Qwen Turbo' }, { id: 'deepseek-r1', name: 'DeepSeek R1' }, { id: 'deepseek-v3', name: 'DeepSeek V3' }]
const model = ref('qwen-plus'); const input = ref(''); const messages = ref([])
const streaming = ref(false); const streamContent = ref(''); const useStream = ref(true)
const chatRef = ref(null)

function renderMarkdown(text) { return marked(text || '', { breaks: true }) }

async function send() {
  if (!input.value.trim()) return
  const msg = input.value.trim(); input.value = ''
  messages.value.push({ role: 'user', content: msg })
  if (useStream.value) { await sendStream(msg) } else { await sendNormal(msg) }
}

async function sendNormal(msg) {
  try { const r = await api.post('/chat/completions', { model: model.value, messages: [...messages.value.slice(-10)] }); const content = r.data.choices?.[0]?.message?.content || ''; messages.value.push({ role: 'assistant', content }); await nextTick(); chatRef.value?.scrollTo(0, chatRef.value.scrollHeight) }
  catch (e) { messages.value.push({ role: 'assistant', content: 'Error: ' + (e.response?.data?.message || e.message) }) }
}

async function sendStream(msg) {
  streaming.value = true; streamContent.value = ''
  try {
    const resp = await fetch('/api/chat/completions', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ model: model.value, messages: [...messages.value.slice(-10)], stream: true }) })
    const reader = resp.body.getReader(); const decoder = new TextDecoder(); let buf = ''
    while (true) { const { done, value } = await reader.read(); if (done) break; buf += decoder.decode(value, { stream: true }); const lines = buf.split('\n'); buf = lines.pop(); for (const line of lines) { if (line.startsWith('data: ') && line !== 'data: [DONE]') { try { const d = JSON.parse(line.slice(6)); const c = d.choices?.[0]?.delta?.content; if (c) streamContent.value += c; await nextTick(); chatRef.value?.scrollTo(0, chatRef.value.scrollHeight) } catch(_) {} } } }
    messages.value.push({ role: 'assistant', content: streamContent.value })
  } catch (e) { messages.value.push({ role: 'assistant', content: 'Error: ' + e.message }) }
  streaming.value = false; streamContent.value = ''
}
</script>
