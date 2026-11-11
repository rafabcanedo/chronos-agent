import { cva } from 'class-variance-authority'
import { cn } from 'cn'
import { CATEGORY_LABELS, STATUS_LABELS } from '@/contants'
import { Card, CardContent } from '@/components/ui/card'
import type { TransactionCardProps } from '@/types'

const statusBadgeVariants = cva(
  "rounded-full px-2 py-0.5 text-xs font-medium",
  {
    variants: {
      status: {
        pending: "bg-yellow-50 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400",
        paid:    "bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400",
        overdue: "bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-400",
      },
    },
  }
)

export function TransactionCard({ transaction }: TransactionCardProps) {
  return (
    <Card>
      <CardContent className="flex flex-col gap-1">
        <div className="flex items-center justify-between">
          <p className="text-sm font-medium">{transaction.description}</p>
          <span className={cn(statusBadgeVariants({ status: transaction.status }))}>
            {STATUS_LABELS[transaction.status]}
          </span>
        </div>
        <p className="text-xs text-muted-foreground">{CATEGORY_LABELS[transaction.category]}</p>
        <p className="text-sm font-semibold mt-2">R$ {transaction.amount.toFixed(2)}</p>
      </CardContent>
    </Card>
  )
}
