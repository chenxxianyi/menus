<template>
  <div class="page home-page">
    <!-- Header -->
    <header class="home-header anim-delay-1">
      <div class="home-greeting">
        <div class="avatar">{{ userInitial }}</div>
        <div>
          <h1 class="greeting-title">早上好，美食家</h1>
          <p class="greeting-sub">今天想吃点什么？</p>
        </div>
      </div>
    </header>

    <!-- Search Box -->
    <div class="search-box anim-delay-2" @click="goToRecipes">
      <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8" />
        <path d="m21 21-4.35-4.35" />
      </svg>
      <span class="search-placeholder">搜索你感兴趣的菜谱...</span>
    </div>

    <!-- Today's Menu Recommendations -->
    <section class="section anim-delay-3">
      <div class="section-top">
        <h2 class="section-heading">今日推荐菜单</h2>
        <span class="section-more">查看周计划</span>
      </div>

      <div class="menu-scroll">
        <!-- Breakfast -->
        <div class="menu-card">
          <div class="menu-card-header">
            <span class="menu-tag breakfast">早餐推荐</span>
            <span class="menu-time">
              <svg class="time-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              10min
            </span>
          </div>
          <h3 class="menu-title">能量活力早餐</h3>
          <p class="menu-desc">燕麦牛奶 + 水煮蛋 + 全麦面包</p>
          <div class="menu-footer">
            <span class="menu-calories">320 kcal</span>
            <button class="menu-btn">去制作</button>
          </div>
        </div>

        <!-- Lunch -->
        <div class="menu-card">
          <div class="menu-card-header">
            <span class="menu-tag lunch">午餐推荐</span>
            <span class="menu-time">
              <svg class="time-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              25min
            </span>
          </div>
          <h3 class="menu-title">经典家常套餐</h3>
          <p class="menu-desc">番茄炒蛋 + 青椒肉丝 + 米饭</p>
          <div class="menu-footer">
            <span class="menu-calories">650 kcal</span>
            <button class="menu-btn">去制作</button>
          </div>
        </div>

        <!-- Dinner -->
        <div class="menu-card">
          <div class="menu-card-header">
            <span class="menu-tag dinner">晚餐推荐</span>
            <span class="menu-time">
              <svg class="time-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              30min
            </span>
          </div>
          <h3 class="menu-title">轻盈低脂晚餐</h3>
          <p class="menu-desc">清蒸鲈鱼 + 蒜蓉西兰花 + 杂粮饭</p>
          <div class="menu-footer">
            <span class="menu-calories">480 kcal</span>
            <button class="menu-btn">去制作</button>
          </div>
        </div>
      </div>
    </section>

    <!-- Quick Access Grid -->
    <section class="section anim-delay-4">
      <div class="quick-grid">
        <button class="quick-item" @click="$router.push('/week-menu')">
          <div class="quick-icon orange">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/>
              <path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/>
            </svg>
          </div>
          <span class="quick-label">按食材推荐</span>
        </button>
        <button class="quick-item">
          <div class="quick-icon red">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2c.5 2.5 2 4.5 2 7a4 4 0 1 1-8 0c0-2.5 1.5-4.5 2-7 1.3 1.5 3 2 4 2s2.7-.5 4-2z"/>
            </svg>
          </div>
          <span class="quick-label">按口味推荐</span>
        </button>
        <button class="quick-item">
          <div class="quick-icon blue">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <circle cx="12" cy="12" r="6"/>
              <circle cx="12" cy="12" r="2"/>
            </svg>
          </div>
          <span class="quick-label">按场景推荐</span>
        </button>
        <button class="quick-item" @click="$router.push('/week-menu')">
          <div class="quick-icon green">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/>
              <line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
          </div>
          <span class="quick-label">一周菜单</span>
        </button>
        <button class="quick-item">
          <div class="quick-icon purple">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </div>
          <span class="quick-label">冰箱剩菜</span>
        </button>
        <button class="quick-item">
          <div class="quick-icon yellow">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
          </div>
          <span class="quick-label">最热菜谱</span>
        </button>
      </div>
    </section>

    <!-- Hot Recipes Grid -->
    <section class="section anim-delay-5">
      <div class="section-top">
        <h2 class="section-heading">热门菜谱</h2>
        <span class="section-more" @click="goToRecipes">更多</span>
      </div>

      <div class="recipe-grid">
        <div
          v-for="recipe in hotRecipes.slice(0, 4)"
          :key="recipe.id"
          class="recipe-card"
          @click="goToRecipe(recipe.id)"
        >
          <div class="recipe-cover">
            <img :src="recipe.cover || heroArt" :alt="recipe.title" />
            <div class="recipe-rating">
              <svg class="star-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              {{ recipe.rating || '4.8' }}
            </div>
          </div>
          <div class="recipe-info">
            <h3 class="recipe-name">{{ recipe.title }}</h3>
            <p class="recipe-meta">
              <span class="recipe-diff">{{ recipe.difficulty }}</span>
              <span class="recipe-time">{{ recipe.cook_time }}min</span>
            </p>
          </div>
        </div>
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

function goToRecipes() {
  router.push('/recipes')
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
  padding-top: var(--sp-4);
  padding-bottom: var(--sp-8);
}

/* ── Header ── */
.home-header {
  margin-bottom: var(--sp-6);
}

.home-greeting {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--color-surface-2);
  border: 2px solid var(--color-surface);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text);
  flex-shrink: 0;
}

.greeting-title {
  color: var(--color-text);
  font-size: var(--text-sm);
  font-weight: 700;
  line-height: 1.2;
}

.greeting-sub {
  color: var(--color-text-3);
  font-size: 10px;
  font-weight: 400;
  margin-top: 2px;
}

/* ── Search Box ── */
.search-box {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  padding: 10px 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-full);
  margin-bottom: var(--sp-6);
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--dur-fast) var(--ease);
}

.search-box:hover {
  box-shadow: var(--shadow-md);
}

.search-icon {
  width: 18px;
  height: 18px;
  color: var(--color-text-3);
  flex-shrink: 0;
}

.search-placeholder {
  color: var(--color-text-3);
  font-size: var(--text-sm);
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

.section-heading {
  color: var(--color-text);
  font-size: var(--text-lg);
  font-weight: 700;
}

.section-more {
  color: var(--color-accent);
  font-size: var(--text-xs);
  font-weight: 500;
  cursor: pointer;
}

/* ── Menu Scroll ── */
.menu-scroll {
  display: flex;
  gap: var(--sp-4);
  overflow-x: auto;
  scroll-snap-type: x mandatory;
  -webkit-overflow-scrolling: touch;
  margin: 0 calc(-1 * var(--sp-5));
  padding: 0 var(--sp-5) var(--sp-2);
}

.menu-scroll::-webkit-scrollbar { display: none; }

.menu-card {
  min-width: 240px;
  background: var(--color-surface);
  border-radius: var(--r-lg);
  padding: var(--sp-4);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  scroll-snap-align: start;
  transition: transform var(--dur-base) var(--ease-out),
              box-shadow var(--dur-base) var(--ease-out);
}

.menu-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.menu-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--sp-3);
}

.menu-tag {
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 10px;
  font-weight: 700;
}

.menu-tag.breakfast {
  background: #E8F5E9;
  color: #07C160;
}

.menu-tag.lunch {
  background: #FFF3E0;
  color: #F57C00;
}

.menu-tag.dinner {
  background: #E3F2FD;
  color: #1976D2;
}

.menu-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: var(--color-text-3);
}

.time-icon {
  width: 12px;
  height: 12px;
}

.menu-title {
  color: var(--color-text);
  font-size: var(--text-sm);
  font-weight: 700;
  margin-bottom: 4px;
}

.menu-desc {
  color: var(--color-text-2);
  font-size: var(--text-xs);
  margin-bottom: var(--sp-3);
}

.menu-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.menu-calories {
  font-size: var(--text-xs);
  font-weight: 700;
  color: #F57C00;
}

.menu-btn {
  padding: 4px 12px;
  background: var(--color-accent);
  color: var(--color-text-inv);
  border: none;
  border-radius: var(--r-full);
  font-size: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.menu-btn:hover {
  background: var(--color-accent-hover);
}

.menu-btn:active {
  transform: scale(0.95);
}

/* ── Quick Grid ── */
.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--sp-3);
}

.quick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--sp-2);
  padding: var(--sp-3);
  background: var(--color-surface);
  border-radius: var(--r-md);
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: transform var(--dur-fast) var(--ease),
              box-shadow var(--dur-fast) var(--ease);
}

.quick-item:hover {
  transform: scale(0.95);
  box-shadow: var(--shadow-md);
}

.quick-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quick-icon svg {
  width: 20px;
  height: 20px;
}

.quick-icon.orange {
  background: #FFF3E0;
  color: #F57C00;
}

.quick-icon.red {
  background: #FFEBEE;
  color: #E53935;
}

.quick-icon.blue {
  background: #E3F2FD;
  color: #1976D2;
}

.quick-icon.green {
  background: #E8F5E9;
  color: #07C160;
}

.quick-icon.purple {
  background: #F3E5F5;
  color: #7B1FA2;
}

.quick-icon.yellow {
  background: #FFFDE7;
  color: #FBC02D;
}

.quick-label {
  font-size: 10px;
  font-weight: 500;
  color: var(--color-text-2);
  text-align: center;
}

/* ── Recipe Grid ── */
.recipe-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--sp-4);
}

.recipe-card {
  background: var(--color-surface);
  border-radius: var(--r-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: transform var(--dur-base) var(--ease-out),
              box-shadow var(--dur-base) var(--ease-out);
}

.recipe-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}

.recipe-card:active {
  transform: translateY(0);
}

.recipe-cover {
  width: 100%;
  height: 112px;
  background: var(--color-surface-2);
  position: relative;
}

.recipe-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.recipe-rating {
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  border-radius: var(--r-full);
  padding: 2px 8px;
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 8px;
  color: var(--color-text-inv);
}

.star-icon {
  width: 10px;
  height: 10px;
  color: #FBC02D;
}

.recipe-info {
  padding: var(--sp-3);
}

.recipe-name {
  color: var(--color-text);
  font-size: var(--text-xs);
  font-weight: 700;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recipe-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 9px;
  color: var(--color-text-3);
}

.recipe-diff {
  font-weight: 500;
}

.recipe-time {
  font-weight: 400;
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
  .menu-card {
    min-width: 200px;
  }
}
</style>
