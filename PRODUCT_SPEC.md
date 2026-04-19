# 百炼工作台 2.0 — 完整产品文档

## 一、产品定位

**一站式 AI 内容创作平台** —— 本地化部署，整合百炼全量 AIGC 能力，覆盖图像、视频、音频、对话四大领域。

---

## 二、功能全景

### 🎨 图像（Image）

| 子功能 | 说明 | 模型 | 输入方式 |
|--------|------|------|---------|
| **文生图** | 文字描述生成图片 | wan2.7-image, wan2.7-image-pro | Prompt |
| **图生图** | 参考图为基础生成新图 | wanx2.1-sketch2image-plus | 参考图 + Prompt |
| **局部重绘** | 选中区域重新生成 | wanx2.1-image-edit-plus | 原图 + 蒙版 + Prompt |
| **扩图** | 向外扩展画面边界 | wanx2.1-image-edit-plus | 原图 + 扩展方向/比例 |
| **图像擦除** | 去除画面中不需要的元素 | wanx2.1-image-edit-plus | 原图 + 擦除区域 |
| **涂鸦生图** | 手绘草图生成精致图片 | wanx2.1-sketch2image-plus | 草图 + Prompt |
| **图像风格迁移** | 将图片转换为不同风格 | wanx-style-transfer | 原图 + 风格 |
| **图像超分放大** | 提升图片分辨率 | wanx-upscale | 原图 |
| **抠图/分割** | 自动抠出主体/分割区域 | wanx-segmentation | 原图 |
| **人像写真** | 生成/编辑人像 | wanx-portrait | 原图 + 参数 |
| **虚拟模特** | 电商商品图生成 | wanx-virtual-model | 商品图 |
| **背景替换** | 一键更换背景 | wanx-background | 原图 + 新背景描述 |

### 🎬 视频（Video）

| 子功能 | 说明 | 模型 | 输入方式 |
|--------|------|------|---------|
| **文生视频** | 文字描述生成视频 | wan2.6-t2v-turbo, wan2.6-t2v-plus | Prompt |
| **图生视频** | 单张参考图驱动生成视频 | wan2.6-i2v-turbo, wan2.6-i2v-plus | 首帧图 + Prompt |
| **首尾帧视频** | 上传首帧和尾帧生成过渡动画 | wan2.6-i2v-turbo | 首帧图 + 尾帧图 |
| **多参考图视频** | 多张参考图保持角色一致性 | wan2.6-i2v-plus | 多张参考图 + Prompt |
| **视频风格迁移** | 视频风格转换 | wanx-video-style | 原视频 + 风格 |
| **视频超分** | 提升视频分辨率 | wanx-video-upsacle | 原视频 |
| **视频编辑** | 视频内容编辑 | wanx-video-edit | 原视频 + 编辑指令 |
| **视频延长** | 延长已有视频 | wanx-video-extend | 原视频 + 时长 |

### 🎙️ 音频（Audio）

| 子功能 | 说明 | 模型 | 输入方式 |
|--------|------|------|---------|
| **语音合成 (TTS)** | 文字转语音 | qwen3-tts-flash, qwen3-tts-flash-realtime | 文本 |
| **语音克隆** | 克隆指定声音生成语音 | qwen3-tts-vc | 样本音频 + 文本 |
| **语音识别 (ASR)** | 语音转文字 | qwen3-asr-flash, qwen3-asr-flash-realtime | 音频文件 |
| **实时语音识别** | 流式识别 | qwen3-asr-flash-realtime | 音频流 |
| **音乐生成** | 文字生成音乐 | Suno/AI 音乐 | 文本描述 |

### 💬 对话（Chat）

| 子功能 | 说明 | 模型 |
|--------|------|------|
| **AI 对话** | 多轮智能对话 | qwen3.6-plus, qwen3.6-flash, qwen3.5-plus, qwen3.5-flash, qwen3-max, qwen-flash |
| **深度研究** | 自动搜索+深度分析 | qwen-deep-research |
| **代码生成** | 智能代码补全/生成 | qwen3-coder-plus, qwen3-coder-flash |

### 👁️ 视觉理解（Vision）

| 子功能 | 说明 | 模型 |
|--------|------|------|
| **图片问答** | 看图回答问题 | qwen3-vl-plus, qwen3-vl-flash |
| **OCR** | 图片文字提取 | qwen-vl-ocr |
| **多模态对话** | 图文对话 | qwen3-vl-plus |

### 🔢 向量 & 📊 排序

| 功能 | 模型 | 说明 |
|------|------|------|
| 文本向量化 | text-embedding-v3 | 语义搜索、相似度计算 |
| 文本排序 | gte-rerank | 搜索结果精排 |

---

## 三、页面设计

### 整体布局
```
┌── 侧边栏（可折叠）──┬──────────────────────────────────────────┐
│ ✦ 百炼工作台       │  当前页面标题        [🌙] [⚙️设置] [👤]  │
│                    ├──────────────────────────────────────────┤
│ 🎨 图像生成        │                                          │
│   ↳ 文生图         │  ┌─ 模型选择 ──────────────────────────┐  │
│   ↳ 图生图         │  │ [wan2.7-image ▾] [wan2.7-image-pro] │  │
│   ↳ 图像编辑       │  └─────────────────────────────────────┘  │
│   ↳ 局部重绘       │                                          │
│   ↳ 扩图           │  ┌─ 输入区 ──────────────────────────┐    │
│   ↳ 擦除           │  │  Prompt: [                       ] │    │
│   ↳ 涂鸦生图       │  │  [📤 上传图片]                     │    │
│   ↳ 风格迁移       │  │  [🎛️ 蒙版绘制]                    │    │
│   ↳ 超分放大       │  └─────────────────────────────────────┘  │
│   ↳ 抠图           │                                          │
│   ↳ 人像写真       │  ┌─ 参数面板 ────────────────────────┐    │
│   ↳ 虚拟模特       │  │  [📐 比例] [🎲 数量] [✨ 优化]    │    │
│   ↳ 背景替换       │  └─────────────────────────────────────┘  │
│ 🎬 视频生成        │                                          │
│   ↳ 文生视频       │         [ 🚀 开始生成 ]                    │
│   ↳ 图生视频       │                                          │
│   ↳ 首尾帧         │  ┌─ 结果展示 ────────────────────────┐    │
│   ↳ 多参考图       │  │  [▶️ 生成结果]                     │    │
│   ↳ 视频编辑       │  └─────────────────────────────────────┘  │
│   ↳ 视频延长       │                                          │
│ 🎙️ 音频处理        │  ┌─── 历史记录 ──────────────────────┐    │
│   ↳ 语音合成       │  │  [img] [img] [img] [img] ...       │    │
│   ↳ 语音识别       │  └─────────────────────────────────────┘  │
│   ↳ 语音克隆       │                                          │
│   ↳ 音乐生成       │                                          │
│ 💬 AI 对话         │                                          │
│ 👁️ 视觉理解        │                                          │
│ 🔢 向量化          │                                          │
│ 📊 排序            │                                          │
│ 📋 任务历史        │                                          │
│ ⚙️ 设置            │                                          │
└────────────────────┴──────────────────────────────────────────┘
```

### 关键设计原则

1. **左侧树形导航**：每个大类可展开，子功能直接跳转，避免顶部 Tab 拥挤
2. **每个子功能独立页面**：不复用一个页面切换模式，而是每个子功能有自己的 URL 和布局
3. **输入区自适应**：根据子功能类型自动显示对应的输入控件
   - 文生图：只有 Prompt
   - 图生图：Prompt + 参考图上传
   - 局部重绘：原图上传 + 蒙版绘制 + Prompt
   - 首尾帧：首帧上传 + 尾帧上传
4. **参数面板可折叠**，每个参数用中文标签 + 说明，不用裸参数名
5. **任务历史全局可见**：每个页面底部 + 独立任务历史页面

---

## 四、模型与 API 映射

### 图像 API 端点

| 子功能 | API 端点 | 模型 |
|--------|---------|------|
| 文生图 | POST /api/v1/services/aigc/text2image/image-synthesis | wan2.7-image, wan2.7-image-pro |
| 图生图 | POST /api/v1/services/aigc/image2image/image-synthesis | wanx2.1-image-edit-plus |
| 局部重绘 | POST /api/v1/services/aigc/image2image/image-synthesis (带 mask) | wanx2.1-image-edit-plus |
| 扩图 | POST /api/v1/services/aigc/image2image/image-synthesis (outpaint) | wanx2.1-image-edit-plus |
| 擦除 | POST /api/v1/services/aigc/image2image/erase | wanx2.1-image-edit-plus |
| 涂鸦生图 | POST /api/v1/services/aigc/image2image/sketch2image | wanx2.1-sketch2image-plus |
| 风格迁移 | POST /api/v1/services/aigc/image2image/style-transfer | wanx-style-transfer |
| 超分放大 | POST /api/v1/services/aigc/image-super-resolution | wanx-upscale |
| 抠图 | POST /api/v1/services/aigc/segmentation/semantic-segmentation | wanx-segmentation |

### 视频 API 端点

| 子功能 | API 端点 | 模型 |
|--------|---------|------|
| 文生视频 | POST /api/v1/services/aigc/text2video/video-synthesis | wan2.6-t2v-turbo, wan2.6-t2v-plus |
| 图生视频 | POST /api/v1/services/aigc/image2video/video-synthesis | wan2.6-i2v-turbo, wan2.6-i2v-plus |
| 首尾帧 | POST /api/v1/services/aigc/image2video/video-synthesis (first/last frame) | wan2.6-i2v-turbo |
| 多参考图 | POST /api/v1/services/aigc/image2video/video-synthesis (multi-ref) | wan2.6-i2v-plus |
| 视频风格迁移 | POST /api/v1/services/aigc/video-generation/video-style-transfer | wanx-video-style |
| 视频超分 | POST /api/v1/services/aigc/video-super-resolution | wanx-video-upscale |

### 音频 API 端点

| 子功能 | API 端点 | 模型 |
|--------|---------|------|
| TTS | POST /api/v1/services/aigc/text2audio/speech-synthesis | qwen3-tts-flash |
| TTS 实时 | POST /api/v1/services/aigc/text2audio/speech-synthesis (stream) | qwen3-tts-flash-realtime |
| 语音克隆 | POST /api/v1/services/aigc/text2audio/speech-synthesis (voice-clone) | qwen3-tts-vc |
| ASR | POST /api/v1/services/aicp/paraformer/async-translation | qwen3-asr-flash |
| 音乐生成 | POST /api/v1/services/aigc/text2audio/music-generation | AI 音乐 |

---

## 五、参数设计模板

所有参数面板遵循统一规范：

### 图像参数
- **画面比例**：按钮组 [1:1] [16:9] [9:16] [4:3] [3:4] [自定义]
- **生成数量**：滑块 1-4
- **提示词优化**：开关
- **随机种子**：数字输入（可选）
- **负向提示词**：文本域（不想出现的元素）

### 图生图特有参数
- **参考图强度**：滑块 0-1（参考图对结果的影响程度）
- **参考图**：拖拽上传

### 局部重绘特有参数
- **蒙版**：图片上绘制/选择区域
- **蒙版模式**：[重绘选中区域] / [重绘非选中区域]

### 扩图特有参数
- **扩展方向**：[上] [下] [左] [右] [四周]
- **扩展比例**：滑块 1.0-2.0

### 视频参数
- **视频时长**：按钮组 [5秒] [10秒] [15秒]
- **分辨率**：下拉 [720P] [1080P]
- **运动幅度**：滑块 0-1（动态感强弱）
- **提示词优化**：开关
- **随机种子**：数字输入

### TTS 参数
- **音色**：下拉（中文男声/中文女声/英文男声/英文女声/童声等）
- **语速**：滑块 0.5x-2.0x
- **采样率**：下拉 [22050] [44100] [48000]
- **格式**：下拉 [WAV] [MP3]

### ASR 参数
- **采样率**：下拉 [8000] [16000] [44100] [48000]
- **语言**：下拉 [中文] [英文] [多语种]
