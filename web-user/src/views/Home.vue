<template>
  <div class="page home-page">
    <header class="home-header anim-delay-1">
      <div class="home-greeting">
        <div class="avatar">{{ userInitial }}</div>
        <div>
          <h1 class="greeting-title">早上好，美食家</h1>
          <p class="greeting-sub">今天想吃点什么？</p>
        </div>
      </div>
      <button class="pantry-btn" type="button" aria-label="查看食材清单" @click="$router.push('/shopping-list')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M6 2h12a2 2 0 0 1 2 2v16a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2Z"/>
          <path d="M4 9h16"/>
          <path d="M8 6h2"/>
          <path d="M8 13h2"/>
        </svg>
      </button>
    </header>

    <div class="search-box anim-delay-2" @click="goToRecipes">
      <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8" />
        <path d="m21 21-4.35-4.35" />
      </svg>
      <span class="search-placeholder">搜索菜名、食材或剩余时间</span>
    </div>

    <section class="couple-entry anim-delay-3" @click="$router.push('/couple')">
      <div class="couple-entry-copy">
        <p class="couple-entry-kicker">双人点餐</p>
        <h2 class="couple-entry-title">TA 想吃什么，一起定下来</h2>
        <p class="couple-entry-desc">绑定双方账号，发布想吃、匹配菜谱、生成采购清单。</p>
      </div>
      <div class="couple-entry-art" aria-hidden="true">
        <span class="couple-avatar self">{{ userInitial }}</span>
        <span class="couple-link">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M8 12h8"/>
            <path d="M13 7l5 5-5 5"/>
          </svg>
        </span>
        <span class="couple-avatar partner">TA</span>
      </div>
    </section>

    <section class="today-plate anim-delay-4" @click="openHero">
      <div class="plate-copy">
        <p class="plate-eyebrow">{{ issueDate }} · 今日主厨建议</p>
        <h2 class="plate-title">{{ heroTitle }}</h2>
        <div class="plate-meta">
          <span>{{ heroTime }}</span>
          <span>{{ heroDifficulty }}</span>
          <span>{{ heroPeople }}</span>
        </div>
      </div>
      <div class="plate-visual" aria-hidden="true">
        <div class="food-bowl">
          <span class="rice-base"></span>
          <span class="food-piece tomato-one"></span>
          <span class="food-piece tomato-two"></span>
          <span class="food-piece egg-one"></span>
          <span class="food-piece egg-two"></span>
          <span class="food-piece leaf-one"></span>
          <span class="food-piece leaf-two"></span>
        </div>
        <div class="steam steam-one"></div>
        <div class="steam steam-two"></div>
      </div>
    </section>

    <section class="meal-tabs anim-delay-5" aria-label="选择餐次">
      <button
        v-for="meal in meals"
        :key="meal.key"
        class="meal-tab"
        :class="{ active: activeMeal === meal.key }"
        type="button"
        @click="activeMeal = meal.key"
      >
        {{ meal.label }}
      </button>
    </section>

    <section class="section anim-delay-5">
      <div class="section-top">
        <div>
          <p class="section-kicker">balanced menu</p>
          <h2 class="section-heading">今日推荐菜单</h2>
        </div>
        <button class="section-more" type="button" @click="$router.push('/week-menu')">查看周计划</button>
      </div>

      <div class="menu-scroll">
        <article
          v-for="menu in visibleMenuCards"
          :key="menu.key"
          class="menu-card"
          :class="menu.key"
        >
          <div class="menu-card-top">
            <span class="menu-tag">{{ menu.tag }}</span>
            <span class="menu-time">
              <svg class="time-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              {{ menu.time }}
            </span>
          </div>
          <div class="menu-illustration" aria-hidden="true">
            <span class="bowl"></span>
            <span class="ingredient dot-a"></span>
            <span class="ingredient dot-b"></span>
            <span class="ingredient dot-c"></span>
          </div>
          <h3 class="menu-title">{{ menu.title }}</h3>
          <p class="menu-desc">{{ menu.desc }}</p>
          <div class="menu-footer">
            <span class="menu-calories">{{ menu.calories }}</span>
            <button class="menu-btn" type="button" @click.stop="goToRecipes">去制作</button>
          </div>
        </article>
      </div>
    </section>

    <section class="section tools-section">
      <div class="quick-grid">
        <button
          v-for="item in quickActions"
          :key="item.label"
          class="quick-item"
          type="button"
          @click="item.path && $router.push(item.path)"
        >
          <div class="quick-icon" :class="item.tone">
            <svg v-if="item.icon === 'sprout'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/>
              <path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/>
            </svg>
            <svg v-else-if="item.icon === 'flame'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2c.5 2.5 2 4.5 2 7a4 4 0 1 1-8 0c0-2.5 1.5-4.5 2-7 1.3 1.5 3 2 4 2s2.7-.5 4-2z"/>
            </svg>
            <svg v-else-if="item.icon === 'target'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <circle cx="12" cy="12" r="6"/>
              <circle cx="12" cy="12" r="2"/>
            </svg>
            <svg v-else-if="item.icon === 'calendar'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/>
              <line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
            <svg v-else-if="item.icon === 'fridge'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 2h12a2 2 0 0 1 2 2v16a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2Z"/>
              <path d="M4 9h16"/>
              <path d="M8 6h2"/>
              <path d="M8 13h2"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
          </div>
          <span class="quick-label">{{ item.label }}</span>
        </button>
      </div>
    </section>

    <section class="section hot-section">
      <div class="section-top">
        <div>
          <p class="section-kicker">popular recipes</p>
          <h2 class="section-heading">热门菜谱</h2>
        </div>
        <button class="section-more" type="button" @click="goToRecipes">更多</button>
      </div>

      <div v-if="hotRecipes.length" class="recipe-grid">
        <article
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
        </article>
      </div>
      <div v-else class="empty-hint">
        <span class="empty-dish"></span>
        <p>后厨还在备菜，稍后再来看看。</p>
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

const heroTitle = computed(() => todayRecommend.value?.title || '为今天留一份温柔的晚餐')
const heroTime = computed(() => todayRecommend.value ? `${todayRecommend.value.cook_time} min` : '20 min')
const heroDifficulty = computed(() => todayRecommend.value?.difficulty || '简单')
const heroPeople = computed(() => todayRecommend.value?.people_count ? `${todayRecommend.value.people_count} 人份` : '2 人份')
const issueDate = computed(() =>
  new Intl.DateTimeFormat('zh-CN', { month: 'long', day: 'numeric' }).format(new Date())
)

const menuCards = [
  {
    key: 'breakfast',
    tag: '早餐推荐',
    time: '10min',
    title: '能量活力早餐',
    desc: '燕麦牛奶 + 水煮蛋 + 全麦面包',
    calories: '320 kcal',
  },
  {
    key: 'lunch',
    tag: '午餐推荐',
    time: '25min',
    title: '经典家常套餐',
    desc: '番茄炒蛋 + 青椒肉丝 + 米饭',
    calories: '650 kcal',
  },
  {
    key: 'dinner',
    tag: '晚餐推荐',
    time: '30min',
    title: '轻盈低脂晚餐',
    desc: '清蒸鲈鱼 + 蒜蓉西兰花 + 杂粮饭',
    calories: '480 kcal',
  },
  {
    key: 'snack',
    tag: '夜宵推荐',
    time: '12min',
    title: '轻负担夜宵',
    desc: '鲜虾豆腐汤 + 小份荞麦面',
    calories: '260 kcal',
  },
]

const quickActions = [
  { label: '按食材推荐', icon: 'sprout', tone: 'orange', path: '/recipes' },
  { label: '按口味推荐', icon: 'flame', tone: 'red', path: '/recipes' },
  { label: '按场景推荐', icon: 'target', tone: 'blue', path: '/recipes' },
  { label: '一周菜单', icon: 'calendar', tone: 'green', path: '/week-menu' },
  { label: '冰箱剩菜', icon: 'fridge', tone: 'purple', path: '/recipes' },
  { label: '最热菜谱', icon: 'trend', tone: 'yellow', path: '/recipes' },
]

const visibleMenuCards = computed(() =>
  menuCards.filter((menu) => menu.key === activeMeal.value)
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
/* Hallmark - pre-emit critique: P4 H4 E4 S4 R4 V4 */
.home-page {
  min-height: 100vh;
  padding-top: var(--sp-5);
  padding-bottom: calc(var(--tab-h) + var(--sp-12));
  background:
    linear-gradient(150deg, var(--color-herb-mist), transparent 34%),
    linear-gradient(24deg, var(--color-broth-soft), transparent 42%),
    var(--color-bg);
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-5);
}

.home-greeting {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.avatar {
  width: 48px;
  height: 48px;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--glass-border);
  border-radius: var(--r-full);
  background:
    radial-gradient(circle at 30% 24%, var(--color-surface), transparent 32%),
    var(--glass-bg);
  box-shadow: var(--glass-shadow);
  color: var(--color-text);
  font-size: 16px;
  font-weight: 700;
  backdrop-filter: blur(var(--glass-blur));
}

.greeting-title {
  color: var(--color-text);
  font-size: var(--text-base);
  font-weight: 700;
  line-height: 1.2;
}

.greeting-sub {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 500;
  margin-top: 2px;
}

.pantry-btn {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border: 1px solid var(--glass-border);
  border-radius: var(--r-full);
  background: var(--glass-bg);
  color: var(--color-text-2);
  box-shadow: var(--glass-shadow);
  cursor: pointer;
  transition:
    color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease),
    box-shadow var(--dur-fast) var(--ease);
  backdrop-filter: blur(var(--glass-blur));
}

.pantry-btn:hover {
  color: var(--color-accent);
  box-shadow: var(--shadow-md);
}

.pantry-btn:active {
  transform: translateY(1px);
}

.pantry-btn svg {
  width: 20px;
  height: 20px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  min-height: 54px;
  padding: 0 var(--sp-5);
  border: 1px solid var(--glass-border);
  border-radius: var(--r-full);
  margin-bottom: var(--sp-5);
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  cursor: pointer;
  transition:
    border-color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease),
    box-shadow var(--dur-fast) var(--ease);
  backdrop-filter: blur(var(--glass-blur));
}

.search-box:hover {
  border-color: var(--color-border-med);
  box-shadow: var(--shadow-md);
}

.search-box:active {
  transform: translateY(1px);
}

.search-icon {
  width: 20px;
  height: 20px;
  color: var(--color-text-3);
  flex-shrink: 0;
}

.search-placeholder {
  color: var(--color-text-3);
  font-size: var(--text-sm);
  font-weight: 600;
}

.couple-entry {
  position: relative;
  min-height: 128px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 132px;
  gap: var(--sp-3);
  align-items: center;
  overflow: hidden;
  padding: var(--sp-5);
  border: 1px solid rgba(255, 255, 255, 0.76);
  border-radius: var(--r-xl);
  margin-bottom: var(--sp-5);
  background:
    radial-gradient(circle at 86% 20%, rgba(255, 255, 255, 0.78), transparent 30%),
    linear-gradient(135deg, var(--color-accent-soft), var(--color-surface));
  box-shadow: 0 18px 42px rgba(7, 193, 96, 0.12);
  cursor: pointer;
  transition:
    transform var(--dur-base) var(--ease-out),
    box-shadow var(--dur-base) var(--ease-out);
}

.couple-entry::before {
  content: "";
  position: absolute;
  right: -30px;
  bottom: -42px;
  width: 140px;
  height: 140px;
  border-radius: var(--r-full);
  background: var(--color-tomato-soft);
}

.couple-entry:hover {
  transform: translateY(-2px);
  box-shadow: 0 22px 48px rgba(7, 193, 96, 0.16);
}

.couple-entry:active {
  transform: translateY(1px);
}

.couple-entry-copy,
.couple-entry-art {
  position: relative;
  z-index: 1;
}

.couple-entry-kicker {
  margin-bottom: var(--sp-1);
  color: var(--color-accent);
  font-size: var(--text-xs);
  font-weight: 800;
}

.couple-entry-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 750;
  line-height: 1.25;
  overflow-wrap: anywhere;
}

.couple-entry-desc {
  max-width: 18em;
  margin-top: var(--sp-2);
  color: var(--color-text-2);
  font-size: var(--text-xs);
  font-weight: 600;
  line-height: 1.55;
}

.couple-entry-art {
  min-height: 78px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.couple-avatar {
  width: 50px;
  height: 50px;
  display: grid;
  place-items: center;
  border: 1px solid var(--glass-border);
  border-radius: var(--r-full);
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  color: var(--color-text);
  font-size: var(--text-xs);
  font-weight: 800;
  backdrop-filter: blur(var(--glass-blur));
}

.couple-avatar.partner {
  background: var(--color-surface);
  color: var(--color-accent);
}

.couple-link {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  margin: 0 -4px;
  border-radius: var(--r-full);
  background: var(--color-accent);
  color: var(--color-text-inv);
  box-shadow: 0 10px 20px rgba(7, 193, 96, 0.24);
}

.couple-link svg {
  width: 17px;
  height: 17px;
}

.today-plate {
  position: relative;
  min-height: 188px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 132px;
  gap: var(--sp-2);
  align-items: center;
  overflow: hidden;
  padding: var(--sp-5);
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: var(--r-xl);
  margin-bottom: var(--sp-4);
  background:
    radial-gradient(circle at 84% 18%, rgba(255, 255, 255, 0.82), transparent 28%),
    linear-gradient(135deg, var(--color-surface), var(--color-broth-soft));
  box-shadow: 0 18px 42px rgba(31, 41, 55, 0.10);
  cursor: pointer;
}

.today-plate::before {
  content: "";
  position: absolute;
  right: -34px;
  bottom: -46px;
  width: 156px;
  height: 156px;
  border-radius: var(--r-full);
  background: var(--color-tomato-soft);
}

.plate-copy {
  position: relative;
  z-index: 1;
  min-width: 0;
}

.plate-eyebrow {
  margin-bottom: var(--sp-2);
  color: var(--color-accent-hover);
  font-size: var(--text-xs);
  font-weight: 700;
}

.plate-title {
  max-width: 9em;
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 700;
  line-height: 1.18;
  overflow-wrap: anywhere;
}

.plate-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--sp-2);
  margin-top: var(--sp-4);
}

.plate-meta span {
  padding: 5px 9px;
  border-radius: var(--r-full);
  background: rgba(255, 255, 255, 0.72);
  color: var(--color-text-2);
  font-size: var(--text-2xs);
  font-weight: 700;
}

.plate-visual {
  position: relative;
  z-index: 1;
  width: 126px;
  height: 126px;
  display: grid;
  place-items: center;
  border-radius: var(--r-full);
  background:
    radial-gradient(circle at 50% 50%, var(--color-surface) 0 46%, transparent 47%),
    conic-gradient(from 20deg, var(--color-accent), var(--color-tomato), var(--color-warning), var(--color-accent));
  box-shadow:
    inset 0 0 0 12px rgba(255, 255, 255, 0.72),
    0 18px 34px rgba(7, 193, 96, 0.18);
}

.food-bowl {
  position: relative;
  width: 86px;
  height: 86px;
  border-radius: var(--r-full);
  background:
    radial-gradient(circle at 50% 48%, rgba(255, 255, 255, 0.98) 0 48%, transparent 49%),
    linear-gradient(145deg, var(--color-broth-soft), var(--color-surface));
  box-shadow:
    inset 0 -10px 0 rgba(245, 158, 11, 0.10),
    inset 0 0 0 8px rgba(255, 255, 255, 0.86);
}

.food-bowl::after {
  content: "";
  position: absolute;
  left: 11px;
  right: 11px;
  bottom: 12px;
  height: 18px;
  border-radius: 0 0 26px 26px;
  background: rgba(249, 115, 22, 0.16);
}

.rice-base {
  position: absolute;
  left: 17px;
  top: 28px;
  width: 52px;
  height: 34px;
  border-radius: 20px;
  background: radial-gradient(circle at 30% 32%, #fff 0 12%, #fff7e8 13% 100%);
  box-shadow: 0 6px 14px rgba(31, 41, 55, 0.08);
}

.food-piece {
  position: absolute;
  display: block;
}

.tomato-one,
.tomato-two {
  width: 20px;
  height: 18px;
  border-radius: 48% 52% 50% 50%;
  background: var(--color-tomato);
}

.tomato-one {
  left: 23px;
  top: 23px;
  transform: rotate(-18deg);
}

.tomato-two {
  right: 22px;
  top: 40px;
  transform: rotate(20deg);
}

.egg-one,
.egg-two {
  width: 22px;
  height: 16px;
  border-radius: var(--r-full);
  background: var(--color-warning-soft);
  box-shadow: inset 7px 0 0 rgba(245, 158, 11, 0.34);
}

.egg-one {
  left: 42px;
  top: 20px;
  transform: rotate(16deg);
}

.egg-two {
  left: 25px;
  top: 47px;
  transform: rotate(-20deg);
}

.leaf-one,
.leaf-two {
  width: 22px;
  height: 12px;
  border-radius: 22px 3px 22px 3px;
  background: var(--color-accent);
}

.leaf-one {
  right: 25px;
  top: 27px;
  transform: rotate(-34deg);
}

.leaf-two {
  left: 41px;
  top: 50px;
  transform: rotate(28deg);
}

.steam {
  position: absolute;
  top: 8px;
  width: 18px;
  height: 34px;
  border-top: 2px solid rgba(255, 255, 255, 0.88);
  border-radius: var(--r-full);
}

.steam-one {
  left: 32px;
  transform: rotate(-22deg);
}

.steam-two {
  right: 34px;
  transform: rotate(18deg);
}

.meal-tabs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--sp-2);
  padding: 5px;
  border: 1px solid var(--glass-border);
  border-radius: var(--r-full);
  margin-bottom: var(--sp-6);
  background: rgba(255, 255, 255, 0.58);
  box-shadow: var(--glass-shadow);
  backdrop-filter: blur(var(--glass-blur));
}

.meal-tab {
  min-height: 36px;
  border: 0;
  border-radius: var(--r-full);
  background: transparent;
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 700;
  cursor: pointer;
  transition:
    background var(--dur-base) var(--ease-spring),
    color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease),
    box-shadow var(--dur-base) var(--ease);
}

.meal-tab:hover {
  color: var(--color-text-2);
}

.meal-tab:active {
  transform: translateY(1px);
}

.meal-tab.active {
  background: var(--color-surface);
  color: var(--color-accent);
  box-shadow: var(--shadow-sm);
}

.section {
  margin-bottom: var(--sp-8);
}

.section-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-4);
}

.section-kicker {
  margin-bottom: 2px;
  color: var(--color-text-3);
  font-size: var(--text-2xs);
  font-weight: 800;
  text-transform: uppercase;
}

.section-heading {
  color: var(--color-text);
  font-size: var(--text-lg);
  font-weight: 700;
  line-height: 1.2;
}

.section-more {
  border: 0;
  background: transparent;
  color: var(--color-accent);
  font-size: var(--text-xs);
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
}

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
  position: relative;
  min-width: 260px;
  overflow: hidden;
  padding: var(--sp-4);
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: var(--r-lg);
  background: var(--glass-bg);
  box-shadow: 0 14px 34px rgba(31, 41, 55, 0.09);
  scroll-snap-align: start;
  backdrop-filter: blur(var(--glass-blur));
  transition:
    transform var(--dur-base) var(--ease-out),
    box-shadow var(--dur-base) var(--ease-out);
}

.menu-card::after {
  content: "";
  position: absolute;
  right: -26px;
  top: 38px;
  width: 98px;
  height: 98px;
  border-radius: var(--r-full);
  opacity: 0.72;
}

.menu-card.breakfast::after {
  background: var(--color-herb-mist);
}

.menu-card.lunch::after {
  background: var(--color-tomato-soft);
}

.menu-card.dinner::after {
  background: var(--color-water-soft);
}

.menu-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.menu-card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--sp-3);
}

.menu-tag {
  padding: 4px 9px;
  border-radius: var(--r-full);
  background: var(--color-accent-soft);
  color: var(--color-accent);
  font-size: 10px;
  font-weight: 700;
}

.menu-card.lunch .menu-tag {
  background: var(--color-tomato-soft);
  color: var(--color-tomato);
}

.menu-card.dinner .menu-tag {
  background: var(--color-water-soft);
  color: var(--color-water);
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

.menu-illustration {
  position: relative;
  width: 84px;
  height: 64px;
  margin-bottom: var(--sp-3);
}

.bowl {
  position: absolute;
  left: 5px;
  bottom: 0;
  width: 70px;
  height: 38px;
  border-radius: 10px 10px 34px 34px;
  background: linear-gradient(180deg, var(--color-surface), var(--color-broth-soft));
  box-shadow: inset 0 -8px 0 rgba(245, 158, 11, 0.10), 0 8px 16px rgba(31, 41, 55, 0.08);
}

.ingredient {
  position: absolute;
  width: 14px;
  height: 14px;
  border-radius: var(--r-full);
}

.dot-a {
  left: 20px;
  top: 10px;
  background: var(--color-accent);
}

.dot-b {
  left: 43px;
  top: 2px;
  background: var(--color-tomato);
}

.dot-c {
  right: 13px;
  top: 18px;
  background: var(--color-warning);
}

.menu-title {
  color: var(--color-text);
  font-size: var(--text-base);
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
  color: var(--color-tomato);
}

.menu-btn {
  min-height: 32px;
  padding: 0 14px;
  background: var(--color-accent);
  color: var(--color-text-inv);
  border: none;
  border-radius: var(--r-full);
  font-size: var(--text-xs);
  font-weight: 700;
  cursor: pointer;
  transition:
    background var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease);
}

.menu-btn:hover {
  background: var(--color-accent-hover);
}

.menu-btn:active {
  transform: translateY(1px);
}

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
  min-height: 106px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: var(--r-lg);
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  cursor: pointer;
  backdrop-filter: blur(var(--glass-blur));
  transition:
    border-color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease),
    box-shadow var(--dur-fast) var(--ease);
}

.quick-item:hover {
  border-color: var(--color-border-med);
  box-shadow: var(--shadow-md);
}

.quick-item:active {
  transform: translateY(1px);
}

.quick-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-full);
}

.quick-icon svg {
  width: 20px;
  height: 20px;
}

.quick-icon.orange {
  background: var(--color-tomato-soft);
  color: var(--color-tomato);
}

.quick-icon.red {
  background: var(--color-berry-soft);
  color: var(--color-berry);
}

.quick-icon.blue {
  background: var(--color-water-soft);
  color: var(--color-water);
}

.quick-icon.green {
  background: var(--color-herb-mist);
  color: var(--color-accent);
}

.quick-icon.purple {
  background: var(--color-plum-soft);
  color: var(--color-plum);
}

.quick-icon.yellow {
  background: var(--color-broth-soft);
  color: var(--color-warning);
}

.quick-label {
  font-size: var(--text-xs);
  font-weight: 700;
  color: var(--color-text-2);
  text-align: center;
}

.recipe-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--sp-4);
}

.recipe-card {
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: var(--r-lg);
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  cursor: pointer;
  backdrop-filter: blur(var(--glass-blur));
  transition:
    transform var(--dur-base) var(--ease-out),
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
  aspect-ratio: 1 / 0.78;
  background: var(--color-broth-soft);
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
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: blur(8px);
  border-radius: var(--r-full);
  padding: 2px 8px;
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: var(--text-2xs);
  color: var(--color-text);
  font-weight: 800;
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
  font-size: var(--text-sm);
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
  gap: var(--sp-2);
  font-size: var(--text-2xs);
  color: var(--color-text-3);
}

.recipe-diff {
  font-weight: 500;
  min-width: 0;
}

.recipe-time {
  font-weight: 400;
  white-space: nowrap;
}

.empty-hint {
  display: grid;
  place-items: center;
  gap: var(--sp-3);
  padding: var(--sp-10) 0;
  text-align: center;
  color: var(--color-text-3);
  font-size: var(--text-sm);
}

.empty-dish {
  width: 58px;
  height: 58px;
  border-radius: var(--r-full);
  background:
    radial-gradient(circle, var(--color-surface) 0 42%, transparent 43%),
    conic-gradient(var(--color-accent), var(--color-tomato), var(--color-warning), var(--color-accent));
  box-shadow: inset 0 0 0 8px rgba(255, 255, 255, 0.78);
}

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
    min-width: 244px;
  }
}

@media (max-width: 360px) {
  .couple-entry {
    grid-template-columns: minmax(0, 1fr) 104px;
    padding: var(--sp-4);
  }

  .couple-avatar {
    width: 42px;
    height: 42px;
  }

  .couple-link {
    width: 30px;
    height: 30px;
  }

  .today-plate {
    grid-template-columns: minmax(0, 1fr) 104px;
    padding: var(--sp-4);
  }

  .plate-visual {
    width: 104px;
    height: 104px;
  }

  .food-bowl {
    width: 70px;
    height: 70px;
  }

  .rice-base {
    left: 13px;
    top: 23px;
    width: 44px;
    height: 28px;
  }

  .tomato-one,
  .tomato-two {
    width: 16px;
    height: 15px;
  }

  .egg-one,
  .egg-two {
    width: 18px;
    height: 13px;
  }

  .leaf-one,
  .leaf-two {
    width: 18px;
    height: 10px;
  }

  .quick-grid {
    gap: var(--sp-2);
  }
}
</style>
