import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ShoppingItem, ShoppingList } from '@/api/shopping'
import { getShoppingLists, updateShoppingList, createShoppingList, deleteShoppingList, generateShoppingListByDish, generateShoppingListByAI } from '@/api/shopping'

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
    const res: any = await createShoppingList({ name, items })
    const newList: ShoppingList = {
      id: res.id,
      name: res.name,
      items: Array.isArray(res.items_json) ? res.items_json : items,
    }
    lists.value.unshift(newList)
    currentList.value = newList
    return newList
  }

  function normalizeItemName(name: string) {
    return String(name || '').trim().replace(/\s+/g, '')
  }

  function mergeAmount(left: string, right: string) {
    const a = String(left || '').trim()
    const b = String(right || '').trim()
    if (!a || a === '按菜谱适量' || a === '适量') return b || a || '适量'
    if (!b || b === '按菜谱适量' || b === '适量' || b === a) return a
    return a + '、' + b
  }

  function mergeItemsIntoCurrent(items: ShoppingItem[]) {
    if (!currentList.value) return { added: 0, merged: 0 }
    let added = 0
    let merged = 0
    const indexMap = new Map<string, number>()
    currentList.value.items.forEach((item, index) => {
      const key = normalizeItemName(item.name)
      if (key) indexMap.set(key, index)
    })

    items.forEach((item) => {
      const key = normalizeItemName(item.name)
      if (!key) return
      const existingIndex = indexMap.get(key)
      if (existingIndex !== undefined) {
        const existing = currentList.value!.items[existingIndex]
        existing.amount = mergeAmount(existing.amount, item.amount)
        existing.checked = false
        if (!existing.category && item.category) existing.category = item.category
        merged += 1
        return
      }
      currentList.value!.items.push({ ...item, checked: false })
      indexMap.set(key, currentList.value!.items.length - 1)
      added += 1
    })
    return { added, merged }
  }

  async function appendItemsToCurrentList(listName: string, items: ShoppingItem[]) {
    if (!items.length) {
      throw new Error('没有可合并的食材')
    }
    if (!currentList.value) {
      await createList(listName, items)
      return { added: items.length, merged: 0, created: true }
    }
    const originalItems = currentList.value.items.map((item) => ({ ...item }))
    const result = mergeItemsIntoCurrent(items)
    try {
      await updateShoppingList(currentList.value.id, {
        name: currentList.value.name,
        items: currentList.value.items,
      })
    } catch (error) {
      currentList.value.items = originalItems
      throw error
    }
    return { ...result, created: false }
  }

  async function generateByDish(dishName: string) {
    const res = await generateShoppingListByDish(dishName, true)
    const items = Array.isArray(res.items) ? res.items : []
    const mergeResult = await appendItemsToCurrentList((res.recipe?.title || dishName) + '采购清单', items)
    return { ...res, merge_result: mergeResult }
  }

  async function generateByAI(dishName: string) {
    const res = await generateShoppingListByAI(dishName, true)
    const items = Array.isArray(res.items) ? res.items : []
    const mergeResult = await appendItemsToCurrentList(dishName + ' AI建议采购清单', items)
    return { ...res, merge_result: mergeResult }
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

  return { lists, currentList, allItems, totalPrice, fetchLists, toggleItemChecked, createList, generateByDish, generateByAI, removeList }
})
