<template>
  <div class="pa-6 pa-md-8" style="max-width:1400px;margin:0 auto">
    <div class="d-flex align-center mb-7">
      <v-avatar color="purple" size="48" class="mr-4 elevation-4"><v-icon size="28">mdi-image-outline</v-icon></v-avatar>
      <div><h1 class="text-h4 font-weight-bold">图片工作站</h1><p class="text-body-1 text-medium-emphasis mt-1">模型定义驱动 · 文生图、多图参考、编辑与连续组图</p></div>
    </div>

    <v-row>
      <v-col cols="12" lg="5">
        <v-card rounded="xl" class="pa-6" elevation="2" style="background:linear-gradient(180deg,rgba(147,51,234,.08),transparent 32%)">
          <v-select v-model="model" :items="models" item-title="name" item-value="id" label="模型" variant="outlined" color="purple" class="mb-2" />
          <v-alert v-if="currentModel" color="purple" variant="tonal" density="compact" class="mb-5">{{ currentModel.description }}</v-alert>
          <v-text-field v-model="modelOverride" label="自定义模型 ID（可选）" hint="沿用所选模型的 API 格式，可调用目录中新出现的兼容模型" persistent-hint variant="outlined" class="mb-3" />

          <template v-for="field in fields" :key="field.name">
            <v-textarea v-if="field.type==='string' && (field.name==='prompt' || field.name==='negative_prompt')" v-model="values[field.name]" :label="field.label + (field.required ? ' *' : '')" :hint="field.description" persistent-hint rows="3" variant="outlined" class="mb-2" />
            <v-textarea v-else-if="field.type==='media_list' || field.type==='json'" v-model="values[field.name]" :label="field.label + (field.required ? ' *' : '')" :placeholder="field.placeholder" :hint="field.description" persistent-hint rows="3" variant="outlined" class="mb-2" />
            <v-select v-else-if="field.type==='select'" v-model="values[field.name]" :items="field.options" item-title="label" item-value="value" :label="field.label" :hint="field.description" persistent-hint variant="outlined" class="mb-2" />
            <v-switch v-else-if="field.type==='bool'" v-model="values[field.name]" :label="field.label" :hint="field.description" persistent-hint color="purple" inset class="mb-2" />
            <v-text-field v-else v-model="values[field.name]" :type="field.type==='int'||field.type==='float' ? 'number' : 'text'" :label="field.label + (field.required ? ' *' : '')" :min="field.min" :max="field.max" :step="field.step" :placeholder="field.placeholder" :hint="field.description" persistent-hint variant="outlined" class="mb-2" />
          </template>

          <v-expansion-panels variant="accordion" class="mb-5">
            <v-expansion-panel title="高级 JSON（覆盖/补充任意文档参数）">
              <v-expansion-panel-text>
                <v-textarea v-model="advancedInput" label="input JSON" rows="4" variant="outlined" hint="会与上方 input 字段合并，同名字段以这里为准" persistent-hint class="mb-3" />
                <v-textarea v-model="advancedParameters" label="parameters JSON" rows="5" variant="outlined" hint="文档新增参数可直接写在这里，无需等待页面更新" persistent-hint />
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>

          <v-btn block color="purple" size="x-large" rounded="lg" :loading="loading" :disabled="!canSubmit" @click="generate">
            <v-icon start>mdi-sparkles</v-icon>{{ loading ? '正在提交…' : '生成图片' }}
          </v-btn>
        </v-card>
      </v-col>

      <v-col cols="12" lg="7">
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4" closable @click:close="error=''">{{ error }}</v-alert>
        <v-card v-if="taskId" rounded="xl" class="pa-5 mb-4" variant="outlined">
          <div class="d-flex align-center flex-wrap ga-3 mb-3"><v-chip :color="statusColor" size="small" label>{{ taskStatus }}</v-chip><code class="text-caption">{{ taskId }}</code></div>
          <v-progress-linear v-if="isWorking" indeterminate color="purple" rounded class="mb-3" />
          <div class="text-caption text-medium-emphasis">{{ isWorking ? '后台自动轮询中，页面可保持打开…' : '任务已结束' }}</div>
        </v-card>
        <v-row v-if="images.length" dense>
          <v-col v-for="(img,i) in images" :key="img+i" cols="12" sm="6"><v-card rounded="xl" elevation="4"><v-img :src="img" cover aspect-ratio="1" @click="previewImage=img;showPreview=true" /><v-card-actions><v-btn :href="img" target="_blank" prepend-icon="mdi-download" color="purple" variant="text">打开原图</v-btn></v-card-actions></v-card></v-col>
        </v-row>
        <v-card v-if="!loading&&!images.length&&!error&&!taskId" rounded="xl" class="pa-16 text-center" variant="outlined"><v-icon size="68" color="purple-lighten-3" class="mb-4">mdi-image-multiple-outline</v-icon><div class="text-h6 text-medium-emphasis">选择模型后，所有受支持参数会自动出现</div><p class="text-body-2 text-medium-emphasis mt-2">参考图请使用公网 URL 或百炼 OSS 临时地址</p></v-card>
      </v-col>
    </v-row>

    <v-dialog v-model="showPreview" max-width="92vw"><v-card rounded="xl" color="transparent" flat><v-img :src="previewImage" max-height="88vh" contain /><div class="d-flex justify-center mt-3"><v-btn icon="mdi-close" @click="showPreview=false" /></div></v-card></v-dialog>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { fetchModels, imageGen, pollTask } from '../api'

const models=ref([]),model=ref(''),modelOverride=ref(''),values=reactive({}),advancedInput=ref('{}'),advancedParameters=ref('{}')
const loading=ref(false),taskId=ref(''),taskStatus=ref(''),images=ref([]),error=ref('')
const showPreview=ref(false),previewImage=ref('');let pollTimer=null,pollCount=0
const currentModel=computed(()=>models.value.find(m=>m.id===model.value));const fields=computed(()=>currentModel.value?.parameters||[])
const isWorking=computed(()=>['PENDING','RUNNING','pending','running'].includes(taskStatus.value))
const statusColor=computed(()=>taskStatus.value==='SUCCEEDED'||taskStatus.value==='succeeded'?'success':taskStatus.value==='FAILED'||taskStatus.value==='failed'?'error':'warning')
const canSubmit=computed(()=>!!model.value&&fields.value.filter(f=>f.required).every(f=>hasValue(values[f.name])))

onMounted(async()=>{try{const r=await fetchModels('image');models.value=r.data.models||[];if(models.value.length)model.value=models.value[0].id}catch(e){error.value='模型目录加载失败'}})
onBeforeUnmount(()=>clearTimeout(pollTimer))
watch(model,resetFields)
function resetFields(){for(const k of Object.keys(values))delete values[k];for(const f of fields.value){if(f.default!==undefined&&f.default!==null)values[f.name]=f.default;else values[f.name]=f.type==='bool'?false:''}advancedInput.value='{}';advancedParameters.value='{}'}
function hasValue(v){return v!==undefined&&v!==null&&String(v).trim()!==''}
function parseObject(text,label){try{const v=JSON.parse(text||'{}');if(!v||Array.isArray(v)||typeof v!=='object')throw new Error();return v}catch{throw new Error(`${label} 必须是 JSON 对象`)}}
function parseJSONValue(value,label){try{return JSON.parse(value)}catch{throw new Error(`${label} 不是有效 JSON`)}}
function lines(value){return String(value||'').split(/\r?\n/).map(v=>v.trim()).filter(Boolean)}
function buildPayload(){const input={},parameters={};let imageURLs=[]
  for(const f of fields.value){let v=values[f.name];if(!hasValue(v)&&f.type!=='bool')continue;if(f.type==='int')v=Number.parseInt(v,10);if(f.type==='float')v=Number(v);if(f.type==='json')v=parseJSONValue(v,f.label)
    if(f.name==='images'){imageURLs=lines(v);continue}(f.scope==='input'?input:parameters)[f.name]=v}
  Object.assign(input,parseObject(advancedInput.value,'input JSON'));Object.assign(parameters,parseObject(advancedParameters.value,'parameters JSON'))
  return{model:model.value,model_override:modelOverride.value.trim(),prompt:input.prompt||'',images:imageURLs,input,parameters}}
async function generate(){loading.value=true;error.value='';images.value=[];taskId.value='';taskStatus.value='';clearTimeout(pollTimer);pollCount=0
  try{const {data}=await imageGen(buildPayload());if(data.task_id){taskId.value=data.task_id;taskStatus.value=String(data.status||'PENDING').toUpperCase();schedulePoll()}else extractImages(data)}catch(e){error.value=e.response?.data?.message||e.message||'请求失败'}finally{loading.value=false}}
function schedulePoll(){if(!taskId.value||!isWorking.value)return;clearTimeout(pollTimer);pollTimer=setTimeout(pollTaskResult,pollCount<10?3000:8000)}
async function pollTaskResult(){if(!taskId.value)return;pollCount++;try{const {data}=await pollTask(taskId.value);taskStatus.value=data.output?.task_status||data.status||'UNKNOWN';if(taskStatus.value==='SUCCEEDED'){extractImages(data);return}if(taskStatus.value==='FAILED'){error.value=data.output?.message||data.message||'任务失败';return}schedulePoll()}catch(e){error.value=e.response?.data?.message||'轮询失败';if(pollCount<60)schedulePoll()}}
function extractImages(data){const out=[];if(Array.isArray(data.data))out.push(...data.data.map(x=>x?.url).filter(Boolean));if(Array.isArray(data.output?.results))out.push(...data.output.results.map(x=>x?.url).filter(Boolean));for(const c of data.output?.choices||[])for(const item of c.message?.content||[])if(item.image)out.push(item.image);images.value=[...new Set(out)];if(!images.value.length)error.value='响应成功，但未找到图片 URL；可在历史记录中查看原始响应'}
</script>
