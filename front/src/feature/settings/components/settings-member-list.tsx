import { Member } from '@/types'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsList, SettingsListAddButton, SettingsListItem } from './settings-list'

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

  return (
    <SettingsList>
      {members?.map((member) => {
        const titleLabel = getTitleLabel(member)
        return (
          <SettingsListItem key={member.id} to={`/settings/member/${member.id}/edit`}>
            <div className="flex flex-col gap-0.5">
              <span className="text-sm">{member.name || member.username}</span>
              {titleLabel && <span className="text-xs text-muted-foreground">{titleLabel}</span>}
            </div>
            <ChevronRight className="size-4 text-muted-foreground" />
          </SettingsListItem>
        )
      })}
      <SettingsListAddButton to="/settings/member/invite" label={t('settings.member.invite')} />
    </SettingsList>
  )
}
