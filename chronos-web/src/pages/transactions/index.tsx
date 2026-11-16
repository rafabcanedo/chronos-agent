import { useState } from 'react'
import { cn } from 'cn'
import { XIcon } from 'lucide-react'
import { mockTransactions } from '@/mocks'
import { CATEGORY_LABELS } from '@/contants'
import { Navbar } from '@/components/dashboard/navbar'
import { AppPagination } from '@/components/chronos/pagination'
import { MonthSelector } from '@/components/transactions/month-selector'
import { StatusFilter } from '@/components/transactions/status-filter'
import { CategoryFilter } from '@/components/transactions/category-filter'
import { StatusBadge } from '@/components/transactions/status-badge'
import { TransactionDetail } from '@/components/transactions/transaction-detail'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Drawer, DrawerClose, DrawerContent, DrawerFooter, DrawerHeader, DrawerTitle } from '@/components/ui/drawer'
import { Button } from '@/components/ui/button'
import type { Transaction } from '@/types'

export default function Transactions() {
  const [selected, setSelected] = useState<Transaction | null>(null)

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main className="flex-1 flex flex-col gap-4 p-6">

        <div className="flex items-center justify-between">
          <h1 className="text-xl font-heading font-semibold">Transactions</h1>
          <MonthSelector />
        </div>

        <div className="flex gap-2">
          <StatusFilter />
          <CategoryFilter />
        </div>

        <div className="flex flex-col gap-4">
          <div className="border rounded-lg overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Description</TableHead>
                  <TableHead>Category</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Due Date</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {mockTransactions.map(transaction => (
                  <TableRow
                    key={transaction.id}
                    className={cn('cursor-pointer', selected?.id === transaction.id && 'bg-muted')}
                    onClick={() => setSelected(transaction)}
                  >
                    <TableCell>{transaction.description}</TableCell>
                    <TableCell className="text-muted-foreground">{CATEGORY_LABELS[transaction.category]}</TableCell>
                    <TableCell className="font-medium">R$ {transaction.amount.toFixed(2)}</TableCell>
                    <TableCell><StatusBadge status={transaction.status} /></TableCell>
                    <TableCell className="text-muted-foreground">{transaction.due_date ?? '—'}</TableCell>
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
              <DrawerTitle>Details</DrawerTitle>
              <DrawerClose render={<Button variant="ghost" size="icon" />}>
                <XIcon className="size-4" />
              </DrawerClose>
            </DrawerHeader>
            {selected && <TransactionDetail transaction={selected} />}
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
