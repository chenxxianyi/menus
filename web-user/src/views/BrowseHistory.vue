<template>
  <main class="history-shell" :style="pageVars">
    <header class="page-header" aria-label="页面顶部">
      <button class="nav-btn" type="button" aria-label="返回" @click="router.back()">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <button class="clear-btn" type="button" :disabled="!hasHistory || loading" @click="clearHistory">清空</button>
      <h1 class="page-title">浏览历史</h1>
      <p class="summary-line">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <circle cx="12" cy="12" r="8.5" />
          <path d="M12 7v5l3 2" />
        </svg>
        <span>最近看过 <strong>{{ total }}</strong> 道菜</span>
      </p>
    </header>

    <section v-if="loading" class="empty-card">
      <div class="loading-spinner" aria-hidden="true"></div>
      <h2>正在读取浏览历史</h2>
      <p>只展示你的真实浏览记录</p>
    </section>

    <section v-else-if="hasHistory" class="history-list" aria-label="浏览历史列表">
      <section
        v-for="group in historyGroups"
        :key="group.group"
        class="group-card"
        :aria-labelledby="`group-${group.group}`"
      >
        <div class="group-title-row">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M12 22c3.7 0 6.8-2.9 6.8-6.9 0-3.2-2-5.3-4.1-7.6-.9-1-1.7-2.1-2.1-3.5-2.4 1.5-3.2 3.8-2.8 6.2-1-.4-1.7-1.2-2.2-2.3C6 9.4 5.2 11.6 5.2 14c0 4.8 3.1 8 6.8 8Z" />
          </svg>
          <h2 :id="`group-${group.group}`">{{ group.group }}</h2>
          <span class="group-line"></span>
        </div>

        <article
          v-for="item in group.items"
          :key="`${group.group}-${item.id}`"
          class="recipe-row"
          :style="recipeStyle(item.recipe)"
          @click="openRecipe(item.recipe)"
        >
          <div class="food-photo" role="img" :aria-label="recipeTitle(item.recipe)">
            <span v-if="!item.recipe.cover">{{ recipeTag(item.recipe) }}</span>
          </div>
          <div class="recipe-body">
            <h3 class="recipe-title">{{ recipeTitle(item.recipe) }}</h3>
            <div class="tag-row">
              <span v-for="tag in recipeTags(item.recipe)" :key="tag" class="tag-chip">{{ tag }}</span>
            </div>
            <p class="favorite-line">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" />
              </svg>
              <span>收藏 {{ formatCount(item.recipe.favorite_count) }}</span>
            </p>
          </div>
          <span class="time">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="12" cy="12" r="8.5" />
              <path d="M12 7v5l3 2" />
            </svg>
            {{ formatTime(item.viewed_at) }}
          </span>
          <svg class="row-arrow" viewBox="0 0 24 24" aria-hidden="true">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </article>
      </section>
    </section>

    <section v-else class="empty-card">
      <div class="empty-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24">
          <path d="M3 12a9 9 0 1 0 3-6.7" />
          <path d="M3 3v6h6" />
          <path d="M12 7v5l3 2" />
        </svg>
      </div>
      <h2>还没有浏览记录</h2>
      <p>去看看今天想吃什么</p>
      <button class="empty-browse-btn" type="button" @click="browseRecipes">去逛菜谱</button>
    </section>

    <button v-if="hasHistory" class="continue-btn" type="button" @click="browseRecipes">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M7 20c0-5 2.5-8 8-9" />
        <path d="M9 20c5.5 0 9-3.5 9-9V4h-7c-5.5 0-9 3.5-9 9 0 2.1.7 3.9 2 5.2" />
      </svg>
      <span class="continue-text">
        <strong>继续逛菜谱</strong>
        <span>发现更多美味灵感</span>
      </span>
    </button>

    <p v-if="hasHistory" class="retention-tip">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
        <path d="m9 12 2 2 4-5" />
      </svg>
      <span>仅保留 60 天内的浏览记录</span>
    </p>

    <div class="toast" :class="{ show: !!toastText }">{{ toastText }}</div>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { clearBrowseHistory, getBrowseHistory } from '@/api/user'

interface HistoryRecipe {
  id: number
  title?: string
  cover?: string
  category_name?: string
  cook_time?: number
  difficulty?: string
  taste?: string
  favorite_count?: number
}

interface BrowseHistoryItem {
  id: number
  viewed_at: string
  recipe: HistoryRecipe
}

interface HistoryGroup {
  group: string
  items: BrowseHistoryItem[]
}

const router = useRouter()
const historyItems = ref<BrowseHistoryItem[]>([])
const total = ref(0)
const loading = ref(true)
const toastText = ref('')
let toastTimer: ReturnType<typeof setTimeout> | null = null

const pageVars = computed(() => ({
  '--history-bg': `url(${kitchenBg})`,
}))

const hasHistory = computed(() => historyItems.value.length > 0)

const historyGroups = computed<HistoryGroup[]>(() => {
  const groups = new Map<string, BrowseHistoryItem[]>()
  historyItems.value.forEach((item) => {
    const label = groupLabel(item.viewed_at)
    const list = groups.get(label) || []
    list.push(item)
    groups.set(label, list)
  })

  return Array.from(groups.entries()).map(([group, items]) => ({ group, items }))
})

function normalizePage(payload: any) {
  if (Array.isArray(payload)) return { list: payload, total: payload.length }
  const list = payload?.list || payload?.items || payload?.data || []
  return {
    list: Array.isArray(list) ? list : [],
    total: Number(payload?.total ?? payload?.count ?? list.length),
  }
}

async function fetchHistory() {
  loading.value = true
  try {
    const res: any = await getBrowseHistory(1, 50)
    const normalized = normalizePage(res)
    historyItems.value = normalized.list.filter((item: BrowseHistoryItem) => item?.recipe?.id)
    total.value = normalized.total
  } catch {
    historyItems.value = []
    total.value = 0
    showToast('浏览历史读取失败')
  } finally {
    loading.value = false
  }
}

function recipeTitle(recipe: HistoryRecipe) {
  return recipe.title || '未命名菜谱'
}

function recipeTag(recipe: HistoryRecipe) {
  return recipe.category_name || recipe.taste || recipe.difficulty || '菜谱'
}

function recipeTags(recipe: HistoryRecipe) {
  const tags = [recipe.category_name, recipe.taste, recipe.cook_time ? `${recipe.cook_time}分钟` : '']
    .filter((tag): tag is string => !!tag)
  return tags.length ? tags.slice(0, 3) : ['暂无标签']
}

function recipeStyle(recipe: HistoryRecipe) {
  if (recipe.cover) {
    return { '--img': `url('${recipe.cover}')` }
  }
  return { '--img': 'none' }
}

function formatCount(count?: number) {
  const value = Number(count || 0)
  if (value >= 10000) return `${(value / 10000).toFixed(1)}w`
  if (value >= 1000) return `${(value / 1000).toFixed(1)}k`
  return String(value)
}

function parseDate(value: string) {
  const date = new Date(value)
  return Number.isFinite(date.getTime()) ? date : new Date()
}

function startOfDay(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate()).getTime()
}

function groupLabel(value: string) {
  const viewed = parseDate(value)
  const today = startOfDay(new Date())
  const day = startOfDay(viewed)
  if (day === today) return '今天'
  if (day === today - 86400000) return '昨天'
  return `${viewed.getMonth() + 1}月${viewed.getDate()}日`
}

function formatTime(value: string) {
  const date = parseDate(value)
  const h = String(date.getHours()).padStart(2, '0')
  const m = String(date.getMinutes()).padStart(2, '0')
  return `${h}:${m}`
}

function showToast(message: string) {
  toastText.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastText.value = ''
  }, 1500)
}

async function clearHistory() {
  if (!hasHistory.value) return
  if (!window.confirm('确定清空浏览历史吗？')) return
  try {
    await clearBrowseHistory()
    historyItems.value = []
    total.value = 0
    showToast('已清空浏览历史')
  } catch {
    showToast('清空失败，请稍后再试')
  }
}

function openRecipe(recipe: HistoryRecipe) {
  if (!recipe.id) return
  router.push(`/recipes/${recipe.id}`)
}

function browseRecipes() {
  router.push('/recipes')
}

onMounted(() => {
  fetchHistory()
})

onUnmounted(() => {
  if (toastTimer) clearTimeout(toastTimer)
})
</script>

<style scoped>
.history-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.82);
  --coral: #e95645;
  --sage: #8fa783;
  --border: rgba(255, 255, 255, 0.62);
  --line: rgba(120, 90, 65, 0.14);
  --shadow: 0 22px 50px rgba(80, 50, 30, 0.16);
  position: relative;
  width: min(100%, 430px);
  min-height: calc(100vh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 150px);
  min-height: calc(100dvh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 150px);
  margin: 0 auto;
  padding: max(52px, env(safe-area-inset-top)) 24px calc(42px + env(safe-area-inset-bottom));
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 237, 205, 0.34), rgba(255, 247, 233, 0.18)),
    var(--history-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.history-shell::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 4%, rgba(255, 255, 255, 0.74), transparent 30%),
    radial-gradient(circle at 88% 15%, rgba(238, 143, 66, 0.24), transparent 34%),
    radial-gradient(circle at 10% 90%, rgba(233, 86, 69, 0.18), transparent 30%),
    linear-gradient(90deg, rgba(255, 239, 214, 0.56), rgba(255, 245, 230, 0.18) 54%, rgba(172, 91, 33, 0.17));
  backdrop-filter: blur(4px) saturate(1.12);
  -webkit-backdrop-filter: blur(4px) saturate(1.12);
}

.page-header,
.history-list,
.empty-card,
.continue-btn,
.retention-tip {
  position: relative;
  z-index: 1;
}

button {
  border: 0;
  font: inherit;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

svg {
  display: block;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.page-header {
  min-height: 146px;
  display: grid;
  place-items: start center;
  align-content: start;
  text-align: center;
}

.nav-btn,
.clear-btn {
  position: absolute;
  top: 0;
  height: 52px;
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 16px;
  color: #4a352a;
  background: rgba(255, 250, 240, 0.86);
  box-shadow:
    0 12px 28px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  transition: transform 180ms ease, opacity 180ms ease, box-shadow 180ms ease;
}

.nav-btn {
  left: 0;
  width: 52px;
  display: grid;
  place-items: center;
}

.nav-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.55;
}

.clear-btn {
  right: 0;
  min-width: 72px;
  padding: 0 18px;
  font-size: 16px;
  font-weight: 760;
  letter-spacing: 0;
}

.page-title {
  margin: 14px 0 0;
  color: var(--text);
  font-size: 29px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
}

.summary-line {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 26px 0 0;
  color: #6b5142;
  font-size: 16px;
  font-weight: 650;
  line-height: 1;
}

.summary-line svg {
  width: 20px;
  height: 20px;
  stroke-width: 2.25;
}

.summary-line strong {
  color: var(--coral);
  font-size: 20px;
  font-weight: 950;
  line-height: 1;
}

.history-list {
  display: grid;
  gap: 22px;
}

.group-card {
  overflow: hidden;
  padding: 22px 20px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 32px;
  background:
    radial-gradient(circle at 10% 0%, rgba(255, 255, 255, 0.48), transparent 42%),
    rgba(255, 250, 240, 0.82);
  box-shadow:
    var(--shadow),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(22px) saturate(1.1);
  -webkit-backdrop-filter: blur(22px) saturate(1.1);
}

.group-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 18px;
}

.group-title-row svg {
  width: 18px;
  height: 18px;
  flex: 0 0 18px;
  color: var(--coral);
  stroke-width: 2.8;
}

.group-title-row h2 {
  margin: 0;
  color: var(--text);
  font-size: 25px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
  white-space: nowrap;
}

.group-line {
  height: 1px;
  flex: 1 1 auto;
  margin-left: 12px;
  background: var(--line);
}

.recipe-row {
  position: relative;
  min-height: 130px;
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 8px 30px 8px 0;
  border-bottom: 1px solid rgba(120, 90, 65, 0.12);
  transition: transform 180ms ease, opacity 180ms ease, background 180ms ease;
}

.recipe-row:last-child {
  border-bottom: 0;
}

.food-photo {
  width: 124px;
  height: 100px;
  display: grid;
  flex: 0 0 124px;
  place-items: center;
  overflow: hidden;
  border-radius: 18px;
  background:
    var(--img) center / cover no-repeat,
    linear-gradient(135deg, #f7dfb8, #e5a15d);
  box-shadow:
    0 12px 24px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.food-photo span {
  width: 58px;
  height: 58px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  color: rgba(255, 255, 255, 0.94);
  background: rgba(233, 86, 69, 0.32);
  font-size: 13px;
  font-weight: 850;
}

.recipe-body {
  min-width: 0;
  flex: 1 1 auto;
  padding-top: 20px;
}

.recipe-title {
  max-width: calc(100% - 40px);
  margin: 0;
  overflow: hidden;
  color: var(--text);
  font-size: 23px;
  font-weight: 950;
  line-height: 1.2;
  letter-spacing: 0;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-row {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  overflow: hidden;
}

.tag-chip {
  height: 28px;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  padding: 0 10px;
  border-radius: 999px;
  color: #8a6857;
  background: rgba(255, 238, 222, 0.72);
  font-size: 13px;
  font-weight: 680;
  line-height: 1;
  white-space: nowrap;
}

.favorite-line {
  display: flex;
  align-items: center;
  gap: 5px;
  margin: 10px 0 0;
  color: #8a7669;
  font-size: 13px;
  font-weight: 620;
  line-height: 1;
}

.favorite-line svg {
  width: 17px;
  height: 17px;
  stroke-width: 2.1;
}

.time {
  position: absolute;
  top: 22px;
  right: 22px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #8a7669;
  font-size: 14px;
  font-weight: 650;
  line-height: 1;
}

.time svg {
  width: 17px;
  height: 17px;
  stroke-width: 2.2;
}

.row-arrow {
  position: absolute;
  top: 50%;
  right: 0;
  width: 20px;
  height: 20px;
  color: var(--muted);
  stroke-width: 2.3;
  transform: translateY(-50%);
}

.continue-btn {
  width: min(72%, 260px);
  min-height: 86px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  margin: 26px auto 0;
  padding: 14px 24px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 28px;
  color: var(--coral);
  background:
    radial-gradient(circle at 10% 0%, rgba(255, 255, 255, 0.45), transparent 42%),
    rgba(255, 250, 240, 0.82);
  box-shadow:
    0 18px 42px rgba(80, 50, 30, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.84);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.continue-btn svg {
  width: 27px;
  height: 27px;
  flex: 0 0 27px;
  stroke-width: 2.4;
}

.continue-text {
  text-align: left;
}

.continue-text strong {
  display: block;
  color: var(--coral);
  font-size: 20px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
  white-space: nowrap;
}

.continue-text span {
  display: block;
  margin-top: 7px;
  color: var(--sub);
  font-size: 14px;
  font-weight: 650;
  line-height: 1;
  white-space: nowrap;
}

.retention-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 22px 0 0;
  color: var(--sub);
  font-size: 15px;
  font-weight: 650;
  line-height: 1.35;
  text-align: center;
}

.retention-tip svg {
  width: 21px;
  height: 21px;
  flex: 0 0 21px;
  color: var(--sage);
  stroke-width: 2.25;
}

.empty-card {
  margin-top: 8px;
  padding: 44px 24px 34px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 32px;
  background:
    radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.5), transparent 42%),
    rgba(255, 250, 240, 0.84);
  box-shadow: var(--shadow), inset 0 1px 0 rgba(255, 255, 255, 0.84);
  text-align: center;
  backdrop-filter: blur(22px) saturate(1.1);
  -webkit-backdrop-filter: blur(22px) saturate(1.1);
}

.empty-icon {
  width: 74px;
  height: 74px;
  display: grid;
  place-items: center;
  margin: 0 auto 20px;
  border-radius: 24px;
  color: var(--coral);
  background: rgba(252, 231, 222, 0.82);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.64);
}

.empty-icon svg {
  width: 40px;
  height: 40px;
  stroke-width: 2.2;
}

.empty-card h2 {
  margin: 0;
  color: var(--text);
  font-size: 23px;
  font-weight: 950;
  line-height: 1;
}

.empty-card p {
  margin: 13px 0 24px;
  color: var(--sub);
  font-size: 15px;
  font-weight: 620;
  line-height: 1.5;
}

.empty-browse-btn {
  width: 100%;
  height: 54px;
  border-radius: 18px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 14px 26px rgba(233, 86, 69, 0.25);
  font-size: 17px;
  font-weight: 860;
  transition: transform 180ms ease, opacity 180ms ease;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  margin: 0 auto 20px;
  border: 3px solid rgba(233, 86, 69, 0.18);
  border-top-color: var(--coral);
  border-radius: 50%;
  animation: spin 0.72s linear infinite;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: calc(28px + env(safe-area-inset-bottom));
  z-index: 5;
  padding: 10px 16px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 999px;
  color: #fff;
  background: rgba(46, 36, 31, 0.78);
  box-shadow: 0 12px 24px rgba(46, 36, 31, 0.18);
  font-size: 14px;
  font-weight: 740;
  opacity: 0;
  pointer-events: none;
  transform: translate(-50%, 12px);
  transition: opacity 180ms ease, transform 180ms ease;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  white-space: nowrap;
}

.toast.show {
  opacity: 1;
  transform: translate(-50%, 0);
}

.nav-btn:active,
.clear-btn:active,
.recipe-row:active,
.continue-btn:active,
.empty-browse-btn:active {
  transform: scale(0.98);
}

.recipe-row:active {
  opacity: 0.88;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (hover: hover) {
  .nav-btn:hover,
  .clear-btn:hover:not(:disabled),
  .continue-btn:hover,
  .empty-browse-btn:hover {
    transform: translateY(-1px);
  }

  .recipe-row:hover {
    opacity: 0.9;
  }
}

@media (max-width: 380px) {
  .history-shell {
    padding-right: 18px;
    padding-left: 18px;
  }

  .group-card {
    padding: 21px 18px;
    border-radius: 30px;
  }

  .recipe-row {
    gap: 14px;
    min-height: 124px;
    padding-right: 24px;
  }

  .food-photo {
    width: 118px;
    height: 94px;
    flex-basis: 118px;
  }

  .recipe-title {
    font-size: 21px;
  }

  .tag-chip {
    padding: 0 9px;
    font-size: 12.5px;
  }

  .time {
    right: 18px;
  }

  .continue-btn {
    width: 76%;
  }
}

@media (max-width: 350px) {
  .history-shell {
    padding-right: 16px;
    padding-left: 16px;
  }

  .page-title {
    font-size: 27px;
  }

  .clear-btn {
    min-width: 68px;
    padding: 0 15px;
  }

  .recipe-row {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
    padding: 8px 24px 16px 0;
  }

  .food-photo {
    width: 100%;
    height: 142px;
    flex-basis: auto;
  }

  .recipe-body {
    width: 100%;
    padding-top: 0;
  }

  .time {
    top: 166px;
  }

  .continue-btn {
    width: 100%;
  }
}

@media (min-width: 431px) {
  .history-shell {
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.18);
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
</style>
