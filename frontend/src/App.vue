<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="sidebar-logo"><div class="logo-icon">N</div><span class="logo-text">Nova</span></div>
      <nav class="sidebar-nav">
        <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="nav-item">
          <span class="nav-icon" v-html="item.icon"></span>
          <span class="nav-label">{{ item.label }}</span>
        </router-link>
      </nav>
      <div class="sidebar-footer"><div class="status-dot"></div><span class="status-text">Bailian Beijing</span></div>
    </aside>
    <main class="main-content"><router-view /></main>
  </div>
</template>

<script setup>
const navItems = [
  { path: '/', label: 'AI Chat', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/></svg>' },
  { path: '/image', label: 'Image Gen', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/></svg>' },
  { path: '/video', label: 'Video Gen', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/></svg>' },
  { path: '/tts', label: 'TTS', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 1a3 3 0 00-3 3v8a3 3 0 006 0V4a3 3 0 00-3-3z"/></svg>' },
  { path: '/toolbox', label: 'Toolbox', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/></svg>' },
  { path: '/history', label: 'History', icon: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>' },
]
</script>

<style>
.app-shell { display: flex; height: 100vh; background: var(--bg-primary); }
.sidebar { width: 220px; background: var(--bg-secondary); border-right: 1px solid var(--border-subtle); display: flex; flex-direction: column; flex-shrink: 0; }
.sidebar-logo { display: flex; align-items: center; gap: 10px; padding: 24px 20px 20px; }
.logo-icon { width: 32px; height: 32px; background: linear-gradient(135deg, var(--accent), #a855f7); border-radius: var(--radius-sm); display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 14px; color: #fff; }
.logo-text { font-size: 18px; font-weight: 700; background: linear-gradient(135deg, var(--text-primary), var(--text-secondary)); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
.sidebar-nav { flex: 1; padding: 8px 12px; display: flex; flex-direction: column; gap: 2px; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 9px 12px; border-radius: var(--radius-sm); color: var(--text-secondary); text-decoration: none; font-size: 13px; font-weight: 500; transition: all var(--transition); }
.nav-item:hover { background: var(--bg-card); color: var(--text-primary); }
.nav-item.router-link-active { background: var(--accent-soft); color: var(--accent); }
.nav-icon { width: 18px; height: 18px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.sidebar-footer { padding: 14px 20px; border-top: 1px solid var(--border-subtle); display: flex; align-items: center; gap: 8px; }
.status-dot { width: 6px; height: 6px; background: var(--success); border-radius: 50%; box-shadow: 0 0 6px var(--success); }
.status-text { font-size: 11px; color: var(--text-muted); }
.main-content { flex: 1; overflow-y: auto; padding: 32px; background: radial-gradient(ellipse at 20% 50%, rgba(99,102,241,0.04) 0%, transparent 50%), radial-gradient(ellipse at 80% 20%, rgba(168,85,247,0.03) 0%, transparent 50%), var(--bg-primary); }
</style>
