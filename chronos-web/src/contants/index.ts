import { cva } from 'class-variance-authority'

export const CHAT_STORAGE_KEY = '@chronos:chat-messages'

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
