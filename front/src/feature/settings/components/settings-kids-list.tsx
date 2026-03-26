import { Kid } from '@/types'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsList, SettingsListAddButton, SettingsListItem } from './settings-list'

interface Props {
  kids?: Kid[]
}

export const SettingsKidsList = ({ kids }: Props) => {
  const { t } = useTranslation()

  return (
    <SettingsList>
      {kids?.map((kid) => (
        <SettingsListItem key={kid.id} to={`/settings/kid/${kid.id}/edit`}>
          <span className="text-sm">{kid.name ?? kid.birth_date ?? `Kid #${kid.id}`}</span>
          <ChevronRight className="size-4 text-muted-foreground" />
        </SettingsListItem>
      ))}
      <SettingsListAddButton to="/settings/kid/new" label={t('settings.kid.add')} />
    </SettingsList>
  )
}
