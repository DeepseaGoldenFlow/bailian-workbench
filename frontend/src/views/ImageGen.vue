<template>
  <div class="pa-8" style="max-width:1300px;margin:0 auto">
    <div class="d-flex align-center mb-8">
      <v-avatar color="purple" size="48" class="mr-4 elevation-4">
        <v-icon size="28">mdi-image-outline</v-icon>
      </v-avatar>
      <div>
        <h1 class="text-h4 font-weight-bold">Image Generation</h1>
        <p class="text-body-1 text-medium-emphasis mt-1">Create stunning visuals with Bailian AI models</p>
      </div>
    </div>

    <v-row>
      <v-col cols="12" md="4">
        <v-card rounded="xl" class="pa-6" elevation="2" style="background:linear-gradient(180deg, rgba(147,51,234,0.08) 0%, transparent 40%)">
          <v-select v-model="model" :items="models" item-title="name" item-value="id" label="Model" variant="outlined" density="comfortable" hide-details class="mb-5" color="purple" />
          <v-textarea v-model="prompt" label="Prompt" variant="outlined" rows="3" density="comfortable" placeholder="Describe the image you want to create..." hide-details class="mb-5" />
          <v-text-field v-if="hasParam('negative_prompt')" v-model="negativePrompt" label="Negative Prompt" variant="outlined" density="comfortable" placeholder="Elements to avoid" hide-details class="mb-5" />
          <v-row dense class="mb-5">
            <v-col cols="7"><v-select v-if="getOptions('size').length" v-model="size" :items="getOptions('size')" item-title="label" item-value="value" label="Size" variant="outlined" density="comfortable" hide-details /></v-col>
            <v-col cols="5"><div class="text-caption text-medium-emphasis mb-1">Count: {{ n }}</div><v-slider v-model="n" min="1" max="4" step="1" hide-details density="compact" color="purple" thumb-label /></v-col>
          </v-row>
          <v-text-field v-if="hasParam('ref_img')" v-model="refImg" label="Reference Image URL" variant="outlined" density="comfortable" placeholder="https://..." hide-details class="mb-5" />
          <div class="mb-5" v-if="hasParam('steps')"><div class="text-caption text-medium-emphasis mb-1">Inference Steps: {{ steps }}</div><v-slider v-model="steps" min="1" max="100" step="1" hide-details density="compact" color="purple" /></div>
          <v-switch v-if="hasParam('prompt_extend')" v-model="promptExtend" label="Prompt Enhance" color="purple" density="compact" hide-details class="mb-5" />
          <v-btn block color="purple" size="x-large" rounded="lg" :loading="loading" :disabled="!prompt" @click="generate" variant="elevated" class="text-none">
            <v-icon start>mdi-sparkles</v-icon> {{ loading ? 'Creating...' : 'Generate Image' }}
          </v-btn>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4" closable><template #text>{{ error }}</template></v-alert>

        <v-card v-if="taskId" rounded="xl" class="pa-5 mb-4" variant="outlined">
          <div class="d-flex align-center ga-3 mb-3"><v-chip :color="statusColor" size="small" label variant="flat">{{ taskStatus }}</v-chip><span class="text-caption text-medium-emphasis">{{ taskId }}</span></div>
          <v-progress-linear v-if="taskStatus === 'PENDING' || taskStatus === 'RUNNING'" indeterminate color="purple" rounded class="mb-3" />
          <v-btn v-if="taskStatus !== 'SUCCEEDED' && taskStatus !== 'FAILED'" size="small" variant="tonal" color="purple" @click="pollTaskResult" :loading="polling">Check Status</v-btn>
        </v-card>

        <v-row v-if="images.length" dense>
          <v-col v-for="(img, i) in images" :key="i" cols="6" sm="4"><v-card rounded="xl" class="cursor-pointer" elevation="4" @click="showPreview = true; previewImage = img"><v-img :src="img" cover aspect-ratio="1" /></v-card></v-col>
        </v-row>

        <v-card v-if="!loading && !images.length && !error && !taskId" rounded="xl" class="pa-16 text-center" variant="outlined">
          <v-icon size="64" color="purple-lighten-3" class="mb-4">mdi-image-plus-outline</v-icon>
          <div class="text-h6 text-medium-emphasis mb-1">Ready to Create</div>
          <p class="text-body-2 text-medium-emphasis">Enter a prompt and choose your model</p>
        </v-card>
      </v-col>
    </v-row>

    <v-dialog v-model="showPreview" max-width="90vw"><v-card rounded="xl" color="transparent" flat><v-img :src="previewImage" max-height="85vh" contain rounded="lg" /><div class="d-flex justify-center mt-3"><v-btn icon="mdi-close" variant="flat" @click="showPreview = false" /></div></v-card></v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchModels, imageGen, pollTask } from '../api'
const models=ref([]);const model=ref('');const prompt=ref('');const negativePrompt=ref('')
const size=ref('1024*1024');const n=ref(1);const steps=ref(30);const promptExtend=ref(true)
const refImg=ref('');const loading=ref(false);const taskId=ref('');const taskStatus=ref('')
const polling=ref(false);const images=ref([]);const error=ref('');const showPreview=ref(false);const previewImage=ref('')
const currentModel=computed(()=>models.value.find(m=>m.id===model.value))
const statusColor=computed(()=>{if(taskStatus.value==='SUCCEEDED')return'success';if(taskStatus.value==='FAILED')return'error';return'warning'})
onMounted(async()=>{try{const r=await fetchModels('image');models.value=r.data.models||[];if(models.value.length)model.value=models.value[0].id}catch(e){error.value='Failed to load models'}})
function hasParam(name){return currentModel.value?.parameters?.some(p=>p.name===name)}
function getOptions(name){return currentModel.value?.parameters?.find(p=>p.name===name)?.options||[]}
async function generate(){loading.value=true;error.value='';images.value=[];taskId.value=''
try{const payload={model:model.value,prompt:prompt.value,size:size.value,n:n.value}
if(hasParam('steps'))payload.steps=steps.value
if(hasParam('prompt_extend'))payload.prompt_extend=promptExtend.value
if(negativePrompt.value)payload.negative_prompt=negativePrompt.value
if(refImg.value){payload.ref_img=refImg.value;payload.ref_strength=0.5}
const r=await imageGen(payload);const data=r.data
if(data.task_id){taskId.value=data.task_id;taskStatus.value=data.status;if(data.status==='PENDING'||data.status==='RUNNING')await pollTaskResult()}
else{extractImages(data)}
}catch(e){error.value=e.response?.data?.message||e.message||'Request failed'}
loading.value=false}
async function pollTaskResult(){if(!taskId.value)return;polling.value=true
try{const r=await pollTask(taskId.value);const data=r.data;taskStatus.value=data.output?.task_status||'UNKNOWN';if(taskStatus.value==='SUCCEEDED')extractTaskImages(data);if(taskStatus.value==='FAILED')error.value=data.output?.message||'Task failed'}
catch(e){error.value='Poll failed'};polling.value=false}
function extractImages(data){const result=[];if(data.output?.results)result.push(...data.output.results.map(r=>r.url).filter(Boolean));if(data.output?.choices){for(const c of data.output.choices){const ct=c.message?.content;if(Array.isArray(ct)){for(const it of ct){if(it.image)result.push(it.image)}}}};images.value=result;if(!result.length){error.value='No image URL in response'}}
function extractTaskImages(data){const r=[];if(data.output?.results)r.push(...data.output.results.map(x=>x.url).filter(Boolean));if(data.output?.video_url)r.push(data.output.video_url);images.value=r}
</script>
