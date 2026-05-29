<template>
  <div class="page recipe-page">
    <div class="hero anim-delay-1">
      <button class="hero-back btn-ghost" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <img :src="recipe?.cover || heroArt" :alt="recipe?.title || 'recipe cover'" class="hero-image" />
      <div class="hero-wash"></div>
      <span v-if="recipe" class="hero-badge chip chip--orange">{{ recipe.difficulty }}</span>
    </div>

    <template v-if="recipe">
      <div class="title-block anim-delay-2">
        <p class="title-eyebrow">Editorial recipe</p>
        <h1 class="recipe-title">{{ recipe.title }}</h1>
        <p class="recipe-sub">{{ recipe.description }}</p>
      </div>

      <div class="inline-stats anim-delay-2 surface">
        <div class="inline-stat">
          <span class="inline-stat-label">时间</span>
          <span class="inline-stat-value">{{ recipe.cook_time }} 分钟</span>
        </div>
        <div class="inline-stat">
          <span class="inline-stat-label">份量</span>
          <span class="inline-stat-value">{{ recipe.people_count }} 人份</span>
        </div>
        <div class="inline-stat">
          <span class="inline-stat-label">风味</span>
          <span class="inline-stat-value">{{ recipe.taste }}</span>
        </div>
        <div class="inline-stat">
          <span class="inline-stat-label">热量</span>
          <span class="inline-stat-value">{{ recipe.nutrition?.calories ?? '--' }}</span>
        </div>
      </div>

      <div v-if="recipe.nutrition" class="nutrition anim-delay-3">
        <h2 class="section-heading">营养成分</h2>
        <div class="nutri-bars">
          <div class="nutri-row">
            <span class="nutri-label">蛋白质</span>
            <div class="nutri-track">
              <div class="nutri-fill" :style="{ width: nutriPct(recipe.nutrition.protein, 60) + '%' }"></div>
            </div>
            <span class="nutri-val">{{ recipe.nutrition.protein }}g</span>
          </div>
          <div class="nutri-row">
            <span class="nutri-label">脂肪</span>
            <div class="nutri-track">
              <div class="nutri-fill nutri-fill--blue" :style="{ width: nutriPct(recipe.nutrition.fat, 60) + '%' }"></div>
            </div>
            <span class="nutri-val">{{ recipe.nutrition.fat }}g</span>
          </div>
          <div class="nutri-row">
            <span class="nutri-label">碳水</span>
            <div class="nutri-track">
              <div class="nutri-fill nutri-fill--green" :style="{ width: nutriPct(recipe.nutrition.carbs, 80) + '%' }"></div>
            </div>
            <span class="nutri-val">{{ recipe.nutrition.carbs }}g</span>
          </div>
          <div class="nutri-row">
            <span class="nutri-label">纤维</span>
            <div class="nutri-track">
              <div class="nutri-fill nutri-fill--purple" :style="{ width: nutriPct(recipe.nutrition.fiber, 30) + '%' }"></div>
            </div>
            <span class="nutri-val">{{ recipe.nutrition.fiber }}g</span>
          </div>
        </div>
      </div>

      <div v-if="ingredients.length" class="ingredients anim-delay-3">
        <h2 class="section-heading">食材清单</h2>
        <div class="list-group">
          <div
            v-for="ing in ingredients"
            :key="ing.name"
            class="list-row ingredient-row"
            @click="ing.checked = !ing.checked"
          >
            <div class="ing-check" :class="{ checked: ing.checked }">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" width="12" height="12">
                <path d="M20 6L9 17l-5-5" />
              </svg>
            </div>
            <div class="list-row-body">
              <div class="list-row-title" :class="{ 'is-done': ing.checked }">{{ ing.name }}</div>
              <div class="list-row-sub">{{ ing.amount }} {{ ing.unit }}</div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="steps.length" class="steps anim-delay-4">
        <h2 class="section-heading">烹饪步骤</h2>
        <div v-for="(step, idx) in steps" :key="idx" class="step">
          <div class="step-num">{{ idx + 1 }}</div>
          <div class="step-body">
            <h3 class="step-title">第 {{ idx + 1 }} 步</h3>
            <p class="step-desc">{{ step.description }}</p>
            <p v-if="step.tip" class="step-tip">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                <circle cx="12" cy="12" r="10" />
                <path d="M12 16v-4M12 8h.01" />
              </svg>
              {{ step.tip }}
            </p>
          </div>
        </div>
      </div>

      <div class="actions anim-delay-5">
        <button class="btn-solid action-primary">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14" />
          </svg>
          加入今日菜单
        </button>
        <button class="btn-outline" :class="{ liked: recipe.is_favorited }" @click="handleFavorite">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
          </svg>
        </button>
      </div>
    </template>

    <div v-else class="empty-hint anim-delay-2">加载中...</div>
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

.hero {
  position: relative;
  min-height: 320px;
  margin: 0 calc(var(--sp-6) * -1) var(--sp-5);
  overflow: hidden;
}

.hero-image {
  width: 100%;
  height: 100%;
  min-height: 320px;
  object-fit: cover;
  display: block;
}

.hero-wash {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, transparent 25%, var(--color-photo-wash) 100%);
}

.hero-back {
  position: absolute;
  left: var(--sp-4);
  top: var(--sp-4);
  z-index: 2;
}

.hero-badge {
  position: absolute;
  left: var(--sp-4);
  bottom: var(--sp-4);
  z-index: 2;
}

.title-block {
  margin-bottom: var(--sp-4);
}

.title-eyebrow {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 600;
  margin-bottom: var(--sp-1);
}

.recipe-title {
  color: var(--color-ink);
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 650;
  line-height: 1.05;
  overflow-wrap: anywhere;
}

.recipe-sub {
  margin-top: var(--sp-2);
  color: var(--color-ink-3);
  font-size: var(--text-sm);
  line-height: 1.65;
}

.inline-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1px;
  padding: 1px;
  margin-bottom: var(--sp-6);
  background: var(--color-rule);
}

.inline-stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: var(--sp-4);
  background: var(--color-surface);
}

.inline-stat-label {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 600;
}

.inline-stat-value {
  color: var(--color-ink);
  font-family: var(--font-display);
  font-size: var(--text-md);
  font-weight: 650;
}

.nutrition,
.ingredients,
.steps {
  margin-bottom: var(--sp-8);
}

.nutri-bars {
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
  width: 60px;
  flex-shrink: 0;
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 600;
}

.nutri-track {
  flex: 1;
  height: 6px;
  overflow: hidden;
  border-radius: var(--r-full);
  background: var(--color-paper-3);
}

.nutri-fill {
  height: 100%;
  border-radius: var(--r-full);
  background: var(--color-orange);
  transition: width var(--dur-slow) var(--ease-out);
}

.nutri-fill--blue { background: var(--color-blue); }
.nutri-fill--green { background: var(--color-green); }
.nutri-fill--purple { background: var(--color-purple); }

.nutri-val {
  width: 36px;
  flex-shrink: 0;
  color: var(--color-ink-2);
  font-size: var(--text-xs);
  font-weight: 700;
  text-align: right;
}

.ingredient-row {
  border-top: 1px solid var(--color-rule);
}

.ing-check {
  width: 22px;
  height: 22px;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border: 1.5px solid var(--color-rule-strong);
  border-radius: 50%;
  transition: background var(--dur-base) var(--ease-out), border-color var(--dur-base) var(--ease-out);
}

.ing-check.checked {
  background: var(--color-ink);
  border-color: var(--color-ink);
  color: var(--color-surface);
}

.ing-check svg {
  display: none;
}

.ing-check.checked svg {
  display: block;
}

.list-row-body {
  flex: 1;
  min-width: 0;
}

.list-row-title {
  color: var(--color-ink);
  font-size: var(--text-base);
  font-weight: 600;
}

.list-row-sub {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  margin-top: 2px;
}

.is-done {
  color: var(--color-ink-3) !important;
  text-decoration: line-through;
}

.step {
  display: grid;
  grid-template-columns: 30px minmax(0, 1fr);
  gap: var(--sp-3);
  padding: var(--sp-4) 0;
  border-top: 1px solid var(--color-rule);
}

.step-num {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--color-ink);
  color: var(--color-surface);
  font-size: 12px;
  font-weight: 700;
}

.step-body {
  min-width: 0;
}

.step-title {
  color: var(--color-ink);
  font-size: var(--text-base);
  font-weight: 650;
  margin-bottom: 4px;
}

.step-desc {
  color: var(--color-ink-2);
  font-size: var(--text-sm);
  line-height: 1.7;
}

.step-tip {
  display: flex;
  align-items: center;
  gap: var(--sp-1);
  margin-top: var(--sp-2);
  padding: var(--sp-2) var(--sp-3);
  border-radius: var(--r-sm);
  background: var(--color-yellow-soft);
  color: var(--color-ink-2);
  font-size: 12px;
  font-weight: 500;
  line-height: 1.45;
}

.actions {
  display: flex;
  gap: var(--sp-3);
}

.action-primary {
  flex: 1;
}

.btn-outline.liked {
  border-color: var(--color-red);
  color: var(--color-red);
}

.empty-hint {
  padding: var(--sp-8) 0;
  color: var(--color-ink-3);
  text-align: center;
  font-size: var(--text-sm);
}

@media (max-width: 420px) {
  .inline-stats {
    grid-template-columns: 1fr;
  }
}
</style>
