import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { House, Settings, Upload } from 'lucide-react'
import { Link, Outlet, useLocation } from 'react-router-dom'

const navs: { path: string; icon: React.ElementType; label: string }[] = [
  {
    path: '/',
    icon: House,
    label: 'Album',
  },
  {
    path: '/upload',
    icon: Upload,
    label: 'Upload',
  },
  {
    path: '/settings',
    icon: Settings,
    label: 'Settings',
  },
]

const NavItem = ({
  path,
  icon,
  label,
  isActive,
}: {
  path: string
  icon: React.ReactNode
  label: string
  isActive: boolean
}) => {
  return (
    <Button asChild variant="ghost" className={cn('h-full', isActive && 'bg-accent text-accent-foreground')}>
      <Link to={path}>
        <div className="flex flex-col items-center">
          {icon}
          {label}
        </div>
      </Link>
    </Button>
  )
}
export const DefaultLayout = () => {
  const location = useLocation()
  const pathname = location.pathname
  return (
    <div className="flex h-screen flex-col">
      <div className="flex-1 overflow-y-hidden">
        <Outlet />
      </div>

      <div className="flex h-[3.5rem] items-center justify-center gap-5 px-5">
        {navs.map((nav) => (
          <NavItem
            key={nav.path}
            path={nav.path}
            icon={<nav.icon />}
            label={nav.label}
            isActive={pathname === nav.path}
          />
        ))}
      </div>
    </div>
  )
}
