import { Card, CardContent } from '@/components/ui/card'
import { mockRecurringBills } from '@/mocks'

export function RecurringStats() {
  const active = mockRecurringBills.filter(b => b.active)
  const total = active.reduce((sum, b) => sum + b.default_amount, 0)
  const avg = active.length > 0 ? total / active.length : 0

  return (
    <div className="grid grid-cols-3 gap-4">
      <Card>
        <CardContent className="flex flex-col gap-1">
          <p className="text-xs text-muted-foreground">Total Monthly</p>
          <p className="text-2xl font-bold">R$ {total.toFixed(2)}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="flex flex-col gap-1">
          <p className="text-xs text-muted-foreground">Active Bills</p>
          <p className="text-2xl font-bold">{active.length}</p>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="flex flex-col gap-1">
          <p className="text-xs text-muted-foreground">Average</p>
          <p className="text-2xl font-bold">R$ {avg.toFixed(2)}</p>
        </CardContent>
      </Card>
    </div>
  )
}
