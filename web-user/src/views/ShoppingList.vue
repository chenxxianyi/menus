<template>
  <main class="shopping-shell" :style="pageVars">
    <div class="status-bar" aria-hidden="true">
      <span>9:41</span>
      <div class="status-icons">
        <span class="cell-bars"><i></i><i></i><i></i><i></i></span>
        <svg class="wifi" viewBox="0 0 24 18">
          <path d="M2.8 5.8a14.5 14.5 0 0 1 18.4 0" />
          <path d="M6.8 9.8a8.6 8.6 0 0 1 10.4 0" />
          <path d="M10.3 13.5a3 3 0 0 1 3.4 0" />
        </svg>
        <span class="battery"></span>
      </div>
    </div>

    <header class="page-header" aria-label="页面顶部">
      <button class="nav-btn back" type="button" aria-label="返回" @click="router.back()">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <div>
        <h1 class="page-title">购物清单</h1>
        <p class="page-subtitle">根据本周菜单生成</p>
      </div>
      <button class="nav-btn share" type="button" aria-label="分享购物清单" @click="shareShoppingList">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <circle cx="18" cy="5" r="3" />
          <circle cx="6" cy="12" r="3" />
          <circle cx="18" cy="19" r="3" />
          <path d="m8.6 10.6 6.8-4.2" />
          <path d="m8.6 13.4 6.8 4.2" />
        </svg>
      </button>
    </header>

    <section class="summary-card" aria-label="价格汇总">
      <div>
        <p class="summary-label">预估总价</p>
        <strong class="summary-price">¥{{ estimatedTotal }}</strong>
      </div>
      <button class="copy-btn" type="button" @click="copyShoppingList">
        <span>一键复制</span>
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <rect x="9" y="9" width="11" height="11" rx="2" />
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
        </svg>
      </button>
    </section>

    <nav class="category-tabs" aria-label="食材分类筛选">
      <button
        v-for="category in categories"
        :key="category.key"
        class="chip"
        :class="{ active: activeCategory === category.key }"
        type="button"
        @click="activeCategory = category.key"
      >
        <span>{{ category.label }}</span>
        <span>{{ category.count }}</span>
      </button>
    </nav>

    <section class="list-card" aria-label="购物清单">
      <template v-if="filteredGroups.length">
        <section
          v-for="group in filteredGroups"
          :key="group.key"
          class="ingredient-group"
          :class="{ collapsed: isCollapsed(group.key) }"
          :style="{ '--group-color': group.color }"
        >
          <button class="group-header" type="button" :aria-label="`展开或收起${group.title}`" @click="toggleGroup(group.key)">
            <span class="group-name">
              <GroupIcon :type="group.icon" />
              <span>{{ group.title }}</span>
              <span class="group-count">{{ group.items.length }}</span>
            </span>
            <svg class="toggle-icon" viewBox="0 0 24 24" aria-hidden="true">
              <path d="m18 15-6-6-6 6" />
            </svg>
          </button>

          <div class="items">
            <div
              v-for="item in group.items"
              :key="item.id"
              class="ingredient-row"
              :class="{ checked: item.checked }"
              @click="toggleItem(item)"
            >
              <button class="check-btn" type="button" :aria-label="`勾选${item.name}`" @click.stop="toggleItem(item)">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="m5 12 4 4L19 6" />
                </svg>
              </button>
              <span class="ingredient-name">{{ item.name }}</span>
              <span class="quantity">{{ item.quantity }}</span>
              <span class="storage-tag" :class="tagClass(item.tag)">{{ item.tag }}</span>
            </div>
          </div>
        </section>
      </template>

      <section v-else class="empty-state">
        <div class="empty-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24">
            <path d="M6 8h12l-1.1 11.2A2 2 0 0 1 14.9 21H9.1a2 2 0 0 1-2-1.8L6 8Z" />
            <path d="M9 8V6a3 3 0 0 1 6 0v2" />
            <path d="M9.5 12.5h5" />
          </svg>
        </div>
        <h2>暂无食材</h2>
        <p>可以手动添加需要采购的食材</p>
        <button class="empty-add" type="button" @click="addIngredient">添加食材</button>
      </section>

      <button class="add-btn" type="button" @click="addIngredient">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M12 5v14" />
          <path d="M5 12h14" />
        </svg>
        <span>手动添加食材</span>
      </button>
    </section>

    <p class="shared-tip">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
        <path d="m9 12 2 2 4-5" />
      </svg>
      <span>绑定后可一起规划菜单，共享购物清单</span>
    </p>

    <div class="toast" :class="{ show: !!toastText }">{{ toastText }}</div>
  </main>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { updateShoppingList, type ShoppingItem } from '@/api/shopping'
import { useShoppingStore } from '@/stores/shopping'

type GroupKey = 'all' | 'vegetables' | 'protein' | 'seasoning' | 'other'
type GroupIconType = 'leaf' | 'egg' | 'bottle' | 'basket'
type StorageTag = '新鲜' | '冷藏' | '常温'

interface DisplayItem {
  id: string
  name: string
  quantity: string
  tag: StorageTag
  checked: boolean
  sourceIndex?: number
}

interface ShoppingGroup {
  key: Exclude<GroupKey, 'all'>
  title: string
  color: string
  icon: GroupIconType
  items: DisplayItem[]
}

const GROUP_META: Record<Exclude<GroupKey, 'all'>, Omit<ShoppingGroup, 'items'>> = {
  vegetables: { key: 'vegetables', title: '蔬菜', color: '#79A35D', icon: 'leaf' },
  protein: { key: 'protein', title: '肉蛋', color: '#F28A2E', icon: 'egg' },
  seasoning: { key: 'seasoning', title: '调味', color: '#9A7957', icon: 'bottle' },
  other: { key: 'other', title: '其他', color: '#8B715E', icon: 'basket' },
}

const GroupIcon = defineComponent({
  name: 'GroupIcon',
  props: {
    type: {
      type: String,
      default: 'basket',
    },
  },
  setup(props) {
    return () => {
      const common = {
        viewBox: '0 0 24 24',
        'aria-hidden': 'true',
      }

      if (props.type === 'leaf') {
        return h('svg', common, [
          h('path', { d: 'M11 20A7 7 0 0 1 4 13c0-6 7-9 15-9 0 8-3 15-9 15' }),
          h('path', { d: 'M5 19c4-4 8-6 14-8' }),
        ])
      }

      if (props.type === 'egg') {
        return h('svg', common, [
          h('path', { d: 'M12 22c4 0 7-3 7-8 0-6-3.5-12-7-12s-7 6-7 12c0 5 3 8 7 8Z' }),
          h('path', { d: 'M9 15c1.6 1.2 4.3 1.2 6 0' }),
        ])
      }

      if (props.type === 'bottle') {
        return h('svg', common, [
          h('path', { d: 'M10 2h4' }),
          h('path', { d: 'M11 2v5l-2.2 2.2A3 3 0 0 0 8 11.3V20a2 2 0 0 0 2 2h4a2 2 0 0 0 2-2v-8.7a3 3 0 0 0-.8-2.1L13 7V2' }),
          h('path', { d: 'M8 14h8' }),
          h('path', { d: 'M11 18h2' }),
        ])
      }

      return h('svg', common, [
        h('path', { d: 'M6 8h12l-1.1 11.2A2 2 0 0 1 14.9 21H9.1a2 2 0 0 1-2-1.8L6 8Z' }),
        h('path', { d: 'M9 8V6a3 3 0 0 1 6 0v2' }),
      ])
    }
  },
})

const router = useRouter()
const shoppingStore = useShoppingStore()

const activeCategory = ref<GroupKey>('all')
const collapsedKeys = ref(new Set<GroupKey>())
const toastText = ref('')

let toastTimer: ReturnType<typeof setTimeout> | null = null

const pageVars = computed(() => ({
  '--shopping-bg': `url(${kitchenBg})`,
}))

const displayGroups = computed<ShoppingGroup[]>(() => {
  const grouped: Record<Exclude<GroupKey, 'all'>, DisplayItem[]> = {
    vegetables: [],
    protein: [],
    seasoning: [],
    other: [],
  }

  shoppingStore.allItems.forEach((item, index) => {
    const key = resolveGroupKey(item)
    grouped[key].push({
      id: `store-${index}-${item.name}`,
      name: item.name || '未命名食材',
      quantity: item.amount || '1份',
      tag: resolveStorageTag(item, key),
      checked: !!item.checked,
      sourceIndex: index,
    })
  })

  return (Object.keys(grouped) as Exclude<GroupKey, 'all'>[])
    .filter((key) => grouped[key].length > 0)
    .map((key) => ({
      ...GROUP_META[key],
      items: grouped[key],
    }))
})

const filteredGroups = computed(() => {
  if (activeCategory.value === 'all') return displayGroups.value
  return displayGroups.value.filter((group) => group.key === activeCategory.value)
})

const totalCount = computed(() => displayGroups.value.reduce((sum, group) => sum + group.items.length, 0))

const categories = computed(() => {
  const base = [
    { key: 'all' as GroupKey, label: '全部', count: totalCount.value },
    { key: 'vegetables' as GroupKey, label: '蔬菜', count: countGroupItems('vegetables') },
    { key: 'protein' as GroupKey, label: '肉蛋', count: countGroupItems('protein') },
    { key: 'seasoning' as GroupKey, label: '调味', count: countGroupItems('seasoning') },
  ]

  const otherCount = countGroupItems('other')
  if (otherCount) {
    base.push({ key: 'other', label: '其他', count: otherCount })
  }

  return base
})

const estimatedTotal = computed(() => {
  const total = shoppingStore.allItems.reduce((sum, item) => sum + (item.checked ? 0 : Number(item.price || 0)), 0)
  return total.toFixed(1)
})

function countGroupItems(key: Exclude<GroupKey, 'all'>) {
  return displayGroups.value.find((group) => group.key === key)?.items.length || 0
}

function resolveGroupKey(item: ShoppingItem): Exclude<GroupKey, 'all'> {
  const text = `${item.category || ''}${item.name || ''}`
  if (/蔬菜|青菜|叶菜|根茎|瓜果|菌菇/.test(text)) return 'vegetables'
  if (/肉|蛋|鱼|虾|鸡|牛|猪|羊/.test(text)) return 'protein'
  if (/调味|调料|油|盐|酱|醋|糖/.test(text)) return 'seasoning'
  return 'other'
}

function resolveStorageTag(item: ShoppingItem, key: Exclude<GroupKey, 'all'>): StorageTag {
  const text = `${item.category || ''}${item.name || ''}`
  if (/冷藏|肉|鱼|虾|鸡|牛|猪|羊/.test(text)) return '冷藏'
  if (/常温|油|盐|酱|醋|糖/.test(text)) return '常温'
  return key === 'vegetables' ? '新鲜' : key === 'protein' ? '冷藏' : '常温'
}

function tagClass(tag: StorageTag) {
  if (tag === '新鲜') return 'fresh'
  if (tag === '冷藏') return 'chilled'
  return 'room'
}

function isCollapsed(key: GroupKey) {
  return collapsedKeys.value.has(key)
}

function toggleGroup(key: GroupKey) {
  const next = new Set(collapsedKeys.value)
  if (next.has(key)) {
    next.delete(key)
  } else {
    next.add(key)
  }
  collapsedKeys.value = next
}

function showToast(message: string) {
  toastText.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastText.value = ''
  }, 1500)
}

async function toggleItem(item: DisplayItem) {
  const wasChecked = item.checked
  if (typeof item.sourceIndex === 'number') {
    await shoppingStore.toggleItemChecked(item.sourceIndex)
  }
  showToast(wasChecked ? '已取消勾选' : '已标记为已购买')
}

function formatShoppingListText() {
  const lines = ['购物清单：']
  displayGroups.value.forEach((group) => {
    if (!group.items.length) return
    lines.push(`${group.title}：`)
    group.items.forEach((item) => {
      lines.push(`- ${item.name} ${item.quantity}`)
    })
  })
  return lines.join('\n')
}

async function copyShoppingList() {
  const text = formatShoppingListText()
  try {
    await navigator.clipboard.writeText(text)
    showToast('已复制购物清单')
  } catch {
    console.log(text)
    showToast('购物清单已生成，可在控制台查看')
  }
}

async function shareShoppingList() {
  const text = formatShoppingListText()
  try {
    if (navigator.share) {
      await navigator.share({
        title: '购物清单',
        text,
      })
      showToast('分享成功')
      return
    }
    await navigator.clipboard.writeText(text)
    showToast('已复制购物清单，可直接发送')
  } catch {
    showToast('分享未完成')
  }
}

async function persistStoreItems() {
  if (!shoppingStore.currentList) return
  try {
    await updateShoppingList(shoppingStore.currentList.id, {
      name: shoppingStore.currentList.name,
      items: shoppingStore.currentList.items,
    })
  } catch {
    showToast('已添加，稍后同步')
  }
}

async function addIngredient() {
  const name = window.prompt('请输入食材名称')
  if (!name?.trim()) return

  const targetKey: Exclude<GroupKey, 'all'> = activeCategory.value === 'all'
    ? 'seasoning'
    : activeCategory.value === 'other'
      ? 'other'
      : activeCategory.value

  const trimmedName = name.trim().slice(0, 12)
  const newItem: ShoppingItem = {
    name: trimmedName,
    amount: '1份',
    emoji: '',
    category: GROUP_META[targetKey].title,
    price: 0,
    checked: false,
  }

  if (shoppingStore.currentList) {
    shoppingStore.currentList.items.push({
      ...newItem,
    })
    await persistStoreItems()
  } else {
    await shoppingStore.createList('我的购物清单', [newItem])
  }

  const next = new Set(collapsedKeys.value)
  next.delete(targetKey)
  collapsedKeys.value = next
  showToast('已添加食材')
}

onMounted(() => {
  shoppingStore.fetchLists()
})

onUnmounted(() => {
  if (toastTimer) clearTimeout(toastTimer)
})
</script>

<style scoped>
.shopping-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.82);
  --coral: #e95645;
  --coral-2: #ef5548;
  --sage: #7ea36a;
  --sage-soft: #eef6e7;
  --orange: #f28a2e;
  --orange-soft: #fff0df;
  --brown-soft: #f7eee4;
  --border: rgba(255, 255, 255, 0.62);
  --line: rgba(143, 111, 86, 0.22);
  position: relative;
  width: min(100%, 430px);
  min-height: calc(100vh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 180px);
  min-height: calc(100dvh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 180px);
  margin: 0 auto;
  padding: max(18px, env(safe-area-inset-top)) 24px calc(34px + env(safe-area-inset-bottom));
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 237, 205, 0.36), rgba(255, 247, 233, 0.18)),
    var(--shopping-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.shopping-shell::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 20% 2%, rgba(255, 255, 255, 0.72), transparent 28%),
    radial-gradient(circle at 88% 13%, rgba(230, 146, 67, 0.22), transparent 32%),
    radial-gradient(circle at 12% 88%, rgba(233, 86, 69, 0.15), transparent 30%),
    linear-gradient(90deg, rgba(255, 239, 214, 0.56), rgba(255, 245, 230, 0.2) 54%, rgba(172, 91, 33, 0.16));
  backdrop-filter: blur(4px) saturate(1.12);
  -webkit-backdrop-filter: blur(4px) saturate(1.12);
}

.status-bar,
.page-header,
.summary-card,
.category-tabs,
.list-card,
.shared-tip {
  position: relative;
  z-index: 1;
}

button {
  border: 0;
  font: inherit;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

svg {
  display: block;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.status-bar {
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px;
  color: #1e1713;
  font-size: 17px;
  font-weight: 850;
  line-height: 1;
}

.status-icons {
  display: flex;
  align-items: center;
  gap: 7px;
}

.cell-bars {
  display: inline-flex;
  align-items: flex-end;
  gap: 3px;
  height: 18px;
}

.cell-bars i {
  width: 4px;
  border-radius: 999px;
  background: currentColor;
}

.cell-bars i:nth-child(1) { height: 7px; }
.cell-bars i:nth-child(2) { height: 10px; }
.cell-bars i:nth-child(3) { height: 13px; }
.cell-bars i:nth-child(4) { height: 16px; }

.wifi {
  width: 23px;
  height: 17px;
}

.wifi path {
  stroke-width: 2.7;
}

.battery {
  position: relative;
  width: 30px;
  height: 16px;
  border: 2px solid currentColor;
  border-radius: 5px;
}

.battery::before {
  content: "";
  position: absolute;
  right: -5px;
  top: 4px;
  width: 3px;
  height: 6px;
  border-radius: 0 2px 2px 0;
  background: currentColor;
}

.battery::after {
  content: "";
  position: absolute;
  left: 3px;
  top: 3px;
  width: 20px;
  height: 6px;
  border-radius: 2px;
  background: currentColor;
}

.page-header {
  min-height: 96px;
  display: grid;
  place-items: center;
  align-content: center;
  text-align: center;
}

.nav-btn {
  position: absolute;
  top: 14px;
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.66);
  border-radius: 16px;
  color: #4a352a;
  background: rgba(255, 250, 240, 0.88);
  box-shadow:
    0 12px 28px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.nav-btn.back {
  left: 0;
}

.nav-btn.share {
  right: 0;
}

.nav-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.55;
}

.page-title {
  margin: 13px 0 0;
  color: var(--text);
  font-size: 28px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
}

.page-subtitle {
  margin: 14px 0 0;
  color: #6f5d51;
  font-size: 16px;
  font-weight: 650;
  line-height: 1;
}

.summary-card {
  height: 128px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-top: 18px;
  padding: 22px 24px;
  border: 1px solid rgba(255, 255, 255, 0.64);
  border-radius: 30px;
  background:
    radial-gradient(circle at 8% 0%, rgba(255, 255, 255, 0.56), transparent 42%),
    rgba(255, 250, 240, 0.82);
  box-shadow:
    0 18px 42px rgba(80, 50, 30, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.summary-label {
  margin: 0;
  color: var(--sub);
  font-size: 16px;
  font-weight: 760;
  line-height: 1;
}

.summary-price {
  display: block;
  margin-top: 18px;
  color: var(--coral);
  font-size: 45px;
  font-weight: 950;
  line-height: 0.92;
  letter-spacing: 0;
}

.copy-btn {
  height: 58px;
  min-width: 138px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-radius: 18px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 14px 28px rgba(233, 86, 69, 0.26);
  font-size: 18px;
  font-weight: 880;
  white-space: nowrap;
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.copy-btn svg {
  width: 24px;
  height: 24px;
  stroke-width: 2.25;
}

.category-tabs {
  display: flex;
  gap: 14px;
  margin-top: 28px;
  overflow-x: auto;
  padding-bottom: 2px;
  scrollbar-width: none;
}

.category-tabs::-webkit-scrollbar {
  display: none;
}

.chip {
  position: relative;
  height: 52px;
  min-width: 86px;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 0 20px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 999px;
  color: #6b5142;
  background: rgba(255, 250, 240, 0.78);
  box-shadow:
    0 14px 26px rgba(80, 50, 30, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.78);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  font-size: 16px;
  font-weight: 820;
  white-space: nowrap;
  transition: transform 180ms ease, color 180ms ease, opacity 180ms ease;
}

.chip.active {
  color: var(--coral);
  font-weight: 920;
}

.chip.active::after {
  content: "";
  position: absolute;
  left: 50%;
  bottom: 0;
  width: 33px;
  height: 3px;
  border-radius: 999px;
  background: var(--coral);
  transform: translateX(-50%);
}

.list-card {
  margin-top: 20px;
  padding: 21px 16px 24px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.64);
  border-radius: 32px;
  background:
    radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.48), transparent 39%),
    rgba(255, 250, 240, 0.84);
  box-shadow:
    0 22px 50px rgba(80, 50, 30, 0.16),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(22px) saturate(1.1);
  -webkit-backdrop-filter: blur(22px) saturate(1.1);
}

.ingredient-group {
  padding: 0 0 18px;
}

.ingredient-group + .ingredient-group {
  padding-top: 20px;
  border-top: 1.5px dashed var(--line);
}

.group-header {
  width: 100%;
  min-height: 44px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2px 8px;
  color: var(--group-color);
  background: transparent;
  text-align: left;
  transition: transform 180ms ease, opacity 180ms ease;
}

.group-name {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 900;
  line-height: 1;
}

.group-name svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.2;
}

.group-count {
  margin-left: 2px;
  font-size: 17px;
  font-weight: 900;
}

.toggle-icon {
  width: 24px;
  height: 24px;
  color: #6b5142;
  stroke-width: 2.3;
  transition: transform 180ms ease;
}

.ingredient-group.collapsed .toggle-icon {
  transform: rotate(180deg);
}

.items {
  display: grid;
  gap: 8px;
}

.ingredient-group.collapsed .items {
  display: none;
}

.ingredient-row {
  min-height: 56px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 13px;
  border: 1px solid rgba(255, 255, 255, 0.42);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.34);
  box-shadow:
    0 8px 16px rgba(80, 50, 30, 0.04),
    inset 0 1px 0 rgba(255, 255, 255, 0.42);
  transition: opacity 180ms ease, transform 180ms ease, background 180ms ease;
}

.ingredient-row.checked {
  opacity: 0.62;
}

.check-btn {
  width: 25px;
  height: 25px;
  display: grid;
  flex: 0 0 25px;
  place-items: center;
  border: 2px solid #7c6b5f;
  border-radius: 7px;
  color: #fff;
  background: transparent;
  transition: transform 180ms ease, background 180ms ease, border-color 180ms ease;
}

.check-btn svg {
  width: 18px;
  height: 18px;
  stroke-width: 3;
  opacity: 0;
  transform: scale(0.7);
  transition: opacity 160ms ease, transform 160ms ease;
}

.ingredient-row.checked .check-btn {
  border-color: var(--coral);
  background: var(--coral);
}

.ingredient-row.checked .check-btn svg {
  opacity: 1;
  transform: scale(1);
}

.ingredient-name {
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
  color: var(--text);
  font-size: 18px;
  font-weight: 880;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ingredient-row.checked .ingredient-name {
  text-decoration: line-through;
  text-decoration-color: rgba(233, 86, 69, 0.45);
  text-decoration-thickness: 2px;
}

.quantity {
  flex: 0 0 auto;
  color: #574137;
  font-size: 17px;
  font-weight: 760;
  white-space: nowrap;
}

.storage-tag {
  height: 28px;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  padding: 0 11px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 820;
  line-height: 1;
  white-space: nowrap;
}

.storage-tag.fresh {
  color: #6f9a57;
  background: var(--sage-soft);
}

.storage-tag.chilled {
  color: var(--orange);
  background: var(--orange-soft);
}

.storage-tag.room {
  color: #8b715e;
  background: var(--brown-soft);
}

.add-btn {
  width: 100%;
  height: 64px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 10px;
  border: 1.5px dashed rgba(160, 120, 90, 0.34);
  border-radius: 20px;
  color: var(--coral);
  background: rgba(255, 255, 255, 0.22);
  font-size: 18px;
  font-weight: 880;
  transition: transform 180ms ease, background 180ms ease, opacity 180ms ease;
}

.add-btn svg {
  width: 27px;
  height: 27px;
  stroke-width: 2.4;
}

.empty-state {
  padding: 42px 22px 32px;
  text-align: center;
}

.empty-icon {
  width: 72px;
  height: 72px;
  display: grid;
  place-items: center;
  margin: 0 auto 18px;
  border-radius: 24px;
  color: var(--coral);
  background: rgba(252, 231, 222, 0.82);
}

.empty-icon svg {
  width: 40px;
  height: 40px;
  stroke-width: 2.1;
}

.empty-state h2 {
  margin: 0;
  color: var(--text);
  font-size: 23px;
  font-weight: 920;
  line-height: 1;
}

.empty-state p {
  margin: 12px 0 22px;
  color: var(--sub);
  font-size: 15px;
  font-weight: 620;
  line-height: 1.5;
}

.empty-add {
  width: 100%;
  height: 54px;
  border-radius: 18px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 14px 26px rgba(233, 86, 69, 0.25);
  font-size: 17px;
  font-weight: 860;
  transition: transform 180ms ease, opacity 180ms ease;
}

.shared-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 25px 0 4px;
  color: var(--sub);
  font-size: 15px;
  font-weight: 650;
  line-height: 1.35;
  text-align: center;
}

.shared-tip svg {
  width: 23px;
  height: 23px;
  flex: 0 0 23px;
  color: var(--sage);
  stroke-width: 2.2;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: calc(28px + env(safe-area-inset-bottom));
  z-index: 5;
  padding: 10px 16px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 999px;
  color: #fff;
  background: rgba(46, 36, 31, 0.78);
  box-shadow: 0 12px 24px rgba(46, 36, 31, 0.18);
  font-size: 14px;
  font-weight: 740;
  opacity: 0;
  pointer-events: none;
  transform: translate(-50%, 12px);
  transition: opacity 180ms ease, transform 180ms ease;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  white-space: nowrap;
}

.toast.show {
  opacity: 1;
  transform: translate(-50%, 0);
}

.nav-btn:active,
.copy-btn:active,
.chip:active,
.group-header:active,
.ingredient-row:active,
.check-btn:active,
.add-btn:active,
.empty-add:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .nav-btn:hover,
  .copy-btn:hover,
  .chip:hover,
  .add-btn:hover,
  .empty-add:hover {
    transform: translateY(-1px);
  }

  .ingredient-row:hover {
    background: rgba(255, 255, 255, 0.44);
  }
}

@media (max-width: 380px) {
  .shopping-shell {
    padding-right: 18px;
    padding-left: 18px;
  }

  .summary-card {
    padding: 20px;
  }

  .summary-price {
    font-size: 42px;
  }

  .copy-btn {
    min-width: 128px;
    font-size: 17px;
  }

  .category-tabs {
    gap: 12px;
  }

  .chip {
    min-width: 80px;
    padding: 0 17px;
  }

  .list-card {
    padding-right: 14px;
    padding-left: 14px;
  }

  .ingredient-row {
    gap: 10px;
    padding-right: 10px;
    padding-left: 12px;
  }

  .ingredient-name {
    font-size: 17px;
  }

  .quantity {
    font-size: 16px;
  }
}

@media (max-width: 350px) {
  .page-title {
    font-size: 26px;
  }

  .nav-btn {
    width: 48px;
    height: 48px;
  }

  .summary-card {
    height: auto;
    min-height: 128px;
    align-items: flex-start;
    flex-direction: column;
  }

  .copy-btn {
    width: 100%;
  }

  .storage-tag {
    padding: 0 8px;
  }
}

@media (min-width: 431px) {
  .shopping-shell {
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.18);
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
</style>
