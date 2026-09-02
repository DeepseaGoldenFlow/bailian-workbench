<template>
  <div class="page-shell studio-page video-studio">
    <header class="page-header">
      <div><div class="page-eyebrow">视频创作</div><h1 class="page-title">让画面开始流动</h1><p class="page-subtitle">从文字、图片或视频参考出发，完成镜头生成、首尾帧控制、编辑、续写与声音合成。</p></div>
      <div class="header-chips"><v-chip variant="tonal" color="teal" prepend-icon="mdi-timer-outline">最长 30 秒</v-chip><v-chip variant="tonal" color="primary" prepend-icon="mdi-volume-high">支持音频</v-chip></div>
    </header>

    <div class="studio-grid">
      <v-card rounded="xl" class="surface-card editor-panel">
        <div class="panel-section"><div class="section-number">01</div><div class="section-copy"><strong>选择视频模型</strong><span>根据素材类型和目标时长选择</span></div></div>
        <v-autocomplete v-model="model" :items="models" item-title="name" item-value="id" label="视频模型" variant="outlined" hide-details class="mt-4" />
        <div v-if="currentModel" class="model-note"><v-icon size="18">mdi-information-outline</v-icon><span>{{ currentModel.description }}</span></div>

        <v-divider class="my-6" />
        <div class="panel-section"><div class="section-number">02</div><div class="section-copy"><strong>描述镜头与素材</strong><span>写清动作、运镜、节奏、对白和声音</span></div></div>
        <div class="fields mt-5"><DynamicField v-for="field in primaryFields" :key="field.name" v-model="values[field.name]" :field="field" /></div>
        <div class="media-help"><v-icon size="17">mdi-link-variant</v-icon><div><strong>媒体填写格式</strong><span><code>first_frame | 图片地址</code>，每行一个；也支持 last_frame、reference_image、reference_video、audio、file、link。</span></div></div>

        <v-divider class="my-6" />
        <div class="panel-section"><div class="section-number">03</div><div class="section-copy"><strong>画面与声音设置</strong><span>控制分辨率、比例、时长与音频</span></div></div>
        <div class="fields field-grid mt-5"><DynamicField v-for="field in settingFields" :key="field.name" v-model="values[field.name]" :field="field" /></div>

        <v-expansion-panels variant="accordion" class="advanced-panel mt-3"><v-expansion-panel rounded="lg"><v-expansion-panel-title><div class="d-flex align-center ga-2"><v-icon size="18">mdi-code-json</v-icon><span>高级参数与自定义模型</span></div></v-expansion-panel-title><v-expansion-panel-text>
          <v-text-field v-model="modelOverride" label="自定义模型 ID" hint="沿用当前模型的请求格式" persistent-hint variant="outlined" class="mb-3" />
          <v-textarea v-model="advancedInput" label="input JSON" rows="5" variant="outlined" hint="可直接粘贴文档中的完整 input" persistent-hint class="mb-3" />
          <v-textarea v-model="advancedParameters" label="parameters JSON" rows="5" variant="outlined" hint="例如智能时长：{&quot;duration&quot;:-1}" persistent-hint />
        </v-expansion-panel-text></v-expansion-panel></v-expansion-panels>
        <v-btn block color="teal" size="x-large" rounded="xl" class="generate-btn mt-6" :loading="loading" :disabled="!canSubmit" prepend-icon="mdi-movie-open-play-outline" @click="generate">生成视频</v-btn>
      </v-card>

      <div class="preview-column">
        <v-alert v-if="error" type="error" variant="tonal" rounded="xl" closable class="mb-4" @click:close="error=''">{{ error }}</v-alert>
        <v-card v-if="taskId" rounded="xl" class="surface-card task-card mb-4"><div class="d-flex align-center ga-3"><div class="task-icon"><v-icon>mdi-movie-roll</v-icon></div><div class="flex-1-1"><strong>{{ isWorking ? '正在制作视频' : '任务已结束' }}</strong><div class="mono text-caption text-medium-emphasis mt-1">{{ taskId }}</div></div><v-chip :color="statusColor" size="small" variant="tonal">{{ statusLabel }}</v-chip></div><v-progress-linear v-if="isWorking" indeterminate color="teal" rounded class="mt-4" /><p v-if="isWorking" class="task-tip">视频生成通常需要几分钟，可以保持页面打开。</p></v-card>
        <v-card v-if="resultUrl" rounded="xl" class="surface-card video-result"><video :src="resultUrl" controls /><div class="video-result__bar"><div><strong>视频生成完成</strong><span>点击播放或打开原视频</span></div><v-btn :href="resultUrl" target="_blank" variant="tonal" color="teal" rounded="xl" prepend-icon="mdi-open-in-new">打开视频</v-btn></div></v-card>
        <v-card v-else rounded="xl" class="surface-card preview-empty">
          <div class="film-stage"><div class="film-frame film-frame--back" /><div class="film-frame"><div class="play-mark"><v-icon size="34">mdi-play</v-icon></div><div class="timeline"><i /><i /><i /><i /><i /></div></div></div>
          <h2>从一个镜头开始你的故事</h2><p>输入创意描述，也可以加入首帧、尾帧、参考图、参考视频或音频。</p>
          <div class="workflow"><span>写下创意</span><v-icon size="15">mdi-arrow-right</v-icon><span>添加素材</span><v-icon size="15">mdi-arrow-right</v-icon><span>生成视频</span></div>
        </v-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed,onBeforeUnmount,onMounted,reactive,ref,watch } from 'vue'
import DynamicField from '../components/DynamicField.vue'
import { fetchModels,pollTask,videoGen } from '../api'
const models=ref([]),model=ref(''),modelOverride=ref(''),values=reactive({}),advancedInput=ref('{}'),advancedParameters=ref('{}')
const loading=ref(false),taskId=ref(''),taskStatus=ref(''),resultUrl=ref(''),error=ref('');let pollTimer=null,pollCount=0
const currentModel=computed(()=>models.value.find(item=>item.id===model.value)),fields=computed(()=>currentModel.value?.parameters||[])
const primaryFields=computed(()=>fields.value.filter(field=>field.scope==='input')),settingFields=computed(()=>fields.value.filter(field=>field.scope!=='input'))
const normalizedStatus=computed(()=>String(taskStatus.value).toUpperCase()),isWorking=computed(()=>['PENDING','RUNNING'].includes(normalizedStatus.value)),statusColor=computed(()=>normalizedStatus.value==='SUCCEEDED'?'success':normalizedStatus.value==='FAILED'?'error':'warning'),statusLabel=computed(()=>normalizedStatus.value==='SUCCEEDED'?'已完成':normalizedStatus.value==='FAILED'?'失败':'处理中')
const canSubmit=computed(()=>{if(!model.value)return false;const requiredOK=fields.value.filter(field=>field.required).every(field=>hasValue(values[field.name]));return requiredOK&&(hasValue(values.prompt)||hasValue(values.media)||hasAdvancedInput())})
onMounted(async()=>{try{const response=await fetchModels('video');models.value=response.data.models||[];if(models.value.length)model.value=models.value[0].id}catch{error.value='模型目录加载失败，请稍后重试'}})
onBeforeUnmount(()=>clearTimeout(pollTimer));watch(model,resetFields)
function resetFields(){for(const key of Object.keys(values))delete values[key];for(const field of fields.value)values[field.name]=field.default!==undefined&&field.default!==null?field.default:field.type==='bool'?false:'';advancedInput.value='{}';advancedParameters.value='{}'}
function hasValue(value){return value!==undefined&&value!==null&&String(value).trim()!==''} function hasAdvancedInput(){try{return Object.keys(JSON.parse(advancedInput.value||'{}')).length>0}catch{return false}}
function parseObject(text,label){try{const value=JSON.parse(text||'{}');if(!value||Array.isArray(value)||typeof value!=='object')throw new Error();return value}catch{throw new Error(`${label} 必须是 JSON 对象`)}} function parseJSONValue(value,label){try{return JSON.parse(value)}catch{throw new Error(`${label} 不是有效 JSON`)}}
function parseMedia(value){return String(value||'').split(/\r?\n/).map(line=>line.trim()).filter(Boolean).map((line,index)=>{const parts=line.split('|').map(part=>part.trim());if(parts.length<2||!parts[0]||!parts.slice(1).join('|'))throw new Error(`第 ${index+1} 行媒体格式应为：类型 | URL`);return{type:parts[0],url:parts.slice(1).join('|')}})}
function buildPayload(){const input={},parameters={};let media=[];for(const field of fields.value){let value=values[field.name];if(!hasValue(value)&&field.type!=='bool')continue;if(field.type==='int')value=Number.parseInt(value,10);if(field.type==='float')value=Number(value);if(field.type==='json')value=parseJSONValue(value,field.label);if(field.name==='media'){media=parseMedia(value);continue}(field.scope==='input'?input:parameters)[field.name]=value}Object.assign(input,parseObject(advancedInput.value,'input JSON'));Object.assign(parameters,parseObject(advancedParameters.value,'parameters JSON'));if(media.length&&!input.media)input.media=media;return{model:model.value,model_override:modelOverride.value.trim(),prompt:input.prompt||'',media,input,parameters}}
async function generate(){loading.value=true;error.value='';taskId.value='';taskStatus.value='';resultUrl.value='';clearTimeout(pollTimer);pollCount=0;try{const{data}=await videoGen(buildPayload());taskId.value=data.task_id||data.output?.task_id||'';taskStatus.value=data.status||data.output?.task_status||'PENDING';if(!taskId.value){extractVideo(data);if(!resultUrl.value)throw new Error('响应中没有任务编号或视频地址')}else schedulePoll()}catch(e){error.value=e.response?.data?.message||e.message||'视频生成请求失败'}finally{loading.value=false}}
function schedulePoll(){if(!taskId.value||!isWorking.value)return;clearTimeout(pollTimer);pollTimer=setTimeout(pollTaskResult,pollCount<6?5000:12000)}
async function pollTaskResult(){pollCount++;try{const{data}=await pollTask(taskId.value);taskStatus.value=data.output?.task_status||data.status||'UNKNOWN';if(normalizedStatus.value==='SUCCEEDED'){extractVideo(data);return}if(normalizedStatus.value==='FAILED'){error.value=data.output?.message||data.message||'视频生成失败';return}schedulePoll()}catch(e){error.value=e.response?.data?.message||'查询任务状态失败';if(pollCount<80)schedulePoll()}}
function extractVideo(data){resultUrl.value=data.output?.video_url||data.output?.results?.[0]?.url||data.data?.[0]?.url||data.url||'';if(!resultUrl.value)error.value='任务成功，但响应中没有找到视频地址'}
</script>

<style scoped>
.studio-grid{display:grid;grid-template-columns:minmax(380px,480px) 1fr;gap:24px;align-items:start}.header-chips{display:flex;gap:8px}.editor-panel{padding:26px}.panel-section{display:flex;align-items:center;gap:12px}.section-number{width:34px;height:34px;flex:0 0 34px;display:grid;place-items:center;border-radius:11px;color:#0d9488;background:rgba(13,148,136,.1);font-size:10px;font-weight:800}.section-copy{display:flex;flex-direction:column}.section-copy strong{font-size:15px}.section-copy span{margin-top:3px;color:rgb(var(--v-theme-on-surface-variant));font-size:10px}.model-note{display:flex;align-items:flex-start;gap:8px;margin-top:11px;padding:11px 12px;border-radius:12px;color:#0f766e;background:rgba(13,148,136,.08);font-size:11px;line-height:1.55}.media-help{display:flex;align-items:flex-start;gap:10px;margin-top:-1px;padding:12px;border-radius:13px;color:rgb(var(--v-theme-on-surface-variant));background:rgb(var(--v-theme-surface-variant));font-size:10px}.media-help>div{display:flex;flex-direction:column}.media-help strong{color:rgb(var(--v-theme-on-surface));font-size:11px}.media-help span{margin-top:4px;line-height:1.55}.field-grid{display:grid;grid-template-columns:1fr 1fr;gap:0 10px}.advanced-panel :deep(.v-expansion-panel){border:1px solid rgba(var(--v-border-color),.1);background:rgba(var(--v-theme-surface-variant),.45)}.generate-btn{box-shadow:0 12px 28px rgba(13,148,136,.22)}.preview-column{position:sticky;top:24px}.preview-empty{min-height:650px;display:flex;flex-direction:column;align-items:center;justify-content:center;padding:45px;text-align:center}.preview-empty h2{margin:32px 0 8px;font-size:20px}.preview-empty>p{max-width:450px;color:rgb(var(--v-theme-on-surface-variant));font-size:12px;line-height:1.7}.film-stage{position:relative;width:320px;height:230px;display:grid;place-items:center}.film-frame{position:relative;width:245px;height:160px;display:flex;flex-direction:column;align-items:center;justify-content:center;overflow:hidden;border:1px solid rgba(13,148,136,.25);border-radius:23px;background:radial-gradient(circle at 70% 20%,rgba(13,148,136,.26),transparent 40%),linear-gradient(145deg,#162b35,#0f172a);box-shadow:0 24px 60px rgba(15,23,42,.25)}.film-frame--back{position:absolute;transform:rotate(-7deg) translate(-18px,3px);opacity:.3}.play-mark{width:58px;height:58px;display:grid;place-items:center;border-radius:50%;color:white;background:rgba(45,212,191,.24);backdrop-filter:blur(5px)}.timeline{position:absolute;left:18px;right:18px;bottom:17px;display:grid;grid-template-columns:repeat(5,1fr);gap:5px}.timeline i{height:4px;border-radius:4px;background:rgba(255,255,255,.18)}.timeline i:first-child{background:#2dd4bf}.workflow{display:flex;align-items:center;gap:9px;margin-top:27px;color:rgb(var(--v-theme-on-surface-variant));font-size:10px}.workflow span{padding:7px 9px;border-radius:9px;background:rgb(var(--v-theme-surface-variant))}.task-card{padding:18px}.task-icon{width:42px;height:42px;display:grid;place-items:center;border-radius:14px;color:#0d9488;background:rgba(13,148,136,.1)}.task-tip{margin:10px 0 0;color:rgb(var(--v-theme-on-surface-variant));font-size:10px}.video-result{overflow:hidden}.video-result video{width:100%;max-height:620px;display:block;background:#000}.video-result__bar{display:flex;align-items:center;justify-content:space-between;gap:15px;padding:16px}.video-result__bar>div{display:flex;flex-direction:column}.video-result__bar strong{font-size:13px}.video-result__bar span{margin-top:3px;color:rgb(var(--v-theme-on-surface-variant));font-size:10px}@media(max-width:1100px){.studio-grid{grid-template-columns:1fr}.preview-column{position:static}.preview-empty{min-height:460px}}@media(max-width:600px){.editor-panel{padding:19px 16px}.field-grid{grid-template-columns:1fr}.header-chips{display:none}.preview-empty{min-height:400px;padding:25px}.film-stage{width:270px}.video-result__bar{align-items:flex-start;flex-direction:column}}
</style>
