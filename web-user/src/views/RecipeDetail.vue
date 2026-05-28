<template>
  <div class="page-container">
    <!-- Hero Banner -->
    <div class="recipe-hero animate-fade-up">
      <div class="recipe-hero-badge">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14l-5-4.87 6.91-1.01L12 2z"/>
        </svg>
        今日推荐
      </div>
      <button class="recipe-hero-back" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div class="recipe-hero-emoji">🍖</div>
    </div>

    <!-- Title -->
    <div class="anim-delay-1" style="margin-bottom: 4px;">
      <h1 style="font-size: 26px; font-weight: 800; letter-spacing: -0.8px; color: #1a1a1a;">红烧排骨</h1>
      <p style="font-size: 14px; color: #a0a0a0; margin-top: 4px;">经典家常菜，软烂入味，色泽红亮</p>
    </div>

    <!-- Info Row -->
    <div class="recipe-info anim-delay-2">
      <div class="recipe-info-item glass-card-sm">
        <div class="recipe-info-icon glass-icon-orange">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
        </div>
        <div class="recipe-info-value">40 分钟</div>
        <div class="recipe-info-label">烹饪时间</div>
      </div>
      <div class="recipe-info-item glass-card-sm">
        <div class="recipe-info-icon glass-icon-blue">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
        </div>
        <div class="recipe-info-value">3 人份</div>
        <div class="recipe-info-label">份量</div>
      </div>
      <div class="recipe-info-item glass-card-sm">
        <div class="recipe-info-icon glass-icon-green">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20V10"/><path d="M18 20V4"/><path d="M6 20v-4"/></svg>
        </div>
        <div class="recipe-info-value">简单</div>
        <div class="recipe-info-label">难度</div>
      </div>
      <div class="recipe-info-item glass-card-sm">
        <div class="recipe-info-icon glass-icon-purple">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/>
            <path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/>
          </svg>
        </div>
        <div class="recipe-info-value">520 kcal</div>
        <div class="recipe-info-label">每份热量</div>
      </div>
    </div>

    <!-- Nutrition -->
    <div class="anim-delay-3" style="margin-bottom: 20px;">
      <div class="section-title">营养成分</div>
      <div class="nutrition-bar glass-card">
        <div class="nutrition-item">
          <div class="nutrition-circle">
            <svg><circle class="track" cx="28" cy="28" r="25"/><circle class="fill-orange" cx="28" cy="28" r="25"/></svg>
            <span class="nutrition-value">35g</span>
          </div>
          <div class="nutrition-label">蛋白质</div>
        </div>
        <div class="nutrition-item">
          <div class="nutrition-circle">
            <svg><circle class="track" cx="28" cy="28" r="25"/><circle class="fill-blue" cx="28" cy="28" r="25"/></svg>
            <span class="nutrition-value">28g</span>
          </div>
          <div class="nutrition-label">脂肪</div>
        </div>
        <div class="nutrition-item">
          <div class="nutrition-circle">
            <svg><circle class="track" cx="28" cy="28" r="25"/><circle class="fill-green" cx="28" cy="28" r="25"/></svg>
            <span class="nutrition-value">42g</span>
          </div>
          <div class="nutrition-label">碳水</div>
        </div>
        <div class="nutrition-item">
          <div class="nutrition-circle">
            <svg><circle class="track" cx="28" cy="28" r="25"/><circle class="fill-purple" cx="28" cy="28" r="25"/></svg>
            <span class="nutrition-value">8g</span>
          </div>
          <div class="nutrition-label">膳食纤维</div>
        </div>
      </div>
    </div>

    <!-- Ingredients -->
    <div class="anim-delay-3" style="margin-bottom: 24px;">
      <div class="section-title">食材清单</div>
      <div class="glass-list">
        <div v-for="ing in ingredients" :key="ing.name" class="glass-list-item" @click="ing.checked = !ing.checked">
          <div class="glass-list-icon" :class="`glass-icon-${ing.color}`">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path v-if="ing.icon === 'leaf'" d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.58-3.25 3.93L12 22"/>
              <path v-if="ing.icon === 'leaf'" d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.58 3.25 3.93"/>
              <circle v-if="ing.icon === 'circle'" cx="12" cy="12" r="10"/>
              <rect v-if="ing.icon === 'rect'" x="3" y="11" width="18" height="11" rx="2"/>
              <path v-if="ing.icon === 'rect'" d="M7 11V7a5 5 0 0 1 10 0v4"/>
              <path v-if="ing.icon === 'cup'" d="M18 8h1a4 4 0 0 1 0 8h-1"/>
              <path v-if="ing.icon === 'cup'" d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"/>
              <path v-if="ing.icon === 'layers'" d="M12 2L2 7l10 5 10-5-10-5z"/>
              <path v-if="ing.icon === 'layers'" d="M2 17l10 5 10-5"/>
            </svg>
          </div>
          <div class="glass-list-text">
            <div class="glass-list-title">{{ ing.name }}</div>
            <div class="glass-list-desc">{{ ing.desc }}</div>
          </div>
          <div class="glass-list-check" :class="{ checked: ing.checked }">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Steps -->
    <div class="anim-delay-4" style="margin-bottom: 24px;">
      <div class="section-title">烹饪步骤</div>
      <div class="step-list">
        <div v-for="(step, idx) in steps" :key="idx" class="step-item">
          <div class="step-num">{{ idx + 1 }}</div>
          <div class="step-content glass-card-sm">
            <div class="step-title">{{ step.title }}</div>
            <div class="step-desc">{{ step.desc }}</div>
            <div v-if="step.tip" class="step-tip">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg>
              {{ step.tip }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="recipe-actions anim-delay-5">
      <button class="btn-primary">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
        加入今日菜单
      </button>
      <button class="btn-secondary" :class="{ liked: isLiked }" @click="isLiked = !isLiked">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
        </svg>
      </button>
      <button class="btn-secondary">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/>
          <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isLiked = ref(false)

const ingredients = ref([
  { name: '排骨', desc: '500g，选肋排口感更佳', color: 'orange', icon: 'leaf', checked: true },
  { name: '冰糖', desc: '30g，用于炒糖色', color: 'green', icon: 'circle', checked: true },
  { name: '生抽 + 老抽', desc: '生抽 2 勺，老抽 1 勺', color: 'blue', icon: 'rect', checked: false },
  { name: '料酒', desc: '1 勺，去腥提鲜', color: 'purple', icon: 'cup', checked: false },
  { name: '姜片 + 八角 + 桂皮', desc: '各适量', color: 'yellow', icon: 'layers', checked: true },
])

const steps = [
  { title: '焯水去腥', desc: '排骨冷水下锅，加姜片和料酒，大火煮开后撇去浮沫，捞出沥干备用。', tip: '冷水下锅能更好去除血水' },
  { title: '炒糖色', desc: '锅中放少许油，加入冰糖小火慢炒，待糖色变成琥珀色时下排骨翻炒均匀。', tip: '' },
  { title: '调味炖煮', desc: '加入生抽、老抽、料酒、八角、桂皮，倒入没过排骨的热水，大火烧开后转小火炖 30 分钟。', tip: '一定要加热水，冷水会让肉质收紧' },
  { title: '收汁装盘', desc: '大火收汁至浓稠，翻炒均匀让每块排骨都裹上酱汁，撒上葱花即可出锅。', tip: '' },
]
</script>

<style scoped>
/* Hero Banner */
.recipe-hero {
  position: relative;
  border-radius: 24px;
  overflow: hidden;
  height: 220px;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #fce4ec 0%, #f8bbd0 50%, #f48fb1 100%);
}

.recipe-hero::after {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 1;
  background: linear-gradient(180deg, rgba(255,255,255,0.15) 0%, rgba(0,0,0,0.25) 100%);
  pointer-events: none;
}

.recipe-hero-emoji {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 80px;
  z-index: 2;
  filter: drop-shadow(0 4px 12px rgba(0,0,0,0.1));
}

.recipe-hero-badge {
  position: absolute;
  top: 16px;
  left: 16px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(12px) saturate(1.3);
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
  color: white;
  border: 1px solid rgba(255,255,255,0.25);
}

.recipe-hero-back {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 2;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: white;
  border: none;
}

.recipe-hero-back svg { width: 18px; height: 18px; }

/* Recipe Info Row */
.recipe-info {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.recipe-info-item {
  flex: 1;
  border-radius: 14px;
  padding: 12px 8px;
  text-align: center;
}

.recipe-info-icon {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 6px;
}

.recipe-info-icon svg { width: 16px; height: 16px; }

.recipe-info-value { font-size: 15px; font-weight: 700; color: #1a1a1a; }
.recipe-info-label { font-size: 11px; color: #a0a0a0; font-weight: 500; margin-top: 1px; }

/* Nutrition */
.nutrition-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  padding: 16px;
}

.nutrition-item {
  flex: 1;
  text-align: center;
}

.nutrition-circle {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  margin: 0 auto 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.nutrition-circle svg {
  width: 56px;
  height: 56px;
  transform: rotate(-90deg);
}

.nutrition-circle circle {
  fill: none;
  stroke-width: 4;
  stroke-linecap: round;
}

.nutrition-circle .track { stroke: rgba(0,0,0,0.06); }
.nutrition-circle .fill-orange { stroke: #ffb347; stroke-dasharray: 157; stroke-dashoffset: 47; }
.nutrition-circle .fill-blue { stroke: #a8d8ea; stroke-dasharray: 157; stroke-dashoffset: 63; }
.nutrition-circle .fill-green { stroke: #b8e6c8; stroke-dasharray: 157; stroke-dashoffset: 31; }
.nutrition-circle .fill-purple { stroke: #b8a9e8; stroke-dasharray: 157; stroke-dashoffset: 79; }

.nutrition-value {
  position: absolute;
  font-size: 13px;
  font-weight: 700;
  color: #1a1a1a;
}

.nutrition-label { font-size: 11px; color: #a0a0a0; font-weight: 500; }

/* Steps */
.step-list { margin-bottom: 24px; }

.step-item {
  display: flex;
  gap: 14px;
  margin-bottom: 16px;
}

.step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  flex-shrink: 0;
  background: #1e1e2e;
  color: white;
  margin-top: 2px;
}

.step-content {
  flex: 1;
  border-radius: 14px;
  padding: 14px 16px;
}

.step-title { font-size: 14px; font-weight: 700; color: #1a1a1a; margin-bottom: 4px; }
.step-desc { font-size: 13px; color: #6b6b6b; line-height: 1.6; }

.step-tip {
  margin-top: 8px;
  padding: 8px 12px;
  border-radius: 10px;
  background: rgba(255,243,205,0.5);
  font-size: 12px;
  color: #f59e0b;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.step-tip svg { width: 14px; height: 14px; flex-shrink: 0; }

/* Actions */
.recipe-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 24px;
}

.btn-primary {
  flex: 1;
  padding: 14px;
  border-radius: 16px;
  border: none;
  background: #1e1e2e;
  color: white;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 6px 20px rgba(30,30,46,0.25);
  transition: transform 0.15s ease;
}

.btn-primary:active { transform: scale(0.97); }
.btn-primary svg { width: 20px; height: 20px; }

.btn-secondary {
  width: 50px;
  height: 50px;
  border-radius: 16px;
  border: 1px solid rgba(255,255,255,0.5);
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(20px) saturate(1.4);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 8px 32px rgba(0,0,0,0.06), inset 0 1px 0 rgba(255,255,255,0.8);
  transition: transform 0.15s ease;
  flex-shrink: 0;
}

.btn-secondary:active { transform: scale(0.92); }
.btn-secondary svg { width: 22px; height: 22px; color: #6b6b6b; }
.btn-secondary.liked svg { color: #ef4444; fill: #ef4444; }
</style>
