<template>
  <div class="page">
    <h2>图片生成</h2>
    <el-form label-width="100px" style="max-width:700px">
      <el-form-item label="模型">
        <el-select v-model="model" style="width:100%">
          <el-option label="Qwen-Image-2.0-Pro (推荐)" value="qwen-image-2.0-pro" />
          <el-option label="Wanx2.1-T2I-Turbo" value="wanx2.1-t2i-turbo" />
        </el-select>
      </el-form-item>
      <el-form-item label="提示词">
        <el-input v-model="prompt" type="textarea" :rows="3" placeholder="描述你想要的图片..." />
      </el-form-item>
      <el-form-item label="反向提示词">
        <el-input v-model="negativePrompt" type="textarea" :rows="2" placeholder="不希望出现的内容（可选）..." />
      </el-form-item>
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="尺寸">
            <el-select v-model="size" style="width:100%">
              <el-option label="1024x1024 (1:1)" value="1024*1024" />
              <el-option label="720x1280 (竖屏)" value="720*1280" />
              <el-option label="1280x720 (横屏)" value="1280*720" />
              <el-option label="1920x1080 (16:9)" value="1920*1080" />
              <el-option label="1024x768 (4:3)" value="1024*768" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="数量">
            <el-input-number v-model="n" :min="1" :max="4" style="width:100%" />
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
          <el-form-item label="采样步数">
            <el-input-number v-model="steps" :min="0" :max="100" :step="5" style="width:100%" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="提示词扩展">
            <el-switch v-model="promptExtend" active-text="开" inactive-text="关" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="参考图片">
        <div class="ref-upload" @dragover.prevent @drop.prevent="onRefDrop" @click="triggerRefUpload" @paste="onRefPaste" tabindex="0">
          <input ref="refFileInput" type="file" accept="image/*" style="display:none" @change="onRefFileChange" />
          <div v-if="!refPreview" class="ref-placeholder">拖拽图片、点击或 Ctrl+V 粘贴</div>
          <img v-else :src="refPreview" class="ref-preview-img" />
        </div>
        <el-input v-model="refImg" placeholder="或直接输入图片URL" clearable style="margin-top:8px" />
      </el-form-item>
      <el-form-item v-if="refImg" label="参考强度">
        <el-slider v-model="refStrength" :min="0" :max="1" :step="0.1" show-input style="width:100%" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="generate" :loading="loading">生成图片</el-button>
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
    </div>

    <div v-if="results.length" class="results">
      <div v-for="(img, i) in results" :key="i" class="img-card">
        <el-image :src="img" fit="contain" style="max-width:400px;max-height:400px" :preview-src-list="results" />
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
import { imageGen } from '../api'

const route = useRoute()
const router = useRouter()

const model = ref('qwen-image-2.0-pro')
const prompt = ref('')
const negativePrompt = ref('')
const size = ref('1024*1024')
const n = ref(1)
const seedStr = ref('')
const steps = ref(0)
const promptExtend = ref(true)
const refImg = ref('')
const refStrength = ref(0.5)
const refPreview = ref('')
const refFileInput = ref(null)
const loading = ref(false)
const taskId = ref('')
const status = ref('')
const polling = ref(false)
const results = ref([])
const raw = ref('')

const statusColor = computed(() => {
  switch (status.value) {
    case 'SUCCEEDED': return 'success'
    case 'FAILED': return 'danger'
    case 'RUNNING': return 'warning'
    default: return 'info'
  }
})

async function generate() {
  if (!prompt.value.trim()) return
  loading.value = true
  results.value = []
  taskId.value = ''
  raw.value = ''
  try {
    const body = {
      model: model.value,
      prompt: prompt.value,
      size: size.value,
      n: n.value,
    }
    if (negativePrompt.value) body.negative_prompt = negativePrompt.value
    if (seedStr.value) body.seed = parseInt(seedStr.value)
    if (steps.value > 0) body.steps = steps.value
    body.prompt_extend = promptExtend.value
    if (refImg.value) {
      body.ref_img = refImg.value
      body.ref_strength = refStrength.value
    }
    const { data } = await imageGen(body)
    raw.value = JSON.stringify(data, null, 2)
    taskId.value = data.task_id
    status.value = data.status
    if (data.task_id) {
      saveActiveTask(data.task_id, 'image')
      router.replace({ query: { task: data.task_id } })
      setTimeout(() => pollTask(), 3000)
    }
  } catch (e) {
    raw.value = e.response?.data || e.message
  }
  loading.value = false
}

function triggerRefUpload() {
  refFileInput.value?.click()
}
function onRefFileChange(e) {
  const file = e.target.files?.[0]
  if (file) loadRefImage(file)
}
function onRefDrop(e) {
  const file = e.dataTransfer?.files?.[0]
  if (file) loadRefImage(file)
}
function onRefPaste(e) {
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      loadRefImage(item.getAsFile())
      break
    }
  }
}
function loadRefImage(file) {
  if (!file || !file.type.startsWith('image/')) return
  const reader = new FileReader()
  reader.onload = () => {
    refPreview.value = reader.result
    refImg.value = reader.result
  }
  reader.readAsDataURL(file)
}

async function pollTask() {
  if (!taskId.value) return
  polling.value = true
  try {
    const resp = await fetch('/api/tasks/' + taskId.value)
    const data = await resp.json()
    raw.value = JSON.stringify(data, null, 2)
    status.value = data.output?.task_status || 'unknown'
    if (data.output?.results) {
      results.value = data.output.results.map(r => r.url).filter(Boolean)
    }
    if (status.value === 'PENDING' || status.value === 'RUNNING') {
      setTimeout(() => pollTask(), 5000)
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
  const taskFromUrl = route.query.task
  if (taskFromUrl) {
    taskId.value = taskFromUrl
    pollTask()
    return
  }
  try {
    const saved = JSON.parse(localStorage.getItem(ACTIVE_TASK_KEY))
    if (saved && saved.id && saved.type === 'image' && (Date.now() - saved.savedAt < 86400000)) {
      taskId.value = saved.id
      router.replace({ query: { task: saved.id } })
      pollTask()
    }
  } catch {}
})
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.results { display: flex; flex-wrap: wrap; gap: 16px; margin-top: 16px; }
.img-card { border: 1px solid #e0e0e0; border-radius: 8px; padding: 8px; background: #fff; }
.task-info { margin-top: 16px; padding: 16px; background: #f5f7fa; border-radius: 8px; }
pre { white-space: pre-wrap; font-size: 12px; max-height: 400px; overflow-y: auto; }
.ref-upload { border: 2px dashed #dcdfe6; border-radius: 8px; padding: 16px; text-align: center; cursor: pointer; transition: border-color .2s; }
.ref-upload:hover { border-color: #409eff; }
.ref-placeholder { color: #909399; font-size: 13px; }
.ref-preview-img { max-width: 200px; max-height: 150px; border-radius: 4px; }
</style>