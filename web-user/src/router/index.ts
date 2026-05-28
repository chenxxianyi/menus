import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
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

export default router
