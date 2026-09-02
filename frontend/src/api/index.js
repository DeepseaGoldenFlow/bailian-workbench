import axios from 'axios'

const api = axios.create({ baseURL: '/api', timeout: 300000 })

export function chat(payload) { return api.post('/chat/completions', payload) }
export function chatStream(payload) {
  return fetch('/api/chat/completions', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
}
export function fetchModels(category) {
  return api.get('/models', { params: category ? { category } : {} }).catch(async error => {
    if (!import.meta.env.PROD) throw error
    const { data } = await axios.get(`${import.meta.env.BASE_URL}models.json`)
    const models = category ? data.models.filter(model => model.category === category) : data.models
    return { data: { models } }
  })
}
export function imageGen(payload) { return api.post('/images/generate', payload) }
export function videoGen(payload) { return api.post('/videos/generate', payload) }
export function pollTask(taskId) { return api.get(`/tasks/${taskId}`) }
export function tts(payload) { return api.post('/audio/speech', payload, { responseType: 'blob' }) }
export function asr(formData) { return api.post('/audio/transcribe', formData) }
export function translate(payload) { return api.post('/toolbox/translate', payload) }
export function ocr(payload) { return api.post('/toolbox/ocr', payload) }
export function document_(payload) { return api.post('/toolbox/document', payload) }
export function chatHistory() { return api.get('/history/chat') }
export function genHistory(type) { return api.get('/history/generations', { params: { type } }) }
export function deleteChat(id) { return api.delete(`/history/chat/${id}`) }
export function deleteGen(id) { return api.delete(`/history/generations/${id}`) }
export function health() { return api.get('/health') }

export default api
