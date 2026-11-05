export const Category = {
  Market: 'market',
  Restaurants: 'restaurants',
  Housing: 'housing',
  Transport: 'transport',
  Leisure: 'leisure',
  Other: 'other',
} as const

export type Category = typeof Category[keyof typeof Category]

export const TransactionStatus = {
  Pending: 'pending',
  Paid: 'paid',
  Overdue: 'overdue',
} as const

export type TransactionStatus = typeof TransactionStatus[keyof typeof TransactionStatus]
