import type { Recipe } from './recipe'

export type RecommendMode = 'ingredients' | 'taste' | 'scene' | 'fridge'

export interface IngredientOption {
  id: number
  name: string
  category?: string
}

export interface RecommendByIngredientsRequest {
  ingredients: string[]
  mode?: 'ingredients' | 'fridge'
  limit?: number
}

export interface RecommendRecipeResult {
  recipe: Partial<Recipe>
  match_rate: number
  matched_ingredients: string[]
  missing_ingredients: string[]
  reason?: string
}

export interface RecommendResultPayload {
  list: RecommendRecipeResult[]
  total: number
}

export interface SceneOption {
  key: string
  title: string
  description: string
  meal_type: string
  people_count: number
  cook_time_preference?: string
  health_goal?: string
}
