import { MessageCircleIcon } from 'lucide-react'
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty'

export function ChatEmptyState() {
  return (
    <Empty className="h-full">
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <MessageCircleIcon />
        </EmptyMedia>
        <EmptyTitle>Hey! I&apos;m Chronos Agent</EmptyTitle>
        <EmptyDescription>
          Tell me what you paid, spent, or want to register.
        </EmptyDescription>
      </EmptyHeader>
    </Empty>
  )
}
