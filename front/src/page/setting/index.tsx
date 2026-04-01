import { SettingsAlbumList } from '@/feature/settings/components/settings-album-list'
import { SettingsGroupList } from '@/feature/settings/components/settings-group-list'
import { SettingsKidsList } from '@/feature/settings/components/settings-kids-list'
import { SettingsMemberList } from '@/feature/settings/components/settings-member-list'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

const MenuSubTitle = ({ title }: { title: string }) => {
  return (
    <div className="flex items-center justify-between px-4 pb-3 pt-7">
      <span className="text-[0.8rem]">{title}</span>
    </div>
  )
}


export const SettingsPage = () => {
  const { t } = useTranslation()
  const { data } = useGetSettingsData()

  return (
    <div className="h-full bg-neutral-100">
      <div className="max-w-[50rem] mx-auto h-full overflow-y-auto bg-neutral-200 pb-[5rem]">
        <div className="flex items-center justify-center gap-2 px-4 py-2 border-b bg-white">
          <span className="text-lg font-bold">{t('settings.title')}</span>
        </div>
        {/* 메뉴 */}
        <section>
          <MenuSubTitle title={t('settings.app.title')} />
          <Link to="/settings/app">
            <div className="flex items-center justify-between px-4 py-3 hover:bg-accent bg-white">
              <span className="text-sm">{t('settings.app.language.title')}</span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </div>
          </Link>
        </section>

        <section>
          <MenuSubTitle title={t('settings.account.title')} />

          <Link to="/settings/account">
            <div className="flex items-center justify-between px-4 py-3 hover:bg-accent bg-white">
              <span className="text-sm">{t('settings.account.account')}</span>
              <ChevronRight className="size-4 text-muted-foreground" />
            </div>
          </Link>
        </section>

        {/* 앨범 그룹 목록 */}
        <section>
          <MenuSubTitle title={t('settings.group.title')} />
          <SettingsGroupList groups={data?.groups} />
        </section>

        {/* 가족 목록 */}
        <section>
          <MenuSubTitle title={t('settings.member.title')} />
          <SettingsMemberList members={data?.members} />
        </section>

        {/* 앨범 목록 */}
        <section>
          <MenuSubTitle title={t('settings.album.title')} />
          <SettingsAlbumList albums={data?.albums} />
        </section>

        {/* 아이 목록 */}
        <section>
          <MenuSubTitle title={t('settings.kid.title')} />
          <SettingsKidsList kids={data?.kids} />
        </section>
      </div>
    </div>
  )
}
