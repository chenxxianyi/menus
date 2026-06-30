import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo } from '@/types/user'
import { login as loginApi } from '@/api/auth'
import { getUserInfo } from '@/api/user'
import type { LoginRequest } from '@/types/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(data: LoginRequest) {
    const res: any = await loginApi(data)
    token.value = res.token
    localStorage.setItem('token', res.token)
    await fetchUserInfo()
  }

  async function fetchUserInfo() {
    try {
      const res: any = await getUserInfo()
      userInfo.value = res
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return { token, userInfo, isLoggedIn, login, fetchUserInfo, logout }
})
