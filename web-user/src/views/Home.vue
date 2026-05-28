<template>
  <div class="page-container">
    <!-- Header -->
    <header class="app-header animate-fade-up">
      <div class="header-left">
        <div class="avatar">王</div>
        <div>
          <div class="greeting-text">今天想吃什么</div>
          <div class="greeting-name">小王</div>
        </div>
      </div>
      <button class="btn-icon" aria-label="搜索">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
        </svg>
      </button>
    </header>

    <!-- Hero Card -->
    <section class="hero animate-fade-up anim-delay-1" @click="goToRecipe(1)">
      <div class="hero-bg-shapes"></div>
      <div class="hero-content">
        <div class="hero-badge">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14l-5-4.87 6.91-1.01L12 2z"/>
          </svg>
          今日推荐
        </div>
        <div class="hero-title">红烧排骨</div>
        <div class="hero-sub">经典家常菜，软烂入味</div>
      </div>
      <div class="hero-avatars">
        <div class="mini-avatar" style="background: #ffb347;">王</div>
        <div class="mini-avatar" style="background: #a8d8ea;">李</div>
        <div class="mini-avatar more">+4</div>
      </div>
    </section>

    <!-- Meal Type Selector -->
    <section class="animate-fade-up anim-delay-2">
      <div class="meal-selector">
        <div
          v-for="meal in meals"
          :key="meal.key"
          class="meal-item"
          :class="{ active: activeMeal === meal.key }"
          @click="activeMeal = meal.key"
        >
          <div class="meal-circle" :class="{ active: activeMeal === meal.key }">
            <svg v-if="meal.key === 'breakfast'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
              <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
            </svg>
            <svg v-else-if="meal.key === 'lunch'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
              <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
            </svg>
            <svg v-else-if="meal.key === 'dinner'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
            </svg>
          </div>
          <span class="meal-label">{{ meal.label }}</span>
        </div>
      </div>
    </section>

    <!-- Plan Cards -->
    <section class="plan-section animate-fade-up anim-delay-3">
      <div class="section-title">推荐菜单</div>
      <div class="plan-cards">
        <div class="plan-card orange" @click="goToRecipe(2)">
          <div>
            <div class="plan-tag glass-tag glass-tag-orange">主菜</div>
            <div class="plan-name">糖醋里脊</div>
            <div class="plan-meta">30 分钟 · 3 人份<br>酸甜可口</div>
          </div>
          <div class="plan-footer">
            <div class="plan-avatar">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <div>
              <div class="plan-author-label">难度</div>
              <div class="plan-author-name">简单</div>
            </div>
          </div>
        </div>
        <div class="plan-right">
          <div class="plan-card blue" @click="goToRecipe(3)">
            <div>
              <div class="plan-tag glass-tag glass-tag-blue">配菜</div>
              <div class="plan-name">清炒时蔬</div>
              <div class="plan-meta">15 分钟 · 2 人份</div>
            </div>
          </div>
          <div class="plan-card pink" @click="goToRecipe(4)">
            <div class="plan-tag glass-tag glass-tag-purple">汤品</div>
            <div class="pink-icons">
              <div class="pink-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 8h1a4 4 0 0 1 0 8h-1"/>
                  <path d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"/>
                </svg>
              </div>
              <div class="pink-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2v4M4.93 4.93l2.83 2.83M20 12h-4M4.93 19.07l2.83-2.83M12 18v4M19.07 4.93l-2.83 2.83M4 12H0"/>
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Quick Actions -->
    <section class="actions animate-fade-up anim-delay-4">
      <div class="glass-list">
        <div class="glass-list-item" @click="$router.push('/week-menu')">
          <div class="glass-icon glass-icon-green">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/>
              <path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/>
            </svg>
          </div>
          <div class="action-text">
            <div class="action-title">AI 智能推荐</div>
            <div class="action-desc">告诉我你有什么食材</div>
          </div>
          <div class="action-arrow">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
          </div>
        </div>
        <div class="glass-list-item" @click="$router.push('/week-menu')">
          <div class="glass-icon glass-icon-purple">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
          </div>
          <div class="action-text">
            <div class="action-title">一周菜单</div>
            <div class="action-desc">自动生成 7 天食谱</div>
          </div>
          <div class="action-arrow">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
          </div>
        </div>
        <div class="glass-list-item" @click="$router.push('/shopping-list')">
          <div class="glass-icon glass-icon-yellow">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/>
              <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
            </svg>
          </div>
          <div class="action-text">
            <div class="action-title">购物清单</div>
            <div class="action-desc">根据菜单自动生成</div>
          </div>
          <div class="action-arrow">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const activeMeal = ref('breakfast')

const meals = [
  { key: 'breakfast', label: '早餐' },
  { key: 'lunch', label: '午餐' },
  { key: 'dinner', label: '晚餐' },
  { key: 'snack', label: '夜宵' },
]

function goToRecipe(id: number) {
  router.push(`/recipes/${id}`)
}
</script>

<style scoped>
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ffb347, #f59e0b);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 18px;
  box-shadow: 0 4px 12px rgba(245,158,11,0.3), inset 0 1px 0 rgba(255,255,255,0.3);
  border: 2px solid rgba(255,255,255,0.4);
}

.greeting-text {
  font-size: 14px;
  color: #6b6b6b;
  font-weight: 500;
}

.greeting-name {
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

/* Hero Card */
.hero {
  position: relative;
  border-radius: 28px;
  overflow: hidden;
  height: 180px;
  margin-bottom: 24px;
  background: linear-gradient(135deg, #c9b8f0 0%, #a894d8 50%, #8b7bc6 100%);
  cursor: pointer;
  transition: transform 0.2s ease;
  box-shadow: 0 8px 40px rgba(139,123,198,0.25);
}

.hero::after {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 1;
  background: linear-gradient(180deg, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0.05) 40%, transparent 100%);
  pointer-events: none;
}

.hero:active { transform: scale(0.98); }

.hero-content {
  position: relative;
  z-index: 2;
  padding: 24px;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  width: fit-content;
  padding: 6px 14px;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(12px) saturate(1.3);
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
  color: white;
  margin-bottom: 12px;
  border: 1px solid rgba(255,255,255,0.25);
  box-shadow: inset 0 1px 0 rgba(255,255,255,0.25);
}

.hero-title {
  font-size: 28px;
  font-weight: 800;
  color: white;
  line-height: 1.15;
  letter-spacing: -0.8px;
  margin-bottom: 4px;
}

.hero-sub {
  font-size: 14px;
  color: rgba(255,255,255,0.85);
  font-weight: 500;
}

.hero-avatars {
  position: absolute;
  right: 24px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 2;
  display: flex;
}

.mini-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 3px solid rgba(255,255,255,0.5);
  margin-left: -12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12);
}

.mini-avatar:first-child { margin-left: 0; }

.mini-avatar.more {
  background: rgba(255,255,255,0.25);
  backdrop-filter: blur(10px);
  font-size: 11px;
  border: 1px solid rgba(255,255,255,0.3);
}

.hero-bg-shapes {
  position: absolute;
  inset: 0;
  z-index: 1;
  overflow: hidden;
}

.hero-bg-shapes::before {
  content: '';
  position: absolute;
  width: 140px;
  height: 140px;
  background: rgba(255,255,255,0.2);
  border-radius: 50%;
  top: -30px;
  right: -20px;
  filter: blur(2px);
}

.hero-bg-shapes::after {
  content: '';
  position: absolute;
  width: 100px;
  height: 100px;
  background: rgba(255,255,255,0.12);
  border-radius: 50%;
  bottom: -25px;
  right: 60px;
  filter: blur(2px);
}

/* Meal Selector */
.meal-selector {
  display: flex;
  gap: 12px;
  margin-bottom: 28px;
  justify-content: space-between;
}

.meal-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.meal-item:active { transform: scale(0.92); }

.meal-circle {
  width: 100%;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(20px) saturate(1.3);
  color: #6b6b6b;
  border: 1px solid rgba(255,255,255,0.5);
  box-shadow: 0 8px 32px rgba(0,0,0,0.06), inset 0 1px 0 rgba(255,255,255,0.8);
}

.meal-circle.active {
  background: #1e1e2e;
  color: white;
  box-shadow: 0 6px 20px rgba(30,30,46,0.3);
}

.meal-circle svg { width: 22px; height: 22px; }

.meal-label {
  font-size: 12px;
  font-weight: 600;
  color: #6b6b6b;
}

.meal-item.active .meal-label {
  color: #1a1a1a;
}

/* Plan Cards */
.plan-section { margin-bottom: 28px; }

.plan-cards {
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: 12px;
  align-items: stretch;
}

.plan-card {
  border-radius: 24px;
  padding: 18px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: transform 0.2s ease;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0,0,0,0.08);
}

.plan-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: linear-gradient(180deg, rgba(255,255,255,0.35) 0%, rgba(255,255,255,0.1) 50%, transparent 100%);
  pointer-events: none;
  z-index: 1;
}

.plan-card:active { transform: scale(0.96); }

.plan-card.orange {
  background: linear-gradient(180deg, rgba(255,224,168,0.65) 0%, rgba(255,201,102,0.7) 100%);
  backdrop-filter: blur(16px) saturate(1.3);
  border: 1px solid rgba(255,255,255,0.5);
}

.plan-card.blue {
  background: linear-gradient(180deg, rgba(204,232,249,0.6) 0%, rgba(168,216,240,0.65) 100%);
  backdrop-filter: blur(16px) saturate(1.3);
  border: 1px solid rgba(255,255,255,0.5);
}

.plan-card.pink {
  background: linear-gradient(180deg, rgba(245,208,240,0.55) 0%, rgba(232,176,224,0.6) 100%);
  backdrop-filter: blur(16px) saturate(1.3);
  border: 1px solid rgba(255,255,255,0.5);
  border-radius: 20px;
  padding: 14px 18px;
}

.plan-right {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.plan-tag {
  display: inline-block;
  padding: 5px 14px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
  width: fit-content;
  margin-bottom: 14px;
}

.plan-card.orange .plan-tag { color: #b45309; }
.plan-card.blue .plan-tag { color: #1d4ed8; }
.plan-card.pink .plan-tag { color: #9333ea; }

.plan-name {
  font-size: 26px;
  font-weight: 800;
  color: #1a1a1a;
  letter-spacing: -0.8px;
  margin-bottom: 12px;
  line-height: 1.15;
  position: relative;
  z-index: 2;
}

.plan-meta {
  font-size: 14px;
  color: #6b6b6b;
  line-height: 1.8;
  position: relative;
  z-index: 2;
}

.plan-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: auto;
  padding-top: 14px;
  position: relative;
  z-index: 2;
}

.plan-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255,255,255,0.35);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid rgba(255,255,255,0.35);
}

.plan-avatar svg { width: 18px; height: 18px; color: #a0a0a0; }

.plan-author-label {
  font-size: 11px;
  color: #a0a0a0;
  font-weight: 500;
}

.plan-author-name {
  font-size: 13px;
  color: #6b6b6b;
  font-weight: 600;
}

.pink-icons {
  display: flex;
  gap: 10px;
}

.pink-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255,255,255,0.5);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #a0a0a0;
  border: 1px solid rgba(255,255,255,0.45);
}

.pink-icon svg { width: 18px; height: 18px; }

/* Actions */
.actions { margin-bottom: 28px; }

.action-text { flex: 1; }

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  letter-spacing: -0.3px;
}

.action-desc {
  font-size: 13px;
  color: #a0a0a0;
  margin-top: 2px;
}

.action-arrow {
  color: #a0a0a0;
  flex-shrink: 0;
}

.action-arrow svg { width: 20px; height: 20px; }
</style>
