export interface CoupleInfo {
  couple_id: number
  couple_name: string
  invite_code: string
  partner: CouplePartner
  bound_at: string
}

export interface CouplePartner {
  id: number
  nickname: string
  avatar: string
}

export interface CoupleOrder {
  id: number
  couple_id: number
  user_id: number
  recipe_id: number | null
  dish_name: string
  meal_type: string
  meal_date: string
  note: string
  status: number
  created_at: string
  user?: CouplePartner
  recipe?: {
    id: number
    title: string
    cover?: string
    cook_time?: number
    difficulty?: string
    people_count?: number
  }
}

export interface ShoppingListItem {
  name: string
  amount: string
  category: string
  emoji?: string
  price?: number
  checked?: boolean
  status?: 'pending' | 'bought' | 'owned'
}

export interface CoupleMenuDish {
  recipe_id: number
  name: string
  reason: string
  source: 'overlap' | 'compromise' | string
  cook_time?: number
  difficulty?: string
}

export interface GenerateShoppingListResult {
  orders: CoupleOrder[]
  agreed_dishes?: CoupleMenuDish[]
  compromise_dishes?: CoupleMenuDish[]
  shopping_list: ShoppingListItem[]
  shared_lists?: Array<{
    id: number
    name: string
    user_id: number
  }>
  saved_shared?: boolean
  total_items: number
}
