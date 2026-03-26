import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Member } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export const useUpdateMemberGroup = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, groupId }: { id: number; groupId: number }) => {
      const response = await MyAxiosWithAuth.put<Member>(API_ROUTES.USER.UPDATE_GROUP(id), { group_id: groupId })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
