<template>
  <div class="page recipes-page">
    <header class="recipes-header anim-delay-1">
      <div>
        <p class="recipes-eyebrow">Recipe library</p>
        <h1 class="recipes-title">菜谱</h1>
      </div>
    </header>

    <form class="search-box anim-delay-2" @submit.prevent="handleSearch">
      <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8" />
        <path d="m21 21-4.35-4.35" />
      </svg>
      <input
        v-model.trim="keyword"
        class="search-input"
        type="search"
        placeholder="搜索菜名、食材或口味"
      />
      <button v-if="keyword" class="clear-btn" type="button" aria-label="清空搜索" @click="clearSearch">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6 6 18" />
          <path d="m6 6 12 12" />
        </svg>
      </button>
    </form>

    <div class="filter-row anim-delay-3" aria-label="菜谱分类">
      <button
        class="filter-chip"
        :class="{ active: activeCategory === undefined }"
        type="button"
        @click="selectCategory(undefined)"
      >
        全部
      </button>
      <button
        v-for="category in categories"
        :key="category.id"
        class="filter-chip"
        :class="{ active: activeCategory === category.id }"
        type="button"
        @click="selectCategory(category.id)"
      >
        {{ category.name }}
      </button>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
    </div>

    <div v-else-if="recipes.length" class="recipe-grid anim-delay-4">
      <RecipeCard
        v-for="recipe in recipes"
        :id="recipe.id"
        :key="recipe.id"
        :title="recipe.title"
        :cover="recipe.cover || heroArt"
        :cook-time="recipe.cook_time || 0"
        :difficulty="recipe.difficulty || '简单'"
        :favorite-count="recipe.favorite_count || 0"
        @click="router.push(`/recipes/${recipe.id}`)"
      />
    </div>

    <EmptyState
      v-else
      class="anim-delay-4"
      type="no-result"
      text="没有找到符合条件的菜谱"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import RecipeCard from '@/components/RecipeCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import { getCategories, getRecipes } from '@/api/recipe'
import heroArt from '@/assets/hero.png'

interface RecipeListItem {
  id: number
  title: string
  cover?: string
  cook_time?: number
  difficulty?: string
  favorite_count?: number
}

interface CategoryItem {
  id: number
  name: string
}

const router = useRouter()
const keyword = ref('')
const activeCategory = ref<number | undefined>(undefined)
const categories = ref<CategoryItem[]>([])
const recipes = ref<RecipeListItem[]>([])
const loading = ref(true)

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
  fetchRecipes()
}

function clearSearch() {
  keyword.value = ''
  fetchRecipes()
}

function selectCategory(categoryId?: number) {
  activeCategory.value = categoryId
  fetchRecipes()
}

onMounted(() => {
  fetchCategories()
  fetchRecipes()
})
</script>

<style scoped>
.recipes-page {
  min-height: 100vh;
  padding-top: var(--sp-5);
  padding-bottom: var(--sp-16);
}

.recipes-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: var(--sp-5);
}

.recipes-eyebrow {
  margin-bottom: var(--sp-1);
  color: var(--color-accent);
  font-size: var(--text-xs);
  font-weight: 700;
}

.recipes-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 700;
  line-height: 1.1;
}

.search-box {
  min-height: 48px;
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  padding: 0 var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.search-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  color: var(--color-text-3);
}

.search-input {
  width: 100%;
  min-width: 0;
  border: 0;
  background: transparent;
  color: var(--color-text);
  font-size: var(--text-sm);
}

.search-input::placeholder {
  color: var(--color-text-3);
}

.search-input:focus {
  outline: none;
}

.clear-btn {
  width: 28px;
  height: 28px;
  display: grid;
  flex-shrink: 0;
  place-items: center;
  border: 0;
  border-radius: var(--r-full);
  background: var(--color-surface-2);
  color: var(--color-text-3);
  cursor: pointer;
  transition: color var(--dur-fast) var(--ease), transform var(--dur-fast) var(--ease);
}

.clear-btn:hover {
  color: var(--color-text-2);
}

.clear-btn:active {
  transform: translateY(1px);
}

.clear-btn svg {
  width: 15px;
  height: 15px;
}

.filter-row {
  display: flex;
  gap: var(--sp-2);
  margin: var(--sp-4) calc(-1 * var(--sp-5)) var(--sp-5);
  padding: 0 var(--sp-5) var(--sp-1);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.filter-row::-webkit-scrollbar {
  display: none;
}

.filter-chip {
  flex-shrink: 0;
  min-height: 34px;
  padding: 0 var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-full);
  background: var(--color-surface);
  color: var(--color-text-2);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease);
}

.filter-chip:hover {
  border-color: var(--color-border-med);
}

.filter-chip:active {
  transform: translateY(1px);
}

.filter-chip.active {
  border-color: var(--color-accent);
  background: var(--color-accent-soft);
  color: var(--color-accent);
}

.recipe-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--sp-4);
}

.loading-state {
  display: flex;
  justify-content: center;
  padding: var(--sp-16) 0;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-text);
  border-radius: var(--r-full);
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (min-width: 768px) {
  .recipes-page {
    max-width: 640px;
    padding-right: var(--sp-8);
    padding-left: var(--sp-8);
  }

  .filter-row {
    margin-right: calc(-1 * var(--sp-8));
    margin-left: calc(-1 * var(--sp-8));
    padding-right: var(--sp-8);
    padding-left: var(--sp-8);
  }
}
</style>
