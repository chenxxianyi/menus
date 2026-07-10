import api from './index'

export function getHomeData(mealType?: string) {
  return api.get('/home', { params: mealType ? { meal_type: mealType } : undefined })
}

export function getBanners() {
  return api.get('/banners')
}
