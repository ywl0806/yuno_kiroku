import { queryClient } from '@/lib/query-client'
import { deleteMediaItem } from '@/service/media-item-service'
import { useMutation } from '@tanstack/react-query'

export const useDeleteMediaItem = () => {
  return useMutation({
    mutationFn: (mediaItemId: string) => deleteMediaItem(mediaItemId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['search-media-items'] })
      queryClient.invalidateQueries({ queryKey: ['mediaItems'] })
    },
  })
}
