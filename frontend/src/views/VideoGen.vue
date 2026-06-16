<template>
  <v-container fluid class="pa-6" style="max-width:1200px">
    <h1 class="text-h4 font-weight-bold mb-1">Video Generation</h1>
    <p class="text-body-1 text-medium-emphasis mb-6">HappyHorse video generation via Bailian</p>

    <v-row>
      <v-col cols="12" md="4">
        <v-card variant="outlined" class="pa-4" rounded="xl">
          <v-select v-model="model" :items="models" item-title="name" item-value="id" label="Model" variant="outlined" density="comfortable" hide-details class="mb-4" />

          <v-textarea v-model="prompt" label="Prompt" variant="outlined" rows="3" density="comfortable" :placeholder="isEdit ? 'Describe the edit intent...' : 'Describe the video...'" hide-details class="mb-4" />

          <v-text-field v-if="isEdit" v-model="videoUrl" label="Video URL" variant="outlined" density="comfortable" placeholder="MP4/MOV, 3-60s, max 100MB" hide-details class="mb-4" />

          <v-text-field v-if="isI2V" v-model="firstFrame" label="First Frame Image URL" variant="outlined" density="comfortable" placeholder="Image URL for first frame" hide-details class="mb-4" />

          <v-textarea v-if="isR2V || isEdit" v-model="refImagesStr" :label="'Reference Images' + (isR2V ? ' (1-9)' : ' (0-5)')" variant="outlined" rows="2" density="comfortable" placeholder="URLs separated by comma or newline" hide-details class="mb-4" />

          <v-select v-model="resolution" :items="['1080P','720P']" label="Resolution" variant="outlined" density="comfortable" hide-details class="mb-4" />

          <div class="mb-4" v-if="!isEdit">
            <div class="text-caption text-medium-emphasis mb-1">Duration: {{ duration }}s</div>
            <v-slider v-model="duration" min="3" max="15" step="1" hide-details density="compact" color="primary" thumb-label />
          </div>

          <v-select v-if="!isEdit" v-model="ratio" :items="ratios" label="Aspect Ratio" variant="outlined" density="comfortable" hide-details class="mb-4" />

          <v-text-field v-model="seed" label="Seed (optional)" variant="outlined" density="comfortable" type="number" placeholder="Random" hide-details class="mb-4" />

          <v-switch v-model="watermark" label="Watermark" color="primary" density="compact" hide-details class="mb-4" />

          <v-select v-model="audioSetting" :items="['auto','origin']" label="Audio" variant="outlined" density="comfortable" hide-details class="mb-4" />

          <v-btn block color="primary" size="large" rounded="lg" :loading="loading" :disabled="!canSubmit" @click="generate" variant="flat">Submit</v-btn>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4" closable>{{ error }}</v-alert>

        <v-card v-if="taskId" variant="outlined" class="pa-4 mb-4" rounded="xl">
          <div class="d-flex align-center ga-3 mb-2">
            <v-chip :color="statusColor" size="small" label>{{ taskStatus }}</v-chip>
            <span class="text-caption text-medium-emphasis">{{ taskId }}</span>
          </div>
          <v-progress-linear v-if="taskStatus === 'PENDING' || taskStatus === 'RUNNING'" indeterminate color="primary" class="mb-2" />
          <v-btn v-if="taskStatus !== 'SUCCEEDED' && taskStatus !== 'FAILED'" size="small" variant="text" color="primary" @click="pollTaskResult" :loading="polling">Check Status</v-btn>
        </v-card>

        <v-card v-if="resultUrl" variant="outlined" rounded="xl">
          <video :src="resultUrl" controls style="width:100%;max-height:500px" />
          <v-card-text>
            <v-btn :href="resultUrl" target="_blank" variant="text" color="primary" prepend-icon="mdi-download">Download</v-btn>
            <span class="text-caption text-medium-emphasis ml-2">Link valid for 24 hours</span>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchModels, videoGen, pollTask } from '../api'

const models = ref([]); const model = ref(''); const prompt = ref('')
const videoUrl = ref(''); const firstFrame = ref(''); const refImagesStr = ref('')
const resolution = ref('1080P'); const duration = ref(5); const ratio = ref('16:9')
const seed = ref(null); const watermark = ref(true); const audioSetting = ref('auto')
const loading = ref(false); const taskId = ref(''); const taskStatus = ref('')
const polling = ref(false); const resultUrl = ref(''); const error = ref('')
const ratios = ['16:9', '9:16', '1:1', '4:3', '3:4']

const isEdit = computed(() => model.value === 'happyhorse-1.0-video-edit')
const isI2V = computed(() => model.value === 'happyhorse-1.0-i2v')
const isR2V = computed(() => model.value === 'happyhorse-1.0-r2v')
const canSubmit = computed(() => { if (!model.value || !prompt.value) return false; if (isEdit.value && !videoUrl.value) return false; if (isI2V.value && !firstFrame.value) return false; if (isR2V.value && !refImagesStr.value) return false; return true })
const statusColor = computed(() => { if (taskStatus.value === 'SUCCEEDED') return 'success'; if (taskStatus.value === 'FAILED') return 'error'; return 'warning' })

onMounted(async () => { try { const r = await fetchModels('video'); models.value = r.data.models || []; if (models.value.length) model.value = models.value[0].id } catch (e) { error.value = 'Failed to load models' } })

function parseRefImages() { if (!refImagesStr.value) return []; return refImagesStr.value.split(/[,\n]/).map(s => s.trim()).filter(Boolean) }

async function generate() {
  loading.value = true; error.value = ''; taskId.value = ''; resultUrl.value = ''
  try {
    const payload = { model: model.value, prompt: prompt.value, resolution: resolution.value, seed: seed.value || undefined, watermark: watermark.value, audio_setting: audioSetting.value }
    if (!isEdit.value) { payload.duration = duration.value; payload.ratio = ratio.value }
    if (isEdit.value) { payload.video_url = videoUrl.value; payload.ref_images = parseRefImages() }
    if (isI2V.value) payload.first_frame = firstFrame.value
    if (isR2V.value) payload.ref_images = parseRefImages()
    const r = await videoGen(payload); const data = r.data
    taskId.value = data.task_id; taskStatus.value = data.status
    if (data.status === 'PENDING' || data.status === 'RUNNING') await pollTaskResult()
  } catch (e) { error.value = e.response?.data?.message || e.message || 'Request failed' }
  loading.value = false
}

async function pollTaskResult() { if (!taskId.value) return; polling.value = true; try { const r = await pollTask(taskId.value); const data = r.data; taskStatus.value = data.output?.task_status || 'UNKNOWN'; if (taskStatus.value === 'SUCCEEDED') resultUrl.value = data.output?.video_url || ''; if (taskStatus.value === 'FAILED') error.value = data.output?.message || 'Task failed' } catch (e) { error.value = 'Poll failed' }; polling.value = false }
</script>
