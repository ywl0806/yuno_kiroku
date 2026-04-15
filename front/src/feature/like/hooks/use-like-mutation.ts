import { queryClient } from '@/lib/query-client'
import { likeMediaItem, unlikeMediaItem } from '@/service/like-service'
import { useMutation } from '@tanstack/react-query'

export const useLikeMutation = () => {
  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['search-media-items'] })
    queryClient.invalidateQueries({ queryKey: ['mediaItems'] })
    console.log('hogehoge')
  }

  const likeMutation = useMutation({
    mutationFn: (mediaItemId: number) => likeMediaItem(mediaItemId),
    onSuccess: invalidate,
  })

  const unlikeMutation = useMutation({
    mutationFn: (mediaItemId: number) => unlikeMediaItem(mediaItemId),
    onSuccess: invalidate,
  })

  const toggle = (mediaItemId: number, currentIsLiked: boolean) => {
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
