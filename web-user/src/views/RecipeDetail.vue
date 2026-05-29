<template>
  <div class="recipe-page">
    <!-- Hero -->
    <div class="hero anim-delay-1">
      <button class="hero-back btn-ghost" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <img :src="recipe?.cover || heroArt" :alt="recipe?.title || 'recipe'" class="hero-image" />
      <div class="hero-wash"></div>
      <span v-if="recipe" class="hero-diff chip">{{ recipe.difficulty }}</span>
    </div>

    <template v-if="recipe">
      <div class="page-content">
        <!-- Title -->
        <div class="title-block anim-delay-2">
          <p class="title-eyebrow">Editorial recipe</p>
          <h1 class="recipe-title">{{ recipe.title }}</h1>
          <p class="recipe-desc">{{ recipe.description }}</p>
        </div>

        <!-- Stats -->
        <div class="stats-grid anim-delay-2">
          <div class="stat-cell">
            <span class="stat-label">时间</span>
            <span class="stat-value">{{ recipe.cook_time }} 分钟</span>
          </div>
          <div class="stat-cell">
            <span class="stat-label">份量</span>
            <span class="stat-value">{{ recipe.people_count }} 人份</span>
          </div>
          <div class="stat-cell">
            <span class="stat-label">风味</span>
            <span class="stat-value">{{ recipe.taste }}</span>
          </div>
          <div class="stat-cell">
            <span class="stat-label">热量</span>
            <span class="stat-value">{{ recipe.nutrition?.calories ?? '--' }}</span>
          </div>
        </div>

        <!-- Nutrition -->
        <div v-if="recipe.nutrition" class="block anim-delay-3">
          <h2 class="block-title">营养成分</h2>
          <div class="nutri-list">
            <div class="nutri-row">
              <span class="nutri-label">蛋白质</span>
              <div class="nutri-track"><div class="nutri-fill" :style="{ width: nutriPct(recipe.nutrition.protein, 60) + '%' }"></div></div>
              <span class="nutri-val">{{ recipe.nutrition.protein }}g</span>
            </div>
            <div class="nutri-row">
              <span class="nutri-label">脂肪</span>
              <div class="nutri-track"><div class="nutri-fill nutri-fill--muted" :style="{ width: nutriPct(recipe.nutrition.fat, 60) + '%' }"></div></div>
              <span class="nutri-val">{{ recipe.nutrition.fat }}g</span>
            </div>
            <div class="nutri-row">
              <span class="nutri-label">碳水</span>
              <div class="nutri-track"><div class="nutri-fill nutri-fill--warm" :style="{ width: nutriPct(recipe.nutrition.carbs, 80) + '%' }"></div></div>
              <span class="nutri-val">{{ recipe.nutrition.carbs }}g</span>
            </div>
            <div class="nutri-row">
              <span class="nutri-label">纤维</span>
              <div class="nutri-track"><div class="nutri-fill nutri-fill--light" :style="{ width: nutriPct(recipe.nutrition.fiber, 30) + '%' }"></div></div>
              <span class="nutri-val">{{ recipe.nutrition.fiber }}g</span>
            </div>
          </div>
        </div>

        <!-- Ingredients -->
        <div v-if="ingredients.length" class="block anim-delay-3">
          <h2 class="block-title">食材清单</h2>
          <div class="list-group">
            <div
              v-for="ing in ingredients"
              :key="ing.name"
              class="list-row ing-row"
              @click="ing.checked = !ing.checked"
            >
              <div class="ing-check" :class="{ checked: ing.checked }">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" width="12" height="12">
                  <path d="M20 6L9 17l-5-5" />
                </svg>
              </div>
              <div class="list-row-body">
                <div class="list-row-title" :class="{ done: ing.checked }">{{ ing.name }}</div>
                <div class="list-row-sub">{{ ing.amount }} {{ ing.unit }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Steps -->
        <div v-if="steps.length" class="block anim-delay-4">
          <h2 class="block-title">烹饪步骤</h2>
          <div v-for="(step, idx) in steps" :key="idx" class="step">
            <div class="step-num">{{ idx + 1 }}</div>
            <div class="step-body">
              <p class="step-desc">{{ step.description }}</p>
              <p v-if="step.tip" class="step-tip">{{ step.tip }}</p>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="actions anim-delay-5">
          <button class="btn-solid action-primary">加入今日菜单</button>
          <button class="btn-outline" :class="{ liked: recipe.is_favorited }" @click="handleFavorite">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
            </svg>
          </button>
        </div>
      </div>
    </template>

    <div v-else class="empty-state anim-delay-2">加载中...</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getRecipeDetail, toggleFavorite, removeFavorite } from '@/api/recipe'
import heroArt from '@/assets/hero.png'

const route = useRoute()
const recipe = ref<any>(null)

const ingredients = computed(() => {
  const raw = recipe.value?.ingredients
  if (!raw) return []
  return Array.isArray(raw) ? raw.map((i: any) => ({ ...i, checked: false })) : []
})

const steps = computed(() => {
  const raw = recipe.value?.steps
  if (!raw) return []
  return Array.isArray(raw) ? raw : []
})

function nutriPct(val: number, max: number) {
  return Math.min(100, Math.round((val / max) * 100))
}

async function handleFavorite() {
  if (!recipe.value) return
  try {
    if (recipe.value.is_favorited) {
      await removeFavorite(recipe.value.id)
      recipe.value.is_favorited = false
    } else {
      await toggleFavorite(recipe.value.id)
      recipe.value.is_favorited = true
    }
  } catch {
    // ignore
  }
}

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) return
  try {
    const res: any = await getRecipeDetail(id)
    recipe.value = res
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.recipe-page {
  padding-bottom: var(--sp-8);
}

/* ── Hero ── */
.hero {
  position: relative;
  min-height: 300px;
  margin: 0 calc(var(--sp-5) * -1) var(--sp-5);
  overflow: hidden;
}

.hero-image {
  width: 100%;
  height: 100%;
  min-height: 300px;
  object-fit: cover;
  display: block;
}

.hero-wash {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, transparent 40%, rgba(0, 0, 0, 0.25) 100%);
}

.hero-back {
  position: absolute;
  left: var(--sp-4);
  top: var(--safe-top);
  z-index: 2;
}

.hero-diff {
  position: absolute;
  left: var(--sp-4);
  bottom: var(--sp-4);
  z-index: 2;
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-border);
  color: var(--color-text);
}

.page-content {
  padding: 0 var(--sp-5);
}

/* ── Title ── */
.title-block {
  margin-bottom: var(--sp-5);
}

.title-eyebrow {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
  margin-bottom: var(--sp-1);
}

.recipe-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 650;
  line-height: 1.08;
  letter-spacing: -0.01em;
}

.recipe-desc {
  margin-top: var(--sp-2);
  color: var(--color-text-3);
  font-size: var(--text-sm);
  line-height: 1.65;
}

/* ── Stats Grid ── */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  overflow: hidden;
  margin-bottom: var(--sp-6);
}

.stat-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--sp-4);
  background: var(--color-surface);
}

.stat-cell:nth-child(2),
.stat-cell:nth-child(4) {
  border-left: 1px solid var(--color-border);
}

.stat-cell:nth-child(3),
.stat-cell:nth-child(4) {
  border-top: 1px solid var(--color-border);
}

.stat-label {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
}

.stat-value {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-md);
  font-weight: 650;
}

/* ── Blocks ── */
.block {
  margin-bottom: var(--sp-8);
}

.block-title {
  margin-bottom: var(--sp-4);
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 650;
}

/* ── Nutrition ── */
.nutri-list {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3);
}

.nutri-row {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.nutri-label {
  width: 56px;
  flex-shrink: 0;
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
}

.nutri-track {
  flex: 1;
  height: 5px;
  overflow: hidden;
  border-radius: 3px;
  background: var(--color-surface-3);
}

.nutri-fill {
  height: 100%;
  border-radius: 3px;
  background: var(--color-text);
  transition: width var(--dur-slow) var(--ease-out);
}

.nutri-fill--muted { background: var(--color-text-3); }
.nutri-fill--warm { background: var(--color-accent); }
.nutri-fill--light { background: var(--color-surface-3); border: 1px solid var(--color-border); }

.nutri-val {
  width: 36px;
  flex-shrink: 0;
  color: var(--color-text-2);
  font-size: var(--text-xs);
  font-weight: 700;
  text-align: right;
}

/* ── Ingredients ── */
.ing-row {
  cursor: pointer;
}

.ing-check {
  width: 22px;
  height: 22px;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border: 1.5px solid var(--color-border-med);
  border-radius: 50%;
  transition: all var(--dur-base) var(--ease);
}

.ing-check.checked {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
}

.ing-check svg { display: none; }
.ing-check.checked svg { display: block; }

.done {
  color: var(--color-text-3) !important;
  text-decoration: line-through;
}

/* ── Steps ── */
.step {
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr);
  gap: var(--sp-3);
  padding: var(--sp-4) 0;
  border-top: 1px solid var(--color-border);
}

.step-num {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--color-surface-2);
  color: var(--color-text-2);
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

.step-body {
  min-width: 0;
}

.step-desc {
  color: var(--color-text-2);
  font-size: var(--text-sm);
  line-height: 1.7;
}

.step-tip {
  margin-top: var(--sp-2);
  padding: var(--sp-2) var(--sp-3);
  border-radius: var(--r-sm);
  background: var(--color-warning-soft);
  color: var(--color-text-2);
  font-size: var(--text-xs);
  line-height: 1.5;
}

/* ── Actions ── */
.actions {
  display: flex;
  gap: var(--sp-3);
  padding-bottom: var(--sp-8);
}

.action-primary {
  flex: 1;
}

.liked {
  border-color: var(--color-error);
  color: var(--color-error);
}

/* ── Empty ── */
.empty-state {
  padding: var(--sp-12) 0;
  text-align: center;
  color: var(--color-text-3);
  font-size: var(--text-sm);
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .page-content {
    max-width: 640px;
    margin: 0 auto;
  }
}

@media (max-width: 420px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-cell:nth-child(2),
  .stat-cell:nth-child(4) {
    border-left: 0;
  }

  .stat-cell:nth-child(2) {
    border-top: 1px solid var(--color-border);
  }
}
</style>
