import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Home', meta: { title: '工作台' }, component: () => import('../views/Home.vue') },
  { path: '/chat', name: 'Chat', meta: { title: '智能对话' }, component: () => import('../views/Chat.vue') },
  { path: '/image', name: 'ImageGen', meta: { title: '图片创作' }, component: () => import('../views/ImageGen.vue') },
  { path: '/video', name: 'VideoGen', meta: { title: '视频创作' }, component: () => import('../views/VideoGen.vue') },
  { path: '/tts', name: 'TTS', meta: { title: '声音合成' }, component: () => import('../views/TTS.vue') },
  { path: '/toolbox', name: 'Toolbox', meta: { title: '智能工具' }, component: () => import('../views/Toolbox.vue') },
  { path: '/history', name: 'History', meta: { title: '创作记录' }, component: () => import('../views/History.vue') },
]

const router = createRouter({ history: createWebHashHistory(), routes })
router.afterEach(to => { document.title = `${to.meta.title || '工作台'} · 百炼工作站` })
export default router
