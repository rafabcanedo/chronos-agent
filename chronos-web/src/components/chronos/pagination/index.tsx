import { cn } from 'cn'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import type { AppPaginationProps } from '@/types'

export function AppPagination({ meta, onPageChange }: AppPaginationProps) {
  const { current_page, total_pages } = meta

  const isFirst = current_page <= 1
  const isLast  = current_page >= total_pages

  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            onClick={() => !isFirst && onPageChange(current_page - 1)}
            className={cn(isFirst && 'pointer-events-none opacity-50')}
          />
        </PaginationItem>

        <PaginationItem>
          <PaginationLink isActive>{current_page}</PaginationLink>
        </PaginationItem>

        {total_pages > 1 && current_page < total_pages && (
          <>
            <PaginationItem><PaginationEllipsis /></PaginationItem>
            <PaginationItem>
              <PaginationLink onClick={() => onPageChange(total_pages)}>
                {total_pages}
              </PaginationLink>
            </PaginationItem>
          </>
        )}

        <PaginationItem>
          <PaginationNext
            onClick={() => !isLast && onPageChange(current_page + 1)}
            className={cn(isLast && 'pointer-events-none opacity-50')}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  )
}
