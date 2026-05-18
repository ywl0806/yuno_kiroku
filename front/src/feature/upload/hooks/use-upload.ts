import { API_ROUTES } from '@/consts/api-route'
import { UPLOAD_STATUS } from '@/enums'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { UPLOAD_MEDIA_ITEM_ERROR_CODE, UploadMediaItem, UploadMediaItemError } from '@/types'
import { UploadStatus } from '@/enums'
import { AxiosError } from 'axios'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { unstable_useBlocker as useBlocker } from 'react-router-dom'
import { UPLOAD_BATCH_STORAGE_KEY, StoredBatch, UploadBatchStatusResponse } from '../constants'

export function useUpload() {
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

  const setIsUploadingWithTimeout = useCallback((timeout: number) => {
    const timeoutId = setTimeout(() => setIsUploading(false), timeout)
    return () => clearTimeout(timeoutId)
  }, [])

  const pollUploadBatchStatus = async (uploadBatchId: number) => {
    if (isUploading) return

    const maxAttempts = 60
    let attempts = 0

    const poll = async () => {
      try {
        const response = await MyAxiosWithAuth.get<UploadBatchStatusResponse>(
          `${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH_STATUS}?upload_batch_id=${uploadBatchId}`,
        )

        const statusMap = new Map<string, { status: UploadStatus; failure_reason?: string }>()
        response.data.statuses.forEach((s) => statusMap.set(s.id, { status: s.upload_status, failure_reason: s.failure_reason }))

        setMediaItems((prev) => {
          const next = [...prev]
          next.forEach((item) => {
            const info = statusMap.get(item.media_item_id ?? '')
            if (info) {
              item.status = info.status
              item.failure_reason = info.failure_reason
            } else {
              item.status = UPLOAD_STATUS.PENDING
            }
          })
          return next
        })

        if (response.data.is_completed) {
          localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
          setIsUploadingWithTimeout(5000)
          return
        }

        attempts++
        if (attempts < maxAttempts) {
          setTimeout(poll, 5000)
        } else {
          localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
          setIsUploadingWithTimeout(5000)
        }
      } catch {
        console.error('상태 조회 실패')
        localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
        setIsUploadingWithTimeout(5000)
      }
    }

    poll()
  }

  const getOrCreateBatchId = async (albumId: string): Promise<number | null> => {
    if (currentUploadBatchId) return currentUploadBatchId
    try {
      const res = await MyAxiosWithAuth.post(`${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH}?album_id=${albumId}`)
      const id: number = res.data.id
      setCurrentUploadBatchId(id)
      return id
    } catch {
      console.error('배치 생성 실패')
      return null
    }
  }

  // 단건 재업로드 전용
  const uploadMediaItem = async (
    mediaItem: UploadMediaItem,
    albumId: string,
    uploadBatchId: number,
    index: number,
  ) => {
    try {
      updatePhotoStatus(index, UPLOAD_STATUS.PENDING)

      const presignedRes = await MyAxiosWithAuth.post(API_ROUTES.MEDIA_ITEM.PRESIGNED_URL, {
        album_id: albumId,
        upload_batch_id: uploadBatchId,
        file_name: mediaItem.file.name,
        content_type: mediaItem.file.type || 'image/jpeg',
      })
      const { media_item_id: mediaItemId, presigned_url: presignedUrl } = presignedRes.data as {
        media_item_id: string
        presigned_url: string
      }

      await fetch(presignedUrl, {
        method: 'PUT',
        body: mediaItem.file,
        headers: { 'Content-Type': mediaItem.file.type || 'image/jpeg' },
      })

      setMediaItems((prev) => {
        const next = [...prev]
        next[index] = { ...next[index], status: UPLOAD_STATUS.PROCESSING, media_item_id: mediaItemId }
        return next
      })
    } catch (error) {
      if (error instanceof AxiosError) {
        updatePhotoStatus(index, UPLOAD_STATUS.FAILED, {
          code: UPLOAD_MEDIA_ITEM_ERROR_CODE.INTERNAL_SERVER_ERROR,
          message: error.response?.data.message ?? 'Internal server error',
        })
      }
    }
  }

  const handleUploadPhotos = async (albumId: string) => {
    setIsUploaded(true)
    setIsUploading(true)

    const uploadBatchId = await getOrCreateBatchId(albumId)
    if (!uploadBatchId) {
      setIsUploading(false)
      return
    }

    type PresignedItem = { mediaItemId: string; presignedUrl: string; index: number }
    let presignedItems: PresignedItem[]
    try {
      const presignedRes = await MyAxiosWithAuth.post(API_ROUTES.MEDIA_ITEM.PRESIGNED_URLS, {
        album_id: albumId,
        upload_batch_id: uploadBatchId,
        files: mediaItems.map((item) => ({
          file_name: item.file.name,
          content_type: item.file.type || 'image/jpeg',
        })),
      })
      type BatchPresignedResponse = {
        success: { media_item_id: string; presigned_url: string; index: number }[]
        failed: { file_name: string; index: number; reason?: string }[]
      }
      const { success, failed } = presignedRes.data as BatchPresignedResponse
      if (failed.length > 0) {
        setMediaItems((prev) => {
          const next = [...prev]
          failed.forEach(({ index, reason }) => {
            next[index] = {
              ...next[index],
              status: UPLOAD_STATUS.FAILED,
              failure_reason: reason,
              error: { code: UPLOAD_MEDIA_ITEM_ERROR_CODE.INTERNAL_SERVER_ERROR, message: 'presigned URL 발급 실패' },
            }
          })
          return next
        })
      }
      presignedItems = success.map((item) => ({
        mediaItemId: item.media_item_id,
        presignedUrl: item.presigned_url,
        index: item.index,
      }))
    } catch {
      console.error('presigned URL 배치 요청 실패')
      setMediaItems((prev) =>
        prev.map((item) => ({
          ...item,
          status: UPLOAD_STATUS.FAILED,
          error: { code: UPLOAD_MEDIA_ITEM_ERROR_CODE.INTERNAL_SERVER_ERROR, message: 'presigned URL 요청 실패' },
        })),
      )
      setIsUploading(false)
      return
    }

    await Promise.all(
      presignedItems.map(async ({ index, mediaItemId, presignedUrl }) => {
        const mediaItem = mediaItems[index]
        updatePhotoStatus(index, UPLOAD_STATUS.PENDING)
        try {
          await fetch(presignedUrl, {
            method: 'PUT',
            body: mediaItem.file,
            headers: { 'Content-Type': mediaItem.file.type || 'image/jpeg' },
          })
          setMediaItems((prev) => {
            const next = [...prev]
            next[index] = { ...next[index], status: UPLOAD_STATUS.PROCESSING, media_item_id: mediaItemId }
            return next
          })
        } catch {
          updatePhotoStatus(index, UPLOAD_STATUS.FAILED, {
            code: UPLOAD_MEDIA_ITEM_ERROR_CODE.INTERNAL_SERVER_ERROR,
            message: 'S3 업로드 실패',
          })
        }
      }),
    )

    localStorage.setItem(
      UPLOAD_BATCH_STORAGE_KEY,
      JSON.stringify({ batchId: uploadBatchId, totalCount: mediaItems.length } satisfies StoredBatch),
    )

    pollUploadBatchStatus(uploadBatchId)
  }

  const reUploadPhoto = async (index: number, albumId: string) => {
    const mediaItem = mediaItems[index]
    if (mediaItem.status !== UPLOAD_STATUS.FAILED && mediaItem.status !== UPLOAD_STATUS.DUPLICATE) return

    const uploadBatchId = await getOrCreateBatchId(albumId)
    if (!uploadBatchId) return

    try {
      await uploadMediaItem(mediaItem, albumId, uploadBatchId, index)
      pollUploadBatchStatus(uploadBatchId)
    } catch {
      console.error('재업로드 실패')
    }
  }

  const clearPhotos = () => {
    setMediaItems([])
    setIsUploaded(false)
    setCurrentUploadBatchId(null)
    localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
  }

  const uploadCompleted = useMemo(
    () =>
      mediaItems.filter(
        (m) =>
          m.status === UPLOAD_STATUS.COMPLETED ||
          m.status === UPLOAD_STATUS.PROCESSING ||
          m.status === UPLOAD_STATUS.FAILED ||
          m.status === UPLOAD_STATUS.DUPLICATE,
      ).length,
    [mediaItems],
  )

  const progress = useMemo(() => {
    if (mediaItems.length === 0) return 0
    return Math.round((uploadCompleted / mediaItems.length) * 100)
  }, [uploadCompleted, mediaItems])

  const statusCounts = useMemo(
    () => ({
      uploading: mediaItems.filter((m) => m.status === UPLOAD_STATUS.PENDING).length,
      processing: mediaItems.filter((m) => m.status === UPLOAD_STATUS.PROCESSING).length,
      completed: mediaItems.filter((m) => m.status === UPLOAD_STATUS.COMPLETED).length,
      failed: mediaItems.filter((m) => m.status === UPLOAD_STATUS.FAILED).length,
      duplicate: mediaItems.filter((m) => m.status === UPLOAD_STATUS.DUPLICATE).length,
    }),
    [mediaItems],
  )

  const uploadPhase = useMemo((): 'done' | 'uploading' | 'processing' => {
    if (progress === 100) return 'done'
    if (statusCounts.uploading > 0) return 'uploading'
    return 'processing'
  }, [progress, statusCounts.uploading])

  useEffect(() => {
    if (uploadPhase !== 'uploading') return
    const handler = (e: BeforeUnloadEvent) => e.preventDefault()
    window.addEventListener('beforeunload', handler)
    return () => window.removeEventListener('beforeunload', handler)
  }, [uploadPhase])

  const blocker = useBlocker(uploadPhase === 'uploading' && isUploading)

  return {
    mediaItems,
    setMediaItems,
    isUploading,
    isUploaded,
    handleUploadPhotos,
    reUploadPhoto,
    clearPhotos,
    progress,
    uploadCompleted,
    statusCounts,
    uploadPhase,
    blocker,
  }
}
