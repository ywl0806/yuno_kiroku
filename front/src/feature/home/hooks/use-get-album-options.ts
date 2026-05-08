import { API_ROUTES } from '@/consts/api-route'
import { AlbumOption } from '@/feature/home/types'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useQuery } from '@tanstack/react-query'

export const useGetAlbumOptions = () => {
  return useQuery({
    queryKey: ['album-options'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<AlbumOption[]>(API_ROUTES.ALBUM.OPTIONS)
      return response.data
    },
    staleTime: 1000 * 60 * 1,
  })
}
