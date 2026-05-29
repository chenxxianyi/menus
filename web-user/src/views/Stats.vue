<template>
  <div class="page">
    <!-- Header -->
    <header class="header anim-delay-1">
      <div style="width:36px;"></div>
      <div style="flex:1;text-align:center;">
        <div class="header-title">营养统计</div>
        <div class="header-sub">本周数据</div>
      </div>
      <button class="btn-ghost" aria-label="设置">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
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
        <div class="streak-visual">
          <div v-for="i in 7" :key="i" class="streak-pip" :class="{ done: i <= 6, today: i === 7 }"></div>
        </div>
      </div>
    </div>

    <!-- Overview numbers -->
    <div class="overview anim-delay-2">
      <div class="overview-item">
        <div class="overview-num overview-num--accent">1,850</div>
        <div class="overview-unit">kcal</div>
        <div class="overview-change overview-change--down">-120</div>
      </div>
      <div class="overview-item">
        <div class="overview-num overview-num--blue">98g</div>
        <div class="overview-unit">蛋白质</div>
        <div class="overview-change overview-change--up">+15</div>
      </div>
      <div class="overview-item">
        <div class="overview-num overview-num--green">420g</div>
        <div class="overview-unit">蔬菜</div>
        <div class="overview-change overview-change--up">+80</div>
      </div>
      <div class="overview-item">
        <div class="overview-num overview-num--purple">1.8L</div>
        <div class="overview-unit">饮水</div>
        <div class="overview-change overview-change--down">-0.2</div>
      </div>
    </div>

    <!-- Weekly chart -->
    <div class="chart-section anim-delay-3">
      <div class="chart-header">
        <h2 class="section-heading" style="margin-bottom:0;">每周热量</h2>
        <div class="chart-tabs">
          <button class="chart-tab active">周</button>
          <button class="chart-tab">月</button>
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

    <!-- Macro breakdown -->
    <div class="chart-section anim-delay-4">
      <h2 class="section-heading">营养素分布</h2>
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
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-accent);"></span><span class="legend-text">碳水化合物</span><span class="legend-pct">42%</span></div>
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-blue);"></span><span class="legend-text">蛋白质</span><span class="legend-pct">30%</span></div>
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-green);"></span><span class="legend-text">脂肪</span><span class="legend-pct">20%</span></div>
          <div class="legend-row"><span class="legend-dot" style="background:var(--color-purple);"></span><span class="legend-text">膳食纤维</span><span class="legend-pct">8%</span></div>
        </div>
      </div>
    </div>

    <!-- Meal distribution -->
    <div class="chart-section anim-delay-5">
      <h2 class="section-heading">三餐热量分布</h2>
      <div class="bar-chart">
        <div v-for="bar in mealBars" :key="bar.label" class="bar-col">
          <span class="bar-val">{{ bar.value }}</span>
          <div class="bar-track">
            <div class="bar-fill" :class="`bar-fill--${bar.color}`" :style="{ height: bar.height }"></div>
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
  { label: '早餐', value: 420, height: '55%', color: 'accent' },
  { label: '午餐', value: 680, height: '90%', color: 'blue' },
  { label: '晚餐', value: 520, height: '68%', color: 'purple' },
  { label: '加餐', value: 230, height: '30%', color: 'green' },
]
</script>

<style scoped>
/* ── Streak ── */
.streak {
  background: var(--color-dark);
  border-radius: var(--r-xl);
  padding: var(--sp-5);
  margin-bottom: var(--sp-6);
  color: var(--color-ink);
}

.streak-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.streak-label {
  font-size: var(--text-sm);
  color: var(--color-ink-3);
  font-weight: 500;
  margin-bottom: var(--sp-1);
}

.streak-value {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: 800;
  letter-spacing: 0;
  line-height: 1;
}

.streak-unit {
  font-size: var(--text-base);
  font-weight: 500;
  color: var(--color-ink-3);
}

.streak-visual {
  display: flex;
  gap: 4px;
}

.streak-pip {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background: var(--color-paper-3);
}

.streak-pip.done { background: var(--color-accent); }
.streak-pip.today { border: 1.5px solid var(--color-accent); background: var(--color-accent-soft); }

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
  font-family: var(--font-display);
  font-size: var(--text-md);
  font-weight: 700;
  letter-spacing: 0;
}

.overview-num--accent { color: var(--color-accent); }
.overview-num--blue { color: var(--color-blue); }
.overview-num--green { color: var(--color-green); }
.overview-num--purple { color: var(--color-purple); }

.overview-unit {
  font-size: 11px;
  color: var(--color-ink-3);
  font-weight: 500;
  margin-top: 1px;
}

.overview-change {
  font-size: 11px;
  font-weight: 600;
  margin-top: 2px;
}

.overview-change--up { color: var(--color-green); }
.overview-change--down { color: var(--color-red); }

/* ── Chart section ── */
.chart-section {
  margin-bottom: var(--sp-8);
}

.chart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sp-4);
}

.chart-tabs {
  display: flex;
  gap: 2px;
  background: var(--color-paper-2);
  border-radius: var(--r-sm);
  padding: 2px;
}

.chart-tab {
  padding: 3px 10px;
  border-radius: 6px;
  border: none;
  background: transparent;
  font-size: 11px;
  font-weight: 600;
  color: var(--color-ink-3);
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-out);
}

.chart-tab.active {
  background: var(--color-surface);
  color: var(--color-ink);
  box-shadow: 0 2px 8px var(--color-shadow);
}

/* ── Bar chart ── */
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
  border-radius: 4px 4px 2px 2px;
  background: var(--color-accent);
  transition: height var(--dur-slow) var(--ease-out);
  min-height: 3px;
}

.bar-fill--accent { background: var(--color-accent); }
.bar-fill--blue { background: var(--color-blue); }
.bar-fill--green { background: var(--color-green); }
.bar-fill--purple { background: var(--color-purple); }

.bar-val {
  font-size: 10px;
  font-weight: 600;
  color: var(--color-ink-2);
}

.bar-lbl {
  font-size: 11px;
  color: var(--color-ink-3);
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
.donut-track { stroke: var(--color-paper-3); }
.donut-seg--1 { stroke: var(--color-accent); stroke-dasharray: 92 220; }
.donut-seg--2 { stroke: var(--color-blue); stroke-dasharray: 66 220; stroke-dashoffset: -92; }
.donut-seg--3 { stroke: var(--color-green); stroke-dasharray: 44 220; stroke-dashoffset: -158; }
.donut-seg--4 { stroke: var(--color-purple); stroke-dasharray: 18 220; stroke-dashoffset: -202; }

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
  color: var(--color-ink-2);
}

.legend-pct {
  font-size: var(--text-xs);
  font-weight: 700;
  color: var(--color-ink);
  margin-left: auto;
}
</style>
