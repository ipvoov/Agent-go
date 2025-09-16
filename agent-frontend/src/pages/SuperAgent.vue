<template>
  <div class="card" style="padding:0; height: calc(100vh - 140px); display:flex;">
    <ChatRoom title="AI 超级智能体" endpoint="/agentStream" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ChatRoom from '../shared/ChatRoom.vue'

const route = useRoute()

onMounted(() => {
  // 如果从首页带着 q 过来，自动填入并发送
  const q = (route.query.q || '').toString().trim()
  if (q) {
    const input = document.querySelector('input.input')
    if (input) {
      input.value = q
      const event = new Event('input', { bubbles: true })
      input.dispatchEvent(event)
    }
    const btn = document.querySelector('form button[type="submit"]')
    if (btn) btn.click()
  }
})
</script>

<style scoped>
</style>

