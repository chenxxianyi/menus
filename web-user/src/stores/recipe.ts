import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Recipe, Category } from '@/types/recipe'
import { getCategories, getRecipeDetail, toggleFavorite } from '@/api/recipe'

export const useRecipeStore = defineStore('recipe', () => {
  const categories = ref<Category[]>([])
  const currentRecipe = ref<Recipe | null>(null)
  const favorites = ref<number[]>([])

  async function fetchCategories() {
    try {
      const res: any = await getCategories()
      categories.value = res
    } catch {
      // ignore
    }
  }

  async function fetchRecipeDetail(id: number) {
    try {
      const res: any = await getRecipeDetail(id)
      currentRecipe.value = res
    } catch {
      // ignore
    }
  }

  async function toggleFavoriteStatus(id: number) {
    try {
      await toggleFavorite(id)
      const idx = favorites.value.indexOf(id)
      if (idx > -1) {
        favorites.value.splice(idx, 1)
      } else {
        favorites.value.push(id)
      }
    } catch {
      // ignore
    }
  }

  return { categories, currentRecipe, favorites, fetchCategories, fetchRecipeDetail, toggleFavoriteStatus }
})
