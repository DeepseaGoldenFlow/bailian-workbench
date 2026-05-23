import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Chat', component: () => import('../views/Chat.vue') },
  { path: '/image', name: 'ImageGen', component: () => import('../views/ImageGen.vue') },
  { path: '/video', name: 'VideoGen', component: () => import('../views/VideoGen.vue') },
  { path: '/tts', name: 'TTS', component: () => import('../views/TTS.vue') },
  { path: '/asr', name: 'ASR', component: () => import('../views/ASR.vue') },
  { path: '/toolbox', name: 'Toolbox', component: () => import('../views/Toolbox.vue') },
  { path: '/history', name: 'History', component: () => import('../views/History.vue') },
]

export default createRouter({ history: createWebHashHistory(), routes })