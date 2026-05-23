<template>
  <div class="page">
    <h2>历史记录</h2>
    <el-tabs v-model="tab">
      <el-tab-pane label="对话历史" name="chat">
        <el-table :data="chatList" stripe>
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="role" label="角色" width="80" />
          <el-table-column prop="model" label="模型" width="120" />
          <el-table-column prop="content" label="内容" show-overflow-tooltip />
          <el-table-column prop="created_at" label="时间" width="180" />
          <el-table-column label="操作" width="80">
            <template #default="scope">
              <el-button type="danger" size="small" @click="delChat(scope.row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="生成历史" name="gen">
        <el-select v-model="genType" style="margin-bottom:12px" @change="loadGen">
          <el-option label="全部" value="" />
          <el-option label="图片" value="image" />
          <el-option label="视频" value="video" />
          <el-option label="音频" value="audio" />
          <el-option label="翻译" value="translate" />
          <el-option label="OCR" value="ocr" />
          <el-option label="代码" value="code" />
          <el-option label="文档" value="document" />
          <el-option label="ASR" value="asr" />
        </el-select>
        <el-table :data="genList" stripe>
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="type" label="类型" width="70" />
          <el-table-column prop="model" label="模型" width="150" show-overflow-tooltip />
          <el-table-column prop="prompt" label="提示词" show-overflow-tooltip />
          <el-table-column label="状态" width="110">
            <template #default="scope">
              <el-tag :type="statusTag(scope.row.status)" size="small">{{ scope.row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180" />
          <el-table-column label="操作" width="160">
            <template #default="scope">
              <el-button v-if="scope.row.status==='PENDING'||scope.row.status==='RUNNING'" type="primary" size="small" @click="viewTask(scope.row)">查看</el-button>
              <el-button v-else-if="scope.row.result" type="success" size="small" @click="viewTask(scope.row)">结果</el-button>
              <el-button type="danger" size="small" @click="delGen(scope.row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="detailVisible" title="任务详情" width="700px">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="任务ID">{{ detail.task_id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTag(detail.status)">{{ detail.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="模型">{{ detail.model }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ detail.type }}</el-descriptions-item>
        <el-descriptions-item label="提示词" :span="2">{{ detail.prompt }}</el-descriptions-item>
      </el-descriptions>
      <div v-if="detail.status==='PENDING'||detail.status==='RUNNING'" style="margin-top:12px">
        <el-button @click="refreshDetail" :loading="detailPolling" type="primary">刷新状态</el-button>
      </div>
      <div v-if="detailVideo" style="margin-top:12px">
        <video :src="detailVideo" controls style="max-width:100%;max-height:400px" />
      </div>
      <div v-if="detailImages.length" class="detail-imgs">
        <el-image v-for="(img,i) in detailImages" :key="i" :src="img" fit="contain" style="max-width:300px;max-height:300px" :preview-src-list="detailImages" />
      </div>
      <div v-if="detailText" style="margin-top:12px;background:#f5f7fa;padding:12px;border-radius:4px;white-space:pre-wrap">{{ detailText }}</div>
      <div v-if="detailRaw" style="margin-top:12px">
        <el-collapse><el-collapse-item title="原始响应"><pre>{{ detailRaw }}</pre></el-collapse-item></el-collapse>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { chatHistory, genHistory, deleteChat, deleteGen } from '../api'

const tab = ref('chat')
const chatList = ref([])
const genList = ref([])
const genType = ref('')
const detailVisible = ref(false)
const detail = ref({})
const detailVideo = ref('')
const detailImages = ref([])
const detailText = ref('')
const detailRaw = ref('')
const detailPolling = ref(false)
let pollTimer = null

function statusTag(s) {
  switch (s) {
    case 'SUCCEEDED': return 'success'
    case 'FAILED': return 'danger'
    case 'RUNNING': return 'warning'
    case 'PENDING': return 'info'
    default: return ''
  }
}

async function loadChat() {
  try { const { data } = await chatHistory(); chatList.value = data } catch {}
}
async function loadGen() {
  try { const { data } = await genHistory(genType.value); genList.value = data; scheduleAutoPoll() } catch {}
}
async function delChat(id) {
  try { await deleteChat(id); await loadChat() } catch {}
}
async function delGen(id) {
  try { await deleteGen(id); await loadGen() } catch {}
}

function viewTask(row) {
  detail.value = row
  detailVideo.value = ''
  detailImages.value = []
  detailText.value = ''
  detailRaw.value = ''
  if (row.result) {
    try {
      const parsed = JSON.parse(row.result)
      detailRaw.value = JSON.stringify(parsed, null, 2)
      if (typeof parsed === 'string') {
        detailText.value = parsed
      } else {
        const url = parsed?.output?.video_url
        if (url) detailVideo.value = url
        const results = parsed?.output?.results
        if (results) detailImages.value = results.map(r => r.url).filter(Boolean)
        const choices = parsed?.output?.choices
        if (choices) {
          const imgs = []
          for (const c of choices) {
            for (const item of c?.message?.content || []) {
              if (item.image) imgs.push(item.image)
            }
          }
          if (imgs.length) detailImages.value = imgs
        }
      }
    } catch {
      detailText.value = row.result
    }
  }
  detailVisible.value = true
  if (row.task_id) {
    refreshDetail()
  }
}

async function refreshDetail() {
  if (!detail.value.task_id) return
  detailPolling.value = true
  try {
    const resp = await fetch('/api/tasks/' + detail.value.task_id)
    const data = await resp.json()
    detailRaw.value = JSON.stringify(data, null, 2)
    const newStatus = data.output?.task_status
    if (newStatus) {
      detail.value.status = newStatus
      const row = genList.value.find(r => r.task_id === detail.value.task_id)
      if (row) {
        row.status = newStatus
        if (newStatus === 'SUCCEEDED' || newStatus === 'FAILED') {
          row.result = JSON.stringify(data)
        }
      }
    }
    const url = data.output?.video_url
    if (url) detailVideo.value = url
    const results = data.output?.results
    if (results) detailImages.value = results.map(r => r.url).filter(Boolean)
  } catch {}
  detailPolling.value = false
}

function scheduleAutoPoll() {
  clearInterval(pollTimer)
  const active = genList.value.filter(r => r.status === 'PENDING' || r.status === 'RUNNING')
  if (active.length) {
    pollTimer = setInterval(async () => {
      for (const row of active) {
        try {
          const resp = await fetch('/api/tasks/' + row.task_id)
          const data = await resp.json()
          const s = data.output?.task_status
          if (s) {
            row.status = s
            if (s === 'SUCCEEDED' || s === 'FAILED') {
              row.result = JSON.stringify(data)
            }
          }
        } catch {}
      }
      if (!genList.value.some(r => r.status === 'PENDING' || r.status === 'RUNNING')) {
        clearInterval(pollTimer)
      }
    }, 15000)
  }
}

onMounted(() => { loadChat(); loadGen() })
onUnmounted(() => clearInterval(pollTimer))
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.detail-imgs { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 12px; }
pre { white-space: pre-wrap; font-size: 12px; max-height: 300px; overflow-y: auto; }
</style>