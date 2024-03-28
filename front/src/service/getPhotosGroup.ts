import { MyAxios } from '@/lib/myAxios'
import { Photo } from '@/types/photo'

export type GetPhotosGroupOptions = {
  from: Date
  to: Date
}

export type PhotosGroup = {
  _id: {
    year: number
    month: number
  }
  photos: Photo[]
}

export const getPhotosGroup = async (options: GetPhotosGroupOptions): Promise<PhotosGroup[]> => {
  const response = await MyAxios.get('/photo/group', {
    params: {
      from: options.from.toISOString(),
      to: options.to.toISOString(),
    },
  })

  return response.data
}
