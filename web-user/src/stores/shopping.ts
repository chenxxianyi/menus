import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ShoppingItem } from '@/api/shopping'
import { getShoppingList, addShoppingItem, updateShoppingItem, deleteShoppingItem } from '@/api/shopping'

export const useShoppingStore = defineStore('shopping', () => {
  const items = ref<ShoppingItem[]>([])

  async function fetchItems() {
    try {
      const res: any = await getShoppingList()
      items.value = res
    } catch {
      // ignore
    }
  }

  async function addItem(data: Partial<ShoppingItem>) {
    try {
      const res: any = await addShoppingItem(data)
      items.value.push(res)
    } catch {
      // ignore
    }
  }

  async function removeItem(id: number) {
    try {
      await deleteShoppingItem(id)
      items.value = items.value.filter(item => item.id !== id)
    } catch {
      // ignore
    }
  }

  async function toggleChecked(id: number) {
    const item = items.value.find(i => i.id === id)
    if (!item) return
    try {
      await updateShoppingItem(id, { checked: !item.checked })
      item.checked = !item.checked
    } catch {
      // ignore
    }
  }

  return { items, fetchItems, addItem, removeItem, toggleChecked }
})
