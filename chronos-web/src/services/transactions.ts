import { api } from '@/lib/api'
import type { GetTransactionsParams, PaginatedResponse, Transaction, TransactionsSummary } from '@/types'

export function getTransactions(params: GetTransactionsParams = {}): Promise<PaginatedResponse<Transaction>> {
  const queryParams = new URLSearchParams()

  const entries: Record<string, string | undefined> = {
    month:     params.month || undefined,
    status:    params.status !== 'all' ? params.status : undefined,
    category:  params.category !== 'all' ? params.category : undefined,
    page:      params.page ? String(params.page) : undefined,
    page_size: params.page_size ? String(params.page_size) : undefined,
  }

  Object.entries(entries).forEach(([key, value]) => {
    if (value !== undefined) queryParams.set(key, value)
  })

  const queryString = queryParams.toString()
  return api.getPaginated<Transaction>(`/api/transactions${queryString ? `?${queryString}` : ''}`)
}

export function getTransactionsSummary(month?: string): Promise<TransactionsSummary> {
  const queryString = month ? `?month=${month}` : ''
  return api.get<TransactionsSummary>(`/api/transactions/summary${queryString}`)
}
