import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useDeleteAlbum = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: string) => {
      await MyAxiosWithAuth.delete(API_ROUTES.ALBUM.DELETE(id))
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
