import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ACTIVE_FILTER_ITEMS } from '@/contants'

export function ActiveFilter() {
  return (
    <Select defaultValue="all" items={ACTIVE_FILTER_ITEMS}>
      <SelectTrigger className="w-32">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {ACTIVE_FILTER_ITEMS.map(item => (
            <SelectItem key={item.value} value={item.value}>{item.label}</SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
