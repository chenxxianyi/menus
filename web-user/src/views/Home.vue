<template>
  <div class="page home-page">
    <!-- Header -->
    <header class="home-header anim-delay-1">
      <div class="home-header-top">
        <span class="home-brand-name">Menu Journal</span>
        <button class="search-btn" aria-label="搜索">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
        </button>
      </div>
      <div class="home-greeting">
        <div class="avatar">{{ userInitial }}</div>
        <div>
          <p class="greeting-sub">今天想吃什么</p>
          <h1 class="greeting-title">温柔一点，慢一点</h1>
        </div>
      </div>
    </header>

    <!-- Hero -->
    <section class="hero anim-delay-2" :class="{ clickable: !!todayRecommend }" @click="openHero">
      <div class="hero-inner">
        <div class="hero-image-wrap">
          <img :src="heroImage" :alt="heroTitle" class="hero-image" />
        </div>
        <div class="hero-content">
          <h2 class="hero-title">{{ heroTitle }}</h2>
          <span class="hero-time">{{ heroTime }}</span>
        </div>
      </div>
      <div class="hero-bottom">
        <span class="hero-tags">{{ heroDifficulty }} · {{ heroPeople }} · {{ heroTaste }}</span>
        <div class="hero-badge">
          <span class="hero-badge-text">{{ hasRecommend ? '今日推荐' : '今日灵感' }}</span>
          <span class="hero-badge-date">{{ issueDate }}</span>
        </div>
      </div>
    </section>

    <!-- Meal Nav -->
    <nav class="meal-nav anim-delay-3" aria-label="餐次切换">
      <button
        v-for="meal in meals"
        :key="meal.key"
        class="meal-pill"
        :class="{ active: activeMeal === meal.key }"
        @click="activeMeal = meal.key"
      >
        {{ meal.label }}
      </button>
    </nav>

    <!-- Hot Recipes -->
    <section class="section anim-delay-3">
      <div class="section-top">
        <h2 class="section-heading">热门菜谱</h2>
        <span class="section-more">Curated</span>
      </div>

      <template v-if="hotRecipes.length">
        <div class="recipe-scroll">
          <article
            v-for="recipe in hotRecipes.slice(0, 6)"
            :key="recipe.id"
            class="recipe-card"
            @click="goToRecipe(recipe.id)"
          >
            <div class="recipe-cover">
              <img :src="recipe.cover || heroArt" :alt="recipe.title" />
            </div>
            <h3 class="recipe-name">{{ recipe.title }}</h3>
            <p class="recipe-meta">
              <span class="recipe-diff">{{ recipe.difficulty }}</span>
              <span class="recipe-time">{{ recipe.cook_time }}分钟</span>
            </p>
          </article>
        </div>
      </template>
      <div v-else class="empty-hint">
        <p>还没有上架的推荐。</p>
      </div>
    </section>

    <!-- Quick Access -->
    <section class="section anim-delay-4">
      <div class="section-top">
        <h2 class="section-heading">快速入口</h2>
      </div>
      <div class="quick-list">
        <button class="quick-row" @click="$router.push('/week-menu')">
          <span class="quick-index">01</span>
          <span class="quick-body">
            <span class="quick-title">AI 智能推荐</span>
            <span class="quick-sub">告诉我你有什么食材</span>
          </span>
          <svg class="quick-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </button>
        <button class="quick-row" @click="$router.push('/week-menu')">
          <span class="quick-index">02</span>
          <span class="quick-body">
            <span class="quick-title">一周菜单</span>
            <span class="quick-sub">自动生成 7 天食谱</span>
          </span>
          <svg class="quick-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </button>
        <button class="quick-row" @click="$router.push('/shopping-list')">
          <span class="quick-index">03</span>
          <span class="quick-body">
            <span class="quick-title">购物清单</span>
            <span class="quick-sub">根据菜单自动生成</span>
          </span>
          <svg class="quick-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
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
  padding-top: var(--sp-2);
}

/* ── Header ── */
.home-header {
  margin-bottom: var(--sp-6);
}

.home-header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-4);
}

.home-brand-name {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-3);
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.search-btn {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
  transition: background var(--dur-fast) var(--ease), border-color var(--dur-fast) var(--ease), transform var(--dur-fast) var(--ease);
}

.search-btn:hover {
  border-color: var(--color-border-med);
}

.search-btn:active {
  background: var(--color-surface-2);
  transform: translateY(1px);
}

.search-btn svg {
  width: 18px;
  height: 18px;
}

.home-greeting {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--color-text);
  color: var(--color-text-inv);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 17px;
  font-weight: 700;
  flex-shrink: 0;
}

.greeting-sub {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 500;
}

.greeting-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 650;
  line-height: 1.1;
}

/* ── Hero ── */
.hero {
  margin-bottom: var(--sp-6);
  border-radius: var(--r-lg);
  background: var(--color-surface-2);
  overflow: hidden;
}

.hero.clickable {
  cursor: pointer;
}

.hero-inner {
  display: flex;
  gap: var(--sp-4);
  padding: var(--sp-4);
}

.hero-image-wrap {
  width: 42%;
  min-height: 130px;
  border-radius: var(--r-sm);
  overflow: hidden;
  flex-shrink: 0;
  background: var(--color-surface-3);
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
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 650;
  line-height: 1.2;
  overflow-wrap: anywhere;
}

.hero-time {
  color: var(--color-text-3);
  font-size: var(--text-xs);
}

.hero-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
}

.hero-tags {
  color: var(--color-text-3);
  font-size: var(--text-xs);
}

.hero-badge {
  display: inline-flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
  padding: var(--sp-2) var(--sp-3);
  border-radius: var(--r-sm);
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-border);
}

.hero-badge-text {
  color: var(--color-text);
  font-size: var(--text-xs);
  font-weight: 600;
}

.hero-badge-date {
  color: var(--color-text-3);
  font-size: var(--text-2xs);
}

/* ── Meal Nav ── */
.meal-nav {
  display: flex;
  gap: var(--sp-2);
  margin-bottom: var(--sp-8);
}

.meal-pill {
  flex: 1;
  min-height: 40px;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease), color var(--dur-fast) var(--ease), border-color var(--dur-fast) var(--ease), transform var(--dur-fast) var(--ease);
}

.meal-pill:hover {
  border-color: var(--color-border-med);
}

.meal-pill.active {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
}

.meal-pill:active {
  transform: translateY(1px);
}

/* ── Sections ── */
.section {
  margin-bottom: var(--sp-8);
}

.section-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-4);
}

.section-more {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 500;
}

/* ── Recipe Scroll ── */
.recipe-scroll {
  display: flex;
  gap: var(--sp-3);
  overflow-x: auto;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;
  margin: 0 calc(-1 * var(--sp-5));
  padding: 0 var(--sp-5) var(--sp-2);
}

.recipe-scroll::-webkit-scrollbar { display: none; }

.recipe-card {
  flex: 0 0 132px;
  scroll-snap-align: start;
  cursor: pointer;
}

.recipe-cover {
  width: 100%;
  aspect-ratio: 1;
  border-radius: var(--r-sm);
  overflow: hidden;
  background: var(--color-surface-2);
}

.recipe-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.recipe-name {
  margin-top: var(--sp-2);
  color: var(--color-text);
  font-size: var(--text-sm);
  font-weight: 600;
  line-height: 1.3;
}

.recipe-meta {
  margin-top: 2px;
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  color: var(--color-text-3);
  font-size: var(--text-xs);
}

.recipe-diff {
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--color-accent-soft);
  color: var(--color-accent);
  font-weight: 600;
  font-size: var(--text-2xs);
}

/* ── Quick Access ── */
.quick-list {
  border-top: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
}

.quick-row {
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
  transition: background var(--dur-fast) var(--ease);
}

.quick-row:hover {
  background: var(--color-surface-2);
}

.quick-row:active {
  background: var(--color-surface-3);
}

.quick-row + .quick-row {
  border-top: 1px solid var(--color-border);
}

.quick-index {
  width: 32px;
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
  letter-spacing: 0.06em;
  flex-shrink: 0;
}

.quick-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.quick-title {
  color: var(--color-text);
  font-size: var(--text-base);
  font-weight: 600;
}

.quick-sub {
  margin-top: 1px;
  color: var(--color-text-3);
  font-size: var(--text-xs);
}

.quick-arrow {
  width: 16px;
  height: 16px;
  color: var(--color-text-3);
  flex-shrink: 0;
}

/* ── Empty ── */
.empty-hint {
  padding: var(--sp-8) 0;
  text-align: center;
  color: var(--color-text-3);
  font-size: var(--text-sm);
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .page {
    max-width: 640px;
    padding: 0 var(--sp-8);
  }

  .hero-inner {
    padding: var(--sp-6);
  }

  .recipe-card {
    flex: 0 0 160px;
  }
}

@media (min-width: 1024px) {
  .page {
    max-width: 800px;
  }
}

@media (min-width: 1440px) {
  .page {
    max-width: 960px;
  }
}

@media (max-width: 420px) {
  .hero-inner {
    flex-direction: column;
  }

  .hero-image-wrap {
    width: 100%;
    min-height: 150px;
  }
}
</style>
