import { Member } from '@/types'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'
import { SettingsList, SettingsListAddButton } from './settings-list'

interface Props {
  members?: Member[]
}

export const SettingsMemberList = ({ members }: Props) => {
  const { t } = useTranslation()

  const getTitleLabel = (member: Member) => {
    if (!member.family_title) return null
    if (member.family_title === 'custom') return member.custom_family_title || null
    return t(`settings.member.familyTitle.${member.family_title}`)
  }

  const getInitial = (member: Member) => {
    const name = member.name || member.username || ''
    return name.charAt(0).toUpperCase() || '?'
  }

  const avatarColors = [
    'bg-rose-100 text-rose-600',
    'bg-amber-100 text-amber-600',
    'bg-emerald-100 text-emerald-600',
    'bg-sky-100 text-sky-600',
    'bg-violet-100 text-violet-600',
    'bg-orange-100 text-orange-600',
  ]

  return (
    <SettingsList>
      {members?.map((member, i) => {
        const titleLabel = getTitleLabel(member)
        const colorClass = avatarColors[i % avatarColors.length]
        return (
          <Link
            key={member.id}
            to={`/settings/member/${member.id}/edit`}
            className="group relative flex items-center justify-between px-4 py-3 transition-colors hover:bg-stone-50 active:bg-stone-100"
          >
            <div className="flex items-center gap-3">
              <div className={`flex size-9 flex-shrink-0 items-center justify-center rounded-full text-sm font-semibold ${colorClass}`}>
                {getInitial(member)}
              </div>
              <div className="flex flex-col gap-0.5">
                <span className="text-sm font-medium text-stone-800">{member.name || member.username}</span>
                {titleLabel && (
                  <span className="text-xs text-stone-400">{titleLabel}</span>
                )}
              </div>
            </div>
            <ChevronRight className="size-4 flex-shrink-0 text-stone-300 transition-transform group-hover:translate-x-0.5" />
            <span className="absolute inset-x-4 bottom-0 h-px bg-stone-100 group-last:hidden" />
          </Link>
        )
      })}
      <SettingsListAddButton to="/settings/member/invite" label={t('settings.member.invite')} />
    </SettingsList>
  )
}
