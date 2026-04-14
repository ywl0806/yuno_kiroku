import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Member } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

type UpdateFamilyTitleParams = {
  id: number
  groupId: number
  familyTitle: string
  customFamilyTitle: string
}

export const useUpdateMember = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, groupId, familyTitle, customFamilyTitle }: UpdateFamilyTitleParams) => {
      const response = await MyAxiosWithAuth.put<Member>(API_ROUTES.USER.UPDATE(id), {
        group_id: groupId,
        family_title: familyTitle,
        custom_family_title: customFamilyTitle,
      })
      return response.data
    },
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
      queryClient.invalidateQueries({ queryKey: ['member', id] })
    },
  })
}
