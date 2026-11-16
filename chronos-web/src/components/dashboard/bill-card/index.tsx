import { CATEGORY_LABELS } from '@/contants'
import { Card, CardContent } from '@/components/ui/card'
import type { BillCardProps } from '@/types'

export function BillCard({ bill }: BillCardProps) {
  return (
    <Card>
      <CardContent className="flex flex-col gap-1">
        <p className="text-sm font-medium">{bill.name}</p>
        <p className="text-xs text-muted-foreground">{CATEGORY_LABELS[bill.category]}</p>
        <div className="flex items-center justify-between mt-2">
          <span className="text-sm font-semibold">R$ {bill.default_amount.toFixed(2)}</span>
          <span className="text-xs text-muted-foreground">Day {bill.billing_day}</span>
        </div>
      </CardContent>
    </Card>
  )
}
