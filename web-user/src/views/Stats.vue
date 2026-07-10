<template>
  <div class="stats-shell">
    <main class="stats-phone">
      <header class="stats-header">
        <button class="back-btn" type="button" aria-label="返回" @click="router.back()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6" />
          </svg>
        </button>
        <div>
          <p>Nutrition</p>
          <h1>营养统计</h1>
        </div>
      </header>

      <section class="empty-card">
        <div class="empty-icon" aria-hidden="true">
          <svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 34c0-9 6-16 16-20 2 10-1 20-10 24" />
            <path d="M14 34c7-1 13-6 18-15" />
            <path d="M8 38h32" />
          </svg>
        </div>
        <h2>暂无营养统计数据</h2>
        <p>当前后端还没有提供营养记录和统计接口，所以这里不会展示前端写死的热量、蛋白质或周报图表。</p>
        <button type="button" @click="router.push('/recipes')">去浏览菜谱</button>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.stats-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --coral: #e95645;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: clip;
  color: var(--text);
  background:
    radial-gradient(circle at 12% 8%, rgba(255, 255, 255, 0.86), transparent 32%),
    linear-gradient(180deg, #fff2dc 0%, #f7e4c7 46%, #ead0aa 100%);
}

.stats-phone {
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: max(24px, env(safe-area-inset-top)) 24px 126px;
}

.stats-header {
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr);
  align-items: center;
  gap: 14px;
  padding-top: 18px;
}

.back-btn {
  width: 48px;
  height: 48px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: 18px;
  color: #6f6258;
  background: rgba(255, 250, 240, 0.78);
  box-shadow: 0 12px 28px rgba(80, 50, 30, 0.12);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  cursor: pointer;
}

.back-btn svg {
  width: 24px;
  height: 24px;
  stroke-width: 2.5;
}

.stats-header p {
  margin: 0 0 2px;
  color: var(--coral);
  font-size: 12px;
  font-weight: 850;
}

.stats-header h1 {
  margin: 0;
  color: var(--text);
  font-size: 30px;
  font-weight: 950;
  line-height: 1.12;
  letter-spacing: 0;
}

.empty-card {
  min-height: 420px;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 18px;
  margin-top: 28px;
  padding: 36px 28px;
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius);
  background: var(--card-surface);
  box-shadow: var(--card-shadow);
  text-align: center;
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.empty-icon {
  width: 96px;
  height: 96px;
  display: grid;
  place-items: center;
  border-radius: var(--card-radius-feature);
  color: #6f8b65;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.86), rgba(241, 233, 208, 0.66)),
    rgba(143, 167, 131, 0.18);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.88);
}

.empty-icon svg {
  width: 54px;
  height: 54px;
  stroke-width: 2.4;
}

.empty-card h2 {
  margin: 8px 0 0;
  color: var(--text);
  font-size: 25px;
  font-weight: 950;
  line-height: 1.2;
}

.empty-card p {
  max-width: 300px;
  margin: 0;
  color: var(--sub);
  font-size: 15px;
  font-weight: 680;
  line-height: 1.7;
}

.empty-card button {
  min-height: 46px;
  margin-top: 4px;
  padding: 0 22px;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: linear-gradient(135deg, #ef6153, #e94c3c);
  box-shadow: 0 14px 28px rgba(233, 86, 69, 0.24);
  font-size: 15px;
  font-weight: 850;
  cursor: pointer;
  transition: transform 180ms ease, opacity 180ms ease;
}

.back-btn:active,
.empty-card button:active {
  transform: scale(0.98);
}

@media (max-width: 380px) {
  .stats-phone {
    padding-left: 20px;
    padding-right: 20px;
  }

  .empty-card {
    min-height: 390px;
    padding: 32px 22px;
  }
}
</style>
