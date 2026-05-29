import api from './index'

export interface ShoppingItem {
  name: string
  amount: string
  emoji: string
  category: string
  price: number
  checked: boolean
}

export interface ShoppingList {
  id: number
  name: string
  items: ShoppingItem[]
}

export function getShoppingLists() {
  return api.get('/shopping-list')
}

export function createShoppingList(data: { name: string; items: ShoppingItem[] }) {
  return api.post('/shopping-list', { name: data.name, items_json: data.items })
}

export function updateShoppingList(id: number, data: { name: string; items: ShoppingItem[] }) {
  return api.put(`/shopping-list/${id}`, { name: data.name, items_json: data.items })
}

export function deleteShoppingList(id: number) {
  return api.delete(`/shopping-list/${id}`)
}
