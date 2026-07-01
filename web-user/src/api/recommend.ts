import api, { AI_GENERATION_TIMEOUT } from './index'
import type {
  IngredientOption,
  RecommendByIngredientsRequest,
  RecommendResultPayload,
} from '@/types/recommend'

export function getIngredientOptions(params?: { keyword?: string; category?: string }) {
  return api.get('/ingredients/options', { params }) as Promise<{ list: IngredientOption[] }>
}

export function recommendByIngredients(data: RecommendByIngredientsRequest) {
  return api.post('/recommend/by-ingredients', data) as Promise<RecommendResultPayload>
}

export function recommendByScene(data: {
  scene: string
  meal_type: string
  people_count: number
  cook_time_preference?: string
  health_goal?: string
  taste_preference?: string[]
  avoid_ingredients?: string[]
}) {
  return api.post('/recommend/menu', data)
}

export function recommendSceneByAI(data: {
  scene: string
  meal_type: string
  people_count: number
  cook_time_preference?: string
  health_goal?: string
  taste_preference?: string[]
  avoid_ingredients?: string[]
}) {
  return api.post('/recommend/menu-ai', data, { timeout: AI_GENERATION_TIMEOUT })
}
