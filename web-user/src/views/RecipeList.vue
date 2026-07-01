<template>
  <div class="recipe-shell" :style="pageVars">
    <div class="recipe-warm-overlay" aria-hidden="true"></div>

    <main class="recipe-phone">
      <header class="recipe-header">
        <p class="recipe-eyebrow">Recipe library</p>
        <h1 class="recipe-title">{{ pageTitle }}</h1>

        <form class="recipe-search" role="search" @submit.prevent="handleSearch">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="10.8" cy="10.8" r="7.2" />
            <path d="m16 16 4.2 4.2" />
          </svg>
          <input
            v-model.trim="keyword"
            type="search"
            placeholder="搜索菜名、食材或口味"
            aria-label="搜索菜名、食材或口味"
          />
          <button v-if="keyword" class="clear-btn" type="button" aria-label="清空搜索" @click="clearSearch">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6 6 18" />
              <path d="m6 6 12 12" />
            </svg>
          </button>
        </form>
      </header>

      <div v-if="activeTaste || activeSort === 'hot'" class="active-filters" aria-label="当前筛选">
        <span v-if="activeTaste">口味：{{ activeTaste }}</span>
        <span v-if="activeSort === 'hot'">最热排序</span>
        <button type="button" @click="clearRouteFilters">清除筛选</button>
      </div>

      <nav class="sort-tabs" aria-label="排序方式">
        <button type="button" :class="{ active: activeSort === 'latest' }" @click="selectSort('latest')">最新发布</button>
        <button type="button" :class="{ active: activeSort === 'hot' }" @click="selectSort('hot')">最热菜谱</button>
      </nav>

      <nav class="recipe-chips" aria-label="菜谱分类">
        <button
          class="recipe-chip"
          :class="{ active: activeCategory === undefined }"
          type="button"
          @click="selectCategory(undefined)"
        >
          全部
        </button>
        <button
          v-for="category in categories"
          :key="category.id"
          class="recipe-chip"
          :class="{ active: activeCategory === category.id }"
          type="button"
          @click="selectCategory(category.id)"
        >
          {{ category.name }}
        </button>
      </nav>

      <section v-if="loading" class="recipe-state" aria-label="加载中">
        <div class="loading-spinner"></div>
      </section>

      <section v-else-if="recipes.length" class="recipe-grid" aria-label="菜谱列表">
        <article
          v-for="recipe in recipes"
          :key="recipe.id"
          class="recipe-card"
          @click="openRecipe(recipe)"
        >
          <div class="recipe-card-media">
            <img v-if="recipe.cover" :src="recipe.cover" :alt="recipeTitle(recipe)" loading="lazy" />
            <div v-else class="no-cover" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 18h16" />
                <path d="M6 18v-6a6 6 0 0 1 12 0v6" />
                <path d="M9 10c.8-1.2 1.8-1.8 3-1.8s2.2.6 3 1.8" />
              </svg>
              <span>暂无图片</span>
            </div>
            <button
              class="favorite-btn"
              :class="{ active: recipe.is_favorited }"
              type="button"
              :aria-label="recipe.is_favorited ? '取消收藏菜谱' : '收藏菜谱'"
              @click.stop="handleFavorite(recipe)"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" />
              </svg>
            </button>
          </div>

          <div class="recipe-card-body">
            <h2>{{ recipeTitle(recipe) }}</h2>
            <p>{{ recipeDescription(recipe) }}</p>
            <div class="recipe-meta">
              <span>
                <svg class="clock-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <circle cx="12" cy="12" r="8.5" />
                  <path d="M12 7v5l3 2" />
                </svg>
                {{ recipeTime(recipe) }}
              </span>
              <i aria-hidden="true"></i>
              <span>
                <svg class="level-icon" :class="{ orange: isMediumRecipe(recipe) }" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <path d="M4 20V9M10 20V4M16 20v-8M22 20H2" />
                </svg>
                {{ recipeTag(recipe) }}
              </span>
            </div>
          </div>
        </article>
      </section>

      <section v-else class="recipe-empty">
        <h2>没有找到符合条件的菜谱</h2>
        <p>可以换个关键词，或在后台维护更多菜谱数据。</p>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCategories, getRecipes, removeFavorite, toggleFavorite } from '@/api/recipe'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'

interface RecipeListItem {
  id: number
  title?: string
  cover?: string
  description?: string
  cook_time?: number
  difficulty?: string
  favorite_count?: number
  is_favorited?: boolean
  category_id?: number
  category_name?: string
  taste?: string
}

interface CategoryItem {
  id: number
  name: string
}

const router = useRouter()
const route = useRoute()
const keyword = ref('')
const activeCategory = ref<number | undefined>(undefined)
const activeTaste = ref('')
const activeSort = ref<'latest' | 'hot'>('latest')
const categories = ref<CategoryItem[]>([])
const recipes = ref<RecipeListItem[]>([])
const loading = ref(true)
let syncingFromRoute = false

const pageVars = computed(() => ({
  '--recipe-bg': `url(${kitchenBg})`,
}))

const pageTitle = computed(() => {
  if (activeSort.value === 'hot') return '最热菜谱'
  if (activeTaste.value) return `${activeTaste.value}菜谱`
  return '菜谱'
})

function normalizeList(payload: any) {
  if (Array.isArray(payload)) return payload
  return payload?.list || payload?.items || []
}

async function fetchRecipes() {
  loading.value = true
  try {
    const res: any = await getRecipes({
      keyword: keyword.value || undefined,
      category_id: activeCategory.value,
      taste: activeTaste.value || undefined,
      sort: activeSort.value,
      page: 1,
      page_size: 24,
    })
    recipes.value = normalizeList(res)
  } catch {
    recipes.value = []
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const res: any = await getCategories()
    categories.value = normalizeList(res)
  } catch {
    categories.value = []
  }
}

function handleSearch() {
  syncQueryAndFetch()
}

function clearSearch() {
  keyword.value = ''
  syncQueryAndFetch()
}

function selectCategory(categoryId?: number) {
  activeCategory.value = categoryId
  syncQueryAndFetch()
}

function selectSort(sort: 'latest' | 'hot') {
  activeSort.value = sort
  syncQueryAndFetch()
}

function clearRouteFilters() {
  activeTaste.value = ''
  activeSort.value = 'latest'
  syncQueryAndFetch()
}

function firstQueryValue(value: unknown) {
  return Array.isArray(value) ? String(value[0] || '') : String(value || '')
}

function applyRouteQuery() {
  syncingFromRoute = true
  keyword.value = firstQueryValue(route.query.keyword)
  activeTaste.value = firstQueryValue(route.query.taste)
  const sort = firstQueryValue(route.query.sort)
  activeSort.value = sort === 'hot' ? 'hot' : 'latest'
  const category = Number(firstQueryValue(route.query.category_id))
  activeCategory.value = Number.isFinite(category) && category > 0 ? category : undefined
  syncingFromRoute = false
}

function buildQuery() {
  const query: Record<string, string> = {}
  if (keyword.value) query.keyword = keyword.value
  if (activeCategory.value) query.category_id = String(activeCategory.value)
  if (activeTaste.value) query.taste = activeTaste.value
  if (activeSort.value !== 'latest') query.sort = activeSort.value
  return query
}

function syncQueryAndFetch() {
  if (!syncingFromRoute) {
    router.replace({ path: '/recipes', query: buildQuery() })
  }
  fetchRecipes()
}

function openRecipe(recipe: RecipeListItem) {
  router.push(`/recipes/${recipe.id}`)
}

async function handleFavorite(recipe: RecipeListItem) {
  const wasFavorited = recipe.is_favorited === true
  const previousCount = recipe.favorite_count || 0
  recipe.is_favorited = !wasFavorited
  recipe.favorite_count = Math.max(0, previousCount + (wasFavorited ? -1 : 1))

  try {
    if (wasFavorited) {
      await removeFavorite(recipe.id)
    } else {
      await toggleFavorite(recipe.id)
    }
  } catch {
    recipe.is_favorited = wasFavorited
    recipe.favorite_count = previousCount
  }
}

function recipeTitle(recipe: RecipeListItem) {
  return recipe.title || '未命名菜谱'
}

function recipeDescription(recipe: RecipeListItem) {
  return recipe.description || recipe.category_name || recipe.taste || '后端暂无描述'
}

function recipeTime(recipe: RecipeListItem) {
  return recipe.cook_time ? `${recipe.cook_time} 分钟` : '时间待补充'
}

function recipeTag(recipe: RecipeListItem) {
  return recipe.category_name || recipe.difficulty || recipe.taste || '待补充'
}

function isMediumRecipe(recipe: RecipeListItem) {
  const value = `${recipe.difficulty || ''}${recipe.category_name || ''}${recipe.taste || ''}`
  return value.includes('中') || value.includes('难') || value.includes('热')
}

onMounted(() => {
  applyRouteQuery()
  fetchCategories()
  fetchRecipes()
})

watch(
  () => route.query,
  () => {
    applyRouteQuery()
    fetchRecipes()
  }
)
</script>

<style scoped>
.recipe-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.82);
  --coral: #e95645;
  --orange: #e89a45;
  --sage: #79a35d;
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

.recipe-warm-overlay {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 20% 8%, rgba(255, 255, 255, 0.66), transparent 33%),
    linear-gradient(90deg, rgba(255, 246, 232, 0.44), rgba(255, 246, 232, 0.1) 55%, rgba(151, 92, 45, 0.08));
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.recipe-phone {
  position: relative;
  z-index: 1;
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: max(47px, env(safe-area-inset-top)) 21px calc(var(--tab-h, 82px) + 70px);
}

.recipe-header {
  margin-bottom: 24px;
}

.active-filters,
.sort-tabs {
  display: flex;
  gap: 10px;
  margin-top: 16px;
  overflow-x: auto;
  scrollbar-width: none;
}

.active-filters::-webkit-scrollbar,
.sort-tabs::-webkit-scrollbar {
  display: none;
}

.active-filters span,
.active-filters button,
.sort-tabs button {
  min-height: 42px;
  flex: 0 0 auto;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 999px;
  background: rgba(255, 250, 240, 0.82);
  color: #3a2a24;
  box-shadow: 0 10px 20px rgba(98, 68, 42, 0.09);
  font-size: 14px;
  font-weight: 850;
}

.active-filters span,
.active-filters button {
  display: inline-flex;
  align-items: center;
  padding: 0 14px;
}

.active-filters button,
.sort-tabs button {
  cursor: pointer;
}

.active-filters button {
  color: var(--coral);
}

.sort-tabs button {
  padding: 0 18px;
}

.sort-tabs button.active {
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
}

.recipe-eyebrow {
  margin: 0 0 8px;
  color: var(--orange);
  font-size: 17px;
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: 0;
}

.recipe-title {
  margin: 0;
  color: #30221d;
  font-size: clamp(46px, 12vw, 56px);
  font-weight: 950;
  line-height: 0.96;
  letter-spacing: 0;
}

.recipe-search {
  height: 66px;
  display: flex;
  align-items: center;
  gap: 15px;
  margin-top: 27px;
  padding: 0 18px 0 22px;
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: 999px;
  background: rgba(255, 251, 242, 0.77);
  box-shadow:
    0 14px 30px rgba(90, 60, 35, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(20px) saturate(1.1);
  -webkit-backdrop-filter: blur(20px) saturate(1.1);
}

.recipe-search > svg {
  width: 31px;
  height: 31px;
  flex: 0 0 auto;
  color: #5f4c43;
  stroke-width: 2.4;
}

.recipe-search input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: 0;
  color: #49372f;
  background: transparent;
  font-size: 18px;
  font-weight: 650;
}

.recipe-search input::placeholder {
  color: #8c7b70;
  opacity: 1;
}

.clear-btn {
  width: 38px;
  height: 38px;
  display: grid;
  flex: 0 0 38px;
  place-items: center;
  border: 0;
  border-radius: 50%;
  color: #806f64;
  background: rgba(255, 255, 255, 0.56);
  cursor: pointer;
  transition: color 180ms ease, transform 180ms ease, background 180ms ease;
}

.clear-btn svg {
  width: 18px;
  height: 18px;
  stroke-width: 2.2;
}

.recipe-chips {
  display: flex;
  gap: 14px;
  margin: 32px -21px 0;
  padding: 0 21px 3px;
  overflow-x: auto;
  scrollbar-width: none;
  -webkit-overflow-scrolling: touch;
}

.recipe-chips::-webkit-scrollbar {
  display: none;
}

.recipe-chip {
  flex: 0 0 auto;
  height: 48px;
  padding: 0 24px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 999px;
  color: #3a2a24;
  background: rgba(255, 250, 240, 0.82);
  box-shadow:
    0 12px 24px rgba(98, 68, 42, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.82);
  font-size: 17px;
  font-weight: 850;
  cursor: pointer;
  white-space: nowrap;
  transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease, color 180ms ease;
}

.recipe-chip.active {
  border-color: rgba(255, 255, 255, 0.45);
  color: #fff;
  background: linear-gradient(135deg, #ff7568, var(--coral));
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.26);
}

.recipe-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px 14px;
  margin-top: 29px;
}

.recipe-card {
  min-width: 0;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.6);
  border-radius: 24px;
  background: var(--cream);
  box-shadow: 0 18px 36px rgba(80, 50, 28, 0.15);
  cursor: pointer;
  backdrop-filter: blur(18px) saturate(1.08);
  -webkit-backdrop-filter: blur(18px) saturate(1.08);
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.recipe-card-media {
  position: relative;
  height: clamp(132px, 34vw, 152px);
  overflow: hidden;
  background: #f4dfc5;
}

.recipe-card-media::after {
  content: "";
  position: absolute;
  inset: auto 0 0;
  height: 50%;
  background: linear-gradient(180deg, transparent, rgba(92, 54, 25, 0.08));
}

.recipe-card-media img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.no-cover {
  height: 100%;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 7px;
  color: #9b887a;
  background:
    radial-gradient(circle at 34% 24%, rgba(255, 255, 255, 0.74), transparent 38%),
    linear-gradient(135deg, #fbecd8, #e9caa7);
}

.no-cover svg {
  width: 38px;
  height: 38px;
  stroke-width: 1.8;
}

.no-cover span {
  font-size: 13px;
  font-weight: 800;
}

.favorite-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 2;
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.74);
  border-radius: 15px;
  color: #8b7b70;
  background: rgba(255, 250, 240, 0.88);
  box-shadow: 0 10px 20px rgba(71, 48, 31, 0.12);
  cursor: pointer;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  transition: transform 180ms ease, color 180ms ease, background 180ms ease;
}

.favorite-btn svg {
  width: 25px;
  height: 25px;
  fill: transparent;
  stroke-width: 2.1;
}

.favorite-btn.active {
  color: #fff;
  background: linear-gradient(135deg, #ff7669, var(--coral));
}

.favorite-btn.active svg {
  fill: currentColor;
}

.recipe-card-body {
  min-height: 126px;
  display: flex;
  flex-direction: column;
  padding: 14px 16px 16px;
  background: linear-gradient(180deg, rgba(255, 250, 239, 0.78), rgba(255, 249, 237, 0.9));
}

.recipe-card h2 {
  margin: 0;
  color: #33241e;
  font-size: clamp(20px, 5.1vw, 23px);
  font-weight: 950;
  line-height: 1.14;
  letter-spacing: 0;
}

.recipe-card p {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  min-height: 43px;
  margin: 9px 0 12px;
  color: var(--sub);
  font-size: 15px;
  font-weight: 560;
  line-height: 1.42;
}

.recipe-meta {
  display: grid;
  grid-template-columns: 1fr 1px 1fr;
  align-items: center;
  gap: 11px;
  margin-top: auto;
  color: #554238;
  font-size: 14.5px;
  font-weight: 760;
  white-space: nowrap;
}

.recipe-meta span {
  min-width: 0;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.recipe-meta i {
  width: 1px;
  height: 18px;
  background: rgba(141, 112, 92, 0.18);
}

.recipe-meta svg {
  width: 16px;
  height: 16px;
  flex: 0 0 16px;
  stroke-width: 2.35;
}

.clock-icon {
  color: var(--coral);
}

.level-icon {
  color: var(--sage);
}

.level-icon.orange {
  color: var(--orange);
}

.recipe-state,
.recipe-empty {
  min-height: 260px;
  display: grid;
  place-items: center;
  margin-top: 29px;
  padding: 26px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 26px;
  background: rgba(255, 250, 240, 0.78);
  box-shadow: 0 18px 36px rgba(80, 50, 28, 0.12);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
}

.recipe-empty {
  align-content: center;
  text-align: center;
}

.recipe-empty h2 {
  margin: 0 0 8px;
  color: #33241e;
  font-size: 21px;
  font-weight: 900;
}

.recipe-empty p {
  margin: 0;
  color: #7a6a5f;
  font-size: 14px;
  line-height: 1.5;
}

.loading-spinner {
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

.clear-btn:active,
.recipe-chip:active,
.favorite-btn:active,
.recipe-card:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .recipe-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 22px 42px rgba(80, 50, 28, 0.18);
  }

  .recipe-chip:hover,
  .favorite-btn:hover,
  .clear-btn:hover {
    transform: translateY(-1px);
  }

  .clear-btn:hover {
    background: rgba(255, 255, 255, 0.72);
  }
}

@media (max-width: 380px) {
  .recipe-phone {
    padding-left: 18px;
    padding-right: 18px;
  }

  .recipe-chips {
    margin-left: -18px;
    margin-right: -18px;
    padding-left: 18px;
    padding-right: 18px;
  }

  .recipe-grid {
    gap: 16px 12px;
  }

  .recipe-card-body {
    padding-left: 13px;
    padding-right: 13px;
  }

  .recipe-meta {
    gap: 8px;
    font-size: 13.5px;
  }
}

@media (max-width: 350px) {
  .recipe-grid {
    grid-template-columns: 1fr;
  }

  .recipe-card-media {
    height: 176px;
  }
}

@media (min-width: 431px) {
  .recipe-shell {
    background-color: #ead7bd;
  }

  .recipe-phone {
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
