<template>
  <div class="recommend-shell" :style="pageVars">
    <div class="recommend-warm-overlay" aria-hidden="true"></div>

    <main class="recommend-phone">
      <header class="recommend-topbar">
        <button class="nav-btn" type="button" aria-label="返回" @click="router.back()">
          <svg viewBox="0 0 24 24"><path d="m15 18-6-6 6-6" /></svg>
        </button>
        <div>
          <h1>{{ modeMeta.title }}</h1>
        </div>
      </header>

      <section class="hero-card">
        <span class="hero-badge">
          <svg viewBox="0 0 24 24"><path d="M12 3 13.7 8.3 19 10l-5.3 1.7L12 17l-1.7-5.3L5 10l5.3-1.7L12 3Z" /></svg>
          {{ modeMeta.badge }}
        </span>
        <h2>{{ modeMeta.heading }}</h2>
      </section>

      <section class="intent-panel" aria-label="推荐入口">
        <button
          v-for="intent in intentOptions"
          :key="intent.mode"
          type="button"
          :class="{ active: mode === intent.mode }"
          @click="goIntent(intent)"
        >
          <strong>{{ intent.title }}</strong>
        </button>
      </section>

      <section class="preference-panel" aria-label="本次推荐偏好">
        <div class="panel-title">
          <div>
            <h2>{{ preferenceSummary }}</h2>
          </div>
          <button type="button" @click="router.push('/user/preferences')">偏好</button>
        </div>
        <div class="preference-controls">
          <label>
            <span>人数</span>
            <select v-model.number="temporaryPreference.people_count">
              <option v-for="count in peopleOptions" :key="count" :value="count">{{ count }}人</option>
            </select>
          </label>
          <label>
            <span>做饭时间</span>
            <select v-model="temporaryPreference.cook_time_preference">
              <option value="">不限</option>
              <option v-for="item in cookTimeOptions" :key="item" :value="item">{{ item }}</option>
            </select>
          </label>
          <label>
            <span>健康目标</span>
            <select v-model="temporaryPreference.health_goal">
              <option value="">不限</option>
              <option v-for="item in healthGoalOptions" :key="item" :value="item">{{ item }}</option>
            </select>
          </label>
        </div>
        <div v-if="userPreference.taste_preference.length || userPreference.avoid_ingredients.length" class="preference-tags">
          <span v-for="taste in userPreference.taste_preference.slice(0, 4)" :key="'taste-' + taste">{{ taste }}</span>
          <span v-for="avoid in userPreference.avoid_ingredients.slice(0, 3)" :key="'avoid-' + avoid">不吃{{ avoid }}</span>
        </div>
      </section>

      <section v-if="isIngredientMode" class="panel">
        <div class="panel-title">
          <h2>{{ mode === 'fridge' ? '录入冰箱里的食材' : '选择已有食材' }}</h2>
          <button v-if="selectedIngredients.length" type="button" @click="clearIngredients">清空</button>
        </div>

        <form class="ingredient-input" @submit.prevent="addCustomIngredient">
          <input
            v-model.trim="ingredientKeyword"
            type="search"
            :placeholder="mode === 'fridge' ? '例如：米饭、鸡蛋、青菜' : '搜索或手动输入食材'"
            aria-label="输入食材"
            @focus="loadIngredientOptions"
          />
          <button type="submit" :disabled="!ingredientKeyword || selectedIngredients.length >= 20">添加</button>
        </form>

        <div class="chip-cloud" aria-label="已选择食材">
          <button
            v-for="name in selectedIngredients"
            :key="name"
            class="selected-chip"
            type="button"
            @click="removeIngredient(name)"
          >
            {{ name }}
            <svg viewBox="0 0 24 24"><path d="M18 6 6 18M6 6l12 12" /></svg>
          </button>
          <span v-if="!selectedIngredients.length" class="hint-chip">添加食材</span>
        </div>

        <div class="option-grid">
          <button
            v-for="item in visibleIngredientOptions"
            :key="item.id || item.name"
            type="button"
            class="option-chip"
            :class="{ active: selectedIngredients.includes(item.name) }"
            @click="toggleIngredient(item.name)"
          >
            <span>{{ item.name }}</span>
            <small>{{ item.category || '常用' }}</small>
          </button>
        </div>

        <p v-if="!ingredientOptions.length && !ingredientOptionsLoading" class="data-empty">
          数据库暂无可选食材。你仍可手动输入食材，推荐结果只会来自真实菜谱数据。
        </p>

        <button class="submit-btn" type="button" :disabled="!selectedIngredients.length || loading" @click="submitIngredients">
          <span v-if="loading" class="mini-spinner" aria-hidden="true"></span>
          <span>{{ loading ? '正在推荐...' : mode === 'fridge' ? '看看能做什么' : '开始推荐' }}</span>
        </button>
      </section>

      <section v-else-if="mode === 'taste'" class="panel">
        <div class="panel-title">
          <h2>选择今天想吃的味道</h2>
          <button type="button" @click="router.push('/recipes')">看全部</button>
        </div>
        <div v-if="tasteOptions.length" class="taste-grid">
          <button v-for="taste in tasteOptions" :key="taste" type="button" @click="chooseTaste(taste)">
            <strong>{{ taste }}</strong>
          </button>
        </div>
        <p v-else class="data-empty">数据库暂无口味选项，请先在后台维护菜谱口味字段。</p>
      </section>

      <section v-else-if="mode === 'scene' || mode === 'new'" class="panel">
        <div class="panel-title">
          <h2>{{ mode === 'new' ? '换点没吃过的新菜' : '按生活场景安排' }}</h2>
          <button type="button" @click="clearSceneResult">重选</button>
        </div>
        <div class="scene-list">
          <button
            v-for="scene in scenes"
            :key="scene.key"
            type="button"
            :class="{ active: selectedScene?.key === scene.key }"
            @click="chooseScene(scene)"
          >
            <strong>{{ scene.title }}</strong>
          </button>
        </div>
        <div class="scene-actions">
          <button v-if="mode === 'scene'" class="submit-btn" type="button" :disabled="!selectedScene || loading || aiSceneLoading" @click="submitScene">
            <span v-if="loading" class="mini-spinner" aria-hidden="true"></span>
            <span>{{ loading ? '正在搭配...' : '从菜谱库推荐' }}</span>
          </button>
          <button class="submit-btn ai-submit-btn" type="button" :disabled="!selectedScene || loading || aiSceneLoading" @click="submitAIScene">
            <span v-if="aiSceneLoading" class="mini-spinner" aria-hidden="true"></span>
            <span>{{ aiSceneLoading ? aiProgressLabel : 'AI 按偏好生成新菜' }}</span>
          </button>
        </div>
        <div v-if="aiSceneLoading || aiProgressStep > 0" class="ai-progress" aria-label="AI 生成进度">
          <span
            v-for="(step, index) in aiProgressSteps"
            :key="step"
            :class="{ active: index <= aiProgressStep, current: index === aiProgressStep && aiSceneLoading }"
          >
            {{ step }}
          </span>
        </div>
      </section>

      <section v-else class="panel shortcut-panel">
        <div class="panel-title">
          <h2>{{ modeMeta.title }}</h2>
        </div>
        <button class="submit-btn" type="button" @click="runShortcutMode">{{ mode === 'week' ? '去安排本周菜单' : '去情侣点餐' }}</button>
      </section>

      <p v-if="message" class="message" role="status">{{ message }}</p>
      <p v-if="error" class="error" role="alert">{{ error }}</p>

      <section v-if="isIngredientMode && results.length" class="result-list" aria-label="推荐结果">
        <article v-for="item in results" :key="item.recipe.id || item.recipe.title" class="result-card">
          <div class="result-media" @click="openRecipe(item.recipe.id)">
            <img v-if="item.recipe.cover" :src="item.recipe.cover" :alt="recipeTitle(item.recipe)" loading="lazy" />
            <div v-else class="no-cover">暂无图片</div>
          </div>
          <div class="result-body">
            <div class="result-head">
              <h2>{{ recipeTitle(item.recipe) }}</h2>
              <span>{{ Math.round(item.match_rate * 100) }}%</span>
            </div>
            <div class="match-row">
              <strong>已匹配</strong>
              <span>{{ item.matched_ingredients.join('、') || '暂无' }}</span>
            </div>
            <div class="match-row missing">
              <strong>{{ item.missing_ingredients.length ? '还缺' : '可直接做' }}</strong>
              <span>{{ item.missing_ingredients.join('、') || '食材已齐' }}</span>
            </div>
            <div class="result-actions">
              <button type="button" @click="openRecipe(item.recipe.id)">查看做法</button>
              <button type="button" @click="addRecipeToShopping(item.recipe.id, item.recipe.title || '')">加入清单</button>
              <button type="button" :disabled="coupleLoading" @click="sendDishToCouple(item.recipe.id, item.recipe.title || '')">和 TA 吃</button>
              <button
                v-if="mode === 'fridge' && item.missing_ingredients.length"
                type="button"
                @click="addMissingToShopping(item.missing_ingredients)"
              >
                缺少食材加入清单
              </button>
            </div>
          </div>
        </article>
      </section>

      <section v-if="sceneResult" class="scene-result">
        <div class="scene-summary">
          <span>{{ sceneSourceLabel }} · {{ sceneResult.menu_name || selectedScene?.title }}</span>
          <h2>{{ sceneResult.reason || '为你搭配好了这一餐' }}</h2>
          <div class="scene-result-actions">
            <button type="button" :disabled="shoppingLoading" @click="addSceneShoppingList">
              {{ shoppingLoading ? '加入中...' : '整组加入清单' }}
            </button>
            <button type="button" :disabled="menuSaving" @click="saveSceneMenu">
              {{ menuSaving ? '保存中...' : '保存菜单' }}
            </button>
            <button type="button" @click="mode === 'new' ? submitAIScene() : submitScene()">换一组</button>
          </div>
        </div>
        <article v-for="dish in sceneDishes" :key="dish.recipe_id || dish.name" class="dish-card" @click="openRecipe(dish.recipe_id)">
          <div>
            <strong>{{ dish.name || '未命名菜谱' }}</strong>
            <p>{{ dish.type || '推荐菜' }} · {{ dish.cook_time ? dish.cook_time + '分钟' : '时间待补充' }} · {{ dish.difficulty || '难度待补充' }}</p>
            <small v-if="dish.reason || dish.steps_summary">推荐理由：{{ dish.reason || dish.steps_summary }}</small>
          </div>
          <button class="couple-dish-btn" type="button" :disabled="coupleLoading" @click.stop="sendDishToCouple(dish.recipe_id, dish.name)">和 TA 吃</button>
          <svg viewBox="0 0 24 24"><path d="m9 18 6-6-6-6" /></svg>
        </article>
      </section>

      <section v-if="hasSubmitted && !loading && isIngredientMode && !results.length && !error" class="empty-panel">
        <h2>暂时没有匹配菜谱</h2>
        <p>可以减少食材限制，或先在后台补充更多菜谱食材数据。</p>
      </section>
    </main>

    <Teleport to="body">
      <Transition name="leave-dialog">
        <div
          v-if="leaveDialogOpen"
          class="leave-dialog-backdrop"
          role="presentation"
          @click.self="resolveLeaveConfirmation(false)"
        >
          <section
            class="leave-dialog-card"
            role="dialog"
            aria-modal="true"
            aria-labelledby="leave-dialog-title"
            aria-describedby="leave-dialog-description"
            @keydown.esc="resolveLeaveConfirmation(false)"
          >
            <span class="leave-dialog-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24"><path d="M12 3 3 20h18L12 3Z" /><path d="M12 9v5" /><path d="M12 17.5h.01" /></svg>
            </span>
            <div>
              <p>推荐还没保存</p>
              <h2 id="leave-dialog-title">要退出这次推荐吗？</h2>
              <span id="leave-dialog-description">退出后，这组 AI 推荐和当前选择将消失，无法恢复。</span>
            </div>
            <div class="leave-dialog-actions">
              <button ref="leaveStayButton" type="button" @click="resolveLeaveConfirmation(false)">继续查看</button>
              <button type="button" @click="resolveLeaveConfirmation(true)">退出推荐</button>
            </div>
          </section>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router'
import { createCoupleOrder } from '@/api/couple'
import { formatLocalDate } from '@/utils/date'
import { saveUserMenu } from '@/api/menu'
import { getIngredientOptions, recommendByIngredients, recommendByScene, recommendSceneByAI } from '@/api/recommend'
import { getRecipeFilterOptions } from '@/api/recipe'
import { getPreferences } from '@/api/user'
import { updateShoppingList } from '@/api/shopping'
import { useShoppingStore } from '@/stores/shopping'
import type { IngredientOption, RecommendMode, RecommendRecipeResult, SceneOption } from '@/types/recommend'
import type { ShoppingItem } from '@/api/shopping'
import { trackEvent } from '@/utils/event'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'

const route = useRoute()
const router = useRouter()
const shoppingStore = useShoppingStore()

const mode = computed<RecommendMode>(() => {
  const value = String(route.params.mode || 'scene')
  return ['ingredients', 'taste', 'scene', 'fridge', 'new', 'week', 'couple'].includes(value) ? value as RecommendMode : 'scene'
})
const isIngredientMode = computed(() => mode.value === 'ingredients' || mode.value === 'fridge')

const ingredientKeyword = ref('')
const ingredientOptions = ref<IngredientOption[]>([])
const ingredientOptionsLoading = ref(false)
const selectedIngredients = ref<string[]>([])
const results = ref<RecommendRecipeResult[]>([])
const tasteOptions = ref<string[]>([])
const selectedScene = ref<SceneOption | null>(null)
const sceneResult = ref<any>(null)
const loading = ref(false)
const aiSceneLoading = ref(false)
const shoppingLoading = ref(false)
const coupleLoading = ref(false)
const menuSaving = ref(false)
const aiProgressStep = ref(0)
let aiProgressTimer: ReturnType<typeof window.setInterval> | null = null
const error = ref('')
const message = ref('')
const hasSubmitted = ref(false)
const sceneMenuSaved = ref(false)
const draftRestored = ref(false)
const leaveDialogOpen = ref(false)
const leaveStayButton = ref<HTMLButtonElement | null>(null)
let leaveConfirmationPromise: Promise<boolean> | null = null
let resolveLeaveDialog: ((confirmed: boolean) => void) | null = null
const userPreference = ref({
  taste_preference: [] as string[],
  avoid_ingredients: [] as string[],
  health_goal: '',
  cook_time_preference: '',
  people_count: 2,
})
const temporaryPreference = ref({
  people_count: 2,
  health_goal: '',
  cook_time_preference: '',
})

const pageVars = computed(() => ({ '--recommend-bg': `url(${kitchenBg})` }))
const peopleOptions = [1, 2, 3, 4, 5, 6, 8, 10]
const cookTimeOptions = ['15分钟内', '30分钟内', '45分钟内', '都可以']
const healthGoalOptions = ['普通', '减脂', '增肌', '控糖', '儿童营养']
const recommendDraftPrefix = 'recommend-draft:v1:'

interface RecommendDraft {
  version: 1
  mode: RecommendMode
  selectedScene: SceneOption | null
  sceneResult: any
  selectedIngredients: string[]
  results: RecommendRecipeResult[]
  hasSubmitted: boolean
  temporaryPreference: typeof temporaryPreference.value
  message: string
  menuSaved: boolean
  scrollY: number
}

const intentOptions: { mode: RecommendMode; title: string; path?: string }[] = [
  { mode: 'scene', title: '不知道吃什么' },
  { mode: 'ingredients', title: '现有食材' },
  { mode: 'taste', title: '按口味找' },
  { mode: 'week', title: '本周菜单', path: '/week-menu' },
  { mode: 'new', title: '新菜灵感' },
  { mode: 'couple', title: '和 TA 吃', path: '/couple' },
]

const aiProgressSteps = ['读取偏好', '生成菜谱', '校验食材', '保存菜谱库']
const aiProgressLabel = computed(() => aiProgressSteps[Math.min(aiProgressStep.value, aiProgressSteps.length - 1)] + '...')

const scenes: SceneOption[] = [
  { key: 'quick_meal', title: '快手一餐', description: '优先 30 分钟内、步骤简单的菜谱', meal_type: 'dinner', people_count: 2, cook_time_preference: '30分钟内' },
  { key: 'family_dinner', title: '家庭聚餐', description: '适合 4 人以上，菜品搭配更完整', meal_type: 'dinner', people_count: 4 },
  { key: 'fat_loss', title: '减脂轻食', description: '偏低脂、低卡或高蛋白', meal_type: 'lunch', people_count: 1, health_goal: '减脂' },
  { key: 'treat_guest', title: '宴客招待', description: '优先热度高、卖相稳的菜', meal_type: 'dinner', people_count: 4 },
  { key: 'late_night', title: '夜宵', description: '分量适中，烹饪时间短', meal_type: 'dinner', people_count: 1, cook_time_preference: '20分钟内' },
]

const modeMeta = computed(() => {
  const map = {
    ingredients: {
      eyebrow: 'Use what you have',
      title: '用现有食材做饭',
      badge: '食材匹配',
      heading: '把手头食材变成一顿饭',
    },
    taste: {
      eyebrow: 'Taste finder',
      title: '按口味找菜',
      badge: '口味筛选',
      heading: '今天想吃什么味道？',
    },
    scene: {
      eyebrow: 'Meal decision',
      title: '我不知道吃什么',
      badge: '场景搭配',
      heading: '让场景替你做选择',
    },
    fridge: {
      eyebrow: 'Fridge rescue',
      title: '冰箱库存',
      badge: '剩菜拯救',
      heading: '先消耗冰箱里的库存',
    },
    new: {
      eyebrow: 'AI new dish',
      title: '换点没吃过的新菜',
      badge: 'AI 生成',
      heading: '按你的偏好生成新菜',
    },
    week: {
      eyebrow: 'Weekly plan',
      title: '安排本周菜单',
      badge: '周菜单',
      heading: '提前安排一周怎么吃',
    },
    couple: {
      eyebrow: 'Together',
      title: '我和 TA 一起决定',
      badge: '情侣点餐',
      heading: '两个人一起选菜',
    },
  }
  return map[mode.value]
})

const visibleIngredientOptions = computed(() => {
  const keyword = ingredientKeyword.value.trim()
  const source = ingredientOptions.value
  if (!keyword) return source.slice(0, 12)
  return source.filter((item) => item.name.includes(keyword)).slice(0, 12)
})

const sceneDishes = computed(() => {
  return Array.isArray(sceneResult.value?.dishes) ? sceneResult.value.dishes : []
})

const sceneSourceLabel = computed(() => {
  if (sceneResult.value?.source === 'ai' || mode.value === 'new') return 'AI 生成'
  return '菜谱库推荐'
})

const preferenceSummary = computed(() => {
  const parts = [
    `${temporaryPreference.value.people_count || userPreference.value.people_count || 2}人`,
    temporaryPreference.value.cook_time_preference || userPreference.value.cook_time_preference || '时间不限',
    temporaryPreference.value.health_goal || userPreference.value.health_goal || '目标不限',
  ]
  return parts.join(' · ')
})

const hasUnsavedSceneResult = computed(() => !!sceneResult.value && !sceneMenuSaved.value)

function recommendDraftKey(value: RecommendMode = mode.value) {
  return recommendDraftPrefix + value
}

function persistRecommendationDraft() {
  if (!sceneResult.value && !results.value.length) return
  const draft: RecommendDraft = {
    version: 1,
    mode: mode.value,
    selectedScene: selectedScene.value,
    sceneResult: sceneResult.value,
    selectedIngredients: selectedIngredients.value,
    results: results.value,
    hasSubmitted: hasSubmitted.value,
    temporaryPreference: temporaryPreference.value,
    message: message.value,
    menuSaved: sceneMenuSaved.value,
    scrollY: window.scrollY,
  }
  try {
    sessionStorage.setItem(recommendDraftKey(), JSON.stringify(draft))
  } catch {
    // 存储不可用时仍允许查看详情，当前页离开确认继续生效。
  }
}

function clearRecommendationDraft(value: RecommendMode = mode.value) {
  try {
    sessionStorage.removeItem(recommendDraftKey(value))
  } catch {
    // 忽略浏览器禁用存储时的异常。
  }
}

function restoreRecommendationDraft() {
  try {
    const raw = sessionStorage.getItem(recommendDraftKey())
    if (!raw) return false
    const draft = JSON.parse(raw) as Partial<RecommendDraft>
    if (draft.version !== 1 || draft.mode !== mode.value || (!draft.sceneResult && !draft.results?.length)) {
      clearRecommendationDraft()
      return false
    }

    selectedScene.value = draft.selectedScene
      ? scenes.find((scene) => scene.key === draft.selectedScene?.key) || draft.selectedScene
      : null
    sceneResult.value = draft.sceneResult || null
    selectedIngredients.value = Array.isArray(draft.selectedIngredients) ? draft.selectedIngredients : []
    results.value = Array.isArray(draft.results) ? draft.results : []
    hasSubmitted.value = !!draft.hasSubmitted
    temporaryPreference.value = {
      ...temporaryPreference.value,
      ...(draft.temporaryPreference || {}),
    }
    message.value = String(draft.message || '')
    sceneMenuSaved.value = !!draft.menuSaved
    draftRestored.value = true

    nextTick(() => {
      window.requestAnimationFrame(() => {
        window.requestAnimationFrame(() => {
          window.scrollTo({ top: Number(draft.scrollY || 0), behavior: 'auto' })
        })
      })
    })
    return true
  } catch {
    clearRecommendationDraft()
    return false
  }
}

function requestLeaveConfirmation() {
  if (leaveConfirmationPromise) return leaveConfirmationPromise
  leaveDialogOpen.value = true
  leaveConfirmationPromise = new Promise<boolean>((resolve) => {
    resolveLeaveDialog = resolve
  })
  nextTick(() => leaveStayButton.value?.focus())
  return leaveConfirmationPromise
}

function resolveLeaveConfirmation(confirmed: boolean) {
  const resolve = resolveLeaveDialog
  leaveDialogOpen.value = false
  leaveConfirmationPromise = null
  resolveLeaveDialog = null
  resolve?.(confirmed)
}

function isRecipeDetailVisit(to: { name?: unknown; query?: Record<string, unknown> }) {
  return to.name === 'RecipeDetail' && to.query?.from === 'recommend'
}

async function confirmRecommendationDeparture(to: { name?: unknown; query?: Record<string, unknown> }) {
  if (isRecipeDetailVisit(to)) {
    persistRecommendationDraft()
    return true
  }
  if (hasUnsavedSceneResult.value) {
    const confirmed = await requestLeaveConfirmation()
    if (!confirmed) return false
  }
  clearRecommendationDraft()
  return true
}

function handleBeforeUnload(event: BeforeUnloadEvent) {
  if (!hasUnsavedSceneResult.value) return
  clearRecommendationDraft()
  event.preventDefault()
  event.returnValue = ''
}

function normalizeName(name: string) {
  return name.trim().replace(/\s+/g, '')
}

function toggleIngredient(name: string) {
  const value = normalizeName(name)
  if (!value) return
  if (selectedIngredients.value.includes(value)) {
    removeIngredient(value)
    return
  }
  if (selectedIngredients.value.length >= 20) {
    message.value = '最多选择 20 项食材。'
    return
  }
  selectedIngredients.value.push(value)
  persistFridgeIngredients()
}

function addCustomIngredient() {
  const value = normalizeName(ingredientKeyword.value)
  if (!value) return
  toggleIngredient(value)
  ingredientKeyword.value = ''
}

function removeIngredient(name: string) {
  selectedIngredients.value = selectedIngredients.value.filter((item) => item !== name)
  persistFridgeIngredients()
}

function clearIngredients() {
  selectedIngredients.value = []
  persistFridgeIngredients()
}

function persistFridgeIngredients() {
  if (mode.value === 'fridge') {
    localStorage.setItem('last_fridge_ingredients', JSON.stringify(selectedIngredients.value))
  }
}

async function loadIngredientOptions() {
  ingredientOptionsLoading.value = true
  try {
    const res = await getIngredientOptions({ keyword: ingredientKeyword.value || undefined })
    ingredientOptions.value = Array.isArray(res?.list) ? res.list : []
  } catch {
    ingredientOptions.value = []
  } finally {
    ingredientOptionsLoading.value = false
  }
}

async function loadTasteOptions() {
  try {
    const res = await getRecipeFilterOptions()
    tasteOptions.value = Array.isArray(res?.tastes) ? res.tastes.filter(Boolean) : []
  } catch {
    tasteOptions.value = []
  }
}

function normalizeStringArray(value: unknown) {
  if (Array.isArray(value)) {
    return value.map((item) => String(item).trim()).filter(Boolean)
  }
  return []
}

async function loadUserPreference() {
  try {
    const pref: any = await getPreferences()
    userPreference.value = {
      taste_preference: normalizeStringArray(pref?.taste_preference),
      avoid_ingredients: normalizeStringArray(pref?.avoid_ingredients),
      health_goal: String(pref?.health_goal || '').trim(),
      cook_time_preference: String(pref?.cook_time_preference || '').trim(),
      people_count: Number(pref?.people_count || pref?.default_servings || 2) || 2,
    }
    if (!draftRestored.value) {
      temporaryPreference.value = {
        people_count: userPreference.value.people_count || 2,
        health_goal: userPreference.value.health_goal,
        cook_time_preference: userPreference.value.cook_time_preference,
      }
    }
  } catch {
    // 偏好读取失败时使用默认参数，推荐流程仍可继续。
  }
}

function scenePayload(scene: SceneOption) {
  return {
    scene: scene.key,
    meal_type: scene.meal_type,
    people_count: temporaryPreference.value.people_count || scene.people_count,
    cook_time_preference: temporaryPreference.value.cook_time_preference || scene.cook_time_preference,
    health_goal: temporaryPreference.value.health_goal || scene.health_goal,
    taste_preference: userPreference.value.taste_preference,
    avoid_ingredients: userPreference.value.avoid_ingredients,
  }
}

async function submitIngredients() {
  if (!selectedIngredients.value.length || loading.value) return
  loading.value = true
  error.value = ''
  message.value = ''
  hasSubmitted.value = true
  try {
    trackEvent({
      event_name: 'recommend_start',
      entity_type: 'recommend',
      payload: { mode: mode.value, ingredients: selectedIngredients.value },
    })
    const res = await recommendByIngredients({
      ingredients: selectedIngredients.value,
      mode: mode.value === 'fridge' ? 'fridge' : 'ingredients',
      limit: 20,
    })
    results.value = Array.isArray(res?.list) ? res.list : []
  } catch (e: any) {
    results.value = []
    error.value = e?.message || '推荐失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

function chooseTaste(taste: string) {
  router.push({ path: '/recipes', query: { taste, sort: 'latest' } })
}

function chooseScene(scene: SceneOption) {
  selectedScene.value = scene
  sceneResult.value = null
  sceneMenuSaved.value = false
  clearRecommendationDraft()
  message.value = ''
  error.value = ''
}

function goIntent(intent: { mode: RecommendMode; path?: string }) {
  if (intent.path) {
    router.push(intent.path)
    return
  }
  router.push('/recommend/' + intent.mode)
}

function runShortcutMode() {
  if (mode.value === 'week') {
    router.push('/week-menu')
    return
  }
  router.push('/couple')
}

function startAIProgress() {
  stopAIProgress()
  aiProgressStep.value = 0
  aiProgressTimer = window.setInterval(() => {
    if (aiProgressStep.value < aiProgressSteps.length - 2) {
      aiProgressStep.value += 1
    }
  }, 1400)
}

function finishAIProgress(success: boolean) {
  stopAIProgress()
  aiProgressStep.value = success ? aiProgressSteps.length - 1 : 0
}

function stopAIProgress() {
  if (aiProgressTimer) {
    window.clearInterval(aiProgressTimer)
    aiProgressTimer = null
  }
}

async function submitScene() {
  if (!selectedScene.value || loading.value || aiSceneLoading.value) return
  loading.value = true
  sceneMenuSaved.value = false
  clearRecommendationDraft()
  error.value = ''
  message.value = ''
  try {
    trackEvent({
      event_name: 'recommend_start',
      entity_type: 'recommend',
      payload: { mode: 'scene', scene: selectedScene.value.key },
    })
    sceneResult.value = await recommendByScene(scenePayload(selectedScene.value))
    message.value = '已从真实菜谱库中为你搭配场景菜单。'
  } catch (e: any) {
    sceneResult.value = null
    error.value = e?.message || '场景推荐失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

async function submitAIScene() {
  if (!selectedScene.value || loading.value || aiSceneLoading.value) return
  aiSceneLoading.value = true
  sceneMenuSaved.value = false
  clearRecommendationDraft()
  startAIProgress()
  error.value = ''
  message.value = ''
  try {
    trackEvent({
      event_name: 'recommend_start',
      entity_type: 'recommend',
      payload: { mode: 'ai_scene', scene: selectedScene.value.key },
    })
    sceneResult.value = await recommendSceneByAI(scenePayload(selectedScene.value))
    finishAIProgress(true)
    message.value = 'AI 已根据你的偏好生成新菜，并同步到菜谱库。'
  } catch (e: any) {
    sceneResult.value = null
    finishAIProgress(false)
    error.value = e?.message || 'AI 场景推荐失败，请确认 AI 配置后重试。'
  } finally {
    aiSceneLoading.value = false
  }
}

function clearSceneResult() {
  selectedScene.value = null
  sceneResult.value = null
  sceneMenuSaved.value = false
  clearRecommendationDraft()
  error.value = ''
  message.value = ''
}

function recipeTitle(recipe: { title?: string }) {
  return recipe.title || '未命名菜谱'
}

function openRecipe(id?: number) {
  if (id) {
    trackEvent({
      event_name: 'recommend_result_click',
      entity_type: 'recipe',
      entity_id: id,
      payload: { mode: mode.value },
    })
    persistRecommendationDraft()
    router.push({
      name: 'RecipeDetail',
      params: { id },
      query: { from: 'recommend', recommendMode: mode.value },
    })
  }
}

function inferCategory(name: string) {
  if (/鸡|鸭|鱼|虾|肉|牛|猪|蛋|排骨/.test(name)) return '肉蛋水产'
  if (/米|面|粥|饭|馒头|粉/.test(name)) return '主食'
  if (/盐|糖|酱|醋|油|姜|蒜|葱|胡椒|淀粉/.test(name)) return '调味'
  return '蔬果'
}

async function addMissingToShopping(names: string[]) {
  const unique = Array.from(new Set(names.map(normalizeName).filter(Boolean)))
  if (!unique.length) return
  const items: ShoppingItem[] = unique.map((name) => ({
    name,
    amount: '按菜谱适量',
    emoji: '食材',
    category: inferCategory(name),
    price: 0,
    checked: false,
  }))
  if (shoppingStore.currentList) {
    const existingNames = new Set(shoppingStore.currentList.items.map((item) => normalizeName(item.name)))
    const nextItems = items.filter((item) => !existingNames.has(normalizeName(item.name)))
    if (!nextItems.length) {
      message.value = '这些缺少食材已经在当前购物清单里了。'
      return
    }
    shoppingStore.currentList.items.push(...nextItems)
    await updateShoppingList(shoppingStore.currentList.id, {
      name: shoppingStore.currentList.name,
      items: shoppingStore.currentList.items,
    })
    message.value = '缺少食材已合并到当前购物清单。'
    return
  }
  await shoppingStore.createList('缺少食材采购', items)
  message.value = '缺少食材已加入购物清单。'
}

async function addRecipeToShopping(recipeId?: number, name?: string) {
  const value = String(name || '').trim()
  if ((!recipeId && !value) || shoppingLoading.value) return
  shoppingLoading.value = true
  message.value = ''
  error.value = ''
  try {
    if (recipeId) {
      await shoppingStore.generateByRecipe(recipeId, value || '菜谱')
    } else {
      await shoppingStore.generateByDish(value)
    }
    trackEvent({
      event_name: 'add_shopping_list',
      entity_type: recipeId ? 'recipe' : 'dish',
      entity_id: recipeId || 0,
      payload: { source: 'recommend', title: value },
    })
    message.value = '「' + value + '」的食材已合并到购物清单。'
  } catch (err) {
    error.value = err instanceof Error ? err.message : '加入购物清单失败，请稍后重试。'
  } finally {
    shoppingLoading.value = false
  }
}

async function addSceneShoppingList() {
  if (shoppingLoading.value) return
  const recipeIds = sceneDishes.value
    .map((dish: any) => Number(dish?.recipe_id || 0))
    .filter((id: number) => Number.isFinite(id) && id > 0)
  if (recipeIds.length) {
    shoppingLoading.value = true
    message.value = ''
    error.value = ''
    try {
      const title = sceneResult.value?.menu_name || selectedScene.value?.title || '推荐采购清单'
      await shoppingStore.generateByRecipes(recipeIds, title + '采购清单')
      trackEvent({
        event_name: 'add_shopping_list',
        entity_type: 'recommend_menu',
        payload: { source: 'recommend', recipe_ids: recipeIds, title },
      })
      message.value = '推荐菜单食材已合并到购物清单。'
    } catch (err) {
      error.value = err instanceof Error ? err.message : '生成采购清单失败，请稍后重试。'
    } finally {
      shoppingLoading.value = false
    }
    return
  }
  const names = Array.isArray(sceneResult.value?.shopping_list) ? sceneResult.value.shopping_list : []
  if (names.length) {
    await addMissingToShopping(names)
    return
  }
  const firstDish = sceneDishes.value.find((dish: any) => dish?.name)
  if (!firstDish?.name) {
    error.value = '当前推荐结果没有可加入清单的食材。'
    return
  }
  await addRecipeToShopping(undefined, firstDish.name)
}

async function saveSceneMenu() {
  if (!sceneResult.value || menuSaving.value) return
  menuSaving.value = true
  message.value = ''
  error.value = ''
  try {
    const title = sceneResult.value?.menu_name || selectedScene.value?.title || '推荐菜单'
    await saveUserMenu({
      name: title,
      menu_type: sceneResult.value?.source === 'ai' || mode.value === 'new' ? 'ai' : 'daily',
      meal_type: scenePayload(selectedScene.value || scenes[0]).meal_type,
      people_count: temporaryPreference.value.people_count,
      health_goal: temporaryPreference.value.health_goal,
      dishes: sceneDishes.value,
      shopping_list: sceneResult.value?.shopping_list || [],
      reason: sceneResult.value?.reason || '',
    })
    trackEvent({
      event_name: 'save_menu',
      entity_type: 'recommend_menu',
      payload: { source: sceneResult.value?.source || mode.value, title },
    })
    message.value = '菜单已保存到“我的菜单”。'
    sceneMenuSaved.value = true
    persistRecommendationDraft()
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存菜单失败，请稍后重试。'
  } finally {
    menuSaving.value = false
  }
}

async function sendDishToCouple(recipeId?: number, name?: string) {
  const title = String(name || '').trim()
  if (!title || coupleLoading.value) return
  coupleLoading.value = true
  message.value = ''
  error.value = ''
  try {
    await createCoupleOrder({
      dish_name: title,
      recipe_id: recipeId,
      meal_type: 'dinner',
      meal_date: formatLocalDate(),
      note: '从推荐结果加入情侣点餐',
    })
    trackEvent({
      event_name: 'couple_order_create',
      entity_type: recipeId ? 'recipe' : 'dish',
      entity_id: recipeId || 0,
      payload: { source: 'recommend', title, mode: mode.value },
    })
    message.value = '已加入情侣点餐，等待 TA 确认。'
  } catch (err) {
    const text = err instanceof Error ? err.message : '加入情侣点餐失败，请稍后重试。'
    error.value = text
    if (text.includes('绑定')) {
      setTimeout(() => router.push('/couple/bind'), 450)
    }
  } finally {
    coupleLoading.value = false
  }
}

watch(mode, () => {
  stopAIProgress()
  aiProgressStep.value = 0
  draftRestored.value = false
  error.value = ''
  message.value = ''
  results.value = []
  sceneResult.value = null
  sceneMenuSaved.value = false
  hasSubmitted.value = false
  selectedScene.value = null
  if (mode.value === 'fridge') {
    try {
      selectedIngredients.value = JSON.parse(localStorage.getItem('last_fridge_ingredients') || '[]')
    } catch {
      selectedIngredients.value = []
    }
  } else {
    selectedIngredients.value = []
  }
}, { immediate: true })

onBeforeRouteLeave((to) => confirmRecommendationDeparture(to))

onBeforeRouteUpdate((to) => {
  const nextMode = String(to.params.mode || 'scene')
  if (to.name === 'Recommend' && nextMode === mode.value) return true
  return confirmRecommendationDeparture(to)
})

onMounted(() => {
  restoreRecommendationDraft()
  loadIngredientOptions()
  loadTasteOptions()
  loadUserPreference()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onBeforeUnmount(() => {
  stopAIProgress()
  window.removeEventListener('beforeunload', handleBeforeUnload)
})
</script>

<style scoped>
.recommend-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --cream: rgba(255, 250, 240, 0.84);
  --coral: #e95645;
  --sage: #6f9856;
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: clip;
  color: var(--text);
  background: linear-gradient(180deg, rgba(255, 246, 231, 0.42), rgba(249, 228, 199, 0.58)), var(--recommend-bg) center top / cover fixed;
  font-family: var(--font-body);
}

.recommend-warm-overlay {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background: radial-gradient(circle at 20% 8%, rgba(255, 255, 255, 0.66), transparent 33%), linear-gradient(90deg, rgba(255, 246, 232, 0.44), rgba(255, 246, 232, 0.1) 55%, rgba(151, 92, 45, 0.08));
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.recommend-phone {
  position: relative;
  z-index: 1;
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: max(38px, env(safe-area-inset-top)) 21px calc(var(--tab-h, 82px) + 70px);
}

.recommend-topbar {
  display: grid;
  grid-template-columns: 54px 1fr;
  align-items: center;
  gap: 14px;
}

.recommend-topbar h1 {
  margin: 0;
  font-size: 29px;
  font-weight: 950;
  line-height: 1.08;
}

svg {
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
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
}

.nav-btn svg {
  width: 28px;
  height: 28px;
}

.hero-card,
.panel,
.preference-panel,
.result-card,
.scene-result,
.empty-panel {
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius);
  background: var(--card-surface);
  box-shadow: var(--card-shadow);
  backdrop-filter: blur(var(--card-blur));
  -webkit-backdrop-filter: blur(var(--card-blur));
}

.hero-card {
  margin-top: 24px;
  padding: 20px;
  border-radius: var(--card-radius-feature);
  box-shadow: var(--card-shadow-feature);
}

.hero-badge {
  width: max-content;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  min-height: 32px;
  padding: 0 13px;
  border-radius: 999px;
  color: var(--coral);
  background: rgba(255, 255, 255, 0.62);
  font-size: 13px;
  font-weight: 850;
}

.hero-badge svg {
  width: 18px;
  height: 18px;
}

.hero-card h2 {
  margin: 15px 0 0;
  font-size: 25px;
  font-weight: 950;
  line-height: 1.12;
}

.intent-panel {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin-top: 14px;
}

.intent-panel button {
  min-height: 58px;
  padding: 10px;
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius-inner);
  background: var(--card-surface-muted);
  color: #3b2b24;
  text-align: left;
  cursor: pointer;
}

.intent-panel button.active {
  border-color: rgba(233, 86, 69, 0.42);
  background: var(--card-surface-strong);
  box-shadow: inset 0 0 0 2px rgba(233, 86, 69, 0.16);
}

.intent-panel strong {
  display: block;
  font-size: 13px;
  font-weight: 950;
  line-height: 1.2;
}

.panel {
  margin-top: 18px;
  padding: 20px;
}

.preference-panel {
  margin-top: 12px;
  padding: 15px;
}

.panel-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 15px;
}

.panel-title h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 950;
}

.panel-title span {
  display: block;
  margin-bottom: 4px;
  color: #c26d3d;
  font-size: 12px;
  font-weight: 850;
}

.panel-title button {
  border: 0;
  background: transparent;
  color: var(--coral);
  font-size: 14px;
  font-weight: 850;
  cursor: pointer;
}

.preference-controls {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 9px;
}

.preference-controls label {
  min-width: 0;
  display: grid;
  gap: 6px;
}

.preference-controls label > span {
  color: #806f64;
  font-size: 12px;
  font-weight: 820;
}

.preference-controls select {
  width: 100%;
  min-height: 42px;
  min-width: 0;
  padding: 0 9px;
  border: 1px solid rgba(143, 111, 86, 0.14);
  border-radius: var(--card-radius-inner);
  color: #3b2b24;
  background: rgba(255, 255, 255, 0.66);
  font: inherit;
  font-size: 13px;
  font-weight: 800;
  outline: none;
}

.preference-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.preference-tags span {
  min-height: 30px;
  display: inline-flex;
  align-items: center;
  padding: 0 10px;
  border-radius: 999px;
  color: #6d7653;
  background: rgba(226, 235, 203, 0.72);
  font-size: 12px;
  font-weight: 850;
}

.ingredient-input {
  display: grid;
  grid-template-columns: 1fr 76px;
  gap: 10px;
}

.ingredient-input input,
.ingredient-input button,
.submit-btn {
  min-height: 48px;
  border: 0;
  border-radius: var(--card-radius-inner);
  font-size: 16px;
}

.ingredient-input input {
  min-width: 0;
  padding: 0 16px;
  color: #352721;
  background: rgba(255, 255, 255, 0.68);
  outline: none;
}

.ingredient-input button,
.submit-btn,
.result-actions button:last-child {
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.24);
  font-weight: 900;
  cursor: pointer;
}

button:disabled {
  opacity: 0.62;
  cursor: not-allowed;
}

.chip-cloud,
.option-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
}

.selected-chip,
.hint-chip,
.option-chip {
  min-height: 42px;
  border: 1px solid rgba(255, 255, 255, 0.68);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.58);
  color: #3b2b24;
  font-weight: 850;
}

.selected-chip {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 0 13px;
  cursor: pointer;
}

.selected-chip svg {
  width: 15px;
  height: 15px;
}

.hint-chip {
  display: inline-flex;
  align-items: center;
  padding: 0 14px;
  color: #8b7a70;
}

.option-chip {
  display: grid;
  align-content: center;
  gap: 1px;
  padding: 7px 14px;
  cursor: pointer;
}

.option-chip small {
  color: #907c70;
  font-size: 11px;
}

.option-chip.active {
  color: #fff;
  background: linear-gradient(135deg, #86aa6a, var(--sage));
}

.option-chip.active small {
  color: rgba(255, 255, 255, 0.82);
}

.submit-btn {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  margin-top: 18px;
}

.scene-actions {
  display: grid;
  gap: 10px;
  margin-top: 18px;
}

.scene-actions .submit-btn {
  margin-top: 0;
}

.ai-submit-btn {
  color: #3a2a24;
  border: 1px solid rgba(233, 86, 69, 0.22);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 238, 214, 0.9));
  box-shadow: 0 12px 24px rgba(80, 50, 30, 0.12);
}

.ai-submit-btn .mini-spinner {
  border-color: rgba(233, 86, 69, 0.22);
  border-top-color: var(--coral);
}

.ai-progress {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 6px;
  margin-top: 12px;
}

.ai-progress span {
  min-height: 32px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.54);
  color: #8a7162;
  font-size: 11px;
  font-weight: 850;
  text-align: center;
}

.ai-progress span.active {
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
}

.ai-progress span.current {
  box-shadow: 0 0 0 3px rgba(233, 86, 69, 0.16);
}

.shortcut-panel p {
  margin: 0;
  color: var(--sub);
  font-size: 14px;
  line-height: 1.58;
}

.mini-spinner {
  width: 17px;
  height: 17px;
  border: 2px solid rgba(255, 255, 255, 0.38);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.taste-grid,
.scene-list {
  display: grid;
  gap: 12px;
}

.taste-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.taste-grid button,
.scene-list button {
  min-height: 58px;
  padding: 14px;
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius-inner);
  text-align: left;
  color: #33241e;
  background: rgba(255, 255, 255, 0.58);
  cursor: pointer;
}

.taste-grid strong,
.scene-list strong {
  display: block;
  font-size: 16px;
  font-weight: 950;
}

.scene-list button.active {
  border-color: rgba(233, 86, 69, 0.4);
  box-shadow: inset 0 0 0 2px rgba(233, 86, 69, 0.22);
}

.message,
.error,
.data-empty {
  margin: 14px 2px 0;
  padding: 12px 14px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 800;
}

.message {
  color: #416834;
  background: rgba(226, 241, 206, 0.8);
}

.error {
  color: #9b2f24;
  background: rgba(255, 225, 219, 0.86);
}

.data-empty {
  color: #7a6a5f;
  background: rgba(255, 255, 255, 0.58);
  line-height: 1.5;
}

.result-list {
  display: grid;
  gap: 16px;
  margin-top: 18px;
}

.result-card {
  overflow: hidden;
}

.result-media {
  height: 172px;
  background: #f4dfc5;
  cursor: pointer;
}

.result-media img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.no-cover {
  height: 100%;
  display: grid;
  place-items: center;
  color: #9b887a;
  font-weight: 900;
  background: linear-gradient(135deg, #fbecd8, #e9caa7);
}

.result-body {
  padding: 18px;
}

.result-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.result-head h2 {
  margin: 0;
  font-size: 23px;
  font-weight: 950;
  line-height: 1.18;
}

.result-head span {
  min-width: 58px;
  height: 34px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: #fff;
  background: var(--sage);
  font-weight: 950;
}

.result-body p,
.match-row {
  color: var(--sub);
  font-size: 14px;
  line-height: 1.5;
}

.match-row {
  display: grid;
  grid-template-columns: 54px 1fr;
  gap: 10px;
  margin-top: 8px;
}

.match-row strong {
  color: #3d2e27;
}

.match-row.missing strong {
  color: var(--coral);
}

.result-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(104px, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.result-actions button {
  min-height: 44px;
  border: 0;
  border-radius: var(--card-radius-inner);
  color: #3a2a24;
  background: rgba(255, 255, 255, 0.7);
  font-weight: 900;
  cursor: pointer;
}

.scene-result,
.empty-panel {
  margin-top: 18px;
  padding: 18px;
}

.scene-summary span {
  color: var(--coral);
  font-weight: 900;
}

.scene-summary h2 {
  margin: 8px 0 14px;
  font-family: var(--font-story);
  font-size: 21px;
  font-weight: 500;
  line-height: 1.55;
  letter-spacing: 0.015em;
}

.scene-result-actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 8px;
}

.scene-result-actions button {
  min-height: 42px;
  border: 0;
  border-radius: var(--card-radius-inner);
  color: #3a2a24;
  background: rgba(255, 255, 255, 0.72);
  font-weight: 900;
  cursor: pointer;
}

.scene-result-actions button:first-child {
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.18);
}

.dish-card {
  min-height: 68px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 12px 0;
  border-top: 1px solid rgba(138, 109, 87, 0.14);
  cursor: pointer;
}

.couple-dish-btn {
  min-width: 72px;
  min-height: 36px;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
  font-size: 12px;
  font-weight: 900;
  cursor: pointer;
}

.couple-dish-btn:disabled {
  cursor: wait;
  opacity: 0.62;
}

.dish-card strong {
  font-size: 17px;
  font-weight: 950;
}

.dish-card p {
  margin: 5px 0 0;
  color: var(--sub);
  font-size: 13px;
}

.dish-card small {
  display: none;
  margin-top: 5px;
  color: #8a7162;
  font-size: 12px;
  line-height: 1.4;
}

.dish-card svg {
  width: 22px;
  height: 22px;
  flex: 0 0 22px;
}

.empty-panel {
  text-align: center;
}

.empty-panel h2 {
  margin: 0 0 8px;
}

.empty-panel p {
  margin: 0;
  color: var(--sub);
  line-height: 1.5;
}

.leave-dialog-backdrop {
  position: fixed;
  inset: 0;
  z-index: 300;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(47, 35, 28, 0.42);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.leave-dialog-card {
  width: min(100%, 354px);
  padding: 24px;
  border: 1px solid rgba(255, 255, 255, 0.88);
  border-radius: var(--card-radius-feature);
  color: #342720;
  background: rgba(255, 250, 241, 0.97);
  box-shadow: var(--card-shadow-dialog);
  font-family: var(--font-body);
  text-align: center;
}

.leave-dialog-icon {
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  margin: 0 auto 14px;
  border-radius: 18px;
  color: #e95645;
  background: #ffe8df;
}

.leave-dialog-icon svg {
  width: 26px;
  height: 26px;
}

.leave-dialog-card p {
  margin: 0 0 5px;
  color: #d95a48;
  font-size: 13px;
  font-weight: 750;
}

.leave-dialog-card h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 850;
  line-height: 1.3;
}

.leave-dialog-card div > span {
  display: block;
  margin-top: 10px;
  color: #806f64;
  font-size: 14px;
  line-height: 1.6;
}

.leave-dialog-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-top: 22px;
}

.leave-dialog-actions button {
  min-height: 46px;
  border: 0;
  border-radius: var(--card-radius-inner);
  color: #4a3931;
  background: #f2e8dd;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
}

.leave-dialog-actions button:last-child {
  color: #fff;
  background: linear-gradient(135deg, #ff7469, #e95645);
  box-shadow: 0 10px 24px rgba(233, 86, 69, 0.22);
}

.leave-dialog-enter-active,
.leave-dialog-leave-active {
  transition: opacity 180ms ease;
}

.leave-dialog-enter-active .leave-dialog-card,
.leave-dialog-leave-active .leave-dialog-card {
  transition: transform 220ms cubic-bezier(0.16, 1, 0.3, 1), opacity 180ms ease;
}

.leave-dialog-enter-from,
.leave-dialog-leave-to {
  opacity: 0;
}

.leave-dialog-enter-from .leave-dialog-card,
.leave-dialog-leave-to .leave-dialog-card {
  opacity: 0;
  transform: translateY(10px) scale(0.97);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (hover: hover) {
  button:hover,
  .dish-card:hover {
    transform: translateY(-1px);
  }
}

@media (max-width: 380px) {
  .recommend-phone {
    padding-left: 18px;
    padding-right: 18px;
  }

  .taste-grid,
  .preference-controls,
  .result-actions,
  .scene-result-actions {
    grid-template-columns: 1fr;
  }

  .intent-panel {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .ai-progress {
    grid-template-columns: repeat(2, minmax(0, 1fr));
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
