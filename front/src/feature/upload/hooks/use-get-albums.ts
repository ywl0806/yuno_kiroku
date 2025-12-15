import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Album } from '@/types'
import { useQuery } from '@tanstack/react-query'
import { useCallback } from 'react'

export const useGetAlbums = () => {
  const fetchFunc = useCallback(async () => {
    const response = await MyAxiosWithAuth.get<Album[]>(API_ROUTES.ALBUM.LIST_FOR_WRITE)
    return response.data
  }, [])

  return useQuery({
    queryKey: ['albums'],
    queryFn: fetchFunc,
  })
}
