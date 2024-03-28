import { Photo } from '@/types/photo'
import { FC } from 'react'

type Props = {
  photos: Photo[]
}

export const PhotoGrid: FC<Props> = ({ photos }) => {
  return (
    <div className="grid grid-cols-3 gap-1 p-1">
      {photos.map((photo) => (
        <img key={photo._id} src={photo.thumbnail_url} alt={photo.file_name} className="h-full w-full object-cover" />
      ))}
    </div>
  )
}
