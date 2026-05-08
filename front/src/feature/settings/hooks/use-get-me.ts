import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Me } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetMe = () => {
  return useQuery({
    queryKey: ['me'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<Me>(API_ROUTES.USER.ME)
      return response.data
    },
    staleTime: 1000 * 60 * 5,
  })
}
