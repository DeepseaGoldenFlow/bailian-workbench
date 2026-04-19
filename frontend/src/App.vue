<template>
  <div class="app-root">
    <aside class="sidebar" :class="{ collapsed: isCollapsed }">
      <div class="sidebar-header">
        <div class="logo-area">
          <div class="logo-icon">✦</div>
          <span class="logo-text" v-show="!isCollapsed">百炼工作台</span>
        </div>
        <el-icon class="collapse-btn" @click="isCollapsed = !isCollapsed">
          <Fold v-if="!isCollapsed" /><Expand v-else />
        </el-icon>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapsed"
        :collapse-transition="false"
        class="side-menu"
        background-color="transparent"
        text-color="var(--text-secondary)"
        active-text-color="var(--gradient-start)"
        @select="handleMenuSelect"
      >
        <!-- 智能对话 -->
        <el-sub-menu index="chat">
          <template #title>
            <el-icon><ChatDotRound /></el-icon>
            <span>智能对话</span>
          </template>
          <el-menu-item index="/chat">
            <el-icon><ChatLineSquare /></el-icon>
            <span>AI 对话</span>
          </el-menu-item>
          <el-menu-item index="/vision">
            <el-icon><View /></el-icon>
            <span>视觉理解</span>
          </el-menu-item>
        </el-sub-menu>

        <!-- 图像工厂 -->
        <el-sub-menu index="image">
          <template #title>
            <el-icon><Picture /></el-icon>
            <span>图像工厂</span>
          </template>
          <el-menu-item index="/image/gen">
            <el-icon><MagicStick /></el-icon>
            <span>AI 生图</span>
          </el-menu-item>
          <el-menu-item index="/image/edit">
            <el-icon><Edit /></el-icon>
            <span>图像编辑</span>
          </el-menu-item>
          <el-menu-item index="/image/tryon">
            <el-icon><Crop /></el-icon>
            <span>虚拟试衣</span>
          </el-menu-item>
        </el-sub-menu>

        <!-- 视频工厂 -->
        <el-sub-menu index="video">
          <template #title>
            <el-icon><VideoCamera /></el-icon>
            <span>视频工厂</span>
          </template>
          <el-menu-item index="/video/t2v">
            <el-icon><Promotion /></el-icon>
            <span>文生视频</span>
          </el-menu-item>
          <el-menu-item index="/video/i2v">
            <el-icon><PictureFilled /></el-icon>
            <span>图生视频</span>
          </el-menu-item>
          <el-menu-item index="/video/ref">
            <el-icon><Connection /></el-icon>
            <span>参考生视频</span>
          </el-menu-item>
          <el-menu-item index="/video/edit">
            <el-icon><Film /></el-icon>
            <span>视频编辑</span>
          </el-menu-item>
          <el-menu-item index="/video/digital-human">
            <el-icon><UserFilled /></el-icon>
            <span>数字人</span>
          </el-menu-item>
          <el-menu-item index="/video/animate">
            <el-icon><Sunny /></el-icon>
            <span>动画动效</span>
          </el-menu-item>
        </el-sub-menu>

        <!-- 音频中心 -->
        <el-sub-menu index="audio">
          <template #title>
            <el-icon><Headset /></el-icon>
            <span>音频中心</span>
          </template>
          <el-menu-item index="/audio/tts">
            <el-icon><Microphone /></el-icon>
            <span>语音合成</span>
          </el-menu-item>
          <el-menu-item index="/audio/asr">
            <el-icon><Monitor /></el-icon>
            <span>语音识别</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>

      <div class="sidebar-footer">
        <el-tooltip :content="isDark ? '切换亮色模式' : '切换暗色模式'" placement="right">
          <el-icon class="theme-toggle" @click="toggleTheme">
            <Sunny v-if="isDark" /><Moon v-else />
          </el-icon>
        </el-tooltip>
      </div>
    </aside>

    <main class="main-content">
      <div class="content-gradient"></div>
      <div class="content-inner">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Fold, Expand, Sunny, Moon, ChatDotRound, ChatLineSquare, View,
  Picture, MagicStick, Edit, Crop,
  VideoCamera, Promotion, PictureFilled, Connection, Film, UserFilled,
  Headset, Microphone, Monitor
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const isDark = ref(localStorage.getItem('theme') !== 'light')
const isCollapsed = ref(false)

const activeMenu = computed(() => route.path)

const toggleTheme = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

const handleMenuSelect = (index) => {
  router.push(index)
}

onMounted(() => {
  document.documentElement.classList.toggle('dark', isDark.value)
})
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { width: 100%; height: 100%; overflow: hidden; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif; }

html.dark {
  --app-bg: #0a0e1a;
  --sidebar-bg: rgba(15, 20, 35, 0.95);
  --sidebar-border: rgba(255,255,255,0.06);
  --content-bg: #0a0e1a;
  --card-bg: rgba(255,255,255,0.03);
  --card-border: rgba(255,255,255,0.08);
  --card-hover: rgba(255,255,255,0.06);
  --text-primary: #e8eaf0;
  --text-secondary: #8b90a0;
  --accent-glow: rgba(99,102,241,0.3);
  --gradient-start: #6366f1;
  --gradient-end: #8b5cf6;
  --btn-gradient: linear-gradient(135deg, #6366f1, #8b5cf6);
  --btn-hover: linear-gradient(135deg, #818cf8, #a78bfa);
}

html:not(.dark) {
  --app-bg: #f0f2f5;
  --sidebar-bg: #ffffff;
  --sidebar-border: #e8e8e8;
  --content-bg: #f0f2f5;
  --card-bg: rgba(255,255,255,0.8);
  --card-border: rgba(0,0,0,0.06);
  --card-hover: rgba(0,0,0,0.04);
  --text-primary: #1a1a2e;
  --text-secondary: #6b7280;
  --accent-glow: rgba(99,102,241,0.15);
  --gradient-start: #6366f1;
  --gradient-end: #8b5cf6;
  --btn-gradient: linear-gradient(135deg, #6366f1, #8b5cf6);
  --btn-hover: linear-gradient(135deg, #818cf8, #a78bfa);
}

.app-root { display: flex; height: 100vh; width: 100vw; background: var(--app-bg); color: var(--text-primary); }

/* Sidebar */
.sidebar {
  width: 240px;
  min-width: 240px;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--sidebar-border);
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4,0,0.2,1), min-width 0.3s;
  backdrop-filter: blur(20px);
  z-index: 100;
}
.sidebar.collapsed { width: 64px; min-width: 64px; }

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--sidebar-border);
  min-height: 60px;
}
.logo-area { display: flex; align-items: center; gap: 10px; overflow: hidden; }
.logo-icon {
  width: 32px; height: 32px;
  display: flex; align-items: center; justify-content: center;
  background: var(--btn-gradient);
  border-radius: 8px;
  font-size: 16px;
  flex-shrink: 0;
  box-shadow: 0 4px 12px var(--accent-glow);
}
.logo-text {
  font-size: 15px; font-weight: 700;
  background: var(--btn-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  white-space: nowrap;
}
.collapse-btn { cursor: pointer; font-size: 18px; color: var(--text-secondary); transition: all 0.2s; flex-shrink: 0; }
.collapse-btn:hover { color: var(--text-primary); }

/* el-menu overrides */
.side-menu {
  flex: 1;
  padding: 8px;
  overflow-y: auto;
  border-right: none !important;
}
.side-menu.el-menu--collapse { width: 48px; }
.side-menu:not(.el-menu--collapse) { width: 224px; }

html.dark .side-menu .el-sub-menu__title,
html.dark .side-menu .el-menu-item {
  background: transparent !important;
  color: var(--text-secondary) !important;
  border-radius: 10px !important;
  margin-bottom: 2px;
  height: 40px !important;
  line-height: 40px !important;
}
html.dark .side-menu .el-sub-menu__title:hover,
html.dark .side-menu .el-menu-item:hover {
  background: var(--card-hover) !important;
  color: var(--text-primary) !important;
}
html.dark .side-menu .el-menu-item.is-active {
  background: rgba(99,102,241,0.15) !important;
  color: var(--gradient-start) !important;
}

/* Sub-menu dropdown */
html.dark .side-menu .el-menu {
  background: rgba(255,255,255,0.02) !important;
}

/* Popper dropdown */
html.dark .el-menu--popup {
  background: #1a1f35 !important;
  border-color: var(--card-border) !important;
}
html.dark .el-menu--popup .el-menu-item {
  color: var(--text-secondary) !important;
}
html.dark .el-menu--popup .el-menu-item:hover {
  background: var(--card-hover) !important;
  color: var(--text-primary) !important;
}
html.dark .el-menu--popup .el-menu-item.is-active {
  background: rgba(99,102,241,0.15) !important;
  color: var(--gradient-start) !important;
}

.sidebar-footer { padding: 12px 16px; border-top: 1px solid var(--sidebar-border); }
.theme-toggle { cursor: pointer; font-size: 20px; color: var(--text-secondary); transition: all 0.2s; }
.theme-toggle:hover { color: var(--text-primary); }

/* Main Content */
.main-content { flex: 1; position: relative; overflow: hidden; display: flex; flex-direction: column; }
.content-gradient {
  position: absolute; top: -200px; right: -200px;
  width: 600px; height: 600px;
  background: radial-gradient(circle, var(--accent-glow) 0%, transparent 70%);
  pointer-events: none; z-index: 0;
}
.content-inner { flex: 1; overflow-y: auto; overflow-x: hidden; padding: 24px; position: relative; z-index: 1; }

/* Transition */
.fade-slide-enter-active { transition: all 0.25s ease; }
.fade-slide-leave-active { transition: all 0.15s ease; }
.fade-slide-enter-from { opacity: 0; transform: translateY(8px); }
.fade-slide-leave-to { opacity: 0; transform: translateY(-8px); }

/* Element Plus dark overrides */
html.dark .el-card { background: var(--card-bg) !important; border-color: var(--card-border) !important; }
html.dark .el-input__wrapper, html.dark .el-textarea__inner, html.dark .el-select__wrapper {
  background: rgba(255,255,255,0.04) !important;
  box-shadow: 0 0 0 1px var(--card-border) inset !important;
}
html.dark .el-input__wrapper:hover, html.dark .el-textarea__inner:hover {
  box-shadow: 0 0 0 1px rgba(255,255,255,0.15) inset !important;
}
html.dark .el-input__inner, html.dark .el-textarea__inner { color: var(--text-primary) !important; }
html.dark .el-select-dropdown, html.dark .el-popper, html.dark .el-popover {
  background: #1a1f35 !important; border-color: var(--card-border) !important; color: var(--text-primary) !important;
}
html.dark .el-select-dropdown__item { color: var(--text-primary) !important; }
html.dark .el-select-dropdown__item:hover { background: rgba(255,255,255,0.08) !important; }
html.dark .el-dialog { background: #1a1f35 !important; }
html.dark .el-table {
  --el-table-bg-color: transparent !important;
  --el-table-tr-bg-color: transparent !important;
  --el-table-header-bg-color: rgba(255,255,255,0.03) !important;
  --el-table-border-color: var(--card-border) !important;
}
html.dark .el-radio-button__inner {
  background: rgba(255,255,255,0.04) !important;
  border-color: var(--card-border) !important;
  color: var(--text-secondary) !important;
}
html.dark .el-radio-button__original-radio:checked + .el-radio-button__inner {
  background: var(--btn-gradient) !important;
  border-color: transparent !important;
  color: #fff !important;
}
html.dark .el-slider__runway { background-color: rgba(255,255,255,0.1) !important; }
html.dark .el-message-box { background: #1a1f35 !important; border-color: var(--card-border) !important; }
html.dark .el-message-box__title { color: var(--text-primary) !important; }
html.dark .el-message-box__content { color: var(--text-secondary) !important; }
html.dark .el-upload { color: var(--text-primary) !important; }
html.dark .el-tabs__item { color: var(--text-secondary) !important; }
html.dark .el-tabs__item.is-active { color: var(--gradient-start) !important; }
html.dark .el-tabs__active-bar { background-color: var(--gradient-start) !important; }
html.dark .el-tabs__nav-wrap::after { background-color: var(--card-border) !important; }
html.dark .el-switch__core { background-color: rgba(255,255,255,0.15) !important; border-color: var(--card-border) !important; }

/* Scrollbar */
::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.2); }
</style>
