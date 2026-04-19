<template>
  <div class="page-container">
    <div class="page-header">
      <h2>📊 文本向量</h2>
      <p class="page-sub">将文本转换为向量表示，用于相似度计算和语义搜索</p>
    </div>

    <div class="glass-card" style="margin-bottom: 20px;">
      <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

      <div class="config-row">
        <div class="field" style="flex: 1; margin-right: 16px;">
          <label class="field-label">选择模型</label>
          <el-select v-model="form.model" class="full-width">
            <el-option label="text-embedding-v3 (推荐)" value="text-embedding-v3" />
            <el-option label="text-embedding-v2" value="text-embedding-v2" />
          </el-select>
        </div>

        <div class="field" style="flex: 1; margin-right: 16px;">
          <label class="field-label">向量类型</label>
          <el-radio-group v-model="form.vector_type">
            <el-radio-button value="query">查询文本 (Query)</el-radio-button>
            <el-radio-button value="document">文档文本 (Document)</el-radio-button>
          </el-radio-group>
          <p class="field-desc">Query 用于搜索词，Document 用于被搜索的文档内容</p>
        </div>
      </div>
    </div>

    <!-- Input Texts -->
    <div class="glass-card" style="margin-bottom: 20px;">
      <div class="card-title-row">
        <h3 class="card-title"><el-icon><List /></el-icon> 输入文本列表</h3>
        <el-button @click="addText" type="primary" :icon="Plus" size="small">添加文本</el-button>
      </div>

      <div v-for="(item, idx) in textItems" :key="idx" class="text-item">
        <el-input v-model="item.text" type="textarea" :rows="2" resize="none"
          :placeholder="`输入第 ${idx + 1} 段文本...`" />
        <el-button @click="removeText(idx)" :icon="Delete" circle size="small" class="remove-btn" />
      </div>

      <el-button @click="computeEmbeddings" type="primary" :loading="loading" class="generate-btn" size="large">
        <el-icon><Connection /></el-icon> 计算向量
      </el-button>
    </div>

    <!-- Results -->
    <div v-if="results.length > 0" class="results-section">
      <div class="glass-card" style="margin-bottom: 20px;">
        <h3 class="card-title"><el-icon><DataAnalysis /></el-icon> 向量结果</h3>
        <el-table :data="results" stripe>
          <el-table-column label="#" width="60" align="center">
            <template #default="{ $index }">{{ $index + 1 }}</template>
          </el-table-column>
          <el-table-column label="原文" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">{{ row.text }}</template>
          </el-table-column>
          <el-table-column label="向量维度" width="120" align="center">
            <template #default="{ row }">{{ row.dimension }}</template>
          </el-table-column>
          <el-table-column label="向量预览" min-width="250">
            <template #default="{ row }">
              <span class="vector-preview">[{{ row.preview }}...]</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80" align="center">
            <template #default="{ row }">
              <el-button link type="primary" @click="copyVector(row)" size="small"><el-icon><CopyDocument /></el-icon></el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Similarity Matrix -->
      <div v-if="results.length >= 2" class="glass-card">
        <h3 class="card-title"><el-icon><Grid /></el-icon> 相似度矩阵 (余弦相似度)</h3>
        <div class="similarity-grid">
          <div class="sim-header" v-for="(r, ci) in results" :key="'h'+ci">
            <span class="sim-label">{{ ci + 1 }}</span>
          </div>
          <template v-for="(row, ri) in similarityMatrix" :key="'r'+ri">
            <div class="sim-label-cell"><span>{{ ri + 1 }}</span></div>
            <div v-for="(val, ci) in row" :key="'c'+ci" class="sim-cell"
              :style="{ background: getSimColor(val) }">
              <span class="sim-value">{{ val.toFixed(3) }}</span>
            </div>
          </template>
        </div>
      </div>
    </div>

    <el-empty v-if="results.length === 0 && !loading" description="输入文本后计算向量 📊" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { EditPen, Plus, Delete, Connection, CopyDocument, List, DataAnalysis, Grid } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({ model: 'text-embedding-v3', vector_type: 'query' })
const textItems = ref([{ text: '' }, { text: '' }])
const loading = ref(false)
const results = ref([])
const similarityMatrix = ref([])

const addText = () => { textItems.value.push({ text: '' }) }
const removeText = (idx) => { if (textItems.value.length > 1) textItems.value.splice(idx, 1) }

const getSimColor = (val) => {
  const alpha = val * 0.6 + 0.1
  return `rgba(99, 102, 241, ${alpha})`
}

const computeSimilarity = (v1, v2) => {
  let dot = 0, n1 = 0, n2 = 0
  for (let i = 0; i < v1.length; i++) {
    dot += v1[i] * v2[i]
    n1 += v1[i] * v1[i]
    n2 += v2[i] * v2[i]
  }
  return dot / (Math.sqrt(n1) * Math.sqrt(n2) || 1)
}

const computeEmbeddings = async () => {
  const texts = textItems.value.map(t => t.text).filter(t => t.trim())
  if (texts.length === 0) return ElMessage.warning('请至少输入一段文本')

  loading.value = true
  results.value = []
  similarityMatrix.value = []

  try {
    // Compute embeddings for each text
    const vectors = []
    for (const text of texts) {
      const res = await fetch('/api/compatible-mode/v1/embeddings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ model: form.value.model, input: text })
      })

      if (!res.ok) { const e = await res.json(); throw new Error(e.message || '计算失败') }
      const data = await res.json()
      const embedding = data.data?.[0]?.embedding
      if (!embedding) throw new Error('未获取到向量数据')
      vectors.push(embedding)

      results.value.push({
        text: text.slice(0, 100),
        dimension: embedding.length,
        vector: embedding,
        preview: embedding.slice(0, 6).map(v => v.toFixed(4)).join(', '),
      })
    }

    // Compute similarity matrix
    const n = vectors.length
    similarityMatrix.value = Array.from({ length: n }, (_, i) =>
      Array.from({ length: n }, (_, j) => computeSimilarity(vectors[i], vectors[j]))
    )

    ElMessage.success('🎉 向量计算完成！')
  } catch (e) {
    ElMessage.error('计算失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

const copyVector = async (row) => {
  try {
    await navigator.clipboard.writeText(JSON.stringify(row.vector))
    ElMessage.success('已复制向量数据')
  } catch { ElMessage.error('复制失败') }
}
</script>

<style scoped>
.page-container { max-width: 1200px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.glass-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; padding: 20px; backdrop-filter: blur(10px); }
.card-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }
.card-title-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }

.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.field-desc { font-size: 11px; color: var(--text-secondary); margin-top: 4px; }
.full-width { width: 100%; }
.config-row { display: flex; gap: 16px; }

.text-item { display: flex; gap: 8px; margin-bottom: 12px; align-items: flex-start; }
.remove-btn { flex-shrink: 0; margin-top: 4px; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.results-section { animation: msgIn 0.3s ease; }

.vector-preview { font-family: 'Fira Code', monospace; font-size: 12px; color: var(--gradient-start); }

.similarity-grid { display: grid; gap: 2px; }
.sim-header { text-align: center; padding: 8px; }
.sim-label { font-size: 13px; font-weight: 700; color: var(--text-secondary); }
.sim-label-cell { display: flex; align-items: center; justify-content: center; padding: 8px; }
.sim-cell { border-radius: 6px; padding: 12px; text-align: center; transition: all 0.2s; }
.sim-cell:hover { transform: scale(1.05); }
.sim-value { font-size: 13px; font-weight: 600; color: #fff; font-variant-numeric: tabular-nums; }

@media (max-width: 768px) { .config-row { flex-direction: column; } }
@keyframes msgIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
