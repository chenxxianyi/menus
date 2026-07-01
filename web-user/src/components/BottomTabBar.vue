<template>
  <nav class="tabbar" aria-label="主导航">
    <router-link
      v-for="tab in tabs"
      :key="tab.path"
      :to="tab.path"
      class="tab"
      :class="{ active: isActive(tab.path) }"
      :aria-current="isActive(tab.path) ? 'page' : undefined"
    >
      <div class="tab-icon">
        <svg v-if="tab.icon === 'home'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 10.5 12 3l9 7.5"/>
          <path d="M5.5 9.5V21h13V9.5"/>
          <path d="M9.5 21v-6h5v6"/>
        </svg>
        <svg v-else-if="tab.icon === 'book'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M5 4.5A2.5 2.5 0 0 1 7.5 2H20v18H7.5A2.5 2.5 0 0 0 5 22V4.5Z"/>
          <path d="M5 4.5A2.5 2.5 0 0 0 2.5 2H4"/>
          <path d="M8 7h8M8 11h6"/>
        </svg>
        <svg v-else-if="tab.icon === 'menu'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="6" y="3" width="12" height="18" rx="2"/>
          <path d="M9 8h6M9 12h6M9 16h4"/>
          <path d="M8 3V1.8M16 3V1.8"/>
        </svg>
        <svg v-else-if="tab.icon === 'cart'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="9" cy="20" r="1"/>
          <circle cx="18" cy="20" r="1"/>
          <path d="M3 4h2l2.4 10.4a2 2 0 0 0 2 1.6h7.8a2 2 0 0 0 1.9-1.4L21 8H6"/>
        </svg>
        <svg v-else-if="tab.icon === 'user'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="7.5" r="4"/>
          <path d="M4.5 21a7.5 7.5 0 0 1 15 0"/>
        </svg>
      </div>
      <span class="tab-label">{{ tab.label }}</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const tabs = [
  { path: '/', icon: 'home', label: '首页' },
  { path: '/recipes', icon: 'book', label: '菜谱' },
  { path: '/week-menu', icon: 'menu', label: '菜单推荐' },
  { path: '/shopping-list', icon: 'cart', label: '购物清单' },
  { path: '/user', icon: 'user', label: '我的' },
]

function isActive(path: string) {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<style scoped>
.tabbar {
  position: fixed;
  bottom: max(16px, var(--safe-bottom));
  left: 50%;
  transform: translateX(-50%);
  width: min(calc(100% - 58px), 374px);
  height: 74px;
  display: flex;
  align-items: center;
  padding: 8px 13px;
  border: 1px solid rgba(255, 255, 255, 0.84);
  border-radius: 999px;
  background: rgba(255, 249, 238, 0.84);
  box-shadow:
    0 20px 42px rgba(77, 44, 22, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(22px) saturate(1.14);
  -webkit-backdrop-filter: blur(22px) saturate(1.14);
  z-index: 100;
}

.tab {
  flex: 1;
  min-width: 0;
  min-height: 44px;
  height: 58px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  border-radius: 999px;
  color: #897b70;
  text-decoration: none;
  transition:
    color 180ms ease,
    background 180ms ease,
    transform 160ms ease,
    box-shadow 180ms ease;
  -webkit-tap-highlight-color: transparent;
}

.tab.active {
  color: #df4438;
}

.tab:active {
  transform: translateY(1px);
}

.tab-icon {
  width: 24px;
  height: 24px;
  display: grid;
  place-items: center;
}

.tab-icon svg {
  width: 24px;
  height: 24px;
}

.tab-label {
  font-size: 12px;
  font-weight: 700;
  line-height: 1.1;
  letter-spacing: 0;
  white-space: nowrap;
}

@media (max-width: 360px) {
  .tabbar {
    width: calc(100% - 34px);
  }

  .tab-label {
    font-size: 10px;
  }
}
</style>
