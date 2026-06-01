<template>
  <div class="page preferences-page">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">偏好设置</h1>
      <button class="save-btn" @click="handleSave" :disabled="saving">
        {{ saving ? '保存中...' : '保存' }}
      </button>
    </header>

    <div class="form-section">
      <!-- 口味偏好 -->
      <div class="form-group">
        <label class="form-label">口味偏好</label>
        <div class="tag-grid">
          <button
            v-for="t in tasteOptions"
            :key="t"
            class="tag-btn"
            :class="{ active: form.taste_preference.includes(t) }"
            @click="toggleTaste(t)"
          >
            {{ t }}
          </button>
        </div>
      </div>

      <!-- 饮食目标 -->
      <div class="form-group">
        <label class="form-label">饮食目标</label>
        <div class="tag-grid">
          <button
            v-for="g in goalOptions"
            :key="g.value"
            class="tag-btn"
            :class="{ active: form.health_goal === g.value }"
            @click="form.health_goal = g.value"
          >
            {{ g.label }}
          </button>
        </div>
      </div>

      <!-- 忌口食材 -->
      <div class="form-group">
        <label class="form-label">忌口食材</label>
        <div class="tag-input-wrap">
          <div class="tag-list">
            <span v-for="(item, idx) in form.avoid_ingredients" :key="idx" class="tag-chip">
              {{ item }}
              <button class="tag-remove" @click="form.avoid_ingredients.splice(idx, 1)">×</button>
            </span>
          </div>
          <input
            v-model="avoidInput"
            class="tag-input"
            placeholder="输入食材后回车添加"
            @keyup.enter="addAvoid"
          />
        </div>
      </div>

      <!-- 常用人数 -->
      <div class="form-group">
        <label class="form-label">常用用餐人数</label>
        <div class="people-picker">
          <button class="people-btn" @click="form.people_count = Math.max(1, form.people_count - 1)">-</button>
          <span class="people-num">{{ form.people_count }}</span>
          <button class="people-btn" @click="form.people_count = Math.min(10, form.people_count + 1)">+</button>
        </div>
      </div>
    </div>

    <div v-if="showToast" class="toast">保存成功</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getPreferences, updatePreferences } from '@/api/user'

const saving = ref(false)
const showToast = ref(false)
const avoidInput = ref('')

const tasteOptions = ['咸鲜', '酸甜', '麻辣', '清淡', '香辣', '酸辣', '甜味', '原味']
const goalOptions = [
  { value: 'normal', label: '普通' },
  { value: 'lose', label: '减脂' },
  { value: 'gain', label: '增肌' },
  { value: 'sugar', label: '控糖' },
  { value: 'child', label: '儿童营养' },
]

const form = reactive({
  taste_preference: [] as string[],
  health_goal: 'normal',
  avoid_ingredients: [] as string[],
  people_count: 2,
})

function toggleTaste(t: string) {
  const idx = form.taste_preference.indexOf(t)
  if (idx >= 0) {
    form.taste_preference.splice(idx, 1)
  } else {
    form.taste_preference.push(t)
  }
}

function addAvoid() {
  const val = avoidInput.value.trim()
  if (val && !form.avoid_ingredients.includes(val)) {
    form.avoid_ingredients.push(val)
  }
  avoidInput.value = ''
}

async function handleSave() {
  saving.value = true
  try {
    await updatePreferences(form as any)
    showToast.value = true
    setTimeout(() => { showToast.value = false }, 2000)
  } catch {
    // ignore
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const res: any = await getPreferences()
    if (res) {
      form.taste_preference = res.taste_preference || []
      form.health_goal = res.health_goal || 'normal'
      form.avoid_ingredients = res.avoid_ingredients || []
      form.people_count = res.people_count || 2
    }
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.preferences-page {
  min-height: 100vh;
  background: var(--color-bg);
  padding-bottom: var(--sp-8);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
}

.back-btn {
  width: 34px; height: 34px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.back-btn:hover { background: var(--color-surface-2); }
.back-btn:active { background: var(--color-surface-3); }
.back-btn svg { width: 18px; height: 18px; }

.page-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

.save-btn {
  padding: var(--sp-2) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.save-btn:hover { background: var(--color-surface-2); }
.save-btn:active { background: var(--color-surface-3); }
.save-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.form-section {
  padding: 0 var(--sp-4);
}

.form-group {
  margin-bottom: var(--sp-6);
}

.form-label {
  display: block;
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--sp-3);
}

.tag-grid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--sp-2);
}

.tag-btn {
  padding: var(--sp-2) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease);
}

.tag-btn:hover {
  border-color: var(--color-border-med);
}

.tag-btn.active {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
}

.tag-input-wrap {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  padding: var(--sp-3);
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--sp-2);
  margin-bottom: var(--sp-2);
}

.tag-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--sp-1);
  padding: 2px 8px;
  background: var(--color-accent-soft);
  color: var(--color-accent);
  border-radius: var(--r-sm);
  font-size: var(--text-xs);
  font-weight: 600;
}

.tag-remove {
  background: none;
  border: none;
  color: var(--color-accent);
  font-size: 14px;
  cursor: pointer;
  padding: 0 2px;
}

.tag-input {
  width: 100%;
  border: none;
  background: transparent;
  color: var(--color-text);
  font-size: var(--text-sm);
  outline: none;
}

.tag-input::placeholder {
  color: var(--color-text-3);
}

.people-picker {
  display: flex;
  align-items: center;
  gap: var(--sp-4);
}

.people-btn {
  width: 36px;
  height: 36px;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-lg);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.people-btn:hover { background: var(--color-surface-2); }
.people-btn:active { background: var(--color-surface-3); }

.people-num {
  font-family: var(--font-outlier);
  font-size: var(--text-xl);
  font-weight: 750;
  color: var(--color-text);
  min-width: 32px;
  text-align: center;
}

.toast {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: var(--sp-4) var(--sp-6);
  background: var(--color-text);
  color: var(--color-text-inv);
  border-radius: var(--r-md);
  font-size: var(--text-sm);
  font-weight: 600;
  z-index: 100;
  animation: fadeInOut 1.5s ease;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
  15% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  85% { opacity: 1; }
  100% { opacity: 0; }
}

@media (min-width: 768px) {
  .preferences-page { max-width: 640px; margin: 0 auto; }
}
</style>
