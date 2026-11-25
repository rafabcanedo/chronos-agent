import { api } from '@/lib/api'
import type { LoginPayload, LoginResponse } from '@/types'

export function login(payload: LoginPayload): Promise<LoginResponse> {
  return api.post<LoginResponse>('/api/auth/login', payload)
}
