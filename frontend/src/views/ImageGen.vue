<template>
  <div class="page-shell studio-page">
    <header class="page-header">
      <div><div class="page-eyebrow">图片创作</div><h1 class="page-title">把脑海里的画面变成现实</h1><p class="page-subtitle">覆盖文生图、参考图编辑、多图融合与连续组图；模型支持的参数会自动呈现。</p></div>
      <div class="header-chips"><v-chip variant="tonal" color="purple" prepend-icon="mdi-image-size-select-large">最高 4K</v-chip><v-chip variant="tonal" color="primary" prepend-icon="mdi-tune-variant">完整参数</v-chip></div>
    </header>

    <div class="studio-grid">
      <v-card rounded="xl" class="surface-card editor-panel">
        <div class="panel-section">
          <div class="section-number">01</div><div class="section-copy"><strong>选择模型</strong><span>不同模型擅长不同的画面风格和编辑任务</span></div>
        </div>
        <v-autocomplete v-model="model" :items="models" item-title="name" item-value="id" label="图片模型" variant="outlined" hide-details class="mt-4" />
        <div v-if="currentModel" class="model-note"><v-icon size="18">mdi-information-outline</v-icon><span>{{ currentModel.description }}</span></div>

        <v-divider class="my-6" />
        <div class="panel-section"><div class="section-number">02</div><div class="section-copy"><strong>描述画面</strong><span>越具体的描述，越容易得到稳定结果</span></div></div>
        <div class="fields mt-5"><DynamicField v-for="field in primaryFields" :key="field.name" v-model="values[field.name]" :field="field" /></div>

        <v-divider class="my-6" />
        <div class="panel-section"><div class="section-number">03</div><div class="section-copy"><strong>生成设置</strong><span>调整尺寸、数量、随机种子与模型选项</span></div></div>
        <div class="fields field-grid mt-5"><DynamicField v-for="field in settingFields" :key="field.name" v-model="values[field.name]" :field="field" /></div>

        <v-expansion-panels variant="accordion" class="advanced-panel mt-3">
          <v-expansion-panel rounded="lg">
            <v-expansion-panel-title><div class="d-flex align-center ga-2"><v-icon size="18">mdi-code-json</v-icon><span>高级参数与自定义模型</span></div></v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-text-field v-model="modelOverride" label="自定义模型 ID" hint="沿用当前模型的请求格式，调用目录中新出现的兼容模型" persistent-hint variant="outlined" class="mb-3" />
              <v-textarea v-model="advancedInput" label="input JSON" rows="4" variant="outlined" hint="覆盖或补充任意 input 字段" persistent-hint class="mb-3" />
              <v-textarea v-model="advancedParameters" label="parameters JSON" rows="5" variant="outlined" hint="覆盖或补充任意 parameters 字段" persistent-hint />
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
        <v-btn block color="primary" size="x-large" rounded="xl" class="generate-btn mt-6" :loading="loading" :disabled="!canSubmit" prepend-icon="mdi-creation-outline" @click="generate">生成图片</v-btn>
      </v-card>

      <div class="preview-column">
        <v-alert v-if="error" type="error" variant="tonal" rounded="xl" closable class="mb-4" @click:close="error=''">{{ error }}</v-alert>
        <v-card v-if="taskId" rounded="xl" class="surface-card task-card mb-4"><div class="d-flex align-center ga-3"><div class="task-icon"><v-icon>mdi-progress-wrench</v-icon></div><div class="flex-1-1"><strong>{{ isWorking ? '正在创作图片' : '任务已结束' }}</strong><div class="mono text-caption text-medium-emphasis mt-1">{{ taskId }}</div></div><v-chip :color="statusColor" size="small" variant="tonal">{{ statusLabel }}</v-chip></div><v-progress-linear v-if="isWorking" indeterminate color="primary" rounded class="mt-4" /></v-card>
        <div v-if="images.length" class="image-results"><v-card v-for="(image,index) in images" :key="image+index" rounded="xl" class="result-image surface-card"><v-img :src="image" cover aspect-ratio="1" @click="previewImage=image;showPreview=true" /><div class="result-image__bar"><span>作品 {{ index+1 }}</span><v-btn :href="image" target="_blank" icon="mdi-open-in-new" variant="text" size="small" /></div></v-card></div>
        <v-card v-else rounded="xl" class="surface-card preview-empty">
          <div class="preview-art"><div class="preview-art__glow" /><div class="preview-frame"><v-icon size="64">mdi-image-plus-outline</v-icon><span>AI IMAGE</span></div></div>
          <h2>你的作品会出现在这里</h2><p>选择模型、描述画面并调整参数，然后点击“生成图片”。</p>
          <div class="preview-tips"><span><v-icon size="15">mdi-lightbulb-outline</v-icon>描述主体与环境</span><span><v-icon size="15">mdi-camera-outline</v-icon>补充光线与镜头</span><span><v-icon size="15">mdi-palette-outline</v-icon>指定风格和色彩</span></div>
        </v-card>
      </div>
    </div>
    <v-dialog v-model="showPreview" max-width="92vw"><v-card rounded="xl" color="transparent" flat><v-img :src="previewImage" max-height="88vh" contain /><v-btn icon="mdi-close" class="dialog-close" @click="showPreview=false" /></v-card></v-dialog>
  </div>
</template>

<script setup>
import { computed,onBeforeUnmount,onMounted,reactive,ref,watch } from 'vue'
import DynamicField from '../components/DynamicField.vue'
import { fetchModels,imageGen,pollTask } from '../api'
const models=ref([]),model=ref(''),modelOverride=ref(''),values=reactive({}),advancedInput=ref('{}'),advancedParameters=ref('{}')
const loading=ref(false),taskId=ref(''),taskStatus=ref(''),images=ref([]),error=ref(''),showPreview=ref(false),previewImage=ref('');let pollTimer=null,pollCount=0
const currentModel=computed(()=>models.value.find(item=>item.id===model.value)),fields=computed(()=>currentModel.value?.parameters||[])
const primaryFields=computed(()=>fields.value.filter(field=>['prompt','negative_prompt','images','ref_img'].includes(field.name)||field.scope==='input'&&!['prompt','negative_prompt','images','ref_img'].includes(field.name)))
const settingFields=computed(()=>fields.value.filter(field=>!primaryFields.value.includes(field)))
const normalizedStatus=computed(()=>String(taskStatus.value).toUpperCase()),isWorking=computed(()=>['PENDING','RUNNING'].includes(normalizedStatus.value)),statusColor=computed(()=>normalizedStatus.value==='SUCCEEDED'?'success':normalizedStatus.value==='FAILED'?'error':'warning'),statusLabel=computed(()=>normalizedStatus.value==='SUCCEEDED'?'已完成':normalizedStatus.value==='FAILED'?'失败':'处理中')
const canSubmit=computed(()=>!!model.value&&fields.value.filter(field=>field.required).every(field=>hasValue(values[field.name])))
onMounted(async()=>{try{const response=await fetchModels('image');models.value=response.data.models||[];if(models.value.length)model.value=models.value[0].id}catch{error.value='模型目录加载失败，请稍后重试'}})
onBeforeUnmount(()=>clearTimeout(pollTimer));watch(model,resetFields)
function resetFields(){for(const key of Object.keys(values))delete values[key];for(const field of fields.value)values[field.name]=field.default!==undefined&&field.default!==null?field.default:field.type==='bool'?false:'';advancedInput.value='{}';advancedParameters.value='{}'}
function hasValue(value){return value!==undefined&&value!==null&&String(value).trim()!==''}
function parseObject(text,label){try{const value=JSON.parse(text||'{}');if(!value||Array.isArray(value)||typeof value!=='object')throw new Error();return value}catch{throw new Error(`${label} 必须是 JSON 对象`)}}
function parseJSONValue(value,label){try{return JSON.parse(value)}catch{throw new Error(`${label} 不是有效 JSON`)}}
function lines(value){return String(value||'').split(/\r?\n/).map(item=>item.trim()).filter(Boolean)}
function buildPayload(){const input={},parameters={};let imageURLs=[];for(const field of fields.value){let value=values[field.name];if(!hasValue(value)&&field.type!=='bool')continue;if(field.type==='int')value=Number.parseInt(value,10);if(field.type==='float')value=Number(value);if(field.type==='json')value=parseJSONValue(value,field.label);if(field.name==='images'){imageURLs=lines(value);continue}(field.scope==='input'?input:parameters)[field.name]=value}Object.assign(input,parseObject(advancedInput.value,'input JSON'));Object.assign(parameters,parseObject(advancedParameters.value,'parameters JSON'));return{model:model.value,model_override:modelOverride.value.trim(),prompt:input.prompt||'',images:imageURLs,input,parameters}}
async function generate(){loading.value=true;error.value='';images.value=[];taskId.value='';taskStatus.value='';clearTimeout(pollTimer);pollCount=0;try{const{data}=await imageGen(buildPayload());const newTask=data.task_id||data.output?.task_id;if(newTask){taskId.value=newTask;taskStatus.value=data.status||data.output?.task_status||'PENDING';schedulePoll()}else extractImages(data)}catch(e){error.value=e.response?.data?.message||e.message||'图片生成请求失败'}finally{loading.value=false}}
function schedulePoll(){if(!taskId.value||!isWorking.value)return;clearTimeout(pollTimer);pollTimer=setTimeout(pollTaskResult,pollCount<10?3000:8000)}
async function pollTaskResult(){if(!taskId.value)return;pollCount++;try{const{data}=await pollTask(taskId.value);taskStatus.value=data.output?.task_status||data.status||'UNKNOWN';if(normalizedStatus.value==='SUCCEEDED'){extractImages(data);return}if(normalizedStatus.value==='FAILED'){error.value=data.output?.message||data.message||'图片生成失败';return}schedulePoll()}catch(e){error.value=e.response?.data?.message||'查询任务状态失败';if(pollCount<60)schedulePoll()}}
function extractImages(data){const urls=[];if(Array.isArray(data.data))urls.push(...data.data.map(item=>item?.url).filter(Boolean));if(Array.isArray(data.output?.results))urls.push(...data.output.results.map(item=>item?.url).filter(Boolean));for(const choice of data.output?.choices||[])for(const item of choice.message?.content||[])if(item.image)urls.push(item.image);images.value=[...new Set(urls)];if(!images.value.length)error.value='任务成功，但响应中没有找到图片地址'}
</script>

<style scoped>
.studio-grid{display:grid;grid-template-columns:minmax(380px,480px) 1fr;gap:24px;align-items:start}.header-chips{display:flex;gap:8px}.editor-panel{padding:26px}.panel-section{display:flex;align-items:center;gap:12px}.section-number{width:34px;height:34px;flex:0 0 34px;display:grid;place-items:center;border-radius:11px;color:rgb(var(--v-theme-primary));background:rgba(var(--v-theme-primary),.1);font-size:10px;font-weight:800}.section-copy{display:flex;flex-direction:column}.section-copy strong{font-size:15px}.section-copy span{margin-top:3px;color:rgb(var(--v-theme-on-surface-variant));font-size:10px}.model-note{display:flex;align-items:flex-start;gap:8px;margin-top:11px;padding:11px 12px;border-radius:12px;color:rgb(var(--v-theme-primary));background:rgba(var(--v-theme-primary),.07);font-size:11px;line-height:1.55}.field-grid{display:grid;grid-template-columns:1fr 1fr;gap:0 10px}.field-grid :deep(.dynamic-field:has(textarea)){grid-column:1/-1}.advanced-panel :deep(.v-expansion-panel){border:1px solid rgba(var(--v-border-color),.1);background:rgba(var(--v-theme-surface-variant),.45)}.generate-btn{box-shadow:0 12px 28px rgba(var(--v-theme-primary),.24)}.preview-column{position:sticky;top:24px}.preview-empty{min-height:650px;display:flex;flex-direction:column;align-items:center;justify-content:center;padding:45px;text-align:center}.preview-empty h2{margin:28px 0 8px;font-size:20px}.preview-empty>p{max-width:440px;color:rgb(var(--v-theme-on-surface-variant));font-size:12px;line-height:1.7}.preview-art{position:relative;width:280px;height:220px;display:grid;place-items:center}.preview-art__glow{position:absolute;width:210px;height:210px;border-radius:50%;background:rgba(124,58,237,.14);filter:blur(22px)}.preview-frame{position:relative;width:200px;height:170px;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:18px;border:1px solid rgba(var(--v-theme-primary),.2);border-radius:28px;color:rgb(var(--v-theme-primary));background:linear-gradient(145deg,rgba(var(--v-theme-primary),.10),rgba(var(--v-theme-secondary),.06));transform:rotate(-3deg);box-shadow:18px 18px 0 rgba(var(--v-theme-primary),.055)}.preview-frame span{font-size:9px;letter-spacing:.3em}.preview-tips{display:flex;flex-wrap:wrap;justify-content:center;gap:8px;margin-top:25px}.preview-tips span{display:flex;align-items:center;gap:6px;padding:8px 10px;border-radius:10px;color:rgb(var(--v-theme-on-surface-variant));background:rgb(var(--v-theme-surface-variant));font-size:10px}.task-card{padding:18px}.task-icon{width:42px;height:42px;display:grid;place-items:center;border-radius:14px;color:rgb(var(--v-theme-primary));background:rgba(var(--v-theme-primary),.1)}.image-results{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:15px}.result-image{overflow:hidden}.result-image :deep(.v-img){cursor:zoom-in}.result-image__bar{display:flex;align-items:center;justify-content:space-between;padding:8px 10px 8px 14px;font-size:11px}.dialog-close{position:absolute;right:10px;top:10px}@media(max-width:1100px){.studio-grid{grid-template-columns:1fr}.preview-column{position:static}.preview-empty{min-height:460px}}@media(max-width:600px){.editor-panel{padding:19px 16px}.field-grid{grid-template-columns:1fr}.header-chips{display:none}.image-results{grid-template-columns:1fr}.preview-empty{min-height:400px;padding:25px}}
</style>
