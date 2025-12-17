import { PhotoDetailSwipeDialog } from '@/components/blocks/photo-detail-swipe-dialog'
import { PhotoGrid } from '@/components/blocks/photo-grid'
import { useGetPhotos } from '@/feature/home/hooks/use-get-photos'
import { FC, useMemo, useState } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'

type Props = {
  year: number
  month: number
}

export const PhotoGridContainer: FC<Props> = ({ year, month }) => {
  const { data: photos } = useGetPhotos({ year, month })
  const [detailViewIndex, setDetailViewIndex] = useState<number>(0)
  const [openDetailView, setOpenDetailView] = useState<boolean>(false)

  const photoAlbum: AlbumPhoto[] = useMemo(() => {
    if (!photos) return []
    return photos.map((photo) => {
      return {
        key: photo.id,
        src: photo.thumbnail_url,
        width: photo.thumbnail_width,
        height: photo.thumbnail_height,
        alt: photo.file_name,
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
