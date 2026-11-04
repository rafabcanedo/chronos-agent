import { Category, TransactionStatus } from '@/enums'

export interface Transaction {
  id: string
  recurring_bill_id: string | null
  description: string
  amount: number
  category: Category
  due_date: string | null
  status: TransactionStatus
  paid_at: string | null
  created_at: string
  updated_at: string
}

export interface RecurringBill {
  id: string
  name: string
  default_amount: number
  billing_day: number
  category: Category
  active: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: string
  username: string
  email: string
  created_at: string
}

export interface ApiResponse<T> {
  data: T
}

export interface ApiError {
  error: string
}
