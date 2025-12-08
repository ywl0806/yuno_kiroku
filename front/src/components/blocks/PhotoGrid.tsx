import { Photo } from '@/types/photo'
import { resizeImageView } from '@/utils/calculateImageSize'
import { FC, useMemo } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import 'react-lazy-load-image-component/src/effects/blur.css'
import PhotoAlbum, { Photo as AlbumPhoto } from 'react-photo-album'

type Props = {
  photos: Photo[]
  onClick?: (index: number) => void
}

export const PhotoGrid: FC<Props> = ({ photos, onClick }) => {
  const photoAlbum: AlbumPhoto[] = useMemo(() => {
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
    <PhotoAlbum
      onClick={(props) => {
        onClick && onClick(props.index)
      }}
      photos={photoAlbum}
      columns={2}
      layout="masonry"
      spacing={2}
      padding={2}
      targetRowHeight={300}
      renderPhoto={({ photo, wrapperStyle, ...props }) => {
        return <LazyLoadImage src={photo.src} alt={photo.alt} style={wrapperStyle} effect="blur" {...props} />
      }}
    />
  )
}
