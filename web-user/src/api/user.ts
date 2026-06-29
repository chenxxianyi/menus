import api from './index'
import type { UserPreferences } from '@/types/user'

export function getUserInfo() {
  return api.get('/user/info')
}

export function updateProfile(data: Partial<{ nickname: string; avatar: string; bio: string }>) {
  return api.put('/user/profile', data)
}

export function getPreferences() {
  return api.get('/user/preferences')
}

export function updatePreferences(data: UserPreferences) {
  return api.put('/user/preferences', data)
}

export function getFavoriteCount() {
  return api.get('/user/favorites/count')
}

export function getUserStats() {
  return api.get('/user/stats')
}

export function getUserFavorites(page = 1, pageSize = 10) {
  return api.get('/user/favorites', { params: { page, page_size: pageSize } })
}
