import api from './index'

export interface UserEventPayload {
  event_name: string
  entity_type?: string
  entity_id?: number
  payload?: Record<string, unknown>
}

export function trackUserEvent(data: UserEventPayload) {
  return api.post('/user/events', data)
}
