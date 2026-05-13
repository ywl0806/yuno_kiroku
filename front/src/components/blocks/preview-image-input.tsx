import { PhotoGrid } from './photo-grid'
import { Button } from '@/components/ui/button'
import { UPLOAD_STATUS } from '@/enums'
import { useUploadPhoto } from '@/providers/upload-photo-provider'
import { UploadMediaItem } from '@/types'
import { BrushCleaning, Check, Loader2, Play, Plus, RefreshCcw, X } from 'lucide-react'
import { useMemo, useState, useEffect, useRef, FC, Dispatch, SetStateAction, RefObject } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'
import { useTranslation } from 'react-i18next'

interface ImageDimensions {
  width: number
  height: number
}
const getColumnCount = (width: number) => {
  return Math.abs(Math.floor(width / 300))
}

const getMediaDimensions = (file: File): Promise<ImageDimensions> => {
  if (file.type.startsWith('video/')) {
    return new Promise((resolve) => {
      const video = document.createElement('video')
      const url = URL.createObjectURL(file)
      video.onloadedmetadata = () => {
        URL.revokeObjectURL(url)
        resolve({ width: video.videoWidth || 16, height: video.videoHeight || 9 })
      }
      video.onerror = () => {
        URL.revokeObjectURL(url)
        resolve({ width: 16, height: 9 })
      }
      video.src = url
    })
  }
  return new Promise((resolve, reject) => {
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      URL.revokeObjectURL(url)
      resolve({ width: img.naturalWidth, height: img.naturalHeight })
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
  albumId: string
}

export const PreviewImageInput: FC<Props> = ({ inputRef, images, setImages, albumId }) => {
  const [imageDimensions, setImageDimensions] = useState<ImageDimensions[]>([])
  const { reUploadPhoto, isUploaded, isUploading, progress, clearPhotos, mediaItems } = useUploadPhoto()
  const imgInputRef = inputRef ?? useRef<HTMLInputElement>(null)
  const { t } = useTranslation()
  // 이미지 크기를 비동기로 로드
  useEffect(() => {
    const loadDimensions = async () => {
      const dimensions = await Promise.all(images.map((image) => getMediaDimensions(image.file)))
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
        ...Array.from(files).map((file) => ({ file, src: URL.createObjectURL(file), status: UPLOAD_STATUS.PENDING })),
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
    <div className="flex flex-col gap-4 p-4">
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
            <span>{t('upload.addImage')}</span>
          </Button>
        </div>
        <div className="flex flex-1 items-center justify-end">
          {(!isUploading || progress === 100) && mediaItems.length > 0 && (
            <Button variant="outline" className="rounded-full" onClick={clearPhotos}>
              <BrushCleaning className="size-4" />
              <span>{t('upload.clear')}</span>
            </Button>
          )}
        </div>
        <input ref={imgInputRef} hidden type="file" accept="image/*, video/*" name="images" multiple onChange={handleChange} />
      </div>

      {images.length > 0 && (
        <div className="overflow-y-visible rounded-md border-2 p-2">
          <PhotoGrid
            photos={photoAlbum}
            columnCount={columnCount}
            renderPhoto={(props) => {
              const isVideo = images[props.layout.index]?.file.type.startsWith('video/')
              return (
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
                {isUploaded && (() => {
                  const status = images[props.layout.index].status
                  const overlayBg =
                    status === UPLOAD_STATUS.COMPLETED
                      ? 'bg-green-200/40'
                      : status === UPLOAD_STATUS.FAILED
                        ? 'bg-red-200/40'
                        : 'bg-black/40'
                  return (
                    <div className={`absolute inset-0 flex items-center justify-center ${overlayBg}`}>
                      {status === UPLOAD_STATUS.PENDING && (
                        <Loader2 className="animate-spin text-white drop-shadow" />
                      )}

                      {status === UPLOAD_STATUS.PROCESSING && (
                        <div className="flex flex-col items-center gap-1.5">
                          <Loader2 className="animate-spin text-white drop-shadow" />
                          <span className="rounded-full bg-black/50 px-2 py-0.5 text-xs text-white">처리 중</span>
                        </div>
                      )}

                      {status === UPLOAD_STATUS.COMPLETED && (
                        <Check className="size-7 text-white drop-shadow" strokeWidth={3} />
                      )}

                      {status === UPLOAD_STATUS.DUPLICATE && (
                        <div className="flex flex-col items-center gap-2">
                          <p className="rounded-md bg-white/90 px-2 py-1 text-xs text-slate-600">{t('upload.alreadyUploaded')}</p>
                          <Button
                            variant="default"
                            size="sm"
                            className="rounded-full"
                            onClick={() => reUploadPhoto(props.layout.index, albumId)}
                          >
                            <RefreshCcw className="size-3" />
                            <span>{t('upload.uploadAnyway')}</span>
                          </Button>
                        </div>
                      )}

                      {status === UPLOAD_STATUS.FAILED && (
                        <div className="flex flex-col items-center gap-2">
                          <p className="rounded-md bg-white/90 px-2 py-1 text-xs text-red-500">{t('upload.uploadFailed')}</p>
                          <Button
                            size="sm"
                            className="rounded-full bg-white text-gray-700 hover:bg-gray-100"
                            onClick={() => reUploadPhoto(props.layout.index, albumId)}
                          >
                            <RefreshCcw className="size-3" />
                            <span>{t('upload.reUpload')}</span>
                          </Button>
                        </div>
                      )}
                    </div>
                  )
                })()}
                {isVideo && !isUploaded && (
                  <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
                    <div className="rounded-full bg-black/40 p-2">
                      <Play className="size-5 fill-white text-white" />
                    </div>
                  </div>
                )}
                {isVideo ? (
                  <video
                    src={images[props.layout.index].src}
                    style={{ width: props.layout.width, height: props.layout.height, objectFit: 'cover', display: 'block' }}
                    muted
                    playsInline
                    preload="metadata"
                  />
                ) : (
                  props.renderDefaultPhoto()
                )}
              </div>
              )
            }}
          />
        </div>
      )}
    </div>
  )
}
