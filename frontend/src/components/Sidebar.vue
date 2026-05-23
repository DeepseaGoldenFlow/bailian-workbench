<template>
  <div class="sidebar">
    <div class="logo">百炼工作台</div>
    <el-menu :default-active="route.path" router background-color="#1a1a2e" text-color="#a0a0c0" active-text-color="#409eff">
      <el-menu-item index="/">
        <el-icon><ChatDotRound /></el-icon>
        <span>AI 对话</span>
      </el-menu-item>
      <el-menu-item index="/image">
        <el-icon><Picture /></el-icon>
        <span>图片生成</span>
      </el-menu-item>
      <el-menu-item index="/video">
        <el-icon><VideoCamera /></el-icon>
        <span>视频生成</span>
      </el-menu-item>
      <el-menu-item index="/tts">
        <el-icon><Microphone /></el-icon>
        <span>语音合成</span>
      </el-menu-item>
      <el-menu-item index="/asr">
        <el-icon><Reading /></el-icon>
        <span>语音识别</span>
      </el-menu-item>
      <el-menu-item index="/toolbox">
        <el-icon><Switch /></el-icon>
        <span>工具箱</span>
      </el-menu-item>
      <el-menu-item index="/history">
        <el-icon><Clock /></el-icon>
        <span>历史记录</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
const route = useRoute()
</script>

<style scoped>
.logo { color: #fff; font-size: 20px; font-weight: bold; padding: 20px 16px; text-align: center; border-bottom: 1px solid #2a2a4e; }
.el-menu { border-right: none; }
</style>