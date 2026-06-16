<template>
  <v-app>
    <v-navigation-drawer permanent rail>
      <div class="d-flex flex-column align-center pt-4 pb-2">
        <v-avatar color="primary" size="40" class="mb-6">
          <span class="text-h6 font-weight-bold">N</span>
        </v-avatar>
      </div>
      <v-list density="compact" nav>
        <v-list-item v-for="item in navItems" :key="item.path" :to="item.path" :prepend-icon="item.icon" :title="item.label" color="primary" rounded="lg" class="mb-1" />
      </v-list>
      <template #append>
        <div class="d-flex justify-center pb-4">
          <v-icon size="14" color="success">mdi-circle</v-icon>
        </div>
      </template>
    </v-navigation-drawer>
    <v-main style="background: linear-gradient(180deg, #0d1117 0%, #161b22 100%);">
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
const navItems = [
  { path: '/', label: 'Home', icon: 'mdi-home-outline' },
  { path: '/chat', label: 'Chat', icon: 'mdi-chat-outline' },
  { path: '/image', label: 'Image', icon: 'mdi-image-outline' },
  { path: '/video', label: 'Video', icon: 'mdi-video-outline' },
  { path: '/tts', label: 'Voice', icon: 'mdi-microphone' },
  { path: '/toolbox', label: 'Tools', icon: 'mdi-tools' },
  { path: '/history', label: 'History', icon: 'mdi-history' },
]
</script>
