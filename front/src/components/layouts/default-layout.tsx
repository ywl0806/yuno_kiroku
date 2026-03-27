import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { House, Search, Settings, Upload } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link, Outlet, useLocation } from 'react-router-dom'

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
    <Button asChild variant="ghost" className={cn('h-full', isActive && 'bg-accent text-accent-foreground')} draggable={false}>
      <Link to={path}>
        <div className="flex flex-col items-center">
          {icon}
          {label}
        </div>
      </Link>
    </Button>
  )
}
const navConfig: { path: string; icon: React.ElementType; labelKey: string }[] = [
  { path: '/', icon: House, labelKey: 'nav.album' },
  { path: '/search', icon: Search, labelKey: 'nav.search' },
  { path: '/upload', icon: Upload, labelKey: 'nav.upload' },
  { path: '/settings', icon: Settings, labelKey: 'nav.settings' },
]

export const DefaultLayout = () => {
  const { t } = useTranslation()
  const location = useLocation()
  const pathname = location.pathname
  return (
    <div className="flex h-screen flex-col">
      <div className="flex-1 overflow-y-hidden">
        <Outlet />
      </div>

      <div className="flex h-[3.5rem] items-center justify-center gap-5 px-5">
        {navConfig.map((nav) => (
          <NavItem
            key={nav.path}
            path={nav.path}
            icon={<nav.icon />}
            label={t(nav.labelKey)}
            isActive={pathname === nav.path}
          />
        ))}
      </div>
    </div>
  )
}
