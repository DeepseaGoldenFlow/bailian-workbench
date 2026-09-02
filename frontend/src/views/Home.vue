<template>
  <div class="page-shell home-page">
    <section class="hero">
      <div class="hero__content">
        <div class="hero__badge"><span class="status-dot" /> 阿里云百炼模型已接入</div>
        <h1>把灵感，变成<br><span>看得见的作品</span></h1>
        <p>从一句描述开始，完成图片、视频、声音与文本创作。所有模型参数，都集中在一个清晰、顺手的工作台里。</p>
        <div class="d-flex flex-wrap ga-3 mt-7">
          <v-btn to="/image" color="primary" size="large" rounded="xl" prepend-icon="mdi-image-multiple-outline">开始创作图片</v-btn>
          <v-btn to="/video" variant="tonal" color="primary" size="large" rounded="xl" prepend-icon="mdi-movie-open-play-outline">制作视频</v-btn>
        </div>
      </div>
      <div class="hero__visual" aria-hidden="true">
        <div class="orb orb--one" /><div class="orb orb--two" />
        <div class="creation-window">
          <div class="creation-window__top"><i /><i /><i /><span>创作预览</span></div>
          <div class="creation-window__canvas"><v-icon size="78">mdi-creation-outline</v-icon><div class="canvas-line canvas-line--wide" /><div class="canvas-line" /></div>
          <div class="creation-window__footer"><span>4K</span><span>16:9</span><span>有声视频</span></div>
        </div>
      </div>
    </section>

    <section class="stats-row">
      <div v-for="stat in stats" :key="stat.label" class="stat-item"><v-icon :color="stat.color" size="25">{{ stat.icon }}</v-icon><div><strong>{{ stat.value }}</strong><span>{{ stat.label }}</span></div></div>
    </section>

    <section class="mt-10">
      <div class="d-flex align-end justify-space-between flex-wrap ga-3 mb-5">
        <div><div class="page-eyebrow">创作入口</div><h2 class="section-heading">今天想做点什么？</h2></div>
        <v-btn to="/history" variant="text" color="primary" append-icon="mdi-arrow-right">查看创作记录</v-btn>
      </div>
      <v-row>
        <v-col v-for="card in cards" :key="card.title" cols="12" sm="6" lg="4">
          <v-card :to="card.to" rounded="xl" class="capability-card h-100" link>
            <div class="capability-card__icon" :style="{ color: card.color, background: card.tint }"><v-icon size="28">{{ card.icon }}</v-icon></div>
            <div class="capability-card__copy"><h3>{{ card.title }}</h3><p>{{ card.desc }}</p></div>
            <v-icon class="capability-card__arrow">mdi-arrow-top-right</v-icon>
          </v-card>
        </v-col>
      </v-row>
    </section>

    <v-row class="mt-6">
      <v-col cols="12" lg="7">
        <v-card rounded="xl" class="surface-card pa-6 h-100">
          <div class="d-flex align-center justify-space-between mb-5"><div><h2 class="section-title">最近创作</h2><p class="section-caption mt-1">继续上一次的灵感</p></div><v-btn icon="mdi-refresh" variant="text" size="small" :loading="loading" @click="loadData" /></div>
          <div v-if="recent.length" class="recent-list">
            <div v-for="item in recent.slice(0,5)" :key="`${item.type}-${item.id}`" class="recent-item">
              <div class="recent-item__icon"><v-icon size="18">{{ typeIcon(item.type) }}</v-icon></div>
              <div class="recent-item__copy"><strong>{{ item.prompt || item.content || '未命名创作' }}</strong><span>{{ typeLabel(item.type) }} · {{ item.model || '默认模型' }}</span></div>
              <v-chip v-if="item.status" size="x-small" variant="tonal" :color="statusColor(item.status)">{{ statusLabel(item.status) }}</v-chip>
            </div>
          </div>
          <div v-else class="mini-empty"><v-icon size="34">mdi-inbox-outline</v-icon><span>还没有创作记录，从上面的入口开始吧</span></div>
        </v-card>
      </v-col>
      <v-col cols="12" lg="5">
        <v-card rounded="xl" class="surface-card pa-6 h-100 tip-card">
          <div class="tip-card__icon"><v-icon>mdi-lightbulb-on-outline</v-icon></div>
          <h2 class="section-title mt-5">写好提示词的小窍门</h2>
          <p class="section-caption mt-2">按“主体 + 环境 + 风格 + 光线 + 镜头”的顺序描述，通常能得到更稳定的结果。</p>
          <div class="prompt-example mt-5">“一只橘猫坐在雨夜的便利店窗边，日系电影感，暖色灯光，35mm 镜头，浅景深”</div>
          <v-btn to="/image" variant="text" color="primary" class="mt-3 px-0" append-icon="mdi-arrow-right">试试这个提示词</v-btn>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import api, { fetchModels } from '../api'

const recent = ref([]), imageModels = ref([]), videoModels = ref([]), loading = ref(false)
const modelCount = computed(() => imageModels.value.length + videoModels.value.length || 29)
const stats = computed(() => [
  { label: '创作模型', value: modelCount.value, icon: 'mdi-cube-outline', color: 'primary' },
  { label: '图片模型', value: imageModels.value.length || 14, icon: 'mdi-image-outline', color: 'purple' },
  { label: '视频模型', value: videoModels.value.length || 15, icon: 'mdi-video-outline', color: 'teal' },
  { label: '参数模式', value: '完整', icon: 'mdi-tune-variant', color: 'orange' },
])
const cards = [
  { title: '智能对话', desc: '与通义千问、DeepSeek 深度交流，支持流式回复。', icon: 'mdi-message-processing-outline', to: '/chat', color: '#2563eb', tint: 'rgba(37,99,235,.10)' },
  { title: '图片创作', desc: '文生图、参考图编辑、连续组图与最高 4K 输出。', icon: 'mdi-image-multiple-outline', to: '/image', color: '#9333ea', tint: 'rgba(147,51,234,.10)' },
  { title: '视频创作', desc: '文生、首尾帧、参考、编辑、续写与有声视频。', icon: 'mdi-movie-open-play-outline', to: '/video', color: '#0d9488', tint: 'rgba(13,148,136,.10)' },
  { title: '声音合成', desc: '多种自然音色，把文案快速转换成高质量语音。', icon: 'mdi-waveform', to: '/tts', color: '#ea580c', tint: 'rgba(234,88,12,.10)' },
  { title: '智能工具', desc: '翻译、图片文字识别与长文档分析集中处理。', icon: 'mdi-creation-outline', to: '/toolbox', color: '#db2777', tint: 'rgba(219,39,119,.10)' },
  { title: '创作记录', desc: '统一查看、预览和管理曾经生成的全部内容。', icon: 'mdi-history', to: '/history', color: '#4f46e5', tint: 'rgba(79,70,229,.10)' },
]
const typeMap = { chat:['智能对话','mdi-message-processing-outline'], image:['图片','mdi-image-outline'], video:['视频','mdi-video-outline'], audio:['声音','mdi-waveform'], translate:['翻译','mdi-translate'], ocr:['文字识别','mdi-text-recognition'], document:['文档','mdi-file-document-outline'] }
function typeLabel(type) { return typeMap[type]?.[0] || '其他' }
function typeIcon(type) { return typeMap[type]?.[1] || 'mdi-creation-outline' }
function statusLabel(status) { return ['SUCCEEDED','completed','succeeded'].includes(status) ? '已完成' : ['FAILED','failed'].includes(status) ? '失败' : '处理中' }
function statusColor(status) { return ['SUCCEEDED','completed','succeeded'].includes(status) ? 'success' : ['FAILED','failed'].includes(status) ? 'error' : 'warning' }
async function loadData() {
  loading.value = true
  const [images, videos, history] = await Promise.allSettled([fetchModels('image'), fetchModels('video'), api.get('/history', { params: { limit: '5' } })])
  if (images.status === 'fulfilled') imageModels.value = images.value.data.models || []
  if (videos.status === 'fulfilled') videoModels.value = videos.value.data.models || []
  if (history.status === 'fulfilled') recent.value = history.value.data.entries || []
  loading.value = false
}
onMounted(loadData)
</script>

<style scoped>
.hero { position:relative; min-height:420px; display:grid; grid-template-columns:1.05fr .95fr; align-items:center; overflow:hidden; padding:54px 58px; border:1px solid rgba(var(--v-border-color),.1); border-radius:32px; background:linear-gradient(135deg,rgba(var(--v-theme-surface),.98),rgba(var(--v-theme-primary),.08)); box-shadow:0 28px 80px rgba(30,41,59,.08); }
.hero__content { position:relative; z-index:2; }.hero__badge { width:fit-content; display:flex; align-items:center; gap:10px; padding:8px 13px; border-radius:999px; color:rgb(var(--v-theme-on-surface-variant)); background:rgba(var(--v-theme-success),.08); font-size:12px; font-weight:650; }
.hero h1 { margin:25px 0 18px; font-size:clamp(42px,5vw,70px); line-height:1.06; letter-spacing:-.055em; }.hero h1 span { color:transparent; background:linear-gradient(90deg,#4f46e5,#7c3aed 52%,#0284c7); background-clip:text; -webkit-background-clip:text; }.hero p { max-width:600px; color:rgb(var(--v-theme-on-surface-variant)); font-size:17px; line-height:1.8; }
.hero__visual { position:relative; min-height:320px; display:grid; place-items:center; }.orb { position:absolute; border-radius:50%; filter:blur(2px); }.orb--one { width:260px; height:260px; right:8%; top:5%; background:rgba(124,58,237,.18); }.orb--two { width:180px; height:180px; left:5%; bottom:0; background:rgba(14,165,233,.15); }
.creation-window { position:relative; z-index:1; width:min(430px,92%); overflow:hidden; border:1px solid rgba(255,255,255,.28); border-radius:24px; background:rgba(15,23,42,.88); box-shadow:0 32px 70px rgba(30,41,59,.3); transform:rotate(2deg); }.creation-window__top { display:flex; align-items:center; gap:6px; padding:14px 16px; color:#94a3b8; font-size:11px; }.creation-window__top i { width:7px; height:7px; border-radius:50%; background:#475569; }.creation-window__top span { margin-left:auto; }.creation-window__canvas { height:205px; display:flex; flex-direction:column; align-items:center; justify-content:center; color:#c4b5fd; background:radial-gradient(circle at 50% 20%,rgba(124,58,237,.38),transparent 50%),linear-gradient(145deg,#1e293b,#111827); }.canvas-line { width:34%; height:7px; margin-top:12px; border-radius:99px; background:rgba(255,255,255,.15); }.canvas-line--wide { width:54%; margin-top:22px; }.creation-window__footer { display:flex; gap:8px; padding:14px; }.creation-window__footer span { padding:5px 9px; border-radius:8px; color:#cbd5e1; background:#1e293b; font-size:10px; }
.stats-row { display:grid; grid-template-columns:repeat(4,1fr); margin-top:22px; padding:20px 26px; border:1px solid rgba(var(--v-border-color),.1); border-radius:22px; background:rgba(var(--v-theme-surface),.72); }.stat-item { display:flex; align-items:center; justify-content:center; gap:13px; min-height:54px; border-right:1px solid rgba(var(--v-border-color),.1); }.stat-item:last-child { border-right:0; }.stat-item div { display:flex; flex-direction:column; }.stat-item strong { font-size:22px; line-height:1; }.stat-item span { margin-top:5px; color:rgb(var(--v-theme-on-surface-variant)); font-size:11px; }.section-heading { margin-top:5px; font-size:28px; letter-spacing:-.03em; }
.capability-card { position:relative; display:flex; gap:16px; min-height:146px; padding:22px; border:1px solid rgba(var(--v-border-color),.1); background:rgba(var(--v-theme-surface),.78); box-shadow:0 12px 30px rgba(15,23,42,.04); transition:transform .2s,box-shadow .2s,border-color .2s; }.capability-card:hover { transform:translateY(-4px); border-color:rgba(var(--v-theme-primary),.22); box-shadow:0 20px 42px rgba(15,23,42,.09); }.capability-card__icon { width:50px; height:50px; flex:0 0 50px; display:grid; place-items:center; border-radius:16px; }.capability-card__copy { padding-right:14px; }.capability-card__copy h3 { margin:2px 0 8px; font-size:17px; }.capability-card__copy p { margin:0; color:rgb(var(--v-theme-on-surface-variant)); font-size:13px; line-height:1.65; }.capability-card__arrow { position:absolute; right:18px; top:18px; color:rgb(var(--v-theme-on-surface-variant)); opacity:.6; }
.recent-item { display:flex; align-items:center; gap:13px; padding:13px 0; border-bottom:1px solid rgba(var(--v-border-color),.08); }.recent-item:last-child { border:0; }.recent-item__icon { width:38px; height:38px; flex:0 0 38px; display:grid; place-items:center; border-radius:12px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.09); }.recent-item__copy { min-width:0; flex:1; display:flex; flex-direction:column; }.recent-item__copy strong { overflow:hidden; font-size:13px; text-overflow:ellipsis; white-space:nowrap; }.recent-item__copy span { margin-top:4px; color:rgb(var(--v-theme-on-surface-variant)); font-size:11px; }.mini-empty { min-height:210px; display:flex; flex-direction:column; align-items:center; justify-content:center; gap:12px; color:rgb(var(--v-theme-on-surface-variant)); font-size:13px; }.tip-card { position:relative; overflow:hidden; }.tip-card::after { content:''; position:absolute; right:-70px; bottom:-90px; width:220px; height:220px; border-radius:50%; background:rgba(var(--v-theme-warning),.07); }.tip-card__icon { width:48px; height:48px; display:grid; place-items:center; border-radius:15px; color:rgb(var(--v-theme-warning)); background:rgba(var(--v-theme-warning),.1); }.prompt-example { position:relative; z-index:1; padding:16px; border-left:3px solid rgb(var(--v-theme-warning)); border-radius:4px 14px 14px 4px; color:rgb(var(--v-theme-on-surface-variant)); background:rgba(var(--v-theme-warning),.06); font-size:13px; line-height:1.7; }
@media (max-width:1100px) { .hero { grid-template-columns:1fr; }.hero__visual { display:none; } }
@media (max-width:700px) { .hero { min-height:auto; padding:36px 24px; border-radius:24px; }.hero h1 { font-size:42px; }.stats-row { grid-template-columns:repeat(2,1fr); gap:16px; }.stat-item { border-right:0; justify-content:flex-start; }.stat-item:nth-child(odd) { border-right:1px solid rgba(var(--v-border-color),.1); } }
</style>
