<template>
  <div class="pa-6 pa-md-8" style="max-width:1400px;margin:0 auto">
    <div class="d-flex align-center mb-7"><v-avatar color="teal" size="48" class="mr-4 elevation-4"><v-icon size="28">mdi-video-outline</v-icon></v-avatar><div><h1 class="text-h4 font-weight-bold">视频工作站</h1><p class="text-body-1 text-medium-emphasis mt-1">文生、首尾帧、参考、编辑、续写和有声视频统一入口</p></div></div>
    <v-row>
      <v-col cols="12" lg="5">
        <v-card rounded="xl" class="pa-6" elevation="2" style="background:linear-gradient(180deg,rgba(0,150,136,.09),transparent 32%)">
          <v-select v-model="model" :items="models" item-title="name" item-value="id" label="模型" variant="outlined" color="teal" class="mb-2" />
          <v-alert v-if="currentModel" color="teal" variant="tonal" density="compact" class="mb-5">{{ currentModel.description }}</v-alert>
          <v-text-field v-model="modelOverride" label="自定义模型 ID（可选）" hint="沿用所选模型的 API 格式，可调用目录中新出现的兼容模型" persistent-hint variant="outlined" class="mb-3" />

          <template v-for="field in fields" :key="field.name">
            <v-textarea v-if="field.type==='string' && field.name==='prompt'" v-model="values[field.name]" :label="field.label + (field.required?' *':'')" :hint="field.description" persistent-hint rows="4" variant="outlined" class="mb-2" />
            <v-textarea v-else-if="field.type==='media_list'||field.type==='json'" v-model="values[field.name]" :label="field.label + (field.required?' *':'')" :placeholder="field.placeholder" :hint="field.description" persistent-hint rows="4" variant="outlined" class="mb-2" />
            <v-select v-else-if="field.type==='select'" v-model="values[field.name]" :items="field.options" item-title="label" item-value="value" :label="field.label" :hint="field.description" persistent-hint variant="outlined" class="mb-2" />
            <v-switch v-else-if="field.type==='bool'" v-model="values[field.name]" :label="field.label" :hint="field.description" persistent-hint color="teal" inset class="mb-2" />
            <v-text-field v-else v-model="values[field.name]" :type="field.type==='int'||field.type==='float'?'number':'text'" :label="field.label + (field.required?' *':'')" :min="field.min" :max="field.max" :step="field.step" :hint="field.description" persistent-hint variant="outlined" class="mb-2" />
          </template>

          <v-expansion-panels variant="accordion" class="mb-5"><v-expansion-panel title="高级 JSON（覆盖/补充任意文档参数）"><v-expansion-panel-text>
            <v-textarea v-model="advancedInput" label="input JSON" rows="5" variant="outlined" hint="可直接粘贴文档中的 input，支持任意 media 类型" persistent-hint class="mb-3" />
            <v-textarea v-model="advancedParameters" label="parameters JSON" rows="5" variant="outlined" hint="例如智能时长：{&quot;duration&quot;:-1}；同名字段以这里为准" persistent-hint />
          </v-expansion-panel-text></v-expansion-panel></v-expansion-panels>

          <v-btn block color="teal" size="x-large" rounded="lg" :loading="loading" :disabled="!canSubmit" @click="generate"><v-icon start>mdi-sparkles</v-icon>{{ loading?'正在提交…':'生成视频' }}</v-btn>
        </v-card>
      </v-col>
      <v-col cols="12" lg="7">
        <v-alert v-if="error" type="error" variant="tonal" class="mb-4" closable @click:close="error=''">{{ error }}</v-alert>
        <v-card v-if="taskId" rounded="xl" class="pa-5 mb-4" variant="outlined"><div class="d-flex align-center flex-wrap ga-3 mb-3"><v-chip :color="statusColor" size="small" label>{{ taskStatus }}</v-chip><code class="text-caption">{{ taskId }}</code></div><v-progress-linear v-if="isWorking" indeterminate color="teal" rounded class="mb-3" /><div class="text-caption text-medium-emphasis">{{ isWorking?'视频通常需要数分钟，正在自动轮询…':'任务已结束' }}</div></v-card>
        <v-card v-if="resultUrl" rounded="xl" variant="outlined" overflow-hidden><video :src="resultUrl" controls style="width:100%;max-height:620px;background:#000;display:block" /><v-card-actions><v-btn :href="resultUrl" target="_blank" variant="text" color="teal" prepend-icon="mdi-download">打开原视频</v-btn></v-card-actions></v-card>
        <v-card v-if="!loading&&!resultUrl&&!error&&!taskId" rounded="xl" class="pa-16 text-center" variant="outlined"><v-icon size="68" color="teal-lighten-3" class="mb-4">mdi-video-plus-outline</v-icon><div class="text-h6 text-medium-emphasis">Wan 3.0 可在同一个模型完成所有视频任务</div><p class="text-body-2 text-medium-emphasis mt-2">媒体格式：<code>类型 | URL</code>，每行一个；也可粘贴完整 input JSON</p></v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { fetchModels, pollTask, videoGen } from '../api'

const models=ref([]),model=ref(''),modelOverride=ref(''),values=reactive({}),advancedInput=ref('{}'),advancedParameters=ref('{}')
const loading=ref(false),taskId=ref(''),taskStatus=ref(''),resultUrl=ref(''),error=ref('');let pollTimer=null,pollCount=0
const currentModel=computed(()=>models.value.find(m=>m.id===model.value));const fields=computed(()=>currentModel.value?.parameters||[])
const isWorking=computed(()=>['PENDING','RUNNING','pending','running'].includes(taskStatus.value))
const statusColor=computed(()=>taskStatus.value==='SUCCEEDED'||taskStatus.value==='succeeded'?'success':taskStatus.value==='FAILED'||taskStatus.value==='failed'?'error':'warning')
const canSubmit=computed(()=>{if(!model.value)return false;const requiredOk=fields.value.filter(f=>f.required).every(f=>hasValue(values[f.name]));const promptOrMedia=hasValue(values.prompt)||hasValue(values.media)||hasAdvancedInput();return requiredOk&&promptOrMedia})

onMounted(async()=>{try{const r=await fetchModels('video');models.value=r.data.models||[];if(models.value.length)model.value=models.value[0].id}catch{error.value='模型目录加载失败'}})
onBeforeUnmount(()=>clearTimeout(pollTimer));watch(model,resetFields)
function resetFields(){for(const k of Object.keys(values))delete values[k];for(const f of fields.value){if(f.default!==undefined&&f.default!==null)values[f.name]=f.default;else values[f.name]=f.type==='bool'?false:''}advancedInput.value='{}';advancedParameters.value='{}'}
function hasValue(v){return v!==undefined&&v!==null&&String(v).trim()!==''}
function hasAdvancedInput(){try{return Object.keys(JSON.parse(advancedInput.value||'{}')).length>0}catch{return false}}
function parseObject(text,label){try{const v=JSON.parse(text||'{}');if(!v||Array.isArray(v)||typeof v!=='object')throw new Error();return v}catch{throw new Error(`${label} 必须是 JSON 对象`)}}
function parseJSONValue(value,label){try{return JSON.parse(value)}catch{throw new Error(`${label} 不是有效 JSON`)}}
function parseMedia(value){return String(value||'').split(/\r?\n/).map(line=>line.trim()).filter(Boolean).map((line,index)=>{const parts=line.split('|').map(v=>v.trim());if(parts.length<2||!parts[0]||!parts.slice(1).join('|'))throw new Error(`第 ${index+1} 行媒体格式应为：类型 | URL`);return{type:parts[0],url:parts.slice(1).join('|')}})}
function buildPayload(){const input={},parameters={};let media=[]
  for(const f of fields.value){let v=values[f.name];if(!hasValue(v)&&f.type!=='bool')continue;if(f.type==='int')v=Number.parseInt(v,10);if(f.type==='float')v=Number(v);if(f.type==='json')v=parseJSONValue(v,f.label);if(f.name==='media'){media=parseMedia(v);continue}(f.scope==='input'?input:parameters)[f.name]=v}
  Object.assign(input,parseObject(advancedInput.value,'input JSON'));Object.assign(parameters,parseObject(advancedParameters.value,'parameters JSON'));if(media.length&&!input.media)input.media=media
  return{model:model.value,model_override:modelOverride.value.trim(),prompt:input.prompt||'',media,input,parameters}}
async function generate(){loading.value=true;error.value='';taskId.value='';taskStatus.value='';resultUrl.value='';clearTimeout(pollTimer);pollCount=0
  try{const {data}=await videoGen(buildPayload());taskId.value=data.task_id||data.output?.task_id||'';taskStatus.value=String(data.status||data.output?.task_status||'PENDING').toUpperCase();if(!taskId.value){extractVideo(data);if(!resultUrl.value)throw new Error('响应中没有 task_id 或视频 URL')}else schedulePoll()}catch(e){error.value=e.response?.data?.message||e.message||'请求失败'}finally{loading.value=false}}
function schedulePoll(){if(!taskId.value||!isWorking.value)return;clearTimeout(pollTimer);pollTimer=setTimeout(pollTaskResult,pollCount<6?5000:12000)}
async function pollTaskResult(){pollCount++;try{const {data}=await pollTask(taskId.value);taskStatus.value=data.output?.task_status||data.status||'UNKNOWN';if(taskStatus.value==='SUCCEEDED'||taskStatus.value==='succeeded'){extractVideo(data);return}if(taskStatus.value==='FAILED'||taskStatus.value==='failed'){error.value=data.output?.message||data.message||'任务失败';return}schedulePoll()}catch(e){error.value=e.response?.data?.message||'轮询失败';if(pollCount<80)schedulePoll()}}
function extractVideo(data){resultUrl.value=data.output?.video_url||data.output?.results?.[0]?.url||data.data?.[0]?.url||data.url||'';if(!resultUrl.value)error.value='任务成功，但未找到视频 URL；可在历史记录中查看原始响应'}
</script>
