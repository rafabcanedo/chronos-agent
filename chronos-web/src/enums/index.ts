export const Category = {
  Mercado: 'mercado',
  Restaurantes: 'restaurantes',
  Moradia: 'moradia',
  Transporte: 'transporte',
  Lazer: 'lazer',
  Outros: 'outros',
} as const

export type Category = typeof Category[keyof typeof Category]

export const TransactionStatus = {
  Pending: 'pending',
  Paid: 'paid',
  Overdue: 'overdue',
} as const

export type TransactionStatus = typeof TransactionStatus[keyof typeof TransactionStatus]
