<template>
  <div class="page">
    <h2>语音识别 (ASR)</h2>
    <el-upload drag :auto-upload="false" :on-change="onFileChange" accept="audio/*">
      <el-icon><UploadFilled /></el-icon>
      <div>拖拽或点击上传音频文件</div>
    </el-upload>
    <div v-if="file" style="margin-top:12px">已选: {{ file.name }}</div>
    <el-button type="primary" @click="transcribe" :loading="loading" style="margin-top:12px">开始识别</el-button>
    <div v-if="result" class="result">
      <el-input v-model="result" type="textarea" :rows="6" readonly />
    </div>
    <div v-if="raw" style="margin-top:12px">
      <el-collapse><el-collapse-item title="API 响应"><pre>{{ raw }}</pre></el-collapse-item></el-collapse>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { asr } from '../api'

const file = ref(null)
const loading = ref(false)
const result = ref('')
const raw = ref('')

function onFileChange(uploadFile) {
  file.value = uploadFile.raw
}

async function transcribe() {
  if (!file.value) return
  loading.value = true
  result.value = ''
  raw.value = ''
  try {
    const fd = new FormData()
    fd.append('file', file.value)
    const { data } = await asr(fd)
    raw.value = JSON.stringify(data, null, 2)
    result.value = data.output?.text || data.output?.results?.[0]?.transcription_url || JSON.stringify(data)
  } catch (e) {
    raw.value = e.response?.data || e.message
  }
  loading.value = false
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.result { margin-top: 16px; }
pre { white-space: pre-wrap; font-size: 12px; }
</style>