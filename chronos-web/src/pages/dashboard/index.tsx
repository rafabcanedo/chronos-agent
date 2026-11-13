import { Link } from 'react-router'
import { mockRecurringBills, mockTransactions } from '@/mocks'
import { DashboardNavbar } from '@/components/dashboard/DashboardNavbar'
import { StatusBar } from '@/components/dashboard/StatusBar'
import { BillCard } from '@/components/dashboard/BillCard'
import { TransactionCard } from '@/components/dashboard/TransactionCard'

export default function Dashboard() {
  return (
    <div className="min-h-screen flex flex-col">
      <DashboardNavbar />
      <main className="flex-1 flex flex-col gap-6 p-6">
        <StatusBar />
        <section className="flex flex-col gap-3 border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-semibold">Recurring Bills</h2>
            <Link to="/recurring" className="text-xs text-muted-foreground hover:text-foreground transition-colors">
              See all
            </Link>
          </div>
          <div className="grid grid-rows-2 grid-flow-col auto-cols-[calc((100%-2*12px)/3)] overflow-x-auto gap-3 pb-1">
            {mockRecurringBills.map(bill => (
              <BillCard key={bill.id} bill={bill} />
            ))}
          </div>
        </section>

        <section className="flex flex-col gap-3 border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-semibold">Transactions</h2>
            <Link to="/transactions" className="text-xs text-muted-foreground hover:text-foreground transition-colors">
              See all
            </Link>
          </div>
          <div className="grid grid-rows-2 grid-flow-col auto-cols-[calc((100%-2*12px)/3)] overflow-x-auto gap-3 pb-1">
            {mockTransactions.map(transaction => (
              <TransactionCard key={transaction.id} transaction={transaction} />
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}
