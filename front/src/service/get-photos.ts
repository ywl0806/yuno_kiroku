import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Photo } from '@/types'

export type GetPhotosParams = {
  from: Date
  to: Date
}

export type Photos = Photo[]

export const getPhotos = async (options: GetPhotosParams): Promise<Photos> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.PHOTO.LIST, {
    params: {
      from: options.from.toISOString(),
      to: options.to.toISOString(),
    },
  })

  return response.data
}
