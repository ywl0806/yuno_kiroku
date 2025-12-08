import { MyAxiosWithAuth } from '@/lib/myAxios'
import { Photo } from '@/types/photo'

export type GetPhotosParams = {
  from: Date
  to: Date
}

export type Photos = Photo[]

export const getPhotos = async (options: GetPhotosParams): Promise<Photos> => {
  const response = await MyAxiosWithAuth.get('/photo', {
    params: {
      from: options.from.toISOString(),
      to: options.to.toISOString(),
    },
  })

  return response.data
}
