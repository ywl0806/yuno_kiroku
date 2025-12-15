import { MyAxiosWithAuth } from '../lib/my-axios'
import { API_ROUTES } from '@/constants/api-route'

export type PhotoRange = {
  year: number
  month: number
}

export const getPhotosRange = async (): Promise<PhotoRange[]> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.PHOTO.RANGE)

  return response.data
}
