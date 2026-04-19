<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🔊 语音识别</h2>
      <p class="page-sub">将音频文件转换为文字</p>
    </div>

    <div class="page-grid">
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="Paraformer 语音识别 (paraformer-realtime-v2)" value="paraformer-realtime-v2" />
            <el-option label="Paraformer 离线识别 (paraformer-v2)" value="paraformer-v2" />
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">上传音频文件</label>
          <el-upload :auto-upload="false" :show-file-list="false" accept="audio/*"
            @change="handleFileUpload" drag>
            <div v-if="!form.fileName" class="upload-area">
              <el-icon class="upload-icon"><UploadFilled /></el-icon>
              <p>拖拽音频文件到此处或点击上传</p>
              <span class="upload-hint">支持 MP3, WAV, FLAC, OGG, M4A, PCM 等格式</span>
            </div>
            <div v-else class="file-info">
              <el-icon><Document /></el-icon>
              <div>
                <p class="file-name">{{ form.fileName }}</p>
                <p class="file-size">{{ form.fileSize }}</p>
              </div>
              <el-button link type="danger" @click.stop="removeFile"><el-icon><Delete /></el-icon></el-button>
            </div>
          </el-upload>
        </div>

        <div class="field">
          <label class="field-label">采样率</label>
          <el-select v-model="form.sample_rate" class="full-width">
            <el-option label="自动检测" value="0" />
            <el-option label="8000 Hz (电话音质)" value="8000" />
            <el-option label="16000 Hz (标准)" value="16000" />
            <el-option label="44100 Hz (CD音质)" value="44100" />
            <el-option label="48000 Hz (专业)" value="48000" />
          </el-select>
        </div>

        <el-button @click="recognize" type="primary" :loading="loading" class="generate-btn" size="large"
          :disabled="!form.file">
          <el-icon><Microphone /></el-icon> 开始识别
        </el-button>
      </div>

      <div class="result-card glass-card">
        <h3 class="card-title">
          <el-icon><Document /></el-icon> 识别结果
          <el-button v-if="resultText" link type="primary" @click="copyResult" class="ml-auto"><el-icon><CopyDocument /></el-icon> 复制</el-button>
        </h3>

        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">正在识别语音...</p>
        </div>

        <div v-if="resultText" class="result-text">
          {{ resultText }}
        </div>

        <el-empty v-if="!resultText && !loading" description="上传音频文件后开始识别 📝" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditPen, Microphone, Document, UploadFilled, CopyDocument, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({
  model: 'paraformer-v2',
  file: null,
  fileName: '',
  fileSize: '',
  sample_rate: '0',
})

const loading = ref(false)
const progress = ref(0)
const resultText = ref('')

const handleFileUpload = (uploadFile) => {
  form.value.file = uploadFile.raw
  form.value.fileName = uploadFile.name
  form.value.fileSize = formatFileSize(uploadFile.size)
}

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1048576).toFixed(1) + ' MB'
}

const removeFile = () => {
  form.value.file = null; form.value.fileName = ''; form.value.fileSize = ''
}

const recognize = async () => {
  if (!form.value.file) return ElMessage.warning('请先上传音频文件')
  loading.value = true
  progress.value = 20
  resultText.value = ''

  try {
    const formData = new FormData()
    formData.append('audio', form.value.file)
    formData.append('model', form.value.model)
    if (parseInt(form.value.sample_rate) > 0) formData.append('sample_rate', form.value.sample_rate)

    const res = await fetch('/api/v1/services/audio/asr/transcription', {
      method: 'POST',
      body: formData,
    })

    if (!res.ok) { const e = await res.json(); throw new Error(e.message || '识别失败') }

    progress.value = 80
    const data = await res.json()
    resultText.value = data.output?.text || data.text || '未识别到内容'
    progress.value = 100
    ElMessage.success('🎉 语音识别完成！')
  } catch (e) {
    ElMessage.error('识别失败: ' + e.message)
  } finally {
    loading.value = false
    setTimeout(() => { progress.value = 0 }, 2000)
  }
}

const copyResult = async () => {
  try { await navigator.clipboard.writeText(resultText.value); ElMessage.success('已复制') }
  catch { ElMessage.error('复制失败') }
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
.ml-auto { margin-left: auto; }

.field { margin-bottom: 16px; }
.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.full-width { width: 100%; }

.upload-area { border: 2px dashed var(--card-border); border-radius: 10px; padding: 24px; text-align: center; cursor: pointer; transition: all 0.2s; }
.upload-area:hover { border-color: var(--gradient-start); }
.upload-icon { font-size: 32px; color: var(--text-secondary); margin-bottom: 8px; }
.upload-area p { font-size: 13px; color: var(--text-secondary); }
.upload-hint { font-size: 11px; color: var(--text-secondary); opacity: 0.7; }

.file-info { display: flex; align-items: center; gap: 12px; padding: 12px; background: rgba(99,102,241,0.1); border-radius: 8px; }
.file-name { font-size: 14px; font-weight: 500; }
.file-size { font-size: 12px; color: var(--text-secondary); }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

.result-text { padding: 16px; background: rgba(255,255,255,0.03); border-radius: 10px; line-height: 1.8; font-size: 15px; min-height: 120px; white-space: pre-wrap; }

@media (max-width: 768px) { .page-grid { grid-template-columns: 1fr; } }
</style>
