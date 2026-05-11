import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'

export const likeMediaItem = async (mediaItemId: string): Promise<void> => {
  await MyAxiosWithAuth.post(API_ROUTES.MEDIA_ITEM.LIKE(mediaItemId))
}

export const unlikeMediaItem = async (mediaItemId: string): Promise<void> => {
  await MyAxiosWithAuth.delete(API_ROUTES.MEDIA_ITEM.LIKE(mediaItemId))
}

export const isMediaItemLiked = async (mediaItemId: string): Promise<boolean> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.MEDIA_ITEM.LIKE(mediaItemId))
  return response.data.is_liked
}
