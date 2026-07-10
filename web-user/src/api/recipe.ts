import api, { AI_GENERATION_TIMEOUT } from './index'
import type { RecipeFeedbackStatus, RecipeQuery } from '@/types/recipe'

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

export function setRecipeFeedback(id: number, type: keyof RecipeFeedbackStatus, source = 'detail') {
  return api.post(`/recipes/${id}/feedback`, { type, source }) as Promise<{ feedback: RecipeFeedbackStatus }>
}

export function deleteRecipeFeedback(id: number, type: keyof RecipeFeedbackStatus) {
  return api.delete(`/recipes/${id}/feedback/${type}`) as Promise<{ feedback: RecipeFeedbackStatus }>
}

export function getUserRecipeFeedback() {
  return api.get('/user/recipe-feedback')
}

export function getCategories() {
  return api.get('/categories')
}

export function getRecipeFilterOptions() {
  return api.get('/recipes/filter-options') as Promise<{ tastes: string[] }>
}

export function generateRecipeByAI(dishName: string): Promise<{ recipe: any; created: boolean }> {
  return api.post(
    '/recipes/generate-by-ai',
    { dish_name: dishName },
    { timeout: AI_GENERATION_TIMEOUT },
  ) as unknown as Promise<{ recipe: any; created: boolean }>
}
