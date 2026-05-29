<template>
  <div class="page">
    <header class="header anim-delay-1">
      <button class="btn-ghost" aria-label="返回" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div class="header-center">
        <div class="header-title">一周菜单</div>
      </div>
      <button class="btn-ghost" aria-label="生成" @click="generateWeekMenu">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
      </button>
    </header>

    <!-- Week Selector -->
    <div class="week-bar anim-delay-2">
      <button
        v-for="(day, idx) in weekDays"
        :key="idx"
        class="week-day"
        :class="{ active: activeDayIdx === idx }"
        @click="activeDayIdx = idx"
      >
        <span class="week-day-label">{{ day.label }}</span>
        <span class="week-day-num">{{ day.date }}</span>
        <span v-if="weekMenuData[idx]" class="week-day-dot"></span>
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="empty-state anim-delay-3">
      <div class="loading-spinner"></div>
      <p>生成中...</p>
    </div>

    <!-- Empty -->
    <div v-else-if="!weekMenuData.length" class="empty-state anim-delay-3">
      <p class="empty-text">点击右上角 + 生成一周菜单</p>
    </div>

    <!-- Timeline -->
    <div v-else class="timeline anim-delay-3">
      <template v-if="currentDayMeals.length">
        <div v-for="(meal, idx) in currentDayMeals" :key="idx" class="tl-item">
          <div class="tl-left">
            <span class="tl-time">{{ mealTime(meal.type) }}</span>
            <span class="tl-dot" :class="meal.type"></span>
            <span v-if="idx < currentDayMeals.length - 1" class="tl-line"></span>
          </div>
          <div class="tl-content">
            <span class="tl-type" :class="meal.type">{{ mealTypeLabel(meal.type) }}</span>
            <div v-for="(dish, di) in meal.dishes" :key="di" class="tl-dish">
              <h3 class="tl-dish-name">{{ dish.name }}</h3>
              <p class="tl-dish-meta">{{ dish.cook_time }}分钟 · {{ dish.difficulty }}</p>
            </div>
          </div>
        </div>
      </template>
      <div v-else class="empty-state">今日暂无推荐</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import api from '@/api/index'

const activeDayIdx = ref(new Date().getDay() === 0 ? 6 : new Date().getDay() - 1)
const loading = ref(false)
const weekMenuData = ref<any[]>([])

const weekDays = [
  { label: '一', date: '' },
  { label: '二', date: '' },
  { label: '三', date: '' },
  { label: '四', date: '' },
  { label: '五', date: '' },
  { label: '六', date: '' },
  { label: '日', date: '' },
]

;(() => {
  const now = new Date()
  const dayOfWeek = now.getDay() === 0 ? 7 : now.getDay()
  const monday = new Date(now)
  monday.setDate(now.getDate() - dayOfWeek + 1)
  for (let i = 0; i < 7; i++) {
    const d = new Date(monday)
    d.setDate(monday.getDate() + i)
    weekDays[i].date = String(d.getDate())
  }
})()

const currentDayMeals = computed(() => {
  const day = weekMenuData.value[activeDayIdx.value]
  if (!day?.meals) return []
  return day.meals
})

function mealTime(type: string) {
  const map: Record<string, string> = { breakfast: '07:30', lunch: '12:00', dinner: '18:30' }
  return map[type] || ''
}

function mealTypeLabel(type: string) {
  const map: Record<string, string> = { breakfast: '早餐', lunch: '午餐', dinner: '晚餐' }
  return map[type] || type
}

async function generateWeekMenu() {
  loading.value = true
  try {
    const res: any = await api.post('/recommend/week-menu', {
      people_count: 2,
      meal_type: 'lunch',
      taste_preference: [],
      health_goal: '',
      avoid_ingredients: [],
      existing_ingredients: [],
      cook_time_preference: '',
    })
    weekMenuData.value = Array.isArray(res) ? res : []
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

onMounted(() => {})
</script>

<style scoped>
.header-center {
  flex: 1;
  text-align: center;
}

/* ── Week Bar ── */
.week-bar {
  display: flex;
  gap: var(--sp-1);
  margin-bottom: var(--sp-6);
}

.week-day {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: var(--sp-2) 0;
  border-radius: var(--r-sm);
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all var(--dur-base) var(--ease);
  position: relative;
  -webkit-tap-highlight-color: transparent;
}

.week-day.active {
  background: var(--color-text);
}

.week-day:active { transform: scale(0.95); }

.week-day-label {
  font-size: 11px;
  font-weight: 500;
  color: var(--color-text-3);
}

.week-day.active .week-day-label {
  color: rgba(255, 255, 255, 0.6);
}

.week-day-num {
  font-size: var(--text-md);
  font-weight: 700;
  color: var(--color-text);
  line-height: 1.2;
}

.week-day.active .week-day-num {
  color: var(--color-text-inv);
}

.week-day-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--color-accent);
  position: absolute;
  bottom: 4px;
}

/* ── Empty ── */
.empty-state {
  text-align: center;
  padding: var(--sp-12) 0;
  color: var(--color-text-3);
  font-size: var(--text-sm);
}

.empty-text {
  color: var(--color-text-3);
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-surface-3);
  border-top-color: var(--color-text-3);
  border-radius: 50%;
  margin: 0 auto var(--sp-3);
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Timeline ── */
.timeline {
  padding-bottom: var(--sp-8);
}

.tl-item {
  display: flex;
  gap: var(--sp-4);
}

.tl-left {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 40px;
  flex-shrink: 0;
}

.tl-time {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-3);
  white-space: nowrap;
}

.tl-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  margin: var(--sp-2) 0;
}

.tl-dot.breakfast { background: var(--color-accent); }
.tl-dot.lunch { background: var(--color-text); }
.tl-dot.dinner { background: var(--color-text-2); }

.tl-line {
  width: 1.5px;
  flex: 1;
  background: var(--color-border);
  min-height: 16px;
}

.tl-content {
  flex: 1;
  padding: 0 0 var(--sp-6);
}

.tl-type {
  font-size: 11px;
  font-weight: 600;
  margin-bottom: var(--sp-2);
  display: inline-block;
}

.tl-type.breakfast { color: var(--color-accent); }
.tl-type.lunch { color: var(--color-text); }
.tl-type.dinner { color: var(--color-text-2); }

.tl-dish {
  margin-bottom: var(--sp-2);
}

.tl-dish:last-child {
  margin-bottom: 0;
}

.tl-dish-name {
  font-family: var(--font-display);
  font-size: var(--text-md);
  font-weight: 650;
  color: var(--color-text);
  line-height: 1.3;
}

.tl-dish-meta {
  font-size: var(--text-xs);
  color: var(--color-text-3);
  margin-top: 2px;
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .page {
    max-width: 640px;
  }
}

@media (min-width: 1024px) {
  .page {
    max-width: 800px;
  }
}
</style>
