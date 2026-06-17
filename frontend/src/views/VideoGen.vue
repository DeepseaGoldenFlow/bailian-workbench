<template>
  <div class="pa-8" style="max-width:1300px;margin:0 auto">
    <div class="d-flex align-center mb-8">
      <v-avatar color="teal" size="48" class="mr-4 elevation-4"><v-icon size="28">mdi-video-outline</v-icon></v-avatar>
      <div><h1 class="text-h4 font-weight-bold">Video Generation</h1><p class="text-body-1 text-medium-emphasis mt-1">HappyHorse video creation via Bailian</p></div>
    </div>
    <v-row>
      <v-col cols="12" md="4">
        <v-card rounded="xl" class="pa-6" elevation="2" style="background:linear-gradient(180deg, rgba(0,150,136,0.08) 0%, transparent 40%)">
          <v-select v-model="model" :items="models" item-title="name" item-value="id" label="Model" variant="outlined" density="comfortable" hide-details class="mb-5" color="teal" />
          <v-textarea v-model="prompt" label="Prompt" variant="outlined" rows="3" density="comfortable" :placeholder="isEdit ? 'Describe edit intent...' : 'Describe the video...'" hide-details class="mb-5" />
          <v-text-field v-if="isEdit" v-model="videoUrl" label="Video URL" variant="outlined" density="comfortable" placeholder="MP4/MOV, 3-60s" hide-details class="mb-5" />
          <v-text-field v-if="isI2V" v-model="firstFrame" label="First Frame URL" variant="outlined" density="comfortable" hide-details class="mb-5" />
          <v-textarea v-if="isR2V || isEdit" v-model="refImagesStr" :label="'Ref Images' + (isR2V ? ' (1-9)' : ' (0-5)')" variant="outlined" rows="2" density="comfortable" hide-details class="mb-5" />
          <v-select v-model="resolution" :items="['1080P','720P']" label="Resolution" variant="outlined" density="comfortable" hide-details class="mb-5" />
          <div class="mb-5" v-if="!isEdit"><div class="text-caption text-medium-emphasis mb-1">Duration: {{ duration }}s</div><v-slider v-model="duration" min="3" max="15" step="1" hide-details density="compact" color="teal" thumb-label /></div>
          <v-select v-if="!isEdit" v-model="ratio" :items="['16:9','9:16','1:1','4:3','3:4']" label="Ratio" variant="outlined" density="comfortable" hide-details class="mb-5" />
          <v-text-field v-model="seed" label="Seed (optional)" variant="outlined" density="comfortable" type="number" hide-details class="mb-5" />
          <v-switch v-model="watermark" label="Watermark" color="teal" density="compact" hide-details class="mb-5" />
          <v-select v-model="audioSetting" :items="['auto','origin']" label="Audio" variant="outlined" density="comfortable" hide-details class="mb-5" />
          <v-btn block color="teal" size="x-large" rounded="lg" :loading="loading" :disabled="!canSubmit" @click="generate" variant="elevated"><v-icon start>mdi-sparkles</v-icon> Submit</v-btn>
        </v-card>
      </v-col>
      <v-col cols="12" md="8">
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4" closable>{{ error }}</v-alert>
        <v-card v-if="taskId" rounded="xl" class="pa-5 mb-4" variant="outlined">
          <div class="d-flex align-center ga-3 mb-3"><v-chip :color="statusColor" size="small" label>{{ taskStatus }}</v-chip><span class="text-caption text-medium-emphasis">{{ taskId }}</span></div>
          <v-progress-linear v-if="taskStatus === 'PENDING' || taskStatus === 'RUNNING'" indeterminate color="teal" rounded class="mb-3" />
          <v-btn v-if="taskStatus !== 'SUCCEEDED' && taskStatus !== 'FAILED'" size="small" variant="tonal" color="teal" @click="pollTaskResult" :loading="polling">Check Status</v-btn>
        </v-card>
        <v-card v-if="resultUrl" rounded="xl" variant="outlined">
          <video :src="resultUrl" controls style="width:100%;max-height:500px;border-radius:16px" />
          <v-card-text><v-btn :href="resultUrl" target="_blank" variant="text" color="teal" prepend-icon="mdi-download">Download</v-btn></v-card-text>
        </v-card>
        <v-card v-if="!loading && !resultUrl && !error && !taskId" rounded="xl" class="pa-16 text-center" variant="outlined">
          <v-icon size="64" color="teal-lighten-3" class="mb-4">mdi-video-plus-outline</v-icon><div class="text-h6 text-medium-emphasis mb-1">Ready to Create</div><p class="text-body-2 text-medium-emphasis">Choose a model and describe your video</p>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchModels, videoGen, pollTask } from '../api'
const models=ref([]);const model=ref('');const prompt=ref('');const videoUrl=ref('');const firstFrame=ref('')
const refImagesStr=ref('');const resolution=ref('1080P');const duration=ref(5);const ratio=ref('16:9')
const seed=ref(null);const watermark=ref(true);const audioSetting=ref('auto')
const loading=ref(false);const taskId=ref('');const taskStatus=ref('');const polling=ref(false)
const resultUrl=ref('');const error=ref('')
const isEdit=computed(()=>model.value==='happyhorse-1.0-video-edit')
const isI2V=computed(()=>model.value==='happyhorse-1.0-i2v')
const isR2V=computed(()=>model.value==='happyhorse-1.0-r2v')
const canSubmit=computed(()=>{if(!model.value||!prompt.value)return false;if(isEdit.value&&!videoUrl.value)return false;if(isI2V.value&&!firstFrame.value)return false;if(isR2V.value&&!refImagesStr.value)return false;return true})
onMounted(async()=>{try{const r=await fetchModels('video');models.value=r.data.models||[];if(models.value.length)model.value=models.value[0].id}catch(e){error.value='Failed to load models'}})
const statusColor=computed(()=>{if(taskStatus.value==='SUCCEEDED')return'success';if(taskStatus.value==='FAILED')return'error';return'warning'})
function parseRefImages(){if(!refImagesStr.value)return[];var NL=String.fromCharCode(10);return refImagesStr.value.replace(RegExp(NL,"g"),",").split(",").map(s=>s.trim()).filter(Boolean)}
async function generate(){loading.value=true;error.value='';taskId.value='';resultUrl.value=''
try{const payload={model:model.value,prompt:prompt.value,resolution:resolution.value,seed:seed.value||undefined,watermark:watermark.value,audio_setting:audioSetting.value}
if(!isEdit.value){payload.duration=duration.value;payload.ratio=ratio.value}
if(isEdit.value){payload.video_url=videoUrl.value;payload.ref_images=parseRefImages()}
if(isI2V.value)payload.first_frame=firstFrame.value
if(isR2V.value)payload.ref_images=parseRefImages()
const r=await videoGen(payload);const data=r.data;taskId.value=data.task_id;taskStatus.value=data.status
if(data.status==='PENDING'||data.status==='RUNNING')await pollTaskResult()}
catch(e){error.value=e.response?.data?.message||e.message||'Request failed'};loading.value=false}
async function pollTaskResult(){if(!taskId.value)return;polling.value=true
try{const r=await pollTask(taskId.value);const data=r.data;taskStatus.value=data.output?.task_status||'UNKNOWN';if(taskStatus.value==='SUCCEEDED')resultUrl.value=data.output?.video_url||'';if(taskStatus.value==='FAILED')error.value=data.output?.message||'Task failed'}
catch(e){error.value='Poll failed'};polling.value=false}
</script>