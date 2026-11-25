import { Card, CardContent } from '@/components/ui/card'
import type { RecurringStatsProps } from '@/types'

export function RecurringStats({ bills }: RecurringStatsProps) {
  const activeBills = bills.filter(bill => bill.active)
  const totalMonthly = activeBills.reduce((sum, bill) => sum + bill.default_amount, 0)
  const average = activeBills.length > 0 ? totalMonthly / activeBills.length : 0

  return (
    <div className="grid grid-cols-3 gap-4">
      <Card>
        <CardContent className="flex flex-col gap-1">
          <p className="text-xs text-muted-foreground">Total Monthly</p>
          <p className="text-2xl font-bold">R$ {totalMonthly.toFixed(2)}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="flex flex-col gap-1">
          <p className="text-xs text-muted-foreground">Active Bills</p>
          <p className="text-2xl font-bold">{activeBills.length}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="flex flex-col gap-1">
          <p className="text-xs text-muted-foreground">Average</p>
          <p className="text-2xl font-bold">R$ {average.toFixed(2)}</p>
        </CardContent>
      </Card>
    </div>
  )
}
