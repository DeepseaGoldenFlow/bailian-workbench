<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" :permanent="!mobile" :rail="rail && !mobile" :width="264" :rail-width="76" class="app-sidebar">
      <div class="brand" :class="{ 'brand--rail': rail && !mobile }">
        <div class="brand__mark"><span>百</span></div>
        <div v-if="!rail || mobile" class="brand__copy"><strong>百炼工作站</strong><span>灵感创作中心</span></div>
        <v-btn v-if="!mobile" :icon="rail ? 'mdi-chevron-right' : 'mdi-chevron-left'" variant="text" size="small" class="ml-auto" :aria-label="rail ? '展开侧栏' : '收起侧栏'" @click="rail = !rail" />
      </div>
      <div v-if="!rail || mobile" class="sidebar-label">创作空间</div>
      <v-list nav density="comfortable" class="nav-list">
        <v-list-item v-for="item in navItems" :key="item.path" :to="item.path" :prepend-icon="item.icon" :title="item.label" color="primary" rounded="xl" class="nav-item" @click="mobile && (drawer = false)" />
      </v-list>
      <template #append>
        <div class="sidebar-footer" :class="{ 'sidebar-footer--rail': rail && !mobile }">
          <div v-if="!rail || mobile" class="sidebar-status"><span class="status-dot" /><div><strong>工作站已就绪</strong><span>模型目录自动同步</span></div></div>
          <v-btn :icon="isDark ? 'mdi-white-balance-sunny' : 'mdi-weather-night'" variant="tonal" size="small" :aria-label="isDark ? '切换浅色模式' : '切换深色模式'" @click="toggleTheme" />
        </div>
      </template>
    </v-navigation-drawer>

    <v-app-bar v-if="mobile" flat class="mobile-bar">
      <v-app-bar-nav-icon aria-label="打开导航" @click="drawer = !drawer" />
      <div class="brand__mark brand__mark--small"><span>百</span></div>
      <v-app-bar-title class="font-weight-bold">百炼工作站</v-app-bar-title>
      <v-btn :icon="isDark ? 'mdi-white-balance-sunny' : 'mdi-weather-night'" variant="text" :aria-label="isDark ? '切换浅色模式' : '切换深色模式'" @click="toggleTheme" />
    </v-app-bar>
    <v-main class="app-main"><router-view /></v-main>
  </v-app>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useDisplay, useTheme } from 'vuetify'

const { mobile } = useDisplay()
const drawer = ref(true)
const rail = ref(false)
const theme = useTheme()
const isDark = computed(() => theme.global.current.value.dark)
const navItems = [
  { path: '/', label: '工作台', icon: 'mdi-view-dashboard-outline' },
  { path: '/chat', label: '智能对话', icon: 'mdi-message-processing-outline' },
  { path: '/image', label: '图片创作', icon: 'mdi-image-multiple-outline' },
  { path: '/video', label: '视频创作', icon: 'mdi-movie-open-play-outline' },
  { path: '/tts', label: '声音合成', icon: 'mdi-waveform' },
  { path: '/toolbox', label: '智能工具', icon: 'mdi-creation-outline' },
  { path: '/history', label: '创作记录', icon: 'mdi-history' },
]
function toggleTheme() {
  const next = isDark.value ? 'light' : 'dark'
  theme.global.name.value = next
  localStorage.setItem('nova-theme', next)
}
</script>
