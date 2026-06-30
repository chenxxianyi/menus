<template>
  <main class="favorites-shell" :style="pageVars">
    <header class="favorites-header" aria-label="页面顶部">
      <button class="back-btn" type="button" aria-label="返回" @click="router.back()">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <h1 class="page-title">我的收藏</h1>
      <p class="favorite-count">已收藏 <strong>{{ displayTotal }}</strong> 道菜</p>
    </header>

    <nav class="category-shell" aria-label="收藏分类">
      <button
        v-for="category in categories"
        :key="category"
        class="category-tab"
        :class="{ active: activeCategory === category }"
        type="button"
        @click="activeCategory = category"
      >
        {{ category }}
      </button>
    </nav>

    <section v-if="loading" class="state-card" aria-label="加载中">
      <div class="loading-spinner"></div>
    </section>

    <section v-else-if="filteredFavorites.length" class="recipe-list" aria-label="收藏菜谱列表">
      <article
        v-for="recipe in filteredFavorites"
        :key="recipe.id"
        class="recipe-card"
        @click="openRecipe(recipe)"
      >
        <div class="food-photo" :style="photoStyle(recipe)" role="img" :aria-label="recipeTitle(recipe)">
          <span v-if="!recipe.cover" class="photo-fallback">{{ recipeTag(recipe) }}</span>
        </div>

        <div class="recipe-body">
          <button
            class="favorite-btn active"
            type="button"
            :aria-label="'取消收藏 ' + recipeTitle(recipe)"
            @click.stop="toggleRecipeFavorite(recipe)"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" />
            </svg>
          </button>

          <h2 class="recipe-title">{{ recipeTitle(recipe) }}</h2>
          <p class="recipe-desc">{{ recipeDescription(recipe) }}</p>

          <div class="recipe-meta">
            <span>
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <circle cx="12" cy="12" r="8.5" />
                <path d="M12 7v5l3 2" />
              </svg>
              {{ recipeTime(recipe) }}
            </span>
            <i aria-hidden="true"></i>
            <span>
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M4 20V9" />
                <path d="M10 20V4" />
                <path d="M16 20v-8" />
                <path d="M22 20H2" />
              </svg>
              {{ recipeDifficulty(recipe) }}
            </span>
          </div>

          <span class="recipe-tag">{{ recipeTag(recipe) }}</span>
        </div>
      </article>
    </section>

    <section v-else class="empty-card">
      <div class="empty-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24">
          <path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" />
        </svg>
      </div>
      <h2>还没有收藏的菜谱</h2>
      <p>去菜谱库发现更多适合你的美味</p>
      <button class="browse-btn" type="button" @click="router.push('/recipes')">去浏览菜谱</button>
    </section>

    <p class="sync-tip">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
        <path d="m9 12 2 2 4-5" />
      </svg>
      <span>收藏的菜谱将同步到你的账号</span>
    </p>

    <div class="toast" :class="{ show: !!toastText }">{{ toastText }}</div>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getUserFavorites } from '@/api/user'
import { removeFavorite, toggleFavorite } from '@/api/recipe'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'

type CategoryName = '全部' | '最近收藏' | '下饭菜' | '快手菜' | '家常菜'

interface FavoriteRecipe {
  id: number
  title?: string
  cover?: string
  description?: string
  cook_time?: number
  difficulty?: string
  category_name?: string
  taste?: string
  favorite_count?: number
  is_favorited?: boolean
  created_at?: string
}

const router = useRouter()
const categories: CategoryName[] = ['全部', '最近收藏', '下饭菜', '快手菜', '家常菜']

const favorites = ref<FavoriteRecipe[]>([])
const activeCategory = ref<CategoryName>('全部')
const loading = ref(true)
const total = ref(0)
const toastText = ref('')

let toastTimer: ReturnType<typeof setTimeout> | null = null

const pageVars = computed(() => ({
  '--favorites-bg': `url(${kitchenBg})`,
}))

const displayTotal = computed(() => total.value || favorites.value.length || 0)

const filteredFavorites = computed(() => {
  if (activeCategory.value === '全部') return favorites.value
  if (activeCategory.value === '最近收藏') return favorites.value.slice(0, 6)
  return favorites.value.filter((recipe) => recipeCategories(recipe).includes(activeCategory.value))
})

function normalizeList(payload: any) {
  if (Array.isArray(payload)) return { list: payload, total: payload.length }
  const list = payload?.list || payload?.items || payload?.data || []
  const count = Number(payload?.total ?? payload?.count ?? list.length)
  return { list: Array.isArray(list) ? list : [], total: count }
}

async function fetchFavorites() {
  loading.value = true
  try {
    const res: any = await getUserFavorites(1, 50)
    const normalized = normalizeList(res)
    favorites.value = normalized.list.map((recipe: FavoriteRecipe) => ({
      ...recipe,
      is_favorited: true,
    }))
    total.value = normalized.total
  } catch {
    favorites.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function showToast(message: string) {
  toastText.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastText.value = ''
  }, 1400)
}

function recipeTitle(recipe: FavoriteRecipe) {
  return recipe.title || '未命名菜谱'
}

function recipeDescription(recipe: FavoriteRecipe) {
  return recipe.description || recipe.category_name || recipe.taste || '后端暂无描述，点击查看详细做法。'
}

function recipeTime(recipe: FavoriteRecipe) {
  return recipe.cook_time ? `${recipe.cook_time} 分钟` : '时间待补充'
}

function recipeDifficulty(recipe: FavoriteRecipe) {
  return recipe.difficulty || '难度待补充'
}

function recipeTag(recipe: FavoriteRecipe) {
  return recipe.category_name || recipe.taste || recipe.difficulty || '家常菜'
}

function recipeCategories(recipe: FavoriteRecipe): CategoryName[] {
  const text = `${recipe.category_name || ''}${recipe.taste || ''}${recipe.title || ''}${recipe.description || ''}`
  const result: CategoryName[] = []
  if (text.includes('下饭')) result.push('下饭菜')
  if (text.includes('快手') || (recipe.cook_time || 999) <= 25) result.push('快手菜')
  if (text.includes('家常') || !result.length) result.push('家常菜')
  return result
}

function photoStyle(recipe: FavoriteRecipe) {
  if (recipe.cover) {
    return { backgroundImage: `url(${recipe.cover})` }
  }
  return {
    backgroundImage: fallbackBackground(recipe),
  }
}

function fallbackBackground(recipe: FavoriteRecipe) {
  const tag = recipeTag(recipe)
  if (tag.includes('下饭')) {
    return 'radial-gradient(circle at 58% 45%, #d84f2f 0 10%, transparent 11%), radial-gradient(circle at 43% 50%, #f0a23c 0 12%, transparent 13%), radial-gradient(circle at 36% 40%, #8f3f25 0 8%, transparent 9%), linear-gradient(135deg, #f4c378, #a85a2c)'
  }
  if (tag.includes('快手')) {
    return 'radial-gradient(ellipse at 50% 58%, #d9c2a3 0 28%, transparent 29%), radial-gradient(ellipse at 56% 53%, #f6f1da 0 18%, transparent 19%), linear-gradient(135deg, #e8c590, #b97a43)'
  }
  return 'radial-gradient(circle at 50% 48%, #f2bf35 0 22%, transparent 23%), radial-gradient(circle at 62% 38%, #e6533e 0 18%, transparent 19%), linear-gradient(135deg, #ffe0a6, #d86a34)'
}

function openRecipe(recipe: FavoriteRecipe) {
  router.push(`/recipes/${recipe.id}`)
}

async function toggleRecipeFavorite(recipe: FavoriteRecipe) {
  const previous = recipe.is_favorited !== false
  recipe.is_favorited = !previous
  favorites.value = favorites.value.filter((item) => item.id !== recipe.id)
  total.value = Math.max(0, total.value - 1)
  showToast('已取消收藏')

  try {
    if (previous) {
      await removeFavorite(recipe.id)
    } else {
      await toggleFavorite(recipe.id)
    }
  } catch {
    recipe.is_favorited = previous
    favorites.value = [recipe, ...favorites.value]
    total.value += 1
    showToast('操作失败，已恢复')
  }
}

onMounted(() => {
  fetchFavorites()
})
</script>

<style scoped>
.favorites-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.8);
  --coral: #e95645;
  --coral-2: #ef5548;
  --coral-soft: #fce7de;
  --sage: #8fa783;
  --border: rgba(255, 255, 255, 0.62);
  position: relative;
  width: min(100%, 430px);
  min-height: calc(100vh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 160px);
  min-height: calc(100dvh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 160px);
  margin: 0 auto;
  padding: max(52px, env(safe-area-inset-top)) 20px calc(36px + env(safe-area-inset-bottom));
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 237, 205, 0.34), rgba(255, 247, 233, 0.22)),
    var(--favorites-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.favorites-shell::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 22% 8%, rgba(255, 255, 255, 0.7), transparent 31%),
    radial-gradient(circle at 86% 16%, rgba(225, 126, 55, 0.22), transparent 33%),
    radial-gradient(circle at 12% 93%, rgba(233, 86, 69, 0.18), transparent 29%),
    linear-gradient(90deg, rgba(255, 239, 214, 0.55), rgba(255, 245, 230, 0.16) 55%, rgba(172, 91, 33, 0.18));
  backdrop-filter: blur(4px) saturate(1.12);
  -webkit-backdrop-filter: blur(4px) saturate(1.12);
}

.favorites-header,
.category-shell,
.recipe-list,
.state-card,
.empty-card,
.sync-tip {
  position: relative;
  z-index: 1;
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

.favorites-header {
  min-height: 88px;
  display: grid;
  place-items: center;
  align-content: center;
  text-align: center;
}

.back-btn {
  position: absolute;
  left: 0;
  top: 0;
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 16px;
  color: #4a352a;
  background: rgba(255, 250, 240, 0.86);
  box-shadow:
    0 12px 28px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.back-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.55;
}

.page-title {
  margin: 0;
  color: #2e241f;
  font-size: 28px;
  font-weight: 900;
  line-height: 1;
  letter-spacing: 0;
}

.favorite-count {
  margin: 15px 0 0;
  color: var(--sub);
  font-size: 16px;
  font-weight: 650;
  line-height: 1;
}

.favorite-count strong {
  color: var(--coral);
  font-size: 19px;
  font-weight: 950;
}

.category-shell {
  height: 66px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2px;
  margin-top: 34px;
  padding: 6px 9px;
  overflow-x: auto;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 999px;
  background: rgba(255, 250, 240, 0.78);
  box-shadow:
    0 14px 30px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.76);
  backdrop-filter: blur(18px) saturate(1.08);
  -webkit-backdrop-filter: blur(18px) saturate(1.08);
  scrollbar-width: none;
}

.category-shell::-webkit-scrollbar {
  display: none;
}

.category-tab {
  position: relative;
  height: 52px;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  padding: 0 11px;
  border-radius: 999px;
  color: #3a2a24;
  background: transparent;
  font-size: 15.5px;
  font-weight: 820;
  white-space: nowrap;
  transition: color 180ms ease, transform 180ms ease;
}

.category-tab.active {
  color: var(--coral);
  font-weight: 920;
}

.category-tab.active::after {
  content: "";
  position: absolute;
  left: 50%;
  bottom: 5px;
  width: 24px;
  height: 3px;
  border-radius: 999px;
  background: var(--coral);
  transform: translateX(-50%);
}

.recipe-list {
  display: grid;
  gap: 18px;
  margin-top: 24px;
}

.recipe-card {
  position: relative;
  min-height: 198px;
  display: grid;
  grid-template-columns: 142px minmax(0, 1fr);
  gap: 16px;
  padding: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 30px;
  background:
    radial-gradient(circle at 5% 0%, rgba(255, 255, 255, 0.42), transparent 42%),
    rgba(255, 250, 240, 0.8);
  box-shadow:
    0 18px 42px rgba(80, 50, 30, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.82);
  cursor: pointer;
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
  transition: transform 180ms ease, opacity 180ms ease, box-shadow 180ms ease;
}

.food-photo {
  width: 142px;
  min-height: 174px;
  display: grid;
  place-items: center;
  border-radius: 22px;
  background-color: #efd3ac;
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
  box-shadow:
    0 12px 24px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.photo-fallback {
  width: 74px;
  height: 74px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  color: rgba(255, 255, 255, 0.92);
  background: rgba(233, 86, 69, 0.32);
  font-size: 15px;
  font-weight: 850;
}

.recipe-body {
  position: relative;
  min-width: 0;
  padding: 9px 2px 6px 0;
}

.favorite-btn {
  position: absolute;
  top: 14px;
  right: 14px;
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 16px;
  color: var(--coral);
  background: rgba(255, 250, 240, 0.78);
  box-shadow: 0 10px 20px rgba(80, 50, 30, 0.08);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  transition: transform 180ms ease, opacity 180ms ease, color 180ms ease;
}

.favorite-btn svg {
  width: 25px;
  height: 25px;
  stroke-width: 2;
}

.favorite-btn.active svg {
  fill: currentColor;
}

.recipe-title {
  max-width: calc(100% - 50px);
  margin: 0;
  color: var(--text);
  font-size: 22px;
  font-weight: 950;
  line-height: 1.2;
  letter-spacing: 0;
}

.recipe-desc {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  margin: 12px 0 0;
  color: var(--sub);
  font-size: 13.5px;
  font-weight: 560;
  line-height: 1.62;
}

.recipe-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 9px;
  margin-top: 17px;
  color: #7d6d61;
  font-size: 13.5px;
  font-weight: 680;
  white-space: nowrap;
}

.recipe-meta span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.recipe-meta svg {
  width: 18px;
  height: 18px;
  stroke-width: 2.2;
}

.recipe-meta i {
  width: 1px;
  height: 17px;
  background: rgba(122, 106, 95, 0.3);
}

.recipe-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-top: 18px;
  padding: 7px 15px;
  border-radius: 999px;
  color: var(--coral);
  background: var(--coral-soft);
  font-size: 14px;
  font-weight: 820;
  line-height: 1;
}

.sync-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 26px 0 4px;
  color: var(--sub);
  font-size: 15px;
  font-weight: 650;
  text-align: center;
}

.sync-tip svg {
  width: 23px;
  height: 23px;
  flex: 0 0 23px;
  color: var(--sage);
  stroke-width: 2.2;
}

.state-card,
.empty-card {
  margin-top: 24px;
  padding: 42px 24px 34px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 30px;
  background: rgba(255, 250, 240, 0.82);
  box-shadow:
    0 18px 42px rgba(80, 50, 30, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.82);
  text-align: center;
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.state-card {
  min-height: 220px;
  display: grid;
  place-items: center;
}

.loading-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid rgba(233, 86, 69, 0.18);
  border-top-color: var(--coral);
  border-radius: 50%;
  animation: spin 0.72s linear infinite;
}

.empty-icon {
  width: 72px;
  height: 72px;
  display: grid;
  place-items: center;
  margin: 0 auto 20px;
  border-radius: 24px;
  color: var(--coral);
  background: rgba(252, 231, 222, 0.82);
}

.empty-icon svg {
  width: 40px;
  height: 40px;
  stroke-width: 2.2;
}

.empty-card h2 {
  margin: 0;
  color: var(--text);
  font-size: 23px;
  font-weight: 920;
}

.empty-card p {
  margin: 11px 0 24px;
  color: var(--sub);
  font-size: 15px;
  line-height: 1.55;
}

.browse-btn {
  width: 100%;
  height: 54px;
  border-radius: 18px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 14px 26px rgba(233, 86, 69, 0.25);
  font-size: 17px;
  font-weight: 860;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: calc(28px + env(safe-area-inset-bottom));
  z-index: 5;
  padding: 10px 16px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 999px;
  color: #fff;
  background: rgba(46, 36, 31, 0.78);
  box-shadow: 0 12px 24px rgba(46, 36, 31, 0.18);
  font-size: 14px;
  font-weight: 740;
  opacity: 0;
  pointer-events: none;
  transform: translate(-50%, 12px);
  transition: opacity 180ms ease, transform 180ms ease;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
}

.toast.show {
  opacity: 1;
  transform: translate(-50%, 0);
}

.back-btn:active,
.category-tab:active,
.recipe-card:active,
.favorite-btn:active,
.browse-btn:active {
  transform: scale(0.98);
}

.recipe-card:active {
  opacity: 0.9;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (hover: hover) {
  .back-btn:hover,
  .category-tab:hover,
  .favorite-btn:hover,
  .browse-btn:hover {
    transform: translateY(-1px);
  }

  .recipe-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 22px 46px rgba(80, 50, 30, 0.18);
  }
}

@media (max-width: 380px) {
  .favorites-shell {
    padding-left: 18px;
    padding-right: 18px;
  }

  .category-tab {
    padding: 0 9px;
    font-size: 14.5px;
  }

  .recipe-card {
    grid-template-columns: 136px minmax(0, 1fr);
    gap: 13px;
    padding: 12px;
  }

  .food-photo {
    width: 136px;
    min-height: 170px;
  }

  .recipe-title {
    font-size: 20px;
  }

  .recipe-desc {
    font-size: 13px;
    line-height: 1.58;
  }

  .recipe-meta {
    gap: 7px;
    font-size: 13px;
  }
}

@media (max-width: 350px) {
  .recipe-card {
    grid-template-columns: 1fr;
  }

  .food-photo {
    width: 100%;
    min-height: 180px;
  }

  .recipe-body {
    padding: 2px 46px 6px 2px;
  }
}

@media (min-width: 431px) {
  .favorites-shell {
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.18);
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
