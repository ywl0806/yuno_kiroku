import { PhotoDetailSwipeDialog } from '@/components/blocks/photo-detail-swipe-dialog'
import { PhotoGrid } from '@/components/blocks/photo-grid'
import { useGetMediaItems } from '@/feature/home/hooks/use-get-media-items'
import { FC, useEffect, useMemo, useState } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'

type Props = {
  year: number
  month: number
  date: string
}

export const PhotoGridContainer: FC<Props> = ({ year, month, date }) => {
  const { data: photos, refetch, isFetched } = useGetMediaItems({ year, month, enabled: false })
  const [detailViewIndex, setDetailViewIndex] = useState<number>(0)
  const [openDetailView, setOpenDetailView] = useState<boolean>(false)

  useEffect(() => {
    if (!date || isFetched) return
    const [y, m] = date.split('-')
    if (parseInt(y) === year && parseInt(m) === month) {
      refetch()
    }
  }, [date, year, month, refetch])

  const photoAlbum: AlbumPhoto[] = useMemo(() => {
    if (!photos) return []

    return photos.map((mediaItem) => {
      return {
        key: mediaItem.id,
        src: mediaItem.thumbnail_url,
        width: mediaItem.thumbnail_width,
        height: mediaItem.thumbnail_height,
        alt: mediaItem.file_name,
      }
    })
  }, [photos])
  return (
    <div className="scroll-container h-full w-full overflow-y-auto">
      <PhotoGrid
        photos={photoAlbum}
        onClick={(index) => {
          setDetailViewIndex(index)
          setOpenDetailView(true)
        }}
      />
      <PhotoDetailSwipeDialog
        photos={photos ?? []}
        index={detailViewIndex}
        setIndex={setDetailViewIndex}
        open={openDetailView}
        onClose={() => setOpenDetailView(false)}
      />
    </div>
  )
}
