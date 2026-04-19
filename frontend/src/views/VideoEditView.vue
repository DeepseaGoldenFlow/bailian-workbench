<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎬 视频编辑</h2>
      <p class="page-sub">对已有视频进行智能重绘、局部编辑与画面延展</p>
    </div>

    <div class="main-layout">
      <!-- 左侧参数区 -->
      <div class="param-panel glass-card">
        <h3 class="panel-title"><el-icon><EditPen /></el-icon> 编辑参数</h3>

        <!-- 原视频上传 -->
        <div class="field">
          <label class="field-label">原视频 <span class="required">*</span></label>
          <el-upload
            action=""
            :auto-upload="false"
            :show-file-list="false"
            accept="video/mp4"
            @change="handleVideoUpload"
            drag
            class="upload-dragger"
          >
            <div v-if="!videoUrl" class="upload-placeholder">
              <el-icon class="upload-icon"><VideoCamera /></el-icon>
              <p>拖拽视频到此处</p>
              <p class="upload-hint">支持 MP4 格式</p>
            </div>
            <video v-else :src="videoUrl" controls class="upload-preview-video" />
          </el-upload>
          <el-button v-if="videoUrl" link type="danger" size="small" @click="clearVideo">移除视频</el-button>
        </div>

        <!-- 参考图上传 -->
        <div class="field">
          <label class="field-label">参考图片 <span class="optional">（可选，用于风格重绘，最多 5 张）</span></label>
          <el-upload
            v-model:file-list="refFileList"
            action=""
            :auto-upload="false"
            list-type="picture-card"
            :limit="5"
            accept="image/jpeg,image/png,image/webp"
            :on-exceed="handleRefExceed"
            :on-remove="handleRefRemove"
            :before-upload="beforeRefUpload"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <p class="field-desc">上传风格参考图片，AI 将参考其视觉风格对视频进行重绘</p>
        </div>

        <!-- 编辑指令 -->
        <div class="field">
          <label class="field-label">编辑指令 <span class="required">*</span></label>
          <el-input
            v-model="prompt"
            type="textarea"
            :rows="4"
            resize="none"
            placeholder="描述你想要的编辑效果，例如：将视频画面风格改为水彩画风格，或延展画面右侧的风景"
          />
          <p class="field-desc">用自然语言描述你希望视频产生的变化</p>
        </div>

        <!-- 掩码图上传 -->
        <div class="field">
          <label class="field-label">掩码图片 <span class="optional">（可选，指定修改区域）</span></label>
          <el-upload
            v-model:file-list="maskFileList"
            action=""
            :auto-upload="false"
            list-type="picture-card"
            :limit="1"
            accept="image/jpeg,image/png,image/webp"
            :on-exceed="handleMaskExceed"
            :on-remove="handleMaskRemove"
            :before-upload="beforeMaskUpload"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <p class="field-desc">上传掩码指定需要编辑的区域，白色区域将被修改，黑色区域保持原样</p>
        </div>

        <!-- 分辨率 -->
        <div class="field">
          <label class="field-label">输出分辨率</label>
          <el-select v-model="resolution" class="full-width">
            <el-option label="480P（快速，适合预览）" value="480P" />
            <el-option label="720P（高清，推荐）" value="720P" />
            <el-option label="1080P（超清，画质最佳）" value="1080P" />
          </el-select>
          <p class="field-desc">选择输出视频的分辨率，分辨率越高耗时越长</p>
        </div>

        <el-button
          @click="submitVideoEdit"
          type="primary"
          :loading="loading"
          class="submit-btn"
          size="large"
        >
          <el-icon><VideoCamera /></el-icon>
          开始编辑
        </el-button>
      </div>

      <!-- 右侧结果区 -->
      <div class="result-panel glass-card">
        <h3 class="panel-title"><el-icon><VideoCameraFilled /></el-icon> 编辑结果</h3>

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
              <el-icon><RefreshRight /></el-icon> 重新编辑
            </el-button>
          </div>
        </div>

        <el-empty v-if="!resultVideoUrl && !loading && !error" description="上传视频并设置参数，开始编辑吧 🎥" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { EditPen, VideoCamera, VideoCameraFilled, Download, RefreshRight, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// ==================== 状态 ====================
const videoUrl = ref('')
const videoBase64 = ref('')
const videoFile = ref(null)

const prompt = ref('')
const resolution = ref('720P')

// 参考图
const refFileList = ref([])
const refImageFiles = ref([])

// 掩码图
const maskFileList = ref([])
const maskImageFile = ref(null)

// 加载状态
const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const error = ref('')
const resultVideoUrl = ref('')

let pollTimer = null

// ==================== 视频上传 ====================
function handleVideoUpload(file) {
  const raw = file.raw
  if (!raw) return
  const isMp4 = raw.type === 'video/mp4'
  if (!isMp4) {
    ElMessage.error('仅支持 MP4 格式视频')
    return
  }
  // 限制 500MB
  if (raw.size > 500 * 1024 * 1024) {
    ElMessage.error('视频大小不能超过 500MB')
    return
  }
  videoFile.value = raw
  const reader = new FileReader()
  reader.onload = (e) => {
    videoBase64.value = e.target.result
    videoUrl.value = e.target.result
  }
  reader.readAsDataURL(raw)
}

function clearVideo() {
  videoUrl.value = ''
  videoBase64.value = ''
  videoFile.value = null
}

// ==================== 参考图上传 ====================
function beforeRefUpload(file) {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG / PNG / WebP 格式')
    return false
  }
  refImageFiles.value.push(file)
  return false
}

function handleRefExceed() {
  ElMessage.warning('最多上传 5 张参考图片')
}

function handleRefRemove(file) {
  const target = file.raw || file
  const idx = refImageFiles.value.indexOf(target)
  if (idx > -1) refImageFiles.value.splice(idx, 1)
}

// ==================== 掩码图上传 ====================
function beforeMaskUpload(file) {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG / PNG / WebP 格式')
    return false
  }
  maskImageFile.value = file
  return false
}

function handleMaskExceed() {
  ElMessage.warning('仅支持上传 1 张掩码图片')
}

function handleMaskRemove() {
  maskImageFile.value = null
}

// ==================== 提交编辑 ====================
async function submitVideoEdit() {
  if (!videoBase64.value) {
    return ElMessage.warning('请上传原视频')
  }
  if (!prompt.value.trim()) {
    return ElMessage.warning('请输入编辑指令')
  }

  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'
  error.value = ''
  resultVideoUrl.value = ''

  try {
    // 构建 ref_images
    const refImages = []
    for (const file of refImageFiles.value) {
      const b64 = await fileToBase64(file)
      refImages.push({ image: b64 })
    }

    // 构建请求体
    const input = {
      input_video: videoBase64.value,
      prompt: prompt.value.trim(),
    }
    if (refImages.length > 0) {
      input.ref_images = refImages
    }
    if (maskImageFile.value) {
      input.mask = await fileToBase64(maskImageFile.value)
    }

    const parameters = {
      resolution: resolution.value,
    }

    const res = await fetch('/api/video/edit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: 'wan2.7-videoedit',
        input,
        parameters,
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
    progressText.value = '任务已提交，正在编辑中（通常 3-10 分钟）...'

    // 开始轮询
    await pollTask(taskId)
  } catch (e) {
    error.value = `编辑失败: ${e.message}`
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

// ==================== 轮询任务 ====================
function pollTask(taskId) {
  return new Promise((resolve, reject) => {
    let pollCount = 0
    const maxPolls = 200 // 最长 200 * 5s = ~16 分钟

    const doPoll = async () => {
      try {
        await new Promise(r => { pollTimer = setTimeout(r, 5000) })
        pollCount++

        const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
        if (!pollRes.ok) throw new Error('轮询任务状态失败')

        const pollData = await pollRes.json()
        const status = pollData.output?.task_status || pollData.status

        progress.value = Math.min(15 + pollCount * 2, 95)
        progressText.value = `编辑中... 已等待 ${pollCount * 5}秒 [${status}]`

        if (status === 'SUCCEEDED') {
          progress.value = 100
          progressText.value = '视频编辑完成！'

          let videoUrl = pollData.output?.video_url || pollData.output?.results?.[0]?.url
          if (!videoUrl) throw new Error('未获取到视频 URL')

          // 处理本地存储路径
          if (videoUrl.startsWith('/data/bailian/storage')) {
            videoUrl = videoUrl.replace('/data/bailian/storage', '/api/files')
          }

          resultVideoUrl.value = videoUrl
          ElMessage.success('🎉 视频编辑成功！')
          resolve()
          return
        }

        if (status === 'FAILED') {
          const errMsg = pollData.output?.message || pollData.message || '未知错误'
          throw new Error(errMsg)
        }

        if (pollCount >= maxPolls) {
          throw new Error('轮询超时，请稍后到任务列表查看结果')
        }

        pollTimer = setTimeout(doPoll, 0)
      } catch (e) {
        reject(e)
      }
    }

    doPoll()
  })
}

// ==================== 工具函数 ====================
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

function downloadVideo() {
  if (resultVideoUrl.value) {
    const a = document.createElement('a')
    a.href = resultVideoUrl.value
    a.download = `video_edit_${Date.now()}.mp4`
    a.click()
  }
}

function resetForm() {
  resultVideoUrl.value = ''
  prompt.value = ''
  clearVideo()
  refFileList.value = []
  refImageFiles.value = []
  maskFileList.value = []
  maskImageFile.value = null
}

// 组件卸载时清理定时器
onUnmounted(() => {
  if (pollTimer) clearTimeout(pollTimer)
})
</script>

<style scoped>
/* ==================== 布局 ==================== */
.page-container { max-width: 1400px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.main-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; align-items: start; }

/* ==================== 卡片 ==================== */
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
.field { margin-bottom: 16px; }
.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}
.required { color: #ef4444; margin-left: 2px; }
.optional { font-weight: 400; color: var(--text-secondary); font-size: 12px; }
.field-desc {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 4px;
  line-height: 1.4;
}
.full-width { width: 100%; }

/* ==================== 上传区域 ==================== */
.upload-dragger { width: 100%; }
.upload-dragger :deep(.el-upload) { width: 100%; }
.upload-dragger :deep(.el-upload-dragger) {
  background: transparent;
  border: 2px dashed var(--card-border);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s;
}
.upload-dragger :deep(.el-upload-dragger:hover) {
  border-color: var(--gradient-start);
}
.upload-placeholder { text-align: center; }
.upload-icon { font-size: 36px; color: var(--text-secondary); margin-bottom: 8px; }
.upload-placeholder p { font-size: 13px; color: var(--text-secondary); margin: 0; }
.upload-hint { font-size: 11px !important; margin-top: 4px !important; }
.upload-preview-video {
  max-width: 100%;
  max-height: 200px;
  border-radius: 8px;
}

/* 参考图 / 掩码上传 */
.field :deep(.el-upload--picture-card) {
  width: 80px;
  height: 80px;
}
.field :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 80px;
  height: 80px;
}

/* ==================== 提交按钮 ==================== */
.submit-btn {
  width: 100%;
  margin-top: 12px;
  border-radius: 10px;
  background: var(--btn-gradient);
  border: none;
}
.submit-btn:hover { background: var(--btn-hover); }

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
.video-player {
  width: 100%;
  border-radius: 12px;
  background: #000;
  max-height: 450px;
}
.video-actions { display: flex; gap: 8px; margin-top: 12px; }

/* ==================== 响应式 ==================== */
@media (max-width: 1024px) {
  .main-layout { grid-template-columns: 1fr; }
}
@media (max-width: 640px) {
  .page-header h2 { font-size: 20px; }
}
</style>
