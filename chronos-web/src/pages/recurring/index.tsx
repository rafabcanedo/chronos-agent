import { useState } from 'react'
import { cn } from 'cn'
import { XIcon } from 'lucide-react'
import { mockRecurringBills } from '@/mocks'
import { CATEGORY_LABELS } from '@/contants'
import { Navbar } from '@/components/dashboard/navbar'
import { AppPagination } from '@/components/chronos/pagination'
import { RecurringStats } from '@/components/recurring/recurring-stats'
import { BillSearch } from '@/components/recurring/bill-search'
import { CategoryFilter } from '@/components/transactions/category-filter'
import { ActiveFilter } from '@/components/recurring/active-filter'
import { ActiveBadge } from '@/components/recurring/active-badge'
import { BillDetail } from '@/components/recurring/bill-detail'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Drawer, DrawerClose, DrawerContent, DrawerFooter, DrawerHeader, DrawerTitle } from '@/components/ui/drawer'
import { Button } from '@/components/ui/button'
import type { RecurringBill } from '@/types'

export default function Recurring() {
  const [selected, setSelected] = useState<RecurringBill | null>(null)

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main className="flex-1 flex flex-col gap-4 p-6">
        <h1 className="text-xl font-heading font-semibold">Recurring Bills</h1>

        <RecurringStats />

        <div className="flex gap-2">
          <BillSearch />
          <CategoryFilter />
          <ActiveFilter />
        </div>

        <div className="flex flex-col gap-4">
          <div className="border rounded-lg overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Category</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Billing Day</TableHead>
                  <TableHead>Active</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {mockRecurringBills.map(bill => (
                  <TableRow
                    key={bill.id}
                    className={cn('cursor-pointer', selected?.id === bill.id && 'bg-muted')}
                    onClick={() => setSelected(bill)}
                  >
                    <TableCell className="font-medium">{bill.name}</TableCell>
                    <TableCell className="text-muted-foreground">{CATEGORY_LABELS[bill.category]}</TableCell>
                    <TableCell className="font-medium">R$ {bill.default_amount.toFixed(2)}</TableCell>
                    <TableCell className="text-muted-foreground">Day {bill.billing_day}</TableCell>
                    <TableCell><ActiveBadge active={bill.active} /></TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          <AppPagination />
        </div>

        <Drawer open={!!selected} onOpenChange={(open) => !open && setSelected(null)} swipeDirection="right">
          <DrawerContent>
            <DrawerHeader className="flex-row items-center justify-between">
              <DrawerTitle>Bill Details</DrawerTitle>
              <DrawerClose render={<Button variant="ghost" size="icon" />}>
                <XIcon className="size-4" />
              </DrawerClose>
            </DrawerHeader>
            {selected && <BillDetail bill={selected} />}
            <DrawerFooter>
              <DrawerClose render={<Button variant="outline" className="w-full" />}>
                Back to list
              </DrawerClose>
            </DrawerFooter>
          </DrawerContent>
        </Drawer>
      </main>
    </div>
  )
}
