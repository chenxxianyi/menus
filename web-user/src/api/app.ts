import api from './index'

export interface AboutFeature {
  title: string
  description: string
  icon?: string
  color?: string
  bg?: string
}

export interface AboutInfo {
  name?: string
  subtitle?: string
  description?: string
  slogan?: string
  version?: string
  email?: string
  wechat?: string
  terms_url?: string
  privacy_url?: string
  features?: AboutFeature[]
}

export function getAboutInfo() {
  return api.get('/about') as Promise<AboutInfo>
}
