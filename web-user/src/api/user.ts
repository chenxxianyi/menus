import api from './index'
import type { PreferenceStatus, UserPreferences } from '@/types/user'

export function getUserInfo() {
  return api.get('/user/info')
}

export function updateProfile(data: Partial<{ nickname: string; avatar: string; bio: string; gender: number }>) {
  return api.put('/user/profile', data)
}

export function uploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/upload/avatar', formData) as Promise<{ url: string }>
}

export function getPreferences() {
  return api.get('/user/preferences')
}

export function getPreferenceStatus() {
  return api.get('/user/preferences/status') as Promise<PreferenceStatus>
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

export function getBrowseHistory(page = 1, pageSize = 20) {
  return api.get('/user/history', { params: { page, page_size: pageSize } })
}

export function clearBrowseHistory() {
  return api.delete('/user/history')
}
