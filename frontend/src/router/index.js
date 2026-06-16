import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Home', component: () => import('../views/Home.vue') },
  { path: '/chat', name: 'Chat', component: () => import('../views/Chat.vue') },
  { path: '/image', name: 'ImageGen', component: () => import('../views/ImageGen.vue') },
  { path: '/video', name: 'VideoGen', component: () => import('../views/VideoGen.vue') },
  { path: '/tts', name: 'TTS', component: () => import('../views/TTS.vue') },
  { path: '/toolbox', name: 'Toolbox', component: () => import('../views/Toolbox.vue') },
  { path: '/history', name: 'History', component: () => import('../views/History.vue') },
]

export default createRouter({ history: createWebHashHistory(), routes })
