import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useMutation } from '@tanstack/react-query'

type CreateInviteResponse = {
  token: string
  invite_url: string
  expires_at: string
}

export const useCreateInvite = () => {
  return useMutation({
    mutationFn: async (groupId: number) => {
      const response = await MyAxiosWithAuth.post<CreateInviteResponse>(API_ROUTES.INVITE.CREATE, {
        group_id: groupId,
      })
      return response.data
    },
  })
}
