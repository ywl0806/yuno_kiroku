import { Progress } from '@/components/ui/progress'
import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import {
  UPLOAD_MEDIA_ITEM_ERROR_CODE,

  UploadMediaItem,
  UploadMediaItemError,
} from '@/types'
import { UploadStatus, UPLOAD_STATUS } from '@/enums'
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
  setMediaItems: () => { },
  handleUploadPhotos: () => Promise.resolve(),
  clearPhotos: () => { },
  progress: 0,
  isUploading: false,
  reUploadPhoto: () => { },
  isUploaded: false,
})

export const UploadPhotoProvider: FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mediaItems, setMediaItems] = useState<UploadMediaItem[]>([])
  const [isUploading, setIsUploading] = useState(false)
  const [isUploaded, setIsUploaded] = useState(false)
  const [currentUploadBatchId, setCurrentUploadBatchId] = useState<number | null>(null)

  const updatePhotoStatus = (index: number, status: UploadStatus, error?: UploadMediaItemError) => {
    setMediaItems((prev) => {
      prev[index] = { ...prev[index], status, error }

      return [...prev]
    })
  }
  const uploadMediaItem = async (
    mediaItem: UploadMediaItem,
    albumId: number,
    uploadBatchId: number,
    index: number,
    retry?: boolean,
  ) => {
    try {
      updatePhotoStatus(index, UPLOAD_STATUS.PENDING)
      const formData = new FormData()
      formData.append('file', mediaItem.file)

      const response = await MyAxiosWithAuth.post(
        `${API_ROUTES.MEDIA_ITEM.UPLOAD}?album_id=${albumId}&upload_batch_id=${uploadBatchId}&retry=${retry ? '1' : '0'
        }`,
        formData,
      )
      const mediaItemId = response.data?.media_item_id as number | undefined
      setMediaItems((prev) => {
        const next = [...prev]
        next[index] = {
          ...next[index],
          status: UPLOAD_STATUS.COMPLETED,
          ...(mediaItemId != null && { media_item_id: mediaItemId }),
        }
        return next
      })

      if (retry) {
        pollUploadBatchStatus(uploadBatchId)
        return
      }
      return response.data
    } catch (error) {
      if (error instanceof AxiosError) {
        updatePhotoStatus(index, UPLOAD_STATUS.FAILED, {
          code: UPLOAD_MEDIA_ITEM_ERROR_CODE.INTERNAL_SERVER_ERROR,
          message: error.response?.data.message ?? 'Internal server error',
        })

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

    // 배치 생성 (기존 배치가 없을 때만)
    let uploadBatchId: number
    if (currentUploadBatchId) {
      uploadBatchId = currentUploadBatchId
    } else {
      try {
        const batchResponse = await MyAxiosWithAuth.post(`${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH}?album_id=${albumId}`)
        uploadBatchId = batchResponse.data.id
        setCurrentUploadBatchId(uploadBatchId)
      } catch (error) {
        console.error('배치 생성 실패:', error)
        setIsUploading(false)
        return
      }
    }

    // 파일 업로드 (원본만 빠르게 업로드)
    const process = mediaItems.map(
      (mediaItem, index) => () => uploadMediaItem(mediaItem, albumId, uploadBatchId, index),
    )

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

    // 상태 polling 시작
    pollUploadBatchStatus(uploadBatchId)
  }

  const pollUploadBatchStatus = async (uploadBatchId: number) => {

    if (isUploading) {
      return
    }

    const maxAttempts = 60 // 최대 5분 (5초마다 polling)
    let attempts = 0

    const poll = async () => {
      try {
        const response = await MyAxiosWithAuth.get<UploadBatchStatusResponse>(
          `${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH_STATUS}?upload_batch_id=${uploadBatchId}`,
        )

        const statuseMap = new Map<number, UploadStatus>()
        response.data.statuses.forEach((status) => {
          statuseMap.set(status.id, status.upload_status)
        })

        setMediaItems((prev) => {
          const next = [...prev]
          next.forEach((mediaItem) => {
            mediaItem.status = statuseMap.get(Number(mediaItem.media_item_id)) ?? UPLOAD_STATUS.PENDING
          })
          return next
        })
        if (response.data.is_completed) {
          setIsUploadingWithTimeout(5000)
          return
        }

        attempts++
        if (attempts < maxAttempts) {
          setTimeout(poll, 5000) // 5초마다 polling
        } else {
          setIsUploadingWithTimeout(5000)
        }
      } catch (error) {
        console.error('상태 조회 실패:', error)
        setIsUploadingWithTimeout(5000)
      }
    }

    poll()
  }

  const clearPhotos = () => {
    setMediaItems([])
    setIsUploaded(false)
    setCurrentUploadBatchId(null)
  }
  const uploadCompleted = useMemo(() => {
    return mediaItems.filter(
      (mediaItem) => mediaItem.status === UPLOAD_STATUS.COMPLETED || mediaItem.status === UPLOAD_STATUS.FAILED || mediaItem.status === UPLOAD_STATUS.DUPLICATE,
    ).length
  }, [mediaItems])

  const progress = useMemo(() => {
    if (mediaItems.length === 0) {
      return 0
    }
    return Math.round((uploadCompleted / mediaItems.length) * 100)
  }, [uploadCompleted, mediaItems])

  const reUploadPhoto = async (index: number, albumId: number) => {
    const mediaItem = mediaItems[index]
    if (mediaItem.status === UPLOAD_STATUS.FAILED || mediaItem.status === UPLOAD_STATUS.DUPLICATE) {
      // 재업로드 시 같은 배치 ID 사용 (기존 배치가 없으면 새로 생성)
      let uploadBatchId: number
      if (currentUploadBatchId) {
        uploadBatchId = currentUploadBatchId
      } else {
        try {
          const batchResponse = await MyAxiosWithAuth.post(`${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH}?album_id=${albumId}`)
          uploadBatchId = batchResponse.data.id
          setCurrentUploadBatchId(uploadBatchId)
        } catch (error) {
          console.error('배치 생성 실패:', error)
          return
        }
      }

      try {
        await uploadMediaItem(mediaItem, albumId, uploadBatchId, index, true)
        pollUploadBatchStatus(uploadBatchId)
      } catch (error) {
        console.error('재업로드 실패:', error)
      }
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

type UploadBatchStatusResponse = {
  statuses: {
    id: number
    upload_status: UploadStatus
  }[]
  is_completed: boolean
}