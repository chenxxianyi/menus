import api from './index'
import type { RecipeQuery } from '@/types/recipe'

export function getRecipes(params: RecipeQuery) {
  return api.get('/recipes', { params })
}

export function getRecipeDetail(id: number) {
  return api.get(`/recipes/${id}`)
}

export function toggleFavorite(id: number) {
  return api.post(`/recipes/${id}/favorite`)
}

export function removeFavorite(id: number) {
  return api.delete(`/recipes/${id}/favorite`)
}

export function getCategories() {
  return api.get('/categories')
}

export function getRecipeFilterOptions() {
  return api.get('/recipes/filter-options') as Promise<{ tastes: string[] }>
}
