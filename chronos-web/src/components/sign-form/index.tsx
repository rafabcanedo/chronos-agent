import { useForm, FormProvider } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { HookFormTextInput } from '@/components/form/HookFormTextInput'
import { HookFormPasswordInput } from '@/components/form/HookFormPasswordInput'
import { Button } from '@/components/ui/button'

const schema = z.object({
  username: z.string().min(1, 'Username required'),
  password: z.string().min(1, 'Password required'),
})

type SignInFormData = z.infer<typeof schema>

export function SignInForm() {
  const methods = useForm<SignInFormData>({
    resolver: zodResolver(schema),
  })

  const onSubmit = (data: SignInFormData) => {
    console.log(data)
  }

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

        <HookFormTextInput name="username" title="Username" label="Your username" />
        <HookFormPasswordInput name="password" title="Password" label="Your password" />

        <Button type="submit" className="w-full mt-2">
          Login
        </Button>
      </form>
    </FormProvider>
  )
}
