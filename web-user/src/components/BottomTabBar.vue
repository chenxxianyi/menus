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
        <svg v-else-if="tab.icon === 'book'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/>
          <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>
        </svg>
        <svg v-else-if="tab.icon === 'utensils'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/>
          <path d="M7 2v20"/>
          <path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3zm0 0v7"/>
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
  { path: '/recipes', icon: 'book', label: '菜谱' },
  { path: '/week-menu', icon: 'utensils', label: '菜单推荐' },
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
  bottom: max(12px, var(--safe-bottom));
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 24px);
  max-width: calc(var(--max-w) - 24px);
  height: 72px;
  display: flex;
  align-items: center;
  padding: 8px;
  border: 1px solid var(--glass-border);
  border-radius: 24px;
  background: var(--glass-bg);
  box-shadow: 0 18px 44px rgba(31, 41, 55, 0.16);
  backdrop-filter: blur(var(--glass-blur));
  -webkit-backdrop-filter: blur(var(--glass-blur));
  z-index: 100;
}

.tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  min-width: 0;
  height: 56px;
  border-radius: 18px;
  padding: 0;
  transition:
    background var(--dur-base) var(--ease-spring),
    color var(--dur-base) var(--ease),
    transform var(--dur-fast) var(--ease),
    box-shadow var(--dur-base) var(--ease);
  color: var(--color-text-3);
  text-decoration: none;
  -webkit-tap-highlight-color: transparent;
  cursor: pointer;
}

.tab:hover {
  color: var(--color-text-2);
}

.tab.active {
  background: var(--color-surface);
  color: var(--color-accent);
  box-shadow: var(--shadow-sm);
}

.tab-label {
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 0;
}

.tab:active {
  transform: translateY(1px);
}

.tab-icon {
  width: 22px;
  height: 22px;
}

.tab-icon svg {
  width: 22px;
  height: 22px;
}
</style>
