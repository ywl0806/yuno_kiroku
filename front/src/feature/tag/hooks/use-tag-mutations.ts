import { addTagToMediaItem, createTag, deleteTag, removeTagFromMediaItem } from '@/service/tag-service'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useTagMutations = (mediaItemId: string | null) => {
  const queryClient = useQueryClient()

  const invalidateTags = () => {
    queryClient.invalidateQueries({ queryKey: ['media-item-tags', mediaItemId] })
  }

  const addTagMutation = useMutation({
    mutationFn: (tagId: number) => addTagToMediaItem(mediaItemId!, tagId),
    onSuccess: invalidateTags,
  })

  const removeTagMutation = useMutation({
    mutationFn: (tagId: number) => removeTagFromMediaItem(mediaItemId!, tagId),
    onSuccess: invalidateTags,
  })

  const createTagMutation = useMutation({
    mutationFn: ({ name }: { name: string }) => createTag(name),
    onSuccess: (newTag) => {
      queryClient.invalidateQueries({ queryKey: ['tags'] })
      if (mediaItemId !== null) {
        addTagMutation.mutate(newTag.id)
      }
    },
  })

  const deleteTagMutation = useMutation({
    mutationFn: (tagId: number) => deleteTag(tagId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] })
      queryClient.invalidateQueries({ queryKey: ['media-item-tags', mediaItemId] })
    },
  })

  return {
    addTag: addTagMutation.mutate,
    removeTag: removeTagMutation.mutate,
    createTag: createTagMutation.mutate,
    deleteTag: deleteTagMutation.mutate,
    isPending: addTagMutation.isPending || removeTagMutation.isPending || createTagMutation.isPending,
  }
}
