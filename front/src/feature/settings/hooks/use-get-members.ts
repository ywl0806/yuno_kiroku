import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Member } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetMembers = () => {
  return useQuery({
    queryKey: ['members'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<Member[]>(API_ROUTES.USER.MEMBERS)
      return response.data
    },
    staleTime: 1000 * 60 * 1,
  })
}
