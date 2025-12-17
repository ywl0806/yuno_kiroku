import { Progress } from '@/components/ui/progress'
import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { UPLOAD_STATUS, UploadPhoto, UploadStatus } from '@/types'
import { createContext, Dispatch, FC, SetStateAction, useCallback, useContext, useMemo, useState } from 'react'

type UploadPhotoContextType = {
  photos: UploadPhoto[]
  setPhotos: Dispatch<SetStateAction<UploadPhoto[]>>
  handleUploadPhotos: (albumId: number) => Promise<void>
  clearPhotos: () => void
  progress: number
  isUploading: boolean
}

export const UploadPhotoContext = createContext<UploadPhotoContextType>({
  photos: [],
  setPhotos: () => {},
  handleUploadPhotos: () => Promise.resolve(),
  clearPhotos: () => {},
  progress: 0,
  isUploading: false,
})

export const UploadPhotoProvider: FC<{ children: React.ReactNode }> = ({ children }) => {
  const [photos, setPhotos] = useState<UploadPhoto[]>([])
  const [isUploading, setIsUploading] = useState(false)

  const updatePhotoStatus = (index: number, status: UploadStatus) => {
    setPhotos((prev) => {
      prev[index] = { ...prev[index], status }

      return [...prev]
    })
  }
  const uploadPhoto = async (photo: UploadPhoto, albumId: number, index: number) => {
    try {
      updatePhotoStatus(index, UPLOAD_STATUS.PENDING)
      const formData = new FormData()
      formData.append('file', photo.file)

      const response = await MyAxiosWithAuth.post(`${API_ROUTES.PHOTO.UPLOAD}?album_id=${albumId}`, formData)
      updatePhotoStatus(index, UPLOAD_STATUS.SUCCESS)
      return response.data
    } catch (error) {
      updatePhotoStatus(index, UPLOAD_STATUS.ERROR)
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
    setIsUploading(true)
    const process = photos.map((photo, index) => uploadPhoto(photo, albumId, index))
    await Promise.all(process)
    setIsUploadingWithTimeout(5000)
  }

  const clearPhotos = () => {
    setPhotos([])
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

  return (
    <UploadPhotoContext.Provider value={{ photos, setPhotos, handleUploadPhotos, clearPhotos, progress, isUploading }}>
      {isUploading && (
        <div className="fixed left-0 right-0 top-0 px-4">
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
