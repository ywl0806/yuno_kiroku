import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Group } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetGroups = () => {
  return useQuery({
    queryKey: ['groups'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<Group[]>(API_ROUTES.GROUP.LIST)
      return response.data
    },
    staleTime: 1000 * 60 * 1,
  })
}
