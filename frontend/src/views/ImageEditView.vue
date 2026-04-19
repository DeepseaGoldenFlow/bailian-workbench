<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🖼️ 图像编辑</h2>
      <p class="page-sub">对已有图像进行智能编辑和修改</p>
    </div>

    <div class="page-grid">
      <!-- ==================== 左侧参数区 ==================== -->
      <div class="param-card glass-card">
        <h3 class="card-title">
          <el-icon><MagicStick /></el-icon> 编辑参数
        </h3>

        <!-- 功能选择 -->
        <div class="field">
          <label class="field-label">编辑功能</label>
          <el-select v-model="form.function" class="full-width" @change="onFunctionChange">
            <el-option label="风格迁移" value="style_transfer" />
            <el-option label="编辑描述" value="description_edit" />
            <el-option label="智能扩图" value="expand" />
            <el-option label="超分辨率" value="super_resolution" />
          </el-select>
          <p class="field-desc">{{ functionDesc }}</p>
        </div>

        <!-- 参考图上传 -->
        <div class="field">
          <label class="field-label">参考图片 <span class="optional">（{{ refImageMin }}~{{ refImageMax }} 张）</span></label>
          <el-upload
            v-model:file-list="refFileList"
            list-type="picture-card"
            :limit="refImageMax"
            accept="image/jpeg,image/png,image/webp"
            :before-upload="beforeRefUpload"
            :on-exceed="handleRefExceed"
            :on-remove="handleRefRemove"
            :auto-upload="false"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <p class="field-desc">上传需要编辑的参考图片，支持 JPG / PNG / WebP 格式</p>
        </div>

        <!-- 编辑指令（风格迁移 / 编辑描述 / 智能扩图需要） -->
        <div class="field" v-if="showPromptField">
          <label class="field-label">编辑指令</label>
          <el-input
            v-model="form.prompt"
            type="textarea"
            :rows="3"
            resize="none"
            :placeholder="promptPlaceholder"
          />
          <p class="field-desc">用自然语言描述你希望图片产生的变化</p>
        </div>

        <!-- 蒙版上传（仅编辑描述功能需要） -->
        <div class="field" v-if="showMaskField">
          <label class="field-label">蒙版图片 <span class="optional">（可选）</span></label>
          <el-upload
            v-model:file-list="maskFileList"
            list-type="picture-card"
            :limit="1"
            accept="image/jpeg,image/png,image/webp"
            :before-upload="beforeMaskUpload"
            :on-exceed="handleMaskExceed"
            :on-remove="handleMaskRemove"
            :auto-upload="false"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <p class="field-desc">上传蒙版指定编辑区域，白色为编辑区域，黑色为保留区域</p>
        </div>

        <!-- 超分辨率参数 -->
        <div class="field" v-if="form.function === 'super_resolution'">
          <div class="param-label-row">
            <span class="field-label">放大倍数</span>
            <span class="param-value">{{ form.parameters.upscale }}x</span>
          </div>
          <el-slider
            v-model="form.parameters.upscale"
            :min="1"
            :max="4"
            :step="1"
            :show-tooltip="false"
            :marks="{ 1: '1x', 2: '2x', 3: '3x', 4: '4x' }"
          />
          <p class="field-desc">选择图片分辨率提升倍数</p>
        </div>

        <!-- 智能扩图参数 -->
        <div class="field" v-if="form.function === 'expand'">
          <div class="param-label-row">
            <span class="field-label">扩展比例</span>
            <span class="param-value">{{ form.parameters.expand_ratio }}%</span>
          </div>
          <el-slider
            v-model="form.parameters.expand_ratio"
            :min="10"
            :max="100"
            :step="10"
            :show-tooltip="false"
            :marks="{ 10: '10%', 50: '50%', 100: '100%' }"
          />
          <p class="field-desc">画面各方向扩展的比例，数值越大扩展范围越多</p>
        </div>

        <el-button
          @click="submitEdit"
          type="primary"
          :loading="loading"
          class="generate-btn"
          size="large"
        >
          <el-icon><MagicStick /></el-icon> 开始编辑
        </el-button>
      </div>

      <!-- ==================== 右侧结果区 ==================== -->
      <div class="result-card glass-card">
        <h3 class="card-title">
          <el-icon><PictureFilled /></el-icon> 编辑结果
        </h3>

        <!-- 参考图预览 -->
        <div v-if="refPreviewUrls.length > 0" class="ref-preview-section">
          <span class="ref-preview-label">参考图</span>
          <div class="ref-preview-grid">
            <div
              v-for="(url, idx) in refPreviewUrls"
              :key="'ref-' + idx"
              class="ref-preview-item"
              @click="previewUrl = url; showPreview = true"
            >
              <img :src="url" :alt="'参考图 ' + (idx + 1)" />
            </div>
          </div>
        </div>

        <!-- 进度条 -->
        <div v-if="loading" class="progress-section">
          <el-progress :percentage="progress" :stroke-width="8" striped striped-flow />
          <p class="progress-text">{{ progressText }}</p>
        </div>

        <!-- 错误信息 -->
        <div v-if="error" class="error-section">
          <el-alert :title="error" type="error" show-icon :closable="false" />
        </div>

        <!-- 结果图片 -->
        <div v-if="resultImage" class="result-image-section">
          <div class="result-image-wrapper" @click="previewUrl = resultImage; showPreview = true">
            <img :src="resultImage" alt="编辑结果" class="result-image" />
            <div class="result-overlay">
              <el-button link @click.stop="downloadResult">⬇ 下载原图</el-button>
            </div>
          </div>
        </div>

        <el-empty
          v-if="!resultImage && !loading && !error"
          description="还没有编辑结果，上传参考图开始吧 🖼️"
        />
      </div>
    </div>

    <!-- 图片预览弹窗 -->
    <el-dialog v-model="showPreview" title="预览图片" width="80%" class="preview-dialog">
      <img :src="previewUrl" style="width: 100%; border-radius: 8px;" alt="预览" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import { MagicStick, PictureFilled, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// ==================== 表单数据 ====================
const form = ref({
  function: 'style_transfer',
  prompt: '',
  parameters: {
    upscale: 2,
    expand_ratio: 50,
  },
})

// ==================== 功能描述映射 ====================
const functionDescMap = {
  style_transfer: '将参考图的画面风格迁移到目标风格，例如改为水彩画或油画风格',
  description_edit: '根据文字指令对图片进行局部或全局编辑，支持蒙版指定区域',
  expand: '智能扩展画面边界，自动补全缺失区域，保持画面自然连贯',
  super_resolution: '提升图片分辨率，增强细节清晰度，适合老照片修复和画质提升',
}

const functionDesc = computed(() => functionDescMap[form.value.function] || '')

const promptPlaceholderMap = {
  style_transfer: '描述你想要的目标风格，例如：将画面改为中国水墨画风格',
  description_edit: '描述你想要的编辑效果，例如：把画面中的天空换成夕阳',
  expand: '描述需要扩展的画面内容，例如：向四周扩展，添加更多的背景环境',
  super_resolution: '可描述增强方向（可选），例如：增强人物面部细节',
}

const promptPlaceholder = computed(() => promptPlaceholderMap[form.value.function] || '')

// ==================== 条件显示 ====================
const showPromptField = computed(() => {
  const fn = form.value.function
  return fn === 'style_transfer' || fn === 'description_edit' || fn === 'expand'
})

const showMaskField = computed(() => {
  return form.value.function === 'description_edit'
})

// ==================== 参考图上传 ====================
const refFileList = ref([])
const refImageFiles = ref([])
const refImageMin = 1
const refImageMax = 4

// 参考图预览 URL
const refPreviewUrls = ref([])

function beforeRefUpload(file) {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG / PNG / WebP 格式')
    return false
  }
  refImageFiles.value.push(file)
  // 生成预览 URL
  const url = URL.createObjectURL(file)
  refPreviewUrls.value.push(url)
  return false // 阻止默认上传，我们自行处理
}

function handleRefExceed() {
  ElMessage.warning(`最多上传 ${refImageMax} 张参考图片`)
}

function handleRefRemove(file) {
  const idx = refImageFiles.value.indexOf(file.raw || file)
  if (idx > -1) {
    // 释放预览 URL
    if (refPreviewUrls.value[idx]) {
      URL.revokeObjectURL(refPreviewUrls.value[idx])
    }
    refImageFiles.value.splice(idx, 1)
    refPreviewUrls.value.splice(idx, 1)
  }
}

// ==================== 蒙版上传 ====================
const maskFileList = ref([])
const maskImageFile = ref(null)

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
  ElMessage.warning('仅支持上传 1 张蒙版图片')
}

function handleMaskRemove() {
  maskImageFile.value = null
}

// ==================== 功能切换时重置蒙版 ====================
function onFunctionChange() {
  // 切换到非编辑描述功能时，清空蒙版
  if (form.value.function !== 'description_edit') {
    maskFileList.value = []
    maskImageFile.value = null
  }
  // 切换功能时清空之前的结果
  resultImage.value = ''
  error.value = ''
}

// ==================== 提交任务 ====================
const loading = ref(false)
const progress = ref(0)
const progressText = ref('')
const error = ref('')
const resultImage = ref('')
let pollTimer = null

async function submitEdit() {
  if (refImageFiles.value.length < refImageMin) {
    return ElMessage.warning(`请至少上传 ${refImageMin} 张参考图片`)
  }
  if (showPromptField.value && !form.value.prompt.trim()) {
    return ElMessage.warning('请输入编辑指令')
  }

  loading.value = true
  progress.value = 5
  progressText.value = '正在提交任务...'
  error.value = ''
  resultImage.value = ''

  try {
    // 将参考图转为 base64
    const refImagesBase64 = []
    for (const file of refImageFiles.value) {
      const b64 = await fileToBase64(file)
      refImagesBase64.push({ image: b64 })
    }

    const input = {
      function: form.value.function,
      ref_images: refImagesBase64,
    }

    if (showPromptField.value && form.value.prompt.trim()) {
      input.prompt = form.value.prompt
    }

    if (maskImageFile.value) {
      const maskB64 = await fileToBase64(maskImageFile.value)
      input.mask_image = { image: maskB64 }
    }

    // 根据功能构建 parameters
    const parameters = {}
    if (form.value.function === 'super_resolution') {
      parameters.upscale = form.value.parameters.upscale
    }
    if (form.value.function === 'expand') {
      parameters.expand_ratio = form.value.parameters.expand_ratio
    }

    const res = await fetch('/api/image/edit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: 'wan2.5-i2i-preview',
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
    progressText.value = '任务已提交，正在编辑中...'

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

// ==================== 任务轮询 ====================
function pollTask(taskId) {
  return new Promise((resolve, reject) => {
    let pollCount = 0
    const maxPolls = 200 // 最长 200 * 3s = 10 分钟

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
          pollTimer = setTimeout(doPoll, 3000)
          return
        }

        if (status === 'SUCCEEDED') {
          progress.value = 100
          progressText.value = '编辑完成！'
          // 提取结果图 URL
          const results = pollData.output?.results || []
          if (results.length > 0 && results[0].url) {
            resultImage.value = results[0].url
          } else if (pollData.output?.result_url) {
            resultImage.value = pollData.output.result_url
          } else {
            // 尝试从 results 中找 base64 或其他字段
            const r = results[0] || {}
            resultImage.value = r.url || r.image || ''
          }
          ElMessage.success('🎉 图片编辑成功！')
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
          throw new Error('任务超时（等待超过 10 分钟）')
        }

        progress.value = Math.min(20 + pollCount * 3, 95)
        progressText.value = `编辑中... 已等待 ${pollCount * 3}秒`

        pollTimer = setTimeout(doPoll, 3000)
      } catch (e) {
        reject(e)
      }
    }

    doPoll()
  })
}

// ==================== 工具函数 ====================
/** 文件转 base64 data URL */
function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

/** 下载结果图 */
function downloadResult() {
  if (!resultImage.value) return
  const a = document.createElement('a')
  a.href = resultImage.value
  a.download = `edited_${Date.now()}.png`
  a.target = '_blank'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

// ==================== 预览弹窗 ====================
const showPreview = ref(false)
const previewUrl = ref('')

// ==================== 清理 ====================
onUnmounted(() => {
  if (pollTimer) clearTimeout(pollTimer)
  // 释放预览 URL
  refPreviewUrls.value.forEach(url => URL.revokeObjectURL(url))
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
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.page-sub {
  font-size: 14px;
  color: var(--text-secondary);
}

.page-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

/* ==================== 玻璃拟态卡片 ==================== */
.glass-card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 16px;
  padding: 20px;
  backdrop-filter: blur(10px);
  transition: box-shadow 0.3s;
}

.glass-card:hover {
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
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

.optional {
  font-weight: 400;
  color: var(--text-secondary);
  font-size: 12px;
}

.param-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.param-value {
  font-size: 13px;
  font-weight: 700;
  color: var(--gradient-start);
  font-variant-numeric: tabular-nums;
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

/* ==================== 参考图预览 ==================== */
.ref-preview-section {
  margin-bottom: 16px;
}

.ref-preview-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.ref-preview-grid {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.ref-preview-item {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid var(--card-border);
}

.ref-preview-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

/* ==================== 结果图片 ==================== */
.result-image-section {
  margin-top: 8px;
}

.result-image-wrapper {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
}

.result-image {
  width: 100%;
  border-radius: 12px;
  transition: transform 0.3s;
}

.result-image-wrapper:hover .result-image {
  transform: scale(1.02);
}

.result-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.6));
  padding: 12px;
  opacity: 0;
  transition: opacity 0.3s;
  display: flex;
  justify-content: center;
}

.result-image-wrapper:hover .result-overlay {
  opacity: 1;
}

/* ==================== 预览弹窗 ==================== */
.preview-dialog :deep(.el-dialog) {
  background: var(--card-bg);
}

/* ==================== 响应式 ==================== */
@media (max-width: 768px) {
  .page-grid {
    grid-template-columns: 1fr;
  }
}
</style>
