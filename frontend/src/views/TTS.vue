<template>
  <v-container fluid class="pa-6" style="max-width:600px">
    <h1 class="text-h4 font-weight-bold mb-1">Text to Speech</h1>
    <p class="text-body-1 text-medium-emphasis mb-6">Convert text to natural speech</p>
    <v-card variant="outlined" rounded="xl" class="pa-4">
      <v-textarea v-model="text" label="Text" variant="outlined" rows="4" density="comfortable" hide-details class="mb-4" />
      <v-select v-model="voice" :items="voices" label="Voice" variant="outlined" density="comfortable" hide-details class="mb-4" />
      <v-select v-model="format" :items="['mp3','wav','opus']" label="Format" variant="outlined" density="comfortable" hide-details class="mb-4" />
      <v-btn block color="primary" size="large" rounded="lg" :loading="loading" :disabled="!text" @click="generate" variant="flat">Generate Speech</v-btn>
    </v-card>
    <v-card v-if="audioUrl" variant="outlined" rounded="xl" class="mt-4 pa-4">
      <audio :src="audioUrl" controls style="width:100%" />
    </v-card>
    <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mt-4" closable>{{ error }}</v-alert>
  </v-container>
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
