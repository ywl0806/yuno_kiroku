import { SettingsAlbumList } from '@/feature/settings/components/settings-album-list'
import { SettingsGroupList } from '@/feature/settings/components/settings-group-list'
import { SettingsKidsList } from '@/feature/settings/components/settings-kids-list'
import { SettingsMemberList } from '@/feature/settings/components/settings-member-list'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { BookImage, ChevronRight, FolderKanban, Globe, Smile, User, Users } from 'lucide-react'
import { ComponentType } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

interface MenuSectionProps {
  title: string
  icon: ComponentType<{ className?: string }>
  iconColor: string
  children: React.ReactNode
}

const MenuSection = ({ title, icon: Icon, iconColor, children }: MenuSectionProps) => (
  <section className="px-4">
    <div className="flex items-center gap-2 px-1 pb-2 pt-6">
      <div className={`flex size-6 items-center justify-center rounded-lg ${iconColor}`}>
        <Icon className="size-3.5" />
      </div>
      <span className="text-[0.7rem] font-semibold uppercase tracking-widest text-stone-400">{title}</span>
    </div>
    {children}
  </section>
)

interface MenuLinkItemProps {
  to: string
  label: string
}

const MenuLinkItem = ({ to, label }: MenuLinkItemProps) => (
  <Link
    to={to}
    className="group relative flex items-center justify-between px-4 py-3.5 transition-colors hover:bg-stone-50 active:bg-stone-100"
  >
    <span className="text-sm font-medium text-stone-800">{label}</span>
    <ChevronRight className="size-4 flex-shrink-0 text-stone-300 transition-transform group-hover:translate-x-0.5" />
  </Link>
)

export const SettingsPage = () => {
  const { t } = useTranslation()
  const { data } = useGetSettingsData()

  return (
    <div className="h-full bg-stone-50">
      <div className="mx-auto h-full max-w-[50rem] overflow-y-auto pb-[5rem]">
        {/* 헤더 */}
        <div className="sticky top-0 z-10 bg-stone-50/80 px-5 pb-3 pt-5 backdrop-blur-sm">
          <h1 className="text-2xl font-bold tracking-tight text-stone-900">{t('settings.title')}</h1>
        </div>

        {/* 앱 설정 */}
        <MenuSection title={t('settings.app.title')} icon={Globe} iconColor="bg-sky-100 text-sky-600">
          <div className="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04]">
            <MenuLinkItem to="/settings/app" label={t('settings.app.language.title')} />
          </div>
        </MenuSection>

        {/* 계정 */}
        <MenuSection title={t('settings.account.title')} icon={User} iconColor="bg-violet-100 text-violet-600">
          <div className="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04]">
            <MenuLinkItem to="/settings/account" label={t('settings.account.account')} />
          </div>
        </MenuSection>

        {/* 앨범 그룹 */}
        <MenuSection title={t('settings.group.title')} icon={FolderKanban} iconColor="bg-emerald-100 text-emerald-600">
          <SettingsGroupList groups={data?.groups} />
        </MenuSection>

        {/* 가족 목록 */}
        <MenuSection title={t('settings.member.title')} icon={Users} iconColor="bg-rose-100 text-rose-500">
          <SettingsMemberList members={data?.members} />
        </MenuSection>

        {/* 앨범 목록 */}
        <MenuSection title={t('settings.album.title')} icon={BookImage} iconColor="bg-indigo-100 text-indigo-600">
          <SettingsAlbumList albums={data?.albums} />
        </MenuSection>

        {/* 아이 목록 */}
        <MenuSection title={t('settings.kid.title')} icon={Smile} iconColor="bg-amber-100 text-amber-600">
          <SettingsKidsList kids={data?.kids} />
        </MenuSection>
      </div>
    </div>
  )
}
