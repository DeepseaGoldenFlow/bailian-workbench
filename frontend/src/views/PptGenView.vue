<template>
  <div class="page-container">
    <div class="page-header">
      <h2>📊 PPT 智能生成</h2>
      <p class="page-sub">输入主题，AI 自动生成大纲并排版成 PPTX 文件</p>
    </div>

    <div class="main-layout">
      <div class="input-panel glass-card">
        <div class="field">
          <label class="field-label">PPT 主题 <span class="required">*</span></label>
          <el-input v-model="form.topic" type="textarea" :rows="4" resize="none"
            placeholder="例如: 2026年AI行业发展趋势分析，或介绍我们的新产品百炼工作台..." />
        </div>
        
        <div class="params-grid">
          <div class="field">
            <label class="field-label">幻灯片页数</label>
            <el-input-number v-model="form.slideCount" :min="3" :max="20" controls-position="right" class="full-width" />
          </div>
          <div class="field">
            <label class="field-label">设计风格</label>
            <el-select v-model="form.style" class="full-width">
              <el-option label="商务简约 (默认)" value="商务简约" />
              <el-option label="科技感深色" value="科技感" />
              <el-option label="教育学术风" value="教育学术" />
              <el-option label="创意活泼" value="创意活泼" />
            </el-select>
          </div>
        </div>

        <el-button @click="generate" type="primary" :loading="loading" class="submit-btn" size="large">
          <el-icon><MagicStick /></el-icon> 开始生成 PPT
        </el-button>
      </div>

      <div class="result-panel glass-card">
        <h3 class="panel-title"><el-icon><Document /></el-icon> 生成结果</h3>
        
        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">{{ progressText }}</p>
        </div>

        <div v-if="downloadUrl" class="result-box">
          <el-icon class="result-icon"><Document /></el-icon>
          <p class="result-title">PPT 已生成完毕</p>
          <el-button @click="downloadPpt" type="primary" class="download-btn">
            <el-icon><Download /></el-icon> 下载 .pptx 文件
          </el-button>
        </div>

        <el-empty v-if="!downloadUrl && !loading" description="输入主题开始创作 📝" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { MagicStick, Document, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({
  topic: '',
  slideCount: 8,
  style: '商务简约'
})
const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const downloadUrl = ref('')

const generate = async () => {
  if (!form.value.topic.trim()) return ElMessage.warning('请输入主题')
  if (form.value.topic.length > 500) return ElMessage.warning('主题不能超过 500 字')
  
  loading.value = true
  progress.value = 10
  progressText.value = '正在构思大纲... (AI 思考中)'
  downloadUrl.value = ''

  try {
    // 注意：这里调用的是我们后端写的 /api/ppt/generate 接口
    const res = await fetch('/api/ppt/generate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        topic: form.value.topic,
        slide_count: form.value.slideCount,
        style: form.value.style
      })
    })

    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || '生成失败')
    }

    const data = await res.json()
    progress.value = 100
    progressText.value = '完成！'
    downloadUrl.value = data.download_url
    
    ElMessage.success('🎉 PPT 生成成功！')
  } catch (e) {
    ElMessage.error('生成失败: ' + e.message)
  } finally {
    loading.value = false
    setTimeout(() => { progress.value = 0 }, 2000)
  }
}

const downloadPpt = () => {
  if (downloadUrl.value) {
    window.open(downloadUrl.value, '_blank')
  }
}
</script>

<style scoped>
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.main-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; align-items: start; }

.glass-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; padding: 20px; backdrop-filter: blur(10px); }
.panel-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }

.field { margin-bottom: 16px; }
.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.required { color: #ef4444; margin-left: 2px; }
.full-width { width: 100%; }
.params-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }

.submit-btn { width: 100%; margin-top: 12px; border-radius: 10px; background: var(--btn-gradient); border: none; }

.progress-section { padding: 16px 0; }
.progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; }

.result-box { text-align: center; padding: 20px 0; }
.result-icon { font-size: 48px; color: var(--gradient-start); margin-bottom: 12px; }
.result-title { font-size: 16px; font-weight: 600; color: var(--text-primary); margin-bottom: 16px; }
.download-btn { width: 80%; font-size: 16px; padding: 12px; border-radius: 12px; }

@media (max-width: 1024px) { .main-layout { grid-template-columns: 1fr; } }
</style>
