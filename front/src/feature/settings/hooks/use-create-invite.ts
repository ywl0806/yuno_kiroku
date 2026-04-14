import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useMutation } from '@tanstack/react-query'

type CreateInviteParams = {
  groupId: number
  familyTitle: string
  customFamilyTitle: string
}

type CreateInviteResponse = {
  token: string
  invite_url: string
  expires_at: string
}

export const useCreateInvite = () => {
  return useMutation({
    mutationFn: async ({ groupId, familyTitle, customFamilyTitle }: CreateInviteParams) => {
      const response = await MyAxiosWithAuth.post<CreateInviteResponse>(API_ROUTES.INVITE.CREATE, {
        group_id: groupId,
        family_title: familyTitle,
        custom_family_title: customFamilyTitle,
      })
      return response.data
    },
  })
}
