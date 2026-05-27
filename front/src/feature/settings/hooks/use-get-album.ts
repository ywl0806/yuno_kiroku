import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { AlbumWithPermissions } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetAlbum = (albumId: string) => {
  return useQuery({
    queryKey: ['album', albumId],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<AlbumWithPermissions>(API_ROUTES.ALBUM.BY_ID(albumId))
      return response.data
    },
    enabled: !!albumId,
  })
}
