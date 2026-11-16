import { PencilIcon, Trash2Icon, XIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from '@/components/ui/drawer'
import type { DetailDrawerProps } from '@/types'

export function DetailDrawer({ open, onOpenChange, title, children }: DetailDrawerProps) {
  return (
    <Drawer open={open} onOpenChange={onOpenChange} swipeDirection="right">
      <DrawerContent>
        <DrawerHeader className="flex-row items-center justify-between">
          <DrawerTitle>{title}</DrawerTitle>
          <DrawerClose render={<Button variant="ghost" size="icon" />}>
            <XIcon className="size-4" />
          </DrawerClose>
        </DrawerHeader>
        {children}
        <DrawerFooter>
          <DrawerClose render={<Button variant="outline" className="w-full" />}>
            Back to list
          </DrawerClose>
          <Button variant="outline" className="w-full">
            <PencilIcon className="size-4" />
            Edit
          </Button>
          <Button variant="outline" className="w-full text-destructive hover:text-destructive border-destructive/40 hover:bg-destructive/5">
            <Trash2Icon className="size-4" />
            Delete
          </Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
