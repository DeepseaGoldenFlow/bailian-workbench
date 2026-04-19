<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎙️ 语音合成</h2>
      <p class="page-sub">将文字转换为自然流畅的语音</p>
    </div>

    <div class="page-grid">
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="🤖 CosyVoice V3.5 Plus" value="cosyvoice-v3.5-plus" />
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">输入文本</label>
          <el-input v-model="form.text" type="textarea" :rows="5" resize="none"
            placeholder="请输入要转换为语音的文字内容，例如：你好，欢迎使用百炼工作台，这是一个强大的 AI 工具平台。" />
          <span class="char-count">{{ form.text.length }} / 3000 字</span>
        </div>

        <div class="field">
          <label class="field-label">音色选择</label>
          <el-select v-model="form.voice" class="full-width" filterable>
            <el-option label="🧑 云霄 (中文男声)" value="longxiaochun" />
            <el-option label="👩 小云 (中文女声)" value="xiaoyun" />
            <el-option label="🧒 小刚 (中文童声)" value="xiaogang" />
            <el-option label="👨 男性通用 (英文)" value="male" />
            <el-option label="👩 女性通用 (英文)" value="female" />
          </el-select>
        </div>

        <div class="field">
          <div class="param-label-row">
            <span class="field-label">语速</span>
            <span class="param-value">{{ form.speed }}x</span>
          </div>
          <el-slider v-model="form.speed" :min="0.5" :max="2.0" :step="0.1" :show-tooltip="false"
            :marks="{0.5:'0.5x', 1.0:'1x', 1.5:'1.5x', 2.0:'2x'}" />
        </div>

        <div class="field">
          <label class="field-label">采样率</label>
          <el-select v-model="form.sample_rate" class="full-width">
            <el-option label="22050 Hz (标准)" value="22050" />
            <el-option label="44100 Hz (高清)" value="44100" />
            <el-option label="48000 Hz (无损)" value="48000" />
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">输出格式</label>
          <el-select v-model="form.format" class="full-width">
            <el-option label="WAV (无损)" value="wav" />
            <el-option label="MP3 (压缩)" value="mp3" />
            <el-option label="PCM (原始)" value="pcm" />
          </el-select>
        </div>

        <el-button @click="synthesize" type="primary" :loading="loading" class="generate-btn" size="large">
          <el-icon><Microphone /></el-icon> 开始合成
        </el-button>
      </div>

      <div class="result-card glass-card">
        <h3 class="card-title"><el-icon><Headset /></el-icon> 合成结果</h3>

        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">正在合成语音...</p>
        </div>

        <div v-if="audioUrl" class="audio-player-section">
          <div class="waveform-icon">🎵</div>
          <audio :src="audioUrl" controls class="audio-player" />
          <div class="audio-actions">
            <el-button @click="downloadAudio" type="primary" :icon="Download">下载音频</el-button>
          </div>
        </div>

        <el-empty v-if="!audioUrl && !loading" description="输入文本后开始合成语音 🎵" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditPen, Microphone, Headset, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({
  model: 'cosyvoice-v3.5-plus',
  text: '',
  voice: 'longxiaochun',
  speed: 1.0,
  sample_rate: '22050',
  format: 'wav',
})

const loading = ref(false)
const progress = ref(0)
const audioUrl = ref('')

const synthesize = async () => {
  if (!form.value.text.trim()) return ElMessage.warning('请输入要合成的文本')
  if (form.value.text.length > 3000) return ElMessage.warning('文本不能超过 3000 字')
  loading.value = true
  progress.value = 30

  try {
    const res = await fetch('/api/compatible-mode/v1/audio/speech', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: form.value.model,
        input: form.value.text,
        voice: form.value.voice,
        speed: form.value.speed,
        sample_rate: parseInt(form.value.sample_rate),
        response_format: form.value.format,
      })
    })

    if (!res.ok) {
      const e = await res.json().catch(() => ({ message: '合成失败' }))
      throw new Error(e.message || `HTTP ${res.status}`)
    }

    progress.value = 80
    const blob = await res.blob()
    audioUrl.value = URL.createObjectURL(blob)
    progress.value = 100
    ElMessage.success('🎉 语音合成完成！')
  } catch (e) {
    ElMessage.error('合成失败: ' + e.message)
  } finally {
    loading.value = false
    setTimeout(() => { progress.value = 0 }, 2000)
  }
}

const downloadAudio = () => {
  const ext = form.value.format === 'pcm' ? 'bin' : form.value.format
  const a = document.createElement('a')
  a.href = audioUrl.value
  a.download = `speech_${Date.now()}.${ext}`
  a.click()
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
.char-count { font-size: 11px; color: var(--text-secondary); float: right; margin-top: 4px; }
.full-width { width: 100%; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

.audio-player-section { text-align: center; padding: 20px 0; }
.waveform-icon { font-size: 48px; margin-bottom: 16px; }
.audio-player { width: 100%; margin: 12px 0; }
.audio-actions { display: flex; justify-content: center; gap: 12px; }

@media (max-width: 768px) { .page-grid { grid-template-columns: 1fr; } }
</style>
