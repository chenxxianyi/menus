<template>
  <div class="page">
    <header class="header anim-delay-1">
      <button class="btn-ghost" aria-label="返回" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div class="header-center">
        <div class="header-title">一周菜单</div>
      </div>
      <button class="btn-ghost" aria-label="生成" :disabled="loading" @click="generateWeekMenu">
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

    <section class="planner-hero anim-delay-3">
      <div>
        <p class="hero-kicker">{{ weekMenuData.length ? activeDayName : '智能周菜单' }}</p>
        <h1 class="hero-title">{{ weekMenuData.length ? '今天这样吃' : '一键安排三餐' }}</h1>
        <p class="hero-copy">
          {{ weekMenuData.length ? '看菜谱、攒采购清单、照着步骤做饭。' : '点击右上角 +，根据菜谱库生成一周早中晚餐。' }}
        </p>
      </div>
      <div class="hero-stats" aria-label="菜单概览">
        <div>
          <strong>{{ currentDishCount }}</strong>
          <span>道菜</span>
        </div>
        <div>
          <strong>{{ currentDayIngredients.length }}</strong>
          <span>食材</span>
        </div>
      </div>
    </section>

    <div v-if="loading" class="empty-state anim-delay-4">
      <div class="loading-spinner"></div>
      <p>正在搭配一周三餐...</p>
    </div>

    <div v-else-if="error" class="empty-state anim-delay-4">
      <p class="empty-title">生成失败</p>
      <p class="empty-text">{{ error }}</p>
      <button class="retry-btn" @click="generateWeekMenu">重新生成</button>
    </div>

    <div v-else-if="!weekMenuData.length" class="empty-state empty-panel anim-delay-4">
      <p class="empty-title">还没有周菜单</p>
      <p class="empty-text">点右上角的 +，系统会按早餐、午餐、晚餐生成一周计划。</p>
      <button class="retry-btn" @click="generateWeekMenu">生成一周菜单</button>
    </div>

    <template v-else>
      <section class="day-tools anim-delay-4">
        <div class="tool-copy">
          <span>今日采购</span>
          <strong>{{ currentDayIngredients.length ? currentDayIngredients.slice(0, 4).join('、') : '暂无食材' }}</strong>
        </div>
        <button class="tool-btn" :disabled="!currentDayIngredients.length || shoppingSaving" @click="createDayShoppingList">
          {{ shoppingSaving ? '生成中' : '加入采购清单' }}
        </button>
      </section>

      <p v-if="shoppingMessage" class="inline-message anim-delay-4">{{ shoppingMessage }}</p>

      <div v-if="isGeneratedButEmpty" class="data-note anim-delay-4">
        已经生成餐次，但后端菜谱库没有返回菜品。请先确认数据库里有启用状态的菜谱数据。
      </div>

      <div class="timeline anim-delay-5">
        <div v-for="(meal, idx) in currentDayMeals" :key="idx" class="tl-item">
          <div class="tl-left">
            <span class="tl-time">{{ mealTime(meal.type) }}</span>
            <span class="tl-dot" :class="meal.type"></span>
            <span v-if="hasNextMeal(idx)" class="tl-line"></span>
          </div>
          <div class="tl-content">
            <div class="meal-head">
              <span class="tl-type" :class="meal.type">{{ mealTypeLabel(meal.type) }}</span>
              <span class="meal-reason">{{ meal.reason || '按当前偏好推荐' }}</span>
            </div>
            <template v-if="meal.dishes.length">
              <article v-for="(dish, di) in meal.dishes" :key="dish.recipe_id || di" class="tl-dish">
                <div class="dish-main">
                  <div class="dish-mark" :class="meal.type">{{ dish.type || mealTypeLabel(meal.type) }}</div>
                  <div class="dish-body">
                    <h3 class="tl-dish-name">{{ dish.name }}</h3>
                    <p class="tl-dish-meta">{{ dish.cook_time || '--' }}分钟 · {{ dish.difficulty || '家常' }}</p>
                  </div>
                </div>
                <div v-if="dish.ingredients?.length" class="ingredient-row">
                  <span v-for="item in dish.ingredients.slice(0, 5)" :key="item" class="ingredient-chip">{{ item }}</span>
                </div>
                <div class="dish-actions">
                  <button class="dish-btn" :disabled="!dish.recipe_id" @click="openRecipe(dish.recipe_id)">看做法</button>
                  <button class="dish-btn dish-btn--quiet" @click="copyDishIngredients(dish)">复制食材</button>
                </div>
              </article>
            </template>
            <div v-else class="meal-empty">
              <strong>{{ mealTypeLabel(meal.type) }}还没有推荐菜</strong>
              <span>后端返回了餐次，但没有菜品内容。可以重新生成，或先去菜谱库新增菜谱。</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index'
import { useShoppingStore } from '@/stores/shopping'
import type { ShoppingItem } from '@/api/shopping'

const router = useRouter()
const shoppingStore = useShoppingStore()
const activeDayIdx = ref(new Date().getDay() === 0 ? 6 : new Date().getDay() - 1)
const loading = ref(false)
const error = ref('')
const weekMenuData = ref<any[]>([])
const shoppingSaving = ref(false)
const shoppingMessage = ref('')

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

const activeDayName = computed(() => `周${weekDays[activeDayIdx.value]?.label || ''}`)

const currentDishCount = computed(() => {
  return currentDayMeals.value.reduce((sum: number, meal: any) => sum + (meal.dishes?.length || 0), 0)
})

const currentDayIngredients = computed(() => {
  const names = new Set<string>()
  currentDayMeals.value.forEach((meal: any) => {
    ;(meal.shopping_list || []).forEach((item: string) => item && names.add(item))
    ;(meal.dishes || []).forEach((dish: any) => {
      ;(dish.ingredients || []).forEach((item: string) => item && names.add(item))
    })
  })
  return Array.from(names)
})

const isGeneratedButEmpty = computed(() => {
  return weekMenuData.value.length > 0 && weekMenuData.value.every((day) => {
    return (day.meals || []).every((meal: any) => !meal.dishes?.length)
  })
})

function mealTime(type: string) {
  const map: Record<string, string> = { breakfast: '07:30', lunch: '12:00', dinner: '18:30' }
  return map[type] || ''
}

function mealTypeLabel(type: string) {
  const map: Record<string, string> = { breakfast: '早餐', lunch: '午餐', dinner: '晚餐' }
  return map[type] || type
}

function hasNextMeal(idx: number | string) {
  return Number(idx) < currentDayMeals.value.length - 1
}

function normalizeWeekMenu(data: any[]) {
  const mealTypes = ['breakfast', 'lunch', 'dinner']
  return data.map((day) => ({
    ...day,
    meals: (day.meals || []).map((meal: any, idx: number) => ({
      ...meal,
      type: meal.type || meal.meal_type || mealTypes[idx] || 'lunch',
      dishes: meal.dishes || [],
      shopping_list: meal.shopping_list || [],
    })),
  }))
}

function openRecipe(id?: number) {
  if (!id) return
  router.push(`/recipes/${id}`)
}

function inferCategory(name: string) {
  if (/鸡|鸭|鱼|肉|排骨|虾|蛋|牛|猪/.test(name)) return '肉蛋水产'
  if (/米|面|粉|饭|面包|馒头/.test(name)) return '主食'
  if (/盐|糖|酱|醋|油|姜|蒜|葱|胡椒|淀粉/.test(name)) return '调味'
  if (/奶|酸奶|芝士/.test(name)) return '乳品'
  return '蔬果'
}

function toShoppingItems(names: string[]): ShoppingItem[] {
  return names.map((name) => ({
    name,
    amount: '按菜谱适量',
    emoji: '食材',
    category: inferCategory(name),
    price: 0,
    checked: false,
  }))
}

async function createDayShoppingList() {
  if (!currentDayIngredients.value.length || shoppingSaving.value) return
  shoppingSaving.value = true
  shoppingMessage.value = ''
  try {
    await shoppingStore.createList(`${activeDayName.value}菜单采购`, toShoppingItems(currentDayIngredients.value))
    shoppingMessage.value = '已生成采购清单，可以去购物清单页面查看。'
    setTimeout(() => {
      router.push('/shopping-list')
    }, 350)
  } catch (e: any) {
    shoppingMessage.value = e?.message || '采购清单生成失败，请稍后再试。'
  } finally {
    shoppingSaving.value = false
  }
}

async function copyDishIngredients(dish: any) {
  const text = (dish.ingredients || []).join('、')
  if (!text) {
    shoppingMessage.value = '这个菜谱暂时没有食材数据。'
    return
  }
  try {
    await navigator.clipboard?.writeText(text)
    shoppingMessage.value = `已复制「${dish.name}」食材。`
  } catch {
    shoppingMessage.value = text
  }
}

async function generateWeekMenu() {
  loading.value = true
  error.value = ''
  shoppingMessage.value = ''
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
    const list = Array.isArray(res) ? normalizeWeekMenu(res) : []
    weekMenuData.value = list
    if (!list.length) {
      error.value = '暂时没有可用菜谱，请先确认后端已有菜谱数据。'
    }
  } catch (e: any) {
    error.value = e?.message || '生成失败，请确认已登录并且后端服务正常运行。'
  } finally {
    loading.value = false
  }
}

onMounted(() => {})
</script>

<style scoped>
/* Hallmark · component: weekly meal planner · genre: warm utility · theme: kitchen ledger */
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

.btn-ghost:disabled {
  opacity: 0.45;
  cursor: wait;
}

.week-day-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--color-accent);
  position: absolute;
  bottom: 4px;
}

.planner-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: var(--sp-4);
  align-items: end;
  padding: var(--sp-5);
  margin-bottom: var(--sp-5);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  background:
    linear-gradient(135deg, var(--color-broth-soft), var(--color-surface) 58%),
    var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.hero-kicker {
  margin-bottom: var(--sp-1);
  color: var(--color-tomato);
  font-size: var(--text-xs);
  font-weight: 800;
}

.hero-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 750;
  line-height: 1.12;
}

.hero-copy {
  max-width: 24em;
  margin-top: var(--sp-2);
  color: var(--color-text-2);
  font-size: var(--text-sm);
  line-height: 1.65;
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(2, 64px);
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: color-mix(in srgb, var(--color-surface) 78%, transparent);
}

.hero-stats div {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 64px;
}

.hero-stats div + div {
  border-left: 1px solid var(--color-border);
}

.hero-stats strong {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-xl);
  line-height: 1;
}

.hero-stats span {
  margin-top: 3px;
  color: var(--color-text-3);
  font-size: var(--text-2xs);
  font-weight: 700;
}

.day-tools {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--sp-3);
  padding: var(--sp-4);
  margin-bottom: var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
}

.tool-copy {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.tool-copy span {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 700;
}

.tool-copy strong {
  overflow: hidden;
  color: var(--color-text);
  font-size: var(--text-sm);
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tool-btn,
.dish-btn,
.retry-btn {
  transition:
    transform var(--dur-fast) var(--ease),
    background var(--dur-base) var(--ease),
    border-color var(--dur-base) var(--ease),
    color var(--dur-base) var(--ease);
}

.tool-btn {
  min-height: 38px;
  flex-shrink: 0;
  padding: 0 var(--sp-4);
  border: 0;
  border-radius: var(--r-sm);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-xs);
  font-weight: 800;
  cursor: pointer;
}

.tool-btn:hover {
  background: var(--color-text-2);
}

.tool-btn:active,
.dish-btn:active,
.retry-btn:active {
  transform: translateY(1px);
}

.tool-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.inline-message {
  margin: 0 0 var(--sp-4);
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-accent-soft);
  border-radius: var(--r-sm);
  background: var(--color-accent-soft);
  color: var(--color-accent);
  font-size: var(--text-xs);
  font-weight: 700;
}

.data-note {
  margin-bottom: var(--sp-4);
  padding: var(--sp-4);
  border: 1px solid var(--color-warning-soft);
  border-radius: var(--r-md);
  background: var(--color-warning-soft);
  color: var(--color-text-2);
  font-size: var(--text-sm);
  line-height: 1.6;
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
  line-height: 1.6;
}

.empty-title {
  margin-bottom: var(--sp-2);
  color: var(--color-text);
  font-size: var(--text-md);
  font-weight: 750;
}

.empty-panel {
  padding: var(--sp-10) var(--sp-5);
  border: 1px dashed var(--color-border-med);
  border-radius: var(--r-lg);
  background: var(--color-surface);
}

.retry-btn {
  margin-top: var(--sp-4);
  min-height: 40px;
  padding: 0 var(--sp-5);
  border: 0;
  border-radius: var(--r-full);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-sm);
  font-weight: 700;
  cursor: pointer;
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
  min-width: 0;
}

.meal-head {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  margin-bottom: var(--sp-2);
}

.tl-type {
  font-size: 11px;
  font-weight: 600;
  display: inline-block;
  flex-shrink: 0;
}

.tl-type.breakfast { color: var(--color-accent); }
.tl-type.lunch { color: var(--color-text); }
.tl-type.dinner { color: var(--color-text-2); }

.meal-reason {
  min-width: 0;
  overflow: hidden;
  color: var(--color-text-3);
  font-size: var(--text-2xs);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tl-dish {
  padding: var(--sp-4);
  margin-bottom: var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.tl-dish:last-child {
  margin-bottom: 0;
}

.dish-main {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr);
  gap: var(--sp-3);
  align-items: center;
}

.dish-mark {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-sm);
  background: var(--color-surface-2);
  color: var(--color-text-2);
  font-size: var(--text-2xs);
  font-weight: 800;
}

.dish-mark.breakfast {
  background: var(--color-broth-soft);
  color: var(--color-tomato);
}

.dish-mark.lunch {
  background: var(--color-accent-soft);
  color: var(--color-accent);
}

.dish-mark.dinner {
  background: var(--color-water-soft);
  color: var(--color-water);
}

.dish-body {
  min-width: 0;
}

.tl-dish-name {
  font-family: var(--font-display);
  font-size: var(--text-md);
  font-weight: 650;
  color: var(--color-text);
  line-height: 1.3;
  overflow-wrap: anywhere;
}

.tl-dish-meta {
  font-size: var(--text-xs);
  color: var(--color-text-3);
  margin-top: 2px;
}

.ingredient-row {
  display: flex;
  flex-wrap: wrap;
  gap: var(--sp-2);
  margin-top: var(--sp-3);
}

.ingredient-chip {
  max-width: 100%;
  padding: 5px 9px;
  border-radius: var(--r-full);
  background: var(--color-surface-2);
  color: var(--color-text-2);
  font-size: var(--text-2xs);
  font-weight: 700;
}

.dish-actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--sp-2);
  margin-top: var(--sp-3);
}

.dish-btn {
  min-height: 36px;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-xs);
  font-weight: 800;
  cursor: pointer;
}

.dish-btn:hover {
  background: var(--color-text-2);
}

.dish-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.dish-btn--quiet {
  background: var(--color-surface);
  color: var(--color-text-2);
}

.dish-btn--quiet:hover {
  border-color: var(--color-border-med);
  background: var(--color-surface-2);
}

.meal-empty {
  display: flex;
  flex-direction: column;
  gap: var(--sp-1);
  padding: var(--sp-4);
  border: 1px dashed var(--color-border-med);
  border-radius: var(--r-md);
  background: var(--color-surface);
}

.meal-empty strong {
  color: var(--color-text);
  font-size: var(--text-sm);
}

.meal-empty span {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  line-height: 1.55;
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

@media (max-width: 420px) {
  .planner-hero {
    grid-template-columns: 1fr;
    align-items: stretch;
  }

  .hero-stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .day-tools {
    align-items: stretch;
    flex-direction: column;
  }

  .tool-btn {
    width: 100%;
  }

  .tl-item {
    gap: var(--sp-3);
  }

  .tl-left {
    width: 34px;
  }

  .dish-actions {
    grid-template-columns: 1fr;
  }
}
</style>
