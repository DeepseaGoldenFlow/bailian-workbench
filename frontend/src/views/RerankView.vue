<template>
  <div class="page-container">
    <div class="page-header">
      <h2>🔀 文本排序</h2>
      <p class="page-sub">根据查询词对多个文档进行相关性排序</p>
    </div>

    <div class="glass-card" style="margin-bottom: 20px;">
      <h3 class="card-title"><el-icon><EditPen /></el-icon> 参数配置</h3>

      <div class="field">
        <label class="field-label">选择模型</label>
        <el-select v-model="form.model" class="full-width">
          <el-option label="gte-rerank (推荐)" value="gte-rerank" />
        </el-select>
      </div>

      <div class="field">
        <label class="field-label">查询词 (Query)</label>
        <el-input v-model="form.query" placeholder="输入搜索查询词，例如：如何学习 Python 编程" />
      </div>

      <div class="field">
        <div class="param-label-row">
          <span class="field-label">返回数量</span>
          <span class="param-value">Top {{ form.top_n }}</span>
        </div>
        <el-slider v-model="form.top_n" :min="1" :max="20" :step="1" :show-tooltip="false"
          :marks="{1:'1', 5:'5', 10:'10', 20:'20'}" />
      </div>
    </div>

    <!-- Documents Input -->
    <div class="glass-card" style="margin-bottom: 20px;">
      <div class="card-title-row">
        <h3 class="card-title"><el-icon><List /></el-icon> 文档列表</h3>
        <el-button @click="addDoc" type="primary" :icon="Plus" size="small">添加文档</el-button>
      </div>

      <div v-for="(doc, idx) in documents" :key="idx" class="doc-item">
        <span class="doc-index">{{ idx + 1 }}</span>
        <el-input v-model="doc.text" type="textarea" :rows="2" resize="none"
          :placeholder="`输入第 ${idx + 1} 篇文档内容...`" />
        <el-button @click="removeDoc(idx)" :icon="Delete" circle size="small" class="remove-btn" />
      </div>

      <el-button @click="rerank" type="primary" :loading="loading" class="generate-btn" size="large"
        :disabled="!form.query || documents.length < 2">
        <el-icon><Rank /></el-icon> 开始排序
      </el-button>
    </div>

    <!-- Results -->
    <div v-if="results.length > 0" class="results-section">
      <div class="glass-card">
        <h3 class="card-title"><el-icon><TrendCharts /></el-icon> 排序结果</h3>

        <div v-for="(item, idx) in results" :key="idx" class="rank-item"
          :style="{ animationDelay: idx * 0.1 + 's' }">
          <div class="rank-badge" :class="'rank-' + (idx + 1)">{{ idx + 1 }}</div>
          <div class="rank-content">
            <p class="rank-text">{{ item.text }}</p>
            <div class="rank-meta">
              <span class="rank-original">原始序号: #{{ item.originalIndex + 1 }}</span>
              <div class="score-bar">
                <span class="score-label">相关性</span>
                <div class="score-track">
                  <div class="score-fill" :style="{ width: (item.score * 100).toFixed(0) + '%' }"></div>
                </div>
                <span class="score-value">{{ (item.score * 100).toFixed(1) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-if="results.length === 0 && !loading" description="输入查询词和至少2篇文档后开始排序 🔀" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { EditPen, Plus, Delete, Rank, List, TrendCharts } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const form = ref({ model: 'gte-rerank', query: '', top_n: 10 })
const documents = ref([{ text: '' }, { text: '' }])
const loading = ref(false)
const results = ref([])

const addDoc = () => { documents.value.push({ text: '' }) }
const removeDoc = (idx) => { if (documents.value.length > 2) documents.value.splice(idx, 1) }

const rerank = async () => {
  if (!form.value.query) return ElMessage.warning('请输入查询词')
  const docs = documents.value.map(d => d.text).filter(t => t.trim())
  if (docs.length < 2) return ElMessage.warning('请至少输入2篇文档')

  loading.value = true
  results.value = []

  try {
    const res = await fetch('/api/compatible-mode/v1/rerank', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: form.value.model,
        query: form.value.query,
        documents: docs,
        top_n: form.value.top_n,
      })
    })

    if (!res.ok) { const e = await res.json(); throw new Error(e.message || '排序失败') }
    const data = await res.json()

    const reranked = data.output?.results || data.results || []
    results.value = reranked.map((r, i) => ({
      index: i,
      originalIndex: r.index,
      text: docs[r.index]?.slice(0, 300) || '',
      score: r.relevance_score ?? r.score ?? 0,
    })).sort((a, b) => b.score - a.score)

    ElMessage.success('🎉 排序完成！')
  } catch (e) {
    ElMessage.error('排序失败: ' + e.message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.page-container { max-width: 1000px; margin: 0 auto; width: 100%; }
.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--text-secondary); }

.glass-card { background: var(--card-bg); border: 1px solid var(--card-border); border-radius: 16px; padding: 20px; backdrop-filter: blur(10px); }
.card-title { display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600; margin-bottom: 16px; }
.card-title-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }

.field { margin-bottom: 16px; }
.field-label { display: block; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 6px; }
.param-label-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.param-value { font-size: 13px; font-weight: 700; color: var(--gradient-start); font-variant-numeric: tabular-nums; }
.full-width { width: 100%; }

.doc-item { display: flex; gap: 8px; margin-bottom: 12px; align-items: flex-start; }
.doc-index { width: 24px; height: 24px; border-radius: 6px; background: var(--card-hover); display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 700; flex-shrink: 0; margin-top: 4px; }
.remove-btn { flex-shrink: 0; margin-top: 4px; }

.generate-btn { width: 100%; margin-top: 8px; border-radius: 10px; }

.results-section { animation: msgIn 0.3s ease; }

.rank-item { display: flex; gap: 14px; padding: 14px; border-radius: 12px; margin-bottom: 10px; border: 1px solid var(--card-border); animation: rankIn 0.4s ease both; }
.rank-badge { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 800; flex-shrink: 0; }
.rank-1 { background: linear-gradient(135deg, #f59e0b, #eab308); color: #fff; }
.rank-2 { background: linear-gradient(135deg, #94a3b8, #64748b); color: #fff; }
.rank-3 { background: linear-gradient(135deg, #b45309, #92400e); color: #fff; }
.rank-badge:not(.rank-1):not(.rank-2):not(.rank-3) { background: var(--card-hover); color: var(--text-secondary); }

.rank-content { flex: 1; min-width: 0; }
.rank-text { font-size: 14px; line-height: 1.6; margin-bottom: 8px; word-break: break-word; }
.rank-meta { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.rank-original { font-size: 12px; color: var(--text-secondary); }

.score-bar { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 150px; }
.score-label { font-size: 11px; color: var(--text-secondary); white-space: nowrap; }
.score-track { flex: 1; height: 6px; background: var(--card-hover); border-radius: 3px; overflow: hidden; }
.score-fill { height: 100%; background: linear-gradient(90deg, var(--gradient-start), var(--gradient-end)); border-radius: 3px; transition: width 0.6s ease; }
.score-value { font-size: 12px; font-weight: 700; color: var(--gradient-start); white-space: nowrap; }

@keyframes rankIn { from { opacity: 0; transform: translateX(-20px); } to { opacity: 1; transform: translateX(0); } }
@keyframes msgIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
