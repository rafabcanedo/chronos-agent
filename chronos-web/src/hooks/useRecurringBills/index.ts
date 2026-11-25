import { useQuery } from '@tanstack/react-query'
import { getRecurringBills } from '@/services/recurring'
import type { GetRecurringParams } from '@/types'

export function useRecurringBills(params: GetRecurringParams = {}) {
  const query = useQuery({
    queryKey: ['recurring', params],
    queryFn: () => getRecurringBills(params),
  })

  return {
    bills: query.data?.data ?? [],
    meta: query.data?.meta,
    isLoading: query.isLoading,
    isError: query.isError,
    error: query.error,
  }
}
