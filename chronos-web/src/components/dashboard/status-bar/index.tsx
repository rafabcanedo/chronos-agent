import { cva } from 'class-variance-authority'
import { cn } from 'cn'
import { useTransactionsSummary } from '@/hooks/useTransactions'
import { STATUS_LABELS } from '@/contants'

const statusPillVariants = cva(
  "inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium",
  {
    variants: {
      status: {
        pending: "border-yellow-200 bg-yellow-50 text-yellow-700 dark:border-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400",
        paid:    "border-green-200 bg-green-50 text-green-700 dark:border-green-800 dark:bg-green-900/20 dark:text-green-400",
        overdue: "border-red-200 bg-red-50 text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400",
      },
    },
  }
)

export function StatusBar() {
  const { summary } = useTransactionsSummary()

  const pending = summary?.pending ?? { count: 0, total: 0 }
  const paid    = summary?.paid    ?? { count: 0, total: 0 }
  const overdue = summary?.overdue ?? { count: 0, total: 0 }

  return (
    <div className="flex items-center gap-3">
      <span className="text-sm font-medium text-muted-foreground">Status</span>
      <div className="flex gap-2">
        <span className={cn(statusPillVariants({ status: 'pending' }))}>
          {STATUS_LABELS.pending} <strong>{pending.count}</strong>
        </span>
        <span className={cn(statusPillVariants({ status: 'paid' }))}>
          {STATUS_LABELS.paid} <strong>{paid.count}</strong>
        </span>
        <span className={cn(statusPillVariants({ status: 'overdue' }))}>
          {STATUS_LABELS.overdue} <strong>{overdue.count}</strong>
        </span>
      </div>
    </div>
  )
}
