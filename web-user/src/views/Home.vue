<template>
  <div class="home-shell" :style="homeVars">
    <div class="home-warm-overlay" aria-hidden="true"></div>

    <main class="home-phone">
      <header class="home-header">
        <div class="greeting">
          <div class="greeting-copy">
            <h1>早安，今天吃什么？</h1>
            <p>今天也好好吃饭</p>
          </div>
        </div>
        <button class="calendar-btn" type="button" aria-label="打开一周菜单" @click="router.push('/week-menu')">
          <svg viewBox="0 0 24 24"><rect x="4" y="5.5" width="16" height="15" rx="3"/><path d="M8 3.5v4M16 3.5v4M4 10h16M8 14h.01M12 14h.01M16 14h.01M8 17h.01M12 17h.01"/></svg>
        </button>
      </header>

      <div class="search-box" role="search" @click="goToSearch">
        <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="10.8" cy="10.8" r="7.2"/><path d="m16 16 4.2 4.2"/></svg>
        <input v-model="keyword" type="search" placeholder="搜索菜谱、食材或菜系" aria-label="搜索菜谱、食材或菜系" @keyup.enter="goToSearch" />
      </div>

      <section class="glass-card decision-panel" aria-labelledby="decision-title">
        <div class="section-title decision-title">
          <div>
            <h2 id="decision-title">今天怎么吃</h2>
          </div>
          <button type="button" @click="router.push('/recommend/scene')">更多</button>
        </div>
        <div class="main-actions" aria-label="主要用餐行动">
          <button v-for="item in mainActions" :key="item.label" class="main-action" type="button" @click="goQuickAction(item)">
            <span class="main-icon" :class="item.tone" aria-hidden="true">
              <svg v-if="item.icon === 'scene'" viewBox="0 0 48 48"><path d="M8 38 24 13l16 25"/><path d="M16 38V27l8-14 8 14v11M24 13v25"/></svg>
              <svg v-else-if="item.icon === 'ingredients'" viewBox="0 0 48 48"><path d="M25 22c-4-9-13-8-18-6 1 8 8 13 18 10"/><path d="M26 21c1-8 7-13 15-13-1 9-6 14-15 14"/><circle cx="27" cy="30" r="8"/></svg>
              <svg v-else viewBox="0 0 48 48"><rect x="11" y="12" width="26" height="28" rx="5"/><path d="M17 8v8M31 8v8M16 22h16M16 29h10"/></svg>
            </span>
            <strong>{{ item.label }}</strong>
            <small>{{ item.description }}</small>
          </button>
        </div>
      </section>

      <section class="glass-card chef-card" :class="{ empty: !todayRecipe }" @click="todayRecipe ? goToRecipe(todayRecipe) : goToRecipes()">
        <div v-if="todayRecipe?.cover" class="chef-photo" :style="coverStyle(todayRecipe.cover)" aria-hidden="true"></div>
        <div v-else class="chef-photo placeholder-photo" aria-hidden="true">
          <svg viewBox="0 0 24 24"><path d="M6 16h12M8 16v2.5a2.5 2.5 0 0 0 2.5 2.5h3a2.5 2.5 0 0 0 2.5-2.5V16"/><path d="M8 12a4 4 0 0 1 8 0M7 13h10"/></svg>
        </div>
        <button v-if="todayRecipe" class="save-btn" type="button" aria-label="收藏菜谱" @click.stop>
          <svg viewBox="0 0 24 24"><path d="M6 4h12v17l-6-4-6 4V4Z"/></svg>
        </button>
        <div class="chef-copy">
          <span class="chef-kicker"><svg viewBox="0 0 24 24"><path d="M6 16h12M8 16v2.5a2.5 2.5 0 0 0 2.5 2.5h3a2.5 2.5 0 0 0 2.5-2.5V16"/><path d="M8 12a4 4 0 0 1 8 0M7 13h10"/></svg>今日主厨建议</span>
          <h2>{{ todayRecipe ? recipeTitle(todayRecipe) : '暂无今日推荐' }}</h2>
          <p v-if="todayRecipe && recipeDescription(todayRecipe)">{{ recipeDescription(todayRecipe) }}</p>
          <span v-if="todayRecipe" class="sage-badge"><i></i>{{ activeMeal }} · {{ recipeDifficulty(todayRecipe) }}</span>
          <div v-if="todayRecipe" class="chef-actions">
            <button class="chef-action-primary" type="button" :disabled="shoppingActionLoading" @click.stop="addRecipeToShopping(todayRecipe)">
              {{ shoppingActionLoading ? '加入中...' : '加入清单' }}
            </button>
            <button class="chef-action-secondary" type="button" @click.stop="goToRecipe(todayRecipe)">查看做法</button>
            <button class="chef-action-secondary" type="button" @click.stop="router.push('/recommend/scene')">换一组</button>
          </div>
        </div>
      </section>

      <p v-if="homeMessage" class="home-message" role="status">{{ homeMessage }}</p>
      <p v-if="homeError" class="home-error" role="alert">{{ homeError }}</p>

      <nav class="meal-tabs" aria-label="餐段切换">
        <button v-for="tab in mealTabs" :key="tab" :class="{ active: activeMeal === tab }" type="button" @click="activeMeal = tab">{{ tab }}</button>
      </nav>

      <section v-if="menuRecipe" class="glass-card recipe-note" @click="goToRecipe(menuRecipe)">
        <div v-if="menuRecipe.cover" class="note-photo" :style="coverStyle(menuRecipe.cover)" aria-hidden="true"></div>
        <div v-else class="note-photo no-cover" aria-hidden="true"><span>暂无图片</span></div>
        <div class="note-copy">
          <button class="heart-btn" type="button" aria-label="收藏菜谱" @click.stop><svg viewBox="0 0 24 24"><path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z"/></svg></button>
          <h2>{{ recipeTitle(menuRecipe) }}</h2>
          <div class="note-meta">
            <span><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8.5"/><path d="M12 7v5l3 2"/></svg>{{ recipeTime(menuRecipe) }}</span>
            <span><svg viewBox="0 0 24 24"><path d="M4 20V9M10 20V4M16 20v-8M22 20H2"/></svg>{{ recipeDifficulty(menuRecipe) }}</span>
            <span><svg viewBox="0 0 24 24"><circle cx="12" cy="7" r="4"/><path d="M5 21a7 7 0 0 1 14 0"/></svg>{{ recipePeople(menuRecipe) }}</span>
          </div>
          <button class="primary-pill note-btn" type="button" @click.stop="goToRecipe(menuRecipe)">查看做法 <svg viewBox="0 0 24 24"><path d="m9 18 6-6-6-6"/></svg></button>
        </div>
      </section>

      <section v-else class="glass-card empty-card">
        <h2>今日菜单</h2>
        <p>{{ loading ? '正在读取后端推荐...' : emptyText }}</p>
      </section>

      <section class="glass-card quick-panel" aria-labelledby="quick-title">
        <div class="section-title quick-title">
          <h2 id="quick-title">更多选择</h2>
        </div>
        <div class="quick-grid" aria-label="更多用餐入口">
          <button v-for="item in secondaryActions" :key="item.label" class="quick-item" type="button" @click="goQuickAction(item)">
            <span class="quick-icon" :class="item.tone" aria-hidden="true">
              <svg v-if="item.icon === 'ingredients'" viewBox="0 0 48 48"><path d="M25 22c-4-9-13-8-18-6 1 8 8 13 18 10"/><path d="M26 21c1-8 7-13 15-13-1 9-6 14-15 14"/><circle cx="27" cy="30" r="8"/></svg>
              <svg v-else-if="item.icon === 'taste'" viewBox="0 0 48 48"><path d="M24 40c9-6 12-17 8-29-4 6-11 7-14 12-4 7-1 13 6 17Z"/><path d="M28 14c2-2 5-3 8-3"/></svg>
              <svg v-else-if="item.icon === 'scene'" viewBox="0 0 48 48"><path d="M8 38 24 13l16 25"/><path d="M16 38V27l8-14 8 14v11M24 13v25"/></svg>
              <svg v-else-if="item.icon === 'week'" viewBox="0 0 48 48"><rect x="11" y="12" width="26" height="28" rx="5"/><path d="M17 8v8M31 8v8M16 22h16M16 29h10"/></svg>
              <svg v-else-if="item.icon === 'fridge'" viewBox="0 0 48 48"><rect x="13" y="7" width="22" height="34" rx="4"/><path d="M13 21h22M19 14h3M19 28h3"/></svg>
              <svg v-else viewBox="0 0 48 48"><path d="M24 41c8-4 13-11 13-19 0-7-4-12-9-16-1 8-8 10-11 15-4 7-1 15 7 20Z"/><path d="M24 41c-3-5-1-9 4-13"/></svg>
            </span>
            <span>{{ item.label }}</span>
          </button>
        </div>
      </section>

      <section class="glass-card couple-card" role="button" tabindex="0" @click="router.push('/couple')">
        <div class="couple-photo" aria-hidden="true"></div>
        <div class="couple-mask" aria-hidden="true"></div>
        <div class="couple-copy">
          <h2>和 TA 一起决定</h2>
          <button class="primary-pill" type="button" @click.stop="router.push('/couple')">去选菜 <svg viewBox="0 0 24 24"><path d="m9 18 6-6-6-6"/></svg></button>
        </div>
      </section>

      <section class="glass-card popular-panel">
        <div class="section-title">
          <h2>热门菜谱</h2>
          <button type="button" @click="goHotRecipes">查看更多</button>
        </div>
        <div v-if="popularRecipes.length" class="popular-list">
          <article v-for="recipe in popularRecipes" :key="recipeKey(recipe)" class="popular-card" @click="goToRecipe(recipe)">
            <div v-if="recipe.cover" class="popular-photo" :style="coverStyle(recipe.cover)" aria-hidden="true"></div>
            <div v-else class="popular-photo no-cover" aria-hidden="true"><span>暂无图片</span></div>
            <div class="popular-copy">
              <h3>{{ recipeTitle(recipe) }}</h3>
              <span>{{ recipeTaste(recipe) }}</span>
              <p>{{ recipeTime(recipe) }} · {{ recipePeople(recipe) }}</p>
              <strong><svg viewBox="0 0 24 24"><path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z"/></svg>{{ recipeHeat(recipe) }}</strong>
            </div>
          </article>
        </div>
        <div v-else class="popular-empty">{{ loading ? '正在读取热门菜谱...' : '后端暂无热门菜谱' }}</div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getHomeData } from '@/api/home'
import { useShoppingStore } from '@/stores/shopping'
import { trackEvent } from '@/utils/event'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import coupleImage from '@/assets/home/couple-dining.jpg'

type HomeRecipe = {
  id?: number
  title?: string
  cover?: string
  description?: string
  cook_time?: number
  difficulty?: string
  taste?: string
  people_count?: number
  servings?: number
  favorite_count?: number
  category_name?: string
}

type QuickAction = {
  label: string
  icon: string
  tone: string
  path: string
  description?: string
  query?: Record<string, string>
}

const router = useRouter()
const shoppingStore = useShoppingStore()

const keyword = ref('')
const activeMeal = ref('晚餐')
const loading = ref(true)
const shoppingActionLoading = ref(false)
const homeMessage = ref('')
const homeError = ref('')
const todayRecommend = ref<HomeRecipe | null>(null)
const hotRecipes = ref<HomeRecipe[]>([])

const mealTabs = ['早餐', '午餐', '晚餐', '夜宵']
const emptyText = '暂无推荐菜谱'

const homeVars = computed(() => ({
  '--home-bg': 'url(' + kitchenBg + ')',
  '--couple-img': 'url(' + coupleImage + ')',
}))

const todayRecipe = computed(() => todayRecommend.value)
const menuRecipe = computed(() => hotRecipes.value[0] || todayRecommend.value || null)
const popularRecipes = computed(() => hotRecipes.value.slice(0, 2))

const mainActions: QuickAction[] = [
  { label: '快速推荐', description: '马上决定', icon: 'scene', tone: 'sage', path: '/recommend/scene' },
  { label: '现有食材', description: '冰箱优先', icon: 'ingredients', tone: 'tomato', path: '/recommend/ingredients' },
  { label: '本周菜单', description: '三餐安排', icon: 'week', tone: 'wheat', path: '/week-menu' },
]

const secondaryActions: QuickAction[] = [
  { label: '按口味找菜', icon: 'taste', tone: 'pepper', path: '/recommend/taste' },
  { label: '冰箱库存', icon: 'fridge', tone: 'sage', path: '/recommend/fridge' },
  { label: '换点新菜', icon: 'scene', tone: 'coral', path: '/recommend/new' },
  { label: '最热菜谱', icon: 'hot', tone: 'coral', path: '/recipes', query: { sort: 'hot' } },
]

function coverStyle(cover: string) {
  return { backgroundImage: 'url(' + cover + ')' }
}

function recipeKey(recipe: HomeRecipe) {
  return String(recipe.id || recipe.title || Math.random())
}

function recipeTitle(recipe: HomeRecipe) {
  return recipe.title || '未命名菜谱'
}

function recipeDescription(recipe: HomeRecipe) {
  return recipe.description || recipe.category_name || recipe.taste || ''
}

function recipeTime(recipe: HomeRecipe) {
  return recipe.cook_time ? String(recipe.cook_time) + '分钟' : '时间待补充'
}

function recipePeople(recipe: HomeRecipe) {
  const count = recipe.people_count || recipe.servings
  return count ? String(count) + '人份' : '人份待补充'
}

function recipeDifficulty(recipe: HomeRecipe) {
  return recipe.difficulty || '难度待补充'
}

function recipeTaste(recipe: HomeRecipe) {
  return recipe.taste || recipe.category_name || recipeDifficulty(recipe)
}

function recipeHeat(recipe: HomeRecipe) {
  const count = recipe.favorite_count || 0
  if (count >= 10000) return (count / 10000).toFixed(1) + 'w'
  if (count >= 1000) return (count / 1000).toFixed(1) + 'k'
  return String(count)
}

function goToRecipe(recipe: HomeRecipe | null) {
  if (recipe?.id) {
    router.push('/recipes/' + recipe.id)
    return
  }
  goToRecipes()
}

function goToRecipes() {
  router.push('/recipes')
}

function goHotRecipes() {
  router.push({ path: '/recipes', query: { sort: 'hot' } })
}

function goToSearch() {
  const value = keyword.value.trim()
  if (value) {
    router.push({ path: '/recipes', query: { keyword: value } })
    return
  }
  goToRecipes()
}

function goQuickAction(item: QuickAction) {
  trackEvent({
    event_name: 'home_action_click',
    entity_type: 'entry',
    payload: { label: item.label, path: item.path },
  })
  router.push(item.query ? { path: item.path, query: item.query } : item.path)
}

async function addRecipeToShopping(recipe: HomeRecipe | null) {
  if (!recipe || shoppingActionLoading.value) return
  const title = recipeTitle(recipe)
  shoppingActionLoading.value = true
  homeMessage.value = ''
  homeError.value = ''
  try {
    if (recipe.id) {
      await shoppingStore.generateByRecipe(recipe.id, title)
    } else {
      await shoppingStore.generateByDish(title)
    }
    trackEvent({
      event_name: 'add_shopping_list',
      entity_type: 'recipe',
      entity_id: recipe.id || 0,
      payload: { source: 'home', title },
    })
    homeMessage.value = '「' + title + '」的食材已合并到购物清单。'
  } catch (error) {
    homeError.value = error instanceof Error ? error.message : '加入购物清单失败，请稍后重试'
  } finally {
    shoppingActionLoading.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const res: any = await getHomeData()
    todayRecommend.value = res?.today_recommend || null
    hotRecipes.value = Array.isArray(res?.hot_recipes) ? res.hot_recipes : []
  } catch {
    todayRecommend.value = null
    hotRecipes.value = []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.home-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --coral: #e84f3f;
  --coral-2: #e95645;
  --surface: var(--card-surface);
  --surface-strong: var(--card-surface-strong);
  --surface-border: var(--card-border);
  --surface-shadow: var(--card-shadow);
  --surface-radius: var(--card-radius);
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: clip;
  color: var(--text);
  background: linear-gradient(180deg, rgba(255, 248, 236, 0.22), rgba(248, 229, 201, 0.44)), var(--home-bg) center / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.home-warm-overlay {
  position: fixed;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, rgba(255, 250, 242, 0.38), rgba(117, 78, 46, 0.08));
  z-index: 0;
}

.home-phone {
  position: relative;
  z-index: 1;
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: max(28px, env(safe-area-inset-top)) 24px calc(var(--tab-h, 82px) + 78px);
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  margin-top: 0;
}

.greeting {
  display: flex;
  align-items: center;
  min-width: 0;
}

.greeting-copy {
  min-width: 0;
}

.greeting h1 {
  margin: 0;
  color: var(--text);
  font-size: 23px;
  font-weight: 850;
  line-height: 1.16;
  letter-spacing: 0;
  white-space: nowrap;
}

.greeting p {
  margin: 7px 0 0;
  color: var(--sub);
  font-size: 14px;
  line-height: 1.35;
  white-space: nowrap;
}

.calendar-btn,
.save-btn,
.heart-btn {
  display: inline-grid;
  place-items: center;
  border: 1px solid var(--surface-border);
  background: var(--surface-strong);
  color: #593a2e;
  box-shadow: var(--surface-shadow);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  cursor: pointer;
  transition: transform 180ms ease, opacity 180ms ease, box-shadow 180ms ease;
}

.calendar-btn {
  width: 50px;
  height: 50px;
  flex: 0 0 50px;
  border-radius: 16px;
}

svg {
  fill: none;
  stroke: currentColor;
  stroke-width: 1.9;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.calendar-btn svg,
.save-btn svg,
.heart-btn svg {
  width: 25px;
  height: 25px;
}

.search-box {
  position: relative;
  height: 52px;
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 22px;
  padding: 0 17px;
  border: 1px solid var(--surface-border);
  border-radius: 18px;
  background: var(--surface);
  box-shadow: var(--surface-shadow);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.search-box svg {
  width: 25px;
  height: 25px;
  flex: 0 0 25px;
  color: #3c2d27;
  stroke-width: 2.1;
}

.search-box input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--text);
  font-size: 16px;
  line-height: 1;
}

.search-box input::placeholder {
  color: #988b81;
}

.search-box input {
  position: relative;
  z-index: 1;
}

.search-box:focus-within {
  border-color: rgba(232, 79, 63, 0.38);
  box-shadow: 0 0 0 3px rgba(232, 79, 63, 0.1), var(--surface-shadow);
}

.glass-card {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--surface-border);
  border-radius: var(--surface-radius);
  background: var(--surface);
  box-shadow: var(--surface-shadow);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  transition: transform 170ms ease, box-shadow 170ms ease;
}

.couple-card {
  height: 166px;
  margin-top: 16px;
  cursor: pointer;
}

.decision-panel {
  margin-top: 18px;
  padding: 18px 16px 12px;
}

.decision-title {
  margin-bottom: 14px;
}

.main-actions {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0;
}

.main-action {
  min-width: 0;
  min-height: 112px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 12px 10px;
  border: 0;
  border-left: 1px solid var(--surface-border);
  border-radius: 0;
  background: transparent;
  color: #3d2d25;
  text-align: center;
  cursor: pointer;
  transition: transform 170ms ease, background 170ms ease;
}

.main-action:first-child {
  border-left: 0;
}

.main-action strong {
  display: block;
  color: #34251f;
  font-size: 15px;
  font-weight: 950;
  line-height: 1.15;
}

.main-action small {
  display: block;
  color: #806f64;
  font-size: 11px;
  font-weight: 750;
  line-height: 1.35;
}

.main-icon,
.main-icon svg {
  width: 36px;
  height: 36px;
}

.main-icon svg {
  fill: rgba(232, 79, 63, 0.1);
  stroke-width: 2.3;
}

.main-icon.tomato { color: #cf4b38; }
.main-icon.sage { color: #73936d; }
.main-icon.wheat { color: #d48b3d; }

.couple-photo {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 70%;
  background-image: var(--couple-img);
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  transform: scale(1.02);
}

.couple-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, rgba(255, 249, 238, 0.92) 0%, rgba(255, 249, 238, 0.74) 35%, rgba(255, 249, 238, 0.18) 68%, rgba(255, 249, 238, 0.04) 100%), linear-gradient(180deg, rgba(255, 255, 255, 0.16), rgba(255, 242, 223, 0.14));
}

.couple-copy {
  position: relative;
  z-index: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 25px 0 24px 29px;
  box-sizing: border-box;
}

.couple-copy h2 {
  margin: 0;
  color: #32231d;
  font-size: 29px;
  font-weight: 900;
  line-height: 1.1;
  letter-spacing: 0;
}

.couple-copy h2 span {
  color: var(--coral);
  font-weight: 500;
}

.couple-copy .primary-pill {
  margin-top: auto;
}

.primary-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  min-height: 44px;
  padding: 0 20px;
  border: 0;
  border-radius: 14px;
  background: var(--coral);
  color: #fff;
  box-shadow: 0 8px 18px rgba(218, 64, 50, 0.2);
  font-size: 15px;
  font-weight: 800;
  cursor: pointer;
  transition: transform 170ms ease, opacity 170ms ease, box-shadow 170ms ease;
}

.primary-pill svg {
  width: 17px;
  height: 17px;
  stroke-width: 2.5;
}

.chef-card {
  min-height: 250px;
  margin-top: 16px;
  cursor: pointer;
}

.chef-photo {
  position: absolute;
  top: 0;
  right: -18px;
  bottom: 0;
  width: 68%;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
}

.chef-card::before {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, rgba(255, 252, 246, 0.97) 0%, rgba(255, 252, 246, 0.84) 44%, rgba(255, 252, 246, 0.14) 78%);
  z-index: 1;
}

.placeholder-photo,
.no-cover {
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, rgba(255, 252, 246, 0.72), rgba(236, 217, 190, 0.7));
  color: #b09b8c;
}

.placeholder-photo svg {
  width: 74px;
  height: 74px;
  stroke-width: 1.4;
}

.save-btn {
  position: absolute;
  top: 20px;
  right: 18px;
  z-index: 2;
  width: 44px;
  height: 44px;
  border-color: transparent;
  background: transparent;
  color: var(--coral);
  box-shadow: none;
}

.chef-copy {
  position: relative;
  z-index: 2;
  min-height: 250px;
  max-width: 58%;
  padding: 24px 0 72px 27px;
}

.chef-kicker {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  color: #b35a35;
  font-size: 14px;
  font-weight: 800;
}

.chef-kicker svg {
  width: 20px;
  height: 20px;
  color: var(--coral);
}

.chef-copy h2 {
  margin: 15px 0 13px;
  color: #34251f;
  font-size: 31px;
  font-weight: 950;
  line-height: 1.08;
}

.chef-copy p {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;
  margin: 0;
  color: #765f53;
  font-size: 15px;
  line-height: 1.55;
}

.chef-actions {
  position: absolute;
  left: 27px;
  right: -128px;
  bottom: 22px;
  display: grid;
  grid-template-columns: minmax(92px, auto) auto auto;
  align-items: center;
  justify-content: start;
  gap: 8px;
  margin: 0;
}

.chef-actions button {
  min-height: 44px;
  padding: 0 13px;
  border: 1px solid rgba(92, 63, 45, 0.12);
  border-radius: 13px;
  background: rgba(255, 253, 249, 0.9);
  color: #4b352c;
  font-size: 12px;
  font-weight: 850;
  cursor: pointer;
}

.chef-actions .chef-action-primary {
  min-width: 96px;
  color: #fff;
  border-color: transparent;
  background: var(--coral);
  box-shadow: 0 8px 16px rgba(218, 64, 50, 0.18);
}

.chef-actions .chef-action-secondary {
  min-width: 72px;
}

.chef-actions button:disabled {
  cursor: wait;
  opacity: 0.68;
}

.home-message,
.home-error {
  margin: 12px 2px 0;
  padding: 11px 13px;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 800;
  line-height: 1.45;
}

.home-message {
  color: #416834;
  background: rgba(226, 241, 206, 0.82);
}

.home-error {
  color: #9b2f24;
  background: rgba(255, 225, 219, 0.88);
}

.sage-badge {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  margin-top: 18px;
  padding: 6px 12px;
  border-radius: 10px;
  background: rgba(213, 227, 190, 0.74);
  color: #4e7048;
  font-size: 13px;
  font-weight: 800;
}

.sage-badge i {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #7fa06f;
  box-shadow: inset 6px 0 0 rgba(255, 255, 255, 0.72);
}

.meal-tabs {
  height: 56px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  align-items: center;
  gap: 4px;
  margin-top: 16px;
  padding: 5px;
  border: 1px solid var(--surface-border);
  border-radius: 18px;
  background: var(--surface);
  box-shadow: var(--surface-shadow);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.meal-tabs button {
  height: 44px;
  border: 0;
  border-radius: 13px;
  background: transparent;
  color: #4a352b;
  font-size: 17px;
  font-weight: 800;
  cursor: pointer;
  transition: color 180ms ease, background 180ms ease, box-shadow 180ms ease, transform 180ms ease;
}

.meal-tabs button.active {
  background: var(--surface-strong);
  color: var(--coral);
  box-shadow: 0 4px 12px rgba(91, 60, 33, 0.1);
}

.recipe-note {
  min-height: 158px;
  display: grid;
  grid-template-columns: 40% 1fr;
  gap: 15px;
  margin-top: 16px;
  padding: 14px;
  cursor: pointer;
}

.note-photo {
  min-height: 130px;
  border-radius: 14px;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  box-shadow: none;
}

.no-cover span {
  font-size: 13px;
  font-weight: 800;
}

.note-copy {
  position: relative;
  min-height: 130px;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 1px 0 0;
}

.heart-btn {
  position: absolute;
  top: -6px;
  right: -7px;
  width: 44px;
  height: 44px;
  border: 0;
  background: transparent;
  color: #df6a42;
  box-shadow: none;
}

.note-copy h2 {
  margin: 0 34px 10px 0;
  color: #3b2a22;
  font-size: 20px;
  font-weight: 900;
  line-height: 1.15;
}

.note-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 11px;
  color: #7e7067;
  font-size: 12px;
  font-weight: 700;
}

.note-meta span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.note-meta svg,
.popular-copy strong svg {
  width: 14px;
  height: 14px;
}

.note-copy p {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  margin: 13px 0 12px;
  color: #7b6659;
  font-size: 13px;
  line-height: 1.52;
}

.note-btn {
  min-height: 44px;
  align-self: flex-end;
  margin-top: auto;
  padding: 0 16px;
  font-size: 14px;
}

.empty-card {
  margin-top: 16px;
  padding: 22px;
}

.empty-card h2 {
  margin: 0 0 8px;
  color: #3b2a22;
  font-size: 20px;
  font-weight: 900;
}

.empty-card p,
.popular-empty {
  margin: 0;
  color: #806f64;
  font-size: 14px;
  line-height: 1.55;
}

.quick-panel {
  margin-top: 16px;
  padding: 16px;
}

.quick-title {
  margin-bottom: 6px;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.quick-item {
  min-width: 0;
  min-height: 76px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  gap: 10px;
  padding: 10px 8px;
  border: 0;
  border-left: 1px solid var(--surface-border);
  border-top: 1px solid var(--surface-border);
  border-radius: 0;
  background: transparent;
  color: #49372f;
  font-size: 13px;
  font-weight: 800;
  text-align: left;
  cursor: pointer;
  transition: transform 170ms ease, background 170ms ease;
}

.quick-item:nth-child(odd) {
  border-left: 0;
}

.quick-item:nth-child(-n + 2) {
  border-top: 0;
}

.quick-icon,
.quick-icon svg {
  width: 36px;
  height: 36px;
  flex: 0 0 36px;
}

.quick-icon svg {
  fill: rgba(232, 79, 63, 0.1);
  stroke-width: 2.3;
}

.quick-icon.tomato { color: #cf4b38; }
.quick-icon.pepper { color: #df4f3e; }
.quick-icon.sage { color: #73936d; }
.quick-icon.wheat { color: #d48b3d; }
.quick-icon.coral { color: #e84f3f; }

.popular-panel {
  margin-top: 16px;
  padding: 16px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-title h2 {
  margin: 0;
  color: #2f231d;
  font-size: 20px;
  font-weight: 900;
}

.section-title button {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  padding: 0 2px 0 10px;
  border: 0;
  background: transparent;
  color: #74645a;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
}

.section-title button::after {
  content: "›";
  margin-left: 5px;
  font-size: 18px;
  line-height: 0;
}

.popular-list {
  display: grid;
  grid-template-columns: 1fr;
}

.popular-card {
  display: grid;
  grid-template-columns: 92px minmax(0, 1fr);
  min-height: 92px;
  overflow: hidden;
  padding: 12px 0;
  border-top: 1px solid var(--surface-border);
  background: transparent;
  cursor: pointer;
  transition: transform 170ms ease, background 170ms ease;
}

.popular-card:first-child {
  padding-top: 0;
  border-top: 0;
}

.popular-card:last-child {
  padding-bottom: 0;
}

.popular-photo {
  min-height: 92px;
  border-radius: 13px;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
}

.popular-copy {
  min-width: 0;
  padding: 4px 4px 3px 13px;
}

.popular-copy h3 {
  margin: 0;
  color: #3b2a22;
  font-size: 16px;
  font-weight: 900;
  line-height: 1.1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.popular-copy span {
  display: inline-flex;
  max-width: 100%;
  margin-top: 6px;
  padding: 2px 7px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(243, 196, 126, 0.44);
  color: #a85f23;
  font-size: 11px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.popular-copy p {
  margin: 7px 0 0;
  color: #806f64;
  font-size: 12px;
  line-height: 1.2;
  white-space: nowrap;
}

.popular-copy strong {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 7px;
  color: #79685d;
  font-size: 12px;
  font-weight: 800;
}

.primary-pill:active,
.calendar-btn:active,
.main-action:active,
.quick-item:active,
.chef-card:active,
.recipe-note:active,
.couple-card:active,
.popular-card:active {
  transform: scale(0.98);
}

.home-shell button:focus-visible,
.couple-card:focus-visible {
  outline: 3px solid rgba(232, 79, 63, 0.3);
  outline-offset: 2px;
}

@media (hover: hover) {
  .primary-pill:hover,
  .main-action:hover,
  .quick-item:hover,
  .calendar-btn:hover {
    transform: translateY(-1px);
  }

  .main-action:hover,
  .quick-item:hover {
    background: rgba(232, 79, 63, 0.05);
  }

  .popular-card:hover {
    background: rgba(232, 79, 63, 0.035);
  }

  .chef-card:hover,
  .recipe-note:hover,
  .couple-card:hover {
    box-shadow: var(--card-shadow-feature);
  }
}

@media (max-width: 380px) {
  .home-phone {
    padding-left: 18px;
    padding-right: 18px;
  }

  .greeting {
    gap: 12px;
  }

  .greeting h1 {
    font-size: 21px;
  }

  .greeting p {
    font-size: 13px;
  }

  .chef-copy {
    max-width: 61%;
    padding-left: 22px;
  }

  .chef-actions {
    left: 22px;
    right: -112px;
    grid-template-columns: minmax(88px, auto) auto auto;
    gap: 7px;
  }

  .chef-copy h2 {
    font-size: 28px;
  }

  .recipe-note {
    grid-template-columns: 38% 1fr;
    gap: 12px;
    padding: 12px;
  }
}

@media (max-width: 360px) {
  .popular-card {
    grid-template-columns: 84px minmax(0, 1fr);
  }

  .popular-photo {
    min-height: 84px;
  }
}

@media (min-width: 431px) {
  .home-shell {
    background-color: #ead7bd;
  }

  .home-phone {
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.24);
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
