<template>
  <div class="page">
    <!-- Header -->
    <header class="header anim-delay-1">
      <button class="btn-ghost" aria-label="返回" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <div style="flex:1;text-align:center;">
        <div class="header-title">购物清单</div>
        <div class="header-sub">根据本周菜单生成</div>
      </div>
      <button class="btn-ghost" aria-label="分享">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/>
          <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
        </svg>
      </button>
    </header>

    <!-- Price summary -->
    <div class="price-bar anim-delay-2">
      <div>
        <div class="price-label">预估总价</div>
        <div class="price-value">¥ {{ totalPrice }}</div>
      </div>
      <button class="btn-solid">一键复制</button>
    </div>

    <!-- Category filter -->
    <div class="cat-filter anim-delay-2">
      <button
        v-for="cat in categories"
        :key="cat.name"
        class="cat-pill"
        :class="{ active: activeCategory === cat.name }"
        @click="activeCategory = cat.name"
      >
        {{ cat.name }}<span class="cat-count">{{ cat.count }}</span>
      </button>
    </div>

    <!-- Groups -->
    <div v-for="group in filteredGroups" :key="group.name" class="shop-group anim-delay-3">
      <div class="group-header">
        <span class="group-emoji">{{ group.emoji }}</span>
        <span class="group-name">{{ group.name }}</span>
        <span class="group-count">{{ group.items.length }} 项</span>
      </div>
      <div class="list-group">
        <div
          v-for="(item, idx) in group.items"
          :key="idx"
          class="list-row shop-row"
          :class="{ checked: item.checked }"
          @click="toggleItem(item)"
        >
          <div class="shop-check" :class="{ checked: item.checked }">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" width="12" height="12"><path d="M20 6L9 17l-5-5"/></svg>
          </div>
          <span class="shop-emoji">{{ item.emoji }}</span>
          <div class="list-row-body">
            <div class="list-row-title" :class="{ 'is-done': item.checked }">{{ item.name }}</div>
            <div class="list-row-sub">{{ item.amount }}</div>
          </div>
          <span class="shop-price">{{ item.price ? '¥' + item.price : '' }}</span>
        </div>
      </div>
    </div>

    <!-- Add -->
    <button class="add-btn anim-delay-4">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18"><path d="M12 5v14M5 12h14"/></svg>
      手动添加
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useShoppingStore } from '@/stores/shopping'

const shoppingStore = useShoppingStore()
const activeCategory = ref('全部')

const categories = computed(() => {
  const items = shoppingStore.allItems
  const catMap: Record<string, number> = {}
  items.forEach(i => {
    const cat = i.category || '其他'
    catMap[cat] = (catMap[cat] || 0) + 1
  })
  return [
    { name: '全部', count: items.length },
    ...Object.entries(catMap).map(([name, count]) => ({ name, count })),
  ]
})

const groups = computed(() => {
  const items = shoppingStore.allItems
  const catMap: Record<string, any[]> = {}
  items.forEach(item => {
    const cat = item.category || '其他'
    if (!catMap[cat]) catMap[cat] = []
    catMap[cat].push(item)
  })
  return Object.entries(catMap).map(([name, items]) => ({
    name,
    emoji: items[0]?.emoji || '📦',
    category: name,
    items,
  }))
})

const filteredGroups = computed(() => {
  if (activeCategory.value === '全部') return groups.value
  return groups.value.filter(g => g.category === activeCategory.value)
})

const totalPrice = computed(() => {
  return shoppingStore.allItems
    .filter(i => !i.checked)
    .reduce((sum, i) => sum + (i.price || 0), 0)
    .toFixed(1)
})

const totalCount = computed(() => shoppingStore.allItems.length)

function toggleItem(item: any) {
  const idx = shoppingStore.allItems.indexOf(item)
  if (idx > -1) shoppingStore.toggleItemChecked(idx)
}

onMounted(() => {
  shoppingStore.fetchLists()
})
</script>

<style scoped>
/* ── Price bar ── */
.price-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-4);
  background: var(--color-surface);
  border-radius: var(--r-lg);
  margin-bottom: var(--sp-5);
  border: 1px solid var(--color-rule);
  box-shadow: 0 8px 28px var(--color-shadow);
}

.price-label {
  font-size: var(--text-xs);
  color: var(--color-ink-3);
  font-weight: 500;
}

.price-value {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 700;
  color: var(--color-ink);
  letter-spacing: 0;
}

/* ── Category filter ── */
.cat-filter {
  display: flex;
  gap: var(--sp-2);
  margin-bottom: var(--sp-5);
  overflow-x: auto;
  scrollbar-width: none;
  -webkit-overflow-scrolling: touch;
}

.cat-filter::-webkit-scrollbar { display: none; }

.cat-pill {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border-radius: var(--r-full);
  border: 1px solid var(--color-rule);
  background: var(--color-surface);
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-ink-2);
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-out);
  -webkit-tap-highlight-color: transparent;
}

.cat-pill.active {
  background: var(--color-dark);
  color: var(--color-ink);
  border-color: var(--color-dark);
}

.cat-pill:active { transform: scale(0.95); }

.cat-count {
  font-size: 10px;
  font-weight: 700;
  opacity: 0.6;
}

/* ── Groups ── */
.shop-group { margin-bottom: var(--sp-5); }

.group-header {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  margin-bottom: var(--sp-2);
  padding: 0 var(--sp-1);
}

.group-emoji { font-size: 16px; }
.group-name { font-size: var(--text-sm); font-weight: 600; color: var(--color-ink); }
.group-count { font-size: 11px; color: var(--color-ink-3); margin-left: auto; }

/* ── Shop row ── */
.shop-row {
  gap: var(--sp-3);
}

.shop-check {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 1.5px solid var(--color-rule);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all var(--dur-base) var(--ease-out);
}

.shop-check.checked {
  background: var(--color-green);
  border-color: var(--color-green);
}

.shop-check svg { display: none; color: white; }
.shop-check.checked svg { display: block; }

.shop-emoji { font-size: 24px; flex-shrink: 0; }

.shop-price {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-accent);
  flex-shrink: 0;
}

.is-done { text-decoration: line-through; color: var(--color-ink-3) !important; }

/* ── Add button ── */
.add-btn {
  width: 100%;
  padding: var(--sp-4);
  border-radius: var(--r-lg);
  border: 1.5px dashed var(--color-rule);
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-2);
  cursor: pointer;
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-ink-3);
  margin-top: var(--sp-2);
  transition: border-color var(--dur-base) var(--ease-out), color var(--dur-base) var(--ease-out);
}

.add-btn:active { border-color: var(--color-accent); color: var(--color-accent); }
</style>
