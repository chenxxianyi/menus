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

    <section class="dish-generator-card" aria-label="按菜品生成采购清单">
      <div class="dish-generator-copy">
        <span class="dish-generator-badge">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M4 17h16" />
            <path d="M6 17a6 6 0 0 1 12 0" />
            <path d="M12 5v3" />
            <path d="M5 20h14" />
          </svg>
          想吃什么
        </span>
        <h2>输入菜品，合并采购食材</h2>
        <p>例如输入“小鸡炖蘑菇”，系统会匹配真实菜谱并把食材合并到当前清单。</p>
      </div>
      <form class="dish-generator-form" @submit.prevent="generateDishShoppingList">
        <input
          v-model.trim="dishName"
          type="search"
          maxlength="30"
          autocomplete="off"
          enterkeyhint="done"
          placeholder="输入想吃的菜品"
          aria-label="输入想吃的菜品"
          @input="dishGenerateError = ''"
        />
        <button type="submit" :disabled="!canGenerateDishList">
          <span v-if="generatingDishList" class="mini-spinner" aria-hidden="true"></span>
          <span>{{ generatingDishList ? '生成中' : '生成' }}</span>
        </button>
      </form>
      <p v-if="dishGenerateError" class="dish-generator-error" role="alert">{{ dishGenerateError }}</p>
      <button
        v-if="canUseAIForDish"
        class="ai-generate-btn"
        type="button"
        :disabled="generatingAIList"
        @click="generateAIShoppingList"
      >
        <svg v-if="!generatingAIList" viewBox="0 0 24 24" aria-hidden="true">
          <path d="M12 3 13.7 8.3 19 10l-5.3 1.7L12 17l-1.7-5.3L5 10l5.3-1.7L12 3Z" />
          <path d="M19 16v4M21 18h-4" />
        </svg>
        <span v-else class="mini-spinner coral" aria-hidden="true"></span>
        <span>{{ generatingAIList ? 'AI 生成中...' : '使用 AI 生成建议清单' }}</span>
      </button>
      <p v-else class="dish-generator-tip">优先使用真实菜谱食材；同名食材会自动合并，不会覆盖当前清单。</p>
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
      <div v-if="totalCount" class="list-toolbar">
        <div class="list-toolbar-copy">
          <strong>{{ isDeleteMode ? '选择要删除的食材' : '清单明细' }}</strong>
          <span>{{ isDeleteMode ? `已选择 ${selectedDeleteCount} 项` : `共 ${totalCount} 项食材` }}</span>
        </div>
        <div v-if="isDeleteMode" class="delete-toolbar-actions">
          <button
            class="cancel-delete-btn"
            type="button"
            aria-label="退出删除模式"
            @click="toggleDeleteMode"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="m6 6 12 12M18 6 6 18" />
            </svg>
          </button>
          <button
            class="delete-toolbar-confirm-btn"
            type="button"
            :disabled="selectedDeleteCount === 0"
            @click="openDeleteConfirm"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M4 7h16M9 7V4h6v3M7 7l1 13h8l1-13" />
              <path d="M10 11v5M14 11v5" />
            </svg>
            <span>{{ selectedDeleteCount ? `删除 ${selectedDeleteCount}` : '删除' }}</span>
          </button>
        </div>
        <button
          v-else
          class="delete-mode-btn"
          type="button"
          @click="toggleDeleteMode"
        >
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M4 7h16M9 7V4h6v3M7 7l1 13h8l1-13" />
            <path d="M10 11v5M14 11v5" />
          </svg>
          <span>删除</span>
        </button>
      </div>

      <Transition name="delete-actions">
        <div v-if="isDeleteMode" class="delete-action-bar" aria-live="polite">
          <button
            class="select-all-btn"
            type="button"
            :disabled="filteredDeleteIndices.length === 0"
            :aria-pressed="areAllFilteredSelected"
            @click="toggleSelectAllFiltered"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <rect x="4" y="4" width="16" height="16" rx="5" />
              <path v-if="areAllFilteredSelected" d="m8 12 3 3 5-6" />
            </svg>
            <span>{{ areAllFilteredSelected ? '取消全选' : '全选当前分类' }}</span>
          </button>
          <div class="delete-selection-count">
            <span>轻触食材进行选择</span>
            <strong>{{ selectedDeleteCount }} 项已选</strong>
          </div>
        </div>
      </Transition>

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
              :class="{
                checked: item.checked && !isDeleteMode,
                deleting: isDeleteMode,
                'delete-selected': isSelectedForDelete(item),
              }"
              @click="handleIngredientClick(item)"
            >
              <button
                v-if="isDeleteMode"
                class="delete-select-btn"
                :class="{ selected: isSelectedForDelete(item) }"
                type="button"
                :aria-label="`${isSelectedForDelete(item) ? '取消选择' : '选择删除'}${item.name}`"
                :aria-pressed="isSelectedForDelete(item)"
                @click.stop="toggleDeleteSelection(item)"
              >
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="m5 12 4 4L19 6" />
                </svg>
              </button>
              <button v-else class="check-btn" type="button" :aria-label="`勾选${item.name}`" @click.stop="toggleItem(item)">
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
        <p>{{ isDeleteMode ? '当前分类没有可以删除的食材' : '可以手动添加需要采购的食材' }}</p>
        <button v-if="!isDeleteMode" class="empty-add" type="button" @click="openAddIngredientDialog">添加食材</button>
      </section>

      <button v-if="!isDeleteMode" class="add-btn" type="button" @click="openAddIngredientDialog">
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

  <Teleport to="body">
    <Transition name="add-dialog">
      <div
        v-if="addDialogVisible"
        class="add-dialog-mask"
        @click.self="closeAddIngredientDialog"
        @keydown.esc="closeAddIngredientDialog"
        @touchmove.self.prevent
      >
        <section
          class="add-dialog"
          role="dialog"
          aria-modal="true"
          aria-labelledby="add-dialog-title"
          aria-describedby="add-dialog-description"
        >
          <div class="add-dialog-handle" aria-hidden="true"></div>

          <header class="add-dialog-header">
            <span class="add-dialog-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M6 8h12l-1.1 11.2A2 2 0 0 1 14.9 21H9.1a2 2 0 0 1-2-1.8L6 8Z" />
                <path d="M9 8V6a3 3 0 0 1 6 0v2" />
                <path d="M12 11v6M9 14h6" />
              </svg>
            </span>
            <div class="add-dialog-copy">
              <h2 id="add-dialog-title">添加食材</h2>
              <p id="add-dialog-description">记下需要采购的食材，买菜时更省心</p>
            </div>
            <button
              class="add-dialog-close"
              type="button"
              aria-label="关闭添加食材弹框"
              :disabled="addingIngredient"
              @click="closeAddIngredientDialog"
            >
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="m6 6 12 12M18 6 6 18" />
              </svg>
            </button>
          </header>

          <form class="add-dialog-form" @submit.prevent="confirmAddIngredient">
            <label for="ingredient-name">食材名称</label>
            <div class="ingredient-input-wrap">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M12 21c4.4 0 8-3.6 8-8 0-5.2-4.2-9.5-8-11-3.8 1.5-8 5.8-8 11 0 4.4 3.6 8 8 8Z" />
                <path d="M8 15c2.6 0 5.6-2.1 7.5-5.2" />
              </svg>
              <input
                id="ingredient-name"
                ref="ingredientInput"
                v-model="ingredientName"
                type="text"
                maxlength="12"
                autocomplete="off"
                enterkeyhint="done"
                placeholder="例如：西红柿、鸡蛋"
                @input="addDialogError = ''"
              />
              <span class="ingredient-count" aria-hidden="true">{{ ingredientName.length }}/12</span>
            </div>

            <p class="add-category-tip">
              <span aria-hidden="true"></span>
              将添加到“{{ addCategoryLabel }}”分类
            </p>
            <p v-if="addDialogError" class="add-dialog-error" role="alert">{{ addDialogError }}</p>

            <div class="add-dialog-actions">
              <button
                class="add-cancel-btn"
                type="button"
                :disabled="addingIngredient"
                @click="closeAddIngredientDialog"
              >
                取消
              </button>
              <button class="add-confirm-btn" type="submit" :disabled="!canAddIngredient">
                <svg v-if="!addingIngredient" viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M12 5v14M5 12h14" />
                </svg>
                <span>{{ addingIngredient ? '添加中...' : '添加到清单' }}</span>
              </button>
            </div>
          </form>
        </section>
      </div>
    </Transition>

    <Transition name="add-dialog">
      <div
        v-if="deleteConfirmVisible"
        class="add-dialog-mask"
        @click.self="closeDeleteConfirm"
        @keydown.esc="closeDeleteConfirm"
        @touchmove.self.prevent
      >
        <section
          class="delete-dialog"
          role="alertdialog"
          aria-modal="true"
          aria-labelledby="delete-dialog-title"
          aria-describedby="delete-dialog-description"
        >
          <div class="add-dialog-handle" aria-hidden="true"></div>
          <span class="delete-dialog-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24">
              <path d="M4 7h16M9 7V4h6v3M7 7l1 13h8l1-13" />
              <path d="M10 11v5M14 11v5" />
            </svg>
          </span>
          <h2 id="delete-dialog-title">删除所选食材？</h2>
          <p id="delete-dialog-description">将从购物清单中删除 {{ selectedDeleteCount }} 项食材，此操作无法撤销。</p>

          <div class="delete-preview" aria-label="待删除食材">
            <span v-for="(name, index) in selectedDeleteNames.slice(0, 3)" :key="`${name}-${index}`">{{ name }}</span>
            <span v-if="selectedDeleteNames.length > 3">等 {{ selectedDeleteNames.length }} 项</span>
          </div>
          <p v-if="deleteDialogError" class="delete-dialog-error" role="alert">{{ deleteDialogError }}</p>

          <div class="delete-dialog-actions">
            <button class="add-cancel-btn" type="button" :disabled="deletingIngredients" @click="closeDeleteConfirm">
              再想想
            </button>
            <button
              ref="deleteConfirmButton"
              class="delete-confirm-btn"
              type="button"
              :disabled="deletingIngredients"
              @click="confirmDeleteIngredients"
            >
              <svg v-if="!deletingIngredients" viewBox="0 0 24 24" aria-hidden="true">
                <path d="M4 7h16M9 7V4h6v3M7 7l1 13h8l1-13" />
                <path d="M10 11v5M14 11v5" />
              </svg>
              <span>{{ deletingIngredients ? '删除中...' : '确认删除' }}</span>
            </button>
          </div>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { deleteShoppingItems, updateShoppingList, type ShoppingItem } from '@/api/shopping'
import { useShoppingStore } from '@/stores/shopping'

type GroupKey = 'all' | 'vegetables' | 'protein' | 'ingredients' | 'seasoning' | 'other'
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
  ingredients: { key: 'ingredients', title: '配料', color: '#C88768', icon: 'basket' },
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
const addDialogVisible = ref(false)
const ingredientName = ref('')
const ingredientInput = ref<HTMLInputElement | null>(null)
const addingIngredient = ref(false)
const addDialogError = ref('')
const dishName = ref('')
const lastDishNameForAI = ref('')
const generatingDishList = ref(false)
const generatingAIList = ref(false)
const dishGenerateError = ref('')
const canUseAIForDish = ref(false)
const isDeleteMode = ref(false)
const selectedDeleteIndices = ref(new Set<number>())
const deleteConfirmVisible = ref(false)
const deletingIngredients = ref(false)
const deleteConfirmButton = ref<HTMLButtonElement | null>(null)
const deleteDialogError = ref('')

let toastTimer: ReturnType<typeof setTimeout> | null = null

const pageVars = computed(() => ({
  '--shopping-bg': `url(${kitchenBg})`,
}))

const addTargetKey = computed<Exclude<GroupKey, 'all'>>(() => {
  if (activeCategory.value === 'all') return 'ingredients'
  return activeCategory.value
})

const addCategoryLabel = computed(() => GROUP_META[addTargetKey.value].title)
const canAddIngredient = computed(() => ingredientName.value.trim().length > 0 && !addingIngredient.value)
const canGenerateDishList = computed(() => dishName.value.trim().length > 0 && !generatingDishList.value)
const selectedDeleteCount = computed(() => selectedDeleteIndices.value.size)
const filteredDeleteIndices = computed(() => (
  filteredGroups.value
    .flatMap((group) => group.items)
    .map((item) => item.sourceIndex)
    .filter((index): index is number => typeof index === 'number')
))
const areAllFilteredSelected = computed(() => (
  filteredDeleteIndices.value.length > 0
  && filteredDeleteIndices.value.every((index) => selectedDeleteIndices.value.has(index))
))
const selectedDeleteNames = computed(() => (
  [...selectedDeleteIndices.value]
    .sort((a, b) => a - b)
    .map((index) => shoppingStore.allItems[index]?.name)
    .filter((name): name is string => !!name)
))

const displayGroups = computed<ShoppingGroup[]>(() => {
  const grouped: Record<Exclude<GroupKey, 'all'>, DisplayItem[]> = {
    vegetables: [],
    protein: [],
    ingredients: [],
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
    { key: 'ingredients' as GroupKey, label: '配料', count: countGroupItems('ingredients') },
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
  const name = item.name || ''
  const category = item.category || ''

  // 优先按名称纠正历史数据，避免普通配料因旧 category=“调味”继续显示错误。
  if (/酱油|生抽|老抽|蚝油|醋|盐|糖|食用油|香油|料酒|味精|鸡精|胡椒|孜然|辣椒粉|调味酱/.test(name)) return 'seasoning'
  if (/西红柿|番茄|芹菜|香菇|红枣|木耳|豆腐|腐竹|粉丝|花生|芝麻|葱|姜|蒜/.test(name)) return 'ingredients'
  if (/肉|蛋|鱼|虾|鸡|牛|猪|羊/.test(name)) return 'protein'
  if (/青菜|白菜|菠菜|生菜|油菜|菜心|萝卜|土豆|茄子|黄瓜|冬瓜|南瓜|西兰花|豆角|菌菇/.test(name)) return 'vegetables'

  if (/配料/.test(category)) return 'ingredients'
  if (/调味|调料/.test(category)) return 'seasoning'
  if (/肉蛋|肉|蛋|鱼|虾/.test(category)) return 'protein'
  if (/蔬菜|青菜|叶菜|根茎|瓜果|菌菇/.test(category)) return 'vegetables'
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

function toggleDeleteMode() {
  isDeleteMode.value = !isDeleteMode.value
  selectedDeleteIndices.value = new Set()
  deleteConfirmVisible.value = false
  deleteDialogError.value = ''
}

function handleIngredientClick(item: DisplayItem) {
  if (isDeleteMode.value) {
    toggleDeleteSelection(item)
    return
  }
  toggleItem(item)
}

function isSelectedForDelete(item: DisplayItem) {
  return typeof item.sourceIndex === 'number' && selectedDeleteIndices.value.has(item.sourceIndex)
}

function toggleDeleteSelection(item: DisplayItem) {
  if (typeof item.sourceIndex !== 'number') return
  const next = new Set(selectedDeleteIndices.value)
  if (next.has(item.sourceIndex)) {
    next.delete(item.sourceIndex)
  } else {
    next.add(item.sourceIndex)
  }
  selectedDeleteIndices.value = next
}

function toggleSelectAllFiltered() {
  const next = new Set(selectedDeleteIndices.value)
  if (areAllFilteredSelected.value) {
    filteredDeleteIndices.value.forEach((index) => next.delete(index))
  } else {
    filteredDeleteIndices.value.forEach((index) => next.add(index))
  }
  selectedDeleteIndices.value = next
}

function openDeleteConfirm() {
  if (!selectedDeleteCount.value) return
  deleteDialogError.value = ''
  deleteConfirmVisible.value = true
  nextTick(() => deleteConfirmButton.value?.focus())
}

function closeDeleteConfirm() {
  if (deletingIngredients.value) return
  deleteConfirmVisible.value = false
  deleteDialogError.value = ''
}

async function confirmDeleteIngredients() {
  const currentList = shoppingStore.currentList
  if (!currentList || deletingIngredients.value) return

  const indices = [...selectedDeleteIndices.value].sort((a, b) => b - a)
  if (!indices.length) return

  deletingIngredients.value = true
  deleteDialogError.value = ''
  try {
    const result = await deleteShoppingItems(currentList.id, indices)
    currentList.items = Array.isArray(result.items_json) ? result.items_json : []
    selectedDeleteIndices.value = new Set()
    deleteConfirmVisible.value = false
    isDeleteMode.value = false
    showToast(`已删除 ${result.deleted_count} 项食材`)
  } catch (error) {
    deleteDialogError.value = error instanceof Error ? error.message : '删除失败，请稍后重试'
  } finally {
    deletingIngredients.value = false
  }
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
    showToast('已更新，稍后同步')
  }
}

function openAddIngredientDialog() {
  ingredientName.value = ''
  addDialogError.value = ''
  addDialogVisible.value = true
  nextTick(() => ingredientInput.value?.focus())
}

function closeAddIngredientDialog() {
  if (addingIngredient.value) return
  addDialogVisible.value = false
  addDialogError.value = ''
}

async function confirmAddIngredient() {
  const trimmedName = ingredientName.value.trim().slice(0, 12)
  if (!trimmedName || addingIngredient.value) return

  const targetKey = addTargetKey.value
  const newItem: ShoppingItem = {
    name: trimmedName,
    amount: '1份',
    emoji: '',
    category: GROUP_META[targetKey].title,
    price: 0,
    checked: false,
  }

  addingIngredient.value = true
  addDialogError.value = ''
  try {
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
    ingredientName.value = ''
    addDialogVisible.value = false
    showToast('已添加食材')
  } catch {
    addDialogError.value = '添加失败，请稍后重试'
  } finally {
    addingIngredient.value = false
  }
}

async function generateDishShoppingList() {
  const value = dishName.value.trim()
  if (!value || generatingDishList.value) return

  generatingDishList.value = true
  dishGenerateError.value = ''
  canUseAIForDish.value = false
  try {
    const result = await shoppingStore.generateByDish(value)
    dishName.value = ''
    lastDishNameForAI.value = ''
    activeCategory.value = 'all'
    collapsedKeys.value = new Set()
    showToast(formatMergeToast(result.merge_result, result.recipe?.title || value))
  } catch (error) {
    dishGenerateError.value = error instanceof Error ? error.message : '生成失败，请稍后重试'
    lastDishNameForAI.value = value
    canUseAIForDish.value = true
  } finally {
    generatingDishList.value = false
  }
}

async function generateAIShoppingList() {
  const value = (lastDishNameForAI.value || dishName.value).trim()
  if (!value || generatingAIList.value) return

  generatingAIList.value = true
  dishGenerateError.value = ''
  try {
    await shoppingStore.generateByAI(value)
    dishName.value = ''
    lastDishNameForAI.value = ''
    canUseAIForDish.value = false
    activeCategory.value = 'all'
    collapsedKeys.value = new Set()
    showToast('AI建议食材已合并到当前清单')
  } catch (error) {
    dishGenerateError.value = error instanceof Error ? error.message : 'AI 生成失败，请稍后重试'
    canUseAIForDish.value = true
  } finally {
    generatingAIList.value = false
  }
}

function formatMergeToast(result: { added?: number; merged?: number; created?: boolean } | undefined, name: string) {
  if (result?.created) return `已生成「${name}」采购清单`
  const added = result?.added || 0
  const merged = result?.merged || 0
  if (added || merged) return `已合并到当前清单：新增 ${added} 项，合并 ${merged} 项`
  return '清单已是最新，无需重复添加'
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
.dish-generator-card,
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

.dish-generator-card {
  margin-top: 16px;
  padding: 18px;
  border: 1px solid rgba(255, 255, 255, 0.64);
  border-radius: 28px;
  background:
    radial-gradient(circle at 8% 0%, rgba(255, 255, 255, 0.62), transparent 42%),
    linear-gradient(135deg, rgba(255, 250, 240, 0.9), rgba(255, 240, 224, 0.82));
  box-shadow:
    0 18px 38px rgba(80, 50, 30, 0.13),
    inset 0 1px 0 rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.dish-generator-copy {
  display: grid;
  gap: 9px;
}

.dish-generator-badge {
  width: max-content;
  min-height: 31px;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 0 12px;
  border-radius: 999px;
  color: var(--coral);
  background: rgba(255, 255, 255, 0.62);
  font-size: 13px;
  font-weight: 850;
}

.dish-generator-badge svg {
  width: 18px;
  height: 18px;
  stroke-width: 2.2;
}

.dish-generator-copy h2 {
  margin: 0;
  color: var(--text);
  font-size: 21px;
  font-weight: 940;
  line-height: 1.16;
}

.dish-generator-copy p,
.dish-generator-tip,
.dish-generator-error {
  margin: 0;
  font-size: 13px;
  line-height: 1.48;
}

.dish-generator-copy p,
.dish-generator-tip {
  color: var(--sub);
  font-weight: 620;
}

.dish-generator-form {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 86px;
  gap: 10px;
  margin-top: 15px;
}

.dish-generator-form input {
  min-width: 0;
  height: 52px;
  padding: 0 15px;
  border: 1px solid rgba(143, 111, 86, 0.16);
  border-radius: 17px;
  outline: 0;
  color: var(--text);
  background: rgba(255, 255, 255, 0.7);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.78);
  font-size: 16px;
  font-weight: 720;
  transition: border-color 180ms ease, box-shadow 180ms ease, background 180ms ease;
}

.dish-generator-form input:focus {
  border-color: rgba(233, 86, 69, 0.58);
  background: rgba(255, 255, 255, 0.9);
  box-shadow:
    0 0 0 4px rgba(233, 86, 69, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.dish-generator-form input::placeholder {
  color: #a99b90;
  font-weight: 560;
}

.dish-generator-form button {
  height: 52px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  border-radius: 17px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.24);
  font-size: 16px;
  font-weight: 880;
  transition: transform 180ms ease, opacity 180ms ease, box-shadow 180ms ease;
}

.dish-generator-form button:disabled {
  cursor: not-allowed;
  opacity: 0.58;
  box-shadow: none;
}

.mini-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.42);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.72s linear infinite;
}

.dish-generator-tip,
.dish-generator-error {
  margin-top: 11px;
}

.dish-generator-error {
  padding: 10px 12px;
  border: 1px solid rgba(220, 68, 57, 0.16);
  border-radius: 14px;
  color: #c93f35;
  background: rgba(252, 226, 214, 0.56);
  font-weight: 720;
}

.ai-generate-btn {
  width: 100%;
  min-height: 48px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 10px;
  border: 1px solid rgba(233, 86, 69, 0.18);
  border-radius: 16px;
  color: var(--coral);
  background: rgba(255, 255, 255, 0.62);
  box-shadow:
    0 10px 20px rgba(80, 50, 30, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  font-size: 15px;
  font-weight: 860;
  transition: transform 180ms ease, opacity 180ms ease, background 180ms ease;
}

.ai-generate-btn svg {
  width: 20px;
  height: 20px;
  stroke-width: 2.2;
}

.ai-generate-btn:disabled {
  cursor: wait;
  opacity: 0.7;
}

.mini-spinner.coral {
  border-color: rgba(233, 86, 69, 0.18);
  border-top-color: var(--coral);
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

.list-toolbar {
  min-height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin: -3px 2px 16px;
  padding: 0 2px 14px;
  border-bottom: 1px dashed var(--line);
}

.list-toolbar-copy {
  min-width: 0;
  display: grid;
  gap: 5px;
}

.list-toolbar-copy strong {
  color: var(--text);
  font-size: 17px;
  font-weight: 900;
  line-height: 1;
}

.list-toolbar-copy span {
  color: var(--sub);
  font-size: 12px;
  font-weight: 650;
  line-height: 1;
}

.delete-mode-btn {
  min-width: 86px;
  min-height: 44px;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 0 14px;
  border: 1px solid rgba(220, 68, 57, 0.2);
  border-radius: 15px;
  color: #d7463b;
  background: rgba(252, 226, 214, 0.66);
  box-shadow:
    0 8px 18px rgba(220, 68, 57, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  font-size: 14px;
  font-weight: 820;
  transition: transform 180ms ease, color 180ms ease, border-color 180ms ease, background 180ms ease;
}

.delete-mode-btn svg {
  width: 19px;
  height: 19px;
  stroke-width: 2.2;
}

.delete-toolbar-actions {
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 8px;
}

.cancel-delete-btn {
  width: 44px;
  height: 44px;
  display: grid;
  flex: 0 0 44px;
  place-items: center;
  border: 1px solid rgba(143, 111, 86, 0.14);
  border-radius: 15px;
  color: #765f51;
  background: rgba(255, 255, 255, 0.62);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.86);
  transition: transform 180ms ease, color 180ms ease, background 180ms ease;
}

.cancel-delete-btn svg {
  width: 19px;
  height: 19px;
  stroke-width: 2.3;
}

.delete-toolbar-confirm-btn {
  min-width: 92px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 0 14px;
  border-radius: 15px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #db3d32);
  box-shadow: 0 10px 20px rgba(219, 61, 50, 0.22);
  font-size: 14px;
  font-weight: 860;
  white-space: nowrap;
  transition: transform 180ms ease, opacity 180ms ease, box-shadow 180ms ease;
}

.delete-toolbar-confirm-btn svg {
  width: 18px;
  height: 18px;
  stroke-width: 2.2;
}

.delete-toolbar-confirm-btn:disabled {
  cursor: not-allowed;
  opacity: 0.38;
  box-shadow: none;
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

.ingredient-row.deleting {
  cursor: pointer;
}

.ingredient-row.delete-selected {
  border-color: rgba(233, 86, 69, 0.24);
  background: rgba(252, 226, 214, 0.58);
  box-shadow:
    0 8px 18px rgba(233, 86, 69, 0.08),
    inset 0 0 0 1px rgba(255, 255, 255, 0.48);
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

.delete-select-btn {
  width: 27px;
  height: 27px;
  display: grid;
  flex: 0 0 27px;
  place-items: center;
  border: 2px solid #b9a99d;
  border-radius: 50%;
  color: #fff;
  background: rgba(255, 255, 255, 0.56);
  transition: transform 180ms ease, background 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.delete-select-btn svg {
  width: 17px;
  height: 17px;
  stroke-width: 3;
  opacity: 0;
  transform: scale(0.7);
  transition: opacity 160ms ease, transform 160ms ease;
}

.delete-select-btn.selected {
  border-color: var(--coral);
  background: var(--coral);
  box-shadow: 0 5px 12px rgba(233, 86, 69, 0.2);
}

.delete-select-btn.selected svg {
  opacity: 1;
  transform: scale(1);
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

.delete-action-bar {
  min-height: 58px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin: -4px 2px 18px;
  padding: 7px 11px;
  border: 1px solid rgba(233, 86, 69, 0.2);
  border-radius: 20px;
  background:
    radial-gradient(circle at 0% 0%, rgba(255, 255, 255, 0.7), transparent 48%),
    rgba(252, 226, 214, 0.58);
  box-shadow:
    0 12px 24px rgba(80, 50, 30, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.78);
}

.select-all-btn {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 0 10px;
  border: 1px solid rgba(143, 111, 86, 0.14);
  border-radius: 14px;
  color: #765f51;
  background: rgba(255, 255, 255, 0.58);
  font-size: 12px;
  font-weight: 800;
  white-space: nowrap;
  transition: transform 180ms ease, color 180ms ease, opacity 180ms ease, background 180ms ease;
}

.select-all-btn[aria-pressed="true"] {
  color: var(--coral);
  background: rgba(255, 255, 255, 0.76);
}

.select-all-btn:disabled {
  cursor: not-allowed;
  opacity: 0.42;
}

.select-all-btn svg {
  width: 18px;
  height: 18px;
  stroke-width: 2.2;
}

.delete-selection-count {
  min-width: 0;
  display: grid;
  gap: 4px;
  color: var(--sub);
  white-space: nowrap;
  text-align: right;
}

.delete-selection-count strong {
  color: var(--coral);
  font-size: 14px;
  font-weight: 880;
  line-height: 1;
}

.delete-selection-count span {
  overflow: hidden;
  color: var(--sub);
  font-size: 11px;
  font-weight: 650;
  text-overflow: ellipsis;
}

.delete-actions-enter-active,
.delete-actions-leave-active {
  transition: opacity 180ms ease, transform 200ms ease;
}

.delete-actions-enter-from,
.delete-actions-leave-to {
  opacity: 0;
  transform: translateY(8px);
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

.add-dialog-mask {
  --dialog-text: #2e241f;
  --dialog-sub: #7a6a5f;
  --dialog-coral: #e95645;
  position: fixed;
  inset: 0;
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 18px;
  background:
    linear-gradient(180deg, rgba(46, 36, 31, 0.18), rgba(46, 36, 31, 0.42)),
    rgba(46, 36, 31, 0.18);
  backdrop-filter: blur(10px) saturate(1.02);
  -webkit-backdrop-filter: blur(10px) saturate(1.02);
}

.add-dialog {
  width: min(100%, 430px);
  max-height: calc(100dvh - 40px);
  overflow-y: auto;
  padding: 10px 22px 22px;
  border: 1px solid rgba(255, 255, 255, 0.74);
  border-radius: 30px;
  color: var(--dialog-text);
  background:
    radial-gradient(circle at 10% 0%, rgba(255, 255, 255, 0.78), transparent 42%),
    linear-gradient(155deg, rgba(255, 253, 247, 0.97), rgba(255, 244, 228, 0.95));
  box-shadow:
    0 28px 72px rgba(64, 38, 23, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.96);
  overscroll-behavior: contain;
}

.delete-dialog {
  width: min(100%, 390px);
  max-height: calc(100dvh - 40px);
  overflow-y: auto;
  padding: 10px 24px 24px;
  border: 1px solid rgba(255, 255, 255, 0.76);
  border-radius: 30px;
  color: var(--dialog-text);
  background:
    radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.82), transparent 44%),
    linear-gradient(155deg, rgba(255, 253, 247, 0.98), rgba(255, 240, 224, 0.96));
  box-shadow:
    0 28px 72px rgba(64, 38, 23, 0.32),
    inset 0 1px 0 rgba(255, 255, 255, 0.96);
  text-align: center;
  overscroll-behavior: contain;
}

.delete-dialog-icon {
  width: 66px;
  height: 66px;
  display: grid;
  place-items: center;
  margin: 0 auto 17px;
  border: 1px solid rgba(255, 255, 255, 0.82);
  border-radius: 22px;
  color: #dc4439;
  background: rgba(252, 226, 214, 0.84);
  box-shadow:
    0 12px 26px rgba(220, 68, 57, 0.14),
    inset 0 1px 0 rgba(255, 255, 255, 0.88);
}

.delete-dialog-icon svg {
  width: 31px;
  height: 31px;
  stroke-width: 2.15;
}

.delete-dialog h2 {
  margin: 0;
  color: var(--dialog-text);
  font-size: 23px;
  font-weight: 920;
  line-height: 1.15;
}

.delete-dialog > p {
  margin: 10px auto 0;
  max-width: 300px;
  color: var(--dialog-sub);
  font-size: 14px;
  font-weight: 620;
  line-height: 1.55;
}

.delete-preview {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
  margin-top: 18px;
}

.delete-preview span {
  min-height: 30px;
  display: inline-flex;
  align-items: center;
  padding: 0 11px;
  border: 1px solid rgba(143, 111, 86, 0.1);
  border-radius: 999px;
  color: #765f51;
  background: rgba(255, 255, 255, 0.62);
  font-size: 12px;
  font-weight: 740;
}

.delete-dialog-error {
  margin: 14px 0 0 !important;
  padding: 10px 12px;
  border: 1px solid rgba(220, 68, 57, 0.16);
  border-radius: 13px;
  color: #c93f35 !important;
  background: rgba(252, 226, 214, 0.56);
  font-size: 13px !important;
  font-weight: 720 !important;
  line-height: 1.4 !important;
  text-align: left;
}

.delete-dialog-actions {
  display: grid;
  grid-template-columns: minmax(100px, 0.78fr) minmax(0, 1.35fr);
  gap: 12px;
  margin-top: 24px;
}

.delete-confirm-btn {
  min-height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  border-radius: 18px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #db3d32);
  box-shadow: 0 14px 28px rgba(219, 61, 50, 0.26);
  font-size: 16px;
  font-weight: 850;
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.delete-confirm-btn svg {
  width: 22px;
  height: 22px;
  stroke-width: 2.3;
}

.delete-confirm-btn:disabled {
  cursor: not-allowed;
  opacity: 0.52;
  box-shadow: none;
}

.add-dialog-handle {
  width: 44px;
  height: 5px;
  margin: 1px auto 18px;
  border-radius: 999px;
  background: rgba(122, 106, 95, 0.2);
}

.add-dialog-header {
  display: grid;
  grid-template-columns: 54px minmax(0, 1fr) 44px;
  align-items: center;
  gap: 14px;
}

.add-dialog-icon {
  width: 54px;
  height: 54px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.8);
  border-radius: 18px;
  color: var(--dialog-coral);
  background: rgba(252, 226, 214, 0.82);
  box-shadow:
    0 10px 22px rgba(233, 86, 69, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
}

.add-dialog-icon svg {
  width: 29px;
  height: 29px;
  stroke-width: 2.1;
}

.add-dialog-copy {
  min-width: 0;
}

.add-dialog-copy h2 {
  margin: 0;
  color: var(--dialog-text);
  font-size: 23px;
  font-weight: 920;
  line-height: 1.1;
}

.add-dialog-copy p {
  margin: 7px 0 0;
  color: var(--dialog-sub);
  font-size: 13px;
  font-weight: 620;
  line-height: 1.4;
}

.add-dialog-close {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border-radius: 15px;
  color: #8b7d70;
  background: rgba(255, 255, 255, 0.68);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.92);
  transition: transform 180ms ease, color 180ms ease, background 180ms ease, opacity 180ms ease;
}

.add-dialog-close svg {
  width: 21px;
  height: 21px;
  stroke-width: 2.3;
}

.add-dialog-form {
  margin-top: 22px;
}

.add-dialog-form > label {
  display: block;
  margin: 0 0 10px 2px;
  color: #5d4b40;
  font-size: 15px;
  font-weight: 800;
  line-height: 1;
}

.ingredient-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 60px;
  border: 1.5px solid rgba(143, 111, 86, 0.18);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.72);
  box-shadow:
    0 10px 24px rgba(80, 50, 30, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  transition: border-color 180ms ease, box-shadow 180ms ease, background 180ms ease;
}

.ingredient-input-wrap:focus-within {
  border-color: rgba(233, 86, 69, 0.68);
  background: rgba(255, 255, 255, 0.9);
  box-shadow:
    0 0 0 4px rgba(233, 86, 69, 0.1),
    0 12px 28px rgba(80, 50, 30, 0.08);
}

.ingredient-input-wrap > svg {
  position: absolute;
  left: 17px;
  width: 23px;
  height: 23px;
  color: #c88768;
  stroke-width: 2;
  pointer-events: none;
}

.ingredient-input-wrap input {
  width: 100%;
  height: 58px;
  padding: 0 58px 0 51px;
  border: 0;
  outline: 0;
  color: var(--dialog-text);
  background: transparent;
  font: inherit;
  font-size: 17px;
  font-weight: 720;
}

.ingredient-input-wrap input::placeholder {
  color: #a99b90;
  font-weight: 560;
}

.ingredient-count {
  position: absolute;
  right: 16px;
  color: #aa9a8d;
  font-size: 12px;
  font-weight: 700;
  pointer-events: none;
}

.add-category-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 12px 2px 0;
  color: var(--dialog-sub);
  font-size: 13px;
  font-weight: 650;
  line-height: 1.4;
}

.add-category-tip > span {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 50%;
  background: #7ea36a;
  box-shadow: 0 0 0 4px rgba(126, 163, 106, 0.12);
}

.add-dialog-error {
  margin: 10px 2px 0;
  color: #c93f35;
  font-size: 13px;
  font-weight: 700;
  line-height: 1.4;
}

.add-dialog-actions {
  display: grid;
  grid-template-columns: minmax(100px, 0.78fr) minmax(0, 1.42fr);
  gap: 12px;
  margin-top: 22px;
}

.add-cancel-btn,
.add-confirm-btn {
  min-height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  border-radius: 18px;
  font-size: 16px;
  font-weight: 850;
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease, background 180ms ease;
}

.add-cancel-btn {
  color: #6f6054;
  background: rgba(255, 255, 255, 0.66);
  box-shadow:
    inset 0 0 0 1px rgba(143, 111, 86, 0.14),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.add-confirm-btn {
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 14px 28px rgba(233, 86, 69, 0.26);
}

.add-confirm-btn svg {
  width: 22px;
  height: 22px;
  stroke-width: 2.5;
}

.add-dialog-close:disabled,
.add-cancel-btn:disabled,
.add-confirm-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
  box-shadow: none;
}

.add-dialog-enter-active,
.add-dialog-leave-active {
  transition: opacity 200ms ease;
}

.add-dialog-enter-active .add-dialog,
.add-dialog-leave-active .add-dialog,
.add-dialog-enter-active .delete-dialog,
.add-dialog-leave-active .delete-dialog {
  transition: opacity 200ms ease, transform 240ms cubic-bezier(0.16, 1, 0.3, 1);
}

.add-dialog-enter-from,
.add-dialog-leave-to {
  opacity: 0;
}

.add-dialog-enter-from .add-dialog,
.add-dialog-leave-to .add-dialog,
.add-dialog-enter-from .delete-dialog,
.add-dialog-leave-to .delete-dialog {
  opacity: 0;
  transform: translateY(18px) scale(0.98);
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
  bottom: calc(106px + max(16px, var(--safe-bottom)));
  z-index: 110;
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
.delete-select-btn:active,
.add-btn:active,
.empty-add:active,
.delete-mode-btn:active,
.cancel-delete-btn:active,
.delete-toolbar-confirm-btn:active,
.select-all-btn:active,
.add-dialog-close:active,
.add-cancel-btn:active,
.add-confirm-btn:active,
.delete-confirm-btn:active {
  transform: scale(0.98);
}

@media (hover: hover) {
  .nav-btn:hover,
  .copy-btn:hover,
  .chip:hover,
  .add-btn:hover,
  .empty-add:hover,
  .delete-mode-btn:hover,
  .cancel-delete-btn:hover,
  .delete-toolbar-confirm-btn:not(:disabled):hover,
  .select-all-btn:not(:disabled):hover,
  .add-dialog-close:hover,
  .add-cancel-btn:hover,
  .add-confirm-btn:not(:disabled):hover,
  .delete-confirm-btn:not(:disabled):hover {
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
  .delete-mode-btn {
    min-width: 78px;
    padding: 0 11px;
  }

  .delete-toolbar-actions {
    gap: 6px;
  }

  .cancel-delete-btn {
    width: 42px;
    height: 42px;
    flex-basis: 42px;
  }

  .delete-toolbar-confirm-btn {
    min-width: 82px;
    min-height: 42px;
    padding: 0 10px;
  }

  .delete-action-bar {
    gap: 10px;
    padding-right: 9px;
    padding-left: 9px;
  }

  .select-all-btn {
    padding: 0 8px;
  }

  .add-dialog {
    padding-right: 18px;
    padding-left: 18px;
    border-radius: 26px;
  }

  .add-dialog-header {
    grid-template-columns: 50px minmax(0, 1fr) 42px;
    gap: 10px;
  }

  .add-dialog-icon {
    width: 50px;
    height: 50px;
  }

  .add-dialog-copy h2 {
    font-size: 21px;
  }

  .add-dialog-copy p {
    font-size: 12px;
  }

  .add-dialog-actions {
    grid-template-columns: 92px minmax(0, 1fr);
    gap: 10px;
  }

  .delete-dialog {
    padding-right: 18px;
    padding-left: 18px;
    border-radius: 26px;
  }

  .delete-dialog-actions {
    grid-template-columns: 92px minmax(0, 1fr);
    gap: 10px;
  }

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
