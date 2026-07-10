export interface Recipe {
  id: number
  title: string
  cover: string
  description: string
  cook_time: number
  difficulty: '简单' | '中等' | '困难'
  taste: string
  servings: number
  favorite_count: number
  is_favorited: boolean
  feedback?: RecipeFeedbackStatus
  category_id: number
  category_name: string
  health_tags: string[]
  ingredients: Ingredient[]
  seasonings: Seasoning[]
  steps: CookingStep[]
  nutrition: Nutrition
  tips?: string
}

export interface RecipeFeedbackStatus {
  cooked: boolean
  like: boolean
  dislike: boolean
  block: boolean
  normal: boolean
  too_complex: boolean
  too_long: boolean
  hard_to_buy: boolean
}

export interface Ingredient {
  name: string
  amount: string
  unit: string
  emoji?: string
}

export interface Seasoning {
  name: string
  amount: string
}

export interface CookingStep {
  step: number
  description: string
  image?: string
  tip?: string
}

export interface Nutrition {
  calories: number
  protein: number
  fat: number
  carbs: number
  fiber: number
}

export interface Category {
  id: number
  name: string
  icon: string
  count?: number
}

export interface RecipeQuery {
  keyword?: string
  category_id?: number
  taste?: string
  sort?: 'latest' | 'hot'
  cook_time?: string
  difficulty?: string
  health_tags?: string
  page?: number
  page_size?: number
}
