import { SignInForm } from '@/components/sign-form'
import logo from '@/assets/logo-chronos.png'

export default function SignIn() {
  return (
    <div className="min-h-screen flex">
      <div className="flex w-1/2 items-center justify-center border-r border-border">
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
    </div>
  )
}
