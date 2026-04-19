<template>
  <div class="page-container">
    <div class="page-header">
      <h2>✍️ 文生视频</h2>
      <p class="page-sub">通过文字描述直接生成高质量视频</p>
    </div>

    <div class="main-layout">
      <div class="input-panel glass-card">
        <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="万相 2.7 (最新，多镜头叙事)" value="wan2.7-t2v" />
            <el-option label="万相 2.2 Plus (稳定)" value="wan2.2-t2v-plus" />
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">视频描述 (Prompt) <span class="required">*</span></label>
          <el-input v-model="form.prompt" type="textarea" :rows="5" resize="none"
            placeholder="例如: 海浪拍打在礁石上, 夕阳余晖映照海面, 一只海鸥从画面中飞过, 电影质感" />
        </div>

        <div class="field">
          <label class="field-label">视频时长</label>
          <el-radio-group v-model="form.duration">
            <el-radio-button value="5">5 秒</el-radio-button>
            <el-radio-button value="10">10 秒</el-radio-button>
          </el-radio-group>
        </div>

        <div class="field">
          <label class="field-label">分辨率</label>
          <el-select v-model="form.resolution" class="full-width">
            <el-option label="720P (1280×720)" value="720P" />
            <el-option label="1080P (1920×1080)" value="1080P" />
          </el-select>
        </div>

        <div class="field">
          <div class="param-row">
            <label class="field-label">智能优化提示词</label>
            <el-switch v-model="form.auto_prompt" active-text="开启" inactive-text="关闭" />
          </div>
          <p class="field-desc">自动增强描述以获得更好的视频效果</p>
        </div>

        <el-button @click="generate" type="primary" :loading="loading" class="submit-btn" size="large">
          <el-icon><VideoCamera /></el-icon> 开始生成
        </el-button>
      </div>

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

        <el-empty v-if="!resultVideoUrl && !loading" description="还没有生成视频，开始创作吧 🎥" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditPen, VideoCamera, VideoCameraFilled, Download, RefreshRight } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({
  model: 'wan2.7-t2v',
  prompt: '',
  duration: '5',
  resolution: '720P',
  auto_prompt: false,
})

const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const resultVideoUrl = ref('')

const generate = async () => {
  if (!form.value.prompt) { ElMessage.warning('请输入视频描述'); return }

  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'
  resultVideoUrl.value = ''

  try {
    const input = { prompt: form.value.prompt }
    const parameters = {
      duration: parseInt(form.value.duration),
      resolution: form.value.resolution,
    }
    if (form.value.auto_prompt) parameters.auto_prompt = true

    const res = await fetch('/api/v1/services/aigc/video-generation/video-synthesis', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ model: form.value.model, input, parameters }),
    })

    if (!res.ok) { const e = await res.json(); throw new Error(e.message || '提交失败') }
    const data = await res.json()
    const taskId = data.output?.task_id
    if (!taskId) throw new Error('未获取到 task_id')

    progress.value = 15
    progressText.value = '视频生成中，请耐心等待（通常 2-5 分钟）...'

    let status = 'PENDING'
    let pollCount = 0
    while (status !== 'SUCCEEDED' && status !== 'FAILED' && pollCount < 200) {
      await new Promise(r => setTimeout(r, 5000))
      pollCount++

      const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
      const pollData = await pollRes.json()
      status = pollData.output?.task_status

      progress.value = Math.min(15 + pollCount * 2, 95)
      progressText.value = `生成中... 已等待 ${pollCount * 5}秒`

      if (status === 'SUCCEEDED') {
        progress.value = 100
        progressText.value = '视频生成完成！'
        let videoUrl = pollData.output?.video_url
        if (videoUrl) {
          if (videoUrl.startsWith('/data/bailian/storage')) {
            videoUrl = videoUrl.replace('/data/bailian/storage', '/api/files')
          }
          resultVideoUrl.value = videoUrl
          ElMessage.success('🎉 视频生成成功！')
        }
      } else if (status === 'FAILED') {
        throw new Error(pollData.output?.message || '生成失败')
      }
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
}
</script>

<style scoped>
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.main-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; align-items: start; }
.glass-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; padding: 20px; backdrop-filter: blur(10px); }
.card-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }
.panel-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }

.field { margin-bottom: 16px; }
.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.required { color: #ef4444; margin-left: 2px; }
.param-row { display: flex; align-items: center; justify-content: space-between; }
.field-desc { font-size: 11px; color: var(--text-secondary); margin-top: 4px; }
.full-width { width: 100%; }

.submit-btn { width: 100%; margin-top: 8px; border-radius: 10px; background: var(--btn-gradient); border: none; }
.submit-btn:hover { background: var(--btn-hover); }

.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

.video-result { margin-top: 8px; }
.video-player { width: 100%; border-radius: 12px; background: #000; max-height: 450px; }
.video-actions { display: flex; gap: 8px; margin-top: 12px; }

@media (max-width: 1024px) { .main-layout { grid-template-columns: 1fr; } }
</style>
