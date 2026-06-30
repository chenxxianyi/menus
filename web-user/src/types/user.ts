export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  gender?: number
  email?: string
  phone?: string
  bio?: string
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirm_password: string
}

export interface UserPreferences {
  taste_preference: string[]
  health_goal: string
  avoid_ingredients: string[]
  favorite_ingredients: string[]
  cook_time_preference: string
  default_servings: number
}
