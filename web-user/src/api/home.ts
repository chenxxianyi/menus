import api from './index'

export function getHomeData() {
  return api.get('/home')
}

export function getBanners() {
  return api.get('/banners')
}
