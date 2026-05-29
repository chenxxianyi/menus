<template>
  <div class="page">
    <!-- Header -->
    <header class="header anim-delay-1">
      <div style="width:36px;"></div>
      <div style="flex:1;text-align:center;">
        <div class="header-title">我的</div>
      </div>
      <button class="btn-ghost" aria-label="设置">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </button>
    </header>

    <!-- Profile -->
    <div class="profile anim-delay-2">
      <div class="profile-avatar">{{ userInitial }}</div>
      <h2 class="profile-name">{{ userName }}</h2>
      <p class="profile-desc">{{ userStore.userInfo?.username || '' }}</p>
    </div>

    <!-- Stats -->
    <div class="stats-row anim-delay-2">
      <div class="stat-item">
        <div class="stat-num">128</div>
        <div class="stat-lbl">收藏</div>
      </div>
      <div class="stat-sep"></div>
      <div class="stat-item">
        <div class="stat-num">56</div>
        <div class="stat-lbl">做过</div>
      </div>
      <div class="stat-sep"></div>
      <div class="stat-item">
        <div class="stat-num">14</div>
        <div class="stat-lbl">打卡</div>
      </div>
    </div>

    <!-- Favorites group -->
    <div class="list-group anim-delay-3">
      <div class="list-row">
        <div class="icon-circle icon-circle--orange">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">我的收藏</div>
        </div>
        <span class="list-row-badge">128 道</span>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
      <div class="list-row">
        <div class="icon-circle icon-circle--purple">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">我的创作</div>
        </div>
        <span class="list-row-badge">12 道</span>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
      <div class="list-row">
        <div class="icon-circle icon-circle--blue">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">最近浏览</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
    </div>

    <div style="height:var(--sp-4);"></div>

    <!-- Settings group -->
    <div class="list-group anim-delay-4">
      <div class="list-row">
        <div class="icon-circle icon-circle--green">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">个人资料</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
      <div class="list-row">
        <div class="icon-circle icon-circle--yellow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/><path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">口味偏好</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
      <div class="list-row">
        <div class="icon-circle icon-circle--orange">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8h1a4 4 0 0 1 0 8h-1"/><path d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">过敏源设置</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
      <div class="list-row">
        <div class="icon-circle icon-circle--red">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">隐私与安全</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
    </div>

    <div style="height:var(--sp-4);"></div>

    <!-- About group -->
    <div class="list-group anim-delay-5">
      <div class="list-row">
        <div class="icon-circle icon-circle--blue">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">帮助与反馈</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
      <div class="list-row">
        <div class="icon-circle icon-circle--purple">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"/><polyline points="16 6 12 2 8 6"/><line x1="12" y1="2" x2="12" y2="15"/></svg>
        </div>
        <div class="list-row-body">
          <div class="list-row-title">分享给朋友</div>
        </div>
        <svg class="list-row-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
      </div>
    </div>

    <button class="logout-btn anim-delay-5" @click="handleLogout">退出登录</button>

    <p class="version">今天吃什么 v1.0.0</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const userName = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '用户')
const userInitial = computed(() => userName.value.charAt(0))

function handleLogout() {
  userStore.logout()
  router.replace('/login')
}
</script>

<style scoped>
/* ── Profile ── */
.profile {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--sp-6) 0 var(--sp-5);
}

.profile-avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: var(--color-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 800;
  font-size: 28px;
  margin-bottom: var(--sp-3);
}

.profile-name {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 700;
  color: var(--color-ink);
}

.profile-desc {
  font-size: var(--text-sm);
  color: var(--color-ink-3);
  margin-top: var(--sp-1);
}

/* ── Stats ── */
.stats-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-8);
  padding: var(--sp-4) 0;
  margin-bottom: var(--sp-6);
  border-top: 0.5px solid var(--color-rule);
  border-bottom: 0.5px solid var(--color-rule);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
}

.stat-num {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 700;
  color: var(--color-ink);
}

.stat-lbl {
  font-size: 11px;
  color: var(--color-ink-3);
  font-weight: 500;
}

.stat-sep {
  width: 1px;
  height: 24px;
  background: var(--color-rule);
}

/* ── List extras ── */
.list-row-badge {
  padding: 2px 8px;
  border-radius: var(--r-full);
  font-size: 11px;
  font-weight: 600;
  background: var(--color-orange-soft);
  color: var(--color-accent);
}

.list-row-chevron {
  width: 16px;
  height: 16px;
  color: var(--color-ink-3);
  flex-shrink: 0;
}

.list-row-body { flex: 1; min-width: 0; }
.list-row-title { font-size: var(--text-base); font-weight: 500; color: var(--color-ink); }

/* ── Logout ── */
.logout-btn {
  width: 100%;
  padding: var(--sp-3);
  border-radius: var(--r-lg);
  border: 1px solid var(--color-rule);
  background: var(--color-surface);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-red);
  cursor: pointer;
  margin-top: var(--sp-4);
  transition: background var(--dur-fast) var(--ease-out);
}

.logout-btn:active { background: var(--color-paper-2); }

/* ── Version ── */
.version {
  text-align: center;
  font-size: 11px;
  color: var(--color-ink-3);
  padding: var(--sp-6) 0 var(--sp-4);
}
</style>
