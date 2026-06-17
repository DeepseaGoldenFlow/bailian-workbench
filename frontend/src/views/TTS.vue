<template>
  <div class="pa-8" style="max-width:700px;margin:0 auto">
    <div class="d-flex align-center mb-8">
      <v-avatar color="deep-orange" size="48" class="mr-4 elevation-4"><v-icon size="28">mdi-microphone</v-icon></v-avatar>
      <div><h1 class="text-h4 font-weight-bold">Text to Speech</h1><p class="text-body-1 text-medium-emphasis mt-1">Convert text to natural speech</p></div>
    </div>
    <v-card rounded="xl" class="pa-6" elevation="2" style="background:linear-gradient(180deg, rgba(234,88,12,0.06) 0%, transparent 40%)">
      <v-textarea v-model="text" label="Text to speak" variant="outlined" rows="4" density="comfortable" hide-details class="mb-5" color="deep-orange" />
      <v-row dense class="mb-5">
        <v-col cols="8"><v-select v-model="voice" :items="voices" label="Voice" variant="outlined" density="comfortable" hide-details color="deep-orange" /></v-col>
        <v-col cols="4"><v-select v-model="format" :items="['mp3','wav','opus']" label="Format" variant="outlined" density="comfortable" hide-details /></v-col>
      </v-row>
      <v-btn block color="deep-orange" size="x-large" rounded="lg" :loading="loading" :disabled="!text" @click="generate" variant="elevated"><v-icon start>mdi-sparkles</v-icon> Generate Speech</v-btn>
    </v-card>
    <v-card v-if="audioUrl" rounded="xl" variant="outlined" class="mt-4 pa-4"><audio :src="audioUrl" controls style="width:100%" /></v-card>
    <v-alert v-if="error" type="error" variant="tonal" class="mt-4" closable>{{ error }}</v-alert>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { tts } from '../api'
const voices = ['Cherry', 'Emily', 'Sarah', 'Michael', 'David', 'Jessica', 'Amelia', 'Ethan', 'Mia', 'Lucas']
const text = ref(''); const voice = ref('Cherry'); const format = ref('mp3')
const loading = ref(false); const audioUrl = ref(''); const error = ref('')
async function generate() {
  loading.value = true; error.value = ''; audioUrl.value = ''
  try { const r = await tts({ input: text.value, voice: voice.value, format: format.value }); audioUrl.value = URL.createObjectURL(r.data) }
  catch (e) { error.value = e.response?.data?.message || e.message }
  loading.value = false
}
</script>
