<template>
  <div class="page">
    <h2>视频生成</h2>
    <el-form label-width="100px">
      <el-form-item label="模型">
        <el-select v-model="model" placeholder="选择模型" @change="onModelChange" style="width:100%">
          <el-option v-for="m in models" :key="m.id" :label="m.name" :value="m.id">
            <span>{{ m.name }}</span>
            <span style="color:#909399;font-size:12px;margin-left:8px">{{ m.description }}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="提示词">
        <el-input v-model="prompt" type="textarea" :rows="3"
          :placeholder="isEdit ? '描述编辑意图' : '描述视频内容...'" />
      </el-form-item>
      <el-form-item v-if="isEdit" label="视频URL">
        <el-input v-model="videoUrl" placeholder="待编辑的视频URL（MP4/MOV）" />
      </el-form-item>
      <el-form-item v-if="isI2V" label="首帧图片">
        <el-input v-model="firstFrame" placeholder="作为视频首帧的图片URL" />
      </el-form-item>
      <el-form-item v-if="isR2V || isEdit" :label="'参考图片' + (isR2V ? ' (1-9张)' : ' (0-5张)')">
        <el-input v-model="refImagesStr" type="textarea" :rows="2" placeholder="图片URL，多个用逗号或换行分隔" />
      </el-form-item>
      <el-form-item label="分辨率">
        <el-select v-model="resolution">
          <el-option label="1080P" value="1080P" />
          <el-option label="720P" value="720P" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="!isEdit" label="时长(秒)">
        <el-slider v-model="duration" :min="3" :max="15" show-input style="width:250px" />
      </el-form-item>
      <el-form-item v-if="!isEdit" label="宽高比">
        <el-select v-model="ratio">
          <el-option v-for="r in ratios" :key="r" :label="r" :value="r" />
        </el-select>
      </el-form-item>
      <el-form-item label="随机种子">
        <el-input-number v-model="seed" :min="0" :max="2147483647" placeholder="随机" style="width:200px" />
      </el-form-item>
      <el-form-item label="水印">
        <el-switch v-model="watermark" />
        <span style="margin-left:8px;color:#909399;font-size:12px">右下角 Happy Horse 水印</span>
      </el-form-item>
      <el-form-item label="声音">
        <el-select v-model="audioSetting">
          <el-option label="自动" value="auto" />
          <el-option label="保留原声" value="origin" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="generate" :loading="loading" :disabled="!canSubmit">提交生成任务</el-button>
      </el-form-item>
    </el-form>
    <div v-if="taskId" class="task-info">
      <el-alert :title="'任务: ' + taskId" :type="taskStatus === 'SUCCEEDED' ? 'success' : 'info'" show-icon :closable="false">
        <template v-if="taskStatus === 'PENDING' || taskStatus === 'RUNNING'">
          <el-progress :percentage="100" :indeterminate="true" :duration="2" />
          <el-button size="small" @click="pollTaskResult" :loading="polling" style="margin-top:8px">刷新状态</el-button>
        </template>
        <template v-if="taskStatus === 'SUCCEEDED'"><span>生成完成！</span></template>
        <template v-if="taskStatus === 'FAILED'"><span style="color:red">任务失败</span></template>
      </el-alert>
    </div>
    <div v-if="resultUrl" class="result">
      <h3>生成结果</h3>
      <video :src="resultUrl" controls style="max-width:100%;max-height:500px" />
      <p style="margin-top:8px">
        <el-link :href="resultUrl" target="_blank" type="primary">下载视频</el-link>
        <span style="color:#909399;font-size:12px;margin-left:8px">链接24小时有效</span>
      </p>
    </div>
    <div v-if="error" style="margin-top:12px"><el-alert :title="error" type="error" show-icon /></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchModels, videoGen, pollTask } from '../api'

const models = ref([])
const model = ref('')
const prompt = ref('')
const videoUrl = ref('')
const firstFrame = ref('')
const refImagesStr = ref('')
const resolution = ref('1080P')
const duration = ref(5)
const ratio = ref('16:9')
const seed = ref(null)
const watermark = ref(true)
const audioSetting = ref('auto')
const loading = ref(false)
const taskId = ref('')
const taskStatus = ref('')
const taskError = ref('')
const polling = ref(false)
const resultUrl = ref('')
const error = ref('')

const ratios = ['16:9', '9:16', '1:1', '4:3', '3:4']
const isEdit = computed(() => model.value === 'happyhorse-1.0-video-edit')
const isI2V = computed(() => model.value === 'happyhorse-1.0-i2v')
const isR2V = computed(() => model.value === 'happyhorse-1.0-r2v')
const isT2V = computed(() => model.value === 'happyhorse-1.0-t2v')
const canSubmit = computed(() => {
  if (!model.value || !prompt.value) return false
  if (isEdit.value && !videoUrl.value) return false
  if (isI2V.value && !firstFrame.value) return false
  if (isR2V.value && !refImagesStr.value) return false
  return true
})

onMounted(async () => {
  try {
    const { data } = await fetchModels('video')
    models.value = data.models || []
    if (models.value.length) model.value = models.value[0].id
  } catch (e) { error.value = '加载模型列表失败' }
})

function onModelChange() {
  taskId.value = ''; taskStatus.value = ''; resultUrl.value = ''; error.value = ''
}

function parseRefImages() {
  if (!refImagesStr.value) return []
  return refImagesStr.value.replace(/\n/g, ",").split(",").map(s => s.trim()).filter(Boolean)
}

async function generate() {
  loading.value = true; error.value = ''; taskId.value = ''; taskStatus.value = ''; resultUrl.value = ''
  try {
    const payload = {
      model: model.value, prompt: prompt.value,
      resolution: resolution.value,
      seed: seed.value != null ? seed.value : undefined,
      watermark: watermark.value, audio_setting: audioSetting.value,
    }
    if (!isEdit.value) { payload.duration = duration.value; payload.ratio = ratio.value }
    if (isEdit.value) { payload.video_url = videoUrl.value; payload.ref_images = parseRefImages() }
    if (isI2V.value) payload.first_frame = firstFrame.value
    if (isR2V.value) payload.ref_images = parseRefImages()
    const { data } = await videoGen(payload)
    taskId.value = data.task_id; taskStatus.value = data.status
    if (data.status === 'PENDING' || data.status === 'RUNNING') await pollTaskResult()
  } catch (e) { error.value = e.response?.data?.message || e.message }
  loading.value = false
}

async function pollTaskResult() {
  if (!taskId.value) return
  polling.value = true
  try {
    const { data } = await pollTask(taskId.value)
    taskStatus.value = data.output?.task_status || 'UNKNOWN'
    if (taskStatus.value === 'SUCCEEDED') resultUrl.value = data.output?.video_url || ''
    if (taskStatus.value === 'FAILED') taskError.value = data.output?.message || '未知错误'
  } catch (e) { error.value = '轮询失败' }
  polling.value = false
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.result { margin-top: 20px; padding: 16px; background: #f5f7fa; border-radius: 8px; }
.task-info { margin-top: 16px; }
</style>