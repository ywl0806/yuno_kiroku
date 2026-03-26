import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Group } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useUpdateGroup = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, name }: { id: number; name: string }) => {
      const response = await MyAxiosWithAuth.put<Group>(API_ROUTES.GROUP.UPDATE(id), { name })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
