import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { STATUS_LABELS } from '@/contants'
import { TransactionStatus } from '@/enums'

const items = [
  { value: 'all', label: 'All' },
  ...Object.values(TransactionStatus).map(s => ({ value: s, label: STATUS_LABELS[s] })),
]

export function StatusFilter() {
  return (
    <Select defaultValue="all" items={items}>
      <SelectTrigger className="w-36">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {items.map(item => (
            <SelectItem key={item.value} value={item.value}>{item.label}</SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
