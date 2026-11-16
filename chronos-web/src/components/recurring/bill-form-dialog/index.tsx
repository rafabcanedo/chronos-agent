import { useForm, FormProvider, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { CATEGORY_LABELS } from '@/contants'
import { Button } from '@/components/ui/button'
import { TextInput } from '@/components/form/text-input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogClose, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import type { BillFormDialogProps } from '@/types'

const schema = z.object({
  name: z.string().min(1, 'Name required').max(100),
  default_amount: z.string().min(1, 'Amount required').transform(Number).pipe(z.number().min(0, 'Must be positive')),
  billing_day: z.string().min(1, 'Billing day required').transform(Number).pipe(z.number().int().min(1, 'Min 1').max(31, 'Max 31')),
  category: z.string().min(1, 'Category required'),
})

type BillFormData = z.output<typeof schema>
type BillFormInput = z.input<typeof schema>

export function BillFormDialog({ open, onOpenChange }: BillFormDialogProps) {
  const methods = useForm<BillFormInput, unknown, BillFormData>({
    resolver: zodResolver(schema),
    defaultValues: { name: '', default_amount: '', billing_day: '', category: '' },
  })

  const onSubmit = (data: BillFormData) => {
    console.log(data)
    methods.reset()
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New Bill</DialogTitle>
        </DialogHeader>
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)} className="flex flex-col gap-4">
            <TextInput name="name" title="Name" label="Bill name" />
            <TextInput name="default_amount" type="number" title="Default Amount" label="0.00" />
            <TextInput name="billing_day" type="number" title="Billing Day" label="1–31" />
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
