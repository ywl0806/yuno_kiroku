import { PreviewImageInput } from '@/components/blocks/preview-image-input'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { useGetAlbums } from '@/feature/upload/hooks/use-get-albums'
import { getSessionStorage, SESSION_STORAGE_KEY, setSessionStorage } from '@/lib/session-storage'
import { useUploadPhoto } from '@/providers/upload-photo-provider'
import { Album, Check, Upload } from 'lucide-react'
import { FC, useEffect, useMemo, useRef, useState } from 'react'

export const UploadContainer: FC = () => {
  const { data: albums } = useGetAlbums()
  const [albumId, setAlbumId] = useState<number | null>(
    getSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID)
      ? Number(getSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID))
      : null,
  )
  const selectedAlbum = useMemo(() => {
    if (albumId === null) {
      return null
    }
    return albums?.find((album) => album.id === albumId)
  }, [albums, albumId])

  const { photos, setPhotos, handleUploadPhotos, progress, isUploading, clearPhotos } = useUploadPhoto()
  const imgInputRef = useRef<HTMLInputElement>(null)

  // useEffect(() => {
  //   if (progress !== 100) {
  //     imgInputRef.current?.click()
  //   }
  // }, [])

  const handleChangeAlbum = (albumId: number) => {
    setAlbumId(albumId)
    setSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID, albumId.toString())
  }
  return (
    <div className="h-full pt-4">
      {albumId === null && (
        <div className="flex h-full select-none flex-col justify-center gap-2">
          {albums?.map((album) => (
            <Button
              key={album.id}
              variant="outline"
              className=" justify-start"
              onClick={() => handleChangeAlbum(album.id)}
            >
              <Album />
              {album.name}
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
          <PreviewImageInput inputRef={imgInputRef} images={photos} setImages={setPhotos} />
          <div className="absolute bottom-20 left-0 right-0 flex w-full justify-center p-4">
            {progress === 100 ? (
              <Button onClick={clearPhotos}>
                <Check />
                Clear
              </Button>
            ) : (
              !isUploading &&
              progress === 0 &&
              photos.length > 0 && (
                <Button variant="outline" onClick={() => handleUploadPhotos(selectedAlbum.id)}>
                  <Upload />
                  Upload
                </Button>
              )
            )}
          </div>
        </div>
      )}
    </div>
  )
}
