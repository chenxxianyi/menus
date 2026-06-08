<template>
  <div class="page couple-menu">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">合意菜单</h1>
      <div style="width:34px"></div>
    </header>

    <!-- Filters -->
    <div class="filter-section">
      <div class="filter-row">
        <label class="filter-label">日期</label>
        <input v-model="mealDate" class="filter-date" type="date" />
      </div>
      <div class="filter-row">
        <label class="filter-label">餐次</label>
        <div class="meal-options">
          <button
            v-for="m in mealTypes"
            :key="m.value"
            class="meal-btn"
            :class="{ active: mealType === m.value }"
            @click="mealType = m.value"
          >
            {{ m.label }}
          </button>
        </div>
      </div>
      <button class="generate-btn" @click="handleGenerate" :disabled="loading">
        {{ loading ? '生成中...' : '生成食材清单' }}
      </button>
    </div>

    <!-- Result -->
    <div v-if="result" class="result-section">
      <!-- Orders summary -->
      <div class="orders-summary">
        <h3 class="summary-title">已确认的点餐（{{ result.orders.length }}）</h3>
        <div v-for="order in result.orders" :key="order.id" class="summary-item">
          <div class="summary-main">
            <span class="summary-dish">{{ order.dish_name }}</span>
            <span class="summary-note" v-if="order.note">{{ order.note }}</span>
            <button v-if="order.recipe" class="guide-link" @click="router.push(`/recipes/${order.recipe.id}`)">
              查看做饭指南
            </button>
          </div>
          <span v-if="order.recipe" class="summary-recipe">{{ order.recipe.title }}</span>
        </div>
        <div v-if="!result.orders.length" class="summary-empty">暂无已确认的点餐</div>
      </div>

      <!-- Shopping list -->
      <div class="shopping-section" v-if="result.shopping_list.length">
        <div class="shopping-header">
          <h3 class="shopping-title">食材清单</h3>
          <span class="shopping-count">共 {{ result.total_items }} 种</span>
        </div>
        <div class="shopping-list">
          <div v-for="(item, idx) in result.shopping_list" :key="idx" class="shopping-item">
            <span class="item-check">
              <input type="checkbox" v-model="item.checked" />
            </span>
            <span class="item-name" :class="{ checked: item.checked }">{{ item.name }}</span>
            <span class="item-amount">{{ item.amount }}</span>
          </div>
        </div>
        <button class="save-btn" @click="handleSaveToShoppingList">
          保存到购物清单
        </button>
      </div>

      <div v-else class="empty-shopping">
        <p>没有找到需要购买的食材</p>
        <p class="empty-hint">请先发布想吃的菜，并确认一条已匹配菜谱的点餐。</p>
      </div>
    </div>

    <!-- Initial state -->
    <div v-else-if="!loading" class="initial-state">
      <div class="initial-icon">&#128722;</div>
      <p class="initial-text">选择日期和餐次，合并两人的点餐生成食材清单</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCoupleStore } from '@/stores/couple'
import { generateShoppingList } from '@/api/couple'
import { createShoppingList } from '@/api/shopping'
import type { GenerateShoppingListResult, ShoppingListItem } from '@/types/couple'

const router = useRouter()
const coupleStore = useCoupleStore()

const today = new Date().toISOString().split('T')[0]
const mealDate = ref(today)
const mealType = ref('')
const loading = ref(false)
type CheckedShoppingItem = ShoppingListItem & { checked: boolean }
type CoupleMenuResult = Omit<GenerateShoppingListResult, 'shopping_list'> & { shopping_list: CheckedShoppingItem[] }
const result = ref<CoupleMenuResult | null>(null)

const mealTypes = [
  { value: '', label: '全部' },
  { value: 'breakfast', label: '早餐' },
  { value: 'lunch', label: '午餐' },
  { value: 'dinner', label: '晚餐' },
]

async function handleGenerate() {
  loading.value = true
  try {
    const res: any = await generateShoppingList(mealDate.value, mealType.value)
    result.value = {
      ...res,
      shopping_list: (res.shopping_list || []).map((item: ShoppingListItem) => ({ ...item, checked: false })),
    }
  } catch {
    result.value = null
  } finally {
    loading.value = false
  }
}

async function handleSaveToShoppingList() {
  if (!result.value) return
  const items = result.value.shopping_list.map(item => ({
    name: item.name,
    amount: item.amount,
    emoji: '',
    category: item.category || '',
    price: 0,
    checked: item.checked,
  }))
  try {
    await createShoppingList({
      name: `${mealDate.value} 情侣点餐食材`,
      items,
    })
    alert('已保存到购物清单')
  } catch {
    alert('保存失败')
  }
}

onMounted(async () => {
  await coupleStore.fetchCoupleInfo()
  if (!coupleStore.isBound) {
    router.replace('/couple/bind')
  }
})
</script>

<style scoped>
.couple-menu {
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
}

.back-btn svg { width: 18px; height: 18px; }

.page-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

/* Filters */
.filter-section {
  margin: 0 var(--sp-4) var(--sp-4);
  padding: var(--sp-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
}

.filter-row {
  margin-bottom: var(--sp-3);
}

.filter-label {
  display: block;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-3);
  margin-bottom: var(--sp-2);
}

.filter-date {
  width: 100%;
  padding: var(--sp-2) var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface-2);
  color: var(--color-text);
  font-size: var(--text-sm);
  outline: none;
}

.filter-date:focus { border-color: var(--color-text); }

.meal-options {
  display: flex;
  gap: var(--sp-2);
}

.meal-btn {
  flex: 1;
  padding: var(--sp-2) 0;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease);
}

.meal-btn.active {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
}

.generate-btn {
  width: 100%;
  padding: var(--sp-3);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  margin-top: var(--sp-2);
  transition: opacity var(--dur-fast) var(--ease);
}

.generate-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Result */
.result-section {
  padding: 0 var(--sp-4);
}

.orders-summary {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-4);
  margin-bottom: var(--sp-4);
}

.summary-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--sp-3);
}

.summary-item {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  padding: var(--sp-2) 0;
  border-bottom: 1px solid var(--color-border);
}

.summary-item:last-child { border-bottom: none; }

.summary-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.summary-dish {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

.summary-note {
  font-size: var(--text-2xs);
  color: var(--color-text-3);
}

.summary-recipe {
  max-width: 42%;
  margin-left: auto;
  flex-shrink: 0;
  color: var(--color-accent);
  font-size: var(--text-2xs);
  font-weight: 700;
  text-align: right;
}

.guide-link {
  align-self: flex-start;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--color-tomato);
  font-size: var(--text-2xs);
  font-weight: 800;
  cursor: pointer;
}

.summary-empty {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  text-align: center;
  padding: var(--sp-3) 0;
}

/* Shopping list */
.shopping-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-4);
  margin-bottom: var(--sp-4);
}

.shopping-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-3);
}

.shopping-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

.shopping-count {
  font-size: var(--text-2xs);
  color: var(--color-text-3);
}

.shopping-list {
  margin-bottom: var(--sp-4);
}

.shopping-item {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  padding: var(--sp-2) 0;
  border-bottom: 1px solid var(--color-border);
}

.shopping-item:last-child { border-bottom: none; }

.item-check input {
  width: 18px;
  height: 18px;
  accent-color: var(--color-text);
  cursor: pointer;
}

.item-name {
  flex: 1;
  font-size: var(--text-sm);
  color: var(--color-text);
  transition: all var(--dur-fast) var(--ease);
}

.item-name.checked {
  text-decoration: line-through;
  color: var(--color-text-3);
}

.item-amount {
  font-size: var(--text-xs);
  color: var(--color-text-3);
}

.save-btn {
  width: 100%;
  padding: var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.save-btn:active {
  background: var(--color-surface-2);
}

/* Empty states */
.empty-shopping, .initial-state {
  text-align: center;
  padding: var(--sp-12) var(--sp-4);
}

.empty-shopping p, .initial-text {
  color: var(--color-text-3);
  font-size: var(--text-sm);
}

.empty-hint {
  font-size: var(--text-xs) !important;
  margin-top: var(--sp-2);
}

.initial-icon {
  font-size: 48px;
  margin-bottom: var(--sp-3);
}

@media (min-width: 768px) {
  .couple-menu { max-width: 640px; margin: 0 auto; }
}
</style>
