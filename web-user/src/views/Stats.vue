<template>
  <div class="page-container">
    <!-- Header -->
    <header class="app-header animate-fade-up">
      <div style="width:40px;"></div>
      <div class="app-header-center">
        <div class="app-header-title">营养统计</div>
        <div class="app-header-sub">本周数据</div>
      </div>
      <button class="btn-icon" aria-label="设置">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
        </svg>
      </button>
    </header>

    <!-- Streak -->
    <div class="streak-card animate-fade-up anim-delay-1">
      <div class="streak-title">连续记录天数</div>
      <div class="streak-value">14 <span class="streak-unit">天</span></div>
      <div class="streak-bar">
        <div class="streak-day done"></div>
        <div class="streak-day done"></div>
        <div class="streak-day done"></div>
        <div class="streak-day done"></div>
        <div class="streak-day done"></div>
        <div class="streak-day done"></div>
        <div class="streak-day today"></div>
      </div>
    </div>

    <!-- Overview -->
    <div class="overview-row animate-fade-up anim-delay-2">
      <div class="overview-card orange glass-card-sm">
        <div class="overview-card-label">今日热量</div>
        <div class="overview-card-value">1,850 <span class="overview-card-unit">kcal</span></div>
        <div class="overview-card-change down">-120 vs 昨日</div>
      </div>
      <div class="overview-card blue glass-card-sm">
        <div class="overview-card-label">蛋白质摄入</div>
        <div class="overview-card-value">98 <span class="overview-card-unit">g</span></div>
        <div class="overview-card-change up">+15 达标</div>
      </div>
      <div class="overview-card green glass-card-sm">
        <div class="overview-card-label">蔬菜摄入</div>
        <div class="overview-card-value">420 <span class="overview-card-unit">g</span></div>
        <div class="overview-card-change up">+80 达标</div>
      </div>
      <div class="overview-card purple glass-card-sm">
        <div class="overview-card-label">饮水量</div>
        <div class="overview-card-value">1.8 <span class="overview-card-unit">L</span></div>
        <div class="overview-card-change down">-0.2 待补充</div>
      </div>
    </div>

    <!-- Weekly Calories Chart -->
    <div class="chart-card glass-card anim-delay-3">
      <div class="chart-header">
        <div class="chart-title">每周热量</div>
        <div class="chart-toggle">
          <button class="chart-toggle-btn active">周</button>
          <button class="chart-toggle-btn">月</button>
        </div>
      </div>
      <div class="bar-chart">
        <div v-for="bar in weeklyBars" :key="bar.label" class="bar-group">
          <div class="bar-value">{{ bar.value }}</div>
          <div class="bar orange" :style="{ height: bar.height }"></div>
          <div class="bar-label">{{ bar.label }}</div>
        </div>
      </div>
    </div>

    <!-- Macro Breakdown -->
    <div class="chart-card glass-card anim-delay-4">
      <div class="chart-header">
        <div class="chart-title">营养素分布</div>
      </div>
      <div class="donut-wrap">
        <div class="donut">
          <svg viewBox="0 0 100 100">
            <circle class="track" cx="50" cy="50" r="35"/>
            <circle class="seg-1" cx="50" cy="50" r="35"/>
            <circle class="seg-2" cx="50" cy="50" r="35"/>
            <circle class="seg-3" cx="50" cy="50" r="35"/>
            <circle class="seg-4" cx="50" cy="50" r="35"/>
          </svg>
        </div>
        <div class="donut-legend">
          <div class="legend-item">
            <div class="legend-dot orange"></div>
            <span class="legend-text">碳水化合物</span>
            <span class="legend-value">42%</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot blue"></div>
            <span class="legend-text">蛋白质</span>
            <span class="legend-value">30%</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot green"></div>
            <span class="legend-text">脂肪</span>
            <span class="legend-value">20%</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot purple"></div>
            <span class="legend-text">膳食纤维</span>
            <span class="legend-value">8%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Meal Distribution -->
    <div class="chart-card glass-card anim-delay-5">
      <div class="chart-header">
        <div class="chart-title">三餐热量分布</div>
      </div>
      <div class="bar-chart">
        <div v-for="bar in mealBars" :key="bar.label" class="bar-group">
          <div class="bar-value">{{ bar.value }}</div>
          <div class="bar" :class="bar.color" :style="{ height: bar.height }"></div>
          <div class="bar-label">{{ bar.label }}</div>
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
  { label: '早餐', value: 420, height: '55%', color: 'orange' },
  { label: '午餐', value: 680, height: '90%', color: 'blue' },
  { label: '晚餐', value: 520, height: '68%', color: 'purple' },
  { label: '加餐', value: 230, height: '30%', color: 'green' },
]
</script>

<style scoped>
/* Streak */
.streak-card {
  background: linear-gradient(135deg, #1e1e2e 0%, #2d2d44 100%);
  border-radius: 20px;
  padding: 20px;
  margin-bottom: 20px;
  color: white;
  position: relative;
  overflow: hidden;
}

.streak-card::after {
  content: '';
  position: absolute;
  top: -20px;
  right: -20px;
  width: 100px;
  height: 100px;
  background: rgba(255,255,255,0.08);
  border-radius: 50%;
}

.streak-title { font-size: 14px; color: rgba(255,255,255,0.6); font-weight: 500; margin-bottom: 6px; }

.streak-value {
  font-size: 36px;
  font-weight: 800;
  letter-spacing: -1px;
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.streak-unit { font-size: 16px; font-weight: 600; color: rgba(255,255,255,0.5); }

.streak-bar { display: flex; gap: 4px; margin-top: 14px; }

.streak-day {
  flex: 1;
  height: 6px;
  border-radius: 3px;
  background: rgba(255,255,255,0.12);
}

.streak-day.done { background: #ffb347; }
.streak-day.today { background: rgba(255,255,255,0.3); border: 1px solid rgba(255,255,255,0.5); }

/* Overview Cards */
.overview-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 20px;
}

.overview-card {
  position: relative;
  overflow: hidden;
}

.overview-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: linear-gradient(180deg, rgba(255,255,255,0.25) 0%, transparent 60%);
  pointer-events: none;
}

.overview-card-label { font-size: 12px; color: #a0a0a0; font-weight: 500; margin-bottom: 6px; position: relative; z-index: 1; }

.overview-card-value {
  font-size: 28px;
  font-weight: 800;
  letter-spacing: -0.8px;
  position: relative;
  z-index: 1;
}

.overview-card-unit { font-size: 14px; font-weight: 600; color: #a0a0a0; }

.overview-card-change { font-size: 12px; font-weight: 600; margin-top: 4px; position: relative; z-index: 1; }
.overview-card-change.up { color: #16a34a; }
.overview-card-change.down { color: #dc2626; }

.overview-card.orange .overview-card-value { color: #f59e0b; }
.overview-card.blue .overview-card-value { color: #0284c7; }
.overview-card.green .overview-card-value { color: #16a34a; }
.overview-card.purple .overview-card-value { color: #8b7bc6; }

/* Chart */
.chart-card { padding: 18px; margin-bottom: 20px; }

.chart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.chart-title { font-size: 16px; font-weight: 700; color: #1a1a1a; }

.chart-toggle {
  display: flex;
  gap: 4px;
  background: rgba(0,0,0,0.04);
  border-radius: 10px;
  padding: 3px;
}

.chart-toggle-btn {
  padding: 5px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 12px;
  font-weight: 600;
  color: #a0a0a0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.chart-toggle-btn.active {
  background: white;
  color: #1a1a1a;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}

/* Bar Chart */
.bar-chart {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  height: 140px;
  padding-top: 10px;
}

.bar-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  height: 100%;
  justify-content: flex-end;
}

.bar {
  width: 100%;
  border-radius: 8px 8px 4px 4px;
  position: relative;
  transition: height 0.5s cubic-bezier(0.25,0.46,0.45,0.94);
  min-height: 4px;
}

.bar.orange { background: linear-gradient(180deg, #ffb347 0%, #f59e0b 100%); }
.bar.blue { background: linear-gradient(180deg, #a8d8ea 0%, #6cb4d8 100%); }
.bar.green { background: linear-gradient(180deg, #b8e6c8 0%, #7dd3a0 100%); }
.bar.purple { background: linear-gradient(180deg, #b8a9e8 0%, #8b7bc6 100%); }

.bar-label { font-size: 11px; color: #a0a0a0; font-weight: 500; }
.bar-value { font-size: 11px; font-weight: 700; color: #6b6b6b; }

/* Donut */
.donut-wrap {
  display: flex;
  align-items: center;
  gap: 20px;
}

.donut { width: 120px; height: 120px; flex-shrink: 0; }
.donut svg { width: 100%; height: 100%; transform: rotate(-90deg); }
.donut circle { fill: none; stroke-width: 10; stroke-linecap: round; }
.donut .track { stroke: rgba(0,0,0,0.05); }
.donut .seg-1 { stroke: #ffb347; stroke-dasharray: 94 220; }
.donut .seg-2 { stroke: #a8d8ea; stroke-dasharray: 66 220; stroke-dashoffset: -94; }
.donut .seg-3 { stroke: #b8e6c8; stroke-dasharray: 44 220; stroke-dashoffset: -160; }
.donut .seg-4 { stroke: #b8a9e8; stroke-dasharray: 22 220; stroke-dashoffset: -204; }

.donut-legend { display: flex; flex-direction: column; gap: 10px; }

.legend-item { display: flex; align-items: center; gap: 8px; }

.legend-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.legend-dot.orange { background: #ffb347; }
.legend-dot.blue { background: #a8d8ea; }
.legend-dot.green { background: #b8e6c8; }
.legend-dot.purple { background: #b8a9e8; }

.legend-text { font-size: 13px; color: #6b6b6b; }
.legend-value { font-size: 13px; font-weight: 700; color: #1a1a1a; margin-left: auto; }
</style>
