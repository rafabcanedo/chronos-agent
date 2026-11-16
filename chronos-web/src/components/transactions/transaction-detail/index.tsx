import { CATEGORY_LABELS } from '@/contants'
import { StatusBadge } from '@/components/transactions/status-badge'
import type { TransactionDetailProps } from '@/types'

export function TransactionDetail({ transaction }: TransactionDetailProps) {
  return (
    <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-4">
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Description</p>
        <p className="text-sm font-medium">{transaction.description}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Amount</p>
        <p className="text-2xl font-bold">R$ {transaction.amount.toFixed(2)}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Category</p>
        <p className="text-sm">{CATEGORY_LABELS[transaction.category]}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Status</p>
        <StatusBadge status={transaction.status} />
      </div>
      {transaction.due_date && (
        <div>
          <p className="text-xs text-muted-foreground mb-0.5">Due Date</p>
          <p className="text-sm">{transaction.due_date}</p>
        </div>
      )}
      {transaction.paid_at && (
        <div>
          <p className="text-xs text-muted-foreground mb-0.5">Paid At</p>
          <p className="text-sm">{transaction.paid_at}</p>
        </div>
      )}
      {transaction.recurring_bill_id && (
        <div>
          <p className="text-xs text-muted-foreground mb-0.5">Recurring</p>
          <p className="text-sm">Yes</p>
        </div>
      )}
    </div>
  )
}
