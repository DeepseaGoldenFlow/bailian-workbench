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

    <v-list rounded="xl" bg-color="transparent">
      <v-list-item v-for="entry in entries" :key="entry.type + '-' + entry.id" rounded="xl" class="mb-2" :class="entry._expanded ? 'bg-surface-variant' : ''" @click="entry._expanded = !entry._expanded">
        <template #prepend>
          <v-avatar :color="typeColor(entry.type)" size="36" rounded="lg">
            <v-icon size="18">{{ typeIcon(entry.type) }}</v-icon>
          </v-avatar>
        </template>
        <v-list-item-title class="text-body-2">{{ entry.prompt || entry.content || '(empty)' }}</v-list-item-title>
        <v-list-item-subtitle class="d-flex align-center ga-2 mt-1">
          <span>{{ entry.type }}</span>
          <span>{{ entry.model }}</span>
          <v-chip v-if="entry.status" :color="statusDotColor(entry.status)" size="x-small" label>{{ entry.status }}</v-chip>
          <span class="ml-auto">{{ entry.created_at }}</span>
        </v-list-item-subtitle>

        <template #append>
          <v-btn icon="mdi-delete-outline" variant="text" size="small" color="error" @click.stop="deleteEntry(entry)" />
        </template>
      </v-list-item>
    </v-list>

    <v-dialog v-model="showPreview" max-width="90vw">
      <v-img :src="previewUrl" max-height="90vh" contain />
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const entries = ref([]); const loading = ref(false); const error = ref('')
const activeFilter = ref(''); const previewUrl = ref(''); const showPreview = ref(false)

const filters = [
  { label: 'All', value: '' }, { label: 'Chat', value: 'chat' }, { label: 'Images', value: 'image' },
  { label: 'Videos', value: 'video' }, { label: 'Audio', value: 'audio' },
  { label: 'Translate', value: 'translate' }, { label: 'OCR', value: 'ocr' },
]

function typeIcon(t) { const m = { chat: 'mdi-chat', image: 'mdi-image', video: 'mdi-video', audio: 'mdi-microphone', translate: 'mdi-translate', ocr: 'mdi-text-recognition', document: 'mdi-file-document', asr: 'mdi-waveform' }; return m[t] || 'mdi-dots-horizontal' }
function typeColor(t) { const m = { chat: 'primary', image: 'purple', video: 'success', audio: 'warning', translate: 'blue', ocr: 'pink' }; return m[t] || 'grey' }
function statusDotColor(s) { if (s === 'SUCCEEDED' || s === 'completed') return 'success'; if (s === 'FAILED') return 'error'; return 'warning' }

onMounted(() => loadHistory())

async function loadHistory() { loading.value = true; error.value = ''; try { const params = activeFilter.value ? { type: activeFilter.value } : {}; const r = await api.get('/history', { params }); entries.value = (r.data.entries || []).map(e => ({ ...e, _expanded: false })) } catch (e) { error.value = 'Failed to load history' }; loading.value = false }

async function deleteEntry(entry) { try { await api.delete('/history/' + entry.id, { params: { type: entry.type } }); entries.value = entries.value.filter(e => !(e.id === entry.id && e.type === entry.type)) } catch (e) { error.value = 'Delete failed' } }
</script>
