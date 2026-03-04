import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Family } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetFamilies = () => {
  return useQuery({
    queryKey: ['families'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<Family[]>(API_ROUTES.FAMILY.LIST)
      return response.data
    },
    staleTime: 1000 * 60 * 1,
  })
}
