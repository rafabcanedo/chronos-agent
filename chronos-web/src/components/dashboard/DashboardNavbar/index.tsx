import { Link } from 'react-router'
import { ChevronDownIcon, PlusIcon, UserIcon } from 'lucide-react'
import logo from '@/assets/logo-chronos.png'
import { Button } from '@/components/ui/button'
import {
  NavigationMenu,
  NavigationMenuList,
  NavigationMenuItem,
  NavigationMenuTrigger,
  NavigationMenuContent,
  NavigationMenuLink,
} from '@/components/ui/navigation-menu'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'

export function DashboardNavbar() {
  return (
    <header className="flex h-14 items-center justify-between border-b border-border px-6">
      <img src={logo} alt="Chronos" className="h-7 select-none" draggable={false} />

      <div className="flex items-center gap-2">
        <NavigationMenu>
          <NavigationMenuList>
            <NavigationMenuItem>
              <NavigationMenuTrigger>Pages</NavigationMenuTrigger>
              <NavigationMenuContent>
                <ul className="w-40">
                  <li>
                    <NavigationMenuLink render={<Link to="/chat" />}>Chat</NavigationMenuLink>
                  </li>
                  <li>
                    <NavigationMenuLink render={<Link to="/" />}>Home</NavigationMenuLink>
                  </li>
                  <li>
                    <NavigationMenuLink render={<Link to="/about" />}>About</NavigationMenuLink>
                  </li>
                </ul>
              </NavigationMenuContent>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>

        <DropdownMenu>
          <DropdownMenuTrigger render={<Button variant="outline" size="sm" />}>
            Actions <ChevronDownIcon />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem><PlusIcon /> Create transaction</DropdownMenuItem>
            <DropdownMenuItem><PlusIcon /> Create bill</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem><UserIcon /> View profile</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  )
}
