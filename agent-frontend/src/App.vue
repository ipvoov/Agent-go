<template>
  <div class="container">
    <header class="nav glass card" style="margin:16px 0; padding:10px 14px; display:flex; align-items:center; gap:12px; border-radius:999px;">
      <div style="display:flex;align-items:center; gap:10px; font-weight:700;">
        <img src="/src/assets/logo-mark.svg" alt="mark" style="width:24px;height:24px;" />
        AI 应用集合
      </div>
      <div style="margin-left:auto; display:flex; gap:8px;">
        <RouterLink to="/" class="pill">主页</RouterLink>
        <RouterLink to="/love" class="pill">AI 恋爱大师</RouterLink>
        <RouterLink to="/agent" class="pill">AI 超级智能体</RouterLink>
        <button class="pill" @click="toggleTheme">{{ theme === 'dark' ? '切到亮色' : '切到暗色' }}</button>
      </div>
    </header>
    <RouterView />
    <footer style="margin:28px 0; color:var(--muted); font-size:12px; text-align:center;">© {{ new Date().getFullYear() }} AI 应用集合</footer>
  </div>
</template>

<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { ref, onMounted } from 'vue'

const theme = ref('light')
onMounted(() => {
  theme.value = (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) ? 'dark' : 'light'
  applyTheme()
})

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  applyTheme()
}

function applyTheme() {
  document.documentElement.setAttribute('data-theme', theme.value)
}
</script>

<style scoped>
</style>

