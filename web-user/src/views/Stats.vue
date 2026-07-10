<template>
  <div class="stats-shell">
    <main class="stats-phone">
      <header class="stats-header">
        <button class="back-btn" type="button" aria-label="返回" @click="router.back()">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="m15 18-6-6 6-6" /></svg>
        </button>
        <div>
          <p>Food habits</p>
          <h1>饮食统计</h1>
        </div>
      </header>

      <div class="period-switch" role="group" aria-label="统计周期">
        <button v-for="option in periodOptions" :key="option.value" type="button" :class="{ active: period === option.value }" :aria-pressed="period === option.value" @click="period = option.value">
          {{ option.label }}
        </button>
      </div>

      <p v-if="errorText" class="notice error" role="alert">
        {{ errorText }} <button type="button" @click="loadStats">重试</button>
      </p>

      <template v-else-if="stats">
        <section class="overview-card" aria-label="本期做饭概览">
          <div>
            <span>{{ period === 'week' ? '本周' : '近 30 天' }}已做</span>
            <strong>{{ stats.cooked_count }}</strong>
            <small>道菜谱记录</small>
          </div>
          <dl>
            <div><dt>自炊天数</dt><dd>{{ stats.cooked_days }} 天</dd></div>
            <div><dt>连续做饭</dt><dd>{{ stats.current_streak }} 天</dd></div>
          </dl>
        </section>

        <section v-if="stats.cooked_count" class="panel nutrition-panel" aria-labelledby="nutrition-title">
          <div class="panel-title">
            <div><p>按已做菜谱营养数据汇总</p><h2 id="nutrition-title">营养摄入参考</h2></div>
            <span>{{ stats.nutrition_completeness }}% 数据完整</span>
          </div>
          <p v-if="!stats.nutrition_recorded_recipes" class="empty-copy">已做菜谱暂未维护营养字段，暂不展示估算值。</p>
          <div v-else class="nutrition-grid">
            <div><span>热量</span><strong>{{ formatNumber(stats.nutrition.calories) }}</strong><small>kcal</small></div>
            <div><span>蛋白质</span><strong>{{ formatNumber(stats.nutrition.protein) }}</strong><small>g</small></div>
            <div><span>脂肪</span><strong>{{ formatNumber(stats.nutrition.fat) }}</strong><small>g</small></div>
            <div><span>碳水</span><strong>{{ formatNumber(stats.nutrition.carbs) }}</strong><small>g</small></div>
          </div>
          <p v-if="stats.nutrition_completeness < 100" class="data-note">部分已做菜谱缺少营养数据；缺失数据没有按 0 参与计算。</p>
        </section>

        <section v-if="stats.cooked_count" class="panel trend-panel" aria-labelledby="trend-title">
          <div class="panel-title"><div><p>{{ stats.start_date }} 至 {{ stats.end_date }}</p><h2 id="trend-title">做饭趋势</h2></div></div>
          <ul class="trend-list" aria-label="每日做饭次数">
            <li v-for="item in stats.daily" :key="item.date">
              <span class="bar-wrap"><i :style="{ height: barHeight(item.cooked_count) }"></i></span>
              <strong>{{ item.cooked_count || '–' }}</strong>
              <small>{{ shortDate(item.date) }}</small>
            </li>
          </ul>
        </section>

        <section v-if="stats.recent_cooked_recipes.length" class="panel recent-panel" aria-labelledby="recent-title">
          <div class="panel-title"><div><p>基于你标记“做过了”的真实记录</p><h2 id="recent-title">最近做过</h2></div></div>
          <button v-for="recipe in stats.recent_cooked_recipes" :key="recipe.recipe_id" type="button" class="recent-row" @click="router.push(`/recipes/${recipe.recipe_id}`)">
            <span>{{ recipe.title }}</span><small>{{ recipe.cooked_at }}</small>
          </button>
        </section>
      </template>

      <section v-else-if="!loading" class="empty-card">
        <div class="empty-icon" aria-hidden="true"><svg viewBox="0 0 48 48"><path d="M12 34c0-9 6-16 16-20 2 10-1 20-10 24" /><path d="M14 34c7-1 13-6 18-15" /><path d="M8 38h32" /></svg></div>
        <h2>还没有已做记录</h2>
        <p>在菜谱详情点“做过了”，这里就会汇总你的自炊次数和已维护的营养数据。</p>
        <button type="button" @click="router.push('/recipes')">去选一道菜</button>
      </section>

      <p v-if="loading" class="notice" role="status">正在汇总真实饮食记录…</p>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getFoodStats, type FoodStats } from '@/api/user'

const router = useRouter()
const period = ref<'week' | 'month'>('week')
const stats = ref<FoodStats | null>(null)
const loading = ref(false)
const errorText = ref('')
const periodOptions = [
  { value: 'week' as const, label: '本周' },
  { value: 'month' as const, label: '近 30 天' },
]

const maxCookedCount = computed(() => Math.max(1, ...(stats.value?.daily.map((item) => item.cooked_count) || [1])))

function formatNumber(value: number) {
  return Math.round(Number(value || 0) * 10) / 10
}

function shortDate(value: string) {
  return value.slice(5).replace('-', '/')
}

function barHeight(value: number) {
  if (!value) return '4px'
  return `${Math.max(12, Math.round((value / maxCookedCount.value) * 64))}px`
}

async function loadStats() {
  loading.value = true
  errorText.value = ''
  try {
    stats.value = await getFoodStats(period.value)
  } catch (error) {
    stats.value = null
    errorText.value = error instanceof Error ? error.message : '统计数据读取失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

watch(period, loadStats, { immediate: true })
</script>

<style scoped>
.stats-shell { min-height: 100dvh; color: #2e241f; background: linear-gradient(180deg, #fff2dc 0%, #f7e4c7 46%, #ead0aa 100%); }
.stats-phone { width: min(100%, 520px); min-height: 100dvh; margin: 0 auto; padding: max(24px, env(safe-area-inset-top)) 24px 126px; }
.stats-header { display: grid; grid-template-columns: 48px minmax(0, 1fr); align-items: center; gap: 14px; padding-top: 18px; }
.back-btn { width: 48px; height: 48px; display: grid; place-items: center; border: 1px solid rgba(255,255,255,.7); border-radius: 18px; color: #6f6258; background: rgba(255,250,240,.78); box-shadow: 0 12px 28px rgba(80,50,30,.12); cursor: pointer; }
.back-btn svg { width: 24px; stroke: currentColor; fill: none; stroke-width: 2.5; stroke-linecap: round; stroke-linejoin: round; }
.stats-header p, .panel-title p { margin: 0 0 3px; color: #e95645; font-size: 12px; font-weight: 800; }
.stats-header h1, .panel-title h2 { margin: 0; color: #2e241f; font-weight: 900; }
.stats-header h1 { font-size: 30px; }
.period-switch { display: inline-flex; gap: 4px; margin-top: 24px; padding: 4px; border-radius: 14px; background: rgba(255,255,255,.56); }
.period-switch button { min-height: 38px; padding: 0 16px; border: 0; border-radius: 10px; color: #75685f; background: transparent; font: inherit; font-weight: 750; cursor: pointer; }
.period-switch button.active { color: #fff; background: #e95645; box-shadow: 0 6px 14px rgba(233,86,69,.24); }
.overview-card, .panel, .empty-card { margin-top: 16px; border: 1px solid var(--card-border); border-radius: var(--card-radius); background: var(--card-surface); box-shadow: var(--card-shadow); backdrop-filter: blur(18px); }
.overview-card { display: flex; justify-content: space-between; gap: 20px; padding: 24px; background: linear-gradient(135deg, rgba(233,86,69,.96), rgba(244,130,82,.9)); color: #fff; }
.overview-card span, .overview-card small, .overview-card dt { display: block; opacity: .82; font-size: 12px; font-weight: 700; }
.overview-card strong { display: block; margin: 2px 0; font-size: 52px; line-height: 1; }
.overview-card dl { display: grid; align-content: center; gap: 14px; margin: 0; }
.overview-card dd { margin: 2px 0 0; font-size: 18px; font-weight: 850; }
.panel { padding: 20px; }
.panel-title { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; }
.panel-title h2 { font-size: 21px; }
.panel-title > span { padding: 5px 8px; border-radius: 8px; color: #725b49; background: rgba(245,219,189,.8); font-size: 11px; font-weight: 800; }
.nutrition-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px; margin-top: 16px; }
.nutrition-grid div { padding: 13px; border-radius: 14px; background: rgba(255,255,255,.58); }
.nutrition-grid span, .nutrition-grid small { color: #7a6a5f; font-size: 12px; font-weight: 700; }
.nutrition-grid strong { margin: 4px 4px 0 0; color: #3b2a22; font-size: 22px; }
.data-note, .empty-copy { margin: 14px 0 0; color: #7a6a5f; font-size: 12px; line-height: 1.55; }
.trend-list { height: 112px; display: grid; grid-template-columns: repeat(auto-fit, minmax(21px, 1fr)); gap: 6px; align-items: end; margin: 16px 0 0; padding: 0; list-style: none; }
.trend-list li { min-width: 0; display: grid; grid-template-rows: 70px 17px 15px; justify-items: center; gap: 2px; }
.bar-wrap { width: 100%; display: flex; align-items: end; justify-content: center; border-radius: 7px; background: rgba(241,207,169,.38); }
.bar-wrap i { width: 100%; max-width: 18px; border-radius: 7px; background: linear-gradient(#e95645, #f29a61); }
.trend-list strong { font-size: 12px; }.trend-list small { color: #8b7d72; font-size: 10px; white-space: nowrap; }
.recent-row { width: 100%; display: flex; justify-content: space-between; gap: 10px; padding: 14px 0; border: 0; border-top: 1px solid rgba(118,87,63,.12); color: #3b2a22; background: transparent; font: inherit; text-align: left; cursor: pointer; }.recent-row:first-of-type { margin-top: 14px; }.recent-row span { font-weight: 750; }.recent-row small { color: #8b7d72; font-size: 12px; }
.empty-card { min-height: 360px; display: grid; place-items: center; align-content: center; gap: 16px; padding: 32px 28px; text-align: center; }.empty-icon { width: 88px; color: #6f8b65; }.empty-icon svg { width: 100%; fill: none; stroke: currentColor; stroke-width: 2.4; stroke-linecap: round; stroke-linejoin: round; }.empty-card h2 { margin: 0; font-size: 24px; }.empty-card p { max-width: 300px; margin: 0; color: #7a6a5f; line-height: 1.7; }.empty-card button, .notice button { min-height: 44px; padding: 0 18px; border: 0; border-radius: 999px; color: #fff; background: #e95645; font: inherit; font-weight: 800; cursor: pointer; }
.notice { margin: 20px 0; color: #75685f; font-size: 14px; text-align: center; }.notice.error { color: #ba382e; }.notice.error button { min-height: 32px; margin-left: 8px; padding: 0 12px; }
@media (max-width: 380px) { .stats-phone { padding-left: 18px; padding-right: 18px; }.overview-card { padding: 20px; }.overview-card strong { font-size: 44px; }.trend-list { gap: 3px; } }
</style>
