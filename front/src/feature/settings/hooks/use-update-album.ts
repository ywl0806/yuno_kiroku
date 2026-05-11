import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { AlbumGroupPermission, AlbumWithPermissions } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

type UpdateAlbumPayload = {
  id: string
  name: string
  permissions: AlbumGroupPermission[]
}

export const useUpdateAlbum = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, name, permissions }: UpdateAlbumPayload) => {
      const response = await MyAxiosWithAuth.put<AlbumWithPermissions>(API_ROUTES.ALBUM.UPDATE(id), {
        name,
        permissions,
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
