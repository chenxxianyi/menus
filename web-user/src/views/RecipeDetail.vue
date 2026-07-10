<template>
  <main class="recipe-detail-shell" :style="pageVars">
    <div class="warm-overlay" aria-hidden="true"></div>

    <section class="recipe-phone">
      <header class="top-bar" aria-label="页面顶部">
        <button class="nav-btn" type="button" aria-label="返回" @click="router.back()">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="m15 18-6-6 6-6" /></svg>
        </button>
        <div class="top-title">
          <span>Recipe detail</span>
          <strong>{{ recipe?.title || '菜谱详情' }}</strong>
        </div>
        <button
          class="nav-btn"
          type="button"
          :aria-label="recipe?.is_favorited ? '取消收藏菜谱' : '收藏菜谱'"
          :disabled="!recipe || favoriteSaving"
          @click="handleFavorite"
        >
          <svg viewBox="0 0 24 24" aria-hidden="true" :class="{ filled: recipe?.is_favorited }">
            <path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" />
          </svg>
        </button>
      </header>

      <section v-if="loading" class="state-card" aria-label="加载中">
        <span class="spinner" aria-hidden="true"></span>
        <p>正在整理菜谱...</p>
      </section>

      <section v-else-if="errorText" class="state-card error" role="alert">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M12 9v4" />
          <path d="M12 17h.01" />
          <path d="M10.3 4.3 2.8 17.2A2 2 0 0 0 4.5 20h15a2 2 0 0 0 1.7-2.8L13.7 4.3a2 2 0 0 0-3.4 0Z" />
        </svg>
        <h1>菜谱加载失败</h1>
        <p>{{ errorText }}</p>
      </section>

      <template v-else-if="recipe">
        <section class="hero-card" aria-label="菜谱概览">
          <div class="hero-media">
            <img v-if="recipe.cover" :src="recipe.cover" :alt="recipe.title" loading="lazy" />
            <div v-else class="hero-placeholder" aria-hidden="true">
              <svg viewBox="0 0 120 120">
                <path d="M22 82h76" />
                <path d="M31 82a29 29 0 0 1 58 0" />
                <path d="M60 25v16" />
                <path d="M36 98h48" />
                <path d="M44 52c7-7 25-7 32 0" />
              </svg>
              <span>家庭菜谱</span>
            </div>
            <div class="hero-wash"></div>
            <div class="hero-badges">
              <span>{{ recipe.difficulty || '家常' }}</span>
              <span v-if="recipe.taste">{{ recipe.taste }}</span>
            </div>
          </div>

          <div class="hero-copy">
            <p class="eyebrow">今日想吃</p>
            <h1>{{ recipe.title }}</h1>
            <p>{{ recipe.description || '一份适合家庭餐桌的菜谱，照着步骤慢慢来就很好吃。' }}</p>
          </div>
        </section>

        <section class="stats-card" aria-label="菜谱统计">
          <div>
            <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="8.5" /><path d="M12 7v5l3 2" /></svg>
            <span>时间</span>
            <strong>{{ recipe.cook_time || '--' }} 分钟</strong>
          </div>
          <div>
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-8 0v2" /><circle cx="12" cy="7" r="4" /></svg>
            <span>份量</span>
            <strong>{{ recipe.people_count || 2 }} 人份</strong>
          </div>
          <div>
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 21c4.4 0 8-3.6 8-8 0-5.2-4.2-9.5-8-11-3.8 1.5-8 5.8-8 11 0 4.4 3.6 8 8 8Z" /><path d="M8 15c2.6 0 5.6-2.1 7.5-5.2" /></svg>
            <span>风味</span>
            <strong>{{ recipe.taste || '家常' }}</strong>
          </div>
          <div>
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M13 2 5 14h6l-1 8 8-12h-6l1-8Z" /></svg>
            <span>热量</span>
            <strong>{{ nutrition.calories || '--' }}</strong>
          </div>
        </section>

        <section v-if="healthTags.length" class="tag-card" aria-label="健康标签">
          <span v-for="tag in healthTags" :key="tag">{{ tag }}</span>
        </section>

        <section v-if="hasNutrition" class="glass-card" aria-label="营养成分">
          <div class="section-title">
            <span><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 19c4-7 9-10 16-11" /><path d="M5 13c4-1 7-4 9-9 3 4 4 8 2 12-2.4 4.8-7.6 5.5-11 3Z" /></svg></span>
            <div>
              <h2>营养成分</h2>
              <p>菜谱维护或 AI 生成的估算值</p>
            </div>
          </div>
          <div class="nutri-list">
            <div v-for="item in nutritionRows" :key="item.key" class="nutri-row">
              <span>{{ item.label }}</span>
              <div class="nutri-track"><i :style="{ width: item.percent + '%' }"></i></div>
              <strong>{{ item.value }}{{ item.unit }}</strong>
            </div>
          </div>
        </section>

        <section v-if="ingredientGroups.length" class="glass-card" aria-label="食材清单">
          <div class="section-title">
            <span><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 8h12l-1.1 11.2A2 2 0 0 1 14.9 21H9.1a2 2 0 0 1-2-1.8L6 8Z" /><path d="M9 8V6a3 3 0 0 1 6 0v2" /></svg></span>
            <div>
              <h2>食材清单</h2>
              <p>轻触可标记家里已有</p>
            </div>
          </div>
          <div class="ingredient-groups">
            <section v-for="group in ingredientGroups" :key="group.title" class="ingredient-group">
              <h3>{{ group.title }} <span>{{ group.items.length }}</span></h3>
              <button
                v-for="item in group.items"
                :key="item.id"
                class="ingredient-row"
                :class="{ checked: item.checked }"
                type="button"
                @click="toggleIngredientOwned(item)"
              >
                <span class="check-dot"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="m5 12 4 4L19 6" /></svg></span>
                <span class="ingredient-name">{{ item.name }}</span>
                <span class="ingredient-amount">{{ item.amountText }}</span>
                <span v-if="item.checked" class="owned-pill">已有</span>
              </button>
            </section>
          </div>
        </section>

        <section v-if="steps.length" class="glass-card" aria-label="烹饪步骤">
          <div class="section-title">
            <span><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M8 6h13" /><path d="M8 12h13" /><path d="M8 18h13" /><path d="M3 6h.01" /><path d="M3 12h.01" /><path d="M3 18h.01" /></svg></span>
            <div>
              <h2>烹饪步骤</h2>
              <p>按顺序来，别慌，锅会等你的</p>
            </div>
          </div>
          <ol class="steps-list">
            <li v-for="(step, index) in steps" :key="index">
              <span class="step-num">{{ index + 1 }}</span>
              <div class="step-body">
                <p>{{ step.description }}</p>
                <em v-if="step.tip">{{ step.tip }}</em>
              </div>
            </li>
          </ol>
        </section>

        <section v-if="recipe.tips" class="tips-card" aria-label="小贴士">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2v4" /><path d="M12 18v4" /><path d="M4.93 4.93 7.76 7.76" /><path d="M16.24 16.24 19.07 19.07" /><path d="M2 12h4" /><path d="M18 12h4" /><path d="M4.93 19.07 7.76 16.24" /><path d="M16.24 7.76 19.07 4.93" /></svg>
          <div>
            <strong>做饭小贴士</strong>
            <p>{{ recipe.tips }}</p>
          </div>
        </section>

        <section class="feedback-card" aria-label="菜谱反馈">
          <div class="section-title">
            <span><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 20h9" /><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4 12.5-12.5Z" /></svg></span>
            <div>
              <h2>这道菜合适吗</h2>
              <p>你的反馈会影响以后推荐</p>
            </div>
          </div>
          <div class="feedback-actions">
            <button
              v-for="item in feedbackActions"
              :key="item.type"
              type="button"
              :class="{ active: feedbackStatus[item.type] }"
              :disabled="feedbackSaving === item.type"
              @click="toggleRecipeFeedback(item.type)"
            >
              <svg v-if="item.type === 'cooked'" viewBox="0 0 24 24" aria-hidden="true"><path d="M4 17h16" /><path d="M6 17a6 6 0 0 1 12 0" /><path d="M12 6v2" /><path d="M5 20h14" /></svg>
              <svg v-else-if="item.type === 'like'" viewBox="0 0 24 24" aria-hidden="true"><path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" /></svg>
              <svg v-else-if="item.type === 'dislike'" viewBox="0 0 24 24" aria-hidden="true"><path d="M10 15 5 20" /><path d="m5 15 5 5" /><path d="M14 4h4a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-3l-4 5v-5H7a2 2 0 0 1-2-2V8" /></svg>
              <svg v-else viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="m5.7 5.7 12.6 12.6" /></svg>
              <span>{{ item.label }}</span>
            </button>
          </div>
        </section>

        <p v-if="actionMessage" class="action-message" role="status">{{ actionMessage }}</p>
        <p v-if="actionError" class="action-error" role="alert">{{ actionError }}</p>

        <div class="bottom-actions" aria-label="菜谱操作">
          <button class="primary-action" type="button" :disabled="shoppingSaving" @click="addRecipeToShopping">
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 6h15l-1.5 9h-12L6 6Z" /><path d="M6 6 5.2 3H3" /><circle cx="9" cy="20" r="1.5" /><circle cx="18" cy="20" r="1.5" /></svg>
            <span>{{ shoppingSaving ? '加入中...' : '加入清单' }}</span>
          </button>
          <button class="secondary-action" type="button" aria-label="发给 TA 想吃" :disabled="coupleSaving" @click="sendToCouple">
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M22 2 11 13" /><path d="m22 2-7 20-4-9-9-4 20-7Z" /></svg>
          </button>
          <button class="secondary-action" type="button" :class="{ active: recipe.is_favorited }" :disabled="favoriteSaving" @click="handleFavorite">
            <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" /></svg>
          </button>
        </div>
      </template>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createCoupleOrder } from '@/api/couple'
import { formatLocalDate } from '@/utils/date'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { deleteRecipeFeedback, getRecipeDetail, removeFavorite, setRecipeFeedback, toggleFavorite } from '@/api/recipe'
import { useShoppingStore } from '@/stores/shopping'
import type { ShoppingItem } from '@/api/shopping'
import type { RecipeFeedbackStatus } from '@/types/recipe'
import { trackEvent } from '@/utils/event'

interface DisplayIngredient {
  id: string
  name: string
  amountText: string
  category: string
  checked: boolean
}

const route = useRoute()
const router = useRouter()
const shoppingStore = useShoppingStore()
const recipe = ref<any>(null)
const loading = ref(true)
const errorText = ref('')
const favoriteSaving = ref(false)
const shoppingSaving = ref(false)
const coupleSaving = ref(false)
const feedbackSaving = ref<keyof RecipeFeedbackStatus | ''>('')
const actionMessage = ref('')
const actionError = ref('')
const preparedItems = ref<DisplayIngredient[]>([])

const defaultFeedbackStatus: RecipeFeedbackStatus = {
  cooked: false,
  like: false,
  dislike: false,
  block: false,
  normal: false,
  too_complex: false,
  too_long: false,
  hard_to_buy: false,
}
const feedbackStatus = ref<RecipeFeedbackStatus>({ ...defaultFeedbackStatus })
const feedbackActions: { type: keyof RecipeFeedbackStatus; label: string }[] = [
  { type: 'cooked', label: '做过了' },
  { type: 'like', label: '喜欢' },
  { type: 'dislike', label: '不喜欢' },
  { type: 'normal', label: '一般' },
  { type: 'too_complex', label: '太复杂' },
  { type: 'too_long', label: '用时过长' },
  { type: 'hard_to_buy', label: '食材难买' },
  { type: 'block', label: '不再推荐' },
]

const pageVars = computed(() => ({
  '--recipe-bg': 'url(' + kitchenBg + ')',
}))

const nutrition = computed(() => recipe.value?.nutrition || {})
const hasNutrition = computed(() => ['calories', 'protein', 'fat', 'carbs', 'fiber'].some((key) => Number(nutrition.value?.[key]) > 0))
const nutritionRows = computed(() => [
  { key: 'protein', label: '蛋白质', value: numberValue(nutrition.value.protein), unit: 'g', percent: percent(nutrition.value.protein, 60) },
  { key: 'fat', label: '脂肪', value: numberValue(nutrition.value.fat), unit: 'g', percent: percent(nutrition.value.fat, 60) },
  { key: 'carbs', label: '碳水', value: numberValue(nutrition.value.carbs), unit: 'g', percent: percent(nutrition.value.carbs, 90) },
  { key: 'fiber', label: '纤维', value: numberValue(nutrition.value.fiber), unit: 'g', percent: percent(nutrition.value.fiber, 30) },
])
const healthTags = computed(() => normalizeArray(recipe.value?.health_tags).slice(0, 6))
const steps = computed(() => {
  const raw = recipe.value?.steps
  if (!Array.isArray(raw)) return []
  return raw
    .map((item: any) => ({ description: String(item.description || item.desc || '').trim(), tip: String(item.tip || '').trim() }))
    .filter((item) => item.description)
})
const ingredientGroups = computed(() => {
  const groups: Record<string, DisplayIngredient[]> = {}
  preparedItems.value.forEach((item) => {
    const title = groupTitle(item.category, item.name)
    if (!groups[title]) groups[title] = []
    groups[title].push(item)
  })
  return Object.entries(groups).map(([title, items]) => ({ title, items }))
})

function normalizeArray(value: unknown): string[] {
  if (Array.isArray(value)) return value.map((item) => String(item).trim()).filter(Boolean)
  return []
}

function normalizeIngredients(raw: unknown): DisplayIngredient[] {
  if (!Array.isArray(raw)) return []
  return raw
    .map((item: any, index) => {
      const name = String(item.name || '').trim()
      if (!name) return null
      return {
        id: String(index) + '-' + name,
        name,
        amountText: formatAmount(item.amount, item.unit),
        category: String(item.category || '').trim(),
        checked: false,
      }
    })
    .filter((item): item is DisplayIngredient => !!item)
}

function formatAmount(amount: unknown, unit: unknown) {
  const amountText = String(amount ?? '').trim()
  const unitText = String(unit ?? '').trim()
  if (!amountText && !unitText) return '适量'
  if (!unitText || amountText.includes(unitText)) return amountText || unitText
  return amountText + unitText
}

function groupTitle(category: string, name: string) {
  const text = category + name
  if (/调味|调料|盐|糖|酱|醋|油|料酒|胡椒/.test(text)) return '调味'
  if (/肉蛋|水产|肉|鸡|鸭|鱼|虾|牛|猪|羊|蛋|排骨/.test(text)) return '肉蛋水产'
  if (/主食|米|面|粉|饭|馒头/.test(text)) return '主食'
  if (/配料|葱|姜|蒜|香菇|木耳|红枣/.test(text)) return '配料'
  if (/蔬菜|青菜|番茄|土豆|萝卜|玉米|白菜/.test(text)) return '蔬菜'
  return category || '其他'
}

function numberValue(value: unknown) {
  const num = Number(value || 0)
  if (!Number.isFinite(num) || num <= 0) return 0
  return Math.round(num * 10) / 10
}

function percent(value: unknown, max: number) {
  const num = Number(value || 0)
  if (!Number.isFinite(num) || num <= 0) return 0
  return Math.min(100, Math.round((num / max) * 100))
}

async function handleFavorite() {
  if (!recipe.value || favoriteSaving.value) return
  const wasFavorited = recipe.value.is_favorited === true
  favoriteSaving.value = true
  recipe.value.is_favorited = !wasFavorited
  try {
    if (wasFavorited) {
      await removeFavorite(recipe.value.id)
    } else {
      await toggleFavorite(recipe.value.id)
    }
    trackEvent({
      event_name: 'recipe_favorited',
      entity_type: 'recipe',
      entity_id: recipe.value.id,
      payload: { active: !wasFavorited, source: 'detail' },
    })
  } catch {
    recipe.value.is_favorited = wasFavorited
  } finally {
    favoriteSaving.value = false
  }
}

async function addRecipeToShopping() {
  if (!recipe.value?.id || shoppingSaving.value) return
  shoppingSaving.value = true
  actionMessage.value = ''
  actionError.value = ''
  try {
    const ownedItems = preparedItems.value
      .filter((item) => item.checked)
      .map((item): ShoppingItem => ({
        name: item.name,
        amount: item.amountText || '适量',
        emoji: '',
        category: item.category || groupTitle(item.category, item.name),
        price: 0,
        checked: false,
        status: 'owned',
      }))
    await shoppingStore.generateByRecipe(recipe.value.id, recipe.value.title || '菜谱')
    if (ownedItems.length) {
      await shoppingStore.appendItemsToCurrentList((recipe.value.title || '菜谱') + '已有食材', ownedItems)
    }
    trackEvent({
      event_name: 'shopping_list_generated',
      entity_type: 'recipe',
      entity_id: recipe.value.id,
      payload: { source: 'detail', title: recipe.value.title, owned_count: ownedItems.length },
    })
    actionMessage.value = ownedItems.length
      ? `食材已合并到购物清单，${ownedItems.length} 项标记为家里已有。`
      : '食材已合并到购物清单。'
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : '加入购物清单失败，请稍后重试。'
  } finally {
    shoppingSaving.value = false
  }
}

function toggleIngredientOwned(item: DisplayIngredient) {
  item.checked = !item.checked
}

async function sendToCouple() {
  if (!recipe.value?.id || coupleSaving.value) return
  coupleSaving.value = true
  actionMessage.value = ''
  actionError.value = ''
  try {
    await createCoupleOrder({
      dish_name: recipe.value.title || '想吃的菜',
      recipe_id: recipe.value.id,
      meal_type: 'dinner',
      meal_date: formatLocalDate(),
      note: '从菜谱详情发给 TA',
    })
    trackEvent({
      event_name: 'couple_order_create',
      entity_type: 'recipe',
      entity_id: recipe.value.id,
      payload: { source: 'detail', title: recipe.value.title },
    })
    actionMessage.value = '已发给 TA，等待对方确认。'
  } catch (error) {
    const message = error instanceof Error ? error.message : '发送失败，请稍后重试。'
    actionError.value = message
    if (message.includes('绑定')) {
      setTimeout(() => router.push('/couple/bind'), 450)
    }
  } finally {
    coupleSaving.value = false
  }
}

function normalizeFeedbackStatus(value: any): RecipeFeedbackStatus {
  return {
    cooked: value?.cooked === true,
    like: value?.like === true,
    dislike: value?.dislike === true,
    block: value?.block === true,
    normal: value?.normal === true,
    too_complex: value?.too_complex === true,
    too_long: value?.too_long === true,
    hard_to_buy: value?.hard_to_buy === true,
  }
}

async function toggleRecipeFeedback(type: keyof RecipeFeedbackStatus) {
  if (!recipe.value?.id || feedbackSaving.value) return
  if (type === 'block' && !feedbackStatus.value.block) {
    const confirmed = window.confirm('之后将不再为你推荐这道菜，确认继续吗？')
    if (!confirmed) return
  }

  const wasActive = feedbackStatus.value[type]
  feedbackSaving.value = type
  actionMessage.value = ''
  actionError.value = ''
  try {
    const res = wasActive
      ? await deleteRecipeFeedback(recipe.value.id, type)
      : await setRecipeFeedback(recipe.value.id, type, 'detail')
    feedbackStatus.value = normalizeFeedbackStatus(res.feedback)
    trackEvent({
      event_name: 'recipe_feedback_submitted',
      entity_type: 'recipe',
      entity_id: recipe.value.id,
      payload: { type, active: !wasActive, source: 'detail' },
    })
    actionMessage.value = wasActive ? '已取消反馈。' : feedbackCopy(type)
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : '反馈提交失败，请稍后重试。'
  } finally {
    feedbackSaving.value = ''
  }
}

function feedbackCopy(type: keyof RecipeFeedbackStatus) {
  const map: Record<keyof RecipeFeedbackStatus, string> = {
    cooked: '已记录做过了，近期会减少重复推荐。',
    like: '已记录喜欢，后续会多推荐相似口味。',
    dislike: '已记录不喜欢，后续会降低推荐权重。',
    block: '已记录不再推荐，后续会避开这道菜。',
    normal: '已记录“一般”，后续会继续探索更合适的菜谱。',
    too_complex: '已记录“太复杂”，后续会优先推荐更易上手的菜谱。',
    too_long: '已记录“用时过长”，后续会优先推荐更快完成的菜谱。',
    hard_to_buy: '已记录“食材难买”，后续会降低相似菜谱的推荐权重。',
  }
  return map[type]
}

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) {
    errorText.value = '菜谱 ID 无效'
    loading.value = false
    return
  }
  try {
    const res: any = await getRecipeDetail(id)
    recipe.value = res
    feedbackStatus.value = normalizeFeedbackStatus(res?.feedback)
    preparedItems.value = normalizeIngredients(res?.ingredients)
  } catch (error) {
    errorText.value = error instanceof Error ? error.message : '请稍后重试'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.recipe-detail-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.84);
  --coral: #e95645;
  --sage: #79a35d;
  --orange: #f28a2e;
  --line: rgba(143, 111, 86, 0.18);
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 246, 231, 0.42), rgba(249, 228, 199, 0.58)),
    var(--recipe-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.warm-overlay {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 20% 8%, rgba(255, 255, 255, 0.68), transparent 33%),
    radial-gradient(circle at 86% 15%, rgba(242, 138, 46, 0.2), transparent 30%),
    linear-gradient(90deg, rgba(255, 246, 232, 0.5), rgba(255, 246, 232, 0.12) 55%, rgba(151, 92, 45, 0.1));
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.recipe-phone {
  position: relative;
  z-index: 1;
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: max(18px, env(safe-area-inset-top)) 21px calc(var(--tab-h, 64px) + 108px + env(safe-area-inset-bottom));
}

button {
  border: 0;
  font: inherit;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

svg {
  display: block;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.top-bar {
  display: grid;
  grid-template-columns: 52px minmax(0, 1fr) 52px;
  align-items: center;
  gap: 12px;
  margin-bottom: 21px;
}

.nav-btn {
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.68);
  border-radius: 18px;
  color: #5f4c43;
  background: rgba(255, 250, 240, 0.8);
  box-shadow:
    0 12px 24px rgba(80, 50, 30, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.nav-btn:disabled {
  cursor: not-allowed;
  opacity: 0.58;
}

.nav-btn svg {
  width: 24px;
  height: 24px;
  stroke-width: 2.5;
}

.nav-btn svg.filled {
  color: var(--coral);
  fill: rgba(233, 86, 69, 0.18);
}

.top-title {
  min-width: 0;
  display: grid;
  gap: 3px;
  text-align: center;
}

.top-title span {
  color: var(--orange);
  font-size: 12px;
  font-weight: 850;
  text-transform: uppercase;
}

.top-title strong {
  overflow: hidden;
  color: var(--text);
  font-size: 19px;
  font-weight: 920;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hero-card,
.stats-card,
.tag-card,
.glass-card,
.tips-card,
.state-card {
  border: 1px solid var(--card-border);
  background: var(--card-surface);
  box-shadow: var(--card-shadow);
  backdrop-filter: blur(var(--card-blur));
  -webkit-backdrop-filter: blur(var(--card-blur));
}

.hero-card {
  overflow: hidden;
  border-radius: var(--card-radius-feature);
  box-shadow: var(--card-shadow-feature);
}

.hero-media {
  position: relative;
  min-height: 254px;
  overflow: hidden;
  background:
    radial-gradient(circle at 25% 10%, rgba(255, 255, 255, 0.8), transparent 34%),
    linear-gradient(145deg, #fff2df, #ecd2b2);
}

.hero-media img {
  width: 100%;
  height: 284px;
  display: block;
  object-fit: cover;
}

.hero-placeholder {
  min-height: 284px;
  display: grid;
  place-items: center;
  color: #c88768;
}

.hero-placeholder svg {
  width: 142px;
  height: 142px;
  stroke-width: 3;
  filter: drop-shadow(0 18px 24px rgba(142, 86, 45, 0.12));
}

.hero-placeholder span {
  margin-top: -42px;
  color: #9a7957;
  font-size: 14px;
  font-weight: 850;
}

.hero-wash {
  position: absolute;
  inset: auto 0 0;
  height: 45%;
  background: linear-gradient(180deg, transparent, rgba(46, 36, 31, 0.28));
  pointer-events: none;
}

.hero-badges {
  position: absolute;
  left: 18px;
  right: 18px;
  bottom: 18px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.hero-badges span,
.tag-card span {
  min-height: 32px;
  display: inline-flex;
  align-items: center;
  padding: 0 12px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 999px;
  color: #5d4b40;
  background: rgba(255, 250, 240, 0.76);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.82);
  font-size: 12px;
  font-weight: 820;
}

.hero-copy {
  padding: 22px 23px 25px;
}

.eyebrow {
  margin: 0 0 8px;
  color: var(--orange);
  font-size: 13px;
  font-weight: 900;
}

.hero-copy h1 {
  margin: 0;
  color: #30221d;
  font-size: clamp(34px, 9vw, 44px);
  font-weight: 950;
  line-height: 1.02;
  letter-spacing: 0;
}

.hero-copy p {
  margin: 12px 0 0;
  color: var(--sub);
  font-size: 15px;
  font-weight: 650;
  line-height: 1.62;
}

.stats-card {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  overflow: hidden;
  margin-top: 17px;
  border-radius: var(--card-radius);
}

.stats-card > div {
  min-height: 104px;
  display: grid;
  align-content: center;
  gap: 6px;
  padding: 18px;
  border-bottom: 1px solid var(--line);
}

.stats-card > div:nth-child(odd) {
  border-right: 1px solid var(--line);
}

.stats-card > div:nth-child(n + 3) {
  border-bottom: 0;
}

.stats-card svg {
  width: 22px;
  height: 22px;
  color: var(--coral);
  stroke-width: 2.2;
}

.stats-card span {
  color: var(--sub);
  font-size: 12px;
  font-weight: 760;
}

.stats-card strong {
  color: var(--text);
  font-size: 22px;
  font-weight: 920;
  line-height: 1;
}

.tag-card {
  display: flex;
  flex-wrap: wrap;
  gap: 9px;
  margin-top: 17px;
  padding: 15px;
  border-radius: var(--card-radius);
}

.tag-card span {
  color: #6f7f47;
  background: rgba(238, 246, 231, 0.76);
}

.glass-card,
.feedback-card,
.tips-card,
.state-card {
  margin-top: 17px;
  padding: 22px;
  border-radius: var(--card-radius);
}

.feedback-card {
  border: 1px solid rgba(255, 255, 255, 0.68);
  background: var(--cream);
  box-shadow:
    0 18px 42px rgba(80, 50, 28, 0.14),
    inset 0 1px 0 rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 13px;
  margin-bottom: 19px;
}

.section-title > span {
  width: 45px;
  height: 45px;
  display: grid;
  flex: 0 0 45px;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 16px;
  color: var(--coral);
  background: rgba(252, 226, 214, 0.62);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.82);
}

.section-title svg {
  width: 24px;
  height: 24px;
  stroke-width: 2.2;
}

.section-title h2 {
  margin: 0;
  color: var(--text);
  font-size: 22px;
  font-weight: 930;
  line-height: 1;
}

.section-title p {
  margin: 6px 0 0;
  color: var(--sub);
  font-size: 12px;
  font-weight: 650;
  line-height: 1.35;
}

.nutri-list {
  display: grid;
  gap: 15px;
}

.nutri-row {
  display: grid;
  grid-template-columns: 58px minmax(0, 1fr) 48px;
  align-items: center;
  gap: 12px;
}

.nutri-row span {
  color: var(--sub);
  font-size: 13px;
  font-weight: 760;
}

.nutri-track {
  height: 8px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(143, 111, 86, 0.12);
}

.nutri-track i {
  height: 100%;
  display: block;
  border-radius: inherit;
  background: linear-gradient(90deg, var(--sage), var(--coral));
}

.nutri-row strong {
  color: var(--text);
  font-size: 13px;
  font-weight: 850;
  text-align: right;
}

.ingredient-groups {
  display: grid;
  gap: 18px;
}

.ingredient-group + .ingredient-group {
  padding-top: 18px;
  border-top: 1.5px dashed var(--line);
}

.ingredient-group h3 {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 0 0 10px;
  color: #6f8d55;
  font-size: 17px;
  font-weight: 920;
}

.ingredient-group h3 span {
  color: var(--sub);
  font-size: 13px;
  font-weight: 820;
}

.ingredient-row {
  width: 100%;
  min-height: 54px;
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
  color: var(--text);
  background: transparent;
  text-align: left;
}

.ingredient-row + .ingredient-row {
  border-top: 1px solid rgba(143, 111, 86, 0.1);
}

.check-dot {
  width: 26px;
  height: 26px;
  display: grid;
  place-items: center;
  border: 2px solid rgba(143, 111, 86, 0.36);
  border-radius: 50%;
  color: #fff;
  background: rgba(255, 255, 255, 0.56);
  transition: border-color 180ms ease, background 180ms ease, transform 180ms ease;
}

.check-dot svg {
  width: 15px;
  height: 15px;
  stroke-width: 3;
  opacity: 0;
  transform: scale(0.7);
  transition: opacity 160ms ease, transform 160ms ease;
}

.ingredient-row.checked .check-dot {
  border-color: var(--coral);
  background: var(--coral);
}

.ingredient-row.checked .check-dot svg {
  opacity: 1;
  transform: scale(1);
}

.ingredient-name {
  overflow: hidden;
  font-size: 16px;
  font-weight: 850;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ingredient-row.checked .ingredient-name {
  color: var(--muted);
  text-decoration: line-through;
}

.ingredient-amount {
  color: var(--sub);
  font-size: 13px;
  font-weight: 760;
  white-space: nowrap;
}

.owned-pill {
  height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 9px;
  border-radius: 999px;
  color: #5f7e4b;
  background: rgba(126, 163, 106, 0.16);
  font-size: 12px;
  font-weight: 880;
  white-space: nowrap;
}

.steps-list {
  display: grid;
  gap: 18px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.steps-list li {
  display: grid;
  grid-template-columns: 34px minmax(0, 1fr);
  gap: 14px;
  padding-top: 18px;
  border-top: 1.5px dashed var(--line);
}

.steps-list li:first-child {
  padding-top: 0;
  border-top: 0;
}

.step-num {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  color: var(--coral);
  background: rgba(252, 226, 214, 0.62);
  font-size: 15px;
  font-weight: 920;
}

.step-body p {
  margin: 0;
  color: #4a3b33;
  font-size: 15px;
  font-weight: 650;
  line-height: 1.72;
}

.step-body em {
  display: block;
  margin-top: 10px;
  padding: 11px 13px;
  border-radius: 15px;
  color: #8b6530;
  background: rgba(255, 241, 184, 0.72);
  font-size: 13px;
  font-style: normal;
  font-weight: 680;
  line-height: 1.5;
}

.tips-card {
  display: flex;
  gap: 13px;
  color: #6b5142;
  background:
    radial-gradient(circle at 0% 0%, rgba(255, 255, 255, 0.72), transparent 48%),
    rgba(238, 246, 231, 0.78);
}

.tips-card svg {
  width: 24px;
  height: 24px;
  flex: 0 0 24px;
  color: var(--sage);
  stroke-width: 2.2;
}

.tips-card strong {
  color: #587343;
  font-size: 16px;
  font-weight: 920;
}

.tips-card p {
  margin: 7px 0 0;
  color: #6b5142;
  font-size: 14px;
  font-weight: 650;
  line-height: 1.6;
}

.feedback-actions {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.feedback-actions button {
  min-height: 48px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid rgba(143, 111, 86, 0.12);
  border-radius: 16px;
  color: #5d4b40;
  background: rgba(255, 255, 255, 0.54);
  font-size: 14px;
  font-weight: 880;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.66);
}

.feedback-actions button.active {
  color: #fff;
  border-color: transparent;
  background: linear-gradient(135deg, #ff7568, var(--coral));
  box-shadow: 0 10px 20px rgba(233, 86, 69, 0.18);
}

.feedback-actions button:disabled {
  cursor: wait;
  opacity: 0.72;
}

.feedback-actions svg {
  width: 19px;
  height: 19px;
  stroke-width: 2.3;
}

.bottom-actions {
  position: fixed;
  left: 50%;
  bottom: calc(16px + env(safe-area-inset-bottom));
  width: min(calc(100% - 42px), 388px);
  display: grid;
  grid-template-columns: minmax(0, 1fr) 58px 58px;
  gap: 12px;
  transform: translateX(-50%);
  z-index: 20;
}

.action-message,
.action-error {
  margin: 17px 2px 0;
  padding: 12px 14px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 800;
  line-height: 1.45;
}

.action-message {
  color: #416834;
  background: rgba(226, 241, 206, 0.82);
}

.action-error {
  color: #9b2f24;
  background: rgba(255, 225, 219, 0.88);
}

.primary-action,
.secondary-action {
  min-height: 58px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  border-radius: 20px;
  font-size: 16px;
  font-weight: 900;
  box-shadow:
    0 18px 34px rgba(80, 50, 28, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.28);
}

.primary-action {
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
}

.primary-action svg {
  width: 22px;
  height: 22px;
  stroke-width: 2.3;
}

.secondary-action {
  color: #765f51;
  background: rgba(255, 250, 240, 0.94);
  border: 1px solid rgba(255, 255, 255, 0.72);
}

.secondary-action svg {
  width: 24px;
  height: 24px;
  stroke-width: 2.25;
}

.secondary-action.active {
  color: var(--coral);
}

.secondary-action:disabled {
  cursor: wait;
  opacity: 0.7;
}

.primary-action:disabled {
  cursor: wait;
  opacity: 0.72;
}

.state-card {
  min-height: 300px;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 12px;
  text-align: center;
}

.state-card p {
  margin: 0;
  color: var(--sub);
  font-size: 14px;
  font-weight: 700;
}

.state-card h1 {
  margin: 0;
  color: var(--text);
  font-size: 23px;
  font-weight: 930;
}

.state-card.error svg {
  width: 40px;
  height: 40px;
  color: var(--coral);
  stroke-width: 2.2;
}

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid rgba(233, 86, 69, 0.18);
  border-top-color: var(--coral);
  border-radius: 50%;
  animation: spin 0.72s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.nav-btn:active,
.ingredient-row:active,
.primary-action:active,
.secondary-action:active,
.feedback-actions button:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .nav-btn:hover,
  .ingredient-row:hover,
  .primary-action:hover,
  .secondary-action:hover,
  .feedback-actions button:hover:not(:disabled) {
    transform: translateY(-1px);
  }
}

@media (max-width: 350px) {
  .recipe-phone {
    padding-right: 17px;
    padding-left: 17px;
  }

  .hero-copy h1 {
    font-size: 32px;
  }

  .stats-card strong {
    font-size: 19px;
  }

  .glass-card,
  .tips-card {
    padding: 19px;
    border-radius: 26px;
  }

  .ingredient-row {
    grid-template-columns: 28px minmax(0, 1fr);
  }

  .ingredient-amount,
  .owned-pill {
    grid-column: 2;
  }

  .bottom-actions {
    width: min(calc(100% - 34px), 396px);
  }
}
</style>
