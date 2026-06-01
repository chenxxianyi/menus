<template>
  <div class="page couple-orders">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">点餐清单</h1>
      <button class="add-btn" @click="$router.push('/couple/order')">+</button>
    </header>

    <!-- Date filter -->
    <div class="date-filter">
      <button
        v-for="d in dateOptions"
        :key="d.value"
        class="date-btn"
        :class="{ active: selectedDate === d.value }"
        @click="selectedDate = d.value"
      >
        {{ d.label }}
      </button>
    </div>

    <!-- Orders list -->
    <div v-if="filteredOrders.length" class="orders-list">
      <div v-for="order in filteredOrders" :key="order.id" class="order-card">
        <div class="order-header">
          <span class="order-author">{{ order.user?.nickname || '我' }}</span>
          <span class="order-meal-tag">{{ mealLabel(order.meal_type) }}</span>
          <span class="order-status" :class="statusClass(order.status)">{{ statusLabel(order.status) }}</span>
        </div>
        <h3 class="order-dish-name">{{ order.dish_name }}</h3>
        <p v-if="order.note" class="order-note">{{ order.note }}</p>
        <div class="order-actions">
          <button v-if="order.status === 0" class="action-btn confirm" @click="handleStatus(order.id, 1)">确认</button>
          <button v-if="order.status === 0" class="action-btn cancel" @click="handleStatus(order.id, 2)">取消</button>
          <button v-if="order.status === 1" class="action-btn cancel" @click="handleStatus(order.id, 0)">撤回</button>
          <button class="action-btn delete" @click="handleDelete(order.id)">删除</button>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <div class="empty-icon">&#127860;</div>
      <p class="empty-text">还没有点餐记录</p>
      <button class="primary-btn" @click="$router.push('/couple/order')">去点餐</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCoupleStore } from '@/stores/couple'
import { updateCoupleOrderStatus, deleteCoupleOrder } from '@/api/couple'

const router = useRouter()
const coupleStore = useCoupleStore()

const selectedDate = ref('')

const today = new Date()
const dateOptions = computed(() => {
  const opts = [{ value: '', label: '全部' }]
  for (let i = 0; i < 7; i++) {
    const d = new Date(today)
    d.setDate(d.getDate() + i)
    const dateStr = d.toISOString().split('T')[0]
    const label = i === 0 ? '今天' : i === 1 ? '明天' : `${d.getMonth() + 1}/${d.getDate()}`
    opts.push({ value: dateStr, label })
  }
  return opts
})

const filteredOrders = computed(() => {
  if (!selectedDate.value) return coupleStore.orders
  return coupleStore.orders.filter(o => o.meal_date === selectedDate.value)
})

function mealLabel(type: string) {
  const map: Record<string, string> = { breakfast: '早餐', lunch: '午餐', dinner: '晚餐', snack: '夜宵' }
  return map[type] || type
}

function statusLabel(status: number) {
  const map: Record<number, string> = { 0: '待确认', 1: '已确认', 2: '已取消' }
  return map[status] || ''
}

function statusClass(status: number) {
  const map: Record<number, string> = { 0: 'pending', 1: 'confirmed', 2: 'cancelled' }
  return map[status] || ''
}

async function handleStatus(id: number, status: number) {
  try {
    await updateCoupleOrderStatus(id, status)
    await coupleStore.fetchOrders()
  } catch {}
}

async function handleDelete(id: number) {
  if (!confirm('确定删除这条点餐吗？')) return
  try {
    await deleteCoupleOrder(id)
    await coupleStore.fetchOrders()
  } catch {}
}

onMounted(async () => {
  await coupleStore.fetchCoupleInfo()
  if (!coupleStore.isBound) {
    router.replace('/couple/bind')
    return
  }
  await coupleStore.fetchOrders()
})
</script>

<style scoped>
.couple-orders {
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

.add-btn {
  width: 34px; height: 34px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-lg);
  font-weight: 600;
  cursor: pointer;
}

/* Date filter */
.date-filter {
  display: flex;
  gap: var(--sp-2);
  padding: 0 var(--sp-4);
  margin-bottom: var(--sp-4);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.date-filter::-webkit-scrollbar { display: none; }

.date-btn {
  flex-shrink: 0;
  padding: var(--sp-2) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease);
}

.date-btn.active {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
}

/* Orders */
.orders-list {
  padding: 0 var(--sp-4);
}

.order-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-4);
  margin-bottom: var(--sp-3);
}

.order-header {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  margin-bottom: var(--sp-2);
}

.order-author {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text);
}

.order-meal-tag {
  font-size: var(--text-2xs);
  padding: 2px 6px;
  background: var(--color-surface-2);
  border-radius: 4px;
  color: var(--color-text-3);
  font-weight: 600;
}

.order-status {
  margin-left: auto;
  font-size: var(--text-2xs);
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
}

.order-status.pending { background: var(--color-warning-soft); color: var(--color-warning); }
.order-status.confirmed { background: var(--color-success-soft); color: var(--color-success); }
.order-status.cancelled { background: var(--color-error-soft); color: var(--color-error); }

.order-dish-name {
  font-size: var(--text-base);
  font-weight: 650;
  color: var(--color-text);
  margin-bottom: var(--sp-1);
}

.order-note {
  font-size: var(--text-xs);
  color: var(--color-text-3);
  margin-bottom: var(--sp-3);
}

.order-actions {
  display: flex;
  gap: var(--sp-2);
  margin-top: var(--sp-3);
  padding-top: var(--sp-3);
  border-top: 1px solid var(--color-border);
}

.action-btn {
  padding: var(--sp-2) var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  font-size: var(--text-2xs);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease);
}

.action-btn.confirm {
  background: var(--color-success-soft);
  border-color: var(--color-success);
  color: var(--color-success);
}

.action-btn.cancel {
  background: var(--color-warning-soft);
  border-color: var(--color-warning);
  color: var(--color-warning);
}

.action-btn.delete {
  background: var(--color-error-soft);
  border-color: var(--color-error);
  color: var(--color-error);
  margin-left: auto;
}

/* Empty */
.empty-state {
  text-align: center;
  padding: var(--sp-12) var(--sp-4);
}

.empty-icon { font-size: 48px; margin-bottom: var(--sp-3); }

.empty-text {
  color: var(--color-text-3);
  font-size: var(--text-sm);
  margin-bottom: var(--sp-4);
}

.primary-btn {
  padding: var(--sp-3) var(--sp-8);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

@media (min-width: 768px) {
  .couple-orders { max-width: 640px; margin: 0 auto; }
}
</style>
