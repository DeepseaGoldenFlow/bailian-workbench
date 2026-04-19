<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎨 AI 生图</h2>
      <p class="page-sub">通过文字描述生成精美图片</p>
    </div>

    <div class="page-grid">
      <!-- Left: Parameters -->
      <div class="param-card glass-card">
        <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

        <!-- 模型选择 -->
        <div class="field">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="万相 2.7 Pro (最高画质)" value="wan2.7-image-pro" />
            <el-option label="万相 2.7 (快速)" value="wan2.7-image" />
          </el-select>
        </div>

        <!-- 画面描述 -->
        <div class="field">
          <label class="field-label">画面描述 (Prompt)</label>
          <el-input v-model="form.prompt" type="textarea" :rows="4" resize="none"
            placeholder="例如: 一只穿着宇航服的猫咪在月球表面漫步, 背景是地球, 赛博朋克风格, 超高清画质" />
        </div>

        <!-- 不想出现的内容 -->
        <div class="field">
          <label class="field-label">不想出现的内容 (Negative Prompt) <span class="optional">可选</span></label>
          <el-input v-model="form.negative_prompt" type="textarea" :rows="2" resize="none"
            placeholder="例如: 模糊, 低质量, 文字, 水印, 变形的手指" />
        </div>

        <!-- 画面比例 -->
        <div class="field">
          <label class="field-label">画面比例</label>
          <el-radio-group v-model="form.size" class="ratio-group">
            <el-radio-button value="1024*1024">1:1</el-radio-button>
            <el-radio-button value="1620*1080">16:9</el-radio-button>
            <el-radio-button value="1080*1920">9:16</el-radio-button>
            <el-radio-button value="1152*864">4:3</el-radio-button>
            <el-radio-button value="864*1152">3:4</el-radio-button>
          </el-radio-group>
        </div>

        <!-- 生成数量 -->
        <div class="field">
          <div class="param-label-row">
            <span class="field-label">生成数量</span>
            <span class="param-value">{{ form.n }} 张</span>
          </div>
          <el-slider v-model="form.n" :min="1" :max="4" :step="1" :show-tooltip="false" :marks="{1:'1',2:'2',3:'3',4:'4'}" />
        </div>

        <!-- 智能优化提示词 -->
        <div class="field">
          <div class="param-label-row">
            <span class="field-label">智能优化提示词</span>
            <el-switch v-model="form.auto_prompt" active-text="开启" inactive-text="关闭" />
          </div>
          <p class="field-desc">开启后会自动增强你的描述, 生成更高质量的图片</p>
        </div>

        <el-button @click="generate" type="primary" :loading="loading" class="generate-btn" size="large">
          <el-icon><Picture /></el-icon> 开始生成
        </el-button>
      </div>

      <!-- Right: Results -->
      <div class="result-card glass-card">
        <h3 class="card-title"><el-icon><PictureFilled /></el-icon> 生成结果</h3>

        <!-- Progress -->
        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">{{ progressText }}</p>
        </div>

        <!-- Gallery -->
        <div v-if="results.length > 0" class="gallery-grid">
          <div v-for="(img, idx) in results" :key="idx" class="gallery-item">
            <img :src="img.url" @click="previewUrl = img.url; showPreview = true" />
            <div class="gallery-overlay">
              <el-button link @click="downloadFile(img.url)">⬇ 下载</el-button>
            </div>
          </div>
        </div>

        <el-empty v-if="results.length === 0 && !loading" description="还没有生成图片，开始创作吧 ✨" />
      </div>
    </div>

    <el-dialog v-model="showPreview" title="预览图片" width="80%">
      <img :src="previewUrl" style="width: 100%; border-radius: 8px;" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditPen, Picture, PictureFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({
  model: 'wan2.7-image-pro',
  prompt: '',
  negative_prompt: '',
  size: '1024*1024',
  n: 1,
  auto_prompt: false,
})

const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const results = ref([])
const showPreview = ref(false)
const previewUrl = ref('')

const generate = async () => {
  if (!form.value.prompt) return ElMessage.warning('请输入画面描述')
  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'

  try {
    const input = { prompt: form.value.prompt }
    if (form.value.negative_prompt) input.negative_prompt = form.value.negative_prompt

    const params = { size: form.value.size, n: form.value.n }
    if (form.value.auto_prompt) params.auto_prompt = true

    const res = await fetch('/api/v1/services/aigc/text2image/image-synthesis', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ model: form.value.model, input, parameters: params })
    })

    if (!res.ok) { const e = await res.json(); throw new Error(e.message || '提交失败') }
    const data = await res.json()
    const taskId = data.output.task_id

    progress.value = 20
    progressText.value = '任务已提交，正在生成中...'

    let status = 'PENDING'
    let pollCount = 0
    while (status !== 'SUCCEEDED' && status !== 'FAILED' && pollCount < 120) {
      await new Promise(r => setTimeout(r, 3000))
      pollCount++
      const pollRes = await fetch(`/api/v1/tasks/${taskId}`)
      const pollData = await pollRes.json()
      status = pollData.output.task_status

      progress.value = Math.min(20 + pollCount * 3, 95)
      progressText.value = `生成中... 已等待 ${pollCount * 3}秒`

      if (status === 'SUCCEEDED') {
        progress.value = 100
        progressText.value = '生成完成！'
        results.value = pollData.output.results.map(r => ({ url: r.url, taskId }))
        ElMessage.success('🎉 图片生成成功！')
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
.param-value { font-size: 13px; font-weight: 700; color: var(--gradient-start); font-variant-numeric: tabular-nums; }
.field-desc { font-size: 11px; color: var(--text-secondary); margin-top: 4px; }

.full-width { width: 100%; }
.ratio-group { display: flex; flex-wrap: wrap; gap: 8px; }
.ratio-group :deep(.el-radio-button) { margin-bottom: 4px; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

.gallery-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; }
.gallery-item { position: relative; border-radius: 12px; overflow: hidden; aspect-ratio: 1; cursor: pointer; }
.gallery-item img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s; }
.gallery-item:hover img { transform: scale(1.05); }
.gallery-overlay { position: absolute; bottom: 0; left: 0; right: 0; background: linear-gradient(transparent, rgba(0,0,0,0.6)); padding: 8px; opacity: 0; transition: opacity 0.3s; }
.gallery-item:hover .gallery-overlay { opacity: 1; }

@media (max-width: 768px) { .page-grid { grid-template-columns: 1fr; } }
</style>
