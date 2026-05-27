import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Member } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetMember = (id: string) => {
  return useQuery({
    queryKey: ['member', id],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<Member>(API_ROUTES.USER.MEMBER(id))
      return response.data
    },
    staleTime: 1000 * 60 * 5,
    enabled: !!id,
  })
}
