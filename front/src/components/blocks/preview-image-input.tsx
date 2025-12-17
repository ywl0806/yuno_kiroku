import { PhotoGrid } from './photo-grid'
import { Button } from '@/components/ui/button'
import { UPLOAD_STATUS, UploadPhoto } from '@/types'
import { Check, Loader2, Plus, X } from 'lucide-react'
import { useMemo, useState, useEffect, useRef, FC, Dispatch, SetStateAction, RefObject } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'

interface ImageDimensions {
  width: number
  height: number
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
  images: UploadPhoto[]
  setImages: Dispatch<SetStateAction<UploadPhoto[]>>
}

export const PreviewImageInput: FC<Props> = ({ inputRef, images, setImages }) => {
  const [imageDimensions, setImageDimensions] = useState<ImageDimensions[]>([])

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

  return (
    <div className="flex h-full flex-col gap-4 p-4">
      <div className="sticky top-2 z-50 flex items-center justify-center gap-2">
        <Button variant="outline" className="rounded-full" onClick={() => imgInputRef.current?.click()}>
          <Plus />
          <span>이미지를 추가하세요</span>
        </Button>
        <input ref={imgInputRef} hidden type="file" accept="image/*" name="images" multiple onChange={handleChange} />
      </div>

      {images.length > 0 && (
        <div className="overflow-y-visible rounded-md border-2 p-2">
          <PhotoGrid
            photos={photoAlbum}
            renderPhoto={(props) => (
              <div className="relative" key={props.photo.key}>
                <div className="absolute right-2 top-2 z-10">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-full"
                    onClick={() => handleRemove(props.layout.index)}
                    disabled={
                      images[props.layout.index].status === UPLOAD_STATUS.PENDING ||
                      images[props.layout.index].status === UPLOAD_STATUS.SUCCESS
                    }
                  >
                    <X />
                  </Button>
                </div>
                <div className="bg-gray-500/40 absolute inset-0 flex items-center justify-center">
                  {images[props.layout.index].status === UPLOAD_STATUS.PENDING && <Loader2 className="animate-spin" />}
                  {images[props.layout.index].status === UPLOAD_STATUS.SUCCESS && <Check className="text-white" />}
                  {images[props.layout.index].status === UPLOAD_STATUS.ERROR && <X className="text-red-500" />}
                </div>
                {props.renderDefaultPhoto()}
              </div>
            )}
          />
        </div>
      )}
    </div>
  )
}
