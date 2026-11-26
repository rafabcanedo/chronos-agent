import { useState } from 'react'
import { MoonIcon, SunIcon } from 'lucide-react'
import { GITHUB_USERNAME } from '@/contants'
import { Navbar } from '@/components/dashboard/navbar'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

function getInitials(username: string): string {
  return username.slice(0, 2).toUpperCase()
}

export default function Profile() {
  const [isDark, setIsDark] = useState(false)

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />

      <main className="flex-1 flex flex-col gap-6 p-6">
        <h1 className="text-xl font-heading font-semibold">Profile</h1>

        <div className="grid grid-cols-2 gap-8">

          <div className="flex flex-col gap-6">
            <Avatar className="size-36">
              <AvatarImage src={`https://github.com/${GITHUB_USERNAME}.png`} alt={GITHUB_USERNAME} />
              <AvatarFallback className="text-2xl">{getInitials(GITHUB_USERNAME)}</AvatarFallback>
            </Avatar>

            <div className="flex flex-col gap-4">
              <h2 className="text-sm font-semibold">Profile details</h2>

              <div className="flex flex-col gap-1.5">
                <Label htmlFor="name">Name</Label>
                <Input id="name" defaultValue={GITHUB_USERNAME} />
              </div>

              <div className="flex flex-col gap-1.5">
                <Label htmlFor="email">Email</Label>
                <Input id="email" type="email" defaultValue="rafa@email.com" />
              </div>

              <div className="flex flex-col gap-1.5">
                <Label htmlFor="password">Password</Label>
                <Input id="password" type="password" placeholder="••••••••" />
              </div>
            </div>

            <div className="flex gap-2">
              <Button>Save</Button>

              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">Delete</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Delete account</AlertDialogTitle>
                    <AlertDialogDescription>
                      This action cannot be undone. Your account and all associated data will be permanently deleted.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                      Delete
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>

          <div className="flex flex-col gap-4">
            <h2 className="text-sm font-semibold">Tools</h2>

            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="icon"
                onClick={() => setIsDark(prev => !prev)}
              >
                {isDark ? <SunIcon className="size-4" /> : <MoonIcon className="size-4" />}
              </Button>
              <span className="text-sm text-muted-foreground">
                {isDark ? 'Light mode' : 'Dark mode'}
              </span>
            </div>
          </div>

        </div>
      </main>
    </div>
  )
}
