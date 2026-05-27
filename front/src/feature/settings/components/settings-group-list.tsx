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
          <span>{group.is_admin ? t('settings.group.adminName') : group.name}</span>
          {group.is_admin && (
            <span className="ml-2 rounded px-1.5 py-0.5 text-xs bg-amber-50 text-amber-600">
              {t('settings.group.adminBadge')}
            </span>
          )}
        </SettingsListItem>
      ))}
      <SettingsListAddButton to="/settings/family/new" label={t('settings.group.add')} />
    </SettingsList>
  )
}
