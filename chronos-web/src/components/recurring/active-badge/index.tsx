import { cn } from 'cn'
import { activeVariants, ACTIVE_LABELS } from '@/contants'
import type { ActiveBadgeProps } from '@/types'

export function ActiveBadge({ active }: ActiveBadgeProps) {
  return (
    <span className={cn(activeVariants({ active }))}>
      {ACTIVE_LABELS[String(active) as 'true' | 'false']}
    </span>
  )
}
