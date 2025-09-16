<template>
  <div class="home">
    <div class="hero" style="position:relative; overflow:hidden;">
      <div class="blob blob--a"></div>
      <div class="blob blob--b"></div>
      <div class="hero__content glass card" style="border-radius:20px;">
        <div>
          <div class="badge">AI Copilot</div>
          <h1 class="hero__title text-gradient" style="font-size:34px;">与智能体协作，创造超越想象的结果</h1>
          <p class="hero__subtitle">启发、执行与总结，一切尽在一个工作台。</p>
        </div>
        <div class="hero-search">
          <input v-model="quick" placeholder="输入想做的事，回车或点开始" @keyup.enter="goAgent()" />
          <button class="btn" @click="goAgent()">开始</button>
        </div>
        <div style="display:flex; gap:10px; flex-wrap:wrap;">
          <button class="pill" @click="runPrompt('帮我准备一个面试自我介绍')"><span class="dot"></span> 自我介绍</button>
          <button class="pill" @click="runPrompt('请为我生成一份本周工作周报，包含完成事项与下周计划')"><span class="dot"></span> 写周报</button>
          <button class="pill" @click="runPrompt('为新品发布生成一段简洁有力的海报文案')"><span class="dot"></span> 海报文案</button>
          <button class="pill" @click="runPrompt('请为 Go 新手制定两周的学习计划')"><span class="dot"></span> 学习计划</button>
        </div>
      </div>
    </div>

    <div class="features">
      <div class="feature card">
        <div class="feature__icon">💬</div>
        <h4>对话体验</h4>
        <p>左右气泡、自动滚动与流式渲染，媲美主流聊天产品。</p>
      </div>
      <div class="feature card">
        <div class="feature__icon">⚡️</div>
        <h4>极速响应</h4>
        <p>基于 SSE 的毫秒级首字延迟，边生成边显示。</p>
      </div>
      <div class="feature card">
        <div class="feature__icon">🧩</div>
        <h4>可拓展</h4>
        <p>独立会话 ID 与清晰组件化，轻松新增更多智能体。</p>
      </div>
    </div>

    <div class="apps">
      <RouterLink to="/love" class="app card glass shimmer">
        <div class="app__icon">🫶</div>
        <div class="app__meta">
          <h3>AI 恋爱大师</h3>
          <p>情感建议、沟通优化与话术生成。</p>
        </div>
        <span class="app__enter">进入 →</span>
      </RouterLink>

      <RouterLink to="/agent" class="app card glass shimmer">
        <div class="app__icon">🧠</div>
        <div class="app__meta">
          <h3>AI 超级智能体</h3>
          <p>多步推理、任务分解与工具调用。</p>
        </div>
        <span class="app__enter">进入 →</span>
      </RouterLink>
    </div>

    <div class="templates">
      <h4 style="margin:4px 0 8px 0;">常用模板</h4>
      <div class="templates__grid">
        <div class="tpl card glass">
          <div class="tpl__title">工作周报</div>
          <div class="tpl__desc">总结本周关键产出并生成待办</div>
          <RouterLink to="/agent" class="btn" style="align-self:flex-start;">使用</RouterLink>
        </div>
        <div class="tpl card glass">
          <div class="tpl__title">学习计划</div>
          <div class="tpl__desc">为 2 周内掌握 Go 的核心语法制定计划</div>
          <RouterLink to="/agent" class="btn" style="align-self:flex-start;">使用</RouterLink>
        </div>
        <div class="tpl card glass">
          <div class="tpl__title">恋爱话术</div>
          <div class="tpl__desc">根据对象性格生成 5 种回复风格</div>
          <RouterLink to="/love" class="btn" style="align-self:flex-start;">使用</RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink, useRouter } from 'vue-router'
import { ref } from 'vue'

const router = useRouter()
const quick = ref('')

function runPrompt(text) {
  quick.value = text
  goAgent()
}

function goAgent() {
  const text = (quick.value || '').trim()
  router.push({ name: 'agent', query: text ? { q: text } : {} })
}
</script>

<style scoped>
.home { display: grid; gap: 28px; }

.hero { position: relative; }
.hero__bg { position:absolute; inset:-20px 0 -40px 0; background: radial-gradient(40% 60% at 30% 20%, rgba(129,140,248,.25), transparent 60%), radial-gradient(40% 60% at 70% 20%, rgba(56,189,248,.18), transparent 60%); filter: blur(30px); pointer-events:none; }
.hero__content { position: relative; padding: 32px; display:grid; gap: 18px; overflow:hidden; }
.badge { display:inline-block; padding:6px 10px; font-size:12px; border-radius:999px; background: rgba(79,70,229,.12); color: var(--primary); border:1px solid var(--border); }
.hero__title { margin: 0; font-size: 28px; line-height: 1.2; }
.hero__subtitle { margin: 0; color: var(--muted); }
.hero__cta { display:flex; gap:12px; }
.btn--ghost { background: transparent; color: var(--primary); border:1px solid var(--border); }
.hero__stats { display:flex; gap:18px; color: var(--muted); font-size: 13px; }
.hero__stats strong { display:block; color: var(--text); font-size: 16px; }

.features { display:grid; grid-template-columns: repeat(auto-fit, minmax(220px,1fr)); gap:16px; }
.feature { padding:16px; display:grid; gap:8px; transition: transform .2s ease; }
.feature:hover { transform: translateY(-2px); }
.feature__icon { font-size:22px; }

.apps { display:grid; grid-template-columns: repeat(auto-fit, minmax(320px,1fr)); gap:16px; }
.app { display:flex; align-items:center; gap:14px; padding:18px; text-decoration:none; color:inherit; transition: transform .2s ease, box-shadow .2s ease; }
.app:hover { transform: translateY(-3px); box-shadow: 0 6px 20px rgba(0,0,0,.06); }
.app__icon { width:44px; height:44px; border-radius:12px; display:grid; place-items:center; background: rgba(79,70,229,.12); font-size:22px; }
.app__meta h3 { margin:0 0 4px 0; font-size:18px; }
.app__meta p { margin:0; color: var(--muted); font-size:14px; }
.app__enter { margin-left:auto; color: var(--primary); font-weight:600; }

.templates { margin-top: 10px; }
.templates__grid { display:grid; grid-template-columns: repeat(auto-fit, minmax(260px,1fr)); gap:16px; }
.tpl { padding:16px; display:grid; gap:8px; }
.tpl__title { font-weight:600; }
.tpl__desc { color: var(--muted); font-size: 13px; }
</style>

