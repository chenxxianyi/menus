<template>
  <main class="menus-page">
    <header class="page-header">
      <button class="back-btn" type="button" aria-label="返回" @click="router.back()">
        <svg viewBox="0 0 24 24" aria-hidden="true"><path d="m15 18-6-6 6-6" /></svg>
      </button>
      <div>
        <h1>我的菜单</h1>
        <p>保存常吃搭配，下次直接复用</p>
      </div>
      <button class="back-btn" type="button" aria-label="安排本周菜单" @click="router.push('/week-menu')">
        <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 19h16" /><path d="M6 17a6 6 0 0 1 12 0" /><path d="M12 5v3" /></svg>
      </button>
    </header>

    <section v-if="loading" class="state-card">
      <span class="spinner"></span>
      <p>正在整理你的菜单...</p>
    </section>

    <section v-else-if="!menus.length" class="state-card">
      <h2>还没有保存菜单</h2>
      <p>在推荐结果或一周菜单里保存后，会出现在这里。</p>
      <button type="button" @click="router.push('/recommend/scene')">去推荐</button>
    </section>

    <section v-else class="menu-list">
      <article v-for="menu in menus" :key="menu.id" class="menu-card">
        <div class="menu-card-main">
          <span>{{ menuTypeLabel(menu.meal_type) }} · {{ formatDate(menu.created_at) }}</span>
          <h2>{{ menu.menu_name || '我的菜单' }}</h2>
          <p>{{ menu.reason || dishNames(menu).join('、') || '已保存的家庭菜单' }}</p>
          <div class="dish-tags">
            <button v-for="dish in dishNames(menu).slice(0, 4)" :key="dish" type="button">{{ dish }}</button>
          </div>
        </div>
        <div class="menu-actions">
          <button type="button" :disabled="busyId === menu.id" @click="reuseMenu(menu)">生成清单</button>
          <button type="button" @click="router.push(`/recommend/scene`)">再推荐</button>
          <button class="danger" type="button" :disabled="busyId === menu.id" @click="removeMenu(menu)">删除</button>
        </div>
      </article>
    </section>

    <p v-if="message" class="message">{{ message }}</p>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { deleteUserMenu, getUserMenus, reuseUserMenu, type UserMenu } from '@/api/menu'
import { useShoppingStore } from '@/stores/shopping'

const router = useRouter()
const shoppingStore = useShoppingStore()
const menus = ref<UserMenu[]>([])
const loading = ref(false)
const busyId = ref(0)
const message = ref('')

async function loadMenus() {
  loading.value = true
  try {
    const res: any = await getUserMenus({ page: 1, page_size: 30 })
    menus.value = Array.isArray(res?.list) ? res.list : Array.isArray(res) ? res : []
  } finally {
    loading.value = false
  }
}

function dishNames(menu: UserMenu) {
  const names: string[] = []
  const walk = (value: any) => {
    if (!value) return
    if (Array.isArray(value)) {
      value.forEach(walk)
      return
    }
    if (typeof value === 'object') {
      const name = String(value.name || value.title || '').trim()
      if (name && !names.includes(name)) names.push(name)
      Object.values(value).forEach(walk)
    }
  }
  walk(menu.dishes_json)
  return names
}

function menuTypeLabel(type: string) {
  const map: Record<string, string> = {
    daily: '日常菜单',
    weekly: '周菜单',
    couple: '情侣菜单',
    ai: 'AI 菜单',
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐',
  }
  return map[type] || '菜单'
}

function formatDate(value: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return `${date.getMonth() + 1}/${date.getDate()}`
}

async function reuseMenu(menu: UserMenu) {
  busyId.value = menu.id
  message.value = ''
  try {
    const result = await reuseUserMenu(menu.id)
    if (result.recipe_ids?.length) {
      await shoppingStore.generateByRecipes(result.recipe_ids, `${menu.menu_name || '我的菜单'}采购清单`)
      message.value = '已按历史菜单生成购物清单。'
      setTimeout(() => router.push('/shopping-list'), 350)
      return
    }
    message.value = '这个菜单暂时没有可复用的菜谱。'
  } catch (error) {
    message.value = error instanceof Error ? error.message : '复用失败'
  } finally {
    busyId.value = 0
  }
}

async function removeMenu(menu: UserMenu) {
  if (!window.confirm(`删除「${menu.menu_name || '我的菜单'}」吗？`)) return
  busyId.value = menu.id
  try {
    await deleteUserMenu(menu.id)
    menus.value = menus.value.filter((item) => item.id !== menu.id)
    message.value = '已删除菜单。'
  } catch (error) {
    message.value = error instanceof Error ? error.message : '删除失败'
  } finally {
    busyId.value = 0
  }
}

onMounted(loadMenus)
</script>

<style scoped>
.menus-page {
  min-height: 100vh;
  padding: max(18px, env(safe-area-inset-top)) 18px calc(86px + env(safe-area-inset-bottom));
  color: #2e241f;
  background: linear-gradient(180deg, #fff6e8 0%, #f7ead8 100%);
}

.page-header {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr) 44px;
  align-items: center;
  gap: 12px;
  margin-bottom: 18px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
}

.page-header p {
  margin: 4px 0 0;
  color: #7a6a5f;
  font-size: 13px;
}

.back-btn {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(143, 111, 86, 0.18);
  border-radius: 14px;
  color: #5f4c43;
  background: rgba(255, 255, 255, 0.72);
}

.back-btn svg {
  width: 22px;
  height: 22px;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 2.3;
}

.state-card,
.menu-card {
  border: 1px solid rgba(143, 111, 86, 0.14);
  border-radius: 8px;
  background: rgba(255, 250, 240, 0.86);
  box-shadow: 0 14px 30px rgba(80, 50, 28, 0.1);
}

.state-card {
  display: grid;
  gap: 12px;
  justify-items: center;
  padding: 34px 18px;
  text-align: center;
}

.state-card h2,
.state-card p {
  margin: 0;
}

.state-card p {
  color: #7a6a5f;
}

.state-card button,
.menu-actions button {
  min-height: 38px;
  border: 0;
  border-radius: 8px;
  background: #2e241f;
  color: #fffaf0;
  font-size: 13px;
  font-weight: 850;
}

.spinner {
  width: 26px;
  height: 26px;
  border: 3px solid rgba(46, 36, 31, 0.16);
  border-top-color: #e95645;
  border-radius: 50%;
  animation: spin 800ms linear infinite;
}

.menu-list {
  display: grid;
  gap: 12px;
}

.menu-card {
  padding: 16px;
}

.menu-card-main span {
  color: #e95645;
  font-size: 12px;
  font-weight: 850;
}

.menu-card-main h2 {
  margin: 6px 0;
  font-size: 19px;
}

.menu-card-main p {
  margin: 0;
  color: #7a6a5f;
  font-size: 13px;
  line-height: 1.55;
}

.dish-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 7px;
  margin-top: 12px;
}

.dish-tags button {
  min-height: 28px;
  border: 0;
  border-radius: 999px;
  color: #5f4c43;
  background: rgba(233, 86, 69, 0.1);
  font-size: 12px;
  font-weight: 760;
}

.menu-actions {
  display: grid;
  grid-template-columns: 1fr 1fr 72px;
  gap: 8px;
  margin-top: 14px;
}

.menu-actions .danger {
  color: #e95645;
  background: rgba(233, 86, 69, 0.12);
}

.menu-actions button:disabled {
  opacity: 0.55;
}

.message {
  margin: 16px 0 0;
  color: #6f8d55;
  font-size: 13px;
  font-weight: 800;
  text-align: center;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
