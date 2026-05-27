import { queryClient } from '@/lib/query-client'
import { likeMediaItem, unlikeMediaItem } from '@/service/like-service'
import { useMutation } from '@tanstack/react-query'

export const useLikeMutation = () => {
  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['search-media-items'] })
    queryClient.invalidateQueries({ queryKey: ['mediaItems'] })
  }

  const likeMutation = useMutation({
    mutationFn: (mediaItemId: string) => likeMediaItem(mediaItemId),
    onSuccess: invalidate,
  })

  const unlikeMutation = useMutation({
    mutationFn: (mediaItemId: string) => unlikeMediaItem(mediaItemId),
    onSuccess: invalidate,
  })

  const toggle = (mediaItemId: string, currentIsLiked: boolean) => {
    if (currentIsLiked) {
      unlikeMutation.mutate(mediaItemId)
    } else {
      likeMutation.mutate(mediaItemId)
    }
  }

  return {
    toggle,
    isPending: likeMutation.isPending || unlikeMutation.isPending,
  }
}
