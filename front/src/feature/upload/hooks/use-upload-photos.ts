import { API_ROUTES } from '@/constants/api-route'
import { UPLOAD_STATUS, UploadPhoto, UploadPhotoStatus } from '@/feature/home/types'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useState } from 'react'

type UseUploadPhotosProps = {
  photos: UploadPhoto[]
  albumId: number
}

export const useUploadPhotos = ({ photos, albumId }: UseUploadPhotosProps) => {
  const [uploadPhotoStatus, setUploadPhotoStatus] = useState<UploadPhotoStatus>({})

  const uploadPhoto = async (photo: UploadPhoto, index: number) => {
    try {
      setUploadPhotoStatus((prev) => ({ ...prev, [index]: UPLOAD_STATUS.PENDING }))
      const formData = new FormData()
      formData.append('file', photo.file)

      const response = await MyAxiosWithAuth.post(`${API_ROUTES.PHOTO.UPLOAD}?album_id=${albumId}`, formData)
      setUploadPhotoStatus((prev) => ({ ...prev, [index]: UPLOAD_STATUS.SUCCESS }))
      return response.data
    } catch (error) {
      setUploadPhotoStatus((prev) => ({ ...prev, [index]: UPLOAD_STATUS.ERROR }))
      return null
    }
  }

  const handleUploadPhotos = async () => {
    // 모든 사진의 업로드 상태를 IDLE로 초기화
    photos.forEach((_, index) => {
      setUploadPhotoStatus((prev) => ({ ...prev, [index]: UPLOAD_STATUS.IDLE }))
    })

    // 모든 사진을 업로드
    const process = []

    // 비동기적으로 모든 사진을 업로드
    for (let index = 0; index < photos.length; index++) {
      const photo = photos[index]
      const result = await uploadPhoto(photo, index)
      process.push(result)
    }

    const results = await Promise.all(process)
    return results
  }

  return { handleUploadPhotos, uploadPhotoStatus }
}
