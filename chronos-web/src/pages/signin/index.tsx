import { SignInForm } from '@/components/sign-form'
import { Card, CardContent } from '@/components/ui/card'
import logo from '@/assets/logo-chronos.png'

export default function SignIn() {
  return (
    <div className="min-h-screen flex items-center justify-center p-8">
      <Card className="shadow-md w-full max-w-3xl h-112.5">
        <CardContent className="flex p-0 h-full">
          <div className="flex w-1/2 items-center justify-center p-8 border-r border-border">
            <img
              src={logo}
              alt="Chronos"
              className="w-52 select-none"
              draggable={false}
            />
          </div>
          <div className="flex w-1/2 items-center justify-center p-8">
            <SignInForm />
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
