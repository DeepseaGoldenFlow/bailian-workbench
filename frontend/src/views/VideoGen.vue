<template>
  <div class="page" @paste="onGlobalPaste">
    <h2>视频生成</h2>
    <el-form label-width="100px" style="max-width:700px">
      <el-form-item label="模型">
        <el-select v-model="model" style="width:100%">
          <el-option label="Happyhorse-1.0-T2V (文生视频)" value="happyhorse-1.0-t2v" />
          <el-option label="Happyhorse-1.0-I2V (图生视频)" value="happyhorse-1.0-i2v" />
          <el-option label="Happyhorse-1.0-R2V (参考生视频)" value="happyhorse-1.0-r2v" />
          <el-option label="Happyhorse-1.0-Video-Edit (视频编辑)" value="happyhorse-1.0-video-edit" />
        </el-select>
      </el-form-item>
      <el-form-item label="提示词">
        <el-input v-model="prompt" type="textarea" :rows="3" :placeholder="promptPlaceholder" />
      </el-form-item>

      <!-- i2v: first frame -->
      <el-form-item v-if="isI2V" label="首帧图片">
        <div class="upload-area" @dragover.prevent @drop.prevent="onFirstFrameDrop" @click="triggerFFUpload">
          <input ref="ffInput" type="file" accept="image/*" style="display:none" @change="onFFChange" />
          <img v-if="firstFrame" :src="firstFrame" class="preview-img" />
          <span v-else class="upload-hint">拖拽、点击或 Ctrl+V 上传首帧图片 (≥300x300)</span>
        </div>
        <span v-if="imgError" class="img-error">{{ imgError }}</span>
        <el-input v-model="firstFrame" placeholder="或直接输入图片URL/Base64" clearable style="margin-top:8px" />
      </el-form-item>

      <!-- r2v: reference images -->
      <el-form-item v-if="isR2V" label="参考图片">
        <div class="ref-list">
          <div v-for="(img, i) in refImages" :key="i" class="ref-item">
            <img :src="img" class="ref-thumb" />
            <el-button type="danger" size="small" circle @click="refImages.splice(i,1)">X</el-button>
          </div>
        </div>
        <div v-if="refImages.length < 9" class="upload-area" @dragover.prevent @drop.prevent="onRefDrop" @click="triggerRefUpload">
          <input ref="refInput" type="file" accept="image/*" style="display:none" @change="onRefChange" />
          <span class="upload-hint">拖拽、点击或 Ctrl+V 上传参考图</span>
        </div>
        <span v-if="imgError" class="img-error">{{ imgError }}</span>
        <span class="hint">支持 1~9 张参考图，提示词中请用 [Image 1]、[Image 2] 引用。图片要求 ≥300x300 像素</span>
      </el-form-item>

      <!-- video-edit: source video + reference images -->
      <el-form-item v-if="isVideoEdit" label="源视频">
        <div v-if="videoPreview" class="video-preview">
          <video :src="videoPreview" controls style="max-width:300px;max-height:200px" />
        </div>
        <div class="upload-area" @dragover.prevent @drop.prevent="onVideoDrop" @click="triggerVideoUpload">
          <input ref="vidInput" type="file" accept="video/*" style="display:none" @change="onVideoChange" />
          <span class="upload-hint">拖拽或点击上传视频文件</span>
        </div>
        <el-input v-model="videoURL" placeholder="或直接输入视频URL" clearable style="margin-top:8px" />
      </el-form-item>
      <el-form-item v-if="isVideoEdit" label="参考图片">
        <div class="ref-list">
          <div v-for="(img, i) in refImages" :key="i" class="ref-item">
            <img :src="img" class="ref-thumb" />
            <el-button type="danger" size="small" circle @click="refImages.splice(i,1)">X</el-button>
          </div>
        </div>
        <div v-if="refImages.length < 5" class="upload-area" @dragover.prevent @drop.prevent="onRefDrop" @click="triggerRefUpload">
          <input ref="refInput2" type="file" accept="image/*" style="display:none" @change="onRefChange2" />
          <span class="upload-hint">拖拽、点击或 Ctrl+V 上传参考图</span>
        </div>
        <span v-if="imgError" class="img-error">{{ imgError }}</span>
        <span class="hint">可选 0~5 张参考图</span>
      </el-form-item>

      <el-row :gutter="16">
        <el-col :span="8">
          <el-form-item label="分辨率">
            <el-select v-model="resolution" style="width:100%">
              <el-option value="1080P" label="1080P" />
              <el-option value="720P" label="720P" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col v-if="!isVideoEdit" :span="8">
          <el-form-item label="画面比例">
            <el-select v-model="ratio" style="width:100%">
              <el-option value="16:9" label="16:9 (横屏)" />
              <el-option value="9:16" label="9:16 (竖屏)" />
              <el-option value="1:1" label="1:1 (方形)" />
              <el-option value="4:3" label="4:3" />
              <el-option value="3:4" label="3:4" />
              <el-option value="4:5" label="4:5" />
              <el-option value="5:4" label="5:4" />
              <el-option value="9:21" label="9:21" />
              <el-option value="21:9" label="21:9" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item v-if="!isVideoEdit" label="时长(秒)">
            <el-input-number v-model="duration" :min="3" :max="15" :step="1" style="width:100%" />
          </el-form-item>
          <el-form-item v-else label="音频设置">
            <el-select v-model="audioSetting" style="width:100%">
              <el-option value="auto" label="自动" />
              <el-option value="origin" label="保留原声" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :span="8">
          <el-form-item label="随机种子">
            <el-input v-model="seedStr" placeholder="留空随机" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="水印">
            <el-switch v-model="watermark" active-text="开" inactive-text="关" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item>
        <el-button type="primary" @click="generate" :loading="loading">生成视频</el-button>
      </el-form-item>
    </el-form>

    <div v-if="taskId" class="task-info">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="任务ID">{{ taskId }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusColor">{{ status }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
      <el-button @click="pollTask" :loading="polling" type="success" style="margin-top:8px">刷新状态</el-button>
      <div v-if="videoUrl" style="margin-top:16px">
        <video :src="videoUrl" controls style="max-width:640px; max-height:400px; width:100%" />
      </div>
    </div>

    <div v-if="raw" style="margin-top:16px">
      <el-collapse><el-collapse-item title="API 响应"><pre>{{ raw }}</pre></el-collapse-item></el-collapse>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { videoGen } from '../api'

const route = useRoute()
const router = useRouter()

const model = ref('happyhorse-1.0-t2v')
const prompt = ref('')
const resolution = ref('1080P')
const ratio = ref('16:9')
const duration = ref(5)
const seedStr = ref('')
const watermark = ref(false)
const firstFrame = ref('')
const refImages = ref([])
const videoURL = ref('')
const videoPreview = ref('')
const audioSetting = ref('auto')
const loading = ref(false)
const taskId = ref('')
const status = ref('')
const videoUrl = ref('')
const polling = ref(false)
const raw = ref('')

const ffInput = ref(null)
const refInput = ref(null)
const refInput2 = ref(null)
const vidInput = ref(null)

const isI2V = computed(() => model.value === 'happyhorse-1.0-i2v')
const isR2V = computed(() => model.value === 'happyhorse-1.0-r2v')
const isVideoEdit = computed(() => model.value === 'happyhorse-1.0-video-edit')

const promptPlaceholder = computed(() => {
  if (isR2V.value) return '描述你想要的视频，用 [Image 1]、[Image 2] 引用参考图...'
  if (isVideoEdit.value) return '描述编辑意图，如"让视频中的角色穿上条纹毛衣"...'
  return '描述你想要的视频...'
})

const statusColor = computed(() => {
  switch (status.value) {
    case 'SUCCEEDED': return 'success'
    case 'FAILED': return 'danger'
    case 'RUNNING': return 'warning'
    default: return 'info'
  }
})

const imgError = ref('')

function checkImageSize(file, onOk) {
  if (!file || !file.type.startsWith('image/')) return
  const reader = new FileReader()
  reader.onload = () => {
    const img = new Image()
    img.onload = () => {
      if (img.width < 300 || img.height < 300) {
        imgError.value = `图片尺寸 ${img.width}x${img.height} 太小，要求至少 300x300 像素`
        return
      }
      imgError.value = ''
      onOk(reader.result)
    }
    img.src = reader.result
  }
  reader.readAsDataURL(file)
}

// i2v first frame
function triggerFFUpload() { ffInput.value?.click() }
function onFFChange(e) { const f = e.target.files?.[0]; checkImageSize(f, v => firstFrame.value = v) }
function onFirstFrameDrop(e) { const f = e.dataTransfer?.files?.[0]; checkImageSize(f, v => firstFrame.value = v) }

// r2v / video-edit reference images
function triggerRefUpload() { (refInput.value || refInput2.value)?.click() }
function onRefChange(e) { const f = e.target.files?.[0]; checkImageSize(f, addRefImage) }
function onRefChange2(e) { const f = e.target.files?.[0]; checkImageSize(f, addRefImage) }
function onRefDrop(e) { const f = e.dataTransfer?.files?.[0]; checkImageSize(f, addRefImage) }

function addRefImage(url) {
  if (url && refImages.value.length < (isVideoEdit.value ? 5 : 9)) {
    refImages.value.push(url)
  }
}

// video-edit source video
function triggerVideoUpload() { vidInput.value?.click() }
function onVideoChange(e) {
  const f = e.target.files?.[0]
  if (!f || !f.type.startsWith('video/')) return
  const reader = new FileReader()
  reader.onload = () => {
    videoPreview.value = reader.result
    videoURL.value = reader.result
  }
  reader.readAsDataURL(f)
}
function onVideoDrop(e) {
  const f = e.dataTransfer?.files?.[0]
  if (!f || !f.type.startsWith('video/')) return
  const reader = new FileReader()
  reader.onload = () => {
    videoPreview.value = reader.result
    videoURL.value = reader.result
  }
  reader.readAsDataURL(f)
}

function onGlobalPaste(e) {
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (!file) continue
      if (isI2V.value) {
        checkImageSize(file, v => firstFrame.value = v)
      } else if (isR2V.value || isVideoEdit.value) {
        checkImageSize(file, addRefImage)
      }
      break
    }
  }
}

async function generate() {
  if (!prompt.value.trim()) return
  loading.value = true
  taskId.value = ''
  videoUrl.value = ''
  raw.value = ''
  try {
    const body = { model: model.value, prompt: prompt.value, resolution: resolution.value }
    if (!isVideoEdit.value) {
      body.duration = duration.value
      body.ratio = ratio.value
    } else {
      body.audio_setting = audioSetting.value
    }
    if (seedStr.value) body.seed = parseInt(seedStr.value)
    body.watermark = watermark.value
    if (isI2V.value && firstFrame.value) body.first_frame = firstFrame.value
    if ((isR2V.value || isVideoEdit.value) && refImages.value.length) body.ref_images = refImages.value
    if (isVideoEdit.value && videoURL.value) body.video_url = videoURL.value
    const { data } = await videoGen(body)
    raw.value = JSON.stringify(data, null, 2)
    taskId.value = data.task_id
    status.value = data.status
    if (data.task_id) {
      saveActiveTask(data.task_id, 'video')
      router.replace({ query: { task: data.task_id } })
      setTimeout(() => pollTask(), 10000)
    }
  } catch (e) {
    raw.value = e.response?.data || e.message
  }
  loading.value = false
}

async function pollTask() {
  if (!taskId.value) return
  polling.value = true
  try {
    const resp = await fetch('/api/tasks/' + taskId.value)
    const data = await resp.json()
    raw.value = JSON.stringify(data, null, 2)
    status.value = data.output?.task_status || 'unknown'
    const url = data.output?.video_url || data.output?.results?.[0]?.url
    if (url) videoUrl.value = url
    if (status.value === 'PENDING' || status.value === 'RUNNING') {
      setTimeout(() => pollTask(), 15000)
    } else {
      clearActiveTask()
      router.replace({ query: {} })
    }
  } catch (e) {
    raw.value = e.message
  }
  polling.value = false
}

const ACTIVE_TASK_KEY = 'bailian_active_task'
function saveActiveTask(id, type) {
  localStorage.setItem(ACTIVE_TASK_KEY, JSON.stringify({ id, type, savedAt: Date.now() }))
}
function clearActiveTask() {
  localStorage.removeItem(ACTIVE_TASK_KEY)
}
onMounted(() => {
  // 1) URL query takes priority (most reliable)
  const taskFromUrl = route.query.task
  if (taskFromUrl) {
    taskId.value = taskFromUrl
    pollTask()
    return
  }
  // 2) fallback to localStorage
  try {
    const saved = JSON.parse(localStorage.getItem(ACTIVE_TASK_KEY))
    if (saved && saved.id && (Date.now() - saved.savedAt < 86400000)) {
      taskId.value = saved.id
      router.replace({ query: { task: saved.id } })
      pollTask()
    }
  } catch {}
})
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.task-info { background: #f5f7fa; padding: 16px; border-radius: 8px; margin-top: 16px; }
pre { white-space: pre-wrap; font-size: 12px; max-height: 400px; overflow-y: auto; }
.hint { color: #909399; font-size: 12px; display: block; margin-top: 4px; }
.upload-area { border: 2px dashed #dcdfe6; border-radius: 8px; padding: 20px; text-align: center; cursor: pointer; transition: border-color .2s; }
.upload-area:hover { border-color: #409eff; }
.upload-hint { color: #909399; font-size: 13px; }
.preview-img { max-width: 200px; max-height: 150px; border-radius: 4px; }
.ref-list { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 8px; }
.ref-item { position: relative; }
.ref-thumb { width: 80px; height: 80px; object-fit: cover; border-radius: 4px; border: 1px solid #dcdfe6; }
.ref-item .el-button { position: absolute; top: -8px; right: -8px; }
.video-preview { margin-bottom: 8px; }
.img-error { color: #f56c6c; font-size: 12px; display: block; margin-top: 4px; }
</style>