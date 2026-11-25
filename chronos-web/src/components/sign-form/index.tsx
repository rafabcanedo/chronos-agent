import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { useForm, FormProvider } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { login } from '@/services/auth'
import { TOKEN_STORAGE_KEY } from '@/contants'
import { TextInput } from '@/components/form/text-input'
import { PasswordInput } from '@/components/form/password-input'
import { Button } from '@/components/ui/button'

const schema = z.object({
  username: z.string().min(1, 'Username required'),
  password: z.string().min(1, 'Password required'),
})

type SignInFormData = z.infer<typeof schema>

export function SignInForm() {
  const navigate = useNavigate()
  const methods = useForm<SignInFormData>({
    resolver: zodResolver(schema),
  })

  const { mutate, isPending } = useMutation({
    mutationFn: login,
    onSuccess: ({ token }) => {
      localStorage.setItem(TOKEN_STORAGE_KEY, token)
      navigate('/dashboard')
    },
    onError: (error) => {
      methods.setError('root', { message: error.message })
    },
  })

  const onSubmit = (data: SignInFormData) => mutate(data)

  return (
    <FormProvider {...methods}>
      <form
        onSubmit={methods.handleSubmit(onSubmit)}
        className="flex flex-col gap-4 w-full"
      >
        <div className="flex flex-col gap-1 mb-2">
          <h1 className="text-2xl font-heading font-semibold">Signin</h1>
          <p className="text-sm text-muted-foreground">Signin to your Chronos account</p>
        </div>

        <TextInput name="username" title="Username" label="Your username" />
        <PasswordInput name="password" title="Password" label="Your password" />

        {methods.formState.errors.root && (
          <p className="text-xs text-destructive">{methods.formState.errors.root.message}</p>
        )}

        <Button type="submit" className="w-full mt-2" disabled={isPending}>
          {isPending ? 'Signing in...' : 'Login'}
        </Button>
      </form>
    </FormProvider>
  )
}
