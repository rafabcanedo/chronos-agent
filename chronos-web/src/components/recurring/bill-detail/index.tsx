import { CATEGORY_LABELS } from '@/contants'
import { ActiveBadge } from '@/components/recurring/active-badge'
import type { BillDetailProps } from '@/types'

export function BillDetail({ bill }: BillDetailProps) {
  return (
    <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-4">
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Name</p>
        <p className="text-sm font-medium">{bill.name}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Amount</p>
        <p className="text-2xl font-bold">R$ {bill.default_amount.toFixed(2)}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Category</p>
        <p className="text-sm">{CATEGORY_LABELS[bill.category]}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Billing Day</p>
        <p className="text-sm">Day {bill.billing_day}</p>
      </div>
      <div>
        <p className="text-xs text-muted-foreground mb-0.5">Status</p>
        <ActiveBadge active={bill.active} />
      </div>
    </div>
  )
}
