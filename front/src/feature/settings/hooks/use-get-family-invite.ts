import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useQuery } from '@tanstack/react-query'

type FamilyInviteResponse = {
  inviteUrl: string
}

export const useGetFamilyInvite = (familyId: string) => {
  return useQuery({
    queryKey: ['family', familyId, 'invite'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<FamilyInviteResponse>(
        API_ROUTES.FAMILY.INVITE(familyId),
      )
      return response.data
    },
    enabled: !!familyId,
  })
}
