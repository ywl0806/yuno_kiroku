import { PreviewImageInput } from '@/components/blocks/preview-image-input'
import { Button } from '@/components/ui/button'
import { useGetAlbums } from '@/feature/upload/hooks/use-get-albums'
import { getSessionStorage, removeSessionStorage, SESSION_STORAGE_KEY, setSessionStorage } from '@/lib/session-storage'
import { useUploadPhoto } from '@/providers/upload-photo-provider'
import { Album, ArrowLeft, Upload } from 'lucide-react'
import { FC, useMemo, useRef, useState } from 'react'

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

  const { mediaItems, setMediaItems, handleUploadPhotos, clearPhotos, isUploaded } = useUploadPhoto()
  const imgInputRef = useRef<HTMLInputElement>(null)

  const handleChangeAlbum = (albumId: number) => {
    setAlbumId(albumId)
    setSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID, albumId.toString())
  }

  const handleClearPhotos = () => {
    clearPhotos()
    setAlbumId(null)
    removeSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID)
  }
  return (
    <div className="h-full w-full px-5 pt-5">
      {albumId === null && (
        <div className="mx-auto flex h-full max-w-[20rem] select-none flex-col justify-center gap-5">
          {albums?.map((album) => (
            <Button
              key={album.id}
              variant="outline"
              className="h-12 justify-start text-[1rem]"
              onClick={() => handleChangeAlbum(album.id)}
            >
              <Album />
              {album.name}
            </Button>
          ))}
        </div>
      )}
      {selectedAlbum && (
        <div className="h-full w-full overflow-y-auto">
          <div className="flex w-full items-center justify-between ">
            <div className="flex-1">
              <Button variant="outline" size="sm" onClick={handleClearPhotos}>
                <ArrowLeft />
                <span>戻る</span>
              </Button>
            </div>
            <div className="flex flex-1 items-center justify-center gap-2 whitespace-nowrap text-sm">
              <Album className="size-4" />
              {selectedAlbum.name}
            </div>
            <div className="flex-1" />
          </div>
          <PreviewImageInput
            inputRef={imgInputRef}
            images={mediaItems}
            setImages={setMediaItems}
            albumId={selectedAlbum.id}
          />
          <div className="absolute bottom-20 left-0 right-0 flex w-full justify-center gap-5 p-4">
            <div className="flex items-center justify-center">
              {!isUploaded && mediaItems.length > 0 && (
                <Button variant="outline" onClick={() => handleUploadPhotos(selectedAlbum.id)}>
                  <Upload />
                  <span>アップロード</span>
                </Button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
