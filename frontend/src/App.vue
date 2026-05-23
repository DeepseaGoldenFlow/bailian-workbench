<template>
  <el-container class="app-container">
    <el-aside width="220px">
      <Sidebar />
    </el-aside>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import Sidebar from './components/Sidebar.vue'
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'Segoe UI', 'Microsoft YaHei', sans-serif; background: #f0f2f5; }
.app-container { height: 100vh; }
.el-aside { background: #1a1a2e; }
.el-main { padding: 24px; overflow-y: auto; }
</style>