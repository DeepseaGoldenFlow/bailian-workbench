<template>
  <div class="page">
    <header class="page-header"><h1>Image Generation</h1><p class="page-desc">AI image creation via Bailian models</p></header>
    <div class="content-grid">
      <div class="config-panel">
        <div class="field"><label class="field-label">Model</label><select v-model="model" class="select" @change="onModelChange"><option v-for="m in models" :key="m.id" :value="m.id">{{ m.name }}</option></select></div>
        <div class="field"><label class="field-label">Prompt</label><textarea v-model="prompt" class="textarea" rows="3" placeholder="Describe the image you want..."></textarea></div>
        <div class="field" v-if="hasParam('negative_prompt')"><label class="field-label">Neg Prompt</label><input v-model="negativePrompt" class="input" placeholder="Elements to avoid..." /></div>
        <div class="field-row"><div class="field"><label class="field-label">Size</label><select v-model="size" class="select" v-if="getOptions('size').length"><option v-for="o in getOptions('size')" :key="o.value" :value="o.value">{{ o.label }}</option></select></div><div class="field"><label class="field-label">Count</label><input type="range" v-model.number="n" min="1" max="4" class="range" /><span class="range-val">{{ n }}</span></div></div>
        <div class="field" v-if="hasParam('ref_img')"><label class="field-label">Ref Image URL</label><input v-model="refImg" class="input" placeholder="https://..." /></div>
        <div class="field" v-if="hasParam('steps')"><label class="field-label">Steps: {{ steps }}</label><input type="range" v-model.number="steps" min="1" max="100" class="range" /></div>
        <div class="field-row" v-if="hasParam('prompt_extend')"><label class="field-label">Prompt Extend</label><label class="toggle"><input type="checkbox" v-model="promptExtend" /><span class="toggle-slider"></span></label></div>
        <button class="btn-primary" @click="generate" :disabled="loading || !prompt"><span v-if="loading" class="spinner"></span>{{ loading ? 'Generating...' : 'Generate' }}</button>
      </div>
      <div class="result-panel">
        <div v-if="error" class="error-card"><div class="error-icon">!</div><div class="error-text">{{ error }}</div></div>
        <div v-if="taskId" class="task-card"><div class="task-header"><span class="task-label">Task</span><span class="task-id">{{ taskId }}</span><span class="task-status" :class="taskStatus">{{ taskStatus }}</span></div><div v-if="taskStatus === 'PENDING' || taskStatus === 'RUNNING'" class="progress-bar"><div class="progress-indeterminate"></div></div><button v-if="taskStatus !== 'SUCCEEDED' && taskStatus !== 'FAILED'" class="btn-small" @click="pollTaskResult" :disabled="polling">Refresh</button></div>
        <div v-if="images.length" class="image-grid"><div v-for="(img, i) in images" :key="i" class="image-card" @click="previewIndex = i"><img :src="img" loading="lazy" /></div></div>
        <div v-if="!loading && !images.length && !error && !taskId" class="empty-state"><div class="empty-icon">+</div><p>Enter a prompt</p></div>
      </div>
    </div>
    <div v-if="previewIndex !== null" class="preview-overlay" @click="previewIndex = null"><img :src="images[previewIndex]" class="preview-img" /><button class="preview-close">&times;</button></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchModels, imageGen, pollTask } from '../api'

const models = ref([])
const model = ref('')
const prompt = ref('')
const negativePrompt = ref('')
const size = ref('1024*1024')
const n = ref(1)
const steps = ref(30)
const promptExtend = ref(true)
const refImg = ref('')
const loading = ref(false)
const taskId = ref('')
const taskStatus = ref('')
const polling = ref(false)
const images = ref([])
const error = ref('')
const previewIndex = ref(null)

const currentModel = computed(() => models.value.find(m => m.id === model.value))

onMounted(async () => {
  try {
    const r = await fetchModels('image')
    models.value = r.data.models || []
    if (models.value.length) model.value = models.value[0].id
  } catch (e) { error.value = 'Failed to load models' }
})

function hasParam(name) { return currentModel.value?.parameters?.some(p => p.name === name) }
function getOptions(name) { return currentModel.value?.parameters?.find(p => p.name === name)?.options || [] }
function onModelChange() { images.value = []; taskId.value = ''; error.value = '' }

async function generate() {
  loading.value = true; error.value = ''; images.value = []; taskId.value = ''
  try {
    const payload = { model: model.value, prompt: prompt.value, size: size.value, n: n.value }
    if (hasParam('steps')) payload.steps = steps.value
    if (hasParam('prompt_extend')) payload.prompt_extend = promptExtend.value
    if (negativePrompt.value) payload.negative_prompt = negativePrompt.value
    if (refImg.value) { payload.ref_img = refImg.value; payload.ref_strength = 0.5 }
    const r = await imageGen(payload)
    const data = r.data
    if (data.task_id) {
      taskId.value = data.task_id
      taskStatus.value = data.status
      if (data.status === 'PENDING' || data.status === 'RUNNING') await pollTaskResult()
    } else {
      extractImages(data)
    }
  } catch (e) {
    error.value = e.response?.data?.message || e.response?.data?.code || e.message || 'Request failed'
    try { error.value += ' | ' + JSON.stringify(e.response?.data || {}).substring(0, 200) } catch(_) {}
  }
  loading.value = false
}

async function pollTaskResult() {
  if (!taskId.value) return
  polling.value = true
  try {
    const r = await pollTask(taskId.value)
    const data = r.data
    taskStatus.value = data.output?.task_status || 'UNKNOWN'
    if (taskStatus.value === 'SUCCEEDED') extractTaskImages(data)
    if (taskStatus.value === 'FAILED') error.value = data.output?.message || 'Task failed'
  } catch (e) { error.value = 'Poll failed' }
  polling.value = false
}

function extractImages(data) {
  const result = []
  if (data.output?.results) result.push(...data.output.results.map(r => r.url).filter(Boolean))
  if (data.output?.choices) {
    for (const c of data.output.choices) {
      const ct = c.message?.content
      if (Array.isArray(ct)) {
        for (const it of ct) {
          if (it.image) result.push(it.image)
        }
      }
    }
  }
  images.value = result
  if (!result.length) {
    try { error.value = 'No image URL: ' + JSON.stringify(data).substring(0, 300) } catch(_) {}
  }
}

function extractTaskImages(data) {
  const result = []
  if (data.output?.results) result.push(...data.output.results.map(r => r.url).filter(Boolean))
  if (data.output?.video_url) result.push(data.output.video_url)
  images.value = result
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
.page-header { margin-bottom: 32px; }
.page-header h1 { font-size: 28px; font-weight: 700; letter-spacing: -0.5px; }
.page-desc { color: var(--text-secondary); margin-top: 4px; font-size: 14px; }
.content-grid { display: grid; grid-template-columns: 380px 1fr; gap: 24px; align-items: start; }
.config-panel { background: var(--bg-card); border: 1px solid var(--border-card); border-radius: var(--radius-lg); padding: 24px; display: flex; flex-direction: column; gap: 16px; }
.field { display: flex; flex-direction: column; gap: 6px; }
.field-row { display: flex; gap: 16px; align-items: end; }
.field-row .field { flex: 1; }
.field-label { font-size: 12px; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.5px; }
.input, .textarea, .select { background: var(--bg-secondary); border: 1px solid var(--border-card); border-radius: var(--radius-sm); color: var(--text-primary); padding: 10px 12px; font-size: 14px; font-family: var(--font-sans); outline: none; transition: border var(--transition); width: 100%; resize: vertical; }
.input:focus, .textarea:focus, .select:focus { border-color: var(--accent); }
.select { cursor: pointer; }
.range { width: 100%; accent-color: var(--accent); }
.range-val { font-size: 12px; color: var(--text-secondary); margin-left: 8px; }
.toggle { position: relative; display: inline-block; width: 44px; height: 24px; cursor: pointer; }
.toggle input { opacity: 0; width: 0; height: 0; }
.toggle-slider { position: absolute; inset: 0; background: var(--bg-secondary); border-radius: 12px; transition: var(--transition); border: 1px solid var(--border-card); }
.toggle input:checked + .toggle-slider { background: var(--accent); border-color: var(--accent); }
.toggle-slider::after { content: ""; position: absolute; width: 18px; height: 18px; left: 2px; top: 2px; background: #fff; border-radius: 50%; transition: var(--transition); }
.toggle input:checked + .toggle-slider::after { transform: translateX(20px); }
.btn-primary { display: flex; align-items: center; justify-content: center; gap: 8px; padding: 12px 24px; background: linear-gradient(135deg, var(--accent), #818cf8); border: none; border-radius: var(--radius-sm); color: #fff; font-size: 15px; font-weight: 600; cursor: pointer; transition: all var(--transition); width: 100%; }
.btn-primary:hover:not(:disabled) { box-shadow: var(--shadow-glow); transform: translateY(-1px); }
.btn-primary:disabled { opacity: 0.4; cursor: not-allowed; }
.spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: spin 0.8s linear infinite; display: inline-block; }
@keyframes spin { to { transform: rotate(360deg); } }
.btn-small { padding: 6px 16px; background: var(--bg-secondary); border: 1px solid var(--border-card); border-radius: var(--radius-sm); color: var(--text-secondary); cursor: pointer; font-size: 13px; }
.btn-small:hover { background: var(--bg-card-hover); }
.result-panel { display: flex; flex-direction: column; gap: 16px; }
.error-card { display: flex; gap: 12px; padding: 16px; background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.2); border-radius: var(--radius-md); }
.error-icon { width: 24px; height: 24px; background: var(--danger); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; color: #fff; flex-shrink: 0; }
.error-text { font-size: 13px; color: #fca5a5; word-break: break-all; }
.task-card { padding: 16px; background: var(--bg-card); border: 1px solid var(--border-card); border-radius: var(--radius-md); }
.task-header { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.task-label { color: var(--text-muted); }
.task-id { color: var(--text-secondary); font-family: monospace; font-size: 12px; }
.task-status { padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; text-transform: uppercase; }
.task-status.PENDING, .task-status.RUNNING { background: rgba(245,158,11,0.15); color: var(--warning); }
.task-status.SUCCEEDED { background: rgba(34,197,94,0.15); color: var(--success); }
.task-status.FAILED { background: rgba(239,68,68,0.15); color: var(--danger); }
.progress-bar { height: 3px; background: var(--bg-secondary); border-radius: 2px; margin: 12px 0; overflow: hidden; }
.progress-indeterminate { width: 40%; height: 100%; background: var(--accent); border-radius: 2px; animation: progress 1.5s ease-in-out infinite; }
@keyframes progress { 0% { transform: translateX(-100%); } 100% { transform: translateX(350%); } }
.image-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; }
.image-card { border-radius: var(--radius-md); overflow: hidden; cursor: pointer; border: 1px solid var(--border-card); background: var(--bg-card); transition: all var(--transition); }
.image-card:hover { border-color: var(--accent); transform: scale(1.02); }
.image-card img { width: 100%; display: block; }
.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 64px 32px; color: var(--text-muted); }
.empty-icon { width: 48px; height: 48px; border: 2px dashed var(--border-card); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 24px; margin-bottom: 12px; }
.empty-state p { font-size: 14px; }
.preview-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.9); display: flex; align-items: center; justify-content: center; z-index: 100; cursor: pointer; }
.preview-img { max-width: 90vw; max-height: 90vh; border-radius: var(--radius-md); }
.preview-close { position: absolute; top: 24px; right: 24px; background: none; border: none; color: #fff; font-size: 32px; cursor: pointer; }
</style>
