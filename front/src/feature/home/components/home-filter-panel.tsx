import { cn } from '@/lib/utils'
import { useTranslation } from 'react-i18next'
import { useGetIdentityOptions } from '@/feature/home/hooks/use-get-identity-options'
import { useMemo } from 'react'
import { useGetAlbumOptions } from '@/feature/home/hooks/use-get-album-options'

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

  const { data: identitiesOptions } = useGetIdentityOptions()
  const { data: albumOptions } = useGetAlbumOptions()

  const kidsOptions: {
    id: number
    name: string
    image_url: string
  }[] = useMemo(() => {
    const options: {
      id: number
      name: string
      image_url: string
    }[] = []
    for (const identityOption of identitiesOptions ?? []) {
      if (identityOption.kid_id !== null) {
        options.push({ id: identityOption.id, name: identityOption.kid_name ?? '', image_url: identityOption.image_url ?? '' })
      }
    }
    return options
  }, [identitiesOptions])

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
          {albumOptions?.map((albumOption) => (
            <button
              key={albumOption.id}
              type="button"
              onClick={() => handleAlbumSelect(albumOption.id)}
              className={cn(
                'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
                filter.selectedAlbumId === albumOption.id
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-background text-muted-foreground',
              )}
            >
              {albumOption.is_common ? t('settings.album.commonName') : albumOption.name}
            </button>
          ))}
        </div>
      </div>

      {/* 아이 섹션 */}
      {identitiesOptions && identitiesOptions.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('home.filterKid')}</p>
          <div className="flex gap-3">
            {kidsOptions.map((kidOption) => {
              const isSelected = filter.selectedIdentityIds.includes(kidOption.id)
              return (
                <button
                  key={kidOption.id}
                  type="button"
                  onClick={() => handleIdentityToggle(kidOption.id)}
                  className="flex flex-col items-center gap-1"
                >
                  <div
                    className={cn(
                      'h-12 w-12 overflow-hidden rounded-full border-2 transition-colors',
                      isSelected ? 'border-primary' : 'border-transparent',
                    )}
                  >
                    {kidOption.image_url ? (
                      <img src={kidOption.image_url} alt={kidOption.name} className="h-full w-full object-cover" />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center bg-accent text-xs text-muted-foreground">
                        {kidOption.name.charAt(0)}
                      </div>
                    )}
                  </div>
                  <span className={cn('text-xs', isSelected ? 'font-semibold text-primary' : 'text-muted-foreground')}>
                    {kidOption.name}
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
