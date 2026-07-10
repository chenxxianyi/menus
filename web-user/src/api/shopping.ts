import api, { AI_GENERATION_TIMEOUT } from './index'

export interface ShoppingItem {
  name: string
  amount: string
  emoji: string
  category: string
  price: number
  checked: boolean
  status?: 'pending' | 'bought' | 'owned'
}

export interface ShoppingList {
  id: number
  name: string
  items: ShoppingItem[]
}

export interface GenerateShoppingByDishResult {
  list?: {
    id: number
    name: string
    items_json: ShoppingItem[]
  }
  recipe?: {
    id: number
    title: string
    description?: string
  }
  items: ShoppingItem[]
  preview?: boolean
}

export interface GenerateShoppingByRecipesResult {
  list?: {
    id: number
    name: string
    items_json: ShoppingItem[]
  }
  recipes?: Array<{
    id: number
    title: string
  }>
  items: ShoppingItem[]
  preview?: boolean
}

export interface DeleteShoppingItemsResult {
  list_id: number
  deleted_count: number
  items_json: ShoppingItem[]
}

export function getShoppingLists() {
  return api.get('/shopping-list')
}

export function createShoppingList(data: { name: string; items: ShoppingItem[] }) {
  return api.post('/shopping-list', { name: data.name, items_json: data.items })
}

export function updateShoppingList(id: number, data: { name: string; items: ShoppingItem[] }) {
  return api.put(`/shopping-list/${id}`, { name: data.name, items_json: data.items })
}

export function deleteShoppingItems(id: number, indices: number[]): Promise<DeleteShoppingItemsResult> {
  return api.delete(`/shopping-list/${id}/items`, {
    data: { indices },
  }) as unknown as Promise<DeleteShoppingItemsResult>
}

export function deleteShoppingList(id: number) {
  return api.delete(`/shopping-list/${id}`)
}

export function generateShoppingListByDish(dishName: string, preview = false): Promise<GenerateShoppingByDishResult> {
  return api.post('/shopping-list/generate-by-dish', {
    dish_name: dishName,
    preview,
  }) as unknown as Promise<GenerateShoppingByDishResult>
}

export function generateShoppingListByRecipe(recipeId: number, preview = false): Promise<GenerateShoppingByDishResult> {
  return api.post('/shopping-list/generate-by-recipe', {
    recipe_id: recipeId,
    preview,
    merge_current: true,
  }) as unknown as Promise<GenerateShoppingByDishResult>
}

export function generateShoppingListByRecipes(
  recipeIds: number[],
  name: string,
  preview = false,
): Promise<GenerateShoppingByRecipesResult> {
  return api.post('/shopping-list/generate-by-recipes', {
    recipe_ids: recipeIds,
    name,
    preview,
    merge_current: true,
  }) as unknown as Promise<GenerateShoppingByRecipesResult>
}

export function generateShoppingListByAI(
  dishName: string,
  preview = false,
): Promise<GenerateShoppingByDishResult & { source: 'ai' }> {
  return api.post(
    '/shopping-list/generate-by-ai',
    { dish_name: dishName, preview },
    { timeout: AI_GENERATION_TIMEOUT },
  ) as unknown as Promise<GenerateShoppingByDishResult & { source: 'ai' }>
}
