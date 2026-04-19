<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🎨 AI 生图</h2>
      <p class="page-sub">通过文字描述生成精美图片，或对已有图片进行智能编辑</p>
    </div>

    <!-- Tab 切换 -->
    <el-tabs v-model="activeTab" class="gen-tabs" type="border-card">
      <!-- ==================== 文生图 Tab ==================== -->
      <el-tab-pane label="文生图" name="text2img">
        <div class="page-grid">
          <!-- 左侧参数区 -->
          <div class="param-card glass-card">
            <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

            <!-- 模型选择 -->
            <div class="field">
              <label class="field-label">选择模型</label>
              <el-select v-model="t2iForm.model" class="full-width" @change="onModelChange">
                <el-option label="万相 2.7 Pro（最高画质）" value="wan2.7-image-pro" />
                <el-option label="万相 2.7（快速生成）" value="wan2.7-image" />
              </el-select>
            </div>

            <!-- 提示词 -->
            <div class="field">
              <label class="field-label">画面描述（Prompt）</label>
              <el-input
                v-model="t2iForm.prompt"
                type="textarea"
                :rows="4"
                resize="none"
                placeholder="例如：一只穿着宇航服的猫咪在月球表面漫步，背景是地球，赛博朋克风格，超高清画质"
              />
              <p class="field-desc">描述你想要的画面内容，越详细效果越好</p>
            </div>

            <!-- 反向提示词（仅 Pro） -->
            <div class="field" v-if="isProModel">
              <label class="field-label">不想出现的内容（Negative Prompt） <span class="optional">可选</span></label>
              <el-input
                v-model="t2iForm.negative_prompt"
                type="textarea"
                :rows="2"
                resize="none"
                placeholder="例如：模糊，低质量，文字，水印，变形的手指"
              />
              <p class="field-desc">指定画面中不应出现的元素</p>
            </div>

            <!-- 画面比例 -->
            <div class="field">
              <label class="field-label">画面比例</label>
              <el-radio-group v-model="t2iForm.size" class="ratio-group">
                <el-radio-button
                  v-for="opt in sizeOptions"
                  :key="opt.value"
                  :value="opt.value"
                >{{ opt.label }}</el-radio-button>
              </el-radio-group>
              <p class="field-desc">选择生成图片的宽高比例</p>
            </div>

            <!-- 生成数量（仅 Pro） -->
            <div class="field" v-if="isProModel">
              <div class="param-label-row">
                <span class="field-label">生成数量</span>
                <span class="param-value">{{ t2iForm.n }} 张</span>
              </div>
              <el-slider
                v-model="t2iForm.n"
                :min="1"
                :max="4"
                :step="1"
                :show-tooltip="false"
                :marks="{ 1: '1', 2: '2', 3: '3', 4: '4' }"
              />
              <p class="field-desc">一次生成的图片数量，最多 4 张</p>
            </div>

            <!-- 智能优化提示词 -->
            <div class="field">
              <div class="param-label-row">
                <span class="field-label">智能优化提示词</span>
                <el-switch v-model="t2iForm.prompt_extend" active-text="开启" inactive-text="关闭" />
              </div>
              <p class="field-desc">开启后 AI 会自动增强你的描述，生成更高质量的图片</p>
            </div>

            <el-button
              @click="submitText2Image"
              type="primary"
              :loading="t2iLoading"
              class="generate-btn"
              size="large"
            >
              <el-icon><Picture /></el-icon> 开始生成
            </el-button>
          </div>

          <!-- 右侧结果区 -->
          <div class="result-card glass-card">
            <h3 class="card-title"><el-icon><PictureFilled /></el-icon> 生成结果</h3>

            <!-- 进度条 -->
            <div v-if="t2iLoading" class="progress-section">
              <el-progress :percentage="t2iProgress" :stroke-width="8" striped striped-flow />
              <p class="progress-text">{{ t2iProgressText }}</p>
            </div>

            <!-- 错误信息 -->
            <div v-if="t2iError" class="error-section">
              <el-alert :title="t2iError" type="error" show-icon :closable="false" />
            </div>

            <!-- 图片网格 -->
            <div v-if="t2iResults.length > 0" class="gallery-grid">
              <div
                v-for="(img, idx) in t2iResults"
                :key="idx"
                class="gallery-item"
                @click="previewUrl = img.url; showPreview = true"
              >
                <img :src="img.url" :alt="'生成结果 ' + (idx + 1)" />
                <div class="gallery-overlay">
                  <el-button link @click.stop="downloadFile(img.url, idx)">⬇ 下载</el-button>
                </div>
              </div>
            </div>

            <el-empty
              v-if="t2iResults.length === 0 && !t2iLoading && !t2iError"
              description="还没有生成图片，开始创作吧 ✨"
            />
          </div>
        </div>
      </el-tab-pane>

      <!-- ==================== 图像编辑 Tab ==================== -->
      <el-tab-pane label="图像编辑" name="imageEdit">
        <div class="page-grid">
          <!-- 左侧参数区 -->
          <div class="param-card glass-card">
            <h3 class="card-title"><el-icon><MagicStick /></el-icon> 编辑参数</h3>

            <!-- 模型选择 -->
            <div class="field">
              <label class="field-label">选择模型</label>
              <el-select v-model="i2iForm.model" class="full-width">
                <el-option label="万相 2.5 图像编辑" value="wan2.5-i2i-preview" />
              </el-select>
            </div>

            <!-- 功能选择 -->
            <div class="field">
              <label class="field-label">编辑功能</label>
              <el-select v-model="i2iForm.function" class="full-width">
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
                :action="uploadAction"
                list-type="picture-card"
                :limit="refImageMax"
                accept="image/jpeg,image/png,image/webp"
                :on-exceed="handleRefExceed"
                :on-remove="handleRefRemove"
                :before-upload="beforeRefUpload"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
              <p class="field-desc">上传需要编辑的参考图片，支持 JPG / PNG / WebP 格式</p>
            </div>

            <!-- 编辑指令 -->
            <div class="field" v-if="showPromptField">
              <label class="field-label">编辑指令</label>
              <el-input
                v-model="i2iForm.prompt"
                type="textarea"
                :rows="3"
                resize="none"
                placeholder="描述你想要的编辑效果，例如：将画面风格改为水彩画"
              />
              <p class="field-desc">用自然语言描述你希望图片产生的变化</p>
            </div>

            <!-- 蒙版上传（编辑描述功能可选） -->
            <div class="field" v-if="showMaskField">
              <label class="field-label">蒙版图片 <span class="optional">（可选）</span></label>
              <el-upload
                v-model:file-list="maskFileList"
                :action="uploadAction"
                list-type="picture-card"
                :limit="1"
                accept="image/jpeg,image/png,image/webp"
                :on-exceed="handleMaskExceed"
                :on-remove="handleMaskRemove"
                :before-upload="beforeMaskUpload"
              >
                <el-icon><Plus /></el-icon>
              </el-upload>
              <p class="field-desc">上传蒙版指定编辑区域，白色为编辑区域，黑色为保留区域</p>
            </div>

            <el-button
              @click="submitImageEdit"
              type="primary"
              :loading="i2iLoading"
              class="generate-btn"
              size="large"
            >
              <el-icon><MagicStick /></el-icon> 开始编辑
            </el-button>
          </div>

          <!-- 右侧结果区 -->
          <div class="result-card glass-card">
            <h3 class="card-title"><el-icon><PictureFilled /></el-icon> 编辑结果</h3>

            <!-- 进度条 -->
            <div v-if="i2iLoading" class="progress-section">
              <el-progress :percentage="i2iProgress" :stroke-width="8" striped striped-flow />
              <p class="progress-text">{{ i2iProgressText }}</p>
            </div>

            <!-- 错误信息 -->
            <div v-if="i2iError" class="error-section">
              <el-alert :title="i2iError" type="error" show-icon :closable="false" />
            </div>

            <!-- 图片网格 -->
            <div v-if="i2iResults.length > 0" class="gallery-grid">
              <div
                v-for="(img, idx) in i2iResults"
                :key="idx"
                class="gallery-item"
                @click="previewUrl = img.url; showPreview = true"
              >
                <img :src="img.url" :alt="'编辑结果 ' + (idx + 1)" />
                <div class="gallery-overlay">
                  <el-button link @click.stop="downloadFile(img.url, idx)">⬇ 下载</el-button>
                </div>
              </div>
            </div>

            <el-empty
              v-if="i2iResults.length === 0 && !i2iLoading && !i2iError"
              description="还没有编辑结果，上传参考图开始吧 🖼️"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 图片预览弹窗 -->
    <el-dialog v-model="showPreview" title="预览图片" width="80%" class="preview-dialog">
      <img :src="previewUrl" style="width: 100%; border-radius: 8px;" alt="预览" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import { EditPen, Picture, PictureFilled, MagicStick, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// ==================== 通用 ====================
const activeTab = ref('text2img')
const showPreview = ref(false)
const previewUrl = ref('')
const uploadAction = '/api/v1/upload' // 占位，实际用本地文件

const downloadFile = (url, idx) => {
  const a = document.createElement('a')
  a.href = url
  a.download = `generated_${Date.now()}_${idx}.png`
  a.target = '_blank'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

// ==================== 文生图 ====================
// 不同模型对应的比例选项
const sizeOptionsMap = {
  'wan2.7-image-pro': [
    { label: '1:1', value: '1024*1024' },
    { label: '16:9', value: '1620*1080' },
    { label: '9:16', value: '1080*1920' },
    { label: '4:3', value: '1152*864' },
    { label: '3:4', value: '864*1152' },
  ],
  'wan2.7-image': [
    { label: '1:1', value: '1024*1024' },
    { label: '16:9', value: '1620*1080' },
    { label: '9:16', value: '1080*1920' },
    { label: '4:3', value: '1152*864' },
  ],
}

const t2iForm = ref({
  model: 'wan2.7-image-pro',
  prompt: '',
  negative_prompt: '',
  size: '1024*1024',
  n: 1,
  prompt_extend: false,
})

const isProModel = computed(() => t2iForm.value.model === 'wan2.7-image-pro')

const sizeOptions = computed(() => {
  return sizeOptionsMap[t2iForm.value.model] || sizeOptionsMap['wan2.7-image-pro']
})

const t2iLoading = ref(false)
const t2iProgress = ref(0)
const t2iProgressText = ref('')
const t2iError = ref('')
const t2iResults = ref([])
let t2iPollTimer = null

/** 模型切换时重置比例选项 */
function onModelChange() {
  const options = sizeOptions.value
  // 如果当前选中的比例不在新选项中，重置为第一个
  if (!options.find(o => o.value === t2iForm.value.size)) {
    t2iForm.value.size = options[0].value
  }
}

/** 提交文生图任务 */
async function submitText2Image() {
  if (!t2iForm.value.prompt.trim()) {
    return ElMessage.warning('请输入画面描述')
  }

  t2iLoading.value = true
  t2iProgress.value = 5
  t2iProgressText.value = '正在提交任务...'
  t2iError.value = ''
  t2iResults.value = []

  try {
    const input = { prompt: t2iForm.value.prompt }
    if (isProModel.value && t2iForm.value.negative_prompt.trim()) {
      input.negative_prompt = t2iForm.value.negative_prompt
    }

    const parameters = { size: t2iForm.value.size }
    if (isProModel.value) {
      parameters.n = t2iForm.value.n
    }
    parameters.prompt_extend = t2iForm.value.prompt_extend

    const res = await fetch('/api/v1/services/aigc/text2image/image-synthesis', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: t2iForm.value.model,
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

    t2iProgress.value = 20
    t2iProgressText.value = '任务已提交，正在生成中...'

    // 开始轮询
    await pollTask(taskId, 't2i')
  } catch (e) {
    t2iError.value = `生成失败: ${e.message}`
    ElMessage.error(t2iError.value)
  } finally {
    t2iLoading.value = false
    // 延迟清除进度条
    setTimeout(() => {
      if (!t2iLoading.value) {
        t2iProgress.value = 0
        t2iProgressText.value = ''
      }
    }, 2000)
  }
}

/**
 * 通用任务轮询逻辑
 * @param {string} taskId
 * @param {'t2i'|'i2i'} mode
 */
function pollTask(taskId, mode) {
  return new Promise((resolve, reject) => {
    let pollCount = 0
    const maxPolls = 120 // 最长 120 * 3s = 6 分钟

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
          t2iPollTimer = setTimeout(doPoll, 3000)
          return
        }

        if (status === 'SUCCEEDED') {
          if (mode === 't2i') {
            t2iProgress.value = 100
            t2iProgressText.value = '生成完成！'
            t2iResults.value = (pollData.output.results || []).map(r => ({ url: r.url, taskId }))
            ElMessage.success('🎉 图片生成成功！')
          } else {
            i2iProgress.value = 100
            i2iProgressText.value = '编辑完成！'
            i2iResults.value = (pollData.output.results || []).map(r => ({ url: r.url, taskId }))
            ElMessage.success('🎉 图片编辑成功！')
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
          throw new Error('任务超时（等待超过 6 分钟）')
        }

        if (mode === 't2i') {
          t2iProgress.value = Math.min(20 + pollCount * 3, 95)
          t2iProgressText.value = `生成中... 已等待 ${pollCount * 3}秒`
        } else {
          i2iProgress.value = Math.min(20 + pollCount * 3, 95)
          i2iProgressText.value = `编辑中... 已等待 ${pollCount * 3}秒`
        }

        t2iPollTimer = setTimeout(doPoll, 3000)
      } catch (e) {
        reject(e)
      }
    }

    doPoll()
  })
}

// ==================== 图像编辑 ====================
const i2iForm = ref({
  model: 'wan2.5-i2i-preview',
  function: 'style_transfer',
  prompt: '',
})

const refFileList = ref([])
const maskFileList = ref([])
const refImageFiles = ref([])   // 本地 File 对象
const maskImageFile = ref(null) // 本地 File 对象

const refImageMin = 1
const refImageMax = 4

const functionDescMap = {
  style_transfer: '将参考图的画面风格迁移到目标风格',
  description_edit: '根据文字指令对图片进行局部或全局编辑',
  expand: '智能扩展画面边界，补全缺失区域',
  super_resolution: '提升图片分辨率，增强细节清晰度',
}

const functionDesc = computed(() => functionDescMap[i2iForm.value.function] || '')

const showPromptField = computed(() => {
  const fn = i2iForm.value.function
  return fn === 'style_transfer' || fn === 'description_edit' || fn === 'expand'
})

const showMaskField = computed(() => {
  return i2iForm.value.function === 'description_edit'
})

const i2iLoading = ref(false)
const i2iProgress = ref(0)
const i2iProgressText = ref('')
const i2iError = ref('')
const i2iResults = ref([])
let i2iPollTimer = null

// ---- 参考图上传处理 ----
function beforeRefUpload(file) {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG / PNG / WebP 格式')
    return false
  }
  refImageFiles.value.push(file)
  return false // 阻止默认上传，我们自行处理
}

function handleRefExceed() {
  ElMessage.warning(`最多上传 ${refImageMax} 张参考图片`)
}

function handleRefRemove(file) {
  const idx = refImageFiles.value.indexOf(file.raw || file)
  if (idx > -1) refImageFiles.value.splice(idx, 1)
}

// ---- 蒙版上传处理 ----
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

/** 提交图像编辑任务 */
async function submitImageEdit() {
  if (refImageFiles.value.length < refImageMin) {
    return ElMessage.warning(`请至少上传 ${refImageMin} 张参考图片`)
  }
  if (showPromptField.value && !i2iForm.value.prompt.trim()) {
    return ElMessage.warning('请输入编辑指令')
  }

  i2iLoading.value = true
  i2iProgress.value = 5
  i2iProgressText.value = '正在提交任务...'
  i2iError.value = ''
  i2iResults.value = []

  try {
    // 将图片转为 base64
    const refImagesBase64 = []
    for (const file of refImageFiles.value) {
      const b64 = await fileToBase64(file)
      refImagesBase64.push(b64)
    }

    const input = {
      prompt: i2iForm.value.prompt || '',
      ref_images: refImagesBase64,
    }

    if (maskImageFile.value) {
      input.mask = await fileToBase64(maskImageFile.value)
    }

    const parameters = {
      function: i2iForm.value.function,
    }

    const res = await fetch('/api/v1/services/aigc/image2image/image-synthesis', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: i2iForm.value.model,
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

    i2iProgress.value = 20
    i2iProgressText.value = '任务已提交，正在编辑中...'

    await pollTask(taskId, 'i2i')
  } catch (e) {
    i2iError.value = `编辑失败: ${e.message}`
    ElMessage.error(i2iError.value)
  } finally {
    i2iLoading.value = false
    setTimeout(() => {
      if (!i2iLoading.value) {
        i2iProgress.value = 0
        i2iProgressText.value = ''
      }
    }, 2000)
  }
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

/** 切换 Tab 时清理定时器 */
watch(activeTab, () => {
  if (t2iPollTimer) { clearTimeout(t2iPollTimer); t2iPollTimer = null }
  if (i2iPollTimer) { clearTimeout(i2iPollTimer); i2iPollTimer = null }
})

/** 组件卸载时清理 */
onUnmounted(() => {
  if (t2iPollTimer) clearTimeout(t2iPollTimer)
  if (i2iPollTimer) clearTimeout(i2iPollTimer)
})
</script>

<style scoped>
/* ==================== 布局 ==================== */
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

/* Tabs 样式 */
.gen-tabs { margin-top: 8px; }
.gen-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}
.gen-tabs :deep(.el-tabs--border-card) {
  background: transparent;
  border: none;
  box-shadow: none;
}
.gen-tabs :deep(.el-tabs__content) {
  padding: 0;
}

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
.optional { font-weight: 400; color: var(--text-secondary); font-size: 12px; }
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

.full-width { width: 100%; }
.ratio-group { display: flex; flex-wrap: wrap; gap: 8px; }
.ratio-group :deep(.el-radio-button) { margin-bottom: 4px; }

/* 上传区域 */
.field :deep(.el-upload--picture-card) {
  width: 80px;
  height: 80px;
}
.field :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 80px;
  height: 80px;
}

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

/* ==================== 进度 & 错误 ==================== */
.progress-section { padding: 16px 0; }
.progress-text {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 8px;
  text-align: center;
}
.error-section { margin-bottom: 16px; }

/* ==================== 图片网格 ==================== */
.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}
.gallery-item {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  aspect-ratio: 1;
  cursor: pointer;
}
.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}
.gallery-item:hover img { transform: scale(1.05); }
.gallery-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.6));
  padding: 8px;
  opacity: 0;
  transition: opacity 0.3s;
}
.gallery-item:hover .gallery-overlay { opacity: 1; }

/* 预览弹窗 */
.preview-dialog :deep(.el-dialog) {
  background: var(--card-bg);
}

/* ==================== 响应式 ==================== */
@media (max-width: 768px) {
  .page-grid { grid-template-columns: 1fr; }
}
</style>
