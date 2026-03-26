import { Group } from '@/types'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsList, SettingsListAddButton, SettingsListItem } from './settings-list'

interface Props {
  groups?: Group[]
}

export const SettingsGroupList = ({ groups }: Props) => {
  const { t } = useTranslation()

  return (
    <SettingsList>
      {groups?.map((group) => (
        <SettingsListItem key={group.id} to={`/settings/family/${group.id}/edit`}>
          <span className="text-sm">{group.name}</span>
          <ChevronRight className="size-4 text-muted-foreground" />
        </SettingsListItem>
      ))}
      <SettingsListAddButton to="/settings/family/new" label={t('settings.group.add')} />
    </SettingsList>
  )
}
