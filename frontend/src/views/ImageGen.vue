<template>
  <div class="page">
    <h2>图片生成</h2>

    <el-form label-width="100px">
      <el-form-item label="模型">
        <el-select v-model="model" placeholder="选择模型" @change="onModelChange" style="width:100%">
          <el-option v-for="m in models" :key="m.id" :label="m.name" :value="m.id">
            <span>{{ m.name }}</span>
            <span style="color:#909399;font-size:12px;margin-left:8px">{{ m.async ? '异步' : '同步' }}</span>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="提示词">
        <el-input v-model="prompt" type="textarea" :rows="3" placeholder="描述要生成的图像..." />
      </el-form-item>

      <el-form-item v-if="hasParam('negative_prompt')" label="反向提示词">
        <el-input v-model="negativePrompt" placeholder="不希望出现的内容" />
      </el-form-item>

      <el-form-item label="尺寸">
        <el-select v-model="size" v-if="hasParam('size')">
          <el-option v-for="o in getOptions('size')" :key="o.value" :label="o.label" :value="o.value" />
        </el-select>
        <span v-else>1024×1024</span>
      </el-form-item>

      <el-form-item label="数量">
        <el-slider v-model="n" :min="1" :max="4" show-input style="width:200px" />
      </el-form-item>

      <el-form-item v-if="hasParam('seed')" label="随机种子">
        <el-input-number v-model="seed" :min="0" :max="2147483647" placeholder="随机" style="width:200px" />
      </el-form-item>

      <el-form-item v-if="hasParam('steps')" label="推理步数">
        <el-slider v-model="steps" :min="1" :max="100" show-input style="width:200px" />
      </el-form-item>

      <el-form-item v-if="hasParam('prompt_extend')" label="提示词扩展">
        <el-switch v-model="promptExtend" />
      </el-form-item>

      <el-form-item v-if="hasParam('ref_img')" label="参考图">
        <el-input v-model="refImg" placeholder="参考图URL（可选）" />
      </el-form-item>

      <el-form-item v-if="hasParam('ref_strength') && refImg" label="参考强度">
        <el-slider v-model="refStrength" :min="0" :max="1" :step="0.1" show-input style="width:200px" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="generate" :loading="loading" :disabled="!prompt || !model">
          {{ currentModel?.async ? '提交任务' : '生成图片' }}
        </el-button>
      </el-form-item>
    </el-form>

    <div v-if="taskId" class="task-info">
      <el-alert :title="'任务: ' + taskId" :type="taskStatus === 'SUCCEEDED' ? 'success' : 'info'" show-icon :closable="false">
        <template v-if="taskStatus === 'PENDING' || taskStatus === 'RUNNING'">
          <el-progress :percentage="100" :indeterminate="true" :duration="2" />
          <el-button size="small" @click="pollTaskResult" :loading="polling">刷新状态</el-button>
        </template>
        <template v-if="taskStatus === 'SUCCEEDED'">
          <span>生成完成！</span>
        </template>
        <template v-if="taskStatus === 'FAILED'">
          <span style="color:red">任务失败</span>
        </template>
      </el-alert>
    </div>

    <div v-if="images.length" class="result-grid">
      <div v-for="(img, i) in images" :key="i" class="result-item">
        <el-image :src="img" fit="contain" :preview-src-list="images" :initial-index="i" style="width:100%;max-height:400px" />
      </div>
    </div>

    <div v-if="error" style="margin-top:12px">
      <el-alert :title="error" type="error" show-icon />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchModels, imageGen, pollTask } from '../api'

const models = ref([])
const model = ref('')
const prompt = ref('')
const negativePrompt = ref('')
const size = ref('1024*1024')
const n = ref(1)
const seed = ref(null)
const steps = ref(30)
const promptExtend = ref(true)
const refImg = ref('')
const refStrength = ref(0.5)
const loading = ref(false)
const taskId = ref('')
const taskStatus = ref('')
const polling = ref(false)
const images = ref([])
const error = ref('')

const currentModel = computed(() => models.value.find(m => m.id === model.value))

onMounted(async () => {
  try {
    const { data } = await fetchModels('image')
    models.value = data.models || []
    if (models.value.length) model.value = models.value[0].id
  } catch (e) {
    error.value = '加载模型列表失败: ' + (e.response?.data?.message || e.message)
  }
})

function hasParam(name) {
  return currentModel.value?.parameters?.some(p => p.name === name)
}

function getOptions(name) {
  return currentModel.value?.parameters?.find(p => p.name === name)?.options || []
}

function onModelChange() {
  taskId.value = ''
  taskStatus.value = ''
  images.value = []
  error.value = ''
}

async function generate() {
  loading.value = true
  error.value = ''
  taskId.value = ''
  taskStatus.value = ''
  images.value = []
  try {
    const payload = {
      model: model.value,
      prompt: prompt.value,
      size: size.value,
      n: n.value,
    }
    if (seed.value != null) payload.seed = seed.value
    if (hasParam('steps')) payload.steps = steps.value
    if (hasParam('prompt_extend')) payload.prompt_extend = promptExtend.value
    if (hasParam('negative_prompt') && negativePrompt.value) payload.negative_prompt = negativePrompt.value
    if (hasParam('ref_img') && refImg.value) {
      payload.ref_img = refImg.value
      if (refStrength.value) payload.ref_strength = refStrength.value
    }

    const { data } = await imageGen(payload)
    if (data.task_id) {
      taskId.value = data.task_id
      taskStatus.value = data.status
      if (data.status === 'PENDING' || data.status === 'RUNNING') {
        await pollTaskResult()
      }
    } else {
      extractImages(data)
    }
  } catch (e) {
    error.value = e.response?.data?.message || e.message
  }
  loading.value = false
}

async function pollTaskResult() {
  if (!taskId.value) return
  polling.value = true
  try {
    const { data } = await pollTask(taskId.value)
    taskStatus.value = data.output?.task_status || 'UNKNOWN'
    if (taskStatus.value === 'SUCCEEDED') {
      const results = data.output?.results || []
      images.value = results.map(r => r.url).filter(Boolean)
      if (!images.value.length && data.output?.video_url) {
        images.value = [data.output.video_url]
      }
    }
  } catch (e) {
    error.value = '轮询失败: ' + (e.response?.data?.message || e.message)
  }
  polling.value = false
}

function extractImages(data) {
  if (data.output?.results) {
    images.value = data.output.results.map(r => r.url).filter(Boolean)
  } else if (data.output?.choices) {
    images.value = data.output.choices
      .filter(c => c.message?.content)
      .flatMap(c => {
        if (Array.isArray(c.message.content)) {
          return c.message.content.filter(x => x.image).map(x => x.image)
        }
        return []
      })
  }
}
</script>

<style scoped>
.page { max-width: 900px; margin: 0 auto; }
.result-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 16px; margin-top: 20px; }
.result-item { border: 1px solid #e0e0e0; border-radius: 8px; overflow: hidden; padding: 8px; background: #fff; }
.task-info { margin-top: 16px; }
</style>
