import { BottomSheet } from '@/components/ui/bottom-sheet'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { API_ROUTES } from '@/consts/api-route'
import { UPLOAD_STATUS } from '@/enums'
import { useGetAlbums } from '@/feature/upload/hooks/use-get-albums'
import { getSessionStorage, SESSION_STORAGE_KEY, setSessionStorage } from '@/lib/session-storage'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import {
  UPLOAD_MEDIA_ITEM_ERROR_CODE,
  UploadMediaItem,
  UploadMediaItemError,
} from '@/types'
import { UploadStatus } from '@/enums'
import { AxiosError } from 'axios'
import { Album, AlertTriangle, CheckCircle2, ChevronDown, ChevronUp, Loader2, X } from 'lucide-react'
import { createContext, Dispatch, FC, SetStateAction, useCallback, useContext, useEffect, useMemo, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { unstable_useBlocker as useBlocker, useNavigate, useLocation } from 'react-router-dom'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@radix-ui/react-collapsible'

const UPLOAD_BATCH_STORAGE_KEY = 'yuno_upload_batch'

type StoredBatch = {
  batchId: number
  totalCount: number
}

type UploadPhotoContextType = {
  mediaItems: UploadMediaItem[]
  setMediaItems: Dispatch<SetStateAction<UploadMediaItem[]>>
  handleUploadPhotos: (albumId: string) => Promise<void>
  clearPhotos: () => void
  progress: number
  isUploading: boolean
  reUploadPhoto: (index: number, albumId: string) => void
  isUploaded: boolean
  // 앨범 시트
  albumSheetOpen: boolean
  openAlbumSheet: () => void
  closeAlbumSheet: () => void
  selectedAlbumId: string | null
  selectedAlbumName: string | null
  setSelectedAlbumId: (id: string) => void
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
  albumSheetOpen: false,
  openAlbumSheet: () => { },
  closeAlbumSheet: () => { },
  selectedAlbumId: null,
  selectedAlbumName: null,
  setSelectedAlbumId: () => { },
})

export const UploadPhotoProvider: FC<{ children: React.ReactNode }> = ({ children }) => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const { data: albums } = useGetAlbums()

  const [mediaItems, setMediaItems] = useState<UploadMediaItem[]>([])
  const [isUploading, setIsUploading] = useState(false)
  const [isUploaded, setIsUploaded] = useState(false)
  const [currentUploadBatchId, setCurrentUploadBatchId] = useState<number | null>(null)

  // 앨범 시트 상태
  const [albumSheetOpen, setAlbumSheetOpen] = useState(false)
  const [selectedAlbumId, setSelectedAlbumIdState] = useState<string | null>(
    getSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID) || null,
  )

  const openAlbumSheet = useCallback(() => setAlbumSheetOpen(true), [])
  const closeAlbumSheet = useCallback(() => setAlbumSheetOpen(false), [])

  const setSelectedAlbumId = useCallback((id: string) => {
    setSelectedAlbumIdState(id)
    setSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID, id)
    closeAlbumSheet()
    if (location.pathname !== '/upload') {
      navigate('/upload')
    }
  }, [closeAlbumSheet, navigate, location.pathname])

  // 복구 관련 상태
  const [isRecovering, setIsRecovering] = useState(false)
  const [recoveredBatch, setRecoveredBatch] = useState<StoredBatch | null>(null)
  const [recoveryDone, setRecoveryDone] = useState(false)
  const [bannerOpen, setBannerOpen] = useState(true)
  const recoveryPollRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const updatePhotoStatus = (index: number, status: UploadStatus, error?: UploadMediaItemError) => {
    setMediaItems((prev) => {
      prev[index] = { ...prev[index], status, error }
      return [...prev]
    })
  }

  const uploadMediaItem = async (
    mediaItem: UploadMediaItem,
    albumId: string,
    uploadBatchId: number,
    index: number,
    retry?: boolean,
  ) => {
    try {
      updatePhotoStatus(index, UPLOAD_STATUS.PENDING)

      const presignedRes = await MyAxiosWithAuth.post(
        `${API_ROUTES.MEDIA_ITEM.PRESIGNED_URL}?album_id=${albumId}&upload_batch_id=${uploadBatchId}`,
        { file_name: mediaItem.file.name, content_type: mediaItem.file.type || 'image/jpeg' },
      )
      const { media_item_id: mediaItemId, presigned_url: presignedUrl } = presignedRes.data as {
        media_item_id: string
        presigned_url: string
        storage_key: string
      }

      await fetch(presignedUrl, {
        method: 'PUT',
        body: mediaItem.file,
        headers: { 'Content-Type': mediaItem.file.type || 'image/jpeg' },
      })

      setMediaItems((prev) => {
        const next = [...prev]
        next[index] = {
          ...next[index],
          status: UPLOAD_STATUS.PROCESSING,
          media_item_id: mediaItemId,
        }
        return next
      })

      if (retry) {
        pollUploadBatchStatus(uploadBatchId)
        return
      }
      return presignedRes.data
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

  const handleUploadPhotos = async (albumId: string) => {
    setIsUploaded(true)
    setIsUploading(true)

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

    const process = mediaItems.map(
      (mediaItem, index) => () => uploadMediaItem(mediaItem, albumId, uploadBatchId, index),
    )

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

    localStorage.setItem(
      UPLOAD_BATCH_STORAGE_KEY,
      JSON.stringify({ batchId: uploadBatchId, totalCount: mediaItems.length } satisfies StoredBatch),
    )

    pollUploadBatchStatus(uploadBatchId)
  }

  const pollUploadBatchStatus = async (uploadBatchId: number) => {
    if (isUploading) {
      return
    }

    const maxAttempts = 60
    let attempts = 0

    const poll = async () => {
      try {
        const response = await MyAxiosWithAuth.get<UploadBatchStatusResponse>(
          `${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH_STATUS}?upload_batch_id=${uploadBatchId}`,
        )

        const statuseMap = new Map<string, UploadStatus>()
        response.data.statuses.forEach((status) => {
          statuseMap.set(status.id, status.upload_status)
        })

        setMediaItems((prev) => {
          const next = [...prev]
          next.forEach((mediaItem) => {
            mediaItem.status = statuseMap.get(mediaItem.media_item_id ?? '') ?? UPLOAD_STATUS.PENDING
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
      } catch (error) {
        console.error('상태 조회 실패:', error)
        localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
        setIsUploadingWithTimeout(5000)
      }
    }

    poll()
  }

  const pollRecoveredBatch = useCallback(async (batch: StoredBatch) => {
    const maxAttempts = 60
    let attempts = 0

    const poll = async () => {
      try {
        const response = await MyAxiosWithAuth.get<UploadBatchStatusResponse>(
          `${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH_STATUS}?upload_batch_id=${batch.batchId}`,
        )

        if (response.data.is_completed) {
          localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
          setIsRecovering(false)
          setRecoveredBatch(null)
          setRecoveryDone(true)
          setTimeout(() => setRecoveryDone(false), 5000)
          return
        }

        attempts++
        if (attempts < maxAttempts) {
          recoveryPollRef.current = setTimeout(poll, 5000)
        } else {
          localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
          setIsRecovering(false)
          setRecoveredBatch(null)
        }
      } catch (error) {
        console.error('복구 상태 조회 실패:', error)
        localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
        setIsRecovering(false)
        setRecoveredBatch(null)
      }
    }

    poll()
  }, [])

  useEffect(() => {
    const stored = localStorage.getItem(UPLOAD_BATCH_STORAGE_KEY)
    if (!stored) return
    try {
      const batch: StoredBatch = JSON.parse(stored)
      setRecoveredBatch(batch)
      setIsRecovering(true)
      pollRecoveredBatch(batch)
    } catch {
      localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
    }

    return () => {
      if (recoveryPollRef.current) clearTimeout(recoveryPollRef.current)
    }
  }, [])

  const clearPhotos = () => {
    setMediaItems([])
    setIsUploaded(false)
    setCurrentUploadBatchId(null)
    localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
  }

  const uploadCompleted = useMemo(() => {
    return mediaItems.filter(
      (mediaItem) =>
        mediaItem.status === UPLOAD_STATUS.COMPLETED ||
        mediaItem.status === UPLOAD_STATUS.PROCESSING ||
        mediaItem.status === UPLOAD_STATUS.FAILED ||
        mediaItem.status === UPLOAD_STATUS.DUPLICATE,
    ).length
  }, [mediaItems])

  const progress = useMemo(() => {
    if (mediaItems.length === 0) {
      return 0
    }
    return Math.round((uploadCompleted / mediaItems.length) * 100)
  }, [uploadCompleted, mediaItems])

  const statusCounts = useMemo(() => ({
    uploading: mediaItems.filter((m) => m.status === UPLOAD_STATUS.PENDING).length,
    processing: mediaItems.filter((m) => m.status === UPLOAD_STATUS.PROCESSING).length,
    completed: mediaItems.filter((m) => m.status === UPLOAD_STATUS.COMPLETED).length,
    failed: mediaItems.filter((m) => m.status === UPLOAD_STATUS.FAILED).length,
    duplicate: mediaItems.filter((m) => m.status === UPLOAD_STATUS.DUPLICATE).length,
  }), [mediaItems])

  const uploadPhase = useMemo(() => {
    if (progress === 100) return 'done'
    if (statusCounts.uploading > 0) return 'uploading'
    return 'processing'
  }, [progress, statusCounts.uploading])

  useEffect(() => {
    if (uploadPhase !== 'uploading') return
    const handler = (e: BeforeUnloadEvent) => {
      e.preventDefault()
    }
    window.addEventListener('beforeunload', handler)
    return () => window.removeEventListener('beforeunload', handler)
  }, [uploadPhase])

  const blocker = useBlocker(uploadPhase === 'uploading' && isUploading)

  const reUploadPhoto = async (index: number, albumId: string) => {
    const mediaItem = mediaItems[index]
    if (mediaItem.status === UPLOAD_STATUS.FAILED || mediaItem.status === UPLOAD_STATUS.DUPLICATE) {
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

  const selectedAlbum = useMemo(
    () => albums?.find((a) => a.id === selectedAlbumId) ?? null,
    [albums, selectedAlbumId],
  )

  const selectedAlbumName = useMemo(() => {
    if (!selectedAlbum) return null
    return selectedAlbum.is_common ? t('settings.album.commonName') : selectedAlbum.name
  }, [selectedAlbum, t])

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
        albumSheetOpen,
        openAlbumSheet,
        closeAlbumSheet,
        selectedAlbumId,
        selectedAlbumName,
        setSelectedAlbumId,
      }}
    >
      {/* SPA 이동 차단 다이얼로그 */}
      {blocker.state === 'blocked' && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50">
          <div className="mx-4 w-full max-w-sm rounded-2xl bg-white px-6 py-5 shadow-2xl">
            <div className="mb-3 flex items-center gap-2">
              <AlertTriangle className="size-5 shrink-0 text-amber-500" />
              <span className="font-semibold text-gray-900">{t('upload.progress.leaveWarningTitle')}</span>
            </div>
            <p className="mb-5 text-sm text-gray-500">
              {t('upload.progress.leaveWarningBody')}
            </p>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => blocker.reset?.()}
                className="flex-1 rounded-lg border border-gray-200 py-2.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50"
              >
                {t('upload.progress.leaveWarningStay')}
              </button>
              <button
                type="button"
                onClick={() => blocker.proceed?.()}
                className="flex-1 rounded-lg bg-alizarin py-2.5 text-sm font-medium text-white transition-colors hover:bg-alizarin/90"
              >
                {t('upload.progress.leaveWarningLeave')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 업로드/처리 중 진행 배너 */}
      {isUploading && (
        <Collapsible open={bannerOpen} onOpenChange={setBannerOpen}>
          <div className="fixed bottom-[calc(4rem+env(safe-area-inset-bottom)+0.5rem)] left-1/2 z-50 w-[320px] -translate-x-1/2 rounded-2xl border bg-white px-5 py-4 shadow-xl transition-all duration-300">
            <div className="flex items-center gap-2">
              {uploadPhase === 'done' ? (
                <CheckCircle2 className="size-5 shrink-0 text-emerald" />
              ) : (
                <Loader2 className="size-5 shrink-0 animate-spin text-peter-river" />
              )}
              <span className="text-sm font-medium text-gray-800">
                {uploadPhase === 'uploading' && t('upload.progress.uploading')}
                {uploadPhase === 'processing' && t('upload.progress.processing')}
                {uploadPhase === 'done' && t('upload.progress.done')}
              </span>
              <span className="ml-auto text-xs text-gray-400">
                {uploadCompleted} / {mediaItems.length}
              </span>
              <CollapsibleTrigger asChild>
                <button type="button" className="ml-1 rounded p-0.5 text-gray-400 hover:text-gray-600">
                  {bannerOpen ? <ChevronDown className="size-4" /> : <ChevronUp className="size-4" />}
                </button>
              </CollapsibleTrigger>
            </div>
            <CollapsibleContent>
              <Progress value={progress} className="mt-3 h-1.5 w-full" />
              <div className="mt-2 flex gap-3 text-xs text-gray-400">
                {statusCounts.uploading > 0 && <span>{t('upload.progress.statusUploading')} {statusCounts.uploading}</span>}
                {statusCounts.processing > 0 && <span>{t('upload.progress.statusProcessing')} {statusCounts.processing}</span>}
                {statusCounts.completed > 0 && <span className="text-emerald">{t('upload.progress.statusCompleted')} {statusCounts.completed}</span>}
                {statusCounts.failed > 0 && <span className="text-alizarin">{t('upload.progress.statusFailed')} {statusCounts.failed}</span>}
                {statusCounts.duplicate > 0 && <span>{t('upload.progress.statusDuplicate')} {statusCounts.duplicate}</span>}
              </div>
            </CollapsibleContent>
          </div>
        </Collapsible>
      )}

      {/* 복구 중 배너 */}
      {isRecovering && recoveredBatch && (
        <div className="fixed bottom-[calc(4rem+env(safe-area-inset-bottom)+0.5rem)] left-1/2 z-50 w-[320px] -translate-x-1/2 rounded-2xl border bg-white px-5 py-4 shadow-xl transition-all duration-300">
          <div className="flex items-center gap-2">
            <Loader2 className="size-5 shrink-0 animate-spin text-peter-river" />
            <span className="text-sm font-medium text-gray-800">{t('upload.progress.recoveryInProgress')}</span>
            <span className="ml-auto text-xs text-gray-400">{recoveredBatch.totalCount}장</span>
            <button
              type="button"
              onClick={() => {
                if (recoveryPollRef.current) clearTimeout(recoveryPollRef.current)
                localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
                setIsRecovering(false)
                setRecoveredBatch(null)
              }}
              className="ml-1 rounded p-0.5 text-gray-400 hover:text-gray-600"
            >
              <X className="size-3.5" />
            </button>
          </div>
        </div>
      )}

      {/* 복구 완료 배너 */}
      {recoveryDone && (
        <div className="fixed bottom-[calc(4rem+env(safe-area-inset-bottom)+0.5rem)] left-1/2 z-50 w-[320px] -translate-x-1/2 rounded-2xl border bg-white px-5 py-4 shadow-xl transition-all duration-300">
          <div className="flex items-center gap-2">
            <CheckCircle2 className="size-5 shrink-0 text-emerald" />
            <span className="text-sm font-medium text-gray-800">{t('upload.progress.recoveryDone')}</span>
          </div>
        </div>
      )}

      {/* 앨범 선택 바텀 시트 */}
      <BottomSheet
        open={albumSheetOpen}
        onClose={closeAlbumSheet}
        title={t('upload.selectAlbum')}
      >
        <div className="flex flex-col gap-2 px-4 pb-6 pt-2">
          {albums?.map((album) => (
            <Button
              key={album.id}
              variant={selectedAlbum?.id === album.id ? 'default' : 'outline'}
              className="h-12 justify-start text-[1rem]"
              onClick={() => setSelectedAlbumId(album.id)}
            >
              <Album className="size-4 shrink-0" />
              {album.is_common ? t('settings.album.commonName') : album.name}
            </Button>
          ))}
        </div>
      </BottomSheet>

      {children}
    </UploadPhotoContext.Provider>
  )
}

export const useUploadPhoto = () => {
  return useContext(UploadPhotoContext)
}

type UploadBatchStatusResponse = {
  statuses: {
    id: string
    upload_status: UploadStatus
  }[]
  is_completed: boolean
}
