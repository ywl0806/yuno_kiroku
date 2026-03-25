import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { IdentityFaceOption } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetIdentityFaceOptions = (
  onlyNotLinked: boolean = true,
  withKidIds: number[] = [],
  withUserIds: number[] = [],
) => {
  return useQuery({
    queryKey: ['identity-face-options', onlyNotLinked, withKidIds, withUserIds],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<IdentityFaceOption[]>(API_ROUTES.SETTINGS.IDENTITY_FACE_OPTIONS, {
        params: { only_not_linked: onlyNotLinked, with_kid_ids: withKidIds, with_user_ids: withUserIds },
      })
      return response.data
    },
  })
}
