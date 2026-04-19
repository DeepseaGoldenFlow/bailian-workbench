<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎬 AI 生视频</h2>
      <p class="page-sub">通过文字描述或参考图片生成动态视频</p>
    </div>

    <div class="page-grid">
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="万相 2.7 (最新，多镜头叙事)" value="wan2.7-t2v" />
            <el-option label="万相 2.2 Plus (稳定)" value="wan2.2-t2v-plus" />
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">视频描述 (Prompt)</label>
          <el-input v-model="form.prompt" type="textarea" :rows="4" resize="none"
            placeholder="例如: 海浪拍打在礁石上, 夕阳余晖映照海面, 一只海鸥从画面中飞过, 电影质感" />
        </div>

        <div class="field">
          <label class="field-label">参考图片 <span class="optional">可选 - 上传图片控制画面风格</span></label>
          <el-upload action="" :auto-upload="false" :show-file-list="false" accept="image/*"
            @change="handleImageUpload" drag>
            <div v-if="!form.image_url" class="upload-area">
              <el-icon class="upload-icon"><UploadFilled /></el-icon>
              <p>拖拽图片到此处或点击上传</p>
            </div>
            <img v-else :src="form.image_url" class="upload-preview" />
          </el-upload>
          <el-button v-if="form.image_url" link type="danger" size="small" @click="form.image_url = ''; form.image_base64 = ''">移除图片</el-button>
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
          <div class="param-label-row">
            <span class="field-label">智能优化提示词</span>
            <el-switch v-model="form.auto_prompt" active-text="开启" inactive-text="关闭" />
          </div>
          <p class="field-desc">自动增强你的描述以获得更好的视频效果</p>
        </div>

        <el-button @click="generate" type="primary" :loading="loading" class="generate-btn" size="large">
          <el-icon><VideoCamera /></el-icon> 开始生成视频
        </el-button>
      </div>

      <div class="result-card glass-card">
        <h3 class="card-title"><el-icon><VideoCameraFilled /></el-icon> 生成结果</h3>

        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">{{ progressText }}</p>
        </div>

        <div v-if="results.length > 0" class="video-gallery">
          <div v-for="vid in results" :key="vid.url" class="video-item">
            <video :src="vid.url" controls preload="metadata" class="video-player" />
            <div class="video-actions">
              <el-button link @click="downloadFile(vid.url)">⬇ 下载视频</el-button>
            </div>
          </div>
        </div>

        <el-empty v-if="results.length === 0 && !loading" description="还没有生成视频，开始创作吧 🎥" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditPen, VideoCamera, VideoCameraFilled, UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({
  model: 'wan2.7-t2v',
  prompt: '',
  image_url: '',
  image_base64: '',
  duration: '5',
  resolution: '720P',
  auto_prompt: false,
})

const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const results = ref([])

const handleImageUpload = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => { form.value.image_base64 = e.target.result; form.value.image_url = e.target.result }
  reader.readAsDataURL(file.raw)
}

const generate = async () => {
  if (!form.value.prompt) return ElMessage.warning('请输入视频描述')
  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'

  try {
    const input = { prompt: form.value.prompt }
    const parameters = { duration: parseInt(form.value.duration), resolution: form.value.resolution }
    if (form.value.auto_prompt) parameters.auto_prompt = true

    const res = await fetch('/api/video/t2v', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ model: form.value.model, input, parameters })
    })

    if (!res.ok) { const e = await res.json(); throw new Error(e.message || '提交失败') }
    const data = await res.json()
    const taskId = data.output.task_id

    progress.value = 15
    progressText.value = '视频生成中，请耐心等待（通常需要 2-5 分钟）...'

    let status = 'PENDING'
    let pollCount = 0
    while (status !== 'SUCCEEDED' && status !== 'FAILED' && pollCount < 200) {
      await new Promise(r => setTimeout(r, 5000))
      pollCount++
      const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
      const pollData = await pollRes.json()
      status = pollData.output.task_status

      progress.value = Math.min(15 + pollCount * 2, 95)
      progressText.value = `生成中... 已等待 ${pollCount * 5}秒`

      if (status === 'SUCCEEDED') {
        progress.value = 100
        progressText.value = '视频生成完成！'
        let videoUrl = pollData.output.video_url
        if (videoUrl.startsWith('/data/bailian/storage')) {
          videoUrl = videoUrl.replace('/data/bailian/storage', '/api/files')
        }
        results.value = [{ url: videoUrl, taskId }]
        ElMessage.success('🎉 视频生成成功！')
      } else if (status === 'FAILED') {
        throw new Error(pollData.output.message || '生成失败')
      }
    }
  } catch (e) {
    ElMessage.error('生成失败: ' + e.message)
  } finally {
    loading.value = false
    setTimeout(() => { progress.value = 0; progressText.value = '' }, 2000)
  }
}

const downloadFile = (url) => { window.open(url, '_blank') }
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
.optional { font-weight: 400; color: var(--text-secondary); font-size: 12px; }
.param-label-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.field-desc { font-size: 11px; color: var(--text-secondary); margin-top: 4px; }
.full-width { width: 100%; }

.upload-area { border: 2px dashed var(--card-border); border-radius: 10px; padding: 24px; text-align: center; cursor: pointer; transition: all 0.2s; }
.upload-area:hover { border-color: var(--gradient-start); }
.upload-icon { font-size: 32px; color: var(--text-secondary); margin-bottom: 8px; }
.upload-area p { font-size: 13px; color: var(--text-secondary); }
.upload-preview { max-width: 100%; max-height: 150px; border-radius: 8px; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

.video-gallery { display: grid; gap: 16px; }
.video-item { border-radius: 12px; overflow: hidden; }
.video-player { width: 100%; border-radius: 12px; background: #000; max-height: 400px; }
.video-actions { display: flex; gap: 8px; padding: 8px 0; }

@media (max-width: 768px) { .page-grid { grid-template-columns: 1fr; } }
</style>
