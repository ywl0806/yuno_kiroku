import { Album } from '@/types'
import { useTranslation } from 'react-i18next'
import { SettingsList, SettingsListAddButton, SettingsListItem } from './settings-list'

interface Props {
  albums?: Album[]
}

export const SettingsAlbumList = ({ albums }: Props) => {
  const { t } = useTranslation()

  return (
    <SettingsList>
      {albums?.map((album) => {
        if (album.is_common) {
          return (
            <SettingsListItem key={album.id} to={`/settings/album/${album.id}/edit`}>
              <span>{album.is_common ? t('settings.album.commonName') : album.name}</span>
              {album.is_common && (
                <span className="ml-2 rounded px-1.5 py-0.5 text-xs bg-amber-50 text-amber-600">
                  {t('settings.album.commonBadge')}
                </span>
              )}
            </SettingsListItem>
          )
        }
        return (
          <SettingsListItem key={album.id} to={`/settings/album/${album.id}/edit`}>
            <span>{album.name}</span>
          </SettingsListItem>
        )
      })}
      <SettingsListAddButton to="/settings/album/new" label={t('settings.album.add')} />
    </SettingsList>
  )
}
