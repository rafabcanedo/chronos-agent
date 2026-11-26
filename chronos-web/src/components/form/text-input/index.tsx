import { useFormContext } from 'react-hook-form'
import { Input } from '@/components/ui/input'
import type { TextInputProps } from '@/types'

export function TextInput({ name, type = 'text', label, title, disabled }: TextInputProps) {
  const { register, formState: { errors } } = useFormContext()

  return (
    <div className="flex flex-col gap-1">
      <label className="text-xs text-muted-foreground">{title}</label>
      <Input
        placeholder={label}
        type={type}
        {...register(name)}
        disabled={disabled}
        className="disabled:cursor-not-allowed"
      />
      <p className="text-xs text-destructive">
        {errors[name]?.message as string}
      </p>
    </div>
  )
}
