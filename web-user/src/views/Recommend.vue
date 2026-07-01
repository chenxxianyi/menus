<template>
  <div class="recommend-shell" :style="pageVars">
    <div class="recommend-warm-overlay" aria-hidden="true"></div>

    <main class="recommend-phone">
      <header class="recommend-topbar">
        <button class="nav-btn" type="button" aria-label="返回" @click="router.back()">
          <svg viewBox="0 0 24 24"><path d="m15 18-6-6 6-6" /></svg>
        </button>
        <div>
          <p>{{ modeMeta.eyebrow }}</p>
          <h1>{{ modeMeta.title }}</h1>
        </div>
      </header>

      <section class="hero-card">
        <span class="hero-badge">
          <svg viewBox="0 0 24 24"><path d="M12 3 13.7 8.3 19 10l-5.3 1.7L12 17l-1.7-5.3L5 10l5.3-1.7L12 3Z" /></svg>
          {{ modeMeta.badge }}
        </span>
        <h2>{{ modeMeta.heading }}</h2>
        <p>{{ modeMeta.description }}</p>
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
          <span v-if="!selectedIngredients.length" class="hint-chip">至少添加 1 项，最多 20 项</span>
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
            <span>查看数据库中标记为“{{ taste }}”的菜谱</span>
          </button>
        </div>
        <p v-else class="data-empty">数据库暂无口味选项，请先在后台维护菜谱口味字段。</p>
      </section>

      <section v-else class="panel">
        <div class="panel-title">
          <h2>按生活场景安排</h2>
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
            <span>{{ scene.description }}</span>
          </button>
        </div>
        <div class="scene-actions">
          <button class="submit-btn" type="button" :disabled="!selectedScene || loading || aiSceneLoading" @click="submitScene">
            <span v-if="loading" class="mini-spinner" aria-hidden="true"></span>
            <span>{{ loading ? '正在搭配...' : '从菜谱库推荐' }}</span>
          </button>
          <button class="submit-btn ai-submit-btn" type="button" :disabled="!selectedScene || loading || aiSceneLoading" @click="submitAIScene">
            <span v-if="aiSceneLoading" class="mini-spinner" aria-hidden="true"></span>
            <span>{{ aiSceneLoading ? 'AI 正在生成...' : 'AI 按偏好生成新菜' }}</span>
          </button>
        </div>
        <p class="scene-ai-hint">AI 会读取你的偏好设置，生成没吃过也适合你的新菜，并同步到菜谱库。</p>
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
            <p>{{ item.reason || '根据食材匹配推荐' }}</p>
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
          <span>{{ sceneResult.menu_name || selectedScene?.title }}</span>
          <h2>{{ sceneResult.reason || '为你搭配好了这一餐' }}</h2>
        </div>
        <article v-for="dish in sceneDishes" :key="dish.recipe_id || dish.name" class="dish-card" @click="openRecipe(dish.recipe_id)">
          <div>
            <strong>{{ dish.name || '未命名菜谱' }}</strong>
            <p>{{ dish.type || '推荐菜' }} · {{ dish.cook_time ? dish.cook_time + '分钟' : '时间待补充' }} · {{ dish.difficulty || '难度待补充' }}</p>
          </div>
          <svg viewBox="0 0 24 24"><path d="m9 18 6-6-6-6" /></svg>
        </article>
      </section>

      <section v-if="hasSubmitted && !loading && isIngredientMode && !results.length && !error" class="empty-panel">
        <h2>暂时没有匹配菜谱</h2>
        <p>可以减少食材限制，或先在后台补充更多菜谱食材数据。</p>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getIngredientOptions, recommendByIngredients, recommendByScene, recommendSceneByAI } from '@/api/recommend'
import { getRecipeFilterOptions } from '@/api/recipe'
import { updateShoppingList } from '@/api/shopping'
import { useShoppingStore } from '@/stores/shopping'
import type { IngredientOption, RecommendMode, RecommendRecipeResult, SceneOption } from '@/types/recommend'
import type { ShoppingItem } from '@/api/shopping'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'

const route = useRoute()
const router = useRouter()
const shoppingStore = useShoppingStore()

const mode = computed<RecommendMode>(() => {
  const value = String(route.params.mode || 'ingredients')
  return ['ingredients', 'taste', 'scene', 'fridge'].includes(value) ? value as RecommendMode : 'ingredients'
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
const error = ref('')
const message = ref('')
const hasSubmitted = ref(false)

const pageVars = computed(() => ({ '--recommend-bg': `url(${kitchenBg})` }))

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
      eyebrow: 'Ingredient match',
      title: '按食材推荐',
      badge: '食材匹配',
      heading: '把手头食材变成一顿饭',
      description: '选择已有食材，系统会按匹配度展示菜谱，并告诉你还缺什么。',
    },
    taste: {
      eyebrow: 'Taste finder',
      title: '按口味推荐',
      badge: '口味筛选',
      heading: '今天想吃什么味道？',
      description: '选择一个口味后进入菜谱列表，刷新和返回都会保留筛选条件。',
    },
    scene: {
      eyebrow: 'Scene menu',
      title: '按场景推荐',
      badge: '场景搭配',
      heading: '让场景替你做选择',
      description: '快手、聚餐、减脂、宴客、夜宵，不同场景用不同推荐规则。',
    },
    fridge: {
      eyebrow: 'Fridge rescue',
      title: '冰箱剩菜',
      badge: '剩菜拯救',
      heading: '先消耗冰箱里的库存',
      description: '录入现有食材，优先推荐能直接做或只差少量食材的菜谱。',
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

async function submitIngredients() {
  if (!selectedIngredients.value.length || loading.value) return
  loading.value = true
  error.value = ''
  message.value = ''
  hasSubmitted.value = true
  try {
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
  message.value = ''
  error.value = ''
}

async function submitScene() {
  if (!selectedScene.value || loading.value || aiSceneLoading.value) return
  loading.value = true
  error.value = ''
  message.value = ''
  try {
    sceneResult.value = await recommendByScene({
      scene: selectedScene.value.key,
      meal_type: selectedScene.value.meal_type,
      people_count: selectedScene.value.people_count,
      cook_time_preference: selectedScene.value.cook_time_preference,
      health_goal: selectedScene.value.health_goal,
    })
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
  error.value = ''
  message.value = ''
  try {
    sceneResult.value = await recommendSceneByAI({
      scene: selectedScene.value.key,
      meal_type: selectedScene.value.meal_type,
      people_count: selectedScene.value.people_count,
      cook_time_preference: selectedScene.value.cook_time_preference,
      health_goal: selectedScene.value.health_goal,
    })
    message.value = 'AI 已根据你的偏好生成新菜，并同步到菜谱库。'
  } catch (e: any) {
    sceneResult.value = null
    error.value = e?.message || 'AI 场景推荐失败，请确认 AI 配置后重试。'
  } finally {
    aiSceneLoading.value = false
  }
}

function clearSceneResult() {
  selectedScene.value = null
  sceneResult.value = null
  error.value = ''
  message.value = ''
}

function recipeTitle(recipe: { title?: string }) {
  return recipe.title || '未命名菜谱'
}

function openRecipe(id?: number) {
  if (id) router.push('/recipes/' + id)
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

watch(mode, () => {
  error.value = ''
  message.value = ''
  results.value = []
  sceneResult.value = null
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

onMounted(() => {
  loadIngredientOptions()
  loadTasteOptions()
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
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
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

.recommend-topbar p {
  margin: 0 0 5px;
  color: #c26d3d;
  font-size: 14px;
  font-weight: 850;
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
.result-card,
.scene-result,
.empty-panel {
  border: 1px solid rgba(255, 255, 255, 0.66);
  border-radius: 28px;
  background: var(--cream);
  box-shadow: 0 18px 40px rgba(80, 50, 30, 0.15);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.hero-card {
  margin-top: 24px;
  padding: 24px;
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
  margin: 17px 0 10px;
  font-size: 27px;
  font-weight: 950;
  line-height: 1.12;
}

.hero-card p {
  margin: 0;
  color: var(--sub);
  font-size: 15px;
  line-height: 1.6;
}

.panel {
  margin-top: 18px;
  padding: 20px;
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
  font-size: 20px;
  font-weight: 950;
}

.panel-title button {
  border: 0;
  background: transparent;
  color: var(--coral);
  font-size: 14px;
  font-weight: 850;
  cursor: pointer;
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
  border-radius: 16px;
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

.scene-ai-hint {
  margin: 12px 2px 0;
  color: #8a7162;
  font-size: 13px;
  font-weight: 750;
  line-height: 1.5;
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
  min-height: 82px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.66);
  border-radius: 20px;
  text-align: left;
  color: #33241e;
  background: rgba(255, 255, 255, 0.58);
  cursor: pointer;
}

.taste-grid strong,
.scene-list strong {
  display: block;
  font-size: 18px;
  font-weight: 950;
}

.taste-grid span,
.scene-list span {
  display: block;
  margin-top: 7px;
  color: var(--sub);
  font-size: 13px;
  line-height: 1.42;
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
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.result-actions button {
  min-height: 44px;
  border: 0;
  border-radius: 15px;
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
  font-size: 21px;
  line-height: 1.3;
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

.dish-card strong {
  font-size: 17px;
  font-weight: 950;
}

.dish-card p {
  margin: 5px 0 0;
  color: var(--sub);
  font-size: 13px;
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
  .result-actions {
    grid-template-columns: 1fr;
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
