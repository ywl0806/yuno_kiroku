import { PreviewImageInput } from '@/components/blocks/preview-image-input'
import { Button } from '@/components/ui/button'
import { UploadPhoto } from '@/feature/home/types'
import { useGetAlbums } from '@/feature/upload/hooks/use-get-albums'
import { useUploadPhotos } from '@/feature/upload/hooks/use-upload-photos'
import { Album, Upload } from 'lucide-react'
import { FC, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'

type Props = {
  albumId: number | null
}

export const UploadContainer: FC<Props> = ({ albumId }) => {
  const { data: albums } = useGetAlbums()
  const selectedAlbum = useMemo(() => albums?.find((album) => album.id === albumId), [albums, albumId])
  const [images, setImages] = useState<UploadPhoto[]>([])

  const { handleUploadPhotos, uploadPhotoStatus } = useUploadPhotos({ photos: images, albumId: selectedAlbum?.id ?? 0 })

  return (
    <div className="h-full pt-4">
      {albumId === null && (
        <div className="flex h-full select-none flex-col justify-center gap-2">
          {albums?.map((album) => (
            <Button asChild key={album.id} variant="outline" className=" justify-start">
              <Link to={`/upload?album_id=${album.id}`}>
                <Album />
                {album.name}
              </Link>
            </Button>
          ))}
        </div>
      )}
      {selectedAlbum && (
        <div className="h-full overflow-y-auto">
          <div className="flex items-center justify-center gap-2">
            <Album />
            {selectedAlbum.name}
          </div>
          <PreviewImageInput images={images} setImages={setImages} uploadPhotoStatus={uploadPhotoStatus} />
          <div className="absolute bottom-20 left-0 right-0 flex w-full justify-center p-4">
            <Button variant="outline" onClick={handleUploadPhotos}>
              <Upload />
              Upload
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
