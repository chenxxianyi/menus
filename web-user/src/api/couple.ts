import api from './index'

export function getInviteCode() {
  return api.get('/couple/invite-code')
}

export function bindCouple(inviteCode: string) {
  return api.post('/couple/bind', { invite_code: inviteCode })
}

export function getCoupleInfo() {
  return api.get('/couple/info')
}

export function unbindCouple() {
  return api.post('/couple/unbind')
}

export function setCoupleName(coupleName: string) {
  return api.put('/couple/name', { couple_name: coupleName })
}

export function createCoupleOrder(data: {
  dish_name: string
  recipe_id?: number
  meal_type?: string
  meal_date?: string
  note?: string
}) {
  return api.post('/couple/orders', data)
}

export function getCoupleOrders(mealDate?: string) {
  return api.get('/couple/orders', { params: { meal_date: mealDate } })
}

export function updateCoupleOrderStatus(id: number, status: number) {
  return api.put(`/couple/orders/${id}`, { status })
}

export function deleteCoupleOrder(id: number) {
  return api.delete(`/couple/orders/${id}`)
}

export function generateShoppingList(mealDate?: string, mealType?: string, saveShared = false) {
  return api.post('/couple/orders/generate-shopping-list', {
    meal_date: mealDate,
    meal_type: mealType,
    save_shared: saveShared,
  })
}
