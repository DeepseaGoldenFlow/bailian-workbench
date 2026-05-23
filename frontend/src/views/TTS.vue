<template>
  <div class="page">
    <h2>语音合成 (TTS)</h2>
    <el-form label-width="80px" style="max-width:600px">
      <el-form-item label="文本">
        <el-input v-model="input" type="textarea" :rows="4" placeholder="输入要合成语音的文字..." />
      </el-form-item>
      <el-form-item label="音色">
        <el-select v-model="voice">
          <el-option label="Cherry" value="Cherry" />
          <el-option label="Eric" value="Eric" />
          <el-option label="Emily" value="Emily" />
          <el-option label="Luna" value="Luna" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="synthesize" :loading="loading">合成语音</el-button>
      </el-form-item>
    </el-form>
    <div v-if="audioUrl" class="result">
      <audio :src="audioUrl" controls style="width:100%" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { tts } from '../api'

const input = ref('')
const voice = ref('Cherry')
const loading = ref(false)
const audioUrl = ref('')

async function synthesize() {
  if (!input.value.trim()) return
  loading.value = true
  audioUrl.value = ''
  try {
    const { data } = await tts({ input: input.value, voice: voice.value, format: 'mp3' })
    audioUrl.value = URL.createObjectURL(data)
  } catch (e) {
    alert('TTS 错误: ' + (e.response?.data || e.message))
  }
  loading.value = false
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.result { background: #fff; padding: 16px; border-radius: 8px; margin-top: 16px; }
</style>