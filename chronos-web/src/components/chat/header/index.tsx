import { Link } from 'react-router'
import { ArrowLeftIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function ChatHeader() {
  return (
    <div className="flex items-center px-4 py-2 border-b">
      <Button variant="ghost" size="sm" render={<Link to="/dashboard" />}>
        <ArrowLeftIcon />
        Dashboard
      </Button>
    </div>
  )
}
