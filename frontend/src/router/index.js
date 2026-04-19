import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/chat' },
  { path: '/chat', name: 'Chat', component: () => import('../views/ChatView.vue') },
  { path: '/vision', name: 'Vision', component: () => import('../views/VisionView.vue') },

  // 图像工厂
  { path: '/image/gen', name: 'ImageGen', component: () => import('../views/ImageGenView.vue') },
  { path: '/image/edit', name: 'ImageEdit', component: () => import('../views/ImageEditView.vue') },
  { path: '/image/tryon', name: 'ImageTryOn', component: () => import('../views/ImageTryOnView.vue') },

  // 视频工厂
  { path: '/video/t2v', name: 'VideoT2v', component: () => import('../views/VideoT2vView.vue') },
  { path: '/video/i2v', name: 'VideoI2v', component: () => import('../views/VideoI2vView.vue') },
  { path: '/video/ref', name: 'VideoRef', component: () => import('../views/VideoRefView.vue') },
  { path: '/video/edit', name: 'VideoEdit', component: () => import('../views/VideoEditView.vue') },
  { path: '/video/digital-human', name: 'DigitalHuman', component: () => import('../views/DigitalHumanView.vue') },
  { path: '/video/animate', name: 'AnimateMove', component: () => import('../views/AnimateMoveView.vue') },

  // 音频中心
  { path: '/audio/tts', name: 'TTS', component: () => import('../views/TtsView.vue') },
  { path: '/audio/asr', name: 'ASR', component: () => import('../views/AsrView.vue') },
]

const router = createRouter({ history: createWebHistory(), routes })
export default router
