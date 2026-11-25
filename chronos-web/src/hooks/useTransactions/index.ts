import { useQuery } from '@tanstack/react-query'
import { getTransactions, getTransactionsSummary } from '@/services/transactions'
import type { GetTransactionsParams } from '@/types'

export function useTransactions(params: GetTransactionsParams = {}) {
  const query = useQuery({
    queryKey: ['transactions', params],
    queryFn: () => getTransactions(params),
  })

  return {
    transactions: query.data?.data ?? [],
    meta: query.data?.meta,
    isLoading: query.isLoading,
    isError: query.isError,
    error: query.error,
  }
}

export function useTransactionsSummary(month?: string) {
  const query = useQuery({
    queryKey: ['transactions-summary', month],
    queryFn: () => getTransactionsSummary(month),
  })

  return {
    summary: query.data,
    isLoading: query.isLoading,
    isError: query.isError,
  }
}
