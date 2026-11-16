import { useForm, FormProvider, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { CATEGORY_LABELS } from '@/contants'
import { Button } from '@/components/ui/button'
import { TextInput } from '@/components/form/text-input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogClose, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import type { TransactionFormDialogProps } from '@/types'

const schema = z.object({
  description: z.string().min(1, 'Description required'),
  amount: z.string().min(1, 'Amount required').transform(Number).pipe(z.number().min(0.01, 'Must be positive')),
  category: z.string().min(1, 'Category required'),
  status: z.enum(['pending', 'paid']),
})

type TransactionFormData = z.output<typeof schema>
type TransactionFormInput = z.input<typeof schema>

export function TransactionFormDialog({ open, onOpenChange }: TransactionFormDialogProps) {
  const methods = useForm<TransactionFormInput, unknown, TransactionFormData>({
    resolver: zodResolver(schema),
    defaultValues: { description: '', amount: '', category: '', status: 'pending' },
  })

  const onSubmit = (data: TransactionFormData) => {
    console.log(data)
    methods.reset()
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New Transaction</DialogTitle>
        </DialogHeader>
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)} className="flex flex-col gap-4">
            <TextInput name="description" title="Description" label="Description" />
            <TextInput name="amount" type="number" title="Amount" label="0.00" />
            <div className="flex flex-col gap-1">
              <label className="text-xs text-muted-foreground">Category</label>
              <Controller
                name="category"
                control={methods.control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                      {Object.entries(CATEGORY_LABELS).map(([value, label]) => (
                        <SelectItem key={value} value={value}>{label}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              <p className="text-xs text-destructive">
                {methods.formState.errors.category?.message}
              </p>
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs text-muted-foreground">Status</label>
              <Controller
                name="status"
                control={methods.control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="pending">Pending</SelectItem>
                      <SelectItem value="paid">Paid</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
            <div className="flex flex-col gap-2 mt-2">
              <Button type="submit" className="w-full">Create</Button>
              <DialogClose render={<Button variant="outline" className="w-full" type="button" />}>
                Cancel
              </DialogClose>
            </div>
          </form>
        </FormProvider>
      </DialogContent>
    </Dialog>
  )
}
