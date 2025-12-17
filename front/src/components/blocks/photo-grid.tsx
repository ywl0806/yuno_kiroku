import { FC, useEffect, useState } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import 'react-lazy-load-image-component/src/effects/blur.css'
import PhotoAlbum, { Photo as AlbumPhoto, RenderPhoto } from 'react-photo-album'

type Props = {
  photos: AlbumPhoto[]
  onClick?: (index: number) => void
  renderPhoto?: RenderPhoto<AlbumPhoto>
}
const getColumnCount = (width: number) => {
  if (width < 1024) {
    return 2
  }
  return 3
}

export const PhotoGrid: FC<Props> = ({ photos, onClick, renderPhoto }) => {
  const [columnCount, setColumnCount] = useState(getColumnCount(window.innerWidth))

  useEffect(() => {
    const handleResize = () => {
      setColumnCount(getColumnCount(window.innerWidth))
    }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])
  return (
    <PhotoAlbum
      onClick={(props) => {
        onClick && onClick(props.index)
      }}
      photos={photos}
      columns={columnCount}
      layout="masonry"
      spacing={2}
      padding={2}
      targetRowHeight={300}
      renderPhoto={(props) => {
        const { photo, wrapperStyle, ...rest } = props

        return renderPhoto ? (
          renderPhoto(props)
        ) : (
          <LazyLoadImage src={photo.src} alt={photo.alt} effect="blur" onClick={rest.imageProps.onClick} {...rest} />
        )
      }}
    />
  )
}
