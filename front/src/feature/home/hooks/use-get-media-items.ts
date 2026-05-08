import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { MediaItem } from '@/types'
import { useQuery } from '@tanstack/react-query'
import { useCallback } from 'react'

export type UseGetMediaItemsProps = {
  year: number
  month: number
  enabled?: boolean
}

export const useGetMediaItems = ({ year, month, enabled }: UseGetMediaItemsProps) => {
  const fetchFunc = useCallback(async () => {
    const from = new Date(year, month - 1, 1)
    const to = new Date(year, month)
    const response = await MyAxiosWithAuth.get<MediaItem[]>(API_ROUTES.MEDIA_ITEM.LIST, {
      params: {
        from: from.toISOString(),
        to: to.toISOString(),
      },
    })
    return response.data
  }, [year, month])

  return useQuery({
    queryKey: ['mediaItems', year, month],
    queryFn: fetchFunc,
    enabled,
    staleTime: 1000 * 60 * 1,
  })
}
