import { useState } from 'react'
import { MoonIcon, SunIcon } from 'lucide-react'
import { GITHUB_USERNAME } from '@/contants'
import { useTheme } from '@/hooks/useTheme'
import { Navbar } from '@/components/dashboard/navbar'
import { DeleteAccountDialog } from '@/components/profile/delete-account-dialog'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

function getInitials(username: string): string {
  return username.slice(0, 2).toUpperCase()
}

export default function Profile() {
  const { theme, toggleTheme } = useTheme()
  const [deleteOpen, setDeleteOpen] = useState(false)

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

            <div className="flex flex-col gap-4 max-w-sm">
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
              <Button variant="destructive" onClick={() => setDeleteOpen(true)}>Delete</Button>
            </div>
          </div>

          <div className="flex flex-col gap-4">
            <h2 className="text-sm font-semibold">Tools</h2>

            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="icon"
                onClick={toggleTheme}
              >
                {theme === 'dark' ? <SunIcon className="size-4" /> : <MoonIcon className="size-4" />}
              </Button>
              <span className="text-sm text-muted-foreground">
                {theme === 'dark' ? 'Light mode' : 'Dark mode'}
              </span>
            </div>
          </div>

        </div>
      </main>
    <DeleteAccountDialog open={deleteOpen} onOpenChange={setDeleteOpen} />
    </div>
  )
}
