import axios from 'axios'
import router from '@/router'

export const DEFAULT_REQUEST_TIMEOUT = 10000
export const AI_GENERATION_TIMEOUT = 75000

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: DEFAULT_REQUEST_TIMEOUT,
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
    if (axios.isAxiosError(error)) {
      const isTimeout =
        error.code === 'ECONNABORTED' ||
        error.code === 'ETIMEDOUT' ||
        error.message.toLowerCase().includes('timeout')
      if (isTimeout) {
        return Promise.reject(new Error('请求等待时间较长，请稍后重试；如果是 AI 生成，请检查 AI 服务响应速度。'))
      }
      if (!error.response) {
        return Promise.reject(new Error('网络连接异常，请确认后端服务正常运行。'))
      }
    }
    return Promise.reject(error)
  }
)

export default api
