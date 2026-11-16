import { CATEGORY_LABELS } from '@/contants'
import { Card, CardContent } from '@/components/ui/card'
import { StatusBadge } from '@/components/transactions/status-badge'
import type { TransactionCardProps } from '@/types'

export function TransactionCard({ transaction }: TransactionCardProps) {
  return (
    <Card>
      <CardContent className="flex flex-col gap-1">
        <div className="flex items-center justify-between">
          <p className="text-sm font-medium">{transaction.description}</p>
          <StatusBadge status={transaction.status} />
        </div>
        <p className="text-xs text-muted-foreground">{CATEGORY_LABELS[transaction.category]}</p>
        <p className="text-sm font-semibold mt-2">R$ {transaction.amount.toFixed(2)}</p>
      </CardContent>
    </Card>
  )
}
