import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { AlbumGroupPermission, AlbumWithPermissions } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

type CreateAlbumPayload = {
  name: string
  permissions: AlbumGroupPermission[]
}

export const useCreateAlbum = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: CreateAlbumPayload) => {
      const response = await MyAxiosWithAuth.post<AlbumWithPermissions>(API_ROUTES.ALBUM.CREATE, payload)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
