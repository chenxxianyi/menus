import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: () => import('@/components/AppLayout.vue'),
    children: [
      { path: '', name: 'Home', component: () => import('@/views/Home.vue'), meta: { showTabBar: true } },
      { path: 'week-menu', name: 'WeekMenu', component: () => import('@/views/WeekMenu.vue'), meta: { showTabBar: true } },
      { path: 'stats', name: 'Stats', component: () => import('@/views/Stats.vue'), meta: { showTabBar: true } },
      { path: 'recipes/:id', name: 'RecipeDetail', component: () => import('@/views/RecipeDetail.vue') },
      { path: 'shopping-list', name: 'ShoppingList', component: () => import('@/views/ShoppingList.vue') },
      { path: 'user', name: 'UserCenter', component: () => import('@/views/UserCenter.vue'), meta: { showTabBar: true } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (!token && !to.meta.guest) {
    return { name: 'Login' }
  }
})

export default router
