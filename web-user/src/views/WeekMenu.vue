<template>
  <div class="page-container">
    <!-- Header -->
    <header class="app-header animate-fade-up">
      <button class="btn-icon" aria-label="返回" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div class="app-header-center">
        <div class="app-header-title">一周菜单</div>
      </div>
      <button class="btn-icon" aria-label="生成菜单">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
      </button>
    </header>

    <!-- Week Tabs -->
    <div class="week-tabs animate-fade-up anim-delay-1">
      <div
        v-for="day in weekDays"
        :key="day.date"
        class="week-tab"
        :class="{ active: activeDay === day.date }"
        @click="activeDay = day.date"
      >
        <span class="week-day">{{ day.label }}</span>
        <span class="week-date">{{ day.date }}</span>
        <div v-if="day.hasPlan" class="week-dot"></div>
      </div>
    </div>

    <!-- Daily Summary -->
    <div class="summary-bar animate-fade-up anim-delay-2">
      <div class="summary-item orange">
        <div class="summary-value">1,850</div>
        <div class="summary-label">总热量 kcal</div>
      </div>
      <div class="summary-item blue">
        <div class="summary-value">3 餐</div>
        <div class="summary-label">已规划</div>
      </div>
      <div class="summary-item green">
        <div class="summary-value">92%</div>
        <div class="summary-label">营养达标</div>
      </div>
    </div>

    <!-- Meal Timeline -->
    <div class="timeline animate-fade-up anim-delay-3">
      <div v-for="(meal, idx) in dayMeals" :key="idx" class="timeline-item">
        <div class="timeline-dot-wrap">
          <div class="timeline-time">{{ meal.time }}</div>
          <div class="timeline-dot" :class="meal.type"></div>
          <div v-if="idx < dayMeals.length - 1" class="timeline-line"></div>
        </div>
        <div class="timeline-card glass-card" @click="$router.push(`/recipes/${meal.recipeId}`)">
          <div class="timeline-card-header">
            <span class="timeline-card-type" :class="meal.type">{{ meal.typeLabel }}</span>
            <span class="timeline-card-kcal">{{ meal.kcal }} kcal</span>
          </div>
          <div class="timeline-card-name">{{ meal.name }}</div>
          <div class="timeline-card-desc">{{ meal.desc }}</div>
          <div class="timeline-card-tags">
            <span v-for="(tag, tagIdx) in meal.tags" :key="tagIdx" class="timeline-card-tag glass-tag" :class="tag.color">{{ tag.text }}</span>
          </div>
          <div class="timeline-card-images">
            <div v-for="(emoji, emojiIdx) in meal.emojis" :key="emojiIdx" class="timeline-card-img-placeholder" :style="{ background: emoji.bg }">{{ emoji.icon }}</div>
          </div>
        </div>
      </div>

      <!-- Add Meal -->
      <div class="timeline-item">
        <div class="timeline-dot-wrap">
          <div class="timeline-line" style="background:transparent;"></div>
        </div>
        <button class="add-meal-btn">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
          添加加餐
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeDay = ref(14)

const weekDays = [
  { label: '一', date: 12, hasPlan: false },
  { label: '二', date: 13, hasPlan: false },
  { label: '三', date: 14, hasPlan: true },
  { label: '四', date: 15, hasPlan: false },
  { label: '五', date: 16, hasPlan: false },
  { label: '六', date: 17, hasPlan: false },
  { label: '日', date: 18, hasPlan: false },
]

const dayMeals = [
  {
    time: '07:30',
    type: 'breakfast',
    typeLabel: '早餐',
    kcal: 420,
    name: '鸡蛋三明治 + 牛奶',
    desc: '全麦面包搭配煎蛋、生菜、番茄，配一杯温牛奶',
    recipeId: 1,
    tags: [
      { text: '高蛋白', color: 'glass-tag-orange' },
      { text: '低脂', color: 'glass-tag-green' },
    ],
    emojis: [
      { icon: '🥪', bg: 'rgba(255,224,168,0.4)' },
      { icon: '🥛', bg: 'rgba(168,216,234,0.4)' },
    ],
  },
  {
    time: '12:00',
    type: 'lunch',
    typeLabel: '午餐',
    kcal: 680,
    name: '红烧排骨 + 清炒时蔬',
    desc: '经典家常菜搭配，荤素均衡，软烂入味',
    recipeId: 2,
    tags: [
      { text: '荤素搭配', color: 'glass-tag-blue' },
      { text: '30分钟', color: 'glass-tag-orange' },
    ],
    emojis: [
      { icon: '🍖', bg: 'rgba(252,165,165,0.4)' },
      { icon: '🥬', bg: 'rgba(184,230,200,0.4)' },
    ],
  },
  {
    time: '18:30',
    type: 'dinner',
    typeLabel: '晚餐',
    kcal: 520,
    name: '番茄鸡蛋面',
    desc: '酸甜开胃，简单快手，适合晚餐轻食',
    recipeId: 3,
    tags: [
      { text: '快手菜', color: 'glass-tag-purple' },
      { text: '清淡', color: 'glass-tag-green' },
    ],
    emojis: [
      { icon: '🍅', bg: 'rgba(252,165,165,0.4)' },
      { icon: '🍜', bg: 'rgba(255,243,205,0.4)' },
    ],
  },
]
</script>

<style scoped>
/* Week Tabs */
.week-tabs {
  display: flex;
  gap: 6px;
  margin-bottom: 20px;
  overflow-x: auto;
  scrollbar-width: none;
  padding-bottom: 4px;
}

.week-tabs::-webkit-scrollbar { display: none; }

.week-tab {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 14px;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 52px;
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(16px) saturate(1.3);
  border: 1px solid rgba(255,255,255,0.45);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.week-tab.active {
  background: #1e1e2e;
  color: white;
  border-color: transparent;
  box-shadow: 0 6px 20px rgba(30,30,46,0.25);
}

.week-tab:active { transform: scale(0.93); }

.week-day { font-size: 11px; font-weight: 600; color: #a0a0a0; }
.week-tab.active .week-day { color: rgba(255,255,255,0.7); }
.week-date { font-size: 16px; font-weight: 700; color: #1a1a1a; }
.week-tab.active .week-date { color: white; }
.week-dot { width: 4px; height: 4px; border-radius: 50%; background: #ffb347; }

/* Summary Bar */
.summary-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 24px;
}

.summary-item {
  flex: 1;
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(16px) saturate(1.3);
  border-radius: 14px;
  border: 1px solid rgba(255,255,255,0.45);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04), inset 0 1px 0 rgba(255,255,255,0.6);
  padding: 12px 10px;
  text-align: center;
}

.summary-value {
  font-size: 20px;
  font-weight: 800;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.summary-label {
  font-size: 11px;
  color: #a0a0a0;
  font-weight: 500;
  margin-top: 2px;
}

.summary-item.orange .summary-value { color: #f59e0b; }
.summary-item.blue .summary-value { color: #0284c7; }
.summary-item.green .summary-value { color: #16a34a; }

/* Timeline */
.timeline { margin-bottom: 24px; }

.timeline-item {
  display: flex;
  gap: 14px;
  margin-bottom: 20px;
}

.timeline-dot-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  width: 52px;
}

.timeline-time {
  font-size: 11px;
  font-weight: 600;
  color: #a0a0a0;
  white-space: nowrap;
}

.timeline-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.timeline-dot.breakfast { background: #ffb347; box-shadow: 0 2px 8px rgba(255,179,71,0.3); }
.timeline-dot.lunch { background: #a8d8ea; box-shadow: 0 2px 8px rgba(168,216,234,0.3); }
.timeline-dot.dinner { background: #b8a9e8; box-shadow: 0 2px 8px rgba(184,169,232,0.3); }

.timeline-line {
  width: 2px;
  flex: 1;
  background: linear-gradient(180deg, rgba(0,0,0,0.08) 0%, transparent 100%);
  min-height: 20px;
}

.timeline-card {
  flex: 1;
  border-radius: 18px;
  padding: 14px 16px;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.timeline-card:active { transform: scale(0.97); }

.timeline-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.timeline-card-type { font-size: 12px; font-weight: 600; }
.timeline-card-type.breakfast { color: #f59e0b; }
.timeline-card-type.lunch { color: #0284c7; }
.timeline-card-type.dinner { color: #8b7bc6; }

.timeline-card-kcal { font-size: 12px; color: #a0a0a0; font-weight: 500; }

.timeline-card-name {
  font-size: 17px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: -0.3px;
  margin-bottom: 4px;
}

.timeline-card-desc { font-size: 13px; color: #a0a0a0; line-height: 1.5; }

.timeline-card-tags { display: flex; gap: 6px; margin-top: 10px; flex-wrap: wrap; }

.timeline-card-tag { font-size: 11px; padding: 3px 10px; }

.timeline-card-images { display: flex; gap: 8px; margin-top: 10px; }

.timeline-card-img-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  border: 2px solid rgba(255,255,255,0.6);
}

/* Add Button */
.add-meal-btn {
  width: 100%;
  padding: 14px;
  border-radius: 16px;
  border: 2px dashed rgba(0,0,0,0.1);
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #a0a0a0;
  transition: all 0.2s ease;
}

.add-meal-btn:active { background: rgba(0,0,0,0.03); }
.add-meal-btn svg { width: 20px; height: 20px; }
</style>
