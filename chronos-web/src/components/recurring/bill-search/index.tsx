import { SearchIcon } from 'lucide-react'
import { InputGroup, InputGroupAddon, InputGroupText, InputGroupInput } from '@/components/ui/input-group'
import type { BillSearchProps } from '@/types'

export function BillSearch({ value, onChange }: BillSearchProps) {
  return (
    <InputGroup className="w-64">
      <InputGroupAddon>
        <InputGroupText><SearchIcon /></InputGroupText>
      </InputGroupAddon>
      <InputGroupInput
        placeholder="Search bills..."
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    </InputGroup>
  )
}
