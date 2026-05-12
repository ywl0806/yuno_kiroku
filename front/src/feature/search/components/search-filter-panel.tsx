import { cn } from '@/lib/utils'
import { useGetAlbumOptions } from '@/feature/home/hooks/use-get-album-options'
import { useGetIdentityOptions } from '@/feature/home/hooks/use-get-identity-options'
import { useMediaItemsRange } from '@/feature/home/hooks/use-media-items-range'
import { useGetTags } from '@/feature/tag/hooks/use-get-tags'
import { SearchFilter } from '@/feature/search/hooks/use-search-media-items'
import { MediaItemRange } from '@/types'
import { Heart } from 'lucide-react'
import { FC, useMemo } from 'react'
import { useTranslation } from 'react-i18next'

type Props = {
  filter: SearchFilter
  onFilterChange: (filter: SearchFilter) => void
}

export const SearchFilterPanel: FC<Props> = ({ filter, onFilterChange }) => {
  const { t } = useTranslation()
  const { data: albumOptions } = useGetAlbumOptions()
  const { data: identityOptions } = useGetIdentityOptions()
  const { range } = useMediaItemsRange()
  const { data: tags = [] } = useGetTags()

  const kidsOptions = useMemo(() => {
    return (identityOptions ?? [])
      .filter((opt) => opt.kid_id !== null)
      .map((opt) => ({ id: opt.id, name: opt.kid_name ?? '', image_url: opt.image_url ?? '' }))
  }, [identityOptions])

  const rangeLabel = (r: MediaItemRange) => `${r.year}-${r.month}`

  return (
    <div className="space-y-4 p-4">
      {/* 좋아요 토글 */}
      <div>
        <button
          type="button"
          onClick={() => onFilterChange({ ...filter, liked: !filter.liked })}
          className={cn(
            'flex items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-medium transition-colors',
            filter.liked
              ? 'border-red-400 bg-red-50 text-red-500 dark:bg-red-950'
              : 'border-border bg-background text-muted-foreground',
          )}
        >
          <Heart className={cn('size-3.5', filter.liked ? 'fill-red-500 text-red-500' : '')} />
          {t('tag.liked_filter')}
        </button>
      </div>

      {/* 달 선택 */}
      <div>
        <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('search.period')}</p>
        <div className="flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => onFilterChange({ ...filter, selectedRange: null })}
            className={cn(
              'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
              filter.selectedRange === null
                ? 'border-primary bg-primary text-primary-foreground'
                : 'border-border bg-background text-muted-foreground',
            )}
          >
            {t('search.allPeriod')}
          </button>
          {range.map((r) => (
            <button
              key={`${r.year}-${r.month}`}
              type="button"
              onClick={() => onFilterChange({ ...filter, selectedRange: r })}
              className={cn(
                'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
                filter.selectedRange?.year === r.year && filter.selectedRange?.month === r.month
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-background text-muted-foreground',
              )}
            >
              {rangeLabel(r)}
            </button>
          ))}
        </div>
      </div>

      {/* 앨범 선택 */}
      <div>
        <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('home.filterAlbum')}</p>
        <div className="flex gap-2 overflow-x-auto pb-1">
          <button
            type="button"
            onClick={() => onFilterChange({ ...filter, selectedAlbumId: null })}
            className={cn(
              'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
              filter.selectedAlbumId === null
                ? 'border-primary bg-primary text-primary-foreground'
                : 'border-border bg-background text-muted-foreground',
            )}
          >
            {t('home.filterAll')}
          </button>
          {albumOptions?.map((album) => (
            <button
              key={album.id}
              type="button"
              onClick={() => onFilterChange({ ...filter, selectedAlbumId: album.id })}
              className={cn(
                'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
                filter.selectedAlbumId === album.id
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-background text-muted-foreground',
              )}
            >
              {album.is_common ? t('settings.album.commonName') : album.name}
            </button>
          ))}
        </div>
      </div>

      {/* 아이 선택 */}
      {kidsOptions.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('home.filterKid')}</p>
          <div className="flex gap-3">
            {kidsOptions.map((kid) => {
              const isSelected = filter.selectedIdentityIds.includes(kid.id)
              return (
                <button
                  key={kid.id}
                  type="button"
                  onClick={() => {
                    const next = isSelected
                      ? filter.selectedIdentityIds.filter((id) => id !== kid.id)
                      : [...filter.selectedIdentityIds, kid.id]
                    onFilterChange({ ...filter, selectedIdentityIds: next })
                  }}
                  className="flex flex-col items-center gap-1"
                >
                  <div
                    className={cn(
                      'h-12 w-12 overflow-hidden rounded-full border-2 transition-colors',
                      isSelected ? 'border-primary' : 'border-transparent',
                    )}
                  >
                    {kid.image_url ? (
                      <img src={kid.image_url} alt={kid.name} className="h-full w-full object-cover" />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center bg-accent text-xs text-muted-foreground">
                        {kid.name.charAt(0)}
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

      {/* 태그 선택 */}
      {tags.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-semibold text-muted-foreground">{t('tag.filter')}</p>
          <div className="flex flex-wrap gap-2">
            {tags.map((tag) => {
              const isSelected = filter.selectedTagIds.includes(tag.id)
              return (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => {
                    const next = isSelected
                      ? filter.selectedTagIds.filter((id) => id !== tag.id)
                      : [...filter.selectedTagIds, tag.id]
                    onFilterChange({ ...filter, selectedTagIds: next })
                  }}
                  className={cn(
                    'shrink-0 rounded-full border px-3 py-1 text-xs transition-colors',
                    isSelected
                      ? 'border-primary bg-primary text-primary-foreground'
                      : 'border-border bg-background text-muted-foreground',
                  )}
                >
                  {tag.name}
                </button>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}
