<template>
  <div class="page couple-home">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">情侣点餐</h1>
      <div style="width:34px"></div>
    </header>

    <!-- Not bound state -->
    <div v-if="!coupleStore.isBound" class="empty-state">
      <div class="empty-icon">&#128150;</div>
      <p class="empty-text">还没有绑定情侣关系</p>
      <button class="primary-btn" @click="$router.push('/couple/bind')">去绑定</button>
    </div>

    <!-- Bound state -->
    <template v-else>
      <!-- Partner card -->
      <section class="partner-card">
        <div class="partner-avatars">
          <div class="avatar-circle">我</div>
          <div class="heart-icon">&#10084;</div>
          <div class="avatar-circle partner">{{ partnerInitial }}</div>
        </div>
        <h2 class="couple-name">{{ coupleStore.coupleInfo?.couple_name || coupleName }}</h2>
        <p class="partner-nickname">你的另一半：{{ coupleStore.coupleInfo?.partner?.nickname }}</p>
      </section>

      <!-- Quick actions -->
      <section class="quick-actions">
        <button class="action-card" @click="$router.push('/couple/order')">
          <span class="action-icon">&#127860;</span>
          <span class="action-label">我要点餐</span>
          <span class="action-sub">告诉TA你想吃什么</span>
        </button>
        <button class="action-card" @click="$router.push('/couple/orders')">
          <span class="action-icon">&#128203;</span>
          <span class="action-label">点餐清单</span>
          <span class="action-sub">看看双方点的菜</span>
        </button>
        <button class="action-card" @click="$router.push('/couple/menu')">
          <span class="action-icon">&#128722;</span>
          <span class="action-label">合意菜单</span>
          <span class="action-sub">合并生成食材清单</span>
        </button>
      </section>

      <!-- Recent orders -->
      <section class="recent-section" v-if="recentOrders.length">
        <div class="section-header">
          <h3 class="section-title">最近点餐</h3>
          <button class="see-all" @click="$router.push('/couple/orders')">查看全部</button>
        </div>
        <div class="order-list">
          <div v-for="order in recentOrders" :key="order.id" class="order-item">
            <div class="order-left">
              <span class="order-meal">{{ mealLabel(order.meal_type) }}</span>
              <span class="order-dish">{{ order.dish_name }}</span>
            </div>
            <div class="order-right">
              <span class="order-date">{{ order.meal_date }}</span>
              <span class="order-status" :class="statusClass(order.status)">{{ statusLabel(order.status) }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Settings -->
      <section class="settings-section">
        <button class="setting-row" @click="showNameDialog = true">
          <span>设置情侣昵称</span>
          <svg class="setting-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </button>
        <button class="setting-row danger" @click="handleUnbind">
          <span>解除绑定</span>
        </button>
      </section>
    </template>

    <!-- Name dialog -->
    <div v-if="showNameDialog" class="dialog-overlay" @click.self="showNameDialog = false">
      <div class="dialog-box">
        <h3 class="dialog-title">设置情侣昵称</h3>
        <input v-model="nameInput" class="dialog-input" placeholder="如：小明&小红" maxlength="20" />
        <div class="dialog-actions">
          <button class="dialog-cancel" @click="showNameDialog = false">取消</button>
          <button class="dialog-confirm" @click="saveName">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCoupleStore } from '@/stores/couple'
import { setCoupleName, unbindCouple } from '@/api/couple'

const router = useRouter()
const coupleStore = useCoupleStore()

const showNameDialog = ref(false)
const nameInput = ref('')

const partnerInitial = computed(() => {
  const name = coupleStore.coupleInfo?.partner?.nickname || 'TA'
  return name.charAt(0)
})

const coupleName = computed(() => {
  const info = coupleStore.coupleInfo
  if (!info?.couple_name) return '甜蜜情侣'
  return info.couple_name
})

const recentOrders = computed(() => {
  return coupleStore.orders.slice(0, 5)
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

async function saveName() {
  if (!nameInput.value.trim()) return
  try {
    await setCoupleName(nameInput.value.trim())
    await coupleStore.fetchCoupleInfo()
    showNameDialog.value = false
  } catch {}
}

async function handleUnbind() {
  if (!confirm('确定要解除情侣绑定吗？')) return
  try {
    await unbindCouple()
    coupleStore.clear()
    router.replace('/couple/bind')
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
.couple-home {
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

/* Empty state */
.empty-state {
  text-align: center;
  padding: var(--sp-16) var(--sp-4);
}

.empty-icon { font-size: 64px; margin-bottom: var(--sp-4); }

.empty-text {
  color: var(--color-text-3);
  font-size: var(--text-sm);
  margin-bottom: var(--sp-6);
}

.primary-btn {
  padding: var(--sp-3) var(--sp-8);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-base);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease),
              transform var(--dur-fast) var(--ease),
              box-shadow var(--dur-fast) var(--ease);
}

.primary-btn:hover {
  background: var(--color-text-2);
  box-shadow: var(--shadow-md);
}

.primary-btn:active {
  transform: translateY(1px);
  box-shadow: none;
}

/* Partner card */
.partner-card {
  margin: var(--sp-4);
  padding: var(--sp-6);
  background: var(--color-surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--color-border);
  text-align: center;
}

.partner-avatars {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-4);
  margin-bottom: var(--sp-4);
}

.avatar-circle {
  width: 52px; height: 52px;
  border-radius: 50%;
  background: var(--color-text);
  color: var(--color-text-inv);
  display: flex; align-items: center; justify-content: center;
  font-size: var(--text-lg);
  font-weight: 700;
}

.avatar-circle.partner {
  background: var(--color-accent);
}

.heart-icon {
  font-size: 24px;
  color: var(--color-error);
  animation: heartbeat 1.5s infinite;
}

@keyframes heartbeat {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15); }
}

.couple-name {
  font-size: var(--text-lg);
  font-weight: 650;
  color: var(--color-text);
  font-family: var(--font-display);
  margin-bottom: var(--sp-1);
}

.partner-nickname {
  color: var(--color-text-3);
  font-size: var(--text-xs);
}

/* Quick actions */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--sp-3);
  padding: 0 var(--sp-4);
  margin-bottom: var(--sp-6);
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--sp-2);
  padding: var(--sp-5) var(--sp-2);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  cursor: pointer;
  transition: transform var(--dur-base) var(--ease-out),
              box-shadow var(--dur-base) var(--ease-out),
              border-color var(--dur-fast) var(--ease);
}

.action-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-border-med);
}

.action-card:active {
  transform: translateY(0);
  box-shadow: none;
}

.action-icon { font-size: 28px; }

.action-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

.action-sub {
  font-size: var(--text-2xs);
  color: var(--color-text-3);
  text-align: center;
}

/* Recent orders */
.recent-section {
  padding: 0 var(--sp-4);
  margin-bottom: var(--sp-6);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-3);
}

.section-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

.see-all {
  background: none;
  border: none;
  color: var(--color-text-3);
  font-size: var(--text-xs);
  cursor: pointer;
}

.order-list {
  background: var(--color-surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.order-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
  border-bottom: 1px solid var(--color-border);
}

.order-item:last-child { border-bottom: none; }

.order-left {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.order-meal {
  font-size: var(--text-2xs);
  color: var(--color-text-3);
  padding: 2px 6px;
  background: var(--color-surface-2);
  border-radius: 4px;
  font-weight: 600;
}

.order-dish {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
}

.order-right {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
}

.order-date {
  font-size: var(--text-2xs);
  color: var(--color-text-3);
}

.order-status {
  font-size: var(--text-2xs);
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
}

.order-status.pending { background: var(--color-warning-soft); color: var(--color-warning); }
.order-status.confirmed { background: var(--color-success-soft); color: var(--color-success); }
.order-status.cancelled { background: var(--color-error-soft); color: var(--color-error); }

/* Settings */
.settings-section {
  padding: 0 var(--sp-4);
}

.setting-row {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  margin-bottom: var(--sp-2);
  color: var(--color-text);
  font-size: var(--text-sm);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease),
              border-color var(--dur-fast) var(--ease);
}

.setting-row:hover {
  background: var(--color-surface-2);
  border-color: var(--color-border-med);
}

.setting-row:active {
  background: var(--color-surface-3);
}

.setting-row.danger { color: var(--color-error); }

.setting-arrow {
  width: 16px; height: 16px;
  color: var(--color-text-3);
}

/* Dialog */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(28, 25, 21, 0.55);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.dialog-box {
  width: 80%;
  max-width: 320px;
  background: var(--color-surface);
  border-radius: var(--r-lg);
  padding: var(--sp-6);
}

.dialog-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--sp-4);
  text-align: center;
}

.dialog-input {
  width: 100%;
  padding: var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface-2);
  color: var(--color-text);
  font-size: var(--text-sm);
  outline: none;
  margin-bottom: var(--sp-4);
}

.dialog-input:focus { border-color: var(--color-text); }

.dialog-actions {
  display: flex;
  gap: var(--sp-3);
}

.dialog-cancel, .dialog-confirm {
  flex: 1;
  padding: var(--sp-3);
  border-radius: var(--r-sm);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.dialog-cancel {
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  color: var(--color-text-2);
}

.dialog-confirm {
  background: var(--color-text);
  border: none;
  color: var(--color-text-inv);
}

@media (min-width: 768px) {
  .couple-home {
    max-width: 640px;
    margin: 0 auto;
  }
}
</style>
