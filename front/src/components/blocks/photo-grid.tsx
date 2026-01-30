import { FC, useEffect, useState } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import 'react-lazy-load-image-component/src/effects/blur.css'
import PhotoAlbum, { Photo as AlbumPhoto, RenderPhoto } from 'react-photo-album'

type Props = {
  photos: AlbumPhoto[]
  onClick?: (index: number) => void
  renderPhoto?: RenderPhoto<AlbumPhoto>
  columnCount?: number
}
const getColumnCount = (width: number) => {
  return Math.abs(Math.floor(width / 300)) + 1
}

export const PhotoGrid: FC<Props> = ({ photos, onClick, renderPhoto, columnCount: _columnCount }) => {
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
      columns={_columnCount || columnCount}
      layout="masonry"
      spacing={2}
      targetRowHeight={300}
      renderPhoto={(props) => {
        const { photo, ...rest } = props

        return renderPhoto ? (
          renderPhoto(props)
        ) : (
          <LazyLoadImage
            className="p-[2px]"
            src={photo.src}
            alt={photo.alt}
            effect="blur"
            onClick={rest.imageProps.onClick}
            {...rest}
          />
        )
      }}
    />
  )
}
