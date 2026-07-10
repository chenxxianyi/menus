<template>
  <div class="page couple-order">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">发布想吃</h1>
      <div style="width:34px"></div>
    </header>

    <div class="order-form">
      <section class="flow-card">
        <p class="flow-kicker">给做饭的人一个明确请求</p>
        <h2 class="flow-title">说出想吃什么，系统会自动匹配菜谱</h2>
        <p class="flow-desc">对方确认后，可以直接查看做法并生成采购清单。</p>
      </section>

      <!-- Dish name -->
      <div class="form-group">
        <label class="form-label">想吃什么</label>
        <input
          v-model="form.dish_name"
          class="form-input"
          type="text"
          placeholder="输入菜名，如：番茄炒蛋、红烧肉"
          maxlength="100"
        />
      </div>

      <!-- Meal type -->
      <div class="form-group">
        <label class="form-label">餐次</label>
        <div class="meal-options">
          <button
            v-for="m in mealTypes"
            :key="m.value"
            class="meal-btn"
            :class="{ active: form.meal_type === m.value }"
            @click="form.meal_type = m.value"
          >
            {{ m.label }}
          </button>
        </div>
      </div>

      <!-- Meal date -->
      <div class="form-group">
        <label class="form-label">想吃的日期</label>
        <input
          v-model="form.meal_date"
          class="form-input"
          type="date"
        />
      </div>

      <!-- Note -->
      <div class="form-group">
        <label class="form-label">备注（可选）</label>
        <textarea
          v-model="form.note"
          class="form-textarea"
          placeholder="如：少辣、多放醋、不要香菜..."
          maxlength="200"
          rows="3"
        ></textarea>
      </div>

      <!-- Quick recipes -->
      <div class="form-group" v-if="recipes.length">
        <label class="form-label">指定菜谱（可选）</label>
        <p class="form-hint">不选择也可以，提交后会按菜名自动匹配一份菜谱。</p>
        <div class="recipe-grid">
          <button
            v-for="r in recipes"
            :key="r.id"
            class="recipe-option"
            :class="{ selected: form.recipe_id === r.id }"
            @click="selectRecipe(r)"
          >
            <span class="recipe-option-name">{{ r.title }}</span>
            <span class="recipe-option-meta">{{ r.cook_time }}分钟 · {{ r.difficulty }}</span>
          </button>
        </div>
      </div>

      <!-- Submit -->
      <button class="submit-btn" @click="handleSubmit" :disabled="!form.dish_name.trim() || submitting">
        {{ submitting ? '发布中...' : '发布给对方' }}
      </button>
    </div>

    <!-- Success toast -->
    <div v-if="showSuccess" class="toast">已发布，等待对方确认</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createCoupleOrder } from '@/api/couple'
import { getRecipes } from '@/api/recipe'
import { formatLocalDate } from '@/utils/date'

const router = useRouter()

const mealTypes = [
  { value: 'breakfast', label: '早餐' },
  { value: 'lunch', label: '午餐' },
  { value: 'dinner', label: '晚餐' },
  { value: 'snack', label: '夜宵' },
]

const today = formatLocalDate()

const form = reactive({
  dish_name: '',
  meal_type: 'dinner',
  meal_date: today,
  note: '',
  recipe_id: undefined as number | undefined,
})

const recipes = ref<any[]>([])
const submitting = ref(false)
const showSuccess = ref(false)

function selectRecipe(r: any) {
  form.dish_name = r.title
  form.recipe_id = r.id
}

async function handleSubmit() {
  if (!form.dish_name.trim()) return
  submitting.value = true
  try {
    await createCoupleOrder({
      dish_name: form.dish_name.trim(),
      recipe_id: form.recipe_id,
      meal_type: form.meal_type,
      meal_date: form.meal_date,
      note: form.note.trim(),
    })
    showSuccess.value = true
    setTimeout(() => {
      showSuccess.value = false
      router.push('/couple/orders')
    }, 1200)
  } catch {
    // ignore
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const res: any = await getRecipes({ page: 1, page_size: 8 })
    recipes.value = res?.list || res || []
  } catch {}
})
</script>

<style scoped>
.couple-order {
  min-height: 100vh;
  background: var(--color-bg);
  padding-bottom: var(--sp-16);
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
}

.back-btn svg { width: 18px; height: 18px; }

.page-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

.order-form {
  padding: 0 var(--sp-4);
}

.flow-card {
  margin-bottom: var(--sp-5);
  padding: var(--sp-5);
  border: 1px solid var(--glass-border);
  border-radius: var(--r-lg);
  background:
    linear-gradient(135deg, var(--color-surface), var(--color-broth-soft));
  box-shadow: var(--glass-shadow);
}

.flow-kicker {
  margin-bottom: var(--sp-1);
  color: var(--color-accent);
  font-size: var(--text-2xs);
  font-weight: 800;
}

.flow-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 700;
  line-height: 1.25;
}

.flow-desc {
  margin-top: var(--sp-2);
  color: var(--color-text-2);
  font-size: var(--text-sm);
  line-height: 1.6;
}

.form-group {
  margin-bottom: var(--sp-5);
}

.form-label {
  display: block;
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--sp-2);
}

.form-hint {
  margin: calc(-1 * var(--sp-1)) 0 var(--sp-3);
  color: var(--color-text-3);
  font-size: var(--text-xs);
  line-height: 1.5;
}

.form-input {
  width: 100%;
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-base);
  outline: none;
  transition: border-color var(--dur-fast) var(--ease);
}

.form-input:focus { border-color: var(--color-text); }

.form-textarea {
  width: 100%;
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-sm);
  outline: none;
  resize: vertical;
  font-family: inherit;
  transition: border-color var(--dur-fast) var(--ease);
}

.form-textarea:focus { border-color: var(--color-text); }

.meal-options {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--sp-2);
}

.meal-btn {
  padding: var(--sp-3) 0;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-3);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease);
}

.meal-btn.active {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
}

.recipe-grid {
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.recipe-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease);
}

.recipe-option.selected {
  border-color: var(--color-text);
  background: var(--color-surface-2);
}

.recipe-option-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

.recipe-option-meta {
  font-size: var(--text-2xs);
  color: var(--color-text-3);
}

.submit-btn {
  position: fixed;
  bottom: var(--sp-4);
  left: var(--sp-4);
  right: var(--sp-4);
  padding: var(--sp-4);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-base);
  font-weight: 600;
  cursor: pointer;
  z-index: 10;
  transition: opacity var(--dur-fast) var(--ease);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
  animation: fadeInOut 1.2s ease;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
  20% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  80% { opacity: 1; }
  100% { opacity: 0; }
}

@media (min-width: 768px) {
  .couple-order { max-width: 640px; margin: 0 auto; }
  .submit-btn { left: 50%; right: auto; transform: translateX(-50%); max-width: 608px; width: calc(100% - 32px); }
}
</style>
