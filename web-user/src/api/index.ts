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
    return Promise.reject(new Error(toFriendlyError(message)))
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
      const serverMessage = error.response?.data?.message
      if (serverMessage) {
        return Promise.reject(new Error(toFriendlyError(serverMessage)))
      }
    }
    return Promise.reject(new Error(toFriendlyError(error?.message || '操作失败，请稍后重试。')))
  }
)

function toFriendlyError(message: unknown) {
  const text = String(message || '').trim()
  if (!text) return '操作失败，请稍后重试。'
  if (/timeout|timed out|deadline/i.test(text)) return '请求等待时间较长，请稍后重试。'
  if (/network|failed to fetch|connection/i.test(text)) return '网络连接异常，请确认后端服务正常运行。'
  if (/unauthorized|token|jwt/i.test(text)) return '登录状态已过期，请重新登录。'
  if (/AI.*未配置|未配置.*AI/i.test(text)) return 'AI 服务还未配置，请先完成 AI 配置。'
  if (/upstream|unavailable|temporarily/i.test(text)) return 'AI 服务暂时不可用，请稍后重试。'
  if (/invalid.*json|format|格式/i.test(text)) return text.includes('格式') ? text : '返回内容格式异常，请稍后重试。'
  if (/ECONN|ENOTFOUND|socket|AxiosError/i.test(text)) return '服务连接失败，请确认后端服务正常运行。'
  return text
}

export default api
