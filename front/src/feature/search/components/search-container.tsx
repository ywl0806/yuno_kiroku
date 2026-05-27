import { PhotoDetailSwipeDialog } from '@/components/blocks/photo-detail-swipe-dialog'
import { PhotoGrid } from '@/components/blocks/photo-grid'
import { SearchFilterPanel } from '@/feature/search/components/search-filter-panel'
import { SearchFilter, useSearchMediaItems } from '@/feature/search/hooks/use-search-media-items'
import { MediaItem } from '@/types'
import { useMemo, useState } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'
import { useTranslation } from 'react-i18next'

export const SearchContainer = () => {
  const { t } = useTranslation()
  const [filter, setFilter] = useState<SearchFilter>({
    selectedAlbumId: null,
    selectedIdentityIds: [],
    selectedRange: null,
    liked: false,
    selectedTagIds: [],
  })
  const [detailViewIndex, setDetailViewIndex] = useState<number>(0)
  const [openDetailView, setOpenDetailView] = useState<boolean>(false)

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = useSearchMediaItems(filter)

  const photos: MediaItem[] = useMemo(() => data?.pages.flatMap((p) => p.items) ?? [], [data])

  const photoAlbum: AlbumPhoto[] = useMemo(
    () =>
      photos.map((mediaItem) => ({
        key: mediaItem.id,
        src: mediaItem.thumbnail_url,
        width: mediaItem.thumbnail_width,
        height: mediaItem.thumbnail_height,
        alt: mediaItem.file_name,
      })),
    [photos],
  )

  const hasFilter =
    filter.selectedAlbumId !== null ||
    filter.selectedIdentityIds.length > 0 ||
    filter.selectedRange !== null ||
    filter.liked ||
    filter.selectedTagIds.length > 0

  return (
    <div className="flex h-full flex-col">
      <div className="border-b pt-5">
        <p className="px-4 pb-2 text-lg font-bold">{t('search.title')}</p>
        <SearchFilterPanel filter={filter} onFilterChange={setFilter} />
      </div>
      <div className="flex-1 overflow-y-auto pb-[calc(4rem+env(safe-area-inset-bottom))]">
        {hasFilter && photos.length === 0 && !isFetchingNextPage && (
          <p className="mt-10 text-center text-sm text-muted-foreground">{t('search.noResults')}</p>
        )}
        {!hasFilter && (
          <p className="mt-10 text-center text-sm text-muted-foreground">{t('search.selectFilter')}</p>
        )}
        {photos.length > 0 && (
          <>
            <PhotoGrid
              photos={photoAlbum}
              onClick={(index) => {
                setDetailViewIndex(index)
                setOpenDetailView(true)
              }}
            />
            {hasNextPage && (
              <div className="flex justify-center py-4">
                <button
                  type="button"
                  onClick={() => fetchNextPage()}
                  disabled={isFetchingNextPage}
                  className="rounded-full border border-border px-5 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent disabled:opacity-50"
                >
                  {isFetchingNextPage ? t('search.loading') : t('search.loadMore')}
                </button>
              </div>
            )}
            <PhotoDetailSwipeDialog
              photos={photos}
              index={detailViewIndex}
              setIndex={setDetailViewIndex}
              open={openDetailView}
              onClose={() => setOpenDetailView(false)}
            />
          </>
        )}
      </div>
    </div>
  )
}
