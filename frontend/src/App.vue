<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" :rail="rail" permanent>
      <div class="d-flex align-center pa-4" :class="rail ? 'justify-center' : ''">
        <v-avatar size="36" style="background:linear-gradient(135deg, var(--v-primary), #a855f7);flex-shrink:0">
          <span class="text-subtitle-1 font-weight-bold">N</span>
        </v-avatar>
        <div v-if="!rail" class="ml-3 flex-1-1">
          <div class="text-subtitle-1 font-weight-bold lh-1">Nova</div>
          <div class="text-caption text-medium-emphasis">AI Workbench</div>
        </div>
        <v-btn
          :icon="rail ? 'mdi-chevron-right' : 'mdi-chevron-left'"
          variant="text"
          size="32"
          @click.stop="rail = !rail"
          class="collapse-btn"
          :class="rail ? 'mt-6' : 'ml-auto'"
          :title="rail ? 'Expand sidebar' : 'Collapse sidebar'"
        />
      </div>

      <v-divider />

      <v-list density="compact" nav class="mt-2">
        <v-list-item v-for="item in navItems" :key="item.path" :to="item.path" :prepend-icon="item.icon" :title="rail ? '' : item.label" color="primary" rounded="lg" class="mb-1" />
      </v-list>

      <template #append>
        <div class="d-flex justify-center pa-2">
          <v-btn
            :icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'"
            variant="text"
            size="32"
            @click="toggleTheme"
            class="control-btn"
            :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
          />
        </div>
      </template>
    </v-navigation-drawer>

    <v-main :style="{ background: isDark ? 'linear-gradient(180deg, #0d1117 0%, #161b22 100%)' : 'linear-gradient(180deg, #f8f9fa 0%, #ffffff 100%)' }">
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useTheme } from 'vuetify'

const drawer = ref(true)
const rail = ref(false)
const theme = useTheme()
const isDark = ref(true)

const navItems = [
  { path: '/', label: 'Home', icon: 'mdi-home-outline' },
  { path: '/chat', label: 'Chat', icon: 'mdi-chat-outline' },
  { path: '/image', label: 'Image', icon: 'mdi-image-outline' },
  { path: '/video', label: 'Video', icon: 'mdi-video-outline' },
  { path: '/tts', label: 'Voice', icon: 'mdi-microphone' },
  { path: '/toolbox', label: 'Tools', icon: 'mdi-tools' },
  { path: '/history', label: 'History', icon: 'mdi-history' },
]

onMounted(() => {
  const saved = localStorage.getItem('nova-theme')
  if (saved === 'light') { theme.global.name.value = 'light'; isDark.value = false }
})

function toggleTheme() {
  const newTheme = theme.global.current.value.dark ? 'light' : 'dark'
  theme.global.name.value = newTheme
  isDark.value = newTheme === 'dark'
  localStorage.setItem('nova-theme', newTheme)
}
</script>

<style>
.v-navigation-drawer__content { display: flex; flex-direction: column; }
.collapse-btn, .control-btn {
  border-radius: 50% !important;
  opacity: 0.7;
  transition: opacity 0.15s, background 0.15s;
}
.collapse-btn:hover, .control-btn:hover { opacity: 1; background: rgba(128,128,128,0.12); }
.lh-1 { line-height: 1.2; }
</style>
