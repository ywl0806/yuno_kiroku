import { Play } from 'lucide-react'
import { FC, useEffect, useState } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import 'react-lazy-load-image-component/src/effects/blur.css'
import PhotoAlbum, { Photo as AlbumPhoto, RenderPhoto } from 'react-photo-album'

type Photo = AlbumPhoto & {
  isVideo?: boolean
}

type Props = {
  photos: Photo[]
  onClick?: (index: number) => void
  renderPhoto?: RenderPhoto<Photo>
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
      componentsProps={{
        columnContainerProps: {
          className: 'relative',
        }
      }}
      renderPhoto={(props) => {
        const { photo, ...rest } = props

        return renderPhoto ? (
          renderPhoto(props)
        ) : (

          <>
            {photo.isVideo && (
              <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
                <div className="rounded-full bg-black/40 p-2 z-50">
                  <Play className="size-5 fill-white text-white" />
                </div>
              </div>
            )}
            <LazyLoadImage
              className="p-[2px]"
              src={photo.src}
              alt={photo.alt}
              effect="blur"
              onClick={rest.imageProps.onClick}
              {...rest}

            />
          </>
        )
      }}
    />
  )
}
