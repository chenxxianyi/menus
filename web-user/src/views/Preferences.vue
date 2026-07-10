<template>
  <main class="preferences-shell" :style="pageVars">
    <div class="status-bar" aria-hidden="true">
      <span>9:41</span>
      <div class="status-icons">
        <span class="cell-bars"><i></i><i></i><i></i><i></i></span>
        <svg class="wifi" viewBox="0 0 24 18">
          <path d="M2.8 5.8a14.5 14.5 0 0 1 18.4 0" />
          <path d="M6.8 9.8a8.6 8.6 0 0 1 10.4 0" />
          <path d="M10.3 13.5a3 3 0 0 1 3.4 0" />
        </svg>
        <span class="battery"></span>
      </div>
    </div>

    <header class="page-header" aria-label="页面顶部">
      <button class="back-btn" type="button" aria-label="返回" @click="handleBack">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <h1 class="page-title">{{ isGuideMode ? '快速了解你' : '偏好设置' }}</h1>
      <button class="save-btn" type="button" :disabled="saving" @click="handleSave">
        {{ saving ? '保存中' : '保存' }}
      </button>
    </header>

    <section class="settings-card" aria-labelledby="taste-title">
      <div class="card-title-row">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M20.8 4.6a5.4 5.4 0 0 0-7.6 0L12 5.8l-1.2-1.2a5.4 5.4 0 1 0-7.6 7.6L12 21l8.8-8.8a5.4 5.4 0 0 0 0-7.6Z" />
        </svg>
        <h2 class="card-title" id="taste-title">口味偏好</h2>
      </div>
      <p class="card-desc">选择你喜欢的口味，我们会更懂你的喜好</p>
      <div class="taste-grid">
        <button
          v-for="taste in tasteOptions"
          :key="taste"
          class="taste-chip"
          :class="{ active: form.taste_preference.includes(taste) }"
          type="button"
          @click="toggleTaste(taste)"
        >
          {{ taste }}
        </button>
      </div>
    </section>

    <section class="settings-card" aria-labelledby="goal-title">
      <div class="card-title-row">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <circle cx="12" cy="12" r="8" />
          <circle cx="12" cy="12" r="3" />
          <path d="M12 2v3" />
          <path d="M12 19v3" />
          <path d="M2 12h3" />
          <path d="M19 12h3" />
        </svg>
        <h2 class="card-title" id="goal-title">饮食目标</h2>
      </div>
      <p class="card-desc">根据目标推荐更合适的菜谱和搭配</p>
      <div class="goal-scroll">
        <button
          v-for="goal in goalOptions"
          :key="goal.value"
          class="goal-chip"
          :class="{ active: form.health_goal === goal.value }"
          type="button"
          @click="form.health_goal = goal.value"
        >
          {{ goal.label }}
        </button>
      </div>
    </section>

    <section class="settings-card" aria-labelledby="time-title">
      <div class="card-title-row">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <circle cx="12" cy="12" r="8.5" />
          <path d="M12 7v5l3 2" />
        </svg>
        <h2 class="card-title" id="time-title">做饭时间</h2>
      </div>
      <p class="card-desc">告诉我们你通常愿意花多久做一餐</p>
      <div class="goal-scroll">
        <button
          v-for="item in cookTimeOptions"
          :key="item"
          class="goal-chip"
          :class="{ active: form.cook_time_preference === item }"
          type="button"
          @click="form.cook_time_preference = item"
        >
          {{ item }}
        </button>
      </div>
    </section>

    <section class="settings-card" aria-labelledby="avoid-title">
      <div class="card-title-row">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <circle cx="12" cy="12" r="9" />
          <path d="m5.7 5.7 12.6 12.6" />
        </svg>
        <h2 class="card-title" id="avoid-title">忌口食材</h2>
      </div>
      <p class="card-desc">添加你不吃的食材，避免为你推荐</p>
      <label class="avoid-input" aria-label="添加忌口食材">
        <input
          v-model="avoidInput"
          type="text"
          placeholder="输入食材后回车添加"
          @keydown.enter.prevent="addAvoidIngredient"
        />
        <button class="add-ingredient" type="button" aria-label="添加食材" @click="addAvoidIngredient">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M12 5v14" />
            <path d="M5 12h14" />
          </svg>
        </button>
      </label>
      <div v-if="form.avoid_ingredients.length" class="avoid-tags">
        <span v-for="name in form.avoid_ingredients" :key="name" class="ingredient-tag">
          {{ name }}
          <button type="button" :aria-label="`删除${name}`" @click="removeAvoidIngredient(name)">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M18 6 6 18" />
              <path d="m6 6 12 12" />
            </svg>
          </button>
        </span>
      </div>
    </section>

    <section class="settings-card" aria-labelledby="serving-title">
      <div class="card-title-row">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
          <circle cx="9" cy="7" r="4" />
          <path d="M22 21v-2a4 4 0 0 0-3-3.9" />
          <path d="M16 3.1a4 4 0 0 1 0 7.8" />
        </svg>
        <h2 class="card-title" id="serving-title">常用用餐人数</h2>
      </div>
      <p class="card-desc">推荐适合份量的菜谱</p>
      <div class="stepper" aria-label="常用用餐人数">
        <button class="step-btn" type="button" aria-label="减少人数" :disabled="form.people_count <= 1" @click="adjustPeople(-1)">-</button>
        <strong class="serving-count">{{ form.people_count }}</strong>
        <button class="step-btn" type="button" aria-label="增加人数" :disabled="form.people_count >= 12" @click="adjustPeople(1)">+</button>
      </div>
    </section>

    <p class="save-tip">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
        <path d="m9 12 2 2 4-5" />
      </svg>
      <span>保存后会用于个性化推荐，随时可修改</span>
    </p>

    <button v-if="isGuideMode" class="skip-btn" type="button" @click="skipGuide">先跳过，稍后再完善</button>

    <div class="toast" :class="{ show: !!toastText }">{{ toastText }}</div>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { getPreferences, updatePreferences } from '@/api/user'

type HealthGoal = '普通' | '减脂' | '增肌' | '控糖' | '儿童营养'

interface PreferenceForm {
  taste_preference: string[]
  health_goal: HealthGoal
  avoid_ingredients: string[]
  cook_time_preference: string
  people_count: number
}

const router = useRouter()
const route = useRoute()

const tasteOptions = ['咸鲜', '酸甜', '麻辣', '清淡', '香辣', '酸辣', '甜味', '原味']
const cookTimeOptions = ['15分钟内', '30分钟内', '45分钟内', '都可以']
const goalOptions: { value: HealthGoal; label: string; legacy: string[] }[] = [
  { value: '普通', label: '普通', legacy: ['normal', '普通'] },
  { value: '减脂', label: '减脂', legacy: ['lose', '减脂'] },
  { value: '增肌', label: '增肌', legacy: ['gain', '增肌'] },
  { value: '控糖', label: '控糖', legacy: ['sugar', '控糖'] },
  { value: '儿童营养', label: '儿童营养', legacy: ['child', '儿童营养'] },
]

const saving = ref(false)
const avoidInput = ref('')
const toastText = ref('')
let toastTimer: ReturnType<typeof setTimeout> | null = null

const form = reactive<PreferenceForm>({
  taste_preference: [],
  health_goal: '普通',
  avoid_ingredients: [],
  cook_time_preference: '30分钟内',
  people_count: 2,
})

const pageVars = computed(() => ({
  '--preference-bg': `url(${kitchenBg})`,
}))
const isGuideMode = computed(() => route.query.guide === '1')

function normalizeArray(value: unknown) {
  if (Array.isArray(value)) {
    return value.filter((item): item is string => typeof item === 'string' && item.trim().length > 0)
  }
  if (typeof value === 'string') {
    try {
      const parsed = JSON.parse(value)
      if (Array.isArray(parsed)) {
        return parsed.filter((item): item is string => typeof item === 'string' && item.trim().length > 0)
      }
    } catch {
      return value
        .split(/[,，、]/)
        .map((item) => item.trim())
        .filter(Boolean)
    }
  }
  return []
}

function normalizeGoal(value: unknown): HealthGoal {
  const text = typeof value === 'string' ? value : ''
  return goalOptions.find((goal) => goal.legacy.includes(text))?.value || '普通'
}

function normalizePeopleCount(value: unknown) {
  const count = Number(value)
  if (!Number.isFinite(count)) return 2
  return Math.min(12, Math.max(1, Math.round(count)))
}

function showToast(message: string) {
  toastText.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastText.value = ''
  }, 1500)
}

function toggleTaste(taste: string) {
  const index = form.taste_preference.indexOf(taste)
  if (index >= 0) {
    form.taste_preference.splice(index, 1)
  } else {
    form.taste_preference.push(taste)
  }
}

function addAvoidIngredient() {
  const value = avoidInput.value.trim().slice(0, 10)
  if (!value) return
  if (form.avoid_ingredients.includes(value)) {
    showToast('已经添加过该食材')
    avoidInput.value = ''
    return
  }
  form.avoid_ingredients.push(value)
  avoidInput.value = ''
}

function removeAvoidIngredient(name: string) {
  const index = form.avoid_ingredients.indexOf(name)
  if (index >= 0) {
    form.avoid_ingredients.splice(index, 1)
  }
}

function adjustPeople(delta: number) {
  form.people_count = Math.min(12, Math.max(1, form.people_count + delta))
}

async function loadPreferences() {
  try {
    const res: any = await getPreferences()
    if (!res) return
    form.taste_preference = normalizeArray(res.taste_preference)
    form.health_goal = normalizeGoal(res.health_goal)
    form.avoid_ingredients = normalizeArray(res.avoid_ingredients)
    form.cook_time_preference = typeof res.cook_time_preference === 'string' && res.cook_time_preference.trim()
      ? res.cook_time_preference.trim()
      : '30分钟内'
    form.people_count = normalizePeopleCount(res.people_count)
  } catch {
    showToast('偏好读取失败')
  }
}

async function handleSave() {
  if (saving.value) return
  saving.value = true
  try {
    await updatePreferences({
      taste_preference: [...form.taste_preference],
      health_goal: form.health_goal,
      avoid_ingredients: [...form.avoid_ingredients],
      favorite_ingredients: [],
      cook_time_preference: form.cook_time_preference,
      default_servings: form.people_count,
      people_count: form.people_count,
    } as any)
    showToast('偏好已保存')
    if (isGuideMode.value) {
      setTimeout(() => router.replace('/'), 450)
    }
  } catch {
    showToast('保存失败，请稍后再试')
  } finally {
    saving.value = false
  }
}

function skipGuide() {
  router.replace('/')
}

function handleBack() {
  if (isGuideMode.value) {
    router.replace('/')
    return
  }
  router.back()
}

onMounted(() => {
  loadPreferences()
})

onUnmounted(() => {
  if (toastTimer) clearTimeout(toastTimer)
})
</script>

<style scoped>
.preferences-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.8);
  --coral: #e95645;
  --coral-2: #ef5548;
  --brown: #7a5942;
  --brown-soft: #ffebd2;
  --sage: #8fa783;
  --border: rgba(255, 255, 255, 0.62);
  --shadow: var(--card-shadow);
  position: relative;
  width: min(100%, 430px);
  min-height: calc(100vh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 160px);
  min-height: calc(100dvh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 160px);
  margin: 0 auto;
  padding: max(18px, env(safe-area-inset-top)) 24px calc(38px + env(safe-area-inset-bottom));
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 237, 205, 0.32), rgba(255, 247, 233, 0.2)),
    var(--preference-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.preferences-shell::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 5%, rgba(255, 255, 255, 0.74), transparent 29%),
    radial-gradient(circle at 88% 16%, rgba(236, 143, 71, 0.22), transparent 33%),
    radial-gradient(circle at 12% 92%, rgba(233, 86, 69, 0.16), transparent 31%),
    linear-gradient(90deg, rgba(255, 239, 214, 0.55), rgba(255, 245, 230, 0.16) 55%, rgba(172, 91, 33, 0.18));
  backdrop-filter: blur(4px) saturate(1.12);
  -webkit-backdrop-filter: blur(4px) saturate(1.12);
}

.status-bar,
.page-header,
.settings-card,
.save-tip {
  position: relative;
  z-index: 1;
}

button,
input {
  font: inherit;
  -webkit-tap-highlight-color: transparent;
}

button {
  border: 0;
  cursor: pointer;
}

svg {
  display: block;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.status-bar {
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px;
  color: #1e1713;
  font-size: 17px;
  font-weight: 850;
  line-height: 1;
}

.status-icons {
  display: flex;
  align-items: center;
  gap: 7px;
}

.cell-bars {
  display: inline-flex;
  align-items: flex-end;
  gap: 3px;
  height: 18px;
}

.cell-bars i {
  width: 4px;
  border-radius: 999px;
  background: currentColor;
}

.cell-bars i:nth-child(1) { height: 7px; }
.cell-bars i:nth-child(2) { height: 10px; }
.cell-bars i:nth-child(3) { height: 13px; }
.cell-bars i:nth-child(4) { height: 16px; }

.wifi {
  width: 23px;
  height: 17px;
}

.wifi path {
  stroke-width: 2.7;
}

.battery {
  position: relative;
  width: 30px;
  height: 16px;
  border: 2px solid currentColor;
  border-radius: 5px;
}

.battery::before {
  content: "";
  position: absolute;
  right: -5px;
  top: 4px;
  width: 3px;
  height: 6px;
  border-radius: 0 2px 2px 0;
  background: currentColor;
}

.battery::after {
  content: "";
  position: absolute;
  left: 3px;
  top: 3px;
  width: 20px;
  height: 6px;
  border-radius: 2px;
  background: currentColor;
}

.page-header {
  position: relative;
  height: 96px;
  display: grid;
  place-items: center;
  text-align: center;
}

.back-btn,
.save-btn {
  position: absolute;
  top: 14px;
  height: 52px;
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 16px;
  background: rgba(255, 250, 240, 0.86);
  box-shadow:
    0 12px 28px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.back-btn {
  left: 0;
  width: 52px;
  display: grid;
  place-items: center;
  color: #4a352a;
}

.back-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.55;
}

.save-btn {
  right: 0;
  min-width: 74px;
  padding: 0 18px;
  color: var(--coral);
  font-size: 17px;
  font-weight: 880;
  letter-spacing: 0;
}

.save-btn:disabled {
  opacity: 0.58;
  cursor: not-allowed;
}

.page-title {
  margin: 12px 0 0;
  color: var(--text);
  font-size: 28px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
}

.settings-card {
  margin-top: 22px;
  padding: 24px 22px;
  overflow: hidden;
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius);
  background: var(--card-surface);
  box-shadow: var(--shadow);
  backdrop-filter: blur(var(--card-blur));
  -webkit-backdrop-filter: blur(var(--card-blur));
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 11px;
  color: var(--text);
}

.card-title-row svg {
  width: 25px;
  height: 25px;
  flex: 0 0 25px;
  color: var(--coral);
  stroke-width: 2.4;
}

.card-title {
  margin: 0;
  color: var(--text);
  font-size: 25px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
}

.card-desc {
  margin: 13px 0 0;
  color: var(--sub);
  font-size: 16px;
  font-weight: 580;
  line-height: 1.55;
}

.taste-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-top: 24px;
}

.taste-chip,
.goal-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(255, 255, 255, 0.42);
  color: #6b5142;
  background: rgba(255, 248, 236, 0.58);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.54);
  font-weight: 850;
  white-space: nowrap;
  transition: transform 180ms ease, background 180ms ease, color 180ms ease, box-shadow 180ms ease;
}

.taste-chip {
  min-width: 0;
  height: 58px;
  border-radius: var(--card-radius-inner);
  font-size: 18px;
}

.taste-chip.active {
  border-color: transparent;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.24);
}

.goal-scroll {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  overflow-x: auto;
  scrollbar-width: none;
}

.goal-scroll::-webkit-scrollbar {
  display: none;
}

.goal-chip {
  height: 54px;
  flex: 0 0 auto;
  padding: 0 20px;
  border-radius: var(--card-radius-inner);
  font-size: 17px;
}

.goal-chip.active {
  border-color: transparent;
  color: #fff;
  background: var(--brown);
  box-shadow: 0 12px 24px rgba(90, 60, 40, 0.2);
}

.avoid-input {
  height: 62px;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 22px;
  padding: 0 16px 0 18px;
  border: 1px solid rgba(120, 90, 65, 0.16);
  border-radius: var(--card-radius-inner);
  background: rgba(255, 255, 255, 0.4);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.avoid-input input {
  width: 0;
  flex: 1 1 auto;
  border: 0;
  outline: 0;
  color: var(--text);
  background: transparent;
  font-size: 16px;
  font-weight: 650;
}

.avoid-input input::placeholder {
  color: var(--muted);
  font-weight: 650;
}

.add-ingredient {
  width: 34px;
  height: 34px;
  display: grid;
  flex: 0 0 34px;
  place-items: center;
  border: 1px solid rgba(120, 90, 65, 0.24);
  border-radius: 999px;
  color: #9a7957;
  background: rgba(255, 250, 240, 0.34);
  transition: transform 180ms ease, color 180ms ease;
}

.add-ingredient svg {
  width: 19px;
  height: 19px;
  stroke-width: 2.4;
}

.avoid-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 16px;
}

.ingredient-tag {
  height: 42px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  border-radius: 14px;
  color: #7a4b25;
  background: var(--brown-soft);
  font-size: 16px;
  font-weight: 760;
  line-height: 1;
}

.ingredient-tag button {
  width: 20px;
  height: 20px;
  display: grid;
  place-items: center;
  margin-right: -4px;
  color: #7a4b25;
  background: transparent;
}

.ingredient-tag svg {
  width: 17px;
  height: 17px;
  stroke-width: 2.6;
}

.stepper {
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 24px;
  padding: 0 22px;
  border: 1px solid rgba(120, 90, 65, 0.14);
  border-radius: var(--card-radius-inner);
  background: rgba(255, 255, 255, 0.34);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.52);
}

.step-btn {
  width: 48px;
  height: 48px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: #5a4032;
  background: rgba(255, 248, 236, 0.76);
  font-size: 26px;
  font-weight: 850;
  line-height: 1;
  transition: transform 180ms ease, opacity 180ms ease;
}

.step-btn:disabled {
  opacity: 0.42;
  cursor: not-allowed;
}

.serving-count {
  color: var(--text);
  font-size: 40px;
  font-weight: 950;
  line-height: 1;
}

.save-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 22px 0 0;
  color: var(--sub);
  font-size: 15px;
  font-weight: 650;
  text-align: center;
}

.save-tip svg {
  width: 23px;
  height: 23px;
  flex: 0 0 23px;
  color: var(--sage);
  stroke-width: 2.2;
}

.skip-btn {
  width: 100%;
  min-height: 54px;
  margin-top: 18px;
  border: 1px solid rgba(120, 90, 65, 0.16);
  border-radius: var(--card-radius-inner);
  color: #6b5142;
  background: rgba(255, 250, 240, 0.72);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.72);
  font-size: 16px;
  font-weight: 850;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: calc(34px + env(safe-area-inset-bottom));
  z-index: 6;
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

.back-btn:active,
.save-btn:active,
.taste-chip:active,
.goal-chip:active,
.add-ingredient:active,
.ingredient-tag button:active,
.step-btn:active,
.skip-btn:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .back-btn:hover,
  .save-btn:hover,
  .taste-chip:hover,
  .goal-chip:hover,
  .add-ingredient:hover,
  .step-btn:hover:not(:disabled),
  .skip-btn:hover {
    transform: translateY(-1px);
  }
}

@media (max-width: 380px) {
  .preferences-shell {
    padding-right: 18px;
    padding-left: 18px;
  }

  .settings-card {
    padding: 23px 18px;
    border-radius: var(--card-radius);
  }

  .taste-grid {
    gap: 12px;
  }

  .taste-chip {
    height: 56px;
    font-size: 17px;
  }

  .goal-chip {
    padding: 0 18px;
  }
}

@media (max-width: 350px) {
  .taste-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .page-title {
    font-size: 26px;
  }

  .save-btn {
    min-width: 68px;
    padding: 0 15px;
  }
}

@media (min-width: 431px) {
  .preferences-shell {
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
