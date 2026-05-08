import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Me } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useUpdateMe = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: { name: string }) => {
      const response = await MyAxiosWithAuth.put<Me>(API_ROUTES.USER.ME, payload)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['me'] })
    },
  })
}
