import { useState } from 'react'
import { cn } from 'cn'
import { PencilIcon, Trash2Icon } from 'lucide-react'
import { mockTransactions } from '@/mocks'
import { CATEGORY_LABELS } from '@/contants'
import { Navbar } from '@/components/dashboard/navbar'
import { AppPagination } from '@/components/chronos/pagination'
import { DetailDrawer } from '@/components/chronos/detail-drawer'
import { MonthSelector } from '@/components/transactions/month-selector'
import { StatusFilter } from '@/components/transactions/status-filter'
import { CategoryFilter } from '@/components/transactions/category-filter'
import { StatusBadge } from '@/components/transactions/status-badge'
import { TransactionDetail } from '@/components/transactions/transaction-detail'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
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
                    className={cn('cursor-pointer group', selected?.id === transaction.id && 'bg-muted')}
                    onClick={() => setSelected(transaction)}
                  >
                    <TableCell>{transaction.description}</TableCell>
                    <TableCell className="text-muted-foreground">{CATEGORY_LABELS[transaction.category]}</TableCell>
                    <TableCell className="font-medium">R$ {transaction.amount.toFixed(2)}</TableCell>
                    <TableCell><StatusBadge status={transaction.status} /></TableCell>
                    <TableCell className="text-muted-foreground relative">
                      <span className="transition-opacity group-hover:opacity-0">
                        {transaction.due_date ?? '—'}
                      </span>
                      <div className="absolute right-2 top-1/2 -translate-y-1/2 flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                        <Button variant="ghost" size="icon-sm" onClick={(e) => e.stopPropagation()}>
                          <PencilIcon className="size-3.5" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          className="text-destructive hover:text-destructive hover:bg-destructive/10"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <Trash2Icon className="size-3.5" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          <AppPagination />
        </div>

        <DetailDrawer
          open={!!selected}
          onOpenChange={(open) => !open && setSelected(null)}
          title="Details"
        >
          {selected && <TransactionDetail transaction={selected} />}
        </DetailDrawer>

      </main>
    </div>
  )
}
