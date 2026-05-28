<template>
  <div class="page-container">
    <!-- Header -->
    <header class="app-header animate-fade-up">
      <div style="width:40px;"></div>
      <div class="app-header-center">
        <div class="app-header-title">我的</div>
      </div>
      <button class="btn-icon" aria-label="设置">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </button>
    </header>

    <!-- Profile Header -->
    <div class="profile-header animate-fade-up anim-delay-1">
      <div class="profile-avatar">王</div>
      <div class="profile-name">小王</div>
      <div class="profile-desc">美食爱好者 · 厨龄 3 年</div>
    </div>

    <!-- Stats -->
    <div class="profile-stats animate-fade-up anim-delay-2">
      <div class="profile-stat glass-card-sm">
        <div class="profile-stat-value">128</div>
        <div class="profile-stat-label">收藏菜谱</div>
      </div>
      <div class="profile-stat glass-card-sm">
        <div class="profile-stat-value">56</div>
        <div class="profile-stat-label">做过菜品</div>
      </div>
      <div class="profile-stat glass-card-sm">
        <div class="profile-stat-value">14</div>
        <div class="profile-stat-label">连续打卡</div>
      </div>
    </div>

    <!-- Favorites Menu -->
    <div class="menu-card glass-card anim-delay-3">
      <div class="menu-card-title">我的收藏</div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-orange">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">我的收藏</div>
        </div>
        <div class="menu-card-badge">128 道</div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-purple">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">我的创作</div>
        </div>
        <div class="menu-card-badge">12 道</div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-blue">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">最近浏览</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
    </div>

    <!-- Settings Menu -->
    <div class="menu-card glass-card anim-delay-4">
      <div class="menu-card-title">设置</div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-green">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">个人资料</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-yellow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/>
            <path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">口味偏好</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-orange">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 8h1a4 4 0 0 1 0 8h-1"/>
            <path d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">过敏源设置</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-red">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">隐私与安全</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
    </div>

    <!-- Help Menu -->
    <div class="menu-card glass-card anim-delay-5">
      <div class="menu-card-title">关于</div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-blue">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
            <line x1="12" y1="17" x2="12.01" y2="17"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">帮助与反馈</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
      <div class="glass-list-item">
        <div class="glass-icon glass-icon-purple">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"/>
            <polyline points="16 6 12 2 8 6"/>
            <line x1="12" y1="2" x2="12" y2="15"/>
          </svg>
        </div>
        <div class="menu-card-text">
          <div class="menu-card-name">分享给朋友</div>
        </div>
        <div class="menu-card-arrow">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"/></svg>
        </div>
      </div>
    </div>

    <div class="version-text anim-delay-6">今天吃什么 v1.0.0</div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
/* Profile Header */
.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 0 24px;
  text-align: center;
}

.profile-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ffb347, #f59e0b);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 800;
  font-size: 32px;
  margin-bottom: 12px;
  box-shadow: 0 8px 24px rgba(245,158,11,0.25), inset 0 2px 0 rgba(255,255,255,0.3);
  border: 3px solid rgba(255,255,255,0.5);
}

.profile-name {
  font-size: 22px;
  font-weight: 800;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.profile-desc {
  font-size: 14px;
  color: #a0a0a0;
  margin-top: 4px;
}

/* Stats Row */
.profile-stats {
  display: flex;
  gap: 10px;
  margin-bottom: 24px;
}

.profile-stat {
  flex: 1;
  border-radius: 16px;
  padding: 14px 8px;
  text-align: center;
}

.profile-stat-value {
  font-size: 22px;
  font-weight: 800;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.profile-stat-label {
  font-size: 11px;
  color: #a0a0a0;
  font-weight: 500;
  margin-top: 2px;
}

/* Menu Card */
.menu-card {
  overflow: hidden;
  margin-bottom: 16px;
}

.menu-card-title {
  font-size: 13px;
  color: #a0a0a0;
  font-weight: 600;
  padding: 14px 18px 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.menu-card-text { flex: 1; }
.menu-card-name { font-size: 15px; font-weight: 600; color: #1a1a1a; }

.menu-card-arrow { color: #a0a0a0; flex-shrink: 0; }
.menu-card-arrow svg { width: 18px; height: 18px; }

.menu-card-badge {
  padding: 3px 10px;
  border-radius: 9999px;
  font-size: 11px;
  font-weight: 600;
  background: rgba(255,224,168,0.5);
  backdrop-filter: blur(6px);
  border: 1px solid rgba(255,255,255,0.3);
  color: #f59e0b;
}

/* Version */
.version-text {
  text-align: center;
  font-size: 12px;
  color: #a0a0a0;
  padding: 20px 0 8px;
}
</style>
