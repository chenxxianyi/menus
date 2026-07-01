<template>
  <div class="home-shell" :style="homeVars">
    <div class="home-warm-overlay" aria-hidden="true"></div>

    <main class="home-phone">
      <header class="home-header">
        <div class="greeting">
          <div class="greeting-copy">
            <h1>早安，今天吃什么？</h1>
            <p>用心规划每一餐，好好吃饭，好好生活</p>
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

      <section class="glass-card couple-card" role="button" tabindex="0" @click="router.push('/couple')">
        <div class="couple-photo" aria-hidden="true"></div>
        <div class="couple-mask" aria-hidden="true"></div>
        <div class="couple-copy">
          <h2>双人点餐 <span aria-hidden="true">♡</span></h2>
          <p>一起选菜，一起下厨<br />为你们定制美味晚餐</p>
          <button class="primary-pill" type="button" @click.stop="router.push('/couple')">去选菜 <svg viewBox="0 0 24 24"><path d="m9 18 6-6-6-6"/></svg></button>
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
          <p>{{ todayRecipe ? recipeDescription(todayRecipe) : emptyText }}</p>
          <span v-if="todayRecipe" class="sage-badge"><i></i>{{ activeMeal }} · {{ recipeDifficulty(todayRecipe) }}</span>
        </div>
      </section>

      <nav class="meal-tabs" aria-label="餐段切换">
        <button v-for="tab in mealTabs" :key="tab" :class="{ active: activeMeal === tab }" type="button" @click="activeMeal = tab">{{ tab }}</button>
      </nav>

      <section v-if="menuRecipe" class="glass-card recipe-note" @click="goToRecipe(menuRecipe)">
        <div class="binder" aria-hidden="true"><span v-for="n in 6" :key="n"></span></div>
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
          <p>{{ recipeDescription(menuRecipe) }}</p>
          <button class="primary-pill note-btn" type="button" @click.stop="goToRecipe(menuRecipe)">查看做法 <svg viewBox="0 0 24 24"><path d="m9 18 6-6-6-6"/></svg></button>
        </div>
      </section>

      <section v-else class="glass-card empty-card">
        <h2>今日推荐菜单</h2>
        <p>{{ loading ? '正在读取后端推荐...' : emptyText }}</p>
      </section>

      <section class="glass-card quick-grid" aria-label="快捷功能">
        <button v-for="item in quickActions" :key="item.label" class="quick-item" type="button" @click="goQuickAction(item)">
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
  query?: Record<string, string>
}

const router = useRouter()

const keyword = ref('')
const activeMeal = ref('晚餐')
const loading = ref(true)
const todayRecommend = ref<HomeRecipe | null>(null)
const hotRecipes = ref<HomeRecipe[]>([])

const mealTabs = ['早餐', '午餐', '晚餐', '夜宵']
const emptyText = '后端暂无推荐菜谱，请先在后台维护菜谱数据。'

const homeVars = computed(() => ({
  '--home-bg': 'url(' + kitchenBg + ')',
  '--couple-img': 'url(' + coupleImage + ')',
}))

const todayRecipe = computed(() => todayRecommend.value)
const menuRecipe = computed(() => hotRecipes.value[0] || todayRecommend.value || null)
const popularRecipes = computed(() => hotRecipes.value.slice(0, 2))

const quickActions: QuickAction[] = [
  { label: '按食材推荐', icon: 'ingredients', tone: 'tomato', path: '/recommend/ingredients' },
  { label: '按口味推荐', icon: 'taste', tone: 'pepper', path: '/recommend/taste' },
  { label: '按场景推荐', icon: 'scene', tone: 'sage', path: '/recommend/scene' },
  { label: '一周菜单', icon: 'week', tone: 'wheat', path: '/week-menu' },
  { label: '冰箱剩菜', icon: 'fridge', tone: 'sage', path: '/recommend/fridge' },
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
  return recipe.description || '后端暂无描述，点击查看详细做法。'
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
  router.push(item.query ? { path: item.path, query: item.query } : item.path)
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
  background: radial-gradient(circle at 16% 9%, rgba(255, 255, 255, 0.46), transparent 30%), linear-gradient(180deg, rgba(255, 247, 235, 0.2), rgba(96, 61, 31, 0.08));
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
  border: 1px solid rgba(255, 255, 255, 0.78);
  background: rgba(255, 250, 241, 0.88);
  color: #593a2e;
  box-shadow: 0 12px 24px rgba(91, 59, 31, 0.14), inset 0 1px 0 rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  cursor: pointer;
  transition: transform 180ms ease, opacity 180ms ease, box-shadow 180ms ease;
}

.calendar-btn {
  width: 56px;
  height: 56px;
  flex: 0 0 56px;
  border-radius: 19px;
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
  height: 56px;
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 26px;
  padding: 0 20px;
  border: 1px solid rgba(255, 255, 255, 0.66);
  border-radius: 28px;
  background: rgba(255, 250, 241, 0.78);
  box-shadow: 0 18px 36px rgba(98, 68, 42, 0.12), inset 0 1px 0 rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(18px) saturate(1.08);
  -webkit-backdrop-filter: blur(18px) saturate(1.08);
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

.glass-card {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.64);
  border-radius: 28px;
  background: rgba(255, 250, 240, 0.76);
  box-shadow: 0 18px 40px rgba(98, 68, 42, 0.14), inset 0 1px 0 rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(18px) saturate(1.1);
  -webkit-backdrop-filter: blur(18px) saturate(1.1);
}

.couple-card {
  height: 166px;
  margin-top: 26px;
  cursor: pointer;
}

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
  padding: 25px 0 0 29px;
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

.couple-copy p {
  margin: 9px 0 10px;
  color: #735f53;
  font-size: 15px;
  line-height: 1.38;
}

.primary-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  min-height: 38px;
  padding: 0 22px;
  border: 0;
  border-radius: 999px;
  background: linear-gradient(135deg, var(--coral), var(--coral-2));
  color: #fff;
  box-shadow: 0 12px 22px rgba(218, 64, 50, 0.28);
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
  min-height: 260px;
  margin-top: 14px;
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
  background: linear-gradient(90deg, rgba(255, 250, 240, 0.95) 0%, rgba(255, 250, 240, 0.82) 43%, rgba(255, 250, 240, 0.16) 75%), linear-gradient(180deg, rgba(255, 255, 255, 0.1), rgba(246, 225, 198, 0.12));
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
  width: 36px;
  height: 46px;
  border-color: transparent;
  background: transparent;
  color: var(--coral);
  box-shadow: none;
}

.chef-copy {
  position: relative;
  z-index: 2;
  max-width: 58%;
  padding: 24px 0 0 27px;
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

.sage-badge {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  margin-top: 18px;
  padding: 6px 12px;
  border-radius: 999px;
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
  height: 62px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  align-items: center;
  gap: 4px;
  margin-top: 14px;
  padding: 6px;
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 28px;
  background: rgba(255, 249, 239, 0.78);
  box-shadow: 0 16px 34px rgba(96, 66, 40, 0.12), inset 0 1px 0 rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
}

.meal-tabs button {
  height: 48px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: #4a352b;
  font-size: 17px;
  font-weight: 800;
  cursor: pointer;
  transition: color 180ms ease, background 180ms ease, box-shadow 180ms ease, transform 180ms ease;
}

.meal-tabs button.active {
  background: rgba(255, 252, 246, 0.96);
  color: var(--coral);
  box-shadow: 0 10px 24px rgba(91, 60, 33, 0.15), inset 0 1px 0 rgba(255, 255, 255, 0.96);
}

.recipe-note {
  min-height: 164px;
  display: grid;
  grid-template-columns: 43% 1fr;
  gap: 16px;
  margin-top: 14px;
  padding: 15px 16px 15px 30px;
  background: repeating-linear-gradient(0deg, rgba(169, 122, 83, 0.12) 0 1px, transparent 1px 28px), rgba(255, 248, 235, 0.9);
  cursor: pointer;
}

.binder {
  position: absolute;
  left: -5px;
  top: 16px;
  bottom: 16px;
  width: 24px;
  display: grid;
  align-content: space-around;
  z-index: 3;
}

.binder span {
  width: 24px;
  height: 7px;
  border: 2px solid rgba(169, 132, 95, 0.82);
  border-left: 0;
  border-radius: 0 999px 999px 0;
  background: rgba(255, 255, 255, 0.62);
}

.note-photo {
  min-height: 134px;
  border-radius: 18px;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  box-shadow: 0 12px 26px rgba(80, 48, 24, 0.16);
}

.no-cover span {
  font-size: 13px;
  font-weight: 800;
}

.note-copy {
  position: relative;
  min-width: 0;
  padding: 1px 0 0;
}

.heart-btn {
  position: absolute;
  top: -5px;
  right: -4px;
  width: 34px;
  height: 34px;
  border: 0;
  background: transparent;
  color: #df6a42;
  box-shadow: none;
}

.note-copy h2 {
  margin: 0 28px 10px 0;
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
  min-height: 38px;
  float: right;
  padding: 0 18px;
  font-size: 14px;
}

.empty-card {
  margin-top: 14px;
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

.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 13px;
  margin-top: 14px;
  padding: 15px;
}

.quick-item {
  min-height: 90px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 20px;
  background: rgba(255, 253, 247, 0.78);
  box-shadow: 0 10px 22px rgba(98, 68, 42, 0.1);
  color: #49372f;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  transition: transform 170ms ease, background 170ms ease, box-shadow 170ms ease;
}

.quick-icon,
.quick-icon svg {
  width: 38px;
  height: 38px;
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
  margin-top: 14px;
  padding: 15px 15px 16px;
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
  border: 0;
  background: transparent;
  color: #a29387;
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
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 11px;
}

.popular-card {
  display: grid;
  grid-template-columns: 52% 1fr;
  min-height: 98px;
  overflow: hidden;
  border-radius: 18px;
  background: rgba(255, 252, 246, 0.82);
  box-shadow: 0 10px 22px rgba(91, 61, 38, 0.11);
  cursor: pointer;
}

.popular-photo {
  min-height: 98px;
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
}

.popular-copy {
  min-width: 0;
  padding: 11px 8px 9px 10px;
}

.popular-copy h3 {
  margin: 0;
  color: #3b2a22;
  font-size: 15px;
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
  font-size: 10px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.popular-copy p {
  margin: 7px 0 0;
  color: #806f64;
  font-size: 10px;
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
.quick-item:active,
.glass-card:active,
.popular-card:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .primary-pill:hover,
  .quick-item:hover,
  .calendar-btn:hover {
    transform: translateY(-1px);
  }

  .quick-item:hover {
    background: rgba(255, 255, 255, 0.86);
    box-shadow: 0 14px 28px rgba(91, 61, 38, 0.13);
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

  .chef-copy h2 {
    font-size: 28px;
  }

  .recipe-note {
    grid-template-columns: 41% 1fr;
    padding-right: 13px;
  }

  .quick-grid {
    gap: 10px;
  }
}

@media (max-width: 360px) {
  .popular-card {
    grid-template-columns: 1fr;
  }

  .popular-photo {
    min-height: 86px;
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
