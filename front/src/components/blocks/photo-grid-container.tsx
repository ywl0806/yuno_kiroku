import { PhotoDetailSwipeDialog } from './photo-detail-swipe-dialog'
import { PhotoGrid } from './photo-grid'
import { useGetPhotos } from '@/feature/home/hooks/use-get-photos'
import { resizeImageView } from '@/utils/calculateImageSize'
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
      const { width, height } = resizeImageView(photo.width, photo.height, 1500, 1500, photo.orientation)

      return {
        key: photo.id,
        src: photo.thumbnail_url,
        width: width,
        height: height,
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
