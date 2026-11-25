import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { formatMonthDisplay, navigateMonth } from '@/lib/date'
import type { MonthSelectorProps } from '@/types'

export function MonthSelector({ value, onChange }: MonthSelectorProps) {
  return (
    <div className="flex items-center gap-1">
      <Button variant="ghost" size="icon" onClick={() => onChange(navigateMonth(value, -1))}>
        <ChevronLeftIcon className="size-4" />
      </Button>
      <span className="text-sm font-medium w-36 text-center">{formatMonthDisplay(value)}</span>
      <Button variant="ghost" size="icon" onClick={() => onChange(navigateMonth(value, 1))}>
        <ChevronRightIcon className="size-4" />
      </Button>
    </div>
  )
}
