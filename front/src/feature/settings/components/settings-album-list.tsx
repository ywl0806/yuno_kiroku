import { Album } from '@/types'
import { ChevronRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SettingsList, SettingsListAddButton, SettingsListItem } from './settings-list'

interface Props {
  albums?: Album[]
}

export const SettingsAlbumList = ({ albums }: Props) => {
  const { t } = useTranslation()

  return (
    <SettingsList>
      {albums?.map((album) => (
        <SettingsListItem key={album.id} to={`/settings/album/${album.id}/edit`}>
          <span className="text-sm">{album.name}</span>
          <ChevronRight className="size-4 text-muted-foreground" />
        </SettingsListItem>
      ))}
      <SettingsListAddButton to="/settings/album/new" label={t('settings.album.add')} />
    </SettingsList>
  )
}
