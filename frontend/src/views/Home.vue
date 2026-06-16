<template>
  <div class="pa-8">
    <div class="text-center mb-10 pt-6">
      <v-avatar color="primary" size="72" class="mb-4 elevation-8">
        <span class="text-h3 font-weight-bold">N</span>
      </v-avatar>
      <h1 class="text-h3 font-weight-bold mb-2">Nova AI Workbench</h1>
      <p class="text-h6 text-medium-emphasis">Bailian-powered creation studio</p>
    </div>

    <v-row class="mb-2">
      <v-col v-for="card in cards" :key="card.title" cols="12" sm="6" lg="4">
        <v-card :to="card.to" :color="card.color" theme="dark" rounded="xl" class="pa-6" height="200" link
          style="position:relative;overflow:hidden">
          <div style="position:absolute;right:-20px;top:-20px;opacity:0.2">
            <v-icon size="140">{{ card.icon }}</v-icon>
          </div>
          <v-icon size="32" class="mb-4">{{ card.icon }}</v-icon>
          <div class="text-h5 font-weight-bold mb-1">{{ card.title }}</div>
          <div class="text-body-1 opacity-70">{{ card.desc }}</div>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mt-8">
      <v-col cols="12" md="6">
        <v-card rounded="xl" variant="outlined" class="pa-6" height="100%">
          <div class="text-h6 font-weight-bold mb-4">Quick Stats</div>
          <v-row>
            <v-col cols="6" v-for="stat in stats" :key="stat.label">
              <div class="text-h4 font-weight-bold text-primary mb-1">{{ stat.value }}</div>
              <div class="text-body-2 text-medium-emphasis">{{ stat.label }}</div>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card rounded="xl" variant="outlined" class="pa-6" height="100%">
          <div class="text-h6 font-weight-bold mb-4">Recent Activity</div>
          <div v-if="recent.length">
            <div v-for="item in recent.slice(0,5)" :key="item.id" class="d-flex align-center ga-3 py-2" style="border-bottom:1px solid rgba(255,255,255,0.05)">
              <v-avatar :color="typeColor(item.type)" size="28" rounded="lg">
                <v-icon size="14">{{ typeIcon(item.type) }}</v-icon>
              </v-avatar>
              <div class="text-body-2 text-truncate flex-1-1">{{ item.prompt || item.content }}</div>
              <div class="text-caption text-medium-emphasis">{{ item.type }}</div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-medium-emphasis">
            <v-icon size="32" class="mb-2">mdi-clock-outline</v-icon>
            <p class="text-body-2">No activity yet</p>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const recent = ref([])

const cards = [
  { title: 'Chat', desc: 'AI conversation with streaming', icon: 'mdi-chat', to: '/chat', color: '#1a73e8' },
  { title: 'Image Gen', desc: 'Generate stunning visuals', icon: 'mdi-image', to: '/image', color: '#9334e6' },
  { title: 'Video Gen', desc: 'Create motion with HappyHorse', icon: 'mdi-video', to: '/video', color: '#0d9488' },
  { title: 'Voice', desc: 'Text-to-speech synthesis', icon: 'mdi-microphone', to: '/tts', color: '#ea580c' },
  { title: 'Tools', desc: 'Translate, OCR, analyze', icon: 'mdi-tools', to: '/toolbox', color: '#2563eb' },
  { title: 'History', desc: 'All your past creations', icon: 'mdi-history', to: '/history', color: '#4b5563' },
]

const stats = [
  { label: 'Models', value: '11' },
  { label: 'Ready', value: 'Yes' },
]

function typeIcon(t) { const m = { chat:'mdi-chat',image:'mdi-image',video:'mdi-video',audio:'mdi-microphone' }; return m[t]||'mdi-dots-horizontal' }
function typeColor(t) { const m = { chat:'primary',image:'purple',video:'teal',audio:'deep-orange' }; return m[t]||'grey' }

onMounted(async () => {
  try { const r = await api.get('/history', { params: { limit: '10' } }); recent.value = r.data.entries || [] } catch(_) {}
})
</script>
