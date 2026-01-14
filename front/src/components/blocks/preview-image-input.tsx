import { PhotoGrid } from './photo-grid'
import { Button } from '@/components/ui/button'
import { useUploadPhoto } from '@/providers/upload-photo-provider'
import { UPLOAD_MEDIA_ITEM_ERROR_CODE, UPLOAD_STATUS, UploadMediaItem } from '@/types'
import { BrushCleaning, Check, Loader2, Plus, RefreshCcw, X } from 'lucide-react'
import { useMemo, useState, useEffect, useRef, FC, Dispatch, SetStateAction, RefObject } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'

interface ImageDimensions {
  width: number
  height: number
}
const getColumnCount = (width: number) => {
  return Math.abs(Math.floor(width / 300))
}

// 이미지의 실제 크기를 가져오는 헬퍼 함수
const getImageDimensions = (file: File): Promise<ImageDimensions> => {
  return new Promise((resolve, reject) => {
    const img = new Image()
    const url = URL.createObjectURL(file)

    img.onload = () => {
      URL.revokeObjectURL(url)
      resolve({
        width: img.naturalWidth,
        height: img.naturalHeight,
      })
    }

    img.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new Error('Failed to load image'))
    }

    img.src = url
  })
}

type Props = {
  inputRef?: RefObject<HTMLInputElement>
  images: UploadMediaItem[]
  setImages: Dispatch<SetStateAction<UploadMediaItem[]>>
  albumId: number
}

export const PreviewImageInput: FC<Props> = ({ inputRef, images, setImages, albumId }) => {
  const [imageDimensions, setImageDimensions] = useState<ImageDimensions[]>([])
  const { reUploadPhoto, isUploaded, isUploading, progress, clearPhotos, mediaItems } = useUploadPhoto()
  const imgInputRef = inputRef ?? useRef<HTMLInputElement>(null)

  // 이미지 크기를 비동기로 로드
  useEffect(() => {
    const loadDimensions = async () => {
      const dimensions = await Promise.all(images.map((image) => getImageDimensions(image.file)))
      setImageDimensions(dimensions)
    }

    if (images.length > 0) {
      loadDimensions()
    } else {
      setImageDimensions([])
    }
  }, [images])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (isUploaded) return
    const files = e.target.files
    if (files) {
      setImages((prev) => [
        ...prev,
        ...Array.from(files).map((file) => ({ file, src: URL.createObjectURL(file), status: UPLOAD_STATUS.IDLE })),
      ])
    }
  }

  const photoAlbum: AlbumPhoto[] = useMemo(() => {
    return images.map((image, index) => {
      // 이미지 크기가 로드되면 실제 크기 사용, 아니면 기본값 100x100
      const dimensions = imageDimensions[index] || { width: 100, height: 100 }

      return {
        key: image.src,
        src: image.src,
        width: dimensions.width,
        height: dimensions.height,
        alt: 'preview',
      }
    })
  }, [images, imageDimensions])

  const handleRemove = (index: number) => {
    setImages((prev) => prev.filter((_, i) => i !== index))
  }
  const [columnCount, setColumnCount] = useState(getColumnCount(window.innerWidth))

  useEffect(() => {
    const handleResize = () => {
      setColumnCount(getColumnCount(window.innerWidth))
    }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])
  return (
    <div className="flex h-full flex-col gap-4 p-4">
      <div className="sticky top-2 z-50 flex items-center justify-center gap-2">
        <div className="flex-1" />
        <div className="flex flex-1 items-center justify-center">
          <Button
            variant="outline"
            className="rounded-full"
            onClick={() => imgInputRef.current?.click()}
            disabled={isUploaded}
          >
            <Plus />
            <span>이미지를 추가하세요</span>
          </Button>
        </div>
        <div className="flex flex-1 items-center justify-end">
          {(!isUploading || progress === 100) && mediaItems.length > 0 && (
            <Button variant="outline" className="rounded-full" onClick={clearPhotos}>
              <BrushCleaning className="size-4" />
              <span>クリア</span>
            </Button>
          )}
        </div>
        <input ref={imgInputRef} hidden type="file" accept="image/*" name="images" multiple onChange={handleChange} />
      </div>

      {images.length > 0 && (
        <div className="overflow-y-visible rounded-md border-2 p-2">
          <PhotoGrid
            photos={photoAlbum}
            columnCount={columnCount}
            renderPhoto={(props) => (
              <div className="relative" key={props.photo.key}>
                {!isUploaded && (
                  <div className="absolute right-2 top-2 z-10">
                    <Button
                      variant="outline"
                      size="icon"
                      className="size-7 rounded-full"
                      onClick={() => handleRemove(props.layout.index)}
                      disabled={isUploaded}
                    >
                      <X />
                    </Button>
                  </div>
                )}
                {isUploaded && (
                  <div className="bg-gray-500/40 absolute inset-0 flex items-center justify-center">
                    {images[props.layout.index].status === UPLOAD_STATUS.PENDING && (
                      <Loader2 className="animate-spin text-white" />
                    )}
                    {images[props.layout.index].status === UPLOAD_STATUS.SUCCESS && <Check className="text-white" />}
                    {images[props.layout.index].status === UPLOAD_STATUS.ERROR ? (
                      images[props.layout.index].error?.code === UPLOAD_MEDIA_ITEM_ERROR_CODE.DUPLICATE ? (
                        <div className="flex flex-col items-center gap-2">
                          <p className="rounded-md bg-white p-2 text-sm text-slate-500">すでに上げてるかも</p>
                          <Button
                            variant="default"
                            className="rounded-full"
                            onClick={() => reUploadPhoto(props.layout.index, albumId)}
                          >
                            <RefreshCcw />
                            <span>けれどアップロードする</span>
                          </Button>
                        </div>
                      ) : (
                        <div className="flex flex-col items-center gap-2">
                          <X className="text-red-500" />
                          <p className="rounded-full bg-white p-2 text-sm text-red-500">アップロード失敗</p>
                          <p className="rounded-md bg-white/90 p-2 text-sm">
                            {images[props.layout.index].error?.message}
                          </p>
                          <Button
                            variant="ghost"
                            className="rounded-full bg-white"
                            onClick={() => reUploadPhoto(props.layout.index, albumId)}
                          >
                            <RefreshCcw />
                            <span>再アップロードする</span>
                          </Button>
                        </div>
                      )
                    ) : null}
                  </div>
                )}
                {props.renderDefaultPhoto()}
              </div>
            )}
          />
        </div>
      )}
    </div>
  )
}
