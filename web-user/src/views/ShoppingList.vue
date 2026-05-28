<template>
  <div class="page-container">
    <!-- Header -->
    <header class="app-header animate-fade-up">
      <button class="btn-icon" aria-label="返回" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div class="app-header-center">
        <div class="app-header-title">购物清单</div>
        <div class="app-header-sub">根据本周菜单生成</div>
      </div>
      <button class="btn-icon" aria-label="分享">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/>
          <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
        </svg>
      </button>
    </header>

    <!-- Summary -->
    <div class="shop-summary glass-card animate-fade-up anim-delay-1">
      <div class="shop-summary-left">
        <div class="shop-summary-label">预估总价</div>
        <div class="shop-summary-value">¥ 86.5</div>
      </div>
      <button class="shop-summary-btn">一键复制</button>
    </div>

    <!-- Category Tabs -->
    <div class="cat-tabs animate-fade-up anim-delay-2">
      <div
        v-for="cat in categories"
        :key="cat.name"
        class="cat-tab"
        :class="{ active: activeCategory === cat.name }"
        @click="activeCategory = cat.name"
      >
        {{ cat.name }} <span class="cat-tab-count">{{ cat.count }}</span>
      </div>
    </div>

    <!-- Shopping Groups -->
    <div v-for="group in filteredGroups" :key="group.name" class="shop-group animate-fade-up anim-delay-3">
      <div class="shop-group-header">
        <div class="shop-group-icon">{{ group.emoji }}</div>
        <div class="shop-group-name">{{ group.name }}</div>
        <div class="shop-group-count">{{ group.items.length }} 项</div>
      </div>
      <div
        v-for="item in group.items"
        :key="item.name"
        class="shop-card glass-card-sm"
        :class="{ checked: item.checked }"
        @click="item.checked = !item.checked"
      >
        <div class="shop-check">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg>
        </div>
        <div class="shop-emoji">{{ item.emoji }}</div>
        <div class="shop-info">
          <div class="shop-name">{{ item.name }}</div>
          <div class="shop-detail">{{ item.detail }}</div>
        </div>
        <div class="shop-price">{{ item.price }}</div>
      </div>
    </div>

    <!-- Add Button -->
    <button class="add-item-btn anim-delay-4">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14"/></svg>
      手动添加
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const activeCategory = ref('全部')

const categories = [
  { name: '全部', count: 12 },
  { name: '蔬菜', count: 4 },
  { name: '肉类', count: 3 },
  { name: '调料', count: 3 },
  { name: '主食', count: 2 },
]

const groups = ref([
  {
    name: '蔬菜水果',
    emoji: '🥬',
    category: '蔬菜',
    items: [
      { name: '番茄', detail: '3 个，选红透的', price: '¥8.5', emoji: '🍅', checked: false },
      { name: '生菜', detail: '1 颗', price: '¥3.0', emoji: '🥬', checked: false },
      { name: '大葱', detail: '2 根', price: '¥2.0', emoji: '🧅', checked: false },
      { name: '黄瓜', detail: '2 根', price: '¥4.0', emoji: '🥒', checked: true },
    ],
  },
  {
    name: '肉蛋水产',
    emoji: '🥩',
    category: '肉类',
    items: [
      { name: '排骨', detail: '500g，肋排', price: '¥38.0', emoji: '🍖', checked: false },
      { name: '鸡蛋', detail: '10 枚', price: '¥12.0', emoji: '🥚', checked: false },
      { name: '鸡胸肉', detail: '300g', price: '¥11.0', emoji: '🐔', checked: false },
    ],
  },
  {
    name: '调料干货',
    emoji: '🧂',
    category: '调料',
    items: [
      { name: '冰糖', detail: '200g 袋装', price: '¥5.0', emoji: '🫙', checked: false },
      { name: '料酒', detail: '家中已有', price: '¥0', emoji: '🍶', checked: true },
      { name: '挂面', detail: '500g', price: '¥6.0', emoji: '🍜', checked: false },
    ],
  },
])

const filteredGroups = computed(() => {
  if (activeCategory.value === '全部') return groups.value
  return groups.value.filter(g => g.category === activeCategory.value)
})
</script>

<style scoped>
/* Category Tabs */
.cat-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  overflow-x: auto;
  scrollbar-width: none;
}

.cat-tabs::-webkit-scrollbar { display: none; }

.cat-tab {
  flex-shrink: 0;
  padding: 8px 18px;
  border-radius: 9999px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(16px) saturate(1.3);
  border: 1px solid rgba(255,255,255,0.45);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  color: #6b6b6b;
  white-space: nowrap;
}

.cat-tab.active {
  background: #1e1e2e;
  color: white;
  border-color: transparent;
  box-shadow: 0 4px 16px rgba(30,30,46,0.2);
}

.cat-tab:active { transform: scale(0.93); }

.cat-tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  font-size: 11px;
  font-weight: 700;
  margin-left: 4px;
  padding: 0 5px;
}

.cat-tab.active .cat-tab-count { background: rgba(255,255,255,0.2); color: white; }
.cat-tab:not(.active) .cat-tab-count { background: rgba(0,0,0,0.06); color: #a0a0a0; }

/* Shopping Group */
.shop-group { margin-bottom: 20px; }

.shop-group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  padding-left: 4px;
}

.shop-group-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255,255,255,0.4);
}

.shop-group-name { font-size: 14px; font-weight: 700; color: #1a1a1a; }
.shop-group-count { font-size: 12px; color: #a0a0a0; font-weight: 500; margin-left: auto; }

/* Shopping Card */
.shop-card {
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  transition: all 0.2s ease;
  cursor: pointer;
}

.shop-card:active { transform: scale(0.98); }
.shop-card.checked { opacity: 0.5; }

.shop-check {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid #a0a0a0;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.shop-card.checked .shop-check { background: #b8e6c8; border-color: #16a34a; }
.shop-check svg { width: 14px; height: 14px; color: #16a34a; display: none; }
.shop-card.checked .shop-check svg { display: block; }

.shop-emoji { font-size: 28px; flex-shrink: 0; }

.shop-info { flex: 1; min-width: 0; }

.shop-name { font-size: 15px; font-weight: 600; color: #1a1a1a; }
.shop-card.checked .shop-name { text-decoration: line-through; color: #a0a0a0; }

.shop-detail { font-size: 12px; color: #a0a0a0; margin-top: 2px; }

.shop-price { font-size: 14px; font-weight: 700; color: #f59e0b; flex-shrink: 0; }

/* Summary Bar */
.shop-summary {
  padding: 16px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.shop-summary-left { display: flex; flex-direction: column; }

.shop-summary-label { font-size: 12px; color: #a0a0a0; font-weight: 500; }
.shop-summary-value { font-size: 24px; font-weight: 800; color: #1a1a1a; letter-spacing: -0.5px; }

.shop-summary-btn {
  padding: 10px 20px;
  border-radius: 14px;
  border: none;
  background: #1e1e2e;
  color: white;
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(30,30,46,0.2);
  transition: transform 0.15s ease;
}

.shop-summary-btn:active { transform: scale(0.96); }

/* Add Button */
.add-item-btn {
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
  margin-top: 8px;
  transition: all 0.2s ease;
}

.add-item-btn:active { background: rgba(0,0,0,0.03); }
.add-item-btn svg { width: 20px; height: 20px; }
</style>
