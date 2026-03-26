import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { MediaItem } from '@/types'
import { useQuery } from '@tanstack/react-query'
import { useCallback } from 'react'

export type UseGetMediaItemsProps = {
  year: number
  month: number
  enabled?: boolean
  identityIds?: number[]
  albumId?: number
}

export const useGetMediaItems = ({ year, month, enabled, identityIds, albumId }: UseGetMediaItemsProps) => {
  const fetchFunc = useCallback(async () => {
    const from = new Date(year, month - 1, 1)
    const to = new Date(year, month)
    const response = await MyAxiosWithAuth.get<MediaItem[]>(API_ROUTES.MEDIA_ITEM.LIST, {
      params: {
        from: from.toISOString(),
        to: to.toISOString(),
        ...(identityIds && identityIds.length > 0 ? { identity_ids: identityIds } : {}),
        ...(albumId ? { album_id: albumId } : {}),
      },
    })
    return response.data
  }, [year, month, identityIds, albumId])

  return useQuery({
    queryKey: ['mediaItems', year, month, identityIds ?? null, albumId ?? null],
    queryFn: fetchFunc,
    enabled,
    staleTime: 1000 * 60 * 1,
  })
}
