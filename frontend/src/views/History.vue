<template>
  <div class="page">
    <header class="page-header"><h1>Execution History</h1><p class="page-desc">All generations and conversations in one place</p></header>
    <div class="toolbar">
      <div class="filters">
        <button v-for="f in filters" :key="f.value" :class="['filter-btn', { active: activeFilter === f.value }]" @click="activeFilter = f.value; loadHistory()">{{ f.label }}</button>
      </div>
      <button class="btn-refresh" @click="loadHistory" :disabled="loading">Refresh</button>
    </div>
    <div v-if="loading" class="loading-bar"><div class="loading-indeterminate"></div></div>
    <div v-if="error" class="error-card"><div class="error-icon">!</div><div class="error-text">{{ error }}</div></div>
    <div v-if="!loading && !error && entries.length === 0" class="empty-state"><div class="empty-icon">&#9716;</div><p>No history yet</p></div>
    <div class="timeline">
      <div v-for="entry in entries" :key="entry.type + '-' + entry.id" :class="['entry-card', entry.type]">
        <div class="entry-meta">
          <span :class="['type-badge', entry.type]">{{ typeLabel(entry.type) }}</span>
          <span class="entry-model">{{ entry.model || 'unknown' }}</span>
          <span v-if="entry.status" :class="['status-dot-sm', entry.status.toLowerCase()]"></span>
          <span class="entry-time">{{ entry.created_at }}</span>
        </div>
        <div class="entry-body">
          <div class="entry-prompt" @click="entry._expanded = !entry._expanded">{{ entry.prompt || entry.content }}</div>
          <div v-if="entry._expanded" class="entry-detail">
            <div v-if="entry.type === 'chat'" class="detail-section"><div class="detail-label">Message</div><div class="detail-text">{{ entry.content }}</div></div>
            <div v-if="parsedImages(entry.result).length" class="detail-section"><div class="detail-label">Images</div><div class="result-thumbs"><img v-for="(url, i) in parsedImages(entry.result)" :key="i" :src="url" @click="previewUrl = url" loading="lazy" /></div></div>
            <div v-if="parsedVideo(entry.result)" class="detail-section"><div class="detail-label">Video</div><video :src="parsedVideo(entry.result)" controls style="max-width:100%;max-height:300px" /></div>
          </div>
        </div>
        <div class="entry-actions">
          <button class="act-btn" @click="entry._expanded = !entry._expanded">{{ entry._expanded ? 'Collapse' : 'Expand' }}</button>
          <button class="act-btn del" @click="deleteEntry(entry)">Delete</button>
        </div>
      </div>
    </div>
    <div v-if="previewUrl" class="preview-overlay" @click="previewUrl = null"><img :src="previewUrl" class="preview-img" /><button class="preview-close">&times;</button></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
const entries = ref([])
const loading = ref(false)
const error = ref('')
const activeFilter = ref('')
const previewUrl = ref(null)
const filters = [
  { label: 'All', value: '' }, { label: 'Chat', value: 'chat' },
  { label: 'Images', value: 'image' }, { label: 'Videos', value: 'video' },
  { label: 'Audio', value: 'audio' }, { label: 'Translate', value: 'translate' },
  { label: 'OCR', value: 'ocr' },
]
function typeLabel(t) {
  const m = { chat: 'Chat', image: 'Image', video: 'Video', audio: 'Audio', translate: 'Translate', ocr: 'OCR', document: 'Doc', asr: 'ASR' }
  return m[t] || t
}
function parsedImages(result) {
  if (!result) return []
  try {
    const obj = JSON.parse(result)
    if (obj.output?.results) return obj.output.results.map(r => r.url).filter(Boolean)
    if (obj.output?.choices) {
      const urls = []
      for (const c of obj.output.choices) { const ct = c.message?.content; if (Array.isArray(ct)) for (const it of ct) { if (it.image) urls.push(it.image) } }
      return urls
    }
    if (Array.isArray(obj)) return obj.filter(u => typeof u === 'string')
  } catch(_) {}
  return []
}
function parsedVideo(result) {
  if (!result) return null
  try { const obj = JSON.parse(result); return obj.output?.video_url || null } catch(_) { return null }
}
onMounted(() => loadHistory())
async function loadHistory() {
  loading.value = true; error.value = ''
  try {
    const params = activeFilter.value ? { type: activeFilter.value } : {}
    const r = await api.get('/history', { params })
    entries.value = (r.data.entries || []).map(e => ({ ...e, _expanded: false }))
  } catch (e) { error.value = e.response?.data?.message || e.message || 'Failed' }
  loading.value = false
}
async function deleteEntry(entry) {
  try {
    await api.delete('/history/' + entry.id, { params: { type: entry.type } })
    entries.value = entries.value.filter(e => !(e.id === entry.id && e.type === entry.type))
  } catch (e) { error.value = 'Delete failed' }
}
</script>

<style scoped>
.page { max-width: 1000px; margin: 0 auto; }
.page-header { margin-bottom: 24px; }
.page-header h1 { font-size: 28px; font-weight: 700; letter-spacing: -0.5px; }
.page-desc { color: var(--text-secondary); margin-top: 4px; font-size: 14px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 12px; flex-wrap: wrap; }
.filters { display: flex; gap: 6px; flex-wrap: wrap; }
.filter-btn { padding: 6px 14px; background: var(--bg-card); border: 1px solid var(--border-card); border-radius: var(--radius-sm); color: var(--text-secondary); font-size: 12px; font-weight: 500; cursor: pointer; transition: all var(--transition); }
.filter-btn:hover { background: var(--bg-card-hover); color: var(--text-primary); }
.filter-btn.active { background: var(--accent-soft); color: var(--accent); border-color: var(--accent); }
.btn-refresh { padding: 6px 16px; background: var(--accent); border: none; border-radius: var(--radius-sm); color: #fff; font-size: 13px; cursor: pointer; }
.btn-refresh:disabled { opacity: 0.4; }
.loading-bar { height: 2px; background: var(--bg-card); border-radius: 1px; margin-bottom: 16px; overflow: hidden; }
.loading-indeterminate { width: 30%; height: 100%; background: var(--accent); animation: progress 1.5s ease-in-out infinite; }
@keyframes progress { 0% { transform: translateX(-100%); } 100% { transform: translateX(400%); } }
.error-card { display: flex; gap: 12px; padding: 16px; background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.2); border-radius: var(--radius-md); margin-bottom: 16px; }
.error-icon { width: 24px; height: 24px; background: var(--danger); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 700; color: #fff; flex-shrink: 0; }
.error-text { font-size: 13px; color: #fca5a5; }
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 80px 32px; color: var(--text-muted); }
.empty-icon { font-size: 40px; margin-bottom: 12px; opacity: 0.5; }
.timeline { display: flex; flex-direction: column; gap: 8px; }
.entry-card { background: var(--bg-card); border: 1px solid var(--border-card); border-radius: var(--radius-md); padding: 16px; transition: border var(--transition); }
.entry-card:hover { border-color: var(--border-subtle); }
.entry-meta { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; flex-wrap: wrap; }
.type-badge { padding: 2px 8px; border-radius: 4px; font-size: 10px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px; }
.type-badge.chat { background: rgba(99,102,241,0.15); color: var(--accent); }
.type-badge.image { background: rgba(168,85,247,0.15); color: #a855f7; }
.type-badge.video { background: rgba(34,197,94,0.15); color: var(--success); }
.type-badge.audio { background: rgba(245,158,11,0.15); color: var(--warning); }
.type-badge.translate { background: rgba(59,130,246,0.15); color: #3b82f6; }
.type-badge.ocr { background: rgba(236,72,153,0.15); color: #ec4899; }
.entry-model { font-size: 11px; color: var(--text-muted); font-family: monospace; }
.status-dot-sm { width: 6px; height: 6px; border-radius: 50%; }
.status-dot-sm.succeeded { background: var(--success); }
.status-dot-sm.pending, .status-dot-sm.running { background: var(--warning); animation: pulse 1.5s infinite; }
.status-dot-sm.failed { background: var(--danger); }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.3; } }
.entry-time { font-size: 11px; color: var(--text-muted); margin-left: auto; font-family: monospace; }
.entry-prompt { font-size: 13px; color: var(--text-secondary); cursor: pointer; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 600px; }
.entry-prompt:hover { color: var(--text-primary); }
.entry-detail { margin-top: 12px; padding-top: 12px; border-top: 1px solid var(--border-subtle); display: flex; flex-direction: column; gap: 12px; }
.detail-label { font-size: 10px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 4px; }
.detail-text { font-size: 13px; color: var(--text-secondary); line-height: 1.5; }
.result-thumbs { display: flex; gap: 8px; flex-wrap: wrap; }
.result-thumbs img { width: 80px; height: 80px; object-fit: cover; border-radius: var(--radius-sm); cursor: pointer; border: 1px solid var(--border-card); }
.result-thumbs img:hover { border-color: var(--accent); }
.entry-actions { display: flex; gap: 8px; margin-top: 10px; }
.act-btn { padding: 4px 12px; background: none; border: 1px solid var(--border-card); border-radius: var(--radius-sm); color: var(--text-secondary); font-size: 11px; cursor: pointer; transition: all var(--transition); }
.act-btn:hover { background: var(--bg-card-hover); color: var(--text-primary); }
.act-btn.del:hover { background: rgba(239,68,68,0.1); color: var(--danger); border-color: rgba(239,68,68,0.3); }
.preview-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.9); display: flex; align-items: center; justify-content: center; z-index: 100; cursor: pointer; }
.preview-img { max-width: 90vw; max-height: 90vh; border-radius: var(--radius-md); }
.preview-close { position: absolute; top: 24px; right: 24px; background: none; border: none; color: #fff; font-size: 32px; cursor: pointer; }
</style>
