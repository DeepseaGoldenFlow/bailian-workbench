<template>
  <div class="chat">
    <h2>AI 对话</h2>
    <div class="chat-box" ref="chatBox">
      <div v-for="(msg, i) in messages" :key="i" :class="['msg', msg.role]">
        <div class="msg-content" v-html="renderContent(msg.content)"></div>
        <div class="msg-time">{{ msg.time }}</div>
      </div>
      <div v-if="streaming" class="msg assistant">
        <div class="msg-content" v-html="renderContent(streamContent)"></div>
      </div>
    </div>
    <div class="input-row">
      <el-input v-model="input" placeholder="输入消息..." @keydown.enter="send" :disabled="streaming" type="textarea" :rows="3" />
      <el-button type="primary" @click="send" :loading="streaming" style="margin-top:8px">发送</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { chatStream } from '../api'
import { marked } from 'marked'

const messages = ref([])
const input = ref('')
const streaming = ref(false)
const streamContent = ref('')
const chatBox = ref(null)

function renderContent(text) {
  if (!text) return ''
  return marked(text, { breaks: true })
}

function addMsg(role, content) {
  messages.value.push({ role, content, time: new Date().toLocaleTimeString() })
}

async function send() {
  const text = input.value.trim()
  if (!text || streaming.value) return
  input.value = ''
  addMsg('user', text)
  streaming.value = true
  streamContent.value = ''
  await nextTick()
  scrollDown()

  try {
    const resp = await chatStream({
      model: 'qwen-plus',
      messages: [{ role: 'user', content: text }],
      stream: true,
    })
    const reader = resp.body.getReader()
    const decoder = new TextDecoder()
    let full = ''
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      const chunk = decoder.decode(value, { stream: true })
      const lines = chunk.split('\n')
      for (const line of lines) {
        if (line.startsWith('data: ') && line !== 'data: [DONE]') {
          try {
            const json = JSON.parse(line.slice(6))
            const content = json.choices?.[0]?.delta?.content || ''
            full += content
            streamContent.value = full
            await nextTick()
            scrollDown()
          } catch {}
        }
      }
    }
    addMsg('assistant', full)
    streamContent.value = ''
  } catch (e) {
    addMsg('assistant', '错误: ' + e.message)
  }
  streaming.value = false
}

function scrollDown() {
  if (chatBox.value) {
    chatBox.value.scrollTop = chatBox.value.scrollHeight
  }
}
</script>

<style scoped>
.chat { max-width: 900px; margin: 0 auto; }
.chat-box { height: calc(100vh - 260px); overflow-y: auto; border: 1px solid #e0e0e0; border-radius: 8px; padding: 16px; background: #fff; margin: 16px 0; }
.msg { margin-bottom: 16px; }
.msg.user .msg-content { background: #409eff; color: #fff; border-radius: 12px 12px 0 12px; padding: 10px 14px; display: inline-block; max-width: 80%; }
.msg.assistant .msg-content { background: #f0f0f0; border-radius: 12px 12px 12px 0; padding: 10px 14px; display: inline-block; max-width: 80%; }
.msg.user { text-align: right; }
.msg-time { font-size: 11px; color: #999; margin-top: 4px; }
.input-row { display: flex; flex-direction: column; }
</style>