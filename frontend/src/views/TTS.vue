<template>
  <div class="page-shell voice-page">
    <header class="page-header"><div><div class="page-eyebrow">声音合成</div><h1 class="page-title">让文字拥有自然的声音</h1><p class="page-subtitle">选择喜欢的音色，把文案、旁白或故事快速转换成清晰自然的语音。</p></div></header>
    <v-row>
      <v-col cols="12" lg="7">
        <v-card rounded="xl" class="surface-card pa-6 pa-md-8">
          <div class="field-heading"><span>01</span><div><strong>输入文本</strong><small>建议使用完整标点控制停顿节奏</small></div><em>{{ text.length }}/2000</em></div>
          <v-textarea v-model="text" placeholder="在这里输入要转换成语音的内容……" rows="10" auto-grow max-rows="14" variant="solo-filled" flat counter="2000" :maxlength="2000" class="mt-5 voice-textarea" />
          <div class="field-heading mt-7"><span>02</span><div><strong>输出设置</strong><small>选择音色和音频格式</small></div></div>
          <v-row class="mt-3">
            <v-col cols="12" sm="8"><v-select v-model="voice" :items="voices" item-title="label" item-value="value" label="音色" variant="outlined" hide-details /></v-col>
            <v-col cols="12" sm="4"><v-select v-model="format" :items="formats" item-title="label" item-value="value" label="格式" variant="outlined" hide-details /></v-col>
          </v-row>
          <v-btn block color="primary" size="x-large" rounded="xl" class="mt-7" :loading="loading" :disabled="!text.trim()" prepend-icon="mdi-waveform" @click="generate">生成语音</v-btn>
        </v-card>
      </v-col>
      <v-col cols="12" lg="5">
        <v-card rounded="xl" class="surface-card preview-card h-100">
          <div v-if="!audioUrl" class="empty-state h-100"><div><div class="empty-state__icon"><v-icon size="34">mdi-headphones</v-icon></div><h2 class="section-title">试听将在这里出现</h2><p class="section-caption mt-2">输入文字并选择音色后，点击“生成语音”。</p></div></div>
          <div v-else class="audio-result">
            <div class="audio-cover"><div class="audio-cover__rings" /><v-icon size="58">mdi-waveform</v-icon></div>
            <div class="text-center"><h2 class="section-title">语音生成完成</h2><p class="section-caption mt-2">{{ currentVoiceLabel }} · {{ format.toUpperCase() }}</p></div>
            <audio :src="audioUrl" controls class="audio-player" />
            <v-btn :href="audioUrl" :download="`百炼语音.${format}`" variant="tonal" color="primary" rounded="xl" prepend-icon="mdi-download">下载音频</v-btn>
          </div>
        </v-card>
      </v-col>
    </v-row>
    <v-alert v-if="error" type="error" variant="tonal" rounded="xl" class="mt-4" closable @click:close="error=''">{{ error }}</v-alert>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { tts } from '../api'
const voices = [{value:'Cherry',label:'Cherry · 温柔女声'},{value:'Emily',label:'Emily · 清晰女声'},{value:'Sarah',label:'Sarah · 成熟女声'},{value:'Michael',label:'Michael · 沉稳男声'},{value:'David',label:'David · 自然男声'},{value:'Jessica',label:'Jessica · 活力女声'},{value:'Amelia',label:'Amelia · 叙事女声'},{value:'Ethan',label:'Ethan · 青年男声'},{value:'Mia',label:'Mia · 轻柔女声'},{value:'Lucas',label:'Lucas · 磁性男声'}]
const formats = [{value:'mp3',label:'MP3'},{value:'wav',label:'WAV'},{value:'opus',label:'OPUS'}]
const text = ref(''), voice = ref('Cherry'), format = ref('mp3'), loading = ref(false), audioUrl = ref(''), error = ref('')
const currentVoiceLabel = computed(() => voices.find(item => item.value === voice.value)?.label || voice.value)
async function generate() {
  loading.value = true; error.value = ''
  if (audioUrl.value) URL.revokeObjectURL(audioUrl.value)
  audioUrl.value = ''
  try { const response = await tts({ input:text.value.trim(), voice:voice.value, format:format.value }); audioUrl.value = URL.createObjectURL(response.data) }
  catch (e) { error.value = e.response?.data?.message || e.message || '语音生成失败' }
  loading.value = false
}
onBeforeUnmount(() => audioUrl.value && URL.revokeObjectURL(audioUrl.value))
</script>

<style scoped>
.voice-page { width:min(1200px,100%); }.field-heading { display:flex; align-items:center; gap:13px; }.field-heading>span { width:34px; height:34px; display:grid; place-items:center; border-radius:11px; color:rgb(var(--v-theme-primary)); background:rgba(var(--v-theme-primary),.1); font-size:11px; font-weight:800; }.field-heading>div { display:flex; flex-direction:column; }.field-heading strong { font-size:15px; }.field-heading small { margin-top:3px; color:rgb(var(--v-theme-on-surface-variant)); font-size:11px; }.field-heading em { margin-left:auto; color:rgb(var(--v-theme-on-surface-variant)); font-size:11px; font-style:normal; }.voice-textarea :deep(.v-field) { border-radius:18px; }.preview-card { min-height:530px; }.audio-result { min-height:530px; display:flex; flex-direction:column; align-items:center; justify-content:center; gap:24px; padding:38px; }.audio-cover { position:relative; width:180px; height:180px; display:grid; place-items:center; overflow:hidden; border-radius:50%; color:white; background:linear-gradient(145deg,#4f46e5,#8b5cf6 55%,#0ea5e9); box-shadow:0 24px 60px rgba(79,70,229,.3); }.audio-cover__rings { position:absolute; width:100px; height:100px; border:1px solid rgba(255,255,255,.3); border-radius:50%; box-shadow:0 0 0 24px rgba(255,255,255,.07),0 0 0 48px rgba(255,255,255,.04); }.audio-player { width:min(360px,100%); }
</style>
