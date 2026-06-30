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
      { path: 'recipes', name: 'RecipeList', component: () => import('@/views/RecipeList.vue'), meta: { showTabBar: true } },
      { path: 'recipes/:id', name: 'RecipeDetail', component: () => import('@/views/RecipeDetail.vue') },
      { path: 'shopping-list', name: 'ShoppingList', component: () => import('@/views/ShoppingList.vue') },
      { path: 'user', name: 'UserCenter', component: () => import('@/views/UserCenter.vue'), meta: { showTabBar: true } },
      { path: 'user/favorites', name: 'Favorites', component: () => import('@/views/Favorites.vue') },
      { path: 'user/preferences', name: 'Preferences', component: () => import('@/views/Preferences.vue') },
      { path: 'user/history', name: 'BrowseHistory', component: () => import('@/views/BrowseHistory.vue') },
      { path: 'about', name: 'About', component: () => import('@/views/About.vue') },
      { path: 'feedback', name: 'Feedback', component: () => import('@/views/Feedback.vue') },
      { path: 'couple', name: 'CoupleHome', component: () => import('@/views/CoupleHome.vue') },
      { path: 'couple/bind', name: 'CoupleBind', component: () => import('@/views/CoupleBind.vue') },
      { path: 'couple/order', name: 'CoupleOrder', component: () => import('@/views/CoupleOrder.vue') },
      { path: 'couple/orders', name: 'CoupleOrders', component: () => import('@/views/CoupleOrders.vue') },
      { path: 'couple/menu', name: 'CoupleMenu', component: () => import('@/views/CoupleMenu.vue') },
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
