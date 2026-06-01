<template>
  <div class="page stats-page">
    <!-- Header -->
    <header class="stats-header anim-delay-1">
      <div class="header-spacer"></div>
      <div class="header-center">
        <h1 class="header-title">营养统计</h1>
        <p class="header-sub">本周数据</p>
      </div>
      <button class="settings-btn" aria-label="设置">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
      </button>
    </header>

    <!-- Streak -->
    <div class="streak anim-delay-2">
      <div class="streak-row">
        <div>
          <div class="streak-label">连续记录</div>
          <div class="streak-value">14<span class="streak-unit"> 天</span></div>
        </div>
        <div class="streak-pips">
          <div v-for="i in 7" :key="i" class="pip" :class="{ done: i <= 6, today: i === 7 }"></div>
        </div>
      </div>
    </div>

    <!-- Overview -->
    <div class="overview anim-delay-2">
      <div class="overview-item">
        <div class="overview-num overview-num--accent">1,850</div>
        <div class="overview-unit">kcal</div>
        <div class="overview-change overview-change--down">-120</div>
      </div>
      <div class="overview-item">
        <div class="overview-num overview-num--muted">98g</div>
        <div class="overview-unit">蛋白质</div>
        <div class="overview-change overview-change--up">+15</div>
      </div>
      <div class="overview-item">
        <div class="overview-num overview-num--warm">420g</div>
        <div class="overview-unit">蔬菜</div>
        <div class="overview-change overview-change--up">+80</div>
      </div>
      <div class="overview-item">
        <div class="overview-num overview-num--light">1.8L</div>
        <div class="overview-unit">饮水</div>
        <div class="overview-change overview-change--down">-0.2</div>
      </div>
    </div>

    <!-- Weekly Chart -->
    <div class="block anim-delay-3">
      <div class="block-top">
        <h2 class="block-title">每周热量</h2>
        <div class="tab-toggle">
          <button class="tab-btn active">周</button>
          <button class="tab-btn">月</button>
        </div>
      </div>
      <div class="bar-chart">
        <div v-for="bar in weeklyBars" :key="bar.label" class="bar-col">
          <span class="bar-val">{{ bar.value }}</span>
          <div class="bar-track">
            <div class="bar-fill" :style="{ height: bar.height }"></div>
          </div>
          <span class="bar-lbl">{{ bar.label }}</span>
        </div>
      </div>
    </div>

    <!-- Macro -->
    <div class="block anim-delay-4">
      <h2 class="block-title">营养素分布</h2>
      <div class="macro-row">
        <div class="donut">
          <svg viewBox="0 0 100 100">
            <circle class="donut-track" cx="50" cy="50" r="35"/>
            <circle class="donut-seg donut-seg--1" cx="50" cy="50" r="35"/>
            <circle class="donut-seg donut-seg--2" cx="50" cy="50" r="35"/>
            <circle class="donut-seg donut-seg--3" cx="50" cy="50" r="35"/>
            <circle class="donut-seg donut-seg--4" cx="50" cy="50" r="35"/>
          </svg>
        </div>
        <div class="macro-legend">
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-text);"></span><span class="legend-text">碳水化合物</span><span class="legend-pct">42%</span></div>
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-text-3);"></span><span class="legend-text">蛋白质</span><span class="legend-pct">30%</span></div>
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-accent);"></span><span class="legend-text">脂肪</span><span class="legend-pct">20%</span></div>
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-surface-3);"></span><span class="legend-text">膳食纤维</span><span class="legend-pct">8%</span></div>
        </div>
      </div>
    </div>

    <!-- Meal Distribution -->
    <div class="block anim-delay-5">
      <h2 class="block-title">三餐热量分布</h2>
      <div class="bar-chart">
        <div v-for="bar in mealBars" :key="bar.label" class="bar-col">
          <span class="bar-val">{{ bar.value }}</span>
          <div class="bar-track">
            <div class="bar-fill" :class="`bar-fill--${bar.tone}`" :style="{ height: bar.height }"></div>
          </div>
          <span class="bar-lbl">{{ bar.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const weeklyBars = [
  { label: '一', value: 1720, height: '72%' },
  { label: '二', value: 1950, height: '85%' },
  { label: '三', value: 1850, height: '80%' },
  { label: '四', value: 2100, height: '92%' },
  { label: '五', value: 1680, height: '70%' },
  { label: '六', value: 1980, height: '87%' },
  { label: '日', value: 1850, height: '80%' },
]

const mealBars = [
  { label: '早餐', value: 420, height: '55%', tone: 'accent' },
  { label: '午餐', value: 680, height: '90%', tone: 'muted' },
  { label: '晚餐', value: 520, height: '68%', tone: 'dark' },
  { label: '加餐', value: 230, height: '30%', tone: 'light' },
]
</script>

<style scoped>
.stats-page {
  padding-top: var(--sp-4);
}

/* ── Header ── */
.stats-header {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  margin-bottom: var(--sp-6);
}

.header-spacer {
  width: 36px;
}

.header-center {
  flex: 1;
  text-align: center;
}

.header-title {
  font-size: var(--text-base);
  font-weight: 700;
  color: var(--color-text);
}

.header-sub {
  font-size: 10px;
  color: var(--color-text-3);
  margin-top: 2px;
}

.settings-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.settings-btn:hover {
  background: var(--color-surface-2);
}

.settings-btn svg {
  width: 18px;
  height: 18px;
}

/* ── Streak ── */
.streak {
  background: var(--color-accent);
  border-radius: var(--r-lg);
  padding: var(--sp-5);
  margin-bottom: var(--sp-6);
  color: var(--color-text-inv);
}

.streak-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.streak-label {
  font-size: var(--text-sm);
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
  margin-bottom: var(--sp-1);
}

.streak-value {
  font-size: var(--text-3xl);
  font-weight: 800;
  line-height: 1;
}

.streak-unit {
  font-size: var(--text-base);
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.streak-pips {
  display: flex;
  gap: 4px;
}

.pip {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.2);
}

.pip.done { background: rgba(255, 255, 255, 0.9); }
.pip.today { border: 1.5px solid rgba(255, 255, 255, 0.9); background: transparent; }

/* ── Overview ── */
.overview {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--sp-2);
  margin-bottom: var(--sp-6);
}

.overview-item {
  text-align: center;
  padding: var(--sp-3) 0;
}

.overview-num {
  font-size: var(--text-md);
  font-weight: 700;
}

.overview-num--accent { color: var(--color-accent); }
.overview-num--muted { color: var(--color-text-2); }
.overview-num--warm { color: var(--color-text); }
.overview-num--light { color: var(--color-text-3); }

.overview-unit {
  font-size: 11px;
  color: var(--color-text-3);
  font-weight: 500;
  margin-top: 1px;
}

.overview-change {
  font-size: 11px;
  font-weight: 600;
  margin-top: 2px;
}

.overview-change--up { color: var(--color-success); }
.overview-change--down { color: var(--color-error); }

/* ── Block ── */
.block {
  margin-bottom: var(--sp-8);
}

.block-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-4);
}

.block-title {
  color: var(--color-text);
  font-size: var(--text-lg);
  font-weight: 700;
}

.tab-toggle {
  display: flex;
  gap: 2px;
  background: var(--color-surface-2);
  border-radius: var(--r-sm);
  padding: 2px;
}

.tab-btn {
  padding: 4px 12px;
  border-radius: 6px;
  border: none;
  background: transparent;
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-3);
  cursor: pointer;
  transition: all var(--dur-base) var(--ease);
}

.tab-btn.active {
  background: var(--color-surface);
  color: var(--color-text);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

/* ── Bar Chart ── */
.bar-chart {
  display: flex;
  align-items: flex-end;
  gap: var(--sp-2);
  height: 120px;
  padding-top: var(--sp-3);
}

.bar-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  height: 100%;
  justify-content: flex-end;
}

.bar-track {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: flex-end;
}

.bar-fill {
  width: 100%;
  border-radius: 3px 3px 1px 1px;
  background: var(--color-accent);
  transition: height var(--dur-slow) var(--ease-out);
  min-height: 3px;
}

.bar-fill--accent { background: var(--color-accent); }
.bar-fill--muted { background: var(--color-text-3); }
.bar-fill--dark { background: var(--color-text); }
.bar-fill--light { background: var(--color-surface-3); }

.bar-val {
  font-size: 10px;
  font-weight: 600;
  color: var(--color-text-2);
}

.bar-lbl {
  font-size: 11px;
  color: var(--color-text-3);
  font-weight: 500;
}

/* ── Donut ── */
.macro-row {
  display: flex;
  align-items: center;
  gap: var(--sp-6);
}

.donut { width: 110px; height: 110px; flex-shrink: 0; }
.donut svg { width: 100%; height: 100%; transform: rotate(-90deg); }
.donut circle { fill: none; stroke-width: 10; stroke-linecap: round; }
.donut-track { stroke: var(--color-surface-3); }
.donut-seg--1 { stroke: var(--color-text); stroke-dasharray: 92 220; }
.donut-seg--2 { stroke: var(--color-text-3); stroke-dasharray: 66 220; stroke-dashoffset: -92; }
.donut-seg--3 { stroke: var(--color-accent); stroke-dasharray: 44 220; stroke-dashoffset: -158; }
.donut-seg--4 { stroke: var(--color-surface-3); stroke-dasharray: 18 220; stroke-dashoffset: -202; }

.macro-legend {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3);
}

.legend-row {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  flex-shrink: 0;
}

.legend-text {
  font-size: var(--text-xs);
  color: var(--color-text-2);
}

.legend-pct {
  font-size: var(--text-xs);
  font-weight: 700;
  color: var(--color-text);
  margin-left: auto;
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .page {
    max-width: 640px;
  }

  .overview {
    gap: var(--sp-4);
  }
}

@media (min-width: 1024px) {
  .page {
    max-width: 800px;
  }
}
</style>
