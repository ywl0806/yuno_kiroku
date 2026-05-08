import { Group } from '@/types'
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
          {group.name}
        </SettingsListItem>
      ))}
      <SettingsListAddButton to="/settings/family/new" label={t('settings.group.add')} />
    </SettingsList>
  )
}
