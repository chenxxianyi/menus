import { trackUserEvent, type UserEventPayload } from '@/api/event'

export function trackEvent(data: UserEventPayload) {
  if (!localStorage.getItem('token')) return
  trackUserEvent(data).catch(() => {
    // 埋点失败不影响用户主流程。
  })
}
