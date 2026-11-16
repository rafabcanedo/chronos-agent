import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { CATEGORY_LABELS } from '@/contants'
import { Category } from '@/enums'

const items = [
  { value: 'all', label: 'All' },
  ...Object.values(Category).map(c => ({ value: c, label: CATEGORY_LABELS[c] })),
]

export function CategoryFilter() {
  return (
    <Select defaultValue="all" items={items}>
      <SelectTrigger className="w-40">
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
