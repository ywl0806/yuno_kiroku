import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Tag } from '@/types'

export const getTags = async (): Promise<Tag[]> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.TAG.LIST)
  return response.data
}

export const createTag = async (name: string): Promise<Tag> => {
  const response = await MyAxiosWithAuth.post(API_ROUTES.TAG.CREATE, { name })
  return response.data
}

export const deleteTag = async (tagId: number): Promise<void> => {
  await MyAxiosWithAuth.delete(API_ROUTES.TAG.DELETE(tagId))
}

export const addTagToMediaItem = async (mediaItemId: string, tagId: number): Promise<void> => {
  await MyAxiosWithAuth.post(API_ROUTES.MEDIA_ITEM.TAGS(mediaItemId), { tag_id: tagId })
}

export const removeTagFromMediaItem = async (mediaItemId: string, tagId: number): Promise<void> => {
  await MyAxiosWithAuth.delete(API_ROUTES.MEDIA_ITEM.TAG_REMOVE(mediaItemId, tagId))
}

export const getTagsForMediaItem = async (mediaItemId: string): Promise<Tag[]> => {
  const response = await MyAxiosWithAuth.get(API_ROUTES.MEDIA_ITEM.TAGS(mediaItemId))
  return response.data
}
