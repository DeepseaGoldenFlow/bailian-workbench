<template>
  <v-container fluid class="pa-6" style="max-width:1000px">
    <h1 class="text-h4 font-weight-bold mb-1">History</h1>
    <p class="text-body-1 text-medium-emphasis mb-6">All generations and conversations</p>

    <v-row class="mb-4">
      <v-chip-group v-model="activeFilter" mandatory @update:model-value="loadHistory">
        <v-chip v-for="f in filters" :key="f.value" :value="f.value" variant="tonal" size="small">{{ f.label }}</v-chip>
      </v-chip-group>
      <v-spacer />
      <v-btn variant="text" color="primary" @click="loadHistory" :loading="loading" icon="mdi-refresh" />
    </v-row>

    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4" />
    <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4" closable>{{ error }}</v-alert>

    <div v-if="!loading && !error && entries.length === 0" class="text-center py-16 text-medium-emphasis">
      <v-icon size="48" class="mb-2">mdi-history</v-icon>
      <p>No history yet</p>
    </div>

    <v-expansion-panels v-model="expanded" multiple variant="accordion" rounded="xl">
      <v-expansion-panel v-for="entry in entries" :key="entry.type + '-' + entry.id" rounded="xl" class="mb-2" :value="entry.type + '-' + entry.id">
        <v-expansion-panel-title>
          <template #default="{ expanded: isExpanded }">
            <v-avatar :color="typeColor(entry.type)" size="36" rounded="lg" class="mr-3">
              <v-icon size="18">{{ typeIcon(entry.type) }}</v-icon>
            </v-avatar>
            <div class="flex-1-1 mr-2">
              <div class="text-body-2 text-truncate" style="max-width:500px">{{ entry.prompt || entry.content || '(empty)' }}</div>
              <div class="d-flex align-center ga-2 mt-1">
                <span class="text-caption text-medium-emphasis">{{ entry.type }}</span>
                <span class="text-caption text-medium-emphasis">{{ entry.model }}</span>
                <v-chip v-if="entry.status" :color="statusDotColor(entry.status)" size="x-small" label class="ml-1">{{ entry.status }}</v-chip>
                <span class="text-caption text-medium-emphasis ml-auto">{{ entry.created_at }}</span>
              </div>
            </div>
            <v-btn icon="mdi-delete-outline" variant="text" size="small" color="error" @click.stop="deleteEntry(entry)" />
          </template>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div v-if="entry.type === 'image' && parsedImages(entry.result).length" class="mb-3">
            <div class="text-caption text-medium-emphasis mb-2">Generated Images</div>
            <div class="d-flex flex-wrap ga-2">
              <v-card v-for="(url, i) in parsedImages(entry.result)" :key="i" width="120" rounded="lg" class="cursor-pointer" @click="previewUrl = url; showPreview = true">
                <v-img :src="url" cover height="120" />
              </v-card>
            </div>
          </div>
          <div v-if="entry.type === 'video' && parsedVideo(entry.result)" class="mb-3">
            <div class="text-caption text-medium-emphasis mb-2">Generated Video</div>
            <video :src="parsedVideo(entry.result)" controls style="max-width:100%;max-height:300px;border-radius:12px" />
          </div>
          <div v-if="entry.type === 'chat' || entry.type === 'translate' || entry.type === 'ocr' || entry.type === 'document'">
            <div class="text-caption text-medium-emphasis mb-2">Content</div>
            <div class="text-body-2 pa-3 rounded-lg bg-surface" style="white-space:pre-wrap;max-height:300px;overflow-y:auto">{{ entry.content || entry.prompt }}</div>
          </div>
          <div v-if="entry.result" class="mt-3">
            <div class="text-caption text-medium-emphasis mb-1">Raw Response</div>
            <div class="text-caption pa-3 rounded-lg bg-surface font-monospace" style="max-height:200px;overflow-y:auto;white-space:pre-wrap;word-break:break-all">{{ entry.result }}</div>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>

    <v-dialog v-model="showPreview" max-width="90vw">
      <v-img :src="previewUrl" max-height="90vh" contain rounded="lg" />
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const entries = ref([]); const loading = ref(false); const error = ref('')
const activeFilter = ref(''); const expanded = ref([])
const previewUrl = ref(''); const showPreview = ref(false)

const filters = [
  { label: 'All', value: '' }, { label: 'Chat', value: 'chat' }, { label: 'Images', value: 'image' },
  { label: 'Videos', value: 'video' }, { label: 'Audio', value: 'audio' },
  { label: 'Translate', value: 'translate' }, { label: 'OCR', value: 'ocr' },
]

function typeIcon(t) { const m = { chat: 'mdi-chat', image: 'mdi-image', video: 'mdi-video', audio: 'mdi-microphone', translate: 'mdi-translate', ocr: 'mdi-text-recognition', document: 'mdi-file-document', asr: 'mdi-waveform' }; return m[t] || 'mdi-dots-horizontal' }
function typeColor(t) { const m = { chat: 'primary', image: 'purple', video: 'success', audio: 'warning', translate: 'blue', ocr: 'pink' }; return m[t] || 'grey' }
function statusDotColor(s) { if (s === 'SUCCEEDED' || s === 'completed') return 'success'; if (s === 'FAILED') return 'error'; return 'warning' }

function parsedImages(result) {
  if (!result) return []
  try {
    const obj = JSON.parse(result)
    if (obj.output?.results) return obj.output.results.map(r => r.url).filter(Boolean)
    if (obj.output?.choices) { const urls = []; for (const c of obj.output.choices) { const ct = c.message?.content; if (Array.isArray(ct)) { for (const it of ct) { if (it.image) urls.push(it.image) } } }; return urls }
    if (Array.isArray(obj)) return obj.filter(u => typeof u === 'string')
  } catch(_) {}
  return []
}

function parsedVideo(result) {
  if (!result) return null
  try { return JSON.parse(result).output?.video_url || null } catch(_) { return null }
}

onMounted(() => loadHistory())

async function loadHistory() { loading.value = true; error.value = ''; try { const params = activeFilter.value ? { type: activeFilter.value } : {}; const r = await api.get('/history', { params }); entries.value = r.data.entries || [] } catch (e) { error.value = 'Failed to load history' }; loading.value = false }

async function deleteEntry(entry) { try { await api.delete('/history/' + entry.id, { params: { type: entry.type } }); entries.value = entries.value.filter(e => !(e.id === entry.id && e.type === entry.type)) } catch (e) { error.value = 'Delete failed' } }
</script>
