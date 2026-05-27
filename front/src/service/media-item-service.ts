import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'

export const deleteMediaItem = async (mediaItemId: string): Promise<void> => {
  await MyAxiosWithAuth.delete(API_ROUTES.MEDIA_ITEM.DELETE(mediaItemId))
}

export const updateMediaItemAlbum = async (mediaItemId: string, albumId: string): Promise<void> => {
  await MyAxiosWithAuth.patch(API_ROUTES.MEDIA_ITEM.UPDATE_ALBUM(mediaItemId), { album_id: albumId })
}
