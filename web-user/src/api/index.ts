import axios from 'axios'
import router from '@/router'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
})

let handlingUnauthorized = false

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => {
    const { code, message, data } = response.data
    if (code === 0) {
      return data
    }
    return Promise.reject(new Error(message))
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      if (!handlingUnauthorized && router.currentRoute.value.name !== 'Login') {
        handlingUnauthorized = true
        router.replace({ name: 'Login' }).finally(() => {
          handlingUnauthorized = false
        })
      }
    }
    return Promise.reject(error)
  }
)

export default api
