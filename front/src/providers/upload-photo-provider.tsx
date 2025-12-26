import { Progress } from '@/components/ui/progress'
import { API_ROUTES } from '@/constants/api-route'
import { HTTP_STATUS } from '@/constants/http'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import {
  UPLOAD_MEDIA_ITEM_ERROR_CODE,
  UPLOAD_STATUS,
  UploadMediaItem,
  UploadMediaItemError,
  UploadStatus,
} from '@/types'
import { AxiosError } from 'axios'
import { createContext, Dispatch, FC, SetStateAction, useCallback, useContext, useMemo, useState } from 'react'

type UploadPhotoContextType = {
  mediaItems: UploadMediaItem[]
  setMediaItems: Dispatch<SetStateAction<UploadMediaItem[]>>
  handleUploadPhotos: (albumId: number) => Promise<void>
  clearPhotos: () => void
  progress: number
  isUploading: boolean
  reUploadPhoto: (index: number, albumId: number) => void
  isUploaded: boolean
}

export const UploadPhotoContext = createContext<UploadPhotoContextType>({
  mediaItems: [],
  setMediaItems: () => {},
  handleUploadPhotos: () => Promise.resolve(),
  clearPhotos: () => {},
  progress: 0,
  isUploading: false,
  reUploadPhoto: () => {},
  isUploaded: false,
})

export const UploadPhotoProvider: FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mediaItems, setMediaItems] = useState<UploadMediaItem[]>([])
  const [isUploading, setIsUploading] = useState(false)
  const [isUploaded, setIsUploaded] = useState(false)

  const updatePhotoStatus = (index: number, status: UploadStatus, error?: UploadMediaItemError) => {
    setMediaItems((prev) => {
      prev[index] = { ...prev[index], status, error }

      return [...prev]
    })
  }
  const uploadMediaItem = async (mediaItem: UploadMediaItem, albumId: number, index: number, retry?: boolean) => {
    try {
      updatePhotoStatus(index, UPLOAD_STATUS.PENDING)
      const formData = new FormData()
      formData.append('file', mediaItem.file)

      const response = await MyAxiosWithAuth.post(
        `${API_ROUTES.MEDIA_ITEM.UPLOAD}?album_id=${albumId}&retry=${retry ? '1' : '0'}`,
        formData,
      )
      updatePhotoStatus(index, UPLOAD_STATUS.SUCCESS)
      return response.data
    } catch (error) {
      if (error instanceof AxiosError) {
        if (error.response?.status === HTTP_STATUS.CONFLICT) {
          updatePhotoStatus(index, UPLOAD_STATUS.ERROR, {
            code: UPLOAD_MEDIA_ITEM_ERROR_CODE.DUPLICATE,
            message: error.response?.data.message ?? 'Duplicate photo',
          })
        } else {
          updatePhotoStatus(index, UPLOAD_STATUS.ERROR, {
            code: UPLOAD_MEDIA_ITEM_ERROR_CODE.INTERNAL_SERVER_ERROR,
            message: error.response?.data.message ?? 'Internal server error',
          })
        }
      }
      return null
    }
  }

  const setIsUploadingWithTimeout = useCallback((timeout: number) => {
    const timeoutId = setTimeout(() => {
      setIsUploading(false)
    }, timeout)
    return () => clearTimeout(timeoutId)
  }, [])

  const handleUploadPhotos = async (albumId: number) => {
    setIsUploaded(true)
    setIsUploading(true)
    const process = mediaItems.map((mediaItem, index) => () => uploadMediaItem(mediaItem, albumId, index))

    // 동시에 3개씩만 처리 (Promise 풀 방식)
    const concurrency = 3
    const executing: Promise<void>[] = []

    for (const task of process) {
      const promise = task().then(() => {
        executing.splice(executing.indexOf(promise), 1)
      })
      executing.push(promise)

      if (executing.length >= concurrency) {
        await Promise.race(executing)
      }
    }

    await Promise.all(executing)
    setIsUploadingWithTimeout(5000)
  }

  const clearPhotos = () => {
    setMediaItems([])
    setIsUploaded(false)
  }
  const uploadCompleted = useMemo(() => {
    return mediaItems.filter(
      (mediaItem) => mediaItem.status === UPLOAD_STATUS.SUCCESS || mediaItem.status === UPLOAD_STATUS.ERROR,
    ).length
  }, [mediaItems])

  const progress = useMemo(() => {
    if (mediaItems.length === 0) {
      return 0
    }
    return Math.round((uploadCompleted / mediaItems.length) * 100)
  }, [uploadCompleted, mediaItems])

  const reUploadPhoto = (index: number, albumId: number) => {
    const mediaItem = mediaItems[index]
    if (mediaItem.status === UPLOAD_STATUS.ERROR) {
      uploadMediaItem(mediaItem, albumId, index, true)
    }
  }
  return (
    <UploadPhotoContext.Provider
      value={{
        mediaItems,
        setMediaItems,
        handleUploadPhotos,
        clearPhotos,
        progress,
        isUploading,
        reUploadPhoto,
        isUploaded,
      }}
    >
      {isUploading && (
        <div className="fixed left-0 right-0 top-5 rounded-md bg-white/90 px-4">
          <span className="text-gray-500 text-sm">
            {uploadCompleted} / {mediaItems.length}
          </span>
          <Progress value={progress} className="w-full" />
        </div>
      )}
      {children}
    </UploadPhotoContext.Provider>
  )
}

export const useUploadPhoto = () => {
  return useContext(UploadPhotoContext)
}
