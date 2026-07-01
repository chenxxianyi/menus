import api from './index'

export interface UserMenu {
  id: number
  user_id: number
  menu_name: string
  meal_type: string
  people_count: number
  taste: string
  health_goal: string
  dishes_json: any
  shopping_list_json: any
  reason: string
  created_at: string
}

export interface SaveMenuPayload {
  name: string
  menu_type: string
  meal_type?: string
  people_count?: number
  taste?: string
  health_goal?: string
  dishes: any
  shopping_list?: any
  reason?: string
}

export function getUserMenus(params?: { page?: number; page_size?: number }) {
  return api.get('/user/menus', { params })
}

export function saveUserMenu(data: SaveMenuPayload): Promise<UserMenu> {
  return api.post('/user/menus', data) as unknown as Promise<UserMenu>
}

export function getUserMenu(id: number): Promise<UserMenu> {
  return api.get(`/user/menus/${id}`) as unknown as Promise<UserMenu>
}

export function deleteUserMenu(id: number) {
  return api.delete(`/user/menus/${id}`)
}

export function reuseUserMenu(id: number): Promise<{ menu: UserMenu; recipe_ids: number[]; recipes: any[] }> {
  return api.post(`/user/menus/${id}/reuse`) as unknown as Promise<{ menu: UserMenu; recipe_ids: number[]; recipes: any[] }>
}
