<template>
  <div class="page-shell tools-page">
    <header class="page-header"><div><div class="page-eyebrow">智能工具</div><h1 class="page-title">把繁琐的工作交给 AI</h1><p class="page-subtitle">翻译、图片文字识别与文档分析，三个高频工具集中在一个页面。</p></div></header>
    <div class="tool-layout">
      <div class="tool-nav">
        <button v-for="item in tools" :key="item.value" :class="{ active:tab===item.value }" @click="tab=item.value"><span :style="{color:item.color,background:item.tint}"><v-icon size="22">{{ item.icon }}</v-icon></span><div><strong>{{ item.label }}</strong><small>{{ item.desc }}</small></div><v-icon size="17">mdi-chevron-right</v-icon></button>
      </div>
      <v-card rounded="xl" class="surface-card tool-panel">
        <div v-if="tab==='translate'" class="panel-content">
          <div class="panel-heading"><div><h2>文本翻译</h2><p>自动识别源语言，保留原文语气与结构。</p></div><div class="panel-icon"><v-icon>mdi-translate</v-icon></div></div>
          <div class="language-row"><v-select v-model="srcLang" :items="languages" label="源语言" variant="outlined" hide-details /><v-btn icon="mdi-swap-horizontal" variant="tonal" size="small" @click="swapLanguages" /><v-select v-model="tgtLang" :items="languages.filter(l=>l.value!=='auto')" label="目标语言" variant="outlined" hide-details /></div>
          <v-textarea v-model="transText" label="输入要翻译的内容" rows="8" variant="outlined" class="mt-5" />
          <v-btn color="primary" rounded="xl" size="large" :loading="transLoading" :disabled="!transText.trim()" prepend-icon="mdi-auto-fix" @click="doTranslate">开始翻译</v-btn>
          <div v-if="transResult" class="result-box mt-5"><div class="result-box__header"><strong>翻译结果</strong><v-btn icon="mdi-content-copy" variant="text" size="small" @click="copy(transResult)" /></div><div>{{ transResult }}</div></div>
        </div>
        <div v-else-if="tab==='ocr'" class="panel-content">
          <div class="panel-heading"><div><h2>图片文字识别</h2><p>上传截图、照片或扫描件，提取其中的文字。</p></div><div class="panel-icon"><v-icon>mdi-text-recognition</v-icon></div></div>
          <label class="upload-zone"><input type="file" accept="image/*" @change="selectOCRFile"><div class="upload-zone__icon"><v-icon size="30">mdi-cloud-upload-outline</v-icon></div><strong>{{ ocrFile?.name || '点击选择图片，或拖放到这里' }}</strong><span>支持 JPG、PNG、WEBP 等常见格式</span></label>
          <v-btn color="primary" rounded="xl" size="large" class="mt-5" :loading="ocrLoading" :disabled="!ocrFile" prepend-icon="mdi-text-search" @click="doOCR">提取文字</v-btn>
          <div v-if="ocrResult" class="result-box mt-5"><div class="result-box__header"><strong>识别结果</strong><v-btn icon="mdi-content-copy" variant="text" size="small" @click="copy(ocrResult)" /></div><div>{{ ocrResult }}</div></div>
        </div>
        <div v-else class="panel-content">
          <div class="panel-heading"><div><h2>文档分析</h2><p>快速总结、问答、提取关键信息或翻译长文本。</p></div><div class="panel-icon"><v-icon>mdi-file-document-multiple-outline</v-icon></div></div>
          <v-select v-model="docTask" :items="docTasks" item-title="label" item-value="value" label="处理方式" variant="outlined" hide-details />
          <v-textarea v-model="docText" label="粘贴文档内容" rows="9" variant="outlined" class="mt-5" />
          <v-text-field v-if="docTask==='qa'" v-model="docQuestion" label="你想询问什么？" variant="outlined" class="mt-1" />
          <v-btn color="primary" rounded="xl" size="large" :loading="docLoading" :disabled="!docText.trim()" prepend-icon="mdi-file-search-outline" @click="doDocument">开始分析</v-btn>
          <div v-if="docResult" class="result-box mt-5"><div class="result-box__header"><strong>分析结果</strong><v-btn icon="mdi-content-copy" variant="text" size="small" @click="copy(docResult)" /></div><div>{{ docResult }}</div></div>
        </div>
      </v-card>
    </div>
    <v-snackbar v-model="copied" color="success" timeout="1600">已复制到剪贴板</v-snackbar>
    <v-alert v-if="error" type="error" variant="tonal" rounded="xl" class="mt-4" closable @click:close="error=''">{{ error }}</v-alert>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { translate as translateAPI, ocr, document_ as documentAPI } from '../api'
const tab = ref('translate'), error = ref(''), copied = ref(false)
const tools = [{value:'translate',label:'文本翻译',desc:'多语言自然翻译',icon:'mdi-translate',color:'#2563eb',tint:'rgba(37,99,235,.10)'},{value:'ocr',label:'文字识别',desc:'从图片提取文字',icon:'mdi-text-recognition',color:'#db2777',tint:'rgba(219,39,119,.10)'},{value:'document',label:'文档分析',desc:'总结、问答与提取',icon:'mdi-file-document-multiple-outline',color:'#7c3aed',tint:'rgba(124,58,237,.10)'}]
const languages = [{title:'自动识别',value:'auto'},{title:'中文',value:'zh'},{title:'英语',value:'en'},{title:'日语',value:'ja'},{title:'韩语',value:'ko'},{title:'法语',value:'fr'},{title:'德语',value:'de'},{title:'西班牙语',value:'es'}]
const srcLang = ref('auto'), tgtLang = ref('zh'), transText = ref(''), transResult = ref(''), transLoading = ref(false)
function swapLanguages() { if (srcLang.value === 'auto') return; [srcLang.value,tgtLang.value] = [tgtLang.value,srcLang.value] }
async function doTranslate() { transLoading.value=true; error.value=''; try { const r=await translateAPI({text:transText.value,source_lang:srcLang.value,target_lang:tgtLang.value}); transResult.value=r.data.translated_text||r.data.result||'' } catch(e){ error.value=e.response?.data?.message||e.message||'翻译失败' } transLoading.value=false }
const ocrFile = ref(null), ocrResult = ref(''), ocrLoading = ref(false)
function selectOCRFile(event) { ocrFile.value=event.target.files?.[0]||null }
function readBase64(file) { return new Promise((resolve,reject)=>{ const reader=new FileReader(); reader.onload=()=>resolve(String(reader.result).split(',')[1]); reader.onerror=reject; reader.readAsDataURL(file) }) }
async function doOCR() { if(!ocrFile.value)return; ocrLoading.value=true; error.value=''; try { const image_base64=await readBase64(ocrFile.value); const r=await ocr({image_base64}); ocrResult.value=r.data.text||r.data.result||'' } catch(e){ error.value=e.response?.data?.message||e.message||'文字识别失败' } ocrLoading.value=false }
const docTasks=[{label:'生成摘要',value:'summarize'},{label:'针对文档提问',value:'qa'},{label:'提取关键信息',value:'extract'},{label:'翻译文档',value:'translate'}]
const docTask=ref('summarize'),docText=ref(''),docQuestion=ref(''),docResult=ref(''),docLoading=ref(false)
async function doDocument(){ docLoading.value=true;error.value='';try{const r=await documentAPI({text:docText.value,task:docTask.value,question:docQuestion.value});docResult.value=r.data.result||''}catch(e){error.value=e.response?.data?.message||e.message||'文档分析失败'}docLoading.value=false }
async function copy(value){ await navigator.clipboard.writeText(value); copied.value=true }
</script>

<style scoped>
.tools-page { width:min(1220px,100%); }.tool-layout { display:grid; grid-template-columns:260px 1fr; gap:20px; }.tool-nav { display:flex; flex-direction:column; gap:9px; }.tool-nav button { width:100%; display:flex; align-items:center; gap:12px; padding:15px; border:1px solid transparent; border-radius:16px; color:rgb(var(--v-theme-on-surface)); background:transparent; font:inherit; text-align:left; cursor:pointer; }.tool-nav button:hover { background:rgba(var(--v-theme-surface),.65); }.tool-nav button.active { border-color:rgba(var(--v-theme-primary),.18); background:rgba(var(--v-theme-surface),.9); box-shadow:0 12px 30px rgba(15,23,42,.05); }.tool-nav button>span { width:42px; height:42px; flex:0 0 42px; display:grid; place-items:center; border-radius:13px; }.tool-nav button>div { min-width:0; flex:1; display:flex; flex-direction:column; }.tool-nav strong { font-size:14px; }.tool-nav small { margin-top:3px; color:rgb(var(--v-theme-on-surface-variant)); font-size:10px; }.tool-panel { min-height:610px; }.panel-content { padding:30px; }.panel-heading { display:flex; align-items:center; justify-content:space-between; margin-bottom:28px; }.panel-heading h2 { margin:0 0 7px; font-size:23px; }.panel-heading p { margin:0; color:rgb(var(--v-theme-on-surface-variant)); font-size:12px; }.panel-icon { width:48px; height:48px; display:grid; place-items:center; border-radius:15px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.1); }.language-row { display:grid; grid-template-columns:1fr auto 1fr; align-items:center; gap:12px; }.upload-zone { min-height:260px; display:flex; flex-direction:column; align-items:center; justify-content:center; gap:10px; border:1.5px dashed rgba(var(--v-theme-primary),.3); border-radius:20px; background:rgba(var(--v-theme-primary),.035); cursor:pointer; }.upload-zone:hover { background:rgba(var(--v-theme-primary),.07); }.upload-zone input { display:none; }.upload-zone__icon { width:58px; height:58px; display:grid; place-items:center; margin-bottom:5px; border-radius:18px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.1); }.upload-zone strong { font-size:14px; }.upload-zone span { color:rgb(var(--v-theme-on-surface-variant)); font-size:11px; }.result-box { overflow:hidden; border:1px solid rgba(var(--v-border-color),.1); border-radius:16px; background:rgb(var(--v-theme-surface-variant)); }.result-box__header { display:flex; align-items:center; justify-content:space-between; padding:10px 14px; border-bottom:1px solid rgba(var(--v-border-color),.1); font-size:12px; }.result-box>div:last-child { max-height:300px; overflow:auto; padding:16px; white-space:pre-wrap; font-size:13px; line-height:1.75; }
@media(max-width:800px){.tool-layout{grid-template-columns:1fr}.tool-nav{display:grid;grid-template-columns:repeat(3,1fr)}.tool-nav button{justify-content:center;padding:12px}.tool-nav button>div,.tool-nav button>.v-icon{display:none}.panel-content{padding:22px 18px}.language-row{grid-template-columns:1fr}.language-row>.v-btn{transform:rotate(90deg);justify-self:center}}
</style>
