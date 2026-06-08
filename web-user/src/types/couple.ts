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
}

export interface GenerateShoppingListResult {
  orders: CoupleOrder[]
  shopping_list: ShoppingListItem[]
  total_items: number
}
