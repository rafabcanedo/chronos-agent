import { cn } from 'cn'
import { statusVariants, STATUS_LABELS } from '@/contants'
import type { StatusBadgeProps } from '@/types'

export function StatusBadge({ status }: StatusBadgeProps) {
  return (
    <span className={cn(statusVariants({ status }))}>
      {STATUS_LABELS[status]}
    </span>
  )
}
