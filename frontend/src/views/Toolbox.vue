<template>
  <div class="page">
    <h2>工具箱</h2>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="翻译" name="translate">
        <h3>文本翻译</h3>
        <el-input v-model="tl.text" type="textarea" :rows="4" placeholder="输入要翻译的文本..." />
        <div style="margin: 8px 0; display:flex; gap:12px">
          <el-select v-model="tl.source" placeholder="源语言"><el-option label="自动检测" value="auto" /><el-option label="中文" value="zh" /><el-option label="英文" value="en" /><el-option label="日文" value="ja" /></el-select>
          <el-select v-model="tl.target" placeholder="目标语言"><el-option label="中文" value="zh" /><el-option label="英文" value="en" /><el-option label="日文" value="ja" /></el-select>
          <el-button type="primary" @click="doTranslate" :loading="tl.loading">翻译</el-button>
        </div>
        <div v-if="tl.result"><el-input v-model="tl.result" type="textarea" :rows="4" readonly /></div>
      </el-tab-pane>
      <el-tab-pane label="OCR" name="ocr">
        <h3>图片文字识别</h3>
        <div class="ocr-upload" @paste="onOCRPaste" tabindex="0">
          <el-upload :auto-upload="false" :on-change="onOCRFile" accept="image/*" :limit="1">
            <el-button>选择图片</el-button>
          </el-upload>
          <span class="paste-hint">或 Ctrl+V 粘贴图片</span>
        </div>
        <div v-if="ocr.previewUrl" style="margin:8px 0"><el-image :src="ocr.previewUrl" fit="contain" style="max-width:300px;max-height:300px" /></div>
        <el-button type="primary" @click="doOCR" :loading="ocr.loading" style="margin-top:8px">开始识别</el-button>
        <div v-if="ocr.result"><el-input v-model="ocr.result" type="textarea" :rows="6" readonly /></div>
      </el-tab-pane>
      <el-tab-pane label="代码执行" name="code">
        <h3>在线代码执行</h3>
        <el-select v-model="code.lang" style="margin-bottom:8px"><el-option label="Python" value="python" /><el-option label="Bash" value="bash" /></el-select>
        <el-input v-model="code.code" type="textarea" :rows="8" placeholder="输入代码..." />
        <el-button type="primary" @click="doCode" :loading="code.loading" style="margin-top:8px">运行</el-button>
        <div v-if="code.result"><pre>{{ code.result }}</pre></div>
      </el-tab-pane>
      <el-tab-pane label="文档处理" name="doc">
        <h3>文档智能处理</h3>
        <el-select v-model="doc.task" style="margin-bottom:8px"><el-option label="摘要" value="summarize" /><el-option label="问答" value="qa" /><el-option label="提取" value="extract" /><el-option label="翻译" value="translate" /></el-select>
        <el-input v-model="doc.text" type="textarea" :rows="6" placeholder="输入文档内容..." />
        <el-input v-if="doc.task==='qa'" v-model="doc.question" placeholder="输入问题..." style="margin-top:8px" />
        <el-button type="primary" @click="doDoc" :loading="doc.loading" style="margin-top:8px">处理</el-button>
        <div v-if="doc.result"><el-input v-model="doc.result" type="textarea" :rows="6" readonly /></div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { translate, ocr, runCode, document_ } from '../api'

const activeTab = ref('translate')

const tl = reactive({ text: '', source: 'auto', target: 'en', loading: false, result: '' })
async function doTranslate() {
  if (!tl.text.trim()) return
  tl.loading = true; tl.result = ''
  try { const { data } = await translate({ text: tl.text, source_lang: tl.source, target_lang: tl.target }); tl.result = data.translated_text } catch (e) { tl.result = String(e.response?.data || e.message) }
  tl.loading = false
}

const ocrData = reactive({ loading: false, result: '', previewUrl: '', file: null })
async function onOCRFile(uploadFile) {
  ocrData.file = uploadFile.raw
  ocrData.previewUrl = URL.createObjectURL(uploadFile.raw)
}
function loadOCRImage(file) {
  if (!file || !file.type.startsWith('image/')) return
  ocrData.file = file
  ocrData.previewUrl = URL.createObjectURL(file)
}
function onOCRPaste(e) {
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      loadOCRImage(item.getAsFile())
      break
    }
  }
}
async function doOCR() {
  if (!ocrData.file) return
  ocrData.loading = true; ocrData.result = ''
  try {
    const reader = new FileReader()
    const base64 = await new Promise((resolve) => { reader.onload = () => resolve(reader.result.split(',')[1]); reader.readAsDataURL(ocrData.file) })
    const { data } = await ocr({ image_base64: base64 })
    ocrData.result = data.text
  } catch (e) { ocrData.result = String(e.response?.data || e.message) }
  ocrData.loading = false
}

const code = reactive({ lang: 'python', code: '', loading: false, result: '' })
async function doCode() {
  if (!code.code.trim()) return
  code.loading = true; code.result = ''
  try { const { data } = await runCode({ language: code.lang, code: code.code }); code.result = data.output } catch (e) { code.result = String(e.response?.data || e.message) }
  code.loading = false
}

const doc = reactive({ task: 'summarize', text: '', question: '', loading: false, result: '' })
async function doDoc() {
  if (!doc.text.trim()) return
  doc.loading = true; doc.result = ''
  try { const { data } = await document_({ task: doc.task, text: doc.text, question: doc.question }); doc.result = data.result } catch (e) { doc.result = String(e.response?.data || e.message) }
  doc.loading = false
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
pre { white-space: pre-wrap; font-size: 13px; background: #f5f5f5; padding: 12px; border-radius: 4px; margin-top: 8px; }
.ocr-upload { display: flex; align-items: center; gap: 12px; padding: 4px 0; outline: none; }
.paste-hint { color: #909399; font-size: 12px; }
</style>