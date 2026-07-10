<template>
  <div class="meal-shell" :style="pageVars">
    <div class="meal-warm-overlay" aria-hidden="true"></div>

    <main class="meal-phone">
      <header class="meal-topbar">
        <button class="nav-btn" type="button" aria-label="返回" @click="router.back()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6" />
          </svg>
        </button>
        <h1>一周菜单</h1>
        <button class="nav-btn" type="button" aria-label="生成菜单" :disabled="loading" @click="generateWeekMenu">
          <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 5v14" />
            <path d="M5 12h14" />
          </svg>
          <span v-else class="mini-spinner" aria-hidden="true"></span>
        </button>
      </header>

      <nav class="week-strip" aria-label="选择日期">
        <button
          v-for="(day, idx) in weekDays"
          :key="day.key"
          class="day"
          :class="{ active: activeDayIdx === idx }"
          type="button"
          @click="activeDayIdx = idx"
        >
          <small>{{ day.label }}</small>
          <strong>{{ day.date }}</strong>
          <span v-if="weekMenuData[idx]" class="day-dot" aria-hidden="true"></span>
        </button>
      </nav>

      <section class="template-strip" aria-label="常用菜单模板">
        <button
          v-for="template in menuTemplates"
          :key="template.key"
          type="button"
          :class="{ active: selectedTemplateKey === template.key }"
          @click="selectedTemplateKey = template.key"
        >
          <strong>{{ template.title }}</strong>
          <span>{{ template.description }}</span>
        </button>
      </section>

      <section class="glass-card smart-card">
        <div class="smart-content">
          <span class="smart-badge">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 3 13.7 8.3 19 10l-5.3 1.7L12 17l-1.7-5.3L5 10l5.3-1.7L12 3Z" />
              <path d="M19 16v4M21 18h-4" />
            </svg>
            智能周菜单
          </span>
          <h2>{{ weekMenuData.length ? `${activeDayName}这样吃` : '一键安排三餐' }}</h2>
        </div>
        <svg class="leaf-mark" viewBox="0 0 48 32" fill="currentColor" aria-hidden="true">
          <path d="M24 29C19 17 9 18 3 21c3 7 11 11 21 8Z" />
          <path d="M25 28C26 15 35 8 45 7c-1 10-8 19-20 21Z" />
        </svg>
        <div class="stat-pill">
          <div class="stat">
            <svg class="dish-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 17h16" />
              <path d="M6 17a6 6 0 0 1 12 0" />
              <path d="M12 6v2" />
              <path d="M5 20h14" />
            </svg>
            <span><strong>{{ currentDishCount }}</strong>道菜</span>
          </div>
          <i aria-hidden="true"></i>
          <div class="stat">
            <svg class="basket-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 10h12l-1.2 10H7.2L6 10Z" />
              <path d="M9 10a3 3 0 0 1 6 0" />
              <path d="M9 14v2M12 14v2M15 14v2" />
            </svg>
            <span><strong>{{ currentDayIngredients.length }}</strong>食材</span>
          </div>
        </div>
        <button
          v-if="weekMenuData.length"
          class="week-shopping-btn"
          type="button"
          :disabled="!allWeekRecipeIds.length && !allWeekIngredients.length || shoppingSaving"
          @click="createWeekShoppingList"
        >
          {{ shoppingSaving ? '生成中' : '生成整周采购清单' }}
        </button>
        <button
          v-if="weekMenuData.length"
          class="week-shopping-btn save-menu-btn"
          type="button"
          :disabled="menuSaving"
          @click="saveWeekMenu"
        >
          {{ menuSaving ? '保存中' : '保存本周菜单' }}
        </button>
      </section>

      <section v-if="loading" class="glass-card state-card">
        <div class="loading-spinner"></div>
        <h2>正在生成周菜单</h2>
      </section>

      <section v-else-if="error" class="glass-card state-card">
        <div class="empty-art error-art" aria-hidden="true">
          <span class="book"><svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 2 14 9.5 22 12l-8 2.5L12 22l-2-7.5L2 12l8-2.5L12 2Z" /></svg></span>
          <span class="bowl"></span>
          <span class="chopstick"></span>
        </div>
        <h2>生成失败</h2>
        <p>{{ error }}</p>
        <button class="primary-btn" type="button" @click="generateWeekMenu">
          重新生成
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </button>
      </section>

      <section v-else-if="!weekMenuData.length" class="glass-card state-card">
        <EmptyArt />
        <h2>还没有周菜单</h2>
        <button class="primary-btn" type="button" @click="generateWeekMenu">
          生成一周菜单
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </button>
      </section>

      <template v-else>
        <section class="glass-card generated-card">
          <div class="generated-head">
            <div>
              <span>{{ activeDayName }}三餐</span>
              <h2>{{ currentDishCount ? '今日安排好了' : '今日暂无菜品' }}</h2>
            </div>
            <button
              class="shopping-btn"
              type="button"
              :disabled="!currentDayIngredients.length || shoppingSaving"
              @click="createDayShoppingList"
            >
              {{ shoppingSaving ? '生成中' : '加入采购清单' }}
            </button>
          </div>

          <p v-if="shoppingMessage" class="inline-message">{{ shoppingMessage }}</p>
          <p v-if="isGeneratedButEmpty" class="data-note">没有返回菜品，可重新生成。</p>

          <div class="meal-plan-list">
            <article v-for="meal in normalizedCurrentMeals" :key="meal.type" class="plan-meal-card">
              <header>
                <span class="meal-icon" :class="meal.type">
                  <MealIcon :type="meal.type" />
                </span>
                <div>
                  <strong>{{ mealTypeLabel(meal.type) }}</strong>
                  <small>{{ mealTime(meal.type) }}</small>
                </div>
              </header>

              <template v-if="meal.dishes.length">
                <div v-for="dish in meal.dishes" :key="dish.recipe_id || dish.name" class="dish-row">
                  <div>
                    <h3>{{ dish.name || '未命名菜品' }}</h3>
                    <p>{{ dish.cook_time || '--' }} 分钟 · {{ dish.difficulty || '家常' }}</p>
                  </div>
                  <button type="button" :disabled="!dish.recipe_id" @click="openRecipe(dish.recipe_id)">做法</button>
                </div>
                <div class="ingredient-row" v-if="mealIngredients(meal).length">
                  <span v-for="item in mealIngredients(meal).slice(0, 5)" :key="item">{{ item }}</span>
                </div>
              </template>

              <p v-else class="meal-empty-text">暂无菜品</p>
            </article>
          </div>
        </section>
      </template>

    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index'
import { saveUserMenu } from '@/api/menu'
import { getPreferences } from '@/api/user'
import { useShoppingStore } from '@/stores/shopping'
import type { ShoppingItem } from '@/api/shopping'
import type { UserPreferences } from '@/types/user'
import { trackEvent } from '@/utils/event'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import smartMealImage from '@/assets/home/couple-dining.jpg'

type MealType = 'breakfast' | 'lunch' | 'dinner' | string

interface WeekDish {
  recipe_id?: number
  name?: string
  type?: string
  cook_time?: number
  difficulty?: string
  ingredients?: string[]
}

interface WeekMeal {
  type: MealType
  reason?: string
  dishes: WeekDish[]
  shopping_list: string[]
}

interface WeekDay {
  day?: number
  meals: WeekMeal[]
}

interface WeekDate {
  key: string
  label: string
  date: string
}

const router = useRouter()
const shoppingStore = useShoppingStore()
const activeDayIdx = ref(new Date().getDay() === 0 ? 6 : new Date().getDay() - 1)
const loading = ref(false)
const error = ref('')
const weekMenuData = ref<WeekDay[]>([])
const shoppingSaving = ref(false)
const menuSaving = ref(false)
const shoppingMessage = ref('')
const userPreferences = ref<UserPreferences | null>(null)
const selectedTemplateKey = ref('weekday')

const menuTemplates = [
  { key: 'weekday', title: '工作日快手餐', description: '30 分钟内', health_goal: '', people_count: 2, cook_time_preference: '30分钟内' },
  { key: 'weekend', title: '周末聚餐', description: '多 1 道菜', health_goal: '', people_count: 4, cook_time_preference: '都可以' },
  { key: 'fat_loss', title: '减脂周菜单', description: '轻食高蛋白', health_goal: '减脂', people_count: 1, cook_time_preference: '30分钟内' },
  { key: 'kids', title: '儿童营养餐', description: '少辣均衡', health_goal: '儿童营养', people_count: 3, cook_time_preference: '45分钟内' },
]

const selectedTemplate = computed(() => menuTemplates.find((item) => item.key === selectedTemplateKey.value) || menuTemplates[0])

const weekDays = ref(buildWeekDays())

const pageVars = computed(() => ({
  '--meal-bg': `url(${kitchenBg})`,
  '--smart-meal-img': `url(${smartMealImage})`,
}))

const currentDayMeals = computed<WeekMeal[]>(() => {
  const day = weekMenuData.value[activeDayIdx.value]
  if (!day?.meals) return []
  return day.meals
})

const normalizedCurrentMeals = computed<WeekMeal[]>(() => {
  const map = new Map(currentDayMeals.value.map((meal) => [meal.type, meal]))
  return ['breakfast', 'lunch', 'dinner'].map((type) => map.get(type) || {
    type,
    dishes: [],
    shopping_list: [],
  })
})

const activeDayName = computed(() => `周${weekDays.value[activeDayIdx.value]?.label || ''}`)

const currentDishCount = computed(() => {
  return currentDayMeals.value.reduce((sum, meal) => sum + (meal.dishes?.length || 0), 0)
})

const currentDayIngredients = computed(() => {
  const names = new Set<string>()
  currentDayMeals.value.forEach((meal) => {
    ;(meal.shopping_list || []).forEach((item) => item && names.add(item))
    ;(meal.dishes || []).forEach((dish) => {
      ;(dish.ingredients || []).forEach((item) => item && names.add(item))
    })
  })
  return Array.from(names)
})

const allWeekIngredients = computed(() => {
  const names = new Set<string>()
  weekMenuData.value.forEach((day) => {
    ;(day.meals || []).forEach((meal) => {
      ;(meal.shopping_list || []).forEach((item) => item && names.add(item))
      ;(meal.dishes || []).forEach((dish) => {
        ;(dish.ingredients || []).forEach((item) => item && names.add(item))
      })
    })
  })
  return Array.from(names)
})

const allWeekRecipeIds = computed(() => {
  const ids = new Set<number>()
  weekMenuData.value.forEach((day) => {
    ;(day.meals || []).forEach((meal) => {
      ;(meal.dishes || []).forEach((dish) => {
        const id = Number(dish.recipe_id || 0)
        if (Number.isFinite(id) && id > 0) ids.add(id)
      })
    })
  })
  return Array.from(ids)
})

const isGeneratedButEmpty = computed(() => {
  return weekMenuData.value.length > 0 && weekMenuData.value.every((day) => {
    return (day.meals || []).every((meal) => !meal.dishes?.length)
  })
})

function buildWeekDays(): WeekDate[] {
  const labels = ['一', '二', '三', '四', '五', '六', '日']
  const now = new Date()
  const dayOfWeek = now.getDay() === 0 ? 7 : now.getDay()
  const monday = new Date(now)
  monday.setDate(now.getDate() - dayOfWeek + 1)
  return labels.map((label, index) => {
    const d = new Date(monday)
    d.setDate(monday.getDate() + index)
    return {
      key: `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`,
      label,
      date: String(d.getDate()),
    }
  })
}

function mealTime(type: string) {
  const map: Record<string, string> = { breakfast: '07:30', lunch: '12:00', dinner: '18:30' }
  return map[type] || ''
}

function mealTypeLabel(type: string) {
  const map: Record<string, string> = { breakfast: '早餐', lunch: '午餐', dinner: '晚餐' }
  return map[type] || type
}

function mealIngredients(meal: WeekMeal) {
  const names = new Set<string>()
  ;(meal.shopping_list || []).forEach((item) => item && names.add(item))
  ;(meal.dishes || []).forEach((dish) => {
    ;(dish.ingredients || []).forEach((item) => item && names.add(item))
  })
  return Array.from(names)
}

function normalizeWeekMenu(data: any[]): WeekDay[] {
  const mealTypes = ['breakfast', 'lunch', 'dinner']
  return data.map((day) => ({
    ...day,
    meals: (day.meals || []).map((meal: any, idx: number) => ({
      ...meal,
      type: meal.type || meal.meal_type || mealTypes[idx] || 'lunch',
      dishes: Array.isArray(meal.dishes) ? meal.dishes : [],
      shopping_list: Array.isArray(meal.shopping_list) ? meal.shopping_list : [],
    })),
  }))
}

function openRecipe(id?: number) {
  if (!id) return
  router.push(`/recipes/${id}`)
}

function inferCategory(name: string) {
  if (/鸡|鸭|鱼|虾|肉|牛|猪|蛋|排骨/.test(name)) return '肉蛋水产'
  if (/米|面|粥|饭|馒头|粉/.test(name)) return '主食'
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
    trackEvent({
      event_name: 'add_shopping_list',
      entity_type: 'weekly_day',
      payload: { day: activeDayName.value, ingredients: currentDayIngredients.value },
    })
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

async function createWeekShoppingList() {
  if (shoppingSaving.value) return
  shoppingSaving.value = true
  shoppingMessage.value = ''
  try {
    if (allWeekRecipeIds.value.length) {
      await shoppingStore.generateByRecipes(allWeekRecipeIds.value, '本周菜单采购清单')
    } else if (allWeekIngredients.value.length) {
      await shoppingStore.appendItemsToCurrentList('本周菜单采购清单', toShoppingItems(allWeekIngredients.value))
    } else {
      shoppingMessage.value = '当前周菜单还没有可生成采购清单的食材。'
      return
    }
    trackEvent({
      event_name: 'add_shopping_list',
      entity_type: 'weekly_menu',
      payload: { recipe_ids: allWeekRecipeIds.value, ingredients: allWeekIngredients.value },
    })
    shoppingMessage.value = '整周采购清单已合并，可以去购物清单页面查看。'
    setTimeout(() => {
      router.push('/shopping-list')
    }, 350)
  } catch (e: any) {
    shoppingMessage.value = e?.message || '整周采购清单生成失败，请稍后再试。'
  } finally {
    shoppingSaving.value = false
  }
}

async function saveWeekMenu() {
  if (!weekMenuData.value.length || menuSaving.value) return
  menuSaving.value = true
  shoppingMessage.value = ''
  try {
    await saveUserMenu({
      name: `${selectedTemplate.value.title} - 本周菜单`,
      menu_type: 'weekly',
      meal_type: 'weekly',
      people_count: selectedTemplate.value.people_count,
      health_goal: selectedTemplate.value.health_goal,
      dishes: weekMenuData.value,
      shopping_list: allWeekIngredients.value,
      reason: `基于${selectedTemplate.value.title}模板生成`,
    })
    trackEvent({
      event_name: 'save_menu',
      entity_type: 'weekly_menu',
      payload: { template: selectedTemplate.value.key, recipe_ids: allWeekRecipeIds.value },
    })
    shoppingMessage.value = '本周菜单已保存到“我的菜单”。'
  } catch (e: any) {
    shoppingMessage.value = e?.message || '保存菜单失败，请稍后再试。'
  } finally {
    menuSaving.value = false
  }
}

async function generateWeekMenu() {
  loading.value = true
  error.value = ''
  shoppingMessage.value = ''
  try {
    if (!userPreferences.value) {
      await loadUserPreferences()
    }
    const pref = userPreferences.value
    const template = selectedTemplate.value
    const res: any = await api.post('/recommend/week-menu', {
      people_count: template.people_count || pref?.default_servings || 2,
      meal_type: 'lunch',
      taste_preference: pref?.taste_preference || [],
      health_goal: template.health_goal || pref?.health_goal || '',
      avoid_ingredients: pref?.avoid_ingredients || [],
      existing_ingredients: [],
      cook_time_preference: template.cook_time_preference || pref?.cook_time_preference || '',
    })
    const list = Array.isArray(res) ? normalizeWeekMenu(res) : []
    trackEvent({
      event_name: 'recommend_start',
      entity_type: 'weekly_menu',
      payload: { source: 'week_menu', template: template.key },
    })
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

async function loadUserPreferences() {
  try {
    const pref: any = await getPreferences()
    userPreferences.value = {
      taste_preference: Array.isArray(pref?.taste_preference) ? pref.taste_preference : [],
      health_goal: pref?.health_goal || '',
      avoid_ingredients: Array.isArray(pref?.avoid_ingredients) ? pref.avoid_ingredients : [],
      favorite_ingredients: Array.isArray(pref?.favorite_ingredients) ? pref.favorite_ingredients : [],
      cook_time_preference: pref?.cook_time_preference || '',
      default_servings: pref?.people_count || pref?.default_servings || 2,
    }
    shoppingMessage.value = '已读取你的口味、忌口和烹饪时间偏好。'
  } catch {
    userPreferences.value = null
    shoppingMessage.value = '偏好读取失败，本次将使用默认参数生成。'
  }
}

const MealIcon = defineComponent({
  props: { type: { type: String, required: true } },
  setup(props) {
    return () => {
      if (props.type === 'breakfast') {
        return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-linecap': 'round', 'stroke-linejoin': 'round' }, [
          h('circle', { cx: 12, cy: 12, r: 4 }),
          h('path', { d: 'M12 2v3M12 19v3M4.9 4.9 7 7M17 17l2.1 2.1M2 12h3M19 12h3M4.9 19.1 7 17M17 7l2.1-2.1' }),
        ])
      }
      if (props.type === 'lunch') {
        return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-linecap': 'round', 'stroke-linejoin': 'round' }, [
          h('path', { d: 'M12 3v3M5.6 5.6 7.8 7.8M3 12h3M18 12h3M16.2 7.8l2.2-2.2' }),
          h('path', { d: 'M7 17a5 5 0 0 1 10 0' }),
          h('path', { d: 'M4 17h16M6 20h12' }),
        ])
      }
      return h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-linecap': 'round', 'stroke-linejoin': 'round' }, [
        h('path', { d: 'M20 15.3A8 8 0 0 1 8.7 4 8.3 8.3 0 1 0 20 15.3Z' }),
        h('path', { d: 'M18 4v3M19.5 5.5h-3' }),
      ])
    }
  },
})

const EmptyArt = defineComponent({
  setup() {
    return () => h('div', { class: 'empty-art', 'aria-hidden': 'true' }, [
      h('span', { class: 'sprout' }, [
        h('svg', { viewBox: '0 0 32 52', fill: 'none', stroke: 'currentColor', 'stroke-width': 2, 'stroke-linecap': 'round' }, [
          h('path', { d: 'M16 49V22' }),
          h('path', { d: 'M16 25C8 16 2 17 1 18c3 8 9 12 15 7Z', fill: 'currentColor', opacity: '.32' }),
          h('path', { d: 'M16 21C19 10 27 8 31 9c-2 8-8 13-15 12Z', fill: 'currentColor', opacity: '.42' }),
        ]),
      ]),
      h('span', { class: 'book' }, [
        h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
          h('path', { d: 'M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z' }),
        ]),
      ]),
      h('span', { class: 'bowl' }),
      h('span', { class: 'chopstick' }),
      h('span', { class: 'spark' }, [
        h('svg', { viewBox: '0 0 24 24', fill: 'currentColor' }, [
          h('path', { d: 'M12 2 14 9.5 22 12l-8 2.5L12 22l-2-7.5L2 12l8-2.5L12 2Z' }),
        ]),
      ]),
    ])
  },
})

onMounted(() => {
  loadUserPreferences()
})
</script>

<style scoped>
.meal-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.76);
  --coral: #e95645;
  --orange: #e89a45;
  --sage: #7ea36a;
  --border: rgba(255, 255, 255, 0.62);
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 239, 214, 0.5), rgba(231, 171, 105, 0.2)),
    var(--meal-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.meal-warm-overlay {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 50% 9%, rgba(255, 255, 255, 0.76), transparent 34%),
    radial-gradient(circle at 84% 30%, rgba(255, 234, 202, 0.2), transparent 34%),
    linear-gradient(180deg, rgba(255, 243, 226, 0.42), rgba(150, 82, 30, 0.08));
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.meal-phone {
  position: relative;
  z-index: 1;
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: max(49px, env(safe-area-inset-top)) 22px calc(var(--tab-h, 82px) + 74px);
}

.meal-topbar {
  height: 56px;
  display: grid;
  grid-template-columns: 56px 1fr 56px;
  align-items: center;
  gap: 10px;
}

.meal-topbar h1 {
  margin: 0;
  color: #2f211b;
  font-size: 25px;
  font-weight: 950;
  line-height: 1.1;
  text-align: center;
  letter-spacing: 0;
}

.nav-btn {
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.74);
  border-radius: 16px;
  color: #261c18;
  background: rgba(255, 250, 240, 0.92);
  box-shadow: 0 12px 26px rgba(80, 50, 30, 0.14);
  cursor: pointer;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.nav-btn:disabled {
  opacity: 0.72;
  cursor: wait;
}

.nav-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.4;
}

.week-strip {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  align-items: center;
  gap: 8px;
  margin-top: 24px;
}

.template-strip {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-top: 12px;
}

.template-strip button {
  min-height: 58px;
  display: grid;
  gap: 3px;
  justify-items: start;
  padding: 10px 12px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 8px;
  color: var(--text);
  background: rgba(255, 250, 240, 0.72);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.72);
}

.template-strip button.active {
  border-color: rgba(233, 86, 69, 0.34);
  background: rgba(252, 226, 214, 0.74);
}

.template-strip strong {
  font-size: 13px;
  font-weight: 900;
}

.template-strip span {
  color: var(--sub);
  font-size: 11px;
  font-weight: 760;
}

.day {
  position: relative;
  min-width: 0;
  height: 74px;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 7px;
  border: 0;
  border-radius: 18px;
  color: #2e241f;
  background: transparent;
  cursor: pointer;
  transition: transform 180ms ease, background 180ms ease, color 180ms ease, box-shadow 180ms ease;
}

.day small {
  font-size: 15px;
  font-weight: 800;
  line-height: 1;
  color: rgba(46, 36, 31, 0.64);
}

.day strong {
  font-size: 24px;
  font-weight: 950;
  line-height: 1;
}

.day.active {
  color: #fff;
  background: linear-gradient(135deg, #f26356, #e93f35);
  box-shadow: 0 16px 28px rgba(233, 86, 69, 0.28);
}

.day.active small {
  color: #fff;
}

.day-dot {
  position: absolute;
  bottom: 7px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.8;
}

.glass-card {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--border);
  border-radius: 30px;
  background: var(--cream);
  box-shadow: 0 18px 40px rgba(80, 50, 30, 0.16);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.smart-card {
  min-height: 210px;
  margin-top: 33px;
  padding: 29px 22px 23px;
}

.smart-card::before {
  content: "";
  position: absolute;
  inset: 0;
  z-index: 1;
  background:
    radial-gradient(circle at 12% 44%, rgba(255, 255, 255, 0.76), transparent 42%),
    linear-gradient(90deg, rgba(255, 250, 240, 0.98) 0%, rgba(255, 250, 240, 0.95) 42%, rgba(255, 250, 240, 0.7) 62%, rgba(255, 250, 240, 0.2) 100%);
}

.smart-card::after {
  content: "";
  position: absolute;
  inset: 0;
  z-index: 0;
  background: var(--smart-meal-img) 28% center / cover no-repeat;
  opacity: 0.72;
}

.smart-content {
  position: relative;
  z-index: 2;
  max-width: 62%;
}

.smart-badge {
  width: max-content;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  height: 31px;
  padding: 0 13px;
  border-radius: 999px;
  color: var(--coral);
  background: rgba(255, 255, 255, 0.56);
  font-size: 14px;
  font-weight: 850;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.86);
}

.smart-badge svg {
  width: 17px;
  height: 17px;
  fill: rgba(233, 86, 69, 0.16);
  stroke-width: 2.2;
}

.smart-card h2 {
  margin: 27px 0 0;
  color: #30221d;
  font-size: clamp(35px, 9vw, 43px);
  font-weight: 950;
  line-height: 1.1;
  letter-spacing: 0;
}

.leaf-mark {
  position: absolute;
  left: 53%;
  top: 48%;
  z-index: 2;
  width: 38px;
  height: 28px;
  color: rgba(126, 163, 106, 0.36);
  transform: rotate(-14deg);
}

.stat-pill {
  position: relative;
  z-index: 2;
  width: min(100%, 210px);
  height: 62px;
  display: grid;
  grid-template-columns: 1fr 1px 1fr;
  align-items: center;
  margin-top: 22px;
  border-radius: 21px;
  background: rgba(255, 250, 240, 0.74);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.82);
}

.week-shopping-btn {
  position: relative;
  z-index: 2;
  width: min(100%, 236px);
  min-height: 48px;
  margin-top: 18px;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: linear-gradient(135deg, #ef6153, #e94c3c);
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.22);
  font-size: 15px;
  font-weight: 900;
  cursor: pointer;
}

.week-shopping-btn:disabled {
  cursor: not-allowed;
  opacity: 0.58;
  box-shadow: none;
}

.save-menu-btn {
  margin-top: 8px;
  color: var(--coral);
  background: rgba(255, 250, 240, 0.82);
}

.stat-pill i {
  width: 1px;
  height: 34px;
  background: rgba(142, 111, 91, 0.18);
}

.stat {
  display: grid;
  grid-template-columns: 38px 1fr;
  align-items: center;
  gap: 9px;
  padding: 0 15px;
  color: #6f5e54;
}

.stat svg {
  width: 34px;
  height: 34px;
  stroke-width: 1.9;
}

.dish-icon {
  color: var(--coral);
}

.basket-icon {
  color: var(--sage);
}

.stat strong {
  display: block;
  color: #2e241f;
  font-size: 22px;
  font-weight: 950;
  line-height: 1;
}

.stat span {
  display: block;
  margin-top: 4px;
  font-size: 14px;
  font-weight: 760;
  white-space: nowrap;
}

.state-card {
  position: relative;
  isolation: isolate;
  overflow: hidden;
  min-height: 240px;
  display: grid;
  place-items: center;
  align-content: center;
  margin-top: 24px;
  padding: 28px 24px;
  text-align: center;
}

.empty-art {
  position: relative;
  z-index: 1;
  width: 138px;
  height: 96px;
  margin: 0 0 18px;
  overflow: hidden;
  pointer-events: none;
}

.empty-art :deep(*) {
  box-sizing: border-box;
}

.empty-art :deep(svg) {
  display: block;
}

.empty-art :deep(.book) {
  position: absolute;
  left: 28px;
  top: 14px;
  width: 58px;
  height: 62px;
  border: 2px solid rgba(188, 128, 78, 0.52);
  border-radius: 8px;
  background: linear-gradient(135deg, #ffe7c8, #fff9ef);
  box-shadow: 0 8px 15px rgba(120, 73, 32, 0.12);
  transform: rotate(-6deg);
}

.empty-art :deep(.book)::before {
  content: "";
  position: absolute;
  left: 7px;
  top: 0;
  bottom: 0;
  border-left: 2px solid rgba(155, 103, 67, 0.28);
}

.empty-art :deep(.book svg) {
  position: absolute;
  right: 11px;
  top: 17px;
  width: 21px;
  height: 21px;
  color: #ef8061;
}

.empty-art :deep(.bowl) {
  position: absolute;
  right: 18px;
  top: 34px;
  width: 61px;
  height: 36px;
  border-radius: 10px 10px 28px 28px;
  background: linear-gradient(180deg, #fffdf7, #f6d39d);
  box-shadow: inset 0 -8px 0 rgba(198, 124, 56, 0.15), 0 10px 15px rgba(124, 74, 35, 0.14);
}

.empty-art :deep(.bowl)::before {
  content: "";
  position: absolute;
  left: 8px;
  right: 8px;
  top: -6px;
  height: 16px;
  border: 2px solid rgba(191, 138, 87, 0.22);
  border-radius: 50%;
  background: radial-gradient(circle at 50% 46%, #8fb170 0 3px, transparent 4px), #fff8e9;
}

.empty-art :deep(.chopstick) {
  position: absolute;
  right: 9px;
  top: 21px;
  width: 62px;
  height: 3px;
  border-radius: 999px;
  background: #c98e4d;
  transform: rotate(-38deg);
}

.empty-art :deep(.sprout) {
  position: absolute;
  left: 5px;
  bottom: 10px;
  width: 31px;
  height: 48px;
  color: #80a65d;
}

.empty-art :deep(.sprout svg),
.empty-art :deep(.spark svg) {
  width: 100%;
  height: 100%;
}

.empty-art :deep(.spark) {
  position: absolute;
  right: 2px;
  top: 10px;
  width: 18px;
  height: 18px;
  color: #e8a657;
}

.state-card h2 {
  position: relative;
  z-index: 2;
  margin: 0 0 13px;
  color: #30221d;
  font-size: 30px;
  font-weight: 950;
  line-height: 1.16;
}

.primary-btn {
  position: relative;
  z-index: 2;
  height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 11px;
  margin-top: 18px;
  padding: 0 35px;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: linear-gradient(135deg, #ef6153, #e94c3c);
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.25);
  font-size: 19px;
  font-weight: 900;
  cursor: pointer;
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.primary-btn svg {
  width: 22px;
  height: 22px;
  stroke-width: 2.4;
}

.generated-card {
  margin-top: 24px;
  padding: 18px 16px 19px;
}

.generated-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.generated-head span {
  color: var(--coral);
  font-size: 13px;
  font-weight: 850;
}

.generated-head h2 {
  margin: 0;
  color: #2f211b;
  font-size: 22px;
  font-weight: 950;
  line-height: 1.2;
}

.shopping-btn {
  min-height: 40px;
  flex: 0 0 auto;
  padding: 0 14px;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: linear-gradient(135deg, #ef6153, #e94c3c);
  font-size: 13px;
  font-weight: 850;
  cursor: pointer;
}

.shopping-btn:disabled {
  opacity: 0.52;
  cursor: not-allowed;
}

.inline-message,
.data-note {
  margin: 0 0 12px;
  padding: 10px 12px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.56);
  color: #8a5b3d;
  font-size: 13px;
  font-weight: 760;
  line-height: 1.5;
}

.data-note {
  color: #9a5c24;
  background: rgba(255, 236, 195, 0.62);
}

.meal-plan-list {
  display: grid;
  gap: 12px;
}

.plan-meal-card {
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.54);
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.48);
  box-shadow: 0 10px 22px rgba(80, 50, 30, 0.08);
}

.plan-meal-card header {
  min-height: 52px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  color: #3b2a22;
}

.meal-icon {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  color: #f3a12f;
}

.meal-icon svg {
  width: 22px;
  height: 22px;
  stroke-width: 2.1;
}

.meal-icon.dinner {
  color: #d48b3d;
}

.plan-meal-card strong {
  display: block;
  color: #3b2a22;
  font-size: 17px;
  font-weight: 900;
}

.plan-meal-card small {
  color: #8a7a70;
  font-size: 12px;
  font-weight: 800;
}

.dish-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 12px;
  border-top: 1px solid rgba(141, 112, 92, 0.12);
}

.dish-row h3 {
  margin: 0;
  color: #372720;
  font-size: 16px;
  font-weight: 900;
  line-height: 1.2;
}

.dish-row p {
  margin: 5px 0 0;
  color: #7a6659;
  font-size: 13px;
  font-weight: 700;
}

.dish-row button {
  min-width: 54px;
  height: 34px;
  border: 0;
  border-radius: 999px;
  color: var(--coral);
  background: rgba(255, 246, 238, 0.82);
  font-size: 13px;
  font-weight: 850;
  cursor: pointer;
}

.dish-row button:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.ingredient-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 0 12px 12px;
}

.ingredient-row span {
  padding: 4px 8px;
  border-radius: 999px;
  color: #6d7653;
  background: rgba(226, 235, 203, 0.72);
  font-size: 12px;
  font-weight: 800;
}

.meal-empty-text {
  margin: 0;
  padding: 0 12px 14px 50px;
  color: #86746a;
  font-size: 13px;
  line-height: 1.55;
}

.loading-spinner,
.mini-spinner {
  border-radius: 50%;
  animation: spin 0.72s linear infinite;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  margin-bottom: 16px;
  border: 3px solid rgba(233, 86, 69, 0.18);
  border-top-color: var(--coral);
}

.mini-spinner {
  width: 24px;
  height: 24px;
  border: 3px solid rgba(46, 36, 31, 0.12);
  border-top-color: #2e241f;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.nav-btn:active,
.day:active,
.glass-card:active,
.primary-btn:active,
.shopping-btn:active,
.week-shopping-btn:active,
.dish-row button:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .nav-btn:hover,
  .day:hover,
  .primary-btn:hover,
  .shopping-btn:hover,
  .week-shopping-btn:hover:not(:disabled) {
    transform: translateY(-1px);
  }

  .primary-btn:hover {
    box-shadow: 0 15px 28px rgba(233, 86, 69, 0.28);
  }
}

@media (max-width: 380px) {
  .meal-phone {
    padding-left: 18px;
    padding-right: 18px;
  }

  .week-strip {
    gap: 5px;
  }

  .day strong {
    font-size: 22px;
  }

  .smart-card {
    padding-left: 20px;
    padding-right: 20px;
  }

  .smart-content {
    max-width: 66%;
  }

  .smart-card h2 {
    font-size: 34px;
  }

  .stat-pill {
    width: 198px;
  }

}

@media (max-width: 350px) {
  .smart-content {
    max-width: 74%;
  }

  .smart-card::before {
    background:
      radial-gradient(circle at 12% 44%, rgba(255, 255, 255, 0.78), transparent 42%),
      linear-gradient(90deg, rgba(255, 250, 240, 0.98), rgba(255, 250, 240, 0.76));
  }

  .smart-card::after {
    opacity: 0.42;
  }
}

@media (min-width: 431px) {
  .meal-shell {
    background-color: #ead7bd;
  }

  .meal-phone {
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.22);
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
</style>
