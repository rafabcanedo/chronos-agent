import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ACTIVE_FILTER_ITEMS } from '@/contants'
import type { ActiveFilterProps } from '@/types'

export function ActiveFilter({ value, onChange }: ActiveFilterProps) {
  return (
    <Select value={value} onValueChange={(newValue) => onChange(newValue as 'true' | 'false' | 'all')}>
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
