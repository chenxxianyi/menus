import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { CoupleInfo, CoupleOrder } from '@/types/couple'
import { getCoupleInfo, getCoupleOrders } from '@/api/couple'

export const useCoupleStore = defineStore('couple', () => {
  const coupleInfo = ref<CoupleInfo | null>(null)
  const orders = ref<CoupleOrder[]>([])
  const isBound = ref(false)

  async function fetchCoupleInfo() {
    try {
      const res: any = await getCoupleInfo()
      if (res) {
        coupleInfo.value = res
        isBound.value = true
      } else {
        coupleInfo.value = null
        isBound.value = false
      }
    } catch {
      coupleInfo.value = null
      isBound.value = false
    }
  }

  async function fetchOrders(mealDate?: string) {
    try {
      const res: any = await getCoupleOrders(mealDate)
      orders.value = res || []
    } catch {
      orders.value = []
    }
  }

  function clear() {
    coupleInfo.value = null
    isBound.value = false
    orders.value = []
  }

  return { coupleInfo, orders, isBound, fetchCoupleInfo, fetchOrders, clear }
})
