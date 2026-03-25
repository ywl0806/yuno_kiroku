import { useGetIdentityFaceOptions } from '@/feature/settings/hooks/use-get-identity-face-options'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'

export type HomeFilter = {
  selectedAlbumId: number | null
  selectedIdentityIds: number[]
}

type Props = {
  filter: HomeFilter
  onFilterChange: (filter: HomeFilter) => void
}

export const HomeFilterPanel = ({ filter, onFilterChange }: Props) => {
  const { t } = useTranslation()
  const { data: settingsData } = useGetSettingsData()

  const albums = settingsData?.albums ?? []
  const kidsWithIdentity = (settingsData?.kids ?? []).filter((k) => k.identity_id !== null)

  const { data: faceOptions } = useGetIdentityFaceOptions(
    true,
    kidsWithIdentity.map((k) => k.id),
  )

  const handleAlbumSelect = (albumId: number | null) => {
    onFilterChange({ ...filter, selectedAlbumId: albumId })
  }

  const handleIdentityToggle = (identityId: number) => {
    const isSelected = filter.selectedIdentityIds.includes(identityId)
    const next = isSelected
      ? filter.selectedIdentityIds.filter((id) => id !== identityId)
      : [...filter.selectedIdentityIds, identityId]
    onFilterChange({ ...filter, selectedIdentityIds: next })
  }

  return (
    <div className="border-b bg-background px-4 py-3">
      {/* 앨범 섹션 */}
      <div className="mb-3">
        <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('home.filterAlbum')}</p>
        <div className="flex gap-2 overflow-x-auto pb-1">
          <button
            type="button"
            onClick={() => handleAlbumSelect(null)}
            className={cn(
              'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
              filter.selectedAlbumId === null
                ? 'border-primary bg-primary text-primary-foreground'
                : 'border-border bg-background text-muted-foreground',
            )}
          >
            {t('home.filterAll')}
          </button>
          {albums.map((album) => (
            <button
              key={album.id}
              type="button"
              onClick={() => handleAlbumSelect(album.id)}
              className={cn(
                'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
                filter.selectedAlbumId === album.id
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-background text-muted-foreground',
              )}
            >
              {album.name}
            </button>
          ))}
        </div>
      </div>

      {/* 아이 섹션 */}
      {kidsWithIdentity.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('home.filterKid')}</p>
          <div className="flex gap-3">
            {kidsWithIdentity.map((kid) => {
              const isSelected = filter.selectedIdentityIds.includes(kid.identity_id!)
              const faceOption = faceOptions?.find((f) => f.identity_id === kid.identity_id)
              return (
                <button
                  key={kid.id}
                  type="button"
                  onClick={() => handleIdentityToggle(kid.identity_id!)}
                  className="flex flex-col items-center gap-1"
                >
                  <div
                    className={cn(
                      'h-12 w-12 overflow-hidden rounded-full border-2 transition-colors',
                      isSelected ? 'border-primary' : 'border-transparent',
                    )}
                  >
                    {faceOption ? (
                      <img src={faceOption.image_url} alt={kid.name ?? ''} className="h-full w-full object-cover" />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center bg-accent text-xs text-muted-foreground">
                        {kid.name?.charAt(0) ?? '?'}
                      </div>
                    )}
                  </div>
                  <span className={cn('text-xs', isSelected ? 'font-semibold text-primary' : 'text-muted-foreground')}>
                    {kid.name}
                  </span>
                </button>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}
