<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🤖 数字人</h2>
      <p class="page-sub">上传人物图片与驱动音频，生成唇形同步的数字人说话视频</p>
    </div>

    <div class="page-grid">
      <!-- ==================== 左侧参数区 ==================== -->
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><UserFilled /></el-icon> 参数配置</h3>

        <!-- 模型选择 -->
        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width" disabled>
            <el-option label="万相 2.2 数字人（S2V）" value="wan2.2-s2v" />
          </el-select>
          <p class="field-desc">使用人物图片 + 音频驱动生成说话视频</p>
        </div>

        <!-- 人物图片上传 -->
        <div class="field">
          <label class="field-label">人物图片</label>
          <el-upload
            class="avatar-uploader"
            :auto-upload="false"
            :show-file-list="false"
            accept="image/jpeg,image/png,image/webp"
            :before-upload="beforeImageUpload"
            @change="handleImageChange"
          >
            <div v-if="imagePreview" class="upload-preview">
              <img :src="imagePreview" alt="人物预览" />
              <div class="upload-mask">
                <el-icon><RefreshRight /></el-icon>
                <span>更换图片</span>
              </div>
            </div>
            <div v-else class="upload-placeholder">
              <el-icon class="upload-icon"><Plus /></el-icon>
              <span class="upload-text">点击上传人物图片</span>
              <span class="upload-hint">支持 JPG / PNG / WebP，清晰正面人像</span>
            </div>
          </el-upload>
          <p class="field-desc">建议上传清晰正面人像，面部完整无遮挡，效果更佳</p>
        </div>

        <!-- 驱动音频上传 -->
        <div class="field">
          <label class="field-label">驱动音频</label>
          <el-upload
            class="audio-uploader"
            :auto-upload="false"
            :show-file-list="false"
            accept=".mp3,.wav,.m4a,audio/mpeg,audio/wav,audio/mp4,audio/x-m4a"
            :before-upload="beforeAudioUpload"
            @change="handleAudioChange"
          >
            <div v-if="audioFile" class="audio-preview">
              <div class="audio-icon-wrap">
                <el-icon class="audio-icon"><Headset /></el-icon>
              </div>
              <div class="audio-info">
                <span class="audio-name">{{ audioFile.name }}</span>
                <span class="audio-size">{{ formatFileSize(audioFile.size) }}</span>
              </div>
              <el-button link type="danger" @click.stop="clearAudio">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <div v-else class="audio-placeholder">
              <el-icon class="upload-icon"><UploadFilled /></el-icon>
              <span class="upload-text">点击上传驱动音频</span>
              <span class="upload-hint">支持 MP3 / WAV / M4A 格式</span>
            </div>
          </el-upload>
          <p class="field-desc">音频将驱动人物口型与表情，建议语音时长 1~5 分钟</p>
        </div>

        <!-- 分辨率选择 -->
        <div class="field">
          <label class="field-label">输出分辨率</label>
          <el-select v-model="form.resolution" class="full-width">
            <el-option label="480P（快速生成）" value="480P" />
            <el-option label="720P（高清画质）" value="720P" />
          </el-select>
          <p class="field-desc">选择生成视频的分辨率，高分辨率需要更长的生成时间</p>
        </div>

        <el-button
          @click="submitDigitalHuman"
          type="primary"
          :loading="loading"
          :disabled="!canSubmit"
          class="generate-btn"
          size="large"
        >
          <el-icon><VideoCamera /></el-icon> 开始生成
        </el-button>
      </div>

      <!-- ==================== 右侧结果区 ==================== -->
      <div class="result-card glass-card">
        <h3 class="card-title"><el-icon><VideoCameraFilled /></el-icon> 生成结果</h3>

        <!-- 进度条 -->
        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">{{ progressText }}</p>
        </div>

        <!-- 错误信息 -->
        <div v-if="error" class="error-section">
          <el-alert :title="error" type="error" show-icon :closable="false" />
        </div>

        <!-- 视频结果 -->
        <div v-if="videoUrl" class="video-result">
          <div class="video-wrapper">
            <video
              :src="videoUrl"
              controls
              preload="metadata"
              class="result-video"
            />
          </div>
          <div class="video-actions">
            <el-button type="primary" @click="downloadVideo">
              <el-icon><Download /></el-icon> 下载视频
            </el-button>
            <el-button @click="resetForm">
              <el-icon><RefreshRight /></el-icon> 重新生成
            </el-button>
          </div>
        </div>

        <el-empty
          v-if="!videoUrl && !loading && !error"
          description="上传人物图片与音频，开始生成数字人视频 🎬"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import {
  UserFilled, Plus, RefreshRight, Headset, UploadFilled,
  Delete, VideoCamera, VideoCameraFilled, Download
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// ==================== 表单数据 ====================
const form = ref({
  model: 'wan2.2-s2v',
  resolution: '480P',
})

const imageFile = ref(null)
const imagePreview = ref('')
const audioFile = ref(null)

// ==================== 状态 ====================
const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const error = ref('')
const videoUrl = ref('')
let pollTimer = null

// 是否可以提交
const canSubmit = computed(() => {
  return imageFile.value && audioFile.value
})

// ==================== 图片上传处理 ====================
function beforeImageUpload(file) {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG / PNG / WebP 格式的图片')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  return false // 阻止默认上传，自行处理
}

function handleImageChange(uploadFile) {
  const file = uploadFile.raw
  if (!file) return
  imageFile.value = file
  // 生成预览
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target.result
  }
  reader.readAsDataURL(file)
}

// ==================== 音频上传处理 ====================
function beforeAudioUpload(file) {
  const isAudio = [
    'audio/mpeg', 'audio/wav', 'audio/mp4', 'audio/x-m4a',
    'audio/ogg', 'audio/aac'
  ].includes(file.type) || /\.(mp3|wav|m4a|ogg|aac)$/i.test(file.name)
  if (!isAudio) {
    ElMessage.error('仅支持 MP3 / WAV / M4A 格式的音频')
    return false
  }
  const isLt50M = file.size / 1024 / 1024 < 50
  if (!isLt50M) {
    ElMessage.error('音频大小不能超过 50MB')
    return false
  }
  return false
}

function handleAudioChange(uploadFile) {
  const file = uploadFile.raw
  if (!file) return
  audioFile.value = file
}

function clearAudio() {
  audioFile.value = null
}

// ==================== 工具函数 ====================
function formatFileSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

/** 文件转 base64 data URL */
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

// ==================== 提交任务 ====================
async function submitDigitalHuman() {
  if (!imageFile.value) {
    return ElMessage.warning('请上传人物图片')
  }
  if (!audioFile.value) {
    return ElMessage.warning('请上传驱动音频')
  }

  loading.value = true
  progress.value = 5
  progressText.value = '正在读取文件...'
  error.value = ''
  videoUrl.value = ''

  try {
    // 将图片和音频转为 base64
    progressText.value = '正在转换图片...'
    const imageBase64 = await fileToBase64(imageFile.value)

    progressText.value = '正在转换音频...'
    progress.value = 10
    const audioBase64 = await fileToBase64(audioFile.value)

    // 提交任务
    progressText.value = '正在提交任务...'
    progress.value = 15

    const res = await fetch('/api/video/digital-human', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: form.value.model,
        input: {
          image: imageBase64,
          audio: audioBase64,
        },
        parameters: {
          resolution: form.value.resolution,
        },
      }),
    })

    if (!res.ok) {
      const errData = await res.json().catch(() => ({}))
      throw new Error(errData.message || `请求失败 (${res.status})`)
    }

    const data = await res.json()
    const taskId = data.output?.task_id
    if (!taskId) throw new Error('未返回 task_id')

    progress.value = 20
    progressText.value = '任务已提交，正在生成视频...'

    // 开始轮询
    await pollTask(taskId)
  } catch (e) {
    error.value = `生成失败: ${e.message}`
    ElMessage.error(error.value)
  } finally {
    loading.value = false
    // 延迟清除进度条
    setTimeout(() => {
      if (!loading.value) {
        progress.value = 0
        progressText.value = ''
      }
    }, 2000)
  }
}

// ==================== 轮询任务 ====================
function pollTask(taskId) {
  return new Promise((resolve, reject) => {
    let pollCount = 0
    const maxPolls = 300 // 最长 300 * 5s = 25 分钟（视频生成较慢）

    const doPoll = async () => {
      try {
        const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
        if (!pollRes.ok) throw new Error(`轮询失败 (${pollRes.status})`)
        const pollData = await pollRes.json()
        const status = pollData.output?.task_status

        if (!status) {
          pollCount++
          if (pollCount >= maxPolls) {
            throw new Error('任务超时')
          }
          pollTimer = setTimeout(doPoll, 5000)
          return
        }

        if (status === 'SUCCEEDED') {
          progress.value = 100
          progressText.value = '生成完成！'
          // 获取视频 URL
          videoUrl.value = pollData.output?.video_url || pollData.output?.results?.[0]?.url || ''
          if (!videoUrl.value) {
            // 尝试从 results 或其他字段获取
            const results = pollData.output?.results
            if (results && results.length > 0) {
              videoUrl.value = results[0].url || ''
            }
          }
          if (videoUrl.value) {
            ElMessage.success('🎉 数字人视频生成成功！')
          } else {
            // 兜底：尝试从 output 直接拿
            videoUrl.value = pollData.output?.url || ''
          }
          resolve()
          return
        }

        if (status === 'FAILED') {
          const errMsg = pollData.output?.message || pollData.output?.code || '未知错误'
          throw new Error(errMsg)
        }

        // PENDING / RUNNING
        pollCount++
        if (pollCount >= maxPolls) {
          throw new Error('任务超时（等待超过 25 分钟）')
        }

        progress.value = Math.min(20 + pollCount * 2, 95)
        progressText.value = `生成中... 已等待 ${pollCount * 5}秒`

        pollTimer = setTimeout(doPoll, 5000)
      } catch (e) {
        reject(e)
      }
    }

    doPoll()
  })
}

// ==================== 下载 & 重置 ====================
function downloadVideo() {
  if (!videoUrl.value) return
  const a = document.createElement('a')
  a.href = videoUrl.value
  a.download = `digital_human_${Date.now()}.mp4`
  a.target = '_blank'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

function resetForm() {
  videoUrl.value = ''
  error.value = ''
  progress.value = 0
  progressText.value = ''
  imageFile.value = null
  imagePreview.value = ''
  audioFile.value = null
}

// ==================== 组件卸载时清理 ====================
onUnmounted(() => {
  if (pollTimer) clearTimeout(pollTimer)
})
</script>

<style scoped>
/* ==================== 布局 ==================== */
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.page-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }

/* ==================== 卡片 ==================== */
.glass-card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(10px);
}
.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
}

/* ==================== 表单字段 ==================== */
.field { margin-bottom: 16px; }
.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}
.field-desc {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 4px;
  line-height: 1.4;
}

.full-width { width: 100%; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

/* ==================== 图片上传 ==================== */
.avatar-uploader {
  display: block;
}
.upload-preview {
  width: 100%;
  aspect-ratio: 3 / 4;
  max-height: 280px;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  cursor: pointer;
  border: 1px solid var(--card-border);
}
.upload-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.upload-mask {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #fff;
  opacity: 0;
  transition: opacity 0.3s;
  font-size: 13px;
}
.upload-preview:hover .upload-mask { opacity: 1; }

.upload-placeholder {
  width: 100%;
  aspect-ratio: 3 / 4;
  max-height: 280px;
  border-radius: 12px;
  border: 2px dashed var(--card-border);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  transition: border-color 0.3s, background 0.3s;
}
.upload-placeholder:hover {
  border-color: var(--gradient-start);
  background: var(--card-hover);
}
.upload-icon { font-size: 36px; color: var(--text-secondary); }
.upload-text { font-size: 14px; color: var(--text-primary); }
.upload-hint { font-size: 11px; color: var(--text-secondary); }

/* ==================== 音频上传 ==================== */
.audio-uploader {
  display: block;
}
.audio-placeholder {
  width: 100%;
  padding: 24px;
  border-radius: 12px;
  border: 2px dashed var(--card-border);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  transition: border-color 0.3s, background 0.3s;
}
.audio-placeholder:hover {
  border-color: var(--gradient-start);
  background: var(--card-hover);
}

.audio-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid var(--card-border);
  background: var(--card-hover);
}
.audio-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: var(--btn-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.audio-icon { font-size: 24px; color: #fff; }
.audio-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.audio-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.audio-size { font-size: 11px; color: var(--text-secondary); }

/* ==================== 进度 & 错误 ==================== */
.progress-section { padding: 16px 0; }
.progress-text {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 8px;
  text-align: center;
}
.error-section { margin-bottom: 16px; }

/* ==================== 视频结果 ==================== */
.video-result { margin-top: 8px; }
.video-wrapper {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--card-border);
  background: #000;
}
.result-video {
  width: 100%;
  max-height: 480px;
  display: block;
}
.video-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
  justify-content: center;
}

/* ==================== 响应式 ==================== */
@media (max-width: 768px) {
  .page-grid { grid-template-columns: 1fr; }
}
</style>
