import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'

export type MediaItemRange = {
  year: number
  month: number
}

export const getMediaItemRange = async (): Promise<MediaItemRange[]> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.MEDIA_ITEM.RANGE)

  return response.data
}
