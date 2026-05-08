import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Group } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useCreateGroup = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: { name: string }) => {
      const response = await MyAxiosWithAuth.post<Group>(API_ROUTES.GROUP.CREATE, payload)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
