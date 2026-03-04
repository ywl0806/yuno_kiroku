import { Button } from '@/components/ui/button'
import { useGetAlbums } from '@/feature/upload/hooks/use-get-albums'
import { useGetFamilies } from '@/feature/settings/hooks/use-get-families'
import { ChevronRight, Pencil, Plus } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

export const SettingsPage = () => {
  const { t } = useTranslation()
  const { data: families } = useGetFamilies()
  const { data: albums } = useGetAlbums()

  return (
    <div className="h-full overflow-y-auto">
      {/* 메뉴 */}
      <section>
        <Link to="/settings/app">
          <div className="flex items-center justify-between px-4 py-3 hover:bg-accent">
            <span className="text-sm">{t('settings.menu.app')}</span>
            <ChevronRight className="size-4 text-muted-foreground" />
          </div>
        </Link>
        <Link to="/settings/account">
          <div className="flex items-center justify-between px-4 py-3 hover:bg-accent">
            <span className="text-sm">{t('settings.menu.account')}</span>
            <ChevronRight className="size-4 text-muted-foreground" />
          </div>
        </Link>
      </section>

      <div className="mx-4 border-t" />

      {/* 가족 목록 */}
      <section>
        <div className="flex items-center justify-between px-4 py-2">
          <span className="text-xs font-semibold text-muted-foreground">{t('settings.family.title')}</span>
          <Button asChild variant="ghost" size="icon">
            <Link to="/settings/family/new">
              <Plus className="size-4" />
            </Link>
          </Button>
        </div>
        {families?.map((family) => (
          <Link key={family.id} to={`/settings/family/${family.id}/edit`}>
            <div className="flex items-center justify-between px-4 py-3 hover:bg-accent">
              <span className="text-sm">{family.name}</span>
              <Pencil className="size-4 text-muted-foreground" />
            </div>
          </Link>
        ))}
      </section>

      <div className="mx-4 border-t" />

      {/* 앨범 목록 */}
      <section>
        <div className="flex items-center justify-between px-4 py-2">
          <span className="text-xs font-semibold text-muted-foreground">{t('settings.album.title')}</span>
          <Button asChild variant="ghost" size="icon">
            <Link to="/settings/album/new">
              <Plus className="size-4" />
            </Link>
          </Button>
        </div>
        {albums?.map((album) => (
          <Link key={album.id} to={`/settings/album/${album.id}/edit`}>
            <div className="flex items-center justify-between px-4 py-3 hover:bg-accent">
              <span className="text-sm">{album.name}</span>
              <Pencil className="size-4 text-muted-foreground" />
            </div>
          </Link>
        ))}
      </section>
    </div>
  )
}
