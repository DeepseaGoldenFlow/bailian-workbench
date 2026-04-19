<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎞️ 参考生视频</h2>
      <p class="page-sub">基于参考图片和视频，生成角色一致的高质量视频</p>
    </div>

    <div class="page-grid">
      <!-- ==================== 左侧参数区 ==================== -->
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><VideoCamera /></el-icon> 参数配置</h3>

        <!-- 参考图片上传 -->
        <div class="field">
          <label class="field-label">参考图片 <span class="optional">（可选，保持人物/物体特征）</span></label>
          <el-upload
            v-model:file-list="refImageFileList"
            list-type="picture-card"
            :limit="1"
            accept="image/jpeg,image/png,image/webp"
            :on-exceed="handleImageExceed"
            :on-remove="handleImageRemove"
            :before-upload="beforeImageUpload"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <p class="field-desc">上传参考图片以保持人物或物体特征一致性，支持 JPG / PNG / WebP</p>
        </div>

        <!-- 参考视频上传 -->
        <div class="field">
          <label class="field-label">参考视频 <span class="optional">（可选，参考动作和风格）</span></label>
          <el-upload
            v-model:file-list="refVideoFileList"
            list-type="picture-card"
            :limit="1"
            accept="video/mp4,video/quicktime,video/webm"
            :on-exceed="handleVideoExceed"
            :on-remove="handleVideoRemove"
            :before-upload="beforeVideoUpload"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <p class="field-desc">上传参考视频以参考其动作和风格，支持 MP4 / MOV / WebM，最大 50MB</p>
        </div>

        <!-- 画面描述 -->
        <div class="field">
          <label class="field-label">画面描述（Prompt） <span class="optional">（可选）</span></label>
          <el-input
            v-model="form.prompt"
            type="textarea"
            :rows="4"
            resize="none"
            placeholder="例如：人物在阳光明媚的公园里散步，微风吹拂，电影级光影效果"
          />
          <p class="field-desc">用自然语言描述你想要的视频画面内容，越详细效果越好</p>
        </div>

        <!-- 分辨率 -->
        <div class="field">
          <label class="field-label">分辨率</label>
          <el-select v-model="form.resolution" class="full-width">
            <el-option label="480P（标清，生成最快）" value="480P" />
            <el-option label="720P（高清，推荐）" value="720P" />
            <el-option label="1080P（全高清，画质最佳）" value="1080P" />
          </el-select>
          <p class="field-desc">选择生成视频的分辨率，分辨率越高生成时间越长</p>
        </div>

        <!-- 时长 -->
        <div class="field">
          <label class="field-label">视频时长</label>
          <el-select v-model="form.duration" class="full-width">
            <el-option label="5 秒（快速预览）" value="5" />
            <el-option label="10 秒（完整片段）" value="10" />
          </el-select>
          <p class="field-desc">选择生成视频的时长</p>
        </div>

        <el-button
          @click="submitTask"
          type="primary"
          :loading="loading"
          class="generate-btn"
          size="large"
        >
          <el-icon><VideoPlay /></el-icon> 开始生成
        </el-button>
      </div>

      <!-- ==================== 右侧结果区 ==================== -->
      <div class="result-card glass-card">
        <h3 class="card-title"><el-icon><Film /></el-icon> 生成结果</h3>

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
          <video :src="videoUrl" controls preload="metadata" class="video-player" />
          <div class="video-actions">
            <el-button type="primary" link @click="downloadVideo">
              <el-icon><Download /></el-icon> 下载视频
            </el-button>
          </div>
        </div>

        <el-empty
          v-if="!videoUrl && !loading && !error"
          description="上传参考图片或视频，开始创作你的视频吧 🎬"
        />
      </div>
    </div>

    <!-- 视频预览弹窗 -->
    <el-dialog v-model="showPreview" title="预览视频" width="80%" class="preview-dialog">
      <video :src="previewUrl" controls style="width: 100%; border-radius: 8px;" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { VideoCamera, VideoPlay, Film, Plus, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// ==================== 表单数据 ====================
const form = ref({
  prompt: '',
  resolution: '720P',
  duration: '5',
})

// ==================== 参考图片上传 ====================
const refImageFileList = ref([])
const refImageFile = ref(null)

function beforeImageUpload(file) {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG / PNG / WebP 格式')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  refImageFile.value = file
  return false // 阻止默认上传，我们自行处理
}

function handleImageExceed() {
  ElMessage.warning('最多上传 1 张参考图片')
}

function handleImageRemove() {
  refImageFile.value = null
}

// ==================== 参考视频上传 ====================
const refVideoFileList = ref([])
const refVideoFile = ref(null)

function beforeVideoUpload(file) {
  const isVideo = ['video/mp4', 'video/quicktime', 'video/webm'].includes(file.type)
  if (!isVideo) {
    ElMessage.error('仅支持 MP4 / MOV / WebM 格式')
    return false
  }
  const isLt50M = file.size / 1024 / 1024 < 50
  if (!isLt50M) {
    ElMessage.error('视频大小不能超过 50MB')
    return false
  }
  refVideoFile.value = file
  return false
}

function handleVideoExceed() {
  ElMessage.warning('最多上传 1 个参考视频')
}

function handleVideoRemove() {
  refVideoFile.value = null
}

// ==================== 提交与轮询 ====================
const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const error = ref('')
const videoUrl = ref('')
const showPreview = ref(false)
const previewUrl = ref('')
let pollTimer = null

/** 文件转 base64 data URL */
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

/** 提交任务 */
async function submitTask() {
  // 校验：至少上传图片或视频之一
  if (!refImageFile.value && !refVideoFile.value) {
    return ElMessage.warning('请至少上传参考图片或参考视频之一')
  }

  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'
  error.value = ''
  videoUrl.value = ''

  try {
    const input = {}
    input.prompt = form.value.prompt.trim()

    // 转为 base64
    if (refImageFile.value) {
      input.ref_image = await fileToBase64(refImageFile.value)
    }
    if (refVideoFile.value) {
      input.ref_video = await fileToBase64(refVideoFile.value)
    }

    const parameters = {
      resolution: form.value.resolution,
      duration: form.value.duration,
    }

    const res = await fetch('/api/video/ref', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: 'wan2.7-r2v',
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

    progress.value = 20
    progressText.value = '任务已提交，正在生成视频...'

    // 开始轮询
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

/** 轮询任务状态 */
function pollTask(taskId) {
  return new Promise((resolve, reject) => {
    let pollCount = 0
    const maxPolls = 240 // 最长 240 * 5s = 20 分钟（视频生成通常较长）

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
          progressText.value = '视频生成完成！'

          // 从 results 中取视频 URL
          const results = pollData.output?.results || []
          if (results.length > 0 && results[0].url) {
            videoUrl.value = results[0].url
          } else if (pollData.output?.video_url) {
            videoUrl.value = pollData.output.video_url
          } else {
            throw new Error('未返回视频地址')
          }

          ElMessage.success('🎉 视频生成成功！')
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

        progress.value = Math.min(20 + pollCount * 2, 95)
        progressText.value = `视频生成中... 已等待 ${pollCount * 5}秒`

        pollTimer = setTimeout(doPoll, 5000)
      } catch (e) {
        reject(e)
      }
    }

    doPoll()
  })
}

/** 下载视频 */
function downloadVideo() {
  if (!videoUrl.value) return
  const a = document.createElement('a')
  a.href = videoUrl.value
  a.download = `video_${Date.now()}.mp4`
  a.target = '_blank'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

/** 组件卸载时清理 */
onUnmounted(() => {
  if (pollTimer) clearTimeout(pollTimer)
})
</script>

<style scoped>
/* ==================== 布局 ==================== */
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 4px;
}

.page-sub {
  font-size: 14px;
  color: var(--text-secondary, #a0a0b0);
}

.page-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

/* ==================== 玻璃拟态卡片 ==================== */
.glass-card {
  background: var(--card-bg, rgba(30, 30, 45, 0.7));
  border: 1px solid var(--card-border, rgba(255, 255, 255, 0.08));
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-primary, #e8e8f0);
}

/* ==================== 表单字段 ==================== */
.field {
  margin-bottom: 16px;
}

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary, #e8e8f0);
  margin-bottom: 6px;
}

.optional {
  font-weight: 400;
  color: var(--text-secondary, #a0a0b0);
  font-size: 12px;
}

.field-desc {
  font-size: 11px;
  color: var(--text-secondary, #a0a0b0);
  margin-top: 4px;
  line-height: 1.4;
}

.full-width {
  width: 100%;
}

/* 上传区域 */
.field :deep(.el-upload--picture-card) {
  width: 80px;
  height: 80px;
}

.field :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 80px;
  height: 80px;
}

.generate-btn {
  width: 100%;
  margin-top: 8px;
  border-radius: 10px;
}

/* ==================== 进度 & 错误 ==================== */
.progress-section {
  padding: 16px 0;
}

.progress-text {
  font-size: 13px;
  color: var(--text-secondary, #a0a0b0);
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
  max-height: 480px;
  object-fit: contain;
}

.video-actions {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

/* ==================== 预览弹窗 ==================== */
.preview-dialog :deep(.el-dialog) {
  background: var(--card-bg, rgba(30, 30, 45, 0.95));
}

/* ==================== 响应式 ==================== */
@media (max-width: 768px) {
  .page-grid {
    grid-template-columns: 1fr;
  }
}
</style>
