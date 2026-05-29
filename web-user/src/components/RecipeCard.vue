<template>
  <div class="recipe-card glass-card" @click="$emit('click')">
    <div class="recipe-card-cover">
      <img :src="cover" :alt="title" />
      <span class="recipe-card-time">{{ cookTime }}分钟</span>
    </div>
    <div class="recipe-card-info">
      <h3 class="recipe-card-title">{{ title }}</h3>
      <div class="recipe-card-meta">
        <span class="recipe-card-difficulty">{{ difficulty }}</span>
        <span class="recipe-card-fav">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor" stroke="none">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
          {{ favoriteCount }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  id: number
  title: string
  cover: string
  cookTime: number
  difficulty: string
  favoriteCount: number
}>()

defineEmits<{
  click: []
}>()
</script>

<style scoped>
.recipe-card {
  padding: 0;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid var(--color-rule);
  border-radius: var(--r-lg);
  background: var(--color-surface);
  box-shadow: 0 8px 28px var(--color-shadow);
  transition: transform var(--dur-base) var(--ease-out);
}
.recipe-card:active { transform: translateY(1px); }
.recipe-card-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 16/10;
  overflow: hidden;
  border-radius: var(--r-lg) var(--r-lg) 0 0;
  background: var(--color-paper-2);
}
.recipe-card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.recipe-card-time {
  position: absolute;
  bottom: 8px;
  right: 8px;
  padding: 4px 10px;
  border-radius: var(--r-sm);
  background: var(--color-ink);
  backdrop-filter: blur(8px);
  color: var(--color-surface);
  font-size: 11px;
  font-weight: 600;
}
.recipe-card-info { padding: 12px 14px; }
.recipe-card-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--color-ink);
  margin-bottom: 6px;
  letter-spacing: 0;
}
.recipe-card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-ink-3);
}
.recipe-card-difficulty {
  padding: 2px 8px;
  border-radius: var(--r-sm);
  background: var(--color-orange-soft);
  color: var(--color-orange);
  font-weight: 600;
  font-size: 11px;
}
.recipe-card-fav {
  display: flex;
  align-items: center;
  gap: 3px;
  color: var(--color-red);
}
</style>
