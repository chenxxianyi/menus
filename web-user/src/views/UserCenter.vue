<template>
  <div class="page user-page">
    <!-- Profile Header -->
    <div class="profile-header anim-delay-1">
      <div class="profile-avatar">{{ userInitial }}</div>
      <div class="profile-info">
        <h2 class="profile-name">{{ userName }}</h2>
        <p class="profile-sub">加入食刻推荐第 128 天</p>
      </div>
      <svg class="edit-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
      </svg>
    </div>

    <!-- Stats Card -->
    <div class="stats-card anim-delay-2">
      <div class="stat-item">
        <div class="stat-num">{{ favoriteCount }}</div>
        <div class="stat-label">我的收藏</div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-num">5</div>
        <div class="stat-label">定制菜单</div>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <div class="stat-num">3</div>
        <div class="stat-label">购物清单</div>
      </div>
    </div>

    <!-- Menu List -->
    <div class="menu-list anim-delay-3">
      <div class="menu-row" @click="$router.push('/user/favorites')">
        <div class="menu-icon red">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
        </div>
        <span class="menu-label">我的收藏</span>
        <svg class="menu-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m9 18 6-6-6-6"/>
        </svg>
      </div>
      <div class="menu-row" @click="$router.push('/week-menu')">
        <div class="menu-icon blue">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
            <polyline points="10 9 9 9 8 9"/>
          </svg>
        </div>
        <span class="menu-label">我的菜单</span>
        <svg class="menu-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m9 18 6-6-6-6"/>
        </svg>
      </div>
      <div class="menu-row" @click="$router.push('/shopping-list')">
        <div class="menu-icon orange">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="9" cy="21" r="1"/>
            <circle cx="20" cy="21" r="1"/>
            <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
          </svg>
        </div>
        <span class="menu-label">我的购物清单</span>
        <svg class="menu-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m9 18 6-6-6-6"/>
        </svg>
      </div>
    </div>

    <!-- Settings List -->
    <div class="menu-list anim-delay-4">
      <div class="menu-row" @click="$router.push('/user/preferences')">
        <div class="menu-icon gray">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </div>
        <span class="menu-label">饮食偏好设置</span>
        <svg class="menu-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m9 18 6-6-6-6"/>
        </svg>
      </div>
      <div class="menu-row">
        <div class="menu-icon green">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
        </div>
        <span class="menu-label">浏览历史</span>
        <svg class="menu-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m9 18 6-6-6-6"/>
        </svg>
      </div>
      <div class="menu-row">
        <div class="menu-icon gray">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="16" x2="12" y2="12"/>
            <line x1="12" y1="8" x2="12.01" y2="8"/>
          </svg>
        </div>
        <span class="menu-label">关于我们</span>
        <svg class="menu-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m9 18 6-6-6-6"/>
        </svg>
      </div>
    </div>

    <button class="logout-btn anim-delay-5" @click="handleLogout">退出登录</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getFavoriteCount } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()

const favoriteCount = ref(0)

const userName = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '用户')
const userInitial = computed(() => userName.value.charAt(0))

function handleLogout() {
  userStore.logout()
  router.replace('/login')
}

onMounted(async () => {
  // 获取用户信息
  if (!userStore.userInfo) {
    await userStore.fetchUserInfo()
  }
  // 获取收藏数
  try {
    const res: any = await getFavoriteCount()
    favoriteCount.value = res?.count || 0
  } catch {
    favoriteCount.value = 0
  }
})
</script>

<style scoped>
.user-page {
  padding-top: var(--sp-8);
  padding-bottom: var(--sp-8);
}

/* ── Profile Header ── */
.profile-header {
  display: flex;
  align-items: center;
  gap: var(--sp-4);
  margin-bottom: var(--sp-8);
}

.profile-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--color-surface);
  border: 4px solid var(--color-surface);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text);
  flex-shrink: 0;
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: var(--text-lg);
  font-weight: 700;
  color: var(--color-text);
}

.profile-sub {
  font-size: var(--text-xs);
  color: var(--color-text-3);
  margin-top: 4px;
}

.edit-icon {
  width: 24px;
  height: 24px;
  color: var(--color-text-3);
  cursor: pointer;
}

/* ── Stats Card ── */
.stats-card {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background: var(--color-surface);
  border-radius: var(--r-lg);
  padding: var(--sp-4);
  margin-bottom: var(--sp-6);
  box-shadow: var(--shadow-sm);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-num {
  font-size: var(--text-lg);
  font-weight: 700;
  color: var(--color-text);
}

.stat-label {
  font-size: 10px;
  color: var(--color-text-3);
}

.stat-divider {
  width: 1px;
  height: 32px;
  background: var(--color-border);
}

/* ── Menu List ── */
.menu-list {
  background: var(--color-surface);
  border-radius: var(--r-lg);
  overflow: hidden;
  margin-bottom: var(--sp-4);
  box-shadow: var(--shadow-sm);
}

.menu-row {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  padding: var(--sp-4);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.menu-row:hover {
  background: var(--color-surface-2);
}

.menu-row:active {
  background: var(--color-surface-3);
}

.menu-row + .menu-row {
  border-top: 1px solid var(--color-border);
}

.menu-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.menu-icon svg {
  width: 20px;
  height: 20px;
}

.menu-icon.red {
  color: #EF4444;
}

.menu-icon.blue {
  color: #3B82F6;
}

.menu-icon.orange {
  color: #F97316;
}

.menu-icon.gray {
  color: #6B7280;
}

.menu-icon.green {
  color: #10B981;
}

.menu-label {
  flex: 1;
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text);
}

.menu-chevron {
  width: 16px;
  height: 16px;
  color: var(--color-text-3);
}

/* ── Logout ── */
.logout-btn {
  width: 100%;
  padding: var(--sp-3);
  border-radius: var(--r-md);
  border: none;
  background: transparent;
  font-size: var(--text-sm);
  font-weight: 500;
  color: #EF4444;
  cursor: pointer;
  margin-top: var(--sp-4);
  transition: background var(--dur-fast) var(--ease);
}

.logout-btn:hover {
  background: #FEE2E2;
}

.logout-btn:active {
  background: #FECACA;
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .page {
    max-width: 640px;
  }
}
</style>
