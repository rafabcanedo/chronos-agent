import { cva } from 'class-variance-authority'
import { Category, TransactionStatus } from '@/enums'

export const CHAT_STORAGE_KEY = '@chronos:chat-messages'
export const TOKEN_STORAGE_KEY = '@chronos:token'

export const statusVariants = cva('rounded-full px-2 py-0.5 text-xs font-medium', {
  variants: {
    status: {
      pending: 'bg-yellow-50 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400',
      paid:    'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400',
      overdue: 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-400',
    },
  },
})

export const CATEGORY_LABELS: Record<string, string> = {
  market:      'Market',
  restaurants: 'Restaurants',
  housing:     'Housing',
  transport:   'Transport',
  leisure:     'Leisure',
  other:       'Other',
}

export const STATUS_LABELS: Record<string, string> = {
  pending: 'Pending',
  paid:    'Paid',
  overdue: 'Overdue',
}

export const activeVariants = cva('rounded-full px-2 py-0.5 text-xs font-medium', {
  variants: {
    active: {
      true:  'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400',
      false: 'bg-gray-100 text-gray-500 dark:bg-gray-800/30 dark:text-gray-400',
    },
  },
})

export const ACTIVE_LABELS: Record<'true' | 'false', string> = {
  true:  'Active',
  false: 'Inactive',
}

export const STATUS_FILTER_ITEMS = [
  { value: 'all', label: 'All' },
  ...Object.values(TransactionStatus).map(s => ({ value: s, label: STATUS_LABELS[s] })),
]

export const CATEGORY_FILTER_ITEMS = [
  { value: 'all', label: 'All' },
  ...Object.values(Category).map(c => ({ value: c, label: CATEGORY_LABELS[c] })),
]

export const ACTIVE_FILTER_ITEMS = [
  { value: 'all',   label: 'All' },
  { value: 'true',  label: 'Active' },
  { value: 'false', label: 'Inactive' },
]
