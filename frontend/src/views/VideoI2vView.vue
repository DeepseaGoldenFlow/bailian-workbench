<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎬 图生视频</h2>
      <p class="page-sub">Wan 2.7 — 基于图片生成高质量动态视频，支持首帧、首尾帧和视频续写三种模式</p>
    </div>

    <div class="main-layout">
      <!-- Left: Input Panel -->
      <div class="input-panel glass-card">
        <el-tabs v-model="activeTab" class="mode-tabs" stretch>
          <!-- Tab 1: 首帧生视频 -->
          <el-tab-pane label="首帧生视频" name="first-frame">
            <div class="tab-content">
              <div class="field">
                <label class="field-label">首帧图片 <span class="required">*</span></label>
                <el-upload
                  action=""
                  :auto-upload="false"
                  :show-file-list="false"
                  accept="image/*"
                  @change="(f) => handleImageUpload(f, 'firstFrame')"
                  drag
                  class="upload-dragger"
                >
                  <div v-if="!form.firstFrameUrl" class="upload-placeholder">
                    <el-icon class="upload-icon"><UploadFilled /></el-icon>
                    <p>拖拽首帧图片到此处</p>
                    <p class="upload-hint">或点击选择文件</p>
                  </div>
                  <img v-else :src="form.firstFrameUrl" class="upload-preview" />
                </el-upload>
                <el-button v-if="form.firstFrameUrl" link type="danger" size="small" @click="clearImage('firstFrame')">移除</el-button>
              </div>

              <div class="field">
                <label class="field-label">Prompt <span class="required">*</span></label>
                <el-input
                  v-model="form.prompt"
                  type="textarea"
                  :rows="4"
                  resize="none"
                  placeholder="描述你期望的视频内容，例如: 人物缓缓转头看向镜头，微风吹动头发，背景虚化，电影质感"
                />
              </div>
            </div>
          </el-tab-pane>

          <!-- Tab 2: 首尾帧生视频 -->
          <el-tab-pane label="首尾帧生视频" name="first-last-frame">
            <div class="tab-content">
              <div class="field">
                <label class="field-label">首帧图片 <span class="required">*</span></label>
                <el-upload action="" :auto-upload="false" :show-file-list="false" accept="image/*"
                  @change="(f) => handleImageUpload(f, 'firstFrame')" drag class="upload-dragger">
                  <div v-if="!form.firstFrameUrl" class="upload-placeholder">
                    <el-icon class="upload-icon"><UploadFilled /></el-icon>
                    <p>上传首帧图片</p>
                  </div>
                  <img v-else :src="form.firstFrameUrl" class="upload-preview" />
                </el-upload>
                <el-button v-if="form.firstFrameUrl" link type="danger" size="small" @click="clearImage('firstFrame')">移除</el-button>
              </div>

              <div class="field">
                <label class="field-label">尾帧图片 <span class="required">*</span></label>
                <el-upload action="" :auto-upload="false" :show-file-list="false" accept="image/*"
                  @change="(f) => handleImageUpload(f, 'lastFrame')" drag class="upload-dragger">
                  <div v-if="!form.lastFrameUrl" class="upload-placeholder">
                    <el-icon class="upload-icon"><UploadFilled /></el-icon>
                    <p>上传尾帧图片</p>
                  </div>
                  <img v-else :src="form.lastFrameUrl" class="upload-preview" />
                </el-upload>
                <el-button v-if="form.lastFrameUrl" link type="danger" size="small" @click="clearImage('lastFrame')">移除</el-button>
              </div>

              <div class="field">
                <label class="field-label">Prompt <span class="optional">可选</span></label>
                <el-input v-model="form.prompt" type="textarea" :rows="3" resize="none"
                  placeholder="描述首帧到尾帧之间的过渡内容" />
              </div>
            </div>
          </el-tab-pane>

          <!-- Tab 3: 视频续写 -->
          <el-tab-pane label="视频续写" name="video-continue">
            <div class="tab-content">
              <div class="field">
                <label class="field-label">前序视频 <span class="required">*</span></label>
                <el-upload action="" :auto-upload="false" :show-file-list="false" accept="video/*"
                  @change="handleVideoUpload" drag class="upload-dragger">
                  <div v-if="!form.prevVideoUrl" class="upload-placeholder">
                    <el-icon class="upload-icon"><VideoCamera /></el-icon>
                    <p>拖拽前序视频到此处</p>
                    <p class="upload-hint">支持 MP4, WebM 格式</p>
                  </div>
                  <video v-else :src="form.prevVideoUrl" controls class="upload-preview-video" />
                </el-upload>
                <el-button v-if="form.prevVideoUrl" link type="danger" size="small" @click="form.prevVideoUrl = ''; form.prevVideoBase64 = ''">移除</el-button>
              </div>

              <div class="field">
                <label class="field-label">Prompt <span class="optional">可选</span></label>
                <el-input v-model="form.prompt" type="textarea" :rows="3" resize="none"
                  placeholder="描述续写视频的内容方向" />
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>

        <!-- Parameters Panel -->
        <div class="params-section">
          <h4 class="params-title">参数设置</h4>
          <div class="params-grid">
            <div class="field">
              <label class="field-label">分辨率</label>
              <el-select v-model="form.resolution" class="full-width">
                <el-option label="720P (1280×720)" value="720P" />
                <el-option label="1080P (1920×1080)" value="1080P" />
              </el-select>
            </div>
            <div class="field">
              <label class="field-label">时长</label>
              <el-select v-model="form.duration" class="full-width">
                <el-option label="5 秒" value="5" />
                <el-option label="10 秒" value="10" />
              </el-select>
            </div>
            <div class="field">
              <div class="param-row">
                <label class="field-label">自动配音</label>
                <el-switch v-model="form.auto_audio" />
              </div>
            </div>
          </div>
        </div>

        <el-button
          @click="submitGenerate"
          type="primary"
          :loading="loading"
          class="submit-btn"
          size="large"
        >
          <el-icon><VideoCamera /></el-icon>
          开始生成
        </el-button>
      </div>

      <!-- Right: Result Panel -->
      <div class="result-panel glass-card">
        <h3 class="panel-title"><el-icon><VideoCameraFilled /></el-icon> 生成结果</h3>

        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">{{ progressText }}</p>
        </div>

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

        <el-empty v-if="!resultVideoUrl && !loading" description="生成结果将在这里展示 🎥" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { UploadFilled, VideoCamera, VideoCameraFilled, Download, RefreshRight } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('first-frame')
const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const resultVideoUrl = ref('')

const form = ref({
  firstFrameUrl: '',
  firstFrameBase64: '',
  lastFrameUrl: '',
  lastFrameBase64: '',
  prevVideoUrl: '',
  prevVideoBase64: '',
  prompt: '',
  resolution: '720P',
  duration: '5',
  auto_audio: false,
})

const handleImageUpload = (file, target) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    if (target === 'firstFrame') {
      form.value.firstFrameBase64 = e.target.result
      form.value.firstFrameUrl = e.target.result
    } else {
      form.value.lastFrameBase64 = e.target.result
      form.value.lastFrameUrl = e.target.result
    }
  }
  reader.readAsDataURL(file.raw)
}

const handleVideoUpload = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    form.value.prevVideoBase64 = e.target.result
    form.value.prevVideoUrl = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

const clearImage = (target) => {
  if (target === 'firstFrame') {
    form.value.firstFrameUrl = ''
    form.value.firstFrameBase64 = ''
  } else {
    form.value.lastFrameUrl = ''
    form.value.lastFrameBase64 = ''
  }
}

const validateForm = () => {
  if (activeTab.value === 'first-frame') {
    if (!form.value.firstFrameUrl) { ElMessage.warning('请上传首帧图片'); return false }
    if (!form.value.prompt) { ElMessage.warning('请输入 Prompt'); return false }
  } else if (activeTab.value === 'first-last-frame') {
    if (!form.value.firstFrameUrl) { ElMessage.warning('请上传首帧图片'); return false }
    if (!form.value.lastFrameUrl) { ElMessage.warning('请上传尾帧图片'); return false }
  } else if (activeTab.value === 'video-continue') {
    if (!form.value.prevVideoUrl) { ElMessage.warning('请上传前序视频'); return false }
  }
  return true
}

const submitGenerate = async () => {
  if (!validateForm()) return

  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'
  resultVideoUrl.value = ''

  try {
    // Build the payload based on active tab
    let input = {}
    if (activeTab.value === 'first-frame') {
      input = {
        prompt: form.value.prompt,
        img_url: form.value.firstFrameUrl,
      }
    } else if (activeTab.value === 'first-last-frame') {
      input = {
        prompt: form.value.prompt,
        first_frame_url: form.value.firstFrameUrl,
        last_frame_url: form.value.lastFrameUrl,
      }
    } else {
      input = {
        prompt: form.value.prompt,
        video_url: form.value.prevVideoUrl,
      }
    }

    const parameters = {
      resolution: form.value.resolution,
      duration: parseInt(form.value.duration),
      auto_audio: form.value.auto_audio,
    }

    const res = await fetch('/api/video/i2v', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        mode: activeTab.value,
        input,
        parameters,
      }),
    })

    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.message || err.error || '提交失败')
    }

    const data = await res.json()
    const taskId = data.task_id

    if (!taskId) throw new Error('未获取到 task_id')

    progress.value = 15
    progressText.value = '视频生成中，请耐心等待（通常 2-5 分钟）...'

    // Poll for task status
    let status = 'PENDING'
    let pollCount = 0
    const maxPolls = 200

    while (status !== 'SUCCEEDED' && status !== 'FAILED' && pollCount < maxPolls) {
      await new Promise(r => setTimeout(r, 5000))
      pollCount++

      const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
      if (!pollRes.ok) throw new Error('轮询任务状态失败')

      const pollData = await pollRes.json()
      status = pollData.output?.task_status || pollData.status

      progress.value = Math.min(15 + pollCount * 2, 95)
      progressText.value = `生成中... 已等待 ${pollCount * 5}秒 [${status}]`

      if (status === 'SUCCEEDED') {
        progress.value = 100
        progressText.value = '视频生成完成！'

        let videoUrl = pollData.output?.video_url || pollData.output?.results?.[0]?.url
        if (videoUrl) {
          if (videoUrl.startsWith('/data/bailian/storage')) {
            videoUrl = videoUrl.replace('/data/bailian/storage', '/api/files')
          }
          // Also handle /storage/videos/ pattern
          if (videoUrl.startsWith('/storage/videos/')) {
            // keep as is - the backend should serve it
          }
          resultVideoUrl.value = videoUrl
          ElMessage.success('🎉 视频生成成功！')
        } else {
          throw new Error('未获取到视频 URL')
        }
      } else if (status === 'FAILED') {
        throw new Error(pollData.output?.message || pollData.message || '生成失败')
      }
    }

    if (pollCount >= maxPolls) {
      throw new Error('轮询超时，请稍后到任务列表查看结果')
    }
  } catch (e) {
    ElMessage.error('生成失败: ' + e.message)
  } finally {
    loading.value = false
    setTimeout(() => { progress.value = 0; progressText.value = '' }, 2000)
  }
}

const downloadVideo = () => {
  if (resultVideoUrl.value) {
    const a = document.createElement('a')
    a.href = resultVideoUrl.value
    a.download = 'generated_video.mp4'
    a.click()
  }
}

const resetForm = () => {
  resultVideoUrl.value = ''
  form.value.prompt = ''
  form.value.firstFrameUrl = ''
  form.value.firstFrameBase64 = ''
  form.value.lastFrameUrl = ''
  form.value.lastFrameBase64 = ''
  form.value.prevVideoUrl = ''
  form.value.prevVideoBase64 = ''
}
</script>

<style scoped>
.page-container { max-width: 1400px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.main-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; align-items: start; }

.glass-card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(10px);
}

.panel-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }

/* Tabs */
.mode-tabs { margin-bottom: 16px; }

/* Fields */
.field { margin-bottom: 16px; }
.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.required { color: #ef4444; margin-left: 2px; }
.optional { font-weight: 400; color: var(--text-secondary); font-size: 12px; }
.param-row { display: flex; align-items: center; justify-content: space-between; }
.full-width { width: 100%; }

/* Upload */
.upload-dragger { width: 100%; }
.upload-dragger :deep(.el-upload) { width: 100%; }
.upload-dragger :deep(.el-upload-dragger) {
  background: transparent;
  border: 2px dashed var(--card-border);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s;
}
.upload-dragger :deep(.el-upload-dragger:hover) { border-color: var(--gradient-start); }
.upload-placeholder { text-align: center; }
.upload-icon { font-size: 36px; color: var(--text-secondary); margin-bottom: 8px; }
.upload-placeholder p { font-size: 13px; color: var(--text-secondary); margin: 0; }
.upload-hint { font-size: 11px !important; margin-top: 4px !important; }
.upload-preview { max-width: 100%; max-height: 200px; border-radius: 8px; object-fit: contain; }
.upload-preview-video { max-width: 100%; max-height: 200px; border-radius: 8px; }

/* Params */
.params-section { margin-top: 8px; padding-top: 16px; border-top: 1px solid var(--card-border); }
.params-title { font-size: 14px; font-weight: 600; margin-bottom: 12px; color: var(--text-primary); }
.params-grid { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 12px; }

/* Submit */
.submit-btn { width: 100%; margin-top: 12px; border-radius: 10px; background: var(--btn-gradient); border: none; }
.submit-btn:hover { background: var(--btn-hover); }

/* Progress */
.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

/* Video Result */
.video-result { margin-top: 8px; }
.video-player { width: 100%; border-radius: 12px; background: #000; max-height: 450px; }
.video-actions { display: flex; gap: 8px; margin-top: 12px; }

@media (max-width: 1024px) {
  .main-layout { grid-template-columns: 1fr; }
  .params-grid { grid-template-columns: 1fr 1fr; }
}
@media (max-width: 640px) {
  .params-grid { grid-template-columns: 1fr; }
}
</style>
