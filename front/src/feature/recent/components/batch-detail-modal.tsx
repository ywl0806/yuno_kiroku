import { PhotoDetailSwipeDialog } from '@/components/blocks/photo-detail-swipe-dialog'
import { PhotoGrid } from '@/components/blocks/photo-grid'
import { FullScreenModal } from '@/components/ui/full-screen-modal'
import { useGetBatchMediaItems } from '@/feature/recent/hooks/use-get-batch-media-items'
import { MediaItem } from '@/types'
import { ArrowLeft } from 'lucide-react'
import { FC, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Photo as AlbumPhoto } from 'react-photo-album'

type Props = {
  batchId: number
  uploadAt: string
  onClose: () => void
}

export const BatchDetailModal: FC<Props> = ({ batchId, uploadAt, onClose }) => {
  const { t } = useTranslation()
  const [detailViewIndex, setDetailViewIndex] = useState(0)
  const [openDetailView, setOpenDetailView] = useState(false)

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = useGetBatchMediaItems(batchId)

  const photos: MediaItem[] = useMemo(() => data?.pages.flatMap((p) => p.items) ?? [], [data])

  const photoAlbum: AlbumPhoto[] = useMemo(
    () =>
      photos.map((item) => ({
        key: item.id,
        src: item.thumbnail_url,
        width: item.thumbnail_width,
        height: item.thumbnail_height,
        alt: item.file_name,
      })),
    [photos],
  )

  const date = new Date(uploadAt)
  const dateLabel = `${date.getFullYear()}.${String(date.getMonth() + 1).padStart(2, '0')}.${String(date.getDate()).padStart(2, '0')}`

  return (
    <FullScreenModal open>
      <div className="flex h-full flex-col">
        <div className="flex items-center gap-3 border-b px-4 py-3">
          <button type="button" onClick={onClose} className="p-1">
            <ArrowLeft size={20} />
          </button>
          <span className="font-medium">{dateLabel}</span>
        </div>
        <div className="flex-1 overflow-y-auto">
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
                    {isFetchingNextPage ? t('recent.loading') : t('recent.loadMore')}
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
          {!isFetchingNextPage && photos.length === 0 && (
            <p className="mt-10 text-center text-sm text-muted-foreground">{t('recent.noUploads')}</p>
          )}
        </div>
      </div>
    </FullScreenModal>
  )
}
