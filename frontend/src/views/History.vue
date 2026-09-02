<template>
  <div class="page-shell history-page">
    <header class="page-header">
      <div><div class="page-eyebrow">创作记录</div><h1 class="page-title">每一次灵感，都有迹可循</h1><p class="page-subtitle">集中查看图片、视频、对话、语音与工具处理记录。</p></div>
      <v-btn icon="mdi-refresh" variant="tonal" rounded="xl" :loading="loading" aria-label="刷新记录" @click="loadHistory" />
    </header>
    <div class="filter-bar">
      <button v-for="filter in filters" :key="filter.value" :class="{active:activeFilter===filter.value}" @click="activeFilter=filter.value;loadHistory()"><v-icon size="17">{{ filter.icon }}</v-icon>{{ filter.label }}</button>
    </div>
    <v-progress-linear v-if="loading" indeterminate color="primary" rounded class="mb-4" />
    <v-alert v-if="error" type="error" variant="tonal" rounded="xl" class="mb-4" closable @click:close="error=''">{{ error }}</v-alert>

    <div v-if="!loading&&!error&&!entries.length" class="surface-card empty-state rounded-xl"><div><div class="empty-state__icon"><v-icon size="34">mdi-history</v-icon></div><h2 class="section-title">还没有创作记录</h2><p class="section-caption mt-2 mb-5">完成一次对话、图片或视频生成后，记录会出现在这里。</p><v-btn to="/image" color="primary" rounded="xl">开始第一次创作</v-btn></div></div>
    <div v-else class="history-list">
      <v-card v-for="entry in entries" :key="`${entry.type}-${entry.id}`" rounded="xl" class="surface-card history-item">
        <div class="history-item__top" @click="toggleEntry(entry)">
          <div class="history-icon" :style="typeStyle(entry.type)"><v-icon size="21">{{ typeIcon(entry.type) }}</v-icon></div>
          <div class="history-copy"><div class="d-flex align-center ga-2"><strong>{{ entry.prompt || entry.content || '未命名创作' }}</strong><v-chip v-if="entry.status" size="x-small" variant="tonal" :color="statusColor(entry.status)">{{ statusLabel(entry.status) }}</v-chip></div><span>{{ typeLabel(entry.type) }} · {{ entry.model || '默认模型' }} · {{ formatTime(entry.created_at) }}</span></div>
          <v-btn icon="mdi-delete-outline" variant="text" size="small" color="error" aria-label="删除记录" @click.stop="deleteEntry(entry)" />
          <v-icon size="19" class="expand-icon" :class="{rotated:isExpanded(entry)}">mdi-chevron-down</v-icon>
        </div>
        <v-expand-transition>
          <div v-if="isExpanded(entry)" class="history-detail">
            <div v-if="entry.type==='image'&&parsedImages(entry.result).length" class="media-grid"><button v-for="(url,index) in parsedImages(entry.result)" :key="url+index" @click="previewUrl=url;showPreview=true"><v-img :src="url" cover height="180" /></button></div>
            <video v-if="entry.type==='video'&&parsedVideo(entry.result)" :src="parsedVideo(entry.result)" controls class="video-preview" />
            <div v-if="['chat','translate','ocr','document'].includes(entry.type)" class="content-preview">{{ entry.content || entry.prompt }}</div>
            <details v-if="entry.result" class="raw-result"><summary>查看原始响应</summary><pre>{{ prettyResult(entry.result) }}</pre></details>
          </div>
        </v-expand-transition>
      </v-card>
    </div>
    <v-dialog v-model="showPreview" max-width="92vw"><v-card rounded="xl" color="transparent" flat><v-img :src="previewUrl" max-height="88vh" contain /><v-btn icon="mdi-close" class="preview-close" @click="showPreview=false" /></v-card></v-dialog>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '../api'
const entries=ref([]),loading=ref(false),error=ref(''),activeFilter=ref(''),expanded=ref(new Set()),previewUrl=ref(''),showPreview=ref(false)
const filters=[{label:'全部',value:'',icon:'mdi-view-grid-outline'},{label:'对话',value:'chat',icon:'mdi-message-processing-outline'},{label:'图片',value:'image',icon:'mdi-image-outline'},{label:'视频',value:'video',icon:'mdi-video-outline'},{label:'声音',value:'audio',icon:'mdi-waveform'},{label:'翻译',value:'translate',icon:'mdi-translate'},{label:'文字识别',value:'ocr',icon:'mdi-text-recognition'}]
const types={chat:['智能对话','mdi-message-processing-outline','#2563eb','rgba(37,99,235,.1)'],image:['图片创作','mdi-image-outline','#9333ea','rgba(147,51,234,.1)'],video:['视频创作','mdi-video-outline','#0d9488','rgba(13,148,136,.1)'],audio:['声音合成','mdi-waveform','#ea580c','rgba(234,88,12,.1)'],translate:['文本翻译','mdi-translate','#2563eb','rgba(37,99,235,.1)'],ocr:['文字识别','mdi-text-recognition','#db2777','rgba(219,39,119,.1)'],document:['文档分析','mdi-file-document-outline','#7c3aed','rgba(124,58,237,.1)']}
function typeLabel(type){return types[type]?.[0]||'其他'} function typeIcon(type){return types[type]?.[1]||'mdi-creation-outline'} function typeStyle(type){return{color:types[type]?.[2]||'#64748b',background:types[type]?.[3]||'rgba(100,116,139,.1)'}}
function statusLabel(status){return['SUCCEEDED','completed','succeeded'].includes(status)?'已完成':['FAILED','failed'].includes(status)?'失败':'处理中'} function statusColor(status){return['SUCCEEDED','completed','succeeded'].includes(status)?'success':['FAILED','failed'].includes(status)?'error':'warning'}
function formatTime(value){if(!value)return'未知时间';const date=new Date(value);return Number.isNaN(date.getTime())?value:date.toLocaleString('zh-CN',{month:'2-digit',day:'2-digit',hour:'2-digit',minute:'2-digit'})}
function key(entry){return`${entry.type}-${entry.id}`} function isExpanded(entry){return expanded.value.has(key(entry))} function toggleEntry(entry){const next=new Set(expanded.value);next.has(key(entry))?next.delete(key(entry)):next.add(key(entry));expanded.value=next}
function parseResult(result){if(!result)return{};try{return typeof result==='string'?JSON.parse(result):result}catch{return{}}}
function parsedImages(result){const obj=parseResult(result),urls=[];if(Array.isArray(obj.data))urls.push(...obj.data.map(x=>x?.url).filter(Boolean));if(Array.isArray(obj.output?.results))urls.push(...obj.output.results.map(x=>x?.url).filter(Boolean));for(const choice of obj.output?.choices||[])for(const item of choice.message?.content||[])if(item.image)urls.push(item.image);return[...new Set(urls)]}
function parsedVideo(result){const obj=parseResult(result);return obj.output?.video_url||obj.output?.results?.[0]?.url||obj.data?.[0]?.url||obj.url||''} function prettyResult(result){try{return JSON.stringify(parseResult(result),null,2)}catch{return result}}
async function loadHistory(){loading.value=true;error.value='';try{const params=activeFilter.value?{type:activeFilter.value}:{};const response=await api.get('/history',{params});entries.value=response.data.entries||[]}catch(e){error.value=e.response?.data?.message||'暂时无法加载创作记录'}loading.value=false}
async function deleteEntry(entry){try{await api.delete(`/history/${entry.id}`,{params:{type:entry.type}});entries.value=entries.value.filter(item=>key(item)!==key(entry))}catch(e){error.value=e.response?.data?.message||'删除失败'}}
onMounted(loadHistory)
</script>

<style scoped>
.history-page{width:min(1180px,100%)}.filter-bar{display:flex;gap:8px;overflow-x:auto;margin-bottom:20px;padding-bottom:3px}.filter-bar button{display:flex;align-items:center;gap:7px;flex:0 0 auto;padding:9px 14px;border:1px solid rgba(var(--v-border-color),.12);border-radius:12px;color:rgb(var(--v-theme-on-surface-variant));background:rgba(var(--v-theme-surface),.55);font:inherit;font-size:12px;cursor:pointer}.filter-bar button.active{border-color:rgba(var(--v-theme-primary),.24);color:rgb(var(--v-theme-primary));background:rgba(var(--v-theme-primary),.1)}.history-list{display:flex;flex-direction:column;gap:11px}.history-item{overflow:hidden}.history-item__top{display:flex;align-items:center;gap:14px;padding:16px 18px;cursor:pointer}.history-icon{width:44px;height:44px;flex:0 0 44px;display:grid;place-items:center;border-radius:14px}.history-copy{min-width:0;flex:1;display:flex;flex-direction:column}.history-copy strong{max-width:650px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-size:14px}.history-copy span{margin-top:5px;color:rgb(var(--v-theme-on-surface-variant));font-size:10px}.expand-icon{color:rgb(var(--v-theme-on-surface-variant));transition:transform .2s}.expand-icon.rotated{transform:rotate(180deg)}.history-detail{padding:2px 18px 20px 76px}.media-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:10px}.media-grid button{overflow:hidden;padding:0;border:0;border-radius:14px;background:none;cursor:zoom-in}.video-preview{width:min(720px,100%);max-height:430px;border-radius:16px;background:#000}.content-preview{padding:16px;border-radius:14px;background:rgb(var(--v-theme-surface-variant));font-size:13px;line-height:1.75;white-space:pre-wrap}.raw-result{margin-top:13px;color:rgb(var(--v-theme-on-surface-variant));font-size:11px}.raw-result summary{cursor:pointer}.raw-result pre{max-height:280px;overflow:auto;padding:13px;border-radius:12px;background:rgba(0,0,0,.1);white-space:pre-wrap;word-break:break-all}.preview-close{position:absolute;right:10px;top:10px}.rounded-xl{border-radius:24px}@media(max-width:650px){.history-item__top{gap:10px;padding:13px}.history-copy strong{max-width:50vw}.history-detail{padding:2px 13px 16px}.history-item__top>.v-btn{display:none}}
</style>
