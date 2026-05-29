import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ShoppingItem, ShoppingList } from '@/api/shopping'
import { getShoppingLists, updateShoppingList, createShoppingList, deleteShoppingList } from '@/api/shopping'

export const useShoppingStore = defineStore('shopping', () => {
  const lists = ref<ShoppingList[]>([])
  const currentList = ref<ShoppingList | null>(null)

  const allItems = computed(() => currentList.value?.items ?? [])

  const totalPrice = computed(() =>
    allItems.value.reduce((sum, item) => sum + (item.checked ? 0 : item.price), 0)
  )

  async function fetchLists() {
    try {
      const res: any = await getShoppingLists()
      lists.value = Array.isArray(res) ? res.map((l: any) => ({
        id: l.id,
        name: l.name,
        items: Array.isArray(l.items_json) ? l.items_json : [],
      })) : []
      if (lists.value.length > 0 && !currentList.value) {
        currentList.value = lists.value[0]
      }
    } catch {
      // ignore
    }
  }

  async function saveCurrentList() {
    if (!currentList.value) return
    try {
      await updateShoppingList(currentList.value.id, {
        name: currentList.value.name,
        items: currentList.value.items,
      })
    } catch {
      // ignore
    }
  }

  async function toggleItemChecked(index: number) {
    const items = currentList.value?.items
    if (!items || !items[index]) return
    items[index].checked = !items[index].checked
    await saveCurrentList()
  }

  async function createList(name: string, items: ShoppingItem[]) {
    try {
      const res: any = await createShoppingList({ name, items })
      const newList: ShoppingList = {
        id: res.id,
        name: res.name,
        items: Array.isArray(res.items_json) ? res.items_json : items,
      }
      lists.value.push(newList)
      currentList.value = newList
    } catch {
      // ignore
    }
  }

  async function removeList(id: number) {
    try {
      await deleteShoppingList(id)
      lists.value = lists.value.filter(l => l.id !== id)
      if (currentList.value?.id === id) {
        currentList.value = lists.value[0] ?? null
      }
    } catch {
      // ignore
    }
  }

  return { lists, currentList, allItems, totalPrice, fetchLists, toggleItemChecked, createList, removeList }
})
