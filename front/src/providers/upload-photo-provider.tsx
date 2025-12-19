import { Progress } from '@/components/ui/progress'
import { API_ROUTES } from '@/constants/api-route'
import { HTTP_STATUS } from '@/constants/http'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { UPLOAD_PHOTO_ERROR_CODE, UPLOAD_STATUS, UploadPhoto, UploadPhotoError, UploadStatus } from '@/types'
import { AxiosError } from 'axios'
import { createContext, Dispatch, FC, SetStateAction, useCallback, useContext, useMemo, useState } from 'react'

type UploadPhotoContextType = {
  photos: UploadPhoto[]
  setPhotos: Dispatch<SetStateAction<UploadPhoto[]>>
  handleUploadPhotos: (albumId: number) => Promise<void>
  clearPhotos: () => void
  progress: number
  isUploading: boolean
  reUploadPhoto: (index: number, albumId: number) => void
  isUploaded: boolean
}

export const UploadPhotoContext = createContext<UploadPhotoContextType>({
  photos: [],
  setPhotos: () => {},
  handleUploadPhotos: () => Promise.resolve(),
  clearPhotos: () => {},
  progress: 0,
  isUploading: false,
  reUploadPhoto: () => {},
  isUploaded: false,
})

export const UploadPhotoProvider: FC<{ children: React.ReactNode }> = ({ children }) => {
  const [photos, setPhotos] = useState<UploadPhoto[]>([])
  const [isUploading, setIsUploading] = useState(false)
  const [isUploaded, setIsUploaded] = useState(false)

  const updatePhotoStatus = (index: number, status: UploadStatus, error?: UploadPhotoError) => {
    setPhotos((prev) => {
      prev[index] = { ...prev[index], status, error }

      return [...prev]
    })
  }
  const uploadPhoto = async (photo: UploadPhoto, albumId: number, index: number, retry?: boolean) => {
    try {
      updatePhotoStatus(index, UPLOAD_STATUS.PENDING)
      const formData = new FormData()
      formData.append('file', photo.file)

      const response = await MyAxiosWithAuth.post(
        `${API_ROUTES.PHOTO.UPLOAD}?album_id=${albumId}&retry=${retry ? '1' : '0'}`,
        formData,
      )
      updatePhotoStatus(index, UPLOAD_STATUS.SUCCESS)
      return response.data
    } catch (error) {
      if (error instanceof AxiosError) {
        if (error.response?.status === HTTP_STATUS.CONFLICT) {
          updatePhotoStatus(index, UPLOAD_STATUS.ERROR, {
            code: UPLOAD_PHOTO_ERROR_CODE.DUPLICATE,
            message: error.response?.data.message ?? 'Duplicate photo',
          })
        } else {
          updatePhotoStatus(index, UPLOAD_STATUS.ERROR, {
            code: UPLOAD_PHOTO_ERROR_CODE.INTERNAL_SERVER_ERROR,
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
    const process = photos.map((photo, index) => () => uploadPhoto(photo, albumId, index))

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
    setPhotos([])
    setIsUploaded(false)
  }
  const uploadCompleted = useMemo(() => {
    return photos.filter((photo) => photo.status === UPLOAD_STATUS.SUCCESS || photo.status === UPLOAD_STATUS.ERROR)
      .length
  }, [photos])

  const progress = useMemo(() => {
    if (photos.length === 0) {
      return 0
    }
    return Math.round((uploadCompleted / photos.length) * 100)
  }, [uploadCompleted, photos])

  const reUploadPhoto = (index: number, albumId: number) => {
    const photo = photos[index]
    if (photo.status === UPLOAD_STATUS.ERROR) {
      uploadPhoto(photo, albumId, index, true)
    }
  }
  return (
    <UploadPhotoContext.Provider
      value={{ photos, setPhotos, handleUploadPhotos, clearPhotos, progress, isUploading, reUploadPhoto, isUploaded }}
    >
      {isUploading && (
        <div className="fixed left-0 right-0 top-5 rounded-md bg-white/90 px-4">
          <span className="text-gray-500 text-sm">
            {uploadCompleted} / {photos.length}
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
