import { getTagsForMediaItem } from '@/service/tag-service'
import { useQuery } from '@tanstack/react-query'

export const useGetMediaItemTags = (mediaItemId: number | null) => {
  return useQuery({
    queryKey: ['media-item-tags', mediaItemId],
    queryFn: () => getTagsForMediaItem(mediaItemId!),
    enabled: mediaItemId !== null,
    staleTime: 1000 * 60 * 1,
  })
}
