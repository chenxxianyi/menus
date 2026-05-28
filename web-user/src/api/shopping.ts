import api from './index'

export interface ShoppingItem {
  id: number
  name: string
  amount: string
  unit: string
  emoji: string
  category: string
  price: number
  checked: boolean
  recipe_id?: number
}

export function getShoppingList() {
  return api.get('/shopping-list')
}

export function addShoppingItem(data: Partial<ShoppingItem>) {
  return api.post('/shopping-list', data)
}

export function updateShoppingItem(id: number, data: Partial<ShoppingItem>) {
  return api.put(`/shopping-list/${id}`, data)
}

export function deleteShoppingItem(id: number) {
  return api.delete(`/shopping-list/${id}`)
}
