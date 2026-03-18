import { SettingsAlbumList } from '@/feature/settings/components/settings-album-list'
import { SettingsGroupList } from '@/feature/settings/components/settings-group-list'
import { SettingsKidsList } from '@/feature/settings/components/settings-kids-list'
import { SettingsMemberList } from '@/feature/settings/components/settings-member-list'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

export const SettingsPage = () => {
  const { t } = useTranslation()
  const { data } = useGetSettingsData()

  return (
    <div className="mx-auto h-full max-w-[50rem] overflow-y-auto pt-10">
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

      {/* 앨범 그룹 목록 */}
      <section>
        <SettingsGroupList groups={data?.groups} />
      </section>

      <div className="mx-4 border-t" />

      {/* 멤버 목록 */}
      <section>
        <SettingsMemberList members={data?.members} />
      </section>

      <div className="mx-4 border-t" />

      {/* 앨범 목록 */}
      <section>
        <SettingsAlbumList albums={data?.albums} />
      </section>

      <div className="mx-4 border-t" />

      {/* 아이 목록 */}
      <section>
        <SettingsKidsList kids={data?.kids} />
      </section>
    </div>
  )
}
