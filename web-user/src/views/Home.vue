<template>
  <div class="page home-page">
    <header class="home-header anim-delay-1">
      <div class="home-header-top">
        <span class="home-brand-kicker">Menu Journal</span>
        <button class="search-btn" aria-label="搜索">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
        </button>
      </div>
      <div class="home-brand-row">
        <div class="avatar">{{ userInitial }}</div>
        <div class="home-brand-copy">
          <p class="home-brand-eyebrow">今天想吃什么</p>
          <h1 class="home-brand-title">温柔一点，慢一点</h1>
        </div>
      </div>
    </header>

    <section class="hero-card anim-delay-2" :class="{ clickable: !!todayRecommend }" @click="openHero">
      <div class="hero-inner">
        <div class="hero-image-wrap">
          <img
            :src="heroImage"
            :alt="heroTitle"
            class="hero-image"
          />
        </div>
        <div class="hero-content">
          <h2 class="hero-title">{{ heroTitle }}</h2>
          <span class="hero-time">{{ heroTime }}</span>
        </div>
      </div>
      <div class="hero-bottom">
        <span class="hero-footer">{{ heroDifficulty }} · {{ heroPeople }} · {{ heroTaste }}</span>
        <div class="hero-badge">
          <span class="hero-badge-label">{{ hasRecommend ? '今日推荐' : '今日灵感' }}</span>
          <span class="hero-badge-date">{{ issueDate }}</span>
        </div>
      </div>
    </section>

    <nav class="meal-nav anim-delay-3" aria-label="餐次切换">
      <button
        v-for="meal in meals"
        :key="meal.key"
        class="meal-pill"
        :class="{ active: activeMeal === meal.key }"
        @click="activeMeal = meal.key"
      >
        <span class="meal-pill-label">{{ meal.label }}</span>
      </button>
    </nav>

    <section class="journal-section anim-delay-3">
      <div class="section-heading-row">
        <h2 class="section-heading">热门菜谱</h2>
        <span class="section-kicker">Curated</span>
      </div>

      <template v-if="hotRecipes.length">
        <div class="recipe-scroll">
          <article
            v-for="(recipe, idx) in hotRecipes.slice(0, 6)"
            :key="recipe.id"
            class="recipe-scroll-card"
            @click="goToRecipe(recipe.id)"
          >
            <div class="recipe-scroll-cover">
              <img :src="recipe.cover || heroArt" :alt="recipe.title" />
            </div>
            <h3 class="recipe-scroll-title">{{ recipe.title }}</h3>
            <p class="recipe-scroll-meta">
              <span class="recipe-scroll-diff">{{ recipe.difficulty }}</span>
              <span class="recipe-scroll-time">{{ recipe.cook_time }}分钟</span>
            </p>
          </article>
        </div>
      </template>
      <div v-else class="issue-empty">
        <p class="issue-empty-title">还没有上架的推荐。</p>
        <p class="issue-empty-desc">先去一周菜单里生成一个更完整的今天。</p>
      </div>
    </section>

    <section class="tools-section anim-delay-4">
      <div class="section-heading-row">
        <h2 class="section-heading">快速入口</h2>
        <span class="section-kicker">Daily tools</span>
      </div>
      <div class="tools-list">
        <button class="tool-row" @click="$router.push('/week-menu')">
          <span class="tool-index">01</span>
          <span class="tool-copy">
            <span class="tool-title">AI 智能推荐</span>
            <span class="tool-sub">告诉我你有什么食材</span>
          </span>
          <svg class="tool-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </button>
        <button class="tool-row" @click="$router.push('/week-menu')">
          <span class="tool-index">02</span>
          <span class="tool-copy">
            <span class="tool-title">一周菜单</span>
            <span class="tool-sub">自动生成 7 天食谱</span>
          </span>
          <svg class="tool-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </button>
        <button class="tool-row" @click="$router.push('/shopping-list')">
          <span class="tool-index">03</span>
          <span class="tool-copy">
            <span class="tool-title">购物清单</span>
            <span class="tool-sub">根据菜单自动生成</span>
          </span>
          <svg class="tool-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getHomeData } from '@/api/home'
import heroArt from '@/assets/hero.png'

const router = useRouter()
const userStore = useUserStore()
const activeMeal = ref('lunch')
const todayRecommend = ref<any>(null)
const hotRecipes = ref<any[]>([])

const userName = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '用户')
const userInitial = computed(() => userName.value.charAt(0))

const meals = [
  { key: 'breakfast', label: '早餐' },
  { key: 'lunch', label: '午餐' },
  { key: 'dinner', label: '晚餐' },
  { key: 'snack', label: '夜宵' },
]

const hasRecommend = computed(() => !!todayRecommend.value)
const heroImage = computed(() => todayRecommend.value?.cover || heroArt)
const heroTitle = computed(() => todayRecommend.value?.title || '为今天留一份温柔的晚餐')
const heroTime = computed(() => todayRecommend.value ? `${todayRecommend.value.cook_time} min` : '20 min')
const heroDifficulty = computed(() => todayRecommend.value?.difficulty || '简单')
const heroPeople = computed(() => todayRecommend.value?.people_count ? `${todayRecommend.value.people_count} 人份` : '2 人份')
const heroTaste = computed(() => todayRecommend.value?.taste || '温和')
const issueDate = computed(() =>
  new Intl.DateTimeFormat('zh-CN', { month: 'long', day: 'numeric' }).format(new Date())
)

function goToRecipe(id: number) {
  router.push(`/recipes/${id}`)
}

function openHero() {
  if (todayRecommend.value?.id) goToRecipe(todayRecommend.value.id)
}

onMounted(async () => {
  try {
    const res: any = await getHomeData()
    todayRecommend.value = res.today_recommend || null
    hotRecipes.value = res.hot_recipes || []
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.home-page {
  padding-top: var(--sp-1);
}

.home-header {
  padding-bottom: var(--sp-4);
}

.home-header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-3);
}

.home-brand-kicker {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.search-btn {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-rule);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-ink-2);
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  transition: background var(--dur-fast) var(--ease-out), transform var(--dur-fast) var(--ease-out);
}

.search-btn:active {
  background: var(--color-paper-3);
  transform: translateY(1px);
}

.search-btn svg {
  width: 18px;
  height: 18px;
}

.home-brand-row {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--color-ink);
  color: var(--color-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 17px;
  font-weight: 700;
  flex-shrink: 0;
}

.home-brand-copy {
  min-width: 0;
}

.home-brand-eyebrow {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 500;
}

.home-brand-title {
  color: var(--color-ink);
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 650;
  line-height: 1.05;
}

.hero-card {
  margin-bottom: var(--sp-6);
  border-radius: var(--r-xl);
  background: var(--color-paper-2);
  overflow: hidden;
}

.hero-card.clickable {
  cursor: pointer;
}

.hero-inner {
  display: flex;
  gap: var(--sp-4);
  padding: var(--sp-4);
}

.hero-image-wrap {
  width: 44%;
  min-height: 140px;
  border-radius: var(--r-lg);
  overflow: hidden;
  flex-shrink: 0;
  background: var(--color-paper-3);
}

.hero-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.hero-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: var(--sp-2);
  min-width: 0;
}

.hero-title {
  color: var(--color-ink);
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 650;
  line-height: 1.15;
  overflow-wrap: anywhere;
}

.hero-time {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
}

.hero-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
}

.hero-footer {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
}

.hero-badge {
  display: inline-flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  padding: var(--sp-2) var(--sp-3);
  border-radius: var(--r-sm);
  background: color-mix(in oklch, var(--color-surface) 88%, transparent);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--color-rule);
}

.hero-badge-label {
  color: var(--color-ink);
  font-size: var(--text-xs);
  font-weight: 700;
}

.hero-badge-date {
  color: var(--color-ink-3);
  font-size: 11px;
}

.section-heading-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-kicker {
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 600;
}

.meal-nav {
  display: flex;
  gap: var(--sp-3);
  margin-bottom: var(--sp-8);
}

.meal-pill {
  flex: 1;
  min-height: 42px;
  border: 1px solid var(--color-rule);
  border-radius: var(--r-full);
  background: transparent;
  color: var(--color-ink-3);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease-out), color var(--dur-fast) var(--ease-out), border-color var(--dur-fast) var(--ease-out), transform var(--dur-fast) var(--ease-out);
}

.meal-pill.active {
  background: var(--color-ink);
  border-color: var(--color-ink);
  color: var(--color-surface);
}

.meal-pill:active {
  transform: translateY(1px);
}

.meal-pill-label {
  font-size: var(--text-xs);
  font-weight: 600;
}

.journal-section,
.tools-section {
  margin-bottom: var(--sp-8);
}

.recipe-scroll {
  display: flex;
  gap: var(--sp-4);
  overflow-x: auto;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;
  padding-bottom: var(--sp-2);
  margin: 0 calc(-1 * var(--sp-6));
  padding-left: var(--sp-6);
  padding-right: var(--sp-6);
}

.recipe-scroll::-webkit-scrollbar {
  display: none;
}

.recipe-scroll-card {
  flex: 0 0 140px;
  scroll-snap-align: start;
  cursor: pointer;
}

.recipe-scroll-cover {
  width: 100%;
  aspect-ratio: 1;
  border-radius: var(--r-lg);
  overflow: hidden;
  background: var(--color-paper-2);
}

.recipe-scroll-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.recipe-scroll-title {
  margin-top: var(--sp-2);
  color: var(--color-ink);
  font-size: var(--text-sm);
  font-weight: 600;
  line-height: 1.3;
  overflow-wrap: anywhere;
}

.recipe-scroll-meta {
  margin-top: var(--sp-1);
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  color: var(--color-ink-3);
  font-size: var(--text-xs);
}

.recipe-scroll-diff {
  padding: 1px 6px;
  border-radius: var(--r-sm);
  background: var(--color-orange-soft);
  color: var(--color-orange);
  font-weight: 600;
  font-size: 10px;
}

.issue-empty {
  padding: var(--sp-8) 0;
  border-top: 1px solid var(--color-rule);
  border-bottom: 1px solid var(--color-rule);
  text-align: center;
}

.issue-empty-title {
  color: var(--color-ink);
  font-size: var(--text-base);
  font-weight: 600;
}

.issue-empty-desc {
  margin-top: var(--sp-1);
  color: var(--color-ink-3);
  font-size: var(--text-sm);
}

.tools-list {
  border-top: 1px solid var(--color-rule);
  border-bottom: 1px solid var(--color-rule);
}

.tool-row {
  width: 100%;
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  padding: var(--sp-4) 0;
  border: 0;
  background: transparent;
  color: inherit;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

.tool-row + .tool-row {
  border-top: 1px solid var(--color-rule);
}

.tool-index {
  width: 38px;
  color: var(--color-ink-3);
  font-size: var(--text-xs);
  font-weight: 700;
  letter-spacing: 0.08em;
  flex-shrink: 0;
}

.tool-copy {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.tool-title {
  color: var(--color-ink);
  font-size: var(--text-base);
  font-weight: 600;
}

.tool-sub {
  margin-top: 2px;
  color: var(--color-ink-3);
  font-size: var(--text-xs);
}

.tool-arrow {
  width: 18px;
  height: 18px;
  color: var(--color-ink-3);
  flex-shrink: 0;
}

@media (max-width: 420px) {
  .hero-inner {
    flex-direction: column;
  }

  .hero-image-wrap {
    width: 100%;
    min-height: 160px;
  }
}
</style>
