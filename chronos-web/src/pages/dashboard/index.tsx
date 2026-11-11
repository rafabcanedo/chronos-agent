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
        <div className="grid grid-cols-3 gap-3">
          {mockRecurringBills.map(bill => (
            <BillCard key={bill.id} bill={bill} />
          ))}
        </div>
        <div className="grid grid-cols-3 gap-3">
          {mockTransactions.map(transaction => (
            <TransactionCard key={transaction.id} transaction={transaction} />
          ))}
        </div>
      </main>
    </div>
  )
}
