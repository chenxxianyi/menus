<template>
  <button class="recipe-card" type="button" :aria-label="`查看${title}做法`" @click="$emit('click')">
    <div class="recipe-card-cover">
      <img :src="cover" :alt="title" />
      <span class="recipe-card-time">{{ cookTime }}分钟</span>
    </div>
    <div class="recipe-card-info">
      <h3 class="recipe-card-title">{{ title }}</h3>
      <div class="recipe-card-meta">
        <span class="recipe-card-diff">{{ difficulty }}</span>
        <span class="recipe-card-fav">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor" stroke="none">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
          {{ favoriteCount }}
        </span>
      </div>
    </div>
  </button>
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
  width: 100%;
  display: block;
  padding: 0;
  text-align: left;
  font: inherit;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid var(--color-border);
  border-radius: var(--card-radius-inner);
  background: var(--card-surface-strong);
  transition: border-color var(--dur-fast) var(--ease), transform var(--dur-fast) var(--ease), box-shadow var(--dur-fast) var(--ease);
}

.recipe-card:focus-visible {
  outline: 3px solid color-mix(in srgb, var(--color-accent) 45%, transparent);
  outline-offset: 3px;
}

.recipe-card:hover {
  border-color: var(--color-border-med);
  box-shadow: var(--shadow-sm);
}

.recipe-card:active { transform: translateY(1px); }

.recipe-card-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 16/10;
  overflow: hidden;
  background: var(--color-surface-2);
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
  padding: 3px 8px;
  border-radius: 4px;
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-border);
  color: var(--color-text);
  font-size: 11px;
  font-weight: 600;
}

.recipe-card-info { padding: var(--sp-3); }

.recipe-card-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4px;
}

.recipe-card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: var(--text-xs);
  color: var(--color-text-3);
}

.recipe-card-diff {
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--color-accent-soft);
  color: var(--color-accent);
  font-weight: 600;
  font-size: var(--text-2xs);
}

.recipe-card-fav {
  display: flex;
  align-items: center;
  gap: 3px;
  color: var(--color-text-3);
}
</style>
