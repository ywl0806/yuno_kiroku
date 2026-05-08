import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useDeleteKid = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (kidId: number) => {
      await MyAxiosWithAuth.delete(API_ROUTES.KID.DELETE(kidId))
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
