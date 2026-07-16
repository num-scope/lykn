import { apiRequest } from '@/lib/api-client'
import type { LoginResponse, UserAccount } from '@/types/api'

export function login(username: string, password: string) {
  return apiRequest<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
}

export function getCurrentUser() {
  return apiRequest<UserAccount>('/auth/me')
}
