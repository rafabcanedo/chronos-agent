import { SearchIcon } from 'lucide-react'
import { InputGroup, InputGroupAddon, InputGroupText, InputGroupInput } from '@/components/ui/input-group'

export function BillSearch() {
  return (
    <InputGroup className="w-64">
      <InputGroupAddon>
        <InputGroupText><SearchIcon /></InputGroupText>
      </InputGroupAddon>
      <InputGroupInput placeholder="Search bills..." />
    </InputGroup>
  )
}
