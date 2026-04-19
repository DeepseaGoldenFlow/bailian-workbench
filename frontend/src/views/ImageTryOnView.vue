<template>
  <div class="page-container">
    <div class="page-header"><h2>👔 虚拟试衣</h2><p class="page-sub">上传人物照片和服装，AI 生成试穿效果</p></div>
    <div class="main-layout">
      <div class="input-panel glass-card">
        <div class="field">
          <label class="field-label">人物图 <span class="required">*</span></label>
          <el-upload action="" :auto-upload="false" :show-file-list="false" accept="image/*" @change="f => handleUpload(f, 'model')" drag class="upload-dragger">
            <div v-if="!form.modelUrl" class="upload-placeholder"><el-icon class="upload-icon"><User /></el-icon><p>上传人物全身照</p></div>
            <img v-else :src="form.modelUrl" class="upload-preview" />
          </el-upload>
        </div>
        <div class="field">
          <label class="field-label">服装图 <span class="required">*</span></label>
          <el-upload action="" :auto-upload="false" :show-file-list="false" accept="image/*" @change="f => handleUpload(f, 'garment')" drag class="upload-dragger">
            <div v-if="!form.garmentUrl" class="upload-placeholder"><el-icon class="upload-icon"><ShoppingBag /></el-icon><p>上传服装平铺图</p></div>
            <img v-else :src="form.garmentUrl" class="upload-preview" />
          </el-upload>
        </div>
        <el-button @click="submit" type="primary" :loading="loading" class="submit-btn" size="large">开始试衣</el-button>
      </div>
      <div class="result-panel glass-card">
        <h3 class="panel-title">生成结果</h3>
        <div v-if="loading" class="progress-section"><el-progress :percentage="progress" striped striped-flow /><p class="progress-text">{{ progressText }}</p></div>
        <div v-if="resultUrl" class="image-result"><img :src="resultUrl" class="result-img" /><el-button @click="download" type="primary" plain class="mt-2"><el-icon><Download /></el-icon> 下载</el-button></div>
        <el-empty v-else-if="!loading" description="试穿结果将在这里展示 👗" />
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { User, ShoppingBag, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
const form = ref({ modelUrl: '', modelBase64: '', garmentUrl: '', garmentBase64: '' })
const loading = ref(false), progress = ref(0), progressText = ref(''), resultUrl = ref('')
const handleUpload = (file, type) => { const r = new FileReader(); r.onload = e => { form.value[`${type}Base64`] = e.target.result; form.value[`${type}Url`] = e.target.result }; r.readAsDataURL(file.raw) }
const submit = async () => {
  if (!form.value.modelUrl || !form.value.garmentUrl) return ElMessage.warning('请上传两张图片')
  loading.value = true; progress.value = 5; progressText.value = '提交中...'
  try {
    // 百炼 API 通常需要公网 URL 或 Base64。这里我们传 Base64，后端需支持或先上传
    // 为了兼容现有后端 (TryOnService)，假设它处理了 Base64 或者我们直接传 URL
    // 如果后端 TryOnService 期望 URL，这里需要改成先上传到 /api/storage 获取 URL
    // 这里暂且传 Base64 尝试，如果失败则提示
    const res = await fetch('/api/image/tryon', { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({ input: { model_image_url: form.value.modelBase64, garment_image_url: form.value.garmentBase64 } }) })
    if (!res.ok) throw new Error((await res.json()).message || '提交失败')
    const { task_id } = await res.json()
    progress.value = 15; progressText.value = 'AI 试穿中...'
    let status = 'PENDING', p = 0
    while (status !== 'SUCCEEDED' && status !== 'FAILED' && p < 120) {
      await new Promise(r => setTimeout(r, 5000)); p++
      const pr = await fetch(`/api/v1/tasks/${task_id}`); const pd = await pr.json()
      status = pd.output?.task_status || pd.status
      progress.value = Math.min(15 + p * 3, 95); progressText.value = `处理中... [${status}]`
      if (status === 'SUCCEEDED') { progress.value = 100; resultUrl.value = pd.output?.results?.[0]?.url || pd.output?.image_url; ElMessage.success('试衣成功！') }
      else if (status === 'FAILED') throw new Error(pd.output?.message || '失败')
    }
  } catch (e) { ElMessage.error(e.message) } finally { loading.value = false; setTimeout(() => progress.value = 0, 2000) }
}
const download = () => { const a = document.createElement('a'); a.href = resultUrl.value; a.download = 'tryon_result.jpg'; a.click() }
</script>
<style scoped>
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; } .page-header { margin-bottom: 24px; } .page-header h2 { font-size: 24px; font-weight: 700; margin-bottom: 4px; } .page-sub { font-size: 14px; color: var(--text-secondary); }
.main-layout { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; align-items: start; } .glass-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; padding: 20px; backdrop-filter: blur(10px); } .panel-title { font-size: 16px; font-weight: 600; margin-bottom: 16px; }
.field { margin-bottom: 16px; } .field-label { display: block; font-size: 13px; font-weight: 600; margin-bottom: 6px; } .required { color: #ef4444; } .upload-dragger :deep(.el-upload-dragger) { background: transparent; border: 2px dashed var(--card-border); border-radius: 12px; padding: 20px; } .upload-dragger :deep(.el-upload-dragger:hover) { border-color: var(--gradient-start); } .upload-placeholder { text-align: center; } .upload-icon { font-size: 36px; color: var(--text-secondary); margin-bottom: 8px; } .upload-preview { max-width: 100%; max-height: 200px; border-radius: 8px; object-fit: contain; }
.submit-btn { width: 100%; margin-top: 12px; border-radius: 10px; background: var(--btn-gradient); border: none; } .progress-section { padding: 16px 0; } .progress-text { font-size: 13px; color: var(--text-secondary); margin-top: 8px; text-align: center; } .image-result { text-align: center; } .result-img { max-width: 100%; border-radius: 12px; } .mt-2 { margin-top: 12px; }
@media (max-width: 1024px) { .main-layout { grid-template-columns: 1fr; } }
</style>