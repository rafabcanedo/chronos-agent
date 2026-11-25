import { TOKEN_STORAGE_KEY } from '@/contants'
import type { PaginatedResponse, PaginationMeta } from '@/types'

const BASE_URL = import.meta.env.VITE_API_URL

function getToken(): string | null {
  return localStorage.getItem(TOKEN_STORAGE_KEY)
}

async function request(path: string, options: RequestInit = {}): Promise<unknown> {
  const token = getToken()

  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })

  const json = await response.json()

  if (!response.ok) {
    throw new Error((json as { error?: string }).error ?? 'Unexpected error')
  }

  return json
}

export const api = {
  get: async <T>(path: string): Promise<T> => {
    const json = await request(path) as { data: T }
    return json.data
  },

  getPaginated: async <T>(path: string): Promise<PaginatedResponse<T>> => {
    const json = await request(path) as { data: T[]; meta: PaginationMeta }
    return { data: json.data, meta: json.meta }
  },

  post: async <T>(path: string, body: unknown): Promise<T> => {
    const json = await request(path, { method: 'POST', body: JSON.stringify(body) }) as { data: T }
    return json.data
  },
}
