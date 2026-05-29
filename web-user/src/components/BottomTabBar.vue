<template>
  <nav class="tabbar" aria-label="主导航">
    <router-link
      v-for="tab in tabs"
      :key="tab.path"
      :to="tab.path"
      class="tab"
      :class="{ active: isActive(tab.path) }"
    >
      <div class="tab-icon">
        <svg v-if="tab.icon === 'home'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
          <polyline points="9 22 9 12 15 12 15 22"/>
        </svg>
        <svg v-else-if="tab.icon === 'calendar'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
          <line x1="16" y1="2" x2="16" y2="6"/>
          <line x1="8" y1="2" x2="8" y2="6"/>
          <line x1="3" y1="10" x2="21" y2="10"/>
        </svg>
        <svg v-else-if="tab.icon === 'chart'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="20" x2="18" y2="10"/>
          <line x1="12" y1="20" x2="12" y2="4"/>
          <line x1="6" y1="20" x2="6" y2="14"/>
        </svg>
        <svg v-else-if="tab.icon === 'user'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
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
  { path: '/week-menu', icon: 'calendar', label: '一周菜单' },
  { path: '/stats', icon: 'chart', label: '营养' },
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
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - var(--sp-8));
  max-width: calc(var(--max-w) - var(--sp-8));
  height: 56px;
  display: flex;
  align-items: center;
  background: var(--glass-bg);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-border);
  border-radius: var(--r-md);
  box-shadow: var(--glass-shadow);
  padding: 0 var(--sp-2);
  margin-bottom: var(--sp-3);
  z-index: 100;
}

.tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: var(--sp-1) 0;
  border-radius: var(--r-sm);
  transition: color var(--dur-base) var(--ease);
  color: var(--color-text-3);
  text-decoration: none;
  -webkit-tap-highlight-color: transparent;
  position: relative;
}

.tab:hover {
  color: var(--color-text-2);
}

.tab.active {
  color: var(--color-text);
}

.tab-label {
  position: relative;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0;
}

.tab.active .tab-label::after {
  content: '';
  position: absolute;
  left: 0;
  bottom: -4px;
  width: 100%;
  height: 2px;
  border-radius: 1px;
  background: var(--color-accent);
}

.tab:active { transform: translateY(1px); }

.tab-icon {
  width: 22px;
  height: 22px;
}

.tab-icon svg {
  width: 22px;
  height: 22px;
}
</style>
