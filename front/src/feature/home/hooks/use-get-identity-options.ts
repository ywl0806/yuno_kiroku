import { API_ROUTES } from '@/consts/api-route'
import { IdentityOption } from '@/feature/home/types'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useQuery } from '@tanstack/react-query'

export const useGetIdentityOptions = () => {
  return useQuery({
    queryKey: ['identity-options'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<IdentityOption[]>(API_ROUTES.IDENTITY.OPTIONS)
      return response.data
    },

    staleTime: 1000 * 60 * 1,
  })
}
