import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { MediaItem } from '@/types'

export type GetMediaItemsParams = {
  from: Date
  to: Date
}

export type MediaItems = MediaItem[]

export const getMediaItems = async (options: GetMediaItemsParams): Promise<MediaItems> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.MEDIA_ITEM.LIST, {
    params: {
      from: options.from.toISOString(),
      to: options.to.toISOString(),
    },
  })

  return response.data
}
