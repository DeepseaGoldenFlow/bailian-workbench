<template>
  <div class="page-container">
    <div class="page-header">
      <h2>💃 动画动效</h2>
      <p class="page-sub">上传人物图片 + 参考动作视频，一键生成舞蹈 / 动作视频</p>
    </div>

    <div class="main-layout">
      <!-- 左侧：参数配置 -->
      <div class="input-panel glass-card">
        <h3 class="panel-title"><el-icon><Setting /></el-icon> 参数配置</h3>

        <!-- 人物图片上传 -->
        <div class="field">
          <label class="field-label">人物图片 <span class="required">*</span></label>
          <el-upload
            action=""
            :auto-upload="false"
            :show-file-list="false"
            accept="image/jpeg,image/png,image/webp"
            @change="handleImageUpload"
            drag
            class="upload-dragger"
          >
            <div v-if="!imageUrl" class="upload-placeholder">
              <el-icon class="upload-icon"><PictureFilled /></el-icon>
              <p>拖拽人物图片到此处</p>
              <p class="upload-hint">支持 JPG / PNG / WebP，建议全身清晰照</p>
            </div>
            <img v-else :src="imageUrl" class="upload-preview" alt="人物图片预览" />
          </el-upload>
          <el-button v-if="imageUrl" link type="danger" size="small" @click="clearImage">
            <el-icon><Delete /></el-icon> 移除图片
          </el-button>
          <p class="field-desc">上传一张全身清晰的人物图片，姿势自然效果更佳</p>
        </div>

        <!-- 参考视频上传 -->
        <div class="field">
          <label class="field-label">动作参考视频 <span class="required">*</span></label>
          <el-upload
            action=""
            :auto-upload="false"
            :show-file-list="false"
            accept="video/mp4,video/webm,video/quicktime"
            @change="handleVideoUpload"
            drag
            class="upload-dragger"
          >
            <div v-if="!videoUrl" class="upload-placeholder">
              <el-icon class="upload-icon"><VideoCamera /></el-icon>
              <p>拖拽动作参考视频到此处</p>
              <p class="upload-hint">支持 MP4 / WebM / MOV</p>
            </div>
            <video v-else :src="videoUrl" controls class="upload-preview-video" />
          </el-upload>
          <el-button v-if="videoUrl" link type="danger" size="small" @click="clearVideo">
            <el-icon><Delete /></el-icon> 移除视频
          </el-button>
          <p class="field-desc">上传包含目标动作 / 舞蹈的视频，人物将模仿该动作</p>
        </div>

        <!-- 模式选择 -->
        <div class="field">
          <label class="field-label">生成模式</label>
          <el-select v-model="mode" class="full-width">
            <el-option label="⚡ 标准模式（速度快）" value="wan-std" />
            <el-option label="🎯 专业模式（动作更流畅）" value="wan-pro" />
          </el-select>
          <p class="field-desc">
            标准模式生成速度更快，适合快速预览；专业模式生成时间较长但动作更流畅自然
          </p>
        </div>

        <el-button
          @click="submitTask"
          type="primary"
          :loading="loading"
          class="submit-btn"
          size="large"
        >
          <el-icon><VideoCameraFilled /></el-icon>
          开始生成
        </el-button>
      </div>

      <!-- 右侧：生成结果 -->
      <div class="result-panel glass-card">
        <h3 class="panel-title"><el-icon><VideoCameraFilled /></el-icon> 生成结果</h3>

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
        <div v-if="resultVideoUrl" class="video-result">
          <video :src="resultVideoUrl" controls preload="metadata" class="video-player" />
          <div class="video-actions">
            <el-button @click="downloadVideo" type="primary" plain>
              <el-icon><Download /></el-icon> 下载视频
            </el-button>
            <el-button @click="resetForm">
              <el-icon><RefreshRight /></el-icon> 重新生成
            </el-button>
          </div>
        </div>

        <el-empty
          v-if="!resultVideoUrl && !loading"
          description="上传人物图片和参考视频开始生成吧 ✨"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import {
  Setting,
  PictureFilled,
  VideoCamera,
  VideoCameraFilled,
  Delete,
  Download,
  RefreshRight,
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// ==================== 状态 ====================
const imageUrl = ref('')
const imageBase64 = ref('')
const videoUrl = ref('')
const videoBase64 = ref('')
const mode = ref('wan-std')

const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const error = ref('')
const resultVideoUrl = ref('')

let pollTimer = null

// ==================== 上传处理 ====================
function handleImageUpload(file) {
  const raw = file.raw
  if (!raw) return
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(raw.type)
  if (!isImage) {
    return ElMessage.error('仅支持 JPG / PNG / WebP 格式的图片')
  }
  const reader = new FileReader()
  reader.onload = (e) => {
    imageBase64.value = e.target.result
    imageUrl.value = e.target.result
  }
  reader.readAsDataURL(raw)
}

function handleVideoUpload(file) {
  const raw = file.raw
  if (!raw) return
  const isVideo = raw.type.startsWith('video/')
  if (!isVideo) {
    return ElMessage.error('请上传有效的视频文件')
  }
  const reader = new FileReader()
  reader.onload = (e) => {
    videoBase64.value = e.target.result
    videoUrl.value = e.target.result
  }
  reader.readAsDataURL(raw)
}

function clearImage() {
  imageUrl.value = ''
  imageBase64.value = ''
}

function clearVideo() {
  videoUrl.value = ''
  videoBase64.value = ''
}

// ==================== 提交 & 轮询 ====================
async function submitTask() {
  if (!imageBase64.value) {
    return ElMessage.warning('请上传人物图片')
  }
  if (!videoBase64.value) {
    return ElMessage.warning('请上传动作参考视频')
  }

  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'
  error.value = ''
  resultVideoUrl.value = ''

  try {
    const res = await fetch('/api/video/animate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: 'wan2.2-animate-move',
        input: {
          image: imageBase64.value,
          ref_video: videoBase64.value,
        },
        parameters: {
          mode: mode.value,
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

    progress.value = 15
    progressText.value = '任务已提交，正在生成舞蹈视频...'

    await pollTask(taskId)
  } catch (e) {
    error.value = `生成失败: ${e.message}`
    ElMessage.error(error.value)
  } finally {
    loading.value = false
    setTimeout(() => {
      if (!loading.value) {
        progress.value = 0
        progressText.value = ''
      }
    }, 2000)
  }
}

function pollTask(taskId) {
  return new Promise((resolve, reject) => {
    let pollCount = 0
    const maxPolls = 240 // 最长 240 * 5s = 20 分钟

    const doPoll = async () => {
      try {
        const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
        if (!pollRes.ok) throw new Error(`轮询失败 (${pollRes.status})`)
        const pollData = await pollRes.json()
        const status = pollData.output?.task_status

        if (!status) {
          pollCount++
          if (pollCount >= maxPolls) throw new Error('任务超时')
          pollTimer = setTimeout(doPoll, 5000)
          return
        }

        if (status === 'SUCCEEDED') {
          progress.value = 100
          progressText.value = '舞蹈视频生成完成！'

          let vUrl =
            pollData.output?.video_url ||
            pollData.output?.results?.[0]?.url ||
            pollData.output?.video_urls?.[0]

          if (vUrl) {
            if (vUrl.startsWith('/data/bailian/storage')) {
              vUrl = vUrl.replace('/data/bailian/storage', '/api/files')
            }
            resultVideoUrl.value = vUrl
            ElMessage.success('🎉 舞蹈视频生成成功！')
          } else {
            throw new Error('未获取到视频 URL')
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
          throw new Error('任务超时（等待超过 20 分钟）')
        }

        progress.value = Math.min(15 + pollCount * 2, 95)
        progressText.value = `生成中... 已等待 ${pollCount * 5}秒 [${status}]`

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
  if (resultVideoUrl.value) {
    const a = document.createElement('a')
    a.href = resultVideoUrl.value
    a.download = `animate_move_${Date.now()}.mp4`
    a.target = '_blank'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
  }
}

function resetForm() {
  resultVideoUrl.value = ''
  imageUrl.value = ''
  imageBase64.value = ''
  videoUrl.value = ''
  videoBase64.value = ''
  error.value = ''
  progress.value = 0
  progressText.value = ''
}

// ==================== 清理 ====================
function cleanup() {
  if (pollTimer) {
    clearTimeout(pollTimer)
    pollTimer = null
  }
}

onUnmounted(cleanup)
</script>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
.page-header {
  margin-bottom: 24px;
}
.page-header h2 {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}
.page-sub {
  font-size: 14px;
  color: var(--text-secondary);
}

/* ==================== 布局 ==================== */
.main-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  align-items: start;
}

/* ==================== 玻璃卡片 ==================== */
.glass-card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(10px);
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
}

/* ==================== 表单字段 ==================== */
.field {
  margin-bottom: 16px;
}
.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}
.required {
  color: #ef4444;
  margin-left: 2px;
}
.field-desc {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 4px;
  line-height: 1.4;
}
.full-width {
  width: 100%;
}

/* ==================== 上传区域 ==================== */
.upload-dragger {
  width: 100%;
}
.upload-dragger :deep(.el-upload) {
  width: 100%;
}
.upload-dragger :deep(.el-upload-dragger) {
  background: transparent;
  border: 2px dashed var(--card-border);
  border-radius: 12px;
  padding: 24px;
  transition: all 0.2s;
}
.upload-dragger :deep(.el-upload-dragger:hover) {
  border-color: var(--gradient-start);
}
.upload-placeholder {
  text-align: center;
}
.upload-icon {
  font-size: 36px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}
.upload-placeholder p {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0;
}
.upload-hint {
  font-size: 11px !important;
  margin-top: 4px !important;
}
.upload-preview {
  max-width: 100%;
  max-height: 280px;
  border-radius: 8px;
  object-fit: contain;
}
.upload-preview-video {
  max-width: 100%;
  max-height: 240px;
  border-radius: 8px;
}

/* ==================== 提交按钮 ==================== */
.submit-btn {
  width: 100%;
  margin-top: 8px;
  border-radius: 10px;
  background: var(--btn-gradient);
  border: none;
}
.submit-btn:hover {
  background: var(--btn-hover);
}

/* ==================== 进度 & 错误 ==================== */
.progress-section {
  padding: 16px 0;
}
.progress-text {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 8px;
  text-align: center;
}
.error-section {
  margin-bottom: 16px;
}

/* ==================== 视频结果 ==================== */
.video-result {
  margin-top: 8px;
}
.video-player {
  width: 100%;
  border-radius: 12px;
  background: #000;
  max-height: 500px;
}
.video-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

/* ==================== 响应式 ==================== */
@media (max-width: 1024px) {
  .main-layout {
    grid-template-columns: 1fr;
  }
}
</style>
