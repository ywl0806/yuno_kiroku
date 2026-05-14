import { queryClient } from '@/lib/query-client'
import { updateMediaItemAlbum } from '@/service/media-item-service'
import { useMutation } from '@tanstack/react-query'

export const useUpdateMediaItemAlbum = () => {
  return useMutation({
    mutationFn: ({ mediaItemId, albumId }: { mediaItemId: string; albumId: string }) =>
      updateMediaItemAlbum(mediaItemId, albumId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['search-media-items'] })
      queryClient.invalidateQueries({ queryKey: ['mediaItems'] })
    },
  })
}
