<template>
  <div class="page favorites-page">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">我的收藏</h1>
      <div style="width:34px"></div>
    </header>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
    </div>

    <div v-else-if="favorites.length" class="recipe-list">
      <div
        v-for="recipe in favorites"
        :key="recipe.id"
        class="recipe-item"
        @click="$router.push(`/recipes/${recipe.id}`)"
      >
        <div class="recipe-cover">
          <img :src="recipe.cover" :alt="recipe.title" />
        </div>
        <div class="recipe-info">
          <h3 class="recipe-title">{{ recipe.title }}</h3>
          <p class="recipe-meta">{{ recipe.cook_time }}分钟 · {{ recipe.difficulty }}</p>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <div class="empty-icon">&#10084;</div>
      <p class="empty-text">还没有收藏的菜谱</p>
      <button class="primary-btn" @click="$router.push('/')">去浏览菜谱</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUserFavorites } from '@/api/user'

const favorites = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res: any = await getUserFavorites(1, 50)
    favorites.value = res?.list || res || []
  } catch {
    favorites.value = []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.favorites-page {
  min-height: 100vh;
  background: var(--color-bg);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
}

.back-btn {
  width: 34px; height: 34px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.back-btn:hover { background: var(--color-surface-2); }
.back-btn:active { background: var(--color-surface-3); }
.back-btn svg { width: 18px; height: 18px; }

.page-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

.recipe-list {
  padding: 0 var(--sp-4);
}

.recipe-item {
  display: flex;
  gap: var(--sp-3);
  padding: var(--sp-3) 0;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.recipe-item:hover {
  background: var(--color-surface-2);
  margin: 0 calc(-1 * var(--sp-4));
  padding: var(--sp-3) var(--sp-4);
}

.recipe-cover {
  width: 64px;
  height: 64px;
  border-radius: var(--r-sm);
  overflow: hidden;
  flex-shrink: 0;
  background: var(--color-surface-2);
}

.recipe-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.recipe-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-width: 0;
}

.recipe-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 2px;
}

.recipe-meta {
  font-size: var(--text-xs);
  color: var(--color-text-3);
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
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: var(--sp-16) var(--sp-4);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: var(--sp-3);
}

.empty-text {
  color: var(--color-text-3);
  font-size: var(--text-sm);
  margin-bottom: var(--sp-6);
}

.primary-btn {
  padding: var(--sp-3) var(--sp-8);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease),
              transform var(--dur-fast) var(--ease);
}

.primary-btn:hover {
  background: var(--color-text-2);
}

.primary-btn:active {
  transform: translateY(1px);
}

@media (min-width: 768px) {
  .favorites-page { max-width: 640px; margin: 0 auto; }
}
</style>
