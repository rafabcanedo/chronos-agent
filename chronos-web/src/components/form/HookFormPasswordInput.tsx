import { FC, useState } from 'react'
import { useFormContext } from 'react-hook-form'
import { Eye, EyeOff } from 'lucide-react'
import { Input } from '@/components/ui/input'
import type { ITextInput } from '@/types'

export const HookFormPasswordInput: FC<Omit<ITextInput, 'type'>> = ({
  name,
  label,
  title,
  disabled,
}) => {
  const [visible, setVisible] = useState(false)
  const { register, formState: { errors } } = useFormContext()

  return (
    <div className="flex flex-col gap-1">
      <label className="text-xs text-muted-foreground">{title}</label>
      <div className="relative">
        <Input
          placeholder={label}
          type={visible ? 'text' : 'password'}
          {...register(name)}
          disabled={disabled}
          className="pr-9 disabled:cursor-not-allowed"
        />
        <button
          type="button"
          tabIndex={-1}
          onClick={() => setVisible(v => !v)}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
        >
          {visible ? <EyeOff className="size-4" /> : <Eye className="size-4" />}
        </button>
      </div>
      <p className="text-xs text-destructive">
        {errors[name]?.message as string}
      </p>
    </div>
  )
}
