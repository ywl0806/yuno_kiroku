import { Photo } from '@/types/photo'
import { resizeImageView } from '@/utils/calculateImageSize'
import { FC, useMemo } from 'react'
import PhotoAlbum, { Photo as AlbumPhoto } from 'react-photo-album'

type Props = {
  photos: Photo[]
}

export const PhotoGrid: FC<Props> = ({ photos }) => {
  const photoAlbum: AlbumPhoto[] = useMemo(() => {
    return photos.map((photo) => {
      const { width, height } = resizeImageView(photo.width, photo.height, 1500, 1500, photo.orientation)

      return {
        key: photo._id,
        src: photo.thumbnail_url,
        width: width,
        height: height,
        alt: photo.file_name,
      }
    })
  }, [photos])
  return <PhotoAlbum photos={photoAlbum} columns={2} layout="masonry" spacing={2} padding={2} targetRowHeight={300} />
}
