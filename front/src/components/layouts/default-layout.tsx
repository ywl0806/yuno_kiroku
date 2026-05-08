import { cn } from '@/lib/utils'
import { useScrollHide } from '@/hooks/use-scroll-hide'
import { UploadPhotoProvider, useUploadPhoto } from '@/providers/upload-photo-provider'
import { Clock, House, Plus, Search, Settings } from 'lucide-react'
import { useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Link, Outlet, useLocation } from 'react-router-dom'

const navConfig = [
  { path: '/recent', icon: Clock, labelKey: 'nav.recent' },
  { path: '/search', icon: Search, labelKey: 'nav.search' },
  { path: '/', icon: House, labelKey: 'nav.album' },
  { path: '/settings', icon: Settings, labelKey: 'nav.settings' },
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
    <Link to={path} draggable={false} className="flex flex-1 flex-col items-center justify-center gap-1 py-2">
      <div
        className={cn(
          'flex items-center justify-center rounded-full px-5 py-1 transition-all duration-200',
          isActive ? 'bg-primary/15' : 'bg-transparent',
        )}
      >
        <div className={cn('transition-colors duration-200', isActive ? 'text-primary' : 'text-muted-foreground')}>{icon}</div>
      </div>
      <span className={cn('text-[11px] transition-colors duration-200', isActive ? 'font-semibold text-primary' : 'text-muted-foreground')}>
        {label}
      </span>
    </Link>
  )
}

const DefaultLayoutInner = () => {
  const { t } = useTranslation()
  const location = useLocation()
  const pathname = location.pathname
  const { hidden, handleScroll, reset } = useScrollHide()
  const { openAlbumSheet, albumSheetOpen } = useUploadPhoto()
  const isUploadPage = useMemo(() => pathname.startsWith('/upload'), [pathname])
  useEffect(() => {
    reset()
  }, [pathname, reset])

  useEffect(() => {
    const onScroll = (e: Event) => {
      const target = e.target as HTMLElement
      const currentY = 'scrollTop' in target ? target.scrollTop : window.scrollY
      handleScroll(currentY)
    }
    document.addEventListener('scroll', onScroll, { passive: true, capture: true })
    return () => document.removeEventListener('scroll', onScroll, { capture: true })
  }, [handleScroll])

  return (
    <div className="flex h-dvh flex-col">
      <div className="flex-1 overflow-y-hidden pb-[calc(4rem+env(safe-area-inset-bottom))]">
        <Outlet />
      </div>

      <div
        className="fixed bottom-0 left-0 right-0 z-10 transition-transform duration-300 ease-in-out"
        style={{ transform: hidden ? 'translateY(100%)' : 'translateY(0)' }}
      >
        {/* FAB */}
        {!isUploadPage && (
          <div className="pointer-events-none absolute bottom-full left-1/2 mb-3 -translate-x-1/2">
            <button
              type="button"
              onClick={openAlbumSheet}
              draggable={false}
              className={cn(
                'pointer-events-auto flex h-14 w-14 items-center justify-center rounded-full shadow-lg transition-transform duration-150 active:scale-95',
                albumSheetOpen ? 'bg-primary/80' : 'bg-primary',
              )}
            >
              <Plus className="h-6 w-6 text-primary-foreground" strokeWidth={2.5} />
            </button>
          </div>
        )}
        {/* Nav bar */}
        <div className="flex min-h-[4rem] items-stretch border-t border-border bg-background pb-[env(safe-area-inset-bottom)]">
          {navConfig.map((nav) => (
            <NavItem
              key={nav.path}
              path={nav.path}
              icon={<nav.icon size={22} />}
              label={t(nav.labelKey)}
              isActive={pathname === nav.path}
            />
          ))}

        </div>
      </div>
    </div>
  )
}

export const DefaultLayout = () => {
  return (
    <UploadPhotoProvider>
      <DefaultLayoutInner />
    </UploadPhotoProvider>
  )
}
