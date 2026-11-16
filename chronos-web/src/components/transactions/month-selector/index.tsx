import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function MonthSelector() {
  return (
    <div className="flex items-center gap-1">
      <Button variant="ghost" size="icon">
        <ChevronLeftIcon className="size-4" />
      </Button>
      <span className="text-sm font-medium w-36 text-center">November 2026</span>
      <Button variant="ghost" size="icon">
        <ChevronRightIcon className="size-4" />
      </Button>
    </div>
  )
}
