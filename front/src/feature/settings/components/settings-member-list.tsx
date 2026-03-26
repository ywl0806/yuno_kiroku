import { Member } from '@/types'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsList, SettingsListAddButton, SettingsListItem } from './settings-list'

interface Props {
  members?: Member[]
}

export const SettingsMemberList = ({ members }: Props) => {
  const { t } = useTranslation()

  return (
    <SettingsList>
      {members?.map((member) => (
        <SettingsListItem key={member.id} to={`/settings/member/${member.id}/edit`}>
          <span className="text-sm">{member.name || member.username}</span>
          <ChevronRight className="size-4 text-muted-foreground" />
        </SettingsListItem>
      ))}
      <SettingsListAddButton to="/settings/member/invite" label={t('settings.member.invite')} />
    </SettingsList>
  )
}
