import { Category, TransactionStatus } from '@/enums'
import type { Transaction, RecurringBill } from '@/types'

export const mockRecurringBills: RecurringBill[] = [
  { id: '1', name: 'Aluguel',  default_amount: 1200.00, billing_day: 5,  category: Category.Housing,     active: true, created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '2', name: 'Internet', default_amount: 120.00,  billing_day: 10, category: Category.Housing,     active: true, created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '3', name: 'Netflix',  default_amount: 45.90,   billing_day: 15, category: Category.Leisure,     active: true, created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '4', name: 'Spotify',  default_amount: 21.90,   billing_day: 20, category: Category.Leisure,     active: true, created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '5', name: 'Gym',      default_amount: 89.90,   billing_day: 1,  category: Category.Leisure,     active: true, created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '6', name: 'Água',     default_amount: 55.00,   billing_day: 25, category: Category.Housing,     active: true, created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
]

export const mockTransactions: Transaction[] = [
  { id: '1', recurring_bill_id: '2',  description: 'Internet',      amount: 120.00,  category: Category.Housing,     due_date: '2026-11-10', status: TransactionStatus.Paid,    paid_at: '2026-11-10', created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-10T00:00:00Z' },
  { id: '2', recurring_bill_id: '1',  description: 'Aluguel',       amount: 1200.00, category: Category.Housing,     due_date: '2026-11-05', status: TransactionStatus.Overdue,  paid_at: null,         created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '3', recurring_bill_id: null, description: 'Almoço',        amount: 35.00,   category: Category.Restaurants, due_date: null,         status: TransactionStatus.Paid,    paid_at: '2026-11-08', created_at: '2026-11-08T00:00:00Z', updated_at: '2026-11-08T00:00:00Z' },
  { id: '4', recurring_bill_id: '3',  description: 'Netflix',       amount: 45.90,   category: Category.Leisure,     due_date: '2026-11-15', status: TransactionStatus.Pending,  paid_at: null,         created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
  { id: '5', recurring_bill_id: null, description: 'Supermercado',  amount: 280.00,  category: Category.Market,      due_date: null,         status: TransactionStatus.Paid,    paid_at: '2026-11-07', created_at: '2026-11-07T00:00:00Z', updated_at: '2026-11-07T00:00:00Z' },
  { id: '6', recurring_bill_id: '5',  description: 'Gym',           amount: 89.90,   category: Category.Leisure,     due_date: '2026-11-01', status: TransactionStatus.Overdue,  paid_at: null,         created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' },
]

export const mockStatusSummary = {
  pending: {
    count: mockTransactions.filter(transaction => transaction.status === TransactionStatus.Pending).length,
    total: mockTransactions.filter(transaction => transaction.status === TransactionStatus.Pending).reduce((sum, transaction) => sum + transaction.amount, 0),
  },
  paid: {
    count: mockTransactions.filter(transaction => transaction.status === TransactionStatus.Paid).length,
    total: mockTransactions.filter(transaction => transaction.status === TransactionStatus.Paid).reduce((sum, transaction) => sum + transaction.amount, 0),
  },
  overdue: {
    count: mockTransactions.filter(transaction => transaction.status === TransactionStatus.Overdue).length,
    total: mockTransactions.filter(transaction => transaction.status === TransactionStatus.Overdue).reduce((sum, transaction) => sum + transaction.amount, 0),
  },
}
