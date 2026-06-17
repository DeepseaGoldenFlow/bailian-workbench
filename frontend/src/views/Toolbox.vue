<template>
  <div class="pa-8" style="max-width:900px;margin:0 auto">
    <div class="d-flex align-center mb-8">
      <v-avatar color="blue-darken-2" size="48" class="mr-4 elevation-4"><v-icon size="28">mdi-tools</v-icon></v-avatar>
      <div><h1 class="text-h4 font-weight-bold">Toolbox</h1><p class="text-body-1 text-medium-emphasis mt-1">Translation, OCR, and document analysis</p></div>
    </div>
    <v-card rounded="xl" class="pa-6" elevation="2" style="background:linear-gradient(180deg, rgba(37,99,235,0.06) 0%, transparent 40%)">
      <v-tabs v-model="tab" color="blue-darken-2" class="mb-4">
        <v-tab value="translate">Translate</v-tab>
        <v-tab value="ocr">OCR</v-tab>
        <v-tab value="document">Document</v-tab>
      </v-tabs>
      <div v-if="tab === 'translate'">
        <v-row dense class="mb-4">
          <v-col cols="6"><v-select v-model="srcLang" :items="languages" label="From" variant="outlined" density="comfortable" hide-details /></v-col>
          <v-col cols="6"><v-select v-model="tgtLang" :items="languages" label="To" variant="outlined" density="comfortable" hide-details /></v-col>
        </v-row>
        <v-textarea v-model="transText" label="Text" variant="outlined" rows="4" density="comfortable" hide-details class="mb-4" />
        <v-btn color="blue-darken-2" rounded="lg" :loading="transLoading" @click="doTranslate" variant="elevated">Translate</v-btn>
        <v-card v-if="transResult" variant="tonal" rounded="xl" class="mt-4 pa-4"><div class="text-body-1" style="white-space:pre-wrap">{{ transResult }}</div></v-card>
      </div>
      <div v-if="tab === 'ocr'">
        <v-file-input v-model="ocrFile" label="Image File" variant="outlined" density="comfortable" accept="image/*" hide-details class="mb-4" />
        <v-btn color="blue-darken-2" rounded="lg" :loading="ocrLoading" :disabled="!ocrFile" @click="doOCR" variant="elevated">Extract Text</v-btn>
        <v-card v-if="ocrResult" variant="tonal" rounded="xl" class="mt-4 pa-4"><div class="text-body-1" style="white-space:pre-wrap">{{ ocrResult }}</div></v-card>
      </div>
      <div v-if="tab === 'document'">
        <v-select v-model="docTask" :items="['summarize','qa','extract','translate']" label="Task" variant="outlined" density="comfortable" hide-details class="mb-4" />
        <v-textarea v-model="docText" label="Document Text" variant="outlined" rows="6" density="comfortable" hide-details class="mb-4" />
        <v-text-field v-if="docTask === 'qa'" v-model="docQuestion" label="Question" variant="outlined" density="comfortable" hide-details class="mb-4" />
        <v-btn color="blue-darken-2" rounded="lg" :loading="docLoading" :disabled="!docText" @click="doDocument" variant="elevated">Analyze</v-btn>
        <v-card v-if="docResult" variant="tonal" rounded="xl" class="mt-4 pa-4"><div class="text-body-1" style="white-space:pre-wrap">{{ docResult }}</div></v-card>
      </div>
    </v-card>
    <v-alert v-if="error" type="error" variant="tonal" class="mt-4" closable>{{ error }}</v-alert>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { translate as tr, ocr, document_ as doc } from '../api'
const tab = ref('translate'); const error = ref('')
const languages = [{ title: 'Auto', value: 'auto' }, { title: 'Chinese', value: 'zh' }, { title: 'English', value: 'en' }, { title: 'Japanese', value: 'ja' }, { title: 'Korean', value: 'ko' }, { title: 'French', value: 'fr' }, { title: 'German', value: 'de' }, { title: 'Spanish', value: 'es' }]
const srcLang = ref('auto'); const tgtLang = ref('en'); const transText = ref(''); const transResult = ref(''); const transLoading = ref(false)
async function doTranslate() { transLoading.value = true; error.value = ''; try { const r = await tr({ text: transText.value, source_lang: srcLang.value, target_lang: tgtLang.value }); transResult.value = r.data.translated_text } catch (e) { error.value = e.response?.data?.message || e.message }; transLoading.value = false }
const ocrFile = ref(null); const ocrResult = ref(''); const ocrLoading = ref(false)
async function doOCR() { ocrLoading.value = true; error.value = ''; try { const r = new FileReader(); r.onload = async (e) => { const b64 = e.target.result.split(',')[1]; const resp = await ocr({ image_base64: b64 }); ocrResult.value = resp.data.text }; r.readAsDataURL(ocrFile.value) } catch (e) { error.value = e.response?.data?.message || e.message }; ocrLoading.value = false }
const docTask = ref('summarize'); const docText = ref(''); const docQuestion = ref(''); const docResult = ref(''); const docLoading = ref(false)
async function doDocument() { docLoading.value = true; error.value = ''; try { const r = await doc({ text: docText.value, task: docTask.value, question: docQuestion.value }); docResult.value = r.data.result } catch (e) { error.value = e.response?.data?.message || e.message }; docLoading.value = false }
</script>
