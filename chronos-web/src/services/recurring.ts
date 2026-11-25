import { api } from '@/lib/api'
import type { GetRecurringParams, PaginatedResponse, RecurringBill } from '@/types'

export function getRecurringBills(params: GetRecurringParams = {}): Promise<PaginatedResponse<RecurringBill>> {
  const queryParams = new URLSearchParams()

  const entries: Record<string, string | undefined> = {
    active:    params.active !== 'all' ? params.active : undefined,
    category:  params.category !== 'all' ? params.category : undefined,
    name:      params.name || undefined,
    page:      params.page ? String(params.page) : undefined,
    page_size: params.page_size ? String(params.page_size) : undefined,
  }

  Object.entries(entries).forEach(([key, value]) => {
    if (value !== undefined) queryParams.set(key, value)
  })

  const queryString = queryParams.toString()
  return api.getPaginated<RecurringBill>(`/api/recurring${queryString ? `?${queryString}` : ''}`)
}
