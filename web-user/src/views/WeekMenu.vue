<template>
  <div class="page">
    <!-- Header -->
    <header class="header anim-delay-1">
      <button class="btn-ghost" aria-label="返回" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div style="flex:1;text-align:center;">
        <div class="header-title">一周菜单</div>
      </div>
      <button class="btn-ghost" aria-label="生成" @click="generateWeekMenu">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
      </button>
    </header>

    <!-- Week Selector -->
    <div class="week-strip anim-delay-2">
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

    <!-- Loading / Empty -->
    <div v-if="loading" class="empty-hint anim-delay-3">生成中...</div>
    <div v-else-if="!weekMenuData.length" class="empty-hint anim-delay-3">
      <p>点击右上角 + 生成一周菜单</p>
    </div>

    <!-- Meal Timeline -->
    <div v-else class="timeline anim-delay-3">
      <template v-if="currentDayMeals.length">
        <div v-for="(meal, idx) in currentDayMeals" :key="idx" class="tl-item">
          <div class="tl-left">
            <span class="tl-time">{{ mealTime(meal.type) }}</span>
            <span class="tl-dot" :class="meal.type"></span>
            <span v-if="Number(idx) < currentDayMeals.length - 1" class="tl-line"></span>
          </div>
          <div class="tl-card">
            <div class="tl-card-top">
              <span class="tl-card-type" :class="meal.type">{{ mealTypeLabel(meal.type) }}</span>
            </div>
            <div v-for="(dish, di) in meal.dishes" :key="di" style="margin-bottom:var(--sp-2);">
              <h3 class="tl-card-name">{{ dish.name }}</h3>
              <p class="tl-card-desc">{{ dish.cook_time }}分钟 · {{ dish.difficulty }}</p>
            </div>
          </div>
        </div>
      </template>
      <div v-else class="empty-hint">今日暂无推荐</div>
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

// Fill in actual dates for this week
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

onMounted(() => {
  // Auto-load if data exists, otherwise show empty state
})
</script>

<style scoped>
/* ── Week selector ── */
.week-strip {
  display: flex;
  gap: var(--sp-1);
  margin-bottom: var(--sp-5);
}

.week-day {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: var(--sp-2) 0;
  border-radius: var(--r-md);
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-out);
  position: relative;
  -webkit-tap-highlight-color: transparent;
}

.week-day.active {
  background: var(--color-dark);
}

.week-day:active { transform: scale(0.93); }

.week-day-label {
  font-size: 11px;
  font-weight: 500;
  color: var(--color-ink-3);
}

.week-day.active .week-day-label {
  color: var(--color-ink-3);
}

.week-day-num {
  font-size: var(--text-md);
  font-weight: 700;
  color: var(--color-ink);
  line-height: 1.2;
}

.week-day.active .week-day-num {
  color: var(--color-ink);
}

.week-day-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--color-accent);
  position: absolute;
  bottom: 4px;
}

/* ── Stats ── */
.stats-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-5);
  padding: var(--sp-4) 0;
  margin-bottom: var(--sp-6);
  border-bottom: 0.5px solid var(--color-rule);
}

.stat {
  display: flex;
  align-items: baseline;
  gap: var(--sp-1);
}

.stat-value {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 700;
  color: var(--color-ink);
}

.stat-value--accent { color: var(--color-accent); }
.stat-value--green { color: var(--color-green); }

.stat-label {
  font-size: var(--text-xs);
  color: var(--color-ink-3);
  font-weight: 500;
}

.stat-sep {
  width: 1px;
  height: 20px;
  background: var(--color-rule);
}

/* ── Timeline ── */
.tl-item {
  display: flex;
  gap: var(--sp-4);
  min-height: 80px;
}

.tl-left {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 44px;
  flex-shrink: 0;
  gap: var(--sp-1);
}

.tl-time {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-ink-3);
  white-space: nowrap;
}

.tl-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  margin: var(--sp-1) 0;
}

.tl-dot.breakfast { background: var(--color-orange); }
.tl-dot.lunch { background: var(--color-blue); }
.tl-dot.dinner { background: var(--color-purple); }

.tl-line {
  width: 1.5px;
  flex: 1;
  background: var(--color-rule);
  min-height: 16px;
}

.tl-card {
  flex: 1;
  padding: var(--sp-3) 0 var(--sp-5);
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  transition: opacity var(--dur-fast);
}

.tl-card:active { opacity: 0.7; }

.tl-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-1);
}

.tl-card-type {
  font-size: 11px;
  font-weight: 600;
}

.tl-card-type.breakfast { color: var(--color-orange); }
.tl-card-type.lunch { color: var(--color-blue); }
.tl-card-type.dinner { color: var(--color-purple); }

.tl-card-kcal {
  font-size: 11px;
  color: var(--color-ink-3);
  font-weight: 500;
}

.tl-card-name {
  font-family: var(--font-display);
  font-size: var(--text-md);
  font-weight: 700;
  color: var(--color-ink);
  letter-spacing: 0;
  line-height: 1.3;
  margin-bottom: 2px;
}

.tl-card-desc {
  font-size: var(--text-xs);
  color: var(--color-ink-3);
  line-height: 1.5;
}

.tl-card-tags {
  display: flex;
  gap: var(--sp-2);
  margin-top: var(--sp-2);
  flex-wrap: wrap;
}

.tl-add {
  flex: 1;
  padding: var(--sp-4);
  border-radius: var(--r-lg);
  border: 1.5px dashed var(--color-rule);
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-2);
  cursor: pointer;
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-ink-3);
  transition: border-color var(--dur-base) var(--ease-out), color var(--dur-base) var(--ease-out);
  margin-bottom: var(--sp-4);
}

.tl-add:active { border-color: var(--color-accent); color: var(--color-accent); }

.empty-hint {
  text-align: center;
  padding: var(--sp-8) 0;
  font-size: var(--text-sm);
  color: var(--color-ink-3);
}
</style>
